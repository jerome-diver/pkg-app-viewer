# Application pkg-app-viewer

## Introduction
Help to show list of application installed by packager.
Has a TUI (pterm) and a some commandline call (cobra)
Has been written in Go 100%
can read application installed from any other sources or mechanism
0% shell command call
It does not pretend to be a challenger for anyvelocity competition or multi-task orientation but to be friendly usable and easy to use.
It should be able to be used inside a script to generate something in pipe.
## Package manager installed from
- apt
- dnf
- pacman
- dpkg
- Flatpak
- Snap
  
## Other installation mode
- Rust
- Go
- Python-3 (pip)
- Ruby
- source (tar.gz or any other at any places)
  
## Options usefull
You can ask fro specific applications installed, like:
- All applications (no dependencies)
- All application installed by user
- All packages from other non-official repositories
- All packages installed from a package's file source
- You can redirect the output in a file or stdout in many format

## What can it be used for
To create files for backup system's applications (this files can be used later to reinstall applications from fresh OS install).
## Dependencies used
- Go version 1.21
- cobra module
- pterm module
## Main design pattern used
MVC patern application.