package game

import (
	"testing"
)

func TestProvideGameService(t *testing.T) {
	service := ProvideGameService()
	if service == nil {
		t.Errorf("ProvideGameService failed, expected %v, got %v", "not nil", "nil")
	}
}

func TestAddGame(t *testing.T) {
	service := ProvideGameService()
	game := &Game{
		State:   "Test",
		Players: make(map[string]*Player),
	}
	service.AddGame("test", game)
	if _, ok := service.GetGame("test"); !ok {
		t.Errorf("AddGame failed, expected %v, got %v", "true", "false")
	}
}

func TestRemoveGame(t *testing.T) {
	service := ProvideGameService()
	game := &Game{
		State:   "Test",
		Players: make(map[string]*Player),
	}
	service.AddGame("test", game)
	service.RemoveGame("test")
	if _, ok := service.GetGame("test"); ok {
		t.Errorf("RemoveGame failed, expected %v, got %v", "false", "true")
	}
}

func TestUpdateGame(t *testing.T) {
	service := ProvideGameService()
	game := &Game{
		State:   "Test",
		Players: make(map[string]*Player),
	}
	service.AddGame("test", game)
	game.State = "Updated"
	service.UpdateGame("test", game)
	updatedGame, _ := service.GetGame("test")
	if updatedGame.State != "Updated" {
		t.Errorf("UpdateGame failed, expected %v, got %v", "Updated", updatedGame.State)
	}
}
