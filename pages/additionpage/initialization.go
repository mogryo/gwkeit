package additionpage

import (
	"fmt"
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/dto"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/validator"
	"github.com/rivo/tview"
)

func (ap *AdditionPage) initFrame() {
	ap.Frame = tview.NewFrame(ap.grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Add snippet", true, tview.AlignCenter, tcell.ColorWhite)
	ap.Frame.SetBackgroundColor(tcell.ColorDefault)
}

func (ap *AdditionPage) initGridLayout() {
	ap.grid = tview.NewGrid().
		SetRows(14, 0).
		SetColumns(0, 50).
		SetBorders(false)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget("[ctrl+t] Title:", ap.title), 0, 1, false).
		AddItem(uibuilder.NewWidget("[ctrl+d] Description:", ap.description), 0, 3, false).
		AddItem(uibuilder.NewWidget("[ctrl+u] URLs:", ap.urls), 0, 3, false)

	ap.grid.
		AddItem(uibuilder.NewWidget("[ctrl+b] Code:", ap.body), 0, 0, 2, 1, 0, 100, false).
		AddItem(uibuilder.NewWidget("Logs:", ap.logs.View), 0, 1, 1, 1, 0, 100, false).
		AddItem(flex, 1, 1, 1, 1, 0, 100, false)
	ap.grid.SetBackgroundColor(tcell.ColorDefault)
}

func (ap *AdditionPage) initMetadataFields() {
	ap.body = uibuilder.NewTextArea("", "")
	ap.title = uibuilder.NewTextArea("", "")
	ap.description = uibuilder.NewTextArea("", "")
	ap.urls = uibuilder.NewTextArea("", "")
}

func (ap *AdditionPage) initInputCapture() {
	ap.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event

		switch event.Key() {
		case tcell.KeyCtrlB:
			ap.globalDeps.App.SetFocus(ap.body)
			resultEvent = nil
		case tcell.KeyCtrlT:
			ap.globalDeps.App.SetFocus(ap.title)
			resultEvent = nil
		case tcell.KeyCtrlD:
			ap.globalDeps.App.SetFocus(ap.description)
			resultEvent = nil
		case tcell.KeyCtrlU:
			ap.globalDeps.App.SetFocus(ap.urls)
			resultEvent = nil
		case tcell.KeyCtrlS:
			snippetDto := dto.NewSnippetFromFields(
				ap.title.GetText(),
				ap.body.GetText(),
				ap.description.GetText(),
				ap.urls.GetText(),
			)
			validationErrors := validator.ValidateSnippet(snippetDto)

			if len(validationErrors) > 0 {
				allErrorsWithTopic := slices.Insert(validationErrors, 0, "Save failed.")
				ap.logs.AddErrorLogs(allErrorsWithTopic)
			} else {
				snippetId, err := ap.globalDeps.Repo.SaveSnippet(ap.globalDeps.Ctx, snippetDto)
				if err != nil {
					ap.logs.AddErrorLogs([]string{err.Error()})
				} else {
					ap.logs.AddSuccessLogs([]string{fmt.Sprintf("Snippet '%s' saved successfully.", snippetDto.Title)})
					ap.clearFields()
					ap.globalDeps.GoToEditPage(snippetId)
				}
			}
			resultEvent = nil
		case tcell.KeyCtrlN:
			ap.body.SetText("", true)
			ap.title.SetText("", true)
			ap.description.SetText("", true)
			ap.urls.SetText("", true)
			resultEvent = nil
		}

		return resultEvent
	})
}
