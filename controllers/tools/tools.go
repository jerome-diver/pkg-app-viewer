package tool

import (
	"compress/gzip"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	model "github.com/pkg-app-viewer/models"
	view "github.com/pkg-app-viewer/views"
)

var config model.Config
var logging view.Logging

type Box struct {
}

type Gz struct {
	FileName    string
	Raw         *os.File
	ClearReader *gzip.Reader
	ClearBytes  []byte
}

func New() *Box {
	logging = view.GetLogger()
	t := new(Box)
	return t
}

func (t *Box) readGzip(filename string) *Gz {
	gz := new(Gz)
	var err error
	gz.FileName = filename
	logging.Debug("Start function readGzip", slog.String("fname", filename))
	gz.Raw, err = os.Open(filename) // get a io.Reader pointer
	logging.SetError(err)
	logging.CheckError("Can not read gzip file")
	defer gz.Raw.Close()

	gz.ClearReader, err = gzip.NewReader(gz.Raw) // get a gzip.Reader pointer
	logging.SetError(err)
	logging.CheckError("Can not create gzip Reader")
	defer gz.ClearReader.Close()

	gz.ClearBytes, err = io.ReadAll(gz.ClearReader) // get a byte slice
	logging.SetError(err)
	logging.CheckError("Can not get bytes slice from Reader")
	return gz
}

func (t *Box) printMeta(g *Gz) {
	logging.Info("Read file apt history gz log with Meta Header",
		slog.String("File", g.FileName),
		slog.Group("Meta",
			slog.String("Name", g.ClearReader.Name),
			slog.String("Extra", string(g.ClearReader.Extra)),
			slog.String("Comment", g.ClearReader.Comment),
			slog.String("ModTime", g.ClearReader.ModTime.String()),
			slog.String("OS", string(g.ClearReader.OS)),
		))
}

func (t *Box) GetFileContent(model_menu *model.Menu) []byte {
	var err error
	// return readable file content bytes (from clear fille or gz file)
	suffix := filepath.Ext(model_menu.FileName)
	if suffix == ".gz" {
		gz := t.readGzip(model_menu.FileName)
		if model_menu.ShowMeta {
			t.printMeta(gz)
		} else {
			logging.Info("Treated file apt history gz log", slog.String("File", model_menu.FileName))
		}
		return gz.ClearBytes
	} else {
		var clearBytes []byte
		clearBytes, err = os.ReadFile(model_menu.FileName)
		logging.SetError(err)
		logging.CheckError("Can not read file.")
		logging.Info("Treated file apt history log", slog.String("File", model_menu.FileName))
		return clearBytes
	}
}

func (t *Box) GetAptHistoryFilesList(model_menu *model.Menu) []string {
	// concatene and return files content bytes
	var aptHistoryFiles []string
	logging.Info("Treated file apt history log directory", slog.String("Dir", model_menu.DirName))
	err := filepath.Walk(model_menu.DirName, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			var matched bool
			matched, err = regexp.Match(`history\.log.*`, []byte(info.Name()))
			logging.SetError(err)
			logging.CheckError("Can not match history log froom bytes slice")
			if matched {
				fullFileName := model_menu.DirName + "/" + info.Name()
				aptHistoryFiles = slices.Insert(aptHistoryFiles, 0, fullFileName)
			}
		}
		return err
	})
	logging.SetError(err)
	logging.CheckError("Can not go through directory to get history files")
	return aptHistoryFiles
}
