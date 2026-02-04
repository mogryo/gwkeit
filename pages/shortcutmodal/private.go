package shortcutmodal

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (sm *ShortcutModal) addKeyCell(rowIdx int, text string) {
	sm.table.SetCell(
		rowIdx*2,
		0,
		tview.NewTableCell(text).
			SetExpansion(2).
			SetAlign(tview.AlignLeft).
			SetBackgroundColor(tcell.ColorDefault).
			SetAttributes(tcell.AttrBold),
	)
}

func (sm *ShortcutModal) addDescriptionCell(rowIdx int, text string) {
	sm.table.SetCell(
		rowIdx*2,
		1,
		tview.NewTableCell(text).
			SetAlign(tview.AlignLeft).
			SetBackgroundColor(tcell.ColorDefault).
			SetAttributes(tcell.AttrBold),
	)
}

func (sm *ShortcutModal) addSeparatorRow(rowIdx int) {
	sm.table.SetCell(rowIdx*2+1, 0, tview.NewTableCell("").SetSelectable(false))
	sm.table.SetCell(rowIdx*2+1, 1, tview.NewTableCell("").SetSelectable(false))
}
