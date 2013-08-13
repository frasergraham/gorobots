package protocol

import (
	"encoding/json"
	"testing"
)

type result struct {
	b   bool
	msg string
}
type clientIDTest struct {
	clientid ClientID
	expected result
}

var clientIDTests = []clientIDTest{
	{ClientID{Type: "robot"}, result{true, ""}},
	{ClientID{Type: "spectator"}, result{true, ""}},
	{ClientID{Type: "schmarglenoggler"}, result{false, "usergent must be 'robot' or 'spectator'"}},
}

func TestClientIDs(t *testing.T) {
	for _, tt := range clientIDTests {
		v, msg := tt.clientid.Valid()
		actual := result{v, msg}
		if actual.b != tt.expected.b || actual.msg != tt.expected.msg {
			t.Errorf("%+v: expected %v, actual %v", tt.clientid, tt.expected, actual)
		}
	}
}

func TestClientIDParse(t *testing.T) {
	var s ClientID
    err := json.Unmarshal(
		[]byte(`{
            "type": "robot",
            "name": "dummy",
            "id": "24601",
            "useragent": "gorobots.js"
        }`), &s)
	if err != nil {
        t.Errorf("fail to parse: %v", err)
	}
}
