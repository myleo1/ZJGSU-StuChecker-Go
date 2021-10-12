package main

import (
	"ZJGSU-StuChecker-Go/cmd"
	"fmt"
	"github.com/myleo1/go-core-kit/init/initkit"
)

var (
	version   string
	date      string
	goVersion string
)

func main() {
	info := fmt.Sprintf("***ZJGSU-StuChecker-Go %s***\n***BuildDate %s***\n***%s***\n", version, date, goVersion)
	fmt.Print(info)
	initkit.LoadConfig()
	cmd.Execute()
}
