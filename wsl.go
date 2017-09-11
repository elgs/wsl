package wsl

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WSL struct {
	Config *Config
	db     *sql.DB
}

func (this *WSL) connectToDb() error {
	if this.db == nil {
		db, err := sql.Open(this.Config.DbType, this.Config.DbUrl)
		if err != nil {
			return err
		}
		this.db = db
	}
	return nil
}

func (this *WSL) Start() {
	err := this.connectToDb()
	if err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))

		urlPath := strings.Split(r.URL.Path, "/")
		qParams, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			log.Println(err)
			return
		}
		if len(urlPath) < 2 {
			log.Println("Invalid URL.")
			return
		}
		qID := urlPath[1]
		// fmt.Println(qID)
		// fmt.Println(qParams, len(qParams))

		if script, ok := this.Config.Scripts[qID]; ok {
			result, err := this.exec(this.db, script, valuesToMap(&qParams))
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(result)
		} else {
			log.Println("Invalid script.")
			return
		}

	})

	if this.Config.httpEnabled() {
		srv := &http.Server{
			Addr:         this.Config.HttpAddr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func() {
			log.Fatal(srv.ListenAndServe())
			fmt.Println(fmt.Sprint("Listening on http://", this.Config.HttpAddr, "/"))
		}()
	}

	if this.Config.httpsEnabled() {
		srvs := &http.Server{
			Addr:         this.Config.HttpsAddr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func() {
			log.Fatal(srvs.ListenAndServeTLS(this.Config.CertFile, this.Config.KeyFile))
			fmt.Println(fmt.Sprint("Listening on https://", this.Config.HttpsAddr, "/"))
		}()
	}
}
