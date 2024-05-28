/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package table

import (
	"github.com/Catizard/lampghost/cmd/table/add"
	"github.com/Catizard/lampghost/cmd/table/del"
	"github.com/Catizard/lampghost/cmd/table/edit"
	"github.com/Catizard/lampghost/cmd/table/sync"
	"github.com/spf13/cobra"
)

var TableCmd = &cobra.Command{
	Use:   "table",
	Short: "Add or edit difficult table settings",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	TableCmd.AddCommand(sync.SyncCmd)
	TableCmd.AddCommand(add.AddCmd)
	TableCmd.AddCommand(edit.EditCmd)
	TableCmd.AddCommand(del.DelCmd)
}
