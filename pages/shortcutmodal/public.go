package shortcutmodal

import (
	"github.com/gwkeit/apptools"
	"github.com/gwkeit/configuration"
)

func (sm *ShortcutModal) SwitchToPage(shortcutList []apptools.ShortcutDescription) {
	sm.table.Clear()
	for i, entry := range shortcutList {
		sm.addKeyCell(i, entry.Key)
		sm.addDescriptionCell(i, entry.Description)
		sm.addSeparatorRow(i)
	}
	sm.tools.ShowPage(configuration.ShortcutModal)
}

func (sm *ShortcutModal) HideModalPage() {
	sm.tools.HidePage(configuration.ShortcutModal)
}
