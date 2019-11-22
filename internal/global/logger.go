package global

import (
	"os"

	pkglogger "phonebook/pkg/logger"

	"github.com/go-kit/kit/log"
)

func InitLogger() log.Logger {
	logger := pkglogger.NewKitxLogger(os.Stderr, os.Stdout)
	logger = logger.SetCaller(pkglogger.Caller(5, 5))
	return logger
}
