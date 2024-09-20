package submission

import (
	"encoding/base64"
	"encoding/json"
)

type Submission struct {
	LanguageID int    `json:"language_id"`
	SourceCode string `json:"source_code"`
	Input      string `json:"stdin"`
	Output     string `json:"expected_output"`
	Callback   string `json:"callback_url"`
}

type Judgeresp struct {
	TestCaseID     string
	StdOut         string `json:"stdout"`
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
