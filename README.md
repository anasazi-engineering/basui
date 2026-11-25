![Project Logo](basui_logo.png)

# basui

A small HTML-first static compiler CLI written in Go. The included two-screen example uses plain JavaScript, local CSS and a sample JSON fixture. It runs offline once built.

## Quick start

Install Go 1.23 or later. From this directory run:

```sh
go mod tidy
go test ./...
go build -o basui .
./basui build example
python3 -m http.server 8080 --bind 127.0.0.1 --directory example/dist
```

Open http://127.0.0.1:8080/ and select AHU-1. The first Go dependency download needs network access. Python is just a convenient local preview server. Any static HTTP server can host the generated site.

## Site convention

| Source | Purpose | Output |
| --- | --- | --- |
| layouts/base.html | Shared document shell | Wraps every page |
| partials/**/*.html | Optional reusable fragments | Included by their source path |
| pages/**/*.html | One HTML file per screen, including nested screens | Same relative path under dist/ |
| assets/ | Optional CSS, JavaScript, JSON, images, fonts and vendor libraries | Copied unchanged to dist/assets/ |

A site directory defaults to the current directory. Each build replaces its generated dist directory and removes stale files. Build into a new site or move an existing hand-authored dist directory before the first build. The compiler checks its ownership marker before replacing previous output.

## Add a screen

Create pages/alarms.html:

```html
{{define "title"}}Active alarms{{end}}
{{define "content"}}
  <h1>Active alarms</h1>
  <p role="status">No alarms loaded.</p>
{{end}}
```

Call a partial with {{template "partials/nav.html" .}}. The example layout supplies default title and scripts blocks. Each page defines content and can override either default. Each page gets an independent clone of the shared templates. For a new page-specific script, add:

```html
{{define "scripts"}}<script type="module" src="{{.Root}}assets/alarms.js"></script>{{end}}
```

.Root is the relative path prefix back to the site root (empty for top-level pages). Use it for site links and assets to support nested pages and deployment under a URL prefix. Paths inside CSS and ES module imports resolve relative to those files.
