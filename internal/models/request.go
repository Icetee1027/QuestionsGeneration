package models

type GenerateRequest struct {
	Subject      string `json:"subject"`
	Difficulty   string `json:"difficulty"`
	QuestionType string `json:"question_type"`
}
