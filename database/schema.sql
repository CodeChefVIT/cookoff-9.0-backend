

CREATE TABLE "user" (
	"id" UUID NOT NULL UNIQUE,
	"email" TEXT NOT NULL UNIQUE, 
	"regNo" TEXT NOT NULL UNIQUE,
	"password" TEXT NOT NULL,
	"role" TEXT NOT NULL CHECK ("role" IN ('user', 'admin')) REQUIRED, 
	"roundQualified" INTEGER NOT NULL DEFAULT 0,
	"score" INTEGER DEFAULT 0,
	"name" TEXT NOT NULL,
	PRIMARY KEY("id")
);

CREATE TABLE "questions" (
	"id" UUID NOT NULL UNIQUE,
	"description" TEXT,
	"title" TEXT,
	"inputFormat" TEXT,
	"points" INTEGER,
	"round" INTEGER NOT NULL,
	"constraints" TEXT,
	"outputFormat" TEXT,
	PRIMARY KEY("id")
);

CREATE TABLE "submissions" (
	"id" UUID NOT NULL UNIQUE,
	"question_id" UUID NOT NULL,
	"testcases_passed" INTEGER DEFAULT 0,
	"testcases_failed" INTEGER DEFAULT 0,
	"runtime" DECIMAL,
	"sub_time" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	"testcases_id" UUID,
	"language_id" INTEGER NOT NULL,
	"description" TEXT,
	"memory" INTEGER,
	"user_id" UUID,
	"ref_id" TEXT,
	PRIMARY KEY("id")
);

CREATE TABLE "testcases" (
	"id" UUID NOT NULL UNIQUE,
	"expected_output" TEXT,
	"memory" TEXT,
	"input" TEXT,
	"hidden" BOOLEAN,
	"runtime" TIME,
	"question_id" UUID NOT NULL,
	PRIMARY KEY("id")
);

ALTER TABLE "submissions"
ADD FOREIGN KEY("question_id") REFERENCES "questions"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "testcases"
ADD FOREIGN KEY("question_id") REFERENCES "questions"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "submissions"
ADD FOREIGN KEY("testcases_id") REFERENCES "testcases"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE "submissions"
ADD FOREIGN KEY("user_id") REFERENCES "user"("id")
ON UPDATE NO ACTION ON DELETE NO ACTION;
