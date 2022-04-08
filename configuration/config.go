package configuration

import (
	"flag"
	"fmt"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Configuration struct {
	Server Server `yaml:"server"`
	Docker Docker `yaml:"docker"`
}

type Server struct {
	Host              string `yaml:"host"`
	Port              uint16 `yaml:"port"`
	GracefulTimeoutMs uint32 `yaml:"gracefulTimeoutMs"`
}

type Docker struct {
	Protocol string `yaml:"protocol"`
	Registry string `yaml:"registry"`
	Port     uint16 `yaml:"port"`
	Address  string
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func Configure() (*Configuration, error) {
	// parse the flags passed along
	flags, err := parseFlags()
	if err != nil {
		return nil, err
	}

	// load default configuration
	config, err := loadConfigurationFile("./default_configuration.yaml")
	if err != nil {
		return nil, err
	}

	// load configuration from path
	if flags.configPath != "" {
		if loadedConfiguration, err := loadConfigurationFile(flags.configPath); err != nil {
			return nil, err
		} else {
			log.Println("merging loaded configuration with default")
			if err := mergo.Merge(config, loadedConfiguration); err != nil {
				return nil, err
			}
		}
	}

	if err := mergo.Merge(&config.Docker, flags.dockerConfig); err != nil {
		return nil, err
	}

	if config.Docker.Port != 0 {
		config.Docker.Address = fmt.Sprintf("%s%s:%d", config.Docker.Protocol, config.Docker.Registry, config.Docker.Port)
	} else {
		config.Docker.Address = fmt.Sprintf("%s%s", config.Docker.Protocol, config.Docker.Registry)
	}

	if err := validateConfiguration(*config); err != nil {
		return nil, err
	}

	return config, nil
}

func loadConfigurationFile(configPath string) (*Configuration, error) {
	config := &Configuration{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfigurationPath(configPath string) error {
	s, err := os.Stat(configPath)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("[%s] is a directory, not a file", configPath)
	}
	return nil
}

func validateConfiguration(config Configuration) error {
	// validate the provided server configuration
	if err := validateServerConfiguration(config.Server); err != nil {
		return err
	}

	// validate the provided docker configuration
	if err := validateDockerConfiguration(config.Docker); err != nil {
		return err
	}

	// no errors found, return nil
	return nil
}

func validateServerConfiguration(config Server) error {
	if config.Port == 0 {
		return fmt.Errorf("[%d] is an invalid port", config.Port)
	}

	return nil
}

func validateDockerConfiguration(config Docker) error {
	if config.Username == "" {
		return fmt.Errorf("docker username is required")
	}

	if config.Password == "" {
		return fmt.Errorf("docker password is required")
	}

	return nil
}

type Flags struct {
	configPath   string
	dockerConfig Docker
}

func parseFlags() (*Flags, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to configuration file")

	var registryName string
	flag.StringVar(&registryName, "registry-name", "", "name for the docker registry, without port/protocol")

	var registryProtocol string
	flag.StringVar(&registryProtocol, "registry-protocol", "https://", "protocol for the docker registry")

	var registryPort uint
	flag.UintVar(&registryPort, "registry-port", 0, "port for the docker registry, by default will use port of the specified protocol")

	var registryPassword string
	flag.StringVar(&registryPassword, "registry-password", "", "password for the docker registry")

	var registryUsername string
	flag.StringVar(&registryUsername, "registry-username", "", "username for the docker registry")

	flag.Parse()

	if configPath != "" {
		if err := validateConfigurationPath(configPath); err != nil {
			return nil, err
		}
	}

	flags := Flags{
		configPath: configPath,
		dockerConfig: Docker{
			Username: registryUsername,
			Password: registryPassword,
			Registry: registryName,
			Protocol: registryProtocol,
			Port:     uint16(registryPort),
		},
	}

	return &flags, nil
}
