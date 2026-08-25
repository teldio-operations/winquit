package winquit

import (
	"log/slog"
	"sync/atomic"
)

var currentLogger atomic.Pointer[slog.Logger]

// SetLogger directs the debug messages of this package to the passed logger. A
// nil logger turns the messages off.
//
// Without a call to this function, the package uses the logger returned by
// slog.Default() at the time of the message. That logger discards debug
// messages unless the application installs a handler which accepts
// slog.LevelDebug, so the package stays quiet by default.
func SetLogger(l *slog.Logger) {
	if l == nil {
		l = slog.New(slog.DiscardHandler)
	}
	currentLogger.Store(l)
}

func logger() *slog.Logger {
	if l := currentLogger.Load(); l != nil {
		return l
	}
	return slog.Default()
}
