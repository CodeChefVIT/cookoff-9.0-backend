package submission

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Token struct {
	Token string `json:"token"`
}

var ErrNoTestCases = fmt.Errorf("no testcases present for question_id")

type Payload struct {
	Submissions []Submission `json:"submissions"`
}

func CreateSubmission(
	ctx context.Context,
	question_id uuid.UUID,
	language_id int,
	source string,
) ([]byte, []uuid.UUID, error) {
	callback_url := os.Getenv("CALLBACK_URL")
	var testcases_id []uuid.UUID
	testcases, err := database.Queries.GetTestCases(
		ctx,
		db.GetTestCasesParams{QuestionID: question_id, Column2: false},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, fmt.Errorf("no testcases exist for this question")
		}
		return nil, nil, fmt.Errorf("error getting test cases for question_id %d: %v", question_id, err)
	}

	if len(testcases) == 0 {
		return nil, nil, ErrNoTestCases
	}

	payload := Payload{
		Submissions: make([]Submission, len(testcases)),
	}

	runtime_mut, err := RuntimeMut(language_id)
	if err != nil {
		return nil, nil, err
	}

	for i, testcase := range testcases {
		runtime, _ := testcase.Runtime.Float64Value()
		testcases_id = append(testcases_id, testcase.ID)
		payload.Submissions[i] = Submission{
			SourceCode: B64(source),
			LanguageID: language_id,
			Input:      B64(testcase.Input),
			Output:     B64(testcase.ExpectedOutput),
			Runtime:    runtime.Float64 * float64(runtime_mut),
			Callback:   callback_url,
		}
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling payload: %v", err)
	}

	return payloadJSON, testcases_id, nil
}

func StoreTokens(ctx context.Context, subID uuid.UUID, resp []byte, testcases_id []uuid.UUID) error {
	var tokens []Token
	err := json.Unmarshal(resp, &tokens)
	if err != nil {
		return fmt.Errorf("invalid request payload")
	}

	for i, t := range tokens {
		err := Tokens.AddToken(ctx, t.Token, subID.String(), testcases_id[i].String())
		if err != nil {
			return fmt.Errorf("failed to add token")
		}
	}
	return nil
}

func GetSubID(ctx context.Context, token string) (string, string, error) {
	subID, err := Tokens.GetSubID(ctx, token)
	if err != nil {
		return "", "", err
	}
	temp := strings.Split(subID, ":")
	return temp[0], temp[1], nil
}
