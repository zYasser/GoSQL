package views

import (
	"GoSQL/internal/constants"
	"GoSQL/internal/services"
	"context"

	"github.com/rivo/tview"
)

func InitiateCreateProfileView(ctx context.Context) *tview.Form {
	form := tview.NewForm().
		AddInputField("Profile Name", "", 20, nil, nil).
		AddDropDown("Database Type", constants.GetAllDatabaseType(), 0, nil).
		AddInputField("Host", "localhost", 20, nil, nil).
		AddInputField("Port", "", 5, nil, nil).
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddInputField("Database Name", "", 20, nil, nil).
		AddButton("Connect", nil).
		AddButton("Cancel", nil)

	form.SetBorder(true).SetTitle("Database Connection").SetTitleAlign(tview.AlignLeft)

	return form
}

func getDatabaseConfig(form *tview.Form) *services.DatabaseConnectionInput {
	profileName := form.GetFormItemByLabel("Profile Name").(*tview.InputField).GetText()
	_, db := form.GetFormItemByLabel("Database Type").(*tview.DropDown).GetCurrentOption()
	host := form.GetFormItemByLabel("Host").(*tview.InputField).GetText()
	port := form.GetFormItemByLabel("Port").(*tview.InputField).GetText()
	username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	databaseName := form.GetFormItemByLabel("Database Name").(*tview.InputField).GetText()
	return &services.DatabaseConnectionInput{
		ProfileName:  profileName,
		DatabaseType: db,
		Host:         host,
		Port:         port,
		Username:     username,
		Password:     password,
		DatabaseName: databaseName,
	}
}
