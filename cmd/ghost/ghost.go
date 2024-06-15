/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package ghost

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/service"
	ghostTui "github.com/Catizard/lampghost/internal/tui/ghost"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
	"github.com/spf13/cobra"
)

// ghostCmd represents the ghost command
var GhostCmd = &cobra.Command{
	Use:   "ghost self-name ghost-name [--tag] [--table table-name]",
	Short: "Open ghost tui application",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// 1) Load self and ghost
		msg := "Multiple rivals matched, please choose one"
		selfInfo, err := service.RivalInfoService.ChooseOneRival(msg, rival.RivalInfoFilter{Name: null.StringFrom(args[0])})
		if err != nil {
			log.Fatal(err)
		}
		ghostInfo, err := service.RivalInfoService.ChooseOneRival(msg, rival.RivalInfoFilter{Name: null.StringFrom(args[1])})
		if err != nil {
			log.Fatal(err)
		}

		// 2) Load table data
		dthNameLike, err := cmd.Flags().GetString("table-name")
		if err != nil {
			log.Fatal(err)
		}

		filter := difftable.DiffTableHeaderFilter{}
		msg = "Choose one difftable to ghost"
		if len(dthNameLike) > 0 {
			filter.NameLike = null.StringFrom(dthNameLike)
			msg = fmt.Sprintf("Multiple tables matched with %s, choose one:", dthNameLike)
		}
		dth, err := service.DiffTableHeaderService.FindDiffTableHeaderListWithChoices(msg, filter)
		if err != nil {
			panic(err)
		}

		diffTable, err := difftable.ReadDiffTable(dth.DataLocation)
		if err != nil {
			panic(err)
		}

		// 3) Load log data

		// Note: Tags can only be applied on ghost
		if err := service.RivalInfoService.LoadRivalData(selfInfo); err != nil {
			log.Fatal(err)
		}

		var tag *rival.RivalTag
		// If tag flagged
		if b, err := cmd.Flags().GetBool("tag"); err != nil {
			panic(err)
		} else if b {
			// You have to choose a tag first
			tag, err = service.RivalTagService.ChooseOneTag("Choose one tag to ghost", rival.RivalTagFilter{})
			if err != nil {
				panic(err)
			}
			log.Infof("Choosed %s, time=%d\n", tag.TagName, tag.TimeStamp)
		}
		// TODO: make tags back
		if err := service.RivalInfoService.LoadRivalData(ghostInfo); err != nil {
			log.Fatal(err)
		}
		// 4) Open tui application
		ghostTui.OpenGhostTui(dth, diffTable, selfInfo, ghostInfo)
	},
}

func init() {
	GhostCmd.Flags().Bool("tag", false, "When flagged, only the logs that before the chosen tag would be used")
	GhostCmd.Flags().StringP("table-name", "T", "", "Filtering argument for table selection")
}
