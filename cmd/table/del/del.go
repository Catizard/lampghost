/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package del 

import (
	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/spf13/cobra"
)

var DelCmd = &cobra.Command{
	Use:   "del [table_name]",
	Short: "Delete a difficult table",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		dth, err := difftable.QueryDiffTableHeaderByNameWithChoices(name)
		if err != nil {
			panic(err)
		}
		if err := dth.DeleteDiffTableHeader(); err != nil {
			panic(err)
		}
	},
}

func init() {
}
