package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/spf13/cobra"
)

var GetWorkGroupCmd = &cobra.Command{
	Use:   "get-work-group <name>",
	Short: "Show a work group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w := NewWorld()
		err := w.GetWorkGroup(args[0])
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

func (world *World) GetWorkGroup(name string) error {
	r, err := world.athenaClient.GetWorkGroupRequest(&athena.GetWorkGroupInput{
		WorkGroup: aws.String(name),
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
