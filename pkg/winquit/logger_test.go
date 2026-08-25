package winquit

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func debugLogger(w *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func saveLoggerState(t *testing.T) {
	t.Helper()

	previousDefault := slog.Default()
	previousLogger := currentLogger.Load()

	t.Cleanup(func() {
		slog.SetDefault(previousDefault)
		currentLogger.Store(previousLogger)
	})
}

func TestLoggerFallsBackToSlogDefault(t *testing.T) {
	saveLoggerState(t)

	var buf bytes.Buffer
	currentLogger.Store(nil)
	slog.SetDefault(debugLogger(&buf))

	logger().Debug("fallback message")

	if !strings.Contains(buf.String(), "fallback message") {
		t.Fatalf("expected the slog default logger to receive the message, got %q", buf.String())
	}
}

func TestSetLoggerRoutesMessages(t *testing.T) {
	saveLoggerState(t)

	var def, own bytes.Buffer
	slog.SetDefault(debugLogger(&def))
	SetLogger(debugLogger(&own))

	logger().Debug("routed message")

	if !strings.Contains(own.String(), "routed message") {
		t.Fatalf("expected the passed logger to receive the message, got %q", own.String())
	}
	if def.Len() != 0 {
		t.Fatalf("expected the slog default logger to receive nothing, got %q", def.String())
	}
}

func TestSetLoggerNilTurnsMessagesOff(t *testing.T) {
	saveLoggerState(t)

	var def bytes.Buffer
	slog.SetDefault(debugLogger(&def))
	SetLogger(nil)

	logger().Debug("dropped message")
	logger().Info("dropped message")
	logger().Error("dropped message")

	if def.Len() != 0 {
		t.Fatalf("expected no output, got %q", def.String())
	}
	if logger().Enabled(context.Background(), slog.LevelError) {
		t.Fatal("expected a nil logger to report every level as disabled")
	}
}
