/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package add

import (
	"github.com/Catizard/lampghost/internal/difftable"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add [the name of table to be added] ",
	Short: "Add one bms difficult table",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		// 1. Fetch difficult table header
		dth, err := difftable.FetchDiffTableHeader(url)
		if err != nil {
			panic(err)
		}
		if aliasName, err := cmd.LocalFlags().GetString("alias"); err == nil {
			dth.Alias = aliasName
		}
		// 2. Add difficult table header
		err = dth.SaveDiffTableHeader()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	AddCmd.Flags().StringP("alias", "a", "", "difficult table's alias, could be used as name in other commands")
}
