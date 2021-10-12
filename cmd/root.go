package cmd

import (
	"ZJGSU-StuChecker-Go/model"
	"ZJGSU-StuChecker-Go/service/check"
	"github.com/myleo1/go-core-kit/init/initkit"
	"github.com/myleo1/go-core-kit/library/jsonkit"
	"github.com/myleo1/go-core-kit/service/configkit"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "ZJGSU-StuChecker-Go",
	Run: func(cmd *cobra.Command, args []string) {
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
