package components

import (
	"GoSQL/internal/config"
	"GoSQL/internal/ui/router"
	"context"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateProfileFooter(ctx context.Context, mainGrid *tview.Grid) *tview.Grid {
	buttonGrid := tview.NewGrid().
		SetRows(1)

	deleteButton := tview.NewButton("Delete").SetSelectedFunc(func() {
		println("Delete action triggered")
	})
	updateButton := tview.NewButton("Update").SetSelectedFunc(func() {
		println("Update action triggered")
	})
	newButton := tview.NewButton("New").SetSelectedFunc(func() {
		router.NavigatePage(config.CreateProfilePage, -1, ctx)
	})
	uiConfig, ok := ctx.Value("ui-config").(*config.UIConfig)
	if !ok {
		log.Fatalln("Error Failed To Get App ")
	}
	buttonGrid.
		AddItem(newButton, 0, 0, 1, 1, 1, 1, true).
		AddItem(deleteButton, 0, 1, 1, 1, 1, 1, true).
		AddItem(updateButton, 0, 2, 1, 1, 1, 1, true)

	// Set keyboard shortcuts in mainGrid
	mainGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'n':
			uiConfig.App.SetFocus(newButton)
		case 'l':
			uiConfig.App.SetFocus(updateButton)
		case 'm':
			uiConfig.App.SetFocus(deleteButton)
		}
		return event
	})
	return buttonGrid

}
