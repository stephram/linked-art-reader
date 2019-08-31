package utils

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(false)

	switch fmt := os.Getenv("LOG_FORMAT"); fmt {
	case "text":
		log.SetFormatter(&log.TextFormatter{
			// DisableColors: false,
			// ForceColors: true,
			FullTimestamp: true,
		})
	case "json":
		setJSONLogFormat()
	default:
		// log.Printf("unknown LOG_FORMAT value: '%s'", fmt)
		setJSONLogFormat()
	}

	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = log.InfoLevel
		// log.WithError(err).Printf("defaulted to %s", level.String())
	}
	log.SetLevel(level)
	// log.Printf("log level set to %s", log.GetLevel().String())
}

func setJSONLogFormat() {
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint: false,
	})
	log.Info("set JSON log format")
}

// GetLogger needs to be called once to ensure logrus is configured correctly.
func GetLogger() *log.Logger {
	return log.StandardLogger()
}

func convertISODateFormatToDB(dateStr string) time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		log.WithError(err).Errorf("Unable to parse date: %s", dateStr)
	}
	return t
}

func ConvertPythonDateToTime(dateStr string) time.Time {
	layout := "2006-01-02 15:04:05+0000"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		log.WithError(err).Errorf("Unable to parse date: %s", dateStr)
	}
	return t
}
