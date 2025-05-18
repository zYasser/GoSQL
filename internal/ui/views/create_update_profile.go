package views

import (
	"GoSQL/internal/config"
	"GoSQL/internal/services"
	"GoSQL/internal/ui/router"
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CreateUpdateProfileView struct {
	MainFlex *tview.Flex
	Form     *tview.Form
	Id       string
	isUpdate bool
}

func (view *CreateUpdateProfileView) AddProfileToForm(id string) {
	if id == "" {
		return
	}
	profile, err := services.GetProfile(id)
	if err != nil {
		return
	}
	view.Form.GetFormItemByLabel("Profile Name").(*tview.InputField).SetText(profile.ProfileName)
	view.Form.GetFormItemByLabel("Host").(*tview.InputField).SetText(profile.Host)
	view.Form.GetFormItemByLabel("Port").(*tview.InputField).SetText(profile.Port)
	view.Form.GetFormItemByLabel("Username").(*tview.InputField).SetText(profile.Username)
	view.Form.GetFormItemByLabel("Database Name").(*tview.InputField).SetText(profile.DatabaseName)
	view.Id = id
	view.isUpdate = true
}

func InitiateCreateUpdateProfileView(ctx context.Context, pageIdx int) (*CreateUpdateProfileView, func(id string)) {
	// Error message text view
	errorText := tview.NewTextView().
		SetTextColor(tcell.ColorRed).
		SetText("")
	errorText.SetTextAlign(tview.AlignCenter)
	form := tview.NewForm().
		AddInputField("Profile Name", "", 20, nil, nil).
		AddInputField("Host", "localhost", 20, nil, nil).
		AddInputField("Port", "5432", 5, nil, nil).
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddInputField("Database Name", "", 20, nil, nil)
	hiddenField := tview.NewInputField()
	profileView := &CreateUpdateProfileView{
		Id:   "",
		Form: form,
	}
	form.AddButton("Connect", func() {
		// Validate form inputs
		input := getDatabaseConfig(form, profileView.Id)
		if validateInput(input, errorText) {
			if profileView.isUpdate {
				err := services.UpdateProfile(*input, profileView.Id, ctx)
				if err != nil {
					errorText.SetText(err.Error())
					return
				}
			} else {
				err := services.CreateProfile(*input, ctx)
				if err != nil {
					errorText.SetText(err.Error())
					return
				}
			}
		}
		router.NavigatePage(config.Previous, pageIdx, ctx)
	}).
		AddButton("Cancel", func() {
			router.NavigatePage(config.Previous, pageIdx, ctx)
		})

	// Track focus index
	focusIndex := 0
	totalFields := form.GetFormItemCount() + form.GetButtonCount()
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			if focusIndex > 0 {
				focusIndex--
				form.SetFocus(focusIndex)
			}
			return nil

		case tcell.KeyTab:
			if focusIndex < totalFields-1 {
				focusIndex++
				form.SetFocus(focusIndex)
			}
			return nil

		case tcell.KeyDown:
			if focusIndex < totalFields-1 {
				focusIndex++
				form.SetFocus(focusIndex)
			}
			return nil
		case tcell.KeyCtrlX:
			form.SetFocus(form.GetFormItemCount())
			focusIndex = form.GetFormItemCount()
		}

		return event
	})

	form.SetBorder(true).SetTitle("Database Connection").SetTitleAlign(tview.AlignLeft)

	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	mainFlex.AddItem(tview.NewBox(), 0, 1, false)

	horizontalFlex := tview.NewFlex()
	horizontalFlex.AddItem(tview.NewBox(), 0, 1, false)
	horizontalFlex.AddItem(form, 40, 1, true)
	horizontalFlex.AddItem(tview.NewBox(), 0, 1, false)

	mainFlex.AddItem(horizontalFlex, 20, 3, true)

	mainFlex.AddItem(errorText, 3, 1, false)

	mainFlex.AddItem(tview.NewBox(), 0, 1, false)
	mainFlex.AddItem(hiddenField, 0, 0, false)
	profileView.MainFlex = mainFlex
	return profileView, func(id string) {

		profileView.AddProfileToForm(id)
	}
}

// Validate input and display error messages
func validateInput(input *config.DatabaseConnectionInput, errorView *tview.TextView) bool {
	// Basic validation
	if input.ProfileName == "" {
		errorView.SetText("Profile Name is required")
		return false
	}
	if input.Host == "" {
		errorView.SetText("Host is required")
		return false
	}
	if input.Port == "" {
		errorView.SetText("Port is required")
		return false
	}
	if input.Username == "" {
		errorView.SetText("Username is required")
		return false
	}
	if input.DatabaseName == "" {
		errorView.SetText("Database Name is required")
		return false
	}

	// Clear any previous error messages
	errorView.SetText("")
	return true
}

func getDatabaseConfig(form *tview.Form, id string) *config.DatabaseConnectionInput {
	profileName := form.GetFormItemByLabel("Profile Name").(*tview.InputField).GetText()
	host := form.GetFormItemByLabel("Host").(*tview.InputField).GetText()
	port := form.GetFormItemByLabel("Port").(*tview.InputField).GetText()
	username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	databaseName := form.GetFormItemByLabel("Database Name").(*tview.InputField).GetText()
	return &config.DatabaseConnectionInput{
		ID:           id,
		ProfileName:  profileName,
		Host:         host,
		Port:         port,
		Username:     username,
		Password:     password,
		DatabaseName: databaseName,
	}
}
