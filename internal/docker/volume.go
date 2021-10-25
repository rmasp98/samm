package docker

import (
	"context"
	"fmt"

	"github.com/docker/cli/cli/command"
	composetypes "github.com/docker/cli/cli/compose/types"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

func ListVolumes(composeFile string) ([]Volume, error) {
	client, clientErr := NewVolumeClient(composeFile, nil)
	if clientErr != nil {
		return []Volume{}, clientErr
	}

	return client.List()
}

func CreateVolumes(volume, composeFile string) (string, error) {
	client, clientErr := NewVolumeClient(composeFile, nil)
	if clientErr != nil {
		return "", clientErr
	}

	err := client.Create(volume)

	if volume == "" {
		return "Created all volumes", err
	} else {
		return "Created " + volume + " volume", err
	}
}

func DeleteVolumes(volume, composeFile string, force bool) (string, error) {
	client, clientErr := NewVolumeClient(composeFile, nil)
	if clientErr != nil {
		return "", clientErr
	}

	err := client.Delete(volume, force)

	if volume == "" {
		return "Deleted all volumes", err
	} else {
		return "Deleted " + volume + " volume", err
	}
}

type Volume struct {
	Name   string
	Driver string
	Scope  string
	Active bool
}

type VolumeClient struct {
	cli     *command.DockerCli
	volumes map[string]composetypes.VolumeConfig
}

func NewVolumeClient(fileLocation string, cli *command.DockerCli) (VolumeClient, error) {
	if cli == nil {
		var clientErr error
		if cli, clientErr = NewClient(flags.NewClientOptions()); clientErr != nil {
			return VolumeClient{}, clientErr
		}
	}

	config, fileErr := load(fileLocation)
	if fileErr != nil {
		return VolumeClient{}, fileErr
	}

	return VolumeClient{cli: cli, volumes: config.Volumes}, nil
}

func (v *VolumeClient) Create(volName string) error {
	created := false
	for name, vol := range v.volumes {
		if volName == "" || volName == name {
			volCreate := volume.VolumeCreateBody{
				Driver:     vol.Driver,
				DriverOpts: vol.DriverOpts,
				Labels:     vol.Labels,
				Name:       vol.Name,
			}
			if _, err := v.cli.Client().VolumeCreate(context.Background(), volCreate); err != nil {
				return err
			}
			created = true
		}
	}
	if created {
		return nil
	} else {
		return fmt.Errorf("Volume %s does not exist in configuration", volName)
	}
}

func (v *VolumeClient) Delete(volName string, force bool) error {
	deleted := false
	for name := range v.volumes {
		if volName == "" || volName == name {
			if err := v.cli.Client().VolumeRemove(context.Background(), name, force); err != nil {
				return err
			}
			deleted = true
		}
	}

	if deleted == false {
		return fmt.Errorf("Volume %s does not exist in configuration", volName)
	}
	return nil
}

func (v *VolumeClient) List() ([]Volume, error) {
	activeVolumes, err := v.cli.Client().VolumeList(context.Background(), filters.Args{})
	if err != nil {
		return []Volume{}, err
	}

	var volumes []Volume
	for name := range v.volumes {
		active := false
		for _, vol := range activeVolumes.Volumes {
			if name == vol.Name {
				volumes = append(volumes, Volume{vol.Name, vol.Driver, vol.Scope, true})
				active = true
			}
		}

		if active == false {
			volumes = append(volumes, Volume{name, "-", "-", false})
		}
	}

	return volumes, nil
}
