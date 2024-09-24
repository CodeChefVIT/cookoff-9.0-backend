package submission

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var bearer = os.Getenv("JUDGE0_TOKEN")

type Submission struct {
	LanguageID int     `json:"language_id"`
	SourceCode string  `json:"source_code"`
	Input      string  `json:"stdin"`
	Output     string  `json:"expected_output"`
	Runtime    float64 `json:"cpu_time_limit"`
	Callback   string  `json:"callback_url"`
}

type Judgeresp struct {
	TestCaseID     string `json:"testcase_id"`
	StdOut         string `json:"stdout"`
	ExpectedOutput string `json:"expected_output"`
	Input          string `json:"input"`
	Time           string `json:"time"`
	Memory         int    `json:"memory"`
	StdErr         string `json:"stderr"`
	Token          string `json:"token"`
	Message        string `json:"message"`
	Status         Status `json:"status"`
	CompilerOutput string `json:"compile_output"`
}

type Status struct {
	ID          json.Number `json:"id"`
	Description string      `json:"description"`
}

type resultresp struct {
	ID          string `json:"submission_id"`
	QuestionID  string `json:"question_id"`
	Passed      int    `json:"testcases_passed"`
	Failed      int    `json:"testcases_failed"`
	Description string `json:"description"`
}

func B64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func DecodeB64(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func RuntimeMut(language_id int) (int, error) {
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
		return 0, fmt.Errorf("invalid language ID: %d", language_id)
	}
	return runtime_mut, nil
}

func SendToJudge0(judge0Url *url.URL, params url.Values, payload []byte) (*http.Response, error) {
	judge0Url.RawQuery = params.Encode()
	judgereq, err := http.NewRequest("POST", judge0Url.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request to Judge0: %v", err)
	}

	judgereq.Header.Add("Content-Type", "application/json")
	judgereq.Header.Add("Accept", "application/json")
	judgereq.Header.Add("Authorization", fmt.Sprintf("Bearer %v", bearer))

	resp, err := http.DefaultClient.Do(judgereq)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Judge0: %v", err)
	}

	return resp, nil
}
