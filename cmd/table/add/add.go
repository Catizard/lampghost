/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package add

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/Catizard/lampghost/internel/remote"
	"github.com/Catizard/lampghost/internel/vo"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "add bms difficult url, ignore if already exists in table_header.json",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		if !strings.HasSuffix(url, ".json") {
			panic("only .json format url is supported, sorry :(")
		}
		dth := &vo.DiffTableHeader{}
		remote.FetchJson(url, dth)
		fmt.Printf("dataUrl=%s\n", dth.DataUrl)
		fmt.Printf("name=%s\n", dth.Name)
		fileName := fmt.Sprintf("%s.json", dth.Name)
		// If data.json is already here, do nothing
		if _, err := os.Stat(fileName); err == nil {
			panic(fmt.Errorf("%s is already exists, if you want to update data.json, use sync command instead", fileName))
		} else if errors.Is(err, fs.ErrExist) {
			// unexpected...
			panic(err)
		}
		file, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		// download to file
		// TODO: if dataUrl is not start with http...
		resp, err := http.Get(dth.DataUrl)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		io.Copy(file, resp.Body)
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
