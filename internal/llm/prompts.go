package llm

import (
	"encoding/json"
	"fmt"
	c "github.com/evgenyshipko/go-rag-chat-helper/internal/const"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/logger"
)

func GetKeywords(text string) (string, error) {

	messages := make([]GptMessage, 0)

	systemMessage := GptMessage{Text: "Ты - чат-бот, который помогает искать искать информацию в документах it-компании.", Role: c.System}
	messageText := fmt.Sprintf("Контекст: пользователь - это разработчик ПО, который задает вопрос боту-помошнику по документации компании. Вопрос пользователя: %s .Переформулируй вопрос с учетом контекста, и вытащи из него ключевые слова, которые отражают тему запроса. Дообогати похожими ключевыми словами. В ответе укажи только ключевые слова, через запятую", text)
	message := GptMessage{Text: messageText, Role: c.User}

	messages = append(messages, systemMessage)
	messages = append(messages, message)

	res, err := SendGptRequest(GptRequestBody{Messages: messages, ChatId: 2660})
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
