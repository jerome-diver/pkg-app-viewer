package engine

import (
	"bufio"
	"bytes"
	"log/slog"
	"regexp"
	"slices"

	model "github.com/pkg-app-viewer/models"
)

type AptHistory struct {
	Packages []string
}

func NewAptHistory() *AptHistory {
	f := new(AptHistory)
	return f
}

func (f *AptHistory) cleanBytes(ipl []byte) []byte {
	_, after, found := bytes.Cut(ipl, []byte("-y "))
	var first_pass []byte
	if found {
		first_pass = after
	} else {
		first_pass = ipl
	}
	re_parser := regexp.MustCompile(`(\-o [\w\:]*\=[\d\/\w]*)+\s`)
	second_pass := re_parser.ReplaceAllString(string(first_pass), "")
	re_parser = regexp.MustCompile(`(\-(\-\w+)+)+\s`)
	clean := re_parser.ReplaceAllString(string(second_pass), "")
	logging.Debug("cleaned ipl", slog.String("ipl", clean))
	return []byte(clean)
}

func (f *AptHistory) installOccurenceFound(ipl []byte, algorythm func(string) bool) {
	// remove unexpected content words
	ipl_clean := f.cleanBytes(ipl)
	// and split words
	splited_list := bytes.SplitAfter(ipl_clean, []byte(" "))
	// to add to this Packages string slice if missing
	for _, w := range splited_list {
		clean_w := bytes.Trim(w, " ")
		sclean_w := string(clean_w)
		if slices.Contains(f.Packages, sclean_w) {
			logging.Debug("already listed in packages []string", slog.String("name", sclean_w))
			continue
		}
		if len(clean_w) == 0 {
			continue
		}
		if algorythm(sclean_w) {
			f.Packages = append(f.Packages, sclean_w)
			logging.Debug("added to packages []string", slog.String("name", sclean_w))
		}
	}
}

func (f *AptHistory) removedOccurenceFound(ipl []byte) {
	//remove unexpected content words
	ipl_clean := f.cleanBytes(ipl)
	//and split words
	splited_list := bytes.SplitAfter(ipl_clean, []byte(" "))
	//to delete to packages each slices of bytes if present
	for _, w := range splited_list {
		clean_w := bytes.Trim(w, " ")
		if len(clean_w) == 0 { // no empty slices bytes word
			continue
		}
		find_index := slices.Index(f.Packages, string(clean_w))
		if find_index != -1 {
			logging.Debug("found occurence to delete in packages []string", slog.String("name", string(clean_w)))
			last := find_index + 1
			f.Packages = slices.Delete(f.Packages, find_index, last)
		}
	}
}

/*
  - Mode can be: [All, Added, OfficialRepos, OtherRepos, FileSource]
    Methods are:
    All:
    |	apt install - apt remove (or apt-get)
    |	occurence inside history.log files (gz included)

    Added:  (all added bu user)
    |	apt install - apt remove (or apt-get)
    |	occurence inside histoy.log files (gz included)
    |	followed by line contains "Requested-By: "

    OfficialAdded: (all official added by user)
    |   substract OtherRepos result from Added one modes to rich your content

    OtherRepos: (all non official added by user)
    |   find repo.sources.list from /etc/apt/sources.list.d/
    |   search content of each repo from /var/lib/atp/lists/repo_pkg_nameè_Packages files
    |   should find indication after "Package: " tagged line
    |   AND ????

    FileSource: (added by user from deb file)
    |	apt install - apt remove (or apt-get)
    |	occurence inside history.log files (gz included)
    |	but package name should match for a ".deb" file
    |	to rich this, a model.Search type indication is linked

*
*/

func (f *AptHistory) AptPackagesToSearchFor(rawHistory []byte, mode model.SystemOption, args ...string) {
	// args could be a name of User
	logging.Debug("Start (*AptHistory).AptInstalledFromHistory", slog.Int("raw", len(rawHistory)))
	scanner := bufio.NewScanner(bytes.NewReader(rawHistory))
	const maxCapacity = 256 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	scanner.Split(bufio.ScanLines)
	installed_words := []string{
		"apt install ", "apt-get install ",
	}
	removed_words := []string{
		"apt remove ", "apt-get remove ",
	}
	var found bool
	var after []byte // need to be outside because of User management detection
	var to_install bool
	var to_remove bool
scanner:
	for scanner.Scan() {
		b := scanner.Bytes()
		if mode == model.User { // find User managed package
			if _, user, found := bytes.Cut(b, []byte("Requested-By: "+args[0])); found {
				logging.Debug("find apt 'Requested-by: ' occurence (user management)",
					slog.String("user", string(user)),
					slog.String("to Add (previous variable)", string(after)))
				if to_install {
					f.installOccurenceFound(after, mode.Algorythm())
				}
				if to_remove {
					f.removedOccurenceFound(after)
				}
				continue // next scan
			}
		}
		// detect, cut and separate apt install list packages
		for _, install := range installed_words {
			to_install = false
			if _, after, found = bytes.Cut(b, []byte(install)); found {
				to_install = true
				logging.Debug("find apt-get install occurence",
					slog.String("CutAfter", string(after)))
				if mode == model.All || mode == model.FileSource {
					f.installOccurenceFound(after, mode.Algorythm())
				}
				continue scanner // next scan if treated
			}
		}
		// detect, cut and separate apt remove list packages
		for _, remove := range removed_words {
			to_remove = false
			if _, after, found = bytes.Cut(b, []byte(remove)); found {
				to_remove = true
				logging.Debug("find apt remove occurence",
					slog.String("CutAfter", string(after)))
				f.removedOccurenceFound(after)
				continue scanner // next scan if treated
			}
		}
	}
	logging.SetError(scanner.Err())
	logging.CheckError("Scanner error at history time")
}
