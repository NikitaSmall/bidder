package models

import (
	"bidder/util"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Tournament struct holds tournament related data and helps to process it
type Tournament struct {
	TournamentID int  `form:"tournamentId" json:"tournamentId" binding:"required"`
	Deposit      int  `form:"deposit" json:"deposit" binding:"required"`
	Finished     bool `json:"finished"`
}

// TournamentResult struct holds the data required to process result finish
type TournamentResult struct {
	TournamentID string   `form:"tournamentId" json:"tournamentId" binding:"required"`
	Winners      []Winner `form:"winners" json:"winners" binding:"required"`
}

// Winner struct holds winner related data. Helper struct to work with TournamentResult struct
type Winner struct {
	PlayerID string `form:"playerId" json:"playerId" binding:"required"`
	Prize    int    `form:"prize" json:"prize" binding:"required"`
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

	stmt, err := tx.Prepare(`INSERT INTO tournaments (id, deposit, finished) VALUES ($1, $2, $3);`)
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

// Validate function checks the params before execute actual request
func (tr *TournamentResult) Validate() error {
	if len(tr.TournamentID) == 0 {
		return errors.New("TournamentID should not empty!")
	}

	if len(tr.Winners) == 0 {
		return errors.New("Winners array should not be empty!")
	}

	return nil
}

// Finish method tries to finish the tournament and pay the prize for every player
func (tr *TournamentResult) Finish() error {
	tx, err := util.DBConnect.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`SELECT finished FROM tournaments WHERE id = $1 FOR UPDATE;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var finished bool
	if err := stmt.QueryRow(tr.TournamentID).Scan(&finished); err != nil {
		tx.Rollback()
		return err
	}

	if finished {
		tx.Rollback()
		return errors.New("Cannot finish finished tournament")
	}

	for _, winner := range tr.Winners {
		var backers []byte
		var playerID string
		var ids []string

		err = tx.QueryRow(`SELECT p.player_id, ta.backers FROM players AS p
											 JOIN tournament_attendees AS ta ON p.player_id = ta.player_id
											 WHERE p.player_id = $1 AND ta.tournament_id = $2 FOR UPDATE;`,
			winner.PlayerID, tr.TournamentID).Scan(&playerID, &backers)

		if err != nil {
			tx.Rollback()
			return err
		}

		b := string(backers)
		if b == "{}" {
			ids = []string{playerID}
		} else {
			ids = append(strings.Split(b[1:len(b)-1], ","), playerID)
		}

		stmt, err = tx.Prepare(`UPDATE players SET points = points + $1 WHERE player_id = ANY($2);`)
		if err != nil {
			tx.Rollback()
			return err
		}

		var requestIDs []string
		for _, id := range ids {
			requestIDs = append(requestIDs, strconv.Quote(id))
		}

		prize := winner.Prize / len(requestIDs)
		playerIDs := fmt.Sprintf(`{%s}`, strings.Join(requestIDs, `, `))

		if _, err = stmt.Exec(prize, playerIDs); err != nil {
			tx.Rollback()
			return err
		}
	}

	stmt, err = tx.Prepare(`UPDATE tournaments SET finished = true WHERE id = $1;`)
	if err != nil {
		tx.Rollback()
		return err
	}

	if _, err = stmt.Exec(tr.TournamentID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}