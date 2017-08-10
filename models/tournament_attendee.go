package models

import (
	"bidder/util"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// TournamentAttendee struct holds attendee related data and helps to process it
type TournamentAttendee struct {
	TournamentID int      `form:"tournamentId" json:"tournamentId" binding:"required"`
	PlayerID     string   `form:"playerId" json:"playerId" binding:"required"`
	Backers      []string `form:"backerId"`
}

// Validate function checks the params before execute actual request
func (ta *TournamentAttendee) Validate() error {
	if len(ta.PlayerID) == 0 {
		return errors.New("PlayerID should not be empty!")
	}

	if ta.TournamentID < 0 {
		return errors.New("TournamentID should be positive number!")
	}

	validationPlayers := append(ta.Backers, ta.PlayerID)

	uniqValidator := make(map[string]bool)

	for _, id := range validationPlayers {
		uniqValidator[id] = true
	}

	if len(uniqValidator) < len(validationPlayers) {
		return errors.New("Every player should be uniq!")
	}

	return nil
}

// JoinTournament function tries to join the tournament by provided users
func (ta *TournamentAttendee) JoinTournament() error {
	tx, err := util.DBConnect.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`SELECT deposit, finished FROM tournaments WHERE id = $1 FOR UPDATE;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var deposit int
	var finished bool
	if err := stmt.QueryRow(ta.TournamentID).Scan(&deposit, &finished); err != nil {
		tx.Rollback()
		return err
	}

	stmt, err = tx.Prepare(`SELECT player_id FROM tournament_attendees WHERE player_id = $1 AND tournament_id = $2 FOR UPDATE;`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	var playerID string
	if err := stmt.QueryRow(ta.PlayerID, ta.TournamentID).Scan(&playerID); err != sql.ErrNoRows {
		tx.Rollback()
		return errors.New("Cannot join tournament second time")
	}

	if finished {
		tx.Rollback()
		return errors.New("Cannot join to finished tournament")
	}

	var ids []string
	for _, id := range append(ta.Backers, ta.PlayerID) {
		ids = append(ids, strconv.Quote(id))
	}

	playerIDs := fmt.Sprintf(`{%s}`, strings.Join(ids, `, `))
	stmt, err = tx.Prepare(`SELECT player_id, points FROM players WHERE player_id = ANY($1) FOR UPDATE;`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(playerIDs)

	var players []Player
	for rows.Next() {
		var playerID string
		var points int
		err = rows.Scan(&playerID, &points)
		if err != nil {
			tx.Rollback()
			return err
		}

		players = append(players, Player{PlayerID: playerID, Points: points})
	}

	if len(players) != len(append(ta.Backers, ta.PlayerID)) {
		tx.Rollback()
		return errors.New("Not every player could be retrieved")
	}

	stmt, err = tx.Prepare(`UPDATE players SET points = points - $1 WHERE player_id = ANY($2);`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	priceToPay := deposit / len(players)
	if _, err = stmt.Exec(priceToPay, playerIDs); err != nil {
		tx.Rollback()
		return err
	}

	stmt, err = tx.Prepare(`INSERT INTO tournament_attendees (player_id, tournament_id, backers)
                           VALUES ($1, $2, $3);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var b []string
	for _, backer := range ta.Backers {
		b = append(b, strconv.Quote(backer))
	}

	backers := fmt.Sprintf(`{%s}`, strings.Join(b, `, `))
	if _, err = stmt.Exec(ta.PlayerID, ta.TournamentID, backers); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
