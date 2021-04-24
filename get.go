package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
)

func ini() {
	RootCmd.AddCommand(GetCmd)
}

var GetCmd = &cobra.Command{
	Use:   "get [flags] <execution-id>",
	Short: "Get an object from the s3 bucket with execution-id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		w := NewSession()
		b, err := w.GetObject(args[0], []string{".csv", ".txt"})
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("%v\n", b)
	},
}

func (sess *Session) GetObject(id string, exts []string) (string, error) {
	var (
		loc = sess.v.GetString(keyOutputLocation)
		err error
	)

	u, err := url.Parse(loc)
	if err != nil {
		return "", err
	}

	var res *s3.GetObjectOutput
	for _, ext := range exts {
		req := &s3.GetObjectInput{
			Bucket: aws.String(u.Host),
			Key:    aws.String(strings.TrimPrefix(u.Path, "/") + "/" + id + ext),
		}
		logger.Printf("%v", ext)
		res, err = sess.s3Client.GetObject(sess.ctx, req)
		if err != nil {
			logger.Printf("%v", err)
			continue
		}
		break
	}
	if res == nil {
		return "", fmt.Errorf("not found for %v in %v with %v, last-err: %w", id, loc, exts, err)
	}
	defer res.Body.Close()

	bb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bb), nil
}
