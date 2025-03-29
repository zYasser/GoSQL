package views

import (
	"GoSQL/internal/config"
	"GoSQL/internal/services"
	"GoSQL/internal/ui/router"
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitiateCreateProfileView(ctx context.Context, pageIdx int) *tview.Flex {
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

	form.AddButton("Connect", func() {
		// Validate form inputs
		input := getDatabaseConfig(form)
		if validateInput(input, errorText) {
			err := services.CreateProfile(*input, ctx)
			if err != nil {
				errorText.SetText(err.Error())
			}
		}
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

	return mainFlex
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

func getDatabaseConfig(form *tview.Form) *config.DatabaseConnectionInput {
	profileName := form.GetFormItemByLabel("Profile Name").(*tview.InputField).GetText()
	host := form.GetFormItemByLabel("Host").(*tview.InputField).GetText()
	port := form.GetFormItemByLabel("Port").(*tview.InputField).GetText()
	username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	databaseName := form.GetFormItemByLabel("Database Name").(*tview.InputField).GetText()

	return &config.DatabaseConnectionInput{
		ProfileName:  profileName,
		Host:         host,
		Port:         port,
		Username:     username,
		Password:     password,
		DatabaseName: databaseName,
	}
}
