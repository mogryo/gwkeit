package searchpage

import (
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/apptools"
	"github.com/gwkeit/configuration"
	"github.com/gwkeit/gwkeitdb"
	"github.com/gwkeit/slicelib"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/utils"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

var shortcutDescription = []apptools.ShortcutDescription{
	{Key: "ctrl+F", Description: "Focus search field"},
	{Key: "ctrl+L", Description: "Focus result list"},
	{Key: "ctrl+O", Description: "Focus search type"},
	{Key: "ctrl+C", Description: "Copy body"},
	{Key: "ctrl+E/Enter", Description: "Edit snippet"},
	{"}", "Show next result page"},
	{"{", "Show previous result page"},
}

func (sp *SearchPage) initMetadataFields() {
	sp.title = uibuilder.NewTextArea(sp.themeName, "", "")
	sp.title.SetDisabled(true)

	sp.description = uibuilder.NewTextArea(sp.themeName, "", "")
	sp.description.SetDisabled(true)

	sp.urls = uibuilder.NewTextArea(sp.themeName, "", "")
	sp.urls.SetDisabled(true)

	sp.language = uibuilder.NewDropDown(sp.themeName, "")
	sp.language.SetOptions(slices.Concat([]string{""}, configuration.LanguagesStrings), nil)
	sp.language.SetDisabled(true)
}

func (sp *SearchPage) initBody() {
	sp.body = uibuilder.NewTextArea(sp.themeName, "", "")
	sp.body.SetDisabled(true)
}

func (sp *SearchPage) initSearchField() {
	executeSearch := func(text string) {
		splitConditions := strings.Split(strings.TrimSpace(text), " ")
		filteredConditions := slicelib.Filter(splitConditions, func(condition string) bool { return condition != "" })

		sp.resultList.Clear()
		if len(filteredConditions) > 0 {
			sp.foundSnippets = sp.searchCallback(sp.tools.Ctx, filteredConditions)
		} else {
			sp.foundSnippets = []gwkeitdb.Snippet{}
		}

		if len(sp.foundSnippets) == 0 {
			sp.selectedSnippetId = -1
			sp.clearMetadataFields()
		}

		sp.totalFoundAmount = int64(len(sp.foundSnippets))
		sp.totalFoundView.SetText(strconv.FormatInt(sp.totalFoundAmount, 10))
		sp.setResultListPage(1)
	}

	sp.searchField = uibuilder.NewInputField(sp.themeName, "", "").
		SetFieldWidth(0).
		SetAcceptanceFunc(func(text string, keyCode rune) bool {
			return unicode.IsLetter(keyCode) || unicode.IsDigit(keyCode) || unicode.IsSpace(keyCode)
		}).
		SetDoneFunc(func(key tcell.Key) {
			if sp.resultList.GetItemCount() > 0 {
				sp.tools.Focus(sp.resultList)
			}

		}).
		SetChangedFunc(executeSearch)

	typeIndex := slices.Index([]string{"Tags", "Like", "FTS"}, sp.pageConf.GetSearchType())
	sp.searchType = uibuilder.NewDropDown(sp.themeName, "Type: ").
		SetOptions([]string{"Tags", "Like", "FTS"}, func(text string, index int) {
			switch text {
			case "Tags":
				sp.searchCallback = sp.tools.Repo.FindSnippetsByTags
			case "Like":
				sp.searchCallback = sp.tools.Repo.FindSnippetsByLikeTags
			case "FTS":
				sp.searchCallback = sp.tools.Repo.FindSnippetsByFts
			}
			executeSearch(sp.searchField.GetText())
		}).
		SetCurrentOption(utils.IfElse(typeIndex > -1, typeIndex, 0))

	sp.searchBox = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(sp.searchField, 0, 4, false).
		AddItem(sp.searchType, 0, 1, false)

	sp.searchBox.SetBorderPadding(0, 0, 0, 0).SetBackgroundColor(tcell.ColorDefault)
}

func (sp *SearchPage) initResultList() {
	sp.resultList = uibuilder.NewList(sp.themeName).
		ShowSecondaryText(false).
		SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			id, err := strconv.ParseInt(secondaryText, 10, 64)
			sp.selectedSnippetId = id

			if err != nil {
				panic(err)
			}

			snippet := sp.tools.Repo.FindSnippet(sp.tools.Ctx, id)

			sp.title.SetText(snippet.Title, true)
			sp.body.SetText(snippet.Body, true)
			sp.description.SetText(snippet.Description, true)
			sp.urls.SetText(snippet.Url, true)

			if snippet.Language.Valid {
				idx := slices.Index(configuration.LanguagesStrings, snippet.Language.String)
				sp.language.SetCurrentOption(idx + 1)
			} else {
				sp.language.SetCurrentOption(0)
			}
		})
	sp.totalFoundView = uibuilder.NewTextView(sp.themeName, strconv.FormatInt(sp.totalFoundAmount, 10))
	sp.currentPageView = uibuilder.NewTextView(sp.themeName, strconv.FormatInt(sp.currentPage, 10))
}

func (sp *SearchPage) initGridLayout() {
	sp.grid = uibuilder.NewGrid(sp.themeName).
		SetRows(3, 11).
		SetColumns(0, 50).
		AddItem(uibuilder.NewWidget(sp.themeName, "Search:", sp.searchBox), 0, 0, 1, 1, 0, 0, false).
		AddItem(uibuilder.NewWidget(sp.themeName, "Logs:", sp.logs.View), 0, 1, 2, 1, 0, 0, false)

	pageStatsFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget(sp.themeName, "Current page:", sp.currentPageView), 0, 1, false).
		AddItem(uibuilder.NewWidget(sp.themeName, "Total found:", sp.totalFoundView), 0, 1, false)

	resultRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(uibuilder.NewWidget(sp.themeName, "List:", sp.resultList), 0, 3, false).
		AddItem(pageStatsFlex, 0, 1, false)

	sp.grid.AddItem(resultRow, 1, 0, 1, 1, 0, 100, false)

	metadataFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget(sp.themeName, "Title:", sp.title), 3, 1, false).
		AddItem(uibuilder.NewWidget(sp.themeName, "Description:", sp.description), 0, 3, false).
		AddItem(uibuilder.NewWidget(sp.themeName, "URLs:", sp.urls), 0, 2, false).
		AddItem(uibuilder.NewWidget(sp.themeName, "Language:", sp.language), 3, 1, false)

	sp.grid.AddItem(uibuilder.NewWidget(sp.themeName, "Code:", sp.body), 2, 0, 1, 1, 0, 100, false).
		AddItem(metadataFlex, 2, 1, 1, 1, 0, 100, false)
}

func (sp *SearchPage) initInputCapture() {
	sp.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event

		if event.Rune() == '?' {
			sp.tools.GoToPage(configuration.ShortcutModal, shortcutDescription)
			resultEvent = nil
		}

		switch event.Key() {
		case tcell.KeyCtrlF:
			sp.tools.Focus(sp.searchField)
			resultEvent = nil
		case tcell.KeyCtrlL:
			sp.tools.Focus(sp.resultList)
			resultEvent = nil
		case tcell.KeyCtrlC:
			bodyText := sp.body.GetText()
			if strings.TrimSpace(bodyText) == "" {
				sp.logs.AddErrorLogs([]string{"Body is empty. Nothing to copy."})
			} else {
				sp.logs.AddSuccessLogs([]string{"Body copied to clipboard."})
				clipboard.Write(clipboard.FmtText, []byte(sp.body.GetText()))
			}
			resultEvent = nil
		case tcell.KeyCtrlE:
			sp.goToEditPage()
			resultEvent = nil
		case tcell.KeyCtrlO:
			sp.tools.Focus(sp.searchType)
			resultEvent = nil
		}

		if event.Rune() == '}' {
			sp.showNextResultPage()
			resultEvent = nil
		}
		if event.Rune() == '{' {
			sp.showPreviousResultPage()
			resultEvent = nil
		}
		if event.Key() == tcell.KeyEnter && sp.resultList.HasFocus() {
			sp.goToEditPage()
			resultEvent = nil
		}

		return resultEvent
	})
}

func (sp *SearchPage) initFrame() {
	sp.Frame = uibuilder.NewPageFrame(sp.themeName, sp.grid, "Search")
}
