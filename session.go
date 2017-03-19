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
	quit    chan struct{}
}

func NewSession(uid, roomID int, conn *websocket.Conn) *Session {
	sess := &Session{
		UID:     uid,
		RoomID:  roomID,
		Conn:    conn,
		MsgChan: make(chan []byte, 100),
		closed:  false,
		quit:    make(chan struct{}, 1),
	}
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

FOR_LOOP:
	for {
		select {
		case msg = <-sess.MsgChan:
			err = sess.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				glog.Errorf("write msg error: %s", err)
				break FOR_LOOP
			}
		case <-sess.quit:
			break FOR_LOOP
		}
	}

	glog.Infof("try to stop sending msg for uid:%d", sess.UID)
}

func (sess *Session) Close() error {
	glog.Infof("closing uid:%d", sess.UID)

	sess.closed = true
	sess.quit <- struct{}{}

	return sess.Conn.Close()
}

func (sess *Session) Run() {
	go sess.goSendMessage()
}
