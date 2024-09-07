package models

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
