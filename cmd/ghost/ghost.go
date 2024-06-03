/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package ghost

import (
	"fmt"

	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/Catizard/lampghost/internel/rival"
	"github.com/Catizard/lampghost/internel/tui/choose"
	ghostTui "github.com/Catizard/lampghost/internel/tui/ghost"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// ghostCmd represents the ghost command
var GhostCmd = &cobra.Command{
	Use:   "ghost [self] [ghost]",
	Short: "Open ghost tui application",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		selfInfo := queryAndLoadRival(args[0])
		ghostInfo := queryAndLoadRival(args[1])

		// Difficult table header
		// TODO: give difftable argument
		dth, err := difftable.QueryDiffTableHeaderByNameWithChoices("insane")
		if err != nil {
			panic(err)
		}

		// Difficult table data
		diffTable, err := difftable.ReadDiffTable(dth.Name + ".json")
		if err != nil {
			panic(err)
		}

		// If tag flagged
		if b, err := cmd.Flags().GetBool("tag"); err != nil {
			panic(err)
		} else {
			if b {
				// You have to choose a tag first
				tag, err := ghostInfo.ChooseFromAllTags()
				if err != nil {
					panic(err)
				}
				log.Infof("Choosed %s, time=%d\n", tag.TagName, tag.TimeStamp)
			}
		}
		ghostTui.OpenGhostTui(&dth, diffTable, &selfInfo, &ghostInfo)
	},
}

func init() {
	GhostCmd.Flags().Bool("tag", false, "When flagged, only logs before the chosen tag would be used")
}

func queryAndLoadRival(rivalName string) rival.RivalInfo {
	rivalInfoArr, err := rival.QueryRivalInfo(rivalName)
	if err != nil {
		panic(err)
	}
	rivalNameArr := make([]string, 0)
	for _, r := range rivalInfoArr {
		rivalNameArr = append(rivalNameArr, r.String())
	}
	index := choose.OpenChooseTuiSkippable(rivalNameArr, fmt.Sprintf("Multiple rivals named %s, please choose one:", rivalName))
	rivalInfo := rivalInfoArr[index]
	if err := rivalInfo.LoadDataIfNil(); err != nil {
		panic(err)
	}
	return rivalInfo
}
