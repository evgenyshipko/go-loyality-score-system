package services

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/const"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/logger"
	"github.com/evgenyshipko/go-rag-chat-helper/internal/storage"
	"io"
	"strings"
	"unicode/utf8"
)

const (
	ChunkSize = 2000
)

type DocumentService struct {
	storage *storage.SQLStorage
}

func NewDocumentService(storage *storage.SQLStorage) *DocumentService {
	return &DocumentService{storage: storage}
}

func (service *DocumentService) UploadDocument(buffer *bytes.Buffer) error {
	documentText, err := extractTextFromDocx(buffer.Bytes())
	if err != nil {
		logger.Instance.Warnw("extractTextFromDocx", "err", err.Error())
		return err
	}

	chunks := splitDocumentIntoChunks(documentText)

	err = service.storage.SaveChunks(chunks)
	if err != nil {
		logger.Instance.Warnw("storage.SaveChunks", "err", err.Error())
		return err
	}
	return nil
}

func (service *DocumentService) SearchDocument(keywords []string) (string, error) {
	chunks, err := service.storage.SearchChunks(keywords)
	if err != nil {
		return "", err
	}
	// TODO: это первая версия - берем все чанки и возвращаем их. Для больших документов надо изобретать более сложную логику
	res := ""
	for _, chunk := range chunks {
		res += chunk.Text
	}

	logger.Instance.Infow("SearchDocument", "res", res)

	return res, nil
}

func splitDocumentIntoChunks(document string) []constants.DocumentChunk {
	var chunks []constants.DocumentChunk
	words := strings.Fields(document)
	var currentChunk []string
	currentSize := 0
	chunkIndex := 0

	for _, word := range words {
		wordLength := utf8.RuneCountInString(word)
		if currentSize+wordLength+len(currentChunk) > ChunkSize {
			chunks = append(chunks, constants.DocumentChunk{
				Text:       strings.Join(currentChunk, " "),
				ChunkIndex: chunkIndex,
			})
			chunkIndex++
			currentChunk = []string{word}
			currentSize = wordLength
		} else {
			currentChunk = append(currentChunk, word)
			currentSize += wordLength
		}
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, constants.DocumentChunk{
			Text:       strings.Join(currentChunk, " "),
			ChunkIndex: chunkIndex,
		})
	}

	return chunks
}

func extractTextFromDocx(docxData []byte) (string, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(docxData), int64(len(docxData)))
	if err != nil {
		return "", fmt.Errorf("не удалось открыть docx как zip: %v", err)
	}

	var documentXMLFile *zip.File
	for _, f := range zipReader.File {
		if f.Name == "word/document.xml" {
			documentXMLFile = f
			break
		}
	}

	if documentXMLFile == nil {
		return "", fmt.Errorf("файл word/document.xml не найден в архиве")
	}

	rc, err := documentXMLFile.Open()
	if err != nil {
		return "", fmt.Errorf("не удалось открыть document.xml: %v", err)
	}
	defer rc.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, rc)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении document.xml: %v", err)
	}

	return parseDocumentXML(buf.Bytes())
}

func parseDocumentXML(data []byte) (string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var documentText strings.Builder

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("ошибка при разборе XML: %v", err)
		}

		switch se := tok.(type) {
		case xml.CharData:
			documentText.Write(se)
		}
	}

	return documentText.String(), nil
}
