package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8666", "http service address")
var velocity = flag.Float64("velocity", 30, "")
var tick = flag.Int("tick", 33, "")
var weapon_radius = flag.Int("weapon_radius", 35, "")
var verbose = flag.Bool("verbose", false, "")

var delta float64
var g = game{
	register:    make(chan *player),
	unregister:  make(chan *player),
	projectiles: make(map[*projectile]bool),
	splosions:   make(map[*splosion]bool),
	players:     make(map[*player]bool),
	turn:        0,
}

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	delta = float64(*tick) / 1000

	http.Handle("/ws/", websocket.Handler(addPlayer))

	go g.run()

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("unable to start server")
	}
}

func addPlayer(ws *websocket.Conn) {
	id := fmt.Sprintf("robot%d", <-g.id)
	log.Printf("sending robot id: %s", id)
	err := websocket.JSON.Send(ws, NewHandshake(id, true))
	if err != nil {
		log.Fatal(err)
	}
	var conf config

	for {
		err = websocket.JSON.Receive(ws, &conf)
		if err != nil {
			log.Print(err)
			return
		}
		if conf.Stats.valid() {
			_ = websocket.JSON.Send(ws, NewHandshake(id, true))
			break
		} else {
			_ = websocket.JSON.Send(ws, NewHandshake(id, false))
			log.Printf("%s invalid config", id)
		}
	}
	log.Printf("%s eventually sent valid config", id)

	start_pos := position{
		X: rand.Float64() * 800,
		Y: rand.Float64() * 550,
	}
	p := &player{
		Robot: robot{
			Stats: stats{
				Speed:         *velocity,
				Hp:            200,
				WeaponRadius:  35,
				ScannerRadius: 200,
			},
			Position: start_pos,
			MoveTo:   start_pos,
			Id:       id,
			Health:   200,
			Scanners: make([]scanner, 0)},
		send: make(chan *boardstate),
		ws:   ws,
	}
	g.register <- p
	defer func() {
		g.unregister <- p
	}()
	go p.sender()
	p.recv()
	log.Printf("%v has been disconnect from this game\n", p.Robot.Id)
}
