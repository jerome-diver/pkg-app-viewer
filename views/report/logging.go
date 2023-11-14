package reeport

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

type Logging interface {
	SetError(error)
	GetError() error
	CheckError(string) bool
	DebugLevel(string)
	Log(msg string, args ...any)
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type logging struct { // proxy to add methd checkError
	*slog.Logger
	Err     error
	Handler *PrettyHandler
	Level   *slog.LevelVar
	Opt     *slog.HandlerOptions
}

var single_log Logging

func (l logging) SetError(err error) {
	l.Err = err
}

func (l logging) GetError() error {
	return l.Err
}

func (l logging) CheckError(msg string) bool {
	if l.Err != nil {
		if os.IsNotExist(l.Err) {
			l.Warn("file doesn't exist")
		} else {
			l.Error(msg, slog.String("err", l.Err.Error()))
		}
		return true
	}
	return false
}

func (l logging) DebugLevel(debugLevel string) {
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

func (l logging) Log(msg string, args ...any) {
	l.Log(msg, args)
}

func (l logging) Info(msg string, args ...any) {
	l.Info(msg, args)
}

func (l logging) Debug(msg string, args ...any) {
	l.Debug(msg, args)
}

func (l logging) Warn(msg string, args ...any) {
	l.Warn(msg, args)
}

func (l logging) Error(msg string, args ...any) {
	l.Error(msg, args)
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

func GetLogger() Logging {
	if single_log == nil {
		loggingLevel := new(slog.LevelVar)
		opts := slog.HandlerOptions{
			Level: loggingLevel,
		}
		textHandler := NewPrettyHandler(os.Stdout, &opts)
		single_log = &logging{
			slog.New(textHandler),
			nil,
			textHandler,
			loggingLevel,
			&opts,
		}
	}
	return single_log
}
