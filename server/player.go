package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

type player struct {
	ws          *websocket.Conn
	robot       robot
	send        chan *[]robot
	instruction instruction
}

type instruction struct {
	MoveTo position `json:"move_to"`
	FireAt position `json:"fire_at"`
}

func (p *player) sender() {
	for robots := range p.send {
		log.Printf("%s sending %d robots", p.robot.Name, len(*robots))
		err := websocket.JSON.Send(p.ws, *robots)
		if err != nil {
			break
		}
		log.Printf("%s: state sent", p.robot.Name)
	}
	p.ws.Close()
}

func (p *player) recv() {
	for {
		var msg instruction
		log.Printf("recv: %s waiting for a recv", p.robot.Name)
		err := websocket.JSON.Receive(p.ws, &msg)
		if err != nil {
			log.Print(err)
			break
		}
		p.instruction = msg
		log.Printf("recv: %s: %+v %v\n", p.robot.Name, p.instruction, msg)
	}
	p.ws.Close()
}
