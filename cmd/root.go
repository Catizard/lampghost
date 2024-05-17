/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package cmd

import (
	"os"

	"github.com/Catizard/lampghost/cmd/table"
	"github.com/Catizard/lampghost/cmd/version"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "lampghost",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// TODO: if start with no args, open tui for use
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(table.TableCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lampghost.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
