package identifier

import (
	engine "github.com/pkg-app-viewer/controllers/engine"
	model "github.com/pkg-app-viewer/models"
)

type Info interface {
	GetTag() string
	GetValue() any
}

type Manager interface {
	Init(option Options)
	GetInfos() []Info
	GetPackages() []string
	GetFull() engine.InstalledSystemRepos
}

type Options map[string]any

type SystemManager struct {
	engine.InstalledSystemRepos
}

func (m SystemManager) Init(options Options) {
}

func (m SystemManager) GetInfos() []Info {
	data := []Info{}
	return data
}

func (m SystemManager) GetPackages() []string {
	data := []string{}
	return data
}

func (m SystemManager) GetFull() engine.InstalledSystemRepos {
	return m.InstalledSystemRepos
}

func NewManager(manager model.Manager, options Options) Manager {
	var obj Manager
	switch manager {
	case model.Dpkg:
		{
			obj = new(SystemManager)
			obj.Init(options)
		}
	}
	return obj
}
