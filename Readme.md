# Bidder

Bidder is a test application to create tournaments, players,
add them to tournaments and then finish tournaments with a prizes.

## Run

Just use `docker-compose up`, then reach the application at `localhost:3000` (after DB startup).

## Testing

For testing purposes run `docker-compose run web go test -v`. You should get
the detailed report.
(this is possible as tests are stored in the image [for preview purpose])
