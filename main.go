package main

import (
	"log/slog"
	"os"

	controller "github.com/pkg-app-viewer/controllers"
	view "github.com/pkg-app-viewer/views"
	"github.com/spf13/viper"
)

const version = "0.0.1"

var logger *view.Logging
var menu *controller.Menu
var packages *controller.Packages

func init() {
	logger = view.NewLogger()
	config := viper.New()
	var home_dir string
	home_dir, logger.Error = os.UserHomeDir()
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath(home_dir + "/.config/pkg-app-viewer")
	err := config.ReadInConfig()
	if err != nil {
		logger.Log.Error("can not unmarshall config file", slog.String("err", err.Error()))
	}
	menu = controller.NewMenu(logger, version, config)
	packages = controller.NewPackages(menu, logger, config)
}

func main() {

	logger.DebugLevel(menu.Model.Debug)
	packages.RunController()

}
