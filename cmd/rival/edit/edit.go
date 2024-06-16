/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package edit

import (
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/service"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
	"github.com/spf13/cobra"
)

var EditCmd = &cobra.Command{
	Use:   "edit [rival name]",
	Short: "Edit one rival's configuration",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		msg := "Multiple rival matched, please choose one"
		rivalInfo, err := service.RivalInfoService.ChooseOneRival(msg, rival.RivalInfoFilter{NameLike: null.StringFrom(name)})
		if err != nil {
			log.Fatal(err)
		}
		newName, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		scoreLogPath, err := cmd.Flags().GetString("scorelog")
		if err != nil {
			log.Fatal(err)
		}
		songDataPath, err := cmd.Flags().GetString("songdata")
		if err != nil {
			log.Fatal(err)
		}
		userDBPath, err := cmd.Flags().GetString("user")
		if err != nil {
			log.Fatal(err)
		}
		updater := rival.RivalInfoUpdater{
			Id: null.IntFrom(int64(rivalInfo.Id)),
			Name: null.StringFrom(newName),
			ScoreLogPath: null.StringFrom(scoreLogPath),
			SongDataPath: null.StringFrom(songDataPath),
			LR2UserDataPath: null.StringFrom(userDBPath),
		}
		if _, err := service.RivalInfoService.UpdateRivalInfo(rivalInfo.Id, updater); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	EditCmd.Flags().String("name", "", "new name for the rival")
	EditCmd.Flags().String("songdata", "", "path to songdata.db(oraja)")
	EditCmd.Flags().String("scorelog", "", "path to scorelog.db(oraja)")
	EditCmd.Flags().String("user", "", "path to user.db(LR2)")
}
