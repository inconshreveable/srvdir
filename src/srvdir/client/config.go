package client

// XXX: factor this out into a library, this is mostly copied wholesale from ngrok

import (
	"fmt"
	"github.com/inconshreveable/go-tunnel/log"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"os/user"
	"path/filepath"
)

type Configuration struct {
	AuthToken string `yaml:"authtoken,omitempty"`
}

func SaveAuthToken(configPath, authtoken string) (err error) {
	if configPath == "" {
		configPath = defaultPath()
	}

	if authtoken == "" {
		return nil
	}

	// empty configuration by default for the case that we can't read it
	c := new(Configuration)

	// read the configuration
	oldConfigBytes, err := ioutil.ReadFile(configPath)
	if err == nil {
		// unmarshal if we successfully read the configuration file
		if err = goyaml.Unmarshal(oldConfigBytes, c); err != nil {
			return
		}
	}

	// no need to save, the authtoken is already the correct value
	if c.AuthToken == authtoken {
		return
	}

	// update auth token
	c.AuthToken = authtoken

	// rewrite configuration
	newConfigBytes, err := goyaml.Marshal(c)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(configPath, newConfigBytes, 0600)
	return
}

func LoadConfiguration(path string) (config *Configuration, err error) {
	configPath := path
	if configPath == "" {
		configPath = defaultPath()
	}

	log.Info("Reading configuration file %s", configPath)
	configBuf, err := ioutil.ReadFile(configPath)
	if err != nil {
		// failure to read a configuration file is only a fatal error if
		// the user specified one explicitly
		if path != "" {
			err = fmt.Errorf("Failed to read configuration file %s: %v", configPath, err)
			return
		}
	}

	// deserialize/parse the config
	config = new(Configuration)
	if err = goyaml.Unmarshal(configBuf, &config); err != nil {
		err = fmt.Errorf("Error parsing configuration file %s: %v", configPath, err)
		return
	}

	return
}

func defaultPath() string {
	user, err := user.Current()

	// user.Current() does not work on linux when cross compilling because
	// it requires CGO; use os.Getenv("HOME") hack until we compile natively
	homeDir := os.Getenv("HOME")
	if err != nil {
		log.Warn("Failed to get user's home directory: %s. Using $HOME: %s", err.Error(), homeDir)
	} else {
		homeDir = user.HomeDir
	}

	return filepath.Join(homeDir, ".srvdir")
}
