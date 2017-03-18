package chatroom

import (
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type Session struct {
	UID     int
	Conn    *websocket.Conn
	MsgChan chan []byte
	closed  bool
}

func NewSession(uid int, conn *websocket.Conn) *Session {
	sess := &Session{
		UID:     uid,
		Conn:    conn,
		MsgChan: make(chan []byte, 100),
		closed:  false,
	}
	sess.Start()
	return sess
}

func (sess *Session) ReadMessage() (int, []byte, error) {
	return sess.Conn.ReadMessage()
}

func (sess *Session) SendMessage(msg []byte) {
	if !sess.closed {
		sess.MsgChan <- msg
	}
}

func (sess *Session) goSendMessage() {
	var err error
	for msg := range sess.MsgChan {
		err = sess.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			glog.Errorf("write msg error: %s", err)
			sess.Close()
		}
	}
	glog.Infof("try to stop sending msg for uid:%d", sess.UID)
}

func (sess *Session) Close() error {
	close(sess.MsgChan)
	sess.closed = true
	return sess.Conn.Close()
}

func (sess *Session) Start() {
	go sess.goSendMessage()
}
