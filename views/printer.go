package view

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/manifoldco/promptui"
	model "github.com/pkg-app-viewer/models"
	report "github.com/pkg-app-viewer/views/report"
)

var logger report.Logging

type Printer struct {
	Output    model.Output
	Packager  model.ManagerName
	SearchFor model.ManagerOption
	Data      []string
}

func NewPrinter(model *model.Menu) *Printer {
	logger = report.GetLogger()
	p := &Printer{
		Output:    model.Output,
		Packager:  model.ManagerName,
		SearchFor: model.ManagerOption}
	return p
}

func (p *Printer) Write(data []string) {
	p.Data = data
	switch p.Output.Mode {
	case "stdout":
		{
			for _, d := range p.Data {
				fmt.Println(d)
			}
		}
	case "file":
		{
			if p.Output.File == "" {
				currentTime := time.Now()
				p.Output.File = p.Packager.String() +
					"_" + p.SearchFor.String() +
					"_Installed_" + currentTime.Format("2006-01-02_15:04:05") +
					"." + p.Output.Format
			}
			_, err := os.Stat(p.Output.File)
			if err == nil {
				//Choose to overwrite or to append or to cancel
				prompt := promptui.Select{
					Label: "File exist, what do you want to do ?",
					Items: []string{"Overwrite", "Append", "Cancel"},
				}
				_, result, err := prompt.Run()
				logger.SetError(err)
				logger.CheckError("Prompt return bad statement")
				var options int
				switch result {
				case "Overwrite":
					{
						logger.SetError(os.Truncate(p.Output.File, 0))
						logger.CheckError("Truncate Output.File gone bad")
						options = os.O_RDWR | os.O_CREATE
					}
				case "Append":
					options = os.O_APPEND | os.O_CREATE | os.O_WRONLY
				case "Cancel":
					os.Exit(0)
				}
				f, err := os.OpenFile(p.Output.File, options, 0644)
				logger.SetError(err)
				logger.CheckError("Can not open Output.File")
				defer f.Close()
				p.writeInFile(f)
			} else if os.IsNotExist(err) {
				f, err := os.Create(p.Output.File)
				logger.SetError(err)
				logger.CheckError("OutFile doesn't existe and i can not create Output.File")
				defer f.Close()
				p.writeInFile(f)
			} else {
				logger.SetError(err)
				logger.CheckError("Output.File error")
			}
		}
	}
}

func (p *Printer) writeInFile(f *os.File) {
	var s_byte int
	for _, str := range p.Data {
		n, err := f.WriteString(str + "\n")
		logger.SetError(err)
		logger.CheckError("Can  t write in file")
		s_byte += n
	}
	logger.Debug("writed in file",
		slog.String("file", p.Output.File),
		slog.Int("size(byte)", s_byte))
}
