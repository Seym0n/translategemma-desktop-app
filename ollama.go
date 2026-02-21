package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const OllamaURL = "http://localhost:11434/api/generate"
const ModelName = "translategemma:latest"

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func Translate(sourceLang, sourceCode, targetLang, targetCode, text string) (string, error) {
	promptTemplate := `You are a professional %s (%s) to %s (%s) translator. Your goal is to accurately convey the meaning and nuances of the original %s text while adhering to %s grammar, vocabulary, and cultural sensitivities. 
Produce only the %s translation, without any additional explanations or commentary. Please translate the following %s text into %s: 

%s`

	prompt := fmt.Sprintf(promptTemplate,
		sourceLang, sourceCode,
		targetLang, targetCode,
		sourceLang,
		targetLang,
		targetLang,
		sourceLang,
		targetLang,
		text)

	reqBody := OllamaRequest{
		Model:  ModelName,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(OllamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned status: %s, body: %s", resp.Status, string(bodyBytes))
	}

	var response OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return strings.TrimSpace(response.Response), nil
}
