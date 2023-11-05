package model

import (
	"os"

	"github.com/spf13/viper"
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
	SetConfigFile(string)
	GetConfigFile() string
	GetData() ConfigFile
	SetData(ConfigFile)
	SaveConfig() error
}

type config struct {
	*viper.Viper
	File string
	Data ConfigFile
}

var single_config *config

func (c config) SetConfigFile(fullFileName string) {}

func (c config) GetConfigFile() string {
	return c.File
}

func (c config) GetData() ConfigFile {
	return c.Data
}

func (c config) SetData(cf ConfigFile) {
	c.Data = cf
}

func (c config) SaveConfig() error {
	return c.Unmarshal(c.Data)
}

func GetConfig() Config {
	if single_config == nil {
		var err error
		var home_dir string
		single_config = &config{
			viper.New(),
			"",
			ConfigFile{},
		}
		home_dir, _ = os.UserHomeDir()
		single_config.SetConfigName("config")
		single_config.SetConfigType("yaml")
		single_config.AddConfigPath(home_dir + "/.config/pkg-app-viewer")
		err = single_config.ReadInConfig()
		if err != nil {
			panic("can not unmarshall config file")
		}
	}
	return single_config
}
