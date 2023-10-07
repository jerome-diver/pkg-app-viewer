package model

type ConfigFile struct {
	System struct {
		Distribution_ID string
		Packager        []string
	}
	Packager struct {
		Go   string
		Rust struct {
			Rustup string
			Cabal  string
		}
		Python struct {
			Pip    string
			Pipenv string
			Poetry string
		}
		Perl5 struct {
			Cpan string
		}
		Node struct {
			Npm  string
			Yarn string
		}
		Snap    string
		Flatpak string
		Rubygem string
	}
	Source string
}
