package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/spf13/cobra"
)

var ListWorkGroupsCmd = &cobra.Command{
	Use:   "list-work-groups",
	Short: "List all work groups",
	Run: func(cmd *cobra.Command, args []string) {
		w := NewWorld()
		err := w.ListWorkGroups()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func (world *World) ListWorkGroups() error {
	r, err := world.athenaClient.ListWorkGroupsRequest(&athena.ListWorkGroupsInput{}).Send(world.ctx)
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
