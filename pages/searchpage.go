package pages

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/globaldeps"
	"github.com/gwkeit/transform"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

var (
	shortcutRunes = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
)

type SearchPage struct {
	body              *tview.TextArea
	searchField       *tview.InputField
	resultList        *tview.List
	description       *tview.TextArea
	urls              *tview.TextArea
	title             *tview.TextArea
	grid              *tview.Grid
	frame             *tview.Frame
	selectedSnippetId int64
}

func NewSearchPage(
	globalDeps *globaldeps.GlobalDependencies,
	logs *widgets.LogsWidget,
) *SearchPage {
	searchPage := &SearchPage{
		selectedSnippetId: -1,
	}

	searchPage.initMetadataFields()
	searchPage.initBody()
	searchPage.initSearchField(globalDeps)
	searchPage.initResultList(globalDeps)
	searchPage.initGridLayout(logs.View)
	searchPage.initInputCapture(globalDeps, logs)
	searchPage.initFrame()

	return searchPage
}

func (sp *SearchPage) initMetadataFields() {
	sp.title = uibuilder.NewTextArea("", "")
	sp.title.SetDisabled(true)

	sp.description = uibuilder.NewTextArea("", "")
	sp.description.SetDisabled(true)

	sp.urls = uibuilder.NewTextArea("", "")
	sp.urls.SetDisabled(true)
}

func (sp *SearchPage) initBody() {
	sp.body = uibuilder.NewTextArea("", "")
	sp.body.SetDisabled(true)
}

func (sp *SearchPage) initSearchField(
	globalDeps *globaldeps.GlobalDependencies,
) {
	sp.searchField = tview.NewInputField().
		SetLabel("").
		SetFieldWidth(0).
		SetAcceptanceFunc(func(text string, keyCode rune) bool { return true }).
		SetDoneFunc(func(key tcell.Key) {
			if sp.resultList.GetItemCount() == 0 {
				return
			}

			globalDeps.App.SetFocus(sp.resultList)
			index := sp.resultList.GetCurrentItem()
			onSelect := sp.resultList.GetSelectedFunc()
			mainText, secText := sp.resultList.GetItemText(index)
			onSelect(index, mainText, secText, shortcutRunes[index])
		}).
		SetChangedFunc(func(text string) {
			foundSnippets := globalDeps.Repo.FindSnippets(globalDeps.Ctx, strings.Split(text, " "))
			sp.resultList.Clear()
			for i, snippet := range foundSnippets {
				sp.resultList.AddItem(snippet.Title, strconv.FormatInt(snippet.ID, 10), shortcutRunes[i], nil)
			}
		})
	sp.searchField.SetFieldStyle(uibuilder.InputBackgroundStyle)
	sp.searchField.SetBackgroundColor(tcell.ColorDefault)
}

func (sp *SearchPage) initResultList(globalDeps *globaldeps.GlobalDependencies) {
	sp.resultList = tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			id, err := strconv.ParseInt(secondaryText, 10, 64)
			sp.selectedSnippetId = id

			if err != nil {
				panic(err)
			}

			snippet := globalDeps.Repo.FindSnippet(globalDeps.Ctx, id)
			title := snippet.Title
			body := snippet.Body
			tags := globalDeps.Repo.FindSnippetTags(globalDeps.Ctx, id)
			urls := globalDeps.Repo.FindSnippetUrls(globalDeps.Ctx, id)

			sp.title.SetText(title, true)
			sp.body.SetText(body, true)
			sp.description.SetText(transform.TagListToFormDescription(tags), true)
			sp.urls.SetText(transform.UrlListToFormUrls(urls), true)
		}).
		SetSelectedFocusOnly(true)
	sp.resultList.SetBackgroundColor(tcell.ColorDefault)
	sp.resultList.SetMainTextStyle(uibuilder.InputBackgroundStyle)
	sp.resultList.SetShortcutStyle(uibuilder.InputBackgroundStyle.Foreground(tcell.ColorGreen))
}

func (sp *SearchPage) initGridLayout(logsView *tview.TextView) {
	sp.grid = tview.NewGrid().
		SetRows(3, 11).
		SetColumns(0, 50).
		SetBorders(false).
		AddItem(uibuilder.NewWidget("[ctr+f] Search:", sp.searchField), 0, 0, 1, 1, 0, 0, false).
		AddItem(uibuilder.NewWidget("Logs:", logsView), 0, 1, 2, 1, 0, 0, false).
		AddItem(uibuilder.NewWidget("[ctr+l] List:", sp.resultList), 1, 0, 1, 1, 0, 0, false)
	metadataFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget("Title:", sp.title), 0, 1, false).
		AddItem(uibuilder.NewWidget("Description:", sp.description), 0, 3, false).
		AddItem(uibuilder.NewWidget("URLs:", sp.urls), 0, 3, false)

	sp.grid.AddItem(uibuilder.NewWidget("Code:", sp.body), 2, 0, 1, 1, 0, 100, false).
		AddItem(metadataFlex, 2, 1, 1, 1, 0, 100, false)

	sp.grid.SetBackgroundColor(tcell.ColorDefault)
}

func (sp *SearchPage) initInputCapture(
	globalDeps *globaldeps.GlobalDependencies,
	logsWidget *widgets.LogsWidget,
) {
	sp.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event

		switch event.Key() {
		case tcell.KeyCtrlF:
			globalDeps.App.SetFocus(sp.searchField)
			resultEvent = nil
		case tcell.KeyCtrlL:
			globalDeps.App.SetFocus(sp.resultList)
			resultEvent = nil
		case tcell.KeyCtrlC:
			bodyText := sp.body.GetText()
			if strings.TrimSpace(bodyText) == "" {
				logsWidget.AddErrorLogs([]string{"Body is empty. Nothing to copy."})
			} else {
				logsWidget.AddSuccessLogs([]string{"Body copied to clipboard."})
				clipboard.Write(clipboard.FmtText, []byte(sp.body.GetText()))
			}
			resultEvent = nil
		case tcell.KeyCtrlE:
			if sp.selectedSnippetId > -1 {
				globalDeps.GoToEditPage(sp.selectedSnippetId)
			} else {
				logsWidget.AddErrorLogs([]string{"No snippet selected."})
			}
			resultEvent = nil
		}

		return resultEvent
	})
}

func (sp *SearchPage) initFrame() {
	sp.frame = tview.NewFrame(sp.grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Search", true, tview.AlignCenter, tcell.ColorWhite)
	sp.frame.SetBackgroundColor(tcell.ColorDefault)
}

func (sp *SearchPage) showSnippet(snippetId int64, globalDeps *globaldeps.GlobalDependencies) {
	snippet := globalDeps.Repo.FindSnippet(globalDeps.Ctx, snippetId)
	title := snippet.Title
	body := snippet.Body
	tags := globalDeps.Repo.FindSnippetTags(globalDeps.Ctx, snippetId)
	urls := globalDeps.Repo.FindSnippetUrls(globalDeps.Ctx, snippetId)

	sp.title.SetText(title, true)
	sp.body.SetText(body, true)
	sp.description.SetText(transform.TagListToFormDescription(tags), true)
	sp.urls.SetText(transform.UrlListToFormUrls(urls), true)
}

func (sp *SearchPage) SwitchToSearchPage(globalDeps *globaldeps.GlobalDependencies) {
	if sp.selectedSnippetId > -1 {
		sp.showSnippet(sp.selectedSnippetId, globalDeps)
	}
	globalDeps.Pages.SwitchToPage("Main")
	globalDeps.App.SetFocus(sp.searchField)
}
