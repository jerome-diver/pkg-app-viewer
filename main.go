package main

import (
	"log/slog"

	controller "github.com/pkg-app-viewer/controllers"
	view "github.com/pkg-app-viewer/views"
)

const version = "0.0.1"

var logger *slog.Logger
var menu *controller.Menu
var packages *controller.Packages

func init() {
	logger = view.NewLogger()
	menu = controller.NewMenu(logger, version)
	packages = controller.NewPackages(menu.Model, logger)
}

func main() {

	view.DebugLevel(menu.Model.Debug)
	packages.RunController()

}
