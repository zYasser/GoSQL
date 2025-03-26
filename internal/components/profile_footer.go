package components

import "github.com/rivo/tview"

func CreateProfileFooter() *tview.Grid {
	buttonGrid := tview.NewGrid().
		SetRows(1)

	chooseButton := tview.NewButton("Choose").SetSelectedFunc(func() {
		println("Choose action triggered")
	})
	deleteButton := tview.NewButton("Delete").SetSelectedFunc(func() {
		println("Delete action triggered")
	})
	updateButton := tview.NewButton("Update").SetSelectedFunc(func() {
		println("Update action triggered")
	})
	newButton := tview.NewButton("New").SetSelectedFunc(func() {
		println("New action triggered")
	})

	buttonGrid.AddItem(chooseButton, 0, 0, 1, 1, 0, 0, false).
		AddItem(deleteButton, 0, 1, 1, 1, 0, 0, false).
		AddItem(updateButton, 0, 2, 1, 1, 0, 0, false).
		AddItem(newButton, 0, 3, 1, 1, 0, 0, false)

	return buttonGrid

}
