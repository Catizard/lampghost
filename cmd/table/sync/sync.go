/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package sync

import (
	"github.com/Catizard/lampghost/internal/difftable"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var SyncCmd = &cobra.Command{
	Use:   "sync [difficult table's name]",
	Short: "sync one specified difficult table's data.",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		dth, err := difftable.QueryDiffTableHeaderByNameWithChoices(name)
		if err != nil {
			panic(err)
		}
		// Call sync function on specified header
		if err = dth.SyncDifficultTable(); err != nil {
			panic(err)
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
