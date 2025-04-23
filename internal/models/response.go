package models

// 基礎題目結構
type BaseQuestion struct {
	QuestionType  string            `json:"question_type"`
	Question      string            `json:"question"`
	Options       map[string]string `json:"options,omitempty"`
	CorrectAnswer []string          `json:"correct_answer,omitempty"`
	Explanation   string            `json:"explanation"`
}

// 配對題結構
type MatchingQuestion struct {
	QuestionType  string            `json:"question_type"`
	Question      string            `json:"question"`
	Pairs         map[string]string `json:"pairs"`
	CorrectAnswer []string          `json:"correct_answer,omitempty"`
	Explanation   string            `json:"explanation"`
}

// 閱讀題組結構
type ReadingQuestion struct {
	QuestionType  string         `json:"question_type"`
	Passage       string         `json:"passage"`
	Questions     []BaseQuestion `json:"questions"`
	CorrectAnswer []string       `json:"correct_answer,omitempty"`
	Explanation   string         `json:"explanation"`
}

// 題目回應介面
type QuestionResponse interface {
	GetQuestionType() string
}

// 實作 QuestionResponse 介面
func (q BaseQuestion) GetQuestionType() string {
	return q.QuestionType
}

func (q MatchingQuestion) GetQuestionType() string {
	return q.QuestionType
}

func (q ReadingQuestion) GetQuestionType() string {
	return q.QuestionType
}
