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

var rootCmd = &cobra.Command{
	Use: "ZJGSU-StuChecker-Go",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`     _                    _               _
 ___| |_ _   _        ___| |__   ___  ___| | _____ _ __
/ __| __| | | |_____ / __| '_ \ / _ \/ __| |/ / _ \ '__|
\__ \ |_| |_| |_____| (__| | | |  __/ (__|   <  __/ |
|___/\__|\__,_|      \___|_| |_|\___|\___|_|\_\___|_|`)
		initkit.BindFlags(cmd)
		//打卡前准备
		checkerList := &model.CheckerList{}
		_ = jsonkit.ParseObj(jsonkit.ToString(configkit.Get("checker", model.CheckerList{})), checkerList)
		//打卡
		check.BeginCheck(*checkerList)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}
