package controller_test

import (
	"testing"

	. "github.com/pkg-app-viewer/controllers"
	model "github.com/pkg-app-viewer/models"
	view "github.com/pkg-app-viewer/views"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Find_cleanBytes(t *testing.T) {
	to_test := []map[string]string{
		{
			"title":   "Have to find and remove '-y ' and all before",
			"to_test": "text to delete -y ok that's it",
			"want":    "ok that's it",
		},
		{
			"title":   "Have to find and remove '-o <word>=<digit> '",
			"to_test": "text to delete -y -o word=1 ok that's it",
			"want":    "ok that's it",
		},
		{
			"title":   "Have to find and remove any '-(-<word>)+ '",
			"to_test": "text to delete -y --word-many-time --other-one-again ok that's it",
			"want":    "ok that's it",
		},
	}
	l := view.NewLogger()
	f := Finder(l)
	for _, data := range to_test {
		Convey(data["title"], t, func() {
			So(ExportCleanBytes(f, []byte(data["to_test"])),
				ShouldEqual,
				[]byte(data["want"]))
		})
	}
}

func Test_Find_installOccurenceFound(t *testing.T) {
	type TestContent struct {
		title   string       // Test title
		to_test string       // string to test
		search  model.Search // comparator algorythm model [All, Added,OfficialRepos, OtherRepos, FileSource]
		want    []string     // expected to get back
		clean   bool         // initialize Find.Packages ?
	}
	to_test := []TestContent{
		{
			title:   "Have to clean and split to add to this.Packages",
			to_test: "text to delete -y one two three four",
			search:  model.All,
			want:    []string{"one", "two", "three", "four"},
			clean:   true,
		},
		{
			title:   "Have to not add if already exist'",
			to_test: "text to delete -y -o word=1 one five one one",
			search:  model.All,
			want:    []string{"one", "two", "three", "four", "five"},
			clean:   false,
		},
		{
			title:   "Have to add deb file only'",
			to_test: "text to delete -y -o word=1 ./test.deb one one file.deb",
			search:  model.FileSource,
			want:    []string{"./test.deb", "file.deb"},
			clean:   true,
		},
	}
	l := view.NewLogger()
	f := Finder(l)
	for _, data := range to_test {
		Convey(data.title, t, func() {
			if data.clean { // initialize list of Packages found
				f.Packages = []string{}
			}
			ExportInstallOccurenceFound(f,
				[]byte(data.to_test),
				data.search.Algorythm())
			So(f.Packages, ShouldEqual, data.want)
		})
	}
}
