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

const (
	keyShowTablesJson = "show-tables.json"
)

func init() {
	RootCmd.AddCommand(ShowTablesCmd)
	f := ShowTablesCmd.PersistentFlags()
	f.Bool("json", false, "JSON")
	viper.BindPFlag(keyShowTablesJson, f.Lookup("json"))
}

var ShowTablesCmd = &cobra.Command{
	Use:   "show-tables [table-name]",
	Short: "Show all tables in .tsv or table metadata for given table name",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		w := NewSession()
		err := w.ShowTables(name)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func (sess *Session) ShowTables(tablename string) error {
	var (
		catalog = sess.profile.CatalogName()
		dbname  = sess.profile.DatabaseName()
	)
	logger.Printf("catalog-name: %v, database-name: %v", catalog, dbname)
	//
	if tablename != "" {
		r, err := sess.athenaClient.GetTableMetadata(sess.ctx, &athena.GetTableMetadataInput{
			CatalogName:  aws.String(catalog),
			DatabaseName: aws.String(dbname),
			TableName:    aws.String(tablename),
		})
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

	// Show tables.
	r, err := sess.athenaClient.ListTableMetadata(sess.ctx, &athena.ListTableMetadataInput{
		CatalogName:  aws.String(catalog),
		DatabaseName: aws.String(dbname),
	})
	if err != nil {
		return err
	}

	if sess.v.GetBool(keyShowTablesJson) {
		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", string(b))
		return nil
	}

	for _, e := range r.TableMetadataList {
		parts := []string{}
		for _, e := range e.PartitionKeys {
			parts = append(parts, fmt.Sprintf("%v:%v", *e.Name, *e.Type))
		}
		fmt.Printf("%v\t%v\t%v\t%v\t%v\n", *e.Name, *e.TableType, strings.Join(parts, ","), *e.CreateTime, *e.LastAccessTime)
	}
	return nil
}
