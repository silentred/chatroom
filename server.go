package chatroom

import (
	"flag"
	"log"

	"sync"

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
	SessPool   *sync.Pool
	mutex      *sync.RWMutex
}

func NewChatServer() *ChatServer {
	s := &ChatServer{
		HTTPServer: echo.New(),
		Rooms:      make(map[int]*Room),
		mutex:      new(sync.RWMutex),
		SessPool: &sync.Pool{
			New: func() interface{} { return NewSession(0, 0, nil) },
		},
	}
	return s
}

func (server *ChatServer) Start() {
	server.HTTPServer.GET("/chatroom", server.Handle, middleware.Recover())
	log.Fatal(server.HTTPServer.Start(ListenAddr))
}

func (server *ChatServer) HasRoom(id int) bool {
	_, has := server.Rooms[id]
	return has
}

func (server *ChatServer) CreateRoom(id int) *Room {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	if !server.HasRoom(id) {
		server.Rooms[id] = NewRoom(id)
	}
	return server.Rooms[id]
}

func (server *ChatServer) GetRoom(id int) *Room {
	server.mutex.RLock()
	defer server.mutex.RUnlock()
	room, has := server.Rooms[id]
	if has {
		return room
	}
	return nil
}
