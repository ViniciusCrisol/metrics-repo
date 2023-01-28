package log

import "go.uber.org/zap"

var (
	Logger *zap.Logger

	Error  = zap.Error
	String = zap.String
)

func init() {
	Logger, _ = zap.NewProduction()
}
