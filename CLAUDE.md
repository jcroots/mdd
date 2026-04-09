# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

mdd is a local markdown HTTP server written in Go. It serves `.md` files from the current working directory as rendered HTML. Requests to `/` or directories default to `README.md`. Only `.md` files within the working directory are served (path traversal is blocked).

Server listens on `127.0.0.1:1980`.

## Build & Install

```bash
make install    # runs go get && go install
go run main.go  # run without installing
```

There are no tests or linting configured.

## Architecture

Single-file Go application (`main.go`) with one HTTP handler (`serveMarkdown`). Uses `github.com/gomarkdown/markdown` for markdown-to-HTML conversion with common extensions and `HrefTargetBlank` enabled. The rendered HTML is wrapped in a minimal inline-styled page template.
