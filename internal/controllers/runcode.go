package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/CodeChefVIT/cookoff-backend/internal/db"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/auth"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/database"
	httphelpers "github.com/CodeChefVIT/cookoff-backend/internal/helpers/http"
	logger "github.com/CodeChefVIT/cookoff-backend/internal/helpers/logging"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/submission"
	"github.com/CodeChefVIT/cookoff-backend/internal/helpers/validator"
	"github.com/google/uuid"
)

type resp struct {
	Result          []submission.Judgeresp `json:"result"`
	TestCasesPassed int                    `json:"no_testcases_passed"`
}

func RunCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req subreq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := validator.ValidatePayload(w, req); err != nil {
		httphelpers.WriteError(
			w,
			http.StatusNotAcceptable,
			"Please provide values for all required fields.",
		)
		return
	}

	question_id, _ := uuid.Parse(req.QuestionID)
	userID, _ := auth.GetUserID(w, r)

	qualified, err := auth.VerifyRound(ctx, userID, question_id)
	if err != nil {
		httphelpers.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if !qualified {
		httphelpers.WriteError(w, http.StatusForbidden, "User is not qualified for this round")
		return
	}

	testcases, err := database.Queries.GetTestCases(
		ctx,
		db.GetTestCasesParams{QuestionID: question_id, Column2: true},
	)
	if err != nil {
		httphelpers.WriteError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("error getting test cases for question_id %d: %v", question_id, err),
		)
		return
	}

	judge0URL, err := url.Parse(JUDGE0_URI + "/submissions")
	if err != nil {
		httphelpers.WriteError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error parsing Judge0 URL: %v", err),
		)
		return
	}
	params := url.Values{}
	params.Add("base64_encoded", "true")
	params.Add("wait", "true")

	var payload submission.Submission
	response := resp{
		Result: make([]submission.Judgeresp, len(testcases)),
	}

	runtime_mut, err := submission.RuntimeMut(req.LanguageID)
	if err != nil {
		httphelpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	for i, testcase := range testcases {
		runtime, _ := testcase.Runtime.Float64Value()
		payload = submission.Submission{
			LanguageID: req.LanguageID,
			SourceCode: submission.B64(req.SourceCode),
			Input:      submission.B64(testcase.Input),
			Output:     submission.B64(testcase.ExpectedOutput),
			Runtime:    runtime.Float64 * float64(runtime_mut),
		}

		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			httphelpers.WriteError(
				w,
				http.StatusInternalServerError,
				fmt.Sprintf("error marshaling payload: %v", err),
			)
			return
		}

		result, err := submission.SendToJudge0(judge0URL, params, payloadJSON)
		if err != nil {
			logger.Errof("Error sending request to Judge0: %v", err)
			httphelpers.WriteError(
				w,
				http.StatusInternalServerError,
				fmt.Sprintf("Error sending request to Judge0: %v", err),
			)
			return
		}
		defer result.Body.Close()

		var data submission.Judgeresp
		if err = json.NewDecoder(result.Body).Decode(&data); err != nil {
			logger.Errof("Error decoding response from Judge0: %v", err)
			httphelpers.WriteError(
				w,
				http.StatusInternalServerError,
				fmt.Sprintf("Error decoding response from Judge0: %v", err),
			)
			return
		}

		data.StdOut, _ = submission.DecodeB64(data.StdOut)
		data.CompilerOutput, _ = submission.DecodeB64(data.CompilerOutput)
		data.ExpectedOutput = testcase.ExpectedOutput
		data.Message, _ = submission.DecodeB64(data.Message)
		data.StdErr, _ = submission.DecodeB64(data.StdErr)
		data.TestCaseID = testcase.ID.String()
		data.Input = testcase.Input
		response.Result[i] = data
		status, _ := data.Status.ID.Int64()
		if status == 3 {
			response.TestCasesPassed += 1
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errof("Error encoding response: %v", err)
		httphelpers.WriteError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Error encoding response: %v", err),
		)
	}
}
