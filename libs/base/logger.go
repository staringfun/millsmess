// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package base

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
)

type LoggerConfig struct {
	Level             int8 `json:"level" yaml:"level" default:"0"`
	EnablePrettyPrint bool `json:"enablePrettyPrint" yaml:"enablePrettyPrint"`
}

func NewLogger(config LoggerConfig) Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	if config.EnablePrettyPrint {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	logger = logger.Level(zerolog.Level(config.Level))
	return &DefaultLogger{logger}
}

func GetLogger(ctx context.Context) Logger {
	return &DefaultLogger{Logger: *zerolog.Ctx(ctx)}
}

func LogErrorIfNotNil(ctx context.Context, err error, msg string) bool {
	hasError := err != nil
	if hasError {
		GetLogger(ctx).Error().Err(err).Msg(msg)
	}
	return hasError
}

type LogEntry interface {
	Any(string, any) LogEntry
	Int(string, int) LogEntry
	AnErr(string, error) LogEntry
	Str(string, string) LogEntry
	Err(error) LogEntry
	Stack() LogEntry
	WithContext(context.Context) context.Context
	Msgf(string, ...any)
	Msg(string)
}

type Logger interface {
	With() LogEntry
	Trace() LogEntry
	Debug() LogEntry
	Info() LogEntry
	Warn() LogEntry
	Error() LogEntry
	Fatal() LogEntry
	Panic() LogEntry
}

type DefaultLogger struct {
	zerolog.Logger
}

func (l *DefaultLogger) With() LogEntry {
	return &DefaultContextEntry{l.Logger.With()}
}

func (l *DefaultLogger) Trace() LogEntry {
	return &DefaultLogEntry{l.Logger.Trace()}
}

func (l *DefaultLogger) Debug() LogEntry {
	return &DefaultLogEntry{l.Logger.Debug()}
}

func (l *DefaultLogger) Info() LogEntry {
	return &DefaultLogEntry{l.Logger.Info()}
}

func (l *DefaultLogger) Warn() LogEntry {
	return &DefaultLogEntry{l.Logger.Warn()}
}

func (l *DefaultLogger) Error() LogEntry {
	return &DefaultLogEntry{l.Logger.Error()}
}

func (l *DefaultLogger) Fatal() LogEntry {
	return &DefaultLogEntry{l.Logger.Fatal()}
}

func (l *DefaultLogger) Panic() LogEntry {
	return &DefaultLogEntry{l.Logger.Panic()}
}

type DefaultContextEntry struct {
	zerolog.Context
}

func (e *DefaultContextEntry) Err(err error) LogEntry {
	return &DefaultContextEntry{e.Context.Err(err)}
}

func (e *DefaultContextEntry) Stack() LogEntry {
	return &DefaultContextEntry{e.Context.Stack()}
}

func (e *DefaultContextEntry) Str(key, value string) LogEntry {
	return &DefaultContextEntry{e.Context.Str(key, value)}
}

func (e *DefaultContextEntry) AnErr(key string, value error) LogEntry {
	return &DefaultContextEntry{e.Context.AnErr(key, value)}
}

func (e *DefaultContextEntry) Int(key string, value int) LogEntry {
	return &DefaultContextEntry{e.Context.Int(key, value)}
}

func (e *DefaultContextEntry) Any(key string, value any) LogEntry {
	return &DefaultContextEntry{e.Context.Any(key, value)}
}

func (e *DefaultContextEntry) WithContext(ctx context.Context) context.Context {
	return e.Logger().WithContext(ctx)
}

func (e *DefaultContextEntry) Msg(msg string) {
}

func (e *DefaultContextEntry) Msgf(msg string, args ...any) {
}

type DefaultLogEntry struct {
	*zerolog.Event
}

func (e *DefaultLogEntry) Err(err error) LogEntry {
	return &DefaultLogEntry{e.Event.Err(err)}
}

func (e *DefaultLogEntry) Stack() LogEntry {
	return &DefaultLogEntry{e.Event.Stack()}
}

func (e *DefaultLogEntry) Str(key, value string) LogEntry {
	return &DefaultLogEntry{e.Event.Str(key, value)}
}

func (e *DefaultLogEntry) AnErr(key string, value error) LogEntry {
	return &DefaultLogEntry{e.Event.AnErr(key, value)}
}

func (e *DefaultLogEntry) Int(key string, value int) LogEntry {
	return &DefaultLogEntry{e.Event.Int(key, value)}
}

func (e *DefaultLogEntry) Any(key string, value any) LogEntry {
	return &DefaultLogEntry{e.Event.Any(key, value)}
}

func (e *DefaultLogEntry) WithContext(ctx context.Context) context.Context {
	return ctx
}
