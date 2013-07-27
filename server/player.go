package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"math"
)

type payload struct {
	Robots      []robot      `json:"robots"`
	Projectiles []projectile `json:"projectiles"`
}

type player struct {
	ws          *websocket.Conn
	Robot       robot
	send        chan *payload
	Instruction instruction
}

type instruction struct {
	MoveTo position `json:"move_to"`
	FireAt position `json:"fire_at"`
}

func (p *player) sender() {
	for things := range p.send {

		err := websocket.JSON.Send(p.ws, *things)
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
		p.Robot.FireAt = msg.FireAt
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
	newPos := move(p.Robot.Position, p.Robot.MoveTo, *velocity, *delta)
	p.Robot.Position.X = newPos.X
	p.Robot.Position.Y = newPos.Y
}

func (p *projectile) nudge() {
	newPos := move(p.Position, p.MoveTo, *velocity*5, *delta)

	if p.Position.Y-p.MoveTo.Y < 5 && p.Position.X-p.MoveTo.X < 5 {
		delete(g.projectiles, p)

		for player := range g.players {
			if player.Robot.Position.X-p.Position.X < 20 &&
				player.Robot.Position.Y-p.Position.Y < 20 {
				player.Robot.Health -= p.Damage
				log.Printf("Robot %+v is injured", player.Robot)
			}
		}
	}
	p.Position.X = newPos.X
	p.Position.Y = newPos.Y
}

func (p *player) fire() {

	for proj := range g.projectiles {
		if proj.Id == p.Robot.Id {
			return
		}
	}

	proj := &projectile{
		Id:       p.Robot.Id,
		Position: p.Robot.Position,
		MoveTo:   p.Robot.FireAt,
		Damage:   10,
	}
	g.projectiles[proj] = true
}
