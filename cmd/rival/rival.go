/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package rival

import (
	"github.com/Catizard/lampghost/cmd/rival/add"
	"github.com/Catizard/lampghost/cmd/rival/del"
	"github.com/Catizard/lampghost/cmd/rival/edit"
	rinit "github.com/Catizard/lampghost/cmd/rival/rinit"
	"github.com/Catizard/lampghost/cmd/rival/sync"
	"github.com/Catizard/lampghost/cmd/rival/tags"
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
	RivalCmd.AddCommand(del.DelCmd)
	RivalCmd.AddCommand(edit.EditCmd)	
}
