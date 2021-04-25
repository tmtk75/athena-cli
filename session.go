package main

import (
	"context"
	"fmt"
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

	pv := profileViper(*r.Account)
	return &Session{
		ctx:          ctx,
		athenaClient: athena.NewFromConfig(cfg),
		s3Client:     s3.NewFromConfig(cfg),
		v:            viper.GetViper(),
		profile:      &Profile{v: viper.GetViper(), pv: pv},
	}
}

func profileViper(aid string) *viper.Viper {
	findSub := func(parent, key string) *viper.Viper {
		pkey := fmt.Sprintf("%s.%s", parent, key)
		v := viper.Sub(pkey)
		if v == nil {
			log.Fatalf("no found profile, %v", pkey)
		}
		logger.Printf("use a profile, %v", pkey)
		return v
	}

	// use it if given explicitly.
	if p := viper.GetString(keyProfile); p != "" {
		return findSub("profiles", p)
	}

	// New empty viper.
	return findSub("accounts", aid)
}

type Profile struct {
	v  *viper.Viper
	pv *viper.Viper
}

// return b if a is empty.
// return def if b is empty.
func (p *Profile) either(key, def string) string {
	a := p.v.GetString(key)
	b := p.pv.GetString(key)
	if a != "" {
		return a
	}
	if b != "" {
		return b
	}
	return def
}

func (p *Profile) WorkGroup() string {
	return p.either(keyWorkGroup, "primary")
}

func (p *Profile) CatalogName() string {
	return p.either(keyCatalogName, "AwsDataCatalog")
}

func (p *Profile) DatabaseName() string {
	return p.either(keyDatabaseName, "")
}

func (p *Profile) OutputLocation() string {
	return p.either(keyOutputLocation, "")
}
