/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package sync

import (
	"log"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var SyncCmd = &cobra.Command{
	Use:   "sync [difficult table's name]",
	Short: "sync one specified difficult table's data.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("TODO: table sync command")
	},
}

func init() {
}
