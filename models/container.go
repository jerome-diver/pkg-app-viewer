package model

import (
	"regexp"
	"slices"
)

/*
Identify package manager family type

the list of const is the full
list of manager that can be
care about
*/
type ManagerName int8

const (
	None ManagerName = iota
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
)

func (p ManagerName) String() string {
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
	}
	return "unknown"
}

/*
Identify packages mode request filter

	can send back a
	- string
	- specific algorithm
*/
type ManagerOption int8

const (
	All ManagerOption = iota
	User
	System
	Distribution
	Foreign
	FileSource
)

func (s ManagerOption) String() string {
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

func (s ManagerOption) Algorythm() func(string) bool {
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

/*
Contain Repo packages datas as

	Origin name
	Packages as list (slices) of packages
*/
type Repository struct { // Will contain the origin and packages list (uniq)
	Origin   string
	Packages []string
}

/*
	 Main Identity struct handle managers infos

	   Hold 33 different type:
		- System is for system uniq package manager info
		- Isolated is for any independant container package managers
		- Language is for specific language's package managers information
*/
type Identity struct {
	System   ManagersInfos[ManagerName]
	Isolated ManagersInfos[[]ManagerName]
	Language ManagersInfos[[]ManagerName]
}

/*
generic interface to handle identity Manager

	answer is one or any infos format interface to embed
*/
type ManagersInfos[T format] interface {
	GetTypes() T
	AsType(ManagerName) bool
	GetStruct() ManagersInfos[T]
}

type format interface {
	ManagerName | []ManagerName
}

/*
System Identity struct

	will handle ManagersInfos interface
	to hold Sytem package managers informations
	And will have only one type of Manager
*/
type SystemId struct {
	Type ManagerName
	Name string
	Arch string
}

// Handle ManagersInfos[T] interface
func (manager SystemId) AsType(t ManagerName) bool {
	return manager.Type == t
}

func (manager SystemId) GetTypes() ManagerName {
	return manager.Type
}

func (manager SystemId) GetStruct() ManagersInfos[ManagerName] {
	return manager
}

/*
Other packages managers Identity struct

	will handle ManagersInfos interface
	to hold managers informations
	And will have many possible type of Manager
*/
type NoSystemId struct {
	Types []ManagerName
	User  string
}

// Handle ManagersInfos[T] interface
func (manager NoSystemId) AsType(t ManagerName) bool {
	return slices.Contains(manager.Types, t)
}

func (manager NoSystemId) GetTypes() []ManagerName {
	return manager.Types
}

func (manager NoSystemId) GetStruct() ManagersInfos[[]ManagerName] {
	return manager
}
