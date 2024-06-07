/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package del

import (
	"fmt"
	"log"

	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/sqlite/service"
	"github.com/Catizard/lampghost/internal/tui/choose"
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
		filter := difftable.DiffTableHeaderFilter{
			Name: &name,
		}
		msg := fmt.Sprintf("Multiple tables matched with %s, choose one to delete:", name)
		dth, err := service.DiffTableHeaderService.FindDiffTableHeaderListWithChoices(msg, filter)
		if err != nil {
			log.Fatal(err)
		}
		if b := choose.OpenYesOrNoChooseTui(fmt.Sprintf("Delete %s?", dth.String())); b {
			if err := service.DiffTableHeaderService.DeleteDifftableHeader(dth.Id); err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	DelCmd.Flags().StringP("name", "n", "", "Specify the deleting table's name")
}
