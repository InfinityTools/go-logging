package logging

import (
  "fmt"
  "testing"
)

func TestLogging1(t *testing.T) {
  defer func(){ recover(); fmt.Println("Recovered from critical error.") }()
  l := NewLogger()
  l.SetVerbosity(INFO)
  l.SetPrefixLevel(true)
  l.SetPrefixTimestamp(true)
  l.SetPrefixCaller(true)
  l.Logln("This debug message should not be visible.\n")
  l.Info("A simple message with manual newline.\n")
  l.Warnf("A formatted message at log level %d.\n", WARN)
  l.Errorln("A simple message with auto-newline.")
  l.Info("A progress bar")
  l.InfoProgressDot(0, 1, 79 - 14)
  l.OverridePrefix(false, false, false).Infoln("")

  l.Criticalln("This is a critical error.")
}
