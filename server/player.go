package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

type player struct {
	name  string
	ws    *websocket.Conn
	robot robot
	send  chan *[]robot
}

type instruction struct {
	Destination position
	FireAt      position `json:"fire_at"`
}

func (p *player) sender() {
    for robots := range p.send {
		err := websocket.JSON.Send(p.ws, robots)
		if err != nil {
			break
		}
	}
	p.ws.Close()
}

func (p *player) recv() {
	for {
		var msg instruction
        log.Printf("recv: %s: %+v\n", p.name, msg)
		err := websocket.JSON.Receive(p.ws, &msg)
		if err != nil {
            log.Print(err)
			break
		}
	}
	p.ws.Close()
}
