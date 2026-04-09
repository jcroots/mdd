// Copyright (c) Jeremías Casteglione <jeremias.rootstrap@gmail.com>
// See LICENSE file.

package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupTestDir creates a temp directory with test markdown files and
// changes the working directory to it. Returns a cleanup function.
func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })

	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Hello\n\nWorld"), 0644)
	os.WriteFile(filepath.Join(dir, "page.md"), []byte("# Page\n\n**bold**"), 0644)
	os.MkdirAll(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "subdir", "README.md"), []byte("# Sub"), 0644)

	return dir
}

func TestRootServesREADME(t *testing.T) {
	setupTestDir(t)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	serveMarkdown(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "<h1") || !strings.Contains(body, "Hello") {
		t.Fatalf("expected rendered markdown with Hello heading, got: %s", body)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Fatalf("expected text/html content type, got: %s", ct)
	}
}

func TestServesNamedFile(t *testing.T) {
	setupTestDir(t)

	req := httptest.NewRequest("GET", "/page.md", nil)
	w := httptest.NewRecorder()
	serveMarkdown(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "<strong>bold</strong>") {
		t.Fatalf("expected rendered bold text, got: %s", body)
	}
}

func TestDirectoryDefaultsToREADME(t *testing.T) {
	setupTestDir(t)

	req := httptest.NewRequest("GET", "/subdir", nil)
	w := httptest.NewRecorder()
	serveMarkdown(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Sub") {
		t.Fatalf("expected subdir README content")
	}
}

func TestNonMdFileRejected(t *testing.T) {
	setupTestDir(t)

	req := httptest.NewRequest("GET", "/secret.txt", nil)
	w := httptest.NewRecorder()
	serveMarkdown(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPathTraversalBlocked(t *testing.T) {
	setupTestDir(t)

	req := httptest.NewRequest("GET", "/../../../etc/passwd.md", nil)
	w := httptest.NewRecorder()
	serveMarkdown(w, req)

	if w.Code == http.StatusOK {
		t.Fatal("expected non-200 for path traversal attempt")
	}
}

func TestMissingFileReturns404(t *testing.T) {
	setupTestDir(t)

	req := httptest.NewRequest("GET", "/nonexistent.md", nil)
	w := httptest.NewRecorder()
	serveMarkdown(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
