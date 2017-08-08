BEGIN;

-- CREATE TABLE "tournaments" ----------------------------------
CREATE TABLE "public"."tournaments" (
	"id" Serial NOT NULL UNIQUE,
	"deposit" Integer NOT NULL,
	"finished" Boolean DEFAULT FALSE NOT NULL,
 PRIMARY KEY ( "id" ) );
-- -------------------------------------------------------------;

COMMIT;
