package wsl

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type WSL struct {
	Config    *Config
	Scripts   map[string]string
	databases map[string]*sql.DB

	// globalInterceptors is the registry of a list of Interceptors that is
	// guaranteed to be executed in its list order on each query.
	globalInterceptors []Interceptor

	// queryInterceptors is the registry for each named query (qID) of a list of
	// Interceptors that is guaranteed to be executed in its list order on each
	// query.
	queryInterceptors map[string][]Interceptor
}

func NewWithConfigPath(confPath string) (*WSL, error) {
	confFile := flag.String("c", confPath, "configration file path")
	flag.Parse()

	return NewWithConfigJSON(*confFile)
}

func NewWithConfigJSON(confFile string) (*WSL, error) {
	config, err := NewConfig(confFile)
	if err != nil {
		return nil, err
	}
	wsl := &WSL{
		Config:            config,
		Scripts:           map[string]string{},
		databases:         map[string]*sql.DB{},
		queryInterceptors: map[string][]Interceptor{},
	}
	err = wsl.LoadScripts("")
	return wsl, err
}

func (this *WSL) connectToDb(dbName string) error {
	if this.databases[dbName] == nil {
		dbData := this.Config.Databases[dbName].(map[string]interface{})
		db, err := sql.Open(dbData["db_type"].(string), dbData["db_url"].(string))
		if err != nil {
			return err
		}
		this.databases[dbName] = db
	}
	return nil
}

func (this *WSL) Start() {
	if err := this.connectToDb("main"); err != nil {
		log.Println(err)
		return
	}

	if this.databases["audit"] != nil {
		if err := this.connectToDb("audit"); err != nil {
			log.Println(err)
			return
		}
	}

	http.HandleFunc("/", this.defaultHandler)
	// http.HandleFunc("/ws", this.wsHandler)

	if this.Config.httpEnabled() {
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

	if this.Config.httpsEnabled() {
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
}
