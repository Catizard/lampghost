/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package tags

import (
	"fmt"

	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add [tag name] [time]",
	Short: "Add a tag to rival",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
	},
}

func init() {
}
