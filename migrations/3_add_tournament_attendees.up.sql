BEGIN;


-- CREATE TABLE "tournament_attendees" -------------------------
CREATE TABLE "public"."tournament_attendees" (
	"id" Serial NOT NULL,
	"tournament_id" Integer NOT NULL references tournaments(id),
	"player_id" Integer NOT NULL references players(id),
	"finished" Boolean DEFAULT false NOT NULL,

	-- I've read the task with more attention and found one example of passing the prize number
	-- to the system. Thus I comment this `guessing` functional out.
	-- "prize" Integer NOT NULL,

	"backers" Integer[] DEFAULT array[]::Integer[] NOT NULL,
 PRIMARY KEY ( "id" )
 );
-- -------------------------------------------------------------;

COMMIT;
