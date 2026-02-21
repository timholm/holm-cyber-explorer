# STAGE 4: HOLM INTELLIGENCE COMPLEX -- OFFLINE RENDERING & SYSTEMS INTEGRATION

## Static-First Architecture, Canvas/SVG Pipelines, Cache Systems, and Cross-Layer Standards

**Document ID:** STAGE4-HIC-OFFLINE
**Version:** 1.0.0
**Date:** 2026-02-17
**Status:** Ratified
**Classification:** Specialized Systems -- These articles provide the complete technical specification for rendering the Holm Intelligence Complex interface offline, without servers or framework dependencies, and for the naming conventions, event protocols, performance budgets, and integration standards that bind every layer of the HIC into a coherent system. They assume mastery of Stage 2 philosophy, Stage 3 operational procedures, and the HIC visual design language established in prior Stage 4 HIC documents.

---

## How to Read This Document

This document contains four agent specifications -- Agents 16 through 19 -- that belong to Stage 4 of the holm.chat Documentation Institution. Agents 16 through 18 collectively define the offline rendering system: how the Holm Intelligence Complex is built, served, cached, and updated without any server infrastructure. Agent 19 defines the integration standards that every other HIC agent must obey: naming conventions, event protocols, performance budgets, accessibility requirements, and compatibility guarantees.

These are implementation specifications. They contain file trees, JSON schemas, SVG structure definitions, IndexedDB table designs, service worker strategies, and event bus contracts. They are written for the developer who must build the HIC from source files and deliver a working cyberpunk skyscraper interface that loads from a single HTML file on an air-gapped machine.

The HIC is the visual interface layer of the holm.chat Documentation Institution. It renders the institution's document library as a neon-lit skyscraper: floors are document categories, rooms are individual documents, elevators are navigation pathways, and the building's exterior is the entry point. The entire metaphor must function offline, load in under two seconds, render at sixty frames per second, and fit within a five-megabyte asset budget. These four agents specify how.

If you are implementing the HIC for the first time, read Agent 19 (Systems Integration) first. The naming conventions and performance budgets defined there constrain every decision in Agents 16 through 18. Then read Agent 16 (Static-First Architecture) for the foundational file structure. Then Agent 17 (Rendering Pipelines) for the SVG and Canvas systems. Then Agent 18 (Cache and Offline Systems) for service workers, IndexedDB, and USB import.

If you are maintaining an existing HIC implementation, these agents serve as the authoritative reference for every structural decision. When in doubt about a file name, a CSS class, a data attribute, an event name, or a performance target, the answer is in this document.

---

---

# HIC-016 -- Static-First Architecture

**Document ID:** HIC-016
**Domain:** 17 -- Interface & Visualization
**Version:** 1.0.0
**Date:** 2026-02-17
**Status:** Ratified
**Depends On:** ETH-001, CON-001, HIC-001 through HIC-015
**Depended Upon By:** HIC-017, HIC-018, HIC-019. All agents involving HIC rendering, caching, or integration.

---

## 1. Purpose

This agent specifies the static-first architecture of the Holm Intelligence Complex. Static-first means the entire building -- every floor, every room, every navigation pathway, every visual effect -- is navigable from a single HTML file and its associated assets, loaded from the local filesystem with no server, no build step, and no JavaScript framework. The browser opens the file. The building appears. The operator navigates.

This is not progressive enhancement where a server-rendered page gets enriched by JavaScript. There is no server. This is not a single-page application framework that compiles TypeScript into bundles. There is no build toolchain. This is a single HTML file that references vanilla JavaScript, SVG assets, CSS stylesheets, and JSON data files. Every dependency is a static file on disk. Every file is human-readable. Every file can be inspected, modified, and understood by an operator with a text editor.

The air-gapped mandate (CON-001, Section 3.3) makes this architecture non-negotiable. The institution's machines have no internet connection. There is no CDN to fetch libraries from. There is no npm registry. There is no package manager. The HIC must be entirely self-contained, entirely offline, and entirely comprehensible to a single developer working alone without external resources.

## 2. File Structure

The complete HIC lives under a single root directory. Every path in the system is relative to this root. The root can be placed anywhere on the filesystem -- in a user's home directory, on a USB stick, in a dedicated partition. The HTML entry point uses only relative paths, so the entire directory can be moved without breaking anything.

### 2.1 Complete File Tree

```
/hic/
├── index.html                          # Single entry point (< 15KB)
├── manifest.json                       # Web app manifest for PWA metadata
├── sw.js                               # Service worker (offline caching)
│
├── js/
│   ├── hic-core.js                     # Core initialization, event bus, state (< 20KB)
│   ├── hic-nav.js                      # Navigation: floor transitions, zoom (< 12KB)
│   ├── hic-render-svg.js               # SVG rendering pipeline (< 15KB)
│   ├── hic-render-canvas.js            # Canvas fallback renderer (< 12KB)
│   ├── hic-cache.js                    # IndexedDB and caching logic (< 10KB)
│   ├── hic-search.js                   # Full-text search over documents (< 8KB)
│   ├── hic-usb.js                      # USB import/merge logic (< 6KB)
│   └── hic-a11y.js                     # Accessibility: keyboard nav, ARIA (< 5KB)
│
├── css/
│   ├── hic-base.css                    # Reset, layout grid, typography (< 6KB)
│   ├── hic-building.css                # Building exterior styles (< 4KB)
│   ├── hic-floor.css                   # Floor plan layout styles (< 5KB)
│   ├── hic-room.css                    # Room interior styles (< 4KB)
│   └── hic-themes/
│       ├── neon.css                    # Default: cyan/magenta neon on dark (< 4KB)
│       ├── amber.css                   # Amber monochrome terminal (< 3KB)
│       ├── high-contrast.css           # WCAG AAA high contrast (< 3KB)
│       └── print.css                   # Print stylesheet, no glow (< 2KB)
│
├── building.svg                        # Full building elevation (< 80KB)
│
├── floors/
│   ├── F01.floor.json                  # Floor 1 metadata + room definitions
│   ├── F01.floor.svg                   # Floor 1 plan vector graphic
│   ├── F02.floor.json
│   ├── F02.floor.svg
│   ├── ...
│   ├── F17.floor.json
│   ├── F17.floor.svg
│   ├── SB1.floor.json                  # Sub-basement 1
│   ├── SB1.floor.svg
│   ├── SB2.floor.json                  # Sub-basement 2
│   ├── SB2.floor.svg
│   ├── G.floor.json                    # Ground floor (lobby)
│   ├── G.floor.svg
│   ├── R.floor.json                    # Roof (observatory / meta)
│   ├── R.floor.svg
│   ├── A.floor.json                    # Annex (appendices / overflow)
│   └── A.floor.svg
│
├── docs/
│   ├── F01/
│   │   ├── F01-R01.doc.json            # Document content for Floor 1 Room 1
│   │   ├── F01-R02.doc.json
│   │   └── ...
│   ├── F02/
│   │   └── ...
│   └── ...
│
├── icons/
│   ├── hic-sprites.png                 # Combined icon sprite sheet (< 30KB)
│   ├── hic-sprites@2x.png             # Retina sprite sheet (< 60KB)
│   └── hic-sprites.json               # Sprite coordinate map
│
├── fonts/
│   ├── hic-mono.woff2                  # Subsetted monospace font (< 40KB)
│   └── hic-mono-bold.woff2            # Subsetted monospace bold (< 40KB)
│
└── meta/
    ├── version.json                    # Current HIC version + floor checksums
    ├── search-index.json               # Pre-built search index (< 200KB)
    └── changelog.json                  # Version history
```

### 2.2 Asset Budget Breakdown

The total asset budget is five megabytes. This is the hard ceiling for the complete HIC with all floors, all documents, all fonts, all icons, and all code. The budget is allocated as follows:

| Category | Budget | Notes |
|---|---|---|
| HTML entry point | 15 KB | Single file, inline critical CSS |
| JavaScript (all modules) | 88 KB | Vanilla JS, no frameworks, no minifier required |
| CSS (all themes) | 31 KB | Includes all four themes |
| Building SVG | 80 KB | Full elevation, compressed |
| Floor SVGs (22 floors) | 660 KB | Average 30 KB per floor |
| Floor JSON (22 floors) | 110 KB | Average 5 KB per floor |
| Document JSON (all rooms) | 2,500 KB | The actual content; largest allocation |
| Icon sprites | 90 KB | 1x and 2x combined |
| Sprite coordinate map | 3 KB | JSON |
| Fonts | 80 KB | Two weights, subsetted |
| Service worker | 8 KB | Cache management |
| Manifest | 1 KB | PWA metadata |
| Meta files | 210 KB | Version, search index, changelog |
| **Total** | **3,876 KB** | **1,124 KB headroom** |

The 1,124 KB of headroom is reserved for growth. As new documents are added, they consume this headroom. When total assets approach 4,500 KB, the operator must audit document sizes and consider splitting the HIC into volumes (HIC-VOL-1, HIC-VOL-2), each self-contained.

### 2.3 The HTML Entry Point

The `index.html` file is the single entry point for the entire HIC. It contains:

1. The HTML5 doctype and document skeleton.
2. Inline critical CSS for above-the-fold rendering (the building exterior).
3. A `<noscript>` fallback message explaining that JavaScript is required.
4. A loading indicator styled with inline CSS (a pulsing neon rectangle).
5. The main `<div id="hic-viewport">` that receives all rendered content.
6. Script tags loading the JavaScript modules in dependency order with `defer`.
7. Service worker registration.

```html
<!DOCTYPE html>
<html lang="en" data-theme="neon">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Holm Intelligence Complex</title>
  <link rel="manifest" href="manifest.json">
  <style>
    /* Inline critical CSS: background, loading indicator, viewport shell */
    :root { --hic-bg: #0a0a0f; --hic-cyan: #00f0ff; --hic-magenta: #ff00aa; }
    body { margin: 0; background: var(--hic-bg); color: #e0e0e0;
           font-family: 'HIC Mono', 'Courier New', monospace; overflow: hidden; }
    #hic-viewport { width: 100vw; height: 100vh; position: relative; }
    #hic-loader { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%);
                  border: 2px solid var(--hic-cyan); padding: 2rem 3rem;
                  animation: hic-pulse 1.5s ease-in-out infinite; }
    @keyframes hic-pulse {
      0%, 100% { box-shadow: 0 0 10px var(--hic-cyan); opacity: 0.7; }
      50% { box-shadow: 0 0 30px var(--hic-cyan); opacity: 1; }
    }
  </style>
  <link rel="stylesheet" href="css/hic-base.css">
  <link rel="stylesheet" href="css/hic-building.css">
  <link rel="stylesheet" href="css/hic-floor.css">
  <link rel="stylesheet" href="css/hic-room.css">
  <link rel="stylesheet" href="css/hic-themes/neon.css" id="hic-theme-link">
</head>
<body>
  <div id="hic-viewport" role="application" aria-label="Holm Intelligence Complex">
    <div id="hic-loader" aria-live="polite">
      <span>INITIALIZING HIC...</span>
    </div>
  </div>
  <noscript>
    <div style="color: #ff4444; padding: 2rem; font-family: monospace;">
      HIC REQUIRES JAVASCRIPT. Enable scripting in your browser to access this system.
    </div>
  </noscript>
  <script src="js/hic-core.js" defer></script>
  <script src="js/hic-a11y.js" defer></script>
  <script src="js/hic-cache.js" defer></script>
  <script src="js/hic-render-svg.js" defer></script>
  <script src="js/hic-render-canvas.js" defer></script>
  <script src="js/hic-nav.js" defer></script>
  <script src="js/hic-search.js" defer></script>
  <script src="js/hic-usb.js" defer></script>
  <script>
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.register('sw.js').catch(function(err) {
        console.warn('HIC: Service worker registration failed:', err);
      });
    }
  </script>
</body>
</html>
```

### 2.4 No Frameworks, No Build Step

The HIC uses vanilla JavaScript exclusively. No React. No Vue. No Svelte. No Angular. No jQuery. No Webpack. No Vite. No Rollup. No TypeScript. No JSX. No templating engines. No transpilation. No polyfills beyond what the operator writes by hand.

The reasons are structural, not ideological:

1. **Air-gap compatibility.** Frameworks require package managers. Package managers require internet access. The institution has no internet access.
2. **Long-term survival.** Frameworks have release cycles measured in months. The institution operates on a timescale of decades. A React application written in 2026 will not build in 2036 without significant migration effort. Vanilla JavaScript written to web standards will run in 2036 without modification.
3. **Comprehensibility.** A single developer must be able to read, understand, and modify every line of the HIC. Framework abstractions introduce cognitive overhead proportional to the framework's API surface. Vanilla JavaScript has no API surface beyond the browser's DOM and standard library.
4. **File size.** The entire HIC JavaScript budget is 88 KB. React's runtime alone exceeds this. The HIC cannot afford framework overhead.

The JavaScript modules use ES module syntax (`import`/`export`) where supported, with a fallback to script-tag loading order for older browsers. Each module is a single file. Each file is under 20 KB. Each file has a single responsibility.

---

---

# HIC-017 -- SVG and Canvas Rendering Pipelines

**Document ID:** HIC-017
**Domain:** 17 -- Interface & Visualization
**Version:** 1.0.0
**Date:** 2026-02-17
**Status:** Ratified
**Depends On:** HIC-016
**Depended Upon By:** HIC-018, HIC-019

---

## 1. Purpose

This agent specifies the two rendering pipelines that draw the Holm Intelligence Complex in the browser: the SVG pipeline (primary) and the Canvas pipeline (fallback). SVG is the preferred renderer for floor plans and room views because it produces crisp vector graphics at any zoom level, supports CSS styling for the neon theme, and allows individual elements to be interactive DOM nodes. Canvas is the fallback renderer for the building elevation view where dozens of floors must be drawn simultaneously and SVG performance degrades.

## 2. SVG Pipeline

### 2.1 SVG Structure for Floor Plans

Every floor plan is an SVG file that conforms to a strict structure. The structure is not negotiable. Rendering code, accessibility code, and navigation code all depend on this structure being exactly as specified.

```xml
<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg"
     viewBox="0 0 1200 800"
     class="hic-floor"
     data-floor="F07"
     data-version="HIC-1.2.0"
     role="img"
     aria-label="Floor 7: Security and Integrity">

  <!-- Background layer: floor slab, structural grid -->
  <g class="hic-layer-bg" aria-hidden="true">
    <rect class="hic-floor-slab" x="0" y="0" width="1200" height="800"
          fill="#0d0d14" stroke="#1a1a2e" stroke-width="1"/>
    <!-- Structural grid lines (decorative) -->
    <line class="hic-grid" x1="200" y1="0" x2="200" y2="800" stroke="#111122" stroke-width="0.5"/>
    <line class="hic-grid" x1="400" y1="0" x2="400" y2="800" stroke="#111122" stroke-width="0.5"/>
    <!-- ... additional grid lines ... -->
  </g>

  <!-- Walls layer: corridor walls, room boundaries -->
  <g class="hic-layer-walls">
    <!-- Main corridor -->
    <path class="hic-wall hic-wall-corridor" d="M 100,100 L 1100,100 L 1100,700 L 100,700 Z"
          fill="none" stroke="#2a2a4e" stroke-width="3"/>
    <!-- Internal walls separating rooms -->
    <line class="hic-wall hic-wall-internal" x1="400" y1="100" x2="400" y2="700"
          stroke="#2a2a4e" stroke-width="2"/>
    <!-- ... additional walls ... -->
  </g>

  <!-- Rooms layer: interactive room areas -->
  <g class="hic-layer-rooms">
    <g class="hic-room" data-room="F07-R01" data-status="active"
       data-security="standard" tabindex="0" role="button"
       aria-label="Room 1: Air-Gap Architecture">
      <rect class="hic-room-area" x="105" y="105" width="290" height="290"
            fill="#0f0f1a" stroke="var(--hic-cyan)" stroke-width="1.5"/>
      <text class="hic-label hic-label-room-id" x="115" y="130"
            fill="var(--hic-cyan)" font-size="11">F07-R01</text>
      <text class="hic-label hic-label-room-title" x="115" y="150"
            fill="#8888aa" font-size="10">Air-Gap Architecture</text>
      <!-- Glow effect (CSS-animated) -->
      <rect class="hic-glow" x="103" y="103" width="294" height="294"
            fill="none" stroke="var(--hic-cyan)" stroke-width="1"
            opacity="0.3" filter="url(#hic-glow-filter)"/>
    </g>

    <g class="hic-room" data-room="F07-R02" data-status="active"
       data-security="elevated" tabindex="0" role="button"
       aria-label="Room 2: Cryptographic Survival">
      <rect class="hic-room-area" x="405" y="105" width="290" height="290"
            fill="#0f0f1a" stroke="var(--hic-magenta)" stroke-width="1.5"/>
      <text class="hic-label hic-label-room-id" x="415" y="130"
            fill="var(--hic-magenta)" font-size="11">F07-R02</text>
      <text class="hic-label hic-label-room-title" x="415" y="150"
            fill="#8888aa" font-size="10">Cryptographic Survival</text>
      <rect class="hic-glow" x="403" y="103" width="294" height="294"
            fill="none" stroke="var(--hic-magenta)" stroke-width="1"
            opacity="0.3" filter="url(#hic-glow-filter)"/>
    </g>
    <!-- ... additional rooms ... -->
  </g>

  <!-- Doors layer: entry points between corridor and rooms -->
  <g class="hic-layer-doors">
    <rect class="hic-door" data-connects="corridor:F07-R01"
          x="200" y="395" width="40" height="10"
          fill="var(--hic-cyan)" opacity="0.6"/>
    <rect class="hic-door" data-connects="corridor:F07-R02"
          x="500" y="395" width="40" height="10"
          fill="var(--hic-magenta)" opacity="0.6"/>
  </g>

  <!-- Labels layer: floor title, navigation markers -->
  <g class="hic-layer-labels">
    <text class="hic-label hic-label-floor-title" x="600" y="40"
          text-anchor="middle" fill="var(--hic-cyan)" font-size="16"
          font-weight="bold">FLOOR 07 -- SECURITY & INTEGRITY</text>
    <text class="hic-label hic-label-floor-id" x="600" y="60"
          text-anchor="middle" fill="#555577" font-size="10">F07</text>
  </g>

  <!-- Effects layer: glow filters, ambient animations -->
  <defs>
    <filter id="hic-glow-filter" x="-20%" y="-20%" width="140%" height="140%">
      <feGaussianBlur in="SourceGraphic" stdDeviation="4" result="blur"/>
      <feMerge>
        <feMergeNode in="blur"/>
        <feMergeNode in="SourceGraphic"/>
      </feMerge>
    </filter>
    <filter id="hic-glow-strong" x="-30%" y="-30%" width="160%" height="160%">
      <feGaussianBlur in="SourceGraphic" stdDeviation="8" result="blur"/>
      <feMerge>
        <feMergeNode in="blur"/>
        <feMergeNode in="blur"/>
        <feMergeNode in="SourceGraphic"/>
      </feMerge>
    </filter>
  </defs>
</svg>
```

### 2.2 SVG Layer Ordering

The SVG layers are rendered bottom to top in this fixed order:

1. **Background layer** (`hic-layer-bg`): Floor slab, structural grid. Decorative only. Marked `aria-hidden="true"`.
2. **Walls layer** (`hic-layer-walls`): Corridor walls, internal walls, structural boundaries.
3. **Rooms layer** (`hic-layer-rooms`): Interactive room groups. Each room is a `<g>` with `tabindex="0"` and `role="button"`.
4. **Doors layer** (`hic-layer-doors`): Door indicators connecting corridor to rooms.
5. **Labels layer** (`hic-layer-labels`): Floor title, floor ID, navigation markers.
6. **Effects layer** (`<defs>`): Glow filters, gradients, animation definitions.

No additional layers may be added without amending this specification. Rendering code iterates layers by class name. Adding unlisted layers will cause undefined rendering behavior.

### 2.3 SVG Rendering Pipeline Steps

When the operator navigates to a floor, the SVG rendering pipeline executes the following steps in order:

1. **Fetch floor JSON.** Load `floors/{floor_id}.floor.json` from cache or filesystem. Parse metadata.
2. **Fetch floor SVG.** Load `floors/{floor_id}.floor.svg` from cache or filesystem.
3. **Parse SVG.** Insert SVG into a temporary `<div>`, not yet in the DOM.
4. **Bind data attributes.** For each `.hic-room` element, look up the corresponding room in the floor JSON. Set `data-status` (active, locked, empty, draft). Set `data-security` (standard, elevated, restricted).
5. **Apply theme.** The current theme CSS handles color mapping via CSS custom properties. No JavaScript color manipulation.
6. **Bind event listeners.** Attach click, keydown (Enter/Space), and focus handlers to each `.hic-room` element.
7. **Insert into DOM.** Replace the current content of `#hic-viewport` with the new SVG. Remove the old SVG from memory.
8. **Dispatch event.** Fire `hic:floor-enter` on `document` with the floor ID as detail.
9. **Preload adjacent floors.** Trigger background fetch for floors above and below (see Agent 18).

Total time budget for steps 1-8: under 500 milliseconds. Step 9 runs asynchronously.

### 2.4 Building Elevation SVG Structure

The building elevation (`building.svg`) shows the full skyscraper from outside. Each floor is a horizontal band. The operator clicks a floor band to navigate to that floor.

```xml
<svg xmlns="http://www.w3.org/2000/svg"
     viewBox="0 0 400 1200"
     class="hic-building"
     role="navigation"
     aria-label="Holm Intelligence Complex Building Elevation">

  <!-- Sky/atmosphere background -->
  <rect class="hic-sky" x="0" y="0" width="400" height="1200"
        fill="url(#hic-sky-gradient)"/>

  <!-- Building shell -->
  <rect class="hic-building-shell" x="50" y="80" width="300" height="1040"
        fill="#0a0a12" stroke="#1a1a3e" stroke-width="2"/>

  <!-- Floor bands (top to bottom: Roof, F17, F16, ... F01, G, SB1, SB2) -->
  <g class="hic-building-floors">
    <g class="hic-building-floor" data-floor="R" tabindex="0" role="link"
       aria-label="Roof: Observatory">
      <rect x="52" y="82" width="296" height="40" fill="#0d0d16"
            stroke="var(--hic-cyan)" stroke-width="0.5"/>
      <text x="200" y="106" text-anchor="middle" fill="var(--hic-cyan)"
            font-size="9">R -- OBSERVATORY</text>
    </g>
    <g class="hic-building-floor" data-floor="F17" tabindex="0" role="link"
       aria-label="Floor 17: Evolution and Memory">
      <rect x="52" y="124" width="296" height="40" fill="#0d0d16"
            stroke="var(--hic-cyan)" stroke-width="0.5"/>
      <text x="200" y="148" text-anchor="middle" fill="#6688aa"
            font-size="9">F17 -- EVOLUTION & MEMORY</text>
    </g>
    <!-- ... F16 through F01 ... -->
    <g class="hic-building-floor" data-floor="G" tabindex="0" role="link"
       aria-label="Ground Floor: Lobby">
      <rect x="52" y="804" width="296" height="48" fill="#0d0d16"
            stroke="var(--hic-cyan)" stroke-width="1"/>
      <text x="200" y="832" text-anchor="middle" fill="var(--hic-cyan)"
            font-size="10" font-weight="bold">G -- LOBBY</text>
    </g>
    <g class="hic-building-floor" data-floor="SB1" tabindex="0" role="link"
       aria-label="Sub-Basement 1: Deep Archives">
      <rect x="52" y="854" width="296" height="40" fill="#080810"
            stroke="#332244" stroke-width="0.5"/>
      <text x="200" y="878" text-anchor="middle" fill="#6644aa"
            font-size="9">SB1 -- DEEP ARCHIVES</text>
    </g>
    <g class="hic-building-floor" data-floor="SB2" tabindex="0" role="link"
       aria-label="Sub-Basement 2: Foundation">
      <rect x="52" y="896" width="296" height="40" fill="#060610"
            stroke="#332244" stroke-width="0.5"/>
      <text x="200" y="920" text-anchor="middle" fill="#6644aa"
            font-size="9">SB2 -- FOUNDATION</text>
    </g>
  </g>

  <!-- Ambient neon effects -->
  <defs>
    <linearGradient id="hic-sky-gradient" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="#020208"/>
      <stop offset="100%" stop-color="#0a0a1a"/>
    </linearGradient>
  </defs>
</svg>
```

## 3. Canvas Fallback Pipeline

### 3.1 When Canvas Is Used

Canvas replaces SVG under two conditions:

1. **Building elevation view with more than 22 floors.** If the HIC grows beyond the standard floor count, the building elevation SVG becomes too dense for smooth rendering. Canvas draws all floor bands as pixel rectangles with no DOM overhead.
2. **Explicit operator preference.** The operator can set `data-renderer="canvas"` on `#hic-viewport` to force Canvas for all views. This is intended for low-powered devices where SVG DOM manipulation is slow.

Canvas is never the default. SVG is always preferred because it supports CSS theming, accessibility attributes, and resolution-independent rendering.

### 3.2 Canvas Rendering Architecture

The Canvas renderer uses a retained-mode abstraction over the immediate-mode Canvas API. Every drawable element is represented as a plain JavaScript object with position, size, color, and type properties. The renderer maintains an array of these objects and redraws only when state changes.

```javascript
// Canvas element structure (plain object, no classes)
const canvasElement = {
  type: 'floor-band',        // Element type for hit-testing
  id: 'F07',                 // Floor ID for event dispatch
  x: 52, y: 124,             // Position
  w: 296, h: 40,             // Dimensions
  fill: '#0d0d16',           // Fill color
  stroke: '#00f0ff',         // Stroke color
  strokeWidth: 0.5,          // Stroke width
  label: 'F07 -- SECURITY',  // Text label
  labelColor: '#6688aa',     // Label color
  hover: false,              // Current hover state
  focus: false               // Current keyboard focus state
};
```

The Canvas render loop:

1. Clear the canvas.
2. Draw sky gradient (single `fillRect` with `createLinearGradient`).
3. Draw building shell (single `strokeRect`).
4. Iterate floor elements. For each: `fillRect` for background, `strokeRect` for border, `fillText` for label.
5. If a floor has `hover: true`, redraw it with brighter stroke and glow effect (`shadowBlur`, `shadowColor`).
6. If a floor has `focus: true`, draw a focus ring (dashed `strokeRect`).

The render loop uses `requestAnimationFrame` but only schedules a frame when state has changed (dirty flag pattern). No animation runs when the view is static.

### 3.3 Canvas Hit-Testing

Because Canvas elements are not DOM nodes, hit-testing is manual. On `mousemove` and `click`, the renderer:

1. Gets the mouse coordinates relative to the canvas.
2. Iterates the element array in reverse order (top-most first).
3. Tests if the point is within the element's bounding box.
4. Sets hover/focus state on the hit element, clears it on all others.
5. On click, dispatches the same `hic:floor-enter` event as the SVG pipeline.

Keyboard navigation in Canvas mode uses a virtual focus index. Arrow keys increment/decrement the index. Enter activates the focused floor.

---

---

# HIC-018 -- Cache, Offline Systems, and USB Import

**Document ID:** HIC-018
**Domain:** 17 -- Interface & Visualization
**Version:** 1.0.0
**Date:** 2026-02-17
**Status:** Ratified
**Depends On:** HIC-016, HIC-017
**Depended Upon By:** HIC-019

---

## 1. Purpose

This agent specifies every mechanism that allows the HIC to function offline: the service worker caching strategy, the IndexedDB schema for document storage, the cache invalidation protocol, the floor preloading strategy, and the USB import system for loading new documents into an air-gapped HIC installation.

## 2. Service Worker Strategy

### 2.1 Cache Architecture

The service worker manages three named caches:

| Cache Name | Contents | Strategy |
|---|---|---|
| `hic-shell-v{N}` | `index.html`, all JS, all CSS, fonts, sprites, `building.svg` | Cache-first, update in background |
| `hic-floors-v{N}` | All `*.floor.json` and `*.floor.svg` files | Cache-first, version-stamped |
| `hic-docs-v{N}` | All `*.doc.json` files | Cache-first, lazy-loaded |

The version number `{N}` in each cache name is incremented when the corresponding assets change. The current version numbers are stored in `meta/version.json`.

### 2.2 Service Worker Registration and Lifecycle

```javascript
// sw.js -- Service Worker
const SHELL_CACHE = 'hic-shell-v1';
const FLOOR_CACHE = 'hic-floors-v1';
const DOCS_CACHE = 'hic-docs-v1';

const SHELL_ASSETS = [
  './',
  './index.html',
  './js/hic-core.js',
  './js/hic-nav.js',
  './js/hic-render-svg.js',
  './js/hic-render-canvas.js',
  './js/hic-cache.js',
  './js/hic-search.js',
  './js/hic-usb.js',
  './js/hic-a11y.js',
  './css/hic-base.css',
  './css/hic-building.css',
  './css/hic-floor.css',
  './css/hic-room.css',
  './css/hic-themes/neon.css',
  './css/hic-themes/amber.css',
  './css/hic-themes/high-contrast.css',
  './css/hic-themes/print.css',
  './building.svg',
  './icons/hic-sprites.png',
  './icons/hic-sprites@2x.png',
  './icons/hic-sprites.json',
  './fonts/hic-mono.woff2',
  './fonts/hic-mono-bold.woff2',
  './meta/version.json',
  './meta/search-index.json'
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(SHELL_CACHE).then((cache) => cache.addAll(SHELL_ASSETS))
  );
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) => {
      return Promise.all(
        keys.filter((key) => {
          return key !== SHELL_CACHE && key !== FLOOR_CACHE && key !== DOCS_CACHE;
        }).map((key) => caches.delete(key))
      );
    })
  );
  self.clients.claim();
});

self.addEventListener('fetch', (event) => {
  const url = new URL(event.request.url);
  const path = url.pathname;

  // Shell assets: cache-first
  if (SHELL_ASSETS.some((asset) => path.endsWith(asset.replace('./', '')))) {
    event.respondWith(
      caches.match(event.request).then((cached) => cached || fetch(event.request))
    );
    return;
  }

  // Floor assets: cache-first with version check
  if (path.includes('/floors/')) {
    event.respondWith(
      caches.open(FLOOR_CACHE).then((cache) => {
        return cache.match(event.request).then((cached) => {
          if (cached) return cached;
          return fetch(event.request).then((response) => {
            cache.put(event.request, response.clone());
            return response;
          });
        });
      })
    );
    return;
  }

  // Document assets: cache-first, lazy
  if (path.includes('/docs/')) {
    event.respondWith(
      caches.open(DOCS_CACHE).then((cache) => {
        return cache.match(event.request).then((cached) => {
          if (cached) return cached;
          return fetch(event.request).then((response) => {
            cache.put(event.request, response.clone());
            return response;
          });
        });
      })
    );
    return;
  }

  // Default: network with cache fallback
  event.respondWith(
    fetch(event.request).catch(() => caches.match(event.request))
  );
});
```

## 3. IndexedDB Schema

### 3.1 Database Definition

The HIC uses a single IndexedDB database named `hic-store` with the following object stores:

```javascript
// Database: hic-store, version 1
const DB_NAME = 'hic-store';
const DB_VERSION = 1;

function openDatabase() {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION);

    request.onupgradeneeded = (event) => {
      const db = event.target.result;

      // Floor metadata store
      if (!db.objectStoreNames.contains('floors')) {
        const floorStore = db.createObjectStore('floors', { keyPath: 'floorId' });
        floorStore.createIndex('byStatus', 'status', { unique: false });
        floorStore.createIndex('byVersion', 'version', { unique: false });
      }

      // Document content store
      if (!db.objectStoreNames.contains('documents')) {
        const docStore = db.createObjectStore('documents', { keyPath: 'docId' });
        docStore.createIndex('byFloor', 'floorId', { unique: false });
        docStore.createIndex('byModified', 'modified', { unique: false });
        docStore.createIndex('byStatus', 'status', { unique: false });
      }

      // User state store (bookmarks, last position, preferences)
      if (!db.objectStoreNames.contains('userState')) {
        db.createObjectStore('userState', { keyPath: 'key' });
      }

      // Import log store (USB import history)
      if (!db.objectStoreNames.contains('importLog')) {
        const importStore = db.createObjectStore('importLog',
          { keyPath: 'importId', autoIncrement: true });
        importStore.createIndex('byTimestamp', 'timestamp', { unique: false });
      }
    };

    request.onsuccess = (event) => resolve(event.target.result);
    request.onerror = (event) => reject(event.target.error);
  });
}
```

### 3.2 Floor Record Schema

```json
{
  "floorId": "F07",
  "title": "Security & Integrity",
  "version": "HIC-1.2.0",
  "checksum": "sha256:a3f2b8c1d9e0...",
  "status": "active",
  "roomCount": 5,
  "rooms": [
    {
      "roomId": "F07-R01",
      "title": "Air-Gap Architecture",
      "status": "active",
      "security": "standard",
      "docId": "SEC-004",
      "wordCount": 8200,
      "modified": "2026-02-17T00:00:00Z"
    },
    {
      "roomId": "F07-R02",
      "title": "Cryptographic Survival",
      "status": "active",
      "security": "elevated",
      "docId": "SEC-005",
      "wordCount": 7400,
      "modified": "2026-02-17T00:00:00Z"
    }
  ],
  "adjacentFloors": {
    "above": "F08",
    "below": "F06"
  },
  "svgFile": "F07.floor.svg",
  "modified": "2026-02-17T00:00:00Z"
}
```

### 3.3 Document Record Schema

```json
{
  "docId": "F07-R01",
  "floorId": "F07",
  "roomId": "F07-R01",
  "title": "Air-Gap Architecture: Theory and Implementation",
  "documentReference": "SEC-004",
  "version": "1.0.0",
  "status": "ratified",
  "security": "standard",
  "content": {
    "format": "markdown",
    "body": "## 1. Purpose\n\nThis article is the complete technical specification...",
    "sections": [
      { "id": "s1", "title": "Purpose", "offset": 0 },
      { "id": "s2", "title": "Scope", "offset": 1240 },
      { "id": "s3", "title": "Background", "offset": 2860 }
    ]
  },
  "metadata": {
    "author": "Institution",
    "created": "2026-02-16T00:00:00Z",
    "modified": "2026-02-17T00:00:00Z",
    "wordCount": 8200,
    "depends_on": ["ETH-001", "CON-001", "SEC-001"],
    "depended_upon_by": ["SEC-005", "SEC-006"]
  },
  "checksum": "sha256:b4c3d2e1f0a9..."
}
```

## 4. Cache Invalidation Strategy

### 4.1 Version Stamps per Floor

Every floor has an independent version stamp. When a single document on Floor 7 is updated, only Floor 7's version stamp changes. The other 21 floors remain cached. This prevents a single document edit from invalidating the entire cache.

The version stamps are stored in `meta/version.json`:

```json
{
  "hicVersion": "HIC-1.2.0",
  "buildTimestamp": "2026-02-17T14:30:00Z",
  "shell": {
    "version": 3,
    "checksum": "sha256:1a2b3c4d..."
  },
  "floors": {
    "SB2": { "version": 1, "checksum": "sha256:aabb1122..." },
    "SB1": { "version": 1, "checksum": "sha256:ccdd3344..." },
    "G":   { "version": 2, "checksum": "sha256:eeff5566..." },
    "F01": { "version": 1, "checksum": "sha256:11223344..." },
    "F02": { "version": 3, "checksum": "sha256:55667788..." },
    "F03": { "version": 1, "checksum": "sha256:99aabbcc..." },
    "F04": { "version": 1, "checksum": "sha256:ddeeff00..." },
    "F05": { "version": 2, "checksum": "sha256:11335577..." },
    "F06": { "version": 1, "checksum": "sha256:99bb1133..." },
    "F07": { "version": 4, "checksum": "sha256:a3f2b8c1..." },
    "F08": { "version": 1, "checksum": "sha256:22446688..." },
    "F09": { "version": 1, "checksum": "sha256:aaccee00..." },
    "F10": { "version": 2, "checksum": "sha256:11335599..." },
    "F11": { "version": 1, "checksum": "sha256:77bbddff..." },
    "F12": { "version": 1, "checksum": "sha256:22448866..." },
    "F13": { "version": 1, "checksum": "sha256:aaccee11..." },
    "F14": { "version": 1, "checksum": "sha256:33557799..." },
    "F15": { "version": 1, "checksum": "sha256:bbddff22..." },
    "F16": { "version": 1, "checksum": "sha256:44668800..." },
    "F17": { "version": 1, "checksum": "sha256:ccee0022..." },
    "R":   { "version": 1, "checksum": "sha256:55770099..." },
    "A":   { "version": 1, "checksum": "sha256:dd110033..." }
  }
}
```

### 4.2 Invalidation Flow

When the HIC boots, `hic-cache.js` executes this invalidation check:

1. Fetch `meta/version.json` from disk (bypassing cache).
2. Compare `shell.version` against the current `SHELL_CACHE` version. If different, purge the shell cache and re-cache all shell assets.
3. For each floor in `floors`, compare the stored version against the version in IndexedDB. If a floor's version has increased, delete that floor's `.floor.json` and `.floor.svg` from the floor cache and its document records from the docs cache.
4. Any floor whose version matches the cached version is left untouched.
5. Write the new version stamps to IndexedDB under the `userState` store with key `lastVersionCheck`.

This strategy guarantees that stale content is never served after an update, while unchanged floors pay zero invalidation cost.

## 5. Preloading Strategy

### 5.1 Adjacent Floor Preloading

When the operator views a floor, the HIC preloads the floors immediately above and below. This ensures that vertical navigation (going up or down one floor) feels instantaneous.

```javascript
function preloadAdjacentFloors(currentFloorId) {
  const floorOrder = [
    'SB2', 'SB1', 'G',
    'F01', 'F02', 'F03', 'F04', 'F05', 'F06', 'F07', 'F08', 'F09',
    'F10', 'F11', 'F12', 'F13', 'F14', 'F15', 'F16', 'F17',
    'R', 'A'
  ];
  const currentIndex = floorOrder.indexOf(currentFloorId);
  const toPreload = [];

  if (currentIndex > 0) toPreload.push(floorOrder[currentIndex - 1]);
  if (currentIndex < floorOrder.length - 1) toPreload.push(floorOrder[currentIndex + 1]);

  toPreload.forEach((floorId) => {
    // Fetch into cache if not already cached
    const jsonPath = `floors/${floorId}.floor.json`;
    const svgPath = `floors/${floorId}.floor.svg`;

    caches.open(FLOOR_CACHE).then((cache) => {
      cache.match(jsonPath).then((existing) => {
        if (!existing) fetch(jsonPath).then((r) => cache.put(jsonPath, r));
      });
      cache.match(svgPath).then((existing) => {
        if (!existing) fetch(svgPath).then((r) => cache.put(svgPath, r));
      });
    });
  });
}
```

### 5.2 Document Preloading

When the operator views a floor, document JSON files for that floor's rooms are preloaded lazily. The preloader fetches one document at a time with a 100ms delay between requests to avoid blocking the main thread.

```javascript
function preloadFloorDocuments(floorData) {
  const rooms = floorData.rooms || [];
  let index = 0;

  function loadNext() {
    if (index >= rooms.length) return;
    const room = rooms[index++];
    const docPath = `docs/${floorData.floorId}/${room.roomId}.doc.json`;

    caches.open(DOCS_CACHE).then((cache) => {
      cache.match(docPath).then((existing) => {
        if (!existing) {
          fetch(docPath).then((r) => {
            cache.put(docPath, r);
            setTimeout(loadNext, 100);
          }).catch(() => setTimeout(loadNext, 100));
        } else {
          setTimeout(loadNext, 50);
        }
      });
    });
  }

  // Start after a brief delay to prioritize floor render
  setTimeout(loadNext, 300);
}
```

## 6. USB Import System

### 6.1 Purpose

The USB import system allows the operator to update the HIC on an air-gapped machine by loading new or updated documents from a USB stick. The operator inserts the USB, opens the HIC, triggers the import function, selects the import bundle file, and the HIC merges the new content into its local storage.

### 6.2 Import Bundle Format

An import bundle is a single JSON file with the extension `.hic-import.json`. It contains everything needed to update one or more floors:

```json
{
  "bundleVersion": "1.0.0",
  "hicTargetVersion": "HIC-1.2.0",
  "created": "2026-02-17T10:00:00Z",
  "creator": "Operator A",
  "description": "Updated SEC-004, added SEC-009",
  "floors": {
    "F07": {
      "version": 5,
      "checksum": "sha256:newchecksum...",
      "floorJson": { /* complete floor JSON object */ },
      "floorSvg": "<svg>...complete SVG string...</svg>",
      "documents": {
        "F07-R01": { /* complete document JSON object */ },
        "F07-R05": { /* new document JSON object */ }
      }
    }
  },
  "searchIndexPatch": {
    "added": [
      { "docId": "F07-R05", "terms": ["network", "isolation", "protocol"] }
    ],
    "removed": []
  },
  "bundleChecksum": "sha256:bundlechecksum..."
}
```

### 6.3 Import Procedure

The import is triggered by the operator through a UI control (a button labeled "IMPORT FROM USB" in the building lobby view, Floor G). The procedure:

1. **File selection.** The HIC opens a `<input type="file" accept=".hic-import.json">` dialog. The operator selects the import bundle from the USB.
2. **Validation.** Parse the JSON. Verify `bundleChecksum` by computing SHA-256 over the bundle contents (excluding the checksum field itself). If the checksum fails, reject the import with error `HIC-ERR-IMPORT-CHECKSUM`.
3. **Version check.** Compare `hicTargetVersion` with the current HIC version. If the bundle targets a future major version (e.g., bundle targets HIC-2.x but system is HIC-1.x), reject with `HIC-ERR-IMPORT-VERSION`.
4. **Floor-by-floor merge.** For each floor in the bundle:
   a. Compare the bundle's floor version against the local floor version.
   b. If the bundle version is higher, replace the local floor JSON, floor SVG, and all included documents.
   c. If the bundle version is equal or lower, skip that floor (no downgrade).
5. **Search index patch.** Apply the `searchIndexPatch`: add new terms, remove obsolete terms.
6. **Version stamp update.** Update `meta/version.json` with new floor versions and checksums.
7. **Cache invalidation.** Purge cached versions of updated floors from the service worker caches.
8. **Import log.** Write an entry to the `importLog` IndexedDB store recording timestamp, bundle description, floors updated, and result.
9. **Confirmation.** Display a summary: "IMPORT COMPLETE. Floors updated: F07. Documents added: 1. Documents updated: 1."

### 6.4 Conflict Resolution

Conflicts are resolved with a simple rule: **higher version wins.** There is no merge. There is no diff. If the import bundle contains a floor with a higher version number than the local copy, the bundle's version replaces the local version entirely. If the local copy has a higher version (because a different import already updated it), the bundle's version is discarded.

This strategy is intentionally simple. The institution's workflow produces documents on a single authoring machine and distributes them via USB to reader machines. There is no concurrent editing. There is one source of truth. Version numbers increase monotonically. Conflicts mean either the import is outdated (skip it) or the local copy is outdated (replace it).

## 7. Image Optimization

### 7.1 SVG Preferred for All Vector Content

All floor plans, building elevations, room layouts, diagrams, and decorative elements use SVG. SVG files are text-based, diffable, version-controllable, and resolution-independent. They compress well with gzip (typically 70-80% reduction). They can be styled with CSS. They can be manipulated with JavaScript. There is no reason to use raster formats for any structural element of the HIC.

### 7.2 PNG Sprites for Icons Only

The only raster images in the HIC are the icon sprite sheets (`hic-sprites.png` and `hic-sprites@2x.png`). These contain small icons (16x16 and 32x32) for UI controls: zoom in, zoom out, search, menu, close, elevator up, elevator down, lock, unlock. The sprite coordinate map (`hic-sprites.json`) maps icon names to pixel coordinates:

```json
{
  "spriteSheet": "hic-sprites.png",
  "spriteSheet2x": "hic-sprites@2x.png",
  "iconSize": 16,
  "iconSize2x": 32,
  "icons": {
    "zoom-in":      { "x": 0,   "y": 0 },
    "zoom-out":     { "x": 16,  "y": 0 },
    "search":       { "x": 32,  "y": 0 },
    "menu":         { "x": 48,  "y": 0 },
    "close":        { "x": 64,  "y": 0 },
    "elevator-up":  { "x": 80,  "y": 0 },
    "elevator-down":{ "x": 96,  "y": 0 },
    "lock":         { "x": 112, "y": 0 },
    "unlock":       { "x": 128, "y": 0 },
    "bookmark":     { "x": 144, "y": 0 },
    "import":       { "x": 160, "y": 0 },
    "status-active":{ "x": 176, "y": 0 },
    "status-draft": { "x": 192, "y": 0 },
    "status-locked":{ "x": 208, "y": 0 }
  }
}
```

### 7.3 No JPEG, No WebP, No GIF

The HIC contains no photographs, no gradients rendered as rasters, and no animations rendered as rasters. There is no use case for JPEG, WebP, or GIF. If a future requirement introduces photographic content, it must be added as a separate media layer outside the five-megabyte core budget.

## 8. Font Subsetting

### 8.1 Monospace Only

The HIC uses a single monospace font family in two weights: regular and bold. The font is subsetted to include only the characters actually used in the interface and documents:

- ASCII printable characters (U+0020 through U+007E)
- Common punctuation and symbols: en-dash, em-dash, bullet, ellipsis, copyright
- Basic Latin Extended for accented characters used in proper names
- Box-drawing characters (U+2500 through U+257F) for table rendering
- A small set of mathematical operators for technical documents

### 8.2 Subsetting Procedure

The font subsetting is performed offline using `pyftsubset` (part of the `fonttools` Python package):

```bash
pyftsubset SourceFont-Regular.ttf \
  --output-file=hic-mono.woff2 \
  --flavor=woff2 \
  --unicodes="U+0020-007E,U+00A0-00FF,U+2010-2027,U+2032-2037,U+2500-257F,U+2200-2211" \
  --layout-features='kern,liga'

pyftsubset SourceFont-Bold.ttf \
  --output-file=hic-mono-bold.woff2 \
  --flavor=woff2 \
  --unicodes="U+0020-007E,U+00A0-00FF,U+2010-2027,U+2032-2037,U+2500-257F,U+2200-2211" \
  --layout-features='kern,liga'
```

The resulting WOFF2 files are typically 30-40 KB each, well within the 80 KB font budget.

### 8.3 Font Loading

Fonts are loaded via `@font-face` in `hic-base.css` with `font-display: swap` to prevent invisible text during load:

```css
@font-face {
  font-family: 'HIC Mono';
  src: url('../fonts/hic-mono.woff2') format('woff2');
  font-weight: 400;
  font-style: normal;
  font-display: swap;
}

@font-face {
  font-family: 'HIC Mono';
  src: url('../fonts/hic-mono-bold.woff2') format('woff2');
  font-weight: 700;
  font-style: normal;
  font-display: swap;
}
```

The fallback stack is `'HIC Mono', 'Courier New', 'Courier', monospace`. On systems where the custom font fails to load, Courier New provides acceptable rendering.

---

---

# HIC-019 -- Systems Integration Standards

**Document ID:** HIC-019
**Domain:** 17 -- Interface & Visualization
**Version:** 1.0.0
**Date:** 2026-02-17
**Status:** Ratified
**Depends On:** HIC-016, HIC-017, HIC-018
**Depended Upon By:** All future HIC agents. All code that interacts with the HIC.

---

## 1. Purpose

This agent defines the integration standards that bind every layer of the Holm Intelligence Complex into a coherent, predictable, maintainable system. It specifies naming conventions for every identifier in the system, version format, browser compatibility requirements, the event bus protocol, error handling standards, accessibility requirements, and performance budgets.

Every other HIC agent must conform to the standards in this document. When there is a conflict between a decision in another agent and a standard defined here, this agent takes precedence. This is the integration contract.

## 2. Naming Conventions

### 2.1 Floor Identifiers

Floor IDs are short, uppercase string codes. They are used in file names, data attributes, event payloads, and human communication. The complete set:

| Floor ID | Name | Position |
|---|---|---|
| `SB2` | Sub-Basement 2: Foundation | Below ground, deepest |
| `SB1` | Sub-Basement 1: Deep Archives | Below ground |
| `G` | Ground Floor: Lobby | Ground level, entry point |
| `F01` | Floor 1 | Above ground |
| `F02` | Floor 2 | Above ground |
| `F03` | Floor 3 | Above ground |
| `F04` | Floor 4 | Above ground |
| `F05` | Floor 5 | Above ground |
| `F06` | Floor 6 | Above ground |
| `F07` | Floor 7 | Above ground |
| `F08` | Floor 8 | Above ground |
| `F09` | Floor 9 | Above ground |
| `F10` | Floor 10 | Above ground |
| `F11` | Floor 11 | Above ground |
| `F12` | Floor 12 | Above ground |
| `F13` | Floor 13 | Above ground |
| `F14` | Floor 14 | Above ground |
| `F15` | Floor 15 | Above ground |
| `F16` | Floor 16 | Above ground |
| `F17` | Floor 17 | Above ground, highest standard floor |
| `R` | Roof: Observatory | Top of building |
| `A` | Annex | Adjacent structure |

**Rules:**
- Standard floors use `F` followed by two zero-padded digits: `F01` through `F17`.
- Sub-basements use `SB` followed by a digit: `SB1`, `SB2`.
- Ground, Roof, and Annex use single-letter codes: `G`, `R`, `A`.
- Floor IDs are case-sensitive. `F07` is valid. `f07` is not.
- No floor ID may contain spaces, hyphens, underscores, or special characters.

### 2.2 Room Identifiers

Room IDs combine the floor ID with a room number: `{floor_id}-R{room_number}`.

Examples:
- `F07-R01` = Floor 7, Room 1
- `F07-R02` = Floor 7, Room 2
- `G-R01` = Ground Floor, Room 1
- `SB1-R03` = Sub-Basement 1, Room 3
- `R-R01` = Roof, Room 1

**Rules:**
- Room numbers are two zero-padded digits: `R01` through `R99`.
- The separator between floor ID and room number is a single hyphen.
- Room numbering starts at `R01` on each floor. There is no `R00`.
- Room numbers are assigned in spatial order: left to right, top to bottom in the floor plan.
- Maximum 99 rooms per floor. If a floor requires more, it must be split into two floors.

### 2.3 File Naming

| File Type | Pattern | Example |
|---|---|---|
| Floor metadata | `{floor_id}.floor.json` | `F07.floor.json` |
| Floor plan SVG | `{floor_id}.floor.svg` | `F07.floor.svg` |
| Document content | `{floor_id}-R{NN}.doc.json` | `F07-R01.doc.json` |
| Import bundle | `*.hic-import.json` | `2026-02-17-security-update.hic-import.json` |

**Rules:**
- All file names are lowercase except the floor ID prefix, which preserves its canonical case (`F07`, `SB1`, `G`, `R`, `A`).
- Extensions are always lowercase: `.json`, `.svg`, `.css`, `.js`.
- No spaces in file names. Use hyphens for human-readable separation in import bundle names.
- The double extension pattern (`.floor.json`, `.floor.svg`, `.doc.json`) is mandatory. It allows glob patterns like `*.floor.json` to select all floor files.

### 2.4 CSS Class Naming

All HIC CSS classes use the `hic-` prefix. No exceptions.

| Class | Purpose |
|---|---|
| `.hic-floor` | Root SVG element of a floor plan |
| `.hic-room` | Interactive room group within a floor |
| `.hic-door` | Door indicator between corridor and room |
| `.hic-wall` | Wall element (corridor or internal) |
| `.hic-wall-corridor` | Corridor wall subtype |
| `.hic-wall-internal` | Internal dividing wall subtype |
| `.hic-label` | Any text label |
| `.hic-label-room-id` | Room ID text (e.g., "F07-R01") |
| `.hic-label-room-title` | Room title text (e.g., "Air-Gap Architecture") |
| `.hic-label-floor-title` | Floor title text |
| `.hic-label-floor-id` | Floor ID text |
| `.hic-glow` | Neon glow effect rectangle |
| `.hic-building` | Root SVG element of building elevation |
| `.hic-building-floor` | Floor band in building elevation |
| `.hic-building-shell` | Building exterior shell |
| `.hic-sky` | Sky/atmosphere background |
| `.hic-layer-bg` | Background SVG layer |
| `.hic-layer-walls` | Walls SVG layer |
| `.hic-layer-rooms` | Rooms SVG layer |
| `.hic-layer-doors` | Doors SVG layer |
| `.hic-layer-labels` | Labels SVG layer |
| `.hic-grid` | Structural grid line (decorative) |
| `.hic-floor-slab` | Floor slab background rectangle |
| `.hic-room-area` | Room background rectangle (within `.hic-room`) |

**Rules:**
- The prefix `hic-` is mandatory for all classes.
- Multi-word class names use hyphens: `hic-room-area`, not `hic-roomArea` or `hic_room_area`.
- Modifier classes follow BEM-like convention with a double hyphen: `hic-room--active`, `hic-room--locked`, `hic-room--draft`.
- State classes use `is-` prefix after `hic-`: `hic-is-focused`, `hic-is-hovered`, `hic-is-loading`.

### 2.5 Data Attribute Naming

| Attribute | Applied To | Values | Purpose |
|---|---|---|---|
| `data-floor` | `.hic-floor`, `.hic-building-floor` | Floor ID (`F07`, `G`, etc.) | Identifies the floor |
| `data-room` | `.hic-room` | Room ID (`F07-R01`) | Identifies the room |
| `data-security` | `.hic-room` | `standard`, `elevated`, `restricted` | Security classification |
| `data-status` | `.hic-room` | `active`, `locked`, `empty`, `draft` | Room/document status |
| `data-version` | `.hic-floor` | Version string (`HIC-1.2.0`) | Floor data version |
| `data-renderer` | `#hic-viewport` | `svg`, `canvas` | Active rendering pipeline |
| `data-theme` | `<html>` | `neon`, `amber`, `high-contrast` | Active color theme |
| `data-connects` | `.hic-door` | `{from}:{to}` pattern | Door connection |
| `data-zoom` | `#hic-viewport` | Numeric zoom level (`1.0`) | Current zoom level |

## 3. Version Format

The HIC version follows the format `HIC-{major}.{minor}.{patch}`.

- **Major** increments when the file structure, IndexedDB schema, or import bundle format changes in a backward-incompatible way. Major version changes require re-importing all content.
- **Minor** increments when new floors, new features, or new themes are added in a backward-compatible way.
- **Patch** increments when existing content is corrected, bugs are fixed, or performance is improved without structural changes.

Examples:
- `HIC-1.0.0` -- Initial release.
- `HIC-1.1.0` -- Added Annex floor.
- `HIC-1.1.1` -- Fixed glow filter on Floor 7 rooms.
- `HIC-2.0.0` -- Restructured IndexedDB schema (breaking change).

The version string appears in:
- `meta/version.json` (the `hicVersion` field).
- Every floor SVG (the `data-version` attribute on the root `<svg>` element).
- The HTML page title during operation: "HIC-1.2.0 | Floor 7 | Security & Integrity".
- Import bundles (the `hicTargetVersion` field).

## 4. Compatibility Matrix

### 4.1 Browser Support

| Browser | Minimum Version | Notes |
|---|---|---|
| Chromium (Chrome, Edge, Brave) | 90+ | Primary target. Full support. |
| Firefox | 90+ | Full support. |
| Safari | 15+ | Full support. Tested on macOS. |
| Safari iOS | 15+ | Touch navigation. Reduced glow effects. |
| Legacy browsers | Not supported | No IE11. No pre-Chromium Edge. |

The minimum version floor of 90 ensures support for: ES2020+ syntax, CSS custom properties, CSS `gap` in flexbox, `IntersectionObserver`, `ResizeObserver`, `IndexedDB` 2.0, service workers, `<input type="file">` for USB import, SVG 1.1 full support, Canvas 2D context.

### 4.2 Screen Size Support

| Category | Viewport Range | Layout |
|---|---|---|
| Desktop large | 1920px+ | Full building elevation + floor plan side by side |
| Desktop standard | 1280-1919px | Building elevation collapsible, floor plan full width |
| Laptop | 1024-1279px | Building elevation as overlay panel |
| Tablet landscape | 768-1023px | Floor plan only, building nav via menu |
| Tablet portrait | 600-767px | Floor plan with horizontal scroll |
| Mobile | < 600px | Simplified room list, no floor plan SVG |

The SVG floor plans use a fixed `viewBox="0 0 1200 800"` and scale responsively via CSS `width: 100%; height: auto;`. Below 600px viewport width, the SVG floor plan is replaced with a styled `<ul>` list of rooms because SVG interaction targets become too small for touch.

### 4.3 Input Method Support

| Input Method | Support Level | Notes |
|---|---|---|
| Mouse | Full | Click, hover, scroll zoom |
| Keyboard | Full | Tab navigation, Enter/Space activation, arrow keys for floor nav |
| Touch | Full | Tap, pinch zoom, swipe for floor navigation |
| Gamepad | Not supported | No current requirement |
| Voice | Partial | Via browser's built-in accessibility; no custom voice commands |

## 5. Event Bus Specification

### 5.1 Architecture

The HIC event bus uses the browser's native `CustomEvent` system dispatched on the `document` object. No event library. No pub/sub framework. Just `document.dispatchEvent()` and `document.addEventListener()`.

All HIC events use the `hic:` namespace prefix.

### 5.2 Event Catalog

| Event Name | Dispatched When | Detail Payload |
|---|---|---|
| `hic:init` | HIC core initialization complete | `{ version: string, renderer: string }` |
| `hic:floor-enter` | Operator navigates to a floor | `{ floorId: string, fromFloorId: string \| null }` |
| `hic:floor-leave` | Operator leaves a floor | `{ floorId: string, toFloorId: string }` |
| `hic:room-open` | Operator opens a room (clicks/activates) | `{ roomId: string, floorId: string, docId: string }` |
| `hic:room-close` | Operator closes a room view | `{ roomId: string, floorId: string }` |
| `hic:room-focus` | Room receives keyboard/mouse focus | `{ roomId: string, floorId: string }` |
| `hic:zoom-change` | Zoom level changes | `{ level: number, min: number, max: number }` |
| `hic:theme-change` | Theme switches | `{ theme: string, previousTheme: string }` |
| `hic:search-open` | Search panel opens | `{}` |
| `hic:search-results` | Search returns results | `{ query: string, resultCount: number, results: Array }` |
| `hic:search-close` | Search panel closes | `{}` |
| `hic:import-start` | USB import begins | `{ fileName: string }` |
| `hic:import-complete` | USB import succeeds | `{ floorsUpdated: Array, docsAdded: number, docsUpdated: number }` |
| `hic:import-error` | USB import fails | `{ errorCode: string, message: string }` |
| `hic:cache-update` | Cache invalidation completes | `{ floorsInvalidated: Array }` |
| `hic:error` | Any recoverable error | `{ code: string, message: string, context: object }` |

### 5.3 Event Dispatch Example

```javascript
// Dispatching a floor-enter event
function enterFloor(floorId, fromFloorId) {
  document.dispatchEvent(new CustomEvent('hic:floor-enter', {
    bubbles: true,
    detail: {
      floorId: floorId,
      fromFloorId: fromFloorId || null
    }
  }));
}

// Listening for floor-enter events
document.addEventListener('hic:floor-enter', (event) => {
  const { floorId, fromFloorId } = event.detail;
  console.log(`Entering floor ${floorId} from ${fromFloorId}`);
  // Update title bar
  document.title = `HIC-${HIC_VERSION} | ${floorId} | ${getFloorTitle(floorId)}`;
  // Trigger adjacent floor preload
  preloadAdjacentFloors(floorId);
  // Update building elevation highlight
  highlightBuildingFloor(floorId);
});
```

### 5.4 Event Ordering Guarantees

Events are dispatched synchronously within each pipeline. The guaranteed order for a floor navigation is:

1. `hic:floor-leave` (from the old floor)
2. `hic:floor-enter` (to the new floor)
3. `hic:zoom-change` (reset zoom to 1.0 on new floor)

The guaranteed order for a room interaction is:

1. `hic:room-focus` (on focus/hover)
2. `hic:room-open` (on activation)
3. `hic:room-close` (on deactivation)

No event may be dispatched out of this order. If a rendering error prevents a floor from loading, `hic:floor-enter` is not dispatched; `hic:error` is dispatched instead.

## 6. Error Handling Standards

### 6.1 Error Code Registry

All HIC errors use the format `HIC-ERR-{CATEGORY}-{SPECIFIC}`.

| Error Code | Category | Description | Recovery |
|---|---|---|---|
| `HIC-ERR-RENDER-SVG-PARSE` | Render | Floor SVG failed to parse | Show error overlay, offer Canvas fallback |
| `HIC-ERR-RENDER-SVG-MISSING` | Render | Floor SVG file not found | Show "floor under construction" placeholder |
| `HIC-ERR-RENDER-CANVAS-INIT` | Render | Canvas context creation failed | Fatal: display text-only fallback |
| `HIC-ERR-DATA-FLOOR-PARSE` | Data | Floor JSON failed to parse | Show error overlay with raw JSON link |
| `HIC-ERR-DATA-FLOOR-MISSING` | Data | Floor JSON file not found | Show "floor under construction" placeholder |
| `HIC-ERR-DATA-DOC-PARSE` | Data | Document JSON failed to parse | Show error overlay with document ID |
| `HIC-ERR-DATA-DOC-MISSING` | Data | Document JSON file not found | Show "document pending" placeholder |
| `HIC-ERR-CACHE-IDB-OPEN` | Cache | IndexedDB failed to open | Fall back to in-memory storage |
| `HIC-ERR-CACHE-IDB-READ` | Cache | IndexedDB read transaction failed | Retry once, then fall back to filesystem |
| `HIC-ERR-CACHE-IDB-WRITE` | Cache | IndexedDB write transaction failed | Retry once, then warn operator |
| `HIC-ERR-CACHE-SW-REG` | Cache | Service worker registration failed | Continue without offline caching |
| `HIC-ERR-IMPORT-CHECKSUM` | Import | Import bundle checksum mismatch | Reject import, display checksum values |
| `HIC-ERR-IMPORT-VERSION` | Import | Import bundle targets incompatible version | Reject import, display version mismatch |
| `HIC-ERR-IMPORT-PARSE` | Import | Import bundle JSON failed to parse | Reject import, display parse error |
| `HIC-ERR-IMPORT-FLOOR` | Import | Individual floor in bundle failed to merge | Skip floor, continue with remaining floors |
| `HIC-ERR-NAV-FLOOR-INVALID` | Navigation | Requested floor ID does not exist | Navigate to Ground floor (G) |
| `HIC-ERR-NAV-ROOM-INVALID` | Navigation | Requested room ID does not exist | Stay on current floor, highlight error |

### 6.2 Error Display

Errors are displayed as an overlay panel within `#hic-viewport`. The panel has a dark red border, the error code in monospace, the error message in plain language, and a "DISMISS" button. The panel does not block the entire interface -- the operator can dismiss it and continue navigating if the error is non-fatal.

```javascript
function displayError(code, message, context) {
  const overlay = document.createElement('div');
  overlay.className = 'hic-error-overlay';
  overlay.setAttribute('role', 'alert');
  overlay.setAttribute('aria-live', 'assertive');
  overlay.innerHTML = `
    <div class="hic-error-panel">
      <div class="hic-error-code">${escapeHtml(code)}</div>
      <div class="hic-error-message">${escapeHtml(message)}</div>
      <button class="hic-error-dismiss" onclick="this.closest('.hic-error-overlay').remove()">
        DISMISS
      </button>
    </div>
  `;
  document.getElementById('hic-viewport').appendChild(overlay);

  // Also dispatch error event for programmatic handling
  document.dispatchEvent(new CustomEvent('hic:error', {
    bubbles: true,
    detail: { code, message, context: context || {} }
  }));
}
```

### 6.3 Error Logging

All errors are logged to the browser console with the format: `[HIC] {ERROR_CODE}: {message}`. If IndexedDB is available, errors are also written to a `diagnostics` key in the `userState` store as a rotating buffer of the last 50 errors.

## 7. Accessibility Requirements

### 7.1 Keyboard Navigation

Every interactive element in the HIC must be reachable and activatable via keyboard alone.

**Building elevation view:**
- `Tab` moves focus between floor bands, top to bottom.
- `Shift+Tab` moves focus in reverse.
- `Enter` or `Space` activates the focused floor (navigates to floor plan).
- `Home` focuses the Roof (R). `End` focuses Sub-Basement 2 (SB2).
- `Arrow Up` and `Arrow Down` move focus between adjacent floors.

**Floor plan view:**
- `Tab` moves focus between rooms in spatial order (left to right, top to bottom).
- `Shift+Tab` moves focus in reverse.
- `Enter` or `Space` opens the focused room (shows document).
- `Escape` closes the current document and returns to the floor plan.
- `Page Up` navigates to the floor above. `Page Down` navigates to the floor below.
- `Backspace` or `Alt+Left` returns to the building elevation view.

**Document view:**
- Standard document scrolling with arrow keys.
- `Escape` closes the document and returns to the floor plan.
- `Ctrl+F` or `/` opens the search panel.

### 7.2 Screen Reader Support

Every interactive element must have an `aria-label` or accessible text content.

| Element | ARIA Role | ARIA Label |
|---|---|---|
| `#hic-viewport` | `application` | "Holm Intelligence Complex" |
| `.hic-building` | `navigation` | "Building Elevation" |
| `.hic-building-floor` | `link` | "{Floor Name}" |
| `.hic-floor` | `img` | "{Floor Name}" |
| `.hic-room` | `button` | "Room {N}: {Title}" |
| `.hic-door` | Not interactive | `aria-hidden="true"` |
| `.hic-layer-bg` | Not interactive | `aria-hidden="true"` |
| Error overlay | `alert` | Dynamic error message |
| Search panel | `search` | "Search Documents" |

Screen readers must announce:
- Floor transitions: "Entering Floor 7: Security and Integrity."
- Room activation: "Opening Room 1: Air-Gap Architecture."
- Errors: The full error message.
- Import completion: The import summary.

### 7.3 High Contrast Mode

The `high-contrast.css` theme meets WCAG AAA requirements:

- Text contrast ratio: 7:1 minimum against background.
- Interactive element boundaries: 3:1 minimum against adjacent colors.
- Focus indicators: 3px solid outline in high-visibility color (#FFFFFF or #FFFF00).
- No information conveyed by color alone. Room status uses both color and text labels.
- Glow effects are disabled (replaced with solid borders).
- Background is pure black (`#000000`). Text is pure white (`#FFFFFF`).

```css
/* css/hic-themes/high-contrast.css */
:root[data-theme="high-contrast"] {
  --hic-bg: #000000;
  --hic-cyan: #00FFFF;
  --hic-magenta: #FF00FF;
  --hic-text: #FFFFFF;
  --hic-text-secondary: #CCCCCC;
  --hic-border: #FFFFFF;
  --hic-focus: #FFFF00;
  --hic-error: #FF4444;
}

:root[data-theme="high-contrast"] .hic-glow {
  display: none; /* Disable glow effects */
}

:root[data-theme="high-contrast"] .hic-room:focus {
  outline: 3px solid var(--hic-focus);
  outline-offset: 2px;
}

:root[data-theme="high-contrast"] .hic-room[data-status="active"]::after {
  content: '[ACTIVE]';
}

:root[data-theme="high-contrast"] .hic-room[data-status="locked"]::after {
  content: '[LOCKED]';
}

:root[data-theme="high-contrast"] .hic-room[data-status="draft"]::after {
  content: '[DRAFT]';
}
```

### 7.4 Reduced Motion

The HIC respects the `prefers-reduced-motion` media query:

```css
@media (prefers-reduced-motion: reduce) {
  .hic-glow {
    animation: none !important;
    opacity: 0.5;
  }
  .hic-room {
    transition: none !important;
  }
  #hic-loader {
    animation: none !important;
    opacity: 1;
  }
}
```

When reduced motion is active, all CSS animations are disabled, all transitions are instant, and glow effects are rendered at static opacity. The interface remains fully functional; it just does not move.

## 8. Performance Budgets

### 8.1 Rendering Performance

| Metric | Budget | Measurement Method |
|---|---|---|
| Frame rate during interaction | 60 fps | `requestAnimationFrame` timestamp delta |
| Frame rate during idle | 0 fps (no unnecessary repaints) | Performance monitor shows no rAF calls |
| SVG floor plan initial render | < 200ms | `performance.mark()` around SVG insertion |
| Canvas building elevation render | < 50ms | `performance.mark()` around draw loop |

### 8.2 Interaction Response

| Interaction | Response Budget | Definition of "Response" |
|---|---|---|
| Room click/tap | < 100ms | Document panel begins rendering |
| Floor navigation | < 500ms | New floor SVG is visible and interactive |
| Search query (per keystroke) | < 150ms | Results list updates |
| Zoom in/out | < 50ms | Viewport transform applied |
| Theme switch | < 200ms | All colors updated |
| USB import (per floor) | < 2000ms | Floor data merged into IndexedDB |

### 8.3 Load Performance

| Metric | Budget |
|---|---|
| Cold start (first load, no cache) | < 2000ms to interactive |
| Warm start (service worker cache) | < 500ms to interactive |
| Time to first meaningful paint | < 1000ms (building elevation visible) |
| Total blocking time | < 200ms |

### 8.4 Memory Budget

| Metric | Budget | Notes |
|---|---|---|
| JavaScript heap (idle) | < 20 MB | After initial load, no floor open |
| JavaScript heap (one floor open) | < 35 MB | With floor SVG in DOM, document data in memory |
| JavaScript heap (peak during import) | < 50 MB | Import bundle parsed and merged |
| DOM node count (building view) | < 500 | Building elevation SVG elements |
| DOM node count (floor view) | < 2000 | Floor SVG elements + UI controls |

### 8.5 Performance Monitoring

The HIC includes a built-in performance monitor that can be activated by the operator with `Ctrl+Shift+P`. When active, it displays a small overlay in the bottom-right corner showing:

- Current FPS (updated per frame).
- Last interaction response time.
- Current JavaScript heap size (via `performance.memory` where supported).
- Cache status (shell cached, floors cached count, docs cached count).

The monitor is implemented in `hic-core.js` and adds less than 1 KB to the JavaScript budget. It is disabled by default and has zero performance impact when inactive.

```javascript
// Performance monitor (hic-core.js excerpt)
const perfMonitor = {
  active: false,
  frameCount: 0,
  lastFrameTime: 0,
  overlay: null,

  toggle() {
    this.active = !this.active;
    if (this.active) {
      this.overlay = document.createElement('div');
      this.overlay.className = 'hic-perf-monitor';
      this.overlay.setAttribute('aria-hidden', 'true');
      document.getElementById('hic-viewport').appendChild(this.overlay);
      this.tick();
    } else if (this.overlay) {
      this.overlay.remove();
      this.overlay = null;
    }
  },

  tick() {
    if (!this.active) return;
    const now = performance.now();
    const delta = now - this.lastFrameTime;
    this.lastFrameTime = now;
    const fps = delta > 0 ? Math.round(1000 / delta) : 0;

    this.overlay.textContent = `FPS: ${fps} | Heap: ${
      performance.memory
        ? Math.round(performance.memory.usedJSHeapSize / 1048576) + 'MB'
        : 'N/A'
    }`;

    requestAnimationFrame(() => this.tick());
  }
};

document.addEventListener('keydown', (e) => {
  if (e.ctrlKey && e.shiftKey && e.key === 'P') {
    e.preventDefault();
    perfMonitor.toggle();
  }
});
```

## 9. Cross-Agent Integration Checklist

Every HIC agent implementation must pass the following integration checks before it is considered complete:

### 9.1 Naming Compliance

- [ ] All floor IDs use canonical format from Section 2.1.
- [ ] All room IDs use `{floor_id}-R{NN}` format from Section 2.2.
- [ ] All file names follow patterns from Section 2.3.
- [ ] All CSS classes use `hic-` prefix from Section 2.4.
- [ ] All data attributes use names from Section 2.5.
- [ ] Version strings follow `HIC-{major}.{minor}.{patch}` from Section 3.

### 9.2 Event Compliance

- [ ] All navigation dispatches `hic:floor-enter` and `hic:floor-leave`.
- [ ] All room interactions dispatch `hic:room-open` and `hic:room-close`.
- [ ] Zoom changes dispatch `hic:zoom-change`.
- [ ] Event ordering follows guarantees from Section 5.4.
- [ ] All event payloads match the schema from Section 5.2.

### 9.3 Accessibility Compliance

- [ ] All interactive elements have `tabindex="0"` and appropriate `role`.
- [ ] All interactive elements have `aria-label` text.
- [ ] Keyboard navigation works without mouse from Section 7.1.
- [ ] High contrast theme is functional from Section 7.3.
- [ ] Reduced motion is respected from Section 7.4.

### 9.4 Performance Compliance

- [ ] Floor load under 500ms from Section 8.2.
- [ ] Click response under 100ms from Section 8.2.
- [ ] Rendering at 60fps from Section 8.1.
- [ ] Cold start under 2000ms from Section 8.3.
- [ ] Memory under budgets from Section 8.4.
- [ ] Total asset size under 5MB from HIC-016 Section 2.2.

### 9.5 Offline Compliance

- [ ] All assets cached by service worker from HIC-018 Section 2.
- [ ] IndexedDB schema matches HIC-018 Section 3.
- [ ] Cache invalidation uses version stamps from HIC-018 Section 4.
- [ ] USB import follows bundle format from HIC-018 Section 6.
- [ ] Adjacent floor preloading active from HIC-018 Section 5.

---

## Appendix A: Complete Floor JSON Schema

This is the authoritative JSON schema for `.floor.json` files. All floor files must validate against this schema.

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://holm.chat/hic/schemas/floor.schema.json",
  "title": "HIC Floor Metadata",
  "type": "object",
  "required": ["floorId", "title", "version", "checksum", "status", "roomCount", "rooms", "adjacentFloors", "svgFile", "modified"],
  "properties": {
    "floorId": {
      "type": "string",
      "pattern": "^(F[0-1][0-9]|F[0-9]|SB[1-2]|G|R|A)$",
      "description": "Canonical floor identifier"
    },
    "title": {
      "type": "string",
      "maxLength": 100,
      "description": "Human-readable floor title"
    },
    "version": {
      "type": "string",
      "pattern": "^HIC-[0-9]+\\.[0-9]+\\.[0-9]+$",
      "description": "HIC version at last floor update"
    },
    "checksum": {
      "type": "string",
      "pattern": "^sha256:[a-f0-9]{64}$",
      "description": "SHA-256 checksum of floor content"
    },
    "status": {
      "type": "string",
      "enum": ["active", "draft", "archived"],
      "description": "Floor lifecycle status"
    },
    "roomCount": {
      "type": "integer",
      "minimum": 0,
      "maximum": 99,
      "description": "Number of rooms on this floor"
    },
    "rooms": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["roomId", "title", "status", "security", "docId", "wordCount", "modified"],
        "properties": {
          "roomId": {
            "type": "string",
            "pattern": "^(F[0-1][0-9]|F[0-9]|SB[1-2]|G|R|A)-R[0-9]{2}$",
            "description": "Canonical room identifier"
          },
          "title": {
            "type": "string",
            "maxLength": 100,
            "description": "Human-readable room title"
          },
          "status": {
            "type": "string",
            "enum": ["active", "locked", "empty", "draft"],
            "description": "Room content status"
          },
          "security": {
            "type": "string",
            "enum": ["standard", "elevated", "restricted"],
            "description": "Security classification level"
          },
          "docId": {
            "type": "string",
            "description": "Reference to the institution document ID"
          },
          "wordCount": {
            "type": "integer",
            "minimum": 0,
            "description": "Word count of the document in this room"
          },
          "modified": {
            "type": "string",
            "format": "date-time",
            "description": "Last modification timestamp (ISO 8601)"
          }
        }
      }
    },
    "adjacentFloors": {
      "type": "object",
      "properties": {
        "above": { "type": ["string", "null"] },
        "below": { "type": ["string", "null"] }
      },
      "description": "Floor IDs of adjacent floors for navigation"
    },
    "svgFile": {
      "type": "string",
      "pattern": "^.+\\.floor\\.svg$",
      "description": "Filename of the floor plan SVG"
    },
    "modified": {
      "type": "string",
      "format": "date-time",
      "description": "Last modification timestamp of floor data"
    }
  }
}
```

## Appendix B: Complete Document JSON Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://holm.chat/hic/schemas/document.schema.json",
  "title": "HIC Document Content",
  "type": "object",
  "required": ["docId", "floorId", "roomId", "title", "documentReference", "version", "status", "security", "content", "metadata", "checksum"],
  "properties": {
    "docId": {
      "type": "string",
      "pattern": "^(F[0-1][0-9]|F[0-9]|SB[1-2]|G|R|A)-R[0-9]{2}$"
    },
    "floorId": {
      "type": "string",
      "pattern": "^(F[0-1][0-9]|F[0-9]|SB[1-2]|G|R|A)$"
    },
    "roomId": {
      "type": "string",
      "pattern": "^(F[0-1][0-9]|F[0-9]|SB[1-2]|G|R|A)-R[0-9]{2}$"
    },
    "title": { "type": "string", "maxLength": 200 },
    "documentReference": { "type": "string" },
    "version": { "type": "string", "pattern": "^[0-9]+\\.[0-9]+\\.[0-9]+$" },
    "status": { "type": "string", "enum": ["ratified", "draft", "archived", "superseded"] },
    "security": { "type": "string", "enum": ["standard", "elevated", "restricted"] },
    "content": {
      "type": "object",
      "required": ["format", "body"],
      "properties": {
        "format": { "type": "string", "enum": ["markdown", "plaintext"] },
        "body": { "type": "string" },
        "sections": {
          "type": "array",
          "items": {
            "type": "object",
            "required": ["id", "title", "offset"],
            "properties": {
              "id": { "type": "string" },
              "title": { "type": "string" },
              "offset": { "type": "integer", "minimum": 0 }
            }
          }
        }
      }
    },
    "metadata": {
      "type": "object",
      "required": ["author", "created", "modified", "wordCount"],
      "properties": {
        "author": { "type": "string" },
        "created": { "type": "string", "format": "date-time" },
        "modified": { "type": "string", "format": "date-time" },
        "wordCount": { "type": "integer", "minimum": 0 },
        "depends_on": { "type": "array", "items": { "type": "string" } },
        "depended_upon_by": { "type": "array", "items": { "type": "string" } }
      }
    },
    "checksum": {
      "type": "string",
      "pattern": "^sha256:[a-f0-9]{64}$"
    }
  }
}
```

## Appendix C: Complete Import Bundle Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://holm.chat/hic/schemas/import-bundle.schema.json",
  "title": "HIC Import Bundle",
  "type": "object",
  "required": ["bundleVersion", "hicTargetVersion", "created", "creator", "description", "floors", "bundleChecksum"],
  "properties": {
    "bundleVersion": {
      "type": "string",
      "pattern": "^[0-9]+\\.[0-9]+\\.[0-9]+$",
      "description": "Version of the import bundle format itself"
    },
    "hicTargetVersion": {
      "type": "string",
      "pattern": "^HIC-[0-9]+\\.[0-9]+\\.[0-9]+$",
      "description": "HIC version this bundle is compatible with"
    },
    "created": { "type": "string", "format": "date-time" },
    "creator": { "type": "string", "maxLength": 100 },
    "description": { "type": "string", "maxLength": 500 },
    "floors": {
      "type": "object",
      "additionalProperties": {
        "type": "object",
        "required": ["version", "checksum", "floorJson", "floorSvg", "documents"],
        "properties": {
          "version": { "type": "integer", "minimum": 1 },
          "checksum": { "type": "string", "pattern": "^sha256:[a-f0-9]{64}$" },
          "floorJson": { "$ref": "floor.schema.json" },
          "floorSvg": { "type": "string", "description": "Complete SVG source as string" },
          "documents": {
            "type": "object",
            "additionalProperties": { "$ref": "document.schema.json" }
          }
        }
      }
    },
    "searchIndexPatch": {
      "type": "object",
      "properties": {
        "added": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "docId": { "type": "string" },
              "terms": { "type": "array", "items": { "type": "string" } }
            }
          }
        },
        "removed": {
          "type": "array",
          "items": { "type": "string" }
        }
      }
    },
    "bundleChecksum": {
      "type": "string",
      "pattern": "^sha256:[a-f0-9]{64}$"
    }
  }
}
```

---

*End of STAGE4-HIC-OFFLINE. This document is the authoritative specification for Agents 16 through 19 of the Holm Intelligence Complex. All implementation must conform to the structures, schemas, naming conventions, event protocols, performance budgets, and accessibility standards defined herein. When this document is updated, its version stamp in `meta/version.json` must be incremented and all cached copies must be invalidated.*
