package controller

import (
	"fmt"
	"log/slog"

	model "github.com/pkg-app-viewer/models"
)

type Packages struct {
	MenuModel *model.Menu
	Tool      *Tool
	Find      *Find
	logging   *slog.Logger
	Installed *model.Installed
}

func NewPackages(m *model.Menu, logger *slog.Logger) *Packages {
	p := new(Packages)
	p.MenuModel = m
	p.logging = logger
	p.Tool = ToolBox(p.logging)
	p.Find = Finder(p.Tool)
	p.Installed = model.NewInstalled(logger)
	return p
}

func (p *Packages) RunController() {
	switch p.MenuModel.PackageType {
	case model.Apt:
		{
			found := p.Apt(p.MenuModel.PackageSearch)
			for _, r := range found {
				switch p.MenuModel.OutputMode {
				case "stdout":
					fmt.Println(r)
				default:
					fmt.Println(r)
				}
			}
		}
	}
}

func (p *Packages) Apt(searchFor model.Search) []string {
	p.logging.Debug("RUN Packages.Apt from controller (package_reader.go)",
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
