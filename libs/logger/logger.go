package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/maxguuse/birdcord/libs/config"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

type Logger interface {
	Info(input string, fields ...any)
	Error(input string, fields ...any)
	Warn(input string, fields ...any)
	Debug(input string, fields ...any)
	Handler() slog.Handler
}

type logger struct {
	log *slog.Logger
}

func New(service string) func() Logger {
	return func() Logger {
		return newLogger(service)
	}
}

func newLogger(service string) Logger { //nolint: ireturn
	var log *slog.Logger

	replaceAttr := func(groups []string, attr slog.Attr) slog.Attr {
		if attr.Key == slog.SourceKey {
			source := attr.Value.Any().(*slog.Source) //nolint: forcetypeassert

			pathParts := strings.Split(source.File, "/")
			p := pathParts[len(pathParts)-4:]
			newPath := strings.Join(p, "/")

			functionParts := strings.Split(source.Function, "/")
			p = functionParts[len(functionParts)-1:]
			newFunction := strings.Join(p, "/")

			newSource := fmt.Sprintf("%s:%s:%d", newPath, newFunction, source.Line)

			return slog.String(slog.SourceKey, newSource)
		}

		return attr
	}

	if env := os.Getenv("ENVIRONMENT"); env == config.EnvProduction {
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level:       slog.LevelWarn,
					AddSource:   true,
					ReplaceAttr: replaceAttr,
				},
			),
		)
	} else {
		zerologLogger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})

		log = slog.New(
			slogzerolog.Option{
				Level:       slog.LevelDebug,
				Logger:      &zerologLogger,
				AddSource:   true,
				ReplaceAttr: replaceAttr,
			}.NewZerologHandler(),
		)
	}

	if service != "" {
		log = log.With(slog.String("service", service))
	}

	return &logger{
		log: log,
	}
}

func (c *logger) handle(level slog.Level, input string, fields ...any) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, input, pcs[0])
	for _, f := range fields {
		r.Add(f)
	}

	if !c.log.Handler().Enabled(context.Background(), level) {
		return
	}

	_ = c.log.Handler().Handle(context.Background(), r) //nolint: errcheck
}

func (c *logger) Info(input string, fields ...any) {
	c.handle(slog.LevelInfo, input, fields...)
}

func (c *logger) Warn(input string, fields ...any) {
	c.handle(slog.LevelWarn, input, fields...)
}

func (c *logger) Error(input string, fields ...any) {
	c.handle(slog.LevelError, input, fields...)
}

func (c *logger) Debug(input string, fields ...any) {
	c.handle(slog.LevelDebug, input, fields...)
}

func (c *logger) Handler() slog.Handler {
	return c.log.Handler()
}
