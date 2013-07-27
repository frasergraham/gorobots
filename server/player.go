package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
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

func (p *projectile) nudge() {
	switch {
	case p.Position.X < p.MoveTo.X:
		p.Position.X += 5
	case p.Position.X > p.MoveTo.X:
		p.Position.X -= 5
	}
	switch {
	case p.Position.Y < p.MoveTo.Y:
		p.Position.Y += 5
	case p.Position.Y > p.MoveTo.Y:
		p.Position.Y -= 5
	}

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
