# Project: static site generator

**Modules it exercises:** 02 (strings), 03 (interfaces), 06 (concurrency), 07 (fs).
**Rough size:** 500–800 lines. **Difficulty:** ★★☆

Turn a directory of Markdown into a directory of HTML, with templates, an index
page, and a watch mode.

## The spec

```
sitegen -in ./content -out ./public -templates ./templates -watch
```

Input:

```
content/
  index.md
  posts/2026-01-02-hello.md
  static/logo.png
```

Each Markdown file starts with front matter:

```
---
title: Hello, world
date: 2026-01-02
tags: [go, testing]
draft: false
---

# Hello

Some **markdown**.
```

Output: `public/posts/2026-01-02-hello/index.html`, rendered through
`templates/post.html`, plus `public/index.html` listing every non-draft post
newest first, plus every static file copied across unchanged.

## Requirements

1. **Front matter parser.** Your own — a tiny `key: value` parser handling
   strings, dates, booleans and `[a, b]` lists. Errors must name the file and
   line.
2. **Markdown subset**, written by hand: headings, bold/italic, inline code,
   fenced code blocks, links, images, unordered lists, paragraphs. It does not
   have to be CommonMark; it has to be *documented* and *tested*.
3. **HTML escaping.** Text is escaped, code blocks doubly so. Write the test
   that proves `<script>` in a title cannot escape into the output. Use
   `html/template`, not `text/template`, and know why.
4. **Templates** with `html/template`: a base layout, `post.html`, `index.html`,
   and a `{{ define }}`-based partial for the post list.
5. **Incremental + concurrent.** Render pages in parallel (bounded), and in
   watch mode re-render only what changed.
6. **Deterministic output.** Running it twice on the same input produces
   byte-identical files. Sorted map iteration everywhere.

## Milestones

1. Walk `content/`, parse front matter, print the metadata. Use `fs.FS` so the
   tests can use `fstest.MapFS`.
2. Markdown → HTML for one file, tested with golden files.
3. Templates and the output layout.
4. Static copying, the index page, drafts.
5. Parallel rendering with a worker pool; check the output is still identical.
6. Watch mode: poll mtimes every 500 ms (a real fs-notify is a dependency, and
   polling is fine at this scale).

## The interesting problem

**Golden-file testing.** Keep `testdata/in/*.md` and `testdata/want/*.html`, and
a `-update` flag:

```go
var update = flag.Bool("update", false, "rewrite golden files")
```

The test renders, compares, and with `-update` rewrites the expectation. It is
the cheapest way to test anything that produces text, and it makes reviewing a
rendering change trivial: the diff *is* the change.

## Stretch

- Syntax highlighting for fenced code blocks — for Go only, using `go/scanner`
  from the standard library. That is a genuinely fun 60 lines.
- An RSS feed and a `sitemap.xml`.
- A `-serve` flag that runs an HTTP server with live reload over SSE.
