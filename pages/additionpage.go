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

type AdditionPage struct {
	body        *tview.TextArea
	title       *tview.TextArea
	description *tview.TextArea
	urls        *tview.TextArea
	grid        *tview.Grid
	frame       *tview.Frame
}

func NewAdditionPage(globalDeps *globaldeps.GlobalDependencies, logs *widgets.LogsWidget) *AdditionPage {
	additionPage := &AdditionPage{}

	additionPage.body = uibuilder.NewTextArea("", "")
	additionPage.title = uibuilder.NewTextArea("", "")
	additionPage.description = uibuilder.NewTextArea("", "")
	additionPage.urls = uibuilder.NewTextArea("", "")

	additionPage.grid = tview.NewGrid().
		SetRows(14, 0).
		SetColumns(0, 50).
		SetBorders(false)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget("[ctrl+t] Title:", additionPage.title), 0, 1, false).
		AddItem(uibuilder.NewWidget("[ctrl+d] Description:", additionPage.description), 0, 3, false).
		AddItem(uibuilder.NewWidget("[ctrl+u] URLs:", additionPage.urls), 0, 3, false)

	additionPage.grid.AddItem(uibuilder.NewWidget("[ctrl+b] Code:", additionPage.body), 0, 0, 2, 1, 0, 100, false).
		AddItem(uibuilder.NewWidget("Logs:", logs.View), 0, 1, 1, 1, 0, 100, false).
		AddItem(flex, 1, 1, 1, 1, 0, 100, false)

	additionPage.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event

		switch event.Key() {
		case tcell.KeyCtrlB:
			globalDeps.App.SetFocus(additionPage.body)
			resultEvent = nil
		case tcell.KeyCtrlT:
			globalDeps.App.SetFocus(additionPage.title)
			resultEvent = nil
		case tcell.KeyCtrlD:
			globalDeps.App.SetFocus(additionPage.description)
			resultEvent = nil
		case tcell.KeyCtrlU:
			globalDeps.App.SetFocus(additionPage.urls)
			resultEvent = nil
		case tcell.KeyCtrlS:
			snippetDto := dto.NewSnippetFromFields(
				additionPage.title.GetText(),
				additionPage.body.GetText(),
				additionPage.description.GetText(),
				additionPage.urls.GetText(),
			)
			validationErrors := validator.ValidateSnippet(snippetDto)

			if len(validationErrors) > 0 {
				allErrorsWithTopic := slices.Insert(validationErrors, 0, "Save failed.")
				logs.AddErrorLogs(allErrorsWithTopic)
			} else {
				snippetId, err := globalDeps.Repo.SaveSnippet(globalDeps.Ctx, snippetDto)
				if err != nil {
					logs.AddErrorLogs([]string{err.Error()})
				} else {
					logs.AddSuccessLogs([]string{fmt.Sprintf("Snippet '%s' saved successfully.", snippetDto.Title)})
					globalDeps.GoToEditPage(snippetId)
				}
			}
			resultEvent = nil
		case tcell.KeyCtrlN:
			additionPage.body.SetText("", true)
			additionPage.title.SetText("", true)
			additionPage.description.SetText("", true)
			additionPage.urls.SetText("", true)
			resultEvent = nil
		}

		return resultEvent
	})

	additionPage.grid.SetBackgroundColor(tcell.ColorDefault)
	additionPage.frame = tview.NewFrame(additionPage.grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Add snippet", true, tview.AlignCenter, tcell.ColorWhite)
	additionPage.frame.SetBackgroundColor(tcell.ColorDefault)

	return additionPage
}

func (ap *AdditionPage) SwitchToAdditionPage(globalDeps *globaldeps.GlobalDependencies) () {
	globalDeps.Pages.SwitchToPage("Addition")
	globalDeps.App.SetFocus(ap.body)
}
