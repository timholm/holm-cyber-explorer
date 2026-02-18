# HIC Visual Design System

## Neon Aesthetics, Color Theory & Typography Standards

**Classification:** HOLM-VIS-004
**Revision:** 3.1
**Effective Date:** 2026-02-17
**Maintainer:** Visual Systems Division, Floor 47

---

> *"In the void, light is language. Every photon carries meaning.
> The building does not decorate -- it communicates."*
>
> -- Holm Doctrine, Appendix V: On Luminance

---

## Table of Contents

1. [Design Philosophy](#1-design-philosophy)
2. [The Neon Palette](#2-the-neon-palette)
3. [Neon Signage System](#3-neon-signage-system)
4. [Typography Hierarchy](#4-typography-hierarchy)
5. [Holographic Display Standards](#5-holographic-display-standards)
6. [Environmental Lighting](#6-environmental-lighting)
7. [Icon & Symbol Library](#7-icon--symbol-library)
8. [Print Compatibility](#8-print-compatibility)
9. [Appendices](#9-appendices)

---

## 1. Design Philosophy

### 1.1 The Void and the Light

The HIC skyscraper rises from darkness. This is not an aesthetic accident. The
foundational principle of the Holm Intelligence Complex visual system is the
dialectic between **void** and **signal** -- between the black expanse of
ignorance and the precise, burning lines of knowledge that cut through it.

Every surface of the HIC defaults to deep black (`#0a0a0f`). Not the
comfortable dark gray of consumer dark modes, but a true absence -- the color
of a powered-down terminal, of a document yet to be written, of a question
not yet asked. Against this void, information announces itself as light.
Content does not sit passively on a page. It *emits*. It *radiates*. Each
word, each diagram, each status indicator is a neon filament burning against
the dark, demanding to be read.

This philosophy serves a practical purpose: in a documentation system of the
HIC's scale -- hundreds of floors, thousands of documents, millions of data
points -- the eye must be guided with absolute precision. There is no room for
ambiguity. When a crimson light pulses on Floor 23, every operator in the
building knows without checking: something is critically wrong with that
document cluster. When Electric Blue lines trace a path through the atrium,
every reader knows: follow that light to the primary content.

The darkness is not hostile. It is *restful*. It is the negative space that
gives the signals meaning. A neon sign in a well-lit room is kitsch. A neon
sign in total darkness is a beacon.

### 1.2 Form Follows Function

No visual element in the HIC exists for decoration alone. Every glow, every
color shift, every flickering indicator maps directly to a documentation
concept:

| Visual Element | Documentation Concept |
|---|---|
| Brightness of a floor's exterior glow | Completeness of that floor's document corpus |
| Color of a corridor's ambient light | Domain classification of the content within |
| Pulse rate of a status indicator | Time since last review or update |
| Depth of a holographic projection | Importance/priority level of the content |
| Sharpness of a neon edge | Confidence level of the information |
| Presence of flickering | Content flagged for review or suspected decay |

This mapping is absolute. A designer working within the HIC system does not
choose colors based on mood or brand preference. They choose colors based on
what the content *is* and what it *needs to communicate*. The visual language
is a protocol, not a palette.

### 1.3 Digital Brutalism

The HIC rejects polish. It rejects the rounded corners and pastel gradients of
contemporary design trends. Its aesthetic lineage traces through Brutalist
architecture -- raw concrete, exposed structure, honest materials -- and
translates those principles into the digital domain:

**Raw structure is visible.** The grid is not hidden. Column lines, baseline
grids, and alignment markers are part of the visible design. A document in the
HIC shows its skeleton. Margins are explicit. Padding is measured and labeled.
The reader sees the architecture of the information, not just the information
itself.

**Monospace is primary.** In the HIC, proportional typefaces are a concession,
not a default. The monospace grid is the native language of the building. Every
character occupies the same width. Every line aligns. Documents are treated as
what they fundamentally are: structured text. Source code. Data. The monospace
grid is the honest representation of this truth.

**No gratuitous animation.** Motion in the HIC carries meaning. A pulse means
the data is live. A fade means the content is aging. A flicker means attention
is required. There are no loading spinners chosen for delight, no page
transitions designed for wow-factor. If something moves, something is
happening. If nothing is happening, nothing moves.

**Density over whitespace.** The HIC is a skyscraper, not a suburban lawn.
Space is used efficiently. Information density is high. Readers are treated as
professionals who can process dense, well-structured content. Excessive
whitespace is wasted darkness -- void with no signal to justify it.

```
DESIGN AXIOMS
=============
01. Dark is default. Light is information.
02. Color is classification. Never decoration.
03. Motion is status. Never spectacle.
04. Density is respect. Never clutter.
05. Monospace is truth. Proportional is compromise.
06. The grid is visible. Structure is not hidden.
07. Every pixel justifies its photon budget.
```

---

## 2. The Neon Palette

### 2.1 Palette Overview

The HIC palette is divided into seven **Primary Signal Colors** and a set of
**Structural Neutrals**. Each signal color has a defined semantic meaning that
is consistent across every floor, every display, and every document in the
building.

No color is used outside its semantic assignment. If a designer needs "a nice
blue for this header," they do not reach for Electric Blue unless that header
represents primary active knowledge content. The palette is a vocabulary, and
using the wrong word is a lie.

### 2.2 Primary Signal Colors

#### Electric Blue -- The Knowledge Signal

| Property | Value |
|---|---|
| Primary Hex | `#7b8cde` |
| Bright Variant | `#8b9cf0` |
| Dim Variant | `#5a6aad` |
| Glow Variant | `rgba(123, 140, 222, 0.35)` |
| CSS Variable | `--hic-electric-blue` |
| Semantic Role | Primary knowledge, active content, navigation links, highlighted text |
| Usage Zones | Floors 1-10 (Public Atrium), all primary navigation elements |

Electric Blue is the beating heart of the HIC palette. It is the color of a
live hyperlink, a verified fact, an active document. When a reader sees
Electric Blue, they know: this content is current, this path is open, this
knowledge is accessible.

```css
/* Electric Blue application */
.hic-link,
.hic-active-content,
.hic-nav-primary {
  color: #7b8cde;
  text-shadow: 0 0 8px rgba(123, 140, 222, 0.4);
}

.hic-link:hover {
  color: #8b9cf0;
  text-shadow: 0 0 12px rgba(139, 156, 240, 0.6);
}

.hic-glow-border--blue {
  border: 1px solid #7b8cde;
  box-shadow:
    0 0 4px rgba(123, 140, 222, 0.3),
    inset 0 0 4px rgba(123, 140, 222, 0.1);
}
```

Electric Blue operates at a color temperature that the human eye processes as
both trustworthy and energetic. It sits in the blue-violet range that avoids
the coldness of pure blue and the anxiety of pure violet. Prolonged exposure
does not cause fatigue -- a critical property for a color that dominates the
primary reading experience.

#### Amber Warning -- The Caution Signal

| Property | Value |
|---|---|
| Primary Hex | `#c49a3a` |
| Bright Variant | `#d4aa4a` |
| Dim Variant | `#a47a2a` |
| Glow Variant | `rgba(196, 154, 58, 0.35)` |
| CSS Variable | `--hic-amber-warning` |
| Semantic Role | Caution, draft content, review-needed flags, aging indicators |
| Usage Zones | Floors 31-40 (Review & QA), all draft-state markers |

Amber Warning is the color of a document that needs attention but is not yet
in crisis. A draft that has been open too long. A review cycle that is
overdue. A reference that may be stale. Amber says: *proceed, but verify*.

```css
/* Amber Warning application */
.hic-draft-marker,
.hic-review-needed,
.hic-caution-badge {
  color: #c49a3a;
  text-shadow: 0 0 6px rgba(196, 154, 58, 0.4);
}

.hic-status-bar--warning {
  background: linear-gradient(
    90deg,
    rgba(196, 154, 58, 0.1) 0%,
    rgba(196, 154, 58, 0.3) 50%,
    rgba(196, 154, 58, 0.1) 100%
  );
  border-left: 3px solid #c49a3a;
}
```

Amber is deliberately chosen over yellow. Pure yellow against a dark
background creates excessive contrast that reads as alarm rather than caution.
The warm, muted amber communicates concern without panic. It is the difference
between a raised eyebrow and a shout.

#### Crimson Alert -- The Danger Signal

| Property | Value |
|---|---|
| Primary Hex | `#c45a5a` |
| Bright Variant | `#d46a6a` |
| Dim Variant | `#a43a3a` |
| Glow Variant | `rgba(196, 90, 90, 0.35)` |
| CSS Variable | `--hic-crimson-alert` |
| Semantic Role | Critical systems, security warnings, deprecated markers, broken links |
| Usage Zones | Floors 41-50 (Security Division), all critical-state indicators |

Crimson Alert is the color of emergency. Deprecated documentation that is
still being referenced. Security vulnerabilities. Broken links that lead to
dead pages. Missing content that should exist. When Crimson light floods a
section of the HIC, operators mobilize.

```css
/* Crimson Alert application */
.hic-deprecated,
.hic-security-warning,
.hic-critical-badge {
  color: #c45a5a;
  text-shadow: 0 0 10px rgba(196, 90, 90, 0.5);
  animation: hic-pulse-crimson 2s ease-in-out infinite;
}

@keyframes hic-pulse-crimson {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

.hic-alert-border {
  border: 2px solid #c45a5a;
  box-shadow:
    0 0 8px rgba(196, 90, 90, 0.4),
    0 0 16px rgba(196, 90, 90, 0.2);
}
```

Crimson is muted from pure red intentionally. The HIC environment is dark, and
pure red (`#ff0000`) at high saturation causes visual afterimages and
discomfort during extended viewing. The desaturated crimson retains urgency
while remaining tolerable for the operators who may need to stare at a problem
screen for hours while resolving it.

#### Emerald Status -- The Health Signal

| Property | Value |
|---|---|
| Primary Hex | `#5fa85f` |
| Bright Variant | `#6fb86f` |
| Dim Variant | `#4f884f` |
| Glow Variant | `rgba(95, 168, 95, 0.35)` |
| CSS Variable | `--hic-emerald-status` |
| Semantic Role | Verified content, healthy systems, approved documents, passing checks |
| Usage Zones | Floors 11-20 (Verified Archives), all approval indicators |

Emerald Status is the color of confidence. A document that has passed all
reviews. A system check that returned healthy. An approval stamp from a senior
editor. Emerald is the signal that the reader can trust what they are reading
without further verification.

```css
/* Emerald Status application */
.hic-verified,
.hic-approved,
.hic-health-ok {
  color: #5fa85f;
  text-shadow: 0 0 6px rgba(95, 168, 95, 0.3);
}

.hic-badge--approved {
  background: rgba(95, 168, 95, 0.15);
  border: 1px solid #5fa85f;
  color: #6fb86f;
  padding: 2px 8px;
  font-family: var(--hic-mono);
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
}
```

#### Violet Accent -- The Meta Signal

| Property | Value |
|---|---|
| Primary Hex | `#9b59b6` |
| Bright Variant | `#ab69c6` |
| Dim Variant | `#7b3996` |
| Glow Variant | `rgba(155, 89, 182, 0.35)` |
| CSS Variable | `--hic-violet-accent` |
| Semantic Role | Meta-content, commentary, cross-references, annotations, footnotes |
| Usage Zones | Floors 21-30 (Commentary Mezzanine), all meta-layer displays |

Violet Accent is the color of content *about* content. Annotations that an
editor has left on a document. Cross-references that link one domain to
another. Footnotes that provide context without cluttering the primary text.
Violet exists in the margin -- literally and figuratively.

```css
/* Violet Accent application */
.hic-annotation,
.hic-cross-ref,
.hic-footnote,
.hic-meta-content {
  color: #9b59b6;
  text-shadow: 0 0 5px rgba(155, 89, 182, 0.3);
  font-style: italic;
}

.hic-sidebar--meta {
  border-left: 2px solid #9b59b6;
  padding-left: 1rem;
  background: rgba(155, 89, 182, 0.05);
}
```

Violet sits at the far edge of the visible spectrum, and this peripheral
quality maps perfectly to its semantic role. Meta-content is peripheral
information -- important but secondary, visible but not dominant.

#### Cyan Data -- The Archive Signal

| Property | Value |
|---|---|
| Primary Hex | `#00bcd4` |
| Bright Variant | `#26c6da` |
| Dim Variant | `#0097a7` |
| Glow Variant | `rgba(0, 188, 212, 0.35)` |
| CSS Variable | `--hic-cyan-data` |
| Semantic Role | Raw data streams, archives, logs, telemetry, machine-generated content |
| Usage Zones | Floors 51-60 (Data Vaults), all log/archive displays |

Cyan Data is the color of the machine. Raw logs. Data streams. Automated
reports. Content that was generated by systems rather than written by humans.
Cyan says: *this information is precise, but it has not been interpreted.*

```css
/* Cyan Data application */
.hic-log-output,
.hic-data-stream,
.hic-archive-entry {
  color: #00bcd4;
  text-shadow: 0 0 4px rgba(0, 188, 212, 0.3);
  font-family: var(--hic-mono);
  font-size: 0.85rem;
  line-height: 1.4;
}

.hic-data-table {
  border-collapse: collapse;
}

.hic-data-table th,
.hic-data-table td {
  border: 1px solid rgba(0, 188, 212, 0.3);
  padding: 4px 8px;
  font-family: var(--hic-mono);
}

.hic-data-table th {
  background: rgba(0, 188, 212, 0.1);
  color: #26c6da;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.05em;
}
```

#### Hot Pink -- The Emergency Signal

| Property | Value |
|---|---|
| Primary Hex | `#ff1493` |
| Bright Variant | `#ff3ca8` |
| Dim Variant | `#cc0077` |
| Glow Variant | `rgba(255, 20, 147, 0.35)` |
| CSS Variable | `--hic-hot-pink` |
| Semantic Role | Emergency overrides, break-glass procedures, system-level alerts |
| Usage Zones | Penthouse (Emergency Command), break-glass panels throughout |

Hot Pink is the color you should never see. It is reserved for scenarios where
normal protocols have failed and extraordinary measures are required. A
documentation system in total collapse. A security breach requiring immediate
content lockdown. An emergency override that bypasses all review gates to
publish critical information immediately.

```css
/* Hot Pink Emergency application */
.hic-emergency,
.hic-break-glass,
.hic-override-active {
  color: #ff1493;
  text-shadow:
    0 0 10px rgba(255, 20, 147, 0.6),
    0 0 20px rgba(255, 20, 147, 0.3);
  animation: hic-strobe-pink 0.5s ease-in-out infinite;
}

@keyframes hic-strobe-pink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.hic-emergency-panel {
  background: rgba(255, 20, 147, 0.1);
  border: 2px solid #ff1493;
  box-shadow:
    0 0 12px rgba(255, 20, 147, 0.5),
    0 0 24px rgba(255, 20, 147, 0.2),
    inset 0 0 12px rgba(255, 20, 147, 0.1);
  padding: 1.5rem;
}
```

Hot Pink was chosen for maximum physiological impact. It combines the urgency
of red with the unnaturalness of magenta. The human brain does not encounter
this color in nature, and that strangeness triggers heightened attention. In a
building bathed in blues and greens, Hot Pink is an alien invasion -- you
cannot ignore it.

### 2.3 Structural Neutrals

The palette is supported by a set of neutrals that provide structure without
competing with the signal colors.

| Name | Hex | CSS Variable | Usage |
|---|---|---|---|
| Void Black | `#0a0a0f` | `--hic-void` | Primary background, the default state |
| Deep Charcoal | `#12121a` | `--hic-charcoal` | Elevated surfaces, card backgrounds |
| Gunmetal | `#1a1a2e` | `--hic-gunmetal` | Secondary panels, sidebar backgrounds |
| Slate Edge | `#2a2a3e` | `--hic-slate` | Borders, dividers, structural lines |
| Ash Gray | `#4a4a5e` | `--hic-ash` | Disabled text, placeholder content |
| Fog | `#8a8a9e` | `--hic-fog` | Secondary text, timestamps, metadata |
| Bone White | `#c8c8d4` | `--hic-bone` | Primary body text on dark backgrounds |
| Stark White | `#e8e8f0` | `--hic-stark` | High-emphasis text, headings |

```css
:root {
  /* Structural Neutrals */
  --hic-void:     #0a0a0f;
  --hic-charcoal: #12121a;
  --hic-gunmetal: #1a1a2e;
  --hic-slate:    #2a2a3e;
  --hic-ash:      #4a4a5e;
  --hic-fog:      #8a8a9e;
  --hic-bone:     #c8c8d4;
  --hic-stark:    #e8e8f0;

  /* Signal Colors */
  --hic-electric-blue:  #7b8cde;
  --hic-amber-warning:  #c49a3a;
  --hic-crimson-alert:  #c45a5a;
  --hic-emerald-status: #5fa85f;
  --hic-violet-accent:  #9b59b6;
  --hic-cyan-data:      #00bcd4;
  --hic-hot-pink:       #ff1493;
}
```

### 2.4 Zone Color Temperature Map

Each zone of the HIC skyscraper operates under a dominant color temperature
that reflects the nature of the work performed on those floors.

| Floor Range | Zone Name | Dominant Color | Temperature | Rationale |
|---|---|---|---|---|
| B3-B1 | Foundational Infrastructure | Cyan Data | Cool (7500K) | Machine systems, raw data storage |
| 1-10 | Public Atrium | Electric Blue | Neutral-Cool (6000K) | Active knowledge, reader-facing content |
| 11-20 | Verified Archives | Emerald Status | Neutral (5200K) | Stable, trustworthy, reviewed content |
| 21-30 | Commentary Mezzanine | Violet Accent | Warm-Cool (5800K) | Meta-analysis, cross-referencing |
| 31-40 | Review & QA Division | Amber Warning | Warm (4200K) | Content under examination, drafts |
| 41-50 | Security Division | Crimson Alert | Warm-Red (3500K) | Threat assessment, access control |
| 51-60 | Data Vaults | Cyan Data | Cool (8000K) | Deep archives, telemetry, logs |
| 61-70 | Research Spire | Electric Blue + Violet | Mixed (5500K) | Active investigation, experimentation |
| 71-80 | Administrative Tower | Bone White + Slate | Neutral (5000K) | Governance, policy, institutional memory |
| 81-90 | Observation Deck | All colors, muted | Adaptive | Overview dashboards, system-wide status |
| Penthouse | Emergency Command | Hot Pink (dormant) | Emergency (N/A) | Activates only during system crises |

### 2.5 Color Interaction Rules

Colors in the HIC do not exist in isolation. Strict rules govern how they
interact:

**Rule 1: Signal colors never touch.** Two different signal colors must always
be separated by at least 8px of neutral space (Void, Charcoal, or Gunmetal).
Adjacent signal colors create visual noise and semantic confusion.

**Rule 2: Glow radii do not overlap.** If two elements with different-colored
glows are placed near each other, their `box-shadow` and `text-shadow` radii
must be reduced so the glow fields do not blend. Blended glows create
undefined colors that are outside the HIC vocabulary.

**Rule 3: Maximum three signal colors per viewport.** At any given scroll
position, no more than three signal colors should be visible simultaneously.
Exceeding this threshold overwhelms the semantic system. If a layout requires
more, the additional signals must be collapsed behind interaction (tooltips,
expandable panels, tabs).

**Rule 4: Crimson and Hot Pink are mutually exclusive.** Both represent failure
states. If Hot Pink (emergency) is active, all Crimson (alert) indicators are
dimmed to their Dim Variant. The hierarchy of urgency must be unambiguous.

**Rule 5: Emerald is never animated.** Health and approval are stable states.
Animating Emerald signals implies instability, which contradicts its semantic
meaning. Emerald elements are always static.

---

## 3. Neon Signage System

### 3.1 Floor Markers

Every floor of the HIC is identified by a **floor marker** -- a neon sign
element that combines the floor number, zone name, and a status indicator into
a single visual unit.

```
 ______________________________________
|                                      |
|   FLOOR 23  //  COMMENTARY MEZZANINE |
|   ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~  |
|   STATUS: NOMINAL  |  DOCS: 1,847    |
|______________________________________|
```

Floor markers follow this specification:

| Component | Style | Notes |
|---|---|---|
| Floor number | `--hic-stark`, 2rem, bold monospace | Always zero-padded to two digits (FLOOR 07) |
| Zone name | `--hic-fog`, 0.85rem, uppercase monospace | Separated from floor number by `//` delimiter |
| Underline | Zone's dominant signal color, 1px solid | Spans the full width of the marker |
| Status label | Zone's signal color, 0.75rem, uppercase | `NOMINAL`, `REVIEW`, `ALERT`, or `EMERGENCY` |
| Document count | `--hic-ash`, 0.75rem, monospace | Pipe-delimited from status |

```css
.hic-floor-marker {
  background: var(--hic-charcoal);
  border: 1px solid var(--hic-slate);
  padding: 0.75rem 1rem;
  font-family: var(--hic-mono);
  position: relative;
}

.hic-floor-marker::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--zone-color);
  box-shadow: 0 0 8px var(--zone-glow);
}

.hic-floor-marker__number {
  font-size: 2rem;
  font-weight: 700;
  color: var(--hic-stark);
  letter-spacing: 0.05em;
}

.hic-floor-marker__zone {
  font-size: 0.85rem;
  color: var(--hic-fog);
  text-transform: uppercase;
  letter-spacing: 0.1em;
}
```

### 3.2 Wayfinding: Neon Arrows and Path Indicators

Navigation within the HIC is guided by a system of neon **pathlines** --
continuous lines of light that trace routes through the building's corridors,
stairwells, and elevator shafts.

**Pathline Types:**

| Pathline | Color | Width | Pattern | Meaning |
|---|---|---|---|---|
| Primary Route | Electric Blue | 3px | Solid | Main navigation path to destination |
| Alternate Route | Electric Blue (dim) | 2px | Dashed (8px/4px) | Secondary path, longer but available |
| Cross-Reference | Violet Accent | 2px | Dotted (4px/4px) | Link to related content in another zone |
| Restricted Path | Crimson Alert | 2px | Solid | Requires elevated access or credentials |
| Emergency Exit | Hot Pink | 4px | Solid, pulsing | Emergency evacuation or break-glass path |
| Data Pipeline | Cyan Data | 1px | Solid, animated flow | Automated data movement between systems |

**Arrow Glyphs:**

Directional indicators use a standardized arrow system rendered in the
pathline's color:

```
  Standard Arrows (navigation)
  ============================
  UP:        ^       RIGHT:    >       DIAGONAL:  /
             |                 -                  /
             |                 -                 /

  DOWN:      |       LEFT:     -       BRANCH:   --+--
             |                 -                    |
             v                 <                    |

  Elevator Indicators
  ===================
  ASCENDING:   [^]     DESCENDING:  [v]     STOPPED:  [=]
```

```css
.hic-pathline {
  stroke-width: var(--pathline-width, 3px);
  stroke: var(--pathline-color, var(--hic-electric-blue));
  fill: none;
  filter: drop-shadow(0 0 4px var(--pathline-glow));
}

.hic-pathline--alternate {
  stroke-dasharray: 8 4;
  opacity: 0.6;
}

.hic-pathline--data {
  stroke-dasharray: 2 2;
  animation: hic-flow 1s linear infinite;
}

@keyframes hic-flow {
  to { stroke-dashoffset: -4; }
}

.hic-arrow {
  width: 24px;
  height: 24px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-family: var(--hic-mono);
  font-size: 1.2rem;
  color: var(--pathline-color);
  text-shadow: 0 0 6px var(--pathline-glow);
}
```

### 3.3 Status Boards

Status boards are live dashboards mounted at zone entrances and key
intersections throughout the HIC. They provide at-a-glance documentation
health metrics for their zone.

**Standard Status Board Layout:**

```
+================================================================+
|  ZONE: VERIFIED ARCHIVES  //  FLOORS 11-20                     |
|================================================================|
|                                                                |
|  TOTAL DOCUMENTS    2,341    [========XXXXXXXX========]  87%   |
|  LAST REVIEWED      2026-02-14 14:32 UTC                       |
|  PENDING REVIEWS    12       [==XX                    ]   4%   |
|  BROKEN LINKS       3        [X                       ]   1%   |
|  FRESHNESS INDEX    0.94     [=======XXXXXXXXXX=======]  94%   |
|                                                                |
|  RECENT ACTIVITY                                               |
|  ~~~~~~~~~~~~~~                                                |
|  14:32  DOC-2341 reviewed by k.vasquez     [APPROVED]          |
|  14:18  DOC-1899 updated by m.chen         [MODIFIED]          |
|  13:55  DOC-0744 flagged by system         [STALE]             |
|  13:41  DOC-2102 published by a.okafor     [NEW]               |
|                                                                |
+================================================================+
```

**Progress Bar Color Logic:**

| Percentage Range | Bar Color | Behavior |
|---|---|---|
| 90-100% | Emerald Status | Solid, no animation |
| 70-89% | Electric Blue | Solid, no animation |
| 50-69% | Amber Warning | Solid, subtle pulse |
| 25-49% | Amber Warning (bright) | Solid, moderate pulse |
| 0-24% | Crimson Alert | Solid, rapid pulse |

```css
.hic-status-board {
  background: var(--hic-charcoal);
  border: 1px solid var(--hic-slate);
  font-family: var(--hic-mono);
  padding: 1rem;
  position: relative;
  overflow: hidden;
}

.hic-status-board::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--zone-color);
  box-shadow: 0 0 12px var(--zone-glow);
}

.hic-progress-bar {
  height: 8px;
  background: var(--hic-gunmetal);
  border-radius: 0;  /* Digital Brutalism: no rounding */
  overflow: hidden;
}

.hic-progress-bar__fill {
  height: 100%;
  background: var(--bar-color);
  box-shadow: 0 0 6px var(--bar-glow);
  transition: width 0.5s ease-out;
}
```

### 3.4 Warning Lights

Warning lights are persistent visual indicators attached to individual
documents, sections, or entire floors. They override the ambient zone lighting
to communicate that something requires attention.

**Warning Light States:**

| State | Visual | Trigger Condition |
|---|---|---|
| `NOMINAL` | No warning light visible | All checks pass, content is current |
| `AGING` | Amber, slow pulse (4s cycle) | Document has not been reviewed in >90 days |
| `STALE` | Amber, moderate pulse (2s cycle) | Document has not been reviewed in >180 days |
| `BROKEN` | Crimson, steady glow | Contains broken links or missing references |
| `DEPRECATED` | Crimson, slow pulse (3s cycle) | Marked for removal but still referenced |
| `SECURITY` | Crimson, rapid pulse (1s cycle) | Security-related content requiring immediate review |
| `EMERGENCY` | Hot Pink, strobe (0.5s cycle) | System-level crisis affecting this content |

```css
.hic-warning-light {
  width: 8px;
  height: 8px;
  border-radius: 50%;  /* Exception to no-rounding: indicator dots */
  display: inline-block;
  margin-right: 0.5rem;
}

.hic-warning-light--nominal {
  background: var(--hic-emerald-status);
  box-shadow: 0 0 4px rgba(95, 168, 95, 0.4);
}

.hic-warning-light--aging {
  background: var(--hic-amber-warning);
  box-shadow: 0 0 6px rgba(196, 154, 58, 0.4);
  animation: hic-pulse 4s ease-in-out infinite;
}

.hic-warning-light--stale {
  background: var(--hic-amber-warning);
  box-shadow: 0 0 8px rgba(196, 154, 58, 0.5);
  animation: hic-pulse 2s ease-in-out infinite;
}

.hic-warning-light--broken {
  background: var(--hic-crimson-alert);
  box-shadow: 0 0 8px rgba(196, 90, 90, 0.5);
}

.hic-warning-light--deprecated {
  background: var(--hic-crimson-alert);
  box-shadow: 0 0 8px rgba(196, 90, 90, 0.5);
  animation: hic-pulse 3s ease-in-out infinite;
}

.hic-warning-light--security {
  background: var(--hic-crimson-alert);
  box-shadow: 0 0 10px rgba(196, 90, 90, 0.6);
  animation: hic-pulse 1s ease-in-out infinite;
}

.hic-warning-light--emergency {
  background: var(--hic-hot-pink);
  box-shadow: 0 0 14px rgba(255, 20, 147, 0.7);
  animation: hic-strobe-pink 0.5s ease-in-out infinite;
}

@keyframes hic-pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}
```

---

## 4. Typography Hierarchy

### 4.1 Typeface Selection

The HIC uses a deliberately constrained type system. Three typefaces. No more.

| Role | Typeface | Fallback Stack | Rationale |
|---|---|---|---|
| Primary (Mono) | **JetBrains Mono** | `'Fira Code', 'Source Code Pro', 'Courier New', monospace` | Ligature support, excellent glyph clarity at small sizes, designed for extended code reading |
| Display (Mono) | **Share Tech Mono** | `'VT323', 'Press Start 2P', monospace` | Industrial aesthetic for headings and signage, evokes CRT terminal displays |
| Fallback (Sans) | **Inter** | `'IBM Plex Sans', 'Helvetica Neue', sans-serif` | Clean proportional type for when monospace would harm readability (long-form prose) |

```css
:root {
  --hic-mono:    'JetBrains Mono', 'Fira Code', 'Source Code Pro', 'Courier New', monospace;
  --hic-display: 'Share Tech Mono', 'VT323', monospace;
  --hic-sans:    'Inter', 'IBM Plex Sans', 'Helvetica Neue', sans-serif;
}
```

**Why monospace-first?** Documentation is structured information. It has
headers, lists, tables, code blocks, hierarchies. These structures align
naturally on a fixed-width grid. When every character is the same width,
indentation is unambiguous, tables align without fudging, and the reader's eye
can track vertical columns of information. Proportional type is optimized for
flowing prose. The HIC is not prose -- it is *architecture*.

### 4.2 Heading Hierarchy

Headings use the Display typeface (Share Tech Mono) and follow a strict size
and color progression:

| Level | Size | Weight | Color | Tracking | Text Transform | Glow |
|---|---|---|---|---|---|---|
| H1 | 2.4rem | 400 | `--hic-stark` | 0.15em | UPPERCASE | 8px zone color |
| H2 | 1.8rem | 400 | `--hic-stark` | 0.10em | UPPERCASE | 4px zone color |
| H3 | 1.4rem | 400 | `--hic-bone` | 0.08em | UPPERCASE | 2px zone color |
| H4 | 1.1rem | 700 | `--hic-bone` | 0.05em | Title Case | None |
| H5 | 0.95rem | 700 | `--hic-fog` | 0.05em | Title Case | None |
| H6 | 0.85rem | 400 | `--hic-fog` | 0.10em | UPPERCASE | None |

```css
h1, h2, h3 {
  font-family: var(--hic-display);
  text-transform: uppercase;
  color: var(--hic-stark);
  margin-top: 3rem;
  margin-bottom: 1rem;
  position: relative;
}

h1 {
  font-size: 2.4rem;
  letter-spacing: 0.15em;
  text-shadow: 0 0 8px var(--zone-glow);
  padding-bottom: 0.5rem;
  border-bottom: 2px solid var(--zone-color);
}

h1::before {
  content: '// ';
  color: var(--hic-ash);
  font-weight: 400;
}

h2 {
  font-size: 1.8rem;
  letter-spacing: 0.10em;
  text-shadow: 0 0 4px var(--zone-glow);
}

h2::before {
  content: '## ';
  color: var(--hic-ash);
  font-weight: 400;
}

h3 {
  font-size: 1.4rem;
  letter-spacing: 0.08em;
  color: var(--hic-bone);
  text-shadow: 0 0 2px var(--zone-glow);
}

h4, h5, h6 {
  font-family: var(--hic-mono);
}

h4 {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--hic-bone);
  letter-spacing: 0.05em;
}

h5 {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--hic-fog);
  letter-spacing: 0.05em;
}

h6 {
  font-size: 0.85rem;
  font-weight: 400;
  color: var(--hic-fog);
  letter-spacing: 0.10em;
  text-transform: uppercase;
}
```

The `::before` pseudo-elements on H1 and H2 display syntax markers (`//` and
`##`) in muted gray. This reinforces the source-code nature of the document
and provides a visual breadcrumb that helps readers quickly identify heading
levels while scanning.

### 4.3 Body Text

Body text is the workhorse of the HIC. Most reading time is spent here, so
legibility in dark environments is paramount.

```css
body {
  font-family: var(--hic-mono);
  font-size: 0.9rem;
  line-height: 1.7;
  color: var(--hic-bone);
  background: var(--hic-void);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

p {
  margin-bottom: 1rem;
  max-width: 78ch;  /* Monospace optimal line length */
}
```

**Line length is capped at 78 characters.** This matches the traditional
terminal width (80 columns minus 2 for margin) and is the optimal line length
for monospace reading. The eye can track from the end of one line to the
beginning of the next without losing its place.

**Line height is 1.7.** This is generous for body text but necessary in a dark
environment. Dark backgrounds make lines of text appear to compress vertically.
The extra leading counteracts this optical illusion and prevents the "wall of
text" effect that causes reader fatigue.

**Font size is 0.9rem.** Monospace characters are wider than proportional
characters at the same size, so the effective reading size is equivalent to
approximately 1rem of proportional text. The slight reduction prevents lines
from feeling bloated.

### 4.4 Code Blocks: The Terminal Green Aesthetic

Code blocks in the HIC are not merely styled containers. They are *terminal
windows* -- portholes into the building's machine infrastructure. Their visual
treatment evokes CRT phosphor green, the ancestral glow of all digital text.

```css
pre, code {
  font-family: var(--hic-mono);
  font-size: 0.85rem;
  line-height: 1.5;
}

code {
  color: #66ff66;  /* Terminal green for inline code */
  background: rgba(102, 255, 102, 0.05);
  padding: 0.15em 0.4em;
  border: 1px solid rgba(102, 255, 102, 0.15);
}

pre {
  background: #050510;  /* Deeper than Void -- the machine layer */
  border: 1px solid var(--hic-slate);
  border-left: 3px solid #66ff66;
  padding: 1rem 1.25rem;
  overflow-x: auto;
  position: relative;
}

pre::before {
  content: attr(data-lang) ' //';
  position: absolute;
  top: 0;
  right: 0;
  padding: 0.25rem 0.75rem;
  font-size: 0.7rem;
  color: var(--hic-ash);
  background: var(--hic-gunmetal);
  border-left: 1px solid var(--hic-slate);
  border-bottom: 1px solid var(--hic-slate);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

pre code {
  color: #66ff66;
  background: none;
  border: none;
  padding: 0;
  text-shadow: 0 0 2px rgba(102, 255, 102, 0.2);
}
```

**Syntax highlighting palette** (for code within terminal blocks):

| Token Type | Color | Hex |
|---|---|---|
| Keywords | Electric Blue | `#7b8cde` |
| Strings | Amber Warning | `#c49a3a` |
| Numbers | Cyan Data | `#00bcd4` |
| Comments | Ash Gray | `#4a4a5e` |
| Functions | Violet Accent | `#9b59b6` |
| Operators | Bone White | `#c8c8d4` |
| Types/Classes | Emerald Status | `#5fa85f` |
| Errors | Crimson Alert | `#c45a5a` |
| Default | Terminal Green | `#66ff66` |

The syntax highlighting palette reuses the HIC signal colors, maintaining
semantic consistency. Keywords are blue because they are active knowledge
(navigation points in code). Strings are amber because they are literal data
that should be verified. Errors are crimson because they are critical.

### 4.5 Data Tables

Tables in the HIC are rendered as **grid-lined data displays** with subtle
neon edge effects. They combine the density of a spreadsheet with the clarity
of the neon signage system.

```css
table {
  width: 100%;
  border-collapse: collapse;
  font-family: var(--hic-mono);
  font-size: 0.85rem;
  margin: 1.5rem 0;
}

thead {
  border-bottom: 2px solid var(--zone-color);
}

th {
  padding: 0.5rem 0.75rem;
  text-align: left;
  color: var(--zone-color);
  font-weight: 400;
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.08em;
  background: rgba(var(--zone-color-rgb), 0.05);
}

td {
  padding: 0.5rem 0.75rem;
  border-bottom: 1px solid var(--hic-slate);
  color: var(--hic-bone);
}

tr:hover td {
  background: rgba(var(--zone-color-rgb), 0.03);
}

/* Grid glow effect on table borders */
table {
  box-shadow:
    0 0 1px var(--zone-color),
    0 0 4px rgba(var(--zone-color-rgb), 0.1);
}
```

Tables in HIC documents never use rounded corners, padding-heavy cells, or
alternating row colors. The grid lines are the structure. The data is the
content. The zone color on the header row tells you what domain this table
belongs to.

### 4.6 Emphasis and Inline Formatting

| Format | Style | Use Case |
|---|---|---|
| **Bold** | `--hic-stark`, font-weight 700 | Key terms, critical identifiers |
| *Italic* | `--hic-fog`, italic, slight left-indent | Definitions, editorial voice, tangential notes |
| `Code` | Terminal green, border, background | Technical identifiers, commands, file paths |
| ~~Strikethrough~~ | `--hic-ash`, line-through, 0.8 opacity | Superseded information (kept for audit trail) |
| [Link](#) | Electric Blue, underline on hover | Internal and external navigation |
| > Blockquote | Left border in zone color, `--hic-fog` text | Quoted material, callouts, official doctrine |

```css
strong {
  color: var(--hic-stark);
  font-weight: 700;
}

em {
  color: var(--hic-fog);
  font-style: italic;
}

del {
  color: var(--hic-ash);
  text-decoration: line-through;
  opacity: 0.8;
}

a {
  color: var(--hic-electric-blue);
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: border-color 0.2s ease, text-shadow 0.2s ease;
}

a:hover {
  border-bottom-color: var(--hic-electric-blue);
  text-shadow: 0 0 8px rgba(123, 140, 222, 0.4);
}

blockquote {
  border-left: 3px solid var(--zone-color);
  padding: 0.75rem 1rem;
  margin: 1rem 0;
  background: rgba(var(--zone-color-rgb), 0.03);
  color: var(--hic-fog);
  font-style: italic;
}

blockquote p:last-child {
  margin-bottom: 0;
}
```

---

## 5. Holographic Display Standards

### 5.1 Information Layering (Z-Depth System)

The HIC employs a z-depth system that maps content importance to visual
proximity. Content that is more important appears *closer* to the reader --
projected further from the surface of the display, casting stronger shadows,
appearing more vivid.

| Z-Layer | Depth (CSS) | Opacity | Glow Intensity | Content Type |
|---|---|---|---|---|
| Z-0 (Surface) | `z-index: 0` | 0.40 | None | Background textures, ambient patterns |
| Z-1 (Recessed) | `z-index: 10` | 0.60 | Minimal (2px) | Archived content, historical records |
| Z-2 (Flush) | `z-index: 20` | 0.85 | Low (4px) | Standard body content, default state |
| Z-3 (Raised) | `z-index: 30` | 0.95 | Medium (6px) | Active content, currently being read/edited |
| Z-4 (Projected) | `z-index: 40` | 1.00 | High (10px) | Critical alerts, modal dialogs, urgent notices |
| Z-5 (Floating) | `z-index: 50` | 1.00 | Maximum (16px) | Emergency overrides, system-level interrupts |

```css
.hic-depth-0 {
  z-index: 0;
  opacity: 0.4;
  filter: blur(0.5px);
}

.hic-depth-1 {
  z-index: 10;
  opacity: 0.6;
  transform: translateZ(0);
}

.hic-depth-2 {
  z-index: 20;
  opacity: 0.85;
}

.hic-depth-3 {
  z-index: 30;
  opacity: 0.95;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.4);
}

.hic-depth-4 {
  z-index: 40;
  opacity: 1;
  box-shadow:
    0 4px 16px rgba(0, 0, 0, 0.5),
    0 0 10px var(--zone-glow);
}

.hic-depth-5 {
  z-index: 50;
  opacity: 1;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.6),
    0 0 16px var(--zone-glow),
    0 0 32px rgba(var(--zone-color-rgb), 0.15);
}
```

### 5.2 Transparency Levels

Transparency in the HIC is not aesthetic -- it communicates **content
maturity**. The more transparent an element, the less finalized its content.

| Transparency Level | Opacity | Application |
|---|---|---|
| Solid | 1.0 | Published, reviewed, canonical content |
| Near-Solid | 0.85 | Published content, not yet reviewed at current revision |
| Translucent | 0.65 | Draft content, visible but clearly marked as unfinished |
| Ghost | 0.40 | Placeholder content, structure exists but text is pending |
| Phantom | 0.20 | Deleted/deprecated content shown for reference only |

```css
.hic-content--published  { opacity: 1.0; }
.hic-content--unreviewed { opacity: 0.85; }
.hic-content--draft      { opacity: 0.65; border-left: 3px dashed var(--hic-amber-warning); }
.hic-content--placeholder{ opacity: 0.40; border-left: 3px dotted var(--hic-ash); }
.hic-content--deprecated { opacity: 0.20; text-decoration: line-through; pointer-events: none; }
```

### 5.3 Animation Guidelines

Animation in the HIC is functional. Every motion communicates a system state.
There are no decorative transitions.

**Permitted Animations:**

| Animation | CSS | Trigger | Duration | Meaning |
|---|---|---|---|---|
| Pulse | `opacity: 1 -> 0.4 -> 1` | Content is live / being updated | 2-4s | Data is fresh, system is active |
| Slow Fade | `opacity: 1 -> 0.6` | Content aging | 10s+ | Content freshness is decreasing |
| Strobe | `opacity: 1 -> 0.2 -> 1` | Emergency | 0.3-0.5s | Immediate attention required |
| Glow Breathe | `box-shadow intensity +/- 30%` | Active selection | 3s | Element is currently focused/selected |
| Slide In | `translateX(-20px) -> 0` | New content appearing | 0.3s | Fresh content arriving in the system |
| Scan Line | `background-position scroll` | Data processing | 2s | System is processing or indexing |

**Prohibited Animations:**

- Bounce effects
- Elastic/spring physics
- Parallax scrolling
- Loading spinners (use scan-line effect instead)
- Page transitions with fades, slides, or zooms
- Hover effects that change layout (grow, shrink, rotate)
- Any animation exceeding 4 seconds unless it represents an ongoing state

```css
/* Permitted: Pulse for live data */
@keyframes hic-pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

/* Permitted: Glow breathe for active selection */
@keyframes hic-glow-breathe {
  0%, 100% { box-shadow: 0 0 6px var(--zone-glow); }
  50% { box-shadow: 0 0 12px var(--zone-glow); }
}

/* Permitted: Scan line for processing */
@keyframes hic-scan {
  0% { background-position: 0 -100%; }
  100% { background-position: 0 100%; }
}

.hic-processing {
  position: relative;
  overflow: hidden;
}

.hic-processing::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 200%;
  background: linear-gradient(
    180deg,
    transparent 0%,
    rgba(var(--zone-color-rgb), 0.05) 45%,
    rgba(var(--zone-color-rgb), 0.15) 50%,
    rgba(var(--zone-color-rgb), 0.05) 55%,
    transparent 100%
  );
  animation: hic-scan 2s linear infinite;
  pointer-events: none;
}
```

### 5.4 Projection Distances

In the physical HIC building, holographic displays project content at varying
distances from the wall surface. In the digital representation, this maps to
shadow depth, border intensity, and visual weight.

| Projection Class | Physical Distance | Digital Representation |
|---|---|---|
| Flush | 0cm (surface-mounted) | No shadow, 1px border, standard colors |
| Near | 2-5cm from surface | 2px shadow offset, subtle glow |
| Mid | 5-15cm from surface | 4-8px shadow offset, moderate glow, slight scale(1.01) |
| Far | 15-30cm from surface | 12-16px shadow offset, strong glow, scale(1.02) |
| Detached | 30cm+ from surface | 20px+ shadow offset, intense glow, scale(1.03), backdrop-filter: blur |

```css
.hic-projection--flush {
  border: 1px solid var(--hic-slate);
  box-shadow: none;
}

.hic-projection--near {
  border: 1px solid var(--zone-color);
  box-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
}

.hic-projection--mid {
  border: 1px solid var(--zone-color);
  box-shadow:
    4px 4px 8px rgba(0, 0, 0, 0.5),
    0 0 6px var(--zone-glow);
  transform: scale(1.01);
}

.hic-projection--far {
  border: 1px solid var(--zone-color);
  box-shadow:
    8px 8px 16px rgba(0, 0, 0, 0.6),
    0 0 12px var(--zone-glow);
  transform: scale(1.02);
}

.hic-projection--detached {
  border: 2px solid var(--zone-color);
  box-shadow:
    12px 12px 24px rgba(0, 0, 0, 0.7),
    0 0 20px var(--zone-glow),
    0 0 40px rgba(var(--zone-color-rgb), 0.1);
  transform: scale(1.03);
  backdrop-filter: blur(4px);
}
```

---

## 6. Environmental Lighting

### 6.1 Ambient Glow: Content Completeness as Luminance

The exterior of the HIC skyscraper is not uniformly lit. Each floor emits
light in direct proportion to the completeness and health of its documentation
corpus. A fully documented floor blazes with its zone color. A floor with
gaps and outdated content is dim, its neon sputtering.

**Luminance Formula:**

```
FLOOR_LUMINANCE = (DOC_COMPLETENESS * 0.4)
               + (DOC_FRESHNESS * 0.3)
               + (LINK_HEALTH * 0.2)
               + (REVIEW_STATUS * 0.1)
```

Each factor is scored 0.0 to 1.0:

| Factor | Score 1.0 | Score 0.5 | Score 0.0 |
|---|---|---|---|
| DOC_COMPLETENESS | All planned documents exist | Half are written | Section is empty |
| DOC_FRESHNESS | All docs reviewed within 90 days | Average 180+ days since review | No reviews in 365+ days |
| LINK_HEALTH | Zero broken links | < 5% broken links | > 20% broken links |
| REVIEW_STATUS | All docs approved | Mixed approved/pending | All pending or rejected |

**Luminance to CSS Mapping:**

| Luminance Score | Visual State | CSS Treatment |
|---|---|---|
| 0.9 - 1.0 | Blazing | Full color, 12px glow, slight bloom effect |
| 0.7 - 0.89 | Bright | Full color, 6px glow |
| 0.5 - 0.69 | Normal | Slightly desaturated, 3px glow |
| 0.3 - 0.49 | Dim | 50% desaturated, 1px glow, reduced opacity |
| 0.0 - 0.29 | Dark | 80% desaturated, no glow, 0.6 opacity, flicker |

```css
.hic-floor-exterior[data-luminance="blazing"] {
  --floor-saturation: 100%;
  --floor-glow-size: 12px;
  --floor-opacity: 1;
  filter: brightness(1.1) saturate(1.2);
}

.hic-floor-exterior[data-luminance="bright"] {
  --floor-saturation: 100%;
  --floor-glow-size: 6px;
  --floor-opacity: 1;
}

.hic-floor-exterior[data-luminance="normal"] {
  --floor-saturation: 85%;
  --floor-glow-size: 3px;
  --floor-opacity: 0.9;
}

.hic-floor-exterior[data-luminance="dim"] {
  --floor-saturation: 50%;
  --floor-glow-size: 1px;
  --floor-opacity: 0.7;
  filter: saturate(0.5);
}

.hic-floor-exterior[data-luminance="dark"] {
  --floor-saturation: 20%;
  --floor-glow-size: 0;
  --floor-opacity: 0.6;
  filter: saturate(0.2) brightness(0.7);
  animation: hic-flicker 4s ease-in-out infinite;
}

@keyframes hic-flicker {
  0%, 92%, 96%, 100% { opacity: 0.6; }
  93% { opacity: 0.3; }
  95% { opacity: 0.5; }
  97% { opacity: 0.2; }
}
```

### 6.2 Shadow Mapping: Darkness as Absence

In the HIC, shadows are not merely the absence of light. They are information.
A shadow on the building's surface tells an observer exactly where
documentation is missing.

**Shadow Categories:**

| Shadow Type | Appearance | Meaning |
|---|---|---|
| Structural Shadow | Hard edge, consistent depth | Normal architectural shadow, no semantic meaning |
| Coverage Gap | Soft edge, irregular shape | Missing documentation in a category |
| Broken Path | Sharp, fragmented shadow | Broken links severing a documentation pathway |
| Decay Shadow | Gradually deepening darkness | Content aging without review, slowly fading |
| Void Pocket | Complete blackness, no ambient light | Entire section missing, no content exists |

```css
/* Coverage Gap: soft shadow indicating missing docs */
.hic-shadow--coverage-gap {
  box-shadow: inset 0 0 30px rgba(0, 0, 0, 0.6);
  position: relative;
}

.hic-shadow--coverage-gap::after {
  content: 'COVERAGE GAP DETECTED';
  position: absolute;
  bottom: 0.5rem;
  right: 0.5rem;
  font-size: 0.65rem;
  color: var(--hic-amber-warning);
  font-family: var(--hic-mono);
  text-transform: uppercase;
  letter-spacing: 0.1em;
  opacity: 0.6;
}

/* Void Pocket: total absence */
.hic-shadow--void {
  background: var(--hic-void);
  border: 1px dashed var(--hic-ash);
  min-height: 4rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.hic-shadow--void::after {
  content: '[ NO CONTENT ]';
  color: var(--hic-ash);
  font-family: var(--hic-mono);
  font-size: 0.75rem;
  letter-spacing: 0.15em;
}
```

### 6.3 Time-of-Day Cycles: Document Freshness

The HIC operates on a visual time-of-day cycle that represents document
freshness. This is not tied to actual clock time -- it is tied to the age of
the content.

| Cycle Phase | Time Metaphor | Freshness | Visual Treatment |
|---|---|---|---|
| Dawn | 0-7 days since update | Brand new | Warm white tint, high contrast, bright neon |
| Midday | 7-30 days | Current | Standard palette, full saturation |
| Afternoon | 30-90 days | Aging | Slightly warm shift, saturation -10% |
| Dusk | 90-180 days | Stale | Amber shift, saturation -25%, reduced glow |
| Night | 180-365 days | Old | Cool blue-gray shift, saturation -50%, dim |
| Deep Night | 365+ days | Ancient | Near-monochrome, minimal glow, flicker risk |

```css
/* Dawn: freshly published content */
.hic-freshness--dawn {
  --freshness-tint: rgba(255, 248, 240, 0.03);
  --freshness-saturation: 110%;
  --freshness-brightness: 105%;
}

/* Midday: current content, baseline */
.hic-freshness--midday {
  --freshness-tint: transparent;
  --freshness-saturation: 100%;
  --freshness-brightness: 100%;
}

/* Afternoon: starting to age */
.hic-freshness--afternoon {
  --freshness-tint: rgba(196, 154, 58, 0.02);
  --freshness-saturation: 90%;
  --freshness-brightness: 98%;
}

/* Dusk: needs review soon */
.hic-freshness--dusk {
  --freshness-tint: rgba(196, 154, 58, 0.05);
  --freshness-saturation: 75%;
  --freshness-brightness: 90%;
}

/* Night: overdue for review */
.hic-freshness--night {
  --freshness-tint: rgba(90, 90, 120, 0.05);
  --freshness-saturation: 50%;
  --freshness-brightness: 80%;
}

/* Deep Night: dangerously stale */
.hic-freshness--deep-night {
  --freshness-tint: rgba(30, 30, 50, 0.08);
  --freshness-saturation: 30%;
  --freshness-brightness: 65%;
  animation: hic-flicker 6s ease-in-out infinite;
}
```

The time-of-day cycle creates an intuitive visual language: content that
*looks* fresh *is* fresh. Content that looks like it is fading into the night
is, in fact, fading out of relevance. The reader does not need to check a
"last updated" timestamp. The lighting tells them.

### 6.4 Emergency Lighting Protocols

When the HIC enters an emergency state (triggered by critical documentation
failures, security breaches, or system-wide outages), the environmental
lighting shifts to **Emergency Protocol**.

**Emergency Protocol Stages:**

| Stage | Trigger | Lighting Change |
|---|---|---|
| Stage 0 (Normal) | No active emergencies | Standard zone lighting |
| Stage 1 (Elevated) | Single critical alert | Affected floor shifts to Crimson ambient |
| Stage 2 (Building Alert) | Multiple critical alerts across zones | All floors dim to 40%, affected floors in Crimson |
| Stage 3 (Lockdown) | System-wide failure or security breach | All non-essential lighting off. Hot Pink strobes on affected areas. Emergency pathlines activate. |

```css
/* Stage 1: Single floor alert */
.hic-emergency-stage-1 .hic-floor--affected {
  --zone-color: var(--hic-crimson-alert);
  --zone-glow: rgba(196, 90, 90, 0.4);
}

/* Stage 2: Building-wide alert */
.hic-emergency-stage-2 .hic-floor {
  filter: brightness(0.4);
}

.hic-emergency-stage-2 .hic-floor--affected {
  filter: brightness(1);
  --zone-color: var(--hic-crimson-alert);
  --zone-glow: rgba(196, 90, 90, 0.5);
  animation: hic-pulse-crimson 2s ease-in-out infinite;
}

/* Stage 3: Lockdown */
.hic-emergency-stage-3 {
  --hic-void: #000000;  /* True black */
}

.hic-emergency-stage-3 .hic-floor {
  filter: brightness(0);
  transition: filter 0.5s ease;
}

.hic-emergency-stage-3 .hic-floor--affected {
  filter: brightness(1);
  --zone-color: var(--hic-hot-pink);
  --zone-glow: rgba(255, 20, 147, 0.6);
  animation: hic-strobe-pink 0.5s ease-in-out infinite;
}

.hic-emergency-stage-3 .hic-pathline--emergency {
  display: block;
  stroke: var(--hic-hot-pink);
  stroke-width: 4px;
  filter: drop-shadow(0 0 8px rgba(255, 20, 147, 0.6));
  animation: hic-flow 0.5s linear infinite;
}
```

**Fallback Styles (CSS Degradation):**

When the rendering environment does not support CSS animations, filters, or
custom properties, the emergency system degrades as follows:

```css
/* Fallback: no animation support */
@media (prefers-reduced-motion: reduce) {
  .hic-emergency-stage-3 .hic-floor--affected {
    border: 4px solid #ff1493;
    background: rgba(255, 20, 147, 0.15);
    animation: none;
  }

  .hic-warning-light {
    animation: none !important;
  }
}

/* Fallback: no custom property support */
.hic-emergency-fallback {
  background: #1a0010;
  border: 3px solid #ff1493;
  color: #ff1493;
}
```

---

## 7. Icon & Symbol Library

### 7.1 Design Principles for Icons

All icons in the HIC are constructed on a **24x24 unit grid** using a stroke
weight of 1.5px. They are monochrome by default, inheriting the current zone
color. Icons are never filled; they are always outlined. This is consistent
with the neon-line aesthetic of the building -- neon tubes are lines, not
shapes.

```css
.hic-icon {
  width: 24px;
  height: 24px;
  display: inline-block;
  vertical-align: middle;
  stroke: currentColor;
  stroke-width: 1.5;
  fill: none;
  stroke-linecap: square;  /* Sharp ends, not rounded */
  stroke-linejoin: miter;  /* Sharp corners, not rounded */
  filter: drop-shadow(0 0 2px currentColor);
}

.hic-icon--sm { width: 16px; height: 16px; stroke-width: 1.5; }
.hic-icon--md { width: 24px; height: 24px; stroke-width: 1.5; }
.hic-icon--lg { width: 32px; height: 32px; stroke-width: 2; }
.hic-icon--xl { width: 48px; height: 48px; stroke-width: 2; }
```

### 7.2 Domain Icons

Domain icons identify the category of documentation or the function of a
building zone.

| Icon Name | Glyph Description | Usage |
|---|---|---|
| `shield` | Shield outline, no fill | Security domain, access-controlled content |
| `book-open` | Open book, pages visible | Archive domain, reference documentation |
| `gear` | Six-spoke gear outline | Automation, CI/CD, infrastructure docs |
| `circuit` | Circuit board trace pattern | Technical/engineering documentation |
| `flask` | Erlenmeyer flask, bubbles optional | Research, experimental, lab content |
| `terminal` | Rectangle with `>_` prompt | CLI tools, command references, shell docs |
| `globe` | Circle with latitude/longitude lines | Public-facing, external documentation |
| `lock-closed` | Padlock, shackle up, closed | Restricted content, requires authentication |
| `lock-open` | Padlock, shackle up, open | Unlocked/declassified content |
| `database` | Three stacked cylinders | Data storage, schemas, database docs |
| `antenna` | Vertical line with radiating waves | API endpoints, webhooks, integrations |
| `blueprint` | Folded paper with gridlines | Architecture documents, system design |
| `megaphone` | Cone shape with sound waves | Announcements, release notes, changelogs |
| `compass` | Circle with cardinal indicator | Navigation, getting-started guides, onboarding |

```
ICON CONSTRUCTION GRID (24x24)
================================

  Shield:               Book-Open:            Gear:
  +---------+           +---------+           +---------+
  |   /=\   |           |  __|__  |           |    |    |
  |  / | \  |           | /  |  \ |           |  --O--  |
  |  | | |  |           | |  |  | |           |  / | \  |
  |  \ | /  |           | |  |  | |           | /  |  \ |
  |   \|/   |           | \__|__/ |           |    |    |
  +---------+           +---------+           +---------+

  Terminal:             Lock:                 Database:
  +---------+           +---------+           +---------+
  | +-----+ |           |   ___   |           |  /===\  |
  | |>_   | |           |  |   |  |           |  |---|  |
  | |     | |           |  [===]  |           |  |---|  |
  | |     | |           |  |   |  |           |  |---|  |
  | +-----+ |           |  [___]  |           |  \===/  |
  +---------+           +---------+           +---------+
```

### 7.3 Status Symbols

Status symbols communicate the state of a document, system, or process at a
glance. They are always accompanied by the appropriate signal color.

| Symbol | Color | Meaning |
|---|---|---|
| `check` (checkmark) | Emerald Status | Approved, verified, passing |
| `check-double` (double checkmark) | Emerald Status | Reviewed and approved by multiple parties |
| `x-mark` (X) | Crimson Alert | Failed, rejected, broken |
| `warning-triangle` (triangle with !) | Amber Warning | Caution, review needed, potential issue |
| `clock` (clock face) | Amber Warning | Pending, in-progress, time-sensitive |
| `clock-expired` (clock with X) | Crimson Alert | Overdue, past deadline, expired |
| `lock` (padlock) | Crimson Alert | Locked, access restricted, frozen |
| `unlock` (open padlock) | Emerald Status | Unlocked, accessible, editable |
| `eye` (eye symbol) | Electric Blue | Under observation, being monitored |
| `eye-off` (eye with slash) | Ash Gray | Hidden, not visible to public, draft |
| `refresh` (circular arrows) | Cyan Data | Syncing, updating, processing |
| `archive` (box with arrow down) | Cyan Data | Archived, moved to cold storage |
| `pin` (pushpin) | Violet Accent | Pinned, highlighted, bookmarked |
| `link` (chain links) | Electric Blue | Connected, linked, referenced |
| `link-broken` (broken chain) | Crimson Alert | Broken link, severed reference |

```css
/* Status symbol color assignments */
.hic-status-icon--approved    { color: var(--hic-emerald-status); }
.hic-status-icon--failed      { color: var(--hic-crimson-alert); }
.hic-status-icon--warning     { color: var(--hic-amber-warning); }
.hic-status-icon--pending     { color: var(--hic-amber-warning); }
.hic-status-icon--overdue     { color: var(--hic-crimson-alert); }
.hic-status-icon--locked      { color: var(--hic-crimson-alert); }
.hic-status-icon--unlocked    { color: var(--hic-emerald-status); }
.hic-status-icon--monitoring  { color: var(--hic-electric-blue); }
.hic-status-icon--hidden      { color: var(--hic-ash); }
.hic-status-icon--syncing     { color: var(--hic-cyan-data); }
.hic-status-icon--archived    { color: var(--hic-cyan-data); }
.hic-status-icon--pinned      { color: var(--hic-violet-accent); }
.hic-status-icon--linked      { color: var(--hic-electric-blue); }
.hic-status-icon--link-broken { color: var(--hic-crimson-alert); }
```

### 7.4 Navigation Glyphs

Navigation glyphs guide readers through the HIC's spatial documentation
structure. They are always rendered in Electric Blue unless they indicate a
restricted or emergency path.

| Glyph | Representation | Usage |
|---|---|---|
| `arrow-up` | Upward chevron | Navigate to parent section / higher floor |
| `arrow-down` | Downward chevron | Navigate to child section / lower floor |
| `arrow-left` | Left chevron | Navigate back / previous document |
| `arrow-right` | Right chevron | Navigate forward / next document |
| `elevator-up` | Enclosed upward arrow `[^]` | Jump to distant higher section (skip floors) |
| `elevator-down` | Enclosed downward arrow `[v]` | Jump to distant lower section (skip floors) |
| `stairs` | Zigzag line | Sequential navigation (step by step) |
| `door` | Rectangle with handle | Enter a new zone or domain |
| `exit` | Rectangle with outward arrow | Leave current zone, return to index |
| `breadcrumb` | Series of `>` separators | Current position trail |
| `home` | Simple house outline | Return to root / main index |
| `search` | Circle with diagonal line | Open search interface |

```css
.hic-nav-glyph {
  color: var(--hic-electric-blue);
  font-family: var(--hic-mono);
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  text-shadow: 0 0 4px rgba(123, 140, 222, 0.3);
  cursor: pointer;
  transition: text-shadow 0.2s ease;
}

.hic-nav-glyph:hover {
  text-shadow: 0 0 8px rgba(123, 140, 222, 0.5);
}

.hic-nav-glyph--restricted {
  color: var(--hic-crimson-alert);
  text-shadow: 0 0 4px rgba(196, 90, 90, 0.3);
}

.hic-breadcrumb {
  font-family: var(--hic-mono);
  font-size: 0.75rem;
  color: var(--hic-fog);
  letter-spacing: 0.05em;
}

.hic-breadcrumb__separator {
  color: var(--hic-ash);
  margin: 0 0.5rem;
}

.hic-breadcrumb__current {
  color: var(--hic-electric-blue);
  text-shadow: 0 0 4px rgba(123, 140, 222, 0.3);
}
```

**Breadcrumb Example:**

```
HOME > FLOOR 23 > COMMENTARY MEZZANINE > CROSS-REFERENCES > DOC-1847
 ash    ash             fog                    fog           blue
```

### 7.5 The HIC Logo and Construction Grid

The HIC logo is a stylized skyscraper silhouette composed of neon lines,
representing the building itself as seen from a low exterior angle at night.

**Logo Specifications:**

| Property | Value |
|---|---|
| Aspect Ratio | 1:2.5 (width:height) |
| Minimum Size | 32px wide / 80px tall |
| Construction Grid | 8x20 unit grid |
| Stroke Weight | 2px at standard size, scaled proportionally |
| Primary Color | `--hic-electric-blue` at full opacity |
| Glow | 4px blur at 40% opacity in Electric Blue |
| Background | Must be placed on `--hic-void` or darker |
| Clear Space | 1 grid unit (4px at minimum size) on all sides |

**Logo Construction (ASCII Approximation):**

```
    Construction Grid (8x20 units)
    ==============================

         ||
        ||||
       ||  ||
       ||  ||
      |||  |||
      ||    ||
     |||    |||
     ||      ||
     ||      ||
    |||      |||
    ||        ||
    ||   HIC  ||
    ||        ||
    ||________||
   /|||||||||||\
  / ||||||||||||\
 /________________\

  The vertical lines represent the building's
  glass-and-steel frame. The "HIC" text sits
  centered at approximately 60% of the total
  height. The base widens slightly to suggest
  the building's foundation and lobby level.
```

**Logo Usage Rules:**

1. The logo is NEVER rendered in colors other than Electric Blue, Bone White,
   or Stark White.
2. The logo is NEVER placed on backgrounds lighter than `--hic-gunmetal`.
3. The logo ALWAYS has its neon glow active. A logo without glow is dead
   signage -- an unacceptable state.
4. The logo may be animated with a slow `hic-glow-breathe` effect (3s cycle)
   on landing pages and loading screens ONLY.
5. The logo is NEVER rotated, skewed, stretched, or otherwise distorted.

```css
.hic-logo {
  display: inline-block;
  color: var(--hic-electric-blue);
  filter: drop-shadow(0 0 4px rgba(123, 140, 222, 0.4));
}

.hic-logo--landing {
  animation: hic-glow-breathe 3s ease-in-out infinite;
}

.hic-logo--monochrome {
  color: var(--hic-stark);
  filter: drop-shadow(0 0 4px rgba(232, 232, 240, 0.3));
}
```

---

## 8. Print Compatibility

### 8.1 Graceful Degradation Philosophy

The HIC was designed for screens in dark environments. Print is a hostile
medium for neon aesthetics -- there is no backlight, no glow, no animation.
The print compatibility system does not attempt to replicate the HIC aesthetic
on paper. Instead, it translates the semantic system into print-native
equivalents.

The guiding principle: **information survives, atmosphere does not.**

### 8.2 Print Stylesheet

```css
@media print {
  /* Reset the void */
  :root {
    --hic-void:     #ffffff;
    --hic-charcoal: #f5f5f5;
    --hic-gunmetal: #eeeeee;
    --hic-slate:    #cccccc;
    --hic-ash:      #999999;
    --hic-fog:      #666666;
    --hic-bone:     #333333;
    --hic-stark:    #000000;
  }

  body {
    background: white;
    color: #333333;
    font-size: 10pt;
    line-height: 1.5;
  }

  /* Remove all glows and shadows */
  * {
    text-shadow: none !important;
    box-shadow: none !important;
    filter: none !important;
    animation: none !important;
  }

  /* Signal colors become print-safe equivalents */
  .hic-link,
  .hic-active-content { color: #003399; }

  .hic-draft-marker,
  .hic-review-needed  { color: #996600; }

  .hic-deprecated,
  .hic-security-warning { color: #990000; }

  .hic-verified,
  .hic-approved { color: #006600; }

  .hic-annotation,
  .hic-cross-ref { color: #660099; }

  .hic-log-output,
  .hic-data-stream { color: #006666; }

  .hic-emergency,
  .hic-break-glass { color: #cc0066; font-weight: bold; }

  /* Code blocks become bordered boxes */
  pre {
    background: #f5f5f5;
    border: 1px solid #cccccc;
    border-left: 3px solid #333333;
    padding: 8pt;
    font-size: 8pt;
  }

  code {
    background: #f0f0f0;
    border: 1px solid #dddddd;
    color: #333333;
    padding: 1pt 3pt;
  }

  /* Tables get full borders for print clarity */
  table {
    border-collapse: collapse;
    width: 100%;
  }

  th, td {
    border: 1px solid #999999;
    padding: 4pt 6pt;
    font-size: 9pt;
  }

  th {
    background: #eeeeee;
    font-weight: bold;
    color: #000000;
  }

  /* Warning lights become text labels */
  .hic-warning-light {
    display: none;
  }

  .hic-warning-light::after {
    display: inline;
    font-size: 8pt;
    font-weight: bold;
    text-transform: uppercase;
  }

  .hic-warning-light--aging::after    { content: '[AGING]'; color: #996600; }
  .hic-warning-light--stale::after    { content: '[STALE]'; color: #996600; }
  .hic-warning-light--broken::after   { content: '[BROKEN]'; color: #990000; }
  .hic-warning-light--deprecated::after { content: '[DEPRECATED]'; color: #990000; }

  /* Navigation elements hidden in print */
  .hic-pathline,
  .hic-nav-glyph,
  .hic-breadcrumb,
  .hic-floor-marker { display: none; }

  /* Page break management */
  h1, h2, h3 { page-break-after: avoid; }
  pre, table  { page-break-inside: avoid; }
  p           { orphans: 3; widows: 3; }
}
```

### 8.3 High Contrast Mode

For accessibility compliance and users who require maximum readability, the
HIC supports a high-contrast mode that strips away all subtlety in favor of
raw legibility.

```css
@media (prefers-contrast: more) {
  :root {
    --hic-void:     #000000;
    --hic-charcoal: #000000;
    --hic-gunmetal: #111111;
    --hic-slate:    #ffffff;
    --hic-ash:      #ffffff;
    --hic-fog:      #ffffff;
    --hic-bone:     #ffffff;
    --hic-stark:    #ffffff;

    /* High contrast signal colors */
    --hic-electric-blue:  #6699ff;
    --hic-amber-warning:  #ffcc00;
    --hic-crimson-alert:  #ff3333;
    --hic-emerald-status: #33ff33;
    --hic-violet-accent:  #cc66ff;
    --hic-cyan-data:      #00ffff;
    --hic-hot-pink:       #ff00ff;
  }

  body {
    background: #000000;
    color: #ffffff;
    font-size: 1rem;
    line-height: 1.8;
  }

  a {
    color: #6699ff;
    text-decoration: underline;
    text-decoration-thickness: 2px;
  }

  /* All borders become high-visibility */
  table, th, td {
    border: 2px solid #ffffff;
  }

  code {
    border: 2px solid #33ff33;
    padding: 0.2em 0.5em;
  }

  pre {
    border: 2px solid #33ff33;
    border-left: 4px solid #33ff33;
  }

  /* Remove all animations */
  * {
    animation: none !important;
    transition: none !important;
  }

  /* Status indicators become larger and text-backed */
  .hic-warning-light {
    width: 12px;
    height: 12px;
    border: 2px solid currentColor;
  }
}
```

### 8.4 Grayscale Fallbacks

When color reproduction is unavailable (monochrome displays, e-ink devices,
photocopied documents), the HIC signal system falls back to **pattern-based
differentiation**.

| Signal Color | Grayscale Value | Pattern Indicator |
|---|---|---|
| Electric Blue | 60% gray | Solid underline |
| Amber Warning | 50% gray | Dashed underline |
| Crimson Alert | 40% gray, bold | Double underline |
| Emerald Status | 55% gray | No underline, checkmark prefix |
| Violet Accent | 45% gray, italic | Dotted underline |
| Cyan Data | 65% gray | Monospace with bracket prefix `[>]` |
| Hot Pink | 35% gray, bold, uppercase | Triple underline, exclamation prefix |

```css
@media (prefers-color-scheme: light), print {
  .hic-grayscale-mode .hic-link {
    color: #666666;
    text-decoration: underline;
    text-decoration-style: solid;
  }

  .hic-grayscale-mode .hic-caution-badge {
    color: #555555;
    text-decoration: underline;
    text-decoration-style: dashed;
  }

  .hic-grayscale-mode .hic-critical-badge {
    color: #444444;
    font-weight: bold;
    text-decoration: underline;
    text-decoration-style: double;
  }

  .hic-grayscale-mode .hic-verified {
    color: #555555;
  }

  .hic-grayscale-mode .hic-verified::before {
    content: '\2713 ';  /* Checkmark prefix */
  }

  .hic-grayscale-mode .hic-annotation {
    color: #4a4a4a;
    font-style: italic;
    text-decoration: underline;
    text-decoration-style: dotted;
  }

  .hic-grayscale-mode .hic-data-stream {
    color: #666666;
    font-family: var(--hic-mono);
  }

  .hic-grayscale-mode .hic-data-stream::before {
    content: '[>] ';
  }

  .hic-grayscale-mode .hic-emergency {
    color: #333333;
    font-weight: bold;
    text-transform: uppercase;
    text-decoration: underline;
    text-decoration-style: double;
  }

  .hic-grayscale-mode .hic-emergency::before {
    content: '!!! ';
  }
}
```

The grayscale system ensures that even when every photon of color has been
stripped away, the *semantic structure* of the HIC persists. A printed,
photocopied, faxed document from the HIC still communicates urgency levels,
content states, and domain classifications. The neon may be gone, but the
signal survives.

---

## 9. Appendices

### Appendix A: Complete CSS Variable Reference

```css
:root {
  /* ============================== */
  /* STRUCTURAL NEUTRALS            */
  /* ============================== */
  --hic-void:     #0a0a0f;
  --hic-charcoal: #12121a;
  --hic-gunmetal: #1a1a2e;
  --hic-slate:    #2a2a3e;
  --hic-ash:      #4a4a5e;
  --hic-fog:      #8a8a9e;
  --hic-bone:     #c8c8d4;
  --hic-stark:    #e8e8f0;

  /* ============================== */
  /* PRIMARY SIGNAL COLORS          */
  /* ============================== */
  --hic-electric-blue:  #7b8cde;
  --hic-amber-warning:  #c49a3a;
  --hic-crimson-alert:  #c45a5a;
  --hic-emerald-status: #5fa85f;
  --hic-violet-accent:  #9b59b6;
  --hic-cyan-data:      #00bcd4;
  --hic-hot-pink:       #ff1493;

  /* ============================== */
  /* SIGNAL COLOR RGB VALUES        */
  /* (for rgba() usage)             */
  /* ============================== */
  --hic-electric-blue-rgb:  123, 140, 222;
  --hic-amber-warning-rgb:  196, 154, 58;
  --hic-crimson-alert-rgb:  196, 90, 90;
  --hic-emerald-status-rgb: 95, 168, 95;
  --hic-violet-accent-rgb:  155, 89, 182;
  --hic-cyan-data-rgb:      0, 188, 212;
  --hic-hot-pink-rgb:       255, 20, 147;

  /* ============================== */
  /* TYPOGRAPHY                     */
  /* ============================== */
  --hic-mono:    'JetBrains Mono', 'Fira Code', 'Source Code Pro', 'Courier New', monospace;
  --hic-display: 'Share Tech Mono', 'VT323', monospace;
  --hic-sans:    'Inter', 'IBM Plex Sans', 'Helvetica Neue', sans-serif;

  /* ============================== */
  /* SPACING SCALE                  */
  /* ============================== */
  --hic-space-xs:  0.25rem;
  --hic-space-sm:  0.5rem;
  --hic-space-md:  1rem;
  --hic-space-lg:  1.5rem;
  --hic-space-xl:  2rem;
  --hic-space-xxl: 3rem;

  /* ============================== */
  /* ANIMATION DURATIONS            */
  /* ============================== */
  --hic-duration-instant: 0.1s;
  --hic-duration-fast:    0.2s;
  --hic-duration-normal:  0.3s;
  --hic-duration-slow:    0.5s;
  --hic-duration-pulse:   2s;
  --hic-duration-breathe: 3s;
}
```

### Appendix B: WCAG Compliance Matrix

All signal colors have been tested against their typical backgrounds for WCAG
2.1 AA compliance (minimum 4.5:1 contrast ratio for normal text, 3:1 for
large text).

| Color | On Void (#0a0a0f) | On Charcoal (#12121a) | On Gunmetal (#1a1a2e) | AA Normal | AA Large |
|---|---|---|---|---|---|
| Electric Blue (#7b8cde) | 6.2:1 | 5.8:1 | 5.1:1 | PASS | PASS |
| Amber Warning (#c49a3a) | 5.9:1 | 5.5:1 | 4.8:1 | PASS | PASS |
| Crimson Alert (#c45a5a) | 4.6:1 | 4.3:1 | 3.8:1 | PASS* | PASS |
| Emerald Status (#5fa85f) | 5.4:1 | 5.0:1 | 4.4:1 | PASS | PASS |
| Violet Accent (#9b59b6) | 4.1:1 | 3.8:1 | 3.4:1 | FAIL* | PASS |
| Cyan Data (#00bcd4) | 7.1:1 | 6.6:1 | 5.8:1 | PASS | PASS |
| Hot Pink (#ff1493) | 5.3:1 | 4.9:1 | 4.3:1 | PASS | PASS |
| Bone White (#c8c8d4) | 10.2:1 | 9.5:1 | 8.4:1 | PASS | PASS |
| Stark White (#e8e8f0) | 13.8:1 | 12.9:1 | 11.4:1 | PASS | PASS |
| Fog (#8a8a9e) | 5.1:1 | 4.7:1 | 4.2:1 | PASS | PASS |
| Ash (#4a4a5e) | 2.4:1 | 2.2:1 | 1.9:1 | FAIL | FAIL |

**Notes:**
- (*) Crimson Alert on Gunmetal falls below 4.5:1. Use Bright Variant (#d46a6a) for normal text on Gunmetal backgrounds.
- (*) Violet Accent fails AA Normal on all backgrounds. Violet is used exclusively for meta-content that is always accompanied by a border indicator, ensuring discoverability is not color-dependent.
- Ash Gray is used only for disabled/placeholder content and never for actionable text.

### Appendix C: Quick Reference Card

```
+================================================================+
|                                                                |
|  HIC VISUAL DESIGN SYSTEM -- QUICK REFERENCE                  |
|                                                                |
|================================================================|
|                                                                |
|  SIGNAL COLORS                                                 |
|  ~~~~~~~~~~~~~                                                 |
|  #7b8cde  Electric Blue .... Active knowledge, links          |
|  #c49a3a  Amber Warning .... Caution, drafts, review          |
|  #c45a5a  Crimson Alert .... Critical, security, broken       |
|  #5fa85f  Emerald Status ... Verified, approved, healthy      |
|  #9b59b6  Violet Accent .... Meta, commentary, refs           |
|  #00bcd4  Cyan Data ........ Archives, logs, raw data         |
|  #ff1493  Hot Pink ......... Emergency, break-glass           |
|                                                                |
|  NEUTRALS                                                      |
|  ~~~~~~~~                                                      |
|  #0a0a0f  Void Black       #4a4a5e  Ash Gray                  |
|  #12121a  Deep Charcoal    #8a8a9e  Fog                       |
|  #1a1a2e  Gunmetal         #c8c8d4  Bone White                |
|  #2a2a3e  Slate Edge       #e8e8f0  Stark White               |
|                                                                |
|  TYPEFACES                                                     |
|  ~~~~~~~~~                                                     |
|  Primary:  JetBrains Mono                                      |
|  Display:  Share Tech Mono                                     |
|  Fallback: Inter                                               |
|                                                                |
|  RULES                                                         |
|  ~~~~~                                                         |
|  01. Signal colors never touch (8px neutral gap)               |
|  02. Max 3 signal colors per viewport                          |
|  03. Emerald is never animated                                 |
|  04. Crimson and Hot Pink are mutually exclusive               |
|  05. All motion is semantic, never decorative                  |
|                                                                |
+================================================================+
```

---

**Document End**

**Classification:** HOLM-VIS-004 | **Floor:** 47 | **Zone:** Visual Systems Division
**Last Review:** 2026-02-17 | **Status:** APPROVED | **Luminance:** Blazing

```
// END TRANSMISSION
// HOLM INTELLIGENCE COMPLEX
// VISUAL DESIGN SYSTEM v3.1
// THE BUILDING SPEAKS IN LIGHT
```
