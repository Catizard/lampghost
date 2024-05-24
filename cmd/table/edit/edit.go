/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package edit

import (
	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var EditCmd = &cobra.Command{
	Use:   "edit [table_name] [--alias alias_name]",
	Short: "Edit table's config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tableName := args[0]
		// TODO: give choices when need
		dt := difftable.QueryDifficultTableHeaderExactlyOne(tableName)
		aliasName := cmd.Flag("alias").Value.String()

		anyChange := false
		if len(aliasName) > 0 {
			if dt.Alias == aliasName {
				anyChange = true
			}
			dt.Alias = aliasName
		}

		if anyChange {
			panic("TODO: table edit command")
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	EditCmd.Flags().StringP("alias", "a", "", "table's alias name")
}
