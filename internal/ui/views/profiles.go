package views

import (
	"GoSQL/internal/components"
	"GoSQL/internal/config"
	"GoSQL/internal/ui/router"
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ProfileView struct {
	MainGrid       *tview.Grid
	list           *tview.List
	buttonGrid     *tview.Grid
	errorText      *tview.TextView
	RenderFunction func()
}

func InitProfileView(currentPage int, ctx context.Context) *ProfileView {
	uiConfig, _ := ctx.Value("ui-config").(*config.UIConfig)

	mainGrid := tview.NewGrid().
		SetRows(1, 1, 1, 0, 1, 1, 1).
		SetColumns(1, 1, 1, 0, 1, 1, 1)

	list, function := components.InitiateProfileList(ctx)

	buttonGrid := components.CreateProfileFooter(ctx, mainGrid)
	errorText := tview.NewTextView().
		SetTextColor(tcell.ColorRed).
		SetText("").
		SetTextAlign(tview.AlignCenter)

	contentGrid := tview.NewGrid().
		SetRows(0, 1, 1).
		AddItem(list.MainFlex, 0, 0, 1, 1, 0, 0, true).
		AddItem(buttonGrid.MainGrid, 1, 0, 1, 1, 0, 0, false).
		AddItem(errorText, 2, 0, 1, 1, 0, 0, false)

	mainGrid.AddItem(contentGrid, 3, 3, 1, 1, 0, 0, true)
	uiConfig.App.SetFocus(list.MainFlex)

	view := &ProfileView{
		MainGrid:       mainGrid,
		list:           list.List,
		buttonGrid:     buttonGrid.MainGrid,
		errorText:      errorText,
		RenderFunction: function,
	}
	buttonGrid.UpdateButton.SetSelectedFunc(func() {
		id := list.GetSelectedProfile().ID
		router.NavigatePage(config.CreateProfilePage, -1, ctx, id)
	})
	ListFunc(view, *list.Profiles, ctx, currentPage)
	InputCapture(view, uiConfig, buttonGrid,)
	return view
}

func ListFunc(view *ProfileView, profiles []config.DatabaseConnectionInput, ctx context.Context, currentPage int) {
	view.list.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		cfg := profiles[i]
		err := config.ConnectToDb(cfg, ctx, true)
		if err != nil {
			view.errorText.SetText(err.Error())
		} else {

			router.NavigatePage(config.QueryPage, currentPage, ctx)
		}
	})

}

func InputCapture(view *ProfileView, uiConfig *config.UIConfig, ProfileFooterUI *components.ProfileFooterUI) {
	uiConfig.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'n':
			uiConfig.App.SetFocus(ProfileFooterUI.NewButton)
			return event
		case 'u':
			uiConfig.App.SetFocus(ProfileFooterUI.UpdateButton)
			return event
		case 'p':
			uiConfig.App.SetFocus(view.list)
			return event
		}
		return event
	})
}
