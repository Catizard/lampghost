/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package del

import (
	"github.com/Catizard/lampghost/internel/rival"
	"github.com/Catizard/lampghost/internel/tui/choose"
	"github.com/spf13/cobra"
)

var DelCmd = &cobra.Command{
	Use:   "del [name]",
	Short: "Delete a rival",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		rivalInfo, err := rival.QueryExactlyRivalInfo(name)
		if err != nil {
			panic(err)
		}
		secq := choose.OpenYesOrNoChooseTui("Do you really want to delete this rival?")
		if secq {
			rival.DeleteRivalInfo(rivalInfo.Id)
		}
	},
}

func init() {
}
