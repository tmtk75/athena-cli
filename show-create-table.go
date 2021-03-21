package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ShowCreateTableCmd)
}

var ShowCreateTableCmd = &cobra.Command{
	Use:     "show-create-table table-anme",
	Short:   "Show CREATE TABLE DDL for given table",
	Example: ``,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w := NewSession()
		q := fmt.Sprintf("SHOW CREATE TABLE %s", args[0])
		Query(w, q)
	},
}

func Query(w *Session, q string) {
	s, err := w.Query(q)
	if err != nil {
		log.Fatalf("%v", err)
	}

	var m struct {
		Rows []map[string]string
	}
	err = json.Unmarshal([]byte(s), &m)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if len(m.Rows) == 0 {
		log.Fatalf("no rows")
	}
	for k, _ := range m.Rows[0] {
		fmt.Printf("%v\n", k)
	}
	for _, e := range m.Rows {
		for _, v := range e {
			fmt.Printf("%v\n", v)
		}
	}
}
