package model

import "regexp"

type CommandRun int8

const (
	nada CommandRun = iota
	SearchSystem_pm
	SearchIsolated_pm
	SearchLanguage_pm
)

func (c CommandRun) String() string {
	switch c {
	case SearchIsolated_pm:
		return "SearchIsolatedPM"
	case SearchSystem_pm:
		return "SearchSystemPM"
	case SearchLanguage_pm:
		return "SearchLanguagePM"
	}
	return "unknown"
}

type Package int8

const (
	None Package = iota
	Apt
	RPM
	Pacman
	Zypper
	Emerge
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

func (s Search) Algorythm() func(string) bool {
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
