package cmd

import (
	"fmt"
	"github.com/mizuki1412/go-core-kit/init/initkit"
	"github.com/mizuki1412/go-core-kit/library/mathkit"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		initkit.BindFlags(cmd)
		gen(10)
	},
}

//生成deviceId
func gen(num int) {
	for i := 0; i < num; i++ {
		str := "164"
		for ii := 0; ii < 17; ii++ {
			str += cast.ToString(mathkit.RandInt32(0, 10))
		}
		fmt.Println(str)
	}
}
