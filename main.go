package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/gdamore/tcell/v2"
	"github.com/go-errors/errors"
	"github.com/integrii/flaggy"
	"github.com/rivo/tview"
	"github.com/typesense/typesense-go/typesense"
)

const DEFAULT_VERSION = "unversioned"

var (
	commit      string
	version     = DEFAULT_VERSION
	date        string
	buildSource = "unknown"

	hostFlag     = "localhost"
	portFlag     = "8108"
	protocolFlag = "http"
	apiKeyFlag   = "xyz"
	debugFlag    = false
)

type App struct {
	typesenseClient *typesense.Client
	app             *tview.Application
	pages           *tview.Pages
	collectionsView *CollectionsView
	metricsView     *MetricsView
	documentsView   *DocumentsView
}

func main() {
	updateBuildInfo()

	info := fmt.Sprintf(
		"%s\nDate: %s\nBuildSource: %s\nCommit: %s\nOS: %s\nArch: %s",
		version,
		date,
		buildSource,
		commit,
		runtime.GOOS,
		runtime.GOARCH,
	)

	flaggy.SetName("typesense-cli")
	flaggy.SetDescription("A terminal UI for browsing Typesense collections, viewing metrics, and managing documents")
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/humbledshuttler/typesense-cli"

	flaggy.String(&hostFlag, "", "host", "Typesense server host (default: localhost)")
	flaggy.String(&portFlag, "p", "port", "Typesense server port (default: 8108)")
	flaggy.String(&protocolFlag, "", "protocol", "Typesense server protocol (default: http)")
	flaggy.String(&apiKeyFlag, "k", "api-key", "Typesense API key (default: xyz)")
	flaggy.Bool(&debugFlag, "d", "debug", "Enable debug mode")
	flaggy.SetVersion(info)

	flaggy.Parse()

	client := typesense.NewClient(
		typesense.WithServer(fmt.Sprintf("%s://%s:%s", protocolFlag, hostFlag, portFlag)),
		typesense.WithAPIKey(apiKeyFlag),
	)

	app := &App{
		typesenseClient: client,
		app:             tview.NewApplication(),
		pages:           tview.NewPages(),
	}

	app.collectionsView = NewCollectionsView(app)
	app.metricsView = NewMetricsView(app)
	app.documentsView = NewDocumentsView(app)

	app.setupPages()

	if err := app.app.SetRoot(app.pages, true).SetFocus(app.pages).Run(); err != nil {
		newErr := errors.Wrap(err, 0)
		stackTrace := newErr.ErrorStack()

		if debugFlag {
			log.Fatalf("Error running application:\n\n%s", stackTrace)
		} else {
			log.Fatalf("Error running application: %v", err)
		}
		os.Exit(1)
	}
}

func updateBuildInfo() {
	if version == DEFAULT_VERSION {
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range buildInfo.Settings {
				switch setting.Key {
				case "vcs.revision":
					commit = setting.Value
					if len(commit) > 7 {
						version = commit[:7]
					} else {
						version = commit
					}
				case "vcs.time":
					date = setting.Value
				}
			}
		}
	}
}

func (a *App) setupPages() {
	a.pages.AddPage("collections", a.collectionsView.Render(), true, true)
	a.pages.AddPage("metrics", a.metricsView.Render(), true, false)
	a.pages.AddPage("documents", a.documentsView.Render(), true, false)

	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			a.app.Stop()
			return nil
		case tcell.KeyF1:
			a.pages.SwitchToPage("collections")
			return nil
		case tcell.KeyF2:
			a.pages.SwitchToPage("metrics")
			return nil
		case tcell.KeyF3:
			a.pages.SwitchToPage("documents")
			return nil
		}
		return event
	})
}
