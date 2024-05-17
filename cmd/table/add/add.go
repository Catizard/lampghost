/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package add

import (
	"fmt"
	"strings"

	"github.com/Catizard/lampghost/internel/remote"
	"github.com/Catizard/lampghost/internel/vo"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "add bms difficult url, ignore if already exists in table_header.json",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		if !strings.HasSuffix(url, ".json") {
			panic("only .json format url is supported, sorry :(")
		}
		// 1. Fetch difficult table header
		dth := &vo.DiffTableHeader{}
		remote.FetchJson(url, dth)
		if aliasName, err := cmd.LocalFlags().GetString("alias"); err == nil {
			dth.Alias = aliasName
		}
		// 2. Add difficult table header
		err := dth.AddDiffTable()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s loaded", dth.Name)
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
