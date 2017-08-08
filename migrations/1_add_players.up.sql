BEGIN;

-- CREATE TABLE "players" --------------------------------------
CREATE TABLE "public"."players" (
	"player_id" Character Varying( 256 ) NOT NULL UNIQUE,
	"points" Integer NOT NULL CHECK (points >= 0),
 PRIMARY KEY ( "player_id" ) );
-- -------------------------------------------------------------;

COMMIT;
