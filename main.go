package main

import (
	"drbh/partita/websocket"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := fiber.New()

	log.Println("🚚 Initializing services...")
	matchMakingManager := InitializeMatch()
	websocketManager := InitializeWebSocket()

	// initialize multiple background services to run in parallel
	backgroundServiceManager := InitializeBackgroundService()
	backgroundServiceManager2 := InitializeBackgroundService()

	log.Println("🚥 Initializing routes...")
	app.Static("/", "./app/build")

	app.Use("/ws", websocket.UpgradeWebSocket)
	app.Get("/ws/:id", websocketManager.HandleWebSocketConnections)
	app.Get("/matcher", matchMakingManager.Get)

	log.Println("🍔 Starting background processes...")
	backgroundServiceManager.Start()
	backgroundServiceManager2.StartEmitting()

	log.Println("🎧 Starting application...")
	app.Listen(":3000")
}
