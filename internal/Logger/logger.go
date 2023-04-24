package Logger

import (
	"github.com/pavel-one/EdgeGPT-Go/internal/Helpers"
	"go.uber.org/zap"
	"os"
)

func NewLogger(channel string) *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()

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
