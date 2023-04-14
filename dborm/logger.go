package dborm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/open-tdp/go-helper/logman"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Logger struct {
	config *logger.Config
	logger *logman.Logger
}

func NewLogger() logger.Interface {

	logger.Default = &Logger{
		logger: logman.Named("gorm"),
		config: &logger.Config{
			IgnoreRecordNotFoundError: false,
			SlowThreshold:             5 * time.Second,
		},
	}

	return logger.Default

}

func (l Logger) LogMode(level logger.LogLevel) logger.Interface {

	l.config.LogLevel = level
	return l

}

func (l Logger) Info(ctx context.Context, msg string, args ...any) {

	l.logger.Info(fmt.Sprintf(msg, args...))

}

func (l Logger) Warn(ctx context.Context, msg string, args ...any) {

	l.logger.Warn(fmt.Sprintf(msg, args...))

}

func (l Logger) Error(ctx context.Context, msg string, args ...any) {

	l.logger.Error(fmt.Sprintf(msg, args...))

}

func (l Logger) Trace(ctx context.Context, begin time.Time, fn func() (string, int64), err error) {

	cfg := l.config
	sql, rows := fn()
	elapsed := time.Since(begin)

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !cfg.IgnoreRecordNotFoundError):
		l.logger.Error("trace error", "error", err, "sql", sql, "rows", rows, "elapsed", elapsed)
	case elapsed > cfg.SlowThreshold && cfg.SlowThreshold != 0:
		slow := fmt.Sprintf("trace slow sql >= %v", cfg.SlowThreshold)
		l.logger.Warn(slow, "sql", sql, "rows", rows, "elapsed", elapsed)
	default:
		l.logger.Info("trace query", "sql", sql, "rows", rows, "elapsed", elapsed)
	}

}
