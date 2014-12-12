package logging

import (
	"os"
	"os/signal"
)

// Call Logger.Reopen() when a signal is received
// Used for log rotation; the log rotator should signal to this go application
// before destroying the old log file.
// This function loops forever; run as a goroutine
func (l *Logger) ReopenOnSignalLoop(s os.Signal) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, s)
	for {
		if _, ok := <-sigs; !ok {
			panic("Channel closed unexpectedly")
		}
		if err := l.Reopen(); err != nil {
			panic(err)
		}
	}
}
