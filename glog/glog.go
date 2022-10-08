package glog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// Terminal Colors
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Add adds the color to the given string
func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

// Any adds the color to any type
func (c Color) Any(s interface{}) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", uint8(c), s)
}

// Bold adds a bold color to the given string
func (c Color) Bold(s string) string {
	return fmt.Sprintf("\x1b[1;%dm%s\x1b[0m", uint8(c), s)
}

// Color represents a text color.
type Color uint8

// Logger Level
const (
	InfoLevel Level = iota
	WarnLevel
	ErrorLevel
)

var levelTexts = map[Level]string{
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERRO",
}

var levelColors = map[Level]Color{
	InfoLevel:  Green,
	WarnLevel:  Yellow,
	ErrorLevel: Red,
}

// Level of logger
type Level uint8

func (l Level) String() string { return levelTexts[l] }

// Colored return colorized level string
func (l Level) Colored() string {
	return levelColors[l].Bold(l.String())
}

// Prefixer return colorized logger prefix
func (l Level) Prefixer() string {
	return fmt.Sprintf("%s - ", l.Colored())
}

var (
	mtx         sync.Mutex
	logFlags    = log.LstdFlags | log.Lshortfile | log.Lmsgprefix
	infoLogger  = log.New(os.Stdout, InfoLevel.Prefixer(), logFlags)
	warnLogger  = log.New(os.Stdout, WarnLevel.Prefixer(), logFlags)
	errorLogger = log.New(os.Stdout, ErrorLevel.Prefixer(), logFlags)
	loggers     = []*log.Logger{infoLogger, warnLogger, errorLogger}
)

// Glogger -
type Glogger struct {
	level Level
	*log.Logger
}

// New -
// func log.New(out io.Writer, prefix string, flag int) *log.Logger
func New(out io.Writer, lvl Level, flag int) *Glogger {
	return &Glogger{
		level:  lvl,
		Logger: log.New(out, lvl.Colored(), flag),
	}
}

// LevelFunc control how to log message
func (g Glogger) LevelFunc(lvl Level) {
	if g.level < lvl {
		g.SetOutput(ioutil.Discard)
	}
}

// Logger Methods Alias
var (
	Info   = infoLogger.Print
	Infof  = infoLogger.Printf
	Warn   = warnLogger.Print
	Warnf  = warnLogger.Printf
	Error  = errorLogger.Print
	Errorf = errorLogger.Printf
)

var glogers = []*Glogger{
	New(os.Stdout, InfoLevel, logFlags),
	New(os.Stdout, WarnLevel, logFlags),
	New(os.Stdout, ErrorLevel, logFlags),
}

// SetLevel set log level
func SetLevel(lvl Level) {
	mtx.Lock()
	defer mtx.Unlock()

	for _, lg := range loggers {
		lg.SetOutput(os.Stdout)
	}

	if InfoLevel < lvl {
		infoLogger.SetOutput(ioutil.Discard)
	}
	if WarnLevel < lvl {
		warnLogger.SetOutput(ioutil.Discard)
	}
	if ErrorLevel < lvl {
		errorLogger.SetOutput(ioutil.Discard)
	}
}

// SetGlogLevel set gloger Level
func SetGlogLevel(lvl Level) {
	mtx.Lock()
	defer mtx.Unlock()

	for _, g := range glogers {
		g.SetOutput(os.Stdout)
		g.LevelFunc(lvl)
	}
}
