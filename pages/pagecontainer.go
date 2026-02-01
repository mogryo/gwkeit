package pages

import (
	"slices"

	"github.com/gwkeit/configuration"
	"github.com/gwkeit/globaldeps"
	"github.com/gwkeit/pages/additionpage"
	"github.com/gwkeit/pages/allsnippetspage"
	"github.com/gwkeit/pages/editpage"
	"github.com/gwkeit/pages/searchpage"
	"github.com/gwkeit/widgets"
)

type PageContainer struct {
	searchPage      *searchpage.SearchPage
	additionPage    *additionpage.AdditionPage
	editPage        *editpage.EditPage
	allSnippetsPage *allsnippetspage.AllSnippetsPage
	logs            *widgets.LogsWidget
	globalDeps      *globaldeps.GlobalDependencies
}

func NewPageContainer(
	globalDeps *globaldeps.GlobalDependencies,
	appState *configuration.AppConfiguration,
) *PageContainer {
	logs := widgets.NewLogsWidget(globalDeps.App)
	additionPage := additionpage.NewPage(globalDeps, logs)
	searchPage := searchpage.NewPage(globalDeps, &appState.SearchPage, logs)
	editPage := editpage.NewPage(globalDeps, logs)
	allSnippetsPage := allsnippetspage.NewPage(globalDeps, logs, &appState.AllSnippets)

	go func() {
		for {
			select {
			case payload := <-globalDeps.Chan:
				switch payload.PageName {
				case editpage.PageName:
					globalDeps.App.QueueUpdateDraw(func() {
						editPage.SwitchToEditPage(payload.SnippetId)
					})
				case additionpage.PageName:
					globalDeps.App.QueueUpdateDraw(func() {
						additionPage.SwitchToAdditionPage()
					})
				case searchpage.PageName:
					globalDeps.App.QueueUpdateDraw(func() {
						searchPage.SwitchToSearchPage()
					})
				case allsnippetspage.PageName:
					globalDeps.App.QueueUpdateDraw(func() {
						allSnippetsPage.SwitchToSnippetListPage()
					})
				}
			case <-globalDeps.Ctx.Done():
				return
			}
		}
	}()

	pc := &PageContainer{
		additionPage:    additionPage,
		searchPage:      searchPage,
		editPage:        editPage,
		allSnippetsPage: allSnippetsPage,
		globalDeps:      globalDeps,
	}

	globalDeps.Pages.AddPage(additionpage.PageName, pc.additionPage.Frame, true, false)
	globalDeps.Pages.AddPage(searchpage.PageName, pc.searchPage.Frame, true, false)
	globalDeps.Pages.AddPage(editpage.PageName, pc.editPage.Frame, true, false)
	globalDeps.Pages.AddPage(allsnippetspage.PageName, pc.allSnippetsPage.Frame, true, false)

	return pc
}

func (pc *PageContainer) FocusSearchPageSearchField() {
	activePages := pc.globalDeps.Pages.GetPageNames(true)
	if slices.Contains(activePages, searchpage.PageName) {
		pc.searchPage.FocusSearchField()
	}
}
