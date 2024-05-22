/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package ghost

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Catizard/lampghost/internel/difftable"
	"github.com/Catizard/lampghost/internel/ghost"
	"github.com/Catizard/lampghost/internel/rival"
	"github.com/spf13/cobra"
)

// TODO: below implementation is for testing core functions, the real ghost command should open a tui application
// ghostCmd represents the ghost command
var GhostCmd = &cobra.Command{
	Use:   "ghost [rival]",
	Short: "ghost",
	Args: cobra.ExactArgs(1),
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
		scoreLogPath := rivalInfo.ScoreLogPath
		fmt.Printf("scoreLogPath=%s\n", scoreLogPath)
		scoreLogArray, err := ghost.ReadScoreLogFromSqlite(scoreLogPath)
		if err != nil {
			panic(err)
		}
		fmt.Printf("ghost: %d logs loaded\n", len(scoreLogArray))

		scoreLogSha256Map := make(map[string][]ghost.ScoreLog)
		for _, v := range scoreLogArray {
			if _, ok := scoreLogSha256Map[v.Sha256]; !ok {
				scoreLogSha256Map[v.Sha256] = make([]ghost.ScoreLog, 0)
			}
			scoreLogSha256Map[v.Sha256] = append(scoreLogSha256Map[v.Sha256], v)
		}
		
		// Song data
		songDataPath := rivalInfo.SongDataPath
		fmt.Printf("songDataPath=%s\n", songDataPath)
		songDataArray, err := ghost.ReadSongDataFromSqlite(songDataPath)
		if err != nil {
			panic(err)
		}
		fmt.Printf("ghost: %d songs loaded\n", len(songDataArray))
		
		// Difficult table
		dthArr, err := difftable.QueryDifficultTableHeader("insane")
		if err != nil {
			panic(err)
		}
		if len(dthArr) != 1 {
			panic("what")
		}
		dth := dthArr[0]

		diffTableMap, err := difftable.ReadDiffTableLevelMap(dth.Name + ".json")			
		if err != nil {
			panic(err)
		}

		// Merge sha256 info from SongData.db to difficult table
		songDataMd5Map := make(map[string]ghost.SongData)
		for _, songData := range songDataArray {
			songDataMd5Map[songData.Md5] = songData
		}

		for _, arr := range diffTableMap {
			// fmt.Printf("len(arr)=%d\n", len(arr))
			for i, v := range arr {
				arr[i].Sha256 = songDataMd5Map[v.Md5].Sha256
			}
		}

		// Apply score logs on difficult table
		for _, arr := range diffTableMap {
			for i, v := range arr {
				// fmt.Printf("sha256=%s\n", v.Sha256)
				for _, scoreLog := range scoreLogSha256Map[v.Sha256] {
					arr[i].Lamp = max(v.Lamp, scoreLog.Clear)
				}
			}
		}

		sortedKeys := make([]string, 0)
		for k := range diffTableMap {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Slice(sortedKeys, func(i, j int) bool {
			l, errl := strconv.Atoi(sortedKeys[i])
			r, errr := strconv.Atoi(sortedKeys[j])
			if errl == nil && errr == nil {
				return l < r	
			} else {
				return sortedKeys[i] < sortedKeys[j]
			}
		})

		for _, k := range sortedKeys {
			hc := 0;
			for _, v := range diffTableMap[k] {
				if v.Lamp >= 6 {
					hc ++
				}
			}
			fmt.Printf("%s%s=%d/%d\n", dth.Symbol, k, hc, len(diffTableMap[k]))
			// for _, v := range diffTableMap[k] {
			// 	fmt.Printf("%s %d\n", v.Title, v.Lamp)
			// }
		}
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
