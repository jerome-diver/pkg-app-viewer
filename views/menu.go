package view

import "github.com/spf13/cobra"

type Menu struct {
	Root     *cobra.Command
	Apt      *cobra.Command
	FromFile *cobra.Command
	FromDir  *cobra.Command
	Flatpak  *cobra.Command
	Snap     *cobra.Command
	Source   *cobra.Command
	Rust     *cobra.Command
	Go       *cobra.Command
}

func NewMenu() *Menu {
	vm := new(Menu)
	vm.Root = &cobra.Command{Use: "pkg_installed"}
	vm.Apt = &cobra.Command{
		Use:       "Apt [All, Added, OfficialAdded, Manual, OtherRepos]",
		Short:     "Select Apt package type",
		Long:      "Select Apt package manager type and content to get list from",
		ValidArgs: []string{"All", "Added", "OfficialAdded", "Manual", "OtherRepos"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}
	vm.Flatpak = &cobra.Command{
		Use:   "Flatpak",
		Short: "Select Flatpak package type",
		Long:  "Select Flatpak package manager type and content to get list from",
		Args:  cobra.NoArgs,
	}
	vm.Snap = &cobra.Command{
		Use:   "Snap",
		Short: "Select Snap package type",
		Long:  "Select Snap package manager type and content to get list from",
		Args:  cobra.NoArgs,
	}
	vm.Source = &cobra.Command{
		Use:   "Source",
		Short: "Select Source  package type",
		Long:  "Select Source package manager type and content to get list from",
		Args:  cobra.MinimumNArgs(0),
	}
	vm.Rust = &cobra.Command{
		Use:   "Rust",
		Short: "Select Rust package type",
		Long:  "Select Rust package manager type and content to get list from",
		Args:  cobra.NoArgs,
	}
	vm.Go = &cobra.Command{
		Use:   "Go",
		Short: "Select Go package type",
		Long:  "Select Go package manager type and content to get list from",
		Args:  cobra.NoArgs,
	}
	vm.FromFile = &cobra.Command{
		Use:   "fromFile [file-name] ",
		Short: "Indicate file name to extract from",
		Long:  `Go inside the file (decode gz if required), to search for packages installed`,
		Args:  cobra.MinimumNArgs(1),
	}
	vm.FromDir = &cobra.Command{
		Use:   "fromDir [dir-name] ",
		Short: "Indicate directory name to search files from",
		Long:  `Go inside the directory to find any history files, to search for packages installed`,
		Args:  cobra.ExactArgs(1),
	}
	return vm
}
