package model

/*
Menu printer mode format
*/
type Output struct {
	File   string
	Mode   string
	Format string
}

/*
Menu Flag to define
type of manager to use
*/
type CommandRun int8

const (
	System_managers CommandRun = iota
	Isolated_managers
	Language_managers
)

func (c CommandRun) String() string {
	switch c {
	case Isolated_managers:
		return "Isolated Managers"
	case System_managers:
		return "System Managers"
	case Language_managers:
		return "Language Managers"
	}
	return "unknown"
}

/*
Container of Menu flags options
choices done to handle
*/
type Menu struct {
	Command       CommandRun
	ManagerName   ManagerName
	ManagerOption ManagerOption
	FileName      string
	DirName       string
	Mode          string
	Debug         string
	Output        Output
	Interactive   bool
}

func NewMenu() *Menu {
	return new(Menu)
}
