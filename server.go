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

func (this *WSL) defaultHandler(w http.ResponseWriter, r *http.Request) {
	if this.Config.Web.Cors {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.Header().Set("Access-Control-Expose-Headers", "Token")
	}

	if r.Method == "OPTIONS" {
		return
	}

	urlPath := strings.Split(r.URL.Path, "/")
	if len(urlPath) < 2 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Sprint(`{"err":"invalid_url"}`))
		return
	}
	qID := urlPath[1]

	if os.Getenv("env") == "dev" {
		this.LoadScripts(qID)
	}
	script := this.Scripts[qID]

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
	var bodyData map[string]interface{}
	//intentionally ignore the errors
	_ = json.Unmarshal(body, &bodyData)

	paramValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, fmt.Sprint(`{"err":"`, err, `"}`))
		log.Println(err)
		return
	}
	params := valuesToMap(false, paramValues)
	if queryParams, ok := params["params"]; ok {
		if ps, ok := queryParams.(string); ok {
			params["params"] = ConvertStringArrayToInterfaceArray(strings.Split(ps, ","))
		}
	}
	for k, v := range bodyData {
		params[k] = v
	}

	params["__client_ip"] = clientIp

	context := map[string]interface{}{}

	headers := valuesToMap(true, r.Header)
	authHeader := headers["access_token"]

	if authHeader == nil || authHeader == "" {
		authHeader = params["access_token"]
	}

	if authHeader != nil || authHeader == "" {
		context["access_token"] = authHeader
	}
	result, err := this.exec(qID, this.Databases["main"], script, params, context)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, fmt.Sprint(`{"err":"`, err, `"}`))
		log.Println(err)
		return
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, fmt.Sprint(`{"err":"`, err, `"}`))
		log.Println(err)
		return
	}
	jsonString := string(jsonData)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, jsonString)
}

// const (
// 	// Time allowed to write a message to the peer.
// 	writeWait = 10 * time.Second

// 	// Time allowed to read the next pong message from the peer.
// 	pongWait = 60 * time.Second

// 	// Send pings to peer with this period. Must be less than pongWait.
// 	pingPeriod = (pongWait * 9) / 10

// 	// Maximum message size allowed from peer.
// 	// maxMessageSize = 512
// )

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin:     func(r *http.Request) bool { return true },
// }

// func (this *WSL) readWs(conn *websocket.Conn, m chan []byte, clientIp string, ins *WSL) {
// 	defer func() {
// 		log.Println("read connection closed.")
// 		conn.Close()
// 	}()
// 	// conn.SetReadLimit(maxMessageSize)
// 	conn.SetReadDeadline(time.Now().Add(pongWait))
// 	conn.SetPongHandler(func(string) error {
// 		conn.SetReadDeadline(time.Now().Add(pongWait))
// 		// log.Println("pong received.")
// 		return nil
// 	})
// 	for {
// 		_, message, err := conn.ReadMessage()
// 		conn.SetReadDeadline(time.Now().Add(pongWait))
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
// 				log.Println(err)
// 			}
// 			break
// 		}

// 		var input map[string]interface{}

// 		err = json.Unmarshal(message, &input)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		query := input["query"]
// 		if query == nil {
// 			log.Println("Invalid query.")
// 			return
// 		}
// 		qID := query.(string)
// 		script := this.Config.Db.Scripts[qID]

// 		params, err := ConvertMapOfInterfacesToMapOfStrings(input["data"].(map[string]interface{}))
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		params["__client_ip"] = clientIp

// 		context := map[string]interface{}{}

// 		authHeader := input["access_token"]
// 		if authHeader != nil {
// 			context["access_token"] = authHeader
// 		}

// 		result, err := ins.exec(qID, this.db, script, params, context)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		ret := make(map[string]interface{})
// 		ret["data"] = result
// 		if tokenString, ok := context["token"]; ok {
// 			ret["token"] = tokenString.(string)
// 		}

// 		jsonData, err := json.Marshal(ret)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		m <- jsonData
// 	}
// }

// func (this *WSL) writeWs(conn *websocket.Conn, m chan []byte) {
// 	defer func() {
// 		log.Println("write connection closed.")
// 		conn.Close()
// 	}()
// 	for {
// 		select {
// 		case message, ok := <-m:
// 			conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if !ok {
// 				conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}
// 			// log.Println("message received.")
// 			err := conn.WriteMessage(websocket.TextMessage, message)
// 			if err != nil {
// 				log.Println(err)
// 				break
// 			}
// 		case <-time.After(pingPeriod):
// 			// wait for pingPeriod time of inactivitity, then send a ping, disconnect if pong is not received within writeWait.
// 			conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			// log.Println("ping sent.")
// 			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
// 				log.Println(err)
// 				return
// 			}
// 		}
// 	}
// }

// func (this *WSL) wsHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "OPTIONS" {
// 		return
// 	}
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprint(w, err)
// 		log.Println(err)
// 		return
// 	}

// 	sepIndex := strings.LastIndex(r.RemoteAddr, ":")
// 	clientIp := r.RemoteAddr[0:sepIndex]
// 	clientIp = strings.Replace(strings.Replace(clientIp, "[", "", -1), "]", "", -1)

// 	// fmt.Println(clientIp)
// 	// log.Println("Connected")

// 	m := make(chan []byte)
// 	go this.readWs(conn, m, clientIp, this)
// 	go this.writeWs(conn, m)
// }
