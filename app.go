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

	"github.com/elgs/gosplitargs"
	"github.com/elgs/optional"
	"github.com/robfig/cron/v3"
)

type Statement struct {
	Index  int
	Label  string
	Text   string
	Script *Script
}

type Script struct {
	ID           string
	Text         string
	Params       *[]string
	Statements   *[]Statement
	Interceptors *[]Interceptor
}

type App struct {
	Config    *Config
	Databases map[string]*sql.DB

	// globalInterceptors is the registry of a list of Interceptors that is
	// guaranteed to be executed in its list order on each query.
	globalInterceptors []Interceptor

	// jobs is the registry for cron jobs
	Jobs map[string]*Job

	Cron *cron.Cron

	Scripts map[string]*Script
}

func NewApp(config *Config) *App {
	return &App{
		Config:    config,
		Scripts:   map[string]*Script{},
		Databases: map[string]*sql.DB{},
		Jobs:      map[string]*Job{},
		Cron:      cron.New(),
	}
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

func BuildScript(scriptString string) *optional.Optional[*Script] {
	statements, err := gosplitargs.SplitArgs(scriptString, ";", true)
	if err != nil {
		return optional.New[*Script](nil, err)
	}

	script := &Script{
		Text:         scriptString,
		Params:       &[]string{},
		Statements:   &[]Statement{},
		Interceptors: &[]Interceptor{},
	}
	for index, statementString := range statements {
		label, statementSQL := SplitSqlLable(statementString)
		param := ExtractSQLParameter(statementSQL)
		if len(param) > 0 {
			*script.Params = append(*script.Params, param)
		}
		statement := &Statement{
			Index:  index,
			Label:  label,
			Text:   statementString,
			Script: script,
		}
		*script.Statements = append(*script.Statements, *statement)
	}
	return optional.New(script, nil)
}

func (this *App) GetScript(scriptName string, forceReload bool) *optional.Optional[*Script] {
	if script, ok := this.Scripts[scriptName]; ok {
		return optional.New(script, nil)
	}

	data, err := ioutil.ReadFile(path.Join("scripts", scriptName, ".sql"))
	if err != nil {
		return optional.New[*Script](nil, err)
	}
	scriptString := string(data)
	scriptOpt := BuildScript(scriptString)
	if scriptOpt.Error != nil {
		return scriptOpt
	}

	scriptOpt.Data.ID = scriptName
	this.Scripts[scriptName] = scriptOpt.Data
	return scriptOpt
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
