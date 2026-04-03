package editpage

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/apptools"
	"github.com/gwkeit/configuration"
	"github.com/gwkeit/dto"
	"github.com/gwkeit/langdetector"
	"github.com/gwkeit/transform"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/validator"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
)

var once sync.Once

var shortcutDescription = []apptools.ShortcutDescription{
	{"ctrl+B", "Focus code field"},
	{"ctrl+O", "Focus and delete everything in code field"},
	{"ctrl+T", "Focus title field"},
	{"ctrl+D", "Focus description field"},
	{"ctrl+U", "Focus urls field"},
	{"ctrl+L", "Focus language field"},
	{"ctrl+S", "Save snippet"},
	{"ctrl+N", "Discard unsaved changes"},
	{"ctrl+C", "Copy snippet body"},
	{"ctrl+F", "Format the code"},
}

func (ep *EditPage) initMetadataFields() {
	ep.body = uibuilder.NewTextArea(ep.themeName, "", "")
	ep.body.SetChangedFunc(ep.updateCodePreview)
	ep.title = uibuilder.NewTextArea(ep.themeName, "", "")
	ep.description = uibuilder.NewTextArea(ep.themeName, "", "")
	ep.urls = uibuilder.NewTextArea(ep.themeName, "", "")
	ep.language = uibuilder.NewDropDown(ep.themeName, "")
	ep.language.SetOptions(configuration.LanguagesStrings, nil)
	ep.language.SetSelectedFunc(func(_ string, _ int) {
		if !ep.isLangSelectFuncSuppressed.Load() && !ep.isLangManuallySelected.Load() {
			ep.isLangManuallySelected.Store(true)
			ep.logs.AddInfoLogs([]string{"Language detect is disabled"})
		}
		ep.updateCodePreview()
	})
	ep.setLanguageOptionProgrammatically(0)
}

func (ep *EditPage) initLayoutGrid() {
	ep.grid = tview.NewGrid().
		SetRows(14, 0).
		SetColumns(0, 50).
		SetBorders(false)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget(ep.themeName, "Title:", ep.title), 3, 1, false).
		AddItem(uibuilder.NewWidget(ep.themeName, "Description:", ep.description), 0, 3, false).
		AddItem(uibuilder.NewWidget(ep.themeName, "URLs:", ep.urls), 0, 2, false).
		AddItem(uibuilder.NewWidget(ep.themeName, "Language:", ep.language), 3, 1, false)

	ep.grid.AddItem(uibuilder.NewWidget(ep.themeName, "Code:", ep.body), 0, 0, 1, 1, 0, 100, false).
		AddItem(uibuilder.NewWidget(ep.themeName, "Code Preview:", ep.codePreview.View), 1, 0, 1, 1, 0, 100, false).
		AddItem(uibuilder.NewWidget(ep.themeName, "Logs:", ep.logs.View), 0, 1, 1, 1, 0, 100, false).
		AddItem(flex, 1, 1, 1, 1, 0, 100, false)
	ep.grid.SetBackgroundColor(tcell.ColorDefault)
}

func (ep *EditPage) initFrame() {
	ep.Frame = uibuilder.NewPageFrame(ep.themeName, ep.grid, "Edit snippet")
}

func (ep *EditPage) initInputCapture() {
	ep.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event

		if event.Rune() == '?' {
			ep.tools.GoToPage(configuration.ShortcutModal, shortcutDescription)
			resultEvent = nil
		}

		switch event.Key() {
		case tcell.KeyCtrlB:
			ep.tools.Focus(ep.body)
			resultEvent = nil
		case tcell.KeyCtrlT:
			ep.tools.Focus(ep.title)
			resultEvent = nil
		case tcell.KeyCtrlD:
			ep.tools.Focus(ep.description)
			resultEvent = nil
		case tcell.KeyCtrlU:
			ep.tools.Focus(ep.urls)
			resultEvent = nil
		case tcell.KeyCtrlL:
			ep.tools.Focus(ep.language)
			resultEvent = nil
		case tcell.KeyCtrlF:
			bodyText := ep.body.GetText()
			alignedBodyText, isParsed := transform.AlignTextLeft(bodyText)
			if isParsed {
				ep.body.SetText(alignedBodyText, true)
			}
			resultEvent = nil
		case tcell.KeyCtrlO:
			ep.body.SetText("", true)
			ep.tools.Focus(ep.body)
			resultEvent = nil
		case tcell.KeyCtrlS:
			_, selectedLanguage := ep.language.GetCurrentOption()
			snippetDto := dto.NewSnippetFromFields(
				ep.title.GetText(),
				ep.body.GetText(),
				ep.description.GetText(),
				ep.urls.GetText(),
				selectedLanguage,
			)
			validationErrors := validator.ValidateSnippet(snippetDto)

			if len(validationErrors) > 0 {
				allErrorsWithTopic := slices.Insert(validationErrors, 0, "Update failed.")
				ep.logs.AddErrorLogs(allErrorsWithTopic)
			} else {
				err := ep.tools.Repo.UpdateSnippet(
					ep.tools.Ctx,
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
				ep.isLangManuallySelected.Store(false)
				ep.logs.AddInfoLogs([]string{"Language detect is enabled"})
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
			resultEvent = nil
		}

		return resultEvent
	})
}

func (ep *EditPage) initLangDetector() {
	once.Do(func() {
		ticker := time.NewTicker(time.Second)
		go func() {
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					var shouldDetect bool
					var bodyText string

					ep.tools.QueueUpdateDraw(func() {
						frontPage, _ := ep.tools.GetFrontPage()
						shouldDetect = frontPage == configuration.EditPage.String() &&
							!ep.isLangManuallySelected.Load() &&
							!ep.language.HasFocus()
						if shouldDetect {
							bodyText = ep.body.GetText()
						}
					})

					if shouldDetect {
						detectedLang := langdetector.Detect(bodyText)
						langIndex := slices.Index(configuration.LanguagesStrings, detectedLang.String())
						ep.tools.QueueUpdateDraw(func() {
							ep.setLanguageOptionProgrammatically(langIndex + 1)
						})
					}
				case <-ep.tools.Ctx.Done():
					return
				}
			}
		}()
	})
}
