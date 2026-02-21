package main

import (
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var mw *walk.MainWindow
	var sourceCB, targetCB *walk.ComboBox
	var inputTE, outputTE *walk.TextEdit
	var translateBtn *walk.PushButton

	langNames := getLanguageNames()
	
	// Find default indices
	sourceIndex := 0
	targetIndex := 0
	for i, l := range SupportedLanguages {
		if l.Code == "en" {
			sourceIndex = i
		}
		if l.Code == "es" {
			targetIndex = i
		}
	}

	if _, err := (MainWindow{
		AssignTo: &mw,
		Title:    "TranslateGemma Desktop",
		MinSize:  Size{Width: 800, Height: 600},
		Layout:   VBox{},
		Children: []Widget{
			Label{
				Text: "TranslateGemma Desktop App",
				Font: Font{PointSize: 16, Bold: true},
				TextAlignment: AlignCenter,
			},
			Composite{
				Layout: Grid{Columns: 4},
				Children: []Widget{
					Label{Text: "Source Language:", TextAlignment: AlignFar},
					ComboBox{
						AssignTo:     &sourceCB,
						Model:        langNames,
						CurrentIndex: sourceIndex,
					},
					Label{Text: "Target Language:", TextAlignment: AlignFar},
					ComboBox{
						AssignTo:     &targetCB,
						Model:        langNames,
						CurrentIndex: targetIndex,
					},
				},
			},
			Label{Text: "Input Text:"},
			TextEdit{
				AssignTo: &inputTE,
				VScroll:  true,
				MinSize:  Size{Height: 150},
			},
			PushButton{
				AssignTo: &translateBtn,
				Text:     "Translate",
				Font:     Font{PointSize: 10, Bold: true},
				OnClicked: func() {
					srcIdx := sourceCB.CurrentIndex()
					tgtIdx := targetCB.CurrentIndex()
					
					if srcIdx < 0 || tgtIdx < 0 {
						walk.MsgBox(mw, "Error", "Please select source and target languages.", walk.MsgBoxIconError)
						return
					}

					text := inputTE.Text()
					if text == "" {
						walk.MsgBox(mw, "Warning", "Please enter text to translate.", walk.MsgBoxIconWarning)
						return
					}

					sourceLang := SupportedLanguages[srcIdx]
					targetLang := SupportedLanguages[tgtIdx]

					translateBtn.SetEnabled(false)
					translateBtn.SetText("Translating...")
					outputTE.SetText("Translating...")

					go func() {
						translated, err := Translate(sourceLang.Name, sourceLang.Code, targetLang.Name, targetLang.Code, text)
						
						// Update UI on main thread
						mw.Synchronize(func() {
							translateBtn.SetEnabled(true)
							translateBtn.SetText("Translate")
							
							if err != nil {
								walk.MsgBox(mw, "Error", fmt.Sprintf("Translation failed: %v", err), walk.MsgBoxIconError)
								outputTE.SetText("")
							} else {
								outputTE.SetText(translated)
							}
						})
					}()
				},
			},
			Label{Text: "Output Text:"},
			TextEdit{
				AssignTo: &outputTE,
				ReadOnly: true,
				VScroll:  true,
				MinSize:  Size{Height: 150},
			},
		},
	}).Run(); err != nil {
		fmt.Println(err)
	}
}

func getLanguageNames() []string {
	names := make([]string, len(SupportedLanguages))
	for i, l := range SupportedLanguages {
		names[i] = fmt.Sprintf("%s (%s)", l.Name, l.Code)
	}
	return names
}
