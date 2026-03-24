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

var addCmd = &cobra.Command{
	Use:     "stheme",
	Aliases: []string{"settheme"},
	Short:   "Set the theme for the application",
	Long:    "Set the theme for the application",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		themeName := uibuilder.ThemeName(args[0])
		configuration.SetAppTheme(themeName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
