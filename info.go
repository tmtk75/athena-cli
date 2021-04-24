package main

import (
	"fmt"
	"log"

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
		profile := viper.GetString(keyProfile)
		key := fmt.Sprintf("profiles.%s", profile)
		fmt.Printf("profile-key: %s\n", key)

		v := viper.Sub(key)
		if v == nil {
			log.Fatalf("no profile for given name, %s", profile)
		}
		fmt.Printf("catalog-name: %s\n", v.GetString(keyCatalogName))
		fmt.Printf("work-group: %s\n", v.GetString(keyWorkGroup))
		fmt.Printf("database-name: %s\n", v.GetString(keyDatabaseName))
		//fmt.Printf("output-location: %s\n", viper.GetString(keyOutputLocation))
	},
}
