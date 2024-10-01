package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/CodeChefVIT/cookoff-backend/internal/controllers"
	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/submission"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	TypeProcessSubmission = "submission:process"
	SubmissionDoneStatus  = "DONE"
)

// ProcessSubmissionTask processes the submission task based on status.
func ProcessSubmissionTask(ctx context.Context, t *asynq.Task) error {
	var data controllers.Data
	logger.Infof("Processing task: %v", t.Type)
	logger.Infof("Payload: %v", string(t.Payload()))

	if err := json.Unmarshal(t.Payload(), &data); err != nil {
		log.Printf("Error unmarshalling task payload: %v\n", err)
		return err
	}

	timeValue, err := parseTime(data.Time)
	if err != nil {
		log.Println("Error parsing time value: ", err)
		return err
	}

	value, testcase, err := submission.GetSubID(ctx, data.Token)
	if err != nil {
		log.Println("Error getting submission ID from token: ", err)
		return err
	}

	idUUID, err := uuid.Parse(value)
	if err != nil {
		log.Fatalf("Error parsing UUID: %v", err)
	}

	testidUUID, err := uuid.Parse(testcase)
	if err != nil {
		log.Fatalf("Error parsing UUID: %v", err)
	}

	switch data.Status.ID {
	case "1":
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"In Queue",
		)
	case "2":
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"Processing",
		)
	case "3":
		// testcasesPassed++
		err = handleCompilationError(ctx, idUUID, data, int(timeValue*1000), testidUUID, "success")
	case "4":
		// testcasesFailed++
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"wrong answer",
		)
	case "5":
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"Time Limit Exceeded",
		)
	case "6":
		// testcasesFailed++
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"Compilation error",
		)
	case "7", "8", "9", "10", "11", "12":
		// testcasesFailed++
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"Runtime error",
		)
	case "13":
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"Internal Error",
		)
	case "14":
		err = handleCompilationError(
			ctx,
			idUUID,
			data,
			int(timeValue*1000),
			testidUUID,
			"Exec Format Error",
		)
	}

	if err != nil {
		return err
	}

	if err := submission.Tokens.DeleteToken(ctx, data.Token); err != nil {
		log.Println("Error deleting token: ", err)
		return err
	}

	tokenCount, err := submission.Tokens.GetTokenCount(ctx, value)
	if err != nil {
		log.Println("Error getting token count: ", err)
		return err
	}

	fmt.Println("Token :- ", tokenCount)

	if tokenCount == 0 {
		err = submission.UpdateSubmission(ctx, idUUID)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseTime(timeStr string) (float64, error) {
	if timeStr == "" {
		log.Println("Time value is empty, setting time to 0 for this submission.")
		return 0, nil
	}

	timeValue, err := strconv.ParseFloat(timeStr, 64)
	if err != nil {
		return 0, err
	}
	return timeValue, nil
}

func handleCompilationError(
	ctx context.Context,
	idUUID uuid.UUID,
	data controllers.Data,
	time int,
	testcase uuid.UUID,
	status string,
) error {
	subID, err := uuid.NewV7()
	if err != nil {
		log.Println("Error updating submission for compilation error: ", err)
		return err
	}

	err = database.Queries.CreateSubmissionStatus(ctx, db.CreateSubmissionStatusParams{
		ID:           subID,
		SubmissionID: idUUID,
		TestcaseID:   uuid.NullUUID{UUID: testcase, Valid: true},
		Runtime:      pgtype.Numeric{Int: big.NewInt(int64(time)), Valid: true},
		Memory:       pgtype.Numeric{Int: big.NewInt(int64(data.Memory)), Valid: true},
		Description:  &data.Status.Description,
		Status:       status,
	})
	if err != nil {
		log.Println("Error creating submission status error: ", err)
		return err
	}
	return nil
}
