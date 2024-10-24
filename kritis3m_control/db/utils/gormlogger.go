package utils

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

type GormZerologger struct {
	zlog       zerolog.Logger
	logLevel   logger.LogLevel
	slowThresh time.Duration
}

// NewGormZerologger creates a new GORM logger that uses Zerolog
func NewGormZerologger(zlog zerolog.Logger, logLevel logger.LogLevel, slowThreshold time.Duration) *GormZerologger {
	return &GormZerologger{
		zlog:       zlog,
		logLevel:   logLevel,
		slowThresh: slowThreshold,
	}
}

// LogMode sets the log level for the GORM logger
func (g *GormZerologger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *g
	newLogger.logLevel = level
	return &newLogger
}

// Info logs general information using Zerolog
func (g *GormZerologger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel >= logger.Info {
		g.zlog.Info().Msgf(msg, data...)
	}
}

// Warn logs warnings using Zerolog
func (g *GormZerologger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel >= logger.Warn {
		g.zlog.Warn().Msgf(msg, data...)
	}
}

// Error logs errors using Zerolog
func (g *GormZerologger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel >= logger.Error {
		g.zlog.Error().Msgf(msg, data...)
	}
}

// Trace logs the SQL queries, execution time, and possible errors or warnings
func (g *GormZerologger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && g.logLevel >= logger.Error:
		g.zlog.Error().
			Err(err).
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Str("sql", sql).
			Msg("SQL execution error")
	case elapsed > g.slowThresh && g.slowThresh != 0 && g.logLevel >= logger.Warn:
		g.zlog.Warn().
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Str("sql", sql).
			Msg("Slow SQL execution")
	case g.logLevel >= logger.Info:
		g.zlog.Info().
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Str("sql", sql).
			Msg("SQL executed")
	}
}
