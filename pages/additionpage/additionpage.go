package additionpage

import (
	"sync/atomic"

	"github.com/gwkeit/apptools"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/widgets"
	"github.com/gwkeit/widgets/codepreview"
	"github.com/rivo/tview"
)

type AdditionPage struct {
	body                       *tview.TextArea
	title                      *tview.TextArea
	description                *tview.TextArea
	language                   *tview.DropDown
	urls                       *tview.TextArea
	grid                       *tview.Grid
	Frame                      *tview.Frame
	codePreview                *codepreview.CodePreviewWidget
	tools                      *apptools.Tools
	logs                       *widgets.LogsWidget
	isLangManuallySelected     atomic.Bool
	isLangSelectFuncSuppressed atomic.Bool
	themeName                  uibuilder.AppThemeName
}

func NewPage(
	tools *apptools.Tools,
	logs *widgets.LogsWidget,
	appThemeName uibuilder.AppThemeName,
	codeThemeName uibuilder.CodeThemeName,
) *AdditionPage {
	additionPage := &AdditionPage{
		tools:                      tools,
		logs:                       logs,
		isLangManuallySelected:     atomic.Bool{},
		isLangSelectFuncSuppressed: atomic.Bool{},
		themeName:                  appThemeName,
		codePreview:                codepreview.NewCodePreviewWidget(appThemeName, codeThemeName),
	}

	additionPage.isLangManuallySelected.Store(false)
	additionPage.isLangSelectFuncSuppressed.Store(false)
	additionPage.initMetadataFields()
	additionPage.initGridLayout()
	additionPage.initInputCapture()
	additionPage.initFrame()
	additionPage.initLangDetector()

	return additionPage
}
