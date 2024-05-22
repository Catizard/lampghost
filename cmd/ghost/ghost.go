/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package ghost

import (
	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/Catizard/lampghost/internel/ghost"
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

		// Score log
		scoreLogArray, err := ghost.ReadScoreLogFromSqlite(rivalInfo.ScoreLogPath)
		if err != nil {
			panic(err)
		}

		// Song data
		songDataArray, err := ghost.ReadSongDataFromSqlite(rivalInfo.SongDataPath)
		if err != nil {
			panic(err)
		}

		// Difficult table header
		dth := difftable.QueryDifficultTableHeaderExactlyOne("insane")

		// Difficult table data
		diffTable, err := difftable.ReadDiffTable(dth.Name + ".json")
		if err != nil {
			panic(err)
		}
		ghost.OpenGhostTui(&dth, diffTable, songDataArray, scoreLogArray)
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
