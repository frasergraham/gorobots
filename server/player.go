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
		err := websocket.JSON.Receive(p.ws, &msg)
		if err != nil {
			log.Print(err)
			break
		}
		p.Robot.MoveTo = msg.MoveTo
	}
	p.ws.Close()
}

func (p *player) nudge() {
	switch {
	case p.Robot.Position.X < p.Robot.MoveTo.X:
		p.Robot.Position.X += 1
	case p.Robot.Position.X > p.Robot.MoveTo.X:
		p.Robot.Position.X -= 1
	}
	switch {
	case p.Robot.Position.Y < p.Robot.MoveTo.Y:
		p.Robot.Position.Y += 1
	case p.Robot.Position.Y > p.Robot.MoveTo.Y:
		p.Robot.Position.Y -= 1
	}
}
