package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyListDatabasesJson = "list-databases.json"
)

func init() {
	RootCmd.AddCommand(ListDatabasesCmd)
	f := ListDatabasesCmd.PersistentFlags()
	f.Bool("json", false, "JSON")
	viper.BindPFlag(keyListDatabasesJson, f.Lookup("json"))
}

var ListDatabasesCmd = &cobra.Command{
	Use:   "list-databases",
	Short: "List all databases",
	Run: func(cmd *cobra.Command, args []string) {
		w := NewSession()
		err := w.ListDatabases()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func (sess *Session) ListDatabases() error {
	name := sess.profile.CatalogName()
	r, err := sess.athenaClient.ListDatabases(sess.ctx, &athena.ListDatabasesInput{CatalogName: aws.String(name)})
	if err != nil {
		return err
	}

	if sess.v.GetBool(keyListDatabasesJson) {
		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", string(b))
	} else {
		for _, e := range r.DatabaseList {
			if e.Description == nil {
				fmt.Printf("%v\n", *e.Name)
			} else {
				fmt.Printf("%v\t%v\n", *e.Name, *e.Description)
			}
		}
	}
	return nil
}
