package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "v0.1.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Show version information of RealDebrid CLI tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
