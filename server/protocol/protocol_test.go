package protocol

import "testing"

type result struct {
	b   bool
	msg string
}
type clientIDTest struct {
	clientid ClientID
	expected result
}

var clientIDTests = []clientIDTest{
	{ClientID{Useragent: "robot"}, result{true, ""}},
	{ClientID{Useragent: "spectator"}, result{true, ""}},
	{ClientID{Useragent: "schmarglenoggler"}, result{false, "usergent must be 'robot' or 'spectator'"}},
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
