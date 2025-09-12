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
			_ = c.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return err
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func TestTCP_PING_SET_MGET_GET__WILDCARD__SAVE__RESET(t *testing.T) {
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

	if c, err := net.DialTimeout("tcp", tcpAddr, 150*time.Millisecond); err == nil {
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

	c, err := net.DialTimeout("tcp", tcpAddr, 2*time.Second)
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
	readN := func(n int) []string {
		out := make([]string, 0, n)
		for i := 0; i < n; i++ {
			out = append(out, readLine())
		}
		return out
	}
	contains := func(list []string, s string) bool {
		for _, x := range list {
			if x == s {
				return true
			}
		}
		return false
	}

	write("PING")
	if got := readLine(); got != "PONG" {
		t.Fatalf("want PONG, got %q", got)
	}

	write("SET foo hello")
	if got := readLine(); got != "OK" {
		t.Fatalf("want OK, got %q", got)
	}
	write("SET bar bat")
	if got := readLine(); got != "OK" {
		t.Fatalf("want OK, got %q", got)
	}
	write("SET user:1 alice")
	if got := readLine(); got != "OK" {
		t.Fatalf("want OK, got %q", got)
	}
	write("SET user:2 bob")
	if got := readLine(); got != "OK" {
		t.Fatalf("want OK, got %q", got)
	}

	write("MGET foo bar baz")
	got := readN(3)
	exp := []string{"hello", "bat", "baz=not found"}
	for i, want := range exp {
		if got[i] != want {
			t.Fatalf("MGET[%d]: want %q, got %q (all=%v)", i, want, got[i], got)
		}
	}

	write("GET foo")
	if got := readLine(); got != "foo=hello" {
		t.Fatalf("want %q, got %q", "foo=hello", got)
	}

	write("GET user:*")
	lines := readN(2)
	if !(contains(lines, "user:1=alice") && contains(lines, "user:2=bob")) {
		t.Fatalf("GET user:* expected both user:1=alice and user:2=bob, got %v", lines)
	}

	write("MGET foo user:* zoo")
	mixed := readN(4)
	if mixed[0] != "hello" {
		t.Fatalf("MGET mixed[0]: want %q, got %q (all=%v)", "hello", mixed[0], mixed)
	}
	if !(contains(mixed[1:3], "user:1=alice") && contains(mixed[1:3], "user:2=bob")) {
		t.Fatalf("MGET mixed user block mismatch, got %v", mixed[1:3])
	}
	if mixed[3] != "zoo=not found" {
		t.Fatalf("MGET mixed last line: want %q, got %q (all=%v)", "zoo=not found", mixed[3], mixed)
	}

	write("SAVE")
	if got := readLine(); got != "OK" {
		t.Fatalf("want OK, got %q", got)
	}

	write("GET foo")
	if got := readLine(); got != "foo=hello" {
		t.Fatalf("want %q, got %q", "foo=hello", got)
	}

	write("RESET")
	if got := readLine(); got != "OK" {
		t.Fatalf("want OK, got %q", got)
	}

	write("GET foo")
	if got := readLine(); got != "foo=not found" {
		t.Fatalf("want %q, got %q", "foo=not found", got)
	}
}
