package additionpage

import (
	"sync/atomic"

	"github.com/gwkeit/apptools"
	"github.com/gwkeit/widgets"
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
	tools                      *apptools.Tools
	logs                       *widgets.LogsWidget
	isLangManuallySelected     atomic.Bool
	isLangSelectFuncSuppressed atomic.Bool
}

func NewPage(tools *apptools.Tools, logs *widgets.LogsWidget) *AdditionPage {
	additionPage := &AdditionPage{
		tools:                      tools,
		logs:                       logs,
		isLangManuallySelected:     atomic.Bool{},
		isLangSelectFuncSuppressed: atomic.Bool{},
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
