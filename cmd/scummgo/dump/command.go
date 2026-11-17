package dump

import (
	"fmt"
	"os"

	"github.com/jsolana/scummgo"
	"github.com/spf13/cobra"
)

var output string

// Command dumps a ScummGo package (.idx/.dat) as canonical JSON.
var Command = &cobra.Command{
	Use:   "dump [package]",
	Short: "dump a ScummGo package as JSON (omit .idx/.dat suffix)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		base := args[0]
		dump, err := scummgo.DumpPackageFile(base)
		if err != nil {
			return err
		}
		data, err := scummgo.DumpPackageJSON(dump)
		if err != nil {
			return err
		}
		if output == "" || output == "-" {
			_, err = os.Stdout.Write(data)
			return err
		}
		if err := os.WriteFile(output, data, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "wrote %s (%d resources)\n", output, len(dump.Resources))
		return nil
	},
}

func init() {
	Command.Flags().StringVarP(
		&output, "output", "o", "-", "output JSON path (default: stdout)",
	)
}
