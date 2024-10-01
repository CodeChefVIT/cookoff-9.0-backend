-- Modify "submission_results" table
ALTER TABLE "public"."submission_results" ADD COLUMN "testcase_id" uuid NULL, ADD COLUMN "status" text NOT NULL;
-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "score" TYPE numeric, ALTER COLUMN "score" SET NOT NULL, ADD COLUMN "is_banned" boolean NOT NULL DEFAULT false;
