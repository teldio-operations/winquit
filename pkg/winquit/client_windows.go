package winquit

import (
	"os"
	"time"

	"github.com/teldio-operations/winquit/pkg/winquit/win32"
)

func requestQuit(pid int) error {
	threads, err := win32.GetProcThreads(uint32(pid))
	if err != nil {
		return err
	}

	for _, thread := range threads {
		logger().Debug("Closing windows on thread", "thread", thread)
		win32.CloseThreadWindows(uint32(thread))
	}

	return nil
}

func quitProcess(pid int, waitNicely time.Duration) error {
	_ = RequestQuit(pid)

	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil
	}

	done := make(chan bool)

	go func() {
		proc.Wait()
		done <- true
	}()

	select {
	case <-done:
		return nil
	case <-time.After(waitNicely):
	}

	return proc.Kill()
}
