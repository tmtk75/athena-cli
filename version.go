package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Commit string
var Version string

func init() {
	RootCmd.AddCommand(VersionCmd)
}

var VersionCmd = &cobra.Command{
	Use:  "version",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %v, commit: %v\n", Version, Commit)
	},
}
