-- Modify "submissions" table
ALTER TABLE "public"."submissions" DROP COLUMN "testcase_id", ALTER COLUMN "memory" TYPE numeric;
-- Create "submission_results" table
CREATE TABLE "public"."submission_results" ("id" uuid NOT NULL, "submission_id" uuid NOT NULL, "runtime" numeric NOT NULL, "memory" numeric NOT NULL, "description" text NULL, PRIMARY KEY ("id"), CONSTRAINT "submission_results_submission_id_fkey" FOREIGN KEY ("submission_id") REFERENCES "public"."submissions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
