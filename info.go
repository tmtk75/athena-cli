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
		v := sess.v
		fmt.Printf("timeout: %s\n", v.GetString(keyTimeout))
		fmt.Printf("catalog-name: %s\n", v.GetString(keyCatalogName))
		fmt.Printf("work-group: %s\n", v.GetString(keyWorkGroup))
		fmt.Printf("database-name: %s\n", v.GetString(keyDatabaseName))
		fmt.Printf("output-location: %s\n", v.GetString(keyOutputLocation))
	},
}
