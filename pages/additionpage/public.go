package additionpage

import "github.com/gwkeit/configuration"

func (ap *AdditionPage) SwitchToPage() () {
	ap.tools.SwitchToPage(configuration.AdditionPage)
	ap.tools.Focus(ap.body)
}
