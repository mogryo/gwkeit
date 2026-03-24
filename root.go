package main

import (
	"github.com/gwkeit/configuration"
	"github.com/gwkeit/uibuilder"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gwkeit",
	Short: "gekeit is a tool to store code snippets",
	Long:  "gekeit is a tool to store code snippets",
	Run: func(cmd *cobra.Command, args []string) {
		ExecuteTUI()
	},
}

var setThemeCmd = &cobra.Command{
	Use:     "stheme",
	Aliases: []string{"settheme"},
	Short:   "Set the theme for the application",
	Long:    "Set the theme for the application",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		themeName := uibuilder.ThemeName(args[0])

		switch themeName {
		case uibuilder.DefaultTheme, uibuilder.LightTheme, uibuilder.DarkTheme, uibuilder.GreyTheme:
			configuration.SetAppTheme(themeName)
		default:
			cmd.Printf("Invalid theme %q. Allowed: default, light, dark, grey\n", args[0])
		}
	},
}

var readThemeCmd = &cobra.Command{
	Use:     "rtheme",
	Aliases: []string{"readtheme"},
	Short:   "Read the theme for the application",
	Long:    "Read the theme for the application",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		conf := configuration.ReadConfiguration()
		cmd.Printf("Current theme: %s", conf.ThemeName.String())
	},
}

func init() {
	rootCmd.AddCommand(setThemeCmd)
	rootCmd.AddCommand(readThemeCmd)
}
