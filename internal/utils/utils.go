package utils

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)

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

func GetLastPathComponent(url string) string {
	ss := strings.Split(url, "/")
	if len(ss) == 0 {
		return ""
	}
	c := ss[len(ss)-1]
	return c
}

func ConvertToPrettyJSON(object interface{}) string {
	jsonb, _ := json.MarshalIndent(object, "", "\t")
	return string(jsonb)
}

func ConvertToJSON(object interface{}) string {
	jsonb, _ := json.Marshal(object)
	return string(jsonb)
}
