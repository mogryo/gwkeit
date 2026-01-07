package pages

import (
	"fmt"
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/dto"
	"github.com/gwkeit/globaldeps"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/validator"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
)

type EditPage struct {
	snippetId   int64
	body        *tview.TextArea
	title       *tview.TextArea
	description *tview.TextArea
	urls        *tview.TextArea
	grid        *tview.Grid
	frame       *tview.Frame
}

func NewEditPage(globalDeps *globaldeps.GlobalDependencies, logs *widgets.LogsWidget) *EditPage {
	editPage := &EditPage{}

	editPage.body = uibuilder.NewTextArea("", "")
	editPage.title = uibuilder.NewTextArea("", "")
	editPage.description = uibuilder.NewTextArea("", "")
	editPage.urls = uibuilder.NewTextArea("", "")

	editPage.grid = tview.NewGrid().
		SetRows(14, 0).
		SetColumns(0, 50).
		SetBorders(false)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget("[ctrl+t] Title:", editPage.title), 0, 1, false).
		AddItem(uibuilder.NewWidget("[ctrl+d] Description:", editPage.description), 0, 3, false).
		AddItem(uibuilder.NewWidget("[ctrl+u] URLs:", editPage.urls), 0, 3, false)

	editPage.grid.AddItem(uibuilder.NewWidget("[ctrl+b] Code:", editPage.body), 0, 0, 2, 1, 0, 100, false).
		AddItem(uibuilder.NewWidget("Logs:", logs.View), 0, 1, 1, 1, 0, 100, false).
		AddItem(flex, 1, 1, 1, 1, 0, 100, false)

	editPage.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event

		switch event.Key() {
		case tcell.KeyCtrlB:
			globalDeps.App.SetFocus(editPage.body)
			resultEvent = nil
		case tcell.KeyCtrlT:
			globalDeps.App.SetFocus(editPage.title)
			resultEvent = nil
		case tcell.KeyCtrlD:
			globalDeps.App.SetFocus(editPage.description)
			resultEvent = nil
		case tcell.KeyCtrlU:
			globalDeps.App.SetFocus(editPage.urls)
			resultEvent = nil
		case tcell.KeyCtrlS:
			snippetDto := dto.NewSnippetFromFields(
				editPage.title.GetText(),
				editPage.body.GetText(),
				editPage.description.GetText(),
				editPage.urls.GetText(),
			)
			validationErrors := validator.ValidateSnippet(snippetDto)

			if len(validationErrors) > 0 {
				allErrorsWithTopic := slices.Insert(validationErrors, 0, "Update failed.")
				logs.AddErrorLogs(allErrorsWithTopic)
			} else {
				err := globalDeps.Repo.UpdateSnippet(
					globalDeps.Ctx,
					editPage.snippetId,
					snippetDto,
				)
				if err != nil {
					logs.AddErrorLogs([]string{err.Error()})
				} else {
					logs.AddSuccessLogs([]string{fmt.Sprintf("Snippet '%s' updated successfully.", snippetDto.Title)})
				}
			}
			resultEvent = nil
		case tcell.KeyCtrlN:
			if editPage.snippetId > -1 {
				editPage.loadSnippet(editPage.snippetId, globalDeps)
			}
			resultEvent = nil
		}

		return resultEvent
	})

	editPage.grid.SetBackgroundColor(tcell.ColorDefault)

	editPage.frame = tview.NewFrame(editPage.grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Edit snippet", true, tview.AlignCenter, tcell.ColorWhite)
	editPage.frame.SetBackgroundColor(tcell.ColorDefault)

	return editPage
}

func (ep *EditPage) loadSnippet(snippetId int64, globalDeps *globaldeps.GlobalDependencies) {
	ep.snippetId = snippetId
	snippet := globalDeps.Repo.FindSnippet(globalDeps.Ctx, snippetId)

	ep.body.SetText(snippet.Body, true)
	ep.title.SetText(snippet.Title, true)
	ep.description.SetText(snippet.Tags, true)
	ep.urls.SetText(snippet.Urls, true)
}

func (ep *EditPage) SwitchToEditPage(snippetId int64, globalDeps *globaldeps.GlobalDependencies) {
	ep.loadSnippet(snippetId, globalDeps)
	globalDeps.Pages.SwitchToPage("Edit")
	globalDeps.App.SetFocus(ep.body)
}
