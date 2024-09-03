CREATE TABLE "user" (
	"id" UUID NOT NULL UNIQUE,
	"submissions" UUID,
	"email" VARCHAR NOT NULL,
	"regNo" VARCHAR NOT NULL,
	"password" VARCHAR NOT NULL,
	"role" CHAR NOT NULL,
	"roundQualified" INTEGER NOT NULL,
	"score" INTEGER,
	"name" VARCHAR NOT NULL,
	PRIMARY KEY("id")
);


CREATE TABLE "questions" (
	"id" UUID NOT NULL UNIQUE,
	"description" TEXT,
	"title" VARCHAR,
	"inputFormat" TEXT,
	"points" INTEGER,
	"round" SMALLINT NOT NULL,
	"constraints" TEXT,
	"outputFormat" TEXT,
	"testcases" UUID,
	PRIMARY KEY("id")
);


CREATE TABLE "submissions" (
	"id" UUID NOT NULL UNIQUE,
	"question_id" UUID NOT NULL,
	"testcases_passed" INTEGER NOT NULL,
	"testcases_failed" INTEGER NOT NULL,
	"runtime" DECIMAL NOT NULL,
	"sub time" TIMESTAMP NOT NULL,
	"testcases_id" UUID,
	"language_id" INTEGER,
	"description" TEXT,
	"memory" INTEGER,
	"user_id" UUID,
	PRIMARY KEY("id")
);


CREATE TABLE "testcases" (
	"id" UUID NOT NULL UNIQUE,
	"expected_output" TEXT,
	"memory" VARCHAR,
	"input" TEXT,
	"hidden" BOOLEAN,
	"runtime" TIME,
	PRIMARY KEY("id")
);


ALTER TABLE "submissions"
ADD FOREIGN KEY("question_id") REFERENCES "questions"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "questions"
ADD FOREIGN KEY("testcases") REFERENCES "testcases"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "submissions"
ADD FOREIGN KEY("testcases_id") REFERENCES "testcases"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "user"
ADD FOREIGN KEY("submissions") REFERENCES "submissions"("user_id")
ON UPDATE NO ACTION ON DELETE NO ACTION;