/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package add

import (
	"log"

	"github.com/Catizard/lampghost/internal/sqlite/service"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add [the name of table to be added] ",
	Short: "Add one bms difficult table",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		aliasName, err := cmd.LocalFlags().GetString("alias")
		if err != nil {
			log.Fatal(err)
		}
		_, err = service.DiffTableHeaderService.FetchAndSaveDiffTableHeader(url, aliasName)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	AddCmd.Flags().StringP("alias", "a", "", "difficult table's alias, could be used as name in other commands")
}
