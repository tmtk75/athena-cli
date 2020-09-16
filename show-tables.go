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
		w := NewWorld()
		err := w.ShowTables(name)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func (world *World) ShowTables(tablename string) error {
	var (
		catalog = viper.GetString(keyCatalogName)
		dbname  = viper.GetString(keyDatabaseName)
	)
	logger.Printf("catalog-name: %v, database-name: %v", catalog, dbname)
	//
	if tablename != "" {
		r, err := world.athenaClient.GetTableMetadataRequest(&athena.GetTableMetadataInput{
			CatalogName:  aws.String(catalog),
			DatabaseName: aws.String(dbname),
			TableName:    aws.String(tablename),
		}).Send(world.ctx)
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
	r, err := world.athenaClient.ListTableMetadataRequest(&athena.ListTableMetadataInput{
		CatalogName:  aws.String(catalog),
		DatabaseName: aws.String(dbname),
	}).Send(world.ctx)
	if err != nil {
		return err
	}

	if viper.GetBool(keyShowTablesJson) {
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
