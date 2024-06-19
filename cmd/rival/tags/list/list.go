package list

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/service"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list [rival's name]",
	Short: "Prints one rival's tags",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rivalName := args[0]
		if len(rivalName) == 0 {
			log.Fatal("Please input rival's name")
		}
		msg := "Multiple rivals matched, please choose one"
		rivalInfo, err := service.RivalInfoService.ChooseOneRival(msg, rival.RivalInfoFilter{NameLike: null.StringFrom(rivalName)})
		if err != nil {
			log.Fatal(err)
		}
		tags, _, err := service.RivalTagService.FindRivalTagList(rival.RivalTagFilter{RivalId: null.IntFrom(int64(rivalInfo.Id))})
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range tags {
			fmt.Println(v.String())
		}
	},
}

func init() {
	ListCmd.Flags().StringP("rival", "r", "", "rival's name")
}
