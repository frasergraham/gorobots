package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"log"
	"net/http"
    "math/rand"
    "time"
)

var addr = flag.String("addr", ":8666", "http service address")

func main() {
    rand.Seed(time.Now().UnixNano())
	flag.Parse()

	http.Handle("/ws/", websocket.Handler(addPlayer))

	go g.run()

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("unable to start server")
	}
}

func addPlayer(ws *websocket.Conn) {
	id := fmt.Sprintf("robot%d", <-g.id)
	log.Printf("sending robot id: %s", id)
    err := websocket.JSON.Send(ws, handshake{id})
    if err != nil {
        log.Fatal(err)
    }
	p := &player{
		Robot: robot{Id: id},
		send:  make(chan *[]robot),
		ws:    ws,
	}
	g.register <- p
	defer func() {
		g.unregister <- p
	}()
	go p.sender()
	p.recv()
	log.Printf("%v has been disconnect from this game\n", p.Robot.Id)
}
