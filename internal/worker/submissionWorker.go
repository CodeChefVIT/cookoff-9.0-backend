package worker

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"strconv"

	"github.com/CodeChefVIT/cookoff-backend/internal/controllers"
	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/submission"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
)

const TypeProcessSubmission = "submission:process"

func ProcessSubmissionTask(ctx context.Context, t *asynq.Task) error {
	var data controllers.Data
	log.Print("Processing task: ", t.Type)
	log.Println("Payload: ", string(t.Payload()))
	log.Println(data)

	if err := json.Unmarshal(t.Payload(), &data); err != nil {
		log.Printf("Error unmarshalling task payload: %v\n", err)
		return err
	}

	timeValue, err := strconv.ParseFloat(data.Time, 64)
	if err != nil {
		log.Println("Error parsing time value: ", err)
		return err
	}

	value, err := submission.GetSubID(ctx, data.Token)
	if err != nil {
		log.Println("Error getting submission ID from token: ", err)
		return err
	}

	idUUID, err := uuid.Parse(value)
	if err != nil {
		log.Fatalf("Error parsing UUID: %v", err)
	}

	sub, err := database.Queries.GetSubmission(ctx, idUUID)
	if err != nil {
		log.Println("Error retrieving submission: ", err)
		return err
	}

	testcasesPassed := int(sub.TestcasesPassed.Int32)
	testcasesFailed := int(sub.TestcasesFailed.Int32)

	switch data.Status.ID {
	case "3":
		testcasesPassed++
	case "4":
		testcasesFailed++
	}

	err = database.Queries.UpdateSubmission(ctx, db.UpdateSubmissionParams{
		TestcasesPassed: pgtype.Int4{Int32: int32(testcasesPassed), Valid: true},
		TestcasesFailed: pgtype.Int4{Int32: int32(testcasesFailed), Valid: true},
		Runtime:         pgtype.Numeric{Int: big.NewInt(int64(timeValue * 1000)), Valid: true},
		Memory:          pgtype.Int4{Int32: int32(data.Memory), Valid: true},
		ID:              idUUID,
	})

	if err != nil {
		log.Println("Error updating submission: ", err)
		return err
	}

	log.Printf("Submission ID: %v Testcases Passed: %v Testcases Failed: %v\n", value, testcasesPassed, testcasesFailed)
	return nil
}
