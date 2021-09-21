package logging

import (
    "os"
    "github.com/sirupsen/logrus"

    "github.com/jcatana/reaper/config"
)

type Log *logrus.Logger

func NewLogger(cfg *config.Config) *logrus.Logger {
    var log = logrus.New()
    logrus.SetFormatter(&logrus.TextFormatter{
        DisableColors: true,
        FullTimestamp: true,
    })
    logLevel, err := logrus.ParseLevel(cfg.GetLogLevel())
    if err != nil {
        logLevel = logrus.InfoLevel
    }
    log.Out = os.Stdout
    log.SetLevel(logLevel)
    return log
}

