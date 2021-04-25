package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initConfig)

	// Root
	f := RootCmd.PersistentFlags()

	// profile options
	f.String(keyProfile, "", "Profile name")
	f.String(keyProfileWorkGroup, "", "Athena work group")
	f.String(keyProfileOutputLocation, "", "S3 URL for output, for instance, s3://your-bucket-name/path/to")
	f.String(keyProfileDatabaseName, "", "Athena database name")
	f.String(keyProfileCatalogName, "", "data catalog name")

	// global options
	f.Bool(keyVerbose, false, "Work verbosely")
	f.Duration(keyTimeout, time.Second*30, "Timeout ex) 30s, 1m")
	f.Bool(keyDryRun, false, "Dry-run only printing templated query.")

	opts := []struct {
		key string
		env string
	}{
		{key: keyProfile, env: "PROFILE"},
		{key: keyProfileWorkGroup, env: "WORK_GROUP"},
		{key: keyProfileOutputLocation, env: "OUTPUT_LOCATION"},
		{key: keyProfileDatabaseName, env: "DATABASE_NAME"},
		{key: keyProfileCatalogName, env: "CATALOG_NAME"},
		{key: keyVerbose, env: "VERBOSE"},
		{key: keyTimeout, env: "TIMEOUT"},
		{key: keyDryRun},
	}
	for _, e := range opts {
		viper.BindPFlag(e.key, f.Lookup(e.key))
		viper.BindEnv(e.key, e.env)
	}
}

func initConfig() {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.config/athena-cli")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if cerr, b := err.(viper.ConfigFileNotFoundError); !b {
			panic(cerr)
		}
	}

	//
	viper.MergeConfigMap(v.AllSettings())
}

var RootCmd = &cobra.Command{
	Use: "athena-cli",
}

const (
	// global
	keyVerbose = "verbose"
	keyTimeout = "timeout"
	keyDryRun  = "dry-run"
)
