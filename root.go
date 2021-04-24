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

	f.String(keyWorkGroup, "primary", "Athena work group")
	f.String(keyOutputLocation, "", "S3 URL for output, for instance, s3://your-bucket-name/path/to")
	f.String(keyDatabaseName, "", "Athena database name")
	f.String(keyCatalogName, "AwsDataCatalog", "data catalog name")
	f.Bool(keyVerbose, false, "Work verbosely")
	f.String(keyProfile, "", "Profile name")
	f.Duration(keyTimeout, time.Second*5, "Timeout ex) 30s")
	f.Bool(keyDryRun, false, "Dry-run only printing templated query.")

	opts := []struct {
		key string
		env string
	}{
		{key: keyWorkGroup, env: "WORK_GROUP"},
		{key: keyOutputLocation, env: "OUTPUT_LOCATION"},
		{key: keyDatabaseName, env: "DATABASE_NAME"},
		{key: keyCatalogName, env: "CATALOG_NAME"},
		{key: keyVerbose, env: "VERBOSE"},
		{key: keyProfile, env: "PROFILE"},
		{key: keyTimeout, env: "TIMEOUT"},
		{key: keyDryRun},
	}
	for _, e := range opts {
		viper.BindPFlag(e.key, f.Lookup(e.key))
		viper.BindEnv(e.key, e.env)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/athena-cli")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if cerr, b := err.(viper.ConfigFileNotFoundError); !b {
			panic(cerr)
		}
	}
}

var RootCmd = &cobra.Command{
	Use: "athena-cli",
}

const (
	// global
	keyVerbose = "verbose"
	keyProfile = "profile"
	keyTimeout = "timeout"
	keyDryRun  = "dry-run"
	// each profile
	keyWorkGroup      = "work-group"
	keyOutputLocation = "output-location"
	keyDatabaseName   = "database-name"
	keyCatalogName    = "catalog-name"
)
