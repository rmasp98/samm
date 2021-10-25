package docker

import (
	"context"
	"fmt"

	"github.com/docker/cli/cli/command"
	composetypes "github.com/docker/cli/cli/compose/types"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/docker/api/types"
	networktypes "github.com/docker/docker/api/types/network"
)

func ListNetworks(composeFile string) ([]Network, error) {
	client, clientErr := NewNetworkClient(composeFile, nil)
	if clientErr != nil {
		return []Network{}, clientErr
	}

	return client.List()
}

func CreateNetworks(network, composeFile string) (string, error) {
	client, clientErr := NewNetworkClient(composeFile, nil)
	if clientErr != nil {
		return "", clientErr
	}

	warning, err := client.Create(network)

	var output string
	if network == "" {
		output = "Deleted all networks"
	} else {
		output = "Deleted " + network + "\n"
	}

	if warning != "" {
		output += "\n" + warning
	}

	return output, err
}

func DeleteNetworks(network, composeFile string) (string, error) {
	client, clientErr := NewNetworkClient(composeFile, nil)
	if clientErr != nil {
		return "", clientErr
	}

	err := client.Delete(network)
	if network == "" {
		return "Deleted all networks", err
	}
	return "Deleted " + network + "\n", err
}

type Network struct {
	Name   string
	ID     string
	Driver string
	Scope  string
	Active bool
}

type NetworkClient struct {
	cli      *command.DockerCli
	networks map[string]composetypes.NetworkConfig
}

func NewNetworkClient(fileLocation string, cli *command.DockerCli) (NetworkClient, error) {
	if cli == nil {
		var clientErr error
		if cli, clientErr = NewClient(flags.NewClientOptions()); clientErr != nil {
			return NetworkClient{}, clientErr
		}
	}

	config, fileErr := load(fileLocation)
	if fileErr != nil {
		return NetworkClient{}, fileErr
	}

	return NetworkClient{cli: cli, networks: config.Networks}, nil
}

// CreateNetwork creates the network called netName, if an empty string is defined all networks are created
func (n *NetworkClient) Create(netName string) (string, error) {
	created := false
	output := ""
	for name, net := range n.networks {
		if netName == "" || netName == name {
			response, err := n.cli.Client().NetworkCreate(context.Background(), name,
				createNetworkOptions(net))
			if err != nil {
				return "", err
			}
			if response.Warning != "" {
				output += response.Warning + "\n"
			}
			created = true
		}
	}
	if created {
		return output, nil
	} else {
		return "", fmt.Errorf("Network %s does not exist in configuration", netName)
	}
}

func (s *NetworkClient) Delete(netName string) error {
	deleted := false
	for name := range s.networks {
		if netName == "" || netName == name {
			if err := s.cli.Client().NetworkRemove(context.Background(), name); err != nil {
				return err
			}
			deleted = true
		}
	}

	if deleted == false {
		return fmt.Errorf("Network %s does not exist in configuration", netName)
	}
	return nil
}

func (s *NetworkClient) List() ([]Network, error) {
	activeNetworks, err := s.cli.Client().NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return []Network{}, err
	}

	var networks []Network
	for name := range s.networks {
		active := false
		for _, net := range activeNetworks {
			if name == net.Name {
				networks = append(networks, Network{net.Name, net.ID, net.Driver, net.Scope, true})
				active = true
			}
		}

		if active == false {
			networks = append(networks, Network{name, "-", "-", "-", false})
		}
	}

	return networks, nil
}

func createNetworkOptions(netConfig composetypes.NetworkConfig) types.NetworkCreate {
	createOpts := types.NetworkCreate{
		CheckDuplicate: true,
		Labels:         netConfig.Labels,
		Driver:         netConfig.Driver,
		Options:        netConfig.DriverOpts,
		Internal:       netConfig.Internal,
		Attachable:     netConfig.Attachable,
	}

	if netConfig.Ipam.Driver != "" || len(netConfig.Ipam.Config) > 0 {
		createOpts.IPAM = &networktypes.IPAM{}
	}

	if netConfig.Ipam.Driver != "" {
		createOpts.IPAM.Driver = netConfig.Ipam.Driver
	}
	for _, ipamConfig := range netConfig.Ipam.Config {
		config := networktypes.IPAMConfig{
			Subnet: ipamConfig.Subnet,
		}
		createOpts.IPAM.Config = append(createOpts.IPAM.Config, config)
	}

	return createOpts
}
