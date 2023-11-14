package model

import (
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type ConfigFile struct {
	System struct {
		OS_ID     string
		OS_Origin string
		Arch      string
		Packager  []string
	}
	Isolated struct {
		Snap    string
		Flatpak string
	}
	Language struct {
		Go     string
		Rustup string
		Cabal  string
		Python struct {
			Pip    string
			Pipenv string
			Poetry string
		}
		Perl5 struct {
			Cpan string
		}
		Rubygem string
	}
	Source string
}

type Config interface {
	GetConfigFile() string
	ReadConfigFile() error
	WriteConfigFile() error
	GetData() ConfigFile
	SetData(ConfigFile)
}

type config struct {
	*viper.Viper
	path   string
	name   string
	suffix string
	data   ConfigFile
}

var single_config *config

func (c config) GetConfigFile() string {
	return c.path + c.name + c.suffix
}

func (c config) GetData() ConfigFile {
	return c.data
}

func (c config) SetData(cf ConfigFile) {
	c.data = cf
}

func (c config) ReadConfigFile() error {
	return c.Unmarshal(c.data)
}

func (c config) WriteConfigFile() error {
	setting := c.AllSettings()
	bs, err := yaml.Marshal(setting)
	if err != nil {
		return err
	}
	options := os.O_RDWR | os.O_CREATE
	f, err := os.OpenFile(c.GetConfigFile(), options, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// write bs in config file
	for _, str := range bs {
		_, err = f.WriteString(string(str) + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func GetConfig() Config {
	if single_config == nil {
		home_dir, err := os.UserHomeDir()
		if err != nil {
			panic("Can not read home directory")
		}
		path := home_dir + "/.config/pkg-app-viewer"
		name := "config"
		suffix := "yaml"
		single_config = &config{
			viper.New(),
			path,
			name,
			suffix,
			ConfigFile{},
		}
		single_config.SetConfigName(name)
		single_config.SetConfigType(suffix)
		single_config.AddConfigPath(path)
		err = single_config.ReadInConfig()
		if err != nil {
			panic("can not read config file")
		}
		err = single_config.Unmarshal(single_config.data)
		if err != nil {
			panic("can not unmarshall config file")
		}
	}
	return single_config
}
