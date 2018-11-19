package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	filepath "path"
)

const defaultPort = 8080

var (
	currUsr          *user.User
	defautConfigFile string
)

func init() {
	var err error
	currUsr, err = user.Current()
	if err != nil {
		log.Printf("Warning: Failed to get current user so default config file cannot be used. Please provide a path to a config file. Error: %v", err)
	} else {
		defautConfigFile = filepath.Join(currUsr.HomeDir, ".config", "irremote", "irremote.conf.json")
	}
}

type Config struct {
	Token string `json:"token"`
	Port  int    `json:"port"`
}

func (c Config) String() string {
	if s, err := json.Marshal(c); err == nil {
		return string(s)
	} else {
		return fmt.Sprintf("Failed to marshal Config: %v", err)
	}
}

func Default() Config {
	return Config{
		Port: defaultPort,
	}
}

func unmarshal(path string) (Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	return config, err
}

func Load(path string) (Config, error) {
	config := Default()
	var readErr error
	if path != "" {
		log.Printf("Using config %v", path)
		config, readErr = unmarshal(path)
	} else if _, err := os.Stat(defautConfigFile); !os.IsNotExist(err) {
		path = defautConfigFile

		log.Printf("Using config %v", path)
		config, readErr = unmarshal(path)
	} else {
		log.Printf("Using default config: %v", config)
	}

	if readErr != nil {
		return config, readErr
	}

	if config.Port == 0 {
		config.Port = defaultPort
	}

	return config, readErr
}