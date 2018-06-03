/*
Package logging provides functionality for printing log messages at various importance levels.
*/
package logging

import (
  "fmt"
  "io"
  "os"
  "runtime"
  "strings"
  "time"
)

// Available verbosity levels.
const (
  // Print LOG or higher priority messages. LOG is directed to Stdout by default.
  LOG = iota
  // Print only INFO or higher priority messages. INFO is directed to Stdout by default. This is the default verbosity level.
  INFO
  // Print WARN or higher priority messages. WARN is directed to Stderr by default.
  WARN
  // Print ERROR or higher priority messages. ERROR is directed to Stderr by default.
  ERROR
  // Set this verbosity level to print only critical messages. CRITICAL is directed to Stderr by default.
  CRITICAL
)

// A set of predefined timestamp formats. You can also use layouts from the Golang package "time".
// TS_FMT_TIME_MILLI is used by default when timestamp prefix is enabled.
const (
  TS_FMT_DATE               = "2006-01-02"
  TS_FMT_TIME               = "15:04:05"
  TS_FMT_TIME_MILLI         = "15:04:05.000"
  TS_FMT_TIME_MICRO         = "15:04:05.000000"
  TS_FMT_DATETIME           = "2006-01-02 15:04:05"
  TS_FMT_DATETIME_MILLI     = "2006-01-02 15:04:05.000"
  TS_FMT_DATETIME_MICRO     = "2006-01-02 15:04:05.000000"
  TS_FMT_TIME_TZ            = "15:04:05-0700"
  TS_FMT_TIME_TZ_MILLI      = "15:04:05.000-0700"
  TS_FMT_TIME_TZ_MICRO      = "15:04:05.000000-0700"
  TS_FMT_DATETIME_TZ        = "2006-01-02 15:04:05-0700"
  TS_FMT_DATETIME_TZ_MILLI  = "2006-01-02 15:04:05.000-0700"
  TS_FMT_DATETIME_TZ_MICRO  = "2006-01-02 15:04:05.000000-0700"
)

type outputMap  map[int]io.Writer

type Logger struct {
  verbosity     int
  output        outputMap
  overrideStack []bool
  prefixTS      bool
  prefixLevel   bool
  prefixCaller  bool
  fmtTimestamp  string
}

var (
  // The global logger object
  logger  *Logger = NewLogger()
)


// NewLogger returns a new logger object.
func NewLogger() *Logger {
  l := Logger{
    verbosity: INFO,    // Setting reasonable default log level
    output: make(outputMap),  // Maps log levels to Writer objects, such as os.Stdout or a file
    overrideStack: make([]bool, 0, 8),  // Temporarily stores current prefix visibility settings
    prefixTS: false,
    prefixLevel: false,
    prefixCaller: false,
    fmtTimestamp: TS_FMT_TIME_MILLI,
  }
  l.output[LOG]       = os.Stdout
  l.output[INFO]      = os.Stdout
  l.output[WARN]      = os.Stderr
  l.output[ERROR]     = os.Stderr
  l.output[CRITICAL]  = os.Stderr
  return &l
}


// Global returns the global logger object. It is rarely needed to call this function directly, 
// since every Logger function comes with a global counterpart of the same name.
func Global() *Logger {
  return logger
}


// GetVerbosity returns the current verbosity level.
// Only log messages of the current verbosity level or higher will be logged.
func (l *Logger) GetVerbosity() int {
  return l.verbosity
}

// Global logger: GetVerbosity returns the current verbosity level.
// Only log messages of the current verbosity level or higher will be logged.
func GetVerbosity() int { return Global().GetVerbosity() }


// SetVerbosity sets the current verbosity level.
//
// Log messages of the current verbosity level or higher will be logged.
// Supported levels in increasing order of importance: LOG, INFO, WARN, ERROR and CRITICAL.
func (l *Logger) SetVerbosity(level int) {
  if level < LOG { level = LOG }
  if level > CRITICAL { level = CRITICAL }
  l.verbosity = level
}

// Global logger: SetVerbosity sets the current verbosity level.
//
// Log messages of the current verbosity level or higher will be logged.
// Supported levels in increasing order of importance: LOG, INFO, WARN, ERROR and CRITICAL.
func SetVerbosity(level int) { Global().SetVerbosity(level) }


// IncreaseVerbosity increases the current verbosity by one level.
// Does nothing if highest level "CRITICAL" is already set. Returns the new verbosity level.
func (l *Logger) IncreaseVerbosity() int {
  if l.verbosity < CRITICAL { l.verbosity++ }
  return l.verbosity
}

// Global logger: IncreaseVerbosity increases the current verbosity by one level.
// Does nothing if highest level "CRITICAL" is already set. Returns the new verbosity level.
func IncreaseVerbosity() int { return Global().IncreaseVerbosity() }


// DecreaseVerbosity decreases the current verbosity by one level.
// Does nothing if lowest level "LOG" is already set. Returns the new verbosity level.
func (l *Logger) DecreaseVerbosity() int {
  if l.verbosity > LOG { l.verbosity-- }
  return l.verbosity
}

// Global logger: DecreaseVerbosity decreases the current verbosity by one level.
// Does nothing if lowest level "LOG" is already set. Returns the new verbosity level.
func DecreaseVerbosity() int { return Global().DecreaseVerbosity() }


// GetPrefixTimestamp returns whether log messages are prefixed by the current timestamp.
func (l *Logger) GetPrefixTimestamp() bool {
  return l.prefixTS
}

// Global logger: GetPrefixTimestamp returns whether log messages are prefixed by the current timestamp.
func GetPrefixTimestamp() bool { return Global().GetPrefixTimestamp() }


// SetPrefixTimestamp defines whether log messages should be prefixed by the current timestamp.
func (l *Logger) SetPrefixTimestamp(set bool) {
  l.prefixTS = set
}

// Global logger: SetPrefixTimestamp defines whether log messages should be prefixed by the current timestamp.
func SetPrefixTimestamp(set bool) { Global().SetPrefixTimestamp(set) }


// GetTimestampFormat returns the format string for the timestamp prefix.
func (l *Logger) GetTimestampFormat() string {
  return l.fmtTimestamp
}

// Global logger: GetTimestampFormat returns the format string for the timestamp prefix.
func GetTimestampFormat() string { return Global().GetTimestampFormat() }


// SetTimestampFormat sets a new format string for the timestamp prefix.
//
// Use either the TS_FMT_xxx constants, predefined constants from Golang's time package
// or define a custom format on your own. Format description: https://golang.org/pkg/time/#pkg-constants
func (l *Logger) SetTimestampFormat(format string) {
  l.fmtTimestamp = format
}

// Global logger: SetTimestampFormat sets a new format string for the timestamp prefix.
//
// Use either the TS_FMT_xxx constants, predefined constants from Golang's time package
// or define a custom format on your own. Format description: https://golang.org/pkg/time/#pkg-constants
func SetTimestampFormat(format string) { Global().SetTimestampFormat(format) }


// GetPrefixCaller returns whether log messags are prefixed by name and line number of the calling function.
func (l *Logger) GetPrefixCaller() bool {
  return l.prefixCaller
}

// Global logger: GetPrefixCaller returns whether log messags are prefixed by name and line number of the calling function.
func GetPrefixCaller() bool { return Global().GetPrefixCaller() }


// SetPrefixCaller defines whether log messages should be prefixed by name and line number of the calling function.
func (l *Logger) SetPrefixCaller(set bool) {
  l.prefixCaller = set
}

// Global logger: SetPrefixCaller defines whether log messages should be prefixed by name and line number of the calling function.
func SetPrefixCaller(set bool) { Global().SetPrefixCaller(set) }


// GetPrefixLevel returns whether log messages are prefixed by a symbolic name of their level.
func (l *Logger) GetPrefixLevel() bool {
  return l.prefixLevel
}

// Global logger: GetPrefixLevel returns whether log messages are prefixed by a symbolic name of their level.
func GetPrefixLevel() bool { return Global().GetPrefixLevel() }


// SetPrefixLevel defines whether log messages should be prefixed by a symbolic name of their level.
func (l *Logger) SetPrefixLevel(set bool) {
  l.prefixLevel = set
}

// Global logger: SetPrefixLevel defines whether log messages should be prefixed by a symbolic name of their level.
func SetPrefixLevel(set bool) { Global().SetPrefixLevel(set) }


// GetOutput returns the Writer object for messages of the given level.
//
// By default LOG and INFO are written to os.Stdout. WARN, ERROR and CRITICAL are written to os.Stderr.
// Returns nil for unsupported log levels.
func (l *Logger) GetOutput(level int) io.Writer {
  if level < LOG || level > CRITICAL { return nil }
  return l.output[level]
}

// Global logger: GetOutput returns the Writer object for messages of the given level.
//
// By default LOG and INFO are written to os.Stdout. WARN, ERROR and CRITICAL are written to os.Stderr.
// Returns nil for unsupported log levels.
func GetOutput(level int) io.Writer { return Global().GetOutput(level) }


// SetOutput redirects log messages of the given level to the specified Writer object.
//
// By default LOG and INFO are written to os.Stdout. WARN, ERROR and CRITICAL are written to os.Stderr.
// Unsupported log levels are ignored. Specify a nil Writer to restore default output channel for the given level.
// The caller is responsible to close the specified Writer after it is no longer used.
func (l *Logger) SetOutput(level int, writer io.Writer) {
  if level < LOG || level > CRITICAL { return }
  if writer == nil {
    switch level {
      case WARN:
      case ERROR:
      case CRITICAL:
        writer = os.Stderr
      default:
        writer = os.Stdout
    }
  }
  l.output[level] = writer
}

// Global logger: SetOutput redirects log messages of the given level to the specified Writer object.
//
// By default LOG and INFO are written to os.Stdout. WARN, ERROR and CRITICAL are written to os.Stderr.
// Unsupported log levels are ignored. Specify a nil Writer to restore default output channel for the given level.
// The caller is responsible to close the specified Writer after it is no longer used.
func SetOutput(level int, writer io.Writer) { Global().SetOutput(level, writer) }


// OverridePrefix overrides current log prefix settings only for the next call of a log output function. 
//
// Log output functions are: Log/Logf/Logln, Info/Infof/Infoln, Warn/Warnf/Warnln, Error/Errorf/Errorln and
// Critical/Criticalf/Criticalln. Returns the Logger object to allow chaining function calls.
//
// Important: Multiple calls of this function are cumulative.
func (l *Logger) OverridePrefix(showTimestamp, showCaller, showLevel bool) *Logger {
  l.pushOverride(showTimestamp, showCaller, showLevel)
  return l
}

// Global logger: OverridePrefix overrides current log prefix settings only for the next call of a log output function. 
//
// Log output functions are: Log/Logf/Logln, Info/Infof/Infoln, Warn/Warnf/Warnln, Error/Errorf/Errorln and 
// Critical/Criticalf/Criticalln. Returns the global Logger object to allow chaining function calls.
//
// Important: Multiple calls of this function are cumulative.
func OverridePrefix(showTimestamp, showCaller, showLevel bool) *Logger { return Global().OverridePrefix(showTimestamp, showCaller, showLevel) }


// Log prints the LOG message if current verbosity level is set to LOG.
func (l *Logger) Log(msg string) {
  l.logf(l.getOutput(LOG), LOG, msg)
}

// Global logger: Log prints the message if current verbosity level is set to LOG.
func Log(msg string) { Global().Log(msg) }

// Info prints the message if current verbosity level is set to INFO or lower.
func (l *Logger) Info(msg string) {
  l.logf(l.getOutput(INFO), INFO, msg)
}

// Global logger: Info prints the message if current verbosity level is set to INFO or lower.
func Info(msg string) { Global().Info(msg) }

// Warn prints the message if current verbosity level is set to WARN or lower.
func (l *Logger) Warn(msg string) {
  l.logf(l.getOutput(WARN), WARN, msg)
}

// Global logger: Warn prints the message if current verbosity level is set to WARN or lower.
func Warn(msg string) { Global().Warn(msg) }

// Error prints the message if current verbosity level is set to ERROR or lower.
func (l *Logger) Error(msg string) {
  l.logf(l.getOutput(ERROR), ERROR, msg)
}

// Global logger: Error prints the message if current verbosity level is set to ERROR or lower.
func Error(msg string) { Global().Error(msg) }

// Critical invokes a panic with the specified message.
func (l *Logger) Critical(msg string) {
  l.logf(l.getOutput(CRITICAL), CRITICAL, msg)
}

// Global logger: Critical invokes a panic with the specified message.
func Critical(msg string) { Global().Critical(msg) }


// Logf prints the formatted string if current verbosity level is set to LOG.
func (l *Logger) Logf(format string, a ...interface{}) {
  l.logf(l.getOutput(LOG), LOG, format, a...)
}

// Global logger: Logf prints the formatted string if current verbosity level is set to LOG.
func Logf(format string, a ...interface{}) { Global().Logf(format, a...) }

// Infof prints the formatted string if current verbosity level is set to INFO or lower.
func (l *Logger) Infof(format string, a ...interface{}) {
  l.logf(l.getOutput(INFO), INFO, format, a...)
}

// Global logger: Infof prints the formatted string if current verbosity level is set to INFO or lower.
func Infof(format string, a ...interface{}) { Global().Infof(format, a...) }

// Warnf prints the formatted string if current verbosity level is set to WARN or lower.
func (l *Logger) Warnf(format string, a ...interface{}) {
  l.logf(l.getOutput(WARN), WARN, format, a...)
}

// Global logger: Warnf prints the formatted string if current verbosity level is set to WARN or lower.
func Warnf(format string, a ...interface{}) { Global().Warnf(format, a...) }

// Errorf prints the formatted string if current verbosity level is set to ERROR or lower.
func (l *Logger) Errorf(format string, a ...interface{}) {
  l.logf(l.getOutput(ERROR), ERROR, format, a...)
}

// Global logger: Errorf prints the formatted string if current verbosity level is set to ERROR or lower.
func Errorf(format string, a ...interface{}) { Global().Errorf(format, a...) }

// Criticalf invokes a panic with the formatted string.
func (l *Logger) Criticalf(format string, a ...interface{}) {
  l.logf(l.getOutput(CRITICAL), CRITICAL, format, a...)
}

// Global logger: Criticalf invokes a panic with the formatted string.
func Criticalf(format string, a ...interface{}) { Global().Criticalf(format, a...) }


// Logln prints the message and a newline if current verbosity is set to LOG.
func (l *Logger) Logln(msg string) {
  l.logf(l.getOutput(LOG), LOG, "%s\n", msg)
}

// Global logger: Logln prints the message and a newline if current verbosity is set to LOG.
func Logln(msg string) { Global().Logln(msg) }

// Infoln prints the message and a newline if current verbosity is set to INFO or lower.
func (l *Logger) Infoln(msg string) {
  l.logf(l.getOutput(INFO), INFO, "%s\n", msg)
}

// Global logger: Infoln prints the message and a newline if current verbosity is set to INFO or lower.
func Infoln(msg string) { Global().Infoln(msg) }

// Warnln prints the message and a newline if current verbosity is set to WARN or lower.
func (l *Logger) Warnln(msg string) {
  l.logf(l.getOutput(WARN), WARN, "%s\n", msg)
}

// Global logger: Warnln prints the message and a newline if current verbosity is set to WARN or lower.
func Warnln(msg string) { Global().Warnln(msg) }

// Errorln prints the message and a newline if current verbosity is set to ERROR or lower.
func (l *Logger) Errorln(msg string) {
  l.logf(l.getOutput(ERROR), ERROR, "%s\n", msg)
}

// Global logger: Errorln prints the message and a newline if current verbosity is set to ERROR or lower.
func Errorln(msg string) { Global().Errorln(msg) }

// Criticalln invokes a panic with the message and a newline.
func (l *Logger) Criticalln(msg string) {
  l.logf(l.getOutput(CRITICAL), CRITICAL, "%s\n", msg)
}

// Global logger: Criticalln invokes a panic with the message and a newline.
func Criticalln(msg string) { Global().Criticalln(msg) }


// LogProgressDot is a specialized version of the function LogProgress.
//
// It prints zero, one or more instances of "dot" (.) characters based on the given arguments if current
// verbosity level is set to LOG.
func (l *Logger) LogProgressDot(cur, max, progressMax int) {
  l.LogProgress(cur, max, progressMax, ".")
}

// Global logger: LogProgressDot is a specialized version of the function LogProgress.
// It prints zero, one or more instances of "dot" (.) characters based on the given arguments if current
// verbosity level is set to LOG.
func LogProgressDot(cur, max, progressMax int) { Global().LogProgressDot(cur, max, progressMax) }


// LogProgress prints zero, one or more instances of "symbol" when "cur" indicates a new position within the
// range 0 to "max" (exclusive).
//
// "progressMax" specifies the total number of "symbol"s to print when calling this function 
// for each position from 0 to max-1. Does nothing if current verbosity level is higher than LOG.
func (l *Logger) LogProgress(cur, max, progressMax int, symbol string) {
  s := Progress(cur, max, progressMax, symbol)
  if len(s) > 0 {
    l.pushOverride(false, false, false)
    l.logf(l.getOutput(LOG), LOG, s)
  }
}

// Global logger: LogProgress prints zero, one or more instances of "symbol" when "cur" indicates a new position within
// the range 0 to "max" (exclusive).
//
// "progressMax" specifies the total number of "symbol"s to print when calling this function 
// for each position from 0 to max-1. Does nothing if current verbosity level is higher than LOG.
func LogProgress(cur, max, progressMax int, symbol string) { Global().LogProgress(cur, max, progressMax, symbol) }


// InfoProgressDot is a specialized version of the function LogProgress.
//
// It prints zero, one or more instances of "dot" (.) characters based on the given arguments if current verbosity level
// is set to INFO or lower.
func (l *Logger) InfoProgressDot(cur, max, progressMax int) {
  l.InfoProgress(cur, max, progressMax, ".")
}

// Global logger: InfoProgressDot is a specialized version of the function LogProgress.
//
// It prints zero, one or more instances of "dot" (.) characters based on the given arguments if current verbosity level
// is set to INFO or lower.
func InfoProgressDot(cur, max, progressMax int) { Global().InfoProgressDot(cur, max, progressMax) }


// InfoProgress prints zero, one or more instances of "symbol" when "cur" indicates a new position within the range 
// 0 to "max" (exclusive).
//
// "progressMax" specifies the total number of "symbol"s to print when calling this function 
// for each position from 0 to max-1. Does nothing if current verbosity level is higher than INFO.
func (l *Logger) InfoProgress(cur, max, progressMax int, symbol string) {
  s := Progress(cur, max, progressMax, symbol)
  if len(s) > 0 {
    l.pushOverride(false, false, false)
    l.logf(l.getOutput(INFO), INFO, s)
  }
}

// Global logger: InfoProgress prints zero, one or more instances of "symbol" when "cur" indicates a new position within 
// the range 0 to "max" (exclusive).
//
// "progressMax" specifies the total number of "symbol"s to print when calling this function 
// for each position from 0 to max-1. Does nothing if current verbosity level is higher than INFO.
func InfoProgress(cur, max, progressMax int, symbol string) { Global().InfoProgress(cur, max, progressMax, symbol) }


// WarnProgressDot is a specialized version of the function LogProgress.
//
// It prints zero, one or more instances of "dot" (.) characters based on the given arguments if current verbosity level
// is set to WARN or lower.
func (l *Logger) WarnProgressDot(cur, max, progressMax int) {
  l.WarnProgress(cur, max, progressMax, ".")
}

// Global logger: WarnProgressDot is a specialized version of the function LogProgress.
//
// It prints zero, one or more instances of "dot" (.) characters based on the given arguments if current verbosity level
// is set to WARN or lower.
func WarnProgressDot(cur, max, progressMax int) { Global().WarnProgressDot(cur, max, progressMax) }


// WarnProgress prints zero, one or more instances of "symbol" when "cur" indicates a new position within the range 
// 0 to "max" (exclusive).
//
// "progressMax" specifies the total number of "symbol"s to print when calling this function 
// for each position from 0 to max-1. Does nothing if current verbosity level is higher than WARN.
func (l *Logger) WarnProgress(cur, max, progressMax int, symbol string) {
  s := Progress(cur, max, progressMax, symbol)
  if len(s) > 0 {
    l.pushOverride(false, false, false)
    l.logf(l.getOutput(WARN), WARN, s)
  }
}

// Global logger: WarnProgress prints zero, one or more instances of "symbol" when "cur" indicates a new position within 
// the range 0 to "max" (exclusive).
//
// "progressMax" specifies the total number of "symbol"s to print when calling this function 
// for each position from 0 to max-1. Does nothing if current verbosity level is higher than WARN.
func WarnProgress(cur, max, progressMax int, symbol string) { Global().WarnProgress(cur, max, progressMax, symbol) }


// ErrorProgressDot is a specialized version of the function LogProgress.
//
// It prints zero, one or more instances of "dot" (.) characters based on the given arguments if current verbosity level
// is set to ERROR or lower.
func (l *Logger) ErrorProgressDot(cur, max, progressMax int) {
  l.ErrorProgress(cur, max, progressMax, ".")
}

// Global logger: ErrorProgressDot is a specialized version of the function LogProgress.
//
// It prints zero, one or more instances of "dot" (.) characters based on the given arguments if current verbosity level
// is set to ERROR or lower.
func ErrorProgressDot(cur, max, progressMax int) { Global().ErrorProgressDot(cur, max, progressMax) }


// ErrorProgress prints zero, one or more instances of "symbol" when "cur" indicates a new position within the range 
// 0 to "max" (exclusive).
//
// "progressMax" specifies the total number of "symbol"s to print when calling this function for each position from
// 0 to max-1. Does nothing if current verbosity level is higher than ERROR.
func (l *Logger) ErrorProgress(cur, max, progressMax int, symbol string) {
  s := Progress(cur, max, progressMax, symbol)
  if len(s) > 0 {
    l.pushOverride(false, false, false)
    l.logf(l.getOutput(ERROR), ERROR, s)
  }
}

// Global logger: ErrorProgress prints zero, one or more instances of "symbol" when "cur" indicates a new position within 
// the range 0 to "max" (exclusive).
//
// "progressMax" specifies the total number of "symbol"s to print when calling this function for each position from 0 to max-1.
// Does nothing if current verbosity level is higher than ERROR.
func ErrorProgress(cur, max, progressMax int, symbol string) { Global().ErrorProgress(cur, max, progressMax, symbol) }


// ProgressDot is the specialized version of the function Progress.
// It returns a sequence of zero, one or more instances of "dot" (.) based on the given arguments.
func ProgressDot(cur, max, progressMax int) string {
  return Progress(cur, max, progressMax, ".")
}

// Progress returns a sequence of zero, one or more instances of "symbol" when "cur" indicates a new position within 
// the range 0 to "max" (exclusive).
//
// "progressMax" specifies the total number of "symbol"s to return when calling this function for each position from
// 0 to max-1. It is not associated with any Logger objects.
func Progress(cur, max, progressMax int, symbol string) string {
  s := ""
  if max < 1 || progressMax < 1 { return s }
  if cur < 0 || cur >= max { return s }
  i1 := cur * progressMax / max
  i2 := (cur + 1) * progressMax / max
  for ; i2 > i1; i1++ {
    s += symbol
  }
  return s
}


// Used internally. Handles writing log messages.
func (l *Logger) logf(w io.Writer, level int, format string, a ...interface{}) {
  if level > CRITICAL { level = CRITICAL }

  if level >= l.verbosity {
    if level == CRITICAL {
      panic(fmt.Sprintf(format, a...))
    }

    if w == nil { w = l.getOutput(level) }
    prefix := l.getLogPrefix(level)
    msg := fmt.Sprintf(format, a...)
    _, err := fmt.Fprintf(w, "%s%s", prefix, msg)
    if err != nil {
      l.logf(os.Stderr, ERROR, "logging.Logf(): %v", err)
    }
  }

  l.popOverride()
}


// Used internally. Returns the Writer object of the specified log level.
func (l *Logger) getOutput(level int) io.Writer {
  if level < LOG { level = LOG }
  if level > CRITICAL { level = CRITICAL }
  return l.output[level]
}


// Used internally. Returns a log prefix string.
func (l *Logger) getLogPrefix(level int) string {
  var prefix strings.Builder
  if l.prefixTS {
    t := time.Now()
    prefix.WriteString(t.Format(l.fmtTimestamp))
    prefix.WriteString(" ")
  }
  if (l.prefixCaller) {
    pc := make([]uintptr, 16)
    cnt := runtime.Callers(1, pc) // skip runtime.Callers from calling stack
    if cnt > 0 {
      // determine key string that should not be present in the name string of the calling function
      f := runtime.FuncForPC(pc[0])
      key := f.Name()
      pos := strings.Index(key, ".(*Logger)")
      if pos >= 0 {
        key = key[:pos]
      }
      // first function not matching key is our prime candidate
      for i := 1; i < cnt; i++ {
        f := runtime.FuncForPC(pc[i])
        name := f.Name()
        if strings.Index(name, key) < 0 {
          _, line := f.FileLine(pc[i])
          prefix.WriteString(fmt.Sprintf("%s:%d ", name, line))
          break
        }
      }
    }
  }
  if l.prefixLevel {
    prefix.WriteString(l.getLevelString(level))
    prefix.WriteString(" ")
  }
  return prefix.String()
}


// Used internally. Returns a textual representation of the given log level.
func (l *Logger) getLevelString(level int) string {
  var s string
  if level < LOG { level = LOG }
  if level > CRITICAL { level = CRITICAL }
  switch level {
    case LOG:       s = "LOG "
    case INFO:      s = "INFO"
    case WARN:      s = "WARN"
    case ERROR:     s = "ERRO"
    case CRITICAL:  s = "CRIT"
    default:        s = "    "
  }
  return s
}


// Used internally. Pushes given log prefix options to the stack.
func (l *Logger) pushOverride(ts, caller, level bool) {
  l.overrideStack = append(l.overrideStack, l.prefixTS, l.prefixCaller, l.prefixLevel)
  l.prefixTS = ts
  l.prefixCaller = caller
  l.prefixLevel = level
}


// Used internally. Restores most recent log prefix overrides. Does nothing if no overrides were stored.
func (l *Logger) popOverride() {
  if len(l.overrideStack) > 2 {
    idx := len(l.overrideStack) - 3
    l.prefixTS = l.overrideStack[idx]
    l.prefixCaller = l.overrideStack[idx+1]
    l.prefixLevel = l.overrideStack[idx+2]
    l.overrideStack = l.overrideStack[:idx]
  }
}
