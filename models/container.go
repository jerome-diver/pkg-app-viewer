package model

import (
	"log/slog"
)

type Package int8

const (
	None Package = iota
	Apt
	RPM
	Pacman
	Zypper
	Nix
	Flatpak
	Snap
	Go
	Rust
	Rubygem
	Pip
	Source
)

func (p Package) String() string {
	switch p {
	case Apt:
		return "Apt"
	case Flatpak:
		return "Flatpak"
	case Snap:
		return "Snap"
	case Rust:
		return "Rust"
	case Rubygem:
		return "Rubygem"
	case Go:
		return "Go"
	case Pip:
		return "Pip"
	case Source:
		return "Source"
	}
	return "unknown"
}

type Search int8

const (
	All Search = iota
	Added
	OfficialRepos
	OtherRepos
	FileSource
)

func (s Search) String() string {
	switch s {
	case All:
		return "All"
	case Added:
		return "Added"
	case OfficialRepos:
		return "OfficialRepos"
	case OtherRepos:
		return "OtherRepos"
	case FileSource:
		return "FileSource"
	}
	return "unknown"
}

type List struct {
	All        []string
	Added      []string
	Official   []string
	Other      []string
	FileSource []string
}

type Installed struct {
	Apt     List
	Flatpak []string
	Snap    []string
	Rust    []string
	Rubygem []string
	Pip     []string
	Go      []string
	Source  []string
}

var logger *slog.Logger

func NewInstalled(log *slog.Logger) *Installed {
	logger = log
	logger.Debug("Created a new Installed struct of installed packages list")
	return new(Installed)
}
