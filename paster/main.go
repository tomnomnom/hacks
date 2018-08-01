package main

import (
	"github.com/andlabs/ui"
	"github.com/atotto/clipboard"
)

func main() {
	snippets := [][]string{
		{"Upside Down", "🙃"},
		{"Eye Roll", "🙄"},
		{"Shrug", `¯\_(ツ)_/¯`},
	}

	err := ui.Main(func() {

		box := ui.NewVerticalBox()

		for _, s := range snippets {
			label, text := s[0], s[1]
			button := ui.NewButton(label)
			button.OnClicked(func(*ui.Button) {
				clipboard.WriteAll(text)
			})

			box.Append(button, false)
		}

		window := ui.NewWindow("Paster", 100, 100, false)
		window.SetMargined(true)
		window.SetChild(box)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})

		window.Show()
	})

	if err != nil {
		panic(err)
	}
}
