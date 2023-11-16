package engine_test

import (
	"os"
	"testing"

	. "github.com/pkg-app-viewer/controllers/engine"
	model "github.com/pkg-app-viewer/models"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Find_cleanBytes(t *testing.T) {
	to_test := []map[string]string{
		{
			"title":   "Have to find and remove '-y ' and all before\n",
			"to_test": "text to delete -y ok that's it",
			"want":    "ok that's it",
		},
		{
			"title":   "Have to find and remove '-o <word>=<digit> '\n",
			"to_test": "text to delete -y -o word=1 ok that's it",
			"want":    "ok that's it",
		},
		{
			"title":   "Have to find and remove any '-(-<word>)+ '\n",
			"to_test": "text to delete -y --word-many-time --other-one-again ok that's it",
			"want":    "ok that's it",
		},
		{
			"title":   "Have to let go simple word that doesn't match\n",
			"to_test": "simple word",
			"want":    "simple word",
		},
	}
	f := NewAptHistory()
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
		title   string              // Test title
		to_test string              // string to test
		search  model.ManagerOption // comparator algorythm model [All, Added,OfficialRepos, OtherRepos, FileSource]
		want    []string            // expected to get back
		clean   bool                // initialize Find.Packages ?
	}
	to_test := []TestContent{
		{
			title:   "Have to clean and split to add to this.Packages\n",
			to_test: "text to delete -y one two three four",
			search:  model.All,
			want:    []string{"one", "two", "three", "four"},
			clean:   true,
		},
		{
			title:   "Have to not add if already exist\n",
			to_test: "text to delete -y -o word=1 one five one one",
			search:  model.All,
			want:    []string{"one", "two", "three", "four", "five"},
			clean:   false,
		},
		{
			title:   "Have to add deb file only\n",
			to_test: "text to delete -y -o word=1 ./test.deb one one file.deb",
			search:  model.FileSource,
			want:    []string{"./test.deb", "file.deb"},
			clean:   true,
		},
	}
	f := NewAptHistory()
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

func Test_Find_removedOccurenceFound(t *testing.T) {
	type TestContent struct {
		title        string   // Test title
		initPackages []string // initialize Find.Packages state
		to_test      string   // string to test
		want         []string // expected to get back
	}
	to_test := []TestContent{
		{
			title:        "Have to clean and split to remove found occurence from this.Pakages\n",
			initPackages: []string{"one", "two", "./test.deb"},
			to_test:      "text to delete -y one three four ./test.deb",
			want:         []string{"two"},
		},
	}
	f := NewAptHistory()
	for _, data := range to_test {
		Convey(data.title, t, func() {
			f.Packages = data.initPackages
			ExportRemovedOccurenceFound(f, []byte(data.to_test))
			So(f.Packages, ShouldEqual, data.want)
		})
	}
}

func Test_Find_DebianPackagesToSearchFor(t *testing.T) {
	pwd_dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	type TestContent struct {
		title    string
		search   model.ManagerOption
		fileName string
		want     []string
	}
	to_test := []TestContent{
		{
			title:    "Extract Debian packages with mode All from mocked log file\n",
			search:   model.All,
			fileName: pwd_dir + "/../models/mock_var_log_apt_history.log",
			want: []string{
				"anydesk", "siftool", "grub-common",
				"grub2-common", "grub-pc", "btrfs-progs",
				"e2fsprogs", "ntfs-3g", "system76-acpi-dkms",
				"system76-dkms", "system76-io-dkms", "amd-ppt-bin",
				"nvidia-driver-515", "./openrgb_0.9_amd64_bookworm_b5f46e3.deb",
			},
		},
		{
			title:    "Extract Debian packages with mode FileSource from mocked log file\n",
			search:   model.FileSource,
			fileName: pwd_dir + "/../models/mock_var_log_apt_history.log",
			want: []string{
				"./openrgb_0.9_amd64_bookworm_b5f46e3.deb",
			},
		},
		{
			title:    "Extract Debian packages with mode Added from mocked log file\n",
			search:   model.User,
			fileName: pwd_dir + "/../models/mock_var_log_apt_history.log",
			want: []string{
				"anydesk", "siftool", "./openrgb_0.9_amd64_bookworm_b5f46e3.deb",
			},
		},
		{
			title:    "Extract Debian packages with mode OfficialAdded from mocked log file\n",
			search:   model.All,
			fileName: pwd_dir + "/../models/mock_var_log_apt_history.log",
			want: []string{
				"grub-common", "grub2-common", "grub-pc", "btrfs-progs",
				"e2fsprogs", "ntfs-3g", "system76-acpi-dkms",
				"system76-dkms", "system76-io-dkms", "amd-ppt-bin",
				"nvidia-driver-515",
			},
		},
	}
	f := NewAptHistory()
	for _, test := range to_test {
		rawHistory, _ := os.ReadFile(test.fileName)
		Convey(test.title, t, func() {
			f.Packages = []string{}
			f.AptPackagesToSearchFor(rawHistory, test.search)
			So(f.Packages, ShouldEqual, test.want)
		})
	}
}
