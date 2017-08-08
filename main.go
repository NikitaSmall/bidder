package main

import (
	"log"

	"bidder/router"
	"bidder/util"
)

func main() {
	defer func() { util.DBConnect.Close() }()
	log.Println("Welcome to the Bidder app!")

	r := router.New()
	r.Run()
}
