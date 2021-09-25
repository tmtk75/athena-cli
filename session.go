package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/viper"
)

type Session struct {
	ctx          context.Context
	athenaClient *athena.Client
	s3Client     *s3.Client
	v            *viper.Viper
	profile      *Profile
}

func NewSession() *Session {
	if viper.GetBool(keyVerbose) {
		logger.Printf = log.Printf
	}

	logger.Printf("version: %v, commit: %v", Version, Commit)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	to := viper.GetDuration(keyTimeout)
	logger.Printf("timeout: %v", to)
	ctx, cancel := context.WithTimeout(context.Background(), to)
	stssvc := sts.NewFromConfig(cfg)
	r, err := stssvc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalf("failed for get-caller-identity@sts, %v", err)
	}
	logger.Printf("aws-account-id: %v", *r.Account)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		select {
		case <-ch:
			cancel()
		}
	}()

	v := viper.GetViper()
	pv := profileViper(v, *r.Account)
	sess := &Session{
		ctx:          ctx,
		athenaClient: athena.NewFromConfig(cfg),
		s3Client:     s3.NewFromConfig(cfg),
		v:            v,
		profile:      &Profile{v: v, pv: pv},
	}
	logger.Printf("work-group: %s", sess.profile.WorkGroup())
	logger.Printf("catalog-name: %s", sess.profile.CatalogName())
	logger.Printf("database-name: %s", sess.profile.DatabaseName())
	logger.Printf("output-location: %s", sess.profile.OutputLocation())
	return sess
}
