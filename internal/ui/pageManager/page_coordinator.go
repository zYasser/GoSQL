package pageManager

import (
	"GoSQL/internal/config"
	"GoSQL/internal/ui/views"
	"context"

	"github.com/rivo/tview"
)

func InitializePages(ctx context.Context) *config.UIConfig {
	pageNavigatorInstance := &config.UIConfig{
		Main:          tview.NewPages(),
		ViewsIndexMap: make(map[int]config.Page),
	}
	pageIndex := 0
	profilePageView := views.InitProfileView(pageIndex, ctx)
	pageIndex++
	createProfilePageView := views.InitiateCreateProfileView(ctx)
	pageIndex++
	pageNavigatorInstance.Main.AddPage(string(config.ProfilePage), profilePageView, true, true)
	pageNavigatorInstance.Main.AddPage(string(config.CreateProfilePage), createProfilePageView, true, false)

	pageNavigatorInstance.ViewsIndexMap[0] = config.ProfilePage
	pageNavigatorInstance.ViewsIndexMap[1] = config.CreateProfilePage

	return pageNavigatorInstance
}
