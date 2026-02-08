package apptools

import (
	"github.com/gwkeit/configuration"
	"github.com/rivo/tview"
)

func (t *Tools) RegisterSwitchToPage(switchToPage func(pageName string) *tview.Pages) {
	t.switchToPage = switchToPage
}

func (t *Tools) RegisterShowPage(showPage func(pageName string) *tview.Pages) {
	t.showPage = showPage
}

func (t *Tools) RegisterHidePage(hidePage func(pageName string) *tview.Pages) {
	t.hidePage = hidePage
}

func (t *Tools) RegisterGoToPage(goToPage func(pageName configuration.PageName, payload any)) {
	t.goToPage = goToPage
}
