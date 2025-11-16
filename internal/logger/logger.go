package logger

import (
	"context"

	"github.com/atharvamhaske/go-tly/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

var L *Logger //global logger

// newlogger - DI friendly
func NewLogger(cfg *config.Config) (*Logger, error) {
	// here cfg is argument and *Logger and error are return values/types
	var zapCfg zap.Config

	if cfg.App.Env == "dev" {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	zapCfg.EncoderConfig.TimeKey = "timestamp"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	base, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{SugaredLogger: base.Sugar()}, nil
}

func GetLoggerWithCtx(ctx context.Context) *Logger{
	reqID, _ := ctx.Value("request_id").(string)
	userID, _ := ctx.Value("user_id").(string)

	return &Logger {
		SugaredLogger: L.SugaredLogger.With(
			"request_id", reqID,
			"user_id", userID,
		),
	}
}

//additional wrapper if need(directly used from some codebase)
func (l *Logger) Debugf(t string, args ...interface{}) { l.SugaredLogger.Debugf(t, args...) }
func (l *Logger) Infof(t string, args ...interface{})  { l.SugaredLogger.Infof(t, args...) }
func (l *Logger) Warnf(t string, args ...interface{})  { l.SugaredLogger.Warnf(t, args...) }
func (l *Logger) Errorf(t string, args ...interface{}) { l.SugaredLogger.Errorf(t, args...) }
func (l *Logger) Fatalf(t string, args ...interface{}) { l.SugaredLogger.Fatalf(t, args...) }
