package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"

	"flag"

	"os"

	"github.com/gorilla/websocket"
)

var (
	send      bool
	msg       string
	uid       int
	roomID    int
	times     int
	clientNum int
)

func init() {
	flag.BoolVar(&send, "send", false, "SEND: bool")
	flag.StringVar(&msg, "msg", "test", "SEND: msg content")
	flag.IntVar(&times, "times", 1, "SEND: write times")
	flag.IntVar(&uid, "uid", 1, "uid")
	flag.IntVar(&roomID, "room", 1, "roomID")
	flag.IntVar(&clientNum, "cliNum", 1, "RECV: client number")
}

func main() {
	var forever chan int
	flag.Parse()

	if send {
		conn := newConn(uid, roomID)
		sendMsg(conn)
		return
	}

	for index := 0; index < clientNum; index++ {
		conn := newConn(uid+index, roomID)
		go recvMsg(conn, index)
	}
	<-forever
}

func newConn(uid, roomID int) *websocket.Conn {
	u, err := url.Parse(fmt.Sprintf("ws://localhost:1300/chatroom?uid=%d&room_id=%d", uid, roomID))
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
		return nil
	}

	return wsConn
}

func recvMsg(wsConn *websocket.Conn, index int) {
	for {
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			log.Println(index, err)
			break
		}
		fmt.Fprintf(os.Stdout, "msg:%s idx:%d \n", msg, index)
	}
}

func sendMsg(wsConn *websocket.Conn) {
	for i := 0; i < times; i++ {
		err := wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println(err)
			return
		}
	}
	return
}
