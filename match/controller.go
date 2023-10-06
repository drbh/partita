package match

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type MatchmakingController struct {
	MatchmakingService MatchmakingService
}

func NewMatchmakingController(matchmakingService MatchmakingService) MatchmakingController {
	return MatchmakingController{
		MatchmakingService: matchmakingService,
	}
}

func (e *MatchmakingController) Get(c *fiber.Ctx) error {
	match := e.MatchmakingService.GetPotentialMatchesByElo(1, 1)
	fmt.Println(match)
	c.SendString("MatchmakingController")
	return nil
}
