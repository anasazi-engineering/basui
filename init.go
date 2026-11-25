package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func initSite(site string) error {
	for _, dir := range []string{"layouts", "pages", "partials", "assets"} {
		if err := os.MkdirAll(filepath.Join(site, dir), 0755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
	}

	files := []struct {
		name string
		body string
	}{
		{"layouts/base.html", `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{block "title" .}}basui{{end}}</title>
</head>
<body>
  <main>{{template "content" .}}</main>
  {{block "scripts" .}}{{end}}
</body>
</html>
`},
		{"pages/index.html", `{{define "content"}}
  <h1>Welcome to basui</h1>
{{end}}
`},
	}
	for _, file := range files {
		name := filepath.Join(site, filepath.FromSlash(file.name))
		out, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if errors.Is(err, fs.ErrExist) {
			info, statErr := os.Lstat(name)
			if statErr != nil {
				return fmt.Errorf("inspect %s: %w", file.name, statErr)
			}
			if !info.Mode().IsRegular() {
				return fmt.Errorf("%s already exists and is not a regular file", file.name)
			}
			continue
		}
		if err != nil {
			return fmt.Errorf("create %s: %w", file.name, err)
		}
		_, writeErr := out.WriteString(file.body)
		if err := errors.Join(writeErr, out.Close()); err != nil {
			return fmt.Errorf("write %s: %w", file.name, err)
		}
	}
	return nil
}
