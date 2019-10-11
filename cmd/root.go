package cmd

import (
	"flag"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/sequra/s3logsbeat/beater"

	cmd "github.com/elastic/beats/libbeat/cmd"
	"github.com/elastic/beats/libbeat/cmd/instance"
)

// Name of this beat
var Name = "s3logsbeat"

type BeatsRootCmd struct {
	*cmd.BeatsRootCmd
	S3ExportsCmd *cobra.Command
}

// RootCmd to handle beats cli
var RootCmd *BeatsRootCmd

func init() {
	var runFlags = pflag.NewFlagSet(Name, pflag.ExitOnError)
	runFlags.AddGoFlag(flag.CommandLine.Lookup("once"))
	runFlags.AddGoFlag(flag.CommandLine.Lookup("keepsqsmessages"))

	RootCmd = &BeatsRootCmd{
		BeatsRootCmd: cmd.GenRootCmdWithSettings(beater.NewS3logsbeat, instance.Settings{Name: Name, RunFlags: runFlags}),
		S3ExportsCmd: genS3ImportsCmd(Name, "", beater.NewS3importsbeat, nil),
	}
	RootCmd.AddCommand(RootCmd.S3ExportsCmd)
}
