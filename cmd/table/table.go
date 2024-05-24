/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package table

import (
	"github.com/Catizard/lampghost/cmd/table/add"
	"github.com/Catizard/lampghost/cmd/table/edit"
	"github.com/Catizard/lampghost/cmd/table/sync"
	"github.com/spf13/cobra"
)

// tableCmd represents the table command
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
