package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type CodeSubmission struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type Judge0Request struct {
	SourceCode string `json:"source_code"`
	LanguageID string `json:"language_id"`
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

var languageMap = map[string]string{
	"python3":         "71", // Python 3
	"cpp":             "54", // C++
	"javascript":      "63", // JavaScript (Node.js)
	"java":            "62", // Java
	"C++ (GCC 9.2.0)": "54",
	// Add more languages as needed
}

func codeResponse(resp *http.Response, w http.ResponseWriter, client *http.Client) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	var judge0Resp Judge0Response
	err = json.Unmarshal(body, &judge0Resp)
	if err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}

	// Poll for result
	for {
		time.Sleep(2 * time.Second)

		resultReq, err := http.NewRequest("GET", fmt.Sprintf("https://judge0-ce.p.rapidapi.com/submissions/%s", judge0Resp.Token), nil)
		if err != nil {
			http.Error(w, "Failed to create result request", http.StatusInternalServerError)
			return
		}
		resultReq.Header.Set("X-RapidAPI-Key", os.Getenv("RAPIDAPI_KEY"))
		resultReq.Header.Set("X-RapidAPI-Host", os.Getenv("API_HOST"))

		resultResp, err := client.Do(resultReq)
		if err != nil {
			http.Error(w, "Failed to fetch result", http.StatusInternalServerError)
			return
		}
		defer resultResp.Body.Close()

		resultBody, err := io.ReadAll(resultResp.Body)
		if err != nil {
			http.Error(w, "Failed to read result response", http.StatusInternalServerError)
			return
		}

		var judge0Result Judge0Result
		err = json.Unmarshal(resultBody, &judge0Result)
		if err != nil {
			http.Error(w, "Failed to decode result", http.StatusInternalServerError)
			return
		}

		// Check if result is ready
		if judge0Result.Status.Description == "Accepted" || judge0Result.Status.Description == "Compilation Error" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status": judge0Result.Status.Description,
				"output": judge0Result.Stdout,
				"error":  judge0Result.Stderr,
			})
			return
		}
	}
}

func SubmitCode(w http.ResponseWriter, r *http.Request) {
	var submission CodeSubmission
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Map the language name to its ID
	languageID, exists := languageMap[submission.Language]
	if !exists {
		http.Error(w, "Unsupported language", http.StatusBadRequest)
		return
	}

	// Prepare the Judge0 API request
	judge0Req := Judge0Request{
		SourceCode: submission.Code,
		LanguageID: languageID,
	}

	jsonData, err := json.Marshal(judge0Req)
	if err != nil {
		http.Error(w, "Failed to encode request", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", "https://judge0-ce.p.rapidapi.com/submissions", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RapidAPI-Key", os.Getenv("RAPIDAPI_KEY"))
	req.Header.Set("X-RapidAPI-Host", os.Getenv("API_HOST"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Call the codeResponse function with the response and client
	codeResponse(resp, w, client)
}
