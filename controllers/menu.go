package controller

import (
	"log/slog"

	model "github.com/pkg-app-viewer/models"
	view "github.com/pkg-app-viewer/views"

	"github.com/spf13/cobra"
)

type Menu struct {
	Model  *model.Menu
	View   *view.Menu
	logger *slog.Logger
}

func NewMenu(logger *slog.Logger) *Menu {
	menu := new(Menu)
	menu.logger = logger
	menu.View = view.NewMenu()
	menu.Model = model.NewMenu()
	menu.View.Apt.Run = menu.execApt
	menu.View.Rust.Run = menu.execRustType
	menu.View.Go.Run = menu.execGoType
	menu.View.Flatpak.Run = menu.execFlatpakType
	menu.View.Snap.Run = menu.execSnapType
	menu.View.Source.Run = menu.execSourceType
	menu.View.Root.PersistentFlags().BoolVarP(&menu.Model.ShowMeta, "meta", "g", false, "show meta of gz files")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.Debug, "debug", "d", "Error", "debug message printed mode [Error, Warn, Info, Debug]")
	menu.View.Root.PersistentFlags().BoolVarP(&menu.Model.Interactive, "interactive", "i", false, "Interactive terminal mode")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.OutputFile, "outFile", "o", "pkg_list.txt", "Output file name")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.OutputMode, "outMode", "m", "stdout", "Output mode")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.Format, "format", "f", "txt", "Output format type")
	menu.View.Apt.Flags().StringVarP(&menu.Model.DirName, "fromDir", "D", "", "indicate directory to search for apt history log files")
	menu.View.Apt.Flags().StringVarP(&menu.Model.FileName, "fromFile", "F", "", "indicate files to search for apt history log files")
	menu.View.Root.AddCommand(menu.View.Apt, menu.View.Flatpak, menu.View.Snap,
		menu.View.Rust, menu.View.Go, menu.View.Source)
	menu.View.Root.Execute()
	return menu
}

func (m *Menu) execApt(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Apt
	switch arg[0] {
	case "All":
		m.Model.PackageSearch = model.All
	case "Added":
		m.Model.PackageSearch = model.Added
	case "OfficialAdded":
		m.Model.PackageSearch = model.OfficialRepos
	case "OtherRepos":
		m.Model.PackageSearch = model.OtherRepos
	case "FileSource":
		m.Model.PackageSearch = model.FileSource
	default:
		m.Model.PackageSearch = model.All
	}
}

func (m *Menu) execRustType(cmd *cobra.Command, arg []string) {
	m.Model.PackageType = model.Rust
}

func (m *Menu) execGoType(cmd *cobra.Command, arg []string) {
	m.Model.PackageType = model.Go
}

func (m *Menu) execSourceType(cmd *cobra.Command, arg []string) {
	m.Model.PackageType = model.Source
}

func (m *Menu) execSnapType(cmd *cobra.Command, arg []string) {
	m.Model.PackageType = model.Snap
}

func (m *Menu) execFlatpakType(cmd *cobra.Command, arg []string) {
	m.Model.PackageType = model.Flatpak
}
