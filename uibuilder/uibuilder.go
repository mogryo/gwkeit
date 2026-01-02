package uibuilder

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	InputBackgroundStyle = tcell.StyleDefault.Background(tcell.ColorDefault)
)

func NewTextArea(
	label string,
	placeholder string,
) *tview.TextArea {
	textArea := tview.NewTextArea().
		SetPlaceholder(placeholder)

	if label != "" {
		textArea.SetLabel(fmt.Sprintf("%-13s", label))
	}

	textArea.SetTextStyle(InputBackgroundStyle).
		SetPlaceholderStyle(InputBackgroundStyle).
		SetBorderPadding(0, 0, 0, 0).
		SetBackgroundColor(tcell.ColorDefault)

	return textArea
}

func NewWidget(title string, content tview.Primitive) *tview.Flex {
	fieldFlex := tview.NewFlex()
	fieldFlex.
		SetDirection(tview.FlexRow).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(title).
		SetBorder(true).
		SetBackgroundColor(tcell.ColorDefault)
	fieldFlex.AddItem(content, 0, 1, false)

	return fieldFlex
}
