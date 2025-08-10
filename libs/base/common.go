// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package base

import (
	"context"
	"io"
	"time"
)

type RetryStrategy interface {
	Next() time.Duration
}

type MaxAttemptsStrategyConfig struct {
	MaxAttempts    int           `json:"maxAttempts" yaml:"maxAttempts" env:"ATTEMPTS_MAX_ATTEMPTS" default:"10"`
	AttemptTimeout time.Duration `json:"attemptTimeout" yaml:"attemptTimeout" env:"ATTEMPTS_ATTEMPT_TIMEOUT" default:"10ms"`
}

type MaxAttemptsStrategy struct {
	Config   MaxAttemptsStrategyConfig
	attempts int
}

func (s *MaxAttemptsStrategy) Next() time.Duration {
	s.attempts += 1
	if s.attempts >= s.Config.MaxAttempts {
		return -1
	}
	return s.Config.AttemptTimeout
}

type ContextWriter interface {
	Write(p []byte, ctx context.Context) (n int, err error)
}

type ContextWriterWrapper struct {
	Writer io.Writer
}

func (w *ContextWriterWrapper) Write(data []byte, ctx context.Context) (int, error) {
	return w.Writer.Write(data)
}

type RetryWriter struct {
	Writer   ContextWriter
	Strategy RetryStrategy
}

func (w *RetryWriter) Write(data []byte, ctx context.Context) (int, error) {
	var err error
	var n int
	for err = ctx.Err(); err == nil; {
		n, err = w.Writer.Write(data, ctx)
		if err == nil {
			break
		}
		if w.Strategy != nil {
			d := w.Strategy.Next()
			if d < 1 {
				break
			}
			time.Sleep(d)
		}
	}
	return n, err
}

type Closable interface {
	Close() error
}
