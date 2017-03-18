package chatroom

import (
	"strconv"

	"github.com/golang/glog"
	"github.com/labstack/echo"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

func (server *ChatServer) Handle(ctx echo.Context) error {
	var uid, roomID int
	var err error
	var wsConn *websocket.Conn

	wsConn, err = upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return err
	}

	uid, err = strconv.Atoi(ctx.QueryParam("uid"))
	if err != nil {
		return err
	}

	roomID, err = strconv.Atoi(ctx.QueryParam("room_id"))
	if err != nil {
		return err
	}

	sess := NewSession(uid, wsConn)
	defer sess.Close()

	server.AddToRoom(roomID, sess)

	for {
		// read message
		_, msg, err := sess.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				glog.Error(err)
			}
			break
		}
		glog.Infof("msg:%s", string(msg))

		// send back for test
		sess.SendMessage(msg)

		// broadcast message

	}

	return nil
}
