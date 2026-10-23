package main

import (
	"fmt"
	"os"

	"github.com/jsolana/scummgo/cmd/scummgo/dump"
	"github.com/jsolana/scummgo/cmd/scummgo/pack"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "scummgo",
	Short: "scummgo is a point and click toolkit for creating adventure games.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cmd.AddCommand(pack.Command)
	cmd.AddCommand(dump.Command)
}
