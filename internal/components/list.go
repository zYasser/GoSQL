package components

import "github.com/rivo/tview"

type ListProps struct {
	mainText      string
	secondaryText string
	shortcut      rune
	selected      func()
}

func createList(props []ListProps) *tview.List {
	list := tview.NewList()

	for _, prop := range props {
		list.AddItem(prop.mainText, prop.secondaryText, prop.shortcut, prop.selected)
	}

	return list
}
