package game

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"

	"encoding/json"
)

type GameService struct {
	Games      map[string]*Game
	GamesMutex sync.Mutex
}

type Game struct {
	State   string
	Players map[string]*Player
}

type Player struct {
	Name         string
	X, Y, Z      float64
	Rotation     float64
	LastRotation float64
	PathPoints   []PathPoint
	JustSpawned  bool
}

type PathPoint struct {
	X, Y, Z float64
}

var gameServiceInstance *GameService
var once sync.Once

const frontFacing = 2 * math.Pi
const limit = 8.0
const lowerLimit = -limit

func ProvideGameService() *GameService {
	log.Println("ProvideGameService")
	return GetGameServiceInstance()
}

func GetGameServiceInstance() *GameService {
	once.Do(func() {
		gameServiceInstance = &GameService{
			Games:      make(map[string]*Game),
			GamesMutex: sync.Mutex{},
		}
		log.Println("🎮 Successfully connected to Game Service")
	})
	return gameServiceInstance
}

func (e *GameService) PlayerFromConnectionID(connectionID string) *Player {

	// random player position
	x := rand.Float64()*(limit-lowerLimit) + lowerLimit
	z := rand.Float64()*(limit-lowerLimit) + lowerLimit

	return &Player{
		Name:         connectionID,
		X:            x,
		Y:            0,
		Z:            z,
		LastRotation: frontFacing,
		Rotation:     frontFacing,
		PathPoints: []PathPoint{
			{
				X: x,
				Y: 0,
				Z: z,
			},
		},
	}
}

func (e *GameService) PrintAllGames() {
	for key, _ := range e.Games {
		fmt.Println(key)
	}
}

func (e *GameService) AddGame(key string, game *Game) {
	e.GamesMutex.Lock()
	defer e.GamesMutex.Unlock()
	e.Games[key] = game
}

func (e *GameService) GetGame(key string) (*Game, bool) {
	e.GamesMutex.Lock()
	defer e.GamesMutex.Unlock()
	game, ok := e.Games[key]
	return game, ok
}

// JoinGame
func (e *GameService) JoinGame(key string, player *Player) {
	e.GamesMutex.Lock()
	defer e.GamesMutex.Unlock()
	game, ok := e.Games[key]
	if !ok {
		log.Println("Game does not exist")
		return
	}
	game.Players[player.Name] = player
}

// LeaveGame
func (e *GameService) LeaveGame(key string, player *Player) {
	e.GamesMutex.Lock()
	defer e.GamesMutex.Unlock()
	game, ok := e.Games[key]
	if !ok {
		log.Println("Game does not exist")
		return
	}
	delete(game.Players, player.Name)

	// if game is empty, remove it
	if len(game.Players) == 0 {
		delete(e.Games, key)
	}
}

// LeaveAllGames
func (e *GameService) LeaveAllGames(player *Player) {
	// use LeaveGame
	keys := make([]string, 0, len(e.Games))
	for k := range e.Games {
		keys = append(keys, k)
	}
	for _, key := range keys {
		e.LeaveGame(key, player)
	}
	log.Printf("❌ Player %v left all games\n", player.Name)
}

func (e *GameService) RotatePlayer(player *Player, rotation float64) error {
	e.GamesMutex.Lock()
	defer e.GamesMutex.Unlock()
	for _, game := range e.Games {
		for _, gamePlayer := range game.Players {
			if gamePlayer.Name == player.Name {
				gamePlayer.Rotation = rotation
				return nil
			}
		}
	}
	return fmt.Errorf("Player not found")
}

func (e *GameService) RemoveGame(key string) {
	e.GamesMutex.Lock()
	defer e.GamesMutex.Unlock()
	delete(e.Games, key)
}

func (e *GameService) UpdateGame(key string, game *Game) {
	e.GamesMutex.Lock()
	defer e.GamesMutex.Unlock()
	e.Games[key] = game
}

func (e *GameService) GetGames() map[string]*Game {
	e.GamesMutex.Lock()
	defer e.GamesMutex.Unlock()
	return e.Games
}

func (e *GameService) GetAllGames() map[string]*Game {
	e.GamesMutex.Lock()
	defer e.GamesMutex.Unlock()
	return e.Games
}

func (e *GameService) UpdateAllGames(games map[string]*Game) {
	e.GamesMutex.Lock()
	defer e.GamesMutex.Unlock()
	e.Games = games
}

func (e *GameService) GetAllGamesJSON() string {
	e.GamesMutex.Lock()
	defer e.GamesMutex.Unlock()
	return toJSON(e.Games)
}

func toJSON(games map[string]*Game) string {
	json, err := json.Marshal(games)
	if err != nil {
		log.Println(err)
	}
	return string(json)
}
