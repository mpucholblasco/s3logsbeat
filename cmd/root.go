package cmd

import (
	"flag"

	"github.com/spf13/pflag"

	"github.com/mpucholblasco/s3logsbeat/beater"

	cmd "github.com/elastic/beats/libbeat/cmd"
)

// Name of this beat
var Name = "s3logsbeat"

// RootCmd to handle beats cli
var RootCmd *cmd.BeatsRootCmd

func init() {
	var runFlags = pflag.NewFlagSet(Name, pflag.ExitOnError)
	runFlags.AddGoFlag(flag.CommandLine.Lookup("once"))
	runFlags.AddGoFlag(flag.CommandLine.Lookup("keepsqsmessages"))

	RootCmd = cmd.GenRootCmdWithRunFlags(Name, "", beater.New, runFlags)
}
