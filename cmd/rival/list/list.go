/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package list

import (
	"fmt"
	"log"

	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/service"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Prints all rivals",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		rs, _, err  := service.RivalInfoService.FindRivalInfoList(rival.RivalInfoFilter{})
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range rs {
			fmt.Println(v.String())
		}
	},
}

func init() {
}
