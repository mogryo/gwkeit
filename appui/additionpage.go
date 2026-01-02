package appui

import (
	"fmt"
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/repository"
	"github.com/gwkeit/transform"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/utils"
	"github.com/gwkeit/validator"
	"github.com/rivo/tview"
)

type AdditionPageUI struct {
	body        *tview.TextArea
	title       *tview.TextArea
	description *tview.TextArea
	urls        *tview.TextArea
	grid        *tview.Grid
}

func (aui *AppUI) NewAdditionPage() {
	aui.additionPage = &AdditionPageUI{}

	aui.additionPage.body = uibuilder.NewTextArea("", "")
	aui.additionPage.title = uibuilder.NewTextArea("", "")
	aui.additionPage.description = uibuilder.NewTextArea("", "")
	aui.additionPage.urls = uibuilder.NewTextArea("", "")

	aui.additionPage.grid = tview.NewGrid().
		SetRows(14, 0).
		SetColumns(0, 50).
		SetBorders(false)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget("[ctrl+t] Title:", aui.additionPage.title), 0, 1, false).
		AddItem(uibuilder.NewWidget("[ctrl+d] Description:", aui.additionPage.description), 0, 3, false).
		AddItem(uibuilder.NewWidget("[ctrl+u] URLs:", aui.additionPage.urls), 0, 3, false)

	aui.additionPage.grid.AddItem(uibuilder.NewWidget("[ctrl+b] Code:", aui.additionPage.body), 0, 0, 2, 1, 0, 100, false).
		AddItem(uibuilder.NewWidget("Logs:", aui.logs.View), 0, 1, 1, 1, 0, 100, false).
		AddItem(flex, 1, 1, 1, 1, 0, 100, false)

	aui.additionPage.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event

		switch event.Key() {
		case tcell.KeyCtrlB:
			aui.App.SetFocus(aui.additionPage.body)
			resultEvent = nil
		case tcell.KeyCtrlT:
			aui.App.SetFocus(aui.additionPage.title)
			resultEvent = nil
		case tcell.KeyCtrlD:
			aui.App.SetFocus(aui.additionPage.description)
			resultEvent = nil
		case tcell.KeyCtrlU:
			aui.App.SetFocus(aui.additionPage.urls)
			resultEvent = nil
		case tcell.KeyCtrlS:
			tags := transform.FormDescriptionToTagList(aui.additionPage.description.GetText())
			urls := transform.FormUrlsToUrlList(aui.additionPage.urls.GetText())
			body := transform.CleanupBody(aui.additionPage.body.GetText())
			title := transform.CleanupTitle(aui.additionPage.title.GetText())

			tagErrors := validator.ValidateTags(tags)
			urlErrors := validator.ValidateUrls(urls)
			bodyErrors := validator.ValidateBody(body)
			titleErrors := validator.ValidateTitle(title)

			if len(tagErrors) > 0 || len(urlErrors) > 0 || len(bodyErrors) > 0 || len(titleErrors) > 0 {
				concatenatedErrors := utils.Concat(bodyErrors, titleErrors, tagErrors, urlErrors)
				allErrorsWithTopic := slices.Insert(concatenatedErrors, 0, "Save failed.")
				aui.logs.AddErrorLogs(allErrorsWithTopic)
			} else {
				err := aui.repo.SaveSnippet(aui.ctx, &repository.SnippetInput{Title: title, Body: body, Tags: tags, Urls: urls})
				if err != nil {
					aui.logs.AddErrorLogs([]string{err.Error()})
				} else {
					aui.logs.AddSuccessLogs([]string{fmt.Sprintf("Snippet '%s' saved successfully.", title)})
				}
			}
			resultEvent = nil
		}

		return resultEvent
	})

	aui.additionPage.grid.SetBackgroundColor(tcell.ColorDefault)
	aui.pages.AddPage("Addition", aui.additionPage.grid, true, false)
}
