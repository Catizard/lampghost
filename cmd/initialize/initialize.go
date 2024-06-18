/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package initialize

import (
	"log"

	"github.com/Catizard/lampghost/internal/sqlite"
	"github.com/spf13/cobra"
)

// Initialize database that lampghost would use
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init lampghost application's database",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Do nothing, only purpose is to override root's hook
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := sqlite.InitializeDatabase(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
}
