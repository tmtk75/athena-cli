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

var ListDatabasesCmd = &cobra.Command{
	Use:   "list-databases",
	Short: "List all databases",
	Run: func(cmd *cobra.Command, args []string) {
		w := NewWorld()
		err := w.List()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func (world *World) ListDatabases() error {
	var (
		wg = viper.GetString(keyWorkGroup)
	)

	r, err := world.athenaClient.ListQueryExecutionsRequest(&athena.ListQueryExecutionsInput{WorkGroup: aws.String(wg)}).Send(world.ctx)
	if err != nil {
		return err
	}

	all := make([]*athena.QueryExecution, 0)
	count := 0
	for _, e := range r.QueryExecutionIds {
		if count >= viper.GetInt(keyListLimit) {
			break
		}
		r, err := world.athenaClient.GetQueryExecutionRequest(&athena.GetQueryExecutionInput{QueryExecutionId: aws.String(e)}).Send(world.ctx)
		if err != nil {
			log.Printf("%v", err)
		}
		all = append(all, r.QueryExecution)
		count++
	}

	//sort.SliceStable(all, func(i, j int) bool {
	//	//return all[i].ResultConfiguration.
	//})

	b, err := json.Marshal(all)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", string(b))
	return nil
}
