package router

import (
	"GoSQL/internal/config"
	"context"
)

func NavigatePage(page config.Page, currentPageIndex int, ctx context.Context, args ...interface{}) {
	uiConfig := ctx.Value("ui-config").(*config.UIConfig)

	var newPageIndex int
	if page == config.Next {
		newPageIndex = currentPageIndex + 1
	} else if page == config.Previous {
		newPageIndex = currentPageIndex - 1
	} else {
		var newPage config.PageConfig
		for _, v := range uiConfig.ViewsIndexMap {
			if v.Page == string(page) {
				newPage = v
				break
			}

		}
		switchToPageView(newPage, uiConfig, args...)
		return
	}

	newPage, exists := uiConfig.ViewsIndexMap[newPageIndex]
	if !exists {
		return
	}

	switchToPageView(newPage, uiConfig, args...)
}

func switchToPageView(page config.PageConfig, config *config.UIConfig, args ...interface{}) {
	config.Main.SwitchToPage(page.Page)
	page.PageFunc(args...)

}
