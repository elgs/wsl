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

// var addr = flag.String("addr", ":8080", "http service address, default to :8080")
// flag.Parse()

type WSL struct {
	config *Config
	db     *sql.DB
}

func (this *WSL) ConnectToDb() error {
	if this.db == nil {
		db, err := sql.Open(this.config.DbType, this.config.DbUrl)
		if err != nil {
			return err
		}
		this.db = db
	}
	return nil
}

func (this *WSL) Start() {
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
			return
		}
		qID := urlPath[1]
		fmt.Println(qID)
		fmt.Println(qParams, len(qParams))
	})

	this.ConnectToDb()

	if this.config.httpEnabled() {
		srv := &http.Server{
			Addr:         this.config.HttpAddr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func() {
			log.Fatal(srv.ListenAndServe())
			fmt.Println(fmt.Sprint("Listening on http://", this.config.HttpAddr, "/"))
		}()
	}

	if this.config.httpsEnabled() {
		srvs := &http.Server{
			Addr:         this.config.HttpsAddr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func() {
			log.Fatal(srvs.ListenAndServeTLS(this.config.CertFile, this.config.KeyFile))
			fmt.Println(fmt.Sprint("Listening on https://", this.config.HttpsAddr, "/"))
		}()
	}
}
