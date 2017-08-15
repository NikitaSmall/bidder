# Bidder

Bidder is a test application to create tournaments, players,
add them to tournaments and then finish tournaments with a prizes.

## Run

Just use `docker-compose up`, then reach the application at `localhost:3000` (after DB startup).

## Testing

For testing purposes run `docker-compose run web go test`. You should get
the detailed report. You may run this command with `-v` flag to see details.
(this is possible to run the tests without any database installed
as tests are stored in the docker image [for preview purpose])
