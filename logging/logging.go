package logging

import (
	"fmt"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger

// Config config
type Config struct {
	FileName    string
	MaxSizeInMB int
	MaxBackups  int
	MaxDays     int
	Compress    bool
}

// Logging log
type Logging struct {
	*zap.SugaredLogger
	*Config
	*lumberjack.Logger
}

// New init
func New(config ...*Config) *Logging {
	var cfg *Config
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	}

	if cfg == nil {
		cfg = &Config{
			FileName:    fmt.Sprintf("logs/%s.log", time.Now().Format("01-02-2006")),
			MaxBackups:  7,
			MaxDays:     2,
			MaxSizeInMB: 500,
			Compress:    true,
		}
	}

	writerSyncer, rotateLogger := getLogWriter(cfg)

	var encoder = getEncoder()
	var core = zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)
	var logger = zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()

	var instance = &Logging{
		sugarLogger,
		cfg,
		rotateLogger,
	}

	return instance
}

func getLogWriter(config *Config) (zapcore.WriteSyncer, *lumberjack.Logger) {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.FileName,
		MaxSize:    config.MaxSizeInMB,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxDays,
		Compress:   config.Compress,
	}

	return zapcore.AddSync(lumberJackLogger), lumberJackLogger
}

func getEncoder() zapcore.Encoder {
	var encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
