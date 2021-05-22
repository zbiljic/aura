package aurafx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type loggerFxPrinter struct {
	log *zap.SugaredLogger
}

func NewFxLogger(log ...*zap.SugaredLogger) fx.Printer {
	if len(log) == 0 {
		logger, _ := zap.NewProduction()
		log = append(log, logger.Sugar())
	}
	return loggerFxPrinter{log[0]}
}

func (l loggerFxPrinter) Printf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l loggerFxPrinter) String() string {
	return "zap"
}
