package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"

	"github.com/silentred/chatroom"
)

func main() {
	flag.Parse()
	debug()

	s := chatroom.NewChatServer()
	s.Start()
}

func debug() {
	go func() {
		http.HandleFunc("/go", func(w http.ResponseWriter, r *http.Request) {
			num := strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
			w.Write([]byte(num))
		})
		http.ListenAndServe(":6060", nil)
	}()
}
