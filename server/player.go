package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"math"
	"math/rand"
)

type payload struct {
	Robots      []robot      `json:"robots"`
	Projectiles []projectile `json:"projectiles"`
	Splosions   []splosion   `json:"splosions"`
	Reset       bool         `json:"reset"`
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

const threshold = 1e-7

func distance(p1, p2 position) float64 {
	deltaX := float64(p2.X - p1.X)
	deltaY := float64(p2.Y - p1.Y)
	mag := math.Abs(math.Sqrt(deltaX*deltaX + deltaY*deltaY))
	return mag
}

func move(d1, d2 position, velocity float64, timeDelta float64) position {
	deltaX := float64(d2.X - d1.X)
	deltaY := float64(d2.Y - d1.Y)
	mag := math.Sqrt(deltaX*deltaX + deltaY*deltaY)
	if math.Abs(mag) < threshold {
		return d1
	}

	distance_this_frame := velocity * timeDelta
	// log.Printf("%+v  ->  %+v   :  %+v", d1, d2, distance_this_frame)

	if math.Abs(mag) < distance_this_frame {
		return d2
	}

	moveX := 0.0
	moveY := 0.0

	if deltaY > 0 {
		theta := math.Atan(deltaX / deltaY)
		moveX = math.Sin(theta) * distance_this_frame
		moveY = math.Cos(theta) * distance_this_frame
	} else {
		theta := math.Atan(deltaX / deltaY)
		moveX = -math.Sin(theta) * distance_this_frame
		moveY = -math.Cos(theta) * distance_this_frame
	}

	// log.Printf("%+v , %+v", moveX, moveY)
	return position{
		d1.X + moveX,
		d1.Y + moveY,
	}
}

func (p *player) nudge() {
	newPos := move(p.Robot.Position, p.Robot.MoveTo, p.Robot.Stats.Speed, *delta)
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

func (s *splosion) tick() {
	s.Lifespan--
	if s.Lifespan <= 0 {
		delete(g.splosions, s)
	}
}

func (p *projectile) nudge() {
	newPos := move(p.Position, p.MoveTo, float64(p.Speed), *delta)

	hit_player := false
	for player := range g.players {
		if player.Robot.Id == p.Id {
			continue
		}
		dist := distance(player.Robot.Position, p.Position)
		if dist < 5.0 {
			hit_player = true
		}
	}

	if distance(p.Position, p.MoveTo) < 5 || hit_player {
		delete(g.projectiles, p)

		// Spawn a splosion
		splo := &splosion{
			Id:        p.Id,
			Position:  p.Position,
			Radius:    p.Radius,
			MaxDamage: 10,
			MinDamage: 5,
			Lifespan:  8,
		}
		g.splosions[splo] = true

		for player := range g.players {
			dist := distance(player.Robot.Position, p.Position)
			if dist < float64(p.Radius) {

				// TODO map damage Max to Min based on distance from explosion
				if player.Robot.Health > 0 {
					player.Robot.Health -= p.Damage
					// log.Printf("Robot %+v is injured", player.Robot)
					if player.Robot.Health <= 0 {
						// log.Printf("Robot %+v is dead", player.Robot)
					}
				}
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
