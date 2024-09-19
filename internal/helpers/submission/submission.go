package submission

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/google/uuid"
)

type Token struct {
	Token string `json:"token"`
}

type Payload struct {
	Submissions []Submission `json:"submissions"`
}

func CreateSubmission(ctx context.Context, question_id uuid.UUID, language_id int, source string) ([]byte, error) {
	callback_url := os.Getenv("CALLBACK_URL")

	testcases, err := database.Queries.GetTestCases(ctx, db.GetTestCasesParams{QuestionID: question_id, Column2: false})
	if err != nil {
		logger.Errof("Error getting test cases for question_id %d: %v", question_id, err)
		return nil, err
	}
	payload := Payload{
		Submissions: make([]Submission, len(testcases)),
	}

	var runtime_mut int
	switch language_id {
	case 50, 54, 60, 73, 63:
		runtime_mut = 1
	case 51, 62:
		runtime_mut = 2
	case 68:
		runtime_mut = 3
	case 71:
		runtime_mut = 5
	default:
		return nil, errors.New("Invalid language ID")
	}

	for i, testcase := range testcases {
		runtime, _ := testcase.Runtime.Float64Value()
		payload.Submissions[i] = Submission{
			SourceCode: B64(source),
			LanguageID: language_id,
			Input:      B64(*testcase.Input),
			Output:     B64(*testcase.ExpectedOutput),
			Runtime:    runtime.Float64 * float64(runtime_mut),
			Callback:   callback_url,
		}
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		logger.Errof("Error marshaling payload: %v", err)
		return nil, err
	}

	return payloadJSON, nil
}

func StoreTokens(ctx context.Context, subID uuid.UUID, r *http.Response) error {
	var tokens []Token
	err := json.NewDecoder(r.Body).Decode(&tokens)
	if err != nil {
		return errors.New("Invalid request payload")
	}

	for _, t := range tokens {
		err := Tokens.AddToken(ctx, t.Token, subID.String())
		if err != nil {
			return errors.New("Failed to add token")
		}
	}
	return nil
}

func GetSubID(ctx context.Context, token string) (string, error) {
	subID, err := Tokens.GetSubID(ctx, token)
	if err != nil {
		return "", err
	}
	return subID, nil
}
