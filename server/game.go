package main

import (
	"fmt"
	"log"
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
		case p := <-g.register:
			log.Printf("adding player: %+v", p.robot.Name)
			g.players[p] = true
		case p := <-g.unregister:
			delete(g.players, p)
			close(p.send)
		case <-time.Tick(10 * time.Second):
			fmt.Printf("\n\n\n")
			log.Printf("calculating state")

			robots := []robot{}
			for p := range g.players {
				p.robot.Position.X = p.instruction.MoveTo.X
				p.robot.Position.Y = p.instruction.MoveTo.Y
				robots = append(robots, p.robot)
			}

			for p := range g.players {
				p.send <- &robots
			}
		}
	}
}
