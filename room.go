package chatroom

import (
	"sync"
)

type Room struct {
	ID    int
	Users map[int]*Session
	mutex *sync.RWMutex
}

func NewRoom(id int) *Room {
	r := &Room{
		ID:    id,
		Users: make(map[int]*Session),
		mutex: new(sync.RWMutex),
	}
	return r
}

func (r *Room) AddUser(uid int, sess *Session) {
	r.mutex.Lock()
	if !r.HasUser(uid) {
		r.Users[uid] = sess

	}
	r.mutex.Unlock()
}

func (r *Room) RemoveUser(uid int) {
	r.mutex.Lock()
	if r.HasUser(uid) {
		delete(r.Users, uid)
	}
	r.mutex.Unlock()
}

func (r *Room) HasUser(uid int) bool {
	_, has := r.Users[uid]
	return has
}

func (r *Room) GetUser(uid int) *Session {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	sess, has := r.Users[uid]
	if has {
		return sess
	}
	return nil
}

func (r *Room) Broadcast(msg []byte) error {
	r.mutex.RLock()
	for _, sess := range r.Users {
		sess.SendMessage(msg)
	}
	r.mutex.RUnlock()
	return nil
}
