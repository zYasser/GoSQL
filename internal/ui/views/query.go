package views

import (
	"GoSQL/internal/config"
	"GoSQL/internal/services"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitializeQueryView(ctx context.Context, pageIdx int) *tview.Flex {
	uiConfig, _ := ctx.Value("ui-config").(*config.UIConfig)

	mainFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	_ = config.Test(ctx)
	result, _ := services.GetTables(ctx)
	itemList := createItemList(result)
	mainFlex.AddItem(itemList, 0, 1, true)
	showDatabaseList := true
	rightPanel := createRightPanel(uiConfig)
	mainFlex.AddItem(rightPanel, 0, 4, false)
	mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			uiConfig.App.SetFocus(rightPanel)
			if showDatabaseList {
				mainFlex.ResizeItem(itemList, 0, 0)
			} else {
				mainFlex.ResizeItem(itemList, 0, 1)

			}
			showDatabaseList = !showDatabaseList

		}
		return event

	})
	return mainFlex
}
func createItemList(rows map[string][]string) *tview.TreeView {
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	root.SetReference(root)

	for schema, tables := range rows {
		node := tview.NewTreeNode(schema).SetSelectable(true).SetColor(tcell.ColorDarkGreen).SetReference(root)
		for _, table := range tables {
			child := tview.NewTreeNode(table).SetColor(tcell.ColorBlue).SetReference(node)
			node.AddChild(child)
		}
		root.AddChild(node)
	}

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		node.SetExpanded(!node.IsExpanded())
	})

	// Handle Vim-like movement
	var numBuffer string

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentNode := tree.GetCurrentNode()
		if currentNode == nil {
			return event
		}

		switch {
		case event.Key() == tcell.KeyUp:
			if numBuffer != "" {
				steps, _ := strconv.Atoi(numBuffer)

				moveTreeView(tree, currentNode, steps, false)
				numBuffer = "" // Reset buffer

			}

			return event

		case event.Key() == tcell.KeyDown:
			if numBuffer != "" {
				steps, _ := strconv.Atoi(numBuffer)

				moveTreeView(tree, currentNode, steps, true)
				numBuffer = "" // Reset buffer

			}
			return event

		case event.Rune() >= '1' && event.Rune() <= '9':
			// Capture number input
			numBuffer += string(event.Rune())
			return event

		case event.Rune() == 'd' || event.Rune() == 'u':
			// Convert captured number
			if numBuffer == "" {
				numBuffer = "1"
			}
			steps, _ := strconv.Atoi(numBuffer)
			numBuffer = "" // Reset buffer

			// Move up or down
			moveTreeView(tree, currentNode, steps, event.Rune() == 'd')
			return event

		case event.Rune() == rune(tcell.KeyCtrlX):

			if currentNode == root {
				return event
			}
			parent := currentNode.GetReference().(*tview.TreeNode)

			parent.SetExpanded(!parent.IsExpanded())
		}

		return event
	})

	return tree
}

func moveTreeView(tree *tview.TreeView, currentNode *tview.TreeNode, steps int, moveDown bool) {
	parent := currentNode.GetReference().(*tview.TreeNode)
	children := parent.GetChildren()
	index := -1

	// Find current node index
	for i, child := range children {
		if child == currentNode {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	newIndex := index + steps
	if !moveDown {
		newIndex = index - steps
	}

	if newIndex < 0 {
		newIndex = 0
	} else if newIndex >= len(children) {
		newIndex = len(children) - 1
	}

	tree.SetCurrentNode(children[newIndex])
}

func createRightPanel(uiConfig *config.UIConfig) *tview.Flex {
	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	queryInput := createQueryInput()
	rightPanel.AddItem(queryInput, 0, 1, false)

	dataAndDetail := createDataAndDetail(uiConfig)
	rightPanel.AddItem(dataAndDetail, 0, 5, true)

	return rightPanel
}

func createQueryInput() *tview.TextArea {
	queryInput := tview.NewTextArea()
	queryInput.SetPlaceholder("Enter SQL query here...").
		SetTitle("SQL Query").SetBorder(true)

	return queryInput
}

func createDataAndDetail(uiConfig *config.UIConfig) *tview.Flex {
	dataAndDetail := tview.NewFlex().SetDirection(tview.FlexColumn)

	dataTable := createDataTable(uiConfig)
	dataAndDetail.AddItem(dataTable, 0, 2, true)

	detailView := createDetailView(dataTable)
	dataAndDetail.AddItem(detailView, 0, 1, false)

	return dataAndDetail
}

func createDataTable(uiConfig *config.UIConfig) *tview.Table {
	dataTable := tview.NewTable().SetBorders(true).SetSelectable(true, false)
	dataTable.SetTitle("Data").SetBorder(true)

	for row := 0; row < 25; row++ {
		for col := 0; col < 5; col++ {
			text := randomText()
			cell := tview.NewTableCell(text).
				SetTextColor(tcell.ColorWhite).
				SetMaxWidth(40).
				SetAlign(tview.AlignLeft)

			if row == 0 {
				cell = cell.SetAttributes(tcell.AttrBold)
			}

			dataTable.SetCell(row, col, cell)
		}
	}
	dataTable.SetFixed(1, 0)

	dataTable.SetSelectionChangedFunc(func(row int, column int) {
		selectedCell := dataTable.GetCell(row, column)
		fmt.Printf("Selected cell text: %s\n", selectedCell.Text)
	})

	dataTable.SetSelectedFunc(func(row, column int) {
		uiConfig.App.SetFocus(dataTable)
	})

	return dataTable
}

func createDetailView(dataTable *tview.Table) *tview.TextArea {
	detailView := tview.NewTextArea().
		SetPlaceholder("Details of selected data...")
	detailView.SetTitle("Title").SetBorder(true)
	detailView.SetBorderColor(tcell.ColorDarkRed)

	x, y := 0, 0
	dataTable.SetSelectionChangedFunc(func(row int, column int) {
		cell := dataTable.GetCell(x, y)
		if cell != nil {
			dataTable.GetCell(x, y).SetBackgroundColor(tcell.ColorDefault)
		}
		dataTable.GetCell(row, column).SetBackgroundColor(tcell.ColorGray)
		x, y = row, column
		cell = dataTable.GetCell(row, column)
		if cell != nil {
			detailView.SetText(cell.Text, true)
		}
	})
	return detailView
}

func randomText() string {
	texts := []string{
		"Short",
		"Medium length textMedium length textMedium length textMedium length textMedium length textMedium length textMedium length textMedium length textMedium length textMedium length textMedium length textMedium length text",
		"Tiny",
		"Somewhat longer text",
	}
	return texts[rand.Intn(len(texts))]
}

func repeatText(text string, n int) string {
	return strings.Repeat(text+" ", n)
}
