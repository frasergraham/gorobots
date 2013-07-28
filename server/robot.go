package main

type position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type weapon struct {
	Strength float64 `json:"strength"`
	Radius   float64 `json:"radius"`
}

type stats struct {
	Speed         float64 `json:"speed"`
	Hp            int     `json:"hp"`
	WeaponRadius  int     `json:"weapon_radius"`
	ScannerRadius int     `json:"scanner_radius"`
}

type scanner struct {
	Position position `json:"position"`
	Stats    stats    `json:"stats"`
}

type robot struct {
	Id       string    `json:"id"`
	Stats    stats     `json:"stats"`
	Health   int       `json:"health"`
	Position position  `json:"position"`
	MoveTo   position  `json:"move_to"`
	FireAt   position  `json:"fire_at"`
	Scanners []scanner `json:"scanners"`
}

type projectile struct {
	Id       string   `json:"id"`
	Position position `json:"position"`
	MoveTo   position `json:"move_to"`
	Radius   int      `json:"radius"`
	Damage   int      `json:"damage"`
}

type splosion struct {
	Id        string   `json:"id"`
	Position  position `json:"position"`
	Radius    int      `json:"radius"`
	MaxDamage int      `json:"damage"`
	MinDamage int      `json:"damage"`
	Lifespan  int      `json:"lifespan"`
}
