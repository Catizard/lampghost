/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package add

import (
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/service"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add [rival's name]",
	Args:  cobra.ExactArgs(1),
	Short: "Register a rival's info",
	Run: func(cmd *cobra.Command, args []string) {
		rivalName := args[0]
		scoreLogPath, err := cmd.Flags().GetString("scorelog")
		if err != nil {
			log.Error(err)
		}
		songDataPath, err := cmd.Flags().GetString("songdata")
		if err != nil {
			log.Error(err)
		}
		userDBPath, err := cmd.Flags().GetString("user")
		if err != nil {
			log.Error(err)
		}

		rivalInfo := &rival.RivalInfo{
			Name:            rivalName,
			ScoreLogPath:    null.StringFrom(scoreLogPath),
			SongDataPath:    null.StringFrom(songDataPath),
			LR2UserDataPath: null.StringFrom(userDBPath),
		}

		// Set path fields to null, otherwise blank string would be stored into database
		rivalInfo.BlankToNull()

		if err := service.RivalInfoService.InsertRivalInfo(rivalInfo); err != nil {
			panic(err)
		}
	},
}

func init() {
	AddCmd.Flags().String("songdata", "", "path to songdata.db(oraja)")
	AddCmd.Flags().String("scorelog", "", "path to scorelog.db(oraja)")
	AddCmd.Flags().String("user", "", "path to user.db(LR2)")
}
