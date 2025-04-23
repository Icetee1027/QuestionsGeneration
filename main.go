package main

import (
	"log"
	"os"

	"recipe-generator/internal/handlers"
	"recipe-generator/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 載入環境變數
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 初始化 AI 服務
	aiService, err := services.NewAIService()
	if err != nil {
		log.Fatalf("Failed to initialize AI service: %v", err)
	}

	// 初始化處理器
	questionHandler := handlers.NewQuestionHandler(aiService)

	// 設置路由
	router := gin.Default()
	router.POST("/generate-question", questionHandler.GenerateQuestion)

	// 啟動伺服器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
