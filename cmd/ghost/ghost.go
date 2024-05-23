/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package ghost

import (
	"github.com/Catizard/lampghost/internel/difftable"
	ghostTui "github.com/Catizard/lampghost/internel/ghost/tui"
	"github.com/Catizard/lampghost/internel/rival"
	"github.com/spf13/cobra"
)

// TODO: every steps should give choices
// ghostCmd represents the ghost command
var GhostCmd = &cobra.Command{
	Use:   "ghost [rival]",
	Short: "ghost",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rivalName := args[0]
		rivalArr, err := rival.QueryRivalInfo(rivalName)
		if err != nil {
			panic(err)
		}
		if len(rivalArr) == 0 {
			panic("no such a rival")
		}
		if len(rivalArr) > 1 {
			panic("multiple rivals matched")
		}
		rivalInfo := rivalArr[0]

		if err := rivalInfo.LoadRivalScoreLog(); err != nil {
			panic(err)
		}

		if err := rivalInfo.LoadRivalSongData(); err != nil {
			panic(err)
		}

		// Difficult table header
		dth := difftable.QueryDifficultTableHeaderExactlyOne("insane")

		// Difficult table data
		diffTable, err := difftable.ReadDiffTable(dth.Name + ".json")
		if err != nil {
			panic(err)
		}
		ghostTui.OpenGhostTui(&dth, diffTable, &rivalInfo)
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ghostCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ghostCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
