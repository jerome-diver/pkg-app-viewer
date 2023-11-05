package model

import "regexp"

type CommandRun int8

const (
	System_managers CommandRun = iota
	Isolated_managers
	Language_managers
)

func (c CommandRun) String() string {
	switch c {
	case Isolated_managers:
		return "Isolated Managers"
	case System_managers:
		return "System Managers"
	case Language_managers:
		return "Language Managers"
	}
	return "unknown"
}

type Manager int8

const (
	None Manager = iota
	Dpkg
	RPM
	Pacman
	Zypper
	Emerge
	Nix
	Flatpak
	Snap
	Go
	Cargo
	Rustup
	Rubygem
	Pip
	Source
)

func (p Manager) String() string {
	switch p {
	case Dpkg:
		return "Debian like dpkg"
	case RPM:
		return "RedHat like rpm"
	case Pacman:
		return "Archlinux like pacman"
	case Flatpak:
		return "Flatpak manager"
	case Snap:
		return "Snap manager"
	case Cargo:
		return "Haskell cargo"
	case Rustup:
		return "Rust rustup"
	case Rubygem:
		return "Ruby rubygem"
	case Go:
		return "Golang go"
	case Pip:
		return "Python pip"
	case Source:
		return "Source file"
	}
	return "unknown"
}

type SystemOption int8

const (
	All SystemOption = iota
	User
	System
	Distribution
	Foreign
	FileSource
)

func (s SystemOption) String() string {
	switch s {
	case All:
		return "All"
	case User:
		return "User Installed"
	case System:
		return "Official Repos"
	case Foreign:
		return "Foreign Repos"
	case FileSource:
		return "File Source"
	case Distribution:
	}
	return "unknown"
}

func (s SystemOption) Algorythm() func(string) bool {
	switch s {
	case FileSource:
		return func(p string) bool {
			re := regexp.MustCompile(`.*\.deb$`)
			return re.MatchString(p)
		}
	default:
		return func(p string) bool { return true }
	}
}

type SystemRepo struct { // Will contain the origin and packages list (uniq)
	Origin   string
	Packages []string
}
