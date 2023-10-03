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
	menu.View.Apt.Run = menu.execAptType
	menu.View.FromFile.Run = menu.execAptFile
	menu.View.FromDir.Run = menu.execAptDir
	menu.View.Rust.Run = menu.execRustType
	menu.View.Go.Run = menu.execGoType
	menu.View.Flatpak.Run = menu.execFlatpakType
	menu.View.Snap.Run = menu.execSnapType
	menu.View.Source.Run = menu.execSourceType
	menu.View.Root.PersistentFlags().BoolVarP(&menu.Model.ShowMeta, "meta", "m", false, "show meta of gz files")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.Debug, "debug", "d", "Error", "debug message printed mode [Error, Warn, Info, Debug]")
	menu.View.Root.PersistentFlags().BoolVarP(&menu.Model.Interactive, "interactive", "i", false, "Interactive terminal mode")
	menu.View.FromDir.Flags().StringVarP(&menu.Model.DirName, "aptDir", "D", "/var/log/apt", "indicate directory to search for apt history log files")
	menu.View.FromFile.Flags().StringVarP(&menu.Model.FileName, "aptFile", "F", "/var/log/apt/history.log", "indicate files to search for apt history log files")
	menu.View.Root.AddCommand(menu.View.Apt, menu.View.Flatpak, menu.View.Snap,
		menu.View.Rust, menu.View.Go, menu.View.Source)
	menu.View.Apt.AddCommand(menu.View.FromFile, menu.View.FromDir)
	menu.View.Root.Execute()
	return menu
}

func (m *Menu) execAptFile(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read file argument from menu execAptFile cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.AptAdvanced
	m.Model.Mode = "File"
	m.Model.FileName = arg[0]
}

func (m *Menu) execAptDir(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu execAptDir cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.AptAdvanced
	m.Model.Mode = "Directory"
	m.Model.DirName = arg[0]
}

func (m *Menu) execAptType(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	switch arg[0] {
	case "All":
		m.Model.PackageType = model.AptAll
	case "Added":
		m.Model.PackageType = model.AptAdded
	case "OfficialAdded":
		m.Model.PackageType = model.AptOfficialAdded
	case "OtherRepos":
		m.Model.PackageType = model.AptOtherRepos
	case "Manual":
		m.Model.PackageType = model.AptManual
	default:
		m.Model.PackageType = model.AptAll
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
