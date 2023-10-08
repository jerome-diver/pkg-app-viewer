package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"os"

	model "github.com/pkg-app-viewer/models"
	view "github.com/pkg-app-viewer/views"
)

type SearchManagers struct {
	logger   *view.Logging
	Config   *model.ConfigFile
	Managers PackageManagers
}

type PackageManagers struct {
	System struct {
		Apt    bool
		RPM    bool
		Pacman bool
		Zypper bool
		Nix    bool
		Name   string
	}
	Isolated struct {
		Flatpak bool
		Snap    bool
	}
	Language struct {
		Go      bool
		Rustup  bool
		Cabal   bool
		Rubygem bool
		Pip     bool
		Npm     bool
	}
}

func NewSearchManager(logger *view.Logging, config *model.ConfigFile) *SearchManagers {
	s := new(SearchManagers)
	s.logger = logger
	s.Config = config
	return s
}

func file_system_possible() []string {
	return []string{
		"/etc/os-release",
		"/etc/SuSE-release",
		"/etc/gentoo-release",
	}
}

func (s *SearchManagers) analyseSystemDataPM(data []byte) {
	data = bytes.Trim(data, "\"")
	s.logger.Log.Debug("Trimed data found", slog.String("data", string(data)))
	if bytes.Contains(data, []byte("debian")) {
		s.logger.Log.Info("It does contain debian apt package manager")
		s.Managers.System.Apt = true
	}
}

func (s *SearchManagers) SearchSystemPM() bool {
	s.logger.Log.Debug("Start Searching for system package manager depend on system OS")
	const (
		//	system_name_tag = "NAME="
		system_id_tag = "ID_LIKE="
	)
	// search if any file_system_possible file exist
	for _, f := range file_system_possible() {
		s.logger.Log.Debug("Searching inside file", slog.String("file", f))
		_, err := os.Stat(f)
		if err == nil { // does exist, then open and read content in slice byte
			s.logger.CheckError("Can not open file")
			var byte_content []byte
			byte_content, s.logger.Error = os.ReadFile(f)
			s.logger.CheckError("Can not read file")
			s.logger.Log.Debug("Reading file")
			// and match tags to obtain distributions OS type ID
			scanner := bufio.NewScanner(bytes.NewReader(byte_content))
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				line := scanner.Bytes()
				s.logger.Log.Debug("Scanning file", slog.String("line", string(line)))
				if _, after, found := bytes.Cut(line, []byte(system_id_tag)); found {
					s.logger.Log.Debug("found OS id tag", slog.String("after", string(after)))
					s.analyseSystemDataPM(after)
					return true // end of search process
				}
			}
		}
		if os.IsNotExist(err) {
			s.logger.Log.Debug("file doesn't exist")
			continue
		} else {
			s.logger.Error = err
			s.logger.CheckError("Stat error")
		}
	}
	return false
}

func (s *SearchManagers) SearchIsolatedPM() {
	var home string
	home, s.logger.Error = os.UserHomeDir()
	s.logger.CheckError("Can not find Home directory")
	const isolated_snap_dir = "/var/lib/snapd/snaps/"
	var isolated_flatpak_dir = home + ".var/app/"
	fmt.Println(isolated_flatpak_dir)
}

func (s *SearchManagers) SearchLanguagePM() {
	//const (
	// Go has his own tool set to manage packages,
	// But it does exist also Cask.
	//	language_go_dir = home + "go/bin/"
	// Ruby can use a rugygem package manager,
	// but can be use in many different environment
	// to rich Ruby version required
	// we know and recognize these one:
	//   - rbenv
	//   - rvm
	//   - asdf-ruby
	//	language_ruby_rubygem_dir = "/var/lib/gems/" // + rubygem_version + "/gems/"
	//	language_perl5_dir        = []string{home + ".cpan/", home + ".perl5/"}
	// Python can use many package manager
	// it can also be used in virtual environment
	// for specific Python version to use
	// package manager we know about here are:
	//   - pip
	//   - pipx
	//   - pyenv
	//   - poetry
	//   - asdf-python
	//	language_python_dirs = map[string]string{
	//		"pipx":   home + ".local/pipx/shared/bin/",
	//		"pip":    home + ".pip/",
	//		"pyenv":  home + ".pyenv/",
	//		"poetry": home + ".poerty/"}
	//	language_haskell_dirs = map[string]string{
	//		"stack": home + ".stack/",
	//		"cabal": home + ".cabal/bin/"}
	//	language_rust_dir = home + ".cargo/bin/" // Rust use only Cargo
	// Javascript Engine that can run on the system, most known are:
	//  - Node.js (Google Chrome V8 engine, written in C++)
	//  - Deno.js (most secure, Typescript ready and written in Rust)
	//  - Bun.js  (JSC Apple Engine, written in Zig)
	// And to rich one or the other dependencies version or packages,
	// There is, for Node.js, many possible package manager to use
	// and some of them also can run versioning Node.js system.
	// For npm, to find current Node.js dir used, you need to run
	// `npm root -g` (because of nvm, fnm or volta can install
	//	on any place due to version requested)
	// shell command to get the direcory location
	// So, to manage the Node.js version
	// (not only the main packager can be used) there is some
	// alternatives Node.js version manager as:
	//	   - fnm (fast node manager, 65% Rust rewriting)
	//     - nvm (node version manager, shell scripts)
	//	   - asdf-node (version manager with node plugin)
	//     - volta (package and version manager writing in Rust 100%)
	// For Deno.js, package manager is not requested
	// due to his own url call dependencies.
	// But there is some Deno version manager as "dvm":
	//     - dvm (like fnm, thereis many different apps
	//             writing in Ocaml or in Shell or in Go or Rust)
	//     - Trex
	// Bun also has its own package manager.
	// At final time, Javascript Engine are too much
	// complicate to manage due to infinite possiblities
	// and i will maybe manage them (or not).
	//language_javascript_dirs = map[string]any{
	//	"nvm": home + ".nvm/",
	//	"npm": "/usr/local/lib/node_modules/npm/node_modules/"}
	//)

}

func (s *SearchManagers) searchApt() {

}
