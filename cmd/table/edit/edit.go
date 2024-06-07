/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package edit

import (
	"log"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var EditCmd = &cobra.Command{
	Use:   "edit [table_name] [--alias alias_name]",
	Short: "Edit table's config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("TODO: table edit command")
	},
}

func init() {
	EditCmd.Flags().StringP("alias", "a", "", "table's alias name")
}
