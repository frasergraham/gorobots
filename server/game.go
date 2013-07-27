package main

import (
	"log"
	"sort"
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
		case <-time.Tick(10 * time.Millisecond):
			robots := []robot{}
			for p := range g.players {
				p.nudge()
				robots = append(robots, p.Robot)
			}
			sort.Sort(robotSorter{robots})

			for p := range g.players {
				p.send <- &robots
			}
		}
	}
}
