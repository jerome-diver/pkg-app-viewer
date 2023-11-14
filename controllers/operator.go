package controller

import (
	identifier "github.com/pkg-app-viewer/controllers/manager"
	tool "github.com/pkg-app-viewer/controllers/tools"
	model "github.com/pkg-app-viewer/models"
	view "github.com/pkg-app-viewer/views"
	report "github.com/pkg-app-viewer/views/report"
)

type Operator struct {
	MenuModel  *model.Menu
	Printer    *view.Printer
	Tool       *tool.Box
	PackagesId *model.Identity
	Installed  []*identifier.Manager
}

var config model.Config
var logging report.Logging

func NewOperator(menu_model *model.Menu) *Operator {
	logging = report.GetLogger()
	o := new(Operator)
	o.MenuModel = menu_model
	o.Tool = tool.New()
	o.PackagesId = identifier.New()
	o.Installed = []*identifier.Manager{}
	o.Printer = view.NewPrinter(menu_model)
	return o
}

/*
BuildConfig will check config file
and popoulate it after to ask user
how to do if something new or different,
then will reflect accordingly
to Operator.Config *viper.Viper content
and Operator.ConfigFile struct
*/
func (o *Operator) BuildConfig(menu *Menu) {
	// get content of config file
	config = model.GetConfig()
	logging.SetError(config.ReadConfigFile())
	logging.CheckError("Can not build config data")
	// compare o.PackageId *identifier with config data
	// complete config file with relevent identifier
	tool.UpdateConfig(&config, o.PackagesId)
	menu.Config = config.GetData()
}

/*
BuildMenu is a Factory that will recognize:
_ System package manager
_ Isolated (Flatpack or Snap) present
_ Language managers
to build menu possible entries to use.
*/
func (o *Operator) BuildMenu(menu *Menu) {
	logging.SetError(config.WriteConfigFile())
	logging.CheckError("Can not save config data")
	// Define System menu package manager present
	if o.PackagesId.System.GetTypes() == model.Dpkg {

	}
	menu.InitMenuCommand()
}

/*
BuildContent will use an abstract Factory to
build the package list content to put
inside Operator.Installed slices
*/
func (o *Operator) BuildContent() {
}

func (p *Operator) Show() {
	switch p.MenuModel.PackageType {
	case model.Dpkg:
		{
			found := p.Apt(p.MenuModel.PackageOption)
			p.Printer.Write(found)
		}
	}
	switch p.MenuModel.Command {
	case model.System_managers:
		{
		}
	}
}
