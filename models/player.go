package models

import (
	"bidder/util"
	"database/sql"
	"errors"
)

// Player struct holds the player's data and allows to work with it in a handy way
type Player struct {
	PlayerID string `form:"playerId" json:"playerId" binding:"required"`
	Points   int    `form:"points" json:"balance" binding:"required"`
}

// Validate method checks the params before execute actual request
func (p *Player) Validate() error {
	if len(p.PlayerID) == 0 {
		return errors.New("PlayerID should not be empty!")
	}

	if p.Points < 0 {
		return errors.New("Points should be positive number!")
	}

	return nil
}

// Fund method tries to create new player or update the existing one with some points
func (p *Player) Fund() error {
	tx, err := util.DBConnect.Begin()
	if err != nil {
		return err
	}

	err = p.fundPlayer(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Take method removes specified number of points from the player
// As this method has more than one database call, every call is in it's own method
func (p *Player) Take() error {
	tx, err := util.DBConnect.Begin()
	if err != nil {
		return err
	}

	if err = p.checkPoints(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err = p.substractPoints(tx); err != nil {
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

	player, err := findPlayer(tx, playerID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return player, tx.Commit()
}

func (p *Player) fundPlayer(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`INSERT INTO players (player_id, points) VALUES ($1, $2)
                           ON CONFLICT(player_id) DO UPDATE SET points = players.points + EXCLUDED.points;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(p.PlayerID, p.Points); err != nil {
		return err
	}

	return nil
}

func (p *Player) checkPoints(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`SELECT points FROM players WHERE player_id = $1 FOR UPDATE;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var currentPoints int
	if err := stmt.QueryRow(p.PlayerID).Scan(&currentPoints); err != nil {
		return err
	}

	if currentPoints-p.Points < 0 {
		return errors.New("Can't set points number to negative")
	}

	return nil
}

func (p *Player) substractPoints(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`UPDATE players SET points = points - $1 WHERE player_id = $2;`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(p.Points, p.PlayerID)
	if err != nil {
		return err
	}

	return nil
}

func findPlayer(tx *sql.Tx, playerID string) (*Player, error) {
	stmt, err := tx.Prepare(`SELECT player_id, points FROM players WHERE player_id = $1;`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	player := new(Player)
	if err := stmt.QueryRow(playerID).Scan(&player.PlayerID, &player.Points); err != nil {
		return nil, err
	}

	return player, nil
}
