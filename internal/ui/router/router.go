package router

import (
	"GoSQL/internal/config"
	"context"
	"fmt"
)

func NavigatePage(page config.Page, currentPageIndex int, ctx context.Context) {
	uiConfig := ctx.Value("ui-config").(*config.UIConfig)

	var newPageIndex int
	if page == config.Next {
		newPageIndex = currentPageIndex + 1
	} else if page == config.Previous {
		newPageIndex = currentPageIndex - 1
	} else {
		switchToPageView(string(page), uiConfig)
		return
	}

	newPage, exists := uiConfig.ViewsIndexMap[newPageIndex]
	if !exists {
		fmt.Println("No Page Left")
		return
	}

	switchToPageView(string(newPage), uiConfig)
}

func switchToPageView(page string, config *config.UIConfig) {
	config.Main.SwitchToPage(page)
}
