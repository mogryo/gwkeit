package additionpage

func (ap *AdditionPage) SwitchToAdditionPage() () {
	ap.globalDeps.Pages.SwitchToPage(PageName)
	ap.globalDeps.App.SetFocus(ap.body)
}
