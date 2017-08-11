package models

import (
	"bidder/util"
	"fmt"
	"strconv"
	"strings"
)

var resetQueries = []string{
	"DELETE FROM tournament_attendees;",
	"DELETE FROM tournaments;",
	"DELETE FROM players;",
}

// ResetDB function removes all the data from the DataBase, leaving structure.
func ResetDB() error {
	tx, err := util.DBConnect.Begin()
	if err != nil {
		return err
	}

	for _, query := range resetQueries {
		_, err := tx.Exec(query)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func preparePostgresArray(array []string) string {
	var result []string
	for _, item := range array {
		result = append(result, strconv.Quote(item))
	}

	return fmt.Sprintf(`{%s}`, strings.Join(result, `, `))
}
