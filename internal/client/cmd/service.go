package cmd

import (
	"fmt"
	"os"

	"github.com/rmasp98/samm/internal/docker"
	"github.com/spf13/cobra"
)

var (
	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "Manage networks of the mailserver",
		Long:  "TODO",
	}

	serviceListCmd = &cobra.Command{
		Use:   "ls",
		Short: "List services",
		Long:  "TODO",
		Run:   ListServices,
	}

	serviceCreateCmd = &cobra.Command{
		Use:   "create [service]",
		Short: "Create service in compose config",
		Long:  "TODO",
		Run:   CreateServices,
	}

	serviceDeleteCmd = &cobra.Command{
		Use:   "delete [service]",
		Short: "Delete service in compose config",
		Long:  "TODO",
		Run:   DeleteServices,
	}

	serviceRemoveVolumes bool
	serviceForce         bool
)

func init() {
	serviceCmd.AddCommand(serviceListCmd)
	serviceCmd.AddCommand(serviceCreateCmd)
	serviceCmd.AddCommand(serviceDeleteCmd)
	rootCmd.AddCommand(serviceCmd)

	serviceDeleteCmd.PersistentFlags().BoolVarP(
		&serviceRemoveVolumes, "remove-volumes", "v", false, "Delete attached volumes",
	)

	serviceDeleteCmd.PersistentFlags().BoolVarP(
		&serviceForce, "force", "f", false, "Force deletion of service",
	)
}

func ListServices(cmd *cobra.Command, args []string) {
	var services []docker.Service
	if local {
		client, clientErr := docker.NewServiceClient(composeFile, nil)
		if clientErr != nil {
			fmt.Println(clientErr.Error())
			os.Exit(1)
		}

		var err error
		if services, err = client.List(); err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}

	fmt.Println("Name\tID\tActive")
	for _, service := range services {
		fmt.Printf("%s\t%s\t%t\n", service.Name, service.ID, service.Active)
	}
}

func CreateServices(cmd *cobra.Command, args []string) {
	service := ""
	if len(args) == 1 {
		service = args[0]
	}

	var warnings string
	if local {
		client, clientErr := docker.NewServiceClient(composeFile, nil)
		if clientErr != nil {
			fmt.Println(clientErr.Error())
			os.Exit(-1)
		}

		var err error
		if warnings, err = client.Create(service); err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}

	if service == "" {
		fmt.Println("Created all services")
	} else {
		fmt.Printf("Created %s\n", service)
	}
	fmt.Println(warnings)
}

func DeleteServices(cmd *cobra.Command, args []string) {
	service := ""
	if len(args) == 1 {
		service = args[0]
	}

	if local {
		client, clientErr := docker.NewServiceClient(composeFile, nil)
		if clientErr != nil {
			fmt.Println(clientErr.Error())
			os.Exit(-1)
		}

		if err := client.Delete(service, serviceRemoveVolumes, serviceForce); err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}

	if service == "" {
		fmt.Println("Deleted all services")
	} else {
		fmt.Printf("Deleted %s\n", service)
	}
}
