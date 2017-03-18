package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	u, err := url.Parse("ws://localhost:1300/chatroom?uid=1&room_id=1")
	if err != nil {
		log.Println(err)
	}

	rawConn, err := net.Dial("tcp", u.Host)
	if err != nil {
		log.Println(err)
	}

	wsHeaders := http.Header{
		"Origin": {"http://localhost:1300"},
	}

	wsConn, resp, err := websocket.NewClient(rawConn, u, wsHeaders, 1024, 1024)
	if err != nil {
		fmt.Printf("websocket.NewClient Error: %s Resp:%+v", err, resp)
		return
	}

	err = wsConn.WriteMessage(websocket.TextMessage, []byte("test msg"))
	if err != nil {
		log.Println(err)
	}

	_, msg, err := wsConn.ReadMessage()
	if err != nil {
		log.Println(err)
	}
	log.Println("msg: ", string(msg))
}
