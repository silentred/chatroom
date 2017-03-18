package chatroom

type Room struct {
	ID    int
	Users map[int]*Session
}

func NewRoom(id int) *Room {
	return nil
}

func AddUser(uid int) {

}

func RemoveUser(uid int) {

}

func HasUser(uid int) {

}

func GetUser(uid int) {

}

func (r *Room) Broadcast(msg []byte) error {
	for _, sess := range r.Users {
		sess.SendMessage(msg)
	}
	return nil
}
