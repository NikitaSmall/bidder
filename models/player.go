package models

import (
	"bidder/util"
	"errors"
)

// Player struct holds the player's data and allows to work with it in a handy way
type Player struct {
	PlayerID string `form:"playerId" json:"playerId" binding:"required"`
	Points   int    `form:"points" json:"balance" binding:"required"`
}

// Validate function checks the params before execute actual request
func (p *Player) Validate() error {
	if len(p.PlayerID) == 0 {
		return errors.New("PlayerID should not be empty!")
	}

	if p.Points < 0 {
		return errors.New("Points should be positive number!")
	}

	return nil
}

// Fund function tries to create new player or update the existing one with some points
func (p *Player) Fund() error {
	tx, err := util.DBConnect.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO players (player_id, points) VALUES ($1, $2)
                           ON CONFLICT(player_id) DO UPDATE SET points = players.points + EXCLUDED.points;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(p.PlayerID, p.Points); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Take function removes specified number of points from the player
// if the number of points should be negative, that player gets zero points instead.
func (p *Player) Take() error {
	tx, err := util.DBConnect.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`SELECT points FROM players WHERE player_id = $1 FOR UPDATE;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var currentPoints int
	if err := stmt.QueryRow(p.PlayerID).Scan(&currentPoints); err != nil {
		tx.Rollback()
		return err
	}

	stmt, err = tx.Prepare(`UPDATE players SET points = $1 WHERE player_id = $2;`)
	if err != nil {
		return err
	}

	if currentPoints-p.Points < 0 {
		tx.Rollback()
		return errors.New("Can't set points number to negative")
	}

	_, err = stmt.Exec(currentPoints-p.Points, p.PlayerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// FindPlayer function tries to find and return the player from the DataBase
func FindPlayer(playerID string) (*Player, error) {
	tx, err := util.DBConnect.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(`SELECT * FROM players WHERE player_id = $1;`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var id string
	var points int

	if err := stmt.QueryRow(playerID).Scan(&id, &points); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &Player{PlayerID: id, Points: points}, tx.Commit()
}
