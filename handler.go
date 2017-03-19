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
	var room *Room
	var msg []byte

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

	// Session 对象可以放入Pool中重用
	sess := server.SessPool.Get().(*Session)
	sess.UID = uid
	sess.RoomID = roomID
	sess.Conn = wsConn
	sess.closed = false
	go sess.Run()

	room = server.CreateRoom(roomID)
	room.AddUser(uid, sess)
	//defer room.RemoveUser(uid)
	defer func() {
		sess.Close()
		room.RemoveUser(uid)
		server.SessPool.Put(sess)
	}()

	for {
		// read message
		_, msg, err = sess.Conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				glog.Error(err)
			}
			break
		}
		glog.Infof("msg:%s", string(msg))

		// broadcast message
		err = room.Broadcast(msg)
		if err != nil {
			glog.Error(err)
		}
	}

	return nil
}
