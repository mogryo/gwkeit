package editpage

import (
	"github.com/gwkeit/globaldeps"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
)

const PageName = "Edit"

type EditPage struct {
	snippetId   int64
	body        *tview.TextArea
	title       *tview.TextArea
	description *tview.TextArea
	urls        *tview.TextArea
	grid        *tview.Grid
	Frame       *tview.Frame
	globalDeps  *globaldeps.GlobalDependencies
	logs        *widgets.LogsWidget
}

func NewPage(globalDeps *globaldeps.GlobalDependencies, logs *widgets.LogsWidget) *EditPage {
	ep := &EditPage{
		globalDeps: globalDeps,
		logs:       logs,
	}

	ep.initMetadataFields()
	ep.initLayoutGrid()
	ep.initInputCapture()
	ep.initFrame()

	return ep
}
