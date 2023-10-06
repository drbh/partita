// Package background provides the background services for the game
package background

// Importing necessary packages
import (
	"drbh/partita/collision"
	"drbh/partita/connection"
	"drbh/partita/game"
	"drbh/partita/match"
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

const limit = 8.0
const lowerLimit = -limit

// BackgroundServiceInterface defines the methods for the background service
type BackgroundServiceInterface interface {
	Start()
	StartEmitting()
	EmitLocations()
	BuildMatches()
}

// BackgroundService struct holds the services needed for the game
type BackgroundService struct {
	connectionService  *connection.ConnectionService
	matchmakingService match.MatchmakingService
	gameService        *game.GameService
	collisionService   *collision.LineSegmentManager
}

// Global instances of the BackgroundServiceInterface and BackgroundService
var BackgroundServiceInterfaceInstance BackgroundServiceInterface
var BackgroundServiceInstance *BackgroundService
var once sync.Once

// NewBackgroundService initializes a new instance of BackgroundService
func NewBackgroundService(
	connectionService *connection.ConnectionService,
	matchmakingService match.MatchmakingService,
	gameService *game.GameService,
	collisionService *collision.LineSegmentManager,
) BackgroundServiceInterface {
	once.Do(func() {
		BackgroundServiceInstance = &BackgroundService{
			connectionService:  connectionService,
			matchmakingService: matchmakingService,
			gameService:        gameService,
			collisionService:   collisionService,
		}
		log.Println("🍬 Successfully connected to Background Service")
	})
	return BackgroundServiceInstance
}

// Start method starts the background service
func (e *BackgroundService) Start() {
	go e.BuildMatches()
	log.Println("🍟 Successfully started Background Service")
}

// StartEmitting method starts emitting locations
func (e *BackgroundService) StartEmitting() {
	go e.EmitLocations()
	log.Println("🍟 Successfully started Background Service")
}

// Constants for the game
const delta = 0.25
const speed = 0.50
const boundary = 8
const frontFacing = 2 * math.Pi
const leftFacing = -math.Pi / 2
const backFacing = frontFacing + math.Pi
const rightFacing = frontFacing + math.Pi/2

// helper function to create segments from path points
func createSegments(
	player *game.Player,
	playerNextX float64,
	playerNextZ float64,
) []collision.Segment {
	var segments []collision.Segment
	for i := 0; i < len(player.PathPoints)-1; i++ {
		newSegment := collision.NewSegmentFromCoords(
			player.PathPoints[i].X, player.PathPoints[i].Z,
			player.PathPoints[i+1].X, player.PathPoints[i+1].Z,
		)
		segments = append(segments, newSegment)
	}
	segments = append(segments, collision.NewSegmentFromCoords(
		player.PathPoints[len(player.PathPoints)-1].X, player.PathPoints[len(player.PathPoints)-1].Z,
		playerNextX, playerNextZ,
	))
	return segments
}

// EmitLocations method emits the locations of the players
func (e *BackgroundService) EmitLocations() {
	ticker := time.NewTicker(30 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		allGames := e.gameService.GetAllGames()
		for _, currentGame := range allGames {
			for _, player := range currentGame.Players {
				e.processPlayerMovement(player, currentGame)
			}
		}
		e.updateGameStateAndNotifyClients(allGames)
	}
}

// processPlayerMovement processes the movement of a single player
func (e *BackgroundService) processPlayerMovement(player *game.Player, currentGame *game.Game) {
	originalX, originalY, originalZ := player.X, player.Y, player.Z

	// move the player
	nextX, nextZ := e.calculateNextPosition(player)

	// start timer
	start := time.Now()

	if len(player.PathPoints) > 1 {
		playersWhoInitatedCollision := e.checkPlayerCollision(player, nextX, nextZ, currentGame)

		// print that player.Name has collided with other players
		for collider := range playersWhoInitatedCollision {
			log.Printf("%v has collided with %v\n", player.Name, collider)

			var payload map[string]interface{} = map[string]interface{}{
				"command": "playerCollision",
				"name":    collider,
				"with":    player.Name,
				"time":    time.Now().UnixNano() / int64(time.Millisecond),
			}
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				log.Printf("Error marshalling playerCollision payload: %v\n", err)
				return
			}
			e.connectionService.SendToAll(string(payloadBytes))
		}
	}

	// end timer
	elapsed := time.Since(start)

	if false {
		log.Printf("Collision detection took %v\n", elapsed)
	}

	// if player has just spawned, don't add a new path point
	if player.JustSpawned {
		player.JustSpawned = false
		log.Printf("🐘🐘🐘 Player has a total of %v path points\n", len(player.PathPoints))
		return
	}

	// move the player
	player.X = nextX
	player.Z = nextZ

	// check if player hits the boundary and reverse its direction
	e.checkBoundaryCollision(player)

	// only add a new path point if the player has turned
	if player.Rotation != player.LastRotation {
		player.PathPoints = append(player.PathPoints, game.PathPoint{
			X: originalX,
			Y: originalY,
			Z: originalZ,
		})
		log.Printf("Player has a total of %v path points\n", len(player.PathPoints))
	}

	player.LastRotation = player.Rotation
}

// calculateNextPosition calculates the next position of a player
func (e *BackgroundService) calculateNextPosition(player *game.Player) (float64, float64) {
	nextX := math.Round((player.X+math.Sin(player.Rotation)*speed*delta)*10000) / 10000
	nextZ := math.Round((player.Z+math.Cos(player.Rotation)*speed*delta)*10000) / 10000
	return nextX, nextZ
}

// checkPlayerCollision checks for collision between players
func (e *BackgroundService) checkPlayerCollision(player *game.Player, nextX float64, nextZ float64, currentGame *game.Game) map[string]bool {
	playerSegments := createSegments(player, player.X, player.Y) //nextX, nextZ)

	playersToReset := make(map[string]bool)

	for _, otherPlayer := range currentGame.Players {
		if player.Name == otherPlayer.Name || len(player.PathPoints) < 2 || len(otherPlayer.PathPoints) < 2 {
			continue
		}
		otherPlayerNextX, otherPlayerNextZ := e.calculateNextPosition(otherPlayer)

		otherPlayerSegments := createSegments(otherPlayer, otherPlayerNextX, otherPlayerNextZ)
		e.collisionService.ClearAllSegments()
		for _, segment := range playerSegments {
			e.collisionService.AddSegment(segment)
		}
		intersectingSegment := e.collisionService.CheckIntersection(
			otherPlayerSegments[len(otherPlayerSegments)-1],
		)
		if intersectingSegment != nil {
			playersToReset[otherPlayer.Name] = true
			e.resetPlayerPosition(player)
		}
	}

	return playersToReset
}

// resetPlayerPosition resets the position of a player after a collision
func (e *BackgroundService) resetPlayerPosition(player *game.Player) {

	// random player position
	x := rand.Float64()*(limit-lowerLimit) + lowerLimit
	z := rand.Float64()*(limit-lowerLimit) + lowerLimit

	// reset player's path points
	player.PathPoints = []game.PathPoint{
		{X: x, Y: 0.0, Z: z},
		{X: x, Y: 0.0, Z: z},
	}
	player.X = x
	player.Y = 0
	player.Z = z
	player.Rotation = frontFacing
	player.LastRotation = frontFacing
	player.JustSpawned = true
}

// checkBoundaryCollision checks if a player hits the boundary and reverses its direction
func (e *BackgroundService) checkBoundaryCollision(player *game.Player) {
	if math.Abs(player.X) > boundary {
		if player.X > 0 {
			player.Rotation = leftFacing
		} else {
			player.Rotation = rightFacing
		}
		if player.X > 0 {
			player.X = boundary
		} else {
			player.X = -boundary
		}
	}
	if math.Abs(player.Z) > boundary {
		if player.Z > 0 {
			player.Rotation = backFacing
		} else {
			player.Rotation = frontFacing
		}
		if player.Z > 0 {
			player.Z = boundary
		} else {
			player.Z = -boundary
		}
	}
}

// updateGameStateAndNotifyClients updates the game state and notifies all clients
func (e *BackgroundService) updateGameStateAndNotifyClients(allGames map[string]*game.Game) {
	// update the game state
	e.gameService.UpdateAllGames(allGames)

	jsonVersion := e.gameService.GetAllGamesJSON()

	if len(jsonVersion) > 4 {
		e.connectionService.SendToAll(jsonVersion)
	}
}

// BuildMatches method builds matches for the game
func (e *BackgroundService) BuildMatches() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pendingPlayers := e.matchmakingService.GetPendingPlayers()
		var matches []string
		for _, playerId := range pendingPlayers {
			playerElo := e.matchmakingService.GetPlayerElo(playerId)
			match := e.matchmakingService.FindMatch(playerId, playerElo, 100)
			if match != "" && playerId != match {
				matches = append(matches, playerId, match)

				// remove players from pending
				e.matchmakingService.RemovePlayer(playerId)
				e.matchmakingService.RemovePlayer(match)

				// update last played with
				e.matchmakingService.UpdateLastPlayedWith(playerId, match)
				e.matchmakingService.UpdateLastPlayedWith(match, playerId)

			}
		}
		if len(matches) > 0 {
			log.Printf("Matches built: %v\n", matches)
			marshaledMatches, err := json.Marshal(matches)
			if err != nil {
				log.Println(err)
				return
			}
			e.matchmakingService.PublishMatch(string(marshaledMatches))
		}
	}
}
