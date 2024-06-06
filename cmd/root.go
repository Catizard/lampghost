/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package cmd

import (
	"os"

	"github.com/Catizard/lampghost/cmd/ghost"
	"github.com/Catizard/lampghost/cmd/initialize"
	"github.com/Catizard/lampghost/cmd/rival"
	"github.com/Catizard/lampghost/cmd/table"
	"github.com/Catizard/lampghost/cmd/version"
	"github.com/Catizard/lampghost/internal/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "lampghost",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Every command should check if lampghost.db file exist before executing
		// Expect init command, which creates lampghost.db
		// Note: If any command need to define its PersistentPreRun function, common.CheckInitialize should be called
		config.CheckInitialize()
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
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
	rootCmd.AddCommand(rival.RivalCmd)
	rootCmd.AddCommand(ghost.GhostCmd)
	rootCmd.AddCommand(initialize.InitCmd)
}
