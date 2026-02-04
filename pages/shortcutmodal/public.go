package shortcutmodal

import (
	"github.com/gwkeit/globaldeps"
)

func (sm *ShortcutModal) SwitchToShortcutPage(shortcutList []globaldeps.ShortcutDescription) {
	sm.table.Clear()
	for i, entry := range shortcutList {
		sm.addKeyCell(i, entry.Key)
		sm.addDescriptionCell(i, entry.Description)
		sm.addSeparatorRow(i)
	}
	sm.globalDeps.Pages.ShowPage(ModalName)
}

func (sm *ShortcutModal) HideModalPage() {
	sm.globalDeps.Pages.HidePage(ModalName)
}
