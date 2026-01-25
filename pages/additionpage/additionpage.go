package additionpage

import (
	"github.com/gwkeit/globaldeps"
	"github.com/gwkeit/widgets"
	"github.com/rivo/tview"
)

const PageName = "Addition"

type AdditionPage struct {
	body        *tview.TextArea
	title       *tview.TextArea
	description *tview.TextArea
	urls        *tview.TextArea
	grid        *tview.Grid
	Frame       *tview.Frame
	globalDeps  *globaldeps.GlobalDependencies
	logs        *widgets.LogsWidget
}

func NewPage(globalDeps *globaldeps.GlobalDependencies, logs *widgets.LogsWidget) *AdditionPage {
	additionPage := &AdditionPage{
		globalDeps: globalDeps,
		logs:       logs,
	}

	additionPage.initMetadataFields()
	additionPage.initGridLayout()
	additionPage.initInputCapture()
	additionPage.initFrame()

	return additionPage
}
