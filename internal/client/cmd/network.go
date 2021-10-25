package cmd

import (
	"fmt"
	"os"

	"github.com/rmasp98/samm/internal/docker"
	"github.com/spf13/cobra"
)

var (
	netCmd = &cobra.Command{
		Use:   "network",
		Short: "Manage networks of the mailserver",
		Long:  "TODO",
	}

	netListCmd = &cobra.Command{
		Use:   "ls",
		Short: "List networks",
		Long:  "TODO",
		Run:   ListNetworks,
	}

	netCreateCmd = &cobra.Command{
		Use:   "create [network]",
		Short: "Create networks in compose config",
		Long:  "TODO",
		Run:   CreateNetworks,
	}

	netDeleteCmd = &cobra.Command{
		Use:   "delete [network]",
		Short: "Delete networks in compose config",
		Long:  "TODO",
		Run:   DeleteNetworks,
	}
)

func init() {
	netCmd.AddCommand(netListCmd)
	netCmd.AddCommand(netCreateCmd)
	netCmd.AddCommand(netDeleteCmd)
	rootCmd.AddCommand(netCmd)
}

func ListNetworks(cmd *cobra.Command, args []string) {
	var networks []docker.Network
	if local {
		var err error
		if networks, err = docker.ListNetworks(composeFile); err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}

	fmt.Println("Name\tDriver\tScope\tActive")
	for _, network := range networks {
		fmt.Printf("%s\t%s\t%s\t%t\n", network.Name, network.Driver, network.Scope, network.Active)
	}
}

func CreateNetworks(cmd *cobra.Command, args []string) {
	network := ""
	if len(args) == 1 {
		network = args[0]
	}

	var output string
	if local {
		var err error
		if output, err = docker.CreateNetworks(network, composeFile); err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}

	fmt.Println(output)
}

func DeleteNetworks(cmd *cobra.Command, args []string) {
	network := ""
	if len(args) == 1 {
		network = args[0]
	}

	var output string
	if local {
		var err error
		if output, err = docker.DeleteNetworks(network, composeFile); err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}

	fmt.Println(output)
}
