package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func put(t *testing.T, site, name, body string) {
	t.Helper()
	target := filepath.Join(site, name)
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, site, name string) string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(site, name))
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func TestBuild(t *testing.T) {
	site := t.TempDir()
	put(t, site, "layouts/base.html", `{{template "partials/nav.html" .}}{{template "content" .}}{{block "scripts" .}}{{end}}`)
	put(t, site, "partials/nav.html", `<a href="{{.Root}}index.html">Home</a>`)
	put(t, site, "pages/index.html", `{{define "content"}}Home page{{end}}{{define "scripts"}}INDEX SCRIPT{{end}}`)
	put(t, site, "pages/z/device.html", `{{define "content"}}Device page{{end}}`)
	put(t, site, "assets/data.bin", "\x00\xff")
	count, err := buildSite(site)
	if err != nil || count != 2 {
		t.Fatalf("built %d pages: %v", count, err)
	}
	got := read(t, site, "dist/z/device.html")
	if got != `<a href="../index.html">Home</a>Device page` {
		t.Fatalf("unexpected nested page: %q", got)
	}
	if !strings.HasSuffix(read(t, site, "dist/index.html"), "Home pageINDEX SCRIPT") {
		t.Fatal("index scripts missing")
	}
	if read(t, site, "dist/assets/data.bin") != "\x00\xff" {
		t.Fatal("asset bytes changed")
	}
	put(t, site, "pages/z/device.html", `{{define "content"}}{{template "unknown" .}}{{end}}`)
	if _, err := buildSite(site); err == nil {
		t.Fatal("accepted missing partial")
	}
	if read(t, site, "dist/z/device.html") != got {
		t.Fatal("failed build changed output")
	}
	if err := os.Remove(filepath.Join(site, "pages/z/device.html")); err != nil {
		t.Fatal(err)
	}
	if _, err := buildSite(site); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(site, "dist/z/device.html")); !os.IsNotExist(err) {
		t.Fatal("deleted page still in output")
	}
}

func TestUnmanagedOutput(t *testing.T) {
	site := t.TempDir()
	put(t, site, "layouts/base.html", `{{template "content" .}}`)
	put(t, site, "pages/index.html", `{{define "content"}}Home{{end}}`)
	put(t, site, "dist/index.html", "existing")
	if _, err := buildSite(site); err == nil || read(t, site, "dist/index.html") != "existing" {
		t.Fatal("unmanaged output was not protected")
	}
}
