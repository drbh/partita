package connection

import (
	"testing"

	"github.com/gofiber/websocket/v2"
)

func TestAddConnection(t *testing.T) {
	service := ProvideConnectionService()
	conn := &websocket.Conn{}
	service.AddConnection("test", conn)
	if _, ok := service.GetConnection("test"); !ok {
		t.Errorf("AddConnection failed, expected %v, got %v", "true", "false")
	}
}

func TestRemoveConnection(t *testing.T) {
	service := ProvideConnectionService()
	conn := &websocket.Conn{}
	service.AddConnection("test", conn)
	service.RemoveConnection("test")
	if _, ok := service.GetConnection("test"); ok {
		t.Errorf("RemoveConnection failed, expected %v, got %v", "false", "true")
	}
}

func TestUpdateConnection(t *testing.T) {
	service := ProvideConnectionService()
	conn1 := &websocket.Conn{}
	conn2 := &websocket.Conn{}
	service.AddConnection("test", conn1)
	service.UpdateConnection("test", conn2)
	conn, _ := service.GetConnection("test")
	if conn != conn2 {
		t.Errorf("UpdateConnection failed, expected %v, got %v", conn2, conn)
	}
}
