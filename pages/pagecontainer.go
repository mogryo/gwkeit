package pages

import (
	"github.com/gwkeit/apptools"
	"github.com/gwkeit/configuration"
	"github.com/gwkeit/pages/additionpage"
	"github.com/gwkeit/pages/allsnippetspage"
	"github.com/gwkeit/pages/editpage"
	"github.com/gwkeit/pages/searchpage"
	"github.com/gwkeit/pages/shortcutmodal"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
)

type PageContainer struct {
	searchPage      *searchpage.SearchPage
	additionPage    *additionpage.AdditionPage
	editPage        *editpage.EditPage
	allSnippetsPage *allsnippetspage.AllSnippetsPage
	logs            *widgets.LogsWidget
	modal           *shortcutmodal.ShortcutModal
	Pages           *tview.Pages
}

type PageSwitcher struct{}

func NewPageContainer(
	tools *apptools.Tools,
	appState *configuration.AppConfiguration,
) *PageContainer {
	logs := widgets.NewLogsWidget(tools)
	additionPage := additionpage.NewPage(tools, logs)
	searchPage := searchpage.NewPage(tools, &appState.SearchPage, logs)
	editPage := editpage.NewPage(tools, logs)
	allSnippetsPage := allsnippetspage.NewPage(tools, logs, &appState.AllSnippets)
	modalPage := shortcutmodal.NewModal(tools)

	pc := &PageContainer{
		additionPage:    additionPage,
		searchPage:      searchPage,
		editPage:        editPage,
		allSnippetsPage: allSnippetsPage,
		modal:           modalPage,
		Pages:           tview.NewPages(),
	}
	tools.RegisterGoToPage(pc.goToPage)
	tools.RegisterSwitchToPage(pc.Pages.SwitchToPage)
	tools.RegisterShowPage(pc.Pages.ShowPage)
	tools.RegisterHidePage(pc.Pages.HidePage)
	tools.RegisterGetFrontPage(pc.Pages.GetFrontPage)

	pc.Pages.AddPage(configuration.AdditionPage.String(), pc.additionPage.Frame, true, false)
	pc.Pages.AddPage(configuration.SearchPage.String(), pc.searchPage.Frame, true, false)
	pc.Pages.AddPage(configuration.EditPage.String(), pc.editPage.Frame, true, false)
	pc.Pages.AddPage(configuration.AllSnippetsPage.String(), pc.allSnippetsPage.Frame, true, false)
	pc.Pages.AddPage(configuration.ShortcutModal.String(), pc.modal.Frame, true, true)

	return pc
}

func (pc *PageContainer) goToPage(pageName configuration.PageName, payload any) {
	switch pageName {
	case configuration.AdditionPage:
		pc.additionPage.SwitchToPage()
	case configuration.AllSnippetsPage:
		pc.allSnippetsPage.SwitchToPage()
	case configuration.EditPage:
		pc.editPage.SwitchToPage(payload.(int64))
	case configuration.SearchPage:
		pc.searchPage.SwitchToPage()
	case configuration.ShortcutModal:
		pc.modal.SwitchToPage(payload.([]apptools.ShortcutDescription))
	}
}
