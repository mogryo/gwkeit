package codepreview

import (
	"strings"

	"github.com/gwkeit/configuration"
	"github.com/gwkeit/uibuilder"
	"github.com/rivo/tview"
)

const colorReset = "[-]"

type CodePreviewWidget struct {
	View          *tview.TextView
	appThemeName  uibuilder.AppThemeName
	codeThemeName uibuilder.CodeThemeName
}

func NewCodePreviewWidget(
	appThemeName uibuilder.AppThemeName,
	codeThemeName uibuilder.CodeThemeName,
) *CodePreviewWidget {
	cpw := &CodePreviewWidget{
		appThemeName:  appThemeName,
		codeThemeName: codeThemeName,
	}
	cpw.View = uibuilder.NewTextView(appThemeName, "").
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetWrap(true)

	return cpw
}

func (widget *CodePreviewWidget) SetText(
	text string,
	language configuration.Language,
) {
	codeTheme := uibuilder.GetCodeTheme(widget.codeThemeName)
	switch language {
	case configuration.Go:
		widget.View.SetText(addGoDynamicColors(codeTheme, text))
	case configuration.Kotlin:
		widget.View.SetText(addKotlinDynamicColors(codeTheme, text))
	case configuration.Python:
		widget.View.SetText(addPythonDynamicColors(codeTheme, text))
	case configuration.Ruby:
		widget.View.SetText(addRubyDynamicColors(codeTheme, text))
	case configuration.TypeScript:
		widget.View.SetText(addTypeScriptDynamicColors(codeTheme, text))
	default:
		widget.View.SetText(strings.ReplaceAll(text, "[", "[["))
	}
}
