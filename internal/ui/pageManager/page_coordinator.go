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
		ViewsIndexMap: make(map[int]config.PageConfig),
		App:           app,
	}
	db := &config.DbConfig{}
	ctx = context.WithValue(ctx, "ui-config", uiConfig)
	ctx = context.WithValue(ctx, "db", db)

	pageIndex := 0
	// queryPage := views.InitializeQueryView(ctx, pageIndex)
	profilePageView := views.InitProfileView(pageIndex, ctx)
	pageIndex++
	createProfilePageView := views.InitiateCreateProfileView(ctx, pageIndex)
	pageIndex++
	queryPage, queryfunc := views.NewQueryViewPage(ctx, pageIndex)
	pageIndex++

	uiConfig.Main.AddPage(string(config.ProfilePage), profilePageView.MainGrid, true, true)
	uiConfig.Main.AddPage(string(config.CreateProfilePage), createProfilePageView, true, false)
	uiConfig.Main.AddPage(string(config.QueryPage), queryPage, true, false)

	uiConfig.ViewsIndexMap[0] = config.PageConfig{
		Page:     string(config.ProfilePage),
		PageFunc: func() {},
	}
	uiConfig.ViewsIndexMap[1] = config.PageConfig{
		Page:     string(config.CreateProfilePage),
		PageFunc: func() {},
	}
	uiConfig.ViewsIndexMap[2] = config.PageConfig{
		Page:     string(config.QueryPage),
		PageFunc: queryfunc,
	}

	return uiConfig
}
