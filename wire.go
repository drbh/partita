//go:build wireinject
// +build wireinject

// Package main is the entry point of the application.
package main

// Importing necessary packages.
import (
	"drbh/partita/background"
	"drbh/partita/collision"
	"drbh/partita/connection"
	"drbh/partita/game"
	"drbh/partita/match"
	"drbh/partita/redis"
	"drbh/partita/websocket"

	"github.com/google/wire"
)

// SuperSet is a Wire provider set that includes all the providers needed for the application.
var SuperSet = wire.NewSet(
	connection.ProvideConnectionService,
	redis.ProvideMyRedisService,
	game.ProvideGameService,
	// background.ProvideBackgroundService,
)

// InitializeMatch is a Wire provider function that provides an instance of MatchmakingController.
func InitializeMatch() match.MatchmakingController {
	// Wire will use the providers in the Build call to inject the necessary dependencies.
	wire.Build(match.NewMatchmakingController, match.NewMatchmakingService, redis.GetMyRedisServiceInstance)
	// An empty MatchmakingController is returned. Wire will replace this with the actual instance.
	return match.MatchmakingController{}
}

// InitializeWebSocket is a Wire provider function that provides an instance of WebsocketController.
func InitializeWebSocket() websocket.WebsocketController {
	// Wire will use the providers in the Build call to inject the necessary dependencies.
	wire.Build(
		websocket.NewWebsocketController,
		connection.GetConnectionServiceInstance,
		game.GetGameServiceInstance,
		match.NewMatchmakingService, redis.GetMyRedisServiceInstance,
		collision.GetLineSegmentManagerInstance,
	)
	// An empty WebsocketController is returned. Wire will replace this with the actual instance.
	return websocket.WebsocketController{}
}

// InitalizeConnection is a Wire provider function that provides an instance of ConnectionService.
func InitalizeConnection() *connection.ConnectionService {
	// Wire will use the provider in the Build call to inject the necessary dependencies.
	wire.Build(connection.GetConnectionServiceInstance)
	// An empty ConnectionService is returned. Wire will replace this with the actual instance.
	return &connection.ConnectionService{}
}

// InitializeBackgroundService is a Wire provider function that provides an instance of BackgroundService.
func InitializeBackgroundService() background.BackgroundServiceInterface {
	// Wire will use the provider in the Build call to inject the necessary dependencies.
	wire.Build(
		background.NewBackgroundService,
		connection.GetConnectionServiceInstance,
		game.GetGameServiceInstance,
		// TODO: fix that both required below since NewMatchmakingService is not a pointer
		match.NewMatchmakingService, redis.GetMyRedisServiceInstance,
		collision.GetLineSegmentManagerInstance,
	)
	// An empty BackgroundService is returned. Wire will replace this with the actual instance.
	// return &background.BackgroundService{}
	return &background.BackgroundService{}
}
