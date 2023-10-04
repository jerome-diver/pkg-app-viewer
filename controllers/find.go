package controller

import (
	"bufio"
	"bytes"
	"log/slog"
	"regexp"
	"slices"

	model "github.com/pkg-app-viewer/models"
)

type Find struct {
	tool     *Tool
	Packages []string
}

func Finder(tool *Tool) *Find {
	f := new(Find)
	f.tool = tool
	return f
}

func (f *Find) cleanBytes(ipl []byte) []byte {
	// find "-y ", "-o <word>=\d", "-<-word>+"
	re_parser := regexp.MustCompile(`.*(\-y\s)|(\-o\s.*\=\d\s)*(\-(\-\w*)+\s)*(.*)`)
	clean := re_parser.ReplaceAllString(string(ipl), "$5")
	f.tool.logger.Debug("cleaned ipl", slog.String("ipl", clean))
	return []byte(clean)
}

func (f *Find) installOccurenceFound(ipl []byte, comp func(string) bool) {
	//remove unexpected content words
	ipl_clean := f.cleanBytes(ipl)
	//and split words
	splited_list := bytes.SplitAfter(ipl_clean, []byte(" "))
	//to add to packages bytes if not already present
	for _, w := range splited_list {
		clean_w := bytes.Trim(w, " ")
		sclean_w := string(clean_w)
		if slices.Contains(f.Packages, sclean_w) {
			f.tool.logger.Warn("already listed in packages []string", slog.String("name", sclean_w))
			continue
		}
		if len(clean_w) == 0 {
			continue
		}
		if comp(sclean_w) {
			f.Packages = append(f.Packages, sclean_w)
		}
		f.tool.logger.Info("added to packages []string", slog.String("name", sclean_w))
	}
}

func (f *Find) removedOccurenceFound(ipl []byte) {
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
			f.tool.logger.Warn("found occurence to delete in packages []string", slog.String("name", string(clean_w)))
			last := find_index + 1
			f.Packages = slices.Delete(f.Packages, find_index, last)
		}
	}
}

/*
  - Mode can be: [all, added, officialAdded, otherRepos, manual]

    all:
    |	apt install - apt remove
    |	occurence inside histoy.log files (gz included)

    added:
    |	apt install - apt remove
    |	occurence inside histoy.log files (gz included)
    |	followed by line contains "Requested-by:"

    officialAdded:
    |

    otherRepos:
    |

    manual:
    |	apt install - apt remove
    |	occurence inside histoy.log files (gz included)
    |	but package name should match for a "".deb" file
    |	to rich this, a "func(string)bool" is passed through params

*
*/
func (f *Find) AptInstalledFromHistory(rawHistory []byte, mode model.Search) {
	var cmp func(string) bool
	switch mode {
	case model.All:
		cmp = func(p string) bool { return true }
	case model.Added:
		cmp = func(p string) bool { return true }
	case model.FileSource:
		cmp = func(p string) bool {
			re := regexp.MustCompile(`.*\.deb$`)
			ok := re.MatchString(p)
			return ok
		}
	}
	f.tool.logger.Info("Start function AptInstalledFromHistory", slog.Int("raw", len(rawHistory)))
	scanner := bufio.NewScanner(bytes.NewReader(rawHistory))
	const maxCapacity = 512 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	scanner.Split(bufio.ScanLines)
	var ipl []byte
	for scanner.Scan() {
		b := scanner.Bytes()
		if mode == model.Added {
			if _, _, found := bytes.Cut(b, []byte("Requested-by: ")); found {
				f.installOccurenceFound(ipl, cmp)
			}
		}
		// detect, cut and separate apt install list packages
		if _, ipl, found := bytes.Cut(b, []byte("apt-get install ")); found {
			f.tool.logger.Debug("find apt-get install occurence",
				slog.String("CutAfter", string(ipl)))
			if mode == model.All || mode == model.FileSource {
				f.installOccurenceFound(ipl, cmp)
			}
			continue // next scan if treated there as "apt install" line detected
		}
		if _, ipl, found := bytes.Cut(b, []byte("apt install ")); found {
			f.tool.logger.Debug("find apt install occurence",
				slog.String("CutAfter", string(ipl)))
			f.installOccurenceFound(ipl, cmp)
			continue // next scan if treated there as "apt install" line detected
		}
		// detect, cut and separate apt remove list packages
		if _, removed_potential_list, found := bytes.Cut(b, []byte("apt remove ")); found {
			f.tool.logger.Warn("find apt remove occurence",
				slog.String("CutAfter", string(removed_potential_list)))
			f.removedOccurenceFound(removed_potential_list)
			continue // next scan if treated there as "apt remove" line detected
		}

	}
	if scanner.Err() != nil {
		f.tool.logger.Error("SCANNER ERROR", slog.String("msg", scanner.Err().Error()))
	}
}
