package chatroom

import (
	"flag"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	ListenAddr string
)

func init() {
	flag.StringVar(&ListenAddr, "listenAddr", ":1300", "default 0.0.0.0:1300")
}

type ChatServer struct {
	HTTPServer *echo.Echo
	Rooms      map[int]*Room
}

func NewChatServer() *ChatServer {
	s := &ChatServer{
		HTTPServer: echo.New(),
		Rooms:      make(map[int]*Room),
	}
	return s
}

func (server *ChatServer) Start() {
	server.HTTPServer.GET("/chatroom", server.Handle, middleware.Recover())
	log.Fatal(server.HTTPServer.Start(ListenAddr))
}

func (server *ChatServer) HasRoom(id int) bool {
	return false
}

func (server *ChatServer) AddToRoom(roomID int, sess *Session) {

}

func (server *ChatServer) PutMsgToRoom() {

}
