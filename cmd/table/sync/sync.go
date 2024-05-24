/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package sync

import (
	"fmt"

	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var SyncCmd = &cobra.Command{
	Use:   "sync [difficult table's name]",
	Short: "sync one specified difficult table's data.\nHint: alias could be used too",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		matchedArray, err := difftable.QueryDifficultTableHeader(name)
		if err != nil {
			panic(err)
		}
		if len(matchedArray) == 0 {
			panic("no such a difficult table")
		}
		// TODO: give options to choose when multi header was matched
		if len(matchedArray) > 1 {
			panic(fmt.Errorf("more than 1 difftable matched on %s", name))
		}

		// Assume there is only one matched result
		matchedHeader := matchedArray[0]
		// Call sync function on specified header
		if err = matchedHeader.SyncDifficultTable(); err != nil {
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
