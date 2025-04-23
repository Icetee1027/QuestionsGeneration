package handlers

import (
	"net/http"

	"recipe-generator/internal/models"
	"recipe-generator/internal/services"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	aiService *services.AIService
}

func NewQuestionHandler(aiService *services.AIService) *QuestionHandler {
	return &QuestionHandler{
		aiService: aiService,
	}
}

func (h *QuestionHandler) GenerateQuestion(c *gin.Context) {
	var request models.GenerateRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求格式",
		})
		return
	}

	// 驗證必填欄位
	if request.Subject == "" || request.Difficulty == "" || request.QuestionType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少必填欄位",
		})
		return
	}

	// 驗證科目
	validSubjects := map[string]bool{
		"國文": true,
		"英文": true,
		"數學": true,
		"自然": true,
		"社會": true,
	}
	if !validSubjects[request.Subject] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的科目",
		})
		return
	}

	// 驗證難度
	validDifficulties := map[string]bool{
		"簡單": true,
		"普通": true,
		"困難": true,
		"極難": true,
	}
	if !validDifficulties[request.Difficulty] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的難度",
		})
		return
	}

	// 驗證題型
	validQuestionTypes := map[string]bool{
		"單選題":  true,
		"多選題":  true,
		"是非題":  true,
		"填空題":  true,
		"簡答題":  true,
		"配對題":  true,
		"閱讀題組": true,
	}
	if !validQuestionTypes[request.QuestionType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的題型",
		})
		return
	}

	question, err := h.aiService.GenerateQuestion(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "題目生成失敗",
		})
		return
	}

	c.JSON(http.StatusOK, question)
}
