package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

type player struct {
	ws          *websocket.Conn
	Robot       robot
	send        chan *[]robot
	Instruction instruction
}

type instruction struct {
	MoveTo position `json:"move_to"`
	FireAt position `json:"fire_at"`
}

func (p *player) sender() {
	for robots := range p.send {
		log.Printf("%s sending %d robots", p.Robot.Id, len(*robots))
		err := websocket.JSON.Send(p.ws, *robots)
		if err != nil {
			break
		}
		log.Printf("%s: state sent", p.Robot.Id)
	}
	p.ws.Close()
}

func (p *player) recv() {
	for {
		var msg instruction
		log.Printf("recv: %s waiting for a recv", p.Robot.Id)
		err := websocket.JSON.Receive(p.ws, &msg)
		log.Printf("recv: %s: %+v\n", p.Robot.Id, msg)
		if err != nil {
			log.Print(err)
			break
		}
	}
	p.ws.Close()
}
