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

func (c *ClientID) Valid() bool {
	switch c.Useragent {
	case "robot", "spectator":
		return true
	}
	return false
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
