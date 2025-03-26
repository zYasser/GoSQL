package components

import "github.com/rivo/tview"

type listProps struct {
	mainText      string
	secondaryText string
	shortcut      rune
	selected      func()
}

func createList(props []listProps) *tview.List {
	list := tview.NewList()

	for _, prop := range props {
		list.AddItem(prop.mainText, prop.secondaryText, prop.shortcut, prop.selected)
	}

	return list
}
