package components

import "github.com/rivo/tview"

func InitiateProfileList() *tview.List {
	props := []listProps{
		{"List item 1", "Some explanatory text", 'a', nil},
		{"List item 2", "Some explanatory text", 'b', nil},
		{"List item 3", "Some explanatory text", 'c', nil},
		{"List item 4", "Some explanatory text", 'd', nil},
	}

	list := createList(props)
	list.SetBorder(true).SetTitle("Choose Database")
	return list
}
