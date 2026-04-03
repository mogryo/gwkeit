package main

import (
	"fmt"

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
	Use:     "sat",
	Aliases: []string{"setapptheme"},
	Short:   "Set the theme for the application",
	Long:    "Set the theme for the application",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		themeName := uibuilder.AppThemeName(args[0])

		switch themeName {
		case uibuilder.DefaultAppTheme, uibuilder.LightAppTheme, uibuilder.DarkAppTheme, uibuilder.GreyAppTheme:
			configuration.SetAppTheme(themeName)
			return nil
		default:
			return fmt.Errorf("Invalid application theme %q. Allowed: default, light, dark, grey\n", args[0])
		}
	},
}

var readThemeCmd = &cobra.Command{
	Use:     "rat",
	Aliases: []string{"readapptheme"},
	Short:   "Read the theme for the application",
	Long:    "Read the theme for the application",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		conf := configuration.ReadConfiguration()
		cmd.Printf("Current theme: %s", conf.AppThemeName.String())
	},
}

var setCodeThemeCmd = &cobra.Command{
	Use:     "sct",
	Aliases: []string{"setcodetheme"},
	Short:   "Set the code theme for the code preview",
	Long:    "Set the code theme for the code preview",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		themeName := uibuilder.CodeThemeName(args[0])

		switch themeName {
		case uibuilder.DefaultCodeTheme, uibuilder.LightCodeTheme, uibuilder.DarkCodeTheme:
			configuration.SetCodeTheme(themeName)
			return nil
		default:
			return fmt.Errorf("Invalid code theme %q. Allowed: default, light, dark", args[0])
		}
	},
}

var readCodeThemeCmd = &cobra.Command{
	Use:     "rct",
	Aliases: []string{"readcodetheme"},
	Short:   "Read the theme for the code preview",
	Long:    "Read the theme for the code preview",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		conf := configuration.ReadConfiguration()
		cmd.Printf("Current theme: %s", conf.CodeThemeName.String())
	},
}

func init() {
	rootCmd.AddCommand(setThemeCmd)
	rootCmd.AddCommand(readThemeCmd)
	rootCmd.AddCommand(setCodeThemeCmd)
	rootCmd.AddCommand(readCodeThemeCmd)
}
