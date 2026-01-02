package appui

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/repository"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
)

type AppUI struct {
	App          *tview.Application
	pages        *tview.Pages
	searchPage   *SearchPageUI
	additionPage *AdditionPageUI
	logs         *widgets.LogsWidget
	repo         *repository.Repository
	ctx          context.Context
}

func (aui *AppUI) SwitchToSearchPage() {
	aui.pages.SwitchToPage("Main")
}

func (aui *AppUI) SwitchToAdditionPage() {
	aui.pages.SwitchToPage("Addition")
}

func New(ctx context.Context, repo *repository.Repository) *AppUI {
	app := tview.NewApplication()
	aui := &AppUI{
		App:   app,
		logs:  widgets.NewLogsWidget(app),
		pages: tview.NewPages(),
		repo:  repo,
		ctx:   ctx,
	}

	aui.NewAdditionPage()
	aui.NewSearchPage()

	aui.App.SetRoot(aui.pages, true).EnableMouse(true)
	aui.App.SetFocus(aui.searchPage.searchField)

	aui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event
		switch event.Key() {
		case tcell.KeyCtrlA:
			aui.SwitchToAdditionPage()
			aui.App.SetFocus(aui.additionPage.body)
			resultEvent = nil
		case tcell.KeyCtrlW:
			aui.SwitchToSearchPage()
			aui.App.SetFocus(aui.searchPage.searchField)
			resultEvent = nil
		case tcell.KeyCtrlQ:
			aui.App.Stop()
		case tcell.KeyCtrlC:
			resultEvent = tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)
		}
		return resultEvent
	})

	return aui
}

func (aui *AppUI) Run() error {
	return aui.App.Run()
}
