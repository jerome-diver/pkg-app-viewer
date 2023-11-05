package engine

import (
	"fmt"

	model "github.com/pkg-app-viewer/models"
	view "github.com/pkg-app-viewer/views"
)

/*
    I'm using a Strategy Design Pattern and Abstract Factory
	there to define to search the packages list
	depend of Manager Type Option as [System, Isolated, Language]
	and	depends on the option for each type
*/

var config model.Config
var logging view.Logging

type ReposFinder interface { // Abstract Finder for repos
	Find() map[string][]string
	AddPackage(origin, packageName string, computeOrigin func(string) bool)
	IsInstalled(packageName string) bool
}

type ReposHandler interface { // Abstract Handler for repos
	GetRepos() map[string][]string
	GetPackages() []string
}

type SystemManagerHandler interface { // Handle System Manager
	ManagerModel() model.Manager
}

type SystemReposAlgorithm interface { // algorithm to find System repos
	OriginSelector(string) bool
	Show()
}

type SystemRepos interface {
	SystemManagerHandler
	SystemReposAlgorithm
	ReposFinder
	ReposHandler
}

/*--------------------------------------------------------------------------------+
|  Struct inherit interface SystemRepos                                          |
|	  and contains Arch and Name                                                    |
|	  and is composed by Installed as PackagesOfRepo slice                          |
|	  and methods to search by crawling files                                       |
|	  and shared methods that will be used by specialized types through interface   |
+--------------------------------------------------------------------------------*/

type InstalledSystemRepos struct {
	//SystemReposAlgorithm
	Installed      []model.SystemRepo
	Arch           string
	Name           string
	hasBeenChecked []string
}

type SystemReposHandler struct { // embed ReposHandler
	InstalledSystemRepos
}

func (ir SystemReposHandler) GetRepos() map[string][]string {
	repos := map[string][]string{}
	for _, repo := range ir.Installed {
		repos[repo.Origin] = repo.Packages
	}
	return repos
}

func (ir SystemReposHandler) GetPackages() []string {
	packages := []string{}
	for _, repo := range ir.Installed {
		packages = append(packages, repo.Packages...)
	}
	return packages
}

/*
	Methods to be used for all the namespace engine
*/

func cleaningSystemRepos(repos []model.SystemRepo) []model.SystemRepo {
	// remove empty repos from the list
	// of Installed repos type PackagesOfRepos
	index_to_remove := []int{}
	for index, repo := range repos {
		if len(repo.Packages) == 0 {
			index_to_remove = append(index_to_remove, index)
		}
	}
	for i, index := range index_to_remove {
		repos = append(repos[:index-i], repos[index+1-i:]...)
	}
	return repos
}

/* ----------------------------------------------------------------------------------------
   B u i l d er
   Construct the content of packages defined and implemented
-----------------------------------------------------------------------------------------*/

func NewInstalledRepos(name, arch string, rf SystemReposAlgorithm) *InstalledSystemRepos {
	this := InstalledSystemRepos{
		SystemReposAlgorithm: rf,
		Installed:            []model.SystemRepo{},
		Name:                 name,
		Arch:                 arch,
		hasBeenChecked:       []string{},
	}
	initSystemRepos(&this)
	return &this
}

func NewDpkgRepos(system_option model.SystemOption) SystemRepos {
	var inherit SystemRepos
	switch system_option {
	case model.Foreign:
		inherit = new(ForeignDebRepos)
	}
	return inherit
}

func initSystemRepos(system_repos *InstalledSystemRepos) {
	repos := system_repos.Find()
	for k, values := range repos {
		for _, v := range values {
			system_repos.AddPackage(k, v, system_repos.OriginSelector)
		}
	}
	system_repos.Installed = cleaningSystemRepos(system_repos.Installed)
}

func Show(i *InstalledSystemRepos) {
	i.Show()
	for _, repo := range i.Installed {
		fmt.Printf("Origin: %s\n", repo.Origin)
		fmt.Printf("Packages: %s\n\n", repo.Packages)
	}
}

func main() {
	foreign_repos := NewInstalledRepos("Foreign", "amd64", NewDpkgRepos(model.Foreign))
	repos := foreign_repos.GetPackages()
	for _, r := range repos {
		fmt.Println(r)
	}
	Show(foreign_repos)
}
