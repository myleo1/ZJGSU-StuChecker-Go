package cmd

import (
	"ZJGSU-StuChecker-Go/model"
	"ZJGSU-StuChecker-Go/service/check"
	"fmt"
	"github.com/mizuki1412/go-core-kit/init/initkit"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(yzyCmd)
}

var yzyCmd = &cobra.Command{
	Use: "yzy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`    _         _           __   _________   __
   / \  _   _| |_ ___     \ \ / /__  /\ \ / /
  / _ \| | | | __/ _ \ ____\ V /  / /  \ V /
 / ___ \ |_| | || (_) |_____| |  / /_   | |
/_/   \_\__,_|\__\___/      |_| /____|  |_|`)
		initkit.BindFlags(cmd)
		//打卡前准备
		checkerList := &model.CheckerList{}
		_ = jsonkit.ParseObj(jsonkit.ToString(configkit.Get("checker", model.CheckerList{})), checkerList)
		////打卡
		check.BeginYzy(*checkerList)
	},
}
