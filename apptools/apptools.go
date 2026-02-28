package apptools

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/gwkeit/configuration"
	"github.com/gwkeit/repository"
	"github.com/rivo/tview"
)

type ShortcutDescription struct {
	Key         string
	Description string
}

type Tools struct {
	app          *tview.Application
	Repo         *repository.Repository
	Ctx          context.Context
	goToPage     func(pageName configuration.PageName, payload any)
	switchToPage func(pageName string) *tview.Pages
	showPage     func(pageName string) *tview.Pages
	hidePage     func(pageName string) *tview.Pages
	getFrontPage func() (name string, item tview.Primitive)
}

func New(
	ctx context.Context,
	app *tview.Application,
	repo *repository.Repository,
) *Tools {
	return &Tools{
		app:  app,
		Repo: repo,
		Ctx:  ctx,
	}
}

func (t *Tools) GoToPage(pageName configuration.PageName, payload any) {
	if t.goToPage == nil {
		panic("GoToPage function is not registered")
	}
	t.goToPage(pageName, payload)
}

// SwitchToPage tview.Pages.SwitchToPage(pageName string)
func (t *Tools) SwitchToPage(pageName configuration.PageName) {
	t.switchToPage(pageName.String())
}

func (t *Tools) HidePage(pageName configuration.PageName) {
	t.hidePage(pageName.String())
}

func (t *Tools) ShowPage(pageName configuration.PageName) {
	t.showPage(pageName.String())
}

func (t *Tools) Focus(component tview.Primitive) {
	t.app.SetFocus(component)
}

func (t *Tools) QueueEvent(event tcell.Event) {
	t.app.QueueEvent(event)
}

func (t *Tools) RefreshScreen() {
	t.app.Draw()
}

func (t *Tools) QueueUpdateDraw(f func()) {
	t.app.QueueUpdateDraw(f)
}

func (t *Tools) GetFrontPage() (name string, item tview.Primitive) {
	if t.getFrontPage == nil {
		panic("GetFrontPage function is not registered")
	}

	return t.getFrontPage()
}
