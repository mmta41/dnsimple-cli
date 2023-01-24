package config

import (
	"errors"
	"fmt"
	"github.com/google/renameio"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"syscall"
)

func getDir() string {
	envDir := os.Getenv("DNS_CONFIG_DIR")
	if envDir != "" {
		return envDir
	}
	usrConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if usrConfigHome == "" {
		usrConfigHome = os.Getenv("HOME")
		if usrConfigHome == "" {
			usrConfigHome, _ = homedir.Expand("~/.config")
		} else {
			usrConfigHome = filepath.Join(usrConfigHome, ".config")
		}
	}
	return filepath.Join(usrConfigHome, "dns-cli")
}

func getFile() string {
	return path.Join(getDir(), "config.yml")
}

func readFromFile() (*Config, error) {
	return parseConfigFile(getFile())
}

func parseConfigFile(filename string) (*Config, error) {
	data, err := readConfigFile(filename)
	if err != nil {
		return nil, err
	}

	cfg, err := parseConfigData(data)
	if err != nil {
		return nil, err
	}
	return cfg, err
}

func readConfigFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, pathError(err)
	}

	return data, nil
}

func parseConfigData(data []byte) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal(data, cfg)
	return cfg, err
}

func pathError(err error) error {
	var pathError *os.PathError
	if errors.As(err, &pathError) && errors.Is(pathError.Err, syscall.ENOTDIR) {
		if p := findRegularFile(pathError.Path); p != "" {
			return fmt.Errorf("remove or rename regular file `%s` (must be a directory)", p)
		}

	}
	return err
}

func findRegularFile(p string) string {
	for {
		if s, err := os.Stat(p); err == nil && s.Mode().IsRegular() {
			return p
		}
		newPath := path.Dir(p)
		if newPath == p || newPath == "/" || newPath == "." {
			break
		}
		p = newPath
	}
	return ""
}

var WriteConfigFile = func(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	filename := getFile()
	err = os.MkdirAll(path.Dir(filename), 0750)
	if err != nil {
		return pathError(err)
	}
	_, err = ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = writeFile(filename, data, 0600)
	return err
}

func writeFile(filename string, data []byte, perm os.FileMode) error {
	pathToSymlink, err := filepath.EvalSymlinks(filename)
	if err == nil {
		filename = pathToSymlink
	}

	return renameio.WriteFile(filename, data, perm)
}
