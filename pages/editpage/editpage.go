package editpage

import (
	"github.com/gwkeit/apptools"
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
	Frame       *tview.Frame
	tools       *apptools.Tools
	logs        *widgets.LogsWidget
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

	return ep
}
