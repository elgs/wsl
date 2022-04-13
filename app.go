package wsl

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

type App struct {
	Config    *Config
	Databases map[string]*sql.DB

	// globalInterceptors is the registry of a list of Interceptors that is
	// guaranteed to be executed in its list order on each query.
	globalInterceptors []Interceptor

	// queryInterceptors is the registry for each named query (qID) of a list of
	// Interceptors that is guaranteed to be executed in its list order on each
	// query.
	queryInterceptors map[string][]Interceptor

	// jobs is the registry for cron jobs
	Jobs map[string]*Job

	Cron *cron.Cron

	Scripts map[string]string
}

func New(config *Config) *App {
	wsl := &App{
		Config:            config,
		Scripts:           map[string]string{},
		Databases:         map[string]*sql.DB{},
		queryInterceptors: map[string][]Interceptor{},
		Jobs:              map[string]*Job{},
		Cron:              cron.New(),
	}
	return wsl
}

func (this *App) connectToDb(dbName string) error {
	if this.Databases[dbName] == nil {
		dbData := this.Config.Databases[dbName]
		db, err := sql.Open(dbData.Type, dbData.Url)
		if err != nil {
			return err
		}
		this.Databases[dbName] = db
	}
	return nil
}

func (this *App) Start() {
	for dbName := range this.Config.Databases {
		if err := this.connectToDb(dbName); err != nil {
			log.Println(err)
			return
		}
	}

	for _, b := range this.Jobs {
		entryId, err := this.Cron.AddFunc(b.Cron, b.MakeFunc(this))
		if err != nil {
			log.Println(err)
			return
		}
		b.ID = &entryId
	}

	this.Cron.Start()

	http.HandleFunc("/", this.defaultHandler)
	// http.HandleFunc("/ws", this.wsHandler)

	if this.Config.Web.HttpAddr != "" {
		srv := &http.Server{
			Addr:         this.Config.Web.HttpAddr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func() {
			fmt.Println(fmt.Sprint("Listening on http://", this.Config.Web.HttpAddr, "/"))
			log.Fatal(srv.ListenAndServe())
		}()
	}

	if this.Config.Web.HttpsAddr != "" {
		srvs := &http.Server{
			Addr:         this.Config.Web.HttpsAddr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func() {
			fmt.Println(fmt.Sprint("Listening on https://", this.Config.Web.HttpsAddr, "/"))
			log.Fatal(srvs.ListenAndServeTLS(this.Config.Web.CertFile, this.Config.Web.KeyFile))
		}()
	}

	Hook()
}

func Hook() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case sig := <-sigs:
				fmt.Println(sig)
				// cleanup code here
				done <- true
			}
		}
	}()

	<-done
	fmt.Println("Bye!")
}
