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
		sess := NewSession()
		fmt.Printf("session: %v\n", sess)

		var v *viper.Viper
		rootKey := "accounts"

		// use it if given explicitly.
		if p := viper.GetString(keyProfile); p != "" {
			pkey := fmt.Sprintf("%s.%s", rootKey, p)
			v = viper.Sub(pkey)
			if v == nil {
				log.Fatalf("no found profile, %v", pkey)
			}
			logger.Printf("use a profile, %v", pkey)
		} else {
			pkey := fmt.Sprintf("%s.%s", rootKey, sess.accountId)
			v = viper.Sub(pkey)
			if v == nil {
				log.Fatalf("no found profile, %v", pkey)
			}
			logger.Printf("use a profile, %v", pkey)
		}

		fmt.Printf("catalog-name: %s\n", v.GetString(keyCatalogName))
		fmt.Printf("work-group: %s\n", v.GetString(keyWorkGroup))
		fmt.Printf("database-name: %s\n", v.GetString(keyDatabaseName))
		//fmt.Printf("output-location: %s\n", viper.GetString(keyOutputLocation))
	},
}
