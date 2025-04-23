# 國中會考題目生成 API

這是一個使用 Go 和 Gemini AI 開發的國中會考題目生成 API 服務。根據指定的科目、難度和題型，生成相應的會考題目。

## 系統功能

1. **多樣化題型**
   - 支援七種題型：單選題、多選題、是非題、填空題、簡答題、配對題、閱讀題組
   - 根據不同題型提供相應的格式和結構

2. **彈性設定**
   - 支援五種科目：國文、英文、數學、自然、社會
   - 四種難度等級：簡單、普通、困難、極難

3. **標準化輸出**
   - 統一的 JSON 回應格式
   - 清晰的題目結構
   - 完整的選項和答案（適用於選擇題）
   - 詳細的題目詳解

## 使用 Docker 快速部署

### 1. 環境準備
```bash
# 複製環境設定檔
cp .env.example .env

# 編輯 .env 檔案，設定您的 Gemini API 金鑰
# GEMINI_API_KEY=您的API金鑰
```

### 2. 建立與運行
```bash
# 建立 Docker 映像檔
docker build -t question-generator .

# 運行容器（前台模式）
docker run -p 8080:8080 question-generator

# 或使用背景模式運行
docker run -d -p 8080:8080 question-generator
```

### 3. API 使用範例

發送 POST 請求到 `http://localhost:8080/generate-question`：

```json
{
    "subject": "英文",
    "difficulty": "困難",
    "question_type": "單選題"
}
```

您將收到包含題目的 JSON 回應，格式根據題型而有所不同：

#### 單選題回應範例
```json
{
    "question_type": "單選題",
    "question": "以下哪一個是正確的英文句子？",
    "options": {
        "A": "He go to school.",
        "B": "He goes to school.",
        "C": "He going to school.",
        "D": "He went to school."
    },
    "correct_answer": ["B"],
    "explanation": "正確答案是 B。在英文中，當主詞是第三人稱單數（he/she/it）時，動詞需要加上 -s 或 -es。選項 A 缺少 -s，選項 C 使用了進行式但缺少 be 動詞，選項 D 使用了過去式，但題目沒有指定過去時間。"
}
```

### 4. 常見問題排解

1. 如果無法連接 API：
   - 確認容器是否正常運行：`docker ps`
   - 檢查日誌：`docker logs <container_id>`

2. 如果需要停止服務：
   - 查看容器 ID：`docker ps`
   - 停止容器：`docker stop <container_id>`

3. 如果需要更新服務：
   - 停止舊容器
   - 重新建立映像檔
   - 啟動新容器
