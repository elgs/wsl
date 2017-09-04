package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var addr = flag.String("addr", ":8080", "http service address, default to :8080")

func main() {
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

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))

		urlPath := strings.Split(r.URL.Path, "/")
		urlQuery, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(urlPath, len(urlPath))
		fmt.Println(urlQuery, len(urlQuery))
	})

	srv := &http.Server{
		Addr:         *addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		log.Fatal(srv.ListenAndServe())
		log.Println("started!")
	}()

	<-done
	fmt.Println("Bye!")
}
