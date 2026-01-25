package searchpage

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/cond"
	"github.com/gwkeit/slicelib"
	"github.com/gwkeit/uibuilder"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

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

func (sp *SearchPage) initSearchField() {
	executeSearch := func(text string) {
		splitConditions := strings.Split(strings.TrimSpace(text), " ")
		filteredConditions := slicelib.Filter(splitConditions, func(condition string) bool { return condition != "" })

		sp.resultList.Clear()
		if len(filteredConditions) > 0 {
			foundSnippets := sp.searchCallback(sp.globalDeps.Ctx, filteredConditions)
			for i, snippet := range foundSnippets {
				sp.resultList.AddItem(
					snippet.Title,
					strconv.FormatInt(snippet.ID, 10),
					cond.IfElse(i < len(shortcutRunes), shortcutRunes[i], 0),
					nil,
				)
			}
		}
	}

	sp.searchField = tview.NewInputField().
		SetLabel("").
		SetFieldWidth(0).
		SetAcceptanceFunc(func(text string, keyCode rune) bool { return true }).
		SetDoneFunc(func(key tcell.Key) {
			if sp.resultList.GetItemCount() == 0 {
				return
			}

			sp.globalDeps.App.SetFocus(sp.resultList)
			index := sp.resultList.GetCurrentItem()
			onSelect := sp.resultList.GetSelectedFunc()
			mainText, secText := sp.resultList.GetItemText(index)
			onSelect(index, mainText, secText, shortcutRunes[index])
		}).
		SetChangedFunc(executeSearch)
	sp.searchField.SetFieldStyle(uibuilder.InputBackgroundStyle)
	sp.searchField.SetBackgroundColor(tcell.ColorDefault)

	sp.searchType = tview.NewDropDown().
		SetLabel("[ctrl+o] Type: ").
		SetOptions([]string{"Tags", "Like", "FTS"}, func(text string, index int) {
			switch text {
			case "Tags":
				sp.searchCallback = sp.globalDeps.Repo.FindSnippetsByTags
			case "Like":
				sp.searchCallback = sp.globalDeps.Repo.FindSnippetsByLikeTags
			case "FTS":
				sp.searchCallback = sp.globalDeps.Repo.FindSnippetsByFts
			}
			executeSearch(sp.searchField.GetText())
		}).
		SetCurrentOption(0)
	sp.searchType.SetFieldStyle(uibuilder.InputBackgroundStyle).
		SetLabelStyle(uibuilder.InputBackgroundStyle).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetListStyles(
			tcell.StyleDefault.Background(tcell.ColorDefault),
			tcell.StyleDefault.Background(tcell.ColorGreenYellow),
		).
		SetBackgroundColor(tcell.ColorDefault)

	sp.searchBox = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(sp.searchField, 0, 4, false).
		AddItem(sp.searchType, 0, 1, false)

	sp.searchBox.SetBorderPadding(0, 0, 0, 0).SetBackgroundColor(tcell.ColorDefault)
}

func (sp *SearchPage) initResultList() {
	sp.resultList = tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			id, err := strconv.ParseInt(secondaryText, 10, 64)
			sp.selectedSnippetId = id

			if err != nil {
				panic(err)
			}

			snippet := sp.globalDeps.Repo.FindSnippet(sp.globalDeps.Ctx, id)

			sp.title.SetText(snippet.Title, true)
			sp.body.SetText(snippet.Body, true)
			sp.description.SetText(snippet.Description, true)
			sp.urls.SetText(snippet.Url, true)
		}).
		SetSelectedFocusOnly(true)
	sp.resultList.SetBackgroundColor(tcell.ColorDefault)
	sp.resultList.SetMainTextStyle(uibuilder.InputBackgroundStyle)
	sp.resultList.SetShortcutStyle(uibuilder.InputBackgroundStyle.Foreground(tcell.ColorGreen))
}

func (sp *SearchPage) initGridLayout() {
	sp.grid = tview.NewGrid().
		SetRows(3, 11).
		SetColumns(0, 50).
		SetBorders(false).
		AddItem(uibuilder.NewWidget("[ctr+f] Search:", sp.searchBox), 0, 0, 1, 1, 0, 0, false).
		AddItem(uibuilder.NewWidget("Logs:", sp.logs.View), 0, 1, 2, 1, 0, 0, false).
		AddItem(uibuilder.NewWidget("[ctr+l] List:", sp.resultList), 1, 0, 1, 1, 0, 0, false)
	metadataFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget("Title:", sp.title), 0, 1, false).
		AddItem(uibuilder.NewWidget("Description:", sp.description), 0, 3, false).
		AddItem(uibuilder.NewWidget("URLs:", sp.urls), 0, 3, false)

	sp.grid.AddItem(uibuilder.NewWidget("Code:", sp.body), 2, 0, 1, 1, 0, 100, false).
		AddItem(metadataFlex, 2, 1, 1, 1, 0, 100, false)

	sp.grid.SetBackgroundColor(tcell.ColorDefault)
}

func (sp *SearchPage) initInputCapture() {
	sp.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event

		switch event.Key() {
		case tcell.KeyCtrlF:
			sp.globalDeps.App.SetFocus(sp.searchField)
			resultEvent = nil
		case tcell.KeyCtrlL:
			sp.globalDeps.App.SetFocus(sp.resultList)
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
			if sp.selectedSnippetId > -1 {
				sp.globalDeps.GoToEditPage(sp.selectedSnippetId)
			} else {
				sp.logs.AddErrorLogs([]string{"No snippet selected."})
			}
			resultEvent = nil
		case tcell.KeyCtrlO:
			sp.globalDeps.App.SetFocus(sp.searchType)
			resultEvent = nil
		}

		return resultEvent
	})
}

func (sp *SearchPage) initFrame() {
	sp.Frame = tview.NewFrame(sp.grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Search", true, tview.AlignCenter, tcell.ColorWhite)
	sp.Frame.SetBackgroundColor(tcell.ColorDefault)
}
