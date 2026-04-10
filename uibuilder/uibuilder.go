package uibuilder

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	InputBackgroundStyle = tcell.StyleDefault.Background(tcell.ColorDefault)
)

func generateFieldStyle(appTheme *AppThemeConfig) tcell.Style {
	return InputBackgroundStyle.Foreground(appTheme.MainColor)
}

func NewInputField(themeName AppThemeName, label string, placeholder string) *tview.InputField {
	theme := GetAppTheme(themeName)

	field := tview.NewInputField().
		SetLabel(label).
		SetPlaceholder(placeholder)

	field.SetFieldStyle(generateFieldStyle(theme)).
		SetLabelStyle(InputBackgroundStyle).
		SetPlaceholderStyle(InputBackgroundStyle).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetBackgroundColor(tcell.ColorDefault)
	field.SetLabelColor(theme.LabelColor)

	return field
}

func NewTextView(themeName AppThemeName, text string) *tview.TextView {
	theme := GetAppTheme(themeName)

	textView := tview.NewTextView().SetText(text)
	textView.SetBackgroundColor(tcell.ColorDefault)
	textView.SetTextColor(theme.MainColor)

	return textView
}

func NewTextArea(
	themeName AppThemeName,
	label string,
	placeholder string,
) *tview.TextArea {
	theme := GetAppTheme(themeName)

	textArea := tview.NewTextArea().
		SetPlaceholder(placeholder)

	if label != "" {
		textArea.SetLabel(fmt.Sprintf("%-13s", label))
	}

	textArea.SetTextStyle(generateFieldStyle(theme)).
		SetPlaceholderStyle(InputBackgroundStyle).
		SetBorderPadding(0, 0, 0, 0).
		SetBackgroundColor(tcell.ColorDefault)

	return textArea
}

func NewWidget(themeName AppThemeName, title string, content tview.Primitive) *tview.Flex {
	theme := GetAppTheme(themeName)

	fieldFlex := tview.NewFlex()
	fieldFlex.
		SetDirection(tview.FlexRow).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(title).
		SetBorder(true).
		SetBorderColor(theme.BorderColor).
		SetTitleColor(theme.LabelColor).
		SetBackgroundColor(tcell.ColorDefault)
	fieldFlex.AddItem(content, 0, 1, false)

	return fieldFlex
}

func NewDropDown(themeName AppThemeName, title string) *tview.DropDown {
	theme := GetAppTheme(themeName)

	dropDown := tview.NewDropDown().
		SetLabel(title).
		SetLabelColor(theme.LabelColor)

	dropDown.SetFieldStyle(generateFieldStyle(theme)).
		SetLabelStyle(InputBackgroundStyle.Foreground(theme.LabelColor)).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetListStyles(
			tcell.StyleDefault.Background(tcell.ColorDefault),
			tcell.StyleDefault.Background(theme.SelectedColor),
		).
		SetBackgroundColor(tcell.ColorDefault)

	dropDown.SetDisabledStyle(generateFieldStyle(theme))

	return dropDown
}

func NewList(themeName AppThemeName) *tview.List {
	theme := GetAppTheme(themeName)

	list := tview.NewList()
	list.SetBackgroundColor(tcell.ColorDefault)
	list.SetMainTextStyle(generateFieldStyle(theme))
	list.SetShortcutStyle(InputBackgroundStyle.Foreground(theme.SelectedColor))
	list.SetSelectedFocusOnly(true)

	return list
}

func NewTable(themeName AppThemeName, rows int, columns int) *tview.Table {
	theme := GetAppTheme(themeName)

	table := tview.NewTable().
		SetBorders(true)
	table.SetBackgroundColor(tcell.ColorDefault)
	table.SetBordersColor(theme.BorderColor)
	table.SetFixed(rows, columns)
	table.SetSelectedStyle(tcell.StyleDefault.Background(theme.SelectedColor).Foreground(theme.MainColor))

	return table
}

func NewTableCell(themeName AppThemeName, text string) *tview.TableCell {
	theme := GetAppTheme(themeName)

	return tview.NewTableCell(text).
		SetAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorDefault).
		SetAttributes(tcell.AttrBold).SetTextColor(theme.MainColor)
}

func NewPageFrame(themeName AppThemeName, primitive tview.Primitive, text string) *tview.Frame {
	theme := GetAppTheme(themeName)

	frame := tview.NewFrame(primitive).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("[::b]"+text+"[::-]", true, tview.AlignCenter, theme.LabelColor)
	frame.SetBackgroundColor(tcell.ColorDefault)

	return frame
}

func NewGrid(_ AppThemeName) *tview.Grid {
	grid := tview.NewGrid().SetBorders(false)
	grid.SetBackgroundColor(tcell.ColorDefault)

	return grid
}
