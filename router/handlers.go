package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func takeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}

func fundHandler(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}

func announceTournamentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}

func joinTournamentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}

func balanceHandler(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}

func resetHandler(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}

func resultTournamentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}
