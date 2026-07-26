package fsx

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

func testFS() fstest.MapFS {
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	return fstest.MapFS{
		"main.go":           {Data: []byte("package main\nfunc main() {}\n"), ModTime: t0},
		"README.md":         {Data: []byte("# title\nhello world\n"), ModTime: t0.Add(time.Hour)},
		"internal/util.go":  {Data: []byte("package internal\n// hello\n"), ModTime: t0.Add(2 * time.Hour)},
		"internal/data.bin": {Data: []byte{0x00, 0x01, 'h', 'e', 'l', 'l', 'o'}, ModTime: t0},
		"internal/sub/x.go": {Data: []byte("package sub\n"), ModTime: t0},
		"vendor/lib/dep.go": {Data: []byte("package lib\nhello\n"), ModTime: t0},
		"notes":             {Data: []byte("no extension here\n"), ModTime: t0},
	}
}

func TestFindFiles(t *testing.T) {
	fsys := testFS()
	got, err := FindFiles(fsys, ".", ".go")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"internal/sub/x.go", "internal/util.go", "main.go", "vendor/lib/dep.go"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("= %v\nwant %v", got, want)
	}

	got, err = FindFiles(fsys, "internal", ".go")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []string{"internal/sub/x.go", "internal/util.go"}) {
		t.Errorf("subtree = %v", got)
	}

	got, _ = FindFiles(fsys, ".", "")
	if len(got) != 7 {
		t.Errorf("empty ext matched %d files, want 7", len(got))
	}
	if _, err := FindFiles(fsys, "nope", ".go"); err == nil {
		t.Error("a missing root should be an error")
	}
}

func TestGrep(t *testing.T) {
	got, err := Grep(testFS(), "hello")
	if err != nil {
		t.Fatal(err)
	}
	want := []Match{
		{Path: "README.md", Line: 2, Text: "hello world"},
		{Path: "internal/util.go", Line: 2, Text: "// hello"},
		{Path: "vendor/lib/dep.go", Line: 2, Text: "hello"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("= %#v\nwant %#v", got, want)
	}
	if got, _ := Grep(testFS(), "zzz"); len(got) != 0 {
		t.Errorf("no matches expected, got %v", got)
	}
}

func TestWalk(t *testing.T) {
	got, err := Walk(testFS())
	if err != nil {
		t.Fatal(err)
	}
	if got.Files != 7 {
		t.Errorf("Files = %d, want 7", got.Files)
	}
	if got.Dirs != 4 {
		t.Errorf("Dirs = %d, want 4 (internal, internal/sub, vendor, vendor/lib)", got.Dirs)
	}
	if got.ByExt[".go"] != 4 || got.ByExt[".md"] != 1 || got.ByExt[""] != 1 || got.ByExt[".bin"] != 1 {
		t.Errorf("ByExt = %v", got.ByExt)
	}
	var wantBytes int64
	for _, f := range testFS() {
		wantBytes += int64(len(f.Data))
	}
	if got.Bytes != wantBytes {
		t.Errorf("Bytes = %d, want %d", got.Bytes, wantBytes)
	}
	if got.Largest != "main.go" {
		t.Errorf("Largest = %q, want main.go", got.Largest)
	}
	if s, err := Walk(fstest.MapFS{}); err != nil || s.Files != 0 || s.Largest != "" {
		t.Errorf("empty fs = %+v, %v", s, err)
	}
}

func TestCopyTree(t *testing.T) {
	dst := t.TempDir()
	if err := CopyTree(testFS(), dst); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dst, "internal", "sub", "x.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "package sub\n" {
		t.Errorf("copied content = %q", data)
	}
	n := 0
	filepath.WalkDir(dst, func(_ string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			n++
		}
		return nil
	})
	if n != 7 {
		t.Errorf("copied %d files, want 7", n)
	}
}

func TestWriteFileAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	if err := WriteFileAtomic(path, []byte("first"), 0o644); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(path); string(data) != "first" {
		t.Errorf("= %q", data)
	}
	if err := WriteFileAtomic(path, []byte("second"), 0o644); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(path); string(data) != "second" {
		t.Errorf("overwrite = %q", data)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o644 {
		t.Errorf("mode = %v, want 0644", info.Mode().Perm())
	}

	// No temp files may be left behind.
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		t.Errorf("directory contains %v, want just config.json", names)
	}

	if err := WriteFileAtomic(filepath.Join(dir, "nope", "x"), []byte("x"), 0o644); err == nil {
		t.Error("writing into a missing directory should fail")
	}
}

func TestTailLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	var sb strings.Builder
	for i := 1; i <= 1000; i++ {
		sb.WriteString("line ")
		sb.WriteString(strings.Repeat("x", i%7))
		sb.WriteString("\n")
	}
	content := sb.String()
	os.WriteFile(path, []byte(content), 0o644)

	all := strings.Split(strings.TrimSuffix(content, "\n"), "\n")

	got, err := TailLines(path, 3)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, all[len(all)-3:]) {
		t.Errorf("= %q, want %q", got, all[len(all)-3:])
	}

	got, err = TailLines(path, 5000)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1000 {
		t.Errorf("asking for more lines than exist gave %d", len(got))
	}

	os.WriteFile(path, []byte("no trailing newline"), 0o644)
	if got, _ := TailLines(path, 2); !reflect.DeepEqual(got, []string{"no trailing newline"}) {
		t.Errorf("= %q", got)
	}
	os.WriteFile(path, []byte(""), 0o644)
	if got, err := TailLines(path, 2); err != nil || len(got) != 0 {
		t.Errorf("empty file = %q, %v", got, err)
	}
	if _, err := TailLines(filepath.Join(dir, "missing"), 1); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("missing file = %v, want fs.ErrNotExist", err)
	}
}

// TestTailLinesDoesNotReadWholeFile writes 8 MB and asks for two lines.
// Reading it all works but is the wrong answer; the allocation budget catches it.
func TestTailLinesDoesNotReadWholeFile(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "big.log")
	line := strings.Repeat("y", 127) + "\n"
	var sb strings.Builder
	for range 65536 {
		sb.WriteString(line)
	}
	os.WriteFile(path, []byte(sb.String()), 0o644)

	var got []string
	allocs := testing.AllocsPerRun(3, func() {
		v, err := TailLines(path, 2)
		if err != nil {
			t.Fatal(err)
		}
		got = v
	})
	if len(got) != 2 {
		t.Fatalf("got %d lines", len(got))
	}
	if allocs > 100 {
		t.Errorf("%.0f allocations for two lines of an 8 MB file; seek to the end "+
			"and read backwards", allocs)
	}
}

func TestNewest(t *testing.T) {
	path, mt, err := Newest(testFS())
	if err != nil {
		t.Fatal(err)
	}
	if path != "internal/util.go" {
		t.Errorf("= %q, want internal/util.go", path)
	}
	if mt.IsZero() {
		t.Error("mod time is zero")
	}
	if _, _, err := Newest(fstest.MapFS{}); err == nil {
		t.Error("an empty filesystem should report an error")
	}
}
