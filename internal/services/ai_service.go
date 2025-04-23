package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"recipe-generator/internal/models"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewAIService() (*AIService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, fmt.Errorf("failed to create AI client: %v", err)
	}

	model := client.GenerativeModel("models/gemini-2.0-flash-001")
	return &AIService{
		client: client,
		model:  model,
	}, nil
}

func (s *AIService) GenerateQuestion(request *models.GenerateRequest) (models.QuestionResponse, error) {
	ctx := context.Background()

	// 構建提示詞
	prompt := buildPrompt(request)

	// 生成回應
	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	// 取得回傳的 AI 原始內容
	raw, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, fmt.Errorf("AI 回傳格式不是 genai.Text")
	}

	rawString := string(raw)
	fmt.Println("🧠 AI 回傳內容：", rawString)

	// 移除 Markdown 格式的反引號
	cleanJson := strings.TrimPrefix(rawString, "```json\n")
	cleanJson = strings.TrimSuffix(cleanJson, "```")

	// 根據題型解析 JSON
	switch request.QuestionType {
	case "單選題", "多選題", "是非題":
		var question models.BaseQuestion
		if err := json.Unmarshal([]byte(cleanJson), &question); err != nil {
			return nil, fmt.Errorf("failed to parse question: %v", err)
		}
		return question, nil

	case "配對題":
		var question models.MatchingQuestion
		if err := json.Unmarshal([]byte(cleanJson), &question); err != nil {
			return nil, fmt.Errorf("failed to parse matching question: %v", err)
		}
		return question, nil

	case "閱讀題組":
		var question models.ReadingQuestion
		if err := json.Unmarshal([]byte(cleanJson), &question); err != nil {
			return nil, fmt.Errorf("failed to parse reading question: %v", err)
		}
		return question, nil

	case "填空題", "簡答題":
		var question models.BaseQuestion
		if err := json.Unmarshal([]byte(cleanJson), &question); err != nil {
			return nil, fmt.Errorf("failed to parse question: %v", err)
		}
		return question, nil

	default:
		return nil, fmt.Errorf("unsupported question type: %s", request.QuestionType)
	}
}

func buildPrompt(request *models.GenerateRequest) string {
	var promptTemplate string

	// 數學符號格式提醒
	mathSymbolGuide := ""
	if request.Subject == "數學" {
		mathSymbolGuide = `
        如果題目包含數學符號，請使用以下格式：
        - 分數：使用 "a/b" 格式，例如 "1/2"
        - 根號：使用 "sqrt(n)" 格式，例如 "sqrt(2)"
        - 指數：使用 "^" 符號，例如 "x^2"
        - 希臘字母：使用英文名稱，例如 "pi" 代替 "π"
        - 函數：使用英文名稱，例如 "cos" 代替 "cos"
        請避免使用 LaTeX 格式或其他特殊符號。`
	}

	switch request.QuestionType {
	case "單選題":
		promptTemplate = `請以%s為主題，出一個%s難度的單選題。%s
        請使用繁體中文撰寫題目和選項。
        請包含四個選項和正確答案。
        回應格式必須完全符合以下 JSON 結構：
        {
            "question_type": "單選題",
            "question": "題目內容",
            "options": {
                "A": "選項A",
                "B": "選項B",
                "C": "選項C",
                "D": "選項D"
            },
            "correct_answer": ["正確選項"],
            "explanation": "簡單解釋為什麼這個答案是正確的，以及為什麼其他選項是錯誤的"
        }`

	case "多選題":
		promptTemplate = `請以%s為主題，出一個%s難度的多選題。%s
        請使用繁體中文撰寫題目和選項。
        請包含四個選項和正確答案（可能有多個）。
        回應格式必須完全符合以下 JSON 結構：
        {
            "question_type": "多選題",
            "question": "題目內容",
            "options": {
                "A": "選項A",
                "B": "選項B",
                "C": "選項C",
                "D": "選項D"
            },
            "correct_answer": ["正確選項1", "正確選項2"],
            "explanation": "簡單解釋為什麼這些答案是正確的，以及為什麼其他選項是錯誤的"
        }`

	case "是非題":
		promptTemplate = `請以%s為主題，出一個%s難度的是非題。%s
        請使用繁體中文撰寫題目和選項。
        請包含正確答案。
        回應格式必須完全符合以下 JSON 結構：
        {
            "question_type": "是非題",
            "question": "題目內容",
            "options": {
                "T": "是",
                "F": "否"
            },
            "correct_answer": ["正確選項"],
            "explanation": "簡單解釋為什麼這個答案是正確的，以及為什麼另一個選項是錯誤的"
        }`

	case "填空題":
		promptTemplate = `請以%s為主題，出一個%s難度的填空題。%s
        請使用繁體中文撰寫題目和選項。
        請在題目中標示填空位置（使用______）。
        回應格式必須完全符合以下 JSON 結構：
        {
            "question_type": "填空題",
            "question": "題目內容",
            "options": null,
            "correct_answer": ["第一格答案",...,"第n格答案"],
            "explanation": "簡單解釋正確答案的內容和原因"
        }`

	case "簡答題":
		promptTemplate = `請以%s為主題，出一個%s難度的簡答題。%s
        請使用繁體中文撰寫題目和選項。
        回應格式必須完全符合以下 JSON 結構：
        {
            "question_type": "簡答題",
            "question": "題目內容",
            "options": null,
            "correct_answer": ["精簡簡答答案+請看詳解提醒"],
            "explanation": "簡單解釋正確答案的內容和原因"
        }`

	case "配對題":
		promptTemplate = `請以%s為主題，出一個%s難度的配對題。%s
        請使用繁體中文撰寫題目和選項。
        請提供三組配對項目。
        回應格式必須完全符合以下 JSON 結構：
        {
            "question_type": "配對題",
            "question": "請將左側項目與右側項目配對。",
            "pairs": {
                "左側項目1": "右側項目1",
                "左側項目2": "右側項目2",
                "左側項目3": "右側項目3"
            },
            "correct_answer": ["左側項目1", "左側項目2", "左側項目3","右側項目1", "右側項目2", "右側項目3"],
            "explanation": "簡單解釋每個配對的正確關係和原因"
        }`

	case "閱讀題組":
		promptTemplate = `請以%s為主題，出一個%s難度的閱讀題組。%s
        請使用繁體中文撰寫題目和選項。
        請包含一篇短文和兩個相關問題。
        回應格式必須完全符合以下 JSON 結構：
        {
            "question_type": "閱讀題組",
            "passage": "短文內容",
            "questions": [
                {
                    "question": "問題1",
                    "options": {
                        "A": "選項A",
                        "B": "選項B",
                        "C": "選項C",
                        "D": "選項D"
                    }
                },
                {
                    "question": "問題2",
                    "options": {
                        "A": "選項A",
                        "B": "選項B",
                        "C": "選項C",
                        "D": "選項D"
                    }
                }
            ],
            "correct_answer": ["第一題正確選項", "第二題正確選項"],
            "explanation": "簡單解釋每個問題的正確答案和原因，以及如何從文章中找出答案"
        }`
	}

	return fmt.Sprintf(promptTemplate, request.Subject, request.Difficulty, mathSymbolGuide)
}
