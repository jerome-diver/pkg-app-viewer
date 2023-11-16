package identifier

/* Obtain system managers nfos supported by
   3 different families as:
	_ System (official 1st system used package manager)
	_ Isolated 'prensence of Flatpak and Snap)
	_ Language (presence of any recognized language own packager)
*/

import (
	"bufio"
	"bytes"
	"log/slog"
	"os"

	model "github.com/pkg-app-viewer/models"
	report "github.com/pkg-app-viewer/views/report"
)

var config model.Config
var logging report.Logging

// Build new iddentifier instance
func New() *model.Identity {
	logging = report.GetLogger()
	id := new(model.Identity)
	id.System = getSystemInfos()
	id.Isolated = getIsolatedInfos()
	id.Language = getLanguageInfos()
	return id
}

// static constant data to try
func file_system_possible() []string {
	return []string{
		"/etc/os-release",
		"/etc/SuSE-release",
		"/etc/gentoo-release",
	}
}

/*
	 wrapper system translator
		from data bytes found
		to model.Manager
*/
func getSystemManagerType(data []byte) model.ManagerName {
	data = bytes.Trim(data, "\"")
	if bytes.Contains(data, []byte("debian")) {
		return model.Dpkg
	}
	if bytes.Contains(data, []byte("redhat")) {
		return model.RPM
	}
	if bytes.Contains(data, []byte("archlinux")) {
		return model.Pacman
	}
	return model.None
}

func getSystemInfos() model.SystemId {
	var sys_info model.SystemId = model.SystemId{}
	logging.Debug("Start Searching for system package manager depend on system OS")
	const (
		system_tag_name = "NAME="
		system_tag_id   = "ID_LIKE="
	)
	// search if any file_system_possible file exist
	for _, f := range file_system_possible() {
		logging.Debug("Searching inside file", slog.String("file", f))
		_, err := os.Stat(f)
		logging.SetError(err)
		if os.IsNotExist(logging.GetError()) { // file does not exist => next
			continue
		}
		if logging.CheckError("File Stat has a problem") {
			return sys_info
		}
		var byte_content []byte
		byte_content, err = os.ReadFile(f)
		logging.SetError(err)
		if logging.CheckError("Can not read file") {
			return sys_info
		}
		logging.Debug("Reading file")
		/* and match tags to obtain:
		   _ distributions OS type ID
		   _ distributions OS Name    */
		scanner := bufio.NewScanner(bytes.NewReader(byte_content))
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Bytes()
			logging.Debug("Scanning file", slog.String("line", string(line)))
			if _, after, found := bytes.Cut(line, []byte(system_tag_id)); found {
				sys_info.Type = getSystemManagerType(after)
				if sys_info.Name != "" {
					return sys_info
				}
				continue
			}
			if _, after, found := bytes.Cut(line, []byte(system_tag_name)); found {
				logging.Debug("found OS name", slog.String("after", string(after)))
				sys_info.Name = string(after)
				if sys_info.Type != model.None {
					return sys_info
				}
				continue
			}
		}
	}
	return sys_info
}

func getIsolatedInfos() model.NoSystemId {
	var isolated_infos model.NoSystemId = model.NoSystemId{}
	var home string
	home, err := os.UserHomeDir()
	logging.SetError(err)
	if logging.CheckError("Can not find Home directory") {
		return isolated_infos
	}
	const isolated_snap_dir = "/var/lib/snapd/snaps/"
	_, err = os.Stat(isolated_snap_dir)
	logging.SetError(err)
	if !logging.CheckError("Snap directory Stat problem") {
		isolated_infos.Types = append(isolated_infos.Types, model.Snap)
	}
	var isolated_flatpak_dir = home + ".var/app/"
	_, err = os.Stat(isolated_flatpak_dir)
	logging.SetError(err)
	if !logging.CheckError("Flatpak directory Stat problem") {
		isolated_infos.Types = append(isolated_infos.Types, model.Flatpak)
	}
	return isolated_infos
}

func getLanguageInfos() model.NoSystemId {
	var languages_infos model.NoSystemId = model.NoSystemId{}
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
	return languages_infos
}
