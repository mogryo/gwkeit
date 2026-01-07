package pages

import (
	"github.com/gwkeit/globaldeps"
	"github.com/gwkeit/widgets"
)

type PageContainer struct {
	searchPage   *SearchPage
	additionPage *AdditionPage
	editPage     *EditPage
	logs         *widgets.LogsWidget
}

func NewPageContainer(globalDeps *globaldeps.GlobalDependencies) *PageContainer {
	logs := widgets.NewLogsWidget(globalDeps.App)
	additionPage := NewAdditionPage(globalDeps, logs)
	searchPage := NewSearchPage(globalDeps, logs)
	editPage := NewEditPage(globalDeps, logs)

	go func() {
		for {
			select {
			case payload := <-globalDeps.Chan:
				if payload.PageName == "Edit" {
					globalDeps.App.QueueUpdateDraw(func() {
						editPage.SwitchToEditPage(payload.SnippetId, globalDeps)
					})
				}
				if payload.PageName == "Addition" {
					globalDeps.App.QueueUpdateDraw(func() {
						additionPage.SwitchToAdditionPage(globalDeps)
					})
				}
				if payload.PageName == "Main" {
					globalDeps.App.QueueUpdateDraw(func() {
						searchPage.SwitchToSearchPage(globalDeps)
					})
				}
			}
		}
	}()

	pc := &PageContainer{
		additionPage: additionPage,
		searchPage:   searchPage,
		editPage:     editPage,
	}

	globalDeps.Pages.AddPage("Addition", pc.additionPage.frame, true, false)
	globalDeps.Pages.AddPage("Main", pc.searchPage.frame, true, false)
	globalDeps.Pages.AddPage("Edit", pc.editPage.frame, true, false)

	return pc
}
