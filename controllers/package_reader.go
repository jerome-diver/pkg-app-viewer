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
	case model.AptAll:
		{
			p.AptAll()
			for _, r := range p.Installed.AptAll {
				fmt.Println(r)
			}
		}
	case model.AptAdded:
		{
			p.AptAdded()
			for _, r := range p.Installed.AptAdded {
				fmt.Println(r)
			}
		}
	case model.AptManual:
		{
			p.AptManual()
			for _, r := range p.Installed.AptManual {
				fmt.Println(r)
			}
		}
	case model.AptAdvanced:
		{
			switch p.MenuModel.Mode {
			case "File":
				{
					p.AptFromFile()
				}
			case "Directory":
				{
					p.AptFromDir()
				}
			}
			for _, r := range p.Installed.AptAdvanced {
				fmt.Println(r)
			}
		}
	}
}

func (p *Packages) AptAll() {
	p.logging.Debug("RUN Packages.AptAll from controller (package_reader.go)")
	p.MenuModel.DirName = "/var/log/apt"
	filesList := p.Tool.GetAptHistoryFilesList(p.MenuModel)
	for _, file := range filesList {
		p.MenuModel.FileName = file
		clearBytes := p.Tool.GetFileContent(p.MenuModel)
		p.Find.AptInstalledFromHistory(clearBytes, "all")
	}
	p.Installed.AptAll = p.Find.Packages
}

func (p *Packages) AptFromFile() {
	p.logging.Debug("RUN Packages.AptFromFile from controller (package_reader.go)")
	clearBytes := p.Tool.GetFileContent(p.MenuModel)
	p.Find.AptInstalledFromHistory(clearBytes, "all")
	p.Installed.AptAdvanced = p.Find.Packages
}

func (p *Packages) AptFromDir() {
	p.logging.Debug("RUN Packages.AptFromDir from controller (package_reader.go)")
	filesList := p.Tool.GetAptHistoryFilesList(p.MenuModel)
	for _, file := range filesList {
		p.MenuModel.FileName = file
		clearBytes := p.Tool.GetFileContent(p.MenuModel)
		p.Find.AptInstalledFromHistory(clearBytes, "all")
	}
	p.Installed.AptAdvanced = p.Find.Packages
}

func (p *Packages) AptAdded() {
	p.logging.Debug("RUN Packages.AptAdded from controller (package_reader.go)")
	p.MenuModel.DirName = "/var/log/apt"
	filesList := p.Tool.GetAptHistoryFilesList(p.MenuModel)
	for _, file := range filesList {
		p.MenuModel.FileName = file
		clearBytes := p.Tool.GetFileContent(p.MenuModel)
		p.Find.AptInstalledFromHistory(clearBytes, "added")
	}
	p.Installed.AptAdded = p.Find.Packages
}

func (p *Packages) AptManual() {
	p.logging.Debug("RUN Packages.AptManual from controller (package_reader.go)")
	p.MenuModel.DirName = "/var/log/apt"
	filesList := p.Tool.GetAptHistoryFilesList(p.MenuModel)
	for _, file := range filesList {
		p.MenuModel.FileName = file
		clearBytes := p.Tool.GetFileContent(p.MenuModel)
		p.Find.AptInstalledFromHistory(clearBytes, "manual")
	}
	p.Installed.AptManual = p.Find.Packages
}
