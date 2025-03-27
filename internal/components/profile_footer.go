package components

import (
	"GoSQL/internal/config"
	"GoSQL/internal/ui/router"
	"context"

	"github.com/rivo/tview"
)

func CreateProfileFooter(ctx context.Context) *tview.Grid {
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

	buttonGrid.
		AddItem(newButton, 0, 1, 1, 1, 0, 0, false).
		AddItem(deleteButton, 0, 2, 1, 1, 0, 0, false).
		AddItem(updateButton, 0, 3, 1, 1, 0, 0, false)

	return buttonGrid

}
