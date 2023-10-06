package match

import (
	"drbh/partita/redis"
	"log"
	"strconv"
	"time"
)

type MatchmakingService struct {
	RedisService *redis.MyRedisService
}

func NewMatchmakingService(redisService *redis.MyRedisService) MatchmakingService {
	return MatchmakingService{
		RedisService: redisService,
	}
}

func (s *MatchmakingService) GetPlayerElo(playerId string) float64 {
	playersElo, err := s.RedisService.GetPlayerElo(playerId)

	if err != nil {
		log.Println(err)
		return -1.0
	}

	return playersElo
}

func (s *MatchmakingService) AddPlayer(playerId string, playerElo float64) error {
	err := s.RedisService.AddPlayer(playerId, playerElo)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *MatchmakingService) GetPotentialMatchesByElo(playerElo float64, eloThreshold float64) []string {
	matches, err := s.RedisService.Elo(playerElo, eloThreshold)

	if err != nil {
		log.Println(err)
		return nil
	}

	return matches
}

func (s *MatchmakingService) WasRecentlyPlayed(playerId, matchId string) bool {
	timestamp, err := s.RedisService.GetLastPlayedWith(playerId, matchId)

	if err != nil {
		// log.Println(err)
		return false
	}

	// check ig its been 1 minute since the last match
	now := time.Now().Unix()
	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)

	if err != nil {
		log.Println(err)
		return false
	}

	if now-timestampInt > 10 {
		return false
	}

	return true
}

func (s *MatchmakingService) RemovePlayer(playerId string) error {
	err := s.RedisService.RemovePlayer(playerId)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *MatchmakingService) UpdateLastPlayedWith(playerId, matchId string) error {
	err := s.RedisService.UpdateLastPlayedWith(playerId, matchId)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *MatchmakingService) AddToRecentlyPlayed(playerId, matchId string) error {
	err := s.RedisService.AddLastPlayedWith(playerId, matchId)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *MatchmakingService) GetMatch(matchId string) (string, error) {
	match, err := s.RedisService.GetMatch(matchId)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return match, nil
}

func (s *MatchmakingService) IsBlocked(playerId, matchId string) bool {
	isBlocked, err := s.RedisService.IsBlocked(playerId, matchId)
	if err != nil {
		log.Println(err)
		return false
	}

	return isBlocked
}

func (s *MatchmakingService) GetPendingPlayers() []string {

	// check that the connection is still alive
	err := s.RedisService.Rdb.Ping(s.RedisService.Ctx).Err()
	if err != nil {
		log.Println("Redis connection is not alive: ", err)
		return nil
	}

	playerIds, errTwo := s.RedisService.GetPendingPlayers()

	if errTwo != nil {
		log.Println("Error getting pending players: ", errTwo)
		return nil
	}

	return playerIds

}

// contains checks if a slice contains a given string
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// FindMatch finds a suitable match for a player based on Elo and other conditions
// Time complexity: O(n), where n is the number of potential matches
func (s *MatchmakingService) FindMatch(playerId string, playerElo float64, eloThreshold float64) string {
	potentialMatches := s.GetPotentialMatchesByElo(playerElo, eloThreshold)
	pendingPlayers := s.GetPendingPlayers()

	for _, matchedPlayer := range potentialMatches {
		if contains(pendingPlayers, matchedPlayer) && !s.WasRecentlyPlayed(playerId, matchedPlayer) && !s.IsBlocked(playerId, matchedPlayer) {
			return matchedPlayer
		}
	}

	return ""
}

// PublishMatch publishes a match to the Redis channel
func (s *MatchmakingService) PublishMatch(match string) error {
	err := s.RedisService.PublishMatch(match)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// ListenForMatch listens for a match on the Redis channel
func (s *MatchmakingService) ListenForMatch(playerId string, callback func(matchId string)) {
	s.RedisService.ListenForMatches(playerId, callback)
}
