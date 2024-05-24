/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package tags

import (
	tags "github.com/Catizard/lampghost/cmd/rival/tags/add"
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
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
