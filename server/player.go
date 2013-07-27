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
		err := websocket.JSON.Send(p.ws, *robots)
		if err != nil {
			break
		}
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

func move(d1, d2 position, velocity float64, timeDelta float64) position {
    deltaX := float64(d2.X - d1.X) / velocity * timeDelta
	deltaY := float64(d2.Y - d1.Y) / velocity * timeDelta
	return position{
        int(deltaX), int(deltaY),
    }
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
