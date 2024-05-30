/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package build

import (
	"fmt"

	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build tags for one or more rival",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("build called")
	},
}

func init() {
	BuildCmd.Flags().BoolP("all", "a", false, "When given, build every rivals' tags")
	BuildCmd.Flags().StringP("rival", "r", "", "When given, build specified rival's tags")
	BuildCmd.MarkFlagsMutuallyExclusive("all", "rival")
}
