package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

type player struct {
	ws    *websocket.Conn
	robot robot
	send  chan *[]robot
}

type instruction struct {
	MoveTo position `json:"move_to"`
	FireAt position `json:"fire_at"`
}

func (p *player) sender() {
	for robots := range p.send {
		log.Printf("%s sending state: %T, %+v", p.robot.Name, robots, robots)
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
		log.Printf("waiting for a recv")
		err := websocket.JSON.Receive(p.ws, &msg)
		log.Printf("got a recv")
		if err != nil {
			log.Print(err)
			break
		}
		log.Printf("recv: %s: %+v\n", p.robot.Name, msg)
	}
	p.ws.Close()
}
