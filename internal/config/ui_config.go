package config

import (
	"github.com/rivo/tview"
)

type Page string

const (
	CreateProfilePage Page = "CreateProfile"
	DbChoose          Page = "DbChoose"
	ProfilePage       Page = "ProfilePage"
	QueryPage         Page = "QueryPage"
	Next              Page = "Next"
	Previous          Page = "Previous"
)

type UIConfig struct {
	Main          *tview.Pages
	ViewsIndexMap map[int]Page
	CurrentPage   Page
	App           *tview.Application
}
