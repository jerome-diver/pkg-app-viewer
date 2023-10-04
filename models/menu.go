package model

type Menu struct {
	PackageType   Package
	PackageSearch Search
	FileName      string
	DirName       string
	Mode          string
	ShowMeta      bool
	Debug         string
	OutputFile    string
	OutputMode    string
	Format        string
	Interactive   bool
}

func NewMenu() *Menu {
	return new(Menu)
}
