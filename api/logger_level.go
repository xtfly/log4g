package api

import "strings"

// Level type for a logger
type Level int

// Log levels
const (
	Uninitialized Level = iota
	All
	Trace
	Debug
	Info
	Warn
	Error
	Critical
	Off
)

var levelStrings = map[Level]string{
	All:      "ALL",
	Trace:    "TRACE",
	Debug:    "DEBUG",
	Info:     "INFO",
	Warn:     "WARN",
	Error:    "ERROR",
	Critical: "CRITICAL",
	Off:      "OFF",
}

var lvlStrings = map[Level]string{
	All:      "ALL",
	Trace:    "TAC",
	Debug:    "DBG",
	Info:     "INF",
	Warn:     "WRN",
	Error:    "ERR",
	Critical: "CRI",
	Off:      "OFF",
}

// String returns the text for the level.
func (lvl Level) String() string {
	return levelStrings[lvl]
}

// ShortStr returns the short text for the level.
func (lvl Level) ShortStr() string {
	return lvlStrings[lvl]
}

// LevelFrom returns level from string.
func LevelFrom(str string) Level {
	for k, v := range levelStrings {
		if strings.ToUpper(str) == v {
			return k
		}
	}
	return Uninitialized
}
