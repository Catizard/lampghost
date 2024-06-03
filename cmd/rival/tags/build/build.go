/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package build

import (
	"github.com/Catizard/lampghost/internal/difftable"
	"github.com/Catizard/lampghost/internal/rival"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build tags for one or more rival",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		rivalName, err := cmd.Flags().GetString("rival")
		if err != nil {
			panic(err)
		}
		if len(rivalName) != 0 {
			rivalInfo, err := rival.QueryRivalInfoWithChoices(rivalName)
			if err != nil {
				panic(err)
			}
			courseInfoArr, err := difftable.QueryAllCourseInfo()
			if err != nil {
				panic(err)
			}
			err = rivalInfo.BuildTags(courseInfoArr)
			if err != nil {
				panic(err)
			}
		} else {
			log.Fatal("unsupported, too bad...")
		}
	},
}

func init() {
	BuildCmd.Flags().BoolP("all", "a", false, "When given, build every rivals' tags")
	BuildCmd.Flags().StringP("rival", "r", "", "When given, build specified rival's tags")
	BuildCmd.MarkFlagsMutuallyExclusive("all", "rival")
	BuildCmd.MarkFlagsOneRequired("all", "rival")
}