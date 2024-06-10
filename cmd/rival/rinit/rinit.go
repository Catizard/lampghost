/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package rinit

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TODO: unused
// initCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Setup yours scorelog and songdata path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
