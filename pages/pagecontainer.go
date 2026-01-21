package pages

import (
	"github.com/gwkeit/globaldeps"
	"github.com/gwkeit/widgets"
)

type PageContainer struct {
	searchPage      *SearchPage
	additionPage    *AdditionPage
	editPage        *EditPage
	allSnippetsPage *AllSnippetsPage
	logs            *widgets.LogsWidget
}

func NewPageContainer(globalDeps *globaldeps.GlobalDependencies) *PageContainer {
	logs := widgets.NewLogsWidget(globalDeps.App)
	additionPage := NewAdditionPage(globalDeps, logs)
	searchPage := NewSearchPage(globalDeps, logs)
	editPage := NewEditPage(globalDeps, logs)
	allSnippetsPage := NewAllSnippetsPage(globalDeps, logs)

	go func() {
		for {
			select {
			case payload := <-globalDeps.Chan:
				switch payload.PageName {
				case "Edit":
					globalDeps.App.QueueUpdateDraw(func() {
						editPage.SwitchToEditPage(payload.SnippetId)
					})
				case "Addition":
					globalDeps.App.QueueUpdateDraw(func() {
						additionPage.SwitchToAdditionPage()
					})
				case "Main":
					globalDeps.App.QueueUpdateDraw(func() {
						searchPage.SwitchToSearchPage()
					})
				case "AllSnippets":
					globalDeps.App.QueueUpdateDraw(func() {
						allSnippetsPage.SwitchToSnippetListPage()
					})
				}
			}
		}
	}()

	pc := &PageContainer{
		additionPage:    additionPage,
		searchPage:      searchPage,
		editPage:        editPage,
		allSnippetsPage: allSnippetsPage,
	}

	globalDeps.Pages.AddPage("Addition", pc.additionPage.frame, true, false)
	globalDeps.Pages.AddPage("Main", pc.searchPage.frame, true, false)
	globalDeps.Pages.AddPage("Edit", pc.editPage.frame, true, false)
	globalDeps.Pages.AddPage("AllSnippets", pc.allSnippetsPage.frame, true, false)

	return pc
}
