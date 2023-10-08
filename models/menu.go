package model

type Output struct {
	File   string
	Mode   string
	Format string
}

type Menu struct {
	Command       CommandRun
	PackageType   Package
	PackageSearch Search
	FileName      string
	DirName       string
	Mode          string
	ShowMeta      bool
	Debug         string
	Output        Output
	Interactive   bool
}

func NewMenu() *Menu {
	return new(Menu)
}
