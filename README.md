# <img width="32" height="32" alt="appicon" src="https://github.com/user-attachments/assets/fe60060a-6c0c-4435-aa41-2af811cf03f9" /> TranslateGemma Desktop App

A modern Windows Desktop Application for translating text using a local Ollama model (`translategemma:latest`).

<img width="550"  alt="Screenshot 2026-02-25 195701" src="https://github.com/user-attachments/assets/3305b3cc-6163-42fc-ab55-a77a2e82013b" />

## Features

-   **Source & Target Language Selection**: Choose from a wide range of supported languages.
-   **Local Processing**: Uses your local Ollama instance for privacy and speed.
-   **Simple UI**: Clean interface for inputting text and viewing translations.

## Prerequisites

1.  **Ollama**: You must have Ollama installed and running.
    -   Download from [ollama.com](https://ollama.com).
2.  **Model**: You need the `translategemma:latest` or other model variants available in your Ollama instance.
    -   To check available models: `ollama list`.
    -   If `translategemma:latest` or other variants are missing, pull via `ollama pull translategemma:latest`. This process takes a while.

## How to Run

1.  Ensure Ollama is running (`ollama serve` or via system tray) via port `11434`.
2.  Double-click `TranslateGemma.exe`.
3.  Select Source and Target languages.
4.  Enter text in the "Input Text" box.
5.  Click "Translate".
6.  Depending on whether the model is already loaded, you'll see the translation streamed in the output box.

## Platform Support

| Platform | Status |
|----------|--------|
| Windows (`.exe`) | Verified |
| macOS | Not verified |
| Linux | Not verified |

## FAQ / Troubleshooting

<summary><b>Why doesn't it translate immediately?</b></summary>

The first translation after starting the app (or after a period of inactivity) requires Ollama to load the model into memory, which can take several seconds or longer depending on your hardware. Subsequent translations are faster because the model stays loaded. You can verify the model is ready with `ollama list` or by checking the Ollama system tray icon.

</details>
<details>
<summary><b>Why does Windows SmartScreen pop up?</b></summary>

The `.exe` is not code-signed, so Windows SmartScreen flags it as an unrecognized app. This is expected for unsigned binaries distributed outside the Microsoft Store. To proceed, click **More info** → **Run anyway**.

</details>

<details>

## Building from Source

Requirements:
-   Go 1.20+
-   Wait for `rsrc` tool (optional, for manifest embedding).

Steps:
1.  Clone repository.
2.  Run `go mod tidy`.
3.  (Optional) Generate resource file: `rsrc -manifest main.manifest -o rsrc.syso`.
4.  Build: `go build -ldflags="-H windowsgui" -o TranslateGemma.exe`.
