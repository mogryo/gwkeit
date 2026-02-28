package additionpage

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/apptools"
	"github.com/gwkeit/configuration"
	"github.com/gwkeit/dto"
	"github.com/gwkeit/langdetector"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/validator"
	"github.com/rivo/tview"
)

var once sync.Once

var shortcutDescription = []apptools.ShortcutDescription{
	{"ctrl+B", "Focus code field"},
	{"ctrl+T", "Focus title field"},
	{"ctrl+D", "Focus description field"},
	{"ctrl+U", "Focus urls field"},
	{"ctrl+L", "Focus language field"},
	{"ctrl+S", "Save snippet"},
	{"ctrl+N", "Clear fields"},
}

func (ap *AdditionPage) initFrame() {
	ap.Frame = tview.NewFrame(ap.grid).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("Add snippet", true, tview.AlignCenter, tcell.ColorDefault)
	ap.Frame.SetBackgroundColor(tcell.ColorDefault)
}

func (ap *AdditionPage) initGridLayout() {
	ap.grid = tview.NewGrid().
		SetRows(14, 0).
		SetColumns(0, 50).
		SetBorders(false)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(uibuilder.NewWidget("Title:", ap.title), 0, 1, false).
		AddItem(uibuilder.NewWidget("Description:", ap.description), 0, 3, false).
		AddItem(uibuilder.NewWidget("URLs:", ap.urls), 0, 2, false).
		AddItem(uibuilder.NewWidget("Language:", ap.language), 3, 1, false)

	ap.grid.
		AddItem(uibuilder.NewWidget("Code:", ap.body), 0, 0, 2, 1, 0, 100, false).
		AddItem(uibuilder.NewWidget("Logs:", ap.logs.View), 0, 1, 1, 1, 0, 100, false).
		AddItem(flex, 1, 1, 1, 1, 0, 100, false)
	ap.grid.SetBackgroundColor(tcell.ColorDefault)
}

func (ap *AdditionPage) initMetadataFields() {
	ap.body = uibuilder.NewTextArea("", "")
	ap.title = uibuilder.NewTextArea("", "")
	ap.description = uibuilder.NewTextArea("", "")
	ap.urls = uibuilder.NewTextArea("", "")
	ap.language = uibuilder.NewDropDown("")
	ap.language.SetOptions(slices.Concat([]string{""}, configuration.LanguagesStrings), nil)
	ap.language.SetSelectedFunc(func(_ string, _ int) {
		if !ap.isLangSelectFuncSuppressed.Load() {
			ap.isLangManuallySelected.Store(true)
			ap.logs.AddInfoLogs([]string{"Language detect is disabled"})
		}
	})
	ap.setLanguageOptionProgrammatically(0)
}

func (ap *AdditionPage) initInputCapture() {
	ap.grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		resultEvent := event

		if event.Rune() == '?' {
			ap.tools.GoToPage(configuration.ShortcutModal, shortcutDescription)
			resultEvent = nil
		}

		switch event.Key() {
		case tcell.KeyCtrlB:
			ap.tools.Focus(ap.body)
			resultEvent = nil
		case tcell.KeyCtrlT:
			ap.tools.Focus(ap.title)
			resultEvent = nil
		case tcell.KeyCtrlD:
			ap.tools.Focus(ap.description)
			resultEvent = nil
		case tcell.KeyCtrlU:
			ap.tools.Focus(ap.urls)
			resultEvent = nil
		case tcell.KeyCtrlL:
			ap.tools.Focus(ap.language)
			resultEvent = nil
		case tcell.KeyCtrlS:
			_, selectedLanguage := ap.language.GetCurrentOption()
			snippetDto := dto.NewSnippetFromFields(
				ap.title.GetText(),
				ap.body.GetText(),
				ap.description.GetText(),
				ap.urls.GetText(),
				selectedLanguage,
			)
			validationErrors := validator.ValidateSnippet(snippetDto)

			if len(validationErrors) > 0 {
				allErrorsWithTopic := slices.Insert(validationErrors, 0, "Save failed.")
				ap.logs.AddErrorLogs(allErrorsWithTopic)
			} else {
				snippetId, err := ap.tools.Repo.SaveSnippet(ap.tools.Ctx, snippetDto)
				if err != nil {
					ap.logs.AddErrorLogs([]string{err.Error()})
				} else {
					ap.logs.AddSuccessLogs([]string{fmt.Sprintf("Snippet '%s' saved successfully.", snippetDto.Title)})
					ap.clearFields()
					ap.tools.GoToPage(configuration.EditPage, snippetId)
				}
			}
			resultEvent = nil
		case tcell.KeyCtrlN:
			ap.body.SetText("", true)
			ap.title.SetText("", true)
			ap.description.SetText("", true)
			ap.urls.SetText("", true)
			ap.isLangManuallySelected.Store(false)
			ap.setLanguageOptionProgrammatically(0)
			ap.logs.AddInfoLogs([]string{"Language detect is enabled"})
			resultEvent = nil
		}

		return resultEvent
	})
}

func (ap *AdditionPage) initLangDetector() {
	once.Do(func() {
		ticker := time.NewTicker(time.Second)
		go func() {
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					var shouldDetect bool
					var bodyText string

					ap.tools.QueueUpdateDraw(func() {
						frontPage, _ := ap.tools.GetFrontPage()
						shouldDetect = frontPage == configuration.AdditionPage.String() &&
							!ap.isLangManuallySelected.Load() &&
							!ap.language.HasFocus()
						if shouldDetect {
							bodyText = ap.body.GetText()
						}
					})

					if shouldDetect {
						detectedLang := langdetector.Detect(bodyText)
						langIndex := slices.Index(configuration.LanguagesStrings, detectedLang.String())
						ap.tools.QueueUpdateDraw(func() {
							ap.setLanguageOptionProgrammatically(langIndex + 1)
						})
					}
				case <-ap.tools.Ctx.Done():
					return
				}
			}
		}()
	})
}
