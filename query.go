package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	f := QueryCmd.PersistentFlags()

	f.Bool(keyQuerySuppressWait, false, "Suppress waiting for completion of query execution.")
	f.Duration(keyQueryWaitDuration, time.Second*1, "Duration to wait for completion of query execution.")
	f.Bool(keyQueryResultCsv, false, "Print result in CSV (raw format in the s3 bucket).")
	f.String(keyQueryValues, "{}", "A map JSON string for templating.")

	opts := []struct{ key string }{
		{key: keyQuerySuppressWait},
		{key: keyQueryWaitDuration},
		{key: keyQueryResultCsv},
		{key: keyQueryValues},
	}
	for _, e := range opts {
		viper.BindPFlag(e.key, f.Lookup(e.key))
	}
}

const (
	keyQuerySuppressWait = "query.suppress-wait"
	keyQueryWaitDuration = "query.wait-duration"
	keyQueryResultCsv    = "query.result-csv"
	keyQueryValues       = "query.values"
)

var QueryCmd = &cobra.Command{
	Use:   "query [flags] [query-string]",
	Short: "Execute query",
	Example: `  Regarding supported DDL, see https://docs.aws.amazon.com/athena/latest/ug/language-reference.html

  athena-cli query "select * from cloudtrail_logs where useridentity.principalid like '%yourname%'"
`,
	Run: func(cmd *cobra.Command, args []string) {
		w := NewWorld()

		var q string
		if len(args) > 0 {
			q = args[0]
		} else {
			b, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("%v", err)
			}
			q = string(b)
		}

		q, err := templ(q)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if viper.GetBool(keyDryRun) {
			fmt.Printf("%v\n", q)
			return
		}

		err = w.Query(q)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func templ(q string) (string, error) {
	t, err := template.New("").Parse(q)
	if err != nil {
		return "", err
	}
	var a map[string]interface{}
	if err := json.Unmarshal([]byte(viper.GetString(keyQueryValues)), &a); err != nil {
		return "", fmt.Errorf("failed to unmarshal. %w", err)
	}
	logger.Printf("%v: %v", keyQueryValues, a)

	b := bytes.NewBuffer([]byte{})
	err = t.Execute(b, a)
	if err != nil {
		return "", fmt.Errorf("failed to execute template. %w", err)
	}
	return b.String(), nil
}

func (world *World) Query(query string) error {
	var (
		wg     = viper.GetString(keyWorkGroup)
		loc    = viper.GetString(keyOutputLocation)
		dbname = viper.GetString(keyDatabaseName)
	)
	r, err := world.ExecuteQuery(wg, dbname, loc, query)
	if err != nil {
		return err
	}

	if viper.GetBool(keyQuerySuppressWait) {
		fmt.Printf("%v\n", *r.QueryExecutionId)
		return nil
	}

	if err := world.WaitExecution(r.QueryExecutionId); err != nil {
		return err
	}

	logger.Printf("query-execution-id: %v", *r.QueryExecutionId)

	var s string
	if viper.GetBool(keyQueryResultCsv) {
		s, err = world.GetObject(*r.QueryExecutionId, []string{".txt", ".csv"})
	} else {
		s, err = world.GetResult(*r.QueryExecutionId)
		if e, ok := err.(NoRows); ok {
			logger.Printf("%v\n", e)
			return nil
		}

	}
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", s)
	return nil
}

func (world *World) ExecuteQuery(wg, dbname, loc, query string) (*athena.StartQueryExecutionResponse, error) {
	qc := &athena.QueryExecutionContext{Database: aws.String(dbname)}
	rc := &athena.ResultConfiguration{OutputLocation: aws.String(loc)}
	i := athena.StartQueryExecutionInput{
		QueryString:           aws.String(query),
		QueryExecutionContext: qc,
		ResultConfiguration:   rc,
		WorkGroup:             aws.String(wg),
	}
	r, err := world.athenaClient.StartQueryExecutionRequest(&i).Send(world.ctx)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (world *World) WaitExecution(id *string) error {
	for {
		//to, _ := world.ctx.Deadline()
		//log.Printf("to: %v, before: %v", to, to.Before(time.Now()))
		select {
		case <-world.ctx.Done():
			bgctx, _ := context.WithTimeout(context.Background(), time.Second*3) // API in this func works with another context in short timeout.
			_, err := world.athenaClient.StopQueryExecutionRequest(&athena.StopQueryExecutionInput{QueryExecutionId: id}).Send(bgctx)
			if err != nil {
				return fmt.Errorf("failed to stop query execution in %v for %v, %w", world.ctx.Err(), *id, err)
			}
			logger.Printf("stop query execution for %v", *id)
			return fmt.Errorf("%w for execution-id, %v", world.ctx.Err(), *id)
		default:
		}

		d := viper.GetDuration(keyQueryWaitDuration)
		logger.Printf("wait for %v", d)
		time.Sleep(d)

		bgctx, _ := context.WithTimeout(context.Background(), time.Second*3) // API in this func works with another context in short timeout.
		r, err := world.athenaClient.GetQueryExecutionRequest(&athena.GetQueryExecutionInput{QueryExecutionId: id}).Send(bgctx)
		if err != nil {
			return err
		}

		//log.Printf("get-query-execution: %v", r)
		s := r.QueryExecution.Status.State
		if s == "QUEUED" || s == "RUNNING" {
			logger.Printf("query-execution.status.state: %v", s)
			continue
		}

		if s == "SUCCEEDED" {
			return nil
		}

		return fmt.Errorf("%v", *r.QueryExecution.Status.StateChangeReason)
	}
}
