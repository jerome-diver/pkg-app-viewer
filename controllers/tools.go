package controller

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
)

type Tool struct {
	Error  error
	logger *slog.Logger
}

type Gz struct {
	FileName    string
	Raw         *os.File
	ClearReader *gzip.Reader
	ClearBytes  []byte
}

func ToolBox(logger *slog.Logger) *Tool {
	t := new(Tool)
	t.logger = logger
	return t
}

func (t *Tool) checkError() {
	if t.Error != nil {
		t.logger.Error(t.Error.Error())
	}
}

func (t *Tool) readGzip(filename string) *Gz {
	gz := new(Gz)
	gz.FileName = filename
	t.logger.Debug("Start function readGzip", slog.String("fname", filename))
	gz.Raw, t.Error = os.Open(filename) // get a io.Reader pointer
	t.checkError()
	defer gz.Raw.Close()

	gz.ClearReader, t.Error = gzip.NewReader(gz.Raw) // get a gzip.Reader pointer
	t.checkError()
	defer gz.ClearReader.Close()

	gz.ClearBytes, t.Error = io.ReadAll(gz.ClearReader) // get a byte slice
	t.checkError()
	return gz
}

func (t *Tool) printMeta(g *Gz) {
	t.logger.Info("Read file apt history gz log with Meta Header",
		slog.String("File", g.FileName),
		slog.Group("Meta",
			slog.String("Name", g.ClearReader.Name),
			slog.String("Extra", string(g.ClearReader.Extra)),
			slog.String("Comment", g.ClearReader.Comment),
			slog.String("ModTime", g.ClearReader.ModTime.String()),
			slog.String("OS", string(g.ClearReader.OS)),
		))
}

func (t *Tool) GetFileContent(model_menu *model.Menu) []byte {
	// return readable file content bytes (from clear fille or gz file)
	suffix := filepath.Ext(model_menu.FileName)
	if suffix == ".gz" {
		gz := t.readGzip(model_menu.FileName)
		if model_menu.ShowMeta {
			t.printMeta(gz)
		} else {
			t.logger.Info("Treated file apt history gz log", slog.String("File", model_menu.FileName))
		}
		return gz.ClearBytes
	} else {
		var clearBytes []byte
		clearBytes, t.Error = os.ReadFile(model_menu.FileName)
		t.checkError()
		t.logger.Info("Treated file apt history log", slog.String("File", model_menu.FileName))
		return clearBytes
	}
}

func (t *Tool) GetAptHistoryFilesList(model_menu *model.Menu) []string {
	// concatene and return files content bytes
	var aptHistoryFiles []string
	t.logger.Info("Treated file apt history log directory", slog.String("Dir", model_menu.DirName))
	t.Error = filepath.Walk(model_menu.DirName, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			var matched bool
			matched, t.Error = regexp.Match(`history\.log.*`, []byte(info.Name()))
			t.checkError()
			if matched {
				fullFileName := model_menu.DirName + "/" + info.Name()
				aptHistoryFiles = slices.Insert(aptHistoryFiles, 0, fullFileName)
			}
		}
		return t.Error
	})
	t.checkError()
	return aptHistoryFiles
}
