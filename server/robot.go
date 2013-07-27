package main

type position struct {
    X int `json:"x"`
    Y int `json:"y"`
}

type weapon struct {
	Strength float64 `json:"Strength"`
	Radius   float64 `json:"radius"`
}

type stats struct { Speed float64 `json:"speed"`
	HP    int     `json:"hp"`
}

type robot struct {
	Stats       stats    `json:"stats"`
	Health      float64  `json:"health"`
	Position    position `json:"position"`
	Destination position `json:"destination"`
}
