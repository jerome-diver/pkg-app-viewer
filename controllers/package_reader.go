package controller

import (
	"log/slog"

	model "github.com/pkg-app-viewer/models"
	view "github.com/pkg-app-viewer/views"
	"github.com/spf13/viper"
)

type Packages struct {
	MenuModel  *model.Menu
	Printer    *view.Printer
	Tool       *Tool
	Find       *Find
	Config     *viper.Viper
	ConfigFile *model.ConfigFile
	logging    *view.Logging
	Installed  *model.Installed
}

func NewPackages(m *Menu, logger *view.Logging, config *viper.Viper) *Packages {
	p := new(Packages)
	p.MenuModel = m.Model
	p.logging = logger
	p.Config = config
	p.ConfigFile = m.Config
	p.Tool = ToolBox(p.logging)
	p.Find = Finder(logger)
	p.Installed = new(model.Installed)
	p.Printer = view.NewPrinter(logger, m.Model)
	return p
}

func (p *Packages) RunController() {
	switch p.MenuModel.PackageType {
	case model.Apt:
		{
			found := p.Apt(p.MenuModel.PackageSearch)
			p.Printer.Write(found)
		}
	}
	switch p.MenuModel.Command {
	case model.SearchSystem_pm:
		{
			p.SearchSystem_pm()
		}
	}
}

func (p *Packages) SearchSystem_pm() {
	spm := NewSearchManager(p.logging, p.ConfigFile)
	spm.SearchSystemPM()
}

func (p *Packages) Apt(searchFor model.Search) []string {
	p.logging.Log.Debug("RUN Packages.Apt from controller (package_reader.go)",
		slog.String("searchFor", searchFor.String()))
	if p.MenuModel.FileName != "" {
		clearBytes := p.Tool.GetFileContent(p.MenuModel)
		p.Find.AptInstalledFromHistory(clearBytes, searchFor)
	} else {
		if p.MenuModel.DirName == "" {
			p.MenuModel.DirName = "/var/log/apt"
		}
		filesList := p.Tool.GetAptHistoryFilesList(p.MenuModel)
		for _, file := range filesList {
			p.MenuModel.FileName = file
			clearBytes := p.Tool.GetFileContent(p.MenuModel)
			p.Find.AptInstalledFromHistory(clearBytes, searchFor)
		}
	}
	switch searchFor {
	case model.All:
		p.Installed.Apt.All = p.Find.Packages
	case model.Added:
		p.Installed.Apt.Added = p.Find.Packages
	case model.OtherRepos:
		p.Installed.Apt.Other = p.Find.Packages
	case model.OfficialRepos:
		p.Installed.Apt.Official = p.Find.Packages
	case model.FileSource:
		p.Installed.Apt.FileSource = p.Find.Packages
	}
	return p.Find.Packages
}
