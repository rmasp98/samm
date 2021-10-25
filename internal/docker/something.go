package docker

import (
	"fmt"
	"os"

	"github.com/docker/cli/cli/command"
)

func Test() {
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	//dockerCli.Client().NetworkCreate
	fmt.Println(dockerCli.ClientInfo().DefaultVersion)
}
