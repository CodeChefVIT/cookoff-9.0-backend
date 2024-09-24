package submission

import (
	"context"
	"fmt"

	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/google/uuid"
)

func CheckStatus(ctx context.Context, subid uuid.UUID) (bool, error) {
	status, err := database.Queries.GetSubmissionStatusByID(ctx, subid)
	if err != nil {
		logger.Errof(fmt.Sprintf("Error while getting submission status: %v", err.Error()))
		return false, err
	}

	if status == nil {
		return false, nil
	}

	return *status == "done", nil
}

func GetSubResult(ctx context.Context, subid uuid.UUID) (resultresp, error) {
	submission, err := database.Queries.GetSubmissionByID(ctx, subid)
	if err != nil {
		logger.Errof(fmt.Sprintf("Error while getting submission result: %v", err.Error()))
		return resultresp{}, err
	}

	var desc string
	if submission.Description == nil {
		desc = ""
	} else {
		desc = *submission.Description
	}

	resp := resultresp{
		ID:          submission.ID.String(),
		QuestionID:  submission.QuestionID.String(),
		Passed:      int(submission.TestcasesPassed.Int32),
		Failed:      int(submission.TestcasesFailed.Int32),
		Description: desc,
	}

	return resp, nil
}
