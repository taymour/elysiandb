package e2e

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

type getEntry struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
}

type multiGetEntry struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
}

func TestMultiplePUTAndMGET(t *testing.T) {
	client, stop := startTestServer(t)
	defer stop()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodPut)
	req.SetRequestURI("http://test/kv/foo")
	req.SetBodyString("bar")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("PUT failed: %v", err)
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusNoContent {
		t.Fatalf("expected 204, got %d", sc)
	}

	req.Reset()
	resp.Reset()

	req.Header.SetMethod(fasthttp.MethodPut)
	req.SetRequestURI("http://test/kv/baz")
	req.SetBodyString("bat")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("PUT failed: %v", err)
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusNoContent {
		t.Fatalf("expected 204, got %d", sc)
	}

	req.Reset()
	resp.Reset()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI("http://test/kv/mget?keys=foo,baz,qux")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusOK {
		t.Fatalf("expected 200, got %d", sc)
	}

	var arr []multiGetEntry
	if err := json.Unmarshal(resp.Body(), &arr); err != nil {
		t.Fatalf("invalid JSON from /kv/mget: %v (body=%q)", err, resp.Body())
	}
	if len(arr) != 3 {
		t.Fatalf("expected 3 results, got %d (%v)", len(arr), arr)
	}
	if arr[0].Key != "foo" || arr[0].Value == nil || *arr[0].Value != "bar" {
		t.Fatalf("unexpected entry[0]: %+v", arr[0])
	}
	if arr[1].Key != "baz" || arr[1].Value == nil || *arr[1].Value != "bat" {
		t.Fatalf("unexpected entry[1]: %+v", arr[1])
	}
	if arr[2].Key != "qux" || arr[2].Value != nil {
		t.Fatalf("unexpected entry[2]: %+v", arr[2])
	}
}

func TestPUTAndGET(t *testing.T) {
	client, stop := startTestServer(t)
	defer stop()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodPut)
	req.SetRequestURI("http://test/kv/foo")
	req.SetBodyString("bar")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("PUT failed: %v", err)
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusNoContent {
		t.Fatalf("expected 204, got %d", sc)
	}

	req.Reset()
	resp.Reset()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI("http://test/kv/foo")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusOK {
		t.Fatalf("expected 200, got %d", sc)
	}

	var e getEntry
	if err := json.Unmarshal(resp.Body(), &e); err != nil {
		t.Fatalf("invalid JSON from /kv/foo: %v (body=%q)", err, resp.Body())
	}
	if e.Key != "foo" || e.Value == nil || *e.Value != "bar" {
		t.Fatalf("unexpected getEntry: %+v", e)
	}
}

func TestPUTWithTTLExpires(t *testing.T) {
	client, stop := startTestServer(t)
	defer stop()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodPut)
	req.SetRequestURI("http://test/kv/ephemeral?ttl=1")
	req.SetBodyString("tmp")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("put failed: %v", err)
	}
	if resp.StatusCode() != 204 {
		t.Fatalf("expected 204, got %d", resp.StatusCode())
	}

	time.Sleep(1100 * time.Millisecond)

	req.Reset()
	resp.Reset()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI("http://test/kv/ephemeral")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if resp.StatusCode() != 404 {
		t.Fatalf("expected 404 after TTL, got %d", resp.StatusCode())
	}
	var e getEntry
	if err := json.Unmarshal(resp.Body(), &e); err != nil {
		t.Fatalf("invalid JSON body on 404: %v (body=%q)", err, resp.Body())
	}
	if e.Key != "ephemeral" || e.Value != nil {
		t.Fatalf("unexpected 404 body: %+v", e)
	}
}

func TestGETNotFound(t *testing.T) {
	client, stop := startTestServer(t)
	defer stop()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI("http://test/kv/notfound")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if resp.StatusCode() != fasthttp.StatusNotFound {
		t.Fatalf("expected 404 for missing key, got %d", resp.StatusCode())
	}
	var e getEntry
	if err := json.Unmarshal(resp.Body(), &e); err != nil {
		t.Fatalf("invalid JSON body on 404: %v (body=%q)", err, resp.Body())
	}
	if e.Key != "notfound" || e.Value != nil {
		t.Fatalf("unexpected 404 body: %+v", e)
	}
}

func TestPutAndSave(t *testing.T) {
	client, stop := startTestServer(t)
	defer stop()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodPut)
	req.SetRequestURI("http://test/kv/foo")
	req.SetBodyString("bar")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("PUT failed: %v", err)
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusNoContent {
		t.Fatalf("expected 204, got %d", sc)
	}

	req.Reset()
	resp.Reset()
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI("http://test/save")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusNoContent {
		t.Fatalf("expected 204, got %d", sc)
	}

	req.Reset()
	resp.Reset()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI("http://test/kv/foo")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusOK {
		t.Fatalf("expected 200, got %d", sc)
	}
	var e getEntry
	if err := json.Unmarshal(resp.Body(), &e); err != nil {
		t.Fatalf("invalid JSON from /kv/foo: %v (body=%q)", err, resp.Body())
	}
	if e.Key != "foo" || e.Value == nil || *e.Value != "bar" {
		t.Fatalf("unexpected getEntry: %+v", e)
	}
}

func TestPutAndReset(t *testing.T) {
	client, stop := startTestServer(t)
	defer stop()

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodPut)
	req.SetRequestURI("http://test/kv/foo2")
	req.SetBodyString("bar")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("PUT failed: %v", err)
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusNoContent {
		t.Fatalf("expected 204, got %d", sc)
	}

	req.Reset()
	resp.Reset()
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI("http://test/reset")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	if sc := resp.StatusCode(); sc != fasthttp.StatusOK {
		t.Fatalf("expected 200, got %d", sc)
	}

	req.Reset()
	resp.Reset()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI("http://test/kv/foo2")
	if err := client.Do(req, resp); err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if resp.StatusCode() != fasthttp.StatusNotFound {
		t.Fatalf("expected 404 for missing key, got %d", resp.StatusCode())
	}
	var e getEntry
	if err := json.Unmarshal(resp.Body(), &e); err != nil {
		t.Fatalf("invalid JSON body on 404: %v (body=%q)", err, resp.Body())
	}
	if e.Key != "foo2" || e.Value != nil {
		t.Fatalf("unexpected 404 body: %+v", e)
	}
}

func TestStats(t *testing.T) {
	client, stop := startTestServer(t)
	defer stop()

	type StatsDTO struct {
		KeysCount           string `json:"keys_count"`
		ExpirationKeysCount string `json:"expiration_keys_count"`
		UptimeSeconds       string `json:"uptime_seconds"`
		TotalRequests       string `json:"total_requests"`
		Hits                string `json:"hits"`
		Misses              string `json:"misses"`
	}

	{
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		req.Header.SetMethod(fasthttp.MethodPut)
		req.SetRequestURI("http://test/kv/foo")
		req.SetBodyString("bar")
		if err := client.Do(req, resp); err != nil {
			t.Fatalf("PUT failed: %v", err)
		}
		if sc := resp.StatusCode(); sc != fasthttp.StatusNoContent {
			t.Fatalf("expected 204, got %d", sc)
		}
	}

	{
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		req.Header.SetMethod(fasthttp.MethodGet)
		req.SetRequestURI("http://test/kv/foo")
		if err := client.Do(req, resp); err != nil {
			t.Fatalf("GET hit failed: %v", err)
		}
		if sc := resp.StatusCode(); sc != fasthttp.StatusOK {
			t.Fatalf("expected 200, got %d", sc)
		}
	}

	{
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		req.Header.SetMethod(fasthttp.MethodGet)
		req.SetRequestURI("http://test/kv/missing")
		if err := client.Do(req, resp); err != nil {
			t.Fatalf("GET miss failed: %v", err)
		}
		if sc := resp.StatusCode(); sc != fasthttp.StatusNotFound {
			t.Fatalf("expected 404, got %d", sc)
		}
	}

	atoi := func(s string) uint64 {
		var n uint64
		for i := 0; i < len(s); i++ {
			c := s[i]
			if c < '0' || c > '9' {
				t.Fatalf("expected numeric string, got %q", s)
			}
			n = n*10 + uint64(c-'0')
		}
		return n
	}

	getStats := func() StatsDTO {
		req := fasthttp.AcquireRequest()
		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(resp)

		req.Header.SetMethod(fasthttp.MethodGet)
		req.SetRequestURI("http://test/stats")

		if err := client.Do(req, resp); err != nil {
			t.Fatalf("GET /stats failed: %v", err)
		}
		if sc := resp.StatusCode(); sc != fasthttp.StatusOK {
			t.Fatalf("expected 200 from /stats, got %d", sc)
		}

		var dto StatsDTO
		if err := json.Unmarshal(resp.Body(), &dto); err != nil {
			t.Fatalf("invalid JSON from /stats: %v (body=%q)", err, resp.Body())
		}

		_ = atoi(dto.KeysCount)
		_ = atoi(dto.ExpirationKeysCount)
		_ = atoi(dto.UptimeSeconds)
		_ = atoi(dto.TotalRequests)
		_ = atoi(dto.Hits)
		_ = atoi(dto.Misses)

		return dto
	}

	s1 := getStats()

	if kc := atoi(s1.KeysCount); kc != 1 {
		t.Fatalf("keys_count = %d, want 1", kc)
	}
	if ec := atoi(s1.ExpirationKeysCount); ec != 0 {
		t.Fatalf("expiration_keys_count = %d, want 0", ec)
	}
	_ = atoi(s1.UptimeSeconds)

	h1 := atoi(s1.Hits)
	m1 := atoi(s1.Misses)
	tr1 := atoi(s1.TotalRequests)

	if h1 < 1 {
		t.Fatalf("hits = %d, want >= 1 (after GET on existing key)", h1)
	}
	if m1 < 1 {
		t.Fatalf("misses = %d, want >= 1 (after GET on missing key)", m1)
	}
	if tr1 < h1+m1 {
		t.Fatalf("total_requests = %d, want >= hits+misses (%d)", tr1, h1+m1)
	}
}
