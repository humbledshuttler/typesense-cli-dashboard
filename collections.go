package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/typesense/typesense-go/typesense/api"
)

type CollectionsView struct {
	app         *App
	mainFlex    *tview.Flex
	list        *tview.List
	details     *tview.TextView
	statusBar   *tview.TextView
	collections []*api.CollectionResponse
	selectedIdx int
}

func NewCollectionsView(app *App) *CollectionsView {
	cv := &CollectionsView{
		app:         app,
		selectedIdx: 0,
	}

	cv.list = tview.NewList().
		SetSelectedFunc(cv.onSelect).
		SetChangedFunc(cv.onSelectionChange)

	cv.details = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true).
		SetScrollable(true)
	cv.details.SetTitle("Collection Details").SetBorder(true)

	cv.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	cv.mainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				AddItem(cv.list, 0, 1, true).
				AddItem(cv.details, 0, 2, false),
			0, 1, true,
		).
		AddItem(cv.statusBar, 1, 0, false)

	cv.list.SetTitle("Collections").SetBorder(true)
	cv.list.SetInputCapture(cv.handleInput)

	cv.refreshCollections()
	cv.updateStatusBar()

	return cv
}

func (cv *CollectionsView) Render() tview.Primitive {
	return cv.mainFlex
}

func (cv *CollectionsView) refreshCollections() {
	cv.list.Clear()
	cv.collections = []*api.CollectionResponse{}

	collections, err := cv.app.typesenseClient.Collections().Retrieve()
	if err != nil {
		cv.details.SetText(fmt.Sprintf("[red]Error loading collections: %v[white]", err))
		return
	}

	cv.collections = collections
	for _, coll := range collections {
		docCount := "?"
		if coll.NumDocuments != nil {
			docCount = fmt.Sprintf("%d", *coll.NumDocuments)
		}
		cv.list.AddItem(
			fmt.Sprintf("%s (%s docs)", coll.Name, docCount),
			fmt.Sprintf("Fields: %d", len(coll.Fields)),
			0,
			nil,
		)
	}

	if len(cv.collections) > 0 {
		cv.onSelectionChange(0, "", "", rune(0))
	}
}

func (cv *CollectionsView) onSelectionChange(index int, mainText, secondaryText string, shortcut rune) {
	if index < 0 || index >= len(cv.collections) {
		return
	}

	cv.selectedIdx = index
	coll := cv.collections[index]

	var details strings.Builder
	details.WriteString(fmt.Sprintf("[yellow]Name:[white] %s\n\n", coll.Name))

	if coll.NumDocuments != nil {
		details.WriteString(fmt.Sprintf("[yellow]Documents:[white] %d\n", *coll.NumDocuments))
	}

	details.WriteString(fmt.Sprintf("[yellow]Fields:[white] %d\n\n", len(coll.Fields)))

	details.WriteString("[yellow]Fields:\n[white]")
	for _, field := range coll.Fields {
		optional := ""
		if field.Optional != nil && *field.Optional {
			optional = " (optional)"
		}
		facet := ""
		if field.Facet != nil && *field.Facet {
			facet = " [facet]"
		}
		details.WriteString(fmt.Sprintf("  • %s: %s%s%s\n", field.Name, field.Type, optional, facet))
	}

	if coll.DefaultSortingField != nil {
		details.WriteString(fmt.Sprintf("\n[yellow]Default Sorting Field:[white] %s\n", *coll.DefaultSortingField))
	}

	cv.details.SetText(details.String())
	cv.updateStatusBar()
}

func (cv *CollectionsView) onSelect(index int, mainText, secondaryText string, shortcut rune) {
	if index < 0 || index >= len(cv.collections) {
		return
	}

	coll := cv.collections[index]
	cv.app.documentsView.SetCollection(coll.Name)
	cv.app.pages.SwitchToPage("documents")
}

func (cv *CollectionsView) handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'r':
			cv.refreshCollections()
			return nil
		}
	}
	return event
}

func (cv *CollectionsView) updateStatusBar() {
	help := "[::b]r[white]: refresh | [::b]Enter[white]: view documents | [::b]c[white]: collections | [::b]m[white]: metrics | [::b]d[white]: documents | [::b]x[white]: quit"
	if len(cv.collections) > 0 {
		help = fmt.Sprintf("%s | Collection %d/%d", help, cv.selectedIdx+1, len(cv.collections))
	}
	cv.statusBar.SetText(help)
}
