package langdetector

import (
	"slices"
	"strings"

	"github.com/gwkeit/configuration"
)

func Detect(text string) configuration.Language {
	detectedLanguages := make([]configuration.Language, 0)
	if strings.TrimSpace(text) == "" {
		return configuration.Text
	}

	if IsKotlin(text) {
		detectedLanguages = append(detectedLanguages, configuration.Kotlin)
	}
	if IsGo(text) {
		detectedLanguages = append(detectedLanguages, configuration.Go)
	}
	if IsPython(text) {
		detectedLanguages = append(detectedLanguages, configuration.Python)
	}
	if IsRuby(text) {
		detectedLanguages = append(detectedLanguages, configuration.Ruby)
	}
	if IsTypeScript(text) {
		detectedLanguages = append(detectedLanguages, configuration.TypeScript)
	}
	if IsSQLite(text) {
		detectedLanguages = append(detectedLanguages, configuration.SQLite)
	}
	if IsPostgreSQL(text) {
		detectedLanguages = append(detectedLanguages, configuration.PostgreSQL)
	}

	if len(detectedLanguages) == 1 {
		return detectedLanguages[0]
	}

	if len(detectedLanguages) == 2 &&
		slices.Contains(detectedLanguages, configuration.SQLite) &&
		slices.Contains(detectedLanguages, configuration.PostgreSQL) {
		return configuration.PostgreSQL
	}

	return configuration.Text
}
