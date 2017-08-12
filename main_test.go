package main

import (
	"bidder/router"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// host for test server
const HOST = "http://localhost:3001"

// and initialize the server for testing
func init() {
	r := router.New()
	go r.Run(":3001")
	time.Sleep(time.Second)
}

// the set of helper structs and functions to avoid code duplication
// and reduce overall code amount and complexity
type playerBalance struct {
	PlayerId string
	Balance  int
}

type winner struct {
	PlayerId string `json:"playerId,omitempty"`
	Prize    int    `json:"prize,omitempty"`
}

type tournament struct {
	TournamentId string   `json:"tournamentId,omitempty"`
	Winners      []winner `json:"winners,omitempty"`
}

func getRequest(t *testing.T, uri string) (*http.Response, string) {
	response, err := http.Get(HOST + uri)
	if err != nil {
		t.Fatal(err)
	}

	return response, getResponceBody(t, response)
}

func postRequest(t *testing.T, uri string, data interface{}) (*http.Response, string) {
	postJson, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	response, err := http.Post(HOST+uri, "application/json", bytes.NewBuffer(postJson))
	if err != nil {
		t.Fatal(err)
	}

	return response, getResponceBody(t, response)
}

func getResponceBody(t *testing.T, response *http.Response) string {
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func parseJsonPlayerBody(t *testing.T, body string) playerBalance {
	var data playerBalance

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		t.Fatal(err)
	}
	return data
}

func resetDB(t *testing.T) {
	getRequest(t, "/reset")
}

// Actual tests start here:
func TestGeneralEndpoints(t *testing.T) {
	Convey("When I call wrong page", t, func() {
		res, _ := getRequest(t, "/wrong")

		Convey("Then I get 404 status code", func() {
			So(res.StatusCode, ShouldEqual, 404)
		})
	})
}

func TestPlayerRelatedEndpoints(t *testing.T) {
	resetDB(t)

	Convey("Test user balance", t, func() {
		Convey("When I call unexisting player", func() {
			res, _ := getRequest(t, "/balance?playerId=P1")
			Convey("Then I get 404 status code", func() {
				So(res.StatusCode, ShouldEqual, 404)
			})
		})

		Convey("Given I create P1 player with 300 points", func() {
			res, _ := getRequest(t, "/fund?playerId=P1&points=300")

			Convey("And I get 200 status code", func() {
				So(res.StatusCode, ShouldEqual, 200)
			})
		})

		Convey("When I call that player's balance", func() {
			res, body := getRequest(t, "/balance?playerId=P1")
			balanceData := parseJsonPlayerBody(t, body)

			Convey("Then I get 200 status code", func() {
				So(res.StatusCode, ShouldEqual, 200)
			})

			Convey("And his balance is equal 300 points", func() {
				So(balanceData.Balance, ShouldEqual, 300)
			})
		})
	})
}
