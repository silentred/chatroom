package main

import (
	"flag"

	"github.com/silentred/chatroom"
)

func main() {
	flag.Parse()

	s := chatroom.NewChatServer()
	s.Start()
}
