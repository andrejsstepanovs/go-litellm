package audio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/andrejsstepanovs/go-litellm/request"
)

func TranscribeAudio(url, token, filePath, model string, extraBody map[string]any) (*http.Response, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add model field
	err := writer.WriteField("model", model)
	if err != nil {
		return nil, fmt.Errorf("error writing model field: %w", err)
	}

	// Add extra body fields
	for key, value := range extraBody {
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case bool:
			strValue = strconv.FormatBool(v)
		case float64:
			strValue = strconv.FormatFloat(v, 'f', -1, 64)
		case int:
			strValue = strconv.Itoa(v)
		default:
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("error marshaling extra body field %s: %w", key, err)
			}
			strValue = string(jsonBytes)
		}
		err := writer.WriteField(key, strValue)
		if err != nil {
			return nil, fmt.Errorf("error writing extra body field %s: %w", key, err)
		}
	}

	// Add file field
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("error copying file content: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing writer: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	client := &http.Client{}
	return client.Do(req)
}

// Speech generates audio from text using OpenAI-compatible TTS API
func Speech(url, token string, speechRequest request.Speech) (*http.Response, error) {
	requestBody, err := json.Marshal(speechRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling speech request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}
