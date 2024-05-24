/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package rival

import (
	"github.com/Catizard/lampghost/cmd/rival/add"
	rinit "github.com/Catizard/lampghost/cmd/rival/rinit"
	"github.com/Catizard/lampghost/cmd/rival/tags"
	"github.com/Catizard/lampghost/cmd/table/sync"
	"github.com/spf13/cobra"
)

// rivalCmd represents the rival command
var RivalCmd = &cobra.Command{
	Use:   "rival",
	Short: "Add or edit rival settings",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RivalCmd.AddCommand(add.AddCmd)
	RivalCmd.AddCommand(rinit.InitCmd)
	RivalCmd.AddCommand(tags.TagsCmd)
	RivalCmd.AddCommand(sync.SyncCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rivalCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rivalCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
