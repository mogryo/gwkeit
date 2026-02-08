package pages

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/apptools"
	"github.com/gwkeit/configuration"
	"github.com/gwkeit/repository"
	"github.com/rivo/tview"
)

func Run(ctx context.Context, repo *repository.Repository, initialState *configuration.AppConfiguration) error {
	app := tview.NewApplication()
	tools := apptools.New(ctx, app, repo)
	pageContainer := NewPageContainer(tools, initialState)

	app.SetRoot(pageContainer.Pages, true).EnableMouse(true)
	pageContainer.searchPage.SwitchToPage()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event
		switch event.Key() {
		case tcell.KeyCtrlA:
			tools.GoToPage(configuration.AdditionPage, nil)
			resultEvent = nil
		case tcell.KeyCtrlW:
			tools.GoToPage(configuration.SearchPage, nil)
			resultEvent = nil
		case tcell.KeyCtrlQ:
			app.Stop()
		case tcell.KeyCtrlC:
			resultEvent = tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)
			return resultEvent
		case tcell.KeyCtrlK:
			tools.GoToPage(configuration.AllSnippetsPage, nil)
			resultEvent = nil
		}

		return resultEvent
	})

	return app.Run()
}
