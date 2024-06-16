/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package build

import (
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score/loader"
	"github.com/Catizard/lampghost/internal/service"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
	"github.com/spf13/cobra"
)

// TODO: Support LR2, which needs some hack
// TODO: Tags are built on courses and only courses' info would be taken in consideration.
// And so, we cannot build tags like "25% HC insane 8".
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
			msg := "Multiple rivals matched, please choose one"
			filter := rival.RivalInfoFilter{Name: null.StringFrom(rivalName)}
			rivalInfo, err := service.RivalInfoService.ChooseOneRival(msg, filter)
			if err != nil {
				log.Fatal(err)
			}
			if err := loader.LoadRivalData(rivalInfo); err != nil {
				log.Fatal(err)
			}
			courseInfoArr, _, err := service.CourseInfoService.FindCourseInfoList(difftable.CourseInfoFilter{})
			if err != nil {
				log.Fatal(err)
			}
			if err := service.RivalTagService.BuildTags(rivalInfo, courseInfoArr); err != nil {
				log.Fatal(err)
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
