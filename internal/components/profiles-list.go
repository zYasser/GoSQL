package components

import (
	"GoSQL/internal/config"
	"GoSQL/internal/services"
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ProfileListUI struct {
	MainFlex *tview.Flex
	List     *tview.List
	Profiles *[]config.DatabaseConnectionInput
}

func InitiateProfileList(ctx context.Context) (*ProfileListUI, func()) {
	ui := &ProfileListUI{
		MainFlex: tview.NewFlex(),
	}

	ui.renderProfileList(ctx)

	ui.MainFlex.AddItem(ui.List, 0, 1, true)
	ui.List.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlD {
			ui.deleteProfile()
		}
		return event
	})
	return ui, func() {
		ui.renderProfileList(ctx)
	}
}
func (p *ProfileListUI) GetSelectedProfile() config.DatabaseConnectionInput {
	selectedItem := p.List.GetCurrentItem()
	return (*p.Profiles)[selectedItem]
}

func (p *ProfileListUI) deleteProfile() {
	selectedItem := p.List.GetCurrentItem()
	err := services.DeleteProfile((*p.Profiles)[selectedItem].ID)
	if err != nil {
		modal := tview.NewModal().
			SetText(err.Error()).
			AddButtons([]string{"Quit"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Quit" {
					os.Exit(1)
				}
			})
		p.MainFlex.AddItem(modal, 1, 1, true)
		return

	}
	p.List.RemoveItem(selectedItem)
}

func (p *ProfileListUI) renderProfileList(ctx context.Context) {
	if p.List != nil {
		p.List.Clear()
	}
	profiles, err := services.GetProfiles()

	if err != nil {
		modal := tview.NewModal().
			SetText(err.Error()).
			AddButtons([]string{"Quit"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Quit" {
					os.Exit(1)
				}
			})
		p.MainFlex.AddItem(modal, 1, 1, true)
		return
	}

	profileList := make([]config.DatabaseConnectionInput, 0)
	for _, value := range profiles {
		profileList = append(profileList, value)
	}
	sort.Slice(profileList, func(i, j int) bool {
		return profileList[i].ProfileName < profileList[j].ProfileName
	})
	p.Profiles = &profileList
	items := []ListProps{}
	for _, value := range profiles {
		items = append(items, ListProps{
			mainText:      value.ProfileName,
			secondaryText: fmt.Sprintf("Database:%s", value.DatabaseName),
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].mainText < items[j].mainText
	})
	p.List = createList(items, p.List)
	p.List.SetBorder(true).SetTitle("Choose Profile")

}
