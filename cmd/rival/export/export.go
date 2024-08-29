/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package export

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/common/clearType"
	"github.com/Catizard/lampghost/internal/common/source"
	"github.com/Catizard/lampghost/internal/data"
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
	"github.com/Catizard/lampghost/internal/data/score/loader"
	"github.com/Catizard/lampghost/internal/service"
	"github.com/Catizard/lampghost/internal/tui/ghost"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
	"github.com/spf13/cobra"
)

var (
	flBlock = lipgloss.NewStyle().SetString(" ").Background(lipgloss.Color("#FF0000"))
	ezBlock = lipgloss.NewStyle().SetString(" ").Background(lipgloss.Color("#00FF00"))
	nrBlock = lipgloss.NewStyle().SetString(" ").Background(lipgloss.Color("#0000FF"))
	hcBlock = lipgloss.NewStyle().SetString(" ").Background(lipgloss.Color("#00CCFF"))
	blBlock = lipgloss.NewStyle().SetString(" ").Background(lipgloss.Color("#000000"))
)

var ExportCmd = &cobra.Command{
	Use:   "export [rival's name] [--tag] [--table table's name]",
	Short: "Export one rival's lamp status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		msg := "Multiple rivals matched, please choose one"
		rivalInfo, err := service.RivalInfoService.ChooseOneRival(msg, rival.RivalInfoFilter{NameLike: null.StringFrom(args[0])})
		if err != nil {
			log.Fatal(err)
		}

		dthNameLike, err := cmd.Flags().GetString("table-name")
		if err != nil {
			log.Fatal(err)
		}

		filter := difftable.DiffTableHeaderFilter{}
		msg = "Choose one difftable to export"
		if len(dthNameLike) > 0 {
			filter.NameLike = null.StringFrom(dthNameLike)
			msg = fmt.Sprintf("Multiple tables matched with %s, choose one:", dthNameLike)
		}

		dth, err := service.DiffTableHeaderService.FindDiffTableHeaderListWithChoices(msg, filter)
		if err != nil {
			log.Fatal(err)
		}
		if err := dth.LoadData(); err != nil {
			log.Fatal(err)
		}

		var tag *rival.RivalTag = nil
		if b, err := cmd.Flags().GetBool("tag"); err != nil {
			log.Fatal(err)
		} else if b {
			filter := rival.RivalTagFilter{
				RivalId: null.IntFrom(int64(rivalInfo.Id)),
			}
			tag, err = service.RivalTagService.ChooseOneTag("Choose one tag to use", filter)
			if err != nil {
				log.Fatal(err)
			}
			log.Infof("Choosed tag [%s], time=%d", tag.TagName, tag.TimeStamp)
		}
		// TODO: support LR2
		if tag == nil {
			log.Debugf("Try loading [%s]'s data with no tag", rivalInfo.Name)
			if err := loader.LoadRivalData(rivalInfo); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Debugf("Try loading [%s]'s data with tag: [%s-%d]", rivalInfo.Name, tag.TagName, tag.TimeStamp)
			if err := loader.LoadTaggedRivalData(rivalInfo, tag); err != nil {
				log.Fatal(err)
			}
		}

		// Borrow calculate method from GhostTui package
		ghost.MergeLampFromScoreLog(dth.Data, rivalInfo.CommonScoreLog, func(data *difftable.DiffTableData, log *score.CommonScoreLog) {
			data.Lamp = max(data.Lamp, log.Clear)
		})

		// Build sorted levels slice
		sortedLevels := data.BuildSortedLevelList(dth)
		log.Debugf("sorted levels=%v", sortedLevels)

		// Level maps to song info array
		songDataMap := make(map[string][]difftable.DiffTableData)
		for _, v := range dth.Data {
			if _, ok := songDataMap[v.Level]; !ok {
				songDataMap[v.Level] = make([]difftable.DiffTableData, 0)
			}
			songDataMap[v.Level] = append(songDataMap[v.Level], v)
		}

		log.Infof("Exporting %s lamp status on %s", rivalInfo.Name, dth.Name)
		if b, err := cmd.Flags().GetBool("block"); err != nil {
			log.Fatal(err)
		} else if !b {
			for _, v := range sortedLevels {
				fmt.Printf("[%s%v]\n", dth.Symbol, v)
				for _, data := range songDataMap[v] {
					if data.Lamp == 0 {
						continue
					}
					// TODO: Give pattern argument?
					pattern := "[%s] %s"
					if rivalInfo.Prefer.String == source.Oraja {
						fmt.Printf(pattern+"\n", data.Title, clearType.ConvOraja(data.Lamp))
					} else if rivalInfo.Prefer.String == source.LR2 {
						fmt.Printf(pattern+"\n", data.Title, clearType.ConvLR2(data.Lamp))
					} else {
						log.Fatalf("unexpected prefer: %s", rivalInfo.Prefer.String)
					}
				}
			}
		} else {
			for _, v := range sortedLevels {
				fmt.Printf("[%s%v]\n", dth.Symbol, v)
				cntlamp := make(map[int]int)
				all := len(songDataMap[v])
				for _, data := range songDataMap[v] {
					cntlamp[int(data.Lamp)]++
				}
				width := 100
				buf := ""
				if cntlamp[clearType.Hard] > 0 {
					per := float64(cntlamp[clearType.Hard]) / float64(all)
					buf += hcBlock.Width(int(per * float64(width))).String()
				}
				if cntlamp[clearType.Normal] > 0 {
					per := float64(cntlamp[clearType.Normal]) / float64(all)
					buf += nrBlock.Width(int(per * float64(width))).String()
				}
				if cntlamp[clearType.Easy] > 0 {
					per := float64(cntlamp[clearType.Easy]) / float64(all)
					buf += ezBlock.Width(int(per * float64(width))).String()
				}
				if cntlamp[clearType.Failed] > 0 {
					per := float64(cntlamp[clearType.Failed]) / float64(all)
					buf += flBlock.Width(int(per * float64(width))).String()
				}
				fmt.Printf("%s\n", buf)
			}
		}

	},
}

func init() {
	ExportCmd.Flags().Bool("tag", false, "When flagged, only the logs that before the chosen tag would be used")
	ExportCmd.Flags().StringP("table-name", "T", "", "Filtering argument for table selection")
	ExportCmd.Flags().Bool("block", false, "When flagged, display color block other than text info")
	// TODO: Do not ignore NO_PLAY?
}
