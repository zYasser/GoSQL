package pageManager

import (
	"GoSQL/internal/config"
	"GoSQL/internal/ui/views"
	"context"

	"github.com/rivo/tview"
)

func InitializePages(ctx context.Context, app *tview.Application) *config.UIConfig {
	uiConfig := &config.UIConfig{
		Main:          tview.NewPages(),
		ViewsIndexMap: make(map[int]config.Page),
		App:           app,
	}
	db := &config.DbConfig{}
	ctx = context.WithValue(ctx, "ui-config", uiConfig)
	ctx = context.WithValue(ctx, "db", db)

	pageIndex := 0
	profilePageView := views.InitProfileView(pageIndex, ctx)
	pageIndex++
	createProfilePageView := views.InitiateCreateProfileView(ctx, pageIndex)
	pageIndex++
	uiConfig.Main.AddPage(string(config.ProfilePage), profilePageView, true, true)
	uiConfig.Main.AddPage(string(config.CreateProfilePage), createProfilePageView, true, false)

	uiConfig.ViewsIndexMap[0] = config.ProfilePage
	uiConfig.ViewsIndexMap[1] = config.CreateProfilePage

	return uiConfig
}
