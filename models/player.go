package models

import "bidder/util"

// Player struct holds the player's data and allows to work with it in a handy way
type Player struct {
	PlayerID string `form:"playerId" json:"playerId" binding:"required"`
	Points   int    `form:"points" json:"points" binding:"required"`
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
