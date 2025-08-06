// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package base

import "github.com/rs/zerolog"

type LogEntry interface {
	Any(string, any) LogEntry
	Int(string, int) LogEntry
	AnErr(string, error) LogEntry
	Str(string, string) LogEntry
	Stack() LogEntry
	Err(error) LogEntry
	Msgf(string, ...any)
	Msg(string)
}

type Logger interface {
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

func (l *DefaultLogger) Trace() LogEntry {
	return &DefaultLogEntry{l.Logger.Trace()}
}

func (l *DefaultLogger) Debug() LogEntry {
	return &DefaultLogEntry{l.Logger.Trace()}
}

func (l *DefaultLogger) Info() LogEntry {
	return &DefaultLogEntry{l.Logger.Trace()}
}

func (l *DefaultLogger) Warn() LogEntry {
	return &DefaultLogEntry{l.Logger.Trace()}
}

func (l *DefaultLogger) Error() LogEntry {
	return &DefaultLogEntry{l.Logger.Trace()}
}

func (l *DefaultLogger) Fatal() LogEntry {
	return &DefaultLogEntry{l.Logger.Trace()}
}

func (l *DefaultLogger) Panic() LogEntry {
	return &DefaultLogEntry{l.Logger.Trace()}
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
