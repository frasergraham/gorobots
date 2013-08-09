package protocol

// > identify
type IdRequest struct {
	Type string `json:"type"`
}

func NewIdRequest() *IdRequest {
	return &IdRequest{
		Type: "idreq",
	}
}

// < [robot | spectator], name, client-type, game ID
type ClientID struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	RequestedID string `json:"id"`
	Useragent   string `json:"useragent"`
}

func (c *ClientID) Valid() (bool, string) {
	switch c.Useragent {
	case "robot", "spectator":
		return true, ""
	}
	return false, "usergent must be 'robot' or 'spectator'"
}

type BoardSize struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type GameParam struct {
	BoardSize BoardSize `json:"boardsize"`
	Type      string    `json:"type"`
}

// > [OK | FULL | NOT AUTH], board size, game params
func NewGameParam(w, h float64) *GameParam {
	return &GameParam{
		BoardSize: BoardSize{
			Width:  w,
			Height: h,
		},
		Type: "gameparam",
	}
}

type Handshake struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Type    string `json:"type"`
}

func NewHandshake(id string, success bool) *Handshake {
	return &Handshake{
		ID:      id,
		Success: success,
		Type:    "handshake",
	}
}

type Failure struct {
	Reason string `json:"reason"`
	Type   string `json:"type"`
}

func NewFailure(reason string) *Failure {
	return &Failure{
		Reason: reason,
		Type:   "failure",
	}
}
