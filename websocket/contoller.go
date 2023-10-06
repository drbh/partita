package websocket

import (
	"drbh/partita/collision"
	"drbh/partita/connection"
	"drbh/partita/game"
	"drbh/partita/match"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WebsocketController struct {
	connectionService  *connection.ConnectionService
	matchmakingService match.MatchmakingService
	gameService        *game.GameService
	collisionService   *collision.LineSegmentManager
}

func NewWebsocketController(
	connectionService *connection.ConnectionService,
	matchmakingService match.MatchmakingService,
	gameService *game.GameService,
	collisionService *collision.LineSegmentManager,
) WebsocketController {
	return WebsocketController{
		connectionService:  connectionService,
		matchmakingService: matchmakingService,
		gameService:        gameService,
		collisionService:   collisionService,
	}
}

// static method
func UpgradeWebSocket(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (e *WebsocketController) Get(c *fiber.Ctx) error {
	c.SendString("WebsocketController")
	return nil
}

func (e *WebsocketController) HandleWebSocketConnections(c *fiber.Ctx) error {
	handler := func(c *websocket.Conn) {
		// get a unique connection ID from the websocket connection
		// we use the remote address and port to generate a unique ID
		connectionID := fmt.Sprintf("%s-%s", c.RemoteAddr(), c.LocalAddr().String())
		log.Printf("New connection: %v\n", connectionID)
		e.connectionService.AddConnection(connectionID, c)
		defer e.connectionService.RemoveConnection(connectionID)

		player := e.gameService.PlayerFromConnectionID(connectionID)
		defer e.gameService.LeaveAllGames(player)

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				c.Close()
				break
			}

			// print the message to the console
			log.Printf("Received: %s", string(msg))

			if string(msg) == "ping" {
				if err := c.WriteMessage(websocket.TextMessage, []byte("pong")); err != nil {
					log.Printf("Error writing pong: %v", err)
				}
			}

			msgParts := strings.Split(string(msg), ":")
			if len(msgParts) < 1 {
				log.Println("Invalid message format")
				// return
				c.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
				continue
			}

			switch msgParts[0] {
			case "addPlayer":
				if len(msgParts) < 3 {
					log.Println("Invalid message format for addPlayer")
					return
				}
				playerElo, err := strconv.ParseFloat(strings.TrimSpace(msgParts[2]), 64)
				if err != nil {
					log.Println("Invalid playerElo: \"", msgParts[2], "\"")
					return
				}
				log.Printf("Adding player: %v with Elo: %v\n", msgParts[1], playerElo)
				e.matchmakingService.AddPlayer(msgParts[1], playerElo)

				e.matchmakingService.ListenForMatch(msgParts[1], func(matchId string) {
					log.Printf("Match found for: %v\n", msgParts[1])
					c.WriteMessage(websocket.TextMessage, []byte("matchFound:"+matchId))
				})

			// joinGame, gameKey
			case "joinGame":
				if len(msgParts) < 2 {
					log.Println("Invalid message format for joinGame")
					return
				}
				gameKey := msgParts[1]
				log.Printf("Joining game: %v\n", gameKey)
				e.gameService.JoinGame(gameKey, player)

			// leaveGame, gameKey
			case "leaveGame":
				if len(msgParts) < 2 {
					log.Println("Invalid message format for leaveGame")
					return
				}
				gameKey := msgParts[1]
				log.Printf("Leaving game: %v\n", gameKey)
				e.gameService.LeaveGame(gameKey, player)

			// setPlayerName, name
			case "setPlayerName":
				if len(msgParts) < 2 {
					log.Println("Invalid message format for setPlayerName")
					return
				}
				name := msgParts[1]
				player.Name = name
				log.Printf("Setting player name: %v\n", name)
				var response = map[string]interface{}{
					"command": "playerNameSet",
					"name":    name,
				}
				payload, err := json.Marshal(response)
				if err != nil {
					log.Printf("Error marshalling playerNameSet payload: %v\n", err)
					return
				}

				c.WriteMessage(websocket.TextMessage, payload)

			// rotate, direction
			case "rotate":
				if len(msgParts) < 2 {
					log.Println("Invalid message format for rotate")
					return
				}
				rotation := msgParts[1]

				log.Printf("Rotating: %v\n", rotation)

				// convert string rotation to float64
				rotationFloat, err := strconv.ParseFloat(rotation, 64)
				if err != nil {
					log.Printf("Invalid rotation value: %v\n", rotation)
					c.WriteMessage(websocket.TextMessage, []byte("Invalid rotation value"))
					continue
				}

				// Apply rotation to the player's game object
				err = e.gameService.RotatePlayer(player, rotationFloat)
				if err != nil {
					log.Printf("Error rotating player: %v\n", err)
					c.WriteMessage(websocket.TextMessage, []byte("Error rotating player"))
					continue
				}

			case "findGame":
				log.Printf("Finding game for: %v\n", player.Name)

				// place player in matchmaking queue
				e.matchmakingService.AddPlayer(player.Name, 100)

				// listen for match
				e.matchmakingService.ListenForMatch(player.Name, func(matchId string) {
					log.Printf("Match found for: %v\n", player.Name)

					var matchList []string
					json.Unmarshal([]byte(matchId), &matchList)
					gameKeyForMatch := fmt.Sprintf("%v_%v", matchList[0], matchList[1])
					payload, err := json.Marshal(
						map[string]interface{}{
							"command":   "matchFound",
							"matchList": matchList,
							"gameKey":   gameKeyForMatch,
						},
					)
					if err != nil {
						log.Printf("Error marshalling matchFound payload: %v\n", err)
						return
					}

					newGame := &game.Game{
						State:   "new",
						Players: make(map[string]*game.Player),
					}

					e.gameService.AddGame(gameKeyForMatch, newGame)

					c.WriteMessage(websocket.TextMessage, payload)

				})

				log.Printf("Matchmaking queue: %v\n", e.matchmakingService.GetPendingPlayers())

			case "startGame":
				gameKey := msgParts[1]

				log.Printf("Starting game: %v\n", gameKey)

				newGame := &game.Game{
					State:   "new",
					Players: make(map[string]*game.Player),
				}

				e.gameService.AddGame(gameKey, newGame)
				e.gameService.JoinGame(gameKey, player)

			default:
				log.Println("Unknown command:", msgParts[0])
			}

		}
	}

	return websocket.New(handler)(c)
}
