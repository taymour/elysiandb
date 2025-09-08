package log_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	pkglog "github.com/taymour/elysiandb/internal/log"
)

func captureStdout(t *testing.T) (finish func() string) {
	t.Helper()

	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	type result struct {
		s   string
		err error
	}
	done := make(chan result, 1)

	go func() {
		b, err := io.ReadAll(r)
		done <- result{s: string(b), err: err}
	}()

	return func() string {
		_ = w.Close()
		os.Stdout = orig
		res := <-done
		_ = r.Close()
		if res.err != nil && !errors.Is(res.err, os.ErrClosed) {
			t.Logf("capture: read error: %v", res.err)
		}
		return res.s
	}
}

func waitLogs() { time.Sleep(20 * time.Millisecond) }

func TestAsyncEnqueueAndWriteLogs(t *testing.T) {
	finish := captureStdout(t)

	pkglog.Info("hello-info")
	pkglog.Success("hello-success")
	pkglog.Warn("hello-warn")
	pkglog.Error("hello-error")
	pkglog.Debug("hello-debug")

	waitLogs()
	pkglog.WriteLogs()

	out := finish()

	for _, frag := range []string{
		"hello-info", "hello-success", "hello-warn", "hello-error", "hello-debug",
	} {
		if !strings.Contains(out, frag) {
			t.Fatalf("flushed output missing %q.\nGot:\n%s", frag, out)
		}
	}

	finish2 := captureStdout(t)
	pkglog.WriteLogs()
	out2 := finish2()
	if strings.TrimSpace(out2) != "" {
		t.Fatalf("expected empty output after second WriteLogs, got: %q", out2)
	}
}

func TestDirectInfoPrintsImmediately(t *testing.T) {
	finish := captureStdout(t)

	pkglog.DirectInfo("direct-message")

	out := finish()

	if !strings.Contains(out, "INFO") {
		t.Fatalf("expected output to contain 'INFO', got:\n%s", out)
	}
	if !strings.Contains(out, "direct-message") {
		t.Fatalf("expected output to contain 'direct-message', got:\n%s", out)
	}
}

func TestFatalExitsWithNonZero(t *testing.T) {
	if os.Getenv("RUN_FATAL_HELPER") == "1" {
		pkglog.Fatal("boom", fmt.Errorf("kaboom"))
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run", t.Name())
	cmd.Env = append(os.Environ(), "RUN_FATAL_HELPER=1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-nil error (process should have exited), got nil")
	}
	if exitErr, ok := err.(*exec.ExitError); !ok {
		t.Fatalf("expected *exec.ExitError, got %T (%v)", err, err)
	} else if exitErr.ExitCode() == 0 {
		t.Fatalf("expected non-zero exit code, got 0")
	}
}
