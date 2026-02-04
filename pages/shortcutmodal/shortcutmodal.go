package shortcutmodal

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/globaldeps"
	"github.com/rivo/tview"
)

const ModalName = "Modal"

type ShortcutModal struct {
	Frame      *tview.Flex
	table      *tview.Table
	globalDeps *globaldeps.GlobalDependencies
}

func NewModal(globalDeps *globaldeps.GlobalDependencies) *ShortcutModal {
	mp := &ShortcutModal{
		globalDeps: globalDeps,
		table:      tview.NewTable(),
	}
	mp.table.SetBorder(true).SetBackgroundColor(tcell.ColorDefault)

	tableContainer := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(mp.table, 20, 1, true).
		AddItem(nil, 0, 1, false)

	mp.Frame = tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tableContainer, 48, 1, true).
		AddItem(nil, 0, 1, false)

	mp.Frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == '?' || event.Key() == tcell.KeyEscape {
			mp.HideModalPage()
		} else {
			mp.globalDeps.Pages.HidePage(ModalName)
			mp.globalDeps.App.QueueEvent(event)
		}

		return nil
	})

	return mp
}
