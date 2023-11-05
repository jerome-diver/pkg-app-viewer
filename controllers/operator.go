package controller

import (
	"log/slog"

	identifier "github.com/pkg-app-viewer/controllers/manager"
	tool "github.com/pkg-app-viewer/controllers/tools"
	model "github.com/pkg-app-viewer/models"
	view "github.com/pkg-app-viewer/views"
)

type Operator struct {
	MenuModel  *model.Menu
	Printer    *view.Printer
	Tool       *tool.Box
	PackagesId *identifier.Identifier
	Installed  []*identifier.Manager
}

var config model.Config
var logging view.Logging

func NewOperator(m *Menu) *Operator {
	config = model.GetConfig()
	logging = view.GetLogger()
	o := new(Operator)
	o.MenuModel = m.Model
	o.Tool = tool.New()
	o.PackagesId = identifier.New()
	o.Installed = []*identifier.Manager{}
	o.Printer = view.NewPrinter(m.Model)
	return o
}

/*
BuildMenu is a Factory that will recognize:
_ System package manager
_ Isolated (Flatpack or Snap) present
_ Language managers
to build menu possible entries to use.
*/
func (o *Operator) BuildMenu(menu *Menu) {
	config.SetData(menu.Config)
	logging.Err = config.SaveConfig()
	// Define System menu package manager present
	if o.PackagesId.Infos.System.GetTypes() == model.Dpkg {

	}
	menu.InitMenuCommand()
}

/*
BuildContent will use an abstract Factory to
build the pacdkages list content to put
inside Operator.Installed slices
*/
func (o *Operator) BuildContent() {
}

/*
BuildConfig will check config file
and popoulate it after to ask user
how to do if something new or different,
then will reflect accordingly
to Operator.Config *viper.Viper content
and Operator.ConfigFile struct
*/
func (o *Operator) BuildConfig() {
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

func (p *Operator) Apt(searchFor model.SystemOption) []string {
	p.logging.Log.Debug("RUN Operator.Apt from controller (operator.go)",
		slog.String("searchFor", searchFor.String()))
	if p.MenuModel.FileName != "" {
		clearBytes := p.Tool.GetFileContent(p.MenuModel)
		//		p.Find.DebianPackagesToSearchFor(clearBytes, searchFor)
	} else {
		if p.MenuModel.DirName == "" {
			p.MenuModel.DirName = "/var/log/apt"
		}
		filesList := p.Tool.GetAptHistoryFilesList(p.MenuModel)
		for _, file := range filesList {
			p.MenuModel.FileName = file
			clearBytes := p.Tool.GetFileContent(p.MenuModel)
			//			p.Find.DebianPackagesToSearchFor(clearBytes, searchFor)
		}
	}
	return nil
}
