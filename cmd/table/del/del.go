/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package del

import (
	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/spf13/cobra"
)

var DelCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a difficult table",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.LocalFlags().GetString("name")
		if err != nil {
			panic(err)
		}
		var dth difftable.DiffTableHeader
		var qErr error
		if len(name) == 0 {
			// Special case: doesn't specify a name
			// In this case, all tables would be printed
			dth, qErr = difftable.AllDiffTableHeaderWithChoices()
		} else {
			dth, qErr = difftable.QueryDiffTableHeaderByNameWithChoices(name)
		}
		if qErr != nil {
			panic(qErr)
		}

		if err := dth.DeleteDiffTableHeader(); err != nil {
			panic(err)
		}
	},
}

func init() {
	DelCmd.Flags().StringP("name", "n", "", "Specify the deleting table's name")
}
