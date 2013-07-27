package main

import (
	"log"
	"time"
)

type game struct {
	players    map[*player]bool
	register   chan *player
	unregister chan *player
	id         chan int
}

type handshake struct {
	ID string `json:"id"`
}

var g = game{
	register:   make(chan *player),
	unregister: make(chan *player),
	players:    make(map[*player]bool),
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
		case <-time.Tick(1 * time.Second):
			// fmt.Printf("\n\n\n")
			// log.Printf("calculating state")

			robots := []robot{}
			for p := range g.players {
				p.nudge()
				robots = append(robots, p.Robot)
			}

			for p := range g.players {
				p.send <- &robots
			}
		}
	}
}
