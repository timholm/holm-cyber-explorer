# STAGE 4: HOLM INTELLIGENCE COMPLEX -- MASTER BLUEPRINT

## Agents 19-20: Systems Integration, Cross-Layer Standards, Naming Rules, System Synthesis

**Document ID:** STAGE4-HIC-MASTER-BLUEPRINT
**Version:** 3.0.0
**Date:** 2026-02-17
**Status:** Ratified
**Classification:** Master Reference -- This document is the single authoritative specification that binds all other Holm Intelligence Complex documents into one coherent system. When any other HIC specification contradicts this document, this document governs. When this document is silent, the most specific applicable HIC specification governs. When two equally specific specifications disagree, the conflict resolution rules in Section 20.11 apply.

**Supersedes:** STAGE4-HIC-MASTER-BLUEPRINT v2.0.0, v1.0.0

---

## How to Read This Document

This is the keystone. Every other HIC document -- the architectural specifications from Agents 1-2, the visual design language from Agents 3-4, the spatial data model from Agents 5-6, the interaction design from Agents 7-8, the knowledge mapping from Agents 9-10, the rendering pipeline from Agents 11-12, the animation systems from Agents 13-14, the accessibility layer from Agents 15-16, and the performance framework from Agents 17-18 -- all of those documents describe individual subsystems. This document describes how those subsystems connect, communicate, fail, recover, and evolve together.

Agent 19 (Systems Integration) specifies the protocols, naming conventions, event architecture, and error handling standards that allow every layer of the HIC to interoperate without ambiguity. Agent 20 (Chief Architect -- System Synthesis) provides the complete system overview: the master data flow, the boot sequence, the performance budget, the room registry, the build pipeline, and the quality gates that determine when the HIC transitions from specification to implementation.

If you are building the HIC, read this document first. Then read the subsystem specifications. Then return to this document. The first reading gives you the map. The subsystem specifications give you the terrain. The second reading confirms the map matches the terrain.

If you are maintaining the HIC, this document is where you start every investigation. The event catalog tells you what messages flow through the system. The naming conventions tell you how to find any component. The error handling standards tell you what should happen when something breaks. The dependency graph tells you what else might break as a consequence.

If you are extending the HIC, the versioning scheme and compatibility matrix tell you how to add new components without destabilizing existing ones. The system invariants tell you what rules you must never violate. The Phase 2 preparation section tells you where the system is going next.

---

---

# PART ONE: AGENT 19 -- SYSTEMS INTEGRATION

---

## 19.1 Cross-Layer Communication Protocol

The Holm Intelligence Complex comprises ten functional layers. Each layer has a distinct responsibility, a defined interface, and a strict contract governing what it may request from other layers and what it must provide in return.

### 19.1.1 Layer Enumeration

| Layer ID | Layer Name       | Responsibility                                          | Primary Agents |
|----------|------------------|---------------------------------------------------------|----------------|
| L-ARCH   | Architecture     | Building structure, floors, rooms, spatial hierarchy    | 1-2            |
| L-VIS    | Visual           | Neon aesthetic, color system, typography, SVG rendering | 3-4            |
| L-DATA   | Data             | Spatial data model, JSON schemas, document metadata     | 5-6            |
| L-INTR   | Interaction      | Input handling, navigation, zoom, gestures, keyboard    | 7-8            |
| L-KNOW   | Knowledge        | Document mapping, cross-references, search index        | 9-10           |
| L-REND   | Rendering        | SVG pipeline, animation, transitions, compositing       | 11-14          |
| L-ACC    | Accessibility    | ARIA, screen reader support, reduced motion, contrast   | 15-16          |
| L-PERF   | Performance      | Budgets, lazy loading, caching, measurement             | 17-18          |
| L-INTG   | Integration      | Event bus, naming, versioning, error handling            | 19             |
| L-SYNTH  | Synthesis        | System overview, boot, build, deployment, governance    | 20             |

### 19.1.2 Communication Rules

**Rule 1: Layers communicate exclusively through the event bus.** No layer may import functions, objects, or data structures directly from another layer's internal modules. All inter-layer communication passes through published events on the HIC event bus. This rule ensures that any layer can be replaced, upgraded, or disabled without modifying the source code of any other layer.

**Rule 2: Layers may read shared data stores, but only the owning layer may write.** The spatial data model (L-DATA) owns the building structure data. The knowledge layer (L-KNOW) owns the document index. The rendering layer (L-REND) owns the current render state. Any layer may read these stores. Only the owning layer may mutate them. Mutations are announced via events so other layers can react.

**Rule 3: Layer dependencies must be acyclic.** If Layer A depends on Layer B, Layer B must not depend on Layer A. The dependency direction flows downward through the stack: Interaction depends on Data; Data does not depend on Interaction. Rendering depends on Visual; Visual does not depend on Rendering. When bidirectional communication is necessary, it is achieved through events, not through direct dependency.

**Rule 4: Every cross-layer call must specify a timeout.** When one layer requests information from another via the event bus, the requesting layer must specify a maximum wait time. If the response does not arrive within that window, the requesting layer must invoke its fallback behavior as defined in the graceful degradation chain (Section 19.8).

**Rule 5: All cross-layer data is serializable to JSON.** No layer may pass functions, closures, DOM nodes, or other non-serializable objects through the event bus. This constraint ensures that the event bus can be logged, replayed, inspected, and transmitted across contexts (such as between a main thread and a web worker) without loss.

### 19.1.3 Layer Dependency Direction

```
L-INTR (Interaction)
   |
   +----> L-DATA (Data / Spatial Model)
   |         |
   |         +----> L-ARCH (Architecture / Structure)
   |         |
   |         +----> L-KNOW (Knowledge / Documents)
   |
   +----> L-REND (Rendering)
   |         |
   |         +----> L-VIS (Visual / Aesthetic)
   |         |
   |         +----> L-PERF (Performance)
   |
   +----> L-ACC (Accessibility)

L-INTG (Integration) -- observes all layers, owns event bus
L-SYNTH (Synthesis)  -- orchestrates boot, build, deploy
```

The Interaction layer sits at the top because it receives user input and must coordinate across all other layers to produce a response. It never processes data directly -- it delegates to the appropriate layer via events and waits for results.

---

## 19.2 Event Bus Specification

### 19.2.1 Architecture

The HIC event bus is a synchronous publish/subscribe system implemented as a singleton JavaScript module. It supports:

- **Named channels:** Events are published to named channels using dot-notation strings.
- **Wildcard subscriptions:** A subscriber may listen to `navigation.*` to receive all navigation events, or `*` to receive all events system-wide.
- **Priority ordering:** Subscribers declare a priority (integer, lower number equals higher priority). The event bus delivers events to subscribers in priority order within a channel. Default priority is 100.
- **Event metadata:** Every event carries a standard envelope containing the event name, a monotonically increasing sequence number, an ISO 8601 timestamp, the originating layer ID, and an arbitrary payload object.
- **Cancellation:** A subscriber may cancel an event by calling `event.cancel()`. Subsequent subscribers in the priority chain will not receive the event. This mechanism is used by the security layer to block unauthorized navigation.

### 19.2.2 Event Envelope Schema

```json
{
  "event": "navigation.floor.select",
  "seq": 14207,
  "timestamp": "2026-02-17T14:30:00.000Z",
  "source": "L-INTR",
  "payload": {
    "floorId": "F07",
    "previousFloorId": "F03",
    "trigger": "click"
  },
  "cancelled": false
}
```

Every event conforms to this envelope. The `payload` object varies by event type and is specified in the event catalog below. The `source` field must be one of the layer IDs defined in Section 19.1.1. The `seq` field is assigned by the event bus at publish time and is globally unique within a session. The `cancelled` field is set to `true` if any subscriber invokes cancellation; it is `false` by default.

### 19.2.3 Subscription API

```javascript
// Subscribe to a specific event
hic.events.on("navigation.floor.select", handler, { priority: 50 });

// Subscribe to all events in a namespace
hic.events.on("navigation.*", handler);

// Subscribe to all events (system monitor)
hic.events.on("*", handler);

// Unsubscribe
hic.events.off("navigation.floor.select", handler);

// Publish an event
hic.events.emit("navigation.floor.select", {
  floorId: "F07",
  previousFloorId: "F03",
  trigger: "click"
});

// One-time subscription (auto-unsubscribes after first delivery)
hic.events.once("render.transition.end", handler);
```

### 19.2.4 Event Bus Guarantees

1. **Ordering within a channel:** Events published to the same channel are delivered in publication order.
2. **Priority ordering for subscribers:** Within a single event delivery, subscribers are called in priority order (lowest number first).
3. **Synchronous delivery:** The `emit` call does not return until all subscribers have been called (unless a subscriber is explicitly marked as `async`, in which case it is scheduled via `queueMicrotask` and the emit returns immediately after synchronous subscribers complete).
4. **No event loss:** If a subscriber is registered when an event is emitted, it will receive that event. The event bus has no buffer limits, no rate limits, and no silent discard behavior. If no subscribers exist for a channel, the event is silently discarded but still logged if the system monitor is active.

---

## 19.3 Event Catalog

The following sections enumerate every standard event in the HIC system. Custom events may be added by extensions, but they must use the `ext.` namespace prefix (e.g., `ext.plugin.loaded`).

### 19.3.1 Navigation Events

**`navigation.floor.select`**

- **Source:** L-INTR
- **Description:** User has selected a floor. Triggers floor transition animation and data loading.
- **Payload:**

| Field            | Type           | Required | Description                                                    |
|------------------|----------------|----------|----------------------------------------------------------------|
| `floorId`        | String         | Yes      | Target floor identifier (B3, B2, B1, F01-F20, R1, R2)         |
| `previousFloorId`| String or null | Yes      | Floor the user is departing from. Null on initial load.        |
| `trigger`        | String         | Yes      | One of: `"click"`, `"keyboard"`, `"url"`, `"api"`, `"restore"` |

**`navigation.room.enter`**

- **Source:** L-INTR
- **Description:** User has navigated into a specific room. Triggers room zoom and content loading.
- **Payload:**

| Field         | Type   | Required | Description                                           |
|---------------|--------|----------|-------------------------------------------------------|
| `roomId`      | String | Yes      | Full room identifier (see Section 19.5.2)             |
| `floorId`     | String | Yes      | Floor containing the room                             |
| `entryMethod` | String | Yes      | One of: `"zoom"`, `"direct"`, `"search"`, `"link"`   |

**`navigation.room.exit`**

- **Source:** L-INTR
- **Description:** User is leaving a room. Triggers reverse zoom and cleanup.
- **Payload:**

| Field         | Type           | Required | Description                                             |
|---------------|----------------|----------|---------------------------------------------------------|
| `roomId`      | String         | Yes      | Room being exited                                       |
| `floorId`     | String         | Yes      | Floor containing the room                               |
| `exitMethod`  | String         | Yes      | One of: `"zoom-out"`, `"back"`, `"navigate"`, `"close"` |
| `destination` | String or null | No       | Where the user is going next, if known                  |

**`navigation.zoom.change`**

- **Source:** L-INTR
- **Description:** Zoom level has changed. Levels: L0 (city), L1 (building), L2 (floor), L3 (room), L4 (document).
- **Payload:**

| Field               | Type    | Required | Description                                     |
|---------------------|---------|----------|-------------------------------------------------|
| `zoomLevel`         | Integer | Yes      | New zoom level (0-4)                            |
| `previousZoomLevel` | Integer | Yes      | Previous zoom level (0-4)                       |
| `target`            | String  | Yes      | ID of element being zoomed to (floor, room, doc)|
| `origin`            | Object  | No       | `{ x, y }` screen coordinates of zoom origin    |

### 19.3.2 Document Events

**`document.open`**

- **Source:** L-KNOW
- **Description:** A document has been opened for viewing within a room.
- **Payload:**

| Field        | Type    | Required | Description                                           |
|--------------|---------|----------|-------------------------------------------------------|
| `documentId` | String  | Yes      | Document identifier (see Section 19.5.3)              |
| `roomId`     | String  | Yes      | Room containing the document                          |
| `floorId`    | String  | Yes      | Floor containing the room                             |
| `format`     | String  | Yes      | One of: `"markdown"`, `"html"`, `"pdf"`, `"image"`, `"data"` |
| `byteSize`   | Integer | Yes      | Size of the document content in bytes                 |

**`document.close`**

- **Source:** L-KNOW
- **Description:** A document has been closed.
- **Payload:**

| Field        | Type   | Required | Description                                                 |
|--------------|--------|----------|-------------------------------------------------------------|
| `documentId` | String | Yes      | Document identifier                                         |
| `roomId`     | String | Yes      | Room the document was displayed in                          |
| `reason`     | String | Yes      | One of: `"user"`, `"navigation"`, `"error"`, `"security"`  |

### 19.3.3 Security Events

**`security.clearance.check`**

- **Source:** L-DATA
- **Description:** A clearance check is being performed before granting access to a room or document.
- **Payload:**

| Field               | Type    | Required | Description                                       |
|---------------------|---------|----------|---------------------------------------------------|
| `entityId`          | String  | Yes      | ID of the floor, room, or document being accessed |
| `entityType`        | String  | Yes      | One of: `"floor"`, `"room"`, `"document"`         |
| `requiredClearance` | Integer | Yes      | Clearance level required (0-5)                    |
| `userClearance`     | Integer | Yes      | Current user clearance level (0-5)                |

**`security.clearance.denied`**

- **Source:** L-DATA
- **Description:** Access has been denied. The UI must display the appropriate denial animation.
- **Payload:**

| Field               | Type    | Required | Description                                       |
|---------------------|---------|----------|---------------------------------------------------|
| `entityId`          | String  | Yes      | ID of the entity access was denied to             |
| `entityType`        | String  | Yes      | One of: `"floor"`, `"room"`, `"document"`         |
| `requiredClearance` | Integer | Yes      | Clearance level that was required                 |
| `userClearance`     | Integer | Yes      | User clearance level that was insufficient        |
| `reason`            | String  | Yes      | Human-readable denial explanation                 |

### 19.3.4 Rendering Events

**`render.floor.loaded`**

- **Source:** L-REND
- **Description:** A floor's SVG has been fully loaded and rendered into the DOM.
- **Payload:**

| Field          | Type    | Required | Description                                    |
|----------------|---------|----------|------------------------------------------------|
| `floorId`      | String  | Yes      | Floor that was loaded                          |
| `elementCount` | Integer | Yes      | Number of SVG elements rendered on this floor  |
| `renderTimeMs` | Float   | Yes      | Milliseconds taken to render the floor         |

**`render.transition.start`**

- **Source:** L-REND
- **Description:** A visual transition has begun (floor change, zoom, room entry, etc.).
- **Payload:**

| Field            | Type    | Required | Description                                                                                      |
|------------------|---------|----------|--------------------------------------------------------------------------------------------------|
| `transitionType` | String  | Yes      | One of: `"floor-change"`, `"zoom"`, `"room-enter"`, `"room-exit"`, `"document-open"`, `"document-close"` |
| `from`           | String  | Yes      | Origin state identifier (floor ID, room ID, or zoom level)                                       |
| `to`             | String  | Yes      | Destination state identifier                                                                     |
| `durationMs`     | Integer | Yes      | Planned duration of the transition in milliseconds                                               |
| `easing`         | String  | Yes      | CSS easing function (e.g., `"cubic-bezier(0.4, 0, 0.2, 1)"`)                                    |

**`render.transition.end`**

- **Source:** L-REND
- **Description:** A visual transition has completed.
- **Payload:**

| Field              | Type    | Required | Description                                       |
|--------------------|---------|----------|---------------------------------------------------|
| `transitionType`   | String  | Yes      | Same type as the corresponding `transition.start` |
| `from`             | String  | Yes      | Origin state identifier                           |
| `to`               | String  | Yes      | Destination state identifier                      |
| `actualDurationMs` | Float   | Yes      | Actual elapsed time of the transition             |
| `dropped`          | Integer | Yes      | Number of frames dropped during the transition    |

### 19.3.5 System Events

**`system.import.detected`**

- **Source:** L-DATA
- **Description:** New documents have been detected for import from an external source.
- **Payload:**

| Field           | Type    | Required | Description                                    |
|-----------------|---------|----------|------------------------------------------------|
| `sourceType`    | String  | Yes      | One of: `"usb"`, `"federation"`, `"manual"`    |
| `sourceId`      | String  | Yes      | Identifier for the import source               |
| `documentCount` | Integer | Yes      | Number of new documents detected               |
| `byteSize`      | Integer | Yes      | Total size in bytes of incoming documents       |

**`system.cache.updated`**

- **Source:** L-PERF
- **Description:** The caching layer has been updated (floor data, document content, etc.).
- **Payload:**

| Field            | Type    | Required | Description                                       |
|------------------|---------|----------|---------------------------------------------------|
| `cacheType`      | String  | Yes      | One of: `"floor"`, `"room"`, `"document"`, `"search"`, `"svg"` |
| `entriesAdded`   | Integer | Yes      | Number of new entries added to the cache          |
| `entriesEvicted` | Integer | Yes      | Number of entries removed to make room            |
| `sizeBytes`      | Integer | Yes      | Current total size of this cache in bytes         |

**`system.boot.stage`**

- **Source:** L-SYNTH
- **Description:** Boot sequence has reached a new stage. Used by the loading screen.
- **Payload:**

| Field         | Type    | Required | Description                               |
|---------------|---------|----------|-------------------------------------------|
| `stage`       | String  | Yes      | Stage name (e.g., `"render"`, `"events"`) |
| `stageIndex`  | Integer | Yes      | Zero-based index of current stage         |
| `totalStages` | Integer | Yes      | Total number of boot stages               |
| `elapsedMs`   | Float   | Yes      | Milliseconds elapsed since boot began     |

**`system.error`**

- **Source:** Any layer
- **Description:** A system error has occurred. See error handling standards (Section 19.9).
- **Payload:**

| Field         | Type    | Required | Description                                              |
|---------------|---------|----------|----------------------------------------------------------|
| `errorCode`   | String  | Yes      | Structured error code (see Section 19.9.1)               |
| `severity`    | Integer | Yes      | 0 = Critical, 1 = Error, 2 = Warning, 3 = Info          |
| `layerId`     | String  | Yes      | Layer where the error originated                         |
| `message`     | String  | Yes      | Human-readable error description                         |
| `recoverable` | Boolean | Yes      | Whether the system can continue operating                |
| `context`     | Object  | No       | Additional context specific to the error type            |

**`system.ready`**

- **Source:** L-SYNTH
- **Description:** The system has completed boot and is ready for user interaction.
- **Payload:**

| Field              | Type    | Required | Description                                 |
|--------------------|---------|----------|---------------------------------------------|
| `bootTimeMs`       | Float   | Yes      | Total boot duration in milliseconds         |
| `floorsAvailable`  | Integer | Yes      | Number of floors with loaded manifests      |
| `documentsIndexed` | Integer | Yes      | Number of documents in the search index     |

---

## 19.4 Event Sequencing Rules

Certain events must occur in specific sequences. Violating these sequences indicates a bug in the implementation and must be flagged by the integration test suite.

**Floor Navigation Sequence:**

```
navigation.floor.select
  --> security.clearance.check
  --> [security.clearance.denied OR continue]
  --> render.transition.start (type: floor-change)
  --> render.floor.loaded (if floor SVG not already cached)
  --> render.transition.end
```

**Room Entry Sequence:**

```
navigation.room.enter
  --> security.clearance.check
  --> [security.clearance.denied OR continue]
  --> navigation.zoom.change (to L3)
  --> render.transition.start (type: room-enter)
  --> render.transition.end
```

**Document Open Sequence:**

```
document.open
  --> security.clearance.check
  --> [security.clearance.denied OR continue]
  --> navigation.zoom.change (to L4)
  --> render.transition.start (type: document-open)
  --> render.transition.end
```

**Room Exit Sequence:**

```
navigation.room.exit
  --> document.close (if a document was open, reason: "navigation")
  --> navigation.zoom.change (to L2)
  --> render.transition.start (type: room-exit)
  --> render.transition.end
```

**Import Detection Sequence:**

```
system.import.detected
  --> security.clearance.check (for each incoming document)
  --> system.cache.updated (after import processing completes)
```

---

## 19.5 Naming Conventions

Naming is not cosmetic. In a system with hundreds of rooms and thousands of documents, a consistent naming scheme is the difference between navigable architecture and chaos. The following conventions are mandatory. No exceptions. No creative variations. These are machine-parseable identifiers, and every tool in the HIC pipeline depends on them.

### 19.5.1 Floor Identifiers

Floors are identified by a prefix letter and a zero-padded number.

| Identifier | Name                     | Zone       | Description                                                  |
|------------|--------------------------|------------|--------------------------------------------------------------|
| B3         | Deep Archive             | Basement   | Long-term cold storage, cryptographic vaults, disaster recovery |
| B2         | Infrastructure Core      | Basement   | Power systems, environmental controls, hardware maintenance  |
| B1         | Security Operations      | Basement   | Access control, audit systems, intrusion detection           |
| F01        | Public Atrium            | Lower      | Welcome, orientation, public documents, visitor access       |
| F02        | Administration           | Lower      | Governance charter, constitutional documents, org structure  |
| F03        | Ethics & Philosophy      | Lower      | Ethical framework, philosophical foundations, values         |
| F04        | Infrastructure Docs      | Lower      | Power, network, environmental, hardware documentation        |
| F05        | Data Architecture        | Lower      | Data models, schemas, storage formats, migration procedures  |
| F06        | Automation Systems       | Mid        | Build pipelines, CI/CD, scripting, task automation           |
| F07        | Intelligence Operations  | Mid        | SIGINT, OSINT, analysis frameworks, threat modeling          |
| F08        | Knowledge Management     | Mid        | Indexing, search, cross-referencing, taxonomy                |
| F09        | Education & Training     | Mid        | Onboarding, skill development, certification, workshops     |
| F10        | Communication Systems    | Mid        | Internal messaging, external relations, broadcast protocols  |
| F11        | Research Laboratory      | Upper      | Experimental systems, prototypes, R&D documentation          |
| F12        | Federation Hub           | Upper      | Multi-instance coordination, import/export, peering          |
| F13        | Platform Engineering     | Upper      | OS, containers, virtualization, deployment infrastructure    |
| F14        | Quality Assurance        | Upper      | Testing frameworks, audit trails, compliance verification    |
| F15        | Interface Design         | Upper      | UI/UX specifications, design system, accessibility standards |
| F16        | Evolution & Memory       | Upper      | Versioning, changelog, institutional memory, decay tracking  |
| F17        | Advanced Security        | High       | Cryptographic systems, zero-trust architecture, red team ops |
| F18        | Advanced Data Systems    | High       | Distributed storage, replication, conflict resolution        |
| F19        | Systems Integration      | High       | Cross-layer protocols, event bus, naming, testing            |
| F20        | Command Center           | Penthouse  | Chief Architect overview, master dashboards, system synthesis|
| R1         | Observation Deck         | Rooftop    | System health visualization, real-time telemetry             |
| R2         | Broadcast Antenna Array  | Rooftop    | External publication, RSS feeds, public API, federation      |

**Format rule:** Basement floors use `B` followed by a single digit (no zero-padding). Standard floors use `F` followed by a two-digit zero-padded number (`F01` through `F20`). Rooftop floors use `R` followed by a single digit. This convention ensures consistent string sorting: all basements sort before all floors, all floors sort before all rooftops.

### 19.5.2 Room Identifiers

Rooms are identified by a composite string:

```
{floor}-{type}-{name}-{number}
```

**Components:**

- `{floor}`: The floor identifier (e.g., `F07`).
- `{type}`: A three-letter uppercase room type code from the table below.
- `{name}`: An uppercase alphanumeric slug describing the room's specific function. Maximum 12 characters. Hyphens allowed within the name.
- `{number}`: A two-digit zero-padded sequence number, starting at `01`.

**Room Type Codes:**

| Code | Room Type        | Description                                            |
|------|------------------|--------------------------------------------------------|
| LAB  | Laboratory       | Experimental, research, or prototype workspace         |
| ARC  | Archive          | Long-term document storage                             |
| OPS  | Operations       | Active operational workspace                           |
| CMD  | Command          | Command and control, dashboards, monitoring            |
| VLT  | Vault            | High-security storage, encrypted archives              |
| LIB  | Library          | Reference collection, reading room                     |
| WRK  | Workshop         | Hands-on maintenance, repair, fabrication              |
| COM  | Communications   | Messaging, broadcast, relay equipment                  |
| SEC  | Security         | Access control, surveillance, intrusion detection      |
| SRV  | Server           | Computing equipment, processing nodes                  |
| UTL  | Utility          | Power, HVAC, plumbing, environmental support           |
| BRF  | Briefing         | Presentation, orientation, training delivery           |
| QUA  | Quarantine       | Isolation zone for untrusted media or data             |
| HUB  | Hub              | Central routing, distribution, interconnection point   |
| MON  | Monitor          | Observation, telemetry, passive surveillance           |

**Examples:**

- `F07-LAB-SIGINT-01` -- Floor 7, Laboratory, SIGINT analysis, room 1
- `F07-LAB-SIGINT-02` -- Floor 7, Laboratory, SIGINT analysis, room 2
- `B1-SEC-ACCESS-01` -- Basement 1, Security, Access control, room 1
- `F01-LIB-WELCOME-01` -- Floor 1, Library, Welcome/orientation, room 1
- `B3-VLT-CRYPTO-01` -- Basement 3, Vault, Cryptographic storage, room 1
- `R1-MON-HEALTH-01` -- Rooftop 1, Monitor, System health, room 1

### 19.5.3 Document Identifiers

Documents are identified by:

```
D{domain}-{number}-{slug}
```

- `D`: Literal prefix character.
- `{domain}`: Two-digit zero-padded domain number (01-20).
- `{number}`: Three-digit zero-padded sequence number within the domain, starting at `001`.
- `{slug}`: Lowercase alphanumeric slug using hyphens as separators. Maximum 40 characters.

**Examples:**

- `D03-001-threat-model` -- Domain 3 (Security), document 1, threat model
- `D07-014-osint-collection-guide` -- Domain 7 (Intelligence), document 14
- `D01-001-institutional-charter` -- Domain 1 (Governance), document 1, charter
- `D16-003-version-history-policy` -- Domain 16 (Evolution), document 3

### 19.5.4 File Naming Conventions

| File Type                | Pattern                           | Example                          |
|--------------------------|-----------------------------------|----------------------------------|
| Floor blueprint (data)   | `{floor-id}.blueprint.json`       | `F07.blueprint.json`             |
| Floor visual (SVG)       | `{floor-id}.floor.svg`            | `F07.floor.svg`                  |
| Room detail (data)       | `{room-id}.room.json`             | `F07-LAB-SIGINT-01.room.json`   |
| Document content         | `{doc-id}.md`                     | `D03-001-threat-model.md`       |
| Document metadata        | `{doc-id}.meta.json`              | `D03-001-threat-model.meta.json` |
| Building master file     | `building.json`                   | `building.json`                  |
| Floor manifest           | `{floor-id}.manifest.json`        | `F07.manifest.json`              |
| CSS component style      | `hic-{component}.css`             | `hic-elevator.css`               |
| JS module                | `hic.{subsystem}.{function}.js`   | `hic.navigation.zoom.js`        |
| Test file                | `{module}.test.js`                | `hic.navigation.zoom.test.js`   |
| Animation definition     | `{floor-id}.anim.json`            | `F07.anim.json`                  |

### 19.5.5 CSS Class Naming

All CSS classes use the `hic-` prefix:

```
hic-{component}-{modifier}
```

**Examples:**

- `hic-building-exterior` -- The building's outer shell at L0/L1 zoom
- `hic-floor-active` -- The currently selected floor
- `hic-floor-locked` -- A floor the user lacks clearance to access
- `hic-room-highlighted` -- A room with mouse hover or focus
- `hic-room-classified` -- A room requiring elevated clearance
- `hic-neon-glow` -- The standard neon glow effect
- `hic-neon-pulse` -- Animated pulsing neon effect
- `hic-elevator-shaft` -- The vertical navigation column
- `hic-elevator-car` -- The moving elevator indicator
- `hic-document-panel` -- The document reading pane
- `hic-transition-active` -- Applied during any transition animation
- `hic-zoom-l0` through `hic-zoom-l4` -- Applied to body for current zoom level

BEM is not used. The naming system is flat with single-level modifiers. Nesting is handled by DOM hierarchy, not class name encoding.

### 19.5.6 JavaScript Module Naming

All JavaScript modules follow the dot-notation pattern:

```
hic.{subsystem}.{function}
```

**Subsystem registry:**

| Subsystem      | Responsibility                                  |
|----------------|-------------------------------------------------|
| `hic.events`   | Event bus (publish, subscribe, unsubscribe)      |
| `hic.nav`      | Navigation engine (floor select, room enter/exit)|
| `hic.zoom`     | Zoom level management and transitions            |
| `hic.render`   | SVG rendering pipeline and DOM management        |
| `hic.anim`     | Animation sequencing and easing                  |
| `hic.data`     | Spatial data model, building.json access         |
| `hic.docs`     | Document loading, rendering, search              |
| `hic.security` | Clearance checks, access control                 |
| `hic.cache`    | Caching layer, lazy loading, preloading          |
| `hic.a11y`     | Accessibility layer, ARIA management             |
| `hic.perf`     | Performance monitoring and budget enforcement    |
| `hic.boot`     | Boot sequence orchestration                      |
| `hic.state`    | State management, localStorage persistence       |
| `hic.ui`       | UI components (elevator, panels, overlays)       |

---

## 19.6 Versioning Scheme

All HIC definition files use Semantic Versioning 2.0.0:

```
MAJOR.MINOR.PATCH
```

- **MAJOR** increments when backward compatibility is broken.
- **MINOR** increments when new features are added in a backward-compatible manner.
- **PATCH** increments when bugs are fixed without changing the interface.

**Version declaration in every JSON file:**

```json
{
  "hicVersion": "1.0.0",
  "fileType": "floor-blueprint",
  "floorId": "F07"
}
```

The `hicVersion` field indicates the minimum HIC system version required. The system must refuse to load files whose `hicVersion` exceeds its own version.

---

## 19.7 Compatibility Matrix

| Component A           | Component B          | Compatibility Rule             | Notes                                           |
|-----------------------|----------------------|--------------------------------|-------------------------------------------------|
| building.json         | HIC Core             | building.hicVersion <= core    | Building file must not exceed core version       |
| Floor blueprint       | building.json        | Same MAJOR version             | Floor blueprints must match building MAJOR       |
| Room definition       | Floor blueprint      | Same MAJOR.MINOR version       | Rooms are tightly coupled to their floor         |
| Document metadata     | Knowledge layer      | Same MAJOR version             | Document format evolves with knowledge layer     |
| SVG floor file        | Visual layer         | Same MAJOR version             | SVG structure must match visual layer expectations|
| Animation definition  | Rendering layer      | Same MAJOR.MINOR version       | Animation formats are rendering-specific         |
| CSS theme             | Visual layer         | Same MAJOR version             | Theme variables must match visual layer schema   |
| Search index          | Knowledge layer      | Exact version match            | Index format is tightly coupled to search engine |
| Service worker        | HIC Core             | Same MAJOR.MINOR version       | Cache strategy must match asset structure        |

**Enforcement:** The boot sequence performs version checks during stage 3. Incompatibilities trigger `system.error` with code `E-DAT-VER-001` and enter degraded mode.

---

## 19.8 Integration Testing Specification

### 19.8.1 Test Categories

| Category            | Scope                                                     | Trigger         | Pass Rate |
|---------------------|-----------------------------------------------------------|-----------------|-----------|
| Unit                | Single module in isolation                                | Every commit    | 100%      |
| Layer               | All modules within one layer                              | Every commit    | 100%      |
| Cross-layer         | Communication between two layers via event bus            | Every commit    | 100%      |
| Sequence            | Complete event sequences (full floor navigation)          | Every commit    | 100%      |
| Performance         | Budget compliance under synthetic load                    | Nightly         | 95%       |
| Regression          | All rooms render correctly, all documents load            | Nightly         | 100%      |
| Accessibility       | Screen reader navigation, keyboard-only traversal         | Weekly          | 100%      |
| Full system         | Complete boot, navigate every floor, open sample documents| Pre-release     | 100%      |

### 19.8.2 Cross-Layer Test Protocol

1. **Setup:** Initialize the event bus. Register a mock subscriber on the target layer's channel.
2. **Act:** Invoke the source layer's public API (which should emit an event).
3. **Assert:** Verify the mock subscriber received the expected event with the correct payload.
4. **Teardown:** Unsubscribe all mocks. Reset event bus state.

Cross-layer tests must not depend on DOM state. They test the event contract, not the visual result.

### 19.8.3 Test File Naming

```
test/integration/{source-layer}--{target-layer}.test.js
```

Example: `test/integration/L-INTR--L-DATA.test.js`

---

## 19.9 Error Handling Standards

### 19.9.1 Error Code Format

```
E-{LAYER}-{CATEGORY}-{NUMBER}
```

**Layer abbreviations:** ARC, VIS, DAT, INT, KNO, REN, ACC, PRF, ITG, SYN

**Error categories:**

| Category | Meaning                                       |
|----------|-----------------------------------------------|
| LOD      | Load failure (file not found, parse error)     |
| REN      | Render failure (SVG error, paint failure)      |
| NAV      | Navigation failure (invalid target, blocked)   |
| SEC      | Security failure (clearance denied, tamper)    |
| TIM      | Timeout (operation exceeded time budget)       |
| DAT      | Data error (invalid schema, missing field)     |
| VER      | Version mismatch (incompatible file version)   |
| STA      | State error (invalid state transition)         |
| EVT      | Event error (delivery failure, bus error)      |

**Examples:**

- `E-REN-LOD-001` -- Rendering layer, load error: SVG file failed to load
- `E-DAT-VER-001` -- Data layer, version error: building.json version mismatch
- `E-INT-NAV-001` -- Interaction layer, navigation error: target floor does not exist
- `E-KNO-LOD-001` -- Knowledge layer, load error: document file not found
- `E-PRF-TIM-001` -- Performance layer, timeout: render exceeded budget

### 19.9.2 Error Severity Levels

| Level    | Code | System Response                                       |
|----------|------|-------------------------------------------------------|
| Critical | 0    | Display error screen. Log fully. Halt all operations. |
| Error    | 1    | Engage degradation. Log. Notify user.                 |
| Warning  | 2    | Log. No user notification unless repeated 3+ times.   |
| Info     | 3    | Log only. No user notification.                       |

### 19.9.3 Error Handling Rules

1. **Every error must be published to the event bus.** No silent failures.
2. **Every error must include a recovery suggestion** in the `context.fallback` field.
3. **Errors must not propagate as uncaught exceptions.** Each layer catches its own exceptions and converts them to `system.error` events.
4. **Error rate monitoring:** More than 10 errors of severity 1+ within 60 seconds triggers degradation to D3.

---

## 19.10 Graceful Degradation Chain

### 19.10.1 Degradation Levels

| Level | Name           | Description                                                                |
|-------|----------------|----------------------------------------------------------------------------|
| D0    | Full Function  | All systems operational. Full neon aesthetic, animations, interactivity.   |
| D1    | Reduced Visual | Animations disabled. Static SVG rendering. Full navigation preserved.     |
| D2    | Simplified     | SVG rendering replaced with HTML table layout. Navigation by links.       |
| D3    | Text Only      | All visual rendering disabled. Document list with hyperlinks.             |
| D4    | Offline Cache  | No live data. Serve from localStorage/service worker cache only.          |
| D5    | Error Screen   | Nothing works. Display diagnostic information and recovery instructions.  |

### 19.10.2 Degradation Triggers

| Failure Condition                    | Target Level |
|--------------------------------------|--------------|
| Animation engine fails               | D1           |
| SVG rendering engine fails           | D2           |
| building.json fails to load          | D3           |
| JavaScript execution fails entirely  | D3           |
| All network requests fail            | D4           |
| HTML parsing fails                   | D5           |

### 19.10.3 Degradation Rules

1. **Monotonic within a session.** Once degraded, the system does not attempt to return to a higher level. User must reload.
2. **Independently testable.** Each level can be activated via `?degrade=D2` URL parameter.
3. **User notification mandatory at D2 and below.**
4. **Content accuracy preserved at every level.** Degradation removes presentation, never alters document content.

---

---

# PART TWO: AGENT 20 -- CHIEF ARCHITECT -- SYSTEM SYNTHESIS

---

## 20.1 Complete System Overview

The Holm Intelligence Complex is a spatial documentation interface that presents a sovereign, air-gapped documentation system as an explorable cyberpunk skyscraper. The building metaphor is not decorative. It is structural. Every architectural decision -- floors, rooms, elevators, security clearances, neon lighting -- maps directly to a documentation concept.

Twenty agents designed this system:

| Agent | Subsystem                    | Deliverable                                              |
|-------|------------------------------|----------------------------------------------------------|
| 1     | Structural Architecture      | Floor plan, room layout, spatial hierarchy               |
| 2     | Zoning & Clearance           | Security zones, access levels, restricted areas          |
| 3     | Visual Identity              | Neon palette, glow effects, cyberpunk aesthetic          |
| 4     | Typography & Iconography     | Font stack, icon set, readability standards              |
| 5     | Spatial Data Model           | building.json schema, floor manifests, room data         |
| 6     | Document-Space Mapping       | How documents map to rooms, shelving, capacity           |
| 7     | Navigation Engine            | Elevator, floor select, room entry, zoom levels          |
| 8     | Input Handling               | Mouse, keyboard, touch, gamepad, accessibility input     |
| 9     | Knowledge Graph              | Cross-references, related docs, dependency map           |
| 10    | Search & Retrieval           | Full-text search, faceted filtering, result ranking      |
| 11    | SVG Rendering Pipeline       | Floor rendering, room detail, dynamic elements           |
| 12    | Compositing & Layering       | Z-ordering, overlay management, depth effects            |
| 13    | Animation System             | Transitions, micro-interactions, state animations        |
| 14    | Particle & Ambient Effects   | Neon flicker, data rain, ambient atmosphere              |
| 15    | Accessibility Architecture   | ARIA roles, screen reader, reduced motion                |
| 16    | Inclusive Design             | Color blindness, dyslexia support, high contrast         |
| 17    | Performance Architecture     | Budgets, lazy loading, code splitting, caching           |
| 18    | Monitoring & Optimization    | Runtime metrics, bottleneck detection, profiling         |
| 19    | Systems Integration          | Event bus, naming, versioning, error handling            |
| 20    | Chief Architect              | This document. System synthesis. Master reference.       |

---

## 20.2 Master Data Flow Diagram

```
                         +-----------------------------------------------+
                         |         HOLM INTELLIGENCE COMPLEX              |
                         |         Master Data Flow (v3.0)                |
                         +-----------------------------------------------+

  +----------+
  |   USER   |
  |  INPUT   |--+
  +----------+  |
                |  click / key / touch / scroll / gamepad
                v
  +----------------------+        +----------------------+
  |   INPUT HANDLER      |        |   URL / STATE        |
  |   (hic.ui.input)     |        |   RESTORE            |
  |                      |        |   (hic.state)        |
  |  Normalize input     |        |                      |
  |  Detect gesture      |        |  localStorage        |
  |  Map to action       |        |  URL hash            |
  +--------+-------------+        +----------+-----------+
           |                                 |
           |  navigation.floor.select        |  navigation.floor.select
           |  navigation.room.enter          |  (trigger: "restore")
           |  navigation.zoom.change         |
           v                                 v
  +---------------------------------------------------------+
  |                      EVENT BUS                           |
  |                      (hic.events)                        |
  |                                                          |
  |  Publish --> Priority Sort --> Deliver --> Log           |
  +----+-----------+-----------+-----------+--------+--------+
       |           |           |           |        |
       v           v           v           v        v
  +---------+ +---------+ +--------+ +--------+ +------+
  |SECURITY | |NAVIGATION| |KNOWLEDGE| |RENDERING| | A11Y |
  | LAYER   | | ENGINE  | | LAYER  | | PIPELINE| |LAYER |
  |         | |         | |        | |         | |      |
  |Clearance| |Zoom calc| |Doc look| |SVG load | |ARIA  |
  |check    | |Floor    | |Search  | |Compose  | |Focus |
  |Deny/    | |Room     | |X-ref   | |Animate  | |Screen|
  |Allow    | |State    | |Index   | |Paint    | |Reader|
  +----+----+ +----+----+ +---+----+ +----+----+ +--+---+
       |           |           |           |         |
       |           v           v           |         |
       |     +-------------------------+   |         |
       |     |  SPATIAL DATA MODEL     |   |         |
       |     |  (building.json +       |   |         |
       |     |   floor manifests +     |   |         |
       |     |   room definitions)     |   |         |
       |     +-------------------------+   |         |
       |                                   v         v
       |                             +-----------+-------+
       |                             |  SVG DOM OUTPUT    |
       |                             |                    |
       |         denied              |  Floors, Rooms,    |
       |     +--------------+        |  Documents,        |
       +---->| ACCESS CTRL  |        |  Effects, Neon,    |
             | RESPONSE     |        |  ARIA attributes   |
             |              |        +--------+-----------+
             | Lock anim    |                 |
             | Denial msg   |                 v
             | Red flash    |        +------------------+
             +--------------+        |   USER DISPLAY   |
                                     |                  |
                                     |  <svg> viewport  |
                                     |  The HIC as the  |
                                     |  user sees it    |
                                     +------------------+

  +-----------------------------------------------+  +---------------------------+
  |            PERFORMANCE LAYER                   |  |      CACHE LAYER          |
  |            (hic.perf)                          |  |      (hic.cache)          |
  |  Monitors all layers. Enforces budgets.        |  |  Intercepts data fetches. |
  |  Triggers degradation on budget violation.     |  |  Serves from cache.       |
  +-----------------------------------------------+  |  Preloads adjacent floors. |
                                                      +---------------------------+
```

**Flow narrative:** The user provides input through any supported channel. The Input Handler normalizes the raw browser event into a standard action and publishes a navigation event to the Event Bus. The Event Bus delivers the event to all subscribed layers in priority order. The Security Layer intercepts navigation events first (priority 10) and performs clearance checks; if denied, the event is cancelled and the Access Control Response renders a denial animation. If allowed, the Navigation Engine calculates the new state, the Knowledge Layer retrieves required document data, and the Rendering Pipeline produces the SVG output. The Accessibility Layer observes all events and updates ARIA attributes in parallel. State changes are persisted to localStorage.

---

## 20.3 Component Dependency Graph

```
                    +-------------+
                    |  index.html |
                    +------+------+
                           |
                    +------v------+
                    |  hic.boot   |
                    +------+------+
                           |
              +------------+------------+
              |            |            |
       +------v------+ +--v----+ +-----v-------+
       | hic.events  | |hic.ui | | hic.state   |
       +------+------+ +--+----+ +-----+-------+
              |            |            |
    +---------+------------+------------+
    |         |            |
+---v----+ +--v------+ +--v--------+
|hic.data| |hic.nav  | |hic.render |
+---+----+ +---+-----+ +---+-------+
    |          |            |
    |    +-----+      +----+----------+
    |    |            |    |          |
+---v--+ | +--v---+  | +--v----+ +--v----+
|hic.  | | |hic.  |  | |hic.   | |hic.   |
|docs  | | |zoom  |  | |anim   | |perf   |
+------+ | +------+  | +-------+ +-------+
         |           |
    +----v-----+ +---v------+
    |hic.      | |hic.      |
    |security  | |cache     |
    +----------+ +----------+

    +----------+
    |hic.a11y  |  (parallel, no hard dependencies,
    +----------+   observes events only)
```

**Key rules:** `hic.events` is loaded first after `hic.boot`. `hic.a11y` has no hard dependencies and can fail without affecting core functionality. The three core pillars (`hic.data`, `hic.nav`, `hic.render`) communicate through events, not direct imports.

---

## 20.4 Boot Sequence Specification

The boot sequence is deterministic: eight stages, executed in order. Each stage must complete before the next begins.

| Stage | Action                                     | Event Emitted                                           | Failure Level | Budget  |
|-------|--------------------------------------------|---------------------------------------------------------|---------------|---------|
| 0     | Load `index.html`                          | (none -- HTML not yet parsed)                           | D5            | 200ms   |
| 1     | Initialize rendering engine                | `system.boot.stage { stage: "render", stageIndex: 0 }`  | D2            | 100ms   |
| 2     | Initialize event bus                       | `system.boot.stage { stage: "events", stageIndex: 1 }`  | D3            | 10ms    |
| 3     | Load `building.json` (complete structure)  | `system.boot.stage { stage: "data", stageIndex: 2 }`    | D3            | 150ms   |
| 4     | Render building exterior (L0 zoom)         | `system.boot.stage { stage: "exterior", stageIndex: 3 }` | D2           | 200ms   |
| 5     | Load floor manifests (lazy, on demand)     | `system.boot.stage { stage: "floors", stageIndex: 4 }`  | D1 (partial)  | async   |
| 6     | Bind interaction handlers                  | `system.boot.stage { stage: "input", stageIndex: 5 }`   | D2            | 50ms    |
| 7     | Restore navigation state from localStorage| `system.ready`                                          | D0 (fresh)    | 50ms    |

**Stage 0:** Browser fetches and parses the minimal HTML shell. Contains only essential `<head>` elements and a loading screen placeholder.

**Stage 1:** SVG rendering pipeline initializes: root `<svg>` element created, viewport transform matrix set up, SVG filter definitions for neon glow registered.

**Stage 2:** Event bus instantiated. System monitor subscriber registered at priority 0. Error handler subscriber registered at priority 1.

**Stage 3:** Complete building structure fetched and parsed. All 25 floor definitions loaded, validated against schema, version compatibility checked.

**Stage 4:** Building exterior SVG generated from building.json and rendered at L0 zoom. Neon glow effects applied. This is the first visual content beyond the loading screen.

**Stage 5:** Floor manifests loaded lazily. Only the current floor's manifest is loaded synchronously. Others load asynchronously, prioritized by proximity to current floor.

**Stage 6:** Mouse, keyboard, touch, and scroll event listeners registered. Elevator click handler bound. Keyboard shortcuts activated.

**Stage 7:** If localStorage contains saved state (floor, room, zoom level, scroll position), restore it. If empty or corrupted, start at F01 at L1 zoom. Emit `system.ready`. Loading screen fades out.

---

## 20.5 Performance Budget

### 20.5.1 Runtime Budgets

| Operation              | Budget    | Measurement Point                                      |
|------------------------|-----------|--------------------------------------------------------|
| First paint            | < 500ms   | Navigation start to first meaningful pixel             |
| Interactive            | < 1000ms  | Navigation start to first input response               |
| Floor transition       | < 200ms   | `render.transition.start` to `.end`                    |
| Room zoom              | < 150ms   | `navigation.room.enter` to render stable               |
| Document open          | < 100ms   | `document.open` to content visible                     |
| Building.json parse    | < 50ms    | JSON.parse() of building definition                    |
| Floor manifest parse   | < 20ms    | JSON.parse() of single floor manifest                  |
| SVG floor render       | < 100ms   | SVG injection to paint complete                        |
| Event bus delivery     | < 1ms     | emit() to last synchronous subscriber return           |
| localStorage read      | < 5ms     | getItem + JSON.parse                                   |
| localStorage write     | < 5ms     | JSON.stringify + setItem                               |
| Search query           | < 200ms   | Query submission to results displayed                  |
| Idle CPU               | < 2%      | Background CPU when not interacting                    |

### 20.5.2 Asset Size Budgets (gzip)

| Asset                  | Budget    |
|------------------------|-----------|
| Total JS bundle        | < 100KB   |
| Total CSS              | < 30KB    |
| building.json          | < 50KB    |
| Single floor SVG       | < 80KB    |
| Single floor manifest  | < 10KB    |
| Search index           | < 200KB   |
| Total initial payload  | < 300KB   |

### 20.5.3 Animation Frame Budget (60fps = 16.67ms per frame)

| Phase             | Budget   |
|-------------------|----------|
| JavaScript        | < 6ms    |
| Style             | < 2ms    |
| Layout            | < 2ms    |
| Paint             | < 4ms    |
| Composite         | < 2ms    |
| **Total**         | **< 16ms** |

---

## 20.6 Master Folder Structure

```
hic/
+-- index.html                          # Entry point
+-- building.json                       # Complete building structure
+-- manifest.json                       # Web app manifest (PWA)
+-- service-worker.js                   # Offline cache support
|
+-- src/
|   +-- boot/
|   |   +-- hic.boot.js                # Boot orchestrator
|   +-- events/
|   |   +-- hic.events.js              # Event bus
|   +-- state/
|   |   +-- hic.state.js               # localStorage persistence
|   +-- data/
|   |   +-- hic.data.js                # Spatial data model
|   |   +-- hic.data.schema.js         # JSON schema validation
|   +-- nav/
|   |   +-- hic.nav.js                 # Navigation engine
|   |   +-- hic.nav.elevator.js        # Elevator logic
|   +-- zoom/
|   |   +-- hic.zoom.js                # Zoom management
|   +-- render/
|   |   +-- hic.render.js              # SVG rendering pipeline
|   |   +-- hic.render.floor.js        # Floor rendering
|   |   +-- hic.render.room.js         # Room rendering
|   |   +-- hic.render.effects.js      # Neon glow, particles
|   +-- anim/
|   |   +-- hic.anim.js                # Animation sequencer
|   |   +-- hic.anim.transitions.js    # Transition animations
|   |   +-- hic.anim.ambient.js        # Ambient effects
|   +-- docs/
|   |   +-- hic.docs.js                # Document loading/display
|   |   +-- hic.docs.search.js         # Full-text search
|   |   +-- hic.docs.xref.js           # Cross-reference resolver
|   +-- security/
|   |   +-- hic.security.js            # Clearance checks
|   +-- cache/
|   |   +-- hic.cache.js               # Caching and preloading
|   +-- a11y/
|   |   +-- hic.a11y.js                # Accessibility layer
|   +-- perf/
|   |   +-- hic.perf.js                # Performance monitoring
|   +-- ui/
|       +-- hic.ui.input.js            # Input normalization
|       +-- hic.ui.elevator.js         # Elevator component
|       +-- hic.ui.panel.js            # Document panel
|       +-- hic.ui.overlay.js          # Modal/overlay
|       +-- hic.ui.loading.js          # Boot loading screen
|
+-- styles/
|   +-- hic-core.css                   # Reset, variables, base
|   +-- hic-building.css               # Building exterior
|   +-- hic-floor.css                  # Floor layout
|   +-- hic-room.css                   # Room detail
|   +-- hic-neon.css                   # Neon glow effects
|   +-- hic-elevator.css               # Elevator component
|   +-- hic-document.css               # Document panel
|   +-- hic-transition.css             # Animation keyframes
|   +-- hic-a11y.css                   # Accessibility overrides
|   +-- hic-degraded.css               # Degraded mode styles
|   +-- hic-print.css                  # Print stylesheet
|
+-- floors/                            # One directory per floor
|   +-- B3/ B2/ B1/                    # Basements
|   +-- F01/ through F20/             # Standard floors
|   +-- R1/ R2/                        # Rooftops
|   (each containing: {id}.blueprint.json, {id}.manifest.json,
|    {id}.floor.svg, {id}.anim.json)
|
+-- rooms/                             # One .room.json per room (~96 files)
+-- docs/                              # .md + .meta.json per document
|   +-- search-index.json             # Pre-built search index
|
+-- build/
|   +-- build.py                       # Master build script
|   +-- build-floors.py                # Floor SVG generator
|   +-- build-search-index.py          # Search index generator
|   +-- build-manifest.py              # Floor manifest generator
|   +-- validate.py                    # Schema validation
|   +-- optimize.py                    # Asset optimization
|
+-- test/
|   +-- unit/                          # Unit tests
|   +-- integration/                   # Cross-layer tests
|   +-- performance/                   # Budget tests
|   +-- accessibility/                 # A11y compliance tests
|   +-- visual/                        # Visual regression
|   +-- fixtures/                      # Shared test data
|
+-- dist/                              # Build output (gitignored)
+-- tools/                             # Dev utilities
    +-- dev-server.py
    +-- event-inspector.js
    +-- floor-editor.html
    +-- room-calculator.py
```

---

## 20.7 Build Pipeline

```
Markdown Sources + JSON Definitions + SVG Floor Art + JS/CSS Source
                              |
                              v
+-----------------------------------------------------------------------+
|                       build.py (Master Build)                          |
|                                                                        |
|  1. Validate all JSON against schemas          (validate.py)           |
|  2. Parse all markdown, extract frontmatter                            |
|  3. Build full-text search index               (build-search-index.py) |
|  4. Generate floor manifests from room defs    (build-manifest.py)     |
|  5. Optimize SVG floor files                   (optimize.py)           |
|  6. Bundle JavaScript modules --> hic.min.js                           |
|  7. Bundle CSS --> hic.min.css                                         |
|  8. Generate service-worker.js with cache manifest                     |
|  9. Copy all assets to dist/                                           |
| 10. Post-build validation: verify all references resolve               |
| 11. Report: file counts, sizes, budget compliance                      |
+-----------------------------------------------------------------------+
                              |
                              v
                         dist/ (deployable static site)
```

**Build invariants:**

1. **Deterministic.** Same inputs produce byte-identical outputs.
2. **Offline.** Zero network requests during build.
3. **Fast.** Under 30 seconds for 25 floors and 1000 documents.
4. **Fail loud.** Any validation error aborts the build with a clear message.
5. **Idempotent.** Running twice with no source changes produces no changes in dist/.

---

## 20.8 Deployment Checklist

- [ ] All JSON files pass schema validation
- [ ] All cross-references resolve
- [ ] All room IDs match existing room definition files
- [ ] All floor IDs match existing floor directories
- [ ] Build completes without warnings
- [ ] Unit tests pass (100%)
- [ ] Integration tests pass (100%)
- [ ] Performance tests pass (95%+)
- [ ] Accessibility tests pass (100%)
- [ ] Visual regression shows no unexpected changes
- [ ] Total bundle under 180KB gzip
- [ ] First paint under 500ms on reference hardware
- [ ] Interactive under 1000ms on reference hardware
- [ ] Service worker cache manifest is current
- [ ] Offline navigation works with network disabled
- [ ] Version numbers correctly incremented
- [ ] dist/ contains only build outputs

---

## 20.9 Master Room Registry

Complete registry of all rooms across all floors. Clearance: 0=Public, 1=Internal, 2=Controlled, 3=Restricted, 4=Secret, 5=TopSecret. Capacity = max documents before pagination.

### Basement Floors (B3, B2, B1)

| Room ID                    | Floor | Name                        | Type | Clr | Cap |
|----------------------------|-------|-----------------------------|------|-----|-----|
| B3-VLT-CRYPTO-01           | B3    | Cryptographic Vault Alpha   | VLT  | 5   | 30  |
| B3-VLT-CRYPTO-02           | B3    | Cryptographic Vault Beta    | VLT  | 5   | 30  |
| B3-ARC-COLDSTORE-01        | B3    | Cold Storage Archive        | ARC  | 4   | 200 |
| B3-ARC-DISASTER-01         | B3    | Disaster Recovery Archive   | ARC  | 4   | 100 |
| B3-SRV-BACKUP-01           | B3    | Backup Server Room          | SRV  | 4   | 50  |
| B2-UTL-POWER-01            | B2    | Primary Power Control       | UTL  | 3   | 25  |
| B2-UTL-HVAC-01             | B2    | Environmental Control       | UTL  | 3   | 20  |
| B2-WRK-HARDWARE-01         | B2    | Hardware Workshop           | WRK  | 3   | 30  |
| B2-SRV-COMPUTE-01          | B2    | Primary Compute Node        | SRV  | 3   | 40  |
| B2-MON-ENVIRON-01          | B2    | Environmental Monitoring    | MON  | 2   | 15  |
| B1-SEC-ACCESS-01           | B1    | Access Control Center       | SEC  | 4   | 35  |
| B1-SEC-AUDIT-01            | B1    | Audit Operations            | SEC  | 4   | 50  |
| B1-SEC-INTRUSION-01        | B1    | Intrusion Detection Center  | SEC  | 5   | 25  |
| B1-QUA-MEDIA-01            | B1    | Media Quarantine Chamber    | QUA  | 4   | 20  |
| B1-MON-PERIMETER-01        | B1    | Perimeter Monitoring        | MON  | 3   | 15  |

### Lower Floors (F01-F05)

| Room ID                    | Floor | Name                        | Type | Clr | Cap |
|----------------------------|-------|-----------------------------|------|-----|-----|
| F01-LIB-WELCOME-01         | F01   | Welcome Library             | LIB  | 0   | 50  |
| F01-BRF-ORIENTATION-01     | F01   | Visitor Orientation         | BRF  | 0   | 20  |
| F01-LIB-PUBLIC-01          | F01   | Public Documents Library    | LIB  | 0   | 100 |
| F01-HUB-DIRECTORY-01       | F01   | Building Directory Hub      | HUB  | 0   | 25  |
| F02-OPS-GOVERNANCE-01      | F02   | Governance Operations       | OPS  | 1   | 40  |
| F02-ARC-CHARTER-01         | F02   | Charter Archive             | ARC  | 1   | 30  |
| F02-BRF-COUNCIL-01         | F02   | Council Briefing Room       | BRF  | 2   | 20  |
| F02-VLT-CONSTITUTIONAL-01  | F02   | Constitutional Vault        | VLT  | 3   | 15  |
| F03-LIB-ETHICS-01          | F03   | Ethics Library              | LIB  | 1   | 60  |
| F03-LIB-PHILOSOPHY-01      | F03   | Philosophy Library          | LIB  | 1   | 60  |
| F03-BRF-DEBATE-01          | F03   | Ethics Debate Hall          | BRF  | 1   | 20  |
| F03-ARC-VALUES-01          | F03   | Value Statements Archive    | ARC  | 1   | 40  |
| F04-LIB-INFRA-01           | F04   | Infrastructure Library      | LIB  | 2   | 60  |
| F04-WRK-POWER-01           | F04   | Power Systems Workshop      | WRK  | 2   | 30  |
| F04-WRK-NETWORK-01         | F04   | Network Systems Workshop    | WRK  | 2   | 30  |
| F04-OPS-MAINTENANCE-01     | F04   | Maintenance Operations      | OPS  | 2   | 25  |
| F05-LAB-SCHEMA-01          | F05   | Schema Design Lab           | LAB  | 2   | 40  |
| F05-ARC-DATAMODEL-01       | F05   | Data Model Archive          | ARC  | 2   | 50  |
| F05-OPS-MIGRATION-01       | F05   | Data Migration Operations   | OPS  | 2   | 30  |
| F05-LIB-FORMAT-01          | F05   | Format Specification Library| LIB  | 2   | 35  |

### Mid Floors (F06-F10)

| Room ID                    | Floor | Name                        | Type | Clr | Cap |
|----------------------------|-------|-----------------------------|------|-----|-----|
| F06-LAB-AUTOMATION-01      | F06   | Automation Lab              | LAB  | 2   | 40  |
| F06-OPS-BUILD-01           | F06   | Build Pipeline Operations   | OPS  | 2   | 35  |
| F06-OPS-CICD-01            | F06   | CI/CD Operations            | OPS  | 2   | 30  |
| F06-LIB-SCRIPTS-01         | F06   | Script Library              | LIB  | 2   | 50  |
| F07-LAB-SIGINT-01          | F07   | SIGINT Analysis Lab         | LAB  | 4   | 35  |
| F07-LAB-SIGINT-02          | F07   | SIGINT Processing Lab       | LAB  | 4   | 35  |
| F07-LAB-OSINT-01           | F07   | OSINT Collection Lab        | LAB  | 3   | 40  |
| F07-OPS-ANALYSIS-01        | F07   | Analysis Operations         | OPS  | 3   | 30  |
| F07-VLT-CLASSIFIED-01      | F07   | Classified Intel Vault      | VLT  | 5   | 20  |
| F08-OPS-INDEX-01           | F08   | Indexing Operations         | OPS  | 2   | 40  |
| F08-LAB-SEARCH-01          | F08   | Search Algorithm Lab        | LAB  | 2   | 30  |
| F08-LIB-TAXONOMY-01        | F08   | Taxonomy Library            | LIB  | 1   | 50  |
| F08-OPS-XREF-01            | F08   | Cross-Reference Operations  | OPS  | 2   | 35  |
| F09-BRF-TRAINING-01        | F09   | Training Briefing Room      | BRF  | 1   | 25  |
| F09-BRF-TRAINING-02        | F09   | Advanced Training Room      | BRF  | 2   | 25  |
| F09-LIB-CURRICULUM-01      | F09   | Curriculum Library          | LIB  | 1   | 50  |
| F09-LAB-PRACTICE-01        | F09   | Practice Lab                | LAB  | 1   | 30  |
| F10-COM-INTERNAL-01        | F10   | Internal Communications     | COM  | 2   | 30  |
| F10-COM-EXTERNAL-01        | F10   | External Relations          | COM  | 2   | 25  |
| F10-COM-BROADCAST-01       | F10   | Broadcast Control           | COM  | 3   | 20  |
| F10-ARC-MESSAGES-01        | F10   | Message Archive             | ARC  | 2   | 60  |

### Upper Floors (F11-F16)

| Room ID                    | Floor | Name                        | Type | Clr | Cap |
|----------------------------|-------|-----------------------------|------|-----|-----|
| F11-LAB-RESEARCH-01        | F11   | Primary Research Lab        | LAB  | 3   | 40  |
| F11-LAB-PROTOTYPE-01       | F11   | Prototype Lab               | LAB  | 3   | 30  |
| F11-LAB-EXPERIMENT-01      | F11   | Experimental Systems Lab    | LAB  | 3   | 35  |
| F11-ARC-FINDINGS-01        | F11   | Research Findings Archive   | ARC  | 2   | 50  |
| F12-HUB-FEDERATION-01      | F12   | Federation Hub              | HUB  | 3   | 30  |
| F12-OPS-IMPORT-01          | F12   | Import Operations           | OPS  | 3   | 25  |
| F12-OPS-EXPORT-01          | F12   | Export Operations           | OPS  | 3   | 25  |
| F12-COM-PEERING-01         | F12   | Peering Communications      | COM  | 3   | 20  |
| F13-OPS-PLATFORM-01        | F13   | Platform Operations         | OPS  | 2   | 35  |
| F13-SRV-CONTAINER-01       | F13   | Container Server Room       | SRV  | 3   | 30  |
| F13-OPS-DEPLOY-01          | F13   | Deployment Operations       | OPS  | 2   | 25  |
| F13-LIB-PLATFORM-01        | F13   | Platform Library            | LIB  | 2   | 40  |
| F14-OPS-QA-01              | F14   | QA Operations               | OPS  | 2   | 40  |
| F14-LAB-TESTING-01         | F14   | Testing Lab                 | LAB  | 2   | 35  |
| F14-ARC-AUDIT-01           | F14   | Audit Trail Archive         | ARC  | 3   | 60  |
| F14-OPS-COMPLIANCE-01      | F14   | Compliance Operations       | OPS  | 3   | 30  |
| F15-LAB-INTERFACE-01       | F15   | Interface Design Lab        | LAB  | 2   | 35  |
| F15-LIB-DESIGN-01          | F15   | Design System Library       | LIB  | 1   | 50  |
| F15-LAB-A11Y-01            | F15   | Accessibility Lab           | LAB  | 2   | 30  |
| F15-BRF-REVIEW-01          | F15   | Design Review Room          | BRF  | 2   | 20  |
| F16-ARC-CHANGELOG-01       | F16   | Changelog Archive           | ARC  | 1   | 80  |
| F16-OPS-VERSION-01         | F16   | Version Control Operations  | OPS  | 2   | 30  |
| F16-MON-DECAY-01           | F16   | Decay Tracking Monitor      | MON  | 2   | 20  |
| F16-LIB-MEMORY-01          | F16   | Institutional Memory Library| LIB  | 1   | 60  |

### High Floors (F17-F20)

| Room ID                    | Floor | Name                        | Type | Clr | Cap |
|----------------------------|-------|-----------------------------|------|-----|-----|
| F17-VLT-ZEROKEY-01         | F17   | Zero-Knowledge Vault        | VLT  | 5   | 15  |
| F17-LAB-CRYPTO-01          | F17   | Cryptographic Research Lab  | LAB  | 4   | 30  |
| F17-OPS-REDTEAM-01         | F17   | Red Team Operations         | OPS  | 5   | 25  |
| F17-SEC-ZEROTRUST-01       | F17   | Zero Trust Architecture     | SEC  | 4   | 35  |
| F18-LAB-DISTRIB-01         | F18   | Distributed Systems Lab     | LAB  | 3   | 35  |
| F18-OPS-REPLICATION-01     | F18   | Replication Operations      | OPS  | 3   | 30  |
| F18-OPS-CONFLICT-01        | F18   | Conflict Resolution Ops     | OPS  | 3   | 25  |
| F18-SRV-STORAGE-01         | F18   | Distributed Storage Server  | SRV  | 3   | 40  |
| F19-OPS-INTEGRATION-01     | F19   | Integration Operations      | OPS  | 3   | 35  |
| F19-LAB-EVENTBUS-01        | F19   | Event Bus Lab               | LAB  | 3   | 25  |
| F19-OPS-TESTING-01         | F19   | Integration Testing Ops     | OPS  | 3   | 30  |
| F19-MON-COMPAT-01          | F19   | Compatibility Monitor       | MON  | 2   | 20  |
| F20-CMD-OVERVIEW-01        | F20   | Chief Architect Overview    | CMD  | 4   | 50  |
| F20-CMD-DASHBOARD-01       | F20   | Master Dashboard            | CMD  | 3   | 25  |
| F20-OPS-SYNTHESIS-01       | F20   | System Synthesis Operations | OPS  | 4   | 40  |
| F20-BRF-STRATEGY-01        | F20   | Strategic Briefing Room     | BRF  | 4   | 20  |

### Rooftop Floors (R1, R2)

| Room ID                    | Floor | Name                        | Type | Clr | Cap |
|----------------------------|-------|-----------------------------|------|-----|-----|
| R1-MON-HEALTH-01           | R1    | System Health Monitor       | MON  | 2   | 20  |
| R1-MON-TELEMETRY-01        | R1    | Real-Time Telemetry         | MON  | 2   | 15  |
| R1-MON-PANORAMIC-01        | R1    | Panoramic Status Display    | MON  | 1   | 10  |
| R2-COM-BROADCAST-01        | R2    | Broadcast Transmitter       | COM  | 3   | 15  |
| R2-COM-RSS-01              | R2    | RSS Feed Control            | COM  | 2   | 10  |
| R2-COM-API-01              | R2    | Public API Gateway          | COM  | 3   | 20  |
| R2-COM-FEDERATION-01       | R2    | Federation Broadcast        | COM  | 3   | 15  |

### Registry Totals

| Metric                             | Count  |
|------------------------------------|--------|
| Total floors                       | 25     |
| Total rooms                        | 96     |
| Total document capacity            | 3,110  |
| Clearance 0 (Public) rooms         | 4      |
| Clearance 1 (Internal) rooms       | 10     |
| Clearance 2 (Controlled) rooms     | 30     |
| Clearance 3 (Restricted) rooms     | 28     |
| Clearance 4 (Secret) rooms         | 15     |
| Clearance 5 (Top Secret) rooms     | 9      |
| Type LAB rooms                     | 19     |
| Type OPS rooms                     | 22     |
| Type LIB rooms                     | 14     |
| Type ARC rooms                     | 11     |
| Type VLT rooms                     | 6      |
| Type COM rooms                     | 7      |
| Type MON rooms                     | 7      |
| Type BRF rooms                     | 9      |
| Type SRV rooms                     | 4      |
| Type SEC rooms                     | 3      |
| Type WRK rooms                     | 3      |
| Type HUB rooms                     | 3      |
| Type UTL rooms                     | 2      |
| Type CMD rooms                     | 2      |
| Type QUA rooms                     | 1      |

---

## 20.10 System Invariants

These rules must never be violated. If an implementation violates any invariant, it is a bug, even if it faithfully implements a subsystem specification.

**Invariant 1: Every room belongs to exactly one floor.** A room ID contains its floor ID. A room cannot span floors or be moved without a new ID.

**Invariant 2: Every document belongs to exactly one room.** Cross-references create links, not copies. A document's metadata specifies its single canonical room.

**Invariant 3: Clearance is monotonically non-decreasing from container to content.** Floor clearance <= room clearance <= document clearance. Always.

**Invariant 4: The building exterior is always accessible.** Zoom levels L0 and L1 never require clearance. Only floor entry (L2+) triggers clearance checks.

**Invariant 5: Navigation state is always recoverable.** Current state serializes to < 1KB JSON in localStorage. Reload restores the same visual position. Corruption falls back to F01/L1.

**Invariant 6: The event bus never silently drops events.** Registered subscribers always receive published events. Subscriber exceptions do not prevent delivery to other subscribers. Delivery failures emit `system.error`.

**Invariant 7: Every file matches its naming convention.** No exceptions, no legacy names, no temp files. The build pipeline rejects non-conforming files.

**Invariant 8: The build is deterministic and offline.** Same inputs = byte-identical outputs. Zero network requests.

**Invariant 9: Degradation never produces incorrect information.** Wrong data is always worse than no data. A degraded view shows placeholders, never wrong content.

**Invariant 10: Document IDs are permanent.** Once assigned, an ID is never reassigned to a different document. Deleted documents retire their IDs.

**Invariant 11: Room type codes are immutable.** The 15 codes in Section 19.5.2 are fixed. LAB always means Laboratory. VLT always means Vault.

**Invariant 12: Floor ordering is stable.** B3 is always at the bottom. R2 is always at the top. Floors never swap positions.

---

## 20.11 Conflict Resolution Rules

When two HIC specifications disagree, apply these rules in order. The first matching rule resolves the conflict.

**Priority 1: System invariants win over everything.** Including this document. If this document contradicts an invariant, the invariant governs.

**Priority 2: This Master Blueprint wins over all other HIC documents.**

**Priority 3: Security specifications win over non-security specifications.** Security is never negotiable for aesthetics.

**Priority 4: Accessibility specifications win over visual specifications.** The building must be usable by all operators.

**Priority 5: Performance specifications win over animation specifications.** Fast before beautiful.

**Priority 6: Data model specifications win over visual layout specifications.** Layouts adapt to data, not the reverse.

**Priority 7: More specific specifications win over more general ones.** Unless a higher-priority rule is violated.

**Priority 8: When no rule resolves the conflict, Agent 20 decides.** Document the resolution here.

---

## 20.12 Phase 2 Preparation: Implementation Order

Phase 1 is specification. Phase 2 is implementation. Phase 2 begins only after all quality gates (Section 20.13) pass.

### P2-1: Foundation (Weeks 1-2)

- Create `hic/` directory structure per Section 20.6
- Implement `hic.events.js` (event bus with pub/sub, wildcards, priority, cancellation)
- Implement `hic.boot.js` (boot orchestrator emitting stage events)
- Implement `hic.state.js` (localStorage persistence with version checking)
- Unit tests: 100% branch coverage on event bus
- **Exit:** Event bus delivers 10,000 events in < 100ms. Boot emits all 8 stages correctly.

### P2-2: Data Layer (Weeks 3-4)

- Define JSON schemas for building.json, floor blueprints, manifests, room definitions
- Implement `hic.data.js` and `hic.data.schema.js`
- Create `building.json` with all 25 floors
- Create all 25 floor manifest files and all 96 room definition files
- Cross-layer integration tests: event bus + data layer
- **Exit:** All JSON validates. Room lookup < 1ms. Floor manifest parse < 20ms.

### P2-3: Rendering Core (Weeks 5-7)

- Implement `hic.render.js`, `hic.render.floor.js`, `hic.render.room.js`
- Create building exterior SVG (L0/L1)
- Create all 25 floor SVGs (F01 first as reference)
- Implement core CSS: `hic-core.css`, `hic-building.css`, `hic-floor.css`, `hic-room.css`
- **Exit:** All floors render without errors. No SVG exceeds 80KB.

### P2-4: Visual Identity (Weeks 8-9)

- Implement `hic-neon.css` and `hic.render.effects.js`
- Apply neon aesthetic to all floor SVGs
- Implement `hic-transition.css`
- Establish visual regression baselines
- **Exit:** Unmistakably cyberpunk. No frame drops on neon effects.

### P2-5: Navigation (Weeks 10-11)

- Implement `hic.nav.js`, `hic.zoom.js`, `hic.nav.elevator.js`
- Implement `hic.ui.elevator.js`, `hic-elevator.css`, `hic.ui.input.js`
- URL hash routing for deep links and back button
- **Exit:** Every room reachable. Keyboard navigation works. All transitions within budget.

### P2-6: Knowledge Layer (Weeks 12-13)

- Port all documents to `docs/` with correct IDs
- Create `.meta.json` files for all documents
- Implement `hic.docs.js`, `hic.docs.search.js`, `hic.docs.xref.js`
- Build search index
- **Exit:** Every document accessible from its room. Search returns correct results.

### P2-7: Security Layer (Week 14)

- Implement `hic.security.js` (subscribes at priority 10)
- Wire clearance events into navigation sequence
- Implement denial animations
- **Exit:** No room accessible below its clearance level. Denial animations work.

### P2-8: Animation & Polish (Weeks 15-16)

- Implement `hic.anim.js`, `hic.anim.transitions.js`, `hic.anim.ambient.js`
- Performance optimization pass
- Implement `hic.cache.js`
- **Exit:** All animations at 60fps. Cache hit rate > 90% for sequential navigation.

### P2-9: Accessibility (Week 17)

- Implement `hic.a11y.js` and `hic-a11y.css`
- Screen reader testing (VoiceOver, NVDA)
- Keyboard-only navigation testing
- WCAG 2.1 AA audit
- **Exit:** No WCAG AA violations. Every room reachable by keyboard alone.

### P2-10: Build & Deploy (Week 18)

- Implement `build.py` and all sub-scripts
- Implement `service-worker.js`
- Execute full deployment checklist
- Full system test: boot, navigate all floors, enter rooms, open documents, test search, test security, test degradation
- **Exit:** All checklist items pass. HIC works fully offline.

---

## 20.13 Quality Gates

Phase 2 does not begin until every gate passes.

### Gate 1: Specification Completeness

- [ ] All 20 agent specifications exist as ratified documents
- [ ] Every floor has defined purpose, room list, and clearance level
- [ ] Every room has defined type, name, clearance, and capacity
- [ ] Every event has a fully defined payload schema
- [ ] Every naming convention has at least three examples
- [ ] Every error code range is allocated with at least one example per category
- [ ] Master room registry accounts for all rooms in any specification

### Gate 2: Specification Consistency

- [ ] No unresolved contradictions between specifications
- [ ] All 96 room IDs conform to naming convention (Section 19.5.2)
- [ ] All document IDs conform to naming convention (Section 19.5.3)
- [ ] All 25 floor IDs are consistent across all documents
- [ ] Event sequences match event payload definitions
- [ ] Clearance levels obey Invariant 3 (monotonically non-decreasing)

### Gate 3: Technical Feasibility

- [ ] Performance budgets validated on reference hardware with a prototype
- [ ] Search index size projected within 200KB budget
- [ ] SVG neon aesthetic demonstrated within 80KB per-floor budget
- [ ] Event bus benchmarked at < 1ms delivery with 50+ subscribers
- [ ] building.json prototype parses in under 50ms

### Gate 4: Toolchain Readiness

- [ ] Python 3.8+ available
- [ ] Node.js 18+ available
- [ ] Local HTTP server can serve dist/
- [ ] Git configured for version control
- [ ] build.py skeleton exists and executes

### Gate 5: Team Readiness

- [ ] At least one person can explain boot sequence, event bus, and naming from memory
- [ ] At least one person can reproduce the master data flow diagram from memory
- [ ] Conflict resolution rules are understood and accepted by all contributors
- [ ] System invariants reviewed and accepted by all contributors

---

## 20.14 Quick Reference Card

```
HOLM INTELLIGENCE COMPLEX -- QUICK REFERENCE
=============================================

FLOORS:       B3 B2 B1 | F01-F20 | R1 R2         (25 total)
ROOMS:        {floor}-{type}-{name}-{number}       (96 total)
DOCUMENTS:    D{domain}-{number}-{slug}
ROOM TYPES:   LAB ARC OPS CMD VLT LIB WRK COM SEC SRV UTL BRF QUA HUB MON
CLEARANCE:    0=Public 1=Internal 2=Controlled 3=Restricted 4=Secret 5=TopSecret
ZOOM:         L0=City L1=Building L2=Floor L3=Room L4=Document
CSS:          hic-{component}-{modifier}
JS:           hic.{subsystem}.{function}
EVENTS:       navigation.* document.* security.* render.* system.*
BOOT:         HTML>Render>EventBus>Data>Exterior>Floors>Input>State>READY
DEGRADE:      D0=Full D1=NoAnim D2=NoSVG D3=Text D4=Offline D5=Error
PERF:         Paint<500ms TTI<1000ms Floor<200ms Room<150ms Doc<100ms
CONFLICT:     Invariants>Blueprint>Security>A11y>Perf>Data>Specific>General
```

---

## 20.15 Glossary

| Term                  | Definition                                                                    |
|-----------------------|-------------------------------------------------------------------------------|
| HIC                   | Holm Intelligence Complex. The cyberpunk skyscraper documentation interface.  |
| Floor                 | A horizontal building layer corresponding to a documentation domain.          |
| Room                  | A spatial unit within a floor containing related documents.                   |
| Clearance             | Integer (0-5) representing required security authorization.                   |
| Event Bus             | Publish/subscribe messaging system connecting all HIC layers.                 |
| Event Envelope        | Standard JSON wrapper: name, sequence, timestamp, source, payload.            |
| Zoom Level            | Magnification state: L0 city, L1 building, L2 floor, L3 room, L4 document.  |
| Neon Aesthetic        | Cyberpunk visual language: glowing lines, dark backgrounds, electric color.   |
| Degradation Level     | Operational state (D0-D5) after component failure.                            |
| Spatial Data Model    | JSON structure defining building floors, rooms, and document placement.        |
| Boot Sequence         | Eight-stage initialization from HTML to interactive state.                    |
| Invariant             | Rule that must never be violated regardless of subsystem specifications.       |
| SemVer                | Semantic Versioning (MAJOR.MINOR.PATCH) for all definition files.             |
| Layer                 | One of ten functional divisions of the HIC system.                            |
| Cross-reference       | Navigable pointer to a document in another room.                              |

---

## 20.16 Revision History

| Version | Date       | Author                     | Changes                                              |
|---------|------------|----------------------------|------------------------------------------------------|
| 1.0.0   | 2026-02-17 | Agent 20 (Chief Architect) | Initial ratification                                 |
| 2.0.0   | 2026-02-17 | Agent 20 (Chief Architect) | Expanded building manifest and integration points    |
| 3.0.0   | 2026-02-17 | Agent 20 (Chief Architect) | Complete rewrite: Agents 19-20 full specification. Added event bus, event catalog, naming conventions, cross-layer protocol, boot sequence, performance budgets, master room registry, system invariants, conflict resolution, Phase 2 implementation plan, quality gates. |

---

*This document is the master reference for the Holm Intelligence Complex. It is the first document you read and the last document you consult before making any architectural decision. When specifications conflict, apply the rules here. When extending the system, verify against the invariants here. When in doubt, return here. The building stands on these foundations.*

**END OF MASTER BLUEPRINT**
