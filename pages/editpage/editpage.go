package editpage

import (
	"sync/atomic"

	"github.com/gwkeit/apptools"
	"github.com/gwkeit/uibuilder"
	"github.com/gwkeit/widgets"
	"github.com/gwkeit/widgets/codepreview"
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
	themeName uibuilder.AppThemeName,
	codeThemeName uibuilder.CodeThemeName,
) *EditPage {
	ep := &EditPage{
		tools:       tools,
		logs:        logs,
		themeName:   themeName,
		codePreview: codepreview.NewCodePreviewWidget(themeName, codeThemeName),
	}

	ep.initMetadataFields()
	ep.initLayoutGrid()
	ep.initInputCapture()
	ep.initFrame()
	ep.initLangDetector()

	return ep
}
