package model

import (
	"log/slog"
	"slices"
)

type Package int8

const (
	Undefined Package = iota
	AptAll
	AptAdded
	AptOfficialAdded
	AptOtherRepos
	AptManual
	AptAdvanced
	Flatpak
	Snap
	Rust
	Go
	Source
)

func (p Package) String() string {
	switch p {
	case Undefined:
		return "undefined"
	case AptAll:
		return "AptOfficialAll"
	case AptAdded:
		return "AptAdded"
	case AptOfficialAdded:
		return "AptOfficialAdded"
	case AptOtherRepos:
		return "AptOtherRepos"
	case AptManual:
		return "AptManual"
	case AptAdvanced:
		return "AptAdvanced"
	case Flatpak:
		return "Flatpak"
	case Snap:
		return "Snap"
	case Rust:
		return "Rust"
	case Go:
		return "Go"
	case Source:
		return "Source"
	}
	return "unknown"
}

type Installed struct {
	AptAll           []string
	AptAdded         []string
	AptOfficialAdded []string
	AptOtherRepos    []string
	AptManual        []string
	AptAdvanced      []string
	Flatpak          []string
	Snap             []string
	Rust             []string
	Go               []string
	Source           []string
}

var logger *slog.Logger

func NewInstalled(log *slog.Logger) *Installed {
	logger = log
	logger.Debug("Created a new Installed struct of installed packages list")
	return new(Installed)
}

func (i *Installed) sort(list []string) {
	slices.Sort(list)
}

func (i *Installed) GetInstalledPackagesList(p Package) []string {
	switch p {
	case AptAll:
		return i.AptAll
	case AptAdded:
		return i.AptAdded
	case AptOfficialAdded:
		return i.AptOfficialAdded
	case AptOtherRepos:
		return i.AptOtherRepos
	case AptManual:
		return i.AptManual
	case AptAdvanced:
		return i.AptAdvanced
	case Flatpak:
		return i.Flatpak
	case Snap:
		return i.Snap
	case Rust:
		return i.Rust
	case Go:
		return i.Go
	case Source:
		return i.Source
	}
	return []string{""}
}
