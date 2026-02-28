package uibuilder

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	InputBackgroundStyle = tcell.StyleDefault.Background(tcell.ColorDefault)
)

func NewInputField(label string, placeholder string) *tview.InputField {
	field := tview.NewInputField().
		SetLabel(label).
		SetPlaceholder(placeholder)

	field.SetFieldStyle(InputBackgroundStyle).
		SetLabelStyle(InputBackgroundStyle).
		SetPlaceholderStyle(InputBackgroundStyle).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetBackgroundColor(tcell.ColorDefault)
	field.SetLabelStyle(tcell.StyleDefault.Background(tcell.ColorDefault))
	field.SetLabelColor(tcell.ColorDarkSeaGreen)

	return field
}

func NewTextView(text string) *tview.TextView {
	textView := tview.NewTextView().SetText(text)
	textView.SetBackgroundColor(tcell.ColorDefault)

	return textView
}

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

func NewDropDown(title string) *tview.DropDown {
	dropDown := tview.NewDropDown().
		SetLabel(title).
		SetLabelColor(tcell.ColorDarkSeaGreen)

	dropDown.SetFieldStyle(InputBackgroundStyle).
		SetLabelStyle(InputBackgroundStyle).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetListStyles(
			tcell.StyleDefault.Background(tcell.ColorDefault),
			tcell.StyleDefault.Background(tcell.ColorGreenYellow),
		).
		SetBackgroundColor(tcell.ColorDefault)

	return dropDown
}
