package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	f := ListCmd.PersistentFlags()

	f.Int(keyListLimit, 5, "number of limitation to list exeuctions")
	f.Bool(keyJson, false, "print in raw JSON")
	f.Bool(keyHeader, false, "print header in .tsv")

	opts := []struct{ key string }{
		{key: keyListLimit},
		{key: keyJson},
		{key: keyHeader},
	}
	for _, e := range opts {
		viper.BindPFlag(e.key, f.Lookup(e.key))
	}
}

const (
	keyListLimit = "list.limit"
	keyJson      = "list.json"
	keyHeader    = "list.headser"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available executions",
	Run: func(cmd *cobra.Command, args []string) {
		w := NewSession()
		err := w.List()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func (sess *Session) List() error {
	var (
		wg = viper.GetString(keyWorkGroup)
	)

	r, err := sess.athenaClient.ListQueryExecutionsRequest(&athena.ListQueryExecutionsInput{WorkGroup: aws.String(wg)}).Send(sess.ctx)
	if err != nil {
		return err
	}

	all := make([]*athena.QueryExecution, 0)
	count := 0
	for _, e := range r.QueryExecutionIds {
		if count >= viper.GetInt(keyListLimit) {
			break
		}
		r, err := sess.athenaClient.GetQueryExecutionRequest(&athena.GetQueryExecutionInput{QueryExecutionId: aws.String(e)}).Send(sess.ctx)
		if err != nil {
			log.Printf("%v", err)
		}
		all = append(all, r.QueryExecution)
		count++
	}

	if viper.GetBool(keyJson) {
		b, err := json.MarshalIndent(all, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", string(b))
	} else {
		for _, e := range all {
			fmt.Printf("%v\n", strings.Join([]string{"QueryExecutionId", "SubmissionDateTime", "State", "WorkGroup", "StatementType", "Query"}, "\t"))
			start := *e.Status.SubmissionDateTime
			fmt.Printf("%v\t%v\t%v\t%v\t%v\t%q\n", *e.QueryExecutionId, start, e.Status.State, *e.WorkGroup, e.StatementType, *e.Query)
		}
	}
	return nil
}
