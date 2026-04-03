package allsnippetspage

import "github.com/gwkeit/configuration"

func (asp *AllSnippetsPage) SwitchToPage() {
	asp.tools.SwitchToPage(configuration.AllSnippetsPage)
	asp.populateTable(asp.tools.Ctx)
	asp.focusTable()
	asp.clearMetadataFields()
}
