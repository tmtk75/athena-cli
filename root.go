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
	f.Bool(keyVerbose, false, "Work verbosely")
	f.Duration(keyTimeout, time.Second*5, "Timeout ex) 30s")
	f.Bool(keyDryRun, false, "Dry-run only printing templated query.")

	opts := []struct {
		key string
		env string
	}{
		{key: keyWorkGroup, env: "WORK_GROUP"},
		{key: keyOutputLocation, env: "OUTPUT_LOCATION"},
		{key: keyDatabaseName, env: "DATABASE_NAME"},
		{key: keyVerbose, env: "VERBOSE"},
		{key: keyTimeout, env: "TIMEOUT"},
		{key: keyDryRun},
	}
	for _, e := range opts {
		viper.BindPFlag(e.key, f.Lookup(e.key))
		viper.BindEnv(e.key, e.env)
	}

	// Query
	RootCmd.AddCommand(QueryCmd)

	// Get
	RootCmd.AddCommand(GetCmd)

	// Result
	RootCmd.AddCommand(ResultCmd)

	// List
	RootCmd.AddCommand(ListCmd)

	// Version
	RootCmd.AddCommand(VersionCmd)
}

func initConfig() {
	viper.SetConfigName(".athena-cli")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}

var RootCmd = &cobra.Command{
	Use: "athena-cli",
}

const (
	keyVerbose        = "verbose"
	keyWorkGroup      = "work-group"
	keyOutputLocation = "output-location"
	keyDatabaseName   = "database-name"
	keyTimeout        = "timeout"
	keyDryRun         = "dry-run"
)
