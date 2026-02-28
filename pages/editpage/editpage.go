package editpage

import (
	"sync/atomic"

	"github.com/gwkeit/apptools"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
)

type EditPage struct {
	snippetId                  int64
	body                       *tview.TextArea
	title                      *tview.TextArea
	description                *tview.TextArea
	urls                       *tview.TextArea
	language                   *tview.DropDown
	grid                       *tview.Grid
	Frame                      *tview.Frame
	tools                      *apptools.Tools
	logs                       *widgets.LogsWidget
	isLangManuallySelected     atomic.Bool
	isLangSelectFuncSuppressed atomic.Bool
}

func NewPage(tools *apptools.Tools, logs *widgets.LogsWidget) *EditPage {
	ep := &EditPage{
		tools: tools,
		logs:  logs,
	}

	ep.initMetadataFields()
	ep.initLayoutGrid()
	ep.initInputCapture()
	ep.initFrame()
	ep.initLangDetector()

	return ep
}
