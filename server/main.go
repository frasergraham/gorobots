package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"net/http"
    "math/rand"
    "time"
)

var addr = flag.String("addr", ":8666", "http service address")

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	http.Handle("/ws/", websocket.Handler(addPlayer))

    go g.run()

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("unable to start server")
	}
}
