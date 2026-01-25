package editpage

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/dto"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/validator"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

func (ep *EditPage) initMetadataFields() {
	ep.body = uibuilder.NewTextArea("", "")
	ep.title = uibuilder.NewTextArea("", "")
	ep.description = uibuilder.NewTextArea("", "")
	ep.urls = uibuilder.NewTextArea("", "")
}

func (ep *EditPage) initLayoutGrid() {
	ep.grid = tview.NewGrid().
		SetRows(14, 0).
		SetColumns(0, 50).
		SetBorders(false)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget("[ctrl+t] Title:", ep.title), 0, 1, false).
		AddItem(uibuilder.NewWidget("[ctrl+d] Description:", ep.description), 0, 3, false).
		AddItem(uibuilder.NewWidget("[ctrl+u] URLs:", ep.urls), 0, 3, false)

	ep.grid.AddItem(uibuilder.NewWidget("[ctrl+b] Code:", ep.body), 0, 0, 2, 1, 0, 100, false).
		AddItem(uibuilder.NewWidget("Logs:", ep.logs.View), 0, 1, 1, 1, 0, 100, false).
		AddItem(flex, 1, 1, 1, 1, 0, 100, false)
	ep.grid.SetBackgroundColor(tcell.ColorDefault)
}

func (ep *EditPage) initFrame() {
	ep.Frame = tview.NewFrame(ep.grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Edit snippet", true, tview.AlignCenter, tcell.ColorWhite)
	ep.Frame.SetBackgroundColor(tcell.ColorDefault)
}

func (ep *EditPage) initInputCapture() {
	ep.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event

		switch event.Key() {
		case tcell.KeyCtrlB:
			ep.globalDeps.App.SetFocus(ep.body)
			resultEvent = nil
		case tcell.KeyCtrlT:
			ep.globalDeps.App.SetFocus(ep.title)
			resultEvent = nil
		case tcell.KeyCtrlD:
			ep.globalDeps.App.SetFocus(ep.description)
			resultEvent = nil
		case tcell.KeyCtrlU:
			ep.globalDeps.App.SetFocus(ep.urls)
			resultEvent = nil
		case tcell.KeyCtrlS:
			snippetDto := dto.NewSnippetFromFields(
				ep.title.GetText(),
				ep.body.GetText(),
				ep.description.GetText(),
				ep.urls.GetText(),
			)
			validationErrors := validator.ValidateSnippet(snippetDto)

			if len(validationErrors) > 0 {
				allErrorsWithTopic := slices.Insert(validationErrors, 0, "Update failed.")
				ep.logs.AddErrorLogs(allErrorsWithTopic)
			} else {
				err := ep.globalDeps.Repo.UpdateSnippet(
					ep.globalDeps.Ctx,
					ep.snippetId,
					snippetDto,
				)
				if err != nil {
					ep.logs.AddErrorLogs([]string{err.Error()})
				} else {
					ep.logs.AddSuccessLogs([]string{fmt.Sprintf("Snippet '%s' updated successfully.", snippetDto.Title)})
				}
			}
			resultEvent = nil
		case tcell.KeyCtrlN:
			if ep.snippetId > -1 {
				ep.loadSnippet(ep.snippetId)
			}
			resultEvent = nil
		case tcell.KeyCtrlC:
			bodyText := ep.body.GetText()
			if strings.TrimSpace(bodyText) == "" {
				ep.logs.AddErrorLogs([]string{"Body is empty. Nothing to copy."})
			} else {
				ep.logs.AddSuccessLogs([]string{"Body copied to clipboard."})
				clipboard.Write(clipboard.FmtText, []byte(ep.body.GetText()))
			}
		}

		return resultEvent
	})
}
