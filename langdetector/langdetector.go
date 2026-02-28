package langdetector

import (
	"github.com/gwkeit/configuration"
)

func Detect(text string) configuration.Language {
	if IsKotlin(text) {
		return configuration.Kotlin
	}
	if IsGo(text) {
		return configuration.Go
	}
	if IsPython(text) {
		return configuration.Python
	}
	if IsRuby(text) {
		return configuration.Ruby
	}
	if IsTypeScript(text) {
		return configuration.TypeScript
	}

	return configuration.Text
}
