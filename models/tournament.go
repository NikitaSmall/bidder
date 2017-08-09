package models

import (
	"bidder/util"
	"errors"
)

// Tournament struct holds tournament related data and helps to process it
type Tournament struct {
	TournamentID int  `form:"tournamentId" json:"tournamentId" binding:"required"`
	Deposit      int  `form:"deposit" json:"deposit" binding:"required"`
	Finished     bool `json:"finished"`
}

// Validate function checks the params before execute actual request
func (t *Tournament) Validate() error {
	if t.TournamentID <= 0 {
		return errors.New("TournamentID should be positive number!")
	}

	if t.Deposit < 0 {
		return errors.New("Deposit should be positive number!")
	}

	return nil
}

// Announce function tries to create new tournament in the DataBase
func (t *Tournament) Announce() error {
	tx, err := util.DBConnect.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO tournaments (id, deposit, finished) VALUES ($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(t.TournamentID, t.Deposit, t.Finished); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
