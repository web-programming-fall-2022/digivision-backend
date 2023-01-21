package bootstrap

import (
	"github.com/sirupsen/logrus"
)

func AdjustLogLevel(logLevel string) {
	logrusLevel, err := logrus.ParseLevel(logLevel)
	logrus.SetReportCaller(true)
	if err != nil {
		logrus.WithError(err).WithField("logLevel", logLevel).Error("failed to AdjustLogLevel")
		logrus.SetLevel(logrus.DebugLevel)
		return
	}
	logrus.WithField("LogLevel", logLevel).Info("logLevel has been set")
	logrus.SetLevel(logrusLevel)
}
