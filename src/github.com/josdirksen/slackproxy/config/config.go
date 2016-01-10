package config
import (
	"os"
	"encoding/json"
	"fmt"
	"errors"
)

// for now fixed client which points to localhost, should move this to
// external configuration
var config *Configuration

// TODO: use a nested struct
type Environment struct {
	Name 	string
	Host    string
	Tls	bool
	Path 	string
	Ca 	string
	Cert 	string
	Key 	string
}

type Configuration struct {
	Environments    []Environment
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
	fmt.Println(configuration.Environments)

	config = &configuration
}

func GetConfig() *Configuration {
	return config
}

func GetEnvironmentConfig(Environment string) (*Environment, error) {

	for _, env := range config.Environments {
		if (env.Name == Environment) {
			return &env, nil
		}
	}

	return nil, errors.New("Environment not found")
}