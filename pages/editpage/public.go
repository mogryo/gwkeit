package editpage

import "github.com/gwkeit/configuration"

func (ep *EditPage) SwitchToPage(snippetId int64) {
	ep.loadSnippet(snippetId)
	ep.tools.SwitchToPage(configuration.EditPage)
	ep.tools.Focus(ep.body)
}
