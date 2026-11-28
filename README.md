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

## Install

After building the executable, you can copy it to a suitable bin directory, such as `/usr/local/bin`.

```sh
sudo cp basui /usr/local/bin/basui 
```

## Usage


**Usage:**
&emsp;`basui [command]`

**Available Commands:**
&emsp;`build`       Compile layouts, pages and assets into the site's dist directory
&emsp;`init`        Initialize a basui project in the current directory
&emsp;`completion`  Generate the autocompletion script for the specified shell
&emsp;`help`        Help about any command

**Flags:**
&emsp;`-h, --help`   help for basui

Use `basui [command] --help` for more information about a command.

### Example Usage
Initialize a project in the current directory:

```sh
basui init
basui build
```

`init` creates `layouts/`, `pages/`, `partials/`, and `assets/`, along with a minimal `layouts/base.html` and `pages/index.html`. Existing regular files are left unchanged, so the command can also fill in missing project files.

Build the site from `site-directory` and write the generated pages and assets to its `dist` directory. If `site-directory` is omitted, the current directory is used.

```sh
basui build example
basui build
```

The first command builds the example site from `example/`; the second builds the site in the current directory. To view all CLI commands and the build command's options:

```sh
basui --help
basui build --help
```

### Project structure

| Source | Purpose | Output |
| --- | --- | --- |
| layouts/base.html | Shared document shell | Wraps every page |
| partials/**/*.html | Optional reusable fragments | Included by their source path |
| pages/**/*.html | One HTML file per page, including nested pages | Same relative path under dist/ |
| assets/ | Optional CSS, JavaScript, JSON, images, fonts and vendor libraries | Copied unchanged to dist/assets/ |

For example, `pages/index.html` becomes `dist/index.html` and `pages/guides/setup.html` becomes `dist/guides/setup.html`. Only `.html` files in `pages/` and `partials/` are used as templates. `layouts/base.html` and at least one page are required; `partials/` and `assets/` are optional. The `pages/assets/` path is reserved for copied assets and cannot contain pages.

Each build replaces its generated `dist/` directory and removes stale output. Keep your source files outside `dist/`: edits made directly to generated files will be lost on the next build. basui refuses to replace an existing `dist/` directory it did not create.

### How pages and layouts work

The shared layout is the HTML document around every page. It usually contains the document's `<head>`, navigation, and footer. Each page supplies its own content to place inside that layout. At build time, basui combines the layout, any partials, and each page to produce ordinary HTML files.

The layout created by `basui init` uses these named sections:

| Name | How it is used |
| --- | --- |
| `content` | The main content of the page. The starter layout inserts it inside `<main>`. Each page should define it; the build fails if no `content` template is available. |
| `title` | The text inside the document's `<title>`. The starter layout defaults to `basui`; a page can define its own title. |
| `scripts` | Optional page-specific script tags, inserted near the end of `<body>`. The starter layout leaves this empty unless a page defines it. |

`title` and `scripts` work because the starter layout includes them. If you change the layout, a page's definition appears only where the layout includes that name. You can also add your own named sections. Each page's definitions apply only to that page, so one page's title or scripts do not affect another.

### Template syntax

The instructions in double braces are processed when you build the site:

| Syntax | Purpose |
| --- | --- |
| `{{define "name"}}...{{end}}` | Name a section for use elsewhere. Put page content and optional page-specific sections inside definitions. |
| `{{template "name" .}}` | Insert a named section or partial at this location. Keep the `.` so the inserted template can use paths such as `{{.Root}}`. |
| `{{block "name" .}}default{{end}}` | Insert a section with a default. A page can replace the default by defining the same name. An empty block has no default content. |

For example, the layout uses `{{template "content" .}}` to insert a page's `content`. It uses `{{block "title" .}}basui{{end}}` to show `basui` as the title unless the page supplies another. These instructions are different from the browser's HTML `<template>` element.

### Create a page

Create `pages/about.html`:

```html
{{define "title"}}About our site{{end}}
{{define "content"}}
  <h1>About us</h1>
  <p>Welcome to our site.</p>
{{end}}
```

Run `basui build` and open `about.html` on your preview server. The result uses the shared layout, with `About our site` as its title and the page's heading and paragraph inside `<main>`. If you omit the `title` definition, the layout's default is used. Place page HTML inside `content`; text outside the named sections is not inserted into the starter layout.

### Share a layout and partials

Edit `layouts/base.html` to change every page. The following example adds a stylesheet and a reusable navigation partial to the starter layout:

```html
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{block "title" .}}My site{{end}}</title>
  <link rel="stylesheet" href="{{.Root}}assets/site.css">
</head>
<body>
  {{template "partials/nav.html" .}}
  <main>{{template "content" .}}</main>
  {{block "scripts" .}}{{end}}
</body>
</html>
```

Create the referenced file `partials/nav.html`:

```html
<nav aria-label="Main navigation">
  <a href="{{.Root}}index.html">Home</a>
  <a href="{{.Root}}about.html">About</a>
</nav>
```

A partial is reusable HTML named by its path, including `partials/`. Insert it with `{{template "partials/nav.html" .}}` in the layout or in page content. The starter layout does not include navigation by default; if you use the example layout above, create the partial or the build will report a missing template.

### Link to pages and assets

Use `{{.Root}}` at the start of links to pages and assets. It supplies the relative path back to the site's root: it is empty on a top-level page, `../` on a page one folder deep, and so on. This keeps links working on nested pages and when the site is hosted under a URL prefix.

For example, `{{.Root}}assets/site.css` becomes `assets/site.css` on `pages/index.html` and `../assets/site.css` on `pages/guides/setup.html`. Include the `.` in a `template` or `block` instruction so that the inserted section can also use `{{.Root}}`. basui does not provide page metadata or application data to templates.

### Add CSS

Put stylesheets in `assets/`. For example, create `assets/site.css`:

```css
body {
  font-family: sans-serif;
}

nav {
  display: flex;
  gap: 1rem;
}
```

The `<link>` in the layout above loads this stylesheet on every page. basui copies it unchanged to `dist/assets/site.css`; CSS files are not processed as templates. Paths inside CSS, such as `url("images/logo.png")`, are relative to the CSS file. Images, fonts, and other static files can also be placed in `assets/`.

### Add JavaScript

Put JavaScript files in `assets/`. To load a script on every page, add a `<script>` tag to the layout. To load one only on a particular page, add a `scripts` definition to that page:

```html
{{define "scripts"}}
  <script type="module" src="{{.Root}}assets/about.js"></script>
{{end}}
```

For example, create `assets/about.js`:

```js
document.querySelector("h1").textContent = "More about us";
```

This script runs in the visitor's browser on the page that includes it. basui copies JavaScript and other assets unchanged; it does not bundle them or process them as templates. Module imports resolve relative to the JavaScript file. For requests to a file beside a module, you can use `fetch(new URL("./data.json", import.meta.url))`, as in the included example. Other relative `fetch` URLs resolve against the page URL.

### Build, preview, and publish

After editing a page, layout, partial, or asset, run `basui build` again and refresh your browser. basui does not automatically rebuild when files change. To preview the current site's output with Python's static server:

```sh
python3 -m http.server 8080 --bind 127.0.0.1 --directory dist
```

Open http://127.0.0.1:8080/ in your browser. Any static HTTP server can be used; HTTP is preferable to opening files directly, especially for JavaScript modules and `fetch`. Publish the contents of `dist/` to your static web host, preserving the folder structure.

### Troubleshooting

- **Missing content:** Ensure the page defines `{{define "content"}}...{{end}}`.
- **Missing template:** Check the partial's path and name, including the `partials/` prefix, and make sure its file exists.
- **Changes do not appear:** Rebuild, then refresh the page. Check that you edited the source files, not `dist/`.
- **Styles or scripts do not load:** Check the browser's network errors, the asset's location under `assets/`, and the `{{.Root}}` prefix in the HTML.
- **Output directory is not a basui output directory:** Move or rename the existing `dist/` before building; basui will not overwrite unrelated files.
- **Cannot acquire build lock:** Check whether another build is running. If an interrupted build left `.basui-lock` behind, remove it only after confirming no build is active.

Template and asset errors leave the previous generated `dist/` intact. If publishing the new output fails, basui attempts to restore the previous output and reports its retained location if restoration also fails.
