package storage_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/taymour/elysiandb/internal/configuration"
	"github.com/taymour/elysiandb/internal/globals"
	"github.com/taymour/elysiandb/internal/storage"
)

func setTmpConfig(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	globals.SetConfig(&configuration.Config{
		Store: configuration.StoreConfig{
			Folder: tmp,
		},
	})
	return tmp
}

func writeFile(t *testing.T, dir, name string, data []byte) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return p
}

func TestReadFromDB_OK(t *testing.T) {
	dir := setTmpConfig(t)
	content := []byte(`{"foo":"YmFy","bin":"AQID"}`)
	writeFile(t, dir, storage.DataFile, content)

	got, err := storage.ReadFromDB(storage.DataFile)
	if err != nil {
		t.Fatalf("ReadFromDB error: %v", err)
	}
	want := map[string][]byte{
		"foo": []byte("bar"),
		"bin": {0x01, 0x02, 0x03},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestReadFromDB_EmptyFile_ReturnsEmptyMap(t *testing.T) {
	dir := setTmpConfig(t)
	writeFile(t, dir, storage.DataFile, nil)

	got, err := storage.ReadFromDB(storage.DataFile)
	if err != nil {
		t.Fatalf("ReadFromDB error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %#v", got)
	}
}

func TestReadFromDB_FileMissing_Error(t *testing.T) {
	_ = setTmpConfig(t)
	_, err := storage.ReadFromDB("does-not-exist.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestReadFromDB_InvalidJSON_Error(t *testing.T) {
	dir := setTmpConfig(t)
	writeFile(t, dir, storage.DataFile, []byte(`{"foo":`))

	_, err := storage.ReadFromDB(storage.DataFile)
	if err == nil {
		t.Fatal("expected JSON error, got nil")
	}
}

func TestReadExpirationsFromDB_OK(t *testing.T) {
	dir := setTmpConfig(t)
	writeFile(t, dir, storage.ExpirationDataFile, []byte(`{
		"1690000000":["a","b"],
		"1690000002":[]
	}`))

	got, err := storage.ReadExpirationsFromDB(storage.ExpirationDataFile)
	if err != nil {
		t.Fatalf("ReadExpirationsFromDB error: %v", err)
	}

	want := map[int64][]string{
		1690000000: {"a", "b"},
		1690000002: nil,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%#v want=%#v", got, want)
	}

	got[1690000000] = append(got[1690000000], "c")
}

func TestReadExpirationsFromDB_EmptyFile_ReturnsEmptyMap(t *testing.T) {
	dir := setTmpConfig(t)
	writeFile(t, dir, storage.ExpirationDataFile, nil)

	got, err := storage.ReadExpirationsFromDB(storage.ExpirationDataFile)
	if err != nil {
		t.Fatalf("ReadExpirationsFromDB error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %#v", got)
	}
}

func TestReadExpirationsFromDB_FileMissing_Error(t *testing.T) {
	_ = setTmpConfig(t)
	_, err := storage.ReadExpirationsFromDB("does-not-exist.expiration.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestReadExpirationsFromDB_InvalidJSON_Error(t *testing.T) {
	dir := setTmpConfig(t)
	writeFile(t, dir, storage.ExpirationDataFile, []byte(`{"169": ["ok"],`))

	_, err := storage.ReadExpirationsFromDB(storage.ExpirationDataFile)
	if err == nil {
		t.Fatal("expected JSON error, got nil")
	}
}

func TestReadExpirationsFromDB_UnparseableKey_Skipped(t *testing.T) {
	dir := setTmpConfig(t)
	writeFile(t, dir, storage.ExpirationDataFile, []byte(`{
		"not-a-number":["x"],
		"1700000000":["y"]
	}`))

	got, err := storage.ReadExpirationsFromDB(storage.ExpirationDataFile)
	if err != nil {
		t.Fatalf("ReadExpirationsFromDB error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 parsed entry, got %d (%v)", len(got), got)
	}
	if _, ok := got[1700000000]; !ok {
		t.Fatalf("expected key 1700000000 to be present, got %v", got)
	}
}
