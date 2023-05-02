package Logger

import (
	"github.com/pavel-one/EdgeGPT-Go/internal/Helpers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewLogger(channel string) *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.Level = getLevel()

	if Helpers.FindInSlice(os.Args, "chat") {
		cfg.OutputPaths = []string{
			"logs/app.log",
		}
		cfg.ErrorOutputPaths = []string{
			"logs/err.log",
		}
	}

	logger, _ := cfg.Build()
	sugar := logger.Sugar()

	return sugar.Named(channel)
}

func getLevel() zap.AtomicLevel {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		l, err := zapcore.ParseLevel("INFO")
		if err != nil {
			panic(err)
		}

		return zap.NewAtomicLevelAt(l)
	}

	l, err := zapcore.ParseLevel(level)
	if err != nil {
		panic(err)
	}

	return zap.NewAtomicLevelAt(l)
}
