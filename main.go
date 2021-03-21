package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

func init() {
	logger.Printf = func(format string, v ...interface{}) {}
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var logger struct {
	Printf func(format string, v ...interface{})
}

type Session struct {
	ctx          context.Context
	athenaClient *athena.Client
	s3Client     *s3.Client
}

func NewSession() *Session {
	if viper.GetBool(keyVerbose) {
		logger.Printf = log.Printf
	}

	logger.Printf("version: %v, commit: %v", Version, Commit)
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(keyTimeout))
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		select {
		case <-ch:
			cancel()
		}
	}()

	return &Session{
		ctx:          ctx,
		athenaClient: athena.New(cfg),
		s3Client:     s3.New(cfg),
	}
}
