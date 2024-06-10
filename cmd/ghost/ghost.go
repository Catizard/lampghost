/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package ghost

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/sqlite/service"
	ghostTui "github.com/Catizard/lampghost/internal/tui/ghost"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
	"github.com/spf13/cobra"
)

// ghostCmd represents the ghost command
var GhostCmd = &cobra.Command{
	Use:   "ghost [self] [ghost]",
	Short: "Open ghost tui application",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		selfInfo, err := queryAndLoadRival(args[0])
		if err != nil {
			log.Fatal(err)
		}
		ghostInfo, err := queryAndLoadRival(args[1])
		if err != nil {
			log.Fatal(err)
		}

		// Difficult table header
		// TODO: give difftable argument
		dthNameLike := "insane"
		filter := difftable.DiffTableHeaderFilter{
			NameLike: null.StringFrom(dthNameLike),
		}
		msg := fmt.Sprintf("Multiple tables matched with %s, choose one:", dthNameLike)
		dth, err := service.DiffTableHeaderService.FindDiffTableHeaderListWithChoices(msg, filter)
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
		ghostTui.OpenGhostTui(dth, diffTable, selfInfo, ghostInfo)
	},
}

func init() {
	GhostCmd.Flags().Bool("tag", false, "When flagged, only logs before the chosen tag would be used")
}

func queryAndLoadRival(rivalName string) (*rival.RivalInfo, error) {
	msg := "Multiple rivals matched, please choose one"
	filter := rival.RivalInfoFilter{Name: null.StringFrom(rivalName)}
	rivalInfo, err := service.RivalInfoService.ChooseOneRival(msg, filter)
	if err != nil {
		return nil, err
	}
	if err := rivalInfo.LoadDataIfNil(); err != nil {
		return nil, err
	}
	return rivalInfo, nil
}
