package controller

import (
	"fmt"
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
	menu.View.RPM.Run = menu.execRPM
	menu.View.Pacman.Run = menu.execPacman
	menu.View.Zypper.Run = menu.execZypper
	menu.View.Nix.Run = menu.execNix
	menu.View.Rust.Run = menu.execRust
	menu.View.Go.Run = menu.execGo
	menu.View.Flatpak.Run = menu.execFlatpak
	menu.View.Snap.Run = menu.execSnap
	menu.View.Rubygem.Run = menu.execRubygem
	menu.View.Pip.Run = menu.execPip
	menu.View.Source.Run = menu.execSource
	menu.View.Root.PersistentFlags().BoolVarP(&menu.Model.ShowMeta, "meta", "g", false, "show meta of gz files")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.Debug, "debug", "d", "Error", "debug message printed mode [Error, Warn, Info, Debug]")
	menu.View.Root.PersistentFlags().BoolVarP(&menu.Model.Interactive, "interactive", "i", false, "Interactive terminal mode")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.Output.File, "outFile", "o", "", "Output file name")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.Output.Mode, "outMode", "m", "stdout", "Output mode")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.Output.Format, "format", "f", "txt", "Output format type")
	menu.View.Apt.Flags().StringVarP(&menu.Model.DirName, "fromDir", "D", "", "indicate directory to search for apt history log files")
	menu.View.Apt.Flags().StringVarP(&menu.Model.FileName, "fromFile", "F", "", "indicate files to search for apt history log files")
	menu.View.Root.AddCommand(menu.View.Apt, menu.View.RPM, menu.View.Pacman,
		menu.View.Zypper, menu.View.Nix,
		menu.View.Flatpak, menu.View.Snap,
		menu.View.Rust, menu.View.Rubygem, menu.View.Pip,
		menu.View.Go, menu.View.Source)
	menu.View.Root.PersistentPreRunE = menu.validateGeneralFlags
	menu.View.Root.Execute()
	return menu
}

func (m *Menu) validateGeneralFlags(cmd *cobra.Command, arg []string) error {
	// check to validate OutputMode flag
	switch m.Model.Output.Mode {
	case "stdout":
		break
	case "file":
		break
	default:
		return fmt.Errorf("OutputMode unvalid flag %s", m.Model.Output.Mode)
	}
	switch m.Model.Output.Format {
	case "txt":
		break
	case "json":
		break
	case "yaml":
		break
	case "csv":
	default:
		return fmt.Errorf("Output.Format unvalid flag %s", m.Model.Output.Format)
	}
	return nil
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

func (m *Menu) execFlatpak(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Flatpak
}

func (m *Menu) execRPM(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.RPM
}

func (m *Menu) execPacman(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Snap
}

func (m *Menu) execZypper(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Snap
}

func (m *Menu) execNix(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Snap
}

func (m *Menu) execSnap(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Snap
}

func (m *Menu) execRubygem(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Rubygem
}

func (m *Menu) execPip(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Pip
}

func (m *Menu) execRust(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Rust
}

func (m *Menu) execGo(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Go
}

func (m *Menu) execSource(cmd *cobra.Command, arg []string) {
	m.logger.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Source
}
