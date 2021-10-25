package cmd

import (
	"fmt"
	"os"

	"github.com/rmasp98/samm/internal/docker"
	"github.com/spf13/cobra"
)

var (
	volCmd = &cobra.Command{
		Use:   "volume",
		Short: "Manage volumes of the mailserver",
		Long:  "TODO",
	}

	volListCmd = &cobra.Command{
		Use:   "ls",
		Short: "List volumes",
		Long:  "TODO",
		Run:   ListVolumes,
	}

	volCreateCmd = &cobra.Command{
		Use:   "create [volume]",
		Short: "Create volumes in compose config",
		Long:  "TODO",
		Run:   CreateVolumes,
	}

	volDeleteCmd = &cobra.Command{
		Use:   "delete [volume]",
		Short: "Delete volumes in compose config",
		Long:  "TODO",
		Run:   DeleteVolumes,
	}

	volForce bool
)

func init() {
	volCmd.AddCommand(volListCmd)
	volCmd.AddCommand(volCreateCmd)
	volCmd.AddCommand(volDeleteCmd)
	rootCmd.AddCommand(volCmd)

	volDeleteCmd.PersistentFlags().BoolVarP(&volForce, "force", "f", false, "Force the deletion of volume")
}

func ListVolumes(cmd *cobra.Command, args []string) {
	var volumes []docker.Volume
	if local {
		var err error
		if volumes, err = docker.ListVolumes(composeFile); err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}

	fmt.Println("Name\tDriver\tScope\tActive")
	for _, volume := range volumes {
		fmt.Printf("%s\t%s\t%s\t%t\n", volume.Name, volume.Driver, volume.Scope, volume.Active)
	}
}

func CreateVolumes(cmd *cobra.Command, args []string) {
	volume := ""
	if len(args) == 1 {
		volume = args[0]
	}

	var output string
	if local {
		var err error
		if output, err = docker.CreateVolumes(volume, composeFile); err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}

	fmt.Println(output)
}

func DeleteVolumes(cmd *cobra.Command, args []string) {
	volume := ""
	if len(args) == 1 {
		volume = args[0]
	}

	var output string
	if local {
		var err error
		if output, err = docker.DeleteVolumes(volume, composeFile, volForce); err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
	}

	fmt.Println(output)
}
