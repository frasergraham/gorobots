package main

import (
	"math"
)

const threshold = 1e-7

type position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

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
