package uibuilder

import "github.com/gdamore/tcell/v2"

type ThemeName string

const (
	DefaultTheme ThemeName = "default"
	LightTheme   ThemeName = "light"
	DarkTheme    ThemeName = "dark"
	GreyTheme    ThemeName = "grey"
)

func (t ThemeName) String() string { return string(t) }

type AppThemeConfig struct {
	MainColor        tcell.Color `json:"mainColor"`
	SecondaryColor   tcell.Color `json:"secondaryColor"`
	SelectedColor    tcell.Color `json:"selectedColor"`
	BorderColor      tcell.Color `json:"borderColor"`
	LabelColor       tcell.Color `json:"labelColor"`
	SuccessMessage   string      `json:"successMessage"`
	InfoMessage      string      `json:"infoMessage"`
	ErrorMessage     string      `json:"errorMessage"`
	TimestampMessage string      `json:"timestampMessage"`
}

var appThemes = map[string]AppThemeConfig{
	DefaultTheme.String(): {
		MainColor:        tcell.ColorDefault,
		SelectedColor:    tcell.ColorGreen,
		SecondaryColor:   tcell.ColorDefault,
		BorderColor:      tcell.ColorDefault,
		LabelColor:       tcell.ColorDefault,
		SuccessMessage:   "#a6e22e",
		InfoMessage:      "#0087ff",
		ErrorMessage:     "#ff4689",
		TimestampMessage: "grey",
	},
	LightTheme.String(): {
		MainColor:        tcell.ColorBlack,
		SecondaryColor:   tcell.ColorDarkSlateGrey,
		SelectedColor:    tcell.ColorGreen,
		BorderColor:      tcell.ColorDarkGrey,
		LabelColor:       tcell.ColorDarkSeaGreen,
		SuccessMessage:   "#2f8f2f",
		InfoMessage:      "#005fcc",
		ErrorMessage:     "#cc2f5a",
		TimestampMessage: "grey",
	},
	DarkTheme.String(): {
		MainColor:        tcell.ColorWhite,
		SecondaryColor:   tcell.ColorSilver,
		SelectedColor:    tcell.ColorGreen,
		BorderColor:      tcell.ColorWhiteSmoke,
		LabelColor:       tcell.ColorDarkSeaGreen,
		SuccessMessage:   "#8be28b",
		InfoMessage:      "#4da3ff",
		ErrorMessage:     "#ff6b8b",
		TimestampMessage: "#a0a0a0",
	},
	GreyTheme.String(): {
		MainColor:        tcell.ColorWhiteSmoke,
		SelectedColor:    tcell.ColorGreen,
		SecondaryColor:   tcell.ColorGainsboro,
		BorderColor:      tcell.ColorDimGrey,
		LabelColor:       tcell.ColorDarkSeaGreen,
		SuccessMessage:   "#73c973",
		InfoMessage:      "#66aaff",
		ErrorMessage:     "#ff7a96",
		TimestampMessage: "#b0b0b0",
	},
}

func GetTheme(themeName ThemeName) *AppThemeConfig {
	theme, ok := appThemes[themeName.String()]
	if !ok {
		panic("Theme not found")
	}

	return &theme
}
