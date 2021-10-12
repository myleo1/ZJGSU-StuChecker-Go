package cmd

import (
	"github.com/myleo1/go-core-kit/init/initkit"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		initkit.BindFlags(cmd)
	},
}
