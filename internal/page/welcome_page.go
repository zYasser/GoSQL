package page

import (
	"github.com/rivo/tview"
)

// initHomePage creates the initial welcome page with a title and database selection list
func initHomePage(currentPage int) *tview.Flex {
	asciiArt := `
  ________ ________    _________________  .____     
 /  _____/ \_____  \  /   _____/\_____  \ |    |    
/   \  ___  /   |   \ \_____  \  /  / \  \|    |    
\    \_\  \/    |    \/        \/   \_/.  \    |___ 
 \______  /\_______  /_______  /\_____\ \_/_______ \
        \/         \/        \/        \__>       \/
`

	art := tview.NewTextView().
		SetText(asciiArt ).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	welcomeText := tview.NewTextView().
		SetText("Welcome to the Database Manager!\n\nPlease select a database to continue." ).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	dbList := tview.NewList().
		AddItem("PostgreSQL", "Use PostgreSQL database", '1', func() {
			changePage(Next, currentPage)
		}).
		AddItem("MySQL", "Use MySQL database", '2', func() {
			changePage(Next, currentPage)
		}).
		AddItem("SQLite", "Use SQLite database", '3', func() {
			changePage(Next, currentPage)
		})
	// Use Flex layout to center the elements
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(art, 6, 1, false). // ASCII art and welcome text
		AddItem(welcomeText, 6, 1, false). // ASCII art and welcome text

		AddItem(dbList, 0, 2, true)        // List takes remaining space

	return flex
}
