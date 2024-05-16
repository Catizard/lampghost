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

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "add bms difficult url, ignore if already exists in table_header.json",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		if !strings.HasSuffix(url, ".json") {
			panic("only .json format url is supported, sorry :(")
		}
		dth := &vo.DiffTableHeader{}
		remote.FetchJson(url, dth)
		fmt.Printf("dataUrl=%s\n", dth.DataUrl)
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
