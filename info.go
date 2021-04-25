package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(InfoCmd)
}

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print some values like workgroup, database-name, catalog-name, etc",
	Run: func(cmd *cobra.Command, args []string) {
		sess := NewSession()
		fmt.Printf("catalog-name: %s\n", sess.profile.CatalogName())
		fmt.Printf("work-group: %s\n", sess.profile.WorkGroup())
		fmt.Printf("database-name: %s\n", sess.profile.DatabaseName())
		fmt.Printf("output-location: %s\n", sess.profile.OutputLocation())
	},
}
