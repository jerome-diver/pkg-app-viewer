package main

import (
	"log/slog"

	controller "github.com/pkg-app-viewer/controllers"
	view "github.com/pkg-app-viewer/views"
	"github.com/spf13/viper"
)

const version = "0.0.1"

var logger *slog.Logger
var menu *controller.Menu
var packages *controller.Packages

func init() {
	logger = view.NewLogger()
	//viper.SetConfigFile("/home/dge/.config/pkg-app-viewer.config.yml")
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath("/home/dge/.config/pkg-app-viewer")
	err := config.ReadInConfig()
	if err != nil {
		logger.Error("can not unmarshall config file", slog.String("err", err.Error()))
	}
	menu = controller.NewMenu(logger, version, config)
	packages = controller.NewPackages(menu.Model, logger)
}

func main() {

	view.DebugLevel(menu.Model.Debug)
	packages.RunController()

}
