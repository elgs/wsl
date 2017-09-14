package wsl

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type WSL struct {
	Config *Config
	db     *sql.DB
}

func New(confFile string) (*WSL, error) {
	config, err := NewConfig(confFile)
	if err != nil {
		return nil, err
	}
	return &WSL{
		Config: config,
	}, nil
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
		if this.Config.Cors {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
			w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		}

		urlPath := strings.Split(r.URL.Path, "/")
		if len(urlPath) < 2 {
			log.Println("Invalid URL.")
			return
		}
		qID := urlPath[1]

		if script, ok := this.Config.Scripts[qID]; ok {
			sepIndex := strings.LastIndex(r.RemoteAddr, ":")
			clientIp := r.RemoteAddr[0:sepIndex]
			clientIp = strings.Replace(strings.Replace(clientIp, "[", "", -1), "]", "", -1)

			r.ParseForm()

			params := valuesToMap(&r.Form)

			params["__client_ip"] = clientIp

			result, err := this.exec(this.db, script, params)
			if err != nil {
				log.Println(err)
				return
			}
			jsonData, err := json.Marshal(result)
			jsonString := string(jsonData)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprint(w, jsonString)
		} else {
			log.Println("Invalid script:", qID)
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
			fmt.Println(fmt.Sprint("Listening on http://", this.Config.HttpAddr, "/"))
			log.Fatal(srv.ListenAndServe())
		}()
	}

	if this.Config.httpsEnabled() {
		srvs := &http.Server{
			Addr:         this.Config.HttpsAddr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func() {
			fmt.Println(fmt.Sprint("Listening on https://", this.Config.HttpsAddr, "/"))
			log.Fatal(srvs.ListenAndServeTLS(this.Config.CertFile, this.Config.KeyFile))
		}()
	}
}
