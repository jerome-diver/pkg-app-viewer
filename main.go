package main

import (
	controller "github.com/pkg-app-viewer/controllers"
	report "github.com/pkg-app-viewer/views/report"
)

const version = "0.2.1"

var logger report.Logging
var menu *controller.Menu
var control *controller.Operator

func init() {
	menu = controller.NewMenu()
	control = controller.NewOperator(menu.Model)
	control.BuildConfig(menu)
	menu.InitView(version, control.PackagesId)
	control.BuildMenu(menu)
}

func main() {

	logger.DebugLevel(menu.Model.Debug)
	control.BuildContent()

}
