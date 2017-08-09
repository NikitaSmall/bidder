BEGIN;


-- CREATE TABLE "tournament_attendees" -------------------------
CREATE TABLE "public"."tournament_attendees" (
	"id" Serial NOT NULL,
	"tournament_id" Integer NOT NULL references tournaments(id) ON DELETE CASCADE,
	"player_id" Character Varying( 256 ) NOT NULL references players(player_id) ON DELETE CASCADE,
	"finished" Boolean DEFAULT false NOT NULL,

	-- I've read the task with more attention and found one example of passing the prize number
	-- to the system. Thus I comment this `guessing` functional out.
	-- "prize" Integer NOT NULL CHECK (prize >= 0),

	"backers" Character Varying( 256 )[] DEFAULT array[]::Character Varying( 256 )[] NOT NULL,
 PRIMARY KEY ( "id" )
 );
-- -------------------------------------------------------------;

COMMIT;
