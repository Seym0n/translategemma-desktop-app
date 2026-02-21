# TranslateGemma Desktop App

A modern Windows Desktop Application for translating text using a local Ollama model (`translategemma:latest`).

## Features

-   **Source & Target Language Selection**: Choose from a wide range of supported languages.
-   **Local Processing**: Uses your local Ollama instance for privacy and speed.
-   **Simple UI**: Clean interface for inputting text and viewing translations.

## Prerequisites

1.  **Ollama**: You must have Ollama installed and running.
    -   Download from [ollama.com](https://ollama.com).
2.  **Model**: You need the `translategemma:latest` model available in your Ollama instance.
    -   If this is a custom model, ensure you have created/pulled it.
    -   To check available models: `ollama list`.
    -   If you need to pull a base model (e.g., `gemma`), use `ollama pull gemma`.

## How to Run

1.  Ensure Ollama is running (`ollama serve` or via system tray).
2.  Double-click `TranslateGemma.exe`.
3.  Select Source and Target languages.
4.  Enter text in the "Input Text" box.
5.  Click "Translate".

## Building from Source

Requirements:
-   Go 1.20+
-   Wait for `rsrc` tool (optional, for manifest embedding).

Steps:
1.  Clone repository.
2.  Run `go mod tidy`.
3.  (Optional) Generate resource file: `rsrc -manifest main.manifest -o rsrc.syso`.
4.  Build: `go build -ldflags="-H windowsgui" -o TranslateGemma.exe`.
