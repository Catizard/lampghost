/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package add

import (
	"github.com/Catizard/lampghost/internel/rival"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add [rival name] [scorelog.db path]",
	Args:  cobra.ExactArgs(2),
	Short: "Register a rival's info and log",
	Run: func(cmd *cobra.Command, args []string) {
		rivalName := args[0]
		dbFilePath := args[1]

		rivalInfo := &rival.RivalInfo{
			Name:         rivalName,
			ScoreLogPath: dbFilePath,
		}
		err := rival.AddRivalInfo(rivalInfo)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
