package components

import (
	"GoSQL/internal/config"
	"GoSQL/internal/ui/router"
	"context"

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
	buttonGrid.
		AddItem(newButton, 0, 0, 1, 1, 1, 1, true).
		AddItem(updateButton, 0, 1, 1, 1, 1, 1, true)

	return &ProfileFooterUI{
		MainGrid:     buttonGrid,
		UpdateButton: updateButton,
		NewButton:    newButton,
	}

}
