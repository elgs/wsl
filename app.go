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
	"strings"
	"syscall"
	"time"

	"github.com/elgs/gosplitargs"
)

type Statement struct {
	Index        int
	Label        string
	Text         string
	Param        string
	IsQuery      bool
	ShouldExport bool
	Script       *Script
}

type Script struct {
	ID           string
	Text         string
	Statements   *[]Statement
	Interceptors *[]Interceptor
}

type App struct {
	Config             *Config
	Databases          map[string]*sql.DB
	globalInterceptors []Interceptor
	Scripts            map[string]*Script
}

func NewApp(config *Config) *App {
	return &App{
		Config:    config,
		Scripts:   map[string]*Script{},
		Databases: map[string]*sql.DB{},
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

func BuildScript(scriptString string) (*Script, error) {
	statements, err := gosplitargs.SplitArgs(scriptString, ";", true)
	if err != nil {
		return nil, err
	}

	script := &Script{
		Text:         scriptString,
		Statements:   &[]Statement{},
		Interceptors: &[]Interceptor{},
	}
	for index, statementString := range statements {
		if len(strings.TrimSpace(statementString)) == 0 {
			continue
		}
		label, statementSQL := SplitSqlLable(statementString)
		param := ExtractSQLParameter(statementSQL)
		statement := &Statement{
			Index:        index,
			Label:        label,
			Text:         statementString,
			Param:        param,
			Script:       script,
			IsQuery:      IsQuery(statementSQL),
			ShouldExport: ShouldExport(statementSQL),
		}
		*script.Statements = append(*script.Statements, *statement)
	}
	return script, nil
}

func (this *App) GetScript(scriptName string, forceReload bool) (*Script, error) {
	if script, ok := this.Scripts[scriptName]; ok {
		return script, nil
	}

	scriptPath := path.Join("scripts", scriptName+".sql")
	data, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		return nil, err
	}
	scriptString := string(data)
	script, err := BuildScript(scriptString)
	if err != nil {
		return nil, err
	}

	script.ID = scriptName
	this.Scripts[scriptName] = script
	return script, nil
}

func (this *App) Start() {
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
