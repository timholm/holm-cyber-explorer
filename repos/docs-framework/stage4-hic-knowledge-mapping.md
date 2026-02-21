# STAGE 4: HOLM INTELLIGENCE COMPLEX -- KNOWLEDGE MAPPING

## Agents 13-15: Room-to-Document Binding, Clearance Levels, Archive Placement, Memory Vaults

**Document ID:** STAGE4-HIC-KNOWLEDGE-MAPPING
**Version:** 1.0.0
**Date:** 2026-02-17
**Status:** Ratified
**Classification:** Specialized Systems -- These articles provide the complete spatial knowledge-mapping specification for the Holm Intelligence Complex (HIC). They define how every document in the institution occupies a physical room in the skyscraper, how clearance levels govern access, and how the basement archive and memory vault systems preserve institutional knowledge across time. They extend the HIC structural architecture (HIC-001 through HIC-006) into a fully operational document-placement system.

---

## How to Read This Document

This document contains three agent specifications that belong to Stage 4 of the HIC interface layer. Where HIC-001 through HIC-003 defined the skyscraper's structural geometry and HIC-004 through HIC-006 defined navigation and room rendering, Agents 13 through 15 define the *knowledge layer* -- the system that turns empty rooms into living document containers.

Documents in the Holm Intelligence Complex are not in folders. They are not in directories. They are not in trees. They are IN rooms. A document exists in exactly one room in the skyscraper. That room has a floor number, a wing, a door, a shelf position, and a clearance level. When a user navigates to a document, they walk through the skyscraper to get there. The spatial metaphor is not decoration -- it is the primary organizational principle.

Agent 13 specifies how documents bind to rooms. Agent 14 specifies the five-tier clearance system that governs who can enter which rooms and how locked rooms appear from outside. Agent 15 specifies the basement archive and memory vault systems that handle cold storage, institutional memory, and anonymous submissions.

If you are implementing the HIC, read Agent 13 first. It defines the data model. Agent 14 layers access control on top of that model. Agent 15 extends the model downward into the basement floors where time-sensitive and permanent records live.

---

---

# Agent 13: Room-to-Document Binding

**Document ID:** HIC-013
**Agent:** 13 -- Knowledge Mapping: Binding
**Status:** Phase 1
**Date:** 2026-02-17
**Depends On:** HIC-001 (Skyscraper Master Architecture), HIC-002 (Floor Plan Architecture), HIC-005 (Room Rendering Engine)
**Depended Upon By:** HIC-014 (Clearance Levels), HIC-015 (Archive Placement & Memory Vaults), HIC-016 (Art Direction), all components that render or retrieve documents

---

## 1. Purpose

This article defines the complete specification for binding documents to rooms within the Holm Intelligence Complex. Every document in the institution -- every article, every operational procedure, every decision record, every memo -- exists in exactly one room. That room is the document's *canonical location* in the spatial interface. The binding is not metadata attached to the document; it is the document's address in the physical world of the skyscraper.

The binding system answers three questions: Where does this document live? How do I get to it? What else is in the same room?

## 2. Scope

This article covers:

- The binding schema that connects documents to rooms.
- Room capacity limits per room type.
- Document placement rules: which types of documents go in which types of rooms.
- The room manifest file format.
- The cross-reference system for documents that reference each other across rooms.
- The search index specification for finding documents without knowing their room.
- The document lifecycle in spatial terms: intake, processing, and final placement.
- The complete domain-to-floor-to-room mapping table.
- Example manifest files.

This article does not cover the visual rendering of rooms (see HIC-005), the clearance system that restricts access to rooms (see HIC-014), or the archive and vault systems in the basement floors (see HIC-015).

## 3. The Binding Schema

Every document-to-room binding is recorded as a structured object. This is the canonical schema.

### 3.1 Binding Record Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "DocumentRoomBinding",
  "description": "Binds a single document to a single room in the HIC skyscraper.",
  "type": "object",
  "required": [
    "document_id",
    "room_id",
    "shelf_position",
    "classification",
    "format",
    "bound_at",
    "bound_by"
  ],
  "properties": {
    "document_id": {
      "type": "string",
      "description": "The unique institutional document identifier (e.g., SEC-004, GOV-001, D6-009).",
      "pattern": "^[A-Z0-9]+-[A-Z]?[0-9]+$"
    },
    "room_id": {
      "type": "string",
      "description": "The room's unique spatial identifier. Format: F{floor}-{wing}-R{room_number}.",
      "pattern": "^[FB][0-9]+-[NESW]-R[0-9]+$"
    },
    "shelf_position": {
      "type": "object",
      "description": "The physical position of the document within the room.",
      "required": ["shelf", "slot"],
      "properties": {
        "shelf": {
          "type": "integer",
          "minimum": 1,
          "maximum": 20,
          "description": "Shelf number, numbered bottom-to-top."
        },
        "slot": {
          "type": "integer",
          "minimum": 1,
          "maximum": 50,
          "description": "Position within the shelf, numbered left-to-right."
        }
      }
    },
    "classification": {
      "type": "string",
      "enum": ["PUBLIC", "INTERNAL", "CONFIDENTIAL", "SECRET", "COSMIC"],
      "description": "The document's clearance classification."
    },
    "format": {
      "type": "string",
      "enum": ["article", "procedure", "decision_record", "memo", "blueprint", "manifest", "transcript", "log", "report", "template", "index"],
      "description": "The document's format type, which determines its visual appearance on the shelf."
    },
    "bound_at": {
      "type": "string",
      "format": "date-time",
      "description": "ISO 8601 timestamp of when the binding was created."
    },
    "bound_by": {
      "type": "string",
      "description": "Identity of the operator or automation that created the binding."
    },
    "title": {
      "type": "string",
      "description": "Human-readable document title."
    },
    "domain": {
      "type": "integer",
      "minimum": 1,
      "maximum": 20,
      "description": "The institutional domain number."
    },
    "stage": {
      "type": "integer",
      "minimum": 1,
      "maximum": 5,
      "description": "The institutional stage number."
    },
    "cross_references": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["target_document_id", "target_room_id", "relationship"],
        "properties": {
          "target_document_id": {
            "type": "string"
          },
          "target_room_id": {
            "type": "string"
          },
          "relationship": {
            "type": "string",
            "enum": ["depends_on", "depended_upon_by", "references", "supersedes", "superseded_by", "related"]
          }
        }
      },
      "description": "List of cross-references to documents in other rooms."
    },
    "lifecycle_state": {
      "type": "string",
      "enum": ["intake", "processing", "placed", "archived", "sealed"],
      "description": "Current lifecycle state of the document within the spatial system."
    }
  }
}
```

### 3.2 The One-Document-One-Room Rule

A document exists in exactly one room. There are no duplicates. If document SEC-004 is in room F3-N-R12, there is no copy of SEC-004 anywhere else in the skyscraper. If a user in another room needs to reference SEC-004, they follow a cross-reference link -- visually rendered as a glowing conduit line that traces a path through the corridors to room F3-N-R12.

**Rationale:** Duplicate documents create synchronization problems. When a document is updated, every copy must be updated. In a spatial interface, duplication would mean the same document occupies multiple rooms simultaneously, which breaks the spatial metaphor. A document has one home. References point to that home.

### 3.3 Cross-Reference Conduits

When document A in Room X references document B in Room Y, the cross-reference is rendered as a visible neon conduit line. The conduit originates from document A's shelf position, exits the room through the wall, traces through corridors and elevator shafts as necessary, and terminates at document B's shelf position in Room Y. The user can follow the conduit visually or click it to teleport to the target room.

Cross-reference conduit colors indicate the relationship type:

| Relationship | Conduit Color | Line Style |
|---|---|---|
| `depends_on` | Neon amber | Solid, pulsing |
| `depended_upon_by` | Neon amber | Dashed, steady |
| `references` | Neon white | Thin solid |
| `supersedes` | Neon red | Thick solid |
| `superseded_by` | Dim red | Thick dashed |
| `related` | Neon cyan | Dotted |

## 4. Room Types and Capacity Limits

The skyscraper contains four room types, each with a defined capacity limit. The capacity limit is the maximum number of document bindings a room can hold.

| Room Type | Capacity | Typical Use | Shelf Configuration | Visual Footprint |
|---|---|---|---|---|
| **Office** | 50 documents | Active operational documents, current procedures, working references | 5 shelves x 10 slots | Small room, single desk, wall-mounted shelves |
| **Library** | 200 documents | Domain philosophy collections, educational materials, reference sets | 10 shelves x 20 slots | Medium room, reading table, floor-to-ceiling shelves |
| **Vault** | 500 documents | Classified materials, decision records, security documentation | 10 shelves x 50 slots | Large reinforced room, heavy door, sealed cabinets |
| **Archive** | 5000 documents | Historical records, cold storage, version histories, bulk logs | 20 shelves x 250 slots | Cavernous basement room, industrial shelving, dim lighting |

When a room reaches 90% capacity, the interface renders a visual warning: the room's ambient lighting shifts from its normal hue to a warm amber, and a capacity indicator appears above the door. At 100% capacity, no further bindings can be created in that room until a document is relocated or the room is expanded (which requires a governance decision per GOV-001).

## 5. Document Placement Rules

Documents are placed in rooms based on their domain, format, and classification. The placement rules are deterministic: given a document's metadata, there is exactly one correct room (or a small set of acceptable rooms) for it.

### 5.1 Placement Priority

1. **Classification overrides all other rules.** A COSMIC document goes to a COSMIC-cleared vault regardless of domain or format.
2. **Domain determines floor.** Each domain is assigned to a specific floor (see Section 8).
3. **Format determines room type.** Articles and procedures go to offices or libraries. Decision records go to vaults. Logs and historical records go to archives.
4. **Stage determines wing.** Stage 2 (philosophy) documents cluster in the East wing. Stage 3 (operational) documents cluster in the West wing. Stage 4 (specialized) documents cluster in the North wing. Stage 1 (framework) and Stage 5 (meta) documents occupy the South wing.

### 5.2 Format-to-Room-Type Mapping

| Document Format | Primary Room Type | Fallback Room Type |
|---|---|---|
| `article` | Library | Office |
| `procedure` | Office | Library |
| `decision_record` | Vault | -- |
| `memo` | Office | -- |
| `blueprint` | Library | Vault |
| `manifest` | Office | -- |
| `transcript` | Vault | Archive |
| `log` | Archive | Vault |
| `report` | Library | Office |
| `template` | Office | Library |
| `index` | Office (lobby room on floor) | -- |

### 5.3 Placement Prohibition Rules

- A PUBLIC document MUST NOT be placed in a vault.
- A COSMIC document MUST NOT be placed in an office.
- A `log` format document MUST NOT be placed in an office (logs are high-volume; offices are low-capacity).
- A `decision_record` MUST be placed in a vault (decisions require permanent protection).
- Documents with `lifecycle_state: "sealed"` MUST be in a vault or archive. Sealed documents cannot be moved.

## 6. Room Manifest File Format

Every room in the skyscraper has a manifest file that lists every document currently bound to that room. The manifest is the room's source of truth.

### 6.1 Manifest Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "RoomManifest",
  "description": "The canonical manifest for a single room in the HIC skyscraper.",
  "type": "object",
  "required": [
    "room_id",
    "room_type",
    "floor",
    "wing",
    "domain",
    "clearance_level",
    "capacity",
    "documents",
    "manifest_version",
    "last_updated"
  ],
  "properties": {
    "room_id": {
      "type": "string",
      "pattern": "^[FB][0-9]+-[NESW]-R[0-9]+$"
    },
    "room_type": {
      "type": "string",
      "enum": ["office", "library", "vault", "archive"]
    },
    "floor": {
      "type": "string",
      "description": "Floor identifier. F1-F20 for above-ground floors. B1-B3 for basement floors."
    },
    "wing": {
      "type": "string",
      "enum": ["N", "E", "S", "W"],
      "description": "Cardinal wing within the floor."
    },
    "domain": {
      "type": "integer",
      "minimum": 0,
      "maximum": 20,
      "description": "Primary domain served by this room. 0 for cross-domain or shared rooms."
    },
    "clearance_level": {
      "type": "integer",
      "minimum": 0,
      "maximum": 4,
      "description": "Minimum clearance required to enter this room."
    },
    "capacity": {
      "type": "integer",
      "description": "Maximum number of documents this room can hold."
    },
    "occupied": {
      "type": "integer",
      "description": "Current number of documents in this room."
    },
    "documents": {
      "type": "array",
      "items": {
        "$ref": "#/$defs/DocumentEntry"
      }
    },
    "manifest_version": {
      "type": "string",
      "pattern": "^[0-9]+\\.[0-9]+\\.[0-9]+$"
    },
    "last_updated": {
      "type": "string",
      "format": "date-time"
    },
    "room_name": {
      "type": "string",
      "description": "Optional human-readable name for the room (e.g., 'Security Philosophy Library')."
    },
    "ambient_color": {
      "type": "string",
      "description": "Hex color code for the room's default neon ambient lighting.",
      "pattern": "^#[0-9a-fA-F]{6}$"
    }
  },
  "$defs": {
    "DocumentEntry": {
      "type": "object",
      "required": ["document_id", "title", "shelf", "slot", "format", "classification"],
      "properties": {
        "document_id": { "type": "string" },
        "title": { "type": "string" },
        "shelf": { "type": "integer" },
        "slot": { "type": "integer" },
        "format": { "type": "string" },
        "classification": { "type": "string" },
        "cross_references": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "target_document_id": { "type": "string" },
              "target_room_id": { "type": "string" },
              "relationship": { "type": "string" }
            }
          }
        }
      }
    }
  }
}
```

### 6.2 Manifest File Location

Each room manifest is stored at the following path within the HIC data directory:

```
hic-data/floors/{floor_id}/wings/{wing}/rooms/{room_id}/room-manifest.json
```

Example: `hic-data/floors/F3/wings/N/rooms/F3-N-R12/room-manifest.json`

### 6.3 Manifest Integrity

Every manifest file includes a SHA-256 hash of its own contents (excluding the hash field) stored in a companion file:

```
hic-data/floors/F3/wings/N/rooms/F3-N-R12/room-manifest.json.sha256
```

On every room load, the rendering engine verifies the manifest hash. If the hash does not match, the room renders in an error state: the room appears with red flickering lights and a warning overlay reading `MANIFEST INTEGRITY FAILURE`. The user is notified. The room's contents are not displayed until the integrity issue is resolved.

## 7. Search Index Specification

A user must be able to find a document without knowing its room. The search index provides this capability.

### 7.1 Global Search Index Schema

The global search index is a single JSON file that maps every document in the skyscraper to its room. It is rebuilt whenever a document is bound, moved, or unbound.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "HICSearchIndex",
  "type": "object",
  "required": ["version", "built_at", "entries"],
  "properties": {
    "version": { "type": "string" },
    "built_at": { "type": "string", "format": "date-time" },
    "entry_count": { "type": "integer" },
    "entries": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["document_id", "room_id", "floor", "wing", "domain", "title", "format", "classification", "keywords"],
        "properties": {
          "document_id": { "type": "string" },
          "room_id": { "type": "string" },
          "floor": { "type": "string" },
          "wing": { "type": "string" },
          "domain": { "type": "integer" },
          "title": { "type": "string" },
          "format": { "type": "string" },
          "classification": { "type": "string" },
          "keywords": {
            "type": "array",
            "items": { "type": "string" },
            "description": "Extracted keywords for full-text search."
          },
          "abstract": {
            "type": "string",
            "description": "First 500 characters of the document's Purpose section."
          }
        }
      }
    }
  }
}
```

### 7.2 Search Behavior

When a user searches from any location in the skyscraper, the interface presents results as a list of *rooms* with highlighted document entries. Each result shows:

1. The document title and ID.
2. The room name and floor/wing location.
3. The clearance level (with appropriate color indicator).
4. A "Navigate" button that triggers the elevator/corridor animation to the target room.

If the user lacks clearance to access the target room, the result still appears in the list, but the document title is replaced with `[CLASSIFIED -- Level {N} Required]` and the "Navigate" button is replaced with `[ACCESS DENIED]`. The room location is still visible -- the user knows the document exists and where it is, but cannot read its title or contents.

### 7.3 Search Index File Location

```
hic-data/search/global-index.json
hic-data/search/global-index.json.sha256
hic-data/search/domain-index-{domain_number}.json
```

Domain-specific indexes exist for each of the 20 domains plus a domain-0 index for cross-domain documents. These are subsets of the global index, used for faster searching when the user is already on a domain floor.

## 8. Document Lifecycle in Spatial Terms

Every document passes through a spatial lifecycle within the skyscraper. The lifecycle has five states.

### 8.1 Lifecycle States

| State | Spatial Location | Description |
|---|---|---|
| **Intake** | Ground floor (F1), Intake Wing (S) | Document has arrived at the institution but has not been classified or assigned a room. It exists in the Intake Lobby -- a large open room with temporary holding shelves. |
| **Processing** | Ground floor (F1), Processing Wing (W) | Document is being classified, formatted, reviewed, and assigned a clearance level. It sits in a Processing Room -- a brightly lit workspace with inspection tables. |
| **Placed** | Assigned room on assigned floor | Document has been bound to its permanent room. It is on a shelf, in a slot, in a specific room. This is the document's canonical location. |
| **Archived** | Basement floors (B1, B2, B3) | Document has been moved to long-term storage. It remains accessible but is no longer in its original above-ground room. Its original room retains a ghost binding -- a dim placeholder on the shelf that reads "Archived to {archive_room_id}" with a follow link. |
| **Sealed** | Vault or Memory Vault (basement) | Document has been permanently sealed. It cannot be moved, modified, or declassified. It exists in a vault with a time-lock indicator showing the seal date. |

### 8.2 Lifecycle Transitions

```
Intake --> Processing --> Placed --> Archived --> Sealed
                              |          ^
                              |          |
                              +----------+
                              (direct archive for bulk imports)
```

Every transition generates an entry in the document's movement log:

```json
{
  "document_id": "SEC-004",
  "transition": "placed_to_archived",
  "from_room": "F3-N-R12",
  "to_room": "B1-E-R003",
  "timestamp": "2026-08-15T14:30:00Z",
  "operator": "holm-operator-1",
  "reason": "Document superseded by SEC-004 v2.0.0. Original preserved in archive."
}
```

### 8.3 Intake Room Specification

The Intake Room is located at F1-S-R001. It has a capacity of 200 documents (temporary holding). Documents in the Intake Room are rendered with a yellow pulsing border to indicate they are unprocessed. The Intake Room has no clearance requirement -- it is PUBLIC by default, because unclassified documents have not yet been assigned a classification.

### 8.4 Processing Room Specification

The Processing Room is located at F1-W-R001 through F1-W-R005 (five processing stations). Each processing room has a capacity of 50 documents. Documents in processing rooms are rendered with an orange rotating border. A processing room requires Level 1 (INTERNAL) clearance minimum because document classification is being determined and partially-classified documents may be present.

## 9. Complete Domain-to-Floor-to-Room Mapping

The skyscraper has 20 above-ground floors (F1 through F20), one penthouse floor (F21), and three basement floors (B1, B2, B3). Each of the 20 institutional domains occupies one floor. The ground floor (F1) serves as the shared intake, processing, and lobby space. The penthouse (F21) is reserved for COSMIC-level materials.

### 9.1 Floor Assignment Table

| Floor | Domain | Domain Name | Primary Room Types | Default Clearance | Neon Accent Color |
|---|---|---|---|---|---|
| F1 | 0 | LOBBY / INTAKE / PROCESSING | Office, Library | Level 0 (PUBLIC) | Neon Green (#00ff41) |
| F2 | 1 | CONSTITUTION & PHILOSOPHY | Library, Vault | Level 1 (INTERNAL) | White (#ffffff) |
| F3 | 2 | GOVERNANCE & AUTHORITY | Office, Vault | Level 1 (INTERNAL) | Gold (#ffd700) |
| F4 | 3 | SECURITY & INTEGRITY | Vault, Office | Level 2 (CONFIDENTIAL) | Neon Blue (#0080ff) |
| F5 | 4 | INFRASTRUCTURE & POWER | Office, Library | Level 1 (INTERNAL) | Neon Orange (#ff6600) |
| F6 | 5 | PLATFORM & CORE SYSTEMS | Office, Library | Level 1 (INTERNAL) | Neon Teal (#00cccc) |
| F7 | 6 | DATA & ARCHIVES | Library, Archive | Level 1 (INTERNAL) | Neon Violet (#9900ff) |
| F8 | 7 | INTELLIGENCE & ANALYSIS | Vault, Office | Level 2 (CONFIDENTIAL) | Neon Magenta (#ff00ff) |
| F9 | 8 | AUTOMATION & AGENTS | Office, Library | Level 1 (INTERNAL) | Neon Yellow (#ccff00) |
| F10 | 9 | EDUCATION & TRAINING | Library, Office | Level 0 (PUBLIC) | Neon Green (#33ff77) |
| F11 | 10 | USER OPERATIONS | Office, Library | Level 1 (INTERNAL) | Neon Cyan (#00ffff) |
| F12 | 11 | ADMINISTRATION | Office, Vault | Level 1 (INTERNAL) | Neon Pink (#ff66b2) |
| F13 | 12 | DISASTER RECOVERY | Vault, Office | Level 2 (CONFIDENTIAL) | Neon Red (#ff3333) |
| F14 | 13 | EVOLUTION & ADAPTATION | Library, Office | Level 1 (INTERNAL) | Neon Amber (#ffbf00) |
| F15 | 14 | RESEARCH & THEORY | Library, Office | Level 1 (INTERNAL) | Neon White (#e0e0ff) |
| F16 | 15 | ETHICS & SAFEGUARDS | Vault, Library | Level 2 (CONFIDENTIAL) | Neon Indigo (#4400ff) |
| F17 | 16 | INTERFACE & NAVIGATION | Office, Library | Level 1 (INTERNAL) | Neon Aqua (#00ffcc) |
| F18 | 17 | FEDERATION | Office, Vault | Level 2 (CONFIDENTIAL) | Neon Lime (#88ff00) |
| F19 | 18 | IMPORT & QUARANTINE | Office, Vault | Level 2 (CONFIDENTIAL) | Neon Coral (#ff4466) |
| F20 | 19 | QUALITY ASSURANCE | Office, Library | Level 1 (INTERNAL) | Neon Sky (#00aaff) |
| F21 | 20 | INSTITUTIONAL MEMORY | Vault, Library | Level 1 (INTERNAL) | Neon Gold (#ffcc00) |
| **Penthouse** | -- | COSMIC CLEARANCE ONLY | Vault | Level 4 (COSMIC) | Neon Red Pulse (#ff0000) |
| B1 | -- | RECENT ARCHIVE | Archive | Level 1 (INTERNAL) | Dim Cyan (#004444) |
| B2 | -- | HISTORICAL ARCHIVE | Archive | Level 1 (INTERNAL) | Dim Blue (#000044) |
| B3 | -- | DEEP COLD ARCHIVE | Archive, Vault | Level 2 (CONFIDENTIAL) | Dim Violet (#110022) |

### 9.2 Room Count Per Floor

Each above-ground floor (F2 through F21) contains the following standard room layout:

| Wing | Room Range | Room Count | Typical Room Types |
|---|---|---|---|
| North (N) | R001 - R010 | 10 | Stage 4 specialized documents |
| East (E) | R001 - R008 | 8 | Stage 2 philosophy documents |
| South (S) | R001 - R006 | 6 | Stage 1 framework & Stage 5 meta documents |
| West (W) | R001 - R010 | 10 | Stage 3 operational documents |
| **Total per floor** | | **34** | |

The ground floor (F1) has a modified layout:

| Wing | Room Range | Purpose |
|---|---|---|
| North (N) | R001 - R004 | Public Reading Rooms (Level 0) |
| East (E) | R001 - R003 | Information Desk, Directory Kiosks |
| South (S) | R001 | Intake Room (incoming documents) |
| West (W) | R001 - R005 | Processing Rooms |

Basement floors have a different, denser layout (see Agent 15, Section 4).

### 9.3 Room Naming Convention

Every room has a unique identifier and an optional human-readable name.

**Identifier format:** `{floor}-{wing}-R{number}`

Examples:
- `F4-N-R003` -- Security & Integrity, North Wing, Room 3
- `B2-E-R017` -- Historical Archive, East Wing, Room 17
- `F1-S-R001` -- Ground floor Intake Room

**Human-readable names** are optional and stored in the room manifest. Examples:
- `F2-E-R001` = "Constitutional Foundation Library"
- `F4-N-R001` = "Threat Model Vault"
- `F8-W-R003` = "Intelligence Cycle Procedures Office"

## 10. Example Manifest Files

### 10.1 Example: Security Philosophy Library

```json
{
  "room_id": "F4-E-R001",
  "room_type": "library",
  "floor": "F4",
  "wing": "E",
  "domain": 3,
  "clearance_level": 2,
  "capacity": 200,
  "occupied": 4,
  "room_name": "Security Philosophy Library",
  "ambient_color": "#0080ff",
  "manifest_version": "1.0.0",
  "last_updated": "2026-02-17T10:00:00Z",
  "documents": [
    {
      "document_id": "SEC-001",
      "title": "Threat Model and Security Philosophy",
      "shelf": 1,
      "slot": 1,
      "format": "article",
      "classification": "CONFIDENTIAL",
      "cross_references": [
        {
          "target_document_id": "ETH-001",
          "target_room_id": "F2-E-R001",
          "relationship": "depends_on"
        },
        {
          "target_document_id": "SEC-002",
          "target_room_id": "F4-W-R001",
          "relationship": "depended_upon_by"
        }
      ]
    },
    {
      "document_id": "SEC-009",
      "title": "Trust Boundary Definitions",
      "shelf": 1,
      "slot": 2,
      "format": "article",
      "classification": "CONFIDENTIAL",
      "cross_references": []
    },
    {
      "document_id": "SEC-010",
      "title": "Cryptographic Principles and Standards",
      "shelf": 1,
      "slot": 3,
      "format": "article",
      "classification": "CONFIDENTIAL",
      "cross_references": [
        {
          "target_document_id": "SEC-003",
          "target_room_id": "F4-W-R002",
          "relationship": "related"
        }
      ]
    },
    {
      "document_id": "SEC-014",
      "title": "Security Documentation and Classification",
      "shelf": 1,
      "slot": 4,
      "format": "article",
      "classification": "CONFIDENTIAL",
      "cross_references": []
    }
  ]
}
```

### 10.2 Example: Governance Decision Vault

```json
{
  "room_id": "F3-N-R001",
  "room_type": "vault",
  "floor": "F3",
  "wing": "N",
  "domain": 2,
  "clearance_level": 2,
  "capacity": 500,
  "occupied": 2,
  "room_name": "Governance Decision Vault",
  "ambient_color": "#ffd700",
  "manifest_version": "1.0.0",
  "last_updated": "2026-02-17T10:00:00Z",
  "documents": [
    {
      "document_id": "GOV-002",
      "title": "Decision Record Keeping",
      "shelf": 1,
      "slot": 1,
      "format": "decision_record",
      "classification": "CONFIDENTIAL",
      "cross_references": [
        {
          "target_document_id": "GOV-001",
          "target_room_id": "F3-E-R001",
          "relationship": "depends_on"
        }
      ]
    },
    {
      "document_id": "GOV-008",
      "title": "Decision-Making Frameworks",
      "shelf": 1,
      "slot": 2,
      "format": "decision_record",
      "classification": "CONFIDENTIAL",
      "cross_references": [
        {
          "target_document_id": "GOV-002",
          "target_room_id": "F3-N-R001",
          "relationship": "depends_on"
        }
      ]
    }
  ]
}
```

### 10.3 Example: Intake Room (Ground Floor)

```json
{
  "room_id": "F1-S-R001",
  "room_type": "office",
  "floor": "F1",
  "wing": "S",
  "domain": 0,
  "clearance_level": 0,
  "capacity": 200,
  "occupied": 0,
  "room_name": "Document Intake Lobby",
  "ambient_color": "#00ff41",
  "manifest_version": "1.0.0",
  "last_updated": "2026-02-17T10:00:00Z",
  "documents": []
}
```

---

---

# Agent 14: Clearance Levels

**Document ID:** HIC-014
**Agent:** 14 -- Knowledge Mapping: Clearance
**Status:** Phase 1
**Date:** 2026-02-17
**Depends On:** HIC-013 (Room-to-Document Binding), HIC-005 (Room Rendering Engine), SEC-001 (Threat Model and Security Philosophy), SEC-002 (Access Control Procedures)
**Depended Upon By:** HIC-015 (Archive Placement & Memory Vaults), HIC-016 (Art Direction), HIC-017 (Lighting & Atmosphere)

---

## 1. Purpose

This article defines the five-tier clearance system that governs access to rooms in the Holm Intelligence Complex. Every room has a clearance level. Every user has a clearance level. If the user's clearance is lower than the room's clearance, the user cannot enter. The room's door is locked, the contents are hidden, and the room is rendered in a locked-out visual state that communicates "this room exists, but you cannot see inside."

The clearance system is not merely access control. It is a *visual language*. Each clearance level has a distinct neon color, door style, lock icon, and ambient glow. A user walking through the skyscraper can immediately identify the classification of every room they pass by its visual treatment, without reading any text.

## 2. The Five Clearance Tiers

### 2.1 Tier Definitions

```json
{
  "clearance_tiers": [
    {
      "level": 0,
      "name": "PUBLIC",
      "code": "PUB",
      "description": "Open to all users. No credentials required. Lobby, common areas, public-facing documentation, educational materials.",
      "neon_color": "#00ff41",
      "neon_name": "Neon Green",
      "glow_behavior": "steady",
      "glow_intensity": 0.6,
      "door_style": "open_archway",
      "lock_icon": "none",
      "border_style": "thin_solid",
      "ambient_sound": "soft_hum"
    },
    {
      "level": 1,
      "name": "INTERNAL",
      "code": "INT",
      "description": "Standard institutional access. Requires authenticated operator identity. Most working documents, standard procedures, non-sensitive references.",
      "neon_color": "#00ffff",
      "neon_name": "Neon Cyan",
      "glow_behavior": "steady",
      "glow_intensity": 0.7,
      "door_style": "standard_door",
      "lock_icon": "keycard",
      "border_style": "thin_solid",
      "ambient_sound": "soft_hum"
    },
    {
      "level": 2,
      "name": "CONFIDENTIAL",
      "code": "CNF",
      "description": "Restricted access. Requires explicit authorization for the specific domain. Security documentation, intelligence analysis, disaster recovery plans, ethics oversight records.",
      "neon_color": "#0080ff",
      "neon_name": "Neon Blue",
      "glow_behavior": "slow_pulse",
      "glow_intensity": 0.8,
      "door_style": "reinforced_door",
      "lock_icon": "keycard_plus_pin",
      "border_style": "double_solid",
      "ambient_sound": "low_frequency_drone"
    },
    {
      "level": 3,
      "name": "SECRET",
      "code": "SEC",
      "description": "Highly restricted. Requires named authorization from the institutional authority. Cryptographic keys, threat assessments, vulnerability reports, succession plans.",
      "neon_color": "#ff00ff",
      "neon_name": "Neon Magenta",
      "glow_behavior": "medium_pulse",
      "glow_intensity": 0.9,
      "door_style": "vault_door",
      "lock_icon": "biometric_scanner",
      "border_style": "triple_solid",
      "ambient_sound": "resonant_hum"
    },
    {
      "level": 4,
      "name": "COSMIC",
      "code": "CSM",
      "description": "Maximum restriction. Penthouse floor only. Requires physical key ceremony. The most sensitive institutional secrets: master keys, succession triggers, existential threat responses, institutional dead man's switches.",
      "neon_color": "#ff0000",
      "neon_name": "Neon Red",
      "glow_behavior": "rapid_pulse",
      "glow_intensity": 1.0,
      "door_style": "blast_door",
      "lock_icon": "multi_key_ceremony",
      "border_style": "animated_border",
      "ambient_sound": "heartbeat_bass"
    }
  ]
}
```

### 2.2 Visual Summary Table

| Level | Name | Neon Color | Door Style | Lock Icon | Glow | Sound |
|---|---|---|---|---|---|---|
| 0 | PUBLIC | Green (#00ff41) | Open archway | None | Steady, 60% | Soft hum |
| 1 | INTERNAL | Cyan (#00ffff) | Standard door | Keycard | Steady, 70% | Soft hum |
| 2 | CONFIDENTIAL | Blue (#0080ff) | Reinforced door | Keycard + PIN | Slow pulse, 80% | Low drone |
| 3 | SECRET | Magenta (#ff00ff) | Vault door | Biometric | Medium pulse, 90% | Resonant hum |
| 4 | COSMIC | Red (#ff0000) | Blast door | Multi-key | Rapid pulse, 100% | Heartbeat bass |

## 3. Access Control Rules

### 3.1 User Clearance Model

Every authenticated user in the HIC has a clearance profile:

```json
{
  "user_id": "holm-operator-1",
  "global_clearance": 3,
  "domain_clearances": {
    "3": 3,
    "7": 2,
    "12": 2
  },
  "room_overrides": {
    "F4-N-R009": 3,
    "B3-E-R001": 2
  },
  "clearance_granted_at": "2026-02-17T00:00:00Z",
  "clearance_expires_at": null,
  "clearance_granted_by": "institutional-authority"
}
```

### 3.2 Access Decision Algorithm

When a user attempts to enter a room, the system evaluates access in the following order:

1. **Room override check.** If the user has a specific room override for this room, use that clearance level.
2. **Domain clearance check.** If the user has a domain-specific clearance for this room's domain, use that clearance level.
3. **Global clearance check.** Use the user's global clearance level.
4. **Compare.** If the user's effective clearance >= the room's clearance level, access is GRANTED. Otherwise, access is DENIED.

```
effective_clearance = max(
  room_overrides.get(room_id, -1),
  domain_clearances.get(room.domain, -1),
  global_clearance
)

access = effective_clearance >= room.clearance_level
```

### 3.3 The Principle of Need-to-Know

Clearance alone is not sufficient. A user with Level 3 clearance can enter Level 3 rooms, but the institution may revoke access to specific rooms through negative room overrides:

```json
{
  "room_overrides": {
    "F4-N-R009": -1
  }
}
```

A clearance of `-1` means explicitly denied, regardless of global or domain clearance. This implements need-to-know: having the clearance does not automatically grant access to every room at that level.

### 3.4 Anonymous and Unauthenticated Access

Unauthenticated users (visitors to the public-facing HIC interface) have an implicit clearance of 0. They can enter PUBLIC rooms only. They cannot see the names or titles of non-PUBLIC documents in search results. They cannot navigate above the ground floor except to floors that contain PUBLIC reading rooms.

## 4. Clearance Inheritance and Overrides

### 4.1 Floor Default Clearance

Every floor has a default clearance level (defined in the floor assignment table, Agent 13, Section 9.1). This is the *minimum* clearance for any room on that floor.

**Rule:** A room's clearance level MUST be >= its floor's default clearance level.

**Example:** Floor F4 (Security & Integrity) has a default clearance of Level 2 (CONFIDENTIAL). No room on F4 can have a clearance level below 2. A room on F4 can be Level 3 (SECRET) or Level 4 (COSMIC) if it contains sufficiently sensitive material, but it cannot be Level 0 or Level 1.

### 4.2 Room-Level Override

Individual rooms can have a higher clearance than the floor default. This is specified in the room manifest's `clearance_level` field. The room-level clearance always takes precedence over the floor default.

**Example:** Floor F3 (Governance & Authority) has a default clearance of Level 1. But room F3-N-R001 (Governance Decision Vault) has a clearance of Level 2 because it contains decision records that are CONFIDENTIAL.

### 4.3 Document-Level Classification vs. Room-Level Clearance

Documents have a `classification` field. Rooms have a `clearance_level` field. These interact as follows:

- **A room's clearance must be >= the highest classification of any document it contains.** If a CONFIDENTIAL document is placed in a room, the room must be at least Level 2.
- **A document can have a lower classification than its room.** An INTERNAL document can exist in a CONFIDENTIAL room. But a CONFIDENTIAL document cannot exist in an INTERNAL room.
- **If a document's classification is raised, the room's clearance must be checked.** If SEC-001 is reclassified from CONFIDENTIAL to SECRET, its room must be upgraded to Level 3 or the document must be moved to a Level 3 room.

## 5. Classification Marking on Documents

Every document in the HIC displays a classification banner. The banner is rendered at the top of the document viewer and is visible on the document's spine when it sits on a shelf.

### 5.1 Banner Specification

```json
{
  "classification_banners": {
    "PUBLIC": {
      "text": "PUBLIC",
      "background_color": "#00ff41",
      "text_color": "#000000",
      "border": "none",
      "font_weight": "normal"
    },
    "INTERNAL": {
      "text": "INTERNAL",
      "background_color": "#00ffff",
      "text_color": "#000000",
      "border": "none",
      "font_weight": "normal"
    },
    "CONFIDENTIAL": {
      "text": "CONFIDENTIAL",
      "background_color": "#0080ff",
      "text_color": "#ffffff",
      "border": "1px solid #0040cc",
      "font_weight": "bold"
    },
    "SECRET": {
      "text": "SECRET",
      "background_color": "#ff00ff",
      "text_color": "#ffffff",
      "border": "2px solid #cc00cc",
      "font_weight": "bold"
    },
    "COSMIC": {
      "text": "COSMIC -- EYES ONLY",
      "background_color": "#ff0000",
      "text_color": "#ffffff",
      "border": "3px solid #cc0000",
      "font_weight": "bold",
      "animation": "pulse_glow"
    }
  }
}
```

### 5.2 Shelf Spine Rendering

When a document sits on a shelf, its spine (the edge facing outward) displays:
- The classification color as a vertical stripe.
- The document ID in small text.
- The document title (truncated to 40 characters).
- The format icon (a small glyph indicating article, procedure, decision_record, etc.).

For classified documents above the user's clearance, the spine displays only the classification color stripe and the text `[CLASSIFIED]`. No title, no ID.

## 6. Declassification Procedures

A document's classification can be lowered through a formal declassification process. In spatial terms, declassification means moving a document from a higher-clearance room to a lower-clearance room.

### 6.1 Declassification Steps

1. **Initiation.** The institutional authority (or authorized delegate per GOV-001) issues a declassification order specifying the document ID, current classification, target classification, and justification.
2. **Review.** The document is reviewed against D15-003 (Ethical Review Process) to confirm that declassification does not create harm.
3. **Reclassification.** The document's `classification` field is updated.
4. **Room Reassignment.** If the document's new classification is below the current room's clearance level, the document must be moved to an appropriate room. A new binding is created in the target room's manifest, and the old binding is replaced with a forwarding reference.
5. **Audit Entry.** A declassification record is written to the audit trail:

```json
{
  "event": "declassification",
  "document_id": "SEC-004",
  "previous_classification": "SECRET",
  "new_classification": "CONFIDENTIAL",
  "previous_room": "F4-N-R009",
  "new_room": "F4-N-R003",
  "authorized_by": "institutional-authority",
  "timestamp": "2026-06-01T09:00:00Z",
  "justification": "Document content no longer sensitive after architectural change documented in SEC-004 v3.0."
}
```

### 6.2 Classification Escalation

The reverse process -- moving a document to a *higher* classification -- follows the same steps but does not require ethical review. Classification escalation is always permitted. The only constraint is that the document must be moved to a room with sufficient clearance.

## 7. Audit Trail Specification

Every access event in the HIC is logged. The audit trail is stored in a dedicated vault room on each floor.

### 7.1 Audit Event Schema

```json
{
  "event_id": "evt-2026-02-17-00001",
  "event_type": "room_entry",
  "user_id": "holm-operator-1",
  "room_id": "F4-N-R003",
  "clearance_used": 2,
  "access_result": "granted",
  "timestamp": "2026-02-17T14:30:00Z",
  "session_id": "sess-abc123",
  "ip_address": null,
  "notes": null
}
```

### 7.2 Audit Event Types

| Event Type | Description |
|---|---|
| `room_entry` | User entered a room. |
| `room_entry_denied` | User attempted to enter a room but was denied. |
| `document_read` | User opened a document for reading. |
| `document_search` | User performed a search query. |
| `document_placed` | A document was bound to a room. |
| `document_moved` | A document was moved between rooms. |
| `document_archived` | A document was moved to archive. |
| `document_sealed` | A document was permanently sealed. |
| `classification_changed` | A document's classification was changed. |
| `clearance_granted` | A user was granted a clearance level. |
| `clearance_revoked` | A user's clearance was revoked. |
| `manifest_modified` | A room manifest was modified. |
| `integrity_failure` | A manifest or index integrity check failed. |

### 7.3 Audit Trail Storage

Audit logs are stored in a dedicated vault room on the ground floor: `F1-N-R001` ("Audit Vault"). This room has a clearance of Level 3 (SECRET). Audit logs are append-only. They cannot be modified or deleted. The audit vault has no capacity limit -- it uses the archive room type (5000 documents), and when full, a new audit vault room is created.

## 8. Visual Lockout Rendering

When a user approaches a room they lack clearance to enter, the room is rendered in a locked-out state. The purpose of the lockout rendering is to communicate four things: (1) the room exists, (2) it contains something, (3) you cannot enter, and (4) what clearance level you need.

### 8.1 Lockout Visual Specification

| Element | Rendering |
|---|---|
| **Door** | The door style matches the clearance level (see Section 2.1) but is rendered in a desaturated, dim version of the clearance color. |
| **Door surface** | Displays the lock icon for the room's clearance level, glowing at 40% intensity. |
| **Room interior** | Not visible. The door has no window. The room is fully opaque from outside. |
| **Door label** | Displays the room ID and room name (if any). Does NOT display document count or document titles. |
| **Clearance indicator** | A small holographic badge above the door displays `LEVEL {N} -- {NAME}` in the clearance color. |
| **Interaction** | Clicking the door triggers a denial animation: the lock icon flashes red three times, and a text overlay reads `ACCESS DENIED -- LEVEL {N} CLEARANCE REQUIRED`. |
| **Ambient effect** | The area immediately around the locked door has a faint glow in the clearance color, creating a visible aura that warns approaching users before they reach the door. |

### 8.2 Progressive Lockout

For floors above the user's clearance, the entire floor is locked. The elevator will not stop at that floor. The floor number on the elevator panel is displayed in the clearance color with a lock icon next to it. The user can see the floor exists in the building directory, but they cannot reach it.

**Exception:** If a floor contains at least one room the user can access (due to a room override granting access to a specific room on an otherwise restricted floor), the elevator will stop at that floor, but all locked rooms on that floor are rendered in the lockout state.

---

---

# Agent 15: Archive Placement & Memory Vaults

**Document ID:** HIC-015
**Agent:** 15 -- Knowledge Mapping: Archives & Memory
**Status:** Phase 1
**Date:** 2026-02-17
**Depends On:** HIC-013 (Room-to-Document Binding), HIC-014 (Clearance Levels), D6-001 (Data Philosophy), D6-004 (Archive Management Procedures), D20-001 (Institutional Memory Philosophy)
**Depended Upon By:** HIC-016 (Art Direction), HIC-017 (Lighting & Atmosphere), all future HIC components that reference archived or sealed documents

---

## 1. Purpose

This article defines the three basement floors of the Holm Intelligence Complex and the specialized vault types that exist within them. The basement is the institution's long-term memory. Where the above-ground floors contain active, working documentation, the basement contains the documentation that has been retired from active use but must be preserved.

The basement is not a dumping ground. It is a carefully organized system of archives, vaults, and sealed chambers, each with a specific purpose, a specific organizational scheme, and a specific set of rules governing what enters, what leaves, and what is permanently locked.

The spatial metaphor is crucial in the basement. Moving deeper underground corresponds to moving further back in time and further from active use. B1 is recent history -- documents retired in the last few years. B2 is deep history -- documents from earlier eras of the institution's life. B3 is the deep cold -- permanently sealed records that exist as the institution's bedrock memory.

## 2. Basement Floor Organization

### 2.1 Floor B1: Recent Archive

**Clearance:** Level 1 (INTERNAL) minimum
**Ambient:** Dim cyan (#004444). Cool, quiet, slightly humid. The sound of distant mechanical ventilation.
**Purpose:** Documents retired from active floors within the last 5 years. These are still referenced occasionally and may be recalled to active floors if needed.

**Room layout:**

| Wing | Room Range | Count | Room Type | Purpose |
|---|---|---|---|---|
| N | R001 - R020 | 20 | Archive | Domain-specific recent archives (one per domain) |
| E | R001 - R005 | 5 | Archive | Cross-domain recent archives |
| S | R001 - R003 | 3 | Office | Archive intake and processing |
| W | R001 - R010 | 10 | Archive | Overflow and bulk storage |
| **Total** | | **38** | | |

**Room assignment:** Each domain has a dedicated recent archive room in the North wing. Domain 1 is B1-N-R001, Domain 2 is B1-N-R002, and so forth through Domain 20 at B1-N-R020. Cross-domain documents (domain 0) go to the East wing.

### 2.2 Floor B2: Historical Archive

**Clearance:** Level 1 (INTERNAL) minimum, with individual rooms at Level 2 as needed
**Ambient:** Dim blue (#000044). Cold, dark, very quiet. The sound of deep silence with occasional metallic groans from the building's structure.
**Purpose:** Documents from the institution's earlier eras. Documents that have been in B1 for more than 5 years are migrated here. This floor preserves the institution's historical record.

**Room layout:**

| Wing | Room Range | Count | Room Type | Purpose |
|---|---|---|---|---|
| N | R001 - R020 | 20 | Archive | Domain-specific historical archives (one per domain) |
| E | R001 - R010 | 10 | Archive | Cross-domain historical archives |
| S | R001 - R005 | 5 | Vault | Sealed historical records (read-only) |
| W | R001 - R010 | 10 | Archive | Decade-organized archives (one room per decade) |
| **Total** | | **45** | | |

**Decade rooms (West wing):** B2-W-R001 contains documents from the institution's founding decade. B2-W-R002 contains documents from the second decade. And so forth. These rooms provide a chronological browsing experience -- walk from B2-W-R001 to B2-W-R010 and you walk through the institution's history.

### 2.3 Floor B3: Deep Cold Archive

**Clearance:** Level 2 (CONFIDENTIAL) minimum, with vault rooms at Level 3 or Level 4
**Ambient:** Dim violet (#110022). Near-total darkness. Cold. The sound of a low, barely perceptible bass tone -- the heartbeat of the building.
**Purpose:** Permanently sealed records. The institution's bedrock. Documents in B3 are never modified, never moved, and never deleted. They are the institutional equivalent of geological strata -- compressed layers of accumulated knowledge that form the foundation everything else rests on.

**Room layout:**

| Wing | Room Range | Count | Room Type | Purpose |
|---|---|---|---|---|
| N | R001 - R010 | 10 | Vault | Memory Vaults (see Section 5) |
| E | R001 - R005 | 5 | Vault | Decision Log Vault (see Section 6) |
| S | R001 - R005 | 5 | Vault | Oral History Vault (see Section 7) |
| W | R001 - R005 | 5 | Vault | Lessons Learned Vault (see Section 8) |
| Central | R001 - R003 | 3 | Vault | Dead Drop Rooms (see Section 9) |
| **Total** | | **28** | | |

## 3. Archive Indexing System

### 3.1 Archive Index Schema

Each basement floor maintains its own index, parallel to the global search index but optimized for archive-specific queries (date range, domain, era, archival reason).

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "ArchiveIndex",
  "type": "object",
  "required": ["floor", "version", "built_at", "entries"],
  "properties": {
    "floor": {
      "type": "string",
      "enum": ["B1", "B2", "B3"]
    },
    "version": { "type": "string" },
    "built_at": { "type": "string", "format": "date-time" },
    "entry_count": { "type": "integer" },
    "entries": {
      "type": "array",
      "items": {
        "type": "object",
        "required": [
          "document_id",
          "room_id",
          "title",
          "domain",
          "classification",
          "archived_at",
          "archived_from",
          "archive_reason"
        ],
        "properties": {
          "document_id": { "type": "string" },
          "room_id": { "type": "string" },
          "title": { "type": "string" },
          "domain": { "type": "integer" },
          "classification": { "type": "string" },
          "archived_at": {
            "type": "string",
            "format": "date-time",
            "description": "When the document was moved to the archive."
          },
          "archived_from": {
            "type": "string",
            "description": "The room_id the document was moved from."
          },
          "archive_reason": {
            "type": "string",
            "enum": [
              "superseded",
              "retired",
              "periodic_archival",
              "classification_change",
              "domain_reorganization",
              "permanent_seal",
              "operator_decision"
            ]
          },
          "era": {
            "type": "string",
            "description": "The institutional era or decade this document belongs to."
          },
          "recall_permitted": {
            "type": "boolean",
            "description": "Whether this document can be recalled to an active floor."
          },
          "keywords": {
            "type": "array",
            "items": { "type": "string" }
          }
        }
      }
    }
  }
}
```

### 3.2 Archive Index File Locations

```
hic-data/floors/B1/archive-index.json
hic-data/floors/B1/archive-index.json.sha256
hic-data/floors/B2/archive-index.json
hic-data/floors/B2/archive-index.json.sha256
hic-data/floors/B3/archive-index.json
hic-data/floors/B3/archive-index.json.sha256
```

### 3.3 Archive Search Behavior

When a user searches from within the basement, results are sorted by archival date (most recent first) rather than by relevance. The search interface on basement floors has additional filter controls:

- **Date range filter:** Restrict results to documents archived within a specific date range.
- **Domain filter:** Restrict results to a specific domain.
- **Era filter:** Restrict results to a specific institutional era (decade).
- **Seal status filter:** Show only sealed documents, only recallable documents, or all.
- **Archive reason filter:** Filter by the reason the document was archived.

## 4. Complete Basement Room Registry

### 4.1 Floor B1 -- Recent Archive: Full Room Registry

```json
{
  "floor": "B1",
  "floor_name": "Recent Archive",
  "clearance_default": 1,
  "ambient_color": "#004444",
  "rooms": [
    { "room_id": "B1-N-R001", "room_type": "archive", "capacity": 5000, "domain": 1, "room_name": "Constitution Recent Archive" },
    { "room_id": "B1-N-R002", "room_type": "archive", "capacity": 5000, "domain": 2, "room_name": "Governance Recent Archive" },
    { "room_id": "B1-N-R003", "room_type": "archive", "capacity": 5000, "domain": 3, "room_name": "Security Recent Archive" },
    { "room_id": "B1-N-R004", "room_type": "archive", "capacity": 5000, "domain": 4, "room_name": "Infrastructure Recent Archive" },
    { "room_id": "B1-N-R005", "room_type": "archive", "capacity": 5000, "domain": 5, "room_name": "Platform Recent Archive" },
    { "room_id": "B1-N-R006", "room_type": "archive", "capacity": 5000, "domain": 6, "room_name": "Data & Archives Recent Archive" },
    { "room_id": "B1-N-R007", "room_type": "archive", "capacity": 5000, "domain": 7, "room_name": "Intelligence Recent Archive" },
    { "room_id": "B1-N-R008", "room_type": "archive", "capacity": 5000, "domain": 8, "room_name": "Automation Recent Archive" },
    { "room_id": "B1-N-R009", "room_type": "archive", "capacity": 5000, "domain": 9, "room_name": "Education Recent Archive" },
    { "room_id": "B1-N-R010", "room_type": "archive", "capacity": 5000, "domain": 10, "room_name": "Operations Recent Archive" },
    { "room_id": "B1-N-R011", "room_type": "archive", "capacity": 5000, "domain": 11, "room_name": "Administration Recent Archive" },
    { "room_id": "B1-N-R012", "room_type": "archive", "capacity": 5000, "domain": 12, "room_name": "Disaster Recovery Recent Archive" },
    { "room_id": "B1-N-R013", "room_type": "archive", "capacity": 5000, "domain": 13, "room_name": "Evolution Recent Archive" },
    { "room_id": "B1-N-R014", "room_type": "archive", "capacity": 5000, "domain": 14, "room_name": "Research Recent Archive" },
    { "room_id": "B1-N-R015", "room_type": "archive", "capacity": 5000, "domain": 15, "room_name": "Ethics Recent Archive" },
    { "room_id": "B1-N-R016", "room_type": "archive", "capacity": 5000, "domain": 16, "room_name": "Interface Recent Archive" },
    { "room_id": "B1-N-R017", "room_type": "archive", "capacity": 5000, "domain": 17, "room_name": "Federation Recent Archive" },
    { "room_id": "B1-N-R018", "room_type": "archive", "capacity": 5000, "domain": 18, "room_name": "Import Recent Archive" },
    { "room_id": "B1-N-R019", "room_type": "archive", "capacity": 5000, "domain": 19, "room_name": "Quality Recent Archive" },
    { "room_id": "B1-N-R020", "room_type": "archive", "capacity": 5000, "domain": 20, "room_name": "Memory Recent Archive" },
    { "room_id": "B1-E-R001", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Recent Archive A" },
    { "room_id": "B1-E-R002", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Recent Archive B" },
    { "room_id": "B1-E-R003", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Recent Archive C" },
    { "room_id": "B1-E-R004", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Recent Archive D" },
    { "room_id": "B1-E-R005", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Recent Archive E" },
    { "room_id": "B1-S-R001", "room_type": "office", "capacity": 50, "domain": 0, "room_name": "Archive Intake Office" },
    { "room_id": "B1-S-R002", "room_type": "office", "capacity": 50, "domain": 0, "room_name": "Archive Processing Office" },
    { "room_id": "B1-S-R003", "room_type": "office", "capacity": 50, "domain": 0, "room_name": "Archive Recall Staging" },
    { "room_id": "B1-W-R001", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Overflow Archive A" },
    { "room_id": "B1-W-R002", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Overflow Archive B" },
    { "room_id": "B1-W-R003", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Overflow Archive C" },
    { "room_id": "B1-W-R004", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Overflow Archive D" },
    { "room_id": "B1-W-R005", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Overflow Archive E" },
    { "room_id": "B1-W-R006", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Overflow Archive F" },
    { "room_id": "B1-W-R007", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Overflow Archive G" },
    { "room_id": "B1-W-R008", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Overflow Archive H" },
    { "room_id": "B1-W-R009", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Overflow Archive I" },
    { "room_id": "B1-W-R010", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Overflow Archive J" }
  ]
}
```

### 4.2 Floor B2 -- Historical Archive: Full Room Registry

```json
{
  "floor": "B2",
  "floor_name": "Historical Archive",
  "clearance_default": 1,
  "ambient_color": "#000044",
  "rooms": [
    { "room_id": "B2-N-R001", "room_type": "archive", "capacity": 5000, "domain": 1, "room_name": "Constitution Historical Archive" },
    { "room_id": "B2-N-R002", "room_type": "archive", "capacity": 5000, "domain": 2, "room_name": "Governance Historical Archive" },
    { "room_id": "B2-N-R003", "room_type": "archive", "capacity": 5000, "domain": 3, "room_name": "Security Historical Archive" },
    { "room_id": "B2-N-R004", "room_type": "archive", "capacity": 5000, "domain": 4, "room_name": "Infrastructure Historical Archive" },
    { "room_id": "B2-N-R005", "room_type": "archive", "capacity": 5000, "domain": 5, "room_name": "Platform Historical Archive" },
    { "room_id": "B2-N-R006", "room_type": "archive", "capacity": 5000, "domain": 6, "room_name": "Data Historical Archive" },
    { "room_id": "B2-N-R007", "room_type": "archive", "capacity": 5000, "domain": 7, "room_name": "Intelligence Historical Archive" },
    { "room_id": "B2-N-R008", "room_type": "archive", "capacity": 5000, "domain": 8, "room_name": "Automation Historical Archive" },
    { "room_id": "B2-N-R009", "room_type": "archive", "capacity": 5000, "domain": 9, "room_name": "Education Historical Archive" },
    { "room_id": "B2-N-R010", "room_type": "archive", "capacity": 5000, "domain": 10, "room_name": "Operations Historical Archive" },
    { "room_id": "B2-N-R011", "room_type": "archive", "capacity": 5000, "domain": 11, "room_name": "Administration Historical Archive" },
    { "room_id": "B2-N-R012", "room_type": "archive", "capacity": 5000, "domain": 12, "room_name": "Disaster Recovery Historical Archive" },
    { "room_id": "B2-N-R013", "room_type": "archive", "capacity": 5000, "domain": 13, "room_name": "Evolution Historical Archive" },
    { "room_id": "B2-N-R014", "room_type": "archive", "capacity": 5000, "domain": 14, "room_name": "Research Historical Archive" },
    { "room_id": "B2-N-R015", "room_type": "archive", "capacity": 5000, "domain": 15, "room_name": "Ethics Historical Archive" },
    { "room_id": "B2-N-R016", "room_type": "archive", "capacity": 5000, "domain": 16, "room_name": "Interface Historical Archive" },
    { "room_id": "B2-N-R017", "room_type": "archive", "capacity": 5000, "domain": 17, "room_name": "Federation Historical Archive" },
    { "room_id": "B2-N-R018", "room_type": "archive", "capacity": 5000, "domain": 18, "room_name": "Import Historical Archive" },
    { "room_id": "B2-N-R019", "room_type": "archive", "capacity": 5000, "domain": 19, "room_name": "Quality Historical Archive" },
    { "room_id": "B2-N-R020", "room_type": "archive", "capacity": 5000, "domain": 20, "room_name": "Memory Historical Archive" },
    { "room_id": "B2-E-R001", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Historical A" },
    { "room_id": "B2-E-R002", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Historical B" },
    { "room_id": "B2-E-R003", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Historical C" },
    { "room_id": "B2-E-R004", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Historical D" },
    { "room_id": "B2-E-R005", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Historical E" },
    { "room_id": "B2-E-R006", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Historical F" },
    { "room_id": "B2-E-R007", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Historical G" },
    { "room_id": "B2-E-R008", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Historical H" },
    { "room_id": "B2-E-R009", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Historical I" },
    { "room_id": "B2-E-R010", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Cross-Domain Historical J" },
    { "room_id": "B2-S-R001", "room_type": "vault", "capacity": 500, "domain": 0, "room_name": "Sealed Historical Records A", "clearance_level": 2 },
    { "room_id": "B2-S-R002", "room_type": "vault", "capacity": 500, "domain": 0, "room_name": "Sealed Historical Records B", "clearance_level": 2 },
    { "room_id": "B2-S-R003", "room_type": "vault", "capacity": 500, "domain": 0, "room_name": "Sealed Historical Records C", "clearance_level": 2 },
    { "room_id": "B2-S-R004", "room_type": "vault", "capacity": 500, "domain": 0, "room_name": "Sealed Historical Records D", "clearance_level": 2 },
    { "room_id": "B2-S-R005", "room_type": "vault", "capacity": 500, "domain": 0, "room_name": "Sealed Historical Records E", "clearance_level": 2 },
    { "room_id": "B2-W-R001", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Decade 1 Archive (Founding)" },
    { "room_id": "B2-W-R002", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Decade 2 Archive" },
    { "room_id": "B2-W-R003", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Decade 3 Archive" },
    { "room_id": "B2-W-R004", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Decade 4 Archive" },
    { "room_id": "B2-W-R005", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Decade 5 Archive" },
    { "room_id": "B2-W-R006", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Decade 6 Archive" },
    { "room_id": "B2-W-R007", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Decade 7 Archive" },
    { "room_id": "B2-W-R008", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Decade 8 Archive" },
    { "room_id": "B2-W-R009", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Decade 9 Archive" },
    { "room_id": "B2-W-R010", "room_type": "archive", "capacity": 5000, "domain": 0, "room_name": "Decade 10 Archive" }
  ]
}
```

### 4.3 Floor B3 -- Deep Cold Archive: Full Room Registry

```json
{
  "floor": "B3",
  "floor_name": "Deep Cold Archive",
  "clearance_default": 2,
  "ambient_color": "#110022",
  "rooms": [
    { "room_id": "B3-N-R001", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 3, "room_name": "Memory Vault Alpha" },
    { "room_id": "B3-N-R002", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 3, "room_name": "Memory Vault Beta" },
    { "room_id": "B3-N-R003", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 3, "room_name": "Memory Vault Gamma" },
    { "room_id": "B3-N-R004", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 3, "room_name": "Memory Vault Delta" },
    { "room_id": "B3-N-R005", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 3, "room_name": "Memory Vault Epsilon" },
    { "room_id": "B3-N-R006", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 3, "room_name": "Memory Vault Zeta" },
    { "room_id": "B3-N-R007", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 3, "room_name": "Memory Vault Eta" },
    { "room_id": "B3-N-R008", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 3, "room_name": "Memory Vault Theta" },
    { "room_id": "B3-N-R009", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 4, "room_name": "Memory Vault Iota (COSMIC)" },
    { "room_id": "B3-N-R010", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 4, "room_name": "Memory Vault Kappa (COSMIC)" },
    { "room_id": "B3-E-R001", "room_type": "vault", "capacity": 500, "domain": 2, "clearance_level": 3, "room_name": "Decision Log Vault: Governance" },
    { "room_id": "B3-E-R002", "room_type": "vault", "capacity": 500, "domain": 3, "clearance_level": 3, "room_name": "Decision Log Vault: Security" },
    { "room_id": "B3-E-R003", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 3, "room_name": "Decision Log Vault: Cross-Domain" },
    { "room_id": "B3-E-R004", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 3, "room_name": "Decision Log Vault: Emergency" },
    { "room_id": "B3-E-R005", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 4, "room_name": "Decision Log Vault: COSMIC" },
    { "room_id": "B3-S-R001", "room_type": "vault", "capacity": 500, "domain": 20, "clearance_level": 2, "room_name": "Oral History Vault: Founder" },
    { "room_id": "B3-S-R002", "room_type": "vault", "capacity": 500, "domain": 20, "clearance_level": 2, "room_name": "Oral History Vault: Operators" },
    { "room_id": "B3-S-R003", "room_type": "vault", "capacity": 500, "domain": 20, "clearance_level": 2, "room_name": "Oral History Vault: External Voices" },
    { "room_id": "B3-S-R004", "room_type": "vault", "capacity": 500, "domain": 20, "clearance_level": 3, "room_name": "Oral History Vault: Classified Testimonies" },
    { "room_id": "B3-S-R005", "room_type": "vault", "capacity": 500, "domain": 20, "clearance_level": 2, "room_name": "Oral History Vault: Institutional Reflections" },
    { "room_id": "B3-W-R001", "room_type": "vault", "capacity": 500, "domain": 12, "clearance_level": 2, "room_name": "Lessons Learned: Incidents" },
    { "room_id": "B3-W-R002", "room_type": "vault", "capacity": 500, "domain": 13, "clearance_level": 2, "room_name": "Lessons Learned: Migrations" },
    { "room_id": "B3-W-R003", "room_type": "vault", "capacity": 500, "domain": 3, "clearance_level": 3, "room_name": "Lessons Learned: Security Events" },
    { "room_id": "B3-W-R004", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 2, "room_name": "Lessons Learned: General" },
    { "room_id": "B3-W-R005", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 2, "room_name": "Lessons Learned: Near Misses" },
    { "room_id": "B3-C-R001", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 2, "room_name": "Dead Drop Room Alpha" },
    { "room_id": "B3-C-R002", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 2, "room_name": "Dead Drop Room Beta" },
    { "room_id": "B3-C-R003", "room_type": "vault", "capacity": 500, "domain": 0, "clearance_level": 2, "room_name": "Dead Drop Room Gamma" }
  ]
}
```

## 5. Memory Vaults

Memory Vaults are time-locked rooms on B3 that preserve institutional memory. They are the institution's equivalent of a time capsule -- but one that can be read at any time, while never being altered.

### 5.1 Memory Vault Specification

```json
{
  "vault_type": "memory_vault",
  "properties": {
    "seal_type": "permanent",
    "modification_permitted": false,
    "deletion_permitted": false,
    "relocation_permitted": false,
    "access_mode": "read_only",
    "time_lock": {
      "sealed_at": "ISO 8601 timestamp",
      "sealed_by": "operator identity",
      "seal_reason": "free text",
      "seal_ceremony_witnesses": ["list of witness identities"]
    },
    "integrity_verification": {
      "method": "SHA-256 hash of all document contents concatenated",
      "verification_schedule": "monthly",
      "last_verified": "ISO 8601 timestamp",
      "verification_status": "pass | fail | pending"
    }
  }
}
```

### 5.2 Memory Vault Contents

Memory Vaults contain documents that record the institution's identity, values, and turning points. Typical contents include:

- **Founding documents.** The original versions of CON-001, ETH-001, GOV-001, and OPS-001 as they existed at the moment of institutional founding. These are the bedrock. They are never updated in the vault -- if the active versions are amended, the vault preserves the originals.
- **Era-defining decisions.** The major decisions that shaped the institution's direction. Each decision is stored as a complete decision record with context, alternatives considered, reasoning, and outcome.
- **Institutional identity snapshots.** Periodic snapshots of what the institution is, what it values, and how it sees itself. These are written by the operator and sealed at regular intervals (annually recommended).
- **Letters to the future.** Per META-010, letters from current operators to future maintainers. These are personal, reflective, and sealed upon submission.

### 5.3 Visual Rendering of Memory Vaults

Memory Vaults are rendered differently from standard vaults. Their visual treatment communicates permanence and reverence:

- **Door:** Massive blast-style door with a visible time-lock mechanism -- a large circular display showing the seal date and the word `SEALED` in slowly pulsing neon amber.
- **Interior:** Dimly lit with warm amber lighting instead of the cold blues and violets of the surrounding archive. Shelves are rendered in dark wood rather than metal.
- **Documents:** Displayed as bound volumes with gilt edges rather than standard shelf spines. Each volume has a seal icon on its cover.
- **Sound:** A low, warm harmonic tone -- distinctly different from the cold silence of the rest of B3.
- **Walls:** Faint inscriptions on the walls display the names of the operators who sealed each document.

## 6. Decision Log Vault

The Decision Log Vault (B3-E-R001 through B3-E-R005) is a specialized vault system that permanently preserves every major institutional decision.

### 6.1 Decision Log Entry Schema

```json
{
  "decision_id": "DEC-2026-001",
  "title": "Adoption of Five-Tier Clearance System",
  "date": "2026-02-17",
  "domain": 3,
  "classification": "SECRET",
  "decided_by": "institutional-authority",
  "context": "The institution required a formalized access control system for the HIC spatial interface. Multiple models were evaluated.",
  "alternatives_considered": [
    {
      "name": "Three-tier model (public, internal, restricted)",
      "reason_rejected": "Insufficient granularity for COSMIC-level materials."
    },
    {
      "name": "Role-based access without tiers",
      "reason_rejected": "Too complex for single-operator institution. Tiers are simpler."
    }
  ],
  "decision": "Adopt a five-tier clearance system: PUBLIC, INTERNAL, CONFIDENTIAL, SECRET, COSMIC.",
  "reasoning": "Five tiers provide sufficient granularity without excessive complexity. The visual language of neon colors makes clearance immediately apparent in the spatial interface.",
  "consequences_expected": "All rooms must be assigned a clearance level. All documents must be classified. The interface must render lockout states for insufficient clearance.",
  "review_date": "2027-02-17",
  "sealed": true,
  "sealed_at": "2026-02-17T12:00:00Z"
}
```

### 6.2 Decision Log Rules

- Every decision that creates, modifies, or destroys an institutional article, domain, clearance level, or governance process MUST be logged.
- Decision logs are append-only. A decision cannot be modified after sealing.
- If a decision is reversed, a new decision log entry is created referencing the original. The original is not altered.
- Decision logs are stored in the East wing of B3. Governance decisions in B3-E-R001, Security decisions in B3-E-R002, Cross-domain decisions in B3-E-R003, Emergency decisions in B3-E-R004, and COSMIC-level decisions in B3-E-R005.

## 7. Oral History Vault

The Oral History Vault (B3-S-R001 through B3-S-R005) preserves the institution's spoken memory -- transcripts, recordings, and reflective narratives that capture knowledge which does not fit neatly into formal articles.

### 7.1 Oral History Entry Schema

```json
{
  "recording_id": "OH-2026-001",
  "title": "Founder's Reflection: Why Air-Gap",
  "speaker": "holm-founder",
  "date_recorded": "2026-02-17",
  "duration_minutes": 45,
  "format": "transcript",
  "classification": "CONFIDENTIAL",
  "room_id": "B3-S-R001",
  "transcript_document_id": "OH-T-2026-001",
  "audio_file_reference": "oh-audio/OH-2026-001.opus",
  "topics": ["air-gap", "founding motivation", "threat model", "institutional philosophy"],
  "sealed": true,
  "sealed_at": "2026-02-18T00:00:00Z",
  "abstract": "The founder explains the personal and philosophical reasons for choosing an air-gapped architecture, including experiences with data loss, surveillance, and institutional dependency that preceded the founding."
}
```

### 7.2 Oral History Vault Organization

| Vault Room | Purpose | Contents |
|---|---|---|
| B3-S-R001 | Founder | The founder's own reflections, recorded at key moments. |
| B3-S-R002 | Operators | Reflections from operators who maintain the institution over time. |
| B3-S-R003 | External Voices | Transcripts of conversations with people outside the institution who provided insight, challenge, or perspective. |
| B3-S-R004 | Classified Testimonies | Recordings that contain sensitive information (threat details, personal disclosures, security incidents). Level 3 clearance. |
| B3-S-R005 | Institutional Reflections | Periodic recorded reflections on the state of the institution, its health, its direction, and its struggles. |

### 7.3 Oral History Rules

- Every oral history entry MUST have a written transcript. Audio alone is insufficient -- audio formats may become unreadable; text survives.
- Oral histories are sealed upon entry. They cannot be edited after sealing.
- The speaker's identity is recorded but may be anonymized in the transcript if the speaker requests it. The vault record always contains the true identity; the transcript may use a pseudonym.
- Oral histories are never declassified automatically. Declassification requires a formal governance decision.

## 8. Lessons Learned Vault

The Lessons Learned Vault (B3-W-R001 through B3-W-R005) preserves post-incident and post-event analyses. Every time something goes wrong -- or almost goes wrong -- the institution records what happened, why, and what was learned.

### 8.1 Lessons Learned Entry Schema

```json
{
  "lesson_id": "LL-2026-001",
  "title": "Backup Verification Failure Due to Expired Hash Algorithm",
  "date": "2026-03-15",
  "incident_date": "2026-03-10",
  "domain": 6,
  "classification": "INTERNAL",
  "severity": "major",
  "category": "data_integrity",
  "what_happened": "Monthly backup verification failed because the hash algorithm used to verify backup integrity had been deprecated in the latest OS update. The verification script reported all backups as corrupt when they were actually intact.",
  "root_cause": "The verification script hardcoded a dependency on a specific hash implementation that was removed during the OS update. No fallback was specified.",
  "what_we_learned": "All verification scripts must use the institution's cryptographic abstraction layer (SEC-005, Section 4) rather than calling hash functions directly. Direct dependencies on specific implementations are fragile.",
  "actions_taken": [
    "Updated verification scripts to use the cryptographic abstraction layer.",
    "Added a pre-update check that validates all critical scripts against the new OS version.",
    "Created D6-003 Amendment 1 documenting the required abstraction."
  ],
  "related_documents": ["D6-003", "SEC-005", "D5-002"],
  "sealed": true,
  "sealed_at": "2026-03-20T00:00:00Z"
}
```

### 8.2 Lessons Learned Vault Organization

| Vault Room | Purpose | Severity Range |
|---|---|---|
| B3-W-R001 | Incidents | Major and critical incidents that caused data loss, service disruption, or security events. |
| B3-W-R002 | Migrations | Lessons from technology migrations, hardware replacements, and format conversions. |
| B3-W-R003 | Security Events | Security-specific lessons. Level 3 clearance due to sensitive threat information. |
| B3-W-R004 | General | Lessons from day-to-day operations that do not fit other categories. |
| B3-W-R005 | Near Misses | Events that almost caused problems but were caught in time. These are often the most valuable lessons. |

## 9. Dead Drop Rooms

Dead Drop Rooms (B3-C-R001 through B3-C-R003) are anonymous submission points. They allow any authenticated user to deposit a document without attribution. The document's `bound_by` field is set to `anonymous`, and no audit trail connects the document to the submitting user.

### 9.1 Dead Drop Specification

```json
{
  "room_type": "dead_drop",
  "properties": {
    "attribution": "anonymous",
    "audit_trail": "submission_only",
    "submitter_identity_recorded": false,
    "review_required": true,
    "review_authority": "institutional-authority",
    "clearance_to_submit": 1,
    "clearance_to_read": 2,
    "visual_treatment": {
      "door_style": "unmarked_slot",
      "interior_visible": false,
      "room_name_visible": false,
      "ambient_color": "#330033"
    }
  }
}
```

### 9.2 Dead Drop Rules

- Any authenticated user (Level 1+) can submit to a dead drop.
- The submitter's identity is NOT recorded anywhere. The HIC explicitly does not log who submitted a dead drop document. This is the one exception to the audit trail rule.
- Dead drop documents enter a review queue. The institutional authority (or designated reviewer per GOV-001) reads each submission and decides whether to: (a) place it in an appropriate room, (b) escalate it, (c) archive it, or (d) discard it.
- Dead drops exist to enable uncomfortable truths. An operator who sees something wrong -- an ethical violation, a security concern, a governance failure -- can report it without fear of identification. This is the institution's safety valve.
- Dead drop rooms are visually distinct: no room name on the door, no window, just a narrow slot in the wall with a faint purple glow. The interface renders a simple "Submit Document" action when the user approaches.

### 9.3 Dead Drop Visual Treatment

From the outside, dead drop rooms appear as featureless walls with a narrow illuminated slot -- a mail slot in a wall of dark metal. There is no door. There is no handle. The user drags a document to the slot, and it disappears inside. The wall gives no indication of what is inside or how full the room is.

From the reviewer's perspective (Level 2+ clearance with reviewer authorization), the dead drop room interior is accessible through a separate entrance on the opposite wall. The reviewer sees a queue of submissions, each with a timestamp and no other identifying information.

## 10. Archive Retrieval Procedures

Documents in the archive can be recalled to active floors if they are not sealed.

### 10.1 Recall Request Schema

```json
{
  "request_id": "RECALL-2026-001",
  "document_id": "D6-005",
  "current_room": "B1-N-R006",
  "target_room": "F7-W-R003",
  "requested_by": "holm-operator-1",
  "requested_at": "2026-08-01T10:00:00Z",
  "reason": "Document needed for active reference during format migration project.",
  "approval_required": true,
  "approved_by": null,
  "approved_at": null,
  "status": "pending",
  "recall_type": "temporary",
  "return_date": "2026-09-01T00:00:00Z"
}
```

### 10.2 Recall Types

| Type | Description | Return Required |
|---|---|---|
| **Temporary** | Document is recalled for a specific project or task. Must be returned by the specified return date. | Yes, by return_date. |
| **Permanent** | Document is recalled to active service. It will not return to the archive unless re-archived later. | No. |
| **Reference Copy** | A read-only reference is created in the target room. The original remains in the archive. This is the ONLY case where a document reference exists in two rooms -- the archive retains the original, and the active floor gets a reference marker. | N/A -- original never moves. |

### 10.3 Recall Visual Treatment

When a recalled document is in transit from the archive to an active floor, the elevator shaft displays a glowing document icon moving upward through the building. The target room displays an "Incoming" indicator on the shelf slot where the document will be placed. The archive room displays a "Recalled" placeholder on the shelf where the document was, showing where it went and when it is expected to return (for temporary recalls).

## 11. Vault Integrity Verification

All vault rooms (standard vaults, memory vaults, decision log vaults, oral history vaults, lessons learned vaults) undergo periodic integrity verification.

### 11.1 Verification Process

1. **Hash computation.** Every document in the vault has its content hashed with SHA-256. The hash is compared against the stored hash recorded at the time of sealing.
2. **Manifest verification.** The room manifest's SHA-256 hash is recomputed and compared against the stored hash.
3. **Document count verification.** The number of documents in the room is compared against the manifest's `occupied` count.
4. **Cross-reference verification.** Every cross-reference in every document is checked to confirm the target document and target room still exist.
5. **Seal verification.** For sealed documents, the seal timestamp and seal authority are verified against the decision log.

### 11.2 Verification Schedule

| Vault Type | Verification Frequency | Authority |
|---|---|---|
| Standard Vault | Monthly | Automated with manual review |
| Memory Vault | Monthly | Manual only (automated verification runs but a human reviews the results) |
| Decision Log Vault | Weekly | Automated with alert on failure |
| Oral History Vault | Quarterly | Manual |
| Lessons Learned Vault | Quarterly | Automated with manual review |
| Dead Drop Room | Weekly | Automated (content only, no attribution check) |

### 11.3 Integrity Failure Response

If any verification check fails:

1. The room is immediately rendered in an error state: red flickering ambient lighting, a `VAULT INTEGRITY ALERT` banner across the door.
2. An alert is generated in the audit trail with event type `integrity_failure`.
3. The room is locked to all users except the institutional authority.
4. The institutional authority investigates and resolves the discrepancy. If the failure is due to data corruption, the recovery procedures in D6-013 (Disaster Data Recovery) are invoked. If the failure is due to unauthorized modification, the incident response procedures in SEC-011 are invoked.
5. The room is not unlocked until the integrity issue is resolved and a fresh verification passes.

## 12. Physical Metaphor Rules

The basement archive system translates digital concepts into spatial metaphors. The following rules govern how digital concepts are rendered in the spatial interface.

### 12.1 Mapping Table: Digital Concept to Spatial Metaphor

| Digital Concept | Spatial Metaphor | Visual Treatment |
|---|---|---|
| File storage | Shelf slot | Document spine visible on shelf |
| Directory/folder | Room | Walled space with door |
| File system depth | Floor depth (above ground = active, below ground = archived) | Vertical position in building |
| Access control | Locked door | Door style, lock icon, clearance color |
| Read-only | Sealed glass case | Document visible but behind glass, no interaction buttons |
| Append-only (logs) | One-way slot | Documents enter through a wall slot; visible inside but cannot be removed |
| Encryption | Opaque container | Document rendered as a dark, featureless block until decrypted |
| Backup | Shadow copy (ghost outline on shelf) | Faint ghostly duplicate visible on a parallel shelf in the archive |
| Version history | Stack of documents | Multiple versions stacked on the same shelf slot, with the current version on top |
| Search index | Building directory (lobby kiosk) | Glowing terminal in the lobby and on each floor's landing |
| Cross-reference | Neon conduit line | Visible glowing line connecting two documents across rooms/floors |
| Metadata | Document spine label | Text and icons visible on the document's shelf spine |
| Checksum/hash | Integrity seal (wax seal icon) | Red wax seal icon on sealed documents; cracked seal on integrity failure |
| Data corruption | Physical damage (torn pages, burn marks) | Document rendered with visible damage indicators |
| Deletion | Empty shelf slot with dust outline | The slot where a document was shows a faint outline and the text `[REMOVED]` |

### 12.2 The Rule of Spatial Honesty

The spatial metaphor must never lie. If a document exists, it must be visible somewhere in the building (subject to clearance). If a document has been deleted, its absence must be visible. If a room is full, it must look full. If a vault's integrity has failed, the vault must look damaged.

The interface does not hide problems. It does not present a clean facade over a broken system. The spatial metaphor is a window into the institution's real state, not a painting of what the institution wishes it looked like.

### 12.3 The Rule of Spatial Consistency

The same digital concept must always map to the same spatial metaphor. An access-controlled resource is always a locked door, never sometimes a locked door and sometimes a password prompt. A cross-reference is always a neon conduit, never sometimes a conduit and sometimes a hyperlink. Consistency is what makes the spatial metaphor learnable. If the user internalizes the mapping once, they can navigate the entire building without relearning the interface for each room.

---

---

## Cross-Agent Integration Notes

### How Agents 13-15 Interoperate

Agent 13 (Binding) creates the data model. Every document has a binding record. Every room has a manifest. The global search index makes everything findable.

Agent 14 (Clearance) layers access control on top of Agent 13's model. It does not modify the binding or the manifest -- it adds a `clearance_level` field to rooms and a `classification` field to documents, and it defines the rules for when those fields allow or deny access. The lockout rendering is a visual consequence of the access control decision.

Agent 15 (Archives & Memory) extends Agent 13's model downward into the basement. It uses the same binding schema, the same manifest format, and the same room types. But it adds specialized vault types (memory vault, decision log vault, oral history vault, lessons learned vault, dead drop room) with additional properties: seal status, time-lock, anonymous submission, and integrity verification.

### Data Flow

```
Document Created
    |
    v
[F1-S-R001: Intake Room] -- lifecycle_state: "intake"
    |
    v
[F1-W-R00x: Processing Room] -- lifecycle_state: "processing"
    |
    +---> Classification assigned (Agent 14)
    +---> Domain and format determined (Agent 13)
    +---> Target room selected (Agent 13, placement rules)
    |
    v
[F{n}-{wing}-R{nnn}: Target Room] -- lifecycle_state: "placed"
    |
    +---> (years pass, document superseded or retired)
    |
    v
[B1-{wing}-R{nnn}: Recent Archive] -- lifecycle_state: "archived"
    |
    +---> (5+ years pass)
    |
    v
[B2-{wing}-R{nnn}: Historical Archive] -- lifecycle_state: "archived"
    |
    +---> (institutional decision to permanently preserve)
    |
    v
[B3-{wing}-R{nnn}: Deep Cold / Vault] -- lifecycle_state: "sealed"
```

### Manifest File Hierarchy

```
hic-data/
  search/
    global-index.json
    global-index.json.sha256
    domain-index-1.json
    domain-index-2.json
    ...
    domain-index-20.json
  floors/
    F1/
      wings/
        N/
          rooms/
            F1-N-R001/
              room-manifest.json
              room-manifest.json.sha256
            F1-N-R002/
              room-manifest.json
              room-manifest.json.sha256
        E/
          rooms/
            ...
        S/
          rooms/
            F1-S-R001/
              room-manifest.json
              room-manifest.json.sha256
        W/
          rooms/
            ...
    F2/
      wings/
        ...
    ...
    F21/
      wings/
        ...
    B1/
      archive-index.json
      archive-index.json.sha256
      wings/
        N/
          rooms/
            B1-N-R001/
              room-manifest.json
              room-manifest.json.sha256
            ...
        ...
    B2/
      archive-index.json
      archive-index.json.sha256
      wings/
        ...
    B3/
      archive-index.json
      archive-index.json.sha256
      wings/
        ...
  audit/
    audit-log-{year}-{month}.json
    audit-log-{year}-{month}.json.sha256
  clearance/
    user-clearances.json
    user-clearances.json.sha256
```

---

## References

- **HIC-001** -- Skyscraper Master Architecture
- **HIC-002** -- Floor Plan Architecture
- **HIC-005** -- Room Rendering Engine
- **HIC-006** -- Corridor & Hallway Logic
- **SEC-001** -- Threat Model and Security Philosophy
- **SEC-002** -- Access Control Procedures
- **SEC-005** -- Long-Term Cryptographic Survival
- **SEC-011** -- Incident Response Framework
- **D6-001** -- Data Philosophy: What We Keep and Why
- **D6-004** -- Archive Management Procedures
- **D6-013** -- Disaster Data Recovery: When Backups Fail
- **D15-003** -- Ethical Review Process
- **D20-001** -- Institutional Memory Philosophy
- **D20-002** -- Decision Log Operations
- **D20-004** -- Lessons Learned Framework
- **D20-005** -- Oral History Capture Methodology
- **GOV-001** -- Authority Model
- **META-010** -- Letters to the Future Maintainer

---

*This document is part of the Holm Intelligence Complex specification. It provides the knowledge-mapping layer that transforms the skyscraper from a spatial container into a living institutional memory. Every document has a room. Every room has a purpose. Every floor has a domain. Every vault has a seal. The building is the documentation system. The documentation system is the building.*
