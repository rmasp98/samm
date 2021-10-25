package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/cli/cli/command"
	composetypes "github.com/docker/cli/cli/compose/types"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// TODO:
// Start, Restart, Stop, Logs

type Service struct {
	Name   string
	ID     string
	Active bool
}

type ServiceClient struct {
	cli      *command.DockerCli
	services map[string]composetypes.ServiceConfig
}

func NewServiceClient(fileLocation string, cli *command.DockerCli) (ServiceClient, error) {
	if cli == nil {
		var clientErr error
		if cli, clientErr = NewClient(flags.NewClientOptions()); clientErr != nil {
			return ServiceClient{}, clientErr
		}
	}

	config, fileErr := load(fileLocation)
	if fileErr != nil {
		return ServiceClient{}, fileErr
	}

	services, marshalErr := config.Services.MarshalYAML()
	if marshalErr != nil {
		return ServiceClient{}, marshalErr
	}

	return ServiceClient{cli: cli, services: services.(map[string]composetypes.ServiceConfig)}, nil
}

func EnvMapToList(mapping map[string]*string) []string {
	var envList []string
	for key, value := range mapping {
		fmt.Printf("%s=%s", key, *value)
		envList = append(envList, key+*value)
	}
	return envList
}

func GetExposedPorts(ports []string) nat.PortSet {
	mapping := make(map[nat.Port]struct{})
	for _, port := range ports {
		mapping[nat.Port(port)] = struct{}{}
	}
	return mapping
}

func GetHealthCheck(hc *composetypes.HealthCheckConfig) *container.HealthConfig {
	if hc != nil {
		config := container.HealthConfig{}
		config.Test = hc.Test
		if hc.Interval != nil {
			config.Interval = time.Duration(*hc.Interval)
		}
		if hc.Timeout != nil {
			config.Timeout = time.Duration(*hc.Timeout)
		}
		if hc.StartPeriod != nil {
			config.StartPeriod = time.Duration(*hc.StartPeriod)
		}
		if hc.Retries != nil {
			config.Retries = int(*hc.Retries)
		}
		return &config
	}
	return nil
}

func GetVolumes(vols []composetypes.ServiceVolumeConfig) map[string]struct{} {
	volumes := make(map[string]struct{})
	for _, vol := range vols {
		if vol.Type == "Volume" { //TODO find out what this is
			volString := vol.Source + ":" + vol.Target + ":"
			if vol.ReadOnly {
				volString += "ro"
			} else {
				volString += "rw"
			}
			volumes[volString] = struct{}{}
		}
	}
	return volumes
}

func GetBinds(vols []composetypes.ServiceVolumeConfig) []string {
	var binds []string
	for _, vol := range vols {
		if vol.Type == "Bind" { //TODO find out what this is
			// create a bind string
		}
	}
	return binds
}

func (s ServiceClient) Create(name string) (string, error) {
	// TODO: Get details from compose file
	if c, exists := s.services[name]; exists {
		fmt.Println(c)
		fmt.Println("Create config")

		var stopTimeout int
		if c.StopGracePeriod != nil {
			stopTimeout = int(time.Duration(*c.StopGracePeriod).Seconds())
		}

		config := container.Config{
			Hostname:        c.Hostname,
			Domainname:      c.DomainName,
			User:            c.User,
			AttachStdin:     false,
			AttachStdout:    false,
			AttachStderr:    false,
			ExposedPorts:    GetExposedPorts(c.Expose),
			Tty:             c.Tty,
			OpenStdin:       c.StdinOpen,
			StdinOnce:       false,
			Env:             EnvMapToList(c.Environment),
			Cmd:             strslice.StrSlice(c.Command),
			Healthcheck:     GetHealthCheck(c.HealthCheck),
			ArgsEscaped:     false,
			Image:           c.Image,
			Volumes:         GetVolumes(c.Volumes),
			WorkingDir:      c.WorkingDir,
			Entrypoint:      strslice.StrSlice(c.Entrypoint),
			NetworkDisabled: false,
			MacAddress:      c.MacAddress,
			OnBuild:         []string{}, //TODO
			Labels:          c.Labels,
			StopSignal:      c.StopSignal,
			StopTimeout:     &stopTimeout,
			Shell:           strslice.StrSlice{}, //TODO
		}

		fmt.Println(config)
		fmt.Println("Create host")
		fmt.Println(c.Ipc)

		restartPolicy := container.RestartPolicy{}
		if c.Deploy.RestartPolicy != nil {
			restartPolicy.MaximumRetryCount = int(*c.Deploy.RestartPolicy.MaxAttempts)
		}
		host := container.HostConfig{
			Binds:           GetBinds(c.Volumes),
			ContainerIDFile: "", //TODO
			LogConfig:       container.LogConfig{Type: c.Logging.Driver, Config: c.Logging.Options},
			NetworkMode:     container.NetworkMode(c.NetworkMode),
			PortBindings:    nat.PortMap{}, //TODO
			RestartPolicy:   restartPolicy,
			AutoRemove:      false,
			VolumeDriver:    "",         //TODO
			VolumesFrom:     []string{}, //TODO
			CapAdd:          c.CapAdd,
			CapDrop:         c.CapDrop,
			CgroupnsMode:    container.CgroupnsMode(c.CgroupNSMode),
			DNS:             c.DNS,
			DNSOptions:      []string{}, //Not supported in compose v3
			DNSSearch:       c.DNSSearch,
			ExtraHosts:      c.ExtraHosts,
			GroupAdd:        []string{}, //TODO
			IpcMode:         "",         //TODO c.Ipc
			Cgroup:          "",         //TODO
			Links:           c.Links,
			OomScoreAdj:     0,  //TODO
			PidMode:         "", //TODO
			Privileged:      c.Privileged,
			PublishAllPorts: false, //TODO
			ReadonlyRootfs:  c.ReadOnly,
			SecurityOpt:     c.SecurityOpt,
			StorageOpt:      map[string]string{}, //TODO
			Tmpfs:           map[string]string{}, //TODO c.Tmpfs
			UTSMode:         "",                  //TODO
			UsernsMode:      container.UsernsMode(c.UserNSMode),
			Sysctls:         c.Sysctls,
			Runtime:         "",              //TODO
			Mounts:          []mount.Mount{}, //TODO
			Init:            c.Init,

			//Resources figure this out
		}

		if value, err := units.RAMInBytes(c.ShmSize); err == nil {
			host.ShmSize = value
		}

		fmt.Println(host)
		fmt.Println("Create network")
		endpoint := network.EndpointSettings{
			IPAMConfig:          &networktypes.EndpointIPAMConfig{}, //TODO
			Links:               []string{},                         //TODO
			Aliases:             []string{},                         //TODO
			NetworkID:           "",                                 //TODO
			EndpointID:          "",                                 //TODO
			Gateway:             "",                                 //TODO
			IPAddress:           "",                                 //TODO
			IPPrefixLen:         0,                                  //TODO
			IPv6Gateway:         "",                                 //TODO
			GlobalIPv6Address:   "",                                 //TODO
			GlobalIPv6PrefixLen: 0,                                  //TODO
			MacAddress:          "",                                 //TODO
			DriverOpts:          map[string]string{},                //TODO
		}
		network := networktypes.NetworkingConfig{
			map[string]*network.EndpointSettings{"": &endpoint},
		}

		fmt.Println(network)
		fmt.Println("Create platform")
		platform := v1.Platform{
			Architecture: "",
			OS:           "",
			OSVersion:    "",
			OSFeatures:   []string{},
			Variant:      "",
		}

		fmt.Println(platform)

		//	response, err := s.cli.Client().ContainerCreate(
		//		context.Background(), &config, &host, &network, &platform, name,
		//	)
		//
		//	var output string
		//	for _, warning := range response.Warnings {
		//		output += warning + "\n"
		//	}
		//	return strings.TrimRight(output, "\n"), err
	}
	return "", nil
}

func (s ServiceClient) Delete(name string, removeVolumes, force bool) error {
	// TODO: get ID from name in list containers and compose file
	id := ""
	options := types.ContainerRemoveOptions{
		RemoveVolumes: removeVolumes,
		Force:         force,
	}
	if err := s.cli.Client().ContainerRemove(context.Background(), id, options); err != nil {
		return err
	}
	return nil
}

func (s ServiceClient) List() ([]Service, error) {
	services, err := s.getRunningContainers()
	if err != nil {
		return []Service{}, err
	}

	for _, service := range s.services {
		running := false
		for _, container := range services {
			if service.Name == container.Name {
				running = true
			}
		}

		if !running {
			services = append(services, Service{service.Name, "", false})
		}
	}
	return services, nil
}

func (s ServiceClient) getRunningContainers() ([]Service, error) {
	options := types.ContainerListOptions{
		Quiet: true, Size: false, All: true, Latest: false,
		Limit: 0, Before: "", Filters: filters.Args{}, Since: "",
	}
	containers, err := s.cli.Client().ContainerList(context.Background(), options)
	if err != nil {
		return []Service{}, err
	}

	var services []Service
	for _, container := range containers {
		services = append(services, Service{container.Names[0], container.ID, true})
	}
	return services, nil
}
