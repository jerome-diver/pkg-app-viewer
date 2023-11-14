package engine

import (
	"fmt"

	model "github.com/pkg-app-viewer/models"
	report "github.com/pkg-app-viewer/views/report"
)

/*
    I'm using a Strategy Design Pattern and Abstract Factory
	there to define to search the packages list
	depend of Manager Type Option as [System, Isolated, Language]
	and	depends on the option for each type
*/

var config model.Config
var logging report.Logging

type ReposFinder interface { // Abstract Finder for repos
	Find() map[string][]string
	AddPackage(origin, packageName string, computeOrigin func(string) bool)
	IsInstalled(packageName string) bool
	UserInstalled(userName string, packageName string) bool
}

type ReposHandler interface { // Abstract Handler for repos
	//	GetRepos() map[string][]string
	//	GetPackages() []string
	GetInstalled() []model.Repository
	SetReposHandler(ReposHandler)
	UpdateInstalled([]model.Repository)
}

type ManagerHandler interface { // Handle System Manager
	ManagerModel() model.Manager
}

type SystemReposAlgorithm interface { // algorithm to find System repos
	OriginSelector(string) bool
	Show()
	Option() model.SystemOption
}

/*--------------------------------------------------------------------------------+
|  Struct inherit interface SystemRepos                                          |
|	  and contains Arch and Name                                                    |
|	  and is composed by Installed as PackagesOfRepo slice                          |
|	  and methods to search by crawling files                                       |
|	  and shared methods that will be used by specialized types through interface   |
+--------------------------------------------------------------------------------*/

type SystemRepos interface {
	ManagerHandler
	SystemReposAlgorithm
	ReposFinder
	ReposHandler
}

/*
Any inherited SystemRepos will have a
ReposHandler interface as component
(it will handle it indirectly by its component)
this can be:
	InstalledSystemRepos
	InstalledLanguageepos
	InstalledIsolatedRepos
*/

type InstalledLanguageRepos struct {
	Installed []model.Repository // ReposHandler must have it
}

type InstalledIsolatedRepos struct {
	Installed []model.Repository // ReposHandler must have it
}

type InstalledSystemRepos struct {
	Installed      []model.Repository // ReposHandler must have it
	Arch           string
	Name           string
	hasBeenChecked []string
	userInstalled  map[string][]string
}

// Handle ReposHandler interface
func (ir InstalledSystemRepos) GetInstalled() []model.Repository {
	return ir.Installed
}

func (ir InstalledSystemRepos) SetReposHandler(installed ReposHandler) {
	ir = installed.(InstalledSystemRepos)
}

func (ir InstalledSystemRepos) UpdateInstalled(installed []model.Repository) {
	ir.Installed = installed
}

//func (i *InstalledSystemRepos) Show() {
//	for _, repo := range i.Installed {
//		fmt.Printf("Origin: %s\n", repo.Origin)
//		fmt.Printf("Packages: %s\n\n", repo.Packages)
//	}
//}

/*
	Methods to be used for all the namespace engine
*/

func cleaningSystemRepos(repos []model.Repository) []model.Repository {
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

func GetRepos(repos_hanlder ReposHandler) map[string][]string {
	repos := map[string][]string{}
	for _, repo := range repos_hanlder.GetInstalled() {
		repos[repo.Origin] = repo.Packages
	}
	return repos
}

func GetPackages(repos_handler ReposHandler) []string {
	packages := []string{}
	for _, repo := range repos_handler.GetInstalled() {
		packages = append(packages, repo.Packages...)
	}
	return packages
}

/* ----------------------------------------------------------------------------------------
   B u i l d er
   Construct the content of packages defined and implemented
-----------------------------------------------------------------------------------------*/

func NewSystemRepos(name, arch string, option model.SystemOption) SystemRepos {
	var ob SystemRepos
	repos_handler := InstalledSystemRepos{
		Installed:      []model.Repository{},
		Name:           name,
		Arch:           arch,
		hasBeenChecked: []string{},
	}
	switch name {
	case "debian":
		ob = NewDpkgRepos(option, repos_handler)
	}
	initSystemRepos(ob)
	return ob
}

func NewDpkgRepos(system_option model.SystemOption, ir InstalledSystemRepos) SystemRepos {
	var inherit SystemRepos
	switch system_option {
	case model.Foreign:
		inherit = ForeignDebRepos{}
	}
	inherit.SetReposHandler(ir)
	return inherit
}

func initSystemRepos(system_repos SystemRepos) {
	repos := system_repos.Find()
	for k, values := range repos {
		for _, v := range values {
			system_repos.AddPackage(k, v, system_repos.OriginSelector)
		}
	}
	installed := system_repos.GetInstalled()
	clean_repositories := cleaningSystemRepos(installed)
	system_repos.UpdateInstalled(clean_repositories)
}

func main() {
	foreign_repos := NewSystemRepos("debian", "amd64", model.Foreign)
	repos := GetPackages(foreign_repos)
	for _, r := range repos {
		fmt.Println(r)
	}
	foreign_repos.Show()
}
