package services

import (
	"encoding/json"
	"fmt"
	c "github.com/evgenyshipko/go-rag-chat-helper/internal/const"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/logger"
	"github.com/go-resty/resty/v2"
	"os"
	"strings"
)

type LlmService struct{}

func NewLlmService() *LlmService {
	return &LlmService{}
}

func (s *LlmService) GetKeywords(text string) ([]string, error) {
	messages := make([]GptMessage, 0)

	systemMessage := GptMessage{Text: "Ты - чат-бот, который отвечает пользователю на вопросы о компании, основываясь только на документах компании", Role: c.System}
	messageText := fmt.Sprintf("Контекст: пользователь задает вопрос, связанный с дейтельностью компании. Вопрос пользователя: %s .Переформулируй вопрос с учетом контекста, и вытащи из него ключевые слова, которые отражают тему запроса. Дообогати похожими ключевыми словами. В ответе укажи только ключевые слова, через запятую", text)
	message := GptMessage{Text: messageText, Role: c.User}

	messages = append(messages, systemMessage)
	messages = append(messages, message)

	res, err := sendGptRequest(GptRequestBody{Messages: messages, ChatId: 2660})
	if err != nil {
		logger.Instance.Warnw("send gpt request failed", "error", err)
		return []string{}, err
	}

	logger.Instance.Infow("send gpt response", "response", res)

	var data GptResponse

	err = json.Unmarshal(res.Body(), &data)
	if err != nil {
		logger.Instance.Warnw("parse gpt response failed", "error", err)
		return []string{}, err
	}

	logger.Instance.Infow("GetKeywords", "data", data.Text)

	return strings.Split(data.Text, ","), nil
}

func (s *LlmService) GetAnswerBasedOnDocument(query string, document string) (string, error) {
	messages := make([]GptMessage, 0)

	systemMessage := GptMessage{Text: "Ты - чат-бот, который отвечает пользователю на вопросы о компании, основываясь только на документах компании", Role: c.System}
	messageText := fmt.Sprintf("Контекст: пользователь задает вопрос, связанный с дейтельностью компании. Ответь на вопрос используя ТОЛЬКО документацию. Вопрос пользователя: %s. Документация: %s. ", query, document)
	message := GptMessage{Text: messageText, Role: c.User}

	messages = append(messages, systemMessage)
	messages = append(messages, message)

	res, err := sendGptRequest(GptRequestBody{Messages: messages, ChatId: 2660})
	if err != nil {
		logger.Instance.Warnw("send gpt request failed", "error", err)
		return "", err
	}

	logger.Instance.Infow("send gpt response", "response", res)

	var data GptResponse

	err = json.Unmarshal(res.Body(), &data)
	if err != nil {
		logger.Instance.Warnw("parse gpt response failed", "error", err)
		return "", err
	}

	logger.Instance.Infow("parse gpt response", "data", data)

	return data.Text, nil
}

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

func sendGptRequest(body GptRequestBody) (*resty.Response, error) {
	url := os.Getenv("LLM_REQUEST_LINK")

	var headers map[string]string

	logger.Instance.Info(os.Getenv("LLM_REQUEST_HEADERS"))

	err := json.Unmarshal([]byte(os.Getenv("LLM_REQUEST_HEADERS")), &headers)
	if err != nil {
		logger.Instance.Warnw("Ошибка доставания заголовков из переменной LLM_REQUEST_HEADERS", err)
		return nil, err
	}

	logger.Instance.Debugw("send request", "body", body, "headers", headers)

	return resty.New().R().
		SetBody(body).
		SetHeaders(headers).
		Post(url)
}
