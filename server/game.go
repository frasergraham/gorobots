package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type game struct {
	players    map[*player]bool
	register   chan *player
	unregister chan *player
	id         chan int
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
		case c := <-g.register:
			log.Printf("adding player: %+v", c)
			g.players[c] = true
		case c := <-g.unregister:
			delete(g.players, c)
			close(c.send)
		case <-time.Tick(2 * time.Second):
			fmt.Printf("\n\n\n")
			log.Printf("timing out, time to calculate state")

			robots := []robot{}
			for p := range g.players {
				p.robot.Position.X = rand.Intn(300)
				p.robot.Position.Y = rand.Intn(300)
				robots = append(robots, p.robot)
			}

			for p := range g.players {
				p.send <- &robots
			}
		}
	}
}
