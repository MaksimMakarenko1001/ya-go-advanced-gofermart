package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type HTTPInfo struct {
	URI      string        `json:"uri"`
	Method   string        `json:"method"`
	Duration time.Duration `json:"duration"`
	Response ResponseInfo  `json:"response"`
}

type ResponseInfo struct {
	Size   int          `json:"size"`
	Status int          `json:"status"`
	Body   bytes.Buffer `json:"body"`
}

type ZapLogger struct {
	logger *zap.Logger
}

func New(config Config) *ZapLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	logger.Sync()

	lvl, err := zap.ParseAtomicLevel(string(config.Level))
	if err != nil {
		logger.Panic("log level not ok", zap.Error(err))
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		logger.Panic("log config not ok", zap.Error(err))
	}
	logger.Info("zap", zap.String("log level", zl.Level().String()))
	return &ZapLogger{logger: zl}
}

func (zl *ZapLogger) LogHTTP(info HTTPInfo) {
	defer zl.logger.Sync()

	b, err := json.Marshal(info)
	if err != nil {
		zl.logger.Error("log not ok, %w", zap.Error(err))
	}

	infoLabel := "body"
	if info.Response.Status != http.StatusOK {
		infoLabel = "error"
	}

	msg := fmt.Sprint(info.Method, info.URI)

	zl.logger.Info(msg, zap.ByteString(infoLabel, info.Response.Body.Bytes()))
	zl.logger.Debug(msg, zap.ByteString("raw", b))
}

func (zl *ZapLogger) Infof(segment, format string, args ...any) {
	logger := zl.logger.Named(segment)
	logger.Sugar().Infof(format, args...)
}

func (zl *ZapLogger) Errorf(segment, format string, args ...any) {
	logger := zl.logger.Named(segment)
	logger.Sugar().Errorf(format, args...)
}

func (zl *ZapLogger) Debugf(segment, format string, args ...any) {
	logger := zl.logger.Named(segment)
	logger.Sugar().Debugf(format, args...)
}

func (zl *ZapLogger) Panicf(segment, format string, args ...any) {
	logger := zl.logger.Named(segment)
	logger.Sugar().Panicf(format, args...)
}
