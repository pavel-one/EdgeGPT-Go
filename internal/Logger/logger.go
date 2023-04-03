package Logger

import "go.uber.org/zap"

func NewLogger(channel string) *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	logger.With(zap.String("channel", channel))

	return logger.Sugar()
}
