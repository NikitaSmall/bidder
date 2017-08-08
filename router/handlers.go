package router

import (
	"bidder/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func fundHandler(c *gin.Context) {
	var player models.Player

	if err := c.Bind(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationError": err.Error()})
		return
	}

	if err := player.Validate(); err != nil {
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
	var player models.Player

	if err := c.Bind(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationError": err.Error()})
		return
	}

	if err := player.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationError": err.Error()})
		return
	}

	if err := player.Take(); err == nil {
		c.JSON(http.StatusOK, gin.H{"Result": "Player's points were taken succesfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"internalError": err.Error()})
	}
}

func announceTournamentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}

func resultTournamentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}

func joinTournamentHandler(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}

func balanceHandler(c *gin.Context) {
	playerID := c.Query("playerId")

	if player, err := models.FindPlayer(playerID); err == nil {
		c.JSON(http.StatusOK, player)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"notFoundError": "No such player"})
	}
}

func resetHandler(c *gin.Context) {
	if err := models.ResetDB(); err == nil {
		c.JSON(http.StatusOK, gin.H{"Result": "DataBase is in clean state now"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"internalError": err.Error()})
	}
}
