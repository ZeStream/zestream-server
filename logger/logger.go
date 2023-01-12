package logger

import (
	"context"
	"os"
	"sync"
	"zestream-server/constants"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Z = map[string]interface{}

var once sync.Once
var singleton *zap.SugaredLogger

/*
Init initializes a thread-safe singleton logger
This would be called from a main method when the application starts up
This function would ideally, take zap configuration, but is left out
in favor of simplicity using the example logger.
*/
func Init(name string, logLevel string) {
	// once ensures the singleton is initialized only once
	once.Do(func() {
		if name == "" {
			name = "zestream-server"
		}

		// by default, this sets the minimum logging level to info
		cfg := zap.NewProductionConfig()
		cfg.Level.SetLevel(parseLogLevel(logLevel))

		logDir := os.Getenv("LOG_DIR")

		if logDir == "" {
			logDir = "logs"
		}

		MakeDir(logDir, 0755)

		outputFile := logDir + "/" + name + ".log"

		// make the logDir if it does not exist

		cfg.OutputPaths = []string{outputFile}

		cfg.EncoderConfig.TimeKey = "logTime"
		cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

		cfg.EncoderConfig.MessageKey = "message"

		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		builtLogger, _ := cfg.Build(zap.AddCallerSkip(2), zap.AddStacktrace(zapcore.DebugLevel))

		singleton = builtLogger.Sugar()
	})
}

func parseLogLevel(logLevel string) zapcore.Level {
	switch logLevel {
	case "debug":
		{
			return zapcore.DebugLevel
		}
	case "info":
		{
			return zapcore.InfoLevel
		}
	case "error":
		{
			return zapcore.ErrorLevel
		}
	case "warn":
		{
			return zapcore.WarnLevel
		}
	default:
		{
			return zapcore.InfoLevel
		}
	}
}

func MakeDir(path string, mode os.FileMode) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, mode)
	}
}

func Info(ctx context.Context, msg string, data Z) {
	logIt(ctx, zapcore.InfoLevel, msg, data)
}

func Debug(ctx context.Context, msg string, data Z) {
	logIt(ctx, zapcore.DebugLevel, msg, data)
}
func Error(ctx context.Context, msg string, data Z) {
	logIt(ctx, zapcore.ErrorLevel, msg, data)
}

func logIt(ctx context.Context, level zapcore.Level, message string, data Z) {
	if singleton == nil {
		Init("", os.Getenv("LOG_LEVEL"))
	}

	if ctx == nil {
		ctx = context.TODO()
	}

	modifiedArgs := ingestData(ctx, data)

	switch level {
	case zapcore.ErrorLevel:
		{
			singleton.Errorw(message, modifiedArgs...)
			break
		}
	case zapcore.WarnLevel:
		{
			singleton.Warnw(message, modifiedArgs...)
			break
		}
	case zapcore.InfoLevel:
		{
			singleton.Infow(message, modifiedArgs...)
			break
		}
	case zapcore.DebugLevel:
		{
			singleton.Debugw(message, modifiedArgs...)
			break
		}
	}
}

func ingestData(ctx context.Context, data Z) []interface{} {
	if data == nil {
		data = Z{}
	}

	data[constants.TRANSACTION_ID_KEY] = ctx.Value(constants.TRANSACTION_ID_KEY)

	if data[constants.TRANSACTION_ID_KEY] == nil || data[constants.TRANSACTION_ID_KEY] == "" {
		data[constants.TRANSACTION_ID_KEY] = uuid.New().String()
	}

	// organize Data in key, value, key, value... in an array of interface
	argsLen := len(data) * 2 // key + value
	args := make([]interface{}, argsLen)

	i := 0
	for k, v := range data {
		args[i] = k
		args[i+1] = v
		i += 2
	}

	return args
}
