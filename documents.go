package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/typesense/typesense-go/typesense/api"
)

type DocumentsView struct {
	app         *App
	mainFlex    *tview.Flex
	list        *tview.List
	details     *tview.TextView
	statusBar   *tview.TextView
	collection  string
	documents   []map[string]interface{}
	selectedIdx int
	page        int
	perPage     int
}

func NewDocumentsView(app *App) *DocumentsView {
	dv := &DocumentsView{
		app:         app,
		selectedIdx: 0,
		page:        1,
		perPage:     20,
	}

	dv.list = tview.NewList().
		SetSelectedFunc(dv.onSelect).
		SetChangedFunc(dv.onSelectionChange)

	dv.details = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true).
		SetScrollable(true)
	dv.details.SetTitle("Document Details").SetBorder(true)

	dv.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	dv.mainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				AddItem(dv.list, 0, 1, true).
				AddItem(dv.details, 0, 2, false),
			0, 1, true,
		).
		AddItem(dv.statusBar, 1, 0, false)

	dv.list.SetTitle("Documents").SetBorder(true)
	dv.list.SetInputCapture(dv.handleInput)

	dv.updateStatusBar()

	return dv
}

func (dv *DocumentsView) Render() tview.Primitive {
	return dv.mainFlex
}

func (dv *DocumentsView) SetCollection(collectionName string) {
	dv.collection = collectionName
	dv.page = 1
	dv.selectedIdx = 0
	dv.refreshDocuments()
}

func (dv *DocumentsView) refreshDocuments() {
	dv.list.Clear()
	dv.documents = []map[string]interface{}{}

	if dv.collection == "" {
		dv.details.SetText("[yellow]No collection selected[white]")
		dv.updateStatusBar()
		return
	}

	searchParams := &api.SearchCollectionParams{
		Q:       "*",
		QueryBy: "*",
		PerPage: &dv.perPage,
		Page:    &dv.page,
	}

	result, err := dv.app.typesenseClient.Collection(dv.collection).Documents().Search(searchParams)
	if err != nil {
		dv.details.SetText(fmt.Sprintf("[red]Error loading documents: %v[white]", err))
		dv.updateStatusBar()
		return
	}

	hits := *result.Hits
	for _, hit := range hits {
		doc := *hit.Document
		dv.documents = append(dv.documents, doc)

		title := "Document"
		if id, ok := doc["id"].(string); ok {
			title = id
		} else if titleField, ok := doc["title"].(string); ok {
			title = titleField
		}

		secondary := ""
		if len(doc) > 0 {
			keys := make([]string, 0, len(doc))
			for k := range doc {
				keys = append(keys, k)
			}
			secondary = fmt.Sprintf("%d fields", len(keys))
		}

		dv.list.AddItem(title, secondary, 0, nil)
	}

	if len(dv.documents) > 0 {
		dv.onSelectionChange(0, "", "", rune(0))
	} else {
		dv.details.SetText("[yellow]No documents found[white]")
	}

	dv.updateStatusBar()
}

func (dv *DocumentsView) onSelectionChange(index int, mainText, secondaryText string, shortcut rune) {
	if index < 0 || index >= len(dv.documents) {
		return
	}

	dv.selectedIdx = index
	doc := dv.documents[index]

	var details strings.Builder
	details.WriteString("[yellow]Document:[white]\n\n")

	for key, value := range doc {
		valueStr := fmt.Sprintf("%v", value)
		if len(valueStr) > 200 {
			valueStr = valueStr[:200] + "..."
		}
		details.WriteString(fmt.Sprintf("[yellow]%s:[white] %s\n", key, valueStr))
	}

	dv.details.SetText(details.String())
	dv.updateStatusBar()
}

func (dv *DocumentsView) onSelect(index int, mainText, secondaryText string, shortcut rune) {
}

func (dv *DocumentsView) handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'r':
			dv.refreshDocuments()
			return nil
		case 'n':
			if len(dv.documents) == dv.perPage {
				dv.page++
				dv.refreshDocuments()
			}
			return nil
		case 'p':
			if dv.page > 1 {
				dv.page--
				dv.refreshDocuments()
			}
			return nil
		}
	}
	return event
}

func (dv *DocumentsView) updateStatusBar() {
	help := "[::b]r[white]: refresh | [::b]n[white]: next page | [::b]p[white]: prev page | [::b]F1[white]: collections | [::b]F2[white]: metrics | [::b]F3[white]: documents | [::b]Ctrl+C[white]: quit"
	if dv.collection != "" {
		help = fmt.Sprintf("%s | Collection: %s | Page: %d", help, dv.collection, dv.page)
		if len(dv.documents) > 0 {
			help = fmt.Sprintf("%s | Document %d/%d", help, dv.selectedIdx+1, len(dv.documents))
		}
	}
	dv.statusBar.SetText(help)
}
