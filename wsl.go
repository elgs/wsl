package wsl

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

		if r.Method == "OPTIONS" {
			return
		}

		urlPath := strings.Split(r.URL.Path, "/")
		if len(urlPath) < 2 {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, fmt.Sprint(`{"err":"Invalid URL"}`))
			return
		}
		qID := urlPath[1]

		if script, ok := this.Config.Scripts[qID]; ok {
			sepIndex := strings.LastIndex(r.RemoteAddr, ":")
			clientIp := r.RemoteAddr[0:sepIndex]
			clientIp = strings.Replace(strings.Replace(clientIp, "[", "", -1), "]", "", -1)

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, fmt.Sprint(`{"err":"`, err, `"}`))
				log.Println(err)
				return
			}
			var bodyData map[string]string
			json.Unmarshal(body, &bodyData)
			//intentionally ignore the errors

			paramValues, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, fmt.Sprint(`{"err":"`, err, `"}`))
				log.Println(err)
				return
			}
			params := valuesToMap(paramValues)
			for k, v := range bodyData {
				params[k] = v
			}

			params["__client_ip"] = clientIp

			headers := valuesToMap(r.Header)

			result, err := this.exec(qID, this.db, script, params, headers)
			if err != nil {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, fmt.Sprint(`{"err":"`, err, `"}`))
				log.Println(err)
				return
			}
			jsonData, err := json.Marshal(result)
			if err != nil {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, fmt.Sprint(`{"err":"`, err, `"}`))
				log.Println(err)
				return
			}
			jsonString := string(jsonData)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprint(w, jsonString)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, fmt.Sprint(`{"err":"`, `Invalid script:`, qID, `"}`))
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
