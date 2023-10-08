package controller

import (
	"fmt"
	"log/slog"

	model "github.com/pkg-app-viewer/models"
	view "github.com/pkg-app-viewer/views"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Menu struct {
	Config *model.ConfigFile
	Model  *model.Menu
	View   *view.Menu
	logger *view.Logging
}

func NewMenu(logger *view.Logging, version string, config *viper.Viper) *Menu {
	menu := new(Menu)
	menu.logger = logger
	menu.Config = new(model.ConfigFile)
	menu.logger.Error = config.Unmarshal(menu.Config)
	menu.logger.CheckError("can not unmarshall config file")
	menu.logger.Log.Debug("Start Menu building process")
	menu.View = view.NewMenu(menu.Config)
	menu.Model = model.NewMenu()
	menu.View.SearchManager.Run = menu.execSearchPackageManager
	menu.initMenuCommand()
	menu.View.Root.PersistentFlags().BoolVarP(&menu.Model.ShowMeta, "meta", "g", false, "show meta of gz files")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.Debug, "debug", "d", "Error", "debug message printed mode [Error, Warn, Info, Debug]")
	menu.View.Root.PersistentFlags().BoolVarP(&menu.Model.Interactive, "interactive", "i", false, "Interactive terminal mode")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.Output.File, "outFile", "o", "", "Output file name")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.Output.Mode, "outMode", "m", "stdout", "Output mode")
	menu.View.Root.PersistentFlags().StringVarP(&menu.Model.Output.Format, "format", "f", "txt", "Output format type")
	menu.View.Root.PersistentPreRunE = menu.validateGeneralFlags
	menu.View.Root.Version = version
	menu.View.Root.AddCommand(menu.View.SearchManager)
	menu.View.Root.Execute()
	return menu
}
func (m *Menu) initMenuCommand() {
	switch m.Config.System.Distribution_ID {
	case "ubuntu":
		{
			m.View.Apt.Run = m.execApt
			m.View.Apt.Flags().StringVarP(&m.Model.DirName, "fromDir", "D", "", "indicate directory to search for apt history log files")
			m.View.Apt.Flags().StringVarP(&m.Model.FileName, "fromFile", "F", "", "indicate files to search for apt history log files")
			m.View.Root.AddCommand(m.View.Apt)
		}
	case "arch":
		{
			m.View.Pacman.Run = m.execPacman
			m.View.Root.AddCommand(m.View.Pacman)
		}
	case "redhat":
		{
			m.View.RPM.Run = m.execRPM
			m.View.Root.AddCommand(m.View.RPM)
		}
	case "gentoo":
		{
			m.View.Zypper.Run = m.execZypper
			m.View.Root.AddCommand(m.View.Zypper)
		}
	case "nixos":
		{
			m.View.Nix.Run = m.execNix
			m.View.Root.AddCommand((m.View.Nix))
		}
	}
	if m.Config.Packager.Flatpak != "" {
		m.View.Flatpak.Run = m.execFlatpak
		m.View.Root.AddCommand(m.View.Flatpak)
	}
	if m.Config.Packager.Snap != "" {
		m.View.Snap.Run = m.execSnap
		m.View.Root.AddCommand(m.View.Snap)
	}
	if m.Config.Packager.Go != "" {
		m.View.Go.Run = m.execGo
		m.View.Root.AddCommand(m.View.Go)
	}
	if m.Config.Packager.Python.Pip != "" {
		m.View.Pip.Run = m.execPip
		m.View.Root.AddCommand((m.View.Pip))
	}
	if m.Config.Packager.Rubygem != "" {
		m.View.Rubygem.Run = m.execRubygem
		m.View.Root.AddCommand(m.View.Rubygem)
	}
	if m.Config.Packager.Rust.Cabal != "" ||
		m.Config.Packager.Rust.Rustup != "" {
		m.View.Rust.Run = m.execRust
		m.View.Root.AddCommand(m.View.Rust)
	}
	if m.Config.Source != "" {
		m.View.Source.Run = m.execSource
		m.View.Root.AddCommand(m.View.Source)
	}
}

func (m *Menu) validateGeneralFlags(cmd *cobra.Command, arg []string) error {
	// check to validate OutputMode flag
	switch m.Model.Debug {
	case "Debug":
		break
	case "Info":
		break
	case "Warn":
		break
	case "Error":
		break
	default:
		return fmt.Errorf("Debug unvalid flag %s", m.Model.Debug)
	}
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
		break
	default:
		return fmt.Errorf("Output.Format unvalid flag %s", m.Model.Output.Format)
	}
	return nil
}

func (m *Menu) execSearchPackageManager(cmd *cobra.Command, arg []string) {
	m.Model.Command = model.SearchSystem_pm
}

func (m *Menu) execApt(cmd *cobra.Command, arg []string) {
	m.logger.Log.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
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
	m.logger.Log.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Flatpak
}

func (m *Menu) execRPM(cmd *cobra.Command, arg []string) {
	m.logger.Log.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.RPM
}

func (m *Menu) execPacman(cmd *cobra.Command, arg []string) {
	m.logger.Log.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Snap
}

func (m *Menu) execZypper(cmd *cobra.Command, arg []string) {
	m.logger.Log.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Snap
}

func (m *Menu) execNix(cmd *cobra.Command, arg []string) {
	m.logger.Log.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Snap
}

func (m *Menu) execSnap(cmd *cobra.Command, arg []string) {
	m.logger.Log.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Snap
}

func (m *Menu) execRubygem(cmd *cobra.Command, arg []string) {
	m.logger.Log.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Rubygem
}

func (m *Menu) execPip(cmd *cobra.Command, arg []string) {
	m.logger.Log.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Pip
}

func (m *Menu) execRust(cmd *cobra.Command, arg []string) {
	m.logger.Log.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Rust
}

func (m *Menu) execGo(cmd *cobra.Command, arg []string) {
	m.logger.Log.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Go
}

func (m *Menu) execSource(cmd *cobra.Command, arg []string) {
	m.logger.Log.Debug("Read dir argument from menu PackageType cmd", slog.String("arg[0]", arg[0]))
	m.Model.PackageType = model.Source
}
