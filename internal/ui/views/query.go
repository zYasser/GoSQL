package views

import (
	"GoSQL/internal/config"
	"GoSQL/internal/services"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	Result  = "Result"
	Error   = "Error"
	Loading = "Loading"
)

type queryView struct {
	app                    *tview.Application
	table                  *tview.Table
	tree                   *tview.TreeView
	queryInput             *tview.TextArea
	detailView             *tview.TextArea
	mainFlex               *tview.Flex
	rightPanel             *tview.Flex
	dataAndDetail          *tview.Flex
	resultContainer        *tview.Flex
	statusModal            *tview.TextView
	uIConfig               *config.UIConfig
	context                context.Context
	x                      int
	y                      int
	isStatusModalDisplayed bool
	data                   [][]string
	newData                [][]bool
	pk                     string
	pkIndex                int
	tableName              string
}

func initializeQueryView(ctx context.Context, pageIdx int) *queryView {
	uiConfig, _ := ctx.Value("ui-config").(*config.UIConfig)

	qv := &queryView{
		uIConfig: uiConfig,
		context:  ctx,
		app:      uiConfig.App,
	}
	qv.mainFlex = tview.NewFlex().SetDirection(tview.FlexColumn)

	qv.statusModal = qv.createStatusModal()

	// Initialize UI components without data
	qv.tree = qv.createItemTree(nil)
	treeContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	treeContainer.AddItem(qv.tree, 0, 1, true)

	qv.mainFlex.AddItem(treeContainer, 0, 1, true)
	qv.rightPanel = qv.createRightPanel()
	qv.mainFlex.AddItem(qv.rightPanel, 0, 4, false)

	showDatabaseList := true
	qv.mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlD:
			if showDatabaseList {
				treeContainer.RemoveItem(qv.tree)
				qv.mainFlex.ResizeItem(treeContainer, 0, 0)
				qv.uIConfig.App.SetFocus(qv.rightPanel)
			} else {
				treeContainer.AddItem(qv.tree, 0, 1, true)
				qv.mainFlex.ResizeItem(treeContainer, 0, 1)

				qv.uIConfig.App.SetFocus(qv.tree)
			}
			showDatabaseList = !showDatabaseList
			return nil
		}
		return event
	})
	qv.tree.SetSelectedFunc(func(node *tview.TreeNode) {
		if len(node.GetChildren()) == 0 {
			parent := node.GetReference().(*tview.TreeNode)
			query := fmt.Sprintf("\"%s\".\"%s\"", parent.GetText(), node.GetText())

			qv.showStatus("Loading", "Loading table data...")

			data, pk, err := services.FetchTableData(ctx, query)

			if err != nil {
				qv.showStatus("Error", fmt.Sprintf("Failed to fetch table data: %v", err))
			} else {
				qv.tableName = node.GetText()
				qv.pk = pk
				qv.data = data

				qv.pkIndex = qv.getPrimaryKeyIndex()
				qv.x, qv.y = 0, 0
				qv.newData = make([][]bool, len(data))
				for i := range qv.newData {
					qv.newData[i] = make([]bool, len(data[i]))
				}
				qv.hideStatus()
				qv.addDataToTable(data, qv.table)
				uiConfig.App.SetFocus(qv.table)
			}
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})
	qv.statusModal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			qv.hideStatus()
		}
		return nil
	})

	qv.switchComponents()

	return qv
}

func (qv *queryView) LoadDatabaseData() {
	qv.showStatus("Loading", "Loading database tables...")

	result, err := services.GetTables(qv.context)
	if err != nil {
		qv.showStatus("Error", fmt.Sprintf("Failed to get tables: %v", err))
		return
	}
	qv.app.SetFocus(qv.tree)
	// Update the tree with the retrieved data
	root := qv.tree.GetRoot()
	root.ClearChildren()

	for schema, tables := range result {
		node := tview.NewTreeNode(schema).SetSelectable(true).SetColor(tcell.ColorDarkGreen).SetReference(root)
		for _, table := range tables {
			child := tview.NewTreeNode(table).SetColor(tcell.ColorBlue).SetReference(node)
			node.AddChild(child)
		}
		root.AddChild(node)
	}

	qv.hideStatus()
	qv.app.SetFocus(qv.tree)
}
func NewQueryViewPage(ctx context.Context, pageIdx int) (*tview.Flex, func()) {
	qv := initializeQueryView(ctx, pageIdx)

	// Return both the view and a function to load data
	return qv.mainFlex, func() {
		qv.LoadDatabaseData()
	}
}
func (qv *queryView) showStatus(statusType string, message string) {
	qv.statusModal.SetTitle(statusType)
	qv.statusModal.SetText(message)

	// Set color based on status type
	if statusType == "Error" {
		qv.statusModal.SetTextColor(tcell.ColorRed)
		qv.statusModal.SetBorderColor(tcell.ColorRed)

	} else if statusType == "Result" {
		qv.statusModal.SetTextColor(tcell.ColorDefault)
		qv.statusModal.SetBorderColor(tcell.ColorGreen)
	} else {
		qv.statusModal.SetTextColor(tcell.ColorWhite)
		qv.statusModal.SetBorderColor(tcell.ColorYellow)
	}
	if !qv.isStatusModalDisplayed {
		qv.resultContainer.RemoveItem(qv.table)
		qv.resultContainer.AddItem(qv.statusModal, 0, 2, true)
	}
	qv.app.SetFocus(qv.statusModal)
	qv.isStatusModalDisplayed = true
}

func (qv *queryView) hideStatus() {
	if !qv.isStatusModalDisplayed {
		return
	}
	qv.resultContainer.RemoveItem(qv.statusModal)
	qv.resultContainer.AddItem(qv.table, 0, 2, true)

	qv.app.SetFocus(qv.queryInput)

	qv.isStatusModalDisplayed = false
}

func (qv *queryView) switchComponents() {
	qv.mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlT:
			qv.app.SetFocus(qv.tree)
		case tcell.KeyCtrlK:
			qv.app.SetFocus(qv.queryInput)
		case tcell.KeyCtrlU:
			if !qv.isStatusModalDisplayed {
				qv.app.SetFocus(qv.table)
			}
		case tcell.KeyEscape:
			if qv.isStatusModalDisplayed {
				qv.hideStatus()
				return nil
			}
		}
		return event
	})

}

func (qv *queryView) createItemTree(rows map[string][]string) *tview.TreeView {
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	root.SetReference(root)

	// Only populate tree if data is provided
	if rows != nil {
		for schema, tables := range rows {
			node := tview.NewTreeNode(schema).SetSelectable(true).SetColor(tcell.ColorDarkGreen).SetReference(root)
			for _, table := range tables {
				child := tview.NewTreeNode(table).SetColor(tcell.ColorBlue).SetReference(node)
				node.AddChild(child)
			}
			root.AddChild(node)
		}
	}

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
				qv.moveTreeView(tree, currentNode, steps, false)
				numBuffer = ""
			}
			return event

		case event.Key() == tcell.KeyDown:
			if numBuffer != "" {
				steps, _ := strconv.Atoi(numBuffer)
				qv.moveTreeView(tree, currentNode, steps, true)
				numBuffer = ""
			}
			return event

		case event.Rune() >= '1' && event.Rune() <= '9':
			numBuffer += string(event.Rune())
			return event

		case event.Rune() == 'd' || event.Rune() == 'u':
			if numBuffer == "" {
				numBuffer = "1"
			}
			steps, _ := strconv.Atoi(numBuffer)
			numBuffer = ""

			qv.moveTreeView(tree, currentNode, steps, event.Rune() == 'd')
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

func (qv *queryView) moveTreeView(tree *tview.TreeView, currentNode *tview.TreeNode, steps int, moveDown bool) {
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

func (qv *queryView) createRightPanel() *tview.Flex {
	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow)

	qv.queryInput = qv.createQueryInput()
	rightPanel.AddItem(qv.queryInput, 0, 1, false)

	qv.dataAndDetail = qv.createDataAndDetail()
	rightPanel.AddItem(qv.dataAndDetail, 0, 5, true)

	return rightPanel
}

func (qv *queryView) createQueryInput() *tview.TextArea {
	queryInput := tview.NewTextArea()
	queryInput.SetPlaceholder("Enter SQL query here...").
		SetTitle("SQL Query").SetBorder(true)

	queryInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlA {
			query := queryInput.GetText()

			qv.showStatus("Loading", "Executing query...")

			data, message, err := services.ExecuteQuery(qv.context, query)
			if err != nil {
				qv.showStatus(Error, fmt.Sprintf("Query error: %v", err))
			} else if len(message) != 0 {
				qv.showStatus(Result, message)

			} else {
				qv.hideStatus()
				qv.addDataToTable(data, qv.table)
				qv.app.SetFocus(qv.table)
			}
			return nil
		}
		return event
	})

	return queryInput
}

func (qv *queryView) createDataAndDetail() *tview.Flex {
	dataAndDetail := tview.NewFlex().SetDirection(tview.FlexColumn)
	qv.table = qv.createDataTable()

	qv.resultContainer = tview.NewFlex()
	qv.resultContainer.AddItem(qv.table, 0, 2, true)
	qv.detailView = qv.createDetailView()
	dataAndDetail.AddItem(qv.resultContainer, 0, 2, false)
	dataAndDetail.AddItem(qv.detailView, 0, 1, false)

	return dataAndDetail
}

func (qv *queryView) createDataTable() *tview.Table {
	dataTable := tview.NewTable().SetBorders(true).SetSelectable(true, true)
	dataTable.SetTitle("Data").SetBorder(true)
	dataTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlR:
			qv.undoDataChange()
		case tcell.KeyRune:
			switch event.Rune() {
			case 'e': // Capital E
				qv.app.SetFocus(qv.detailView)
			case 's':
				qv.submitDataChange()
			}

		}

		return event

	})
	return dataTable
}

func (qv *queryView) addDataToTable(data [][]string, table *tview.Table) {
	table.Clear()

	for i, values := range data {
		for j, row := range values {
			cell := tview.NewTableCell(row).
				SetTextColor(tcell.ColorWhite).
				SetMaxWidth(40).
				SetAlign(tview.AlignLeft)
			table.SetCell(i, j, cell)
		}
	}
	table.SetFixed(1, 0)
}

func (qv *queryView) createDetailView() *tview.TextArea {

	detailView := tview.NewTextArea().
		SetPlaceholder("Details of selected data...")
	detailView.SetTitle("Details").SetBorder(true)
	detailView.SetBorderColor(tcell.ColorDarkRed)
	detailView.SetFocusFunc(func() {
		detailView.SetText(detailView.GetText(), true)
	})

	qv.table.SetSelectionChangedFunc(func(row int, column int) {
		if qv.isStatusModalDisplayed {
			return
		}

		cellColumnName := qv.table.GetCell(0, column)
		cell := qv.table.GetCell(qv.x, qv.y)
		if cell != nil {
			qv.table.GetCell(qv.x, qv.y).SetBackgroundColor(tcell.ColorDefault)
		}
		qv.table.GetCell(row, column).SetBackgroundColor(tcell.ColorGray)
		qv.x, qv.y = row, column
		cell = qv.table.GetCell(row, column)
		if cell != nil {
			if row != 0 {
				detailView.SetText(cell.Text, true)

				detailView.SetTitle(cellColumnName.Text)

			} else {
				detailView.SetText(cell.Text, true)
			}
		}
	})
	detailView.SetChangedFunc(func() {

		currentCell := qv.table.GetCell(qv.x, qv.y)
		style := tcell.StyleDefault.Foreground(tcell.ColorYellow)
		currentText := detailView.GetText()
		if currentCell.Text != currentText {
			currentCell.SetStyle(style)
			currentCell.SetText(currentText)
			qv.newData[qv.x][qv.y] = true

		}

	})
	detailView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlR:
			qv.undoDataChange()
		case tcell.KeyCtrlS:
			qv.submitDataChange()

		}

		return event

	})
	return detailView
}

func (qv *queryView) undoDataChange() {
	currentCell := qv.table.GetCell(qv.x, qv.y)
	style := qv.table.GetCell(0, 0).Style
	qv.detailView.SetText(qv.data[qv.x][qv.y], true)
	currentCell.SetStyle(style)
	qv.newData[qv.x][qv.y] = false
}

func (qv *queryView) submitDataChange() {
	updateMap := make(map[string]map[string]string) // pkVal -> {col: newVal}

	for i := range qv.newData {
		for j := range qv.newData[i] {
			if qv.newData[i][j] {
				col := strings.Split(qv.table.GetCell(0, j).Text, ":")[0]
				val := qv.table.GetCell(i, j).Text
				pkVal := qv.table.GetCell(i, qv.pkIndex).Text

				if _, ok := updateMap[pkVal]; !ok {
					updateMap[pkVal] = make(map[string]string)
				}
				updateMap[pkVal][col] = val
			}
		}
	}

	if len(updateMap) == 0 {
		qv.showStatus("Info", "No changes to submit")
		return
	}

	updates := make([]services.UpdateQueryParams, 0, len(updateMap))
	for pkVal, values := range updateMap {
		updates = append(updates, services.UpdateQueryParams{
			Key:    pkVal,
			Values: values,
		})
	}

	err := services.GenerateSQL(updates, qv.pk, qv.tableName, qv.context)
	if err != nil {
		qv.showStatus("Error", err.Error())

	}
}
func (qv *queryView) getPrimaryKeyIndex() int {
	if len(qv.data) == 0 || len(qv.data[0]) == 0 {
		return -1 // no data or headers
	}

	for i, column := range qv.data[0] {
		if strings.Split(column, ":")[0] == qv.pk {
			return i
		}
	}
	return -1
}

func (qv *queryView) createStatusModal() *tview.TextView {
	statusModal := tview.NewTextView()
	statusModal.SetBorder(true)

	return statusModal
}
