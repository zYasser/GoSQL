package page

import (
	"GoSQL/internal/components"

	"github.com/rivo/tview"
)

func initProfilePage(currentPage int) *tview.Grid {
	list := components.InitiateProfileList()
	buttonGrid := components.CreateProfileFooter()

	contentGrid := tview.NewGrid().
		SetRows(0, 1, 1).
		AddItem(list, 0, 0, 1, 1, 0, 0, true).
		AddItem(buttonGrid, 1, 0, 1, 1, 0, 0, false)

	mainGrid := tview.NewGrid().
		SetRows(1, 1, 1, 0, 1, 1, 1).
		SetColumns(1, 1, 1, 0, 1, 1, 1)

	mainGrid.AddItem(contentGrid, 3, 3, 1, 1, 0, 0, true)

	return mainGrid
}
