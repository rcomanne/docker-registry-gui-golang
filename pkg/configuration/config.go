package configuration

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

//go:embed default_configuration.yaml
var defaultConfig []byte

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

func (d Docker) HasAuthentication() bool {
	return d.Username != "" || d.Password != ""
}

func Configure() (*Configuration, error) {
	// parse the flags passed along
	flags, err := parseFlags()
	if err != nil {
		return nil, err
	}

	// load default configuration
	config, err := loadDefaultConfiguration()
	if err != nil {
		return nil, err
	}

	// load configuration from path if provided
	if flags.configPath != "" {
		// validate that file exists
		if err := validateConfigurationPath(flags.configPath); err != nil {
			return nil, err
		}
		// load file
		if loadedConfiguration, err := loadConfigurationFile(flags.configPath); err != nil {
			return nil, err
		} else {
			// and merge in into the default config
			log.Println("merging loaded configuration with default")
			if err := mergo.Merge(config, loadedConfiguration); err != nil {
				return nil, err
			}
		}
	}

	// merge docker config from flags into config
	if err := mergo.Merge(&config.Docker, flags.dockerConfig); err != nil {
		return nil, err
	}

	// create and set the docker address
	createDockerAddress(&config.Docker)

	// validate the configuration
	if valid := validateConfiguration(*config); !valid {
		return nil, err
	}

	return config, nil
}

func loadDefaultConfiguration() (*Configuration, error) {
	config := &Configuration{}
	if err := yaml.Unmarshal(defaultConfig, config); err != nil {
		return nil, err
	} else {
		return config, nil
	}
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

func createDockerAddress(config *Docker) {
	if config.Port != 0 {
		config.Address = fmt.Sprintf("%s%s:%d", config.Protocol, config.Registry, config.Port)
	} else {
		config.Address = fmt.Sprintf("%s%s", config.Protocol, config.Registry)
	}
}

func validateConfiguration(config Configuration) bool {
	valid := true

	// validate the provided server configuration
	if v := validateServerConfiguration(config.Server); !v {
		valid = false
	}

	// validate the provided docker configuration
	if v := validateDockerConfiguration(config.Docker); !v {
		valid = false
	}

	// no errors found, return true
	return valid
}

func validateServerConfiguration(config Server) bool {
	valid := true

	if config.Port == 0 {
		valid = false
		log.Printf("[%d] is an invalid port", config.Port)
	}

	return valid
}

func validateDockerConfiguration(config Docker) bool {
	valid := true

	// not every registry requires authentication, just log it
	if config.Username == "" {
		log.Println("no docker username set")
	}

	// not every registry requires authentication, just log it
	if config.Password == "" {
		log.Println("no docker password set")
	}

	if config.Registry == "" {
		valid = false
		log.Println("no docker registry provided to connect to")
	}

	if config.Protocol == "" {
		valid = false
		log.Println("empty docker registry protocol provided")
	}

	return valid
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
