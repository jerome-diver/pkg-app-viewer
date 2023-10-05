# Application pkg-app-viewer

## Introduction
Help to show list of application installed by packager.
Recognize distribution and packager log that can be used.
Has a TUI (pterm) and a some commandline/tags call (cobra) with prompt (promptui).
Has been written in Go 100%.
Can read application installed from any other sources or mechanism and hope to be easy to extend.
0% shell command call (strategy is to search files presence and crawl inside the important one).
It does not pretend to be a challenger for any velocity competition or multi-task orientation but to be friendly usable and easy to use.
It should be able to be used inside a script to generate something in pipe (to backup some lists for example).

## Package manager installed from
- apt, dpkg (debian / Ubuntu like distros)
- rpm, dnf (ReHat distros) 
- pacman, pamac (Arch distros)
- Zypper (Gentoo distros)
- Nix (NixOS distribution)
- Flatpak
- Snap
  
## Other installation mode
- Go
- Rust (cargo, rustup)
- Ruby (rubygem)
- Python (pipenv, pip)
- Perl-5 (cpan)
- Node.js (npm, yarn)
- source (tar.gz or any other at any places)
  
## Options usefull
You can ask fro specific applications installed origin from, like:
- All applications (no dependencies)
- All application installed by user
- All packages from other non-official repositories
- All packages installed from a package's file source
- You can redirect the output in a file or stdout in many format

## What can it be used for
To create files for backup system's applications (this files can be used later to reinstall applications from fresh OS install).

## Dependencies used
- [Go version 1.21](https://go.dev)
- [cobra](https://github.com/spf13/cobra/tree/v1.7.0) module (to get a command/tag menu shell call)
- [pterm](https://github.com/pterm/pterm) module (to add nice TUI)
- [promptui](https://github.com/manifoldco/promptui) module (to add easy prompt)
- [fatih/color](https://github.com/fatih/color) module (to addd color to logging verbose)

## Main design pattern used
MVC patern application.