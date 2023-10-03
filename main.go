package main

import (
	"fmt"
	"log/slog"

	controller "github.com/pkg-app-viewer/controllers"
	model "github.com/pkg-app-viewer/models"
	view "github.com/pkg-app-viewer/views"
)

var logger *slog.Logger
var menu *controller.Menu
var packages *controller.Packages

func init() {
	logger = view.NewLogger()
	menu = controller.NewMenu(logger)
	packages = controller.InitPackages(menu.Model, logger)
}

func main() {

	view.DebugLevel(menu.Model.Debug)

	switch menu.Model.PackageType {
	case model.AptAll:
		{
			packages.AptAll()
			for _, r := range packages.Installed.AptAll {
				fmt.Println(r)
			}
		}
	case model.AptAdded:
		{
			packages.AptAdded()
			for _, r := range packages.Installed.AptAdded {
				fmt.Println(r)
			}
		}
	case model.AptManual:
		{
			packages.AptManual()
			for _, r := range packages.Installed.AptManual {
				fmt.Println(r)
			}
		}
	case model.AptAdvanced:
		{
			switch menu.Model.Mode {
			case "File":
				{
					packages.AptFromFile()
				}
			case "Directory":
				{
					packages.AptFromDir()
				}
			}
			for _, r := range packages.Installed.AptAdvanced {
				fmt.Println(r)
			}
		}
	}

}
