package view

import (
	model "github.com/pkg-app-viewer/models"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

type Menu struct {
	Root          *cobra.Command
	SearchManager *cobra.Command
	Apt           *cobra.Command
	RPM           *cobra.Command
	Pacman        *cobra.Command
	Zypper        *cobra.Command
	Nix           *cobra.Command
	Flatpak       *cobra.Command
	Snap          *cobra.Command
	Rust          *cobra.Command
	Rubygem       *cobra.Command
	Go            *cobra.Command
	Pip           *cobra.Command
	Source        *cobra.Command
}

func NewMenu(id *model.Identity) *Menu {
	vm := new(Menu)
	cobra.EnableCommandSorting = false
	vm.Root = &cobra.Command{
		Use:   "pkg-apt-viewer",
		Short: "Helper in Go to find installed aplicatins packages",
		Long:  "You can find for many packages manager installed application from (see command list to call them)",
	}
	vm.SearchManager = &cobra.Command{
		Use:   "SearchManager",
		Short: "Search package managers",
		Long:  "Search in the system any package managers to update config file and menu",
		Args:  cobra.NoArgs,
	}
	switch id.System.GetTypes() {
	case model.Dpkg:
		vm.Apt = &cobra.Command{
			Use:       "Apt [All, Added, OfficialAdded, FileSource, OtherRepos]",
			Short:     "Select Debian like package type",
			Long:      "Select Debian (dpkg/apt/aptitude) package manager type and content to get list from",
			ValidArgs: []string{"All", "Added", "OfficialAdded", "OtherRepos", "FileSource"},
			Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		}
	case model.RPM:
		vm.RPM = &cobra.Command{
			Use:       "RPM [All, Added, OfficialAdded, FileSource, OtherRepos]",
			Short:     "Select Red Hat like package type",
			Long:      "Select Red Hat (rpm, yumm, dnf) package manager type and content to get list from",
			ValidArgs: []string{"All", "Added", "OfficialAdded", "OtherRepos", "FileSource"},
			Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		}
	case model.Pacman:
		vm.Pacman = &cobra.Command{
			Use:       "Pacman [All, Added, OfficialAdded, FileSource, OtherRepos]",
			Short:     "Select Arch like package type",
			Long:      "Select Archlinux (pacman, pamac, conde) package manager type and content to get list from",
			ValidArgs: []string{"All", "Added", "OfficialAdded", "OtherRepos", "FileSource"},
			Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		}
	case model.Zypper:
		vm.Zypper = &cobra.Command{
			Use:       "Zypper [All, Added, OfficialAdded, FileSource, OtherRepos]",
			Short:     "Select Gentoo like package type",
			Long:      "Select Gentoo (zypper) package manager type and content to get list from",
			ValidArgs: []string{"All", "Added", "OfficialAdded", "OtherRepos", "FileSource"},
			Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		}
	case model.Nix:
		vm.Nix = &cobra.Command{
			Use:       "Nix [All, Added, OfficialAdded, FileSource, OtherRepos]",
			Short:     "Select NixOS like package type",
			Long:      "Select NixOS (nix) package manager type and content to get list from",
			ValidArgs: []string{"All", "Added", "OfficialAdded", "OtherRepos", "FileSource"},
			Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		}
	}
	if slices.Contains(id.Isolated.GetTypes(), model.Flatpak) {
		vm.Flatpak = &cobra.Command{
			Use:   "Flatpak",
			Short: "Select Flatpak package type",
			Long:  "Select Flatpak package manager type and content to get list from",
			Args:  cobra.NoArgs,
		}
	}
	if slices.Contains(id.Isolated.GetTypes(), model.Snap) {
		vm.Snap = &cobra.Command{
			Use:   "Snap",
			Short: "Select Snap package type",
			Long:  "Select Snap package manager type and content to get list from",
			Args:  cobra.NoArgs,
		}
	}
	if slices.Contains(id.Language.GetTypes(), model.Rustup) {
		vm.Rust = &cobra.Command{
			Use:   "Rust",
			Short: "Select Rust package type",
			Long:  "Select Rust package manager type and content to get list from",
			Args:  cobra.NoArgs,
		}
	}
	if slices.Contains(id.Language.GetTypes(), model.Rubygem) {
		vm.Rubygem = &cobra.Command{
			Use:   "Rubygem",
			Short: "Select Rubygem manager package type",
			Long:  "Select Rubygem package manager type and content to get list from",
			Args:  cobra.NoArgs,
		}
	}
	if slices.Contains(id.Language.GetTypes(), model.Go) {
		vm.Go = &cobra.Command{
			Use:   "Go",
			Short: "Select Go package type",
			Long:  "Select Go package manager type and content to get list from",
			Args:  cobra.NoArgs,
		}
	}
	if slices.Contains(id.Language.GetTypes(), model.Pip) {
		vm.Pip = &cobra.Command{
			Use:   "Pip",
			Short: "Select Pip manager package type",
			Long:  "Select Python pip package manager type and content to get list from",
			Args:  cobra.NoArgs,
		}
	}
	return vm
}
