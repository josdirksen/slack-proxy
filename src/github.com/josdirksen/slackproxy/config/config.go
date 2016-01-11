package config
import (
	"os"
	"encoding/json"
	"fmt"
	"errors"
)

var config *Configuration

type Environment struct {
	Name 	string
	Host    string
	Tls	bool
	Path 	string
	Ca 	string
	Cert 	string
	Key 	string
}

type Docker struct {
	Environments    []Environment
}

type Configuration struct {
	Docker Docker
}

func ParseConfig(configFile string) {
	// parse the configuration
	file, _ := os.Open(configFile)
	decoder := json.NewDecoder(file)

	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error parsing config: ", err)
		os.Exit(1);
	}
	fmt.Println(configuration)
	config = &configuration
}

func GetConfig() *Configuration {
	return config
}

func GetDockerEnvironmentConfig(Environment string) (*Environment, error) {

	for _, env := range config.Docker.Environments {
		if (env.Name == Environment) {
			return &env, nil
		}
	}

	return nil, errors.New("Environment not found")
}