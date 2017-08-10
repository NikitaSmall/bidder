package router

import (
	"bidder/models"
	"database/sql"
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
		c.JSON(http.StatusBadRequest, gin.H{"badRequest": err.Error()})
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
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"notFoundError": "No such player"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"badRequest": err.Error()})
		}
	}
}

func announceTournamentHandler(c *gin.Context) {
	var tournament models.Tournament

	if err := c.Bind(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationError": err.Error()})
		return
	}

	if err := tournament.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationError": err.Error()})
		return
	}

	if err := tournament.Announce(); err == nil {
		c.JSON(http.StatusOK, gin.H{"Result": "Tournament announced succesfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"badRequest": err.Error()})
	}
}

func joinTournamentHandler(c *gin.Context) {
	var attendee models.TournamentAttendee

	if err := c.Bind(&attendee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationError": err.Error()})
		return
	}

	if err := attendee.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationError": err.Error()})
		return
	}

	if err := attendee.JoinTournament(); err == nil {
		c.JSON(http.StatusOK, gin.H{"Result": "Attendee joined succesfully"})
	} else {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"notFoundError": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"badRequest": err.Error()})
		}
	}
}

func resultTournamentHandler(c *gin.Context) {
	var result models.TournamentResult

	if err := c.BindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationError": err.Error()})
		return
	}

	if err := result.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validationError": err.Error()})
		return
	}

	if err := result.Finish(); err == nil {
		c.JSON(http.StatusOK, gin.H{"Result": "Tournament finished succesfully"})
	} else {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"notFoundError": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"badRequest": err.Error()})
		}
	}
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
