package protocol

import "testing"

type clientIDTest struct {
	clientid ClientID
	expected bool
}

var clientIDTests = []clientIDTest{
    {ClientID{Useragent: "robot"}, true},
    {ClientID{Useragent: "spectator"}, true},
    {ClientID{Useragent: "schmarglenoggler"}, false},
}

func TestClientIDs(t *testing.T) {
	for _, tt := range clientIDTests {
		actual := tt.clientid.Valid()
		if actual != tt.expected {
			t.Errorf("%+v: expected %t, actual %t", tt.clientid, tt.expected, actual)
		}
	}
}
