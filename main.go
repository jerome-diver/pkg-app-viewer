package main

import (
	controller "github.com/pkg-app-viewer/controllers"
	view "github.com/pkg-app-viewer/views"
)

const version = "0.0.2"

var logger view.Logging
var menu *controller.Menu
var control *controller.Operator

func init() {
	menu = controller.NewMenu(version)
	control = controller.NewOperator(menu)
	control.BuildConfig()
	control.BuildMenu(menu)
}

func main() {

	logger.DebugLevel(menu.Model.Debug)
	control.BuildContent()

}
