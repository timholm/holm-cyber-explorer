# HIC Spatial Data Architecture

## Blueprint Formats, Coordinate Systems & Layer Management

**Document Class:** SPATIAL-CORE
**Clearance Level:** Structural Engineer / Cartographer
**HIC Registry:** HIC-SYS-SPATIAL-0001
**Revision:** 4.0.0
**Last Survey:** 2026-02-17

---

> *"Every wall is a schema boundary. Every corridor, a navigation path.
> The building does not represent the documentation --- the building IS the documentation."*
>
> --- Kael Holm, Founding Architect, on the day the first neon grid flickered to life

---

## Table of Contents

1. [Preamble: The Geometry of Knowledge](#preamble-the-geometry-of-knowledge)
2. [The Coordinate System](#the-coordinate-system)
   - [Axis Definitions](#axis-definitions)
   - [The Origin Point](#the-origin-point)
   - [Coordinate Notation](#coordinate-notation)
   - [Coordinate Resolution and Precision](#coordinate-resolution-and-precision)
   - [Boundary Conditions and Edge Coordinates](#boundary-conditions-and-edge-coordinates)
3. [Blueprint Format Specification](#blueprint-format-specification)
   - [Floor Plan Schema](#floor-plan-schema)
   - [Room Definitions](#room-definitions)
   - [Door Connections](#door-connections)
   - [Window Views](#window-views)
   - [Utility Conduits](#utility-conduits)
4. [Layer System](#layer-system)
   - [Layer 0: Structural Layer](#layer-0-structural-layer)
   - [Layer 1: Utility Layer](#layer-1-utility-layer)
   - [Layer 2: Interior Layer](#layer-2-interior-layer)
   - [Layer 3: Signage Layer](#layer-3-signage-layer)
   - [Layer 4: Atmospheric Layer](#layer-4-atmospheric-layer)
   - [Layer 5: Holographic Layer](#layer-5-holographic-layer)
   - [Layer Compositing and Render Order](#layer-compositing-and-render-order)
5. [Data Encoding](#data-encoding)
   - [Transformation Pipeline](#transformation-pipeline)
   - [Heading Levels as Spatial Partitions](#heading-levels-as-spatial-partitions)
   - [Content Density Calculations](#content-density-calculations)
   - [Media Object Placement](#media-object-placement)
6. [Scale and Proportions](#scale-and-proportions)
   - [Standard Dimensions](#standard-dimensions)
   - [Dynamic Scaling Rules](#dynamic-scaling-rules)
   - [Vertical Proportions](#vertical-proportions)
7. [Collision Detection](#collision-detection)
   - [Namespace Collisions](#namespace-collisions)
   - [Anchor Uniqueness Enforcement](#anchor-uniqueness-enforcement)
   - [Redirect Corridors](#redirect-corridors)
   - [Overlap Resolution Algorithms](#overlap-resolution-algorithms)
8. [Cartographic Standards](#cartographic-standards)
   - [Map Projections](#map-projections)
   - [Legend and Key Standards](#legend-and-key-standards)
   - [Scale Bars and Measurement Units](#scale-bars-and-measurement-units)
   - [Orientation and the Constitutional Compass](#orientation-and-the-constitutional-compass)
9. [Versioning and Time](#versioning-and-time)
   - [Historical Blueprints](#historical-blueprints)
   - [Construction Timeline](#construction-timeline)
   - [Renovation Logs](#renovation-logs)
   - [Demolition Records](#demolition-records)
10. [Appendices](#appendices)

---

## Preamble: The Geometry of Knowledge

The Holm Intelligence Complex is not a metaphor layered on top of a documentation system. It is a spatial data structure that *happens* to render as a cyberpunk skyscraper. Every pixel of neon trim, every hum of a fluorescent corridor light, every holographic sign floating above a doorway --- these are the direct, computable consequences of the underlying documentation data.

This document defines the spatial encoding that makes that possible.

When a contributor writes a markdown file, they are not merely authoring text. They are pouring concrete, running electrical conduit, and hanging signage in a living building. The transformation is deterministic: given the same documentation corpus, the same building will materialize, down to the flicker pattern on the lobby's status board.

The HIC Spatial Data Architecture governs three fundamental concerns:

- **Where** every piece of content exists in three-dimensional space (coordinate system)
- **What** the structural encoding looks like at rest (blueprint format)
- **How** the raw data becomes rendered geometry (transformation pipeline and layers)

Read this document as both a technical specification and a construction manual. The building is always under construction. The cranes never stop.

---

## The Coordinate System

### Axis Definitions

The HIC exists in a right-handed Cartesian coordinate system with three semantic axes. Unlike arbitrary spatial frameworks, each axis carries intrinsic meaning tied to the documentation structure.

#### X-Axis: Document Position (Horizontal Sequence)

The X-axis represents an article's sequential position within its domain floor. Articles are laid out left-to-right in their canonical ordering, much like rooms along a corridor.

| Property         | Value                                              |
| ---------------- | -------------------------------------------------- |
| Direction        | West to East (left to right when facing the tower) |
| Unit             | Article Index (1-based integer)                    |
| Physical Scale   | 1 unit = 1 standard room width (4 meters)          |
| Minimum Value    | 1 (first article in the domain)                    |
| Maximum Value    | Unbounded (grows with domain content)              |
| Padding          | 2m corridor gaps between each article unit         |

```
X-axis layout (Floor 7: Authentication Domain)

  [A001]---[A002]---[A003]---[A004]---[A005]---...
  Login    OAuth    SSO      MFA      Session
  Guide    Setup    Config   Setup    Mgmt

  |--4m--|--2m--|--4m--|--2m--|--4m--|--2m--|...
   room   corr   room   corr   room
```

The X-axis is *not* continuous in the mathematical sense. It is a discrete lattice. Each integer position holds exactly one article. The 2-meter corridor gaps between rooms are implicit and do not consume an index value.

#### Y-Axis: Floor / Domain Number (Vertical Position)

The Y-axis maps directly to the domain hierarchy. Floor 1 is the Constitution --- the foundational governance document. Above it, domains stack according to their registry order.

| Property         | Value                                                      |
| ---------------- | ---------------------------------------------------------- |
| Direction        | Ground to sky (ascending floor numbers)                    |
| Unit             | Floor Number (1-based integer)                             |
| Physical Scale   | 1 unit = 1 floor height (4 meters floor-to-floor)         |
| Minimum Value    | 0 (Sub-basement: system infrastructure, not public-facing) |
| Maximum Value    | Unbounded (new domains add new floors)                     |
| Special Values   | Floor 0 = Sub-basement; Floor 1 = Constitution             |

```
Y-axis layout (cross-section, south face)

  Floor N  +---------------------+  Domain N
  ...      |         ...         |  ...
  Floor 5  +---------------------+  Deployment
  Floor 4  +---------------------+  API Reference
  Floor 3  +---------------------+  Architecture
  Floor 2  +---------------------+  Tutorials
  Floor 1  |=====================|  Constitution (Origin)
  Floor 0  +---------------------+  Sub-basement (Infra)
           #######################  Foundation
```

Floor 0 (the sub-basement) houses system-level infrastructure: build configuration, CI/CD pipeline definitions, template engines, and the rendering machinery itself. It is not publicly navigable but can be accessed by maintainers with structural clearance.

#### Z-Axis: Depth / Detail Level

The Z-axis is the most unusual of the three. It represents information depth --- how far you drill into a piece of content. The surface level (Z=0) is the summary or overview. Increasing Z values take you deeper into raw data, implementation details, and source material.

| Property         | Value                                                           |
| ---------------- | --------------------------------------------------------------- |
| Direction        | Front face to back face (surface to depth)                      |
| Unit             | Depth Level (0-based integer)                                   |
| Physical Scale   | 1 unit = 6 meters (standard room depth)                         |
| Minimum Value    | 0 (surface: title, summary, TL;DR)                             |
| Maximum Value    | 4 (deepest: raw data, source code, unprocessed references)      |

```
Z-axis depth levels

  Z=0  Surface     Title card, abstract, one-line summary
  Z=1  Overview    Introduction, key concepts, prerequisites
  Z=2  Standard    Full article body, diagrams, examples
  Z=3  Detailed    Implementation notes, edge cases, caveats
  Z=4  Raw         Source data, API response dumps, logs
```

Visually, the Z-axis manifests as room depth. When you enter a room from the corridor (Z=0), you see the summary signage on the near wall. Walk deeper into the room, past the main workbenches (Z=2), and you reach the back wall where the raw terminals and data conduits are exposed (Z=4).

### The Origin Point

Every coordinate system needs an origin. In the HIC, the origin is both physically and philosophically fixed:

```
ORIGIN = HIC-F01-A001-D0

  Floor:    1   (The Constitution)
  Article:  1   (First article of the Constitution)
  Depth:    0   (Surface level --- the preamble)
```

This is the main lobby entrance. It is the first thing a visitor sees. The neon sign above the revolving door reads the title of the Constitution's opening article. Every coordinate in the building is measured relative to this point.

The origin is *immutable*. Even if the Constitution is amended, Article 1 remains Article 1. Content may change; the coordinate does not. This is a foundational guarantee of the spatial system --- coordinates are stable addresses, not content-dependent labels.

### Coordinate Notation

The HIC uses a structured notation system for referencing any point in the building:

```
HIC-F{floor}-A{article}-D{depth}
```

#### Format Breakdown

| Component | Format    | Description                          | Example     |
| --------- | --------- | ------------------------------------ | ----------- |
| Prefix    | `HIC`     | Identifies the coordinate as HIC     | `HIC`       |
| Floor     | `F{nn}`   | Zero-padded floor number (2+ digits) | `F01`, `F12`|
| Article   | `A{nnn}`  | Zero-padded article index (3+ digits)| `A001`      |
| Depth     | `D{n}`    | Single digit depth level (0-4)       | `D2`        |

#### Examples

```
HIC-F01-A001-D0    The front door. Constitution, Article 1, surface.
HIC-F01-A001-D4    Deep in the Constitution's raw governance data.
HIC-F03-A015-D2    Architecture floor, 15th article, standard depth.
HIC-F12-A003-D1    Floor 12, 3rd article, overview level.
```

#### Shorthand Notation

For convenience, several shorthand forms are recognized:

```
F3:A15:D2          Compact form (no prefix, colon-delimited)
F3:A15             Implies D0 (surface level)
F3:*               All articles on Floor 3 (wildcard)
F3:A15:*           All depth levels of Floor 3, Article 15
*:*:D0             Surface level of every room in the building
```

#### Coordinate Ranges

Ranges are expressed with double-dot notation:

```
HIC-F03-A001..A020-D0    Articles 1-20 on Floor 3, surface level
HIC-F01..F05-A001-D0     Article 1, Floors 1-5, surface level
HIC-F03-A015-D0..D4      All depths of a single article
```

### Coordinate Resolution and Precision

Not every coordinate maps to a populated room. The building is sparse --- many potential coordinates are empty voids (unwritten articles, non-existent floors). The spatial engine distinguishes between:

- **Populated coordinates**: A room exists here. Content is rendered.
- **Reserved coordinates**: The floor exists but this article slot is empty. Shown as a dark, locked door.
- **Void coordinates**: The floor itself does not exist. No geometry is generated.

Resolution queries return one of three statuses:

```json
{
  "coordinate": "HIC-F03-A015-D2",
  "status": "populated",
  "content_hash": "a4f8c3...",
  "last_modified": "2026-02-10T14:32:00Z"
}
```

```json
{
  "coordinate": "HIC-F03-A099-D0",
  "status": "reserved",
  "reason": "article_index_exceeds_domain_content"
}
```

```json
{
  "coordinate": "HIC-F99-A001-D0",
  "status": "void",
  "reason": "floor_does_not_exist"
}
```

### Boundary Conditions and Edge Coordinates

Several special coordinates carry semantic weight:

| Coordinate Pattern       | Meaning                                         |
| ------------------------ | ----------------------------------------------- |
| `HIC-F{n}-A001-D0`      | The lobby of floor N (domain landing page)      |
| `HIC-F{n}-AMAX-D0`      | The last room on floor N (final article)        |
| `HIC-F01-A001-D0`       | The building entrance (absolute origin)         |
| `HIC-FMAX-A001-D0`      | The rooftop lobby (highest domain)              |
| `HIC-F00-A001-D4`       | Deepest infrastructure point (build system core)|

The `MAX` keyword is dynamically resolved at query time. It always points to the current maximum populated index.

---

## Blueprint Format Specification

### Floor Plan Schema

Every floor of the HIC is described by a structured blueprint file. These blueprints are generated from the documentation source and stored as JSON. They serve as the intermediate representation between raw markdown and rendered spatial geometry.

#### Root Schema

```yaml
# HIC Floor Blueprint Schema v4.0
# File: blueprints/floor-{nn}.blueprint.json

schema_version: "4.0"
floor:
  number: 3                          # Y-axis value
  domain: "architecture"             # Domain identifier
  display_name: "Architecture"       # Human-readable name
  neon_color: "#00FFA3"              # Primary neon accent for this floor
  secondary_color: "#003D28"         # Shadow/background accent
  clearance: "public"                # Access level required

  dimensions:
    width_units: 20                  # Number of article slots (X-axis max)
    depth_units: 5                   # Depth levels available (Z-axis max)
    floor_height_m: 4.0              # Physical floor height in meters
    corridor_width_m: 2.0            # Corridor width in meters

  lobby:
    position: "A001"                 # Always the first article
    width_m: 12.0                    # Lobby width (scales with floor size)
    features:
      - "directory_hologram"         # Floating floor map
      - "elevator_bank"             # Vertical navigation
      - "status_board"              # Domain health indicators

  rooms: []                          # Array of room definitions (see below)
  corridors: []                      # Array of corridor definitions
  utilities: []                      # Array of utility conduit definitions
```

#### Full Floor Blueprint Example

```json
{
  "schema_version": "4.0",
  "floor": {
    "number": 3,
    "domain": "architecture",
    "display_name": "Architecture",
    "neon_color": "#00FFA3",
    "secondary_color": "#003D28",
    "clearance": "public",
    "dimensions": {
      "width_units": 15,
      "depth_units": 5,
      "floor_height_m": 4.0,
      "corridor_width_m": 2.0
    },
    "lobby": {
      "position": "A001",
      "width_m": 12.0,
      "features": ["directory_hologram", "elevator_bank", "status_board"]
    },
    "rooms": [
      {
        "article_index": 1,
        "slug": "architecture-overview",
        "title": "Architecture Overview",
        "coordinate": "HIC-F03-A001",
        "dimensions": { "width_m": 4.0, "depth_m": 6.0, "height_m": 4.0 },
        "content_hash": "a4f8c3de...",
        "word_count": 2400,
        "heading_partitions": 8,
        "doors": [],
        "windows": [],
        "depth_layers": [
          { "z": 0, "type": "surface", "content": "title_card" },
          { "z": 1, "type": "overview", "content": "introduction" },
          { "z": 2, "type": "standard", "content": "full_body" },
          { "z": 3, "type": "detailed", "content": "implementation_notes" },
          { "z": 4, "type": "raw", "content": "source_diagrams" }
        ]
      }
    ],
    "corridors": [
      {
        "id": "main-hall-f03",
        "type": "primary",
        "width_m": 2.0,
        "connects": ["A001", "A015"],
        "lighting": "neon_strip",
        "signage": "overhead_holographic"
      }
    ],
    "utilities": []
  }
}
```

### Room Definitions

A room is the spatial manifestation of a single documentation article. The mapping is strict and deterministic: one article, one room. No article spans multiple rooms. No room contains multiple articles.

#### Room Anatomy

```
+---------------------------------------------+
|  DOOR (from corridor)                       |  Z=0  Surface Wall
|  +-------------------------------------+    |       Title placard,
|  |  TITLE PLACARD    [status light]    |    |       status indicator
|  +-------------------------------------+    |
|                                             |
|  +-------------------------------------+    |  Z=1  Overview Zone
|  |  Introduction text                  |    |       Key concepts,
|  |  Prerequisite callout boxes         |    |       prerequisites
|  |  Navigation breadcrumb strip        |    |
|  +-------------------------------------+    |
|                                             |
|  +-------------------------------------+    |  Z=2  Standard Zone
|  |  Main article body                  |    |       Full content,
|  |  Code blocks (illuminated panels)   |    |       examples,
|  |  Diagrams (wall-mounted displays)   |    |       diagrams
|  |  Tables (workbench surfaces)        |    |
|  +-------------------------------------+    |
|                                             |
|  +-------------------------------------+    |  Z=3  Detailed Zone
|  |  Edge case documentation            |    |       Implementation
|  |  Performance considerations         |    |       notes, caveats,
|  |  Compatibility matrices             |    |       deep context
|  +-------------------------------------+    |
|                                             |
|  +-------------------------------------+    |  Z=4  Raw Zone
|  |  Source data terminals              |    |       Raw data,
|  |  API response dumps                 |    |       source refs,
|  |  Log readouts (scrolling screens)   |    |       unprocessed
|  +-------------------------------------+    |       material
|                                             |
+---------------------------------------------+
      4m wide                6m deep
```

#### Room Properties Schema

```json
{
  "article_index": 5,
  "slug": "event-driven-architecture",
  "title": "Event-Driven Architecture",
  "coordinate": "HIC-F03-A005",
  "dimensions": {
    "width_m": 4.0,
    "depth_m": 6.0,
    "height_m": 4.0
  },
  "content_hash": "b7e2f1aa...",
  "word_count": 3200,
  "heading_partitions": 12,
  "content_density": 133.3,
  "freshness": {
    "last_modified": "2026-02-14T09:00:00Z",
    "age_days": 3,
    "glow_intensity": 0.92
  },
  "doors": [
    {
      "id": "door-f03-a005-to-f03-a002",
      "type": "hyperlink",
      "target": "HIC-F03-A002",
      "label": "See also: Microservices Guide",
      "wall": "west",
      "style": "glass_panel"
    }
  ],
  "windows": [
    {
      "id": "win-f03-a005-preview-f04-a010",
      "type": "preview",
      "target": "HIC-F04-A010",
      "wall": "north",
      "shows": "title_and_summary"
    }
  ],
  "occupancy": {
    "current_visitors": 14,
    "peak_visitors_24h": 89,
    "heat_level": "warm"
  }
}
```

### Door Connections

Doors are the spatial realization of hyperlinks. When one documentation article links to another, a door appears in the source room's wall leading to the target room.

#### Door Types

| Door Type        | Link Type              | Visual Treatment                        |
| ---------------- | ---------------------- | --------------------------------------- |
| `glass_panel`    | Internal same-floor    | Transparent door, see-through to target |
| `steel_fire`     | Internal cross-floor   | Heavy door with floor indicator display |
| `neon_portal`    | External link          | Glowing portal frame, opens to outside  |
| `service_hatch`  | Footnote / reference   | Small hatch at knee level               |
| `emergency_exit` | Deprecated link        | Red-lit exit sign, may lead to void     |
| `revolving`      | Bidirectional link     | Revolving door, works both ways         |

#### Door Placement Algorithm

Doors are placed on room walls according to the target's relative position:

```
Target is to the WEST  (lower article index, same floor)   -> West wall
Target is to the EAST  (higher article index, same floor)   -> East wall
Target is ABOVE        (higher floor number)                 -> Ceiling hatch
Target is BELOW        (lower floor number)                  -> Floor hatch
Target is CROSS-FLOOR  (different floor, any article)        -> North wall
Target is EXTERNAL     (outside the HIC entirely)            -> South wall (exterior)
```

#### Door Schema

```json
{
  "id": "door-{source}-to-{target}",
  "type": "glass_panel | steel_fire | neon_portal | service_hatch | emergency_exit | revolving",
  "source": "HIC-F03-A005-D2",
  "target": "HIC-F03-A002-D0",
  "label": "Human-readable link text",
  "wall": "north | south | east | west | ceiling | floor",
  "position_on_wall": 0.5,
  "style": {
    "frame_color": "#00FFA3",
    "glow": true,
    "glow_intensity": 0.7,
    "glass_opacity": 0.3
  },
  "metadata": {
    "link_type": "inline | reference | nav | breadcrumb",
    "anchor": "#section-name",
    "bidirectional": false,
    "last_verified": "2026-02-17T00:00:00Z",
    "http_status": 200
  }
}
```

#### Broken Link Visualization

When a door's target no longer exists (404), the door enters a failure state:

- The door frame flickers between its original neon color and a deep red (`#FF003C`).
- A holographic caution tape appears across the doorway.
- The door handle sparks when touched, displaying: `TARGET NOT FOUND: {coordinate}`.
- After 30 days in failure state, the door is bricked over (removed from the blueprint) and logged in demolition records.

### Window Views

Windows are preview panels --- they show a glimpse of content in an adjacent or related room without requiring the visitor to leave their current location. They map to documentation concepts like "related articles" sidebars, tooltip previews, and summary cards.

#### Window Types

| Window Type       | Documentation Analog        | Renders                              |
| ----------------- | --------------------------- | ------------------------------------ |
| `clear_glass`     | Related article preview     | Title + first paragraph of target    |
| `frosted`         | Category/tag reference      | Blurred outline, only title visible  |
| `one_way_mirror`  | Analytics dashboard         | Visitor can see out; outside cannot  |
| `display_screen`  | Embedded widget / iframe    | Live-updating content from target    |
| `skylight`        | Parent category view        | Upward view to the domain above      |

#### Window Schema

```json
{
  "id": "win-{source}-{direction}-{target}",
  "type": "clear_glass | frosted | one_way_mirror | display_screen | skylight",
  "source_room": "HIC-F03-A005",
  "target_room": "HIC-F04-A010",
  "wall": "north",
  "dimensions": { "width_m": 1.5, "height_m": 1.0 },
  "position": { "x_offset": 0.3, "y_offset": 1.5 },
  "shows": "title_and_summary | title_only | live_content | category_list",
  "refresh_interval_s": 300,
  "tint_color": "#00FFA3",
  "opacity": 0.8
}
```

### Utility Conduits

Behind every wall, beneath every floor, and above every ceiling, the HIC is laced with utility conduits. These are the invisible (or partially visible) pathways that carry metadata, cross-references, search indices, and dependency information throughout the building.

#### Conduit Types

| Conduit           | Carries                        | Visibility      | Color Code  |
| ----------------- | ------------------------------ | ---------------- | ----------- |
| `power_main`      | Core navigation structure      | Hidden           | `#FFD700`   |
| `data_fiber`      | Search index data              | Hidden           | `#00BFFF`   |
| `cross_ref_pipe`  | Cross-reference links          | Semi-transparent | `#FF6EC7`   |
| `dependency_wire` | Import/dependency declarations | Visible on scan  | `#FF3131`   |
| `metadata_duct`   | Frontmatter, tags, categories  | Hidden           | `#B0B0B0`   |
| `analytics_tube`  | Usage data, visitor counts     | Hidden           | `#8A2BE2`   |

#### Conduit Routing Rules

1. **Same-floor conduits** run horizontally through the corridor walls at ceiling height (3.5m above floor level).
2. **Cross-floor conduits** run vertically through dedicated utility shafts located at the east and west extremities of each floor.
3. **Building-wide conduits** (search index, global navigation) occupy the central utility core --- a vertical shaft running the full height of the building, accessible only from Floor 0.
4. **No conduit may pass through a room's content zone** (Z=2). Conduits route around rooms, never through them. Content integrity is sacred.

#### Conduit Schema

```json
{
  "id": "conduit-{type}-{source}-{target}",
  "type": "power_main | data_fiber | cross_ref_pipe | dependency_wire | metadata_duct | analytics_tube",
  "source": "HIC-F03-A005",
  "target": "HIC-F07-A012",
  "path": [
    { "x": 5, "y": 3, "z": 0, "segment": "horizontal" },
    { "x": 5, "y": 3, "z": 0, "segment": "shaft_east" },
    { "x": 5, "y": 7, "z": 0, "segment": "vertical" },
    { "x": 12, "y": 7, "z": 0, "segment": "horizontal" }
  ],
  "bandwidth": "high | medium | low",
  "latency_ms": 12,
  "color": "#FF6EC7",
  "glow": true,
  "pulse_frequency_hz": 2.0
}
```

When a visitor activates "infrastructure view" (a maintenance-level overlay), the walls become transparent and the conduits light up, pulsing with data flow. The pulse frequency indicates throughput: faster pulses mean heavier cross-referencing between the connected rooms.

---

## Layer System

The HIC renders through a six-layer compositing system. Each layer is independently toggleable, cacheable, and has its own update frequency. The layers stack from bottom (structural foundation) to top (ephemeral overlays).

```
  +-------------------------------------+
  |  Layer 5: HOLOGRAPHIC               |  Floating overlays, search, tooltips
  +-------------------------------------+
  |  Layer 4: ATMOSPHERIC               |  Status glow, alerts, ambient FX
  +-------------------------------------+
  |  Layer 3: SIGNAGE                   |  Labels, breadcrumbs, nav markers
  +-------------------------------------+
  |  Layer 2: INTERIOR                  |  Actual content (text, media, code)
  +-------------------------------------+
  |  Layer 1: UTILITY                   |  Pipes, wires, cross-references
  +-------------------------------------+
  |  Layer 0: STRUCTURAL                |  Walls, floors, ceilings, skeleton
  +-------------------------------------+
```

### Layer 0: Structural Layer

**Purpose:** The building skeleton. Domain hierarchy, floor plates, wall positions, corridor routing.

**Update Frequency:** On documentation tree change (new article added, article deleted, domain restructured).

**Data Sources:**
- Domain registry (floor assignments)
- Article index (room count per floor)
- Navigation tree (corridor layout)

**Render Properties:**

| Property          | Value                                            |
| ----------------- | ------------------------------------------------ |
| Material          | Concrete and steel (dark gray, `#1A1A2E`)        |
| Edge highlighting | Neon trim along structural edges (floor color)   |
| Transparency      | Opaque (alpha = 1.0)                             |
| Shadows           | Cast and receive                                 |
| LOD Levels        | 3 (distant: box; medium: walls; close: textures) |

**Generation Rules:**
1. Each domain in the registry generates one floor plate.
2. Each article in a domain generates one room-sized partition on that floor.
3. The lobby room (A001) is always 3x the standard room width.
4. Corridors connect all rooms on a floor in a linear sequence.
5. Elevator shafts are placed at positions X=1 and X=MAX on every floor.
6. Stairwells (sequential navigation) flank the elevator shafts.

### Layer 1: Utility Layer

**Purpose:** Cross-references, dependencies, data flows, metadata pathways. The nervous system of the building.

**Update Frequency:** On link change (new cross-reference, broken link detected, dependency updated).

**Data Sources:**
- Link graph (all hyperlinks between articles)
- Dependency manifests (import/require declarations)
- Search index (inverted index routing)
- Tag/category associations

**Render Properties:**

| Property          | Value                                            |
| ----------------- | ------------------------------------------------ |
| Material          | Translucent tubes and fiber bundles              |
| Visibility        | Hidden by default; shown in infrastructure view  |
| Color coding      | Per conduit type (see Utility Conduits table)     |
| Animation         | Pulsing glow indicates active data flow          |
| Transparency      | Alpha = 0.0 (hidden) to 0.6 (infrastructure view)|

**Generation Rules:**
1. For every hyperlink `A -> B`, generate a `cross_ref_pipe` from A to B.
2. For every shared tag between articles, generate a `metadata_duct` connecting them through the nearest utility shaft.
3. Search index fibers run from every room to the central core on Floor 0.
4. Dependency wires follow the shortest path through utility shafts, never crossing room interiors.

### Layer 2: Interior Layer

**Purpose:** The actual documentation content. Text, code blocks, images, diagrams, tables --- everything a visitor comes to read.

**Update Frequency:** On content change (article edited, media updated).

**Data Sources:**
- Rendered markdown (HTML output of each article)
- Media assets (images, videos, embedded content)
- Code blocks (syntax-highlighted source)

**Render Properties:**

| Property          | Value                                                  |
| ----------------- | ------------------------------------------------------ |
| Text rendering    | Floating holographic text panels, anti-aliased         |
| Code blocks       | Illuminated wall panels with syntax-colored neon text  |
| Images            | Framed display screens mounted on walls                |
| Tables            | Workbench surfaces with data etched into the surface   |
| Diagrams          | Wall-mounted lightboxes with backlit schematics        |
| Callout boxes     | Freestanding kiosks with colored neon borders          |

**Content-to-Geometry Mapping:**

```
Markdown Element          Spatial Object
-----------------------   ----------------------------------
Paragraph                 Text panel (wall-mounted, 1.2m wide per 80 chars)
Heading (h3+)             Room divider / alcove partition
Code block                Illuminated panel (green neon backlight)
Inline code               Neon-highlighted text segment
Blockquote                Recessed wall alcove with quote etched in glass
Unordered list            Vertical display strip with bullet markers
Ordered list              Numbered display strip (LED counters)
Table                     Workbench surface with grid lines
Image                     Framed display screen
Link                      Door or window (see Door/Window sections)
Horizontal rule           Floor stripe (neon divider line)
Bold text                 Brighter glow intensity (+30%)
Italic text               Slight text rotation (2 degrees)
```

### Layer 3: Signage Layer

**Purpose:** Navigation aids, labels, titles, breadcrumbs, floor indicators, room numbers. Everything that helps a visitor know where they are.

**Update Frequency:** On navigation structure change (reorder, rename, new breadcrumb path).

**Data Sources:**
- Article titles and slugs
- Breadcrumb paths
- Table of contents
- Floor/domain names

**Render Properties:**

| Property          | Value                                              |
| ----------------- | -------------------------------------------------- |
| Material          | Holographic projections and neon-lit physical signs |
| Position          | Above doorways, corridor ceilings, room entrances  |
| Font              | Monospaced, all-caps for titles; proportional for  |
|                   | breadcrumbs                                        |
| Color             | Floor's primary neon color                         |
| Animation         | Gentle hover (0.5cm oscillation, 4s period)        |

**Sign Types:**

| Sign                  | Location             | Content                              |
| --------------------- | -------------------- | ------------------------------------ |
| Room title            | Above room door      | Article title                        |
| Room number           | Door frame, right    | Article coordinate (e.g., `A005`)    |
| Floor indicator       | Elevator interior    | Current floor number and domain name |
| Breadcrumb strip      | Corridor ceiling     | Full navigation path to current room |
| Directory hologram    | Floor lobby          | Interactive map of all rooms on floor|
| You-Are-Here marker   | Various              | Pulsing dot on the nearest map       |
| Exit sign             | Corridor ends        | Points toward lobby / elevator       |
| Warning placard       | Deprecated rooms     | `CAUTION: CONTENT UNDER REVIEW`      |

### Layer 4: Atmospheric Layer

**Purpose:** Ambient environmental effects that communicate status information without requiring direct reading. The building *feels* different based on the health and freshness of its content.

**Update Frequency:** Continuous (recalculated every 60 seconds).

**Data Sources:**
- Content freshness (last modified timestamps)
- Build status (CI/CD pipeline health)
- Error rates (broken links, validation failures)
- Visitor traffic (page view analytics)

**Effect Definitions:**

#### Freshness Glow

Every room emits a glow proportional to its content freshness:

```
Freshness Glow Formula:

  glow_intensity = max(0.1, 1.0 - (age_days / 365))

  age_days = 0     -> intensity 1.0   (brilliant white-blue glow)
  age_days = 30    -> intensity 0.92  (strong glow)
  age_days = 180   -> intensity 0.51  (moderate glow)
  age_days = 365   -> intensity 0.10  (minimum dim glow, never fully dark)
```

| Age Range       | Glow Color    | Hex       | Description                   |
| --------------- | ------------- | --------- | ----------------------------- |
| 0 - 7 days      | Electric blue | `#00D4FF` | Freshly written or updated    |
| 8 - 30 days     | Cyan          | `#00FFD4` | Recent, still current         |
| 31 - 90 days    | Green         | `#00FF6E` | Stable, not stale             |
| 91 - 180 days   | Yellow-green  | `#A8FF00` | Aging, review suggested       |
| 181 - 365 days  | Amber         | `#FFB800` | Stale, review recommended     |
| 365+ days       | Deep red      | `#FF2D00` | Critical staleness, dim glow  |

#### Build Status Lighting

The corridor lighting on each floor reflects the build/CI status of that domain:

```
Build passing   -> Corridor lights steady, floor neon color
Build failing   -> Corridor lights flicker red (#FF003C), 2Hz
Build pending   -> Corridor lights pulse amber (#FFB800), 0.5Hz
Build unknown   -> Corridor lights dim to 30% intensity
```

#### Traffic Heat

High-traffic rooms glow warmer. This is a subtle orange-shift overlay that indicates popular content:

```
Heat Level Calculation:

  heat = log10(visitors_24h + 1) / log10(max_visitors_24h + 1)

  heat < 0.2   -> "cold"     No overlay
  heat < 0.5   -> "warm"     Faint orange tint (alpha 0.05)
  heat < 0.8   -> "hot"      Visible orange tint (alpha 0.15)
  heat >= 0.8  -> "blazing"  Strong orange-red tint (alpha 0.30), heat shimmer FX
```

#### Alert Pulses

Critical issues send visible pulses through the building:

- **Broken link detected:** Red pulse radiates outward from the affected room.
- **Security advisory:** Yellow strobe on the affected floor, all corridors.
- **Deprecation notice:** Slow amber fade-in on the affected room walls.
- **New content published:** Blue ripple expands from the new room outward.

### Layer 5: Holographic Layer

**Purpose:** Ephemeral, user-triggered overlays. Search results, tooltips, preview popups, contextual help. These exist only for the current visitor and are not part of the persistent building state.

**Update Frequency:** On user interaction (real-time).

**Data Sources:**
- Search queries and results
- Hover/focus events
- User preferences and history
- Contextual AI suggestions

**Hologram Types:**

#### Search Result Constellation

When a visitor initiates a search, the results materialize as a constellation of floating cards arranged in relevance order:

```
Search: "authentication flow"

  +---------------------+
  |  * HIC-F07-A003     |  Relevance: 0.97
  |  "Auth Flow Guide"  |  <- Brightest, closest
  +---------------------+
        +---------------------+
        |  . HIC-F07-A001     |  Relevance: 0.84
        |  "Auth Overview"    |
        +---------------------+
              +---------------------+
              |  . HIC-F03-A012     |  Relevance: 0.61
              |  "Security Arch"    |  <- Dimmer, farther
              +---------------------+
```

Each result card:
- Floats at eye level in the current room
- Glows with intensity proportional to relevance score
- Shows a directional arrow pointing toward the target room's physical location
- Can be "grabbed" to navigate instantly to that room

#### Tooltip Projections

When hovering over a door or window, a tooltip hologram appears:

```json
{
  "type": "tooltip",
  "trigger": "hover",
  "target": "door-f03-a005-to-f03-a002",
  "content": {
    "title": "Microservices Guide",
    "summary": "Comprehensive guide to microservice architecture patterns...",
    "word_count": 4200,
    "freshness": "2 days ago",
    "visitors_24h": 156
  },
  "position": "above_trigger",
  "fade_in_ms": 200,
  "fade_out_ms": 150,
  "max_width_m": 2.0
}
```

#### Contextual Navigation Ribbons

Faintly glowing ribbons that trace the path from the current room to recommended next reading:

```
Current: HIC-F03-A005 (Event-Driven Architecture)

  Suggested path ribbons:
  --- cyan ribbon --> HIC-F03-A006 (Message Queues)      "Next in sequence"
  --- gold ribbon --> HIC-F07-A003 (Auth Flow Guide)     "Prerequisite"
  --- pink ribbon --> HIC-F05-A001 (Deployment Overview)  "Related topic"
```

### Layer Compositing and Render Order

Layers are composited bottom-to-top with the following blending rules:

```
Final Pixel = composite(
  Layer 0: Structural     (opaque base)
  Layer 1: Utility        (additive blend, alpha-masked)
  Layer 2: Interior       (alpha-over on structural surfaces)
  Layer 3: Signage        (alpha-over, floating above surfaces)
  Layer 4: Atmospheric    (multiply blend for glow; additive for pulses)
  Layer 5: Holographic    (additive blend, user-local only)
)
```

**Performance Tiers:**

| Tier     | Layers Rendered | Use Case                           |
| -------- | --------------- | ---------------------------------- |
| Minimal  | 0, 3            | Low-bandwidth / text-only clients  |
| Standard | 0, 2, 3         | Default documentation reading      |
| Enhanced | 0, 1, 2, 3, 4   | Full building experience           |
| Complete | 0, 1, 2, 3, 4, 5| Interactive exploration mode        |

---

## Data Encoding

### Transformation Pipeline

The pipeline that converts raw markdown documentation into spatial geometry follows a strict sequence of transformations. Each stage is deterministic: the same input always produces the same output.

```
 +----------+    +----------+    +-----------+    +-----------+    +----------+
 | Markdown |===>|   AST    |===>|  Spatial  |===>| Blueprint |===>| Rendered |
 |  Files   |    |  Parse   |    |  Mapping  |    |   JSON    |    | Geometry |
 +----------+    +----------+    +-----------+    +-----------+    +----------+
   Stage 1         Stage 2         Stage 3          Stage 4          Stage 5
   Source          Structure       Coordinates      Persisted        Visual
```

#### Stage 1: Source Collection

Raw markdown files are collected from the documentation repository. Each file's path determines its floor and article assignment:

```
docs/
|-- constitution/          -> Floor 1
|   |-- 01-preamble.md     -> HIC-F01-A001
|   |-- 02-principles.md   -> HIC-F01-A002
|   +-- 03-governance.md   -> HIC-F01-A003
|-- tutorials/             -> Floor 2
|   |-- 01-getting-started.md -> HIC-F02-A001
|   +-- 02-first-project.md  -> HIC-F02-A002
|-- architecture/          -> Floor 3
|   |-- 01-overview.md     -> HIC-F03-A001
|   +-- ...
```

The directory name maps to a domain (floor). The numeric prefix determines the article index (X position). Files without numeric prefixes are sorted alphabetically and assigned sequential indices.

#### Stage 2: AST Parsing

Each markdown file is parsed into an Abstract Syntax Tree. The AST captures the hierarchical structure of the content:

```json
{
  "type": "document",
  "source": "docs/architecture/01-overview.md",
  "children": [
    {
      "type": "heading",
      "depth": 1,
      "text": "Architecture Overview",
      "children": [
        {
          "type": "heading",
          "depth": 2,
          "text": "System Components",
          "children": [
            {
              "type": "heading",
              "depth": 3,
              "text": "Frontend Layer",
              "children": [
                { "type": "paragraph", "text": "The frontend..." },
                { "type": "code_block", "language": "typescript", "text": "..." }
              ]
            },
            {
              "type": "heading",
              "depth": 3,
              "text": "Backend Services",
              "children": [
                { "type": "paragraph", "text": "Backend services..." }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

#### Stage 3: Spatial Mapping

The AST is transformed into spatial coordinates. This is where the document structure becomes physical geometry:

```
AST Node                  Spatial Result
------------------------  ----------------------------------------
document                  Room envelope (4m x 6m x 4m)
heading depth=1           Floor assignment (already determined by directory)
heading depth=2           Wing partition (divides room into major sections)
heading depth=3           Room partition (alcove or sub-room)
heading depth=4           Alcove niche (recessed area within a section)
paragraph                 Text panel (wall-mounted, depth-positioned)
code_block                Illuminated panel (recessed into wall)
image                     Display screen (wall-mounted frame)
table                     Workbench surface (horizontal display)
list                      Display strip (vertical wall mount)
blockquote                Glass alcove (recessed, quote-style)
link                      Door or window (see Door Connections)
horizontal_rule           Floor neon strip (section divider)
```

#### Stage 4: Blueprint Generation

The spatial mapping is serialized into blueprint JSON (as defined in the Blueprint Format Specification). This is the canonical persisted form of the building's geometry.

#### Stage 5: Render

The blueprint JSON is consumed by the rendering engine, which generates the final visual output across all six layers. The renderer applies materials, lighting, atmospheric effects, and user-specific holographic overlays.

### Heading Levels as Spatial Partitions

The mapping of heading levels to spatial divisions is one of the most critical rules in the transformation pipeline:

```
# Heading 1 (h1)     =   FLOOR
                          The entire floor plate. One h1 per document.
                          Defines the room's top-level identity.

## Heading 2 (h2)    =   WING
                          A major section of the room, divided by
                          a translucent partition wall. Visitors can
                          see through to adjacent wings.

### Heading 3 (h3)   =   ROOM / ALCOVE
                          A distinct sub-space within a wing.
                          Has its own lighting and can be
                          individually addressed.

#### Heading 4 (h4)  =   NICHE
                          A recessed area within an alcove.
                          Typically holds a single focused piece
                          of content (a code example, a table,
                          a specific note).

##### Heading 5 (h5) =   SHELF
                          A labeled shelf or drawer within a niche.
                          Very fine-grained content subdivision.

###### Heading 6 (h6)=   LABEL
                          A physical label on a shelf or drawer.
                          The smallest addressable unit.
```

#### Visual Representation

```
+------------------------------------------------------------------+
|  h1: Architecture Overview                          (ROOM)       |
|                                                                  |
|  +-------------------------+  +-----------------------------+    |
|  | h2: System Components   |  | h2: Design Principles       |    |
|  |         (WING)          |  |         (WING)               |    |
|  |                         |  |                               |    |
|  |  +-------+  +-------+  |  |  +-----------------------+   |    |
|  |  | h3:   |  | h3:   |  |  |  | h3: SOLID Principles  |   |    |
|  |  |Front  |  |Back   |  |  |  |       (ALCOVE)         |   |    |
|  |  |end    |  |end    |  |  |  |  +-----+ +-----+      |   |    |
|  |  |       |  |       |  |  |  |  |h4:  | |h4:  |      |   |    |
|  |  |       |  |       |  |  |  |  |SRP  | |OCP  |      |   |    |
|  |  |       |  |       |  |  |  |  |     | |     |      |   |    |
|  |  +-------+  +-------+  |  |  |  +-----+ +-----+      |   |    |
|  |                         |  |  +-----------------------+   |    |
|  +-------------------------+  +-----------------------------+    |
|                                                                  |
+------------------------------------------------------------------+
```

### Content Density Calculations

Content density determines how "full" a room feels. It affects ambient lighting, echo effects, and visual clutter levels.

#### Formula

```
density = word_count / floor_area_m2

Where:
  word_count   = total words in the article
  floor_area   = room width (4m) x room depth (6m) = 24 m2
```

#### Density Classifications

| Density (words/m2) | Classification | Visual Effect                              |
| ------------------- | -------------- | ------------------------------------------ |
| 0 - 25              | Sparse         | Echoing, empty feel. Minimal wall content. |
| 26 - 75             | Light          | Comfortable spacing. Clean, airy room.     |
| 76 - 150            | Standard       | Well-furnished. Balanced text and space.   |
| 151 - 300           | Dense          | Crowded bookshelves. Walls fully covered.  |
| 301 - 500           | Packed         | Stacked displays. Requires scrolling.      |
| 500+                | Overloaded     | Warning: room should be split.             |

#### Overload Response

When density exceeds 500 words/m2, the spatial engine triggers an **overload warning**:

1. The room's walls begin to "bulge" outward (subtle geometry distortion).
2. A yellow caution hologram appears at the door: `CONTENT DENSITY CRITICAL`.
3. The room's neon trim shifts from the floor color to amber.
4. A suggested split point is calculated (the h2 heading nearest the midpoint of the content) and a dashed line appears on the floor indicating where the room could be divided.

### Media Object Placement

Non-text content elements have specific placement rules:

| Media Type        | Placement                    | Size Calculation                        |
| ----------------- | ---------------------------- | --------------------------------------- |
| Inline image      | Wall-mounted frame           | Aspect ratio preserved, max 2m wide     |
| Block image       | Freestanding display easel   | Aspect ratio preserved, max 3m wide     |
| Video embed       | Recessed wall screen         | 16:9, width = min(3m, room width - 1m)  |
| Code block        | Illuminated wall panel       | Width = room width - 0.5m, height varies|
| Table (small)     | Workbench surface            | 1 row = 0.15m height, cols scale to fit |
| Table (large)     | Floor-mounted display table  | Scrollable, max 3m x 2m                 |
| Mermaid diagram   | Wall-mounted lightbox        | Auto-scaled to fit 2m x 2m frame        |
| Math equation     | Floating holographic display | Centered in room, 0.5m above eye level  |

---

## Scale and Proportions

### Standard Dimensions

The HIC uses a fixed set of standard dimensions. These are not arbitrary --- they are tuned to produce readable, navigable spaces at typical documentation scales.

#### Core Measurements

| Element               | Dimension              | Notes                                    |
| --------------------- | ---------------------- | ---------------------------------------- |
| Standard room width   | 4 meters               | One article = one room                   |
| Standard room depth   | 6 meters               | Five depth zones (Z=0 to Z=4)           |
| Floor height          | 4 meters               | 3.5m clear + 0.5m structural slab       |
| Corridor width        | 2 meters               | Navigation / table of contents           |
| Lobby width           | 3x standard room width | 12 meters minimum                        |
| Elevator shaft        | 3m x 3m                | Vertical navigation between floors       |
| Stairwell             | 2m x 4m                | Sequential floor-to-floor navigation     |
| Utility shaft         | 1m x 1m                | Conduit routing between floors           |
| Door width            | 1.2 meters             | Standard passage                         |
| Door height           | 2.5 meters             | Taller than standard for dramatic effect |
| Window width          | 1.5 meters             | Preview panel                            |
| Window height         | 1.0 meter              | Summary viewport                         |

#### Derived Measurements

```
Floor footprint for N articles:

  width  = (N x 4m) + ((N - 1) x 2m) + 12m + (2 x 3m)
           rooms      corridors        lobby   elevators

  depth  = 6m (constant --- all rooms same depth)

Example: Floor with 15 articles
  width  = (15 x 4) + (14 x 2) + 12 + 6
         = 60 + 28 + 12 + 6
         = 106 meters

  Total floor area = 106 x 6 = 636 m2
```

### Dynamic Scaling Rules

While standard dimensions are fixed, the building's overall envelope is dynamic. It grows and shifts as documentation is added or removed.

#### Floor Width Scaling

Each floor's width is independently determined by the number of articles in that domain. This means the building is not a perfect rectangular prism --- it is an irregular tower, wider at floors with more content and narrower at floors with less.

```
Cross-section (not to scale):

         +------+           Floor 8: 5 articles (narrow)
       +-+------+-+         Floor 7: 10 articles
    +--+----------+--+      Floor 6: 20 articles (wide)
    +----------------+      Floor 5: 18 articles
    +----------------+      Floor 4: 19 articles
   ++----------------++     Floor 3: 22 articles (widest)
   +------------------+     Floor 2: 20 articles
   |==================|     Floor 1: 15 articles (Constitution)
   +------------------+     Floor 0: Sub-basement (full width)
```

The sub-basement (Floor 0) always extends to the maximum width of any floor above it, providing a stable foundation.

#### Lobby Scaling

The lobby on each floor scales with the total number of rooms on that floor:

```
Lobby Width Formula:

  lobby_width = max(12m, 3 x sqrt(num_articles) meters)

  5 articles   -> max(12, 3 x 2.24)  = 12.0m  (minimum)
  15 articles  -> max(12, 3 x 3.87)  = 12.0m  (minimum)
  25 articles  -> max(12, 3 x 5.0)   = 15.0m
  100 articles -> max(12, 3 x 10.0)  = 30.0m
```

Large lobbies feature additional amenities:
- Over 20m: Interactive holographic floor map
- Over 25m: Seating areas (visitor rest zones, conceptually: "reading lists")
- Over 30m: Sub-lobbies with wing indicators (domain sub-categories)

### Vertical Proportions

#### Building Height

```
Total Height Formula:

  height = (num_floors x 4m) + foundation_depth

  Where foundation_depth = 8m (sub-basement + structural foundation)

  10 domains  -> (11 x 4) + 8 = 52 meters
  25 domains  -> (26 x 4) + 8 = 112 meters
  50 domains  -> (51 x 4) + 8 = 212 meters
  100 domains -> (101 x 4) + 8 = 412 meters
```

#### Rooftop

The highest occupied floor is topped by:

- A 2m structural cap (weather protection --- metaphorically, the system boundary)
- A rooftop observation deck (the global search interface and site map)
- A neon spire (the HIC identifier beacon, visible from any distance)

The spire's height scales logarithmically with total building height:

```
spire_height = 10 + (5 x log2(total_floors))

  10 floors -> 10 + (5 x 3.32) = 26.6m spire
  50 floors -> 10 + (5 x 5.64) = 38.2m spire
  100 floors -> 10 + (5 x 6.64) = 43.2m spire
```

#### Aspect Ratio Constraints

To prevent the building from becoming absurdly tall and thin (or short and wide), the following soft constraints apply:

```
Height-to-Width Ratio:

  Recommended:  2:1 to 8:1
  Acceptable:   1:1 to 12:1
  Warning:      Outside acceptable range triggers architectural review

If ratio exceeds 12:1:
  -> Suggest merging low-content domains (combining floors)

If ratio falls below 1:1:
  -> Suggest splitting high-content domains (adding floors)
```

---

## Collision Detection

In a building generated from data, conflicts are inevitable. Two articles might claim the same address. A redirect might point into a wall. A namespace change might leave two rooms trying to occupy the same coordinates. The collision detection system prevents these structural failures.

### Namespace Collisions

A namespace collision occurs when two distinct content items attempt to occupy the same coordinate. This is the spatial equivalent of two files trying to have the same path.

#### Detection

```
Collision Check Algorithm:

  For each (floor, article_index) pair in the blueprint:
    If count(rooms at this coordinate) > 1:
      COLLISION DETECTED
      -> Flag all rooms at this coordinate
      -> Halt rendering for affected coordinate
      -> Emit collision event
```

#### Resolution Strategies

| Strategy          | When Used                        | Action                                    |
| ----------------- | -------------------------------- | ----------------------------------------- |
| Sequence Bump     | Duplicate article index          | Later-added article gets index + 1        |
| Floor Reassignment| Duplicate domain mapping         | Newer domain gets next available floor    |
| Merge             | Intentional consolidation        | Two rooms become one, content combined    |
| Tombstone         | Irreconcilable conflict          | One room is tombstoned (emptied, marked)  |

#### Collision Event Schema

```json
{
  "event": "collision_detected",
  "timestamp": "2026-02-17T03:22:00Z",
  "type": "namespace",
  "coordinate": "HIC-F03-A012",
  "conflicting_sources": [
    {
      "file": "docs/architecture/12-caching.md",
      "hash": "a1b2c3...",
      "created": "2025-06-01T00:00:00Z"
    },
    {
      "file": "docs/architecture/12-caching-v2.md",
      "hash": "d4e5f6...",
      "created": "2026-02-15T00:00:00Z"
    }
  ],
  "resolution": "sequence_bump",
  "resolved_coordinates": {
    "docs/architecture/12-caching.md": "HIC-F03-A012",
    "docs/architecture/12-caching-v2.md": "HIC-F03-A013"
  }
}
```

### Anchor Uniqueness Enforcement

Within a single room, every addressable point (anchor) must have a unique identifier. This is the "no two rooms share an address" rule, applied at the intra-room level: no two sections within a room share an anchor ID.

#### The Problem

Markdown heading-to-anchor generation can produce duplicates:

```markdown
## Configuration
...content...

## Configuration          <-- Duplicate! Same anchor ID generated.
...different content...
```

#### The Solution: Anchor Deduplication

```
Anchor Generation Algorithm:

  1. Convert heading text to slug: lowercase, replace spaces with hyphens,
     strip non-alphanumeric characters.
  2. Check slug against existing anchors in this room.
  3. If collision:
     a. Append "-{n}" where n is the lowest integer that resolves the collision.
     b. Log a warning: "Anchor collision resolved: {slug} -> {slug}-{n}"
  4. Register the anchor in the room's anchor registry.

Example:
  "Configuration"  -> "configuration"
  "Configuration"  -> "configuration-2"    (collision resolved)
  "Configuration"  -> "configuration-3"    (second collision resolved)
```

#### Spatial Manifestation

Deduplicated anchors manifest as slightly different door markings. The original anchor's door has a clean label. The deduplicated door has a small subscript number on its placard, and a faint connecting line to the original door indicating they share a name lineage.

### Redirect Corridors

When content moves from one coordinate to another (a URL redirect), the old location does not simply vanish. Instead, a redirect corridor is constructed --- a physical passageway from the old coordinate to the new one.

#### Redirect Types

| HTTP Status | Building Analog                          | Visual Treatment                           |
| ----------- | ---------------------------------------- | ------------------------------------------ |
| 301         | Permanent renovation (room moved)        | Sealed old door, clear tunnel to new room  |
| 302         | Temporary relocation (room under repair) | Open old door, caution tape, detour signs  |
| 307         | Temporary redirect (preserve method)     | Open old door, holographic "follow me" arrow|
| 308         | Permanent redirect (preserve method)     | Sealed old door, neon arrow to new location|
| 410         | Content permanently gone                 | Bricked-over doorway, memorial placard     |

#### Redirect Corridor Schema

```json
{
  "id": "redirect-f03-a012-to-f03-a025",
  "type": "301",
  "source_coordinate": "HIC-F03-A012",
  "target_coordinate": "HIC-F03-A025",
  "created": "2026-01-15T00:00:00Z",
  "reason": "Content restructured: caching guide split into strategy and implementation",
  "corridor": {
    "length_m": 52.0,
    "width_m": 1.5,
    "lighting": "dim_amber",
    "signage": [
      {
        "position": 0.0,
        "text": "THIS ROOM HAS MOVED -> A025",
        "style": "neon_arrow"
      },
      {
        "position": 0.5,
        "text": "Redirecting...",
        "style": "holographic_pulse"
      },
      {
        "position": 1.0,
        "text": "ARRIVING: Caching Strategy Guide",
        "style": "neon_sign"
      }
    ]
  },
  "expiry": null,
  "visits_since_creation": 342
}
```

#### Redirect Corridor Lifecycle

1. **Creation:** When a redirect is established, a corridor is instantiated immediately.
2. **Active use:** Visitors entering the old coordinate are automatically routed through the corridor. The transit is seamless but visible --- the visitor sees the redirect happening.
3. **Aging:** After 180 days with no visitors, a 301 redirect corridor's lighting dims further and a "HISTORICAL ROUTE" sign appears.
4. **Compaction:** After 365 days with no visitors, the corridor is eligible for compaction. The old room coordinate is released back to the void pool. The corridor is recorded in demolition records and removed from the active blueprint.

### Overlap Resolution Algorithms

Beyond simple namespace collisions, overlaps can occur when dynamic scaling causes one room's geometry to intersect another's. This happens when:

- A room's content density causes visual "bulging"
- A lobby scales up and encroaches on adjacent rooms
- A utility conduit's routing conflicts with a room boundary

#### Overlap Detection

```
For each pair of rooms (A, B) on the same floor:
  If bounding_box(A) INTERSECT bounding_box(B) != EMPTY:
    overlap_area = area(bounding_box(A) INTERSECT bounding_box(B))
    If overlap_area > 0:
      OVERLAP DETECTED
```

#### Resolution Priority

When an overlap is detected, resolution follows a priority hierarchy:

1. **Lobby wins over room.** Lobbies are never compressed; adjacent rooms shift outward.
2. **Older room wins over newer room.** The room with the earlier creation date holds its position; the newer room shifts.
3. **Higher traffic wins.** If creation dates are equal, the room with more visitors in the past 30 days holds position.
4. **In case of true tie:** Both rooms shift outward by half the overlap distance, and the corridor between them widens to fill the gap.

---

## Cartographic Standards

The HIC can be viewed from multiple perspectives. Each perspective is a map projection --- a way of flattening or framing the three-dimensional building into a two-dimensional view suitable for navigation and reference.

### Map Projections

#### Top-Down Projection (Site Map)

The top-down view collapses the Y-axis (floors) and shows each floor as a separate 2D plan. This is equivalent to a traditional site map.

```
Top-Down: Floor 3 (Architecture)

  +============+
  |   LOBBY    |--[A002]--[A003]--[A004]--[A005]--[A006]--...--[A015]
  |   A001     |
  +============+
       |
    [ELEVATOR]
```

**Properties:**
- Each floor rendered independently
- Room sizes proportional to content density
- Doors shown as gaps in walls
- Corridors shown as connecting lines
- Color indicates freshness (atmospheric layer applied)

**Use cases:** Finding a specific room on a known floor. Understanding floor layout.

#### Isometric Projection (Building View)

The isometric view shows the full building as a 3D structure viewed from a 30-degree elevation angle. All three axes are visible.

```
Isometric: HIC Building (simplified)

              /\
             /  \  <- Spire
            /    \
           +------+
          /| F08  |\
         / +------+ \
        /  | F07  |  \
       /   +------+   \
      /    | F06  |    \
     /     +------+     \
    /      | F05  |      \
   /       +------+       \
  /        | F04  |        \
 /         +------+         \
/          | F03  |          \
           +------+
           | F02  |
           |======|
           | F01  | <- Constitution
           +------+
           ########  Foundation
```

**Properties:**
- All floors visible simultaneously
- Floor widths reflect actual article counts (building is not a perfect box)
- Each floor's neon accent color is visible on its exterior edge
- Atmospheric glow visible as a per-floor halo
- Animated: conduit pulses visible on the exterior

**Use cases:** Understanding overall documentation structure. Identifying large vs. small domains. Spotting freshness patterns across the building.

#### Cross-Section Projection (Domain Deep-Dive)

The cross-section cuts through the building along the Z-axis, showing depth layers within a single floor.

```
Cross-Section: Floor 3, Article 5

  +-------------------------------------------+
  |                  CEILING                   |
  |                                            |
  |  Z=0  |  Title Card    |  status light    |
  |  -----|----------------|------------------|
  |  Z=1  |  Introduction  |  prerequisites   |
  |  -----|----------------|------------------|
  |  Z=2  |  Main Content  |  code panels     |
  |       |  diagrams      |  examples        |
  |  -----|----------------|------------------|
  |  Z=3  |  Edge Cases    |  perf notes      |
  |  -----|----------------|------------------|
  |  Z=4  |  Raw Data      |  source refs     |
  |                                            |
  |                  FLOOR                     |
  +-------------------------------------------+
         Door                    Back Wall
       (Corridor)               (Deepest)
```

**Properties:**
- Single room, full depth visible
- Each Z-layer labeled and color-coded
- Content density visible as fill percentage per layer
- Utility conduits visible passing through walls (but not content zone)

**Use cases:** Understanding a single article's depth structure. Identifying which depth layers have content. Planning content depth strategy.

#### Elevation Projection (Facade View)

The elevation shows one face of the building --- typically the south face (the "front" facing visitors).

```
Elevation: South Face

  +---------------------------------------------+
  | F08 | . . . . .                              | 5 rooms
  +-----+-------------------------------------------+
  | F07 | . . . . . . . . . .                   | 10 rooms
  +-----+-------------------------------------------+
  | F06 | . . . . . . . . . . . . . . . . . .   | 18 rooms
  +-----+-------------------------------------------+
  | F05 | . . . . . . . . . . . . . . . . .     | 17 rooms
  +-----+-------------------------------------------+
  | F04 | . . . . . . . . . . . . . . . .       | 16 rooms
  +-----+-------------------------------------------+
  | F03 | . . . . . . . . . . . . . . . . . .   | 18 rooms
  +-----+-------------------------------------------+
  | F02 | . . . . . . . . . . . .               | 12 rooms
  +=====+===========================================+
  | F01 | <> CONSTITUTION <> . . . . . . . .     | 10 rooms
  +-----+-------------------------------------------+

  .  = room window (lit according to freshness)
  <> = main entrance
```

**Properties:**
- Every room visible as a lit window on the facade
- Window color indicates freshness (atmospheric layer)
- Window brightness indicates traffic (heat level)
- Floor labels visible on the left edge
- The Constitution floor has a distinctive border (double line)

### Legend and Key Standards

All HIC maps must include a standardized legend. The legend is divided into sections:

#### Structural Symbols

```
  ===    Double line:  Load-bearing wall (domain boundary)
  ---    Single line:  Partition wall (section boundary)
  - -    Dashed line:  Suggested split point
  / \    Diagonal:     Roof/spire structure
  ###    Filled:       Foundation / sub-basement
```

#### Room Indicators

```
  .      Solid dot:         Populated room (has content)
  o      Open circle:       Reserved room (empty, placeholder)
  x      X mark:            Tombstoned room (content removed)
  <>     Diamond:           Special room (lobby, entrance)
  [~]    Hourglass:         Room under construction (draft content)
```

#### Conduit Markers

```
  ~~~~   Wavy line:         Data flow conduit
  ....   Dotted line:       Metadata duct
  -=-    Dashed double:     Dependency wire
  -o-    Line with dot:     Junction / hub point
```

#### Status Colors

```
  #00D4FF  Electric Blue:    Fresh content (0-7 days)
  #00FFD4  Cyan:             Recent content (8-30 days)
  #00FF6E  Green:            Stable content (31-90 days)
  #A8FF00  Yellow-green:     Aging content (91-180 days)
  #FFB800  Amber:            Stale content (181-365 days)
  #FF2D00  Deep Red:         Critical staleness (365+ days)
  #FF003C  Alert Red:        Build failure / broken link
  #8A2BE2  Purple:           Analytics data flow
  #FF6EC7  Pink:             Cross-reference link
```

### Scale Bars and Measurement Units

Every HIC map includes a scale bar calibrated to the standard room unit:

```
Scale Bar (Standard):

  |<-- 4m -->|<-- 4m -->|<-- 4m -->|
  +-----------+-----------+-----------+
  0     1 room      2 rooms     3 rooms

  1 room = 4m x 6m = 24 m2 floor area
```

#### Measurement Units

| Unit              | Symbol | Definition                            |
| ----------------- | ------ | ------------------------------------- |
| Room Unit (RU)    | `ru`   | 1 standard room = 4m x 6m            |
| Floor Unit (FU)   | `fu`   | 1 floor = 4m vertical height          |
| Corridor Unit     | `cu`   | 1 corridor segment = 2m x 6m         |
| Depth Unit (DU)   | `du`   | 1 depth layer = 1.2m (6m / 5 layers) |

Maps are annotated with both metric (meters) and HIC-native (room units) measurements. The dual notation ensures readability for both spatial engineers and content architects.

### Orientation and the Constitutional Compass

Every map includes an orientation indicator. In the HIC, there is no magnetic north. Instead, the Constitution serves as the fixed reference point.

#### The Constitutional Compass

```
                    ^
                    |
              CONSTITUTION
              (Floor 1, Up)
                    |
  OLDER <-----------+-----------> NEWER
  ARTICLES          |            ARTICLES
  (Lower Index)     |            (Higher Index)
                    |
                DEEPEST
              (Floor Max, Down)
                    v
```

**Rules:**
- **Up** always points toward the Constitution (Floor 1). Even though Floor 1 is at the bottom of the building, it is "up" in the cartographic sense --- it is the origin, the source, the foundation from which all else derives.
- **Right** always points toward higher article indices (newer/later content).
- **Left** always points toward lower article indices (older/earlier content).
- **Down** always points toward the highest floor number (the newest domain, the top of the physical building).

This inversion --- the Constitution being physically at the bottom but cartographically "up" --- is intentional. It reinforces the governance model: everything flows from the Constitution outward.

---

## Versioning and Time

The HIC is not a static structure. It is continuously renovated, expanded, and occasionally partially demolished. The versioning system tracks every change as a construction event.

### Historical Blueprints

Every commit to the documentation repository generates a new building snapshot. These snapshots are stored as historical blueprints and can be replayed to show the building at any point in its history.

#### Snapshot Schema

```json
{
  "snapshot_id": "hic-snap-20260217-a4f8c3de",
  "timestamp": "2026-02-17T14:32:00Z",
  "commit_hash": "a4f8c3de7b9012345678abcdef",
  "commit_message": "Add caching strategy guide to Architecture domain",
  "author": "kael.holm@hic.internal",
  "building_stats": {
    "total_floors": 8,
    "total_rooms": 127,
    "total_doors": 342,
    "total_windows": 89,
    "total_conduits": 1205,
    "total_words": 284000,
    "building_height_m": 40,
    "max_floor_width_m": 106,
    "avg_content_density": 98.3
  },
  "changes": [
    {
      "type": "room_added",
      "coordinate": "HIC-F03-A016",
      "title": "Caching Strategy Guide",
      "word_count": 2400
    }
  ],
  "blueprint_ref": "blueprints/snapshots/20260217-a4f8c3de/"
}
```

#### Historical Replay

The building can be "rewound" to any historical snapshot. The replay system:

1. Loads the blueprint from the target snapshot.
2. Renders the building in its historical state.
3. Applies a sepia-toned atmospheric overlay to indicate "historical view."
4. Displays a timeline bar at the top of the viewport showing the current position in history.
5. Allows scrubbing forward and backward through commits.

Visual treatment for historical replay:

```
Active (current) building:    Full neon color, all layers active
Historical building:          Muted palette, sepia overlay (alpha 0.3)
                              Holographic layer disabled
                              Atmospheric layer frozen (no live updates)
                              Timestamp watermark: "HISTORICAL: 2025-08-14"
```

### Construction Timeline

The construction timeline is a linear visualization of the building's growth over time. Each commit is a construction event.

#### Event Types

| Event Type          | Git Operation       | Visual                                          |
| ------------------- | ------------------- | ----------------------------------------------- |
| Foundation Pour     | Initial commit      | The ground breaks. First floor appears.         |
| New Floor           | New domain added    | A new floor materializes atop the building.     |
| Room Construction   | New article added   | Scaffolding appears, then walls, then content.  |
| Renovation          | Article edited      | Scaffolding wraps the room, then reveals update.|
| Room Expansion      | Major content add   | Walls shift outward, room grows.                |
| Corridor Extension  | Nav restructured    | New corridor segments appear.                   |
| Demolition          | Article deleted     | Controlled demolition: room implodes to void.   |
| Floor Demolition    | Domain removed      | Entire floor removed; upper floors descend.     |
| Relocation          | Article moved       | Room lifts off, travels to new coordinate.      |

#### Timeline Rendering

```
Construction Timeline: HIC Building

  2025-01-15  === FOUNDATION POUR ===
              Floor 1: Constitution (3 articles)

  2025-02-01  | Floor 2: Tutorials added (5 articles)
              |
  2025-02-14  | Floor 1: +2 articles (total: 5)
              |
  2025-03-01  | Floor 3: Architecture added (8 articles)
              |
  2025-03-15  | Floor 3: +4 articles (total: 12)
              | Floor 2: Renovation on A003
              |
  2025-04-01  | Floor 4: API Reference added (20 articles)
              |
  ...         |
              |
  2026-02-17  | Floor 3: +1 article (A016: Caching Strategy)
              | Floor 7: Renovation on A003, A005
              | Total: 8 floors, 127 rooms
              v NOW
```

### Renovation Logs

When content is edited, the change is logged as a renovation. The renovation log includes before-and-after snapshots, making it possible to visualize what changed in a room.

#### Renovation Record Schema

```json
{
  "renovation_id": "reno-20260217-f03-a005",
  "coordinate": "HIC-F03-A005",
  "title": "Event-Driven Architecture",
  "timestamp": "2026-02-17T14:32:00Z",
  "commit_hash": "a4f8c3de...",
  "author": "kael.holm@hic.internal",
  "before": {
    "word_count": 2800,
    "heading_partitions": 10,
    "content_density": 116.7,
    "content_hash": "old_hash..."
  },
  "after": {
    "word_count": 3200,
    "heading_partitions": 12,
    "content_density": 133.3,
    "content_hash": "new_hash..."
  },
  "diff_summary": {
    "lines_added": 45,
    "lines_removed": 12,
    "net_change": "+33 lines",
    "sections_added": ["Error Handling", "Retry Patterns"],
    "sections_removed": [],
    "sections_modified": ["Message Formats"]
  },
  "spatial_impact": {
    "room_size_changed": false,
    "partitions_added": 2,
    "partitions_removed": 0,
    "doors_added": 1,
    "doors_removed": 0,
    "density_change": "+16.6 words/m2"
  }
}
```

#### Diff Visualization

Renovations can be viewed as a split-screen comparison:

```
+------ BEFORE ------++------- AFTER -------+
|                     ||                     |
|  HIC-F03-A005      ||  HIC-F03-A005       |
|  Event-Driven Arch  ||  Event-Driven Arch  |
|                     ||                     |
|  +--------------+   ||  +--------------+   |
|  | Message      |   ||  | Message  [*] |   |
|  | Formats      |   ||  | Formats      |   |
|  +--------------+   ||  +--------------+   |
|                     ||  +--------------+   |
|     (empty)         ||  | Error    [+] |   |
|                     ||  | Handling     |   |
|                     ||  +--------------+   |
|                     ||  +--------------+   |
|     (empty)         ||  | Retry    [+] |   |
|                     ||  | Patterns     |   |
|                     ||  +--------------+   |
|                     ||                     |
|  Density: 116.7     ||  Density: 133.3     |
|  Partitions: 10     ||  Partitions: 12     |
|  Doors: 4           ||  Doors: 5           |
|                     ||                     |
+---------------------++---------------------+

  [*] = Modified section       [+] = New section
```

### Demolition Records

When content is deleted, it is not simply erased from the building. A demolition record is created, preserving the memory of what once stood at that coordinate.

#### Demolition Types

| Type              | Trigger              | Spatial Effect                                     |
| ----------------- | -------------------- | -------------------------------------------------- |
| Room Demolition   | Article deleted      | Room walls collapse inward, leaving an empty lot    |
| Floor Demolition  | Domain removed       | Entire floor removed, upper floors descend by 4m   |
| Partial Demolition| Major section removed| Room shrinks, freed space becomes void              |
| Condemnation      | Content deprecated   | Room remains but is cordoned off, amber warning     |

#### Demolition Record Schema

```json
{
  "demolition_id": "demo-20260215-f03-a020",
  "type": "room_demolition",
  "coordinate": "HIC-F03-A020",
  "title": "Legacy Caching (Deprecated)",
  "timestamp": "2026-02-15T10:00:00Z",
  "commit_hash": "b7e2f1aa...",
  "author": "maint-bot@hic.internal",
  "reason": "Replaced by A016 (Caching Strategy Guide)",
  "redirect": {
    "type": "301",
    "target": "HIC-F03-A016"
  },
  "final_snapshot": {
    "word_count": 1800,
    "heading_partitions": 6,
    "content_hash": "final_hash...",
    "last_visitor": "2026-02-14T23:59:00Z",
    "total_lifetime_visitors": 4521,
    "age_days": 245
  },
  "memorial": {
    "placard_text": "HERE STOOD: Legacy Caching Guide (2025-06-15 to 2026-02-15)",
    "placard_coordinate": "HIC-F03-A020-D0",
    "visible": true,
    "style": "brass_plate"
  }
}
```

#### The Empty Lot

After demolition, the coordinate remains in the building as an empty lot:

- The floor plate is intact but the walls are gone.
- A faint outline (ghosted geometry) shows where the room used to be.
- A brass memorial placard is mounted on a small standing post at the front of the lot.
- The lot emits no light. It is the darkest point on the floor.
- Over time (configurable, default 90 days), the empty lot may be reclaimed by adjacent rooms expanding, or by new content taking the coordinate.

#### The Ruin View

Activating "ruin view" in historical mode renders demolished rooms as crumbling ruins rather than empty lots. The walls are partially standing, shattered glass panes glow faintly with the ghost of their former content, and the room's neon accent flickers erratically before going dark. This is a visualization mode only --- it has no functional impact.

---

## Appendices

### Appendix A: Complete Coordinate Grammar (EBNF)

```ebnf
coordinate     = prefix , "-" , floor , "-" , article , "-" , depth ;
prefix         = "HIC" ;
floor          = "F" , floor_number ;
article        = "A" , article_number ;
depth          = "D" , depth_number ;

floor_number   = digit , digit , { digit } ;   (* 2+ digits, zero-padded *)
article_number = digit , digit , digit , { digit } ;  (* 3+ digits, zero-padded *)
depth_number   = "0" | "1" | "2" | "3" | "4" ;

shorthand      = short_floor , ":" , short_article , [ ":" , short_depth ] ;
short_floor    = "F" , nonzero_digit , { digit } ;
short_article  = "A" , nonzero_digit , { digit } | "*" ;
short_depth    = "D" , depth_number | "*" ;

range          = coordinate_or_short , ".." , coordinate_or_short ;

digit          = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" ;
nonzero_digit  = "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" ;
```

### Appendix B: Default Neon Color Palette by Floor Type

| Floor Purpose       | Primary Neon  | Secondary     | Accent        |
| ------------------- | ------------- | ------------- | ------------- |
| Constitution        | `#FFD700`     | `#3D2E00`     | `#FFFFFF`     |
| Tutorials           | `#00FF85`     | `#003D20`     | `#80FFB2`     |
| Architecture        | `#00FFA3`     | `#003D28`     | `#80FFD1`     |
| API Reference       | `#00BFFF`     | `#002E3D`     | `#80DFFF`     |
| Guides              | `#FF6EC7`     | `#3D1A2F`     | `#FFB6E3`     |
| Deployment          | `#FF3131`     | `#3D0C0C`     | `#FF9898`     |
| Security            | `#8A2BE2`     | `#21093D`     | `#C495F0`     |
| Operations          | `#FF8C00`     | `#3D2200`     | `#FFC580`     |
| Contributing        | `#00CED1`     | `#00323D`     | `#80E6E8`     |
| Changelog           | `#B0B0B0`     | `#2A2A2A`     | `#D8D8D8`     |

### Appendix C: Rendering Performance Budgets

| Layer       | Max Render Time | Max Memory | Update Trigger           |
| ----------- | --------------- | ---------- | ------------------------ |
| Structural  | 200ms           | 50MB       | Documentation tree change|
| Utility     | 150ms           | 30MB       | Link graph change        |
| Interior    | 500ms           | 200MB      | Content change           |
| Signage     | 100ms           | 20MB       | Navigation change        |
| Atmospheric | 50ms/frame      | 10MB       | Continuous (60fps)       |
| Holographic | 30ms/frame      | 15MB       | User interaction         |

### Appendix D: Glossary of Spatial Terms

| Term                | Definition                                                   |
| ------------------- | ------------------------------------------------------------ |
| **Alcove**          | A sub-division of a room, created by an h4 heading           |
| **Conduit**         | A metadata or data pathway running through building structure|
| **Coordinate**      | A three-part address identifying a point in the HIC          |
| **Density**         | Words per square meter of virtual floor space                |
| **Empty Lot**       | A coordinate where content was demolished                    |
| **Freshness Glow**  | The ambient light a room emits based on content age          |
| **Ghost Geometry**  | The faint outline of a demolished room                       |
| **Heat Level**      | Traffic-based intensity overlay on a room                    |
| **Lobby**           | The entrance room of a floor (domain landing page)           |
| **Niche**           | The smallest partitioned space, from an h4 heading           |
| **Populated**       | A coordinate that contains rendered content                  |
| **Redirect Corridor**| A passageway from an old coordinate to a new one            |
| **Reserved**        | A coordinate on an existing floor with no content yet        |
| **Room**            | The spatial unit corresponding to one article                |
| **Ruin View**       | Visual mode showing demolished rooms as crumbling structures |
| **Snapshot**        | A complete building state at a specific commit               |
| **Tombstone**       | A permanent marker at a demolished room's location           |
| **Void**            | A coordinate where no floor exists                           |
| **Wing**            | A major section of a room, created by an h2 heading          |

---

> *"The blueprints are never finished. The last line of this specification
> is a wall that hasn't been built yet. Somewhere, right now, someone is
> writing a paragraph that will become a window in a room that doesn't
> exist yet on a floor that hasn't been poured."*
>
> --- HIC Spatial Engineering Division, Internal Memo

---

**Document Coordinate:** `HIC-F00-A004-D4`
**Classification:** Infrastructure / Spatial Systems
**Maintained by:** HIC Structural Engineering Bureau
**Schema Version:** 4.0.0
**Next Review:** 2026-08-17
