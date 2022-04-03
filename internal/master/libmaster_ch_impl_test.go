package master

import (
	"testing"
)

func TestLibMasterCH_Router(t *testing.T) {

	masterLib := NewLibMasterCH()

	masterLib.AddMember("localhost:8081", []uint64{hashKey("apple"), hashKey("pie")})
	masterLib.AddMember("localhost:8082", []uint64{hashKey("name")})
	masterLib.AddMember("localhost:8083", []uint64{hashKey("k1")})

	if hostport, _ := masterLib.Router("apple"); hostport != "localhost:8081" {
		t.Errorf("Router key = apple, want result = 8081, get = %v", hostport)
	}
	if hostport, _ := masterLib.Router("pie"); hostport != "localhost:8081" {
		t.Errorf("Router key = pie, want result = 8081, get = %v", hostport)
	}
	if hostport, _ := masterLib.Router("name"); hostport != "localhost:8082" {
		t.Errorf("Router key = name, want result = 8082, get = %v", hostport)
	}
	if hostport, _ := masterLib.Router("k1"); hostport != "localhost:8083" {
		t.Errorf("Router key = k1, want result = 8083, get = %v", hostport)
	}

}