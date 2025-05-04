package components

import (
	"GoSQL/internal/config"
	"GoSQL/internal/services"
	"context"
	"fmt"
	"os"

	"github.com/rivo/tview"
)

type ProfileListUI struct {
	MainFlex *tview.Flex
	List     *tview.List
	Profiles *[]config.DatabaseConnectionInput
}

func InitiateProfileList(ctx context.Context) *ProfileListUI {
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
		return &ProfileListUI{MainFlex: main}
	}
	items := []ListProps{}
	for _, value := range profiles {
		items = append(items, ListProps{
			mainText:      value.ProfileName,
			secondaryText: fmt.Sprintf("Database:%s", value.DatabaseName),
		})
	}
	list := createList(items)
	list.SetBorder(true).SetTitle("Choose Profile")
	main.AddItem(list, 0, 1, true)

	return &ProfileListUI{
		MainFlex: main,
		List:     list,
		Profiles: &profiles,
	}
}
