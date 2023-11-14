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
	report "github.com/pkg-app-viewer/views/report"
)

var logging report.Logging

type Box struct {
}

type Gz struct {
	FileName    string
	Raw         *os.File
	ClearReader *gzip.Reader
	ClearBytes  []byte
}

func New() *Box {
	logging = report.GetLogger()
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

func (t *Box) GetFileContent(fileName string) []byte {
	var err error
	// return readable file content bytes (from clear fille or gz file)
	suffix := filepath.Ext(fileName)
	if suffix == ".gz" {
		gz := t.readGzip(fileName)
		return gz.ClearBytes
	} else {
		var clearBytes []byte
		clearBytes, err = os.ReadFile(fileName)
		logging.SetError(err)
		logging.CheckError("Can not read file.")
		logging.Info("Treated file apt history log", slog.String("File", fileName))
		return clearBytes
	}
}

func (t *Box) GetAptHistoryFilesList(dirName string) []string {
	// concatene and return files content bytes
	var aptHistoryFiles []string
	logging.Info("Treated file apt history log directory", slog.String("Dir", dirName))
	err := filepath.Walk(dirName, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			var matched bool
			matched, err = regexp.Match(`history\.log.*`, []byte(info.Name()))
			logging.SetError(err)
			logging.CheckError("Can not match history log froom bytes slice")
			if matched {
				fullFileName := dirName + "/" + info.Name()
				aptHistoryFiles = slices.Insert(aptHistoryFiles, 0, fullFileName)
			}
		}
		return err
	})
	logging.SetError(err)
	logging.CheckError("Can not go through directory to get history files")
	return aptHistoryFiles
}

func UpdateConfig(config *model.Config, identity *model.Identity) {

}
