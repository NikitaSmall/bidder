package router

import (
	"bidder/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func fundHandler(c *gin.Context) {
	var player models.Player
	err := c.Bind(&player)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationError": err.Error()})
		return
	}

	if err := player.Fund(); err == nil {
		c.JSON(http.StatusOK, gin.H{"Result": "Player funded succesfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"internalError": err.Error()})
	}
}

func takeHandler(c *gin.Context) {
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
	if err := models.ResetDB(); err == nil {
		c.JSON(http.StatusOK, gin.H{"Result": "DataBase is in clean state now"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"internalError": err.Error()})
	}
}

func resultTournamentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}
