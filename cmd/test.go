package cmd

import (
	"ZJGSU-StuChecker-Go/model"
	"github.com/mizuki1412/go-core-kit/init/initkit"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"math/rand"
	"time"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		initkit.BindFlags(cmd)
		//convert()
	},
}

//生成deviceId
func gen() string {
	str := "164"
	for ii := 0; ii < 17; ii++ {
		rand.Seed(time.Now().UnixNano())
		str += cast.ToString(rand.Intn(10))
	}
	return str
}

func convert() {
	res := make([]*model.Checker, 0)
	s, _ := filekit.ReadString("test.json")
	temp := gjson.Get(s, "id").Array()
	for _, v := range temp {
		tt := &model.Checker{
			Name:          v.Get("trueName").String(),
			Username:      v.Get("name").String(),
			Password:      v.Get("psswd").String(),
			Campus:        "金沙港",
			WechatPushKey: v.Get("wechatPushKey").String(),
			DeviceId:      gen(),
			Skip:          false,
		}
		res = append(res, tt)
	}
	r := jsonkit.ToString(res)
	_ = filekit.WriteFile("con.json", []byte(r))
}
