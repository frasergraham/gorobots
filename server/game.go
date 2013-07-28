package main

import (
	"log"
	"sort"
	"time"
)

type game struct {
	players     map[*player]bool
	projectiles map[*projectile]bool
	splosions   map[*splosion]bool
	register    chan *player
	unregister  chan *player
	id          chan int
}

type handshake struct {
	ID string `json:"id"`
}

var g = game{
	register:    make(chan *player),
	unregister:  make(chan *player),
	projectiles: make(map[*projectile]bool),
	splosions:   make(map[*splosion]bool),
	players:     make(map[*player]bool),
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

func (g *game) run() {
	g.id = make(chan int)
	go func() {
		for i := 0; ; i++ {
			g.id <- i
		}
	}()

	for {
		select {
		case p := <-g.register:
			log.Printf("adding player: %+v", p.Robot.Id)
			g.players[p] = true
		case p := <-g.unregister:
			delete(g.players, p)
			close(p.send)
		case <-time.Tick(time.Duration(*tick) * time.Millisecond):
			payload := payload{}
			payload.Robots = []robot{}
			payload.Projectiles = []projectile{}

			for p := range g.players {
				if p.Robot.Health > 0 {
					p.nudge()
					if p.Robot.FireAt.X != 0 && p.Robot.FireAt.Y != 0 {
						p.fire()
					}
				}
				payload.Robots = append(payload.Robots, p.Robot)
			}
			sort.Sort(robotSorter{payload.Robots})

			for p := range g.projectiles {
				p.nudge()
				payload.Projectiles = append(payload.Projectiles, *p)
			}

			for s := range g.splosions {
				s.tick()
				payload.Splosions = append(payload.Splosions, *s)
			}

			for p := range g.players {
				p.send <- &payload
			}
		}
	}
}
