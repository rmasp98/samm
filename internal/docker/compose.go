package docker

import (
	"os"
	"path/filepath"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/compose/loader"
	composetypes "github.com/docker/cli/cli/compose/types"
	"github.com/docker/cli/cli/flags"
)

// TODO
// Load docker-compose.yml file
// Check if file latest version
//

func NewClient(options *flags.ClientOptions) (*command.DockerCli, error) {
	cli, cliErr := command.NewDockerCli()
	if cliErr != nil {
		return nil, cliErr
	}

	if initErr := cli.Initialize(options); initErr != nil {
		return nil, initErr
	}
	return cli, nil
}

type StackConfig struct {
	dockerCli *command.DockerCli

	services []composetypes.ServiceConfig
	networks map[string]composetypes.NetworkConfig
	volumes  map[string]composetypes.VolumeConfig
}

func NewStackConfig(fileLocation string, dockerCli *command.DockerCli) (StackConfig, error) {
	if dockerCli == nil {
		var clientErr error
		if dockerCli, clientErr = NewClient(flags.NewClientOptions()); clientErr != nil {
			return StackConfig{}, clientErr
		}
	}

	config, fileErr := load(fileLocation)
	if fileErr != nil {
		return StackConfig{}, fileErr
	}

	return StackConfig{dockerCli: dockerCli, services: config.Services,
		networks: config.Networks, volumes: config.Volumes}, nil
}

func (s *StackConfig) CreateAll() {}

func (s *StackConfig) StartService(name string)             {}
func (s *StackConfig) StopService(name string, remove bool) {}
func (s *StackConfig) RestartService(name string)           {}

func (s *StackConfig) CreateVolume(name string) {}
func (s *StackConfig) DeleteVolume(name string) {}

func load(fileLocation string) (*composetypes.Config, error) {
	fileContents, err := os.ReadFile(fileLocation)
	if err != nil {
		return nil, err
	}

	config, parseErr := loader.ParseYAML(fileContents)
	if parseErr != nil {
		return nil, parseErr
	}

	configFile := composetypes.ConfigFile{Filename: filepath.Base(fileLocation), Config: config}
	configDetails := composetypes.ConfigDetails{
		Version:     "3.8",
		WorkingDir:  filepath.Dir(fileLocation),
		ConfigFiles: []composetypes.ConfigFile{configFile},
		Environment: map[string]string{},
	}

	return loader.Load(configDetails)
}
