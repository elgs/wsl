package wsl

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

type WSL struct {
	Config    *Config
	databases map[string]*sql.DB
}

func New(confFile string) (*WSL, error) {
	config, err := NewConfig(confFile)
	if err != nil {
		return nil, err
	}
	return &WSL{
		Config:    config,
		databases: map[string]*sql.DB{},
	}, nil
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
	err := this.connectToDb("main")
	if err != nil {
		log.Println(err)
		return
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
