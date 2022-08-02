package wsl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func (this *App) defaultHandler(w http.ResponseWriter, r *http.Request) {
	if this.Config.Web.Cors {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		// w.Header().Set("Access-Control-Expose-Headers", "Token")
	}

	if r.Method == "OPTIONS" {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	script, err := this.GetScript(r.URL.Path[1:], os.Getenv("env") == "dev")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Sprintf(`{"err":"%v"}`, err))
		return
	}

	sepIndex := strings.LastIndex(r.RemoteAddr, ":")
	clientIP := r.RemoteAddr[0:sepIndex]
	clientIP = strings.Replace(strings.Replace(clientIP, "[", "", -1), "]", "", -1)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Sprintf(`{"err":"%v"}`, err))
		log.Println(err)
		return
	}
	var bodyData map[string]any
	err = json.Unmarshal(body, &bodyData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Sprintf(`{"err":"%v"}`, err))
		log.Println(err)
		return
	}

	paramValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Sprintf(`{"err":"%v"}`, err))
		log.Println(err)
		return
	}

	params := valuesToMap(false, paramValues)
	for k, v := range bodyData {
		params[k] = v
	}

	authHeader := r.Header.Get("access_token")
	if authHeader == "" && params["access_token"] != nil {
		if token, ok := params["access_token"].(string); ok {
			authHeader = token
		}
	}

	context := &Context{
		App:         this,
		Script:      script,
		AccessToken: authHeader,
		ClientIP:    clientIP,
		Params:      params,
		Opt:         map[string]any{},
	}

	result, err := this.exec(context)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Sprintf(`{"err":"%v"}`, err))
		log.Println(err)
		return
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Sprintf(`{"err":"%v"}`, err))
		log.Println(err)
		return
	}
	jsonString := string(jsonData)
	fmt.Fprint(w, jsonString)
}
