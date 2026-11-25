package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitCommand(t *testing.T) {
	site := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(site); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Error(err)
		}
	})

	var output bytes.Buffer
	cmd := newRootCommand()
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"init"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{"layouts", "pages", "partials", "assets"} {
		info, err := os.Stat(filepath.Join(site, dir))
		if err != nil || !info.IsDir() {
			t.Fatalf("%s is not a directory: %v", dir, err)
		}
	}
	if count, err := buildSite(site); err != nil || count != 1 {
		t.Fatalf("built %d pages: %v", count, err)
	}
	if got := read(t, site, "dist/index.html"); !strings.Contains(got, "<h1>Welcome to basui</h1>") {
		t.Fatalf("starter page not rendered: %q", got)
	}

	put(t, site, "layouts/base.html", "custom layout")
	put(t, site, "pages/index.html", "custom page")
	if err := os.Remove(filepath.Join(site, "assets")); err != nil {
		t.Fatal(err)
	}
	if err := initSite(site); err != nil {
		t.Fatal(err)
	}
	if read(t, site, "layouts/base.html") != "custom layout" || read(t, site, "pages/index.html") != "custom page" {
		t.Fatal("init changed existing files")
	}
	if info, err := os.Stat(filepath.Join(site, "assets")); err != nil || !info.IsDir() {
		t.Fatalf("missing assets directory was not restored: %v", err)
	}
}

func TestInitConflicts(t *testing.T) {
	t.Run("directory", func(t *testing.T) {
		site := t.TempDir()
		put(t, site, "pages", "existing")
		if err := initSite(site); err == nil {
			t.Fatal("accepted a file in place of a directory")
		}
		if read(t, site, "pages") != "existing" {
			t.Fatal("changed conflicting file")
		}
	})
	t.Run("starter file", func(t *testing.T) {
		site := t.TempDir()
		if err := os.MkdirAll(filepath.Join(site, "layouts", "base.html"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := initSite(site); err == nil {
			t.Fatal("accepted a directory in place of a starter file")
		}
	})
}

func TestInitRejectsArguments(t *testing.T) {
	cmd := newRootCommand()
	cmd.SetArgs([]string{"init", "elsewhere"})
	cmd.SetErr(&bytes.Buffer{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("init accepted an unexpected argument")
	}
}
