package components

import (
	"GoSQL/internal/services"
	"os"

	"github.com/rivo/tview"
)

func InitiateProfileList() *tview.Flex {
	profiles, err := services.GetProfiles()
	main := tview.NewFlex()
	if err != nil {
		modal := tview.NewModal().
			SetText(err.Error()).
			AddButtons([]string{"Quit"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Quit" {
					os.Exit(1)
				}
			})
		main.AddItem(modal, 1, 1, true)
		return main

	}
	items := []ListProps{}
	for _, value := range profiles {
		items = append(items, ListProps{
			mainText:      value.ProfileName,
			secondaryText: value.DatabaseName,
		})
	}
	list := createList(items)
	list.SetBorder(true).SetTitle("Choose Profile")

	main.AddItem(list, 0, 1, true)
	return main
}
