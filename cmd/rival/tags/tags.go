/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package tags

import (
	tags "github.com/Catizard/lampghost/cmd/rival/tags/add"
	"github.com/Catizard/lampghost/cmd/rival/tags/build"
	"github.com/spf13/cobra"
)

// tagsCmd represents the tags command
var TagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Commands on rival's tags",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	TagsCmd.AddCommand(tags.AddCmd)
	TagsCmd.AddCommand(build.BuildCmd)
}
