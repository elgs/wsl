package wsl

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
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

	// queryInterceptors is the registry for each named query (queryID) of a list of
	// Interceptors that is guaranteed to be executed in its list order on each
	// query.
	queryInterceptors map[string][]Interceptor

	// jobs is the registry for cron jobs
	Jobs map[string]*Job

	Cron *cron.Cron

	Scripts map[string]string
}

func NewApp(config *Config) *App {
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

func (this *App) GetDB(dbName string) *sql.DB {
	if db, ok := this.Databases[dbName]; ok {
		return db
	}
	dbData := this.Config.Databases[dbName]
	db, err := sql.Open(dbData.Type, dbData.Url)
	if err != nil {
		return nil
	}
	this.Databases[dbName] = db
	return db
}

func (this *App) GetScript(scriptName string, forceReload bool) string {
	if sqlString, ok := this.Scripts[scriptName]; ok {
		return sqlString
	}
	data, err := ioutil.ReadFile(path.Join("scripts", scriptName, ".sql"))
	if err != nil {
		return ""
	}
	sqlString := string(data)
	this.Scripts[scriptName] = sqlString
	return sqlString
}

func (this *App) Start() {
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
