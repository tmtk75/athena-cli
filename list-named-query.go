package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ListNamedQueriesCmd)
}

var ListNamedQueriesCmd = &cobra.Command{
	Use:   "list-named-queries",
	Short: "List named queries",
	Run: func(cmd *cobra.Command, args []string) {
		w := NewSession()
		err := w.ListNamedQueries()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func (sess *Session) ListNamedQueries() error {
	r, err := sess.athenaClient.ListNamedQueries(sess.ctx, &athena.ListNamedQueriesInput{WorkGroup: aws.String(sess.profile.WorkGroup())})
	if err != nil {
		return err
	}
	l := make([]*athena.GetNamedQueryOutput, 0)
	for _, id := range r.NamedQueryIds {
		r, err := sess.athenaClient.GetNamedQuery(sess.ctx, &athena.GetNamedQueryInput{NamedQueryId: aws.String(id)})
		if err != nil {
			return err
		}
		l = append(l, r)
	}
	b, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", string(b))
	return nil
}
