package ESL_Logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Severity string

const (
	SeverityDebug Severity = "DEBUG"
	SeverityInfo  Severity = "INFO"
	SeverityWarn  Severity = "WARN"
	SeverityError Severity = "ERROR"
	SeverityFatal Severity = "FATAL"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func defaultConfig(name string) *Config {
	return &Config{
		Output: "stdout",
		Name:   name,
		Format: "%s %s %s: %s",
	}
}

func prodConfig(name, output string) *Config {
	return &Config{
		Output:     output,
		Name:       name,
		Format:     "[%s] %s - %s: %s",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     30,
	}
}

type Logger struct {
	cnf     *Config
	logChan chan string
	wg      sync.WaitGroup
}

type Config struct {
	Name       string
	Output     string // stdout, <filepath>
	Format     string // log message format
	MaxSize    int    // megabytes
	MaxBackups int
	MaxAge     int // days
}

func NewLogger(name, env string, output ...string) *Logger {
	var config *Config
	if env == "production" && len(output) > 0 {
		config = prodConfig(name, output[0])
	} else {
		config = defaultConfig(name)
	}

	logger := &Logger{
		cnf:     config,
		logChan: make(chan string, 100),
	}

	if env == "production" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   logger.cnf.Output,
			MaxSize:    logger.cnf.MaxSize,
			MaxBackups: logger.cnf.MaxBackups,
			MaxAge:     logger.cnf.MaxAge,
		})
		log.SetFlags(0)
	}

	logger.wg.Add(1)
	go logger.processLogs()

	return logger
}

func (l *Logger) processLogs() {
	defer l.wg.Done()
	for msg := range l.logChan {
		switch l.cnf.Output {
		case "stdout":
			fmt.Println(msg)
		default:
			log.Println(msg)
		}
	}
}

func (l *Logger) LogF(level Severity, format string, args ...any) string {
	color := getColor(level)
	message := fmt.Sprintf(
		l.cnf.Format,
		time.Now().UTC().Format(time.RFC3339Nano),
		l.cnf.Name, level, fmt.Sprintf(format, args...),
	)

	coloredMessage := fmt.Sprintf("%s%s%s", color, message, ColorReset)
	l.logChan <- coloredMessage

	return message
}

func getColor(level Severity) string {
	switch level {
	case SeverityDebug:
		return ColorCyan
	case SeverityInfo:
		return ColorBlue
	case SeverityWarn:
		return ColorYellow
	case SeverityError:
		return ColorRed
	case SeverityFatal:
		return ColorPurple
	default:
		return ColorWhite
	}
}

func (l *Logger) Debug(format string, args ...any) string {
	return l.LogF(SeverityDebug, format, args...)
}

func (l *Logger) Info(format string, args ...any) string {
	return l.LogF(SeverityInfo, format, args...)
}

func (l *Logger) Warn(format string, args ...any) string {
	return l.LogF(SeverityWarn, format, args...)
}

func (l *Logger) ErrorF(format string, args ...any) string {
	return l.LogF(SeverityError, format, args...)
}

func (l *Logger) Fatal(format string, args ...any) {
	l.LogF(SeverityFatal, format, args...)
	os.Exit(1)
}

func (l *Logger) Close() {
	close(l.logChan)
	l.wg.Wait()
}
