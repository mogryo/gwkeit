package searchpage

import "github.com/gwkeit/configuration"

func (sp *SearchPage) SwitchToPage() {
	if sp.selectedSnippetId > -1 {
		sp.showSnippet(sp.selectedSnippetId)
	}
	sp.tools.SwitchToPage(configuration.SearchPage)
	sp.tools.Focus(sp.searchField)
}
