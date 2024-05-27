/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package ghost

import (
	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/Catizard/lampghost/internel/rival"
	ghostTui "github.com/Catizard/lampghost/internel/tui/ghost"
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
		dth := difftable.QueryDifficultTableHeaderExactlyOne("insane")

		// Difficult table data
		diffTable, err := difftable.ReadDiffTable(dth.Name + ".json")
		if err != nil {
			panic(err)
		}
		ghostTui.OpenGhostTui(&dth, diffTable, &selfInfo, &ghostInfo)
	},
}

func init() {
}

func queryAndLoadRival(rivalName string) rival.RivalInfo {
	rivalInfo, err := rival.QueryRivalInfo(rivalName)
	if err != nil {
		panic(err)
	}
	if err := rivalInfo.LoadRivalScoreLog(); err != nil {
		panic(err)
	}
	if err := rivalInfo.LoadRivalSongData(); err != nil {
		panic(err)
	}
	return rivalInfo
}
