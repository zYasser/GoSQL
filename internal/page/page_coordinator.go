package page

import (
	"fmt"

	"github.com/rivo/tview"
)

type Page string

const (
	DbChoose   Page = "DbChoose"
	SecondPage Page = "SecondPage"
	Next       Page = "Next"
	Previous   Page = "Previous"
)

type PageManager struct {
	Main    *tview.Pages
	PageMap map[int]Page
}

var pages *PageManager

func InitPages() *PageManager {

	pages = &PageManager{
		Main:    tview.NewPages(),
		PageMap: make(map[int]Page),
	}
	wp := initHomePage(0)
	sc := initProfilePage(1)

	pages.Main.AddPage(string(DbChoose), wp, true, true)
	pages.Main.AddPage(string(SecondPage), sc, true, false)
	pages.PageMap[0] = DbChoose
	pages.PageMap[1] = SecondPage

	return pages

}

func changePage(page Page, currentPage int) {

	var newIndex int
	if page == Next {
		newIndex = currentPage + 1
	} else if page == Previous {
		newIndex = currentPage - 1
	} else {
		pages.switchToPage(string(page))
		return
	}

	newPage, exists := pages.PageMap[newIndex]
	if !exists {
		fmt.Println("No Page Left")
		return
	}

	pages.switchToPage(string(newPage))
}

func (pm *PageManager) switchToPage(page string) {
	pm.Main.SwitchToPage(page)
}
