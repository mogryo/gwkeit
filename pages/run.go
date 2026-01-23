package pages

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/globaldeps"
	"github.com/gwkeit/repository"
	"github.com/rivo/tview"
)

func Run(ctx context.Context, repo *repository.Repository) error {
	app := tview.NewApplication()
	globalDeps := globaldeps.New(ctx, app, repo)
	pageContainer := NewPageContainer(globalDeps)

	globalDeps.App.SetRoot(globalDeps.Pages, true).EnableMouse(true)
	pageContainer.searchPage.SwitchToSearchPage()
	globalDeps.App.SetFocus(pageContainer.searchPage.searchField)

	globalDeps.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event
		switch event.Key() {
		case tcell.KeyCtrlA:
			pageContainer.additionPage.SwitchToAdditionPage()
			resultEvent = nil
		case tcell.KeyCtrlW:
			pageContainer.searchPage.SwitchToSearchPage()
			resultEvent = nil
		case tcell.KeyCtrlQ:
			globalDeps.App.Stop()
		case tcell.KeyCtrlC:
			resultEvent = tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)
			return resultEvent
		case tcell.KeyCtrlK:
			pageContainer.allSnippetsPage.SwitchToSnippetListPage()
			resultEvent = nil
		}

		return resultEvent
	})

	return globalDeps.App.Run()
}
