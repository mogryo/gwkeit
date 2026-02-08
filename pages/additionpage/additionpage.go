package additionpage

import (
	"github.com/gwkeit/apptools"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
)

type AdditionPage struct {
	body        *tview.TextArea
	title       *tview.TextArea
	description *tview.TextArea
	urls        *tview.TextArea
	grid        *tview.Grid
	Frame       *tview.Frame
	tools       *apptools.Tools
	logs        *widgets.LogsWidget
}

func NewPage(tools *apptools.Tools, logs *widgets.LogsWidget) *AdditionPage {
	additionPage := &AdditionPage{
		tools: tools,
		logs:  logs,
	}

	additionPage.initMetadataFields()
	additionPage.initGridLayout()
	additionPage.initInputCapture()
	additionPage.initFrame()

	return additionPage
}
