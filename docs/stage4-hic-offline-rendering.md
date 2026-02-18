# HIC Offline Rendering Engine

## Static-First Architecture, Build Pipeline & Air-Gap Compliance

> *"The Holm Intelligence Complex does not stream. It does not fetch. It does not phone home.
> Every floor, every corridor, every room exists in poured concrete before the first
> visitor ever walks through the lobby. The neon lights are beautiful, yes --- but the
> building stands in the dark."*

**Document Classification:** INFRASTRUCTURE / BUILD SYSTEMS / OFFLINE COMPLIANCE
**Applies To:** All HIC rendering, compilation, and deployment pipelines
**Last Structural Audit:** Current build cycle
**Air-Gap Status:** ENFORCED

---

## Table of Contents

1. [Static-First Philosophy](#static-first-philosophy)
2. [The Build Pipeline](#the-build-pipeline)
3. [Markdown as Blueprint Language](#markdown-as-blueprint-language)
4. [HTML as Construction Material](#html-as-construction-material)
5. [CSS as Neon Lighting](#css-as-neon-lighting)
6. [Zero-JavaScript Guarantee](#zero-javascript-guarantee)
7. [Air-Gap Compliance Checklist](#air-gap-compliance-checklist)
8. [The SVG Pipeline](#the-svg-pipeline)
9. [USB Import Protocol](#usb-import-protocol)
10. [Build Reproducibility](#build-reproducibility)
11. [Disaster Recovery Builds](#disaster-recovery-builds)
12. [Performance Budget](#performance-budget)

---

## Static-First Philosophy

### The Building Is Made of Concrete, Not Holograms

The Holm Intelligence Complex is not a holographic projection that vanishes when the
projector loses power. It is a skyscraper made of reinforced concrete and structural
steel. Every floor exists whether or not anyone is looking at it. Every room has walls
whether or not the lights are on.

This is the foundational principle of the HIC rendering engine: **the output must
function with zero runtime dependencies**. No server-side rendering on each request.
No client-side JavaScript assembling the page after load. No network calls to populate
content. The HTML file, as it sits on disk, is the complete and final artifact.

The implications are absolute:

- **Every page is a self-contained HTML file.** Open it in any browser, on any
  machine, from a local filesystem or a static server. It renders. It reads. It works.
- **No JavaScript is required for any core functionality.** Navigation, content
  display, table of contents, cross-references --- all of it functions with JavaScript
  completely disabled.
- **No external resources are referenced.** No CDN-hosted fonts. No remotely loaded
  icon libraries. No analytics endpoints. No third-party stylesheets. The HTML file
  and its co-located assets are everything.
- **The building is pre-fabricated.** Every page is fully constructed at build time.
  The "server" (if one exists at all) is a dumb file server. It does not interpret,
  transform, or generate. It serves bytes from disk.

### Progressive Enhancement: Neon Is Optional, Structure Is Mandatory

The HIC is a cyberpunk skyscraper. Its neon lighting is iconic --- pulsing cyan
borders, glowing magenta accents, the soft hum of animated flourishes. But neon
is a layer, not a foundation.

The rendering philosophy follows strict progressive enhancement:

| Layer | Building Metaphor | Technology | Required? |
|-------|-------------------|------------|-----------|
| Structure | Concrete, steel, floors, walls | Semantic HTML | **Yes** |
| Surface | Paint, signage, wayfinding | CSS (inline-capable) | Expected |
| Atmosphere | Neon lights, ambient glow | CSS animations | No |
| Intelligence | Elevator AI, smart lighting | JavaScript | **No** |

If CSS fails to load, the content is still readable, navigable, and complete.
If JavaScript fails to load (or is deliberately blocked), the user loses nothing
of substance. The building stands in the dark.

### Pre-Fabrication Over Runtime Assembly

Traditional web applications are built in real time: a request arrives, a server
assembles a response from databases and templates and API calls, and a page is
delivered --- often incomplete, requiring further client-side assembly.

The HIC rejects this model entirely. The build pipeline runs once, producing a
directory of static files. These files are the building. Deployment is copying
files. Serving is returning files. There is no assembly line running during
visiting hours.

```
TRADITIONAL WEB APP:                    HIC STATIC BUILD:

Request → Server → DB → Template →      Request → File Server → HTML File
  → API calls → Assembly → Response        (that's it)
  → Client JS → More API calls
  → Final render (maybe)
```

This is not a limitation. This is a feature. The fastest page load is the one
that requires no computation. The most reliable page is the one with no moving
parts. The most secure page is the one that executes nothing.

### The Concrete Test

Every design decision in the HIC must pass the Concrete Test:

> If this feature were a physical component of a building, would it exist
> without electricity?

Walls exist without electricity. Doors exist without electricity. Painted signs
exist without electricity. The directory board in the lobby exists without
electricity. These are HTML and CSS.

Automatic sliding doors require electricity. The digital directory kiosk requires
electricity. The mood lighting requires electricity. These are JavaScript.

Build the walls first. Install the doors. Paint the signs. Then, and only then,
consider whether automated sliding doors would improve the experience. If the
power goes out, people can still push the doors open manually.

---

## The Build Pipeline

### Overview: From Blueprints to Building

The HIC build pipeline transforms a directory of markdown files into a complete
static website. It runs locally, requires no network access, has no external
package dependencies, and completes in seconds.

```
┌─────────────────────────────────────────────────────────────┐
│                    THE BUILD PIPELINE                        │
│                                                             │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐              │
│  │ Markdown │───>│ build.py │───>│  Static  │              │
│  │  Source  │    │ compiler │    │   Site   │              │
│  └──────────┘    └──────────┘    └──────────┘              │
│                       │                                     │
│              ┌────────┼────────┐                            │
│              │        │        │                            │
│              v        v        v                            │
│          ┌──────┐ ┌──────┐ ┌──────┐                        │
│          │ HTML │ │ CSS  │ │Assets│                        │
│          │pages │ │ file │ │copied│                        │
│          └──────┘ └──────┘ └──────┘                        │
│                                                             │
│  Network access: NONE                                       │
│  External deps:  NONE                                       │
│  Build time:     SECONDS                                    │
└─────────────────────────────────────────────────────────────┘
```

### Source Materials: Markdown Files as Architectural Blueprints

The raw materials for the HIC are markdown files organized in a directory
structure that mirrors the building's floor plan:

```
content/
├── index.md                    # Lobby / Ground Floor
├── getting-started/
│   ├── index.md                # Getting Started Wing - Entrance
│   ├── installation.md         # Installation Room
│   ├── quickstart.md           # Quick Start Room
│   └── configuration.md        # Configuration Room
├── architecture/
│   ├── index.md                # Architecture Wing - Entrance
│   ├── overview.md             # Overview Room
│   ├── decisions.md            # Decision Records Room
│   └── principles.md           # Principles Room
├── reference/
│   ├── index.md                # Reference Wing - Entrance
│   ├── api.md                  # API Room
│   └── glossary.md             # Glossary Room
└── operations/
    ├── index.md                # Operations Wing - Entrance
    ├── deployment.md           # Deployment Room
    └── troubleshooting.md      # Troubleshooting Room
```

Each markdown file contains:
- **Front matter** (YAML): metadata about the room (title, order, category)
- **Content** (Markdown): the actual documentation
- **Nothing else**: no embedded scripts, no HTML hacks, no remote includes

Example source file:

```markdown
---
title: Installation
order: 1
category: getting-started
description: How to install and set up the system
---

# Installation

This guide covers the installation process...
```

### The Compiler: build.py as the Construction Crew

The entire build system is a single Python file: `build.py`. It has **zero
external dependencies** beyond the Python standard library. No pip install.
No requirements.txt. No virtual environment needed (though one is recommended
for development).

```python
#!/usr/bin/env python3
"""
HIC Build Pipeline
Transforms markdown blueprints into static HTML pages.
Zero external dependencies. Runs with Python 3.8+.
"""

import os
import re
import json
import hashlib
from pathlib import Path
from datetime import datetime

# Standard library markdown processing --- no external packages
# The markdown parser is included in the build system itself
```

Why a single file with no dependencies?

- **Air-gap compliance**: No package registry access needed
- **Reproducibility**: No version conflicts, no dependency resolution
- **Portability**: Runs on any system with Python 3.8+
- **Auditability**: One file to review, one file to trust
- **Disaster recovery**: Trivial to reconstruct the build system

### Step 1: Parse Markdown into AST (Reading the Blueprints)

The first phase reads every markdown file in the content directory and parses
each one into a structured representation.

```
Input:  Raw markdown text
Output: Abstract Syntax Tree (structured document)

Processing:
  1. Extract front matter (YAML between --- delimiters)
  2. Tokenize markdown body into block elements
  3. Parse inline elements within each block
  4. Build heading hierarchy (document outline)
  5. Resolve internal cross-references
  6. Validate structure against expected patterns
```

The parser handles standard markdown plus the following extensions:

| Element | Syntax | AST Node Type |
|---------|--------|---------------|
| Fenced code blocks | ` ``` ` | `CodeBlock` |
| Tables | `\| pipe \| syntax \|` | `Table` |
| Task lists | `- [x] item` | `TaskList` |
| Definition lists | `Term\n: Definition` | `DefinitionList` |
| Admonitions | `> **Note:** text` | `Admonition` |
| Internal links | `[text](./path.md)` | `InternalLink` |

The parser explicitly **rejects**:
- Embedded `<script>` tags
- Inline `<style>` blocks with `@import`
- HTML `<iframe>` elements
- Any element with `src`, `href`, or `url()` pointing to external domains

This is not a limitation of the parser. This is a security boundary.

### Step 2: Apply Templates (Pouring the Concrete)

Each AST is fed into a template engine that produces the final HTML structure.
The template engine is a simple string-interpolation system built into `build.py`
--- not Jinja2, not Mako, not any external templating library.

```
Template structure:

<!DOCTYPE html>
<html lang="en" data-theme="dark">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{page_title} — Holm Intelligence Complex</title>
    <style>{inline_critical_css}</style>
    <link rel="stylesheet" href="{relative_path}/style.css">
    <link rel="icon" href="data:image/svg+xml,...">
</head>
<body>
    <a href="#main-content" class="skip-link">Skip to content</a>
    <header role="banner">{site_header}</header>
    <nav role="navigation" aria-label="Main">{navigation}</nav>
    <main id="main-content" role="main">
        <article>{rendered_content}</article>
    </main>
    <footer role="contentinfo">{site_footer}</footer>
</body>
</html>
```

Key template decisions:
- **Critical CSS is inlined** in the `<head>` for first-paint performance
- **The full stylesheet is linked** as a co-located file (not a CDN)
- **The favicon is a data URI** --- no additional HTTP request
- **No external resources** appear in any template

### Step 3: Generate Navigation (Installing Wayfinding Systems)

The build system generates navigation structures from the content directory
layout and front matter metadata. Navigation is **fully static** --- it is
baked into every HTML page at build time.

```
Navigation generation:

1. Scan all content files
2. Sort by front matter 'order' field (then alphabetically)
3. Build hierarchical navigation tree
4. For each page, render navigation with:
   - Current page highlighted (aria-current="page")
   - Sibling pages listed
   - Parent/child relationships shown
   - Relative paths calculated from current page location
```

The navigation tree is written into every page. There is no JavaScript-driven
menu loading, no AJAX calls to fetch nav data, no client-side routing. Every
page contains its own complete navigation context.

```html
<nav role="navigation" aria-label="Main navigation">
  <ul class="nav-tree">
    <li>
      <a href="../index.html">Home</a>
    </li>
    <li>
      <a href="../getting-started/index.html">Getting Started</a>
      <ul>
        <li aria-current="page">
          <a href="installation.html">Installation</a>
        </li>
        <li>
          <a href="quickstart.html">Quick Start</a>
        </li>
      </ul>
    </li>
    <!-- ... -->
  </ul>
</nav>
```

### Step 4: Write HTML Files (Completing Construction)

The compiled HTML pages are written to the output directory, maintaining the
same structure as the content source:

```
_site/                              # The completed building
├── index.html                      # Lobby
├── style.css                       # Main electrical system
├── getting-started/
│   ├── index.html                  # Wing entrance
│   ├── installation.html           # Room
│   ├── quickstart.html             # Room
│   └── configuration.html          # Room
├── architecture/
│   ├── index.html
│   └── ...
└── reference/
    ├── index.html
    └── ...
```

Each HTML file is:
- **Complete**: renders independently without any other file
- **Self-referencing**: all internal links use relative paths
- **Portable**: the entire `_site/` directory can be copied anywhere

### Step 5: Copy Static Assets (Furnishing the Rooms)

Static assets (the CSS file, any local images, SVG diagrams) are copied from
the `static/` source directory into the output:

```
static/                             _site/
├── style.css            ──>        ├── style.css
├── images/              ──>        ├── images/
│   └── diagram.svg                 │   └── diagram.svg
└── fonts/               ──>        └── fonts/
    └── (empty: we use               └── (empty: we use
         system fonts)                     system fonts)
```

Note the `fonts/` directory: it exists in the project structure as a reminder
that **we deliberately chose not to ship custom fonts**. System fonts render
instantly, require zero network transfer, and respect the user's platform
conventions.

### Build Execution: Zero Network, Zero Wait

The complete build runs with a single command:

```bash
python3 build.py
```

Typical build characteristics:

| Metric | Value |
|--------|-------|
| Network requests | 0 |
| External dependencies | 0 |
| Build time (50 pages) | < 2 seconds |
| Build time (500 pages) | < 15 seconds |
| Memory usage | < 100MB |
| Disk I/O | Sequential writes |
| CPU usage | Single-threaded, brief |

The build can be run:
- On a developer laptop
- On a CI server (with no internet)
- On the air-gapped target machine itself
- On a USB-bootable Linux environment
- Anywhere Python 3.8+ exists

---

## Markdown as Blueprint Language

### Why Markdown

Markdown was chosen as the source format for the HIC's content because it
satisfies every constraint of an air-gapped documentation system:

**Human-readable without tools.** A markdown file is a plain text file. Open it
in `cat`, `less`, `vim`, `nano`, Notepad, or any text editor on any operating
system. The content is immediately legible. No special software is required to
read or write documentation.

**Version-controllable.** Markdown files diff cleanly in Git. Changes are
visible line-by-line. Merge conflicts are resolvable by humans reading plain
text. The entire history of every document is preserved in the repository.

**Tool-agnostic.** Markdown is not owned by any vendor. It is not tied to any
specific editor, platform, or build system. If `build.py` were to be replaced
entirely, the markdown source files would remain valid and useful.

**Universally supported.** Every major code hosting platform renders markdown.
Every major text editor has markdown support. Every developer has encountered
markdown. The learning curve is measured in minutes.

**Constrained by design.** Markdown's limited feature set is a feature, not a
bug. You cannot embed a tracking pixel in markdown. You cannot create a
self-executing script. You cannot produce div soup. The format constrains
authors toward content, not presentation.

### Supported Markdown Extensions

The HIC parser supports a carefully chosen set of markdown extensions beyond
the base CommonMark specification:

#### Fenced Code Blocks with Language Hints

````markdown
```python
def build_site():
    """Compile all pages."""
    for source in get_sources():
        compile_page(source)
```
````

Code blocks are rendered with language-specific syntax classes for CSS-only
syntax highlighting. No JavaScript syntax highlighting libraries are loaded.

#### Tables

```markdown
| Column A | Column B | Column C |
|----------|----------|----------|
| Data 1   | Data 2   | Data 3   |
```

Tables render as proper `<table>` elements with `<thead>` and `<tbody>`.
They are wrapped in a scrollable container for mobile viewports.

#### Task Lists

```markdown
- [x] Completed item
- [ ] Pending item
- [ ] Another pending item
```

Task lists render as styled checkboxes. They are **not interactive** --- this
is static HTML. The checked/unchecked state is set at build time from the
markdown source.

#### Admonitions (Block Quotes with Type)

```markdown
> **Note:** This is important information that supplements the main text.

> **Warning:** This action is irreversible. Proceed with caution.

> **Danger:** This will destroy data. There is no undo.
```

Admonitions are parsed from specially-formatted blockquotes and rendered with
appropriate ARIA roles and visual styling.

#### Definition Lists

```markdown
Term
: The definition of the term, which may span
  multiple lines.

Another Term
: Another definition.
```

### Markdown Limitations (By Design)

The following are **intentionally not supported**:

| Feature | Reason for Exclusion |
|---------|---------------------|
| Embedded HTML `<script>` | Security: no executable content in source |
| Embedded HTML `<iframe>` | Air-gap: no external content embedding |
| `<img>` with external `src` | Air-gap: no external resource loading |
| HTML `<style>` blocks | Separation: styling belongs in CSS pipeline |
| Markdown `include` directives | Complexity: each file must be self-contained |
| LaTeX math rendering | Dependency: requires external library |
| Mermaid diagrams | Dependency: requires JavaScript runtime |

For diagrams, the HIC uses pre-rendered SVG files checked into the repository.
For mathematical notation, plain text or pre-rendered images are used. These
constraints ensure that the build pipeline remains dependency-free.

### The Translation Table: Markdown to Building Components

Every markdown element maps to a structural component of the HIC skyscraper:

```
MARKDOWN ELEMENT          BUILDING COMPONENT         HTML OUTPUT
─────────────────────────────────────────────────────────────────
# Heading 1               Floor / Level              <h1>
## Heading 2              Wing / Section             <h2>
### Heading 3             Room / Subsection          <h3>
#### Heading 4            Alcove / Detail Area       <h4>
##### Heading 5           Shelf / Minor Detail       <h5>
###### Heading 6          Label / Finest Detail      <h6>

Paragraph                 Wall Plaque / Sign         <p>
Unordered list            Directory Board            <ul><li>
Ordered list              Sequential Signage         <ol><li>
Blockquote                Framed Notice              <blockquote>
Code block                Technical Blueprint        <pre><code>
Inline code               Label Plate                <code>
Table                     Data Display Panel         <table>
Horizontal rule           Floor Divider              <hr>
Link                      Corridor / Doorway         <a>
Image                     Window / Viewport          <img> (local)
Bold                      Highlighted Text           <strong>
Italic                    Annotated Text             <em>
```

### Heading Levels as Structural Elements

The heading hierarchy is not merely semantic decoration. It defines the
physical structure of each floor in the HIC:

```
# Installation Guide                    ← THE FLOOR (h1)
                                          One per page. The floor's name.

## Prerequisites                        ← WING (h2)
                                          Major section. A corridor of
                                          the floor dedicated to one topic.

### System Requirements                 ← ROOM (h3)
                                          A specific room within the wing.
                                          Self-contained topic area.

#### Operating System Details            ← ALCOVE (h4)
                                          A detail area within a room.
                                          Supporting information.
```

Rules enforced by the build system:
- **Exactly one `h1` per page.** A floor has one name.
- **No skipped heading levels.** You cannot have an `h1` followed by an `h3`.
  You must pass through the wing (`h2`) to reach a room (`h3`).
- **Heading text must be unique within a page.** Every room on a floor has a
  distinct name for navigation purposes.

---

## HTML as Construction Material

### Semantic HTML: The Right Material for Each Component

The HIC does not use `<div>` as a universal building block. Each structural
component uses the HTML element that carries the correct semantic meaning:

```html
<!-- The building entrance -->
<header role="banner">
  <a href="/" class="hic-logo" aria-label="Holm Intelligence Complex — Home">
    <!-- Inline SVG logo -->
  </a>
  <p class="hic-tagline">Documentation Infrastructure</p>
</header>

<!-- The directory / wayfinding system -->
<nav role="navigation" aria-label="Main navigation">
  <ul class="nav-tree">...</ul>
</nav>

<!-- The current room -->
<main id="main-content" role="main">
  <article>
    <h1>Page Title</h1>
    <!-- Content rendered from markdown -->
  </article>
</main>

<!-- Building information plaque -->
<footer role="contentinfo">
  <p>Holm Intelligence Complex — Air-gapped documentation system</p>
</footer>
```

### No Div Soup: Meaningful Tags for Every Element

Every HTML element in the HIC output carries semantic weight:

| Component | HTML Element | Why This Element |
|-----------|-------------|------------------|
| Page header | `<header>` | Introductory content for the page |
| Navigation | `<nav>` | Navigation landmark |
| Main content | `<main>` | Primary content of the page |
| Article body | `<article>` | Self-contained composition |
| Content sections | `<section>` | Thematic grouping with heading |
| Sidebar | `<aside>` | Tangentially related content |
| Page footer | `<footer>` | Footer for the page |
| Code examples | `<pre><code>` | Preformatted code content |
| Data tables | `<table>` with `<caption>` | Tabular data with description |
| Term definitions | `<dl><dt><dd>` | Definition list structure |
| Timestamps | `<time>` | Machine-readable date/time |
| Abbreviations | `<abbr>` | Abbreviation with expansion |

Elements that **never appear** in HIC output:
- `<div>` used as a semantic substitute (only as a non-semantic wrapper when no better element exists)
- `<span>` without a clear purpose
- `<b>` or `<i>` (use `<strong>` and `<em>` for semantic emphasis)
- `<br>` in running text (use proper paragraph breaks)
- `<font>`, `<center>`, or any deprecated presentational element

### Accessibility Built into the Structure

Accessibility in the HIC is not an afterthought or a compliance checkbox. It
is built into the HTML structure itself, because accessible HTML is simply
**correct HTML**.

```html
<!-- Skip link: first element in body -->
<a href="#main-content" class="skip-link">Skip to main content</a>

<!-- ARIA landmarks mirror the building's physical layout -->
<header role="banner">          <!-- Building entrance sign -->
<nav role="navigation"
     aria-label="Main">         <!-- Directory board -->
<main role="main"
      id="main-content">        <!-- The room you're visiting -->
<footer role="contentinfo">     <!-- Building information desk -->

<!-- Headings form a navigable outline -->
<h1>Floor Name</h1>
  <h2>Wing Name</h2>
    <h3>Room Name</h3>

<!-- Tables include captions and scope -->
<table>
  <caption>System requirements by platform</caption>
  <thead>
    <tr>
      <th scope="col">Platform</th>
      <th scope="col">Minimum Version</th>
    </tr>
  </thead>
  <tbody>...</tbody>
</table>

<!-- Images include descriptive alt text -->
<img src="diagram.svg" alt="Architecture diagram showing the three
  main subsystems and their communication paths">

<!-- Navigation indicates current location -->
<a href="current-page.html" aria-current="page">Current Page</a>

<!-- Interactive elements are keyboard accessible -->
<!-- (All navigation uses standard <a> tags with href) -->
```

### The HTML Is the Building

The ultimate test of HIC's HTML: **disable CSS and JavaScript entirely**.

What remains must be:
- **Readable**: All text content is present and legible
- **Navigable**: All links work, all navigation is functional
- **Structured**: Headings, lists, and tables convey organization
- **Complete**: No content is hidden behind JavaScript-only interactions
- **Ordered**: Content appears in a logical reading order

This is not a theoretical exercise. This is the actual state of the HIC
during a "power outage" --- when CSS files fail to load, when JavaScript is
blocked, when the browser is minimal. The building stands.

---

## CSS as Neon Lighting

### Custom Properties: The Building's Electrical Wiring

The HIC's visual system is built on CSS custom properties (variables), which
serve as the building's central electrical wiring. Change a wire, change
every light connected to it.

```css
:root {
  /* === STRUCTURAL PALETTE === */
  /* The building's core colors: concrete and steel */
  --hic-bg-primary: #0a0a0f;
  --hic-bg-secondary: #12121a;
  --hic-bg-tertiary: #1a1a2e;
  --hic-text-primary: #e0e0e8;
  --hic-text-secondary: #a0a0b0;
  --hic-text-muted: #606070;

  /* === NEON PALETTE === */
  /* The building's signature lighting */
  --hic-neon-cyan: #00fff0;
  --hic-neon-magenta: #ff00aa;
  --hic-neon-amber: #ffaa00;
  --hic-neon-green: #00ff88;
  --hic-neon-red: #ff3355;

  /* === STRUCTURAL DIMENSIONS === */
  --hic-content-width: 72ch;
  --hic-nav-width: 260px;
  --hic-header-height: 60px;
  --hic-spacing-unit: 0.5rem;

  /* === TYPOGRAPHY === */
  --hic-font-body: system-ui, -apple-system, "Segoe UI", Roboto,
                   "Helvetica Neue", Arial, sans-serif;
  --hic-font-mono: ui-monospace, "Cascadia Code", "Source Code Pro",
                   Menlo, Consolas, "DejaVu Sans Mono", monospace;
  --hic-font-size-base: 1rem;
  --hic-line-height-base: 1.6;

  /* === NEON GLOW EFFECTS === */
  --hic-glow-cyan: 0 0 10px rgba(0, 255, 240, 0.3);
  --hic-glow-magenta: 0 0 10px rgba(255, 0, 170, 0.3);
}
```

### Dark Theme as Default: The Building Operates at Night

The HIC defaults to a dark color scheme. This is not merely an aesthetic choice
--- it reflects the operational reality of the system:

1. **The building is underground.** Air-gapped systems operate in secure
   facilities. Dark themes reduce eye strain in low-ambient-light environments.
2. **Neon glows best in darkness.** The cyberpunk aesthetic requires a dark
   canvas for its signature lighting effects to register.
3. **Dark is the new default.** Modern operating systems and browsers
   increasingly default to dark mode. The HIC follows the user's preference
   via `prefers-color-scheme`.

```css
/* Dark theme is the default (no media query needed) */
body {
  background-color: var(--hic-bg-primary);
  color: var(--hic-text-primary);
}

/* Light theme for users who prefer it, or for print */
@media (prefers-color-scheme: light) {
  :root {
    --hic-bg-primary: #fafafa;
    --hic-bg-secondary: #f0f0f2;
    --hic-bg-tertiary: #e4e4e8;
    --hic-text-primary: #1a1a2e;
    --hic-text-secondary: #404058;
    --hic-text-muted: #808098;

    /* Neon tones shift for legibility on light backgrounds */
    --hic-neon-cyan: #0088aa;
    --hic-neon-magenta: #cc0088;
  }
}
```

### No External Fonts: System Font Stacks Only

The HIC loads **zero font files**. Every typeface used is already installed
on the user's operating system.

```
SYSTEM FONT STACK RESOLUTION:

Windows:     Segoe UI → fallback chain
macOS:       San Francisco (system-ui) → fallback chain
Linux:       System default → Roboto / Noto → fallback chain
Android:     Roboto → fallback chain
iOS:         San Francisco → fallback chain

MONOSPACE STACK RESOLUTION:

Windows:     Cascadia Code → Consolas → fallback chain
macOS:       SF Mono (ui-monospace) → Menlo → fallback chain
Linux:       Source Code Pro → DejaVu Sans Mono → fallback chain
```

Benefits of system fonts in an air-gapped context:
- **Zero download**: no font files to transfer via USB or network
- **Zero latency**: fonts are already cached in the OS
- **Zero FOIT/FOUT**: no flash of invisible or unstyled text
- **Zero maintenance**: font updates are handled by the OS vendor
- **Full air-gap compliance**: no external font service contact

### No Images Required

The HIC's visual design requires no raster images. All visual elements are
created with:

- **CSS**: borders, backgrounds, gradients, shadows, shapes
- **Inline SVG**: icons, diagrams, logos (see SVG Pipeline section)
- **Unicode**: symbols and special characters where appropriate

```css
/* Neon border effect - pure CSS, no image */
.hic-card {
  border: 1px solid var(--hic-neon-cyan);
  box-shadow: var(--hic-glow-cyan),
              inset 0 0 20px rgba(0, 255, 240, 0.05);
}

/* Decorative separator - pure CSS */
.hic-divider::after {
  content: "";
  display: block;
  width: 60%;
  height: 1px;
  margin: 2rem auto;
  background: linear-gradient(
    90deg,
    transparent,
    var(--hic-neon-cyan),
    transparent
  );
}
```

### Print Stylesheet: The Building by Daylight

The HIC includes a print stylesheet that transforms the dark cyberpunk
aesthetic into a clean, ink-efficient format:

```css
@media print {
  :root {
    --hic-bg-primary: #ffffff;
    --hic-text-primary: #000000;
    --hic-neon-cyan: #000000;
  }

  /* Remove neon effects */
  * {
    box-shadow: none !important;
    text-shadow: none !important;
  }

  /* Hide navigation (not useful in print) */
  nav, .skip-link, .hic-search {
    display: none !important;
  }

  /* Ensure content fills the page */
  main {
    width: 100% !important;
    margin: 0 !important;
    padding: 0 !important;
  }

  /* Handle page breaks intelligently */
  h1, h2, h3 {
    page-break-after: avoid;
  }

  pre, blockquote, table {
    page-break-inside: avoid;
  }

  /* Show link URLs for reference */
  a[href^="http"]::after {
    content: " (" attr(href) ")";
    font-size: 0.85em;
    color: #666;
  }
}
```

### CSS Layering: Structured Electrical Systems

The HIC's CSS is organized in layers, each building upon the previous:

```
Layer 1: BASE          Reset, box-sizing, root variables
                       The building's foundation wiring.

Layer 2: LAYOUT        Grid structure, navigation layout, content area
                       The conduit runs through walls and ceilings.

Layer 3: COMPONENTS    Cards, code blocks, tables, admonitions
                       Individual fixtures: sconces, panels, displays.

Layer 4: THEME         Colors, shadows, neon effects
                       The neon tubes and ambient glow.

Layer 5: UTILITIES     Spacing helpers, visibility, screen-reader-only
                       Switches and dimmers for fine control.

Each layer depends only on layers below it.
No circular dependencies. No specificity wars.
```

The entire CSS system compiles to a single file under 30KB. There is no CSS
preprocessor (no Sass, no Less, no PostCSS). The source CSS is the output CSS.
Custom properties provide all the abstraction needed. The fewer tools in the
chain, the fewer tools that can break.

---

## Zero-JavaScript Guarantee

### The Building Operates During a Power Outage

The HIC makes a hard guarantee: **every feature that appears in the navigation
and content of the site works with JavaScript completely disabled.**

This is not "graceful degradation." This is not "we try our best." This is a
contractual obligation of the build system. If a feature cannot work without
JavaScript, it does not ship.

### Why Zero-JavaScript Matters

In an air-gapped environment, JavaScript is a liability:

- **Attack surface**: Every line of JavaScript is a potential vector for
  code injection, data exfiltration, or unexpected behavior
- **Audit burden**: JavaScript must be reviewed line-by-line for security;
  the less there is, the less there is to audit
- **Failure modes**: JavaScript can fail silently, producing a page that
  looks correct but is missing content or functionality
- **Reproducibility**: JavaScript execution depends on browser engine, version,
  and configuration --- static HTML does not
- **Accessibility**: Screen readers and assistive technologies handle static
  HTML more reliably than JavaScript-driven interfaces

The zero-JavaScript guarantee eliminates an entire category of risk.

### Navigation Without JavaScript

All navigation in the HIC uses standard HTML mechanisms:

```html
<!-- Page-to-page navigation: standard anchor tags -->
<a href="../getting-started/installation.html">Installation Guide</a>

<!-- On-page navigation: fragment identifiers -->
<a href="#prerequisites">Jump to Prerequisites</a>

<!-- Table of contents: list of fragment links -->
<nav aria-label="Table of contents">
  <ol>
    <li><a href="#overview">Overview</a></li>
    <li><a href="#prerequisites">Prerequisites</a></li>
    <li><a href="#installation">Installation</a></li>
    <li><a href="#verification">Verification</a></li>
  </ol>
</nav>

<!-- Previous/Next navigation: static links -->
<nav aria-label="Page navigation">
  <a href="quickstart.html" rel="prev">Previous: Quick Start</a>
  <a href="configuration.html" rel="next">Next: Configuration</a>
</nav>
```

No JavaScript router. No client-side navigation. No history API manipulation.
Click a link. Load a page. That is how the web works.

### Server-Rendered Table of Contents

The table of contents for each page is generated at build time by `build.py`.
It is baked into the HTML of every page as a static `<nav>` element.

```
Build time:
  1. Parse headings from AST
  2. Generate slug for each heading (id attribute)
  3. Build nested list structure
  4. Inject into page template

Result: A complete, static, clickable table of contents
        with zero JavaScript involvement.
```

### Search: Pre-Indexed or External

Search is the one feature that traditionally requires JavaScript. The HIC
handles this through a layered approach:

**Tier 1: No-JS search (default)**
- A pre-built index page listing all content with descriptions
- Browser's built-in `Ctrl+F` / `Cmd+F` for in-page search
- Structured navigation that serves as a browsable index

**Tier 2: Pre-built search index (progressive enhancement)**
- At build time, a JSON search index is generated
- If JavaScript is available, a lightweight client-side search loads this index
- The search index is a local file --- no external search service
- If JavaScript is unavailable, users fall back to Tier 1

**Tier 3: External tool integration (optional)**
- `grep -r` across the `_site/` directory
- OS-level desktop search indexing
- These operate outside the browser entirely

```
Search Tier Selection:

  JavaScript enabled?
  ├── Yes → Load pre-built JSON index, client-side search
  └── No → Fall back to browsable index page + Ctrl+F
```

### Interactive Features Degrade to Static Equivalents

Any future interactive enhancement must specify its static fallback:

| Interactive Feature | Static Fallback |
|---------------------|-----------------|
| Collapsible sections | All sections expanded by default |
| Copy-to-clipboard on code blocks | Users select and copy manually |
| Theme toggle button | Respects `prefers-color-scheme` automatically |
| Smooth scroll to anchors | Browser's native fragment navigation |
| Search-as-you-type | Pre-built index page |
| Syntax highlighting (JS) | CSS-only language classes on `<code>` |
| Mobile menu toggle | Navigation visible by default, CSS-only collapse |

### The Progressive Enhancement Contract

Any JavaScript added to the HIC in the future must satisfy ALL of the following:

1. **The page is fully functional without it.** No content is hidden behind
   JS-only interactions. No navigation requires JS to operate.
2. **It loads from a local file.** No CDN, no external script source. The
   script is a file in the `_site/` directory.
3. **It has no external dependencies.** No npm packages, no frameworks, no
   libraries that require their own dependencies.
4. **It fails silently.** If the script errors, the page continues to work.
   No white screens. No missing content. No broken navigation.
5. **It is optional.** A build flag can exclude all JavaScript from output.
   The site remains fully functional.
6. **It is small.** Any single script must be under 10KB minified.
7. **It is auditable.** The complete source is in the repository, readable
   by a human in under 10 minutes.
8. **It makes no network requests.** No fetch, no XHR, no WebSocket, no
   beacon, no image ping. Nothing leaves the machine.

---

## Air-Gap Compliance Checklist

### The Compliance Matrix

The following checklist must pass before any HIC build is deployed to the
air-gapped target environment. Failure of any single item blocks deployment.

```
AIR-GAP COMPLIANCE VERIFICATION
================================

[ ] 1. NO EXTERNAL URLs IN PRODUCTION HTML
       Scan: grep -rn "https\?://" _site/ --include="*.html"
       Expected: Zero matches in src, href, url(), @import attributes
       (Matches in visible prose content are acceptable)

[ ] 2. NO CDN DEPENDENCIES
       Scan: grep -rn "cdn\.\|cdnjs\.\|unpkg\.\|jsdelivr\." _site/
       Expected: Zero matches

[ ] 3. NO EXTERNAL FONT LOADING
       Scan: grep -rn "fonts.googleapis\|fonts.gstatic\|use.typekit" _site/
       Expected: Zero matches

[ ] 4. NO ANALYTICS OR TRACKING
       Scan: grep -rn "google-analytics\|gtag\|analytics\.js\|pixel" _site/
       Expected: Zero matches in script contexts

[ ] 5. NO API CALLS OR FETCH REQUESTS
       Scan: grep -rn "fetch(\|XMLHttpRequest\|\.ajax\|axios" _site/
       Expected: Zero matches

[ ] 6. NO WEBSOCKET CONNECTIONS
       Scan: grep -rn "WebSocket\|wss:\|ws:" _site/
       Expected: Zero matches

[ ] 7. NO SERVICE WORKER HARD DEPENDENCIES
       Verify: Site functions identically with service worker
               unregistered or blocked

[ ] 8. NO EXTERNAL IMAGES
       Scan: grep -rn 'src="http' _site/ --include="*.html"
       Expected: Zero matches

[ ] 9. NO EXTERNAL STYLESHEETS
       Scan: grep -rn 'link.*href="http' _site/ --include="*.html"
       Expected: Zero matches

[ ] 10. NO EXTERNAL SCRIPTS
        Scan: grep -rn 'script.*src="http' _site/ --include="*.html"
        Expected: Zero matches

[ ] 11. CONTENT SECURITY POLICY COMPATIBLE
        Verify: All resources load under:
        Content-Security-Policy: default-src 'self'; style-src 'self'
                                 'unsafe-inline'; img-src 'self' data:

[ ] 12. ALL INTERNAL LINKS RESOLVE
        Verify: Every href in the site points to an existing file
        Tool: python3 build.py --verify-links
```

### Automated Compliance Scanning

The build system includes a verification mode that runs all compliance checks
automatically:

```bash
python3 build.py --verify

# Output:
# [PASS] No external URLs in HTML attributes
# [PASS] No CDN references detected
# [PASS] No external font loading
# [PASS] No analytics or tracking scripts
# [PASS] No fetch/XHR/AJAX calls
# [PASS] No WebSocket references
# [PASS] No external image sources
# [PASS] No external stylesheet links
# [PASS] No external script sources
# [PASS] CSP compatibility verified
# [PASS] All internal links resolve (247/247)
#
# AIR-GAP COMPLIANCE: PASSED (12/12 checks)
```

This verification runs as part of every build. A failed check halts the build
and produces an error with the exact file and line number of the violation.

### Content Security Policy

The HIC is designed to operate under a strict Content Security Policy. When
served by a web server, the following CSP header should be applied:

```
Content-Security-Policy:
  default-src 'none';
  script-src 'self';
  style-src 'self' 'unsafe-inline';
  img-src 'self' data:;
  font-src 'self';
  connect-src 'none';
  frame-src 'none';
  object-src 'none';
  base-uri 'self';
  form-action 'self';
```

Explanation of each directive:

| Directive | Value | Reason |
|-----------|-------|--------|
| `default-src` | `'none'` | Deny everything by default |
| `script-src` | `'self'` | Only local scripts (if any exist) |
| `style-src` | `'self' 'unsafe-inline'` | Local CSS + critical inline styles |
| `img-src` | `'self' data:` | Local images + SVG data URIs |
| `font-src` | `'self'` | Local fonts only (currently none) |
| `connect-src` | `'none'` | No XHR/fetch/WebSocket at all |
| `frame-src` | `'none'` | No iframes of any kind |
| `object-src` | `'none'` | No plugins, no embeds |

The `'unsafe-inline'` for styles is required for the critical CSS inlined in
the `<head>`. This is an acceptable tradeoff: inline styles cannot exfiltrate
data, and `connect-src: 'none'` prevents any network communication regardless.

### Why Each Check Matters

Each compliance check maps to a specific threat vector in an air-gapped context:

| Check | Threat Mitigated |
|-------|-----------------|
| No external URLs | Data exfiltration via resource requests |
| No CDN deps | Supply chain compromise via CDN poisoning |
| No external fonts | Tracking via font loading fingerprinting |
| No analytics | Surveillance and usage profiling |
| No fetch/XHR | Data exfiltration via HTTP requests |
| No WebSocket | Persistent covert communication channels |
| No external images | Tracking pixels and data exfiltration |
| No external CSS | CSS-based data exfiltration techniques |
| No external scripts | Arbitrary code execution from untrusted sources |
| CSP enforcement | Defense-in-depth against injection attacks |
| Link verification | Broken references that might prompt external lookups |

---

## The SVG Pipeline

### Icons and Diagrams as Inline SVG

Every icon and diagram in the HIC is rendered as inline SVG directly in the
HTML. No external SVG files are loaded via `<img>` tags for small icons. No
icon font libraries. No sprite sheets loaded from external sources.

```html
<!-- Inline SVG icon example -->
<svg xmlns="http://www.w3.org/2000/svg"
     width="20" height="20"
     viewBox="0 0 24 24"
     fill="none"
     stroke="currentColor"
     stroke-width="2"
     stroke-linecap="round"
     stroke-linejoin="round"
     role="img"
     aria-hidden="true">
  <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
  <polyline points="9 22 9 12 15 12 15 22"/>
</svg>
```

Key principles for HIC SVG:
- **`xmlns` is always included** for compatibility outside HTML5 contexts
- **`viewBox` is always set** for proper scaling at any size
- **`currentColor` is used for stroke/fill** so icons inherit text color
  from their parent element and respond to theme changes automatically
- **Decorative icons use `aria-hidden="true"`** to hide from screen readers
- **Meaningful icons use `role="img"` and `<title>`** for accessibility
- **No external references**: no `xlink:href` to external files, no `use`
  elements pointing to external sprite sheets

### SVG Construction Standards

All SVG in the HIC follows strict construction rules:

```
SVG CONSTRUCTION CHECKLIST:

[ ] No external references (xlink:href, use with external URL)
[ ] No embedded scripts (<script> inside SVG)
[ ] No event handlers (onclick, onload, onmouseover, etc.)
[ ] No external stylesheets (<style> with @import)
[ ] No raster image embeds (<image> with external href)
[ ] No foreignObject elements (potential XSS vector)
[ ] Uses viewBox for responsive scaling
[ ] Uses currentColor where appropriate for theme adaptation
[ ] Minimal path data (optimized, no unnecessary decimal precision)
[ ] Accessible: aria-hidden="true" for decorative,
    role="img" + <title> for meaningful
[ ] No unnecessary metadata (editor artifacts, comments, defaults)
[ ] Valid XML structure (parseable by any XML processor)
```

### The Favicon: Inline SVG Data URI

The HIC favicon is an SVG data URI embedded directly in the HTML `<head>`.
Zero additional HTTP requests. Zero external files.

```html
<link rel="icon" type="image/svg+xml" href="data:image/svg+xml,
  %3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 32 32'%3E
    %3Crect width='32' height='32' rx='4' fill='%230a0a0f'/%3E
    %3Ctext x='16' y='22' text-anchor='middle'
      font-family='monospace' font-size='18'
      fill='%2300fff0'%3EH%3C/text%3E
  %3C/svg%3E">
```

This renders as a dark rounded square with a cyan "H" --- the HIC logo mark
--- directly in the browser tab. It loads instantly because it is part of the
HTML document itself. No network request. No file lookup. The favicon travels
with the page.

### Floor Plan Diagrams: SVG Blueprints

Larger diagrams (such as architecture overviews and floor plans) are stored
as standalone `.svg` files in the `static/images/` directory and referenced
with relative paths:

```html
<figure>
  <img src="../images/architecture-overview.svg"
       alt="HIC architecture overview showing the three main layers:
            content source, build pipeline, and static output"
       loading="lazy"
       width="800"
       height="400">
  <figcaption>HIC Architecture Overview</figcaption>
</figure>
```

These SVG files are:
- **Checked into the Git repository** (version-controlled alongside content)
- **Self-contained** (no external references within the SVG)
- **Optimized** (run through manual optimization or svgo before commit)
- **Styled via CSS custom properties** where possible (using `var()` in
  inline SVG, or class-based styling for referenced SVG)
- **Sized with explicit dimensions** to prevent layout shift during load

### The Building Logo

The HIC logo is a full SVG that represents the skyscraper silhouette with
neon accents. It is inlined in the header template:

```html
<a href="/" class="hic-logo" aria-label="Holm Intelligence Complex — Home">
  <svg xmlns="http://www.w3.org/2000/svg"
       viewBox="0 0 120 40"
       width="120"
       height="40"
       role="img">
    <title>Holm Intelligence Complex</title>
    <!-- Skyscraper silhouette -->
    <rect x="10" y="5" width="20" height="35" fill="var(--hic-bg-tertiary)"/>
    <rect x="12" y="8" width="4" height="3" fill="var(--hic-neon-cyan)"
          opacity="0.8"/>
    <rect x="18" y="8" width="4" height="3" fill="var(--hic-neon-cyan)"
          opacity="0.6"/>
    <rect x="12" y="14" width="4" height="3" fill="var(--hic-neon-cyan)"
          opacity="0.7"/>
    <rect x="18" y="14" width="4" height="3" fill="var(--hic-neon-magenta)"
          opacity="0.5"/>
    <rect x="12" y="20" width="4" height="3" fill="var(--hic-neon-cyan)"
          opacity="0.9"/>
    <rect x="18" y="20" width="4" height="3" fill="var(--hic-neon-cyan)"
          opacity="0.4"/>
    <!-- HIC text -->
    <text x="40" y="28" font-family="var(--hic-font-mono)"
          font-size="16" fill="var(--hic-neon-cyan)"
          font-weight="bold">HIC</text>
  </svg>
</a>
```

The logo uses CSS custom properties for its colors, meaning it automatically
adapts to the current theme. In dark mode, the neon glows against a dark
silhouette. In light mode (or print), the colors shift to maintain contrast
and legibility.

---

## USB Import Protocol

### How Content Enters the Air-Gapped System

The HIC operates in an air-gapped environment. By definition, it has no
network connection to the outside world. New content enters the system via
physical USB storage devices, following a strict import protocol.

This is not optional security theater. This is the only way content moves
into the building. There is no back door. There is no "just this once"
network cable. The USB import protocol is the front door of the HIC, and it
has a guard, a scanner, and a quarantine room.

```
CONTENT FLOW:

  External World          Airlock           HIC Internal
  ──────────────         ─────────         ──────────────
  Author writes    →     USB device    →   Content staging
  markdown on            scanned and       area (quarantine)
  connected              validated              │
  workstation                                   v
                                          Validation checks
                                               │
                                               v
                                          Import approved
                                               │
                                               v
                                          Git commit to
                                          internal repo
                                               │
                                               v
                                          Build triggered
                                               │
                                               v
                                          Site deployed
```

### USB Device Scanning and Quarantine

Every USB device that enters the air-gapped perimeter passes through a
scanning procedure before its contents are accessed:

```
USB SCANNING PROTOCOL:

1. PHYSICAL INSPECTION
   - Device is a known/approved make and model
   - No physical tampering evident
   - Device is labeled with origin and handler
   - Device serial number is logged

2. QUARANTINE MOUNT
   - Mount device read-only on a dedicated scanning workstation
   - Scanning workstation is isolated from the HIC build system
   - Mount options: noexec, nosuid, nodev, ro
   - No auto-run, no auto-mount, no thumbnail generation

3. AUTOMATED SCAN
   - File type verification (only allowed extensions)
   - Content scanning for embedded executables
   - Filename validation (no path traversal attempts)
   - Size limit enforcement per file and total
   - Binary content detection in text files
   - UTF-8 encoding verification for text files

4. MANIFEST COMPARISON
   - Device must contain an import-manifest.json
   - Manifest lists every file with SHA-256 checksum
   - Actual files are checked against manifest
   - Any mismatch halts the import
   - Extra files not in manifest halt the import
```

### File Format Validation

Only specific file formats are accepted through the USB import process:

| Allowed Extension | Purpose | Validation |
|-------------------|---------|------------|
| `.md` | Documentation content | UTF-8 text, valid markdown, no HTML scripts |
| `.svg` | Diagrams and icons | Valid XML, no `<script>`, no external refs |
| `.json` | Configuration/metadata | Valid JSON, schema-validated |
| `.txt` | Plain text notes | UTF-8 text only |
| `.png` | Screenshots (rare) | Valid PNG header, under 500KB |

Explicitly **rejected** formats:
- `.exe`, `.bat`, `.sh`, `.py`, `.js` --- no executables of any kind
- `.html` --- HTML is generated by the build system, not imported
- `.css` --- styles are managed within the build system
- `.zip`, `.tar`, `.gz` --- no archives (to prevent hidden content)
- `.doc`, `.docx`, `.pdf` --- use markdown only
- `.woff`, `.woff2`, `.ttf`, `.otf` --- no font files (system fonts only)
- Any file without an extension
- Any file with a double extension (e.g., `file.md.exe`)

### The Airlock Concept: Staging Area for New Content

Imported content does not go directly into the content source directory.
It enters a staging area --- the "airlock" --- where it is reviewed before
integration:

```
import-staging/                     ← THE AIRLOCK
├── incoming/                       ← Raw files from USB
│   ├── import-manifest.json        ← What was sent
│   └── content/
│       ├── new-guide.md
│       └── images/
│           └── new-diagram.svg
├── validated/                      ← Files that passed checks
│   ├── validation-report.json      ← Check results
│   └── content/
│       ├── new-guide.md
│       └── images/
│           └── new-diagram.svg
└── rejected/                       ← Files that failed checks
    ├── rejection-report.json       ← Failure reasons
    └── (failed files, if any)
```

The airlock workflow:

1. **Receive**: Files are copied from USB to `incoming/`
2. **Validate**: Automated checks run against every file
3. **Report**: Validation results are written to a report file
4. **Review**: A human reviews the validation report and samples content
5. **Accept or Reject**: Files move to `validated/` or `rejected/`
6. **Integrate**: Accepted files are copied to the content source directory
7. **Commit**: Changes are committed to the internal Git repository
8. **Build**: The build pipeline runs, producing updated static output
9. **Verify**: The new build is checked against compliance requirements
10. **Deploy**: The updated site is deployed to the serving infrastructure
11. **Clear**: The airlock is emptied after successful integration

No file ever bypasses the airlock. Even urgent content follows this process.
The process can be fast (minutes for a single file), but it cannot be skipped.

### Checksum Verification

Every file in the import process is checksummed to ensure integrity:

```json
{
  "import_id": "2026-02-17-001",
  "origin": "documentation-team-alpha",
  "handler": "operator-jsmith",
  "created": "2026-02-17T09:30:00Z",
  "files": [
    {
      "path": "content/new-guide.md",
      "sha256": "a1b2c3d4e5f6...",
      "size_bytes": 4523,
      "type": "text/markdown"
    },
    {
      "path": "images/new-diagram.svg",
      "sha256": "f6e5d4c3b2a1...",
      "size_bytes": 12870,
      "type": "image/svg+xml"
    }
  ],
  "total_files": 2,
  "total_size_bytes": 17393,
  "manifest_sha256": "9f8e7d6c5b4a..."
}
```

Verification steps:
1. Compute SHA-256 of each file on the USB device
2. Compare against checksums in the manifest
3. Compute SHA-256 of the manifest file itself
4. Compare against the `manifest_sha256` field (provided out-of-band,
   for example on a separate printed verification slip)
5. Verify file count matches manifest declaration
6. Verify total byte count matches manifest declaration
7. Any mismatch aborts the entire import --- no partial imports

### Import Logging

Every import operation is logged to an append-only import register:

```
IMPORT REGISTER ENTRY:

Import ID:     2026-02-17-001
Timestamp:     2026-02-17T10:15:33Z
Operator:      jsmith
USB Device:    Kingston-DT-32GB-SN7742
Origin:        documentation-team-alpha
Files:         2 (1 markdown, 1 SVG)
Total Size:    17,393 bytes
Validation:    PASSED (all checks)
Review:        APPROVED by operator jsmith
Integration:   COMMITTED as git hash a1b2c3d
Build:         SUCCESSFUL, 248 pages generated
Compliance:    PASSED (12/12 checks)
Status:        COMPLETE
```

The import register is itself a file within the air-gapped system, tracked
in Git, providing a complete audit trail of everything that has ever entered
the HIC. The register is append-only: entries are never modified or deleted.
Every piece of content in the building can be traced back to the USB device
and operator that brought it in.

---

## Build Reproducibility

### Same Inputs, Same Outputs --- Always

The HIC build pipeline is **deterministic**. Given the same source files and
the same version of `build.py`, the output is byte-for-byte identical
regardless of:

- What time the build runs
- What machine the build runs on
- What operating system the build runs on
- How many times the build has run before
- What timezone the machine is set to
- What locale the machine is configured for
- What other software is installed on the machine

This is not aspirational. This is enforced by design decisions at every level
of the build system.

### No Timestamps in Output

The most common source of non-deterministic builds is timestamps. The HIC
eliminates them:

```
WHAT THE HIC DOES NOT INCLUDE IN OUTPUT:

- "Generated on [date]" footers
- "Last modified" timestamps derived from system clock
- Cache-busting query strings with timestamps (?v=1708123456)
- Build timestamps in HTML comments
- Date-based file naming
- "Page generated in X.XXs" performance tags
- Any output that depends on datetime.now()
```

If a "last modified" date is needed for content, it comes from the **Git
history** of the source markdown file, not from the system clock at build
time. This ensures the date is tied to the content change, not to when
someone happened to run the build.

### Git Hash as Version Identifier

Instead of timestamps or sequential version numbers, the HIC uses the Git
commit hash of the source repository as its version identifier:

```html
<!-- In the HTML output footer -->
<footer role="contentinfo">
  <p>Holm Intelligence Complex</p>
  <p class="hic-version">Build: <code>a1b2c3d</code></p>
</footer>
```

The Git hash is:
- **Deterministic**: same source state always produces the same hash
- **Verifiable**: you can check what source state produced a given build
- **Unique**: every distinct source state has a distinct identifier
- **Immutable**: once committed, the hash never changes
- **Cross-referenceable**: links the deployed site to a specific repo state

### Build Manifest

Every build produces a manifest file (`_site/build-manifest.json`) recording
exactly what was built:

```json
{
  "build_version": "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
  "builder_version": "build.py@e4f5a6b7",
  "python_version": "3.11.5",
  "source_file_count": 52,
  "output_file_count": 58,
  "pages": [
    {
      "source": "content/index.md",
      "output": "_site/index.html",
      "source_sha256": "abc123...",
      "output_sha256": "def456...",
      "output_size_bytes": 24680
    }
  ],
  "total_output_size_bytes": 1284500,
  "compliance_check": "PASSED",
  "compliance_checks_passed": 12,
  "compliance_checks_total": 12
}
```

The manifest enables:
- **Verification**: confirm a deployed site matches a known build
- **Diffing**: compare two manifests to see what changed between builds
- **Auditing**: trace any output file back to its source
- **Rollback**: identify exactly which build to restore
- **Monitoring**: track site size growth over time

### Diff-Friendly Output

The HIC build system produces HTML that minimizes unnecessary differences
between builds:

- **Consistent attribute ordering**: attributes are always written in the
  same order (`id` before `class` before `role` before `aria-*`)
- **Consistent indentation**: output HTML uses consistent 2-space indentation
- **Stable sort orders**: navigation items, page lists, and other ordered
  content use stable, deterministic sort criteria
- **No randomized identifiers**: class names and IDs are derived from content,
  not generated randomly
- **No hash-based filenames**: CSS and JS files keep their original names
  (no `style.a1b2c3.css` cache-busting --- the Git hash serves this purpose)

This means that a `diff` between two builds shows only meaningful changes ---
content that was actually modified --- not noise from non-deterministic
formatting.

```bash
# Comparing two builds shows only real changes
diff -r _site_v1/ _site_v2/

# Expected: only files whose source markdown changed
# Not expected: every file changed because of timestamps or formatting
```

### Verifying Reproducibility

To verify that the build is reproducible:

```bash
# Build once
python3 build.py
cp -r _site/ _site_build1/

# Build again (same source, different time)
python3 build.py
cp -r _site/ _site_build2/

# Compare
diff -r _site_build1/ _site_build2/

# Expected output: (empty --- no differences)
```

If the diff is empty, the build is reproducible. If any files differ, the
build system has a bug that must be fixed before deployment.

---

## Disaster Recovery Builds

### Rebuilding the Entire Site from Scratch

In a disaster scenario --- hardware failure, corrupted output, compromised
build environment --- the HIC can be rebuilt from its minimal source
components. This section documents the absolute minimum required.

### Minimum Requirements

To rebuild the HIC from nothing, you need exactly three things:

```
DISASTER RECOVERY REQUIREMENTS:

1. Python 3.8 or later
   - Standard installation, no additional packages
   - Available on every major operating system
   - Often pre-installed on Linux and macOS

2. The content repository
   - The Git repository containing markdown source files
   - Specifically: content/*.md, build.py, templates/, static/
   - This can be restored from any backup of the repository

3. A terminal
   - Any command-line environment capable of running Python
   - bash, zsh, PowerShell, cmd.exe --- any will work
```

That is the complete list. No package manager. No `pip install`. No
`npm install`. No Docker image. No build server. No cloud service.
No internet connection.

### Zero-Dependency Verification

The build system's freedom from external dependencies can be verified:

```bash
# Verify: no imports outside the standard library
python3 -c "
import ast, sys
with open('build.py') as f:
    tree = ast.parse(f.read())
imports = [node for node in ast.walk(tree)
           if isinstance(node, (ast.Import, ast.ImportFrom))]
for imp in imports:
    if isinstance(imp, ast.Import):
        for alias in imp.names:
            print(f'import {alias.name}')
    else:
        print(f'from {imp.module} import ...')
"

# Expected output: only standard library modules
# import os
# import re
# import json
# import hashlib
# from pathlib import Path
# from datetime import datetime
# ...
#
# NOT expected: import markdown, import jinja2, import requests
```

### The Emergency Build

In the most extreme scenario --- a clean machine with only Python and the
source files --- the emergency build procedure is:

```bash
# Step 1: Verify Python is available
python3 --version
# Python 3.x.x (anything 3.8+)

# Step 2: Verify source files are present
ls content/
# index.md  getting-started/  architecture/  reference/  operations/

ls build.py
# build.py

# Step 3: Run the build
python3 build.py

# Step 4: Verify the output
ls _site/
# index.html  style.css  getting-started/  architecture/  ...

# Step 5: Serve locally (Python's built-in HTTP server)
cd _site && python3 -m http.server 8000
# Serving on http://localhost:8000

# Step 6: Open in browser and verify
# Navigate to http://localhost:8000
# Verify: all pages load, navigation works, content is correct
```

Total time from "I have Python and the repo" to "site is running": under
60 seconds.

### Bare Minimum Output

If even the template files are lost, `build.py` can produce a **bare minimum
build** that wraps each markdown file's content in minimal valid HTML:

```bash
python3 build.py --emergency

# This mode:
# - Uses a hardcoded minimal HTML template
# - Skips navigation generation (no template for it)
# - Skips CSS (relies on browser defaults)
# - Produces readable, navigable HTML files
# - Is ugly but functional
```

The emergency template:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{title}</title>
</head>
<body>
  <nav><a href="../index.html">Home</a></nav>
  <main>
    {content}
  </main>
</body>
</html>
```

No CSS. No SVG. No neon. Just the content, in valid HTML, readable in any
browser. The building without paint, without lights --- but still standing.
The concrete is intact.

### Recovery Verification

After a disaster recovery build, the following checks verify that the site
has been correctly reconstructed:

```
RECOVERY VERIFICATION CHECKLIST:

[ ] 1. PAGE COUNT
       Count HTML files in _site/
       Compare against known page count from last build manifest
       All pages present? [ YES / NO ]

[ ] 2. CONTENT INTEGRITY
       For each page, compare source markdown word count
       against rendered HTML text content word count
       All content rendered? [ YES / NO ]

[ ] 3. NAVIGATION
       Click every link in the navigation tree
       All links resolve to existing pages? [ YES / NO ]

[ ] 4. CROSS-REFERENCES
       Check all internal links within content
       All cross-references valid? [ YES / NO ]

[ ] 5. ASSETS
       Verify CSS file is present and loads
       Verify SVG images render correctly
       All assets functional? [ YES / NO ]

[ ] 6. AIR-GAP COMPLIANCE
       Run: python3 build.py --verify
       All compliance checks pass? [ YES / NO ]

[ ] 7. BUILD MANIFEST
       Compare new build manifest against last known good manifest
       Differences are explained by known changes? [ YES / NO ]

[ ] 8. VISUAL SPOT CHECK
       Load the home page, two section indexes, and three content pages
       Visual appearance matches expectations? [ YES / NO ]
```

### Backup Strategy

The HIC's disaster recovery is only as good as its backups. The backup
strategy is simple because the system is simple:

```
WHAT TO BACK UP:

1. The Git repository (contains everything needed to rebuild)
   - Backup frequency: after every commit
   - Backup method: git bundle or tar of .git directory
   - Backup destination: separate air-gapped storage

2. The build manifest from the last known-good deployment
   - Provides a reference point for recovery verification

3. The import register
   - Audit trail of all content that entered the system

WHAT NOT TO BACK UP:

- The _site/ output directory (can be regenerated from source)
- The import-staging/ directory (transient working area)
- Any caches or temporary files
```

---

## Performance Budget

### Why Performance Matters in an Air-Gapped System

"It is air-gapped, so network speed does not matter." This is a misconception.

Air-gapped systems are often served over **local networks** that may be:
- Low-bandwidth internal networks in secure facilities
- Shared infrastructure with other critical systems
- Accessed from low-powered terminals or thin clients
- Subject to strict resource monitoring and quotas

Performance also matters for:
- **Build time**: faster builds mean faster iteration cycles
- **Storage**: smaller output means more headroom on constrained systems
- **Responsiveness**: instant page loads improve documentation usability
- **Reliability**: simpler pages have fewer failure modes
- **Battery life**: on portable devices in field environments, efficiency matters

### Maximum Page Weight Targets

Every page in the HIC must meet the following weight budget:

```
PERFORMANCE BUDGET PER PAGE:

┌─────────────────────────────────────────────────┐
│ Resource          │ Budget    │ Typical  │ Max   │
├───────────────────┼───────────┼──────────┼───────┤
│ HTML document     │ < 100 KB  │ 25 KB    │ 100KB │
│ CSS (shared)      │ < 30 KB   │ 18 KB    │ 30 KB │
│ JavaScript (if    │ < 10 KB   │ 0 KB     │ 10 KB │
│   any, optional)  │           │          │       │
│ SVG assets        │ < 20 KB   │ 5 KB     │ 20 KB │
│ Total per page    │ < 150 KB  │ 48 KB    │ 150KB │
│   (first load)    │           │          │       │
│ Total per page    │ < 100 KB  │ 25 KB    │ 100KB │
│   (cached CSS)    │           │          │       │
└───────────────────┴───────────┴──────────┴───────┘
```

### HTML: Under 100KB Per Page

The HTML document for any single page must not exceed 100KB. This budget
includes all inlined content: the full navigation, the rendered content,
inline SVG icons, and any critical inline CSS.

Strategies to stay within budget:
- **Concise markup**: semantic elements without unnecessary wrappers
- **No inline JavaScript**: zero bytes of JS in the HTML
- **Efficient navigation**: hierarchical collapse for deep nav trees
- **Content discipline**: very long pages should be split into multiple pages

For reference, 100KB of HTML is approximately 2,500 lines of rendered markup,
or roughly 15,000 words of content. Very few documentation pages approach
this limit.

### CSS: Single File, Under 30KB

The entire visual system of the HIC fits in a single CSS file under 30KB.
This is achievable because:

- **No framework overhead**: no Bootstrap, no Tailwind, no utility CSS
  framework adding thousands of unused rules
- **No component library**: styles are purpose-built for the HIC's specific
  components
- **System fonts**: zero bytes of font-face declarations
- **CSS custom properties**: reusable values avoid repetitive declarations
- **Minimal specificity**: simple selectors, no deep nesting, no `!important`
  overrides (except in print styles)

```
CSS SIZE BREAKDOWN (approximate):

Base reset & defaults:       2 KB
Layout (grid, navigation):   4 KB
Typography:                  3 KB
Components:                  6 KB
  - Code blocks:  1.5 KB
  - Tables:       1.0 KB
  - Admonitions:  1.0 KB
  - Cards:        1.0 KB
  - Other:        1.5 KB
Theme (colors, neon):        4 KB
Print styles:                2 KB
Utilities:                   1 KB
                           ──────
Total:                     ~22 KB
Gzip compressed:           ~6 KB
```

### Total Site Size: Linear Scaling

The total size of the HIC output scales linearly with content volume.
There is no exponential growth, no hidden bloat:

```
SITE SIZE SCALING:

Pages    │ Total Size  │ Per Page Average
─────────┼─────────────┼─────────────────
10       │ ~500 KB     │ ~50 KB
50       │ ~2.0 MB     │ ~40 KB
100      │ ~3.8 MB     │ ~38 KB
500      │ ~18 MB      │ ~36 KB
1000     │ ~35 MB      │ ~35 KB

Notes:
- Per-page average decreases slightly at scale because
  the shared CSS is amortized across more pages
- No database, no server application, no runtime overhead
- The entire site fits on the smallest USB drive available
- Even a 1000-page site is smaller than a single
  high-resolution photograph
```

### Load Time Targets

Page load time targets for the HIC on a local network:

| Metric | Target | Notes |
|--------|--------|-------|
| First Contentful Paint | < 200ms | HTML starts rendering immediately |
| Largest Contentful Paint | < 500ms | All content visible |
| Total Page Load | < 1s | Everything including CSS |
| Time to Interactive | < 1s | Identical to load (no JS to execute) |
| Cumulative Layout Shift | 0 | No late-loading resources to cause shift |

These targets are easily achievable because:
- No render-blocking JavaScript
- Critical CSS is inlined
- No external resource fetching
- No DNS resolution for third-party domains
- No TLS handshakes with external servers
- No waiting for API responses to populate content

### No Render-Blocking Resources

The HIC loading sequence is optimized for instant rendering:

```
BROWSER LOADING SEQUENCE:

1. Receive HTML file
   ├── Parse <head>
   │   ├── Read inline critical CSS → apply immediately
   │   ├── Find <link> to style.css → begin loading (non-blocking)
   │   └── Read SVG favicon data URI → apply immediately
   ├── Parse <body>
   │   ├── Render skip link
   │   ├── Render header with inline SVG logo
   │   ├── Render navigation (static HTML, no JS needed)
   │   ├── Render main content
   │   └── Render footer
   └── style.css arrives → apply full styles (enhances appearance)

Time to first meaningful content: ~100ms on local network
Time to fully styled page: ~300ms on local network

JavaScript blocking: ZERO
External resource blocking: ZERO
Layout shift after load: ZERO
```

The critical CSS inlined in the `<head>` ensures that the page is styled
to a readable state before the full stylesheet loads. The full stylesheet
enhances the appearance (adding neon effects, refined spacing) but the page
is already usable without it.

### Budget Enforcement

The performance budget is not advisory. It is enforced by the build system:

```bash
python3 build.py --check-budget

# Output:
# [PASS] index.html: 24.2 KB (budget: 100 KB)
# [PASS] getting-started/installation.html: 31.5 KB (budget: 100 KB)
# [PASS] style.css: 21.8 KB (budget: 30 KB)
# ...
# [WARN] reference/api.html: 89.3 KB (budget: 100 KB, 89% used)
# ...
# BUDGET CHECK: PASSED (all files within budget)
# WARNINGS: 1 file(s) above 80% of budget
```

Files that exceed their budget fail the build. Files above 80% of budget
generate warnings to prompt proactive splitting or optimization.

---

## Appendix: Quick Reference

### Build Commands

```bash
# Standard build
python3 build.py

# Build with compliance verification
python3 build.py --verify

# Build with performance budget check
python3 build.py --check-budget

# Emergency build (minimal templates)
python3 build.py --emergency

# Serve locally for review
python3 -m http.server 8000 --directory _site

# Check for external references (manual scan)
grep -rn "https\?://" _site/ --include="*.html" | grep -v ">http"

# Compare two builds for reproducibility
diff -r _site_build1/ _site_build2/
```

### File Structure

```
hic/
├── content/              # Markdown source files (the blueprints)
│   ├── index.md
│   └── .../
├── templates/            # HTML templates (the molds)
│   ├── base.html
│   ├── page.html
│   └── nav.html
├── static/               # Static assets (the furnishings)
│   ├── style.css
│   └── images/
├── build.py              # The build system (the construction crew)
├── import-staging/       # USB import airlock
│   ├── incoming/
│   ├── validated/
│   └── rejected/
└── _site/                # Build output (the completed building)
    ├── index.html
    ├── style.css
    ├── build-manifest.json
    └── .../
```

### Compliance Quick Check

```bash
# One-line compliance scan
python3 build.py --verify && echo "COMPLIANT" || echo "VIOLATION DETECTED"
```

### Emergency Recovery Steps

```
1. Obtain: Python 3.8+, content repository
2. Run:    python3 build.py
3. Serve:  python3 -m http.server 8000 --directory _site
4. Verify: Browse all sections, check navigation
5. Deploy: Copy _site/ to target server
```

### Key Principles Summary

```
PRINCIPLE                          IMPLEMENTATION
────────────────────────────────────────────────────────────
Static-first                       Pre-built HTML, no runtime assembly
Zero JavaScript (core)             All navigation and content in HTML
Zero external dependencies         No CDN, no fonts, no APIs
Deterministic builds               Same inputs → same outputs always
Air-gap compliant                  No network contact of any kind
Progressive enhancement            Neon optional, concrete mandatory
Single-file build system           build.py with stdlib only
USB import protocol                Controlled ingress with quarantine
Performance budgeted               Enforced per-file size limits
Accessible by default              Semantic HTML, ARIA landmarks
Auditable                          Every file human-readable
Recoverable                        Rebuild from markdown + Python
```

---

> *"The Holm Intelligence Complex stands at the intersection of paranoia and
> pragmatism. Every design decision asks the same question: does this work
> when everything else fails? If the network is gone, if the CDN is down, if
> JavaScript is blocked, if the power flickers --- does the building still
> stand? Does the documentation still serve its readers?*
>
> *The answer is always yes. Because the building is made of concrete.
> The neon is just decoration."*

---

**Document Status:** ACTIVE
**Classification:** INFRASTRUCTURE
**Air-Gap Compliance:** VERIFIED
**Build Dependency Count:** ZERO
