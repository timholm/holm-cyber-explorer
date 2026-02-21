# STAGE 4: HOLM INTELLIGENCE COMPLEX -- SPATIAL DATA MODEL

## Agents 7-9: Blueprint File Format, Coordinate Systems, JSON/SVG Schema, Layer Definitions

The Holm Intelligence Complex (HIC) renders the sovereign documentation system as a navigable cyberpunk neon skyscraper. Every document domain occupies a physical floor. Every operational division inhabits a room. Every data flow traces a visible conduit through corridors and vertical shafts. This stage defines the spatial data model that makes that metaphor computationally real: the file formats that encode floors, the coordinate systems that position elements within them, and the layer architecture that controls what the user sees at every zoom level.

---

## Agent 7: Blueprint File Format

### 7.1 Overview

A **Blueprint** is the canonical machine-readable description of a single floor of the Holm Intelligence Complex. It encodes geometry (walls, doors, corridors), semantics (room names, types, security clearances), and connectivity (which rooms adjoin, which elevators and stairwells link floors vertically). Blueprints are consumed by the SVG renderer, the navigation engine, and the search indexer.

Every floor has exactly one Blueprint JSON file and one or more derived SVG files. The JSON is the source of truth; the SVG is a rendered artifact.

### 7.2 File Naming Convention

```
floor-{FLOOR_ID}.blueprint.json    # Source of truth
floor-{FLOOR_ID}.render.svg        # Derived visual render
floor-{FLOOR_ID}.thumb.svg         # Thumbnail for building overview
```

**FLOOR_ID** follows this scheme:

| Range | Format | Example | Meaning |
|-------|--------|---------|---------|
| Sub-basements | `B3`, `B2`, `B1` | `floor-B2.blueprint.json` | Below-ground infrastructure floors |
| Ground | `00` | `floor-00.blueprint.json` | Lobby / entry level |
| Standard floors | `01` through `99` | `floor-07.blueprint.json` | Domain and operational floors |
| Penthouse | `PH` | `floor-PH.blueprint.json` | Meta-governance level |
| Roof | `RF` | `floor-RF.blueprint.json` | Antenna array / external APIs |

Floor IDs are always zero-padded to two characters for standard floors. Sub-basement IDs use the `B` prefix with a single digit. This ensures correct lexicographic sorting: `B3 < B2 < B1 < 00 < 01 < ... < 99 < PH < RF`.

### 7.3 Metadata Header

Every Blueprint JSON file begins with a `meta` object:

```json
{
  "meta": {
    "schema_version": "1.0.0",
    "floor_id": "07",
    "floor_name": "Intelligence Operations",
    "floor_subtitle": "SIGINT / OSINT / Analysis",
    "building": "HIC-PRIMARY",
    "version": 14,
    "created": "2025-06-15T09:00:00Z",
    "last_modified": "2026-02-17T14:32:11Z",
    "modified_by": "agent:governance-core",
    "security_level": "RESTRICTED",
    "status": "ACTIVE",
    "checksum": "sha256:a4c9e8f..."
  }
}
```

**Field definitions:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `schema_version` | semver string | yes | Blueprint schema version. Current: `1.0.0` |
| `floor_id` | string | yes | Matches filename. Unique within building. |
| `floor_name` | string | yes | Human-readable floor name |
| `floor_subtitle` | string | no | Secondary descriptor shown in UI |
| `building` | string | yes | Building identifier (for multi-tower future) |
| `version` | integer | yes | Monotonically increasing edit counter |
| `created` | ISO 8601 | yes | First creation timestamp |
| `last_modified` | ISO 8601 | yes | Most recent edit timestamp |
| `modified_by` | string | yes | Agent or user who last edited |
| `security_level` | enum | yes | One of: `PUBLIC`, `INTERNAL`, `RESTRICTED`, `CLASSIFIED`, `SOVEREIGN` |
| `status` | enum | yes | One of: `DRAFT`, `ACTIVE`, `ARCHIVED`, `CONDEMNED` |
| `checksum` | string | yes | SHA-256 of the file contents excluding the checksum field itself |

### 7.4 Room Definition Schema

Each room is a JSON object within the `rooms` array:

```json
{
  "id": "room-07-03",
  "name": "OSINT Analysis Lab",
  "type": "OPERATIONS",
  "clearance": "RESTRICTED",
  "bounds": {
    "x": 320,
    "y": 160,
    "width": 256,
    "height": 192
  },
  "polygon": null,
  "color_primary": "#00e5ff",
  "color_accent": "#ff0055",
  "glow_intensity": 0.7,
  "connections": [
    { "target": "room-07-02", "type": "door", "position": "west", "locked": false },
    { "target": "room-07-04", "type": "door", "position": "east", "locked": true },
    { "target": "corridor-07-main", "type": "corridor", "position": "south", "locked": false }
  ],
  "vertical_connections": [
    { "target_floor": "06", "target_room": "room-06-03", "type": "elevator", "shaft_id": "shaft-A" }
  ],
  "contents": {
    "documents": ["intel-ops/osint-methodology.md", "intel-ops/source-evaluation.md"],
    "equipment": ["terminal-bank-alpha", "holotable-07"],
    "agents": ["agent:osint-crawler", "agent:source-validator"]
  },
  "status": {
    "power": "ONLINE",
    "alert_level": "NOMINAL",
    "occupancy": 3,
    "max_occupancy": 12,
    "temperature": "COOL"
  }
}
```

**Room types (enum):**

| Type | Description | Default Color |
|------|-------------|---------------|
| `OPERATIONS` | Active operational workspace | `#00e5ff` (cyan) |
| `COMMAND` | Command and control center | `#ff0055` (hot pink) |
| `ARCHIVE` | Document storage vault | `#8b5cf6` (violet) |
| `LAB` | Research and analysis | `#10b981` (emerald) |
| `CORRIDOR` | Passageway (not a true room) | `#1a1a2e` (dark blue) |
| `ELEVATOR_SHAFT` | Vertical transport | `#f59e0b` (amber) |
| `STAIRWELL` | Emergency vertical access | `#ef4444` (red) |
| `UTILITY` | Power, cooling, network closet | `#6b7280` (gray) |
| `LOBBY` | Reception and entry | `#e2e8f0` (silver) |
| `SECURE_VAULT` | Maximum security enclosure | `#dc2626` (crimson) |
| `SERVER_ROOM` | Compute infrastructure | `#06b6d4` (teal) |
| `CONFERENCE` | Meeting and briefing space | `#3b82f6` (blue) |

**Bounds vs. Polygon:** Most rooms use rectangular `bounds` (x, y, width, height). For irregular shapes, set `bounds` to `null` and provide `polygon` as an array of `[x, y]` coordinate pairs forming a closed shape:

```json
"polygon": [[320, 160], [576, 160], [576, 288], [448, 352], [320, 288]]
```

### 7.5 Wall, Door, and Corridor Data Structures

Walls are defined at the floor level, not per-room, to avoid duplication where two rooms share a wall:

```json
{
  "walls": [
    {
      "id": "wall-07-001",
      "type": "STRUCTURAL",
      "points": [[0, 0], [1024, 0]],
      "thickness": 4,
      "material": "reinforced",
      "color": "#334155"
    },
    {
      "id": "wall-07-002",
      "type": "PARTITION",
      "points": [[320, 0], [320, 384]],
      "thickness": 2,
      "material": "standard",
      "color": "#475569"
    }
  ],
  "doors": [
    {
      "id": "door-07-001",
      "wall_id": "wall-07-002",
      "position_along_wall": 0.35,
      "width": 32,
      "type": "SLIDING",
      "locked": false,
      "clearance_required": "INTERNAL",
      "connects": ["room-07-02", "room-07-03"],
      "color_frame": "#64748b",
      "color_active": "#22c55e"
    }
  ],
  "corridors": [
    {
      "id": "corridor-07-main",
      "path": [[0, 192], [1024, 192]],
      "width": 48,
      "type": "MAIN",
      "color": "#1e293b",
      "glow_strips": true,
      "glow_color": "#0ea5e9"
    }
  ]
}
```

**Wall types:** `STRUCTURAL` (load-bearing, cannot be removed), `PARTITION` (internal divider), `SECURITY` (reinforced, blast-resistant), `GLASS` (transparent, rendered with opacity), `ENERGY` (force-field style, rendered with animation).

**Door types:** `SLIDING`, `HINGED`, `BLAST_DOOR`, `AIRLOCK`, `VIRTUAL` (holographic barrier).

### 7.6 Complete Example: Floor 07 Blueprint

```json
{
  "meta": {
    "schema_version": "1.0.0",
    "floor_id": "07",
    "floor_name": "Intelligence Operations",
    "floor_subtitle": "SIGINT / OSINT / Analysis Nexus",
    "building": "HIC-PRIMARY",
    "version": 14,
    "created": "2025-06-15T09:00:00Z",
    "last_modified": "2026-02-17T14:32:11Z",
    "modified_by": "agent:governance-core",
    "security_level": "RESTRICTED",
    "status": "ACTIVE",
    "checksum": "sha256:a4c9e8f1b3d7..."
  },
  "dimensions": {
    "width": 1024,
    "height": 768,
    "grid_size": 16,
    "unit": "px"
  },
  "rooms": [
    {
      "id": "room-07-01",
      "name": "Floor 07 Lobby",
      "type": "LOBBY",
      "clearance": "PUBLIC",
      "bounds": { "x": 0, "y": 0, "width": 320, "height": 192 },
      "polygon": null,
      "color_primary": "#e2e8f0",
      "color_accent": "#0ea5e9",
      "glow_intensity": 0.3,
      "connections": [
        { "target": "corridor-07-main", "type": "corridor", "position": "south", "locked": false },
        { "target": "room-07-02", "type": "door", "position": "east", "locked": false }
      ],
      "vertical_connections": [
        { "target_floor": "00", "target_room": "room-00-lobby", "type": "elevator", "shaft_id": "shaft-A" }
      ],
      "contents": { "documents": [], "equipment": ["reception-terminal"], "agents": ["agent:greeter"] },
      "status": { "power": "ONLINE", "alert_level": "NOMINAL", "occupancy": 1, "max_occupancy": 8, "temperature": "AMBIENT" }
    },
    {
      "id": "room-07-02",
      "name": "Intel Briefing Room",
      "type": "CONFERENCE",
      "clearance": "INTERNAL",
      "bounds": { "x": 320, "y": 0, "width": 256, "height": 192 },
      "polygon": null,
      "color_primary": "#3b82f6",
      "color_accent": "#f59e0b",
      "glow_intensity": 0.5,
      "connections": [
        { "target": "room-07-01", "type": "door", "position": "west", "locked": false },
        { "target": "room-07-03", "type": "door", "position": "east", "locked": false },
        { "target": "corridor-07-main", "type": "corridor", "position": "south", "locked": false }
      ],
      "vertical_connections": [],
      "contents": { "documents": ["intel-ops/daily-brief.md"], "equipment": ["holoprojector-07A", "secure-display-wall"], "agents": [] },
      "status": { "power": "ONLINE", "alert_level": "NOMINAL", "occupancy": 0, "max_occupancy": 16, "temperature": "AMBIENT" }
    },
    {
      "id": "room-07-03",
      "name": "OSINT Analysis Lab",
      "type": "LAB",
      "clearance": "RESTRICTED",
      "bounds": { "x": 576, "y": 0, "width": 256, "height": 192 },
      "polygon": null,
      "color_primary": "#10b981",
      "color_accent": "#06b6d4",
      "glow_intensity": 0.7,
      "connections": [
        { "target": "room-07-02", "type": "door", "position": "west", "locked": false },
        { "target": "room-07-04", "type": "door", "position": "east", "locked": true },
        { "target": "corridor-07-main", "type": "corridor", "position": "south", "locked": false }
      ],
      "vertical_connections": [
        { "target_floor": "06", "target_room": "room-06-03", "type": "elevator", "shaft_id": "shaft-B" }
      ],
      "contents": { "documents": ["intel-ops/osint-methodology.md", "intel-ops/source-evaluation.md"], "equipment": ["terminal-bank-alpha", "holotable-07"], "agents": ["agent:osint-crawler", "agent:source-validator"] },
      "status": { "power": "ONLINE", "alert_level": "NOMINAL", "occupancy": 3, "max_occupancy": 12, "temperature": "COOL" }
    },
    {
      "id": "room-07-04",
      "name": "SIGINT Processing Center",
      "type": "OPERATIONS",
      "clearance": "CLASSIFIED",
      "bounds": { "x": 832, "y": 0, "width": 192, "height": 192 },
      "polygon": null,
      "color_primary": "#00e5ff",
      "color_accent": "#ff0055",
      "glow_intensity": 0.9,
      "connections": [
        { "target": "room-07-03", "type": "door", "position": "west", "locked": true },
        { "target": "corridor-07-main", "type": "corridor", "position": "south", "locked": true }
      ],
      "vertical_connections": [],
      "contents": { "documents": ["intel-ops/sigint-protocols.md"], "equipment": ["signal-array-console", "decryption-rig-07"], "agents": ["agent:sigint-processor"] },
      "status": { "power": "ONLINE", "alert_level": "ELEVATED", "occupancy": 2, "max_occupancy": 6, "temperature": "COLD" }
    },
    {
      "id": "room-07-05",
      "name": "Analyst Workstations",
      "type": "OPERATIONS",
      "clearance": "INTERNAL",
      "bounds": { "x": 0, "y": 240, "width": 384, "height": 256 },
      "polygon": null,
      "color_primary": "#0ea5e9",
      "color_accent": "#a855f7",
      "glow_intensity": 0.5,
      "connections": [
        { "target": "corridor-07-main", "type": "corridor", "position": "north", "locked": false },
        { "target": "room-07-06", "type": "door", "position": "east", "locked": false },
        { "target": "room-07-08", "type": "door", "position": "south", "locked": false }
      ],
      "vertical_connections": [],
      "contents": { "documents": ["intel-ops/analyst-handbook.md", "intel-ops/report-templates.md"], "equipment": ["workstation-cluster-A", "workstation-cluster-B"], "agents": ["agent:analyst-assist"] },
      "status": { "power": "ONLINE", "alert_level": "NOMINAL", "occupancy": 7, "max_occupancy": 20, "temperature": "AMBIENT" }
    },
    {
      "id": "room-07-06",
      "name": "Data Fusion Chamber",
      "type": "COMMAND",
      "clearance": "RESTRICTED",
      "bounds": { "x": 384, "y": 240, "width": 320, "height": 256 },
      "polygon": null,
      "color_primary": "#ff0055",
      "color_accent": "#00e5ff",
      "glow_intensity": 0.85,
      "connections": [
        { "target": "room-07-05", "type": "door", "position": "west", "locked": false },
        { "target": "corridor-07-main", "type": "corridor", "position": "north", "locked": false },
        { "target": "room-07-07", "type": "door", "position": "east", "locked": true }
      ],
      "vertical_connections": [
        { "target_floor": "08", "target_room": "room-08-06", "type": "elevator", "shaft_id": "shaft-C" },
        { "target_floor": "06", "target_room": "room-06-06", "type": "elevator", "shaft_id": "shaft-C" }
      ],
      "contents": { "documents": ["intel-ops/fusion-protocols.md", "intel-ops/correlation-engine.md", "intel-ops/threat-matrix.md"], "equipment": ["fusion-core-07", "threat-board-main", "3d-mapper"], "agents": ["agent:fusion-engine", "agent:correlator"] },
      "status": { "power": "ONLINE", "alert_level": "NOMINAL", "occupancy": 4, "max_occupancy": 10, "temperature": "COOL" }
    },
    {
      "id": "room-07-07",
      "name": "Secure Evidence Vault",
      "type": "SECURE_VAULT",
      "clearance": "SOVEREIGN",
      "bounds": { "x": 704, "y": 240, "width": 192, "height": 192 },
      "polygon": null,
      "color_primary": "#dc2626",
      "color_accent": "#fbbf24",
      "glow_intensity": 1.0,
      "connections": [
        { "target": "room-07-06", "type": "door", "position": "west", "locked": true }
      ],
      "vertical_connections": [],
      "contents": { "documents": ["intel-ops/classified-index.md"], "equipment": ["faraday-cage", "biometric-scanner", "dead-drop-terminal"], "agents": ["agent:vault-guardian"] },
      "status": { "power": "ISOLATED", "alert_level": "HEIGHTENED", "occupancy": 0, "max_occupancy": 3, "temperature": "COLD" }
    },
    {
      "id": "room-07-08",
      "name": "Server Closet 07-A",
      "type": "SERVER_ROOM",
      "clearance": "RESTRICTED",
      "bounds": { "x": 0, "y": 496, "width": 192, "height": 160 },
      "polygon": null,
      "color_primary": "#06b6d4",
      "color_accent": "#22c55e",
      "glow_intensity": 0.6,
      "connections": [
        { "target": "room-07-05", "type": "door", "position": "north", "locked": false }
      ],
      "vertical_connections": [
        { "target_floor": "B1", "target_room": "room-B1-core", "type": "conduit", "shaft_id": "riser-D" }
      ],
      "contents": { "documents": [], "equipment": ["rack-07A-01", "rack-07A-02", "ups-07A", "cooling-unit-07A"], "agents": ["agent:sysmon-07"] },
      "status": { "power": "ONLINE", "alert_level": "NOMINAL", "occupancy": 0, "max_occupancy": 2, "temperature": "COLD" }
    },
    {
      "id": "room-07-09",
      "name": "Elevator Shaft A",
      "type": "ELEVATOR_SHAFT",
      "clearance": "INTERNAL",
      "bounds": { "x": 896, "y": 240, "width": 128, "height": 128 },
      "polygon": null,
      "color_primary": "#f59e0b",
      "color_accent": "#f59e0b",
      "glow_intensity": 0.4,
      "connections": [],
      "vertical_connections": [
        { "target_floor": "ALL", "target_room": null, "type": "elevator", "shaft_id": "shaft-A" }
      ],
      "contents": { "documents": [], "equipment": ["elevator-car-A"], "agents": [] },
      "status": { "power": "ONLINE", "alert_level": "NOMINAL", "occupancy": 0, "max_occupancy": 8, "temperature": "AMBIENT" }
    }
  ],
  "walls": [
    { "id": "wall-07-N", "type": "STRUCTURAL", "points": [[0, 0], [1024, 0]], "thickness": 4, "material": "reinforced", "color": "#334155" },
    { "id": "wall-07-S", "type": "STRUCTURAL", "points": [[0, 768], [1024, 768]], "thickness": 4, "material": "reinforced", "color": "#334155" },
    { "id": "wall-07-W", "type": "STRUCTURAL", "points": [[0, 0], [0, 768]], "thickness": 4, "material": "reinforced", "color": "#334155" },
    { "id": "wall-07-E", "type": "STRUCTURAL", "points": [[1024, 0], [1024, 768]], "thickness": 4, "material": "reinforced", "color": "#334155" },
    { "id": "wall-07-001", "type": "PARTITION", "points": [[320, 0], [320, 192]], "thickness": 2, "material": "standard", "color": "#475569" },
    { "id": "wall-07-002", "type": "PARTITION", "points": [[576, 0], [576, 192]], "thickness": 2, "material": "standard", "color": "#475569" },
    { "id": "wall-07-003", "type": "SECURITY", "points": [[832, 0], [832, 192]], "thickness": 3, "material": "reinforced", "color": "#991b1b" },
    { "id": "wall-07-004", "type": "PARTITION", "points": [[0, 192], [1024, 192]], "thickness": 2, "material": "standard", "color": "#475569" },
    { "id": "wall-07-005", "type": "PARTITION", "points": [[384, 240], [384, 496]], "thickness": 2, "material": "standard", "color": "#475569" },
    { "id": "wall-07-006", "type": "SECURITY", "points": [[704, 240], [704, 432]], "thickness": 3, "material": "reinforced", "color": "#991b1b" },
    { "id": "wall-07-007", "type": "PARTITION", "points": [[0, 496], [192, 496]], "thickness": 2, "material": "standard", "color": "#475569" }
  ],
  "doors": [
    { "id": "door-07-001", "wall_id": "wall-07-001", "position_along_wall": 0.5, "width": 32, "type": "SLIDING", "locked": false, "clearance_required": "PUBLIC", "connects": ["room-07-01", "room-07-02"], "color_frame": "#64748b", "color_active": "#22c55e" },
    { "id": "door-07-002", "wall_id": "wall-07-002", "position_along_wall": 0.5, "width": 32, "type": "SLIDING", "locked": false, "clearance_required": "INTERNAL", "connects": ["room-07-02", "room-07-03"], "color_frame": "#64748b", "color_active": "#22c55e" },
    { "id": "door-07-003", "wall_id": "wall-07-003", "position_along_wall": 0.5, "width": 32, "type": "BLAST_DOOR", "locked": true, "clearance_required": "CLASSIFIED", "connects": ["room-07-03", "room-07-04"], "color_frame": "#991b1b", "color_active": "#ef4444" },
    { "id": "door-07-004", "wall_id": "wall-07-005", "position_along_wall": 0.4, "width": 32, "type": "SLIDING", "locked": false, "clearance_required": "INTERNAL", "connects": ["room-07-05", "room-07-06"], "color_frame": "#64748b", "color_active": "#22c55e" },
    { "id": "door-07-005", "wall_id": "wall-07-006", "position_along_wall": 0.4, "width": 32, "type": "BLAST_DOOR", "locked": true, "clearance_required": "SOVEREIGN", "connects": ["room-07-06", "room-07-07"], "color_frame": "#991b1b", "color_active": "#dc2626" },
    { "id": "door-07-006", "wall_id": "wall-07-007", "position_along_wall": 0.5, "width": 32, "type": "SLIDING", "locked": false, "clearance_required": "RESTRICTED", "connects": ["room-07-05", "room-07-08"], "color_frame": "#64748b", "color_active": "#22c55e" }
  ],
  "corridors": [
    { "id": "corridor-07-main", "path": [[0, 192], [1024, 192]], "width": 48, "type": "MAIN", "color": "#1e293b", "glow_strips": true, "glow_color": "#0ea5e9" }
  ],
  "data_conduits": [
    { "id": "conduit-07-001", "path": [[96, 560], [96, 768]], "type": "FIBER", "bandwidth": "10Gbps", "color": "#22c55e", "animated": true },
    { "id": "conduit-07-002", "path": [[544, 368], [544, 192], [544, 0]], "type": "FIBER", "bandwidth": "40Gbps", "color": "#06b6d4", "animated": true }
  ]
}
```

### 7.7 SVG Template Specification

The rendered SVG for each floor follows a strict template. The SVG uses a viewBox of `0 0 1024 768` and organizes content into named `<g>` groups corresponding to the layer system defined in Agent 9.

```xml
<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink"
     viewBox="0 0 1024 768"
     width="1024" height="768"
     data-floor-id="07"
     data-schema-version="1.0.0"
     data-security-level="RESTRICTED">

  <defs>
    <!-- Neon glow filter -->
    <filter id="neon-glow" x="-20%" y="-20%" width="140%" height="140%">
      <feGaussianBlur in="SourceGraphic" stdDeviation="3" result="blur"/>
      <feColorMatrix in="blur" type="matrix"
        values="0 0 0 0 0  0 0 0 0 0.9  0 0 0 0 1  0 0 0 1 0" result="glow"/>
      <feMerge>
        <feMergeNode in="glow"/>
        <feMergeNode in="SourceGraphic"/>
      </feMerge>
    </filter>

    <!-- Animated pulse for active elements -->
    <filter id="pulse-glow">
      <feGaussianBlur stdDeviation="4">
        <animate attributeName="stdDeviation" values="2;6;2" dur="2s" repeatCount="indefinite"/>
      </feGaussianBlur>
    </filter>

    <!-- Scanline overlay pattern -->
    <pattern id="scanlines" patternUnits="userSpaceOnUse" width="4" height="4">
      <line x1="0" y1="0" x2="4" y2="0" stroke="rgba(0,0,0,0.15)" stroke-width="1"/>
    </pattern>

    <!-- Grid pattern -->
    <pattern id="grid-16" patternUnits="userSpaceOnUse" width="16" height="16">
      <rect width="16" height="16" fill="none"/>
      <path d="M 16 0 L 0 0 0 16" fill="none" stroke="rgba(100,200,255,0.06)" stroke-width="0.5"/>
    </pattern>
    <pattern id="grid-64" patternUnits="userSpaceOnUse" width="64" height="64">
      <rect width="64" height="64" fill="none"/>
      <path d="M 64 0 L 0 0 0 64" fill="none" stroke="rgba(100,200,255,0.12)" stroke-width="0.5"/>
    </pattern>

    <!-- Gradient for corridor glow strips -->
    <linearGradient id="glow-strip-horiz" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" stop-color="#0ea5e9" stop-opacity="0"/>
      <stop offset="50%" stop-color="#0ea5e9" stop-opacity="0.8"/>
      <stop offset="100%" stop-color="#0ea5e9" stop-opacity="0"/>
    </linearGradient>
  </defs>

  <!-- LAYER 0: Grid Background -->
  <g id="layer-0-grid" class="layer" data-layer="grid">
    <rect width="1024" height="768" fill="#0a0a1a"/>
    <rect width="1024" height="768" fill="url(#grid-16)"/>
    <rect width="1024" height="768" fill="url(#grid-64)"/>
  </g>

  <!-- LAYER 1: Structural Walls -->
  <g id="layer-1-walls" class="layer" data-layer="walls">
    <!-- walls rendered here from blueprint data -->
  </g>

  <!-- LAYER 2: Room Fills -->
  <g id="layer-2-rooms" class="layer" data-layer="rooms">
    <!-- room rectangles/polygons rendered here -->
  </g>

  <!-- Layers 3-9 follow the same structure (see Agent 9) -->

  <!-- Scanline overlay (always on top) -->
  <rect width="1024" height="768" fill="url(#scanlines)" pointer-events="none" opacity="0.5"/>
</svg>
```

### 7.8 Validation Rules

1. **Unique IDs:** Every `id` field within a floor must be unique. Room IDs must be globally unique across all floors (enforced by the `room-{floor}-{seq}` convention).
2. **Bounds within floor:** All room bounds and polygon coordinates must fall within `[0, 0]` to `[dimensions.width, dimensions.height]`.
3. **Grid alignment:** All `x`, `y`, `width`, and `height` values in bounds must be multiples of `grid_size` (default 16).
4. **No overlapping rooms:** Room bounds must not overlap. Validation computes axis-aligned bounding box intersection for all room pairs.
5. **Connection symmetry:** If room A lists room B as a connection, room B must list room A with the same door/corridor ID.
6. **Door placement:** Every door must reference a valid `wall_id`, and its `connects` array must contain exactly two valid room IDs.
7. **Clearance hierarchy:** A door's `clearance_required` must be greater than or equal to the lower of the two connected rooms' clearances. Hierarchy: `PUBLIC < INTERNAL < RESTRICTED < CLASSIFIED < SOVEREIGN`.
8. **Vertical connection validity:** Every `vertical_connections` entry must reference a valid `target_floor` that has a corresponding blueprint file.
9. **Schema version compatibility:** The renderer must reject blueprints with a `schema_version` major version it does not support.
10. **Checksum integrity:** On load, the checksum must be verified by hashing all file contents except the checksum field.

---

## Agent 8: Coordinate Systems

### 8.1 Building Coordinate System (BCS)

The Building Coordinate System establishes the global spatial frame for the entire Holm Intelligence Complex. It uses a left-handed coordinate system consistent with SVG conventions (Y increases downward).

**Origin:** The top-left corner of the ground floor (floor `00`) exterior wall, at the northwest corner of the building footprint.

**Axes:**

- **X-axis:** Positive rightward (east). Range: `0` to `BUILDING_WIDTH` (default 1024 units).
- **Y-axis:** Positive downward (south). Range: `0` to `BUILDING_DEPTH` (default 768 units per floor).
- **Z-axis:** Positive upward. Each floor occupies a Z-band. Floor `00` occupies Z `0` to `FLOOR_HEIGHT`. Floor `01` occupies Z `FLOOR_HEIGHT` to `2 * FLOOR_HEIGHT`. Sub-basements extend into negative Z.

**Units:** All coordinates are in **logical pixels (px)**. One logical pixel equals one SVG user unit. At the default viewport, 1px maps to approximately 0.3 meters of in-world scale, making the building footprint roughly 307m x 230m -- appropriate for a megastructure.

```
BUILDING COORDINATE SYSTEM (BCS)
=================================

  Origin (0,0,0)
    +---------------------------------------------> X (East)
    |                                                1024 px
    |   Floor 00 (Ground)           Z = 0
    |   +--------------------------------------+
    |   |                                      |
    |   |         Ground Floor Plan            |
    |   |         (1024 x 768 px)              |
    |   |                                      |
    |   +--------------------------------------+
    |
    |   Floor 01                    Z = 1 * FLOOR_HEIGHT
    |   +--------------------------------------+
    |   |                                      |
    |   |                                      |
    |   +--------------------------------------+
    |
    v Y (South)
    768 px

  Z-axis (vertical):
    RF  ---- Z = N * FLOOR_HEIGHT
    PH  ---- Z = (N-1) * FLOOR_HEIGHT
    ...
    07  ---- Z = 7 * FLOOR_HEIGHT
    ...
    01  ---- Z = 1 * FLOOR_HEIGHT
    00  ---- Z = 0                    <-- ORIGIN
    B1  ---- Z = -1 * FLOOR_HEIGHT
    B2  ---- Z = -2 * FLOOR_HEIGHT
    B3  ---- Z = -3 * FLOOR_HEIGHT

  FLOOR_HEIGHT = 128 px (logical)
```

### 8.2 Floor-Local Coordinate System (FCS)

Each floor has its own 2D coordinate system. For all standard floors, the FCS is identical to the BCS X/Y plane -- the origin is the top-left corner of the floor's bounding rectangle. The Z component is implicit (determined by the floor ID).

```
FLOOR-LOCAL COORDINATE SYSTEM (FCS)
=====================================

  (0, 0) -----------------------> X (1024)
    |
    |    +--------+  +--------+
    |    | Room 1 |  | Room 2 |
    |    |        |  |        |
    |    +--------+  +--------+
    |
    |    ========================  <- Corridor
    |
    |    +-------------+  +----+
    |    |   Room 5     |  | R6 |
    |    |              |  |    |
    |    +-------------+  +----+
    |
    v Y (768)
```

**Transform from FCS to BCS:**

```
BCS.x = FCS.x
BCS.y = FCS.y
BCS.z = floor_z_offset(floor_id)
```

Where `floor_z_offset` maps floor IDs to Z values:

```javascript
function floor_z_offset(floor_id) {
  if (floor_id.startsWith('B')) {
    return -parseInt(floor_id.substring(1)) * FLOOR_HEIGHT;
  }
  if (floor_id === 'PH') return PENTHOUSE_Z;
  if (floor_id === 'RF') return ROOF_Z;
  return parseInt(floor_id) * FLOOR_HEIGHT;
}
```

### 8.3 Room-Local Coordinate System (RCS)

Within each room, a local coordinate system allows room contents (furniture, equipment, labels) to be positioned relative to the room's own top-left corner rather than the floor origin. This simplifies room templates and makes rooms relocatable.

**Origin:** Top-left corner of the room's bounding box.

**Axes:** Same orientation as FCS (X right, Y down).

**Range:** `(0, 0)` to `(room.bounds.width, room.bounds.height)`.

**Transform from RCS to FCS:**

```
FCS.x = RCS.x + room.bounds.x
FCS.y = RCS.y + room.bounds.y
```

**Transform from RCS to BCS:**

```
BCS.x = RCS.x + room.bounds.x
BCS.y = RCS.y + room.bounds.y
BCS.z = floor_z_offset(floor_id)
```

For polygon-bounded rooms (no rectangular `bounds`), the RCS origin is the minimum X and minimum Y of all polygon vertices:

```javascript
const rcs_origin_x = Math.min(...polygon.map(p => p[0]));
const rcs_origin_y = Math.min(...polygon.map(p => p[1]));
```

### 8.4 Coordinate Transform Summary

```
                +-----+
                | BCS |   Building Coordinate System (global 3D)
                +--+--+
                   |
         floor_z_offset(floor_id) subtracted
                   |
                +--v--+
                | FCS |   Floor-Local Coordinate System (2D per floor)
                +--+--+
                   |
         room.bounds.x/y subtracted
                   |
                +--v--+
                | RCS |   Room-Local Coordinate System (2D per room)
                +-----+

  BCS(x, y, z) = (RCS.x + room.x, RCS.y + room.y, floor_z_offset)
  FCS(x, y)    = (RCS.x + room.x, RCS.y + room.y)
  RCS(x, y)    = (FCS.x - room.x, FCS.y - room.y)
```

### 8.5 Grid Snap Specification

All positioning in the HIC uses a **16px base grid**. This ensures visual alignment, simplifies coordinate math, and produces clean SVG output.

**Grid rules:**

| Property | Grid Multiple | Notes |
|----------|---------------|-------|
| Room origin (x, y) | 16px | Must snap to grid |
| Room dimensions (w, h) | 16px | Must be multiples of 16 |
| Wall endpoints | 16px | Must snap to grid |
| Door positions | 8px | Half-grid allowed for door placement along a wall |
| Furniture placement | 8px | Half-grid for finer positioning |
| Label positions | 1px | Labels float freely for readability |
| Corridor paths | 16px | Waypoints snap to grid |
| Data conduit paths | 8px | Half-grid for routing flexibility |

**Snap function:**

```javascript
function snapToGrid(value, gridSize = 16) {
  return Math.round(value / gridSize) * gridSize;
}

function snapToHalfGrid(value, gridSize = 16) {
  const half = gridSize / 2;
  return Math.round(value / half) * half;
}
```

**Major grid vs. minor grid:** The 16px grid is the minor grid. A 64px major grid (4x minor) provides visual hierarchy in the background pattern and is used for coarse alignment guides.

### 8.6 Viewport Coordinate Mapping

The **Viewport Coordinate System (VCS)** maps the logical floor coordinates to actual screen pixels. It is defined by three parameters: **pan offset**, **zoom scale**, and **screen dimensions**.

```
VCS.x = (FCS.x - pan.x) * zoom + screen.width / 2
VCS.y = (FCS.y - pan.y) * zoom + screen.height / 2
```

Inverse (screen to floor):

```
FCS.x = (VCS.x - screen.width / 2) / zoom + pan.x
FCS.y = (VCS.y - screen.height / 2) / zoom + pan.y
```

This uses a center-anchored pan model: `pan.x` and `pan.y` represent the FCS coordinate that appears at the center of the screen.

### 8.7 Zoom Level Definitions

The HIC supports four named zoom levels with continuous zoom between them. Named levels serve as snap targets and determine layer visibility (see Agent 9).

| Level | Name | Scale Factor | Visible Area | Primary Use |
|-------|------|-------------|--------------|-------------|
| Z0 | Building View | 0.08 - 0.15 | Entire building (all floors stacked) | Navigation, floor selection |
| Z1 | Floor View | 0.25 - 0.60 | Full floor plan | Room overview, corridor navigation |
| Z2 | Room View | 0.80 - 1.50 | Single room + neighbors | Room inspection, document browsing |
| Z3 | Detail View | 2.00 - 4.00 | Equipment/terminal detail | Document reading, terminal interaction |

**Zoom behavior:**

- Scroll wheel increments zoom by `1.15x` per tick (logarithmic scaling).
- Double-click zooms to the next named level, centered on the click point.
- Pinch-to-zoom is supported on touch devices with the same logarithmic curve.
- Minimum zoom: `0.05` (whole building fits in a small widget).
- Maximum zoom: `5.00` (terminal text becomes legible).

**Zoom snap thresholds:** When the user stops zooming within 15% of a named level boundary, the viewport smoothly animates to the nearest named level over 200ms (ease-out curve). This can be disabled in user preferences.

### 8.8 Pan and Scroll Boundaries

Panning is constrained so the user cannot scroll entirely off the floor plan. The constraint is defined as a padded bounding box:

```javascript
const PAN_PADDING = 128; // px of "overshoot" allowed

const panBounds = {
  minX: -PAN_PADDING,
  minY: -PAN_PADDING,
  maxX: floor.dimensions.width + PAN_PADDING,
  maxY: floor.dimensions.height + PAN_PADDING
};

function clampPan(pan, bounds) {
  return {
    x: Math.max(bounds.minX, Math.min(bounds.maxX, pan.x)),
    y: Math.max(bounds.minY, Math.min(bounds.maxY, pan.y))
  };
}
```

At **Building View (Z0)**, pan boundaries expand to encompass the full building envelope including all floors rendered in an isometric or stacked layout. The vertical scroll range becomes `total_floors * FLOOR_HEIGHT + vertical_padding`.

**Elastic overscroll:** When the user drags beyond the pan boundary, the viewport allows up to 32px of elastic overscroll with a rubber-band deceleration curve. On release, it snaps back to the boundary with a 150ms spring animation.

### 8.9 Coordinate Precision Rules

All coordinate values in Blueprint JSON files are stored as **integers** (whole pixel values). No floating-point coordinates are permitted in persistent storage. This eliminates rounding drift and ensures deterministic rendering.

| Context | Precision | Type | Rationale |
|---------|-----------|------|-----------|
| Blueprint JSON coordinates | Integer | `int` | Deterministic, grid-aligned |
| Runtime pan/zoom state | Float (64-bit) | `double` | Smooth animation |
| SVG attribute output | 2 decimal places | `string` | Balance precision and file size |
| Hit-test calculations | Float (64-bit) | `double` | Accurate cursor detection |
| Serialized viewport state | 4 decimal places | `string` | Restore exact view position |

**Anti-aliasing note:** The renderer should use `shape-rendering="crispEdges"` for walls and room outlines (which are grid-aligned) and `shape-rendering="auto"` for diagonal data conduits and decorative elements.

### 8.10 ASCII Diagram of Coordinate Spaces

```
HOLM INTELLIGENCE COMPLEX -- COORDINATE SPACE OVERVIEW
========================================================

BUILDING VIEW (Z0: scale 0.08-0.15)
+---------------------------------------------------+
|                                                     |
|   [ RF ] Roof / Antenna Array                       |
|   [ PH ] Penthouse / Meta-Governance                |
|   [ 20 ] ...                                        |
|   [ .. ] ...                                        |
|   [ 07 ] <<<  Intelligence Operations  >>>          |
|   [ .. ] ...                                        |
|   [ 01 ] ...                                        |
|   [ 00 ] Ground / Lobby                             |
|   [ B1 ] Sub-basement 1 / Infrastructure            |
|   [ B2 ] Sub-basement 2 / Deep Archives             |
|   [ B3 ] Sub-basement 3 / Core Systems              |
|                                                     |
+---------------------------------------------------+
  ^--- Each floor is a clickable strip in this view

FLOOR VIEW (Z1: scale 0.25-0.60) -- Floor 07
+---------------------------------------------------+
| (0,0)                                    (1024,0)  |
|  +--------+--------+--------+------+               |
|  | Lobby  | Brief  | OSINT  | SGNIT|               |
|  | 07-01  | 07-02  | 07-03  | 07-04|               |
|  +--------+--------+--------+------+               |
|  =================CORRIDOR===================      |
|  +--------------+----------+------+ +------+       |
|  | Analysts     | Fusion   | Vault| | Elev |       |
|  | 07-05        | 07-06    | 07-07| | Shaft|       |
|  +--------------+----------+------+ +------+       |
|  +------+                                           |
|  | Svr  |                                           |
|  | 07-08|                                           |
|  +------+                                           |
| (0,768)                                (1024,768)  |
+---------------------------------------------------+

ROOM VIEW (Z2: scale 0.80-1.50) -- Room 07-06
+---------------------------------------------------+
|  DATA FUSION CHAMBER                                |
|  +---------+---------+---------+                    |
|  | Terminal | Fusion  | Threat  |                   |
|  | Bank     | Core    | Board   |                   |
|  |  [T]     | [=FC=]  |  [TB]   |                   |
|  +---------+---------+---------+                    |
|  |         3D Mapper            |                   |
|  |           [3D]               |                   |
|  +------------------------------+                   |
|  << Door to 07-05    Door to 07-07 (LOCKED) >>     |
+---------------------------------------------------+

DETAIL VIEW (Z3: scale 2.00-4.00) -- Fusion Core
+---------------------------------------------------+
|  +-------------------------------------------+     |
|  |  FUSION CORE 07                            |     |
|  |  Status: ONLINE    Threads: 847/1024       |     |
|  |  +---------+  +---------+  +---------+    |     |
|  |  | OSINT   |  | SIGINT  |  | HUMINT  |    |     |
|  |  | Feed    |  | Feed    |  | Feed    |    |     |
|  |  | [====]  |  | [==  ]  |  | [=    ] |    |     |
|  |  +---------+  +---------+  +---------+    |     |
|  |  Correlation confidence: 78.3%             |     |
|  |  Last cycle: 2.3s ago                      |     |
|  +-------------------------------------------+     |
+---------------------------------------------------+
```

---

## Agent 9: Layer Definitions

### 9.1 Layer Stack Overview

The HIC renderer uses a fixed 10-layer stack (Layers 0 through 9). Layers are rendered in ascending order (Layer 0 is the bottommost, Layer 9 is the topmost). Each layer corresponds to an SVG `<g>` group with a unique `id` and `data-layer` attribute.

Layers are conceptually divided into three tiers:

| Tier | Layers | Purpose |
|------|--------|---------|
| **Foundation** | 0-2 | Static structural elements (grid, walls, rooms) |
| **Content** | 3-6 | Dynamic content elements (furniture, labels, security, data) |
| **Interactive** | 7-9 | User interaction and status overlays |

### 9.2 Layer 0: Grid Background

**Purpose:** Provides the spatial reference grid and the dark ambient background that establishes the cyberpunk aesthetic.

**SVG structure:**

```xml
<g id="layer-0-grid" class="layer" data-layer="grid" data-layer-index="0">
  <!-- Deep space background -->
  <rect width="1024" height="768" fill="#0a0a1a"/>

  <!-- Minor grid (16px) -->
  <rect width="1024" height="768" fill="url(#grid-16)" opacity="0.4"/>

  <!-- Major grid (64px) -->
  <rect width="1024" height="768" fill="url(#grid-64)" opacity="0.6"/>

  <!-- Origin crosshair (debug, hidden by default) -->
  <line x1="0" y1="0" x2="32" y2="0" stroke="#ff0055" stroke-width="1" class="debug-only"/>
  <line x1="0" y1="0" x2="0" y2="32" stroke="#00e5ff" stroke-width="1" class="debug-only"/>
</g>
```

**Rendering rules:**
- Always visible at all zoom levels.
- Grid line opacity scales with zoom: at Z0, only the major grid is visible (opacity 0.3). At Z3, both grids render at full opacity.
- Background color is always `#0a0a1a` (near-black with a blue-midnight undertone).

### 9.3 Layer 1: Structural Walls

**Purpose:** Renders all wall segments, including structural, partition, security, glass, and energy walls.

**SVG structure:**

```xml
<g id="layer-1-walls" class="layer" data-layer="walls" data-layer-index="1">
  <!-- Structural walls (thick, high contrast) -->
  <line x1="0" y1="0" x2="1024" y2="0"
        stroke="#334155" stroke-width="4" stroke-linecap="square"
        data-wall-id="wall-07-N" data-wall-type="STRUCTURAL"/>

  <!-- Partition walls (thinner, lower contrast) -->
  <line x1="320" y1="0" x2="320" y2="192"
        stroke="#475569" stroke-width="2" stroke-linecap="square"
        data-wall-id="wall-07-001" data-wall-type="PARTITION"/>

  <!-- Security walls (with glow effect) -->
  <line x1="832" y1="0" x2="832" y2="192"
        stroke="#991b1b" stroke-width="3" stroke-linecap="square"
        filter="url(#neon-glow)" class="security-wall"
        data-wall-id="wall-07-003" data-wall-type="SECURITY"/>

  <!-- Glass walls (semi-transparent) -->
  <!-- <line ... stroke="rgba(148,163,184,0.4)" stroke-width="2" stroke-dasharray="8 4"/> -->

  <!-- Energy walls (animated) -->
  <!-- <line ... stroke="#8b5cf6" stroke-width="2" filter="url(#pulse-glow)"/> -->

  <!-- Door gaps are rendered as breaks in the wall line -->
  <!-- Doors themselves are part of this layer -->
  <g class="door" data-door-id="door-07-001" data-locked="false">
    <rect x="304" y="94" width="32" height="4" fill="#64748b" rx="1"/>
    <rect x="304" y="94" width="32" height="4" fill="#22c55e" opacity="0.6" rx="1">
      <animate attributeName="opacity" values="0.4;0.8;0.4" dur="3s" repeatCount="indefinite"/>
    </rect>
  </g>

  <!-- Locked door (red indicator) -->
  <g class="door locked" data-door-id="door-07-003" data-locked="true">
    <rect x="816" y="94" width="32" height="4" fill="#991b1b" rx="1"/>
    <rect x="816" y="94" width="32" height="4" fill="#ef4444" opacity="0.8" rx="1">
      <animate attributeName="opacity" values="0.6;1.0;0.6" dur="1.5s" repeatCount="indefinite"/>
    </rect>
    <!-- Lock icon -->
    <text x="832" y="88" text-anchor="middle" fill="#ef4444" font-size="8" font-family="monospace">LOCKED</text>
  </g>
</g>
```

**Rendering rules:**
- Visible at Z1, Z2, and Z3. At Z0 (building view), walls are simplified to floor outlines only.
- Security walls always render with the `neon-glow` filter.
- Energy walls always animate. Animation pauses when the layer is not visible.
- Door color indicates state: green = unlocked, red = locked, amber = restricted access.

### 9.4 Layer 2: Room Fills

**Purpose:** Renders the filled interiors of rooms, providing color-coding by room type and visual affordance for clickable/navigable spaces.

**SVG structure:**

```xml
<g id="layer-2-rooms" class="layer" data-layer="rooms" data-layer-index="2">
  <!-- Standard rectangular room -->
  <rect x="0" y="0" width="320" height="192"
        fill="#e2e8f0" fill-opacity="0.08"
        stroke="#e2e8f0" stroke-width="1" stroke-opacity="0.3"
        rx="2" ry="2"
        class="room" data-room-id="room-07-01" data-room-type="LOBBY"
        data-clearance="PUBLIC"/>

  <!-- Room with higher glow intensity -->
  <rect x="576" y="0" width="256" height="192"
        fill="#10b981" fill-opacity="0.12"
        stroke="#10b981" stroke-width="1" stroke-opacity="0.5"
        filter="url(#neon-glow)"
        rx="2" ry="2"
        class="room active" data-room-id="room-07-03" data-room-type="LAB"
        data-clearance="RESTRICTED"/>

  <!-- Polygon room example -->
  <!-- <polygon points="320,160 576,160 576,288 448,352 320,288"
              fill="#10b981" fill-opacity="0.10"
              stroke="#10b981" stroke-width="1"/> -->

  <!-- Corridor fill -->
  <rect x="0" y="192" width="1024" height="48"
        fill="#1e293b" fill-opacity="0.15"
        class="corridor" data-corridor-id="corridor-07-main"/>

  <!-- Corridor glow strips -->
  <line x1="0" y1="193" x2="1024" y2="193"
        stroke="url(#glow-strip-horiz)" stroke-width="1" opacity="0.6"/>
  <line x1="0" y1="239" x2="1024" y2="239"
        stroke="url(#glow-strip-horiz)" stroke-width="1" opacity="0.6"/>
</g>
```

**Rendering rules:**
- Visible at all zoom levels. At Z0, rooms render as solid blocks with higher fill opacity (0.3) for legibility in the miniature view.
- Fill opacity is computed as `base_opacity + (glow_intensity * 0.15)`. Base opacity is `0.06`.
- On hover, room fill opacity increases by `0.08` and the stroke brightens to full opacity with a 150ms CSS transition.
- The currently selected/focused room has a pulsing border animation.

### 9.5 Layer 3: Furniture and Equipment

**Purpose:** Renders in-room objects such as terminal banks, server racks, holographic tables, and other equipment defined in the room's `contents.equipment` array.

**SVG structure:**

```xml
<g id="layer-3-furniture" class="layer" data-layer="furniture" data-layer-index="3">
  <!-- Equipment items are positioned in Room-Local Coordinates, -->
  <!-- then translated by the room's bounds offset.              -->

  <g transform="translate(384, 240)" data-room-id="room-07-06">
    <!-- Fusion Core (centered in room) -->
    <g class="equipment" data-equipment-id="fusion-core-07" transform="translate(112, 48)">
      <rect x="0" y="0" width="96" height="64" rx="4"
            fill="#1a1a2e" stroke="#00e5ff" stroke-width="1.5"/>
      <text x="48" y="36" text-anchor="middle" fill="#00e5ff"
            font-family="monospace" font-size="9">FUSION CORE</text>
      <!-- Activity indicator -->
      <circle cx="48" cy="52" r="4" fill="#22c55e">
        <animate attributeName="r" values="3;5;3" dur="1.5s" repeatCount="indefinite"/>
      </circle>
    </g>

    <!-- Threat Board -->
    <g class="equipment" data-equipment-id="threat-board-main" transform="translate(224, 16)">
      <rect x="0" y="0" width="80" height="48" rx="2"
            fill="#0f172a" stroke="#ff0055" stroke-width="1"/>
      <text x="40" y="28" text-anchor="middle" fill="#ff0055"
            font-family="monospace" font-size="8">THREAT BOARD</text>
    </g>
  </g>
</g>
```

**Rendering rules:**
- Visible at Z2 and Z3 only. Equipment is too small to be meaningful at Z0 and Z1.
- Equipment items are clickable at Z2+ and trigger detail popups or zoom-to-detail actions.
- Active equipment shows animated indicators (pulsing circles, scrolling text).
- Inactive/offline equipment renders with reduced opacity (0.3) and a desaturated color.

### 9.6 Layer 4: Labels and Text

**Purpose:** Renders room names, floor identifiers, coordinate labels, and any other textual annotations.

**SVG structure:**

```xml
<g id="layer-4-labels" class="layer" data-layer="labels" data-layer-index="4">
  <!-- Room name labels -->
  <text x="160" y="96" text-anchor="middle" dominant-baseline="middle"
        fill="#e2e8f0" font-family="'JetBrains Mono', 'Fira Code', monospace"
        font-size="11" font-weight="600" letter-spacing="0.05em"
        class="room-label" data-room-id="room-07-01">
    LOBBY
  </text>

  <!-- Room ID sublabel -->
  <text x="160" y="112" text-anchor="middle" dominant-baseline="middle"
        fill="#94a3b8" font-family="monospace" font-size="7"
        class="room-sublabel" data-room-id="room-07-01">
    07-01 / PUBLIC
  </text>

  <!-- Floor header (shown at Z1) -->
  <text x="512" y="-20" text-anchor="middle"
        fill="#00e5ff" font-family="monospace" font-size="14"
        font-weight="700" letter-spacing="0.15em" filter="url(#neon-glow)"
        class="floor-header">
    FLOOR 07 -- INTELLIGENCE OPERATIONS
  </text>

  <!-- Coordinate markers (debug, togglable) -->
  <g class="coord-labels debug-only" opacity="0.4">
    <text x="0" y="-4" fill="#475569" font-size="6" font-family="monospace">0,0</text>
    <text x="1024" y="-4" text-anchor="end" fill="#475569" font-size="6" font-family="monospace">1024,0</text>
    <text x="0" y="780" fill="#475569" font-size="6" font-family="monospace">0,768</text>
  </g>
</g>
```

**Rendering rules:**
- Room name labels are visible at Z1, Z2, and Z3. At Z0, only floor identifiers are shown.
- Font sizes scale inversely with zoom at Z0 and Z1 to maintain readability (labels rendered at a fixed screen-space size via `vector-effect` or manual scaling).
- At Z2+, sublabels (room IDs, clearance levels) become visible.
- At Z3, coordinate debug labels can be toggled on via a developer setting.
- All text uses monospace fonts exclusively to maintain the cyberpunk terminal aesthetic.

### 9.7 Layer 5: Security Overlays

**Purpose:** Renders security zones, clearance boundaries, surveillance coverage areas, and access restriction indicators.

**SVG structure:**

```xml
<g id="layer-5-security" class="layer" data-layer="security" data-layer-index="5">
  <!-- Security zone highlight for CLASSIFIED rooms -->
  <rect x="832" y="0" width="192" height="192"
        fill="#dc2626" fill-opacity="0.04"
        stroke="#dc2626" stroke-width="1" stroke-dasharray="8 4"
        class="security-zone classified" data-clearance="CLASSIFIED"/>

  <!-- Security zone for SOVEREIGN vault -->
  <rect x="704" y="240" width="192" height="192"
        fill="#dc2626" fill-opacity="0.08"
        stroke="#fbbf24" stroke-width="2" stroke-dasharray="4 2"
        class="security-zone sovereign" data-clearance="SOVEREIGN">
    <animate attributeName="stroke-dashoffset" from="0" to="12" dur="2s" repeatCount="indefinite"/>
  </rect>

  <!-- Surveillance camera coverage cone -->
  <polygon points="512,192 480,240 544,240"
           fill="#f59e0b" fill-opacity="0.06"
           stroke="none" class="camera-cone"/>
  <circle cx="512" cy="192" r="3" fill="#f59e0b" class="camera-node"/>
</g>
```

**Rendering rules:**
- Visible at Z1, Z2, and Z3 when the security overlay toggle is active (default: on for users with RESTRICTED+ clearance).
- Security zone animations run continuously to draw attention.
- Sovereign zones have a unique animated dashed border (marching ants).
- Camera coverage cones render at Z2+ only.
- This layer respects the user's own clearance level: zones above the user's clearance render as opaque black blocks with a "NO ACCESS" label instead of showing the room interior.

### 9.8 Layer 6: Data Flow Lines

**Purpose:** Renders the visible conduits, data pipelines, and information flow paths between rooms, floors, and external systems. These are the "veins" of the building.

**SVG structure:**

```xml
<g id="layer-6-dataflow" class="layer" data-layer="dataflow" data-layer-index="6">
  <!-- Horizontal data conduit -->
  <g class="conduit" data-conduit-id="conduit-07-002" data-bandwidth="40Gbps">
    <!-- Base line -->
    <line x1="544" y1="368" x2="544" y2="0"
          stroke="#06b6d4" stroke-width="2" stroke-opacity="0.4"/>
    <!-- Animated data packets traveling along the conduit -->
    <circle r="2" fill="#06b6d4" opacity="0.9">
      <animateMotion dur="3s" repeatCount="indefinite"
                     path="M544,368 L544,0"/>
    </circle>
    <circle r="2" fill="#06b6d4" opacity="0.9">
      <animateMotion dur="3s" repeatCount="indefinite" begin="1s"
                     path="M544,368 L544,0"/>
    </circle>
    <circle r="2" fill="#06b6d4" opacity="0.9">
      <animateMotion dur="3s" repeatCount="indefinite" begin="2s"
                     path="M544,368 L544,0"/>
    </circle>
  </g>

  <!-- Vertical riser to basement -->
  <g class="conduit riser" data-conduit-id="conduit-07-001">
    <line x1="96" y1="560" x2="96" y2="768"
          stroke="#22c55e" stroke-width="1.5" stroke-opacity="0.3"
          stroke-dasharray="4 4"/>
    <circle r="1.5" fill="#22c55e">
      <animateMotion dur="2s" repeatCount="indefinite"
                     path="M96,560 L96,768"/>
    </circle>
  </g>

  <!-- Conduit junction node -->
  <circle cx="544" cy="192" r="4"
          fill="#0a0a1a" stroke="#06b6d4" stroke-width="1.5"
          class="junction-node"/>
</g>
```

**Rendering rules:**
- Visible at Z1 (faintly, base lines only), Z2 (full conduits with reduced animation), and Z3 (full detail with all animated data packets).
- Not visible at Z0.
- Data packet animation speed corresponds to the conduit's declared `bandwidth`: higher bandwidth means faster animation.
- Conduits carrying encrypted data render in a different color (`#8b5cf6`, purple) with a lock icon at each endpoint.
- Conduits that cross floor boundaries terminate at the floor edge with a directional arrow and a label indicating the target floor.

### 9.9 Layer 7: Interactive Hotspots

**Purpose:** Renders invisible (or subtly indicated) hit areas that respond to mouse hover, click, and touch events. These are the primary interaction targets for navigation.

**SVG structure:**

```xml
<g id="layer-7-hotspots" class="layer" data-layer="hotspots" data-layer-index="7">
  <!-- Room click target (invisible, full room area) -->
  <rect x="0" y="0" width="320" height="192"
        fill="transparent" cursor="pointer"
        class="hotspot room-hotspot"
        data-room-id="room-07-01" data-action="navigate"
        data-target="room-07-01"/>

  <!-- Door interaction target -->
  <rect x="300" y="80" width="40" height="32"
        fill="transparent" cursor="pointer"
        class="hotspot door-hotspot"
        data-door-id="door-07-001" data-action="toggle-door"/>

  <!-- Elevator shaft click target -->
  <rect x="896" y="240" width="128" height="128"
        fill="transparent" cursor="pointer"
        class="hotspot elevator-hotspot"
        data-shaft-id="shaft-A" data-action="change-floor"/>

  <!-- Equipment interaction target (visible at Z3) -->
  <rect x="496" y="288" width="96" height="64"
        fill="transparent" cursor="pointer"
        class="hotspot equipment-hotspot"
        data-equipment-id="fusion-core-07" data-action="inspect"/>
</g>
```

**Rendering rules:**
- Always present in the DOM but visually transparent.
- At Z2+, hotspots show a subtle highlight on hover: a 1px border with the room's accent color at 30% opacity.
- Hotspot areas must exactly match or slightly exceed the visual bounds of the element they target (4px padding recommended for touch targets).
- Hotspots have a `pointer-events="all"` attribute; all layers beneath them have `pointer-events="none"` for their fill areas (strokes may retain pointer events for visual feedback).
- Touch targets must be at least 44x44 CSS pixels at the current zoom level. If a hotspot would be smaller than this, it is expanded to meet the minimum.

**Event dispatch:** Clicking a hotspot dispatches a custom DOM event:

```javascript
hotspot.addEventListener('click', (e) => {
  const event = new CustomEvent('hic:interact', {
    detail: {
      action: hotspot.dataset.action,
      targetId: hotspot.dataset.roomId || hotspot.dataset.doorId || hotspot.dataset.equipmentId,
      coordinates: { x: e.clientX, y: e.clientY },
      floor: currentFloorId
    },
    bubbles: true
  });
  hotspot.dispatchEvent(event);
});
```

### 9.10 Layer 8: Status Indicators

**Purpose:** Renders real-time status information for rooms, equipment, and systems. This includes power states, occupancy indicators, temperature warnings, and network status.

**SVG structure:**

```xml
<g id="layer-8-status" class="layer" data-layer="status" data-layer-index="8">
  <!-- Room power status indicator (top-right corner of room) -->
  <g class="status-indicator power" data-room-id="room-07-03" transform="translate(816, 8)">
    <circle r="4" fill="#22c55e" class="power-online"/>
    <text x="8" y="4" fill="#22c55e" font-size="6" font-family="monospace">ONLINE</text>
  </g>

  <!-- Occupancy badge (bottom-left corner of room) -->
  <g class="status-indicator occupancy" data-room-id="room-07-05" transform="translate(8, 484)">
    <rect x="0" y="0" width="32" height="14" rx="2" fill="#1e293b" stroke="#0ea5e9" stroke-width="0.5"/>
    <text x="16" y="10" text-anchor="middle" fill="#0ea5e9" font-size="7" font-family="monospace">7/20</text>
  </g>

  <!-- Elevated alert indicator (room border pulse) -->
  <rect x="832" y="0" width="192" height="192"
        fill="none" stroke="#f59e0b" stroke-width="2" opacity="0.6"
        class="alert-border elevated" data-room-id="room-07-04">
    <animate attributeName="opacity" values="0.3;0.8;0.3" dur="2s" repeatCount="indefinite"/>
  </rect>

  <!-- Temperature warning badge -->
  <g class="status-indicator temperature" data-room-id="room-07-07" transform="translate(780, 420)">
    <rect x="0" y="0" width="28" height="14" rx="2" fill="#1e293b" stroke="#06b6d4" stroke-width="0.5"/>
    <text x="14" y="10" text-anchor="middle" fill="#06b6d4" font-size="7" font-family="monospace">COLD</text>
  </g>

  <!-- Isolated power (special indicator) -->
  <g class="status-indicator power-isolated" data-room-id="room-07-07" transform="translate(880, 248)">
    <circle r="4" fill="#f59e0b"/>
    <text x="8" y="4" fill="#f59e0b" font-size="6" font-family="monospace">ISOLATED</text>
  </g>
</g>
```

**Rendering rules:**
- Visible at Z1 (power status only, as colored dots), Z2 (all indicators with labels), and Z3 (full detail with numeric values).
- Not visible at Z0.
- Status indicators update via data binding: the renderer polls or subscribes to a status endpoint and updates indicator attributes in real time.
- Color semantics for power status: `#22c55e` (green) = ONLINE, `#f59e0b` (amber) = ISOLATED or STANDBY, `#ef4444` (red) = OFFLINE, `#6b7280` (gray) = UNKNOWN.
- Occupancy badges change color when above 80% capacity (amber) or at max (red).

### 9.11 Layer 9: Alert Overlays

**Purpose:** Renders critical alerts, emergency notifications, lockdown indicators, and breach warnings. This is the highest-priority visual layer and can obscure all content beneath it when active.

**SVG structure:**

```xml
<g id="layer-9-alerts" class="layer" data-layer="alerts" data-layer-index="9">
  <!-- NOMINAL state: layer is empty -->

  <!-- LOCKDOWN example: full-floor red overlay -->
  <!--
  <rect width="1024" height="768" fill="#dc2626" fill-opacity="0.15"
        class="alert-overlay lockdown" data-alert-type="LOCKDOWN">
    <animate attributeName="fill-opacity" values="0.10;0.20;0.10" dur="1s" repeatCount="indefinite"/>
  </rect>
  <text x="512" y="384" text-anchor="middle" dominant-baseline="middle"
        fill="#dc2626" font-family="monospace" font-size="48" font-weight="900"
        letter-spacing="0.3em" filter="url(#neon-glow)"
        class="alert-text lockdown">
    LOCKDOWN
  </text>
  -->

  <!-- BREACH example: room-specific alert -->
  <!--
  <g class="alert-room-breach" data-room-id="room-07-07">
    <rect x="704" y="240" width="192" height="192"
          fill="#dc2626" fill-opacity="0.25"
          stroke="#ff0055" stroke-width="3">
      <animate attributeName="stroke-width" values="2;4;2" dur="0.5s" repeatCount="indefinite"/>
    </rect>
    <text x="800" y="336" text-anchor="middle" fill="#ff0055"
          font-family="monospace" font-size="14" font-weight="700"
          filter="url(#neon-glow)">
      BREACH DETECTED
    </text>
  </g>
  -->

  <!-- EVACUATION example: directional arrows overlaid on corridors -->
  <!--
  <g class="alert-evacuation">
    <polygon points="480,210 512,200 480,200" fill="#f59e0b" opacity="0.8">
      <animateTransform attributeName="transform" type="translate"
                        values="0,0;40,0;0,0" dur="1s" repeatCount="indefinite"/>
    </polygon>
  </g>
  -->
</g>
```

**Alert types and their visual treatments:**

| Alert Type | Visual Treatment | Animation | Duration |
|------------|-----------------|-----------|----------|
| `NOMINAL` | Layer empty | None | Permanent |
| `ADVISORY` | Thin amber border on affected room | Slow pulse (3s) | Until cleared |
| `ELEVATED` | Amber overlay on affected room, 10% opacity | Medium pulse (2s) | Until cleared |
| `LOCKDOWN` | Full-floor red overlay, large text | Fast pulse (1s) | Until cleared |
| `BREACH` | Room-specific red overlay, flashing border | Rapid flash (0.5s) | Until cleared |
| `EVACUATION` | Directional arrows on corridors, exit highlights | Scrolling arrows | Until cleared |
| `SYSTEM_FAILURE` | Blue overlay on affected room, error text | Glitch effect (random opacity) | Until resolved |

**Rendering rules:**
- Visible at all zoom levels when active. Alerts must be visible regardless of current view.
- Alert overlays have `pointer-events="none"` so users can still interact with underlying hotspots during non-critical alerts.
- During LOCKDOWN and BREACH alerts, hotspot interactions are blocked (all doors show as locked, navigation is restricted to the current floor).
- Alerts are managed by a priority queue: if multiple alerts are active simultaneously, the highest-severity alert controls the overlay. Lower-severity alerts are indicated by status badges on Layer 8.
- An alert sound cue (not part of the SVG) is triggered by the application layer when an alert transitions from NOMINAL to any active state.

### 9.12 Layer Visibility Rules Per Zoom Level

| Layer | Name | Z0 (Building) | Z1 (Floor) | Z2 (Room) | Z3 (Detail) |
|-------|------|:-:|:-:|:-:|:-:|
| 0 | Grid Background | Major only | Full | Full | Full |
| 1 | Structural Walls | Outlines only | Full | Full | Full |
| 2 | Room Fills | Solid blocks | Full | Full | Full |
| 3 | Furniture/Equipment | Hidden | Hidden | Full | Full |
| 4 | Labels and Text | Floor IDs only | Room names | Full + sublabels | Full + coords |
| 5 | Security Overlays | Zone colors | Zone borders | Full + cameras | Full + details |
| 6 | Data Flow Lines | Hidden | Base lines | Full conduits | Full + packets |
| 7 | Interactive Hotspots | Floor strips | Room areas | Room + doors | Room + equipment |
| 8 | Status Indicators | Hidden | Power dots | All indicators | Full detail |
| 9 | Alert Overlays | Always (if active) | Always (if active) | Always (if active) | Always (if active) |

**Transition behavior:** When zooming between levels, layers fade in/out over 200ms with CSS opacity transitions. This prevents jarring visual pops.

### 9.13 Layer Rendering Order

Layers are rendered strictly in ascending numeric order. Within each layer, elements are rendered in document order (first in the SVG source = rendered first = appears behind later elements).

```
Screen (top)
  |
  |  Layer 9: Alert Overlays        (topmost, blocks everything when active)
  |  Layer 8: Status Indicators     (badges and dots above content)
  |  Layer 7: Interactive Hotspots  (transparent hit areas)
  |  Layer 6: Data Flow Lines       (animated conduits)
  |  Layer 5: Security Overlays     (zone highlights)
  |  Layer 4: Labels and Text       (room names, floor IDs)
  |  Layer 3: Furniture/Equipment   (desks, servers, terminals)
  |  Layer 2: Room Fills            (colored room interiors)
  |  Layer 1: Structural Walls      (walls, doors, corridors)
  |  Layer 0: Grid Background       (dark grid, always at bottom)
  |
Background
```

**Z-index enforcement:** SVG does not support CSS `z-index`. Rendering order is determined exclusively by DOM order. The renderer must ensure that layer `<g>` groups appear in the SVG in ascending order (Layer 0 first, Layer 9 last). If elements are dynamically added, they must be appended to the correct layer group, never to the root `<svg>` element.

### 9.14 Layer Toggle Controls

Users can toggle individual layers on and off via the Layer Control Panel (LCP), a UI sidebar element. Each toggle controls the `visibility` CSS property of the corresponding `<g>` layer group.

```javascript
const LAYER_CONFIG = [
  { index: 0, name: "Grid",       icon: "grid",     default: true,  locked: true  },
  { index: 1, name: "Walls",      icon: "walls",    default: true,  locked: true  },
  { index: 2, name: "Rooms",      icon: "rooms",    default: true,  locked: false },
  { index: 3, name: "Equipment",  icon: "cpu",      default: true,  locked: false },
  { index: 4, name: "Labels",     icon: "type",     default: true,  locked: false },
  { index: 5, name: "Security",   icon: "shield",   default: true,  locked: false },
  { index: 6, name: "Data Flow",  icon: "activity", default: false, locked: false },
  { index: 7, name: "Hotspots",   icon: "cursor",   default: true,  locked: true  },
  { index: 8, name: "Status",     icon: "monitor",  default: true,  locked: false },
  { index: 9, name: "Alerts",     icon: "alert",    default: true,  locked: true  }
];

function toggleLayer(index, visible) {
  const layer = document.getElementById(`layer-${index}-${LAYER_CONFIG[index].name.toLowerCase()}`);
  if (layer && !LAYER_CONFIG[index].locked) {
    layer.style.visibility = visible ? 'visible' : 'hidden';
  }
}
```

**Locked layers** (Grid, Walls, Hotspots, Alerts) cannot be toggled off by the user. They are essential for navigation and safety.

**Layer presets:** Common layer combinations can be saved and recalled:

| Preset | Active Layers | Use Case |
|--------|---------------|----------|
| Standard | 0-5, 7-9 | Default view, data flow hidden |
| Operations | 0-9 | Full operational view |
| Security Audit | 0, 1, 2, 5, 7, 9 | Security assessment |
| Structural | 0, 1, 4 | Architecture review |
| Data Analysis | 0, 1, 2, 6, 7 | Data flow tracing |
| Minimal | 0, 1, 2, 7 | Clean navigation |

### 9.15 Complete SVG Group Structure

The following shows the complete SVG document structure with all layer groups, demonstrating how they nest within the root `<svg>` element:

```xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 768"
     data-floor-id="07" data-schema-version="1.0.0">

  <defs>
    <!-- All filter definitions, patterns, gradients, clip paths -->
    <filter id="neon-glow">...</filter>
    <filter id="pulse-glow">...</filter>
    <pattern id="scanlines">...</pattern>
    <pattern id="grid-16">...</pattern>
    <pattern id="grid-64">...</pattern>
    <linearGradient id="glow-strip-horiz">...</linearGradient>
    <clipPath id="floor-bounds"><rect width="1024" height="768"/></clipPath>
  </defs>

  <!-- All content clipped to floor bounds -->
  <g clip-path="url(#floor-bounds)">

    <!-- FOUNDATION TIER -->
    <g id="layer-0-grid"      class="layer" data-layer="grid"      data-layer-index="0">...</g>
    <g id="layer-1-walls"     class="layer" data-layer="walls"     data-layer-index="1">...</g>
    <g id="layer-2-rooms"     class="layer" data-layer="rooms"     data-layer-index="2">...</g>

    <!-- CONTENT TIER -->
    <g id="layer-3-furniture" class="layer" data-layer="furniture" data-layer-index="3">...</g>
    <g id="layer-4-labels"    class="layer" data-layer="labels"    data-layer-index="4">...</g>
    <g id="layer-5-security"  class="layer" data-layer="security"  data-layer-index="5">...</g>
    <g id="layer-6-dataflow"  class="layer" data-layer="dataflow"  data-layer-index="6">...</g>

    <!-- INTERACTIVE TIER -->
    <g id="layer-7-hotspots"  class="layer" data-layer="hotspots"  data-layer-index="7">...</g>
    <g id="layer-8-status"    class="layer" data-layer="status"    data-layer-index="8">...</g>
    <g id="layer-9-alerts"    class="layer" data-layer="alerts"    data-layer-index="9">...</g>

  </g>

  <!-- Post-processing overlays (outside clip, always full viewport) -->
  <rect width="1024" height="768" fill="url(#scanlines)" pointer-events="none" opacity="0.4"/>

  <!-- CRT vignette effect -->
  <radialGradient id="vignette">
    <stop offset="60%" stop-color="transparent"/>
    <stop offset="100%" stop-color="rgba(0,0,0,0.4)"/>
  </radialGradient>
  <rect width="1024" height="768" fill="url(#vignette)" pointer-events="none"/>
</svg>
```

---

## Cross-Agent Integration Notes

### Blueprint-to-Renderer Pipeline

1. **Load:** Parse `floor-{ID}.blueprint.json`, validate against schema (Agent 7 rules).
2. **Transform:** Convert all room bounds and wall coordinates to viewport coordinates using the current pan/zoom state (Agent 8 transforms).
3. **Render:** For each layer (Agent 9), iterate over the relevant blueprint data and emit SVG elements into the corresponding `<g>` group.
4. **Animate:** Start CSS and SMIL animations for active elements (glow pulses, data packets, alert flashes).
5. **Bind:** Attach event listeners to Layer 7 hotspots for navigation dispatch.
6. **Update:** Subscribe to status feeds and update Layer 8 indicators and Layer 9 alerts in real time.

### Performance Considerations

- **Culling:** At Z2 and Z3, rooms outside the current viewport are not rendered. The renderer uses the viewport inverse transform (Agent 8, section 8.6) to determine which rooms intersect the visible area.
- **Animation budgets:** At most 50 simultaneous SVG animations are active at any time. When the budget is exceeded, the lowest-priority animations (data packets on distant conduits) are paused.
- **Layer caching:** Foundation-tier layers (0-2) are rasterized into an offscreen canvas when the floor data has not changed, avoiding redundant SVG layout calculations on every frame.
- **Lazy loading:** Blueprint JSON for non-visible floors is not loaded until the user navigates to that floor or enters Building View (Z0).

### File Size Budgets

| File Type | Target Size | Maximum Size |
|-----------|-------------|--------------|
| Blueprint JSON (per floor) | < 50 KB | 200 KB |
| Rendered SVG (per floor) | < 100 KB | 500 KB |
| Thumbnail SVG (per floor) | < 10 KB | 25 KB |
| Total loaded assets (single floor view) | < 300 KB | 1 MB |
| Total loaded assets (building view) | < 2 MB | 5 MB |

---

*This specification is maintained by the HIC Spatial Data Working Group. All amendments require approval from at minimum two of the three responsible agents (7, 8, 9). Schema version changes that break backward compatibility require a major version bump and a 30-day deprecation period for the previous version.*
