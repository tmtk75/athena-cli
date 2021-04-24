package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(InfoCmd)
}

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print some values like workgroup, database-name, catalog-name, etc",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("catalog-name: %s\n", viper.GetString(keyCatalogName))
		fmt.Printf("work-group: %s\n", viper.GetString(keyWorkGroup))
		fmt.Printf("database-name: %s\n", viper.GetString(keyDatabaseName))
		//fmt.Printf("output-location: %s\n", viper.GetString(keyOutputLocation))
	},
}
