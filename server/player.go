package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"math/rand"
)

type player struct {
	ws          *websocket.Conn
	Robot       robot
	send        chan *boardstate
	Instruction instruction
}

type instruction struct {
	MoveTo position `json:"move_to"`
	FireAt position `json:"fire_at"`
	Stats  stats    `json:"stats"`
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
		if msg.MoveTo.X != 0 && msg.MoveTo.Y != 0 {
			p.Robot.MoveTo = msg.MoveTo
		}
		if msg.FireAt.X != 0 && msg.FireAt.Y != 0 {
			p.Robot.FireAt = msg.FireAt
		}
		if msg.Stats.Speed > 0 {
			p.Robot.Stats = msg.Stats
			p.Robot.Health = p.Robot.Stats.Hp
			log.Printf("%+v", p.Robot.Stats)
		}
	}
	p.ws.Close()
}

func (p *player) nudge() {
	newPos := move(p.Robot.Position, p.Robot.MoveTo, p.Robot.Stats.Speed, delta)
	p.Robot.Position.X = newPos.X
	p.Robot.Position.Y = newPos.Y
}

func (p *player) scan() {
	p.Robot.Scanners = make([]scanner, 0)
	for player := range g.players {
		if player.Robot.Id == p.Robot.Id || player.Robot.Health <= 0 {
			continue
		}
		dist := distance(player.Robot.Position, p.Robot.Position)
		if dist < float64(p.Robot.Stats.ScannerRadius) {
			s := scanner{
				Position: position{
					X: player.Robot.Position.X,
					Y: player.Robot.Position.Y,
				},
			}
			p.Robot.Scanners = append(p.Robot.Scanners, s)
		}
	}
}

func (p *player) fire() {

	for proj := range g.projectiles {
		if proj.Id == p.Robot.Id {
			return
		}
	}

	// log.Printf("%v Fired at %v %v", p.Robot.Id, p.Robot.FireAt.X, p.Robot.FireAt.Y)

	proj := &projectile{
		Id:       p.Robot.Id,
		Position: p.Robot.Position,
		MoveTo:   p.Robot.FireAt,
		Damage:   10,
		Radius:   p.Robot.Stats.WeaponRadius,
		Speed:    float64(p.Robot.Stats.Speed * 2),
	}
	g.projectiles[proj] = true
}

func (p *player) reset() {
	start_pos := position{
		X: rand.Float64() * 800,
		Y: rand.Float64() * 550,
	}
	p.Robot.MoveTo = start_pos
	p.Robot.Position = start_pos
	p.Robot.Health = p.Robot.Stats.Hp
}
