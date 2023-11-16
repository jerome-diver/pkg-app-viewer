package engine

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	tool "github.com/pkg-app-viewer/controllers/tools"
	model "github.com/pkg-app-viewer/models"
)

/*
DpkgRepos handle SystemRepos interface
He has an InstalledSystemRepos as ReposHanlder
*/
type DpkgRepos struct {
	InstalledSystemRepos
}

// Handle ManagerHandler interface
func (d DpkgRepos) Model() model.ManagerName {
	return model.Dpkg
}

// Handle ReposFinder interface
func (ir DpkgRepos) Find() map[string][]string {
	// will find repo packages for each repository through apt repo files crawl
	// First here is the glob *Release files with Origin: flag (has repo name and key)
	// send back repo data without cleaning uniq entries and no specific origin
	repos := map[string][]string{}
	aptReposDir := "/var/lib/apt/lists"
	releaseFiles, err := filepath.Glob(filepath.Join(aptReposDir, "*Release"))
	if err != nil {
		fmt.Println("Error to get Release's files:", err)
		return nil
	}
	for _, releaseFile := range releaseFiles {
		f, err := os.Open(releaseFile)
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}
		defer f.Close()
		//fmt.Printf("Reading file: %s\n", releaseFile)
		scanner := bufio.NewScanner(f)
		var currentOrigin string
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Origin:") {
				currentOrigin = strings.TrimSpace(strings.TrimPrefix(line, "Origin:"))
				break
			}
		}
		//fmt.Printf("\tFound origin: %s\n", currentOrigin)
		repos[currentOrigin] = ir.findPackages(currentOrigin, releaseFile)
	}
	return repos
}

func (ir DpkgRepos) findPackages(origin string, rFile string) []string {
	// from release file name, and architecture, will find the glob *Packages files
	// and Package: tag content to return repo packages list (not uniq) of the repo
	packages := []string{}
	pFile := strings.TrimSuffix(rFile, "Release")
	pFile = strings.TrimSuffix(pFile, "In")
	//fmt.Printf("\tSearching for files Glob %s*Packages\n", pFile)
	packagesFiles, err := filepath.Glob(pFile + "*" + ir.Arch + "_Packages")
	if err != nil {
		fmt.Println("Error matching Packages files:", err)
		return nil
	}
	for _, packagesFile := range packagesFiles {
		f, err := os.Open(packagesFile)
		if err != nil {
			fmt.Println("Error opening Packages file:", err)
			continue
		}
		defer f.Close()
		//fmt.Printf("\t\tReading file: %s\n", packagesFile)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Package: ") {
				packageName := strings.TrimSpace(strings.TrimPrefix(line, "Package: "))
				packages = append(packages, packageName)
			}
		}
	}
	return packages
}

func (ir DpkgRepos) AddPackage(origin, packageName string, computeOrigin func(string) bool) {
	// Add package one by one with its origin.
	// will select data depend to the computation code of the specialsed type handled
	//fmt.Printf("\t\t\t\tTry to add package: %s at Origin: %s for lentgh Foreign: %d\n", packageName, origin, len(ir.Foreign))
	if computeOrigin(origin) {
		// Add the package in his place after to check unicity list
		if !slices.ContainsFunc(ir.Installed, func(r model.Repository) bool {
			return r.Origin == origin
		}) {
			ir.Installed = append(ir.Installed, model.Repository{
				Origin:   origin,
				Packages: []string{}})
		}
		for i, repo := range ir.Installed {
			if repo.Origin == origin {
				//fmt.Println("\t\t\t\tFind repo.Origin corresponding.")
				// Can only add uniq package and
				// the one that is present in /var/lib/dpkg/status file
				if !slices.Contains(repo.Packages, packageName) { // uniq...
					if ir.IsInstalled(packageName) {
						repo.Packages = append(repo.Packages, packageName)
						ir.Installed[i] = repo
					}
				}
			}
		}
	}
}

func (ir DpkgRepos) IsInstalled(packageName string) bool {
	// check if package is installed
	if !slices.Contains(ir.hasBeenChecked, packageName) {
		ir.hasBeenChecked = append(ir.hasBeenChecked, packageName)
		file, err := os.Open("/var/lib/dpkg/status")
		if err != nil {
			fmt.Println("Error opening dpkg status file:", err)
			return false
		}
		defer file.Close()
		//fmt.Printf("\t\tReading file: %s\n", packagesFile)
		exist := false
		installed := false
		scanner := bufio.NewScanner(file)
		max := 256 * 1024
		buf := make([]byte, max)
		scanner.Buffer(buf, max)
		for scanner.Scan() {
			line := scanner.Text()
			switch line {
			case "":
				if exist {
					break
				}
			case "Package: " + packageName:
				exist = true
				continue
			case "Status: install ok installed":
				installed = true
				break
			}
			//fmt.Printf("Found \"%s\"\n", packageName)
		}
		return exist && installed
	}
	return false
}

func (ir DpkgRepos) UserInstalled(userName string, packageName string) bool {
	// check if packageName has been installed by userName
	if ir.userInstalled[userName] == nil {
		ir.populateUserInstalled(userName)
	}
	pkg_list := ir.userInstalled[userName]
	return slices.Contains(pkg_list, packageName)
}

func (ir DpkgRepos) populateUserInstalled(userName string) {
	// Have to use Apt search for user mode
	logging.Debug("RUN Dpkg.UserInstalled from engine",
		slog.String("userName", userName))
	aptHistory := NewAptHistory()
	toolBox := tool.New()
	aptDirName := "/var/log/apt"
	filesList := toolBox.GetAptHistoryFilesList(aptDirName)
	for _, file := range filesList {
		clearBytes := toolBox.GetFileContent(file)
		aptHistory.AptPackagesToSearchFor(clearBytes, model.User, userName)
	}
	ir.userInstalled[userName] = aptHistory.Packages
}

/* ----------------------------------------------------------------------------------------
   PopOSRepos handles SystemReposAlgorithm
   interface and is composed by
   DpkgRepos as InstalledSystemRepos
-----------------------------------------------------------------------------------------*/

type PopOSRepos struct { // Origin is: pop-os*
	DpkgRepos
}

func (repos PopOSRepos) OriginSelector(origin string) bool {
	return (strings.Contains(origin, "pop-os") || strings.Contains(origin, "system76"))
}

func (repos PopOSRepos) Show() {
	fmt.Println("+------------------------------+")
	fmt.Println("|          P o P   O S         |")
	fmt.Println("+------------------------------+")
}

func (repos PopOSRepos) Option() model.ManagerOption {
	return model.Distribution
}

/* ----------------------------------------------------------------------------------------
   DebianRepos handles SystemReposAlgorithm
   interface and is composed by
   DpkgRepos as InstalledSystemRepos
-----------------------------------------------------------------------------------------*/

type DebianRepos struct { // Origin is: Debian
	DpkgRepos
}

func (repos DebianRepos) OriginSelector(origin string) bool {
	return strings.Contains(origin, "Debian")
}

func (repos DebianRepos) Show() {
	fmt.Println("+------------------------------+")
	fmt.Println("|          D e b i a n         |")
	fmt.Println("+------------------------------+")
}

func (repos DebianRepos) Option() model.ManagerOption {
	return model.System
}

/* ----------------------------------------------------------------------------------------
   UbuntuRepos handles SystemReposAlgorithm
   interface and is composed by
   DpkgRepos as InstalledSystemRepos
-----------------------------------------------------------------------------------------*/

type UbuntuRepos struct { // Origin is: Ubuntu
	DpkgRepos
}

func (repos UbuntuRepos) OriginSelector(origin string) bool {
	return strings.Contains(origin, "Ubuntu")
}

func (repos UbuntuRepos) Show() {
	fmt.Println("+------------------------------+")
	fmt.Println("|          U b u n t u         |")
	fmt.Println("+------------------------------+")
}

func (repos UbuntuRepos) Option() model.ManagerOption {
	return model.System
}

/* ----------------------------------------------------------------------------------------
   ForeignRepos handle SystemReposAlgorithm
   interface and is composed by
   DpkgRepos as InstalledSystemRepos
-----------------------------------------------------------------------------------------*/

type ForeignDebRepos struct { // Origin is not any: [ Ubuntu*, pop-os*, system*]
	DpkgRepos
}

func (repos ForeignDebRepos) OriginSelector(origin string) bool {
	return !(strings.Contains(origin, "Ubuntu") ||
		strings.Contains(origin, "pop-os") ||
		strings.Contains(origin, "system76"))
}

func (repos ForeignDebRepos) Show() {
	fmt.Println("+---------------------------------+")
	fmt.Println("|          F o r e i g n          |")
	fmt.Println("+---------------------------------+")
}

func (repos ForeignDebRepos) Option() model.ManagerOption {
	return model.Foreign
}
