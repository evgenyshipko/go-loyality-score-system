package llm

import (
	"encoding/json"
	c "github.com/evgenyshipko/go-rag-chat-helper/internal/const"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/logger"
	"github.com/go-resty/resty/v2"
	"os"
)

type GptMessage struct {
	Text string     `json:"text"`
	Role c.LlmRoles `json:"role"`
	File *os.File   `json:"file,omitempty"`
}

type GptRequestBody struct {
	Messages []GptMessage `json:"messages"`
	ChatId   int          `json:"chatId"`
}

type GptResponse struct {
	Text string `json:"text"`
}

func SendGptRequest(body GptRequestBody) (*resty.Response, error) {
	url := os.Getenv("LLM_REQUEST_LINK")

	var headers map[string]string

	err := json.Unmarshal([]byte(os.Getenv("LLM_REQUEST_HEADERS")), &headers)
	if err != nil {
		logger.Instance.Warnw("Ошибка доставания заголловков из переменной LLM_REQUEST_HEADERS", err)
		return nil, err
	}

	logger.Instance.Debugw("send request", "body", body, "headers", headers)

	return resty.New().R().
		SetBody(body).
		SetHeaders(headers).
		Post(url)
}
