package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const OllamaURL = "http://localhost:11434/api/generate"
const DefaultModel = "translategemma:latest"

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// App holds application state and context
type App struct {
	ctx         context.Context
	activeModel string
}

func NewApp() *App {
	return &App{activeModel: DefaultModel}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetModel returns the currently selected model name
func (a *App) GetModel() string {
	return a.activeModel
}

// SetModel updates the active model
func (a *App) SetModel(model string) {
	a.activeModel = model
}

// GetLanguages returns all supported languages as a JSON-serializable slice
func (a *App) GetLanguages() []Language {
	return SupportedLanguages
}

// TranslateStream performs a streaming translation and emits events for each token.
// It emits "translation:chunk" events with partial text and a final "translation:done" event.
func (a *App) TranslateStream(sourceLang, sourceCode, targetLang, targetCode, text string) error {
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
		Model:  a.activeModel,
		Prompt: prompt,
		Stream: true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		runtime.EventsEmit(a.ctx, "translation:error", err.Error())
		return err
	}

	resp, err := http.Post(OllamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		msg := fmt.Sprintf("Failed to connect to Ollama: %v", err)
		runtime.EventsEmit(a.ctx, "translation:error", msg)
		return fmt.Errorf(msg)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		msg := fmt.Sprintf("Ollama returned status: %s, body: %s", resp.Status, string(bodyBytes))
		runtime.EventsEmit(a.ctx, "translation:error", msg)
		return fmt.Errorf(msg)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var chunk OllamaResponse
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue
		}
		if chunk.Response != "" {
			runtime.EventsEmit(a.ctx, "translation:chunk", chunk.Response)
		}
		if chunk.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		msg := fmt.Sprintf("Stream read error: %v", err)
		runtime.EventsEmit(a.ctx, "translation:error", msg)
		return fmt.Errorf(msg)
	}

	runtime.EventsEmit(a.ctx, "translation:done", "")
	return nil
}

// TranslateOnce performs a non-streaming translation (fallback)
func (a *App) TranslateOnce(sourceLang, sourceCode, targetLang, targetCode, text string) (string, error) {
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
		Model:  a.activeModel,
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
