package uibuilder

import "github.com/gdamore/tcell/v2"

type AppThemeName string
type CodeThemeName string

const (
	DefaultAppTheme AppThemeName = "default"
	LightAppTheme   AppThemeName = "light"
	DarkAppTheme    AppThemeName = "dark"
	GreyAppTheme    AppThemeName = "grey"
)

const (
	DefaultCodeTheme CodeThemeName = "default"
	LightCodeTheme   CodeThemeName = "light"
	DarkCodeTheme    CodeThemeName = "dark"
)

func (t AppThemeName) String() string  { return string(t) }
func (t CodeThemeName) String() string { return string(t) }

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

type CodeThemeConfig struct {
	Keyword    string `json:"keyword"`
	String     string `json:"string"`
	Comment    string `json:"comment"`
	Number     string `json:"number"`
	Identifier string `json:"identifier"`
}

var codeThemes = map[CodeThemeName]CodeThemeConfig{
	DefaultCodeTheme: {
		Keyword:    "[#66aaff]",
		String:     "[#e6db74]",
		Comment:    "[grey]",
		Number:     "[#73c973]",
		Identifier: "[white]",
	},
	LightCodeTheme: {
		Keyword:    "[#005fcc]",
		String:     "[#8a6d3b]",
		Comment:    "[#808080]",
		Number:     "[#2f8f2f]",
		Identifier: "[black]",
	},
	DarkCodeTheme: {
		Keyword:    "[#66aaff]",
		String:     "[#e6db74]",
		Comment:    "[#a0a0a0]",
		Number:     "[#73c973]",
		Identifier: "[white]",
	},
}

var appThemes = map[AppThemeName]AppThemeConfig{
	DefaultAppTheme: {
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
	LightAppTheme: {
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
	DarkAppTheme: {
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
	GreyAppTheme: {
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

func GetAppTheme(themeName AppThemeName) *AppThemeConfig {
	theme, ok := appThemes[themeName]
	if !ok {
		panic("Application theme not found")
	}

	return &theme
}

func GetCodeTheme(themeName CodeThemeName) *CodeThemeConfig {
	theme, ok := codeThemes[themeName]
	if !ok {
		panic("Code theme not found")
	}

	return &theme
}
