package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type Colors struct {
	BLACK, RED, GREEN, BROWN, BLUE, PURPLE, CYAN, LIGHT_GRAY, DARK_GRAY, LIGHT_RED, LIGHT_GREEN, YELLOW, LIGHT_BLUE, LIGHT_PURPLE, LIGHT_CYAN, LIGHT_WHITE, BOLD, FAINT, ITALIC, UNDERLINE, BLINK, NEGATIVE, CROSSED, END string
}

var colors = Colors{
	BLACK:       "\033[0;30m",
	RED:         "\033[0;31m",
	GREEN:       "\033[0;32m",
	BROWN:       "\033[0;33m",
	BLUE:        "\033[0;34m",
	PURPLE:      "\033[0;35m",
	CYAN:        "\033[0;36m",
	LIGHT_GRAY:  "\033[0;37m",
	DARK_GRAY:   "\033[1;30m",
	LIGHT_RED:   "\033[1;31m",
	LIGHT_GREEN: "\033[1;32m",
	YELLOW:      "\033[1;33m",
	LIGHT_BLUE:  "\033[1;34m",
	LIGHT_PURPLE: "\033[1;35m",
	LIGHT_CYAN:  "\033[1;36m",
	LIGHT_WHITE: "\033[1;37m",
	BOLD:        "\033[1m",
	FAINT:       "\033[2m",
	ITALIC:      "\033[3m",
	UNDERLINE:   "\033[4m",
	BLINK:       "\033[5m",
	NEGATIVE:    "\033[7m",
	CROSSED:     "\033[9m",
	END:         "\033[0m",
}

func init() {
	if !isTerminal(os.Stdout) {
		resetColors()
	} else if runtime.GOOS == "windows" {
		enableVirtualTerminalProcessing()
	}
}

func isTerminal(f *os.File) bool {
	return strings.Contains(os.Getenv("TERM"), "xterm") || strings.Contains(os.Getenv("TERM"), "screen")
}

func resetColors() {
	colors = Colors{}
}

func enableVirtualTerminalProcessing() {
	// Enable virtual terminal processing on Windows
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleMode")
	handle := kernel32.NewProc("GetStdHandle").Call(uintptr(-11))
	proc.Call(handle, uintptr(7))
}

type Logger struct {
	*log.Logger
}

func NewLogger(prefix string, flag int) *Logger {
	return &Logger{log.New(os.Stdout, prefix, flag)}
}

func (l *Logger) formatMessage(level string, msg string) string {
	var color string
	switch level {
	case "DEBUG":
		color = colors.CYAN
	case "INFO":
		color = colors.BLUE
	case "WARN":
		color = colors.YELLOW
	case "ERROR":
		color = colors.RED
	case "CRITICAL":
		color = colors.LIGHT_RED
	default:
		color = colors.END
	}
	return fmt.Sprintf("%s%s: %s%s", color, level, msg, colors.END)
}

func (l *Logger) Debug(v ...interface{}) {
	l.Println(l.formatMessage("DEBUG", fmt.Sprint(v...)))
}

func (l *Logger) Info(v ...interface{}) {
	l.Println(l.formatMessage("INFO", fmt.Sprint(v...)))
}

func (l *Logger) Warn(v ...interface{}) {
	l.Println(l.formatMessage("WARN", fmt.Sprint(v...)))
}

func (l *Logger) Error(v ...interface{}) {
	l.Println(l.formatMessage("ERROR", fmt.Sprint(v...)))
}

func (l *Logger) Critical(v ...interface{}) {
	l.Println(l.formatMessage("CRITICAL", fmt.Sprint(v...)))
}

var (
	DEBUG = false // set this based on your config
	logger *Logger
)

func init() {
	logger = NewLogger("", log.LstdFlags)
	if DEBUG {
		logger.SetFlags(log.LstdFlags | log.Lshortfile)
	}
}

func main() {
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
	logger.Critical("This is a critical message")
}
