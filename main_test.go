package main

import (
	"bidder/router"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
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
	PlayerID string
	Balance  int
}

type winner struct {
	PlayerID string `json:"playerId,omitempty"`
	Prize    int    `json:"prize,omitempty"`
}

type tournament struct {
	TournamentID string   `json:"tournamentId,omitempty"`
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
	postJSON, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	response, err := http.Post(HOST+uri, "application/json", bytes.NewBuffer(postJSON))
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

func parseJSONPlayerBody(t *testing.T, body string) playerBalance {
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

func TestPlayerBalance(t *testing.T) {
	Convey("Test user balance", t, func() {
		resetDB(t)

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

			Convey("When I call that player's balance", func() {
				res, body := getRequest(t, "/balance?playerId=P1")
				balanceData := parseJSONPlayerBody(t, body)

				Convey("Then I get 200 status code", func() {
					So(res.StatusCode, ShouldEqual, 200)
				})

				Convey("And his balance is equal 300 points", func() {
					So(balanceData.Balance, ShouldEqual, 300)
				})
			})
		})
	})
}

func TestPlayerFund(t *testing.T) {
	Convey("Test fund user", t, func() {
		resetDB(t)

		Convey("When I fund new user with playerID P2 and points 200", func() {
			res, _ := getRequest(t, "/fund?playerId=P2&points=200")

			Convey("And I get 200 status code", func() {
				So(res.StatusCode, ShouldEqual, 200)
			})

			Convey("When I get his balance", func() {
				res, body := getRequest(t, "/balance?playerId=P2")
				balanceData := parseJSONPlayerBody(t, body)

				Convey("Then I get 200 status code", func() {
					So(res.StatusCode, ShouldEqual, 200)
				})

				Convey("And his balance is equal 200 points", func() {
					So(balanceData.Balance, ShouldEqual, 200)
				})

				Convey("And his player ID is equal P2", func() {
					So(balanceData.PlayerID, ShouldEqual, "P2")
				})
			})
		})

		Convey("When I fund existing user with PlayerID P3 and 200 points and 500 points", func() {
			getRequest(t, "/fund?playerId=P3&points=200")
			getRequest(t, "/fund?playerId=P3&points=500")

			Convey("When I get his balance", func() {
				res, body := getRequest(t, "/balance?playerId=P3")
				balanceData := parseJSONPlayerBody(t, body)

				Convey("Then I get 200 status code", func() {
					So(res.StatusCode, ShouldEqual, 200)
				})

				Convey("And his balance is equal 700 points", func() {
					So(balanceData.Balance, ShouldEqual, 700)
				})

				Convey("And his player ID is equal P3", func() {
					So(balanceData.PlayerID, ShouldEqual, "P3")
				})
			})
		})

		Convey("When I fund existing user with playerID P4 and 200 points and -500 points", func() {
			getRequest(t, "/fund?playerId=P4&points=200")
			res, _ := getRequest(t, "/fund?playerId=P4&points=-500")

			Convey("And I get 400 status code", func() {
				So(res.StatusCode, ShouldEqual, 400)
			})

			Convey("When I get his balance", func() {
				res, body := getRequest(t, "/balance?playerId=P4")
				balanceData := parseJSONPlayerBody(t, body)

				Convey("Then I get 200 status code", func() {
					So(res.StatusCode, ShouldEqual, 200)
				})

				Convey("And his balance is equal 700 points", func() {
					So(balanceData.Balance, ShouldEqual, 200)
				})

				Convey("And his player ID is equal P4", func() {
					So(balanceData.PlayerID, ShouldEqual, "P4")
				})
			})
		})
	})
}

func TestPlayerTake(t *testing.T) {
	Convey("Test take user", t, func() {
		resetDB(t)

		Convey("When I try to take points from unexisted player", func() {
			res, _ := getRequest(t, "/take?playerId=P1&points=300")

			Convey("Then I get 404 status code", func() {
				So(res.StatusCode, ShouldEqual, 404)
			})
		})

		Convey("Given I fund 300 points to new player P1", func() {
			getRequest(t, "/fund?playerId=P1&points=300")

			Convey("When I take 200 points from player P1", func() {
				res, _ := getRequest(t, "/take?playerId=P1&points=200")

				Convey("Then I get 200 status code", func() {
					So(res.StatusCode, ShouldEqual, 200)
				})

				Convey("And when I check player P1 balance", func() {
					_, body := getRequest(t, "/balance?playerId=P1")
					balanceData := parseJSONPlayerBody(t, body)

					Convey("Then his balance is equal to 100", func() {
						So(balanceData.Balance, ShouldEqual, 100)
					})
				})
			})

			Convey("When I take too many points from P1", func() {
				res, _ := getRequest(t, "/take?playerId=P1&points=500")

				Convey("Then I get 400 status code", func() {
					So(res.StatusCode, ShouldEqual, 400)
				})
			})

			Convey("When I take negative points from P1", func() {
				res, _ := getRequest(t, "/take?playerId=P1&points=-300")

				Convey("Then I get 400 status code", func() {
					So(res.StatusCode, ShouldEqual, 400)
				})
			})
		})
	})
}

func TestTournamentCreation(t *testing.T) {
	Convey("Test create tournament", t, func() {
		resetDB(t)

		Convey("Given I create P1 player with 2000 points available", func() {
			getRequest(t, "/fund?playerId=P1&points=2000")

			Convey("When I try to create a tournament with negative deposit", func() {
				res, _ := getRequest(t, "/announceTournament?tournamentId=1&deposit=-100")

				Convey("Then I get 400 status code", func() {
					So(res.StatusCode, ShouldEqual, 400)
				})

				Convey("And when I try to join it with P1 player", func() {
					res, _ = getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")

					Convey("Then I get 404 status code", func() {
						So(res.StatusCode, ShouldEqual, 404)
					})
				})
			})

			Convey("When I create the tournament with ID 1 and 1000 deposit", func() {
				res, _ := getRequest(t, "/announceTournament?tournamentId=1&deposit=1000")

				Convey("Then I get 200 status code", func() {
					So(res.StatusCode, ShouldEqual, 200)
				})

				Convey("And when I try to create a tournament with same ID 1", func() {
					res, _ := getRequest(t, "/announceTournament?tournamentId=1&deposit=1000")
					Convey("Then I get 400 status code", func() {
						So(res.StatusCode, ShouldEqual, 400)
					})
				})

				Convey("And when I try to join it with P1 player", func() {
					res, _ := getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")

					Convey("Then I get 200 status code", func() {
						So(res.StatusCode, ShouldEqual, 200)
					})

					Convey("And when I check player P1 balance", func() {
						_, body := getRequest(t, "/balance?playerId=P1")
						balanceData := parseJSONPlayerBody(t, body)

						Convey("Then his balance is equal to 1000", func() {
							So(balanceData.Balance, ShouldEqual, 1000)
						})
					})
				})
			})
		})
	})
}

func TestTournamentJoin(t *testing.T) {
	Convey("Test join tournament", t, func() {
		resetDB(t)

		Convey("Given I set a player P1 with 1500 points", func() {
			getRequest(t, "/fund?playerId=P1&points=1500")

			Convey("When I announce a tournament with 500 points deposit", func() {
				getRequest(t, "/announceTournament?tournamentId=1&deposit=500")

				Convey("And I join player P1 to the tournament", func() {
					res, _ := getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")

					Convey("Then I get 200 status code", func() {
						So(res.StatusCode, ShouldEqual, 200)
					})

					Convey("Then I check player P1 balance", func() {
						_, body := getRequest(t, "/balance?playerId=P1")
						balanceData := parseJSONPlayerBody(t, body)

						Convey("Then his balance is equal to 1000", func() {
							So(balanceData.Balance, ShouldEqual, 1000)
						})
					})

					Convey("When I try to join same tournament with P1 player", func() {
						res, _ := getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")

						Convey("Then I get 400 status code", func() {
							So(res.StatusCode, ShouldEqual, 400)
						})
					})
				})

				Convey("And given I set player P2 with 1000 points", func() {
					getRequest(t, "/fund?playerId=P2&points=1000")

					Convey("And I join player P1 with backer P2 to the tournament", func() {
						res, _ := getRequest(t, "/joinTournament?tournamentId=1&playerId=P1&backerId=P2")

						Convey("Then I get 200 status code", func() {
							So(res.StatusCode, ShouldEqual, 200)
						})

						Convey("Then I check player P1 balance", func() {
							_, body := getRequest(t, "/balance?playerId=P1")
							balanceData := parseJSONPlayerBody(t, body)

							Convey("And his balance is equal to 1250", func() {
								So(balanceData.Balance, ShouldEqual, 1250)
							})
						})

						Convey("Then I check player P2 balance", func() {
							_, body := getRequest(t, "/balance?playerId=P2")
							balanceData := parseJSONPlayerBody(t, body)

							Convey("And his balance is equal to 750", func() {
								So(balanceData.Balance, ShouldEqual, 750)
							})
						})
					})
				})

				Convey("And given I set player P3 with 250 points", func() {
					getRequest(t, "/fund?playerId=P3&points=250")

					Convey("When I try to join player P3 to the tournament", func() {
						res, _ := getRequest(t, "/joinTournament?tournamentId=1&playerId=P3")

						Convey("Then I get 400 status code", func() {
							So(res.StatusCode, ShouldEqual, 400)
						})
					})
				})
			})
		})
	})
}

func TestTournamentResult(t *testing.T) {
	Convey("Test result tournament", t, func() {
		resetDB(t)

		Convey("Given I set player P1 with 1000 points", func() {
			getRequest(t, "/fund?playerId=P1&points=1000")

			Convey("And I announce a tournament with deposit 500", func() {
				getRequest(t, "/announceTournament?tournamentId=1&deposit=500")

				Convey("Then I join the tournament with player P1", func() {
					res, _ := getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")

					Convey("And I get 200 status code", func() {
						So(res.StatusCode, ShouldEqual, 200)
					})

					Convey("And When I result tournament with P1 as a winner with 1000 win", func() {
						winner1 := winner{PlayerID: "P1", Prize: 1000}
						result := tournament{TournamentID: "1", Winners: []winner{winner1}}

						res, _ := postRequest(t, "/resultTournament", result)

						Convey("Then I get 200 status code", func() {
							So(res.StatusCode, ShouldEqual, 200)
						})

						Convey("Then I check player P1 balance", func() {
							_, body := getRequest(t, "/balance?playerId=P1")
							balanceData := parseJSONPlayerBody(t, body)

							Convey("And his balance is equal to 1500", func() {
								So(balanceData.Balance, ShouldEqual, 1500)
							})
						})

						Convey("And when I try to result tournament second time", func() {
							res, _ := postRequest(t, "/resultTournament", result)

							Convey("Then I get 400 status code", func() {
								So(res.StatusCode, ShouldEqual, 400)
							})
						})
					})
				})

				Convey("And I set player P2 with 250 points", func() {
					getRequest(t, "/fund?playerId=P2&points=250")

					Convey("Then I join the tournament with player P1 and backer P2", func() {
						res, _ := getRequest(t, "/joinTournament?tournamentId=1&playerId=P1&backerId=P2")

						Convey("Then I get 200 status code", func() {
							So(res.StatusCode, ShouldEqual, 200)
						})

						Convey("And when I result tournament with P1 as a winner with 1000 win", func() {
							winner1 := winner{PlayerID: "P1", Prize: 1000}
							result := tournament{TournamentID: "1", Winners: []winner{winner1}}

							res, _ := postRequest(t, "/resultTournament", result)

							Convey("Then I get 200 status code", func() {
								So(res.StatusCode, ShouldEqual, 200)
							})

							Convey("Then I check player P1 balance", func() {
								_, body := getRequest(t, "/balance?playerId=P1")
								balanceData := parseJSONPlayerBody(t, body)

								Convey("And his balance is equal to 1250", func() {
									So(balanceData.Balance, ShouldEqual, 1250)
								})
							})

							Convey("Then I check player P2 balance", func() {
								_, body := getRequest(t, "/balance?playerId=P2")
								balanceData := parseJSONPlayerBody(t, body)

								Convey("And his balance is equal to 500", func() {
									So(balanceData.Balance, ShouldEqual, 500)
								})
							})
						})
					})
				})
			})
		})
	})
}

func TestConcurrentFund(t *testing.T) {
	Convey("Test fund endpoint concurrently", t, func() {
		resetDB(t)

		Convey("Given I call fund endpoint 10 times in different goroutines", func() {
			const requestsCount = 10
			responses := make(chan int, requestsCount)

			var wg sync.WaitGroup
			for i := 1; i <= requestsCount; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					res, _ := getRequest(t, "/fund?playerId=P1&points=100")
					responses <- res.StatusCode
				}()
			}
			wg.Wait()

			Convey("When I count every successful response", func() {
				successResponses := 0
				for i := 1; i <= requestsCount; i++ {
					if 200 == <-responses {
						successResponses++
					}
				}

				Convey("Then I get 10 responses with 200 status code", func() {
					So(successResponses, ShouldEqual, requestsCount)
				})

				Convey("And when I check the player's balance", func() {
					_, body := getRequest(t, "/balance?playerId=P1")
					balanceData := parseJSONPlayerBody(t, body)

					Convey("Then his balance should equal 1000", func() {
						So(balanceData.Balance, ShouldEqual, 1000)
					})
				})
			})
		})
	})
}

func TestConcurrentTake(t *testing.T) {
	Convey("Test take endpoint concurrently", t, func() {
		resetDB(t)

		Convey("Given I set player P1 with 100 points", func() {
			getRequest(t, "/fund?playerId=P1&points=100")

			Convey("When I call take endpoint 5 times in different goroutines", func() {
				const requestsCount = 5
				responses := make(chan int, requestsCount)

				var wg sync.WaitGroup
				for i := 1; i <= requestsCount; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						res, _ := getRequest(t, "/take?playerId=P1&points=100")
						responses <- res.StatusCode

					}()
				}
				wg.Wait()

				Convey("Then I count every successful response", func() {
					successResponses := 0
					for i := 1; i <= requestsCount; i++ {
						if 200 == <-responses {
							successResponses++
						}
					}

					Convey("And I get 1 successful response", func() {
						So(successResponses, ShouldEqual, 1)
					})

					Convey("And when I check player P1 balance", func() {
						_, body := getRequest(t, "/balance?playerId=P1")
						balanceData := parseJSONPlayerBody(t, body)

						Convey("Then his balance should equal to 0", func() {
							So(balanceData.Balance, ShouldEqual, 0)
						})
					})
				})
			})
		})
	})
}

func TestConcurrentTournamentAnnounce(t *testing.T) {
	Convey("Test announceTournament endpoint concurrently", t, func() {
		resetDB(t)

		Convey("When I call announceTournament endpoint 5 times in different goroutines", func() {
			const requestsCount = 5
			responses := make(chan int, requestsCount)

			var wg sync.WaitGroup
			for i := 1; i <= requestsCount; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					res, _ := getRequest(t, "/announceTournament?tournamentId=1&deposit=500")
					responses <- res.StatusCode

				}()
			}
			wg.Wait()

			Convey("Then I count every successful response", func() {
				successResponses := 0
				for i := 1; i <= requestsCount; i++ {
					if 200 == <-responses {
						successResponses++
					}
				}

				Convey("And I get 1 successful response", func() {
					So(successResponses, ShouldEqual, 1)
				})
			})
		})
	})
}

func TestConcurrentTournamentJoin(t *testing.T) {
	Convey("Test joinTournament endpoint concurrently", t, func() {
		resetDB(t)

		Convey("Given I set player P1 with 1000 points", func() {
			getRequest(t, "/fund?playerId=P1&points=1000")

			Convey("And given I announce tournament with 500 points deposit", func() {
				getRequest(t, "/announceTournament?tournamentId=1&deposit=500")

				Convey("When I call joinTournament endpoint 5 times in different goroutines", func() {
					const requestsCount = 5
					responses := make(chan int, requestsCount)

					var wg sync.WaitGroup
					for i := 1; i <= requestsCount; i++ {
						wg.Add(1)
						go func() {
							defer wg.Done()
							res, _ := getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")
							responses <- res.StatusCode

						}()
					}
					wg.Wait()

					Convey("Then I count every successful response", func() {
						successResponses := 0
						for i := 1; i <= requestsCount; i++ {
							if 200 == <-responses {
								successResponses++
							}
						}

						Convey("And I get 1 successful response", func() {
							So(successResponses, ShouldEqual, 1)
						})

						Convey("And when I check player P1 balance", func() {
							_, body := getRequest(t, "/balance?playerId=P1")
							balanceData := parseJSONPlayerBody(t, body)

							Convey("Then his balance should equal to 500", func() {
								So(balanceData.Balance, ShouldEqual, 500)
							})
						})
					})
				})
			})
		})
	})
}

func TestConcurrentTournamentResult(t *testing.T) {
	Convey("Test result endpoint concurrently", t, func() {
		resetDB(t)

		Convey("Given I set players P1, P2, P3", func() {
			getRequest(t, "/fund?playerId=P1&points=2000")
			getRequest(t, "/fund?playerId=P2&points=2000")
			getRequest(t, "/fund?playerId=P3&points=1000")

			Convey("And given I announce tournament with 1000 points deposit", func() {
				getRequest(t, "/announceTournament?tournamentId=1&deposit=1000")

				Convey("And given I join those players to the tournament", func() {
					getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")
					getRequest(t, "/joinTournament?tournamentId=1&playerId=P2&backerId=P3")

					Convey("When I call resultTournament endpoint 5 times in different goroutines", func() {
						winner1 := winner{PlayerID: "P1", Prize: 1000}
						winner2 := winner{PlayerID: "P2", Prize: 500}
						result := tournament{TournamentID: "1", Winners: []winner{winner1, winner2}}

						const requestsCount = 5
						responses := make(chan int, requestsCount)

						var wg sync.WaitGroup
						for i := 1; i <= requestsCount; i++ {
							wg.Add(1)
							go func() {
								defer wg.Done()
								res, _ := postRequest(t, "/resultTournament", result)
								responses <- res.StatusCode

							}()
						}
						wg.Wait()

						Convey("Then I count every successful response", func() {
							successResponses := 0
							for i := 1; i <= requestsCount; i++ {
								if 200 == <-responses {
									successResponses++
								}
							}

							Convey("And I get 1 successful response", func() {
								So(successResponses, ShouldEqual, 1)
							})

							Convey("And when I check player P1 balance", func() {
								_, body := getRequest(t, "/balance?playerId=P1")
								balanceData := parseJSONPlayerBody(t, body)

								Convey("Then his balance should equal to 2000", func() {
									So(balanceData.Balance, ShouldEqual, 2000)
								})
							})

							Convey("And when I check player P2 balance", func() {
								_, body := getRequest(t, "/balance?playerId=P2")
								balanceData := parseJSONPlayerBody(t, body)

								Convey("Then his balance should equal to 1750", func() {
									So(balanceData.Balance, ShouldEqual, 1750)
								})
							})

							Convey("And when I check player P3 balance", func() {
								_, body := getRequest(t, "/balance?playerId=P3")
								balanceData := parseJSONPlayerBody(t, body)

								Convey("Then his balance should equal to 750", func() {
									So(balanceData.Balance, ShouldEqual, 750)
								})
							})
						})
					})
				})
			})
		})
	})
}
