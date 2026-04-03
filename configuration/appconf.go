package configuration

import (
	"cmp"
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/gwkeit/uibuilder"
)

type SearchPageConf struct {
	SearchType string `json:"searchType"`
}

type ISearchPageConf interface {
	SetSearchType(searchType string)
	GetSearchType() string
}

type AllSnippetsConf struct {
	PageSize int64 `json:"pageSize"`
}

type IAllSnippetsConf interface {
	SetPageSize(pageSize int64)
	GetPageSize() int64
}

type AppConfiguration struct {
	SearchPage    SearchPageConf          `json:"searchPage"`
	AllSnippets   AllSnippetsConf         `json:"allSnippets"`
	AppThemeName  uibuilder.AppThemeName  `json:"themeName"`
	CodeThemeName uibuilder.CodeThemeName `json:"codeThemeName"`
}

var DefaultAppConf = &AppConfiguration{
	SearchPage: SearchPageConf{
		SearchType: "FTS",
	},
	AllSnippets: AllSnippetsConf{
		PageSize: 10,
	},
	AppThemeName:  uibuilder.DefaultAppTheme,
	CodeThemeName: uibuilder.DefaultCodeTheme,
}

func ReadConfiguration() *AppConfiguration {
	homeDir, _ := os.UserHomeDir()
	configFile, err := os.Open(path.Join(homeDir, AppDirectory, AppStateName))
	if err != nil {
		return DefaultAppConf
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			panic(err)
		}
	}(configFile)

	appState := *DefaultAppConf
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&appState); err == io.EOF {
		return DefaultAppConf
	} else if err != nil {
		panic(err)
	}

	return &appState
}

func updateAppConfigurationFile(updateFunc func(state *AppConfiguration)) *AppConfiguration {
	homeDir, _ := os.UserHomeDir()
	configFile, err := os.OpenFile(path.Join(homeDir, AppDirectory, AppStateName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			panic(err)
		}
	}(configFile)

	var appState AppConfiguration
	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(&appState); err == io.EOF {
		appState = *DefaultAppConf
	} else if err != nil {
		panic(err)
	}

	updateFunc(&appState)

	truncateErr := configFile.Truncate(0)
	_, seekErr := configFile.Seek(0, 0)
	if cmp.Or(truncateErr, seekErr) != nil {
		panic(truncateErr)
	}

	jsonEncoder := json.NewEncoder(configFile)
	if err = jsonEncoder.Encode(&appState); err != nil {
		panic(err)
	}

	return &appState
}

func (s *SearchPageConf) SetSearchType(searchType string) {
	appState := updateAppConfigurationFile(
		func(state *AppConfiguration) {
			state.SearchPage.SearchType = searchType
		},
	)
	s.SearchType = appState.SearchPage.SearchType
}

func (s *SearchPageConf) GetSearchType() string {
	return s.SearchType
}

func (as *AllSnippetsConf) SetPageSize(pageSize int64) {
	appState := updateAppConfigurationFile(
		func(state *AppConfiguration) {
			state.AllSnippets.PageSize = pageSize
		},
	)
	as.PageSize = appState.AllSnippets.PageSize
}

func (as *AllSnippetsConf) GetPageSize() int64 {
	return as.PageSize
}

func SetAppTheme(name uibuilder.AppThemeName) {
	updateAppConfigurationFile(
		func(state *AppConfiguration) {
			state.AppThemeName = name
		},
	)
}

func SetCodeTheme(name uibuilder.CodeThemeName) {
	updateAppConfigurationFile(
		func(state *AppConfiguration) {
			state.CodeThemeName = name
		},
	)
}
