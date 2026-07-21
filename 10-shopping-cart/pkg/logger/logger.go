package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

var Log *zerolog.Logger

type LoggerConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	IsDev      string
}

func InitLogger(cfg LoggerConfig) {
	Log = NewLogger(cfg)
}

func NewLogger(cfg LoggerConfig) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	lvl, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	var writer io.Writer

	if cfg.IsDev == "development" {
		if strings.Contains(cfg.Filename, "app.logs") {
			writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		} else {
			writer = PrettyJSONWriter{Writer: os.Stdout}
		}
	} else {
		writer = &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger
}

type PrettyJSONWriter struct {
	Writer io.Writer
}

func (w PrettyJSONWriter) Write(p []byte) (n int, err error) {
	var prettyJSON bytes.Buffer

	err = json.Indent(&prettyJSON, p, "", "  ")
	if err != nil {
		return w.Writer.Write(p)
	}

	return w.Writer.Write(prettyJSON.Bytes())
}

type contextKey string

const TraceIdKey contextKey = "trace_id"

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIdKey).(string); ok {
		return traceID
	}

	return ""
}
