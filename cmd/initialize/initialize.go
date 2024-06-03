/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package initialize

import (
	"github.com/Catizard/lampghost/internal/config"
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
		config.InitLampGhost()
	},
}

func init() {
}
