package components

import (
	"GoSQL/internal/config"
	"GoSQL/internal/ui/router"
	"context"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ProfileFooterUI struct {
	MainGrid     *tview.Grid
	UpdateButton *tview.Button
	NewButton    *tview.Button
}

func CreateProfileFooter(ctx context.Context, mainGrid *tview.Grid) *ProfileFooterUI {
	buttonGrid := tview.NewGrid().
		SetRows(1)

	updateButton := tview.NewButton("Update").SetSelectedFunc(func() {

	})
	newButton := tview.NewButton("New").SetSelectedFunc(func() {
		router.NavigatePage(config.CreateProfilePage, -1, ctx, "")
	})
	uiConfig, ok := ctx.Value("ui-config").(*config.UIConfig)
	if !ok {
		log.Fatalln("Error Failed To Get App ")
	}
	buttonGrid.
		AddItem(newButton, 0, 0, 1, 1, 1, 1, true).
		AddItem(updateButton, 0, 1, 1, 1, 1, 1, true)

	mainGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'n':
			uiConfig.App.SetFocus(newButton)
		case 'u':
			uiConfig.App.SetFocus(updateButton)

		}
		return event
	})
	return &ProfileFooterUI{
		MainGrid:     buttonGrid,
		UpdateButton: updateButton,
		NewButton:    newButton,
	}

}
