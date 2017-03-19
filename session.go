package chatroom

import (
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type Session struct {
	UID     int
	RoomID  int
	Conn    *websocket.Conn
	MsgChan chan []byte
	closed  bool
}

func NewSession(uid, roomID int, conn *websocket.Conn) *Session {
	sess := &Session{
		UID:     uid,
		RoomID:  roomID,
		Conn:    conn,
		MsgChan: make(chan []byte, 100),
		closed:  false,
	}
	sess.Start()
	return sess
}

func (sess *Session) SendMessage(msg []byte) {
	if !sess.closed {
		sess.MsgChan <- msg
	}
}

func (sess *Session) goSendMessage() {
	var err error
	var msg []byte
	for msg = range sess.MsgChan {
		err = sess.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			glog.Errorf("write msg error: %s", err)
			//sess.Close()
			break
		}
	}
	glog.Infof("try to stop sending msg for uid:%d", sess.UID)
}

func (sess *Session) Close() error {
	glog.Infof("closing uid:%d", sess.UID)

	sess.closed = true
	close(sess.MsgChan)

	return sess.Conn.Close()
}

func (sess *Session) Start() {
	go sess.goSendMessage()
}
