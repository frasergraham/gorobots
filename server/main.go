package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8666", "http service address")

func main() {
	flag.Parse()

	http.Handle("/ws/", websocket.Handler(addPlayer))

	go g.run()

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("unable to start server")
	}
}

func addPlayer(ws *websocket.Conn) {
	name := fmt.Sprintf("robot%d", <-g.id)
	log.Printf("adding robot: %s", name)
	p := &player{
		robot: robot{Name: name},
		send:  make(chan *[]robot),
		ws:    ws,
	}
	g.register <- p
	defer func() {
		g.unregister <- p
	}()
	go p.sender()
	p.recv()
	log.Printf("%v has been disconnect from this game\n", p.robot.Name)
}
