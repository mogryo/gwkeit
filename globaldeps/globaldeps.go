package globaldeps

import (
	"context"

	"github.com/gwkeit/repository"
	"github.com/rivo/tview"
)

type ShortcutDescription struct {
	Key         string
	Description string
}

type DepsChanPayload struct {
	PageName     string
	SnippetId    int64
	ShortcutList []ShortcutDescription
}

type GlobalDependencies struct {
	App   *tview.Application
	Repo  *repository.Repository
	Ctx   context.Context
	Pages *tview.Pages
	Chan  chan DepsChanPayload
}

func New(
	ctx context.Context,
	app *tview.Application,
	repo *repository.Repository,
) *GlobalDependencies {
	return &GlobalDependencies{
		App:   app,
		Repo:  repo,
		Ctx:   ctx,
		Pages: tview.NewPages(),
		Chan:  make(chan DepsChanPayload),
	}
}

func (gd *GlobalDependencies) GoToEditPage(snippetId int64) {
	gd.Chan <- DepsChanPayload{PageName: "Edit", SnippetId: snippetId}
}

func (gd *GlobalDependencies) ShowShortcutModal(shortcutList []ShortcutDescription) {
	gd.Chan <- DepsChanPayload{PageName: "Modal", ShortcutList: shortcutList}
}
