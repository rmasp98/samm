package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "sammcli",
		Short: "CLI tool for managing mailserver",
		Long:  "Write this",
	}

	composeFile string
	configFile  string
	local       bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&composeFile, "compose-file", "k",
		"/home/rmaspero/programming/mailserver/docker-compose.yml",
		"Location of compose file",
	)

	rootCmd.PersistentFlags().StringVarP(
		&configFile, "config", "c",
		"/home/rmaspero/programming/mailserver/samm.conf",
		"Loction of the sammcli configuation file",
	)

	// TODO: change this to false when local stuff finished
	rootCmd.PersistentFlags().BoolVarP(&local, "local", "l", true, "Run commands on local docker daemon")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
