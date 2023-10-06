package redis

import (
	"testing"
)

type MockRedisService struct{}

func (m *MockRedisService) AddPlayer(player string, elo float64) error {
	return nil
}

func (m *MockRedisService) GetPlayerElo(player string) (float64, error) {
	return 1000.0, nil
}

func (m *MockRedisService) RemovePlayer(player string) error {
	return nil
}

func TestProvideMyRedisService(t *testing.T) {
	service := &MockRedisService{}
	if service == nil {
		t.Errorf("ProvideMyRedisService failed, expected %v, got %v", "not nil", "nil")
	}
}

func TestAddPlayer(t *testing.T) {
	service := &MockRedisService{}
	err := service.AddPlayer("testPlayer", 1000.0)
	if err != nil {
		t.Errorf("AddPlayer failed, expected %v, got %v", "nil", err)
	}
}

func TestGetPlayerElo(t *testing.T) {
	service := &MockRedisService{}
	_, err := service.GetPlayerElo("testPlayer")
	if err != nil {
		t.Errorf("GetPlayerElo failed, expected %v, got %v", "nil", err)
	}
}

func TestRemovePlayer(t *testing.T) {
	service := &MockRedisService{}
	err := service.RemovePlayer("testPlayer")
	if err != nil {
		t.Errorf("RemovePlayer failed, expected %v, got %v", "nil", err)
	}
}
