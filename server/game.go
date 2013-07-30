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
	turn        int
}

type handshake struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

type config struct {
	ID    string `json:"id"`
	Stats stats  `json:"stats"`
}

type payload struct {
	Robots      []robot      `json:"robots"`
	Projectiles []projectile `json:"projectiles"`
	Splosions   []splosion   `json:"splosions"`
	Reset       bool         `json:"reset"`
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
			g.turn++
			t0 := time.Now()
			if *verbose {
				log.Printf("\033[2JTurn: %v", g.turn)
				log.Printf("Players: %v", len(g.players))
				log.Printf("Projectiles: %v", len(g.projectiles))
				log.Printf("Explosions: %v", len(g.splosions))
			}
			payload := payload{}
			payload.Robots = []robot{}
			payload.Projectiles = []projectile{}

			robots_remaining := 0

			for p := range g.players {
				if p.Robot.Health > 0 {
					robots_remaining++
					p.scan()
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

			if robots_remaining <= 1 && len(g.players) > 1 {
				for p := range g.players {
					if p.Robot.Health > 0 {
						log.Printf("Robot %v Wins", p.Robot.Id)
					}
					p.reset()
				}
				payload.Reset = true
			} else {
				payload.Reset = false
			}

			t1 := time.Now()
			if *verbose {
				log.Printf("Turn Processes %v\n", t1.Sub(t0))
			}

			for p := range g.players {
				p.send <- &payload
			}

			t1 = time.Now()
			if *verbose {
				log.Printf("Sent Payload %v\n", t1.Sub(t0))
			}

		}
	}
}
