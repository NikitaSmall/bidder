package router

import "github.com/gin-gonic/gin"

// New creates, configures and returns ready to work router
func New() *gin.Engine {
	r := gin.Default()

	r.GET("/take", takeHandler)
	r.GET("/fund", fundHandler)
	r.GET("/announceTournament", announceTournamentHandler)
	r.GET("/joinTournament", joinTournamentHandler)
	r.GET("/balance", balanceHandler)
	r.GET("/reset", resetHandler)

	r.POST("/resultTournament", resultTournamentHandler)

	return r
}
