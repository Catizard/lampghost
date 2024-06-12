/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package add

import (
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/sqlite/service"
	"github.com/guregu/null/v5"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add [rival name] [scorelog.db path] [songData.db path]",
	Args:  cobra.ExactArgs(3),
	Short: "Register a rival's info",
	Run: func(cmd *cobra.Command, args []string) {
		rivalName := args[0]
		scoreLogPath := args[1]
		songDataPath := args[2]

		rivalInfo := &rival.RivalInfo{
			Name:         rivalName,
			ScoreLogPath: null.StringFrom(scoreLogPath),
			SongDataPath: null.StringFrom(songDataPath),
		}

		if err := service.RivalInfoService.InsertRivalInfo(rivalInfo); err != nil {
			panic(err)
		}
	},
}

func init() {
}
