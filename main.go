package main

import (
	"GoSQL/internal/page"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pageManager := page.InitPages()
	if err := app.SetRoot(pageManager.Main, true).SetFocus(pageManager.Main).Run(); err != nil {
		panic(err)
	}

}
