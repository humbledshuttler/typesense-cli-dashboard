package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MetricsView struct {
	app       *App
	mainFlex  *tview.Flex
	metrics   *tview.TextView
	statusBar *tview.TextView
	ticker    *time.Ticker
}

func NewMetricsView(app *App) *MetricsView {
	mv := &MetricsView{
		app: app,
	}

	mv.metrics = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true).
		SetScrollable(true)
	mv.metrics.SetTitle("Server Metrics").SetBorder(true)

	mv.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	mv.mainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mv.metrics, 0, 1, true).
		AddItem(mv.statusBar, 1, 0, false)

	mv.metrics.SetInputCapture(mv.handleInput)

	mv.refreshMetrics()
	mv.ticker = time.NewTicker(5 * time.Second)
	go mv.autoRefresh()

	return mv
}

func (mv *MetricsView) Render() tview.Primitive {
	return mv.mainFlex
}

func (mv *MetricsView) refreshMetrics() {
	var metrics strings.Builder

	health, err := mv.app.typesenseClient.Health(5 * time.Second)
	if err != nil {
		metrics.WriteString(fmt.Sprintf("[red]Error loading health: %v[white]\n", err))
	} else {
		metrics.WriteString("[yellow]Health:[white]\n")
		if health {
			metrics.WriteString("  Status: [green]OK[white]\n")
		} else {
			metrics.WriteString("  Status: [red]UNHEALTHY[white]\n")
		}
	}

	collections, err := mv.app.typesenseClient.Collections().Retrieve()
	if err == nil {
		metrics.WriteString(fmt.Sprintf("\n[yellow]Collections:[white] %d\n", len(collections)))
		var totalDocs int64 = 0
		for _, coll := range collections {
			if coll.NumDocuments != nil {
				totalDocs += *coll.NumDocuments
			}
		}
		metrics.WriteString(fmt.Sprintf("[yellow]Total Documents:[white] %d\n", totalDocs))
	}

	metrics.WriteString(fmt.Sprintf("\n[yellow]Last Updated:[white] %s\n", time.Now().Format("2006-01-02 15:04:05")))

	mv.metrics.SetText(metrics.String())
	mv.updateStatusBar()
}

func (mv *MetricsView) autoRefresh() {
	for range mv.ticker.C {
		mv.app.app.QueueUpdateDraw(func() {
			mv.refreshMetrics()
		})
	}
}

func (mv *MetricsView) handleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'r':
			mv.refreshMetrics()
			return nil
		}
	}
	return event
}

func (mv *MetricsView) updateStatusBar() {
	help := "[::b]r[white]: refresh | [::b]c[white]: collections | [::b]m[white]: metrics | [::b]d[white]: documents | [::b]x[white]: quit"
	mv.statusBar.SetText(help)
}
