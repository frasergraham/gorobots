package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type game struct {
	players    map[*player]bool
	status     chan *robot
	register   chan *player
	unregister chan *player
	id         chan int
}

var g = game{
	status:     make(chan *robot),
	register:   make(chan *player),
	unregister: make(chan *player),
	players:    make(map[*player]bool),
}

func (g *game) run() {
	g.id = make(chan int, 16)
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
		case p := <-g.status:
			log.Printf("%+v", p)
		case <-time.Tick(2 * time.Second):
			log.Printf("timing out, time to calculate state")

			robots := []robot{}
			for p := range g.players {
				p.robot.Position.X = rand.Intn(300)
				p.robot.Position.Y = rand.Intn(300)
				robots = append(robots, p.robot)
			}

			for p := range g.players {
				log.Printf("sending state: %+v", p)
				p.send <- &robots
			}
		}
	}
}

func addPlayer(ws *websocket.Conn) {
	name := fmt.Sprintf("robot%d", <-g.id)
	log.Printf("adding robot: %s", name)
	p := &player{
		name: name,
		send: make(chan *[]robot, 256),
		ws:   ws,
	}
	g.register <- p
	defer func() {
		g.unregister <- p
	}()
	go p.sender()
	p.recv()
	fmt.Printf("%v has been disconnect from this game\n", p)
}
