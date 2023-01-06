package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func CommandVersion() *cobra.Command {
	//TODO : Needs to go to config
	version := "DEV"

	return &cobra.Command{
		Use:   "version",
		Short: "Print the version and exit",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf(
				"Bee Version: %s\nGo Version: %s\nGo OS/ARCH: %s %s\n",
				version,
				runtime.Version(),
				runtime.GOOS,
				runtime.GOARCH,
			)
		},
	}
}
