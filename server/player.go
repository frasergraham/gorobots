package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"math"
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

const threshold = 1e-7

func move(d1, d2 position, velocity float64, timeDelta float64) position {
	deltaX := float64(d2.X - d1.X)
	deltaY := float64(d2.Y - d1.Y)
	mag := math.Sqrt(deltaX*deltaX + deltaY*deltaY)
	if math.Abs(mag) < threshold {
		return d1
	}
	moveX := deltaX / mag * velocity
	moveY := deltaY / mag * velocity
	return position{
		int(float64(d1.X) + moveX),
		int(float64(d1.Y) + moveY),
	}
}

func (p *player) nudge() {
	fmt.Println()
	newPos := move(p.Robot.Position, p.Robot.MoveTo, *velocity, *delta)
	p.Robot.Position.X = newPos.X
	p.Robot.Position.Y = newPos.Y
}
