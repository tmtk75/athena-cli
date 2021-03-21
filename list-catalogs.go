package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/spf13/cobra"
)

var ListCatalogsCmd = &cobra.Command{
	Use:   "list-catalogs",
	Short: "List all catalogs",
	Run: func(cmd *cobra.Command, args []string) {
		w := NewSession()
		err := w.ListCatalogs()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func (sess *Session) ListCatalogs() error {
	r, err := sess.athenaClient.ListDataCatalogsRequest(&athena.ListDataCatalogsInput{}).Send(sess.ctx)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", string(b))
	return nil
}
