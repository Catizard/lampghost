/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package del

import (
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/service"
	"github.com/Catizard/lampghost/internal/tui/choose"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
	"github.com/spf13/cobra"
)

var DelCmd = &cobra.Command{
	Use:   "del [name]",
	Short: "Delete a rival",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		msg := "Multiple rivals matched, please choose one"
		filter := rival.RivalInfoFilter{NameLike: null.StringFrom(name)}
		rivalInfo, err := service.RivalInfoService.ChooseOneRival(msg, filter)
		if err != nil {
			log.Fatal(err)
		}
		secq := choose.OpenYesOrNoChooseTui("Do you really want to delete this rival?")
		if secq {
			service.RivalInfoService.DeleteRivalInfo(rivalInfo.Id)
		}
	},
}

func init() {
}
