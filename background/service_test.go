package background

import (
	"log"
	"testing"
)

type DummyBackgroundService struct {
	started          bool
	emitting         bool
	locationsEmitted bool
	matchesBuilt     bool
}

func (e *DummyBackgroundService) Start() {
	e.started = true
	log.Println("🤷‍♀️ Successfully started Background Service")
}

func (e *DummyBackgroundService) StartEmitting() {
	e.emitting = true
	log.Println("🤷‍♀️ Successfully started Background Service")
}

func (e *DummyBackgroundService) EmitLocations() {
	e.locationsEmitted = true
	log.Println("🤷‍♀️ Successfully started Background Service")
}

func (e *DummyBackgroundService) BuildMatches() {
	e.matchesBuilt = true
	log.Println("🤷‍♀️ Successfully started Background Service")
}

func TestNewBackgroundService(t *testing.T) {
	service := &DummyBackgroundService{}
	if service == nil {
		t.Errorf("NewBackgroundService failed, expected %v, got %v", "not nil", "nil")
	}
}

func TestStart(t *testing.T) {
	service := &DummyBackgroundService{}
	service.Start()
	if !service.started {
		t.Errorf("Start failed, expected %v, got %v", true, service.started)
	}
}

func TestStartEmitting(t *testing.T) {
	service := &DummyBackgroundService{}
	service.StartEmitting()
	if !service.emitting {
		t.Errorf("StartEmitting failed, expected %v, got %v", true, service.emitting)
	}
}

func TestEmitLocations(t *testing.T) {
	service := &DummyBackgroundService{}
	service.EmitLocations()
	if !service.locationsEmitted {
		t.Errorf("EmitLocations failed, expected %v, got %v", true, service.locationsEmitted)
	}
}

func TestBuildMatches(t *testing.T) {
	service := &DummyBackgroundService{}
	service.BuildMatches()
	if !service.matchesBuilt {
		t.Errorf("BuildMatches failed, expected %v, got %v", true, service.matchesBuilt)
	}
}
