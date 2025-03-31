package main

import (
	"GoSQL/internal/ui/pageManager"
	"context"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	context := context.Background()
	app.EnableMouse(true)
	pageManager := pageManager.InitializePages(context , app)
	if err := app.SetRoot(pageManager.Main, true).SetFocus(pageManager.Main).Run(); err != nil {
		panic(err)
	}

}
