package main

import (
	"log"
)

type weapon struct {
	Strength float64 `json:"strength"`
	Radius   float64 `json:"radius"`
}

type stats struct {
	Speed         float64 `json:"speed"`
	Hp            int     `json:"hp"`
	WeaponRadius  int     `json:"weapon_radius"`
	ScannerRadius int     `json:"scanner_radius"`
}

func (s stats) valid() bool {
	total := int(s.Speed) + s.Hp + s.WeaponRadius + s.ScannerRadius
	log.Printf("total: %d", total)
	if total > 500 {
		return false
	}
	return true
}

type scanner struct {
	Position position `json:"position"`
	Stats    stats    `json:"stats"`
}

type robot struct {
	Id       string    `json:"id"`
	Stats    stats     `json:"stats"`
	Health   int       `json:"health"`
	Position position  `json:"position"`
	MoveTo   position  `json:"move_to"`
	FireAt   position  `json:"fire_at"`
	Scanners []scanner `json:"scanners"`
}

type robotSorter struct {
	robots []robot
}

func (s robotSorter) Len() int {
	return len(s.robots)
}

func (s robotSorter) Swap(i, j int) {
	s.robots[i], s.robots[j] = s.robots[j], s.robots[i]
}

func (s robotSorter) Less(i, j int) bool {
	return s.robots[i].Id < s.robots[j].Id
}

type projectile struct {
	Id       string   `json:"id"`
	Position position `json:"position"`
	MoveTo   position `json:"move_to"`
	Radius   int      `json:"radius"`
	Speed    float64  `json:"speed"`
	Damage   int      `json:"damage"`
}

func (p *projectile) nudge() {
	newPos := move(p.Position, p.MoveTo, float64(p.Speed), delta)

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

type splosion struct {
	Id        string   `json:"id"`
	Position  position `json:"position"`
	Radius    int      `json:"radius"`
	MaxDamage int      `json:"damage"`
	MinDamage int      `json:"damage"`
	Lifespan  int      `json:"lifespan"`
}

func (s *splosion) tick() {
	s.Lifespan--
	if s.Lifespan <= 0 {
		delete(g.splosions, s)
	}
}
