package cmd

import (
	"flag"
	"os"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cmd/instance"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func genS3ImportsCmd(name, version string, beatCreator beat.Creator, runFlags *pflag.FlagSet) *cobra.Command {
	runCmd := cobra.Command{
		Use:   "s3imports",
		Short: "Imports s3 objects to ElasticSearch and stops",
		Run: func(cmd *cobra.Command, args []string) {
			err := instance.Run(instance.Settings{Name: name, Version: version}, beatCreator)
			if err != nil {
				os.Exit(1)
			}
		},
	}

	// Run subcommand flags, only available to *beat run
	runCmd.Flags().AddGoFlag(flag.CommandLine.Lookup("N"))
	runCmd.Flags().AddGoFlag(flag.CommandLine.Lookup("httpprof"))
	runCmd.Flags().AddGoFlag(flag.CommandLine.Lookup("cpuprofile"))
	runCmd.Flags().AddGoFlag(flag.CommandLine.Lookup("memprofile"))

	if runFlags != nil {
		runCmd.Flags().AddFlagSet(runFlags)
	}

	return &runCmd
}
