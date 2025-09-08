package tcp

import (
	"bufio"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/taymour/elysiandb/internal/boot"
	"github.com/taymour/elysiandb/internal/configuration"
	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/storage"
)

const tcpAddr = "127.0.0.1:8088"

func waitTCPUp(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		c, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if err == nil {
			c.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return err
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func TestTCP_PING_SET_GET__SAVE__RESET(t *testing.T) {
	tmp := t.TempDir()
	globals.SetConfig(&configuration.Config{
		Store: configuration.StoreConfig{
			Folder: tmp,
			Shards: 8,
		},
	})
	storage.LoadDB()
	boot.BootSaver()
	boot.BootExpirationHandler()
	boot.BootLogger()

	c, err := net.DialTimeout("tcp", tcpAddr, 150*time.Millisecond)
	if err == nil {
		_ = c.Close()
	} else {
		go boot.InitTCP()
		if err := waitTCPUp(tcpAddr, 2*time.Second); err != nil {
			var ne net.Error
			if errors.As(err, &ne) {
				t.Skipf("skipping TCP test: cannot connect to %s (%v)", tcpAddr, err)
			}
			t.Skipf("skipping TCP test: %v", err)
		}
	}

	c, err = net.DialTimeout("tcp", tcpAddr, 2*time.Second)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer c.Close()
	r := bufio.NewReader(c)

	write := func(s string) {
		_ = c.SetWriteDeadline(time.Now().Add(1 * time.Second))
		if _, err := c.Write([]byte(s + "\n")); err != nil {
			t.Fatalf("write %q: %v", s, err)
		}
	}
	readLine := func() string {
		_ = c.SetReadDeadline(time.Now().Add(2 * time.Second))
		l, err := r.ReadString('\n')
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		return l[:len(l)-1]
	}

	write("PING")
	if got := readLine(); got != "PONG" {
		t.Fatalf("want PONG, got %q", got)
	}

	write("SET foo hello")
	if got := readLine(); got != "OK" {
		t.Fatalf("want OK, got %q", got)
	}

	write("GET foo")
	if got := readLine(); got != "hello" {
		t.Fatalf("want hello, got %q", got)
	}

	write("SAVE")
	if got := readLine(); got != "OK" {
		t.Fatalf("want OK, got %q", got)
	}

	write("GET foo")
	if got := readLine(); got != "hello" {
		t.Fatalf("want hello, got %q", got)
	}

	write("RESET")
	if got := readLine(); got != "OK" {
		t.Fatalf("want OK, got %q", got)
	}

	write("GET foo")
	if got := readLine(); got != "Key not found" {
		t.Fatalf("want not found, got %q", got)
	}
}
