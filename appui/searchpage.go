package appui

import (
	"context"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/gwkeitdb"
	"github.com/gwkeit/repository"
	"github.com/gwkeit/transform"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

var (
	shortcutRunes = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
)

type CurrentViewData struct {
	snippets []gwkeitdb.Snippet
	tags     []gwkeitdb.Tag
	urls     []gwkeitdb.Url
}

type SearchPageUI struct {
	body        *tview.TextArea
	searchField *tview.InputField
	resultList  *tview.List
	description *tview.TextArea
	urls        *tview.TextArea
	title       *tview.TextArea
	grid        *tview.Grid
	currentData *CurrentViewData
}

func (aui *AppUI) NewSearchPage() {
	aui.searchPage = &SearchPageUI{
		currentData: &CurrentViewData{},
	}

	aui.searchPage.initMetadataFields()
	aui.searchPage.initBody()
	aui.searchPage.initSearchField(aui.ctx, aui.App, aui.repo)
	aui.searchPage.initResultList(aui.ctx, aui.repo)
	aui.searchPage.initGridLayout(aui.logs.View)
	aui.searchPage.initInputCapture(aui.App, aui.logs)

	aui.pages.AddPage("Main", aui.searchPage.grid, true, true)
}

func (sp *SearchPageUI) initMetadataFields() {
	sp.title = uibuilder.NewTextArea("", "")
	sp.title.SetDisabled(true)

	sp.description = uibuilder.NewTextArea("", "")
	sp.description.SetDisabled(true)

	sp.urls = uibuilder.NewTextArea("", "")
	sp.urls.SetDisabled(true)
}

func (sp *SearchPageUI) initBody() {
	sp.body = uibuilder.NewTextArea("", "")
	sp.body.SetDisabled(true)
}

func (sp *SearchPageUI) initSearchField(
	ctx context.Context,
	app *tview.Application,
	repo *repository.Repository,
) {
	sp.searchField = tview.NewInputField().
		SetLabel("").
		SetFieldWidth(0).
		SetAcceptanceFunc(func(text string, keyCode rune) bool { return true }).
		SetDoneFunc(func(key tcell.Key) {
			if sp.resultList.GetItemCount() == 0 {
				return
			}

			app.SetFocus(sp.resultList)
			index := sp.resultList.GetCurrentItem()
			onSelect := sp.resultList.GetSelectedFunc()
			mainText, secText := sp.resultList.GetItemText(index)
			onSelect(index, mainText, secText, shortcutRunes[index])
		}).
		SetChangedFunc(func(text string) {
			sp.currentData.snippets = repo.FindSnippets(ctx, strings.Split(text, " "))
			sp.resultList.Clear()
			for i, snippet := range sp.currentData.snippets {
				sp.resultList.AddItem(snippet.Title, strconv.FormatInt(snippet.ID, 10), shortcutRunes[i], nil)
			}
		})
	sp.searchField.SetFieldStyle(uibuilder.InputBackgroundStyle)
	sp.searchField.SetBackgroundColor(tcell.ColorDefault)
}

func (sp *SearchPageUI) initResultList(ctx context.Context, repo *repository.Repository) {
	sp.resultList = tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			id, err := strconv.ParseInt(secondaryText, 10, 64)

			if err != nil {
				panic(err)
			}

			title := sp.currentData.snippets[index].Title
			body := sp.currentData.snippets[index].Body
			tags := repo.FindSnippetTags(ctx, id)
			urls := repo.FindSnippetUrls(ctx, id)

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

func (sp *SearchPageUI) initGridLayout(logsView *tview.TextView) {
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

func (sp *SearchPageUI) initInputCapture(app *tview.Application, logsWidget *widgets.LogsWidget) {
	sp.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event

		switch event.Key() {
		case tcell.KeyCtrlF:
			app.SetFocus(sp.searchField)
			resultEvent = nil
		case tcell.KeyCtrlL:
			app.SetFocus(sp.resultList)
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
		}

		return resultEvent
	})
}
