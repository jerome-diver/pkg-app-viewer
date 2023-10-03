package model

type Menu struct {
	PackageType Package
	FileName    string
	DirName     string
	Mode        string
	ShowMeta    bool
	Debug       string
	Output      string
	Interactive bool
}

func NewMenu() *Menu {
	return new(Menu)
}
