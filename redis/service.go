package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// MyRedisService is the main struct for Redis operations
// It contains a Redis client and a context for operations
type MyRedisService struct {
	Rdb *redis.Client
	Ctx context.Context
}

// Singleton instance of MyRedisService
var myRedisServiceInstance *MyRedisService
var once sync.Once

// ProvideMyRedisService returns the singleton instance of MyRedisService
func ProvideMyRedisService() *MyRedisService {
	return GetMyRedisServiceInstance()
}

// GetMyRedisServiceInstance returns the singleton instance of MyRedisService
// It initializes the instance if it hasn't been already
func GetMyRedisServiceInstance() *MyRedisService {
	once.Do(func() {
		myRedisServiceInstance = &MyRedisService{
			Rdb: redis.NewClient(&redis.Options{
				Addr: "localhost:6379",
				DB:   1,
			}),
			Ctx: context.Background(),
		}
		log.Println("🔥 Successfully connected to Redis Service")
	})
	return myRedisServiceInstance
}

// Elo returns a list of players within a certain Elo range
func (s *MyRedisService) Elo(playerElo float64, eloThreshold float64) ([]string, error) {

	lowerBound := playerElo - eloThreshold
	upperBound := playerElo + eloThreshold

	return s.Rdb.ZRangeByScore(s.Ctx, "elo_scores", &redis.ZRangeBy{
		Min: fmt.Sprintf("%f", lowerBound),
		Max: fmt.Sprintf("%f", upperBound),
	}).Result()
}

// GetLastPlayedWith returns the last player a given player played with
func (s *MyRedisService) GetLastPlayedWith(playerId, matchId string) (string, error) {
	return s.Rdb.HGet(s.Ctx, "last_played_with:"+playerId, matchId).Result()
}

// AddLastPlayedWith adds a player to the list of players a given player has played with
func (s *MyRedisService) AddLastPlayedWith(playerId, matchId string) error {
	return s.Rdb.HSet(s.Ctx, "last_played_with:"+playerId, matchId, time.Now().Unix()).Err()
}

// GetMatch returns the match data for a given match ID
func (s *MyRedisService) GetMatch(matchId string) (string, error) {
	return s.Rdb.Get(s.Ctx, "match:"+matchId).Result()
}

// GetPendingPlayers returns a list of players who are waiting for a match
func (s *MyRedisService) GetPendingPlayers() ([]string, error) {
	return s.Rdb.LRange(s.Ctx, "pending_players", 0, -1).Result()
}

// GetPlayerElo returns the Elo score of a given player
func (s *MyRedisService) GetPlayerElo(playerId string) (float64, error) {
	playersElo, err := s.Rdb.ZScore(s.Ctx, "elo_scores", playerId).Result()
	if err != nil {
		log.Println(err)
		return -1.0, err
	}

	return playersElo, nil
}

// IsBlocked checks if a player is blocked from a given match
func (s *MyRedisService) IsBlocked(playerId, matchId string) (bool, error) {
	isBlocked, err := s.Rdb.SIsMember(s.Ctx, "blocked:"+playerId, matchId).Result()
	if err != nil {
		log.Println(err)
		return false, err
	}

	return isBlocked, nil
}

// AddPlayer adds a player to the Elo scores and pending players lists
func (s *MyRedisService) AddPlayer(playerId string, playerElo float64) error {
	err := s.Rdb.ZAdd(s.Ctx, "elo_scores", &redis.Z{Score: playerElo, Member: playerId}).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.Rdb.LPush(s.Ctx, "pending_players", playerId).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// RemovePlayer removes a player from the Elo scores and pending players lists
func (s *MyRedisService) RemovePlayer(playerId string) error {
	err := s.Rdb.ZRem(s.Ctx, "elo_scores", playerId).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.Rdb.LRem(s.Ctx, "pending_players", 0, playerId).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// UpdateLastPlayedWith updates the last played with list for a given player
func (s *MyRedisService) UpdateLastPlayedWith(playerId, matchId string) error {
	err := s.Rdb.HSet(s.Ctx, "last_played_with:"+playerId, matchId, time.Now().Unix()).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// PublishMatch publishes a match to the matches channel
func (s *MyRedisService) PublishMatch(match string) error {
	err := s.Rdb.Publish(s.Ctx, "matches", match).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// ListenForMatches listens for matches and call the callback function when a match is received
func (s *MyRedisService) ListenForMatches(playerId string, callback func(string)) {
	pubsub := s.Rdb.Subscribe(s.Ctx, "matches")
	ch := pubsub.Channel()
	for msg := range ch {

		log.Printf("\nReceived: %s\n", string(msg.Payload))

		var matchDataSlice []string
		json.Unmarshal([]byte(msg.Payload), &matchDataSlice)

		log.Printf("\nMatch data: %v\n", matchDataSlice)

		// check if player is in the match
		for _, player := range matchDataSlice {
			if player == playerId {
				callback(msg.Payload)
				log.Printf("Match found for: %v\n", playerId)
				return
			}
		}

	}
}
