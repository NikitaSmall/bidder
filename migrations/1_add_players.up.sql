BEGIN;

-- CREATE TABLE "players" --------------------------------------
CREATE TABLE "public"."players" (
	"id" Serial NOT NULL UNIQUE,
	"points" Integer NOT NULL,
 PRIMARY KEY ( "id" ) );
-- -------------------------------------------------------------;

COMMIT;
