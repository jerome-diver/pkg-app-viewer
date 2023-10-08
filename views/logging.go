package view

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

type Logging struct { // proxy to add methd checkError
	Error   error
	Log     *slog.Logger
	Handler *PrettyHandler
	Level   *slog.LevelVar
	Opt     *slog.HandlerOptions
}

func (l *Logging) CheckError(msg string) {
	if l.Error != nil {
		l.Log.Error(msg, slog.String("err", l.Error.Error()))
	}
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

func NewPrettyHandler(
	out io.Writer,
	opts *slog.HandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts),
		l:       log.New(out, "", 0),
	}

	return h
}

func NewLogger() *Logging {
	loggingLevel := new(slog.LevelVar)
	opts := slog.HandlerOptions{
		Level: loggingLevel,
	}
	textHandler := NewPrettyHandler(os.Stdout, &opts)
	log := &Logging{
		Log:     slog.New(textHandler),
		Handler: textHandler,
		Level:   loggingLevel,
		Opt:     &opts,
	}
	return log
}

func (l *Logging) DebugLevel(debugLevel string) {
	switch debugLevel {
	case "":
		l.Level.Set(slog.LevelError)
	case "Error":
		l.Level.Set(slog.LevelError)
	case "Warn":
		l.Level.Set(slog.LevelWarn)
	case "Info":
		l.Level.Set(slog.LevelInfo)
	case "Debug":
		l.Level.Set(slog.LevelDebug)
	}
}
