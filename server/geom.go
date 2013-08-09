package main

import (
	v "github.com/frasergraham/govector2d"
)

const threshold = 1e-7

func distance(p1, p2 v.Point2d) float64 {
	return p1.Sub(p2).Mag()
}

func move(d1, d2 v.Point2d, velocity float64, timeDelta float64) v.Point2d {
	v := d2.Sub(d1)
	v_norm := v.Normalize()
	v_scaled := v_norm.Scale(velocity * timeDelta)
	return d1.Add(v_scaled)
}
