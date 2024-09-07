package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL             = "https://trauma.codechefvit.com"
	submissionsEndpoint = "/submissions"
	timeout             = 30 * time.Second
	pollInterval        = 2 * time.Second
)

var languageMap = map[string]string{
	"Python (3.8.1)": "71", // Python 3
	"cpp":            "54", // C++
	"javascript":     "63", // JavaScript (Node.js)
	"java":           "62", // Java
	// Add more languages as needed
}

type CodeSubmission struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type Judge0Request struct {
	SourceCode   string `json:"source_code"`
	LanguageID   string `json:"language_id"`
	Base64Encode bool   `json:"base64_encoded"` // Indicates if the code is base64-encoded
}

type Judge0Response struct {
	Token string `json:"token"`
}

type Judge0Result struct {
	Status struct {
		Description string `json:"description"`
	} `json:"status"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func codeResponse(token string, w http.ResponseWriter, client *http.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			http.Error(w, "Request timed out", http.StatusRequestTimeout)
			return
		default:
			time.Sleep(pollInterval)

			resultURL := fmt.Sprintf("%s%s/%s?base64_encoded=false,fields=stdout,stderr", baseURL, submissionsEndpoint, token)
			resultResp, err := doRequest("GET", resultURL, nil, client)
			if err != nil {
				http.Error(w, "Failed to fetch result", http.StatusInternalServerError)
				return
			}

			var judge0Result Judge0Result
			if err := json.Unmarshal(resultResp, &judge0Result); err != nil {
				http.Error(w, "Failed to decode result", http.StatusInternalServerError)
				return
			}

			if isFinalStatus(judge0Result.Status.Description) {
				respondWithResult(w, judge0Result)
				return
			}
		}
	}
}

func SubmitCode(w http.ResponseWriter, r *http.Request) {
	var submission CodeSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	languageID, exists := languageMap[submission.Language]
	if !exists {
		http.Error(w, "Unsupported language", http.StatusBadRequest)
		return
	}

	encodedCode := base64.StdEncoding.EncodeToString([]byte(submission.Code))
	judge0Req := Judge0Request{
		SourceCode:   encodedCode,
		LanguageID:   languageID,
		Base64Encode: true,
	}

	jsonData, err := json.Marshal(judge0Req)
	if err != nil {
		http.Error(w, "Failed to encode request", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s?base64_encoded=true", baseURL, submissionsEndpoint), bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	body, err := doRequest("POST", req.URL.String(), bytes.NewBuffer(jsonData), client)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}

	var judge0Resp Judge0Response
	if err := json.Unmarshal(body, &judge0Resp); err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}

	codeResponse(judge0Resp.Token, w, client)
}

func doRequest(method, url string, body io.Reader, client *http.Client) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func isFinalStatus(description string) bool {
	switch description {
	case "Accepted", "Compilation Error", "Runtime Error (NZEC)":
		return true
	}
	return false
}

func respondWithResult(w http.ResponseWriter, result Judge0Result) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": result.Status.Description,
		"output": result.Stdout,
		"error":  result.Stderr,
	})
}
