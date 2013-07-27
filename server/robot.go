package main

type position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type weapon struct {
	Strength float64 `json:"strength"`
	Radius   float64 `json:"radius"`
}

type stats struct {
	Speed float64 `json:"speed"`
	Hp    int     `json:"hp"`
}

type robot struct {
	Id       string   `json:"id"`
	Stats    stats    `json:"stats"`
	Health   int      `json:"health"`
	Position position `json:"position"`
	MoveTo   position `json:"move_to"`
	FireAt   position `json:"fire_at"`
}

type projectile struct {
	Id       string   `json:"id"`
	Position position `json:"position"`
	MoveTo   position `json:"move_to"`
	Damage   int      `json:"damage"`
}
