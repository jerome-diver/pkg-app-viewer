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

var loggingLevel *slog.LevelVar
var logging *slog.Logger

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
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
func NewLogger() *slog.Logger {

	loggingLevel = new(slog.LevelVar)
	opts := slog.HandlerOptions{
		Level: loggingLevel,
	}
	textHandler := NewPrettyHandler(os.Stdout, &opts)
	logging = slog.New(textHandler)
	return logging
}

func DebugLevel(debugLevel string) {
	switch debugLevel {
	case "":
		loggingLevel.Set(slog.LevelError)
	case "Error":
		loggingLevel.Set(slog.LevelError)
	case "Warn":
		loggingLevel.Set(slog.LevelWarn)
	case "Info":
		loggingLevel.Set(slog.LevelInfo)
	case "Debug":
		loggingLevel.Set(slog.LevelDebug)
	}
}
