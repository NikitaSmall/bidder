package models

import (
	"bidder/util"
	"database/sql"
	"errors"
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

// JoinTournament method tries to join the tournament by provided users
// As this method has more than one database call, each call is in it's own method.
func (ta *TournamentAttendee) JoinTournament() error {
	tx, err := util.DBConnect.Begin()
	if err != nil {
		return err
	}

	deposit, err := ta.getTournamentDeposit(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = ta.checkUniqAttendee(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = ta.updateAttendeeProfiles(tx, deposit)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = ta.addAttendee(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ta *TournamentAttendee) getTournamentDeposit(tx *sql.Tx) (int, error) {
	stmt, err := tx.Prepare(`SELECT deposit, finished FROM tournaments WHERE id = $1 FOR UPDATE;`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var deposit int
	var finished bool
	if err := stmt.QueryRow(ta.TournamentID).Scan(&deposit, &finished); err != nil {
		return 0, err
	}

	if finished {
		return 0, errors.New("Cannot join to finished tournament")
	}

	return deposit, nil
}

func (ta *TournamentAttendee) checkUniqAttendee(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`SELECT player_id FROM tournament_attendees WHERE player_id = $1 AND tournament_id = $2 FOR UPDATE;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var playerID string
	if err := stmt.QueryRow(ta.PlayerID, ta.TournamentID).Scan(&playerID); err != sql.ErrNoRows {
		return errors.New("Cannot join tournament second time")
	}

	return nil
}

func (ta *TournamentAttendee) updateAttendeeProfiles(tx *sql.Tx, deposit int) error {
	playerIDs := preparePostgresArray(append(ta.Backers, ta.PlayerID))
	stmt, err := tx.Prepare(`SELECT player_id, points FROM players WHERE player_id = ANY($1) FOR UPDATE;`)
	if err != nil {
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
			return err
		}

		players = append(players, Player{PlayerID: playerID, Points: points})
	}

	if len(players) != len(append(ta.Backers, ta.PlayerID)) {
		return errors.New("Not every player could be retrieved")
	}

	stmt, err = tx.Prepare(`UPDATE players SET points = points - $1 WHERE player_id = ANY($2);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	priceToPay := deposit / len(players)
	if _, err = stmt.Exec(priceToPay, playerIDs); err != nil {
		return err
	}

	return nil
}

func (ta *TournamentAttendee) addAttendee(tx *sql.Tx) error {
	stmt, err := tx.Prepare(`INSERT INTO tournament_attendees (player_id, tournament_id, backers)
                           VALUES ($1, $2, $3);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	backers := preparePostgresArray(ta.Backers)
	if _, err = stmt.Exec(ta.PlayerID, ta.TournamentID, backers); err != nil {
		return err
	}

	return nil
}
