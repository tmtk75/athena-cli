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
	RootCmd.AddCommand(ResultCmd)
}

var ResultCmd = &cobra.Command{
	Use:   `result [flags] <execution-id>`,
	Short: "Show query reulst",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w := NewSession()
		b, err := w.GetResult(args[0])
		if err != nil {
			if e, ok := err.(NoRows); ok {
				fmt.Printf("%v\n", e)
				return
			}
			log.Fatalf("%v", err)
		}
		fmt.Printf("%v", b)
	},
}

type NoRows struct {
	executionId string
}

func (r NoRows) Error() string {
	return fmt.Sprintf("no rows for %v", r.executionId)
}

func (sess *Session) GetResult(id string) (string, error) {
	r, err := sess.athenaClient.GetQueryResultsRequest(&athena.GetQueryResultsInput{QueryExecutionId: aws.String(id)}).Send(sess.ctx)
	if err != nil {
		return "", err
	}
	if r.ResultSet == nil {
		return "", fmt.Errorf("no query result for %v", id)
	}
	if len(r.ResultSet.Rows) == 0 {
		return "", NoRows{executionId: id}
	}

	s := []map[string]string{}
	keys := make([]string, len(r.ResultSet.Rows[0].Data))
	for i, k := range r.ResultSet.Rows[0].Data {
		keys[i] = *k.VarCharValue
	}

	for i := 1; i < len(r.ResultSet.Rows); i++ {
		m := make(map[string]string)
		r := r.ResultSet.Rows[i]
		for i, k := range keys {
			if r.Data[i].VarCharValue != nil {
				m[k] = *r.Data[i].VarCharValue
			}
		}
		s = append(s, m)
	}

	m := struct {
		ExecutionId string              `json:"execution_id"`
		Rows        []map[string]string `json:"rows"`
	}{
		ExecutionId: id,
		Rows:        s,
	}
	b, _ := json.Marshal(m)
	return string(b), nil
}
