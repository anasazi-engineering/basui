package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const marker = ".basui-output"

func buildSite(site string) (int, error) {
	lock := filepath.Join(site, ".basui-lock")
	if err := os.Mkdir(lock, 0700); err != nil {
		return 0, fmt.Errorf("acquire build lock: %w", err)
	}
	defer os.Remove(lock)

	// Stage on the same filesystem so publication uses directory renames.
	work, err := os.MkdirTemp(site, ".basui-build-")
	if err != nil {
		return 0, err
	}
	cleanup := true
	defer func() {
		if cleanup {
			os.RemoveAll(work)
		}
	}()
	next := filepath.Join(work, "next")
	if err := os.Mkdir(next, 0755); err != nil {
		return 0, err
	}
	count, err := renderSite(os.DirFS(site), next)
	if err != nil {
		return 0, err
	}
	if err := os.WriteFile(filepath.Join(next, marker), []byte("basui v1\n"), 0644); err != nil {
		return 0, err
	}

	dist := filepath.Join(site, "dist")
	info, err := os.Lstat(dist)
	if errors.Is(err, fs.ErrNotExist) {
		return count, os.Rename(next, dist)
	}
	if err != nil {
		return 0, err
	}
	signature, readErr := os.ReadFile(filepath.Join(dist, marker))
	if !info.IsDir() || readErr != nil || string(signature) != "basui v1\n" {
		return 0, fmt.Errorf("%s is not a basui output directory", dist)
	}
	previous := filepath.Join(work, "previous")
	if err := os.Rename(dist, previous); err != nil {
		return 0, err
	}
	if err := os.Rename(next, dist); err != nil {
		if restoreErr := os.Rename(previous, dist); restoreErr != nil {
			cleanup = false
			return 0, fmt.Errorf("publish: %w; previous output retained at %s: %v", err, previous, restoreErr)
		}
		return 0, err
	}
	return count, nil
}

// walkFiles visits regular files only. Optional directories may be absent.
func walkFiles(src fs.FS, dir string, optional bool, visit func(string) error) error {
	return fs.WalkDir(src, dir, func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			if optional && name == dir && errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if name == dir || !entry.Type().IsRegular() {
			return fmt.Errorf("unexpected file type at %s", name)
		}
		return visit(name)
	})
}

func renderSite(src fs.FS, out string) (int, error) {
	base, err := fs.ReadFile(src, "layouts/base.html")
	if err != nil {
		return 0, err
	}
	shared, err := template.New("base").Option("missingkey=error").Parse(string(base))
	if err != nil {
		return 0, err
	}
	err = walkFiles(src, "partials", true, func(name string) error {
		if path.Ext(name) != ".html" {
			return nil
		}
		body, err := fs.ReadFile(src, name)
		if err != nil {
			return err
		}
		_, err = shared.New(name).Parse(string(body))
		return err
	})
	if err != nil {
		return 0, err
	}
	count := 0
	err = walkFiles(src, "pages", false, func(name string) error {
		if path.Ext(name) != ".html" {
			return nil
		}
		rel := strings.TrimPrefix(name, "pages/")
		if strings.HasPrefix(rel, "assets/") {
			return fmt.Errorf("%s conflicts with the assets directory", name)
		}
		err := renderPage(src, shared, name, rel, out)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		count++
		return nil
	})
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, errors.New("pages directory contains no .html files")
	}
	err = walkFiles(src, "assets", true, func(name string) error {
		dest := filepath.Join(out, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		return copyAsset(src, name, dest)
	})
	return count, err
}

func renderPage(src fs.FS, shared *template.Template, name, rel, out string) error {
	body, err := fs.ReadFile(src, name)
	if err != nil {
		return err
	}
	page, err := shared.Clone()
	if err != nil {
		return err
	}
	if _, err := page.New(name).Parse(string(body)); err != nil {
		return err
	}
	if page.Lookup("content") == nil {
		return errors.New("missing content template")
	}
	dest := filepath.Join(out, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	data := struct{ Root string }{Root: strings.Repeat("../", strings.Count(rel, "/"))}
	err = page.ExecuteTemplate(file, "base", data)
	return errors.Join(err, file.Close())
}
