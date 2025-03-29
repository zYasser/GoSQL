package views

import (
	"GoSQL/internal/components"
	"context"

	"github.com/rivo/tview"
)

func InitProfileView(currentPage int, ctx context.Context) *tview.Grid {
	mainGrid := tview.NewGrid().
		SetRows(1, 1, 1, 0, 1, 1, 1).
		SetColumns(1, 1, 1, 0, 1, 1, 1)

	list := components.InitiateProfileList()
	buttonGrid := components.CreateProfileFooter(ctx, mainGrid)

	contentGrid := tview.NewGrid().
		SetRows(0, 1, 1).
		AddItem(list, 0, 0, 1, 1, 0, 0, true).
		AddItem(buttonGrid, 1, 0, 1, 1, 0, 0, false)

	mainGrid.AddItem(contentGrid, 3, 3, 1, 1, 0, 0, true)

	return mainGrid
}
