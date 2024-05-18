/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package rename

import (
	"fmt"

	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var RenameCmd = &cobra.Command{
	Use:   "rename [difficult table original name]",
	Short: "Give alias name to the specified difficult table, alias name wouldn't be considered",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		// TODO: change signature of query function to support strict query
		headers, err := difftable.QueryDifficultTableHeader(name)
		if err != nil {
			panic(err)
		}
		strictMatched := make([]difftable.DiffTableHeader, 0)
		for _, v := range headers {
			if v.Name == name {
				strictMatched = append(strictMatched, v)
			}
		}
		if len(strictMatched) == 0 {
			panic(fmt.Errorf("no such a table named with %s", name))
		}
		// Should be unique on original name
		if len(strictMatched) > 1 {
			panic(fmt.Errorf(`multiple table matched with %s
			If you edit json file by hand, rollback or clear everything
			If it's not your fault, please report`, name))
		}
		panic("TODO: rename command")
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// renameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// renameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
