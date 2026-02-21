# STAGE 4: HOLM INTELLIGENCE COMPLEX -- INTERACTION MECHANICS

## Click Logic, Hover Logic, Zoom/Pan, and Navigation State

---

> **Agents Covered:** 10 (Click Logic), 11 (Hover Logic), 12 (Zoom/Pan and Navigation State)
>
> **System Context:** The Holm Intelligence Complex (HIC) is a cyberpunk neon skyscraper rendered as an interactive SVG/Canvas interface. The building represents the sovereign intranet's documentation structure. Floors correspond to documentation domains, rooms correspond to individual document collections or subsystems, and objects within rooms correspond to individual documents, links, or actionable endpoints. Every interaction described in this specification maps a physical-spatial metaphor to a navigation action within the intranet.
>
> **Implementation Target:** Browser-based (ES2022+), SVG primary with Canvas fallback, touch-enabled, fully keyboard-accessible.

---

## TABLE OF CONTENTS

1. [Agent 10: Click Logic](#agent-10-click-logic)
   - [10.1 Click Detection Pipeline](#101-click-detection-pipeline)
   - [10.2 Building View Click](#102-building-view-click)
   - [10.3 Floor View Click](#103-floor-view-click)
   - [10.4 Room View Click](#104-room-view-click)
   - [10.5 Double-Click Behavior](#105-double-click-behavior)
   - [10.6 Right-Click Context Menu](#106-right-click-context-menu)
   - [10.7 Click-Outside Deselection](#107-click-outside-deselection)
   - [10.8 Click State Machine](#108-click-state-machine)
   - [10.9 Touch Equivalents](#109-touch-equivalents)
2. [Agent 11: Hover Logic](#agent-11-hover-logic)
   - [11.1 Hover Detection Pipeline](#111-hover-detection-pipeline)
   - [11.2 Building View Hover](#112-building-view-hover)
   - [11.3 Floor View Hover](#113-floor-view-hover)
   - [11.4 Room View Hover](#114-room-view-hover)
   - [11.5 Tooltip Timing and Lifecycle](#115-tooltip-timing-and-lifecycle)
   - [11.6 CSS/SVG Hover State Definitions](#116-csssvg-hover-state-definitions)
   - [11.7 Performance Constraints](#117-performance-constraints)
   - [11.8 Mobile Long-Press Hover Equivalent](#118-mobile-long-press-hover-equivalent)
3. [Agent 12: Zoom/Pan and Navigation State](#agent-12-zoompan-and-navigation-state)
   - [12.1 Zoom Level Definitions](#121-zoom-level-definitions)
   - [12.2 Zoom Transition Animation](#122-zoom-transition-animation)
   - [12.3 Pan Mechanics](#123-pan-mechanics)
   - [12.4 Scroll Wheel Zoom](#124-scroll-wheel-zoom)
   - [12.5 Pinch Zoom on Touch Devices](#125-pinch-zoom-on-touch-devices)
   - [12.6 Navigation State Object](#126-navigation-state-object)
   - [12.7 URL Hash Encoding](#127-url-hash-encoding)
   - [12.8 Browser History Integration](#128-browser-history-integration)
   - [12.9 Keyboard Shortcuts](#129-keyboard-shortcuts)
   - [12.10 Breadcrumb Trail](#1210-breadcrumb-trail)

---

## AGENT 10: CLICK LOGIC

### 10.1 Click Detection Pipeline

Every click event in the HIC passes through a unified detection pipeline before being dispatched to view-specific handlers. The pipeline ensures consistent behavior across input modes and prevents race conditions between click, double-click, and long-press events.

**Pipeline Stages:**

```
STAGE 1: Raw Event Capture
  - Capture mousedown/mouseup (pointer events preferred)
  - Record timestamp, clientX, clientY, button (0=left, 1=middle, 2=right)
  - Convert clientX/clientY to SVG coordinate space via getScreenCTM().inverse()

STAGE 2: Click Classification
  - If time between mousedown and mouseup > 300ms  --> DRAG (abort click)
  - If distance between mousedown and mouseup > 4px --> DRAG (abort click)
  - If second click within 250ms of first click     --> DOUBLE-CLICK
  - If button === 2                                  --> RIGHT-CLICK
  - Otherwise                                        --> SINGLE-CLICK

STAGE 3: Hit Testing
  - Determine current view level (building / floor / room)
  - Run view-specific hit test against SVG coordinate
  - Return hit target or null

STAGE 4: Action Dispatch
  - Map (view_level, click_type, hit_target) to action
  - Execute action through navigation state controller
  - Emit click telemetry event
```

**Coordinate Transformation Function:**

```javascript
function clientToSVG(event, svgElement) {
  const pt = svgElement.createSVGPoint();
  pt.x = event.clientX;
  pt.y = event.clientY;
  const ctm = svgElement.getScreenCTM().inverse();
  return pt.matrixTransform(ctm);
}
```

**Click Classifier Implementation:**

```javascript
const CLICK_CONFIG = {
  MAX_CLICK_DURATION_MS: 300,
  MAX_CLICK_DRIFT_PX: 4,
  DOUBLE_CLICK_WINDOW_MS: 250,
  SINGLE_CLICK_DELAY_MS: 260  // slightly longer than double-click window
};

class ClickClassifier {
  constructor() {
    this._pendingSingle = null;
    this._lastClickTime = 0;
    this._lastClickPos = { x: 0, y: 0 };
  }

  classify(downEvent, upEvent) {
    const duration = upEvent.timestamp - downEvent.timestamp;
    const dx = upEvent.x - downEvent.x;
    const dy = upEvent.y - downEvent.y;
    const drift = Math.sqrt(dx * dx + dy * dy);

    if (duration > CLICK_CONFIG.MAX_CLICK_DURATION_MS) return 'drag';
    if (drift > CLICK_CONFIG.MAX_CLICK_DRIFT_PX) return 'drag';
    if (upEvent.button === 2) return 'right-click';

    const now = performance.now();
    const timeSinceLast = now - this._lastClickTime;

    if (timeSinceLast < CLICK_CONFIG.DOUBLE_CLICK_WINDOW_MS) {
      this._lastClickTime = 0;
      clearTimeout(this._pendingSingle);
      return 'double-click';
    }

    this._lastClickTime = now;
    this._lastClickPos = { x: upEvent.x, y: upEvent.y };
    return 'single-click-pending'; // resolved after SINGLE_CLICK_DELAY_MS
  }
}
```

The `single-click-pending` return value triggers a delayed emission. If no second click arrives within `DOUBLE_CLICK_WINDOW_MS`, the pending click is confirmed and dispatched as a true single click. This delay is imperceptible to users (260ms) but critical for disambiguating single from double clicks.

---

### 10.2 Building View Click

In Building View, the entire HIC skyscraper is visible. Floors are rendered as horizontal strips stacked vertically. Each floor strip has a defined Y-coordinate range.

**Hit Testing -- Y-Coordinate Floor Identification:**

The building is rendered with the following dimensional constants:

```javascript
const BUILDING_GEOMETRY = {
  TOTAL_FLOORS: 20,
  BUILDING_TOP_Y: 50,        // SVG units from top of viewport
  BUILDING_BOTTOM_Y: 950,    // SVG units
  FLOOR_HEIGHT: 45,           // (950 - 50) / 20 = 45 SVG units per floor
  BUILDING_LEFT_X: 200,
  BUILDING_RIGHT_X: 800,
  FLOOR_GAP: 2               // 2px visual gap between floors (neon glow line)
};

function hitTestFloor(svgPoint) {
  const { BUILDING_TOP_Y, BUILDING_BOTTOM_Y, BUILDING_LEFT_X, BUILDING_RIGHT_X,
          FLOOR_HEIGHT, TOTAL_FLOORS, FLOOR_GAP } = BUILDING_GEOMETRY;

  // Check horizontal bounds
  if (svgPoint.x < BUILDING_LEFT_X || svgPoint.x > BUILDING_RIGHT_X) {
    return null; // click outside building horizontal bounds
  }

  // Check vertical bounds
  if (svgPoint.y < BUILDING_TOP_Y || svgPoint.y > BUILDING_BOTTOM_Y) {
    return null; // click above or below building
  }

  // Calculate floor index (0 = ground floor, counted from bottom)
  const relativeY = BUILDING_BOTTOM_Y - svgPoint.y;
  const floorIndex = Math.floor(relativeY / FLOOR_HEIGHT);

  // Check if click is in the gap between floors
  const positionInFloor = relativeY % FLOOR_HEIGHT;
  if (positionInFloor < FLOOR_GAP) {
    return null; // click landed on the neon separator line, not a floor
  }

  if (floorIndex < 0 || floorIndex >= TOTAL_FLOORS) {
    return null;
  }

  return {
    type: 'floor',
    floorIndex: floorIndex,
    floorId: `floor-${floorIndex}`,
    floorName: FLOOR_REGISTRY[floorIndex].name
  };
}
```

**Action on Building View Click:**

When a floor is identified by hit testing, the system executes a zoom transition to Floor View:

1. Flash the clicked floor strip with an intensified neon pulse (100ms, opacity 0.4 to 1.0).
2. Begin zoom animation (300ms ease-in-out) centering on the clicked floor.
3. Update navigation state: `current_view: 'floor'`, `floor_id: floorIndex`.
4. Push new state to browser history.
5. Update breadcrumb trail to `Building > Floor N`.

---

### 10.3 Floor View Click

In Floor View, a single floor is displayed as a horizontal cross-section. Rooms are rendered as polygonal regions (rectangles, L-shapes, or irregular polygons defined by vertex arrays).

**Hit Testing -- Point-in-Polygon:**

Each room is defined by an ordered array of vertices forming a closed polygon:

```javascript
// Room definition example
const room = {
  id: 'room-3-07',
  name: 'Threat Intelligence Archive',
  floorId: 'floor-3',
  vertices: [
    { x: 120, y: 80 },
    { x: 280, y: 80 },
    { x: 280, y: 200 },
    { x: 200, y: 200 },
    { x: 200, y: 160 },
    { x: 120, y: 160 }
  ],
  status: 'active',       // active | locked | maintenance
  glowColor: '#00FFCC',
  documents: [...]
};

function pointInPolygon(point, vertices) {
  let inside = false;
  const n = vertices.length;
  for (let i = 0, j = n - 1; i < n; j = i++) {
    const xi = vertices[i].x, yi = vertices[i].y;
    const xj = vertices[j].x, yj = vertices[j].y;

    const intersect = ((yi > point.y) !== (yj > point.y))
      && (point.x < (xj - xi) * (point.y - yi) / (yj - yi) + xi);

    if (intersect) inside = !inside;
  }
  return inside;
}

function hitTestRoom(svgPoint, floorId) {
  const rooms = ROOM_REGISTRY[floorId];
  // Test in reverse render order (topmost drawn last = tested first)
  for (let i = rooms.length - 1; i >= 0; i--) {
    if (pointInPolygon(svgPoint, rooms[i].vertices)) {
      return rooms[i];
    }
  }
  return null;
}
```

**Action on Floor View Click:**

When a room is identified:

1. Apply highlight effect to room polygon (glow intensity increase, border pulse).
2. Slide in the Room Info Panel from the right edge (250ms ease-out):
   - Room name (large, neon-colored heading).
   - Room status indicator (green pulse = active, amber = locked, red = maintenance).
   - Document count.
   - List of top-level document titles (scrollable).
   - "Enter Room" button.
3. Update navigation state: `room_id: room.id`, maintain `current_view: 'floor'`.
4. Clicking "Enter Room" button or double-clicking the room transitions to Room View.

---

### 10.4 Room View Click

In Room View, the interior of a single room is displayed. Interactive elements are rendered as distinct objects: document icons, terminal screens, data nodes, sub-system portals.

**Hit Testing -- Object Identification:**

Objects in Room View are rendered as SVG groups (`<g>`) with bounding boxes. Hit testing uses `document.elementFromPoint()` combined with data attributes:

```javascript
function hitTestRoomObject(svgPoint, roomId) {
  const objects = OBJECT_REGISTRY[roomId];
  for (const obj of objects) {
    const bbox = obj.svgElement.getBBox();
    if (svgPoint.x >= bbox.x && svgPoint.x <= bbox.x + bbox.width &&
        svgPoint.y >= bbox.y && svgPoint.y <= bbox.y + bbox.height) {
      return obj;
    }
  }
  return null;
}
```

**Action on Room View Click:**

Depending on object type:

| Object Type       | Click Action                                                    |
|--------------------|-----------------------------------------------------------------|
| `document`         | Open document viewer in slide-over panel (right, 60% width)    |
| `sub-system`       | Navigate to linked subsystem (full view transition)             |
| `terminal`         | Open interactive terminal overlay (centered modal, 80% x 70%)  |
| `data-node`        | Expand data node to show connections and metadata               |
| `external-link`    | Confirm dialog, then open external URL in new tab               |
| `locked-element`   | Display access-denied overlay with neon red flash               |

---

### 10.5 Double-Click Behavior

Double-click universally means "zoom in one level" regardless of current view:

| Current View   | Double-Click Target   | Result                                    |
|----------------|-----------------------|-------------------------------------------|
| Building       | Floor strip           | Zoom to Floor View (same as single click) |
| Building       | Empty space           | No action                                 |
| Floor          | Room polygon          | Zoom to Room View (skip panel step)       |
| Floor          | Empty space           | No action                                 |
| Room           | Object                | Open object (same as single click)        |
| Room           | Empty space           | No action                                 |

The double-click zoom animation uses a slightly faster curve than the standard transition: 200ms cubic-bezier(0.25, 0.1, 0.25, 1.0).

---

### 10.6 Right-Click Context Menu

Right-click (or two-finger tap on trackpad, long-press on touch) opens a context menu. The native browser context menu is suppressed via `event.preventDefault()` on the `contextmenu` event.

**Context Menu Structure:**

```
+----------------------------------+
|  [icon] Floor 7: Intel Ops       |  <-- header: target name
|----------------------------------|
|  > View Info                     |  <-- opens info panel
|  > Navigate Here                 |  <-- zooms to target
|  > Bookmark                      |  <-- adds to user bookmarks
|  > Copy Link                     |  <-- copies hash URL to clipboard
|  > Open in New Tab               |  <-- opens hash URL in new tab
|----------------------------------|
|  > Back to Building View         |  <-- resets to root view
+----------------------------------+
```

**Context Menu Configuration:**

```javascript
const CONTEXT_MENU_CONFIG = {
  WIDTH: 240,                  // pixels
  ITEM_HEIGHT: 36,             // pixels
  PADDING: 8,                  // pixels
  BORDER_RADIUS: 4,            // pixels
  BACKGROUND: 'rgba(10, 10, 30, 0.95)',
  BORDER_COLOR: '#00FFCC',
  BORDER_WIDTH: 1,
  TEXT_COLOR: '#E0E0FF',
  HIGHLIGHT_BG: 'rgba(0, 255, 204, 0.15)',
  FONT_FAMILY: '"JetBrains Mono", "Fira Code", monospace',
  FONT_SIZE: '13px',
  DISMISS_ON_CLICK_OUTSIDE: true,
  DISMISS_ON_SCROLL: true,
  DISMISS_ON_ESC: true,
  APPEAR_ANIMATION: 'scale(0.95) -> scale(1.0), opacity 0->1, 120ms ease-out',
  POSITION_STRATEGY: 'prefer-bottom-right, flip-if-overflow'
};
```

---

### 10.7 Click-Outside Deselection

Clicking on empty space (no hit target) at any view level triggers a deselection or zoom-out:

| Current View | Selected State     | Click-Outside Action                       |
|--------------|--------------------|--------------------------------------------|
| Building     | No selection       | No action                                  |
| Floor        | Room highlighted   | Deselect room, close room panel            |
| Floor        | No selection       | Zoom out to Building View                  |
| Room         | Object highlighted | Deselect object                            |
| Room         | No selection       | Zoom out to Floor View                     |

**Deselection triggers:**
- Remove highlight/glow from previously selected element.
- Close any open info panel or tooltip (200ms fade-out).
- Update navigation state to remove selection.
- Do NOT push to browser history for deselection (only for view changes).

---

### 10.8 Click State Machine

The complete click interaction is governed by the following state machine. States represent the current selection/view context, and transitions represent user click actions.

```
+===========================================================================+
|                      HIC CLICK STATE MACHINE                              |
+===========================================================================+

                         +------------------+
                         |   BUILDING_VIEW  |
                         |   (no selection) |
                         +--------+---------+
                                  |
                    click floor   |   click outside
                    ----------->  |  <-----------+
                                  |              |
                         +--------v---------+    |
                         |   FLOOR_VIEW     |----+
                         |   (no selection) |
                         +--------+---------+
                           |              |
             click room    |              |  click outside
             ----------->  |              +--------+
                           |                       |
                  +--------v---------+             |
                  |   FLOOR_VIEW     |             |
                  |  (room selected) |             |
                  +--------+---------+             |
                    |         |    |                |
       click other  |  dbl-   | click              |
       room         |  click  | outside            |
       +--------->  |  room   | ------+            |
       |            |  |      |       |            |
       +------------+  |      |  +----v--------+  |
                        |      |  | FLOOR_VIEW  |--+
                        |      |  |(no select)  |
                        |      |  +-------------+
                        |      |
               +--------v------v--+
               |    ROOM_VIEW     |
               |   (no selection) |
               +--------+---------+
                  |              |
    click object  |              | click outside
    ----------->  |              +--------+
                  |                       |
         +--------v---------+            |
         |    ROOM_VIEW     |            |
         |  (obj selected)  |            |
         +--------+---------+            |
           |         |                   |
  click    |  click  |                   |
  other    |  outside|                   |
  obj      |         |                   |
  +------->+    +----v--------+         |
                | ROOM_VIEW   |---------+
                |(no select)  |
                +-------------+


  UNIVERSAL TRANSITIONS:
  +-------------------------------------------------+
  | Escape key at any level --> zoom out one level   |
  | Breadcrumb click        --> jump to that level   |
  | Context menu "Navigate" --> jump to target       |
  +-------------------------------------------------+
```

**Formal State Definitions:**

```
States = {
  S0: BUILDING_VIEW_IDLE,
  S1: FLOOR_VIEW_IDLE,
  S2: FLOOR_VIEW_ROOM_SELECTED,
  S3: ROOM_VIEW_IDLE,
  S4: ROOM_VIEW_OBJECT_SELECTED,
  S5: DOCUMENT_OPEN,
  S6: CONTEXT_MENU_OPEN
}

Transitions:
  S0 --[click_floor]--> S1       (zoom to floor)
  S0 --[click_outside]--> S0     (no-op)
  S1 --[click_room]--> S2        (highlight room)
  S1 --[click_outside]--> S0     (zoom out to building)
  S2 --[click_other_room]--> S2  (switch room selection)
  S2 --[dblclick_room]--> S3     (zoom to room)
  S2 --[click_enter]--> S3       (zoom to room)
  S2 --[click_outside]--> S1     (deselect room)
  S3 --[click_object]--> S4      (highlight object)
  S3 --[click_outside]--> S1     (zoom out to floor)
  S4 --[click_other_obj]--> S4   (switch object selection)
  S4 --[dblclick_object]--> S5   (open document)
  S4 --[click_outside]--> S3     (deselect object)
  S5 --[close_document]--> S4    (return to room, object still selected)
  ANY --[right_click]--> S6      (open context menu)
  S6 --[menu_action]--> varies   (depends on action chosen)
  S6 --[click_outside]--> prev   (close menu, restore previous state)
  ANY --[escape]--> parent_state (zoom out one level)
```

---

### 10.9 Touch Equivalents

All click interactions have touch equivalents to ensure full functionality on tablet and mobile devices.

| Mouse Action        | Touch Equivalent                                  | Timing         |
|---------------------|---------------------------------------------------|----------------|
| Single click        | Single tap                                        | < 300ms        |
| Double click        | Double tap                                        | < 250ms gap    |
| Right-click         | Long press                                        | 500ms hold     |
| Hover (see Agent 11)| Long press (lighter, 500ms)                       | 500ms hold     |
| Click-drag (pan)    | Single finger drag                                | Immediate      |
| Scroll wheel zoom   | Pinch gesture                                     | Immediate      |

**Touch Event Configuration:**

```javascript
const TOUCH_CONFIG = {
  TAP_MAX_DURATION_MS: 300,
  TAP_MAX_DRIFT_PX: 10,          // more generous than mouse (finger imprecision)
  DOUBLE_TAP_WINDOW_MS: 250,
  LONG_PRESS_THRESHOLD_MS: 500,
  LONG_PRESS_VIBRATION_MS: 10,   // haptic feedback if available
  PINCH_MIN_DISTANCE: 20,        // minimum px distance change to register as pinch
  PAN_THRESHOLD_PX: 8            // minimum movement before pan begins
};
```

**Long-Press Feedback:** When a long press reaches the 500ms threshold, the system provides visual feedback by rendering a radial neon pulse ring (color `#00FFCC`, expanding from 0 to 40px radius, 200ms animation) centered on the touch point, followed by opening the context menu.

**Touch Target Sizing:** All interactive elements must have a minimum touch target of 44x44 CSS pixels, per WCAG 2.5.5. Floor strips in Building View inherently exceed this. Room polygons in Floor View must have their hit-test area expanded to at least 44x44px around the centroid for small rooms.

---

## AGENT 11: HOVER LOGIC

### 11.1 Hover Detection Pipeline

Hover effects provide continuous visual feedback as the user moves their pointer across the HIC interface. The hover pipeline operates independently from the click pipeline but shares the same coordinate transformation and hit-testing infrastructure.

**Pipeline Architecture:**

```
STAGE 1: Pointer Move Capture
  - Listen to pointermove events on the SVG root
  - Throttle to 60fps (16.67ms interval) via requestAnimationFrame
  - Convert pointer coordinates to SVG space

STAGE 2: Hit Testing
  - Run view-specific hit test (same functions as click pipeline)
  - Compare result with previous hover target

STAGE 3: Hover State Transition
  - If target changed: trigger hover-exit on previous, hover-enter on new
  - If target same: no action (already in hover state)
  - If target null: trigger hover-exit on previous

STAGE 4: Visual Effect Application
  - Apply CSS class or SVG attribute changes
  - Start/cancel tooltip timers as needed
```

**Throttled Hover Handler:**

```javascript
class HoverController {
  constructor(svgRoot, hitTester) {
    this._svgRoot = svgRoot;
    this._hitTester = hitTester;
    this._currentTarget = null;
    this._tooltipTimer = null;
    this._rafId = null;
    this._lastPointerEvent = null;

    svgRoot.addEventListener('pointermove', (e) => {
      this._lastPointerEvent = e;
      if (!this._rafId) {
        this._rafId = requestAnimationFrame(() => this._processHover());
      }
    });

    svgRoot.addEventListener('pointerleave', () => {
      this._triggerHoverExit();
    });
  }

  _processHover() {
    this._rafId = null;
    const e = this._lastPointerEvent;
    if (!e) return;

    const svgPoint = clientToSVG(e, this._svgRoot);
    const target = this._hitTester.test(svgPoint);

    if (target === this._currentTarget) return;

    if (this._currentTarget) {
      this._triggerHoverExit();
    }

    if (target) {
      this._triggerHoverEnter(target, svgPoint);
    }
  }

  _triggerHoverEnter(target, position) {
    this._currentTarget = target;
    target.element.classList.add('hic-hover');

    // Start tooltip timer
    this._tooltipTimer = setTimeout(() => {
      this._showTooltip(target, position);
    }, HOVER_CONFIG.TOOLTIP_DELAY_MS);
  }

  _triggerHoverExit() {
    if (this._currentTarget) {
      this._currentTarget.element.classList.remove('hic-hover');
      this._currentTarget.element.classList.add('hic-hover-exit');

      // Remove exit class after fade completes
      setTimeout(() => {
        if (this._currentTarget) {
          this._currentTarget.element.classList.remove('hic-hover-exit');
        }
      }, HOVER_CONFIG.EXIT_FADE_MS);
    }

    clearTimeout(this._tooltipTimer);
    this._hideTooltip();
    this._currentTarget = null;
  }
}

const HOVER_CONFIG = {
  TOOLTIP_DELAY_MS: 150,
  EXIT_FADE_MS: 100,
  TOOLTIP_OFFSET_X: 16,
  TOOLTIP_OFFSET_Y: -8,
  TOOLTIP_MAX_WIDTH: 280,
  THROTTLE_INTERVAL: 16   // ~60fps, handled via rAF
};
```

---

### 11.2 Building View Hover

When the pointer moves over a floor strip in Building View, the floor visually responds with a neon highlight intensification.

**Hover-In Effects:**

1. Floor strip background opacity increases from 0.15 to 0.35 (150ms transition).
2. Neon border glow on the floor strip intensifies: `filter: drop-shadow()` blur radius from 2px to 6px, glow color brightens by 20%.
3. Floor strip expands vertically by 2px (1px each direction) to create a subtle "lift" effect.
4. Tooltip appears after 150ms delay showing:
   - Floor number (e.g., "FLOOR 07")
   - Floor name (e.g., "INTEL OPERATIONS")
   - Document count (e.g., "42 documents")
   - Status indicator (green/amber/red dot)

**Tooltip Markup:**

```html
<div class="hic-tooltip hic-tooltip--floor" role="tooltip" aria-live="polite">
  <div class="hic-tooltip__header">FLOOR 07</div>
  <div class="hic-tooltip__name">INTEL OPERATIONS</div>
  <div class="hic-tooltip__meta">
    <span class="hic-tooltip__status hic-tooltip__status--active"></span>
    <span>42 documents</span>
  </div>
</div>
```

---

### 11.3 Floor View Hover

When the pointer moves over a room polygon in Floor View, the room boundary illuminates.

**Hover-In Effects:**

1. Room polygon stroke width increases from 1px to 2px.
2. Room polygon stroke color shifts to full-brightness variant of the room's assigned glow color.
3. Room interior fill opacity increases from 0.08 to 0.18.
4. An outer glow is applied via SVG filter: `feGaussianBlur stdDeviation="4"` combined with `feComposite`.
5. Room label text (if visible at current zoom) transitions to full opacity.
6. Tooltip appears after 150ms delay showing:
   - Room name (e.g., "Threat Intelligence Archive")
   - Room status with label (e.g., "ACTIVE", "LOCKED", "MAINTENANCE")
   - Document count
   - Last updated timestamp
   - Access level required

**Hover-Exit Effects:**

1. All hover-in effects reverse over 100ms.
2. Room returns to its base visual state.
3. Tooltip fades out over 100ms.

---

### 11.4 Room View Hover

When the pointer moves over interactive elements within a room, those elements provide hover feedback.

**Object-Specific Hover Behaviors:**

| Object Type     | Hover Effect                                                              |
|-----------------|---------------------------------------------------------------------------|
| `document`      | Icon glows brighter, document title appears as tooltip, subtle pulse      |
| `sub-system`    | Portal frame pulses with neon animation, destination name in tooltip      |
| `terminal`      | Screen text flickers, "Click to interact" in tooltip                     |
| `data-node`     | Connection lines brighten, metadata summary in tooltip                   |
| `external-link` | External link icon appears, URL preview in tooltip                       |
| `locked-element`| Red glow pulse, "ACCESS RESTRICTED" in tooltip                           |

**Document Hover Tooltip Content:**

```
+--------------------------------------+
|  [doc-icon]  DOCUMENT                |
|  Signal Intelligence Protocol v3.2   |
|  ----------------------------------- |
|  Last modified: 2026-02-14           |
|  Author: OPS-ANALYST-07             |
|  Size: 24 KB                        |
|  Classification: INTERNAL            |
+--------------------------------------+
```

---

### 11.5 Tooltip Timing and Lifecycle

Tooltips follow a strict lifecycle to avoid flickering and ensure smooth user experience.

**Timing Diagram:**

```
Pointer enters target
  |
  |--- 0ms:   Hover highlight applied immediately
  |
  |--- 150ms: Tooltip delay elapses
  |            Tooltip begins fade-in (80ms, opacity 0 -> 1)
  |
  |--- 230ms: Tooltip fully visible
  |
  ...pointer remains over target...
  |
Pointer exits target
  |
  |--- 0ms:   Hover highlight begins fade-out (100ms)
  |            Tooltip begins fade-out (100ms, opacity 1 -> 0)
  |
  |--- 100ms: Tooltip removed from DOM
  |            Hover highlight removed
  |

SPECIAL CASE: Rapid movement between adjacent targets
  Pointer exits target A, enters target B within 50ms
  |
  |--- 0ms:   Target A hover-exit starts
  |            Target B hover-enter starts
  |            Tooltip repositions to target B (no fade-out/in, just translate)
  |            Tooltip content updates immediately
  |
  (This prevents tooltip flickering when moving between adjacent floors/rooms)
```

**Grace Period Implementation:**

```javascript
const TOOLTIP_GRACE_PERIOD_MS = 50;

// In HoverController, modify _triggerHoverExit:
_triggerHoverExit() {
  this._graceTimer = setTimeout(() => {
    // only truly exit if no new target was entered during grace period
    if (!this._currentTarget) {
      this._performExit();
    }
  }, TOOLTIP_GRACE_PERIOD_MS);
}
```

---

### 11.6 CSS/SVG Hover State Definitions

**CSS Transition Definitions:**

```css
/* === BASE STATES === */

.hic-floor-strip {
  fill-opacity: 0.15;
  stroke-width: 1;
  stroke-opacity: 0.6;
  filter: drop-shadow(0 0 2px var(--floor-glow-color));
  transition:
    fill-opacity 150ms ease-in-out,
    stroke-width 150ms ease-in-out,
    stroke-opacity 150ms ease-in-out,
    filter 150ms ease-in-out,
    transform 150ms ease-in-out;
  transform-origin: center center;
  will-change: fill-opacity, filter, transform;
}

.hic-room-polygon {
  fill-opacity: 0.08;
  stroke: var(--room-glow-color);
  stroke-width: 1;
  stroke-opacity: 0.5;
  filter: none;
  transition:
    fill-opacity 150ms ease-in-out,
    stroke-width 150ms ease-in-out,
    stroke-opacity 150ms ease-in-out,
    filter 200ms ease-in-out;
  will-change: fill-opacity, stroke-width, filter;
}

.hic-room-object {
  opacity: 0.7;
  filter: drop-shadow(0 0 1px var(--object-glow-color));
  transition:
    opacity 120ms ease-in-out,
    filter 120ms ease-in-out,
    transform 120ms ease-in-out;
  will-change: opacity, filter;
}


/* === HOVER STATES === */

.hic-floor-strip.hic-hover {
  fill-opacity: 0.35;
  stroke-width: 1.5;
  stroke-opacity: 1.0;
  filter: drop-shadow(0 0 6px var(--floor-glow-color))
          drop-shadow(0 0 12px var(--floor-glow-color-dim));
  transform: scaleY(1.04);
}

.hic-room-polygon.hic-hover {
  fill-opacity: 0.18;
  stroke-width: 2;
  stroke-opacity: 1.0;
  filter: drop-shadow(0 0 8px var(--room-glow-color))
          drop-shadow(0 0 16px var(--room-glow-color-dim));
}

.hic-room-object.hic-hover {
  opacity: 1.0;
  filter: drop-shadow(0 0 4px var(--object-glow-color))
          drop-shadow(0 0 8px var(--object-glow-color));
  transform: scale(1.05);
}


/* === HOVER EXIT STATES (for fade-out animation) === */

.hic-floor-strip.hic-hover-exit {
  fill-opacity: 0.15;
  stroke-width: 1;
  stroke-opacity: 0.6;
  filter: drop-shadow(0 0 2px var(--floor-glow-color));
  transform: scaleY(1.0);
  transition-duration: 100ms;
}

.hic-room-polygon.hic-hover-exit {
  fill-opacity: 0.08;
  stroke-width: 1;
  stroke-opacity: 0.5;
  filter: none;
  transition-duration: 100ms;
}

.hic-room-object.hic-hover-exit {
  opacity: 0.7;
  filter: drop-shadow(0 0 1px var(--object-glow-color));
  transform: scale(1.0);
  transition-duration: 100ms;
}


/* === TOOLTIP === */

.hic-tooltip {
  position: absolute;
  pointer-events: none;
  background: rgba(10, 10, 30, 0.92);
  border: 1px solid var(--tooltip-border-color, #00FFCC);
  border-radius: 4px;
  padding: 10px 14px;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 12px;
  color: #E0E0FF;
  max-width: 280px;
  box-shadow:
    0 0 8px rgba(0, 255, 204, 0.3),
    0 4px 16px rgba(0, 0, 0, 0.5);
  opacity: 0;
  transform: translateY(4px);
  transition:
    opacity 80ms ease-out,
    transform 80ms ease-out;
  z-index: 9000;
}

.hic-tooltip.hic-tooltip--visible {
  opacity: 1;
  transform: translateY(0);
}

.hic-tooltip.hic-tooltip--exiting {
  opacity: 0;
  transform: translateY(4px);
  transition-duration: 100ms;
}

.hic-tooltip__header {
  font-size: 10px;
  letter-spacing: 2px;
  text-transform: uppercase;
  color: var(--tooltip-border-color, #00FFCC);
  margin-bottom: 4px;
}

.hic-tooltip__name {
  font-size: 14px;
  font-weight: 600;
  color: #FFFFFF;
  margin-bottom: 6px;
}

.hic-tooltip__meta {
  font-size: 11px;
  color: #A0A0CC;
  display: flex;
  align-items: center;
  gap: 6px;
}

.hic-tooltip__status {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.hic-tooltip__status--active {
  background: #00FF88;
  box-shadow: 0 0 6px #00FF88;
}

.hic-tooltip__status--locked {
  background: #FFAA00;
  box-shadow: 0 0 6px #FFAA00;
}

.hic-tooltip__status--maintenance {
  background: #FF3366;
  box-shadow: 0 0 6px #FF3366;
}
```

**SVG Filter Definitions (for glow effects):**

```xml
<defs>
  <filter id="hic-glow-subtle" x="-20%" y="-20%" width="140%" height="140%">
    <feGaussianBlur in="SourceGraphic" stdDeviation="2" result="blur" />
    <feComposite in="SourceGraphic" in2="blur" operator="over" />
  </filter>

  <filter id="hic-glow-intense" x="-30%" y="-30%" width="160%" height="160%">
    <feGaussianBlur in="SourceGraphic" stdDeviation="4" result="blur" />
    <feColorMatrix in="blur" type="saturate" values="2" result="saturated" />
    <feComposite in="SourceGraphic" in2="saturated" operator="over" />
  </filter>

  <filter id="hic-glow-hover" x="-40%" y="-40%" width="180%" height="180%">
    <feGaussianBlur in="SourceAlpha" stdDeviation="6" result="blur" />
    <feFlood flood-color="var(--glow-color, #00FFCC)" flood-opacity="0.6" result="color" />
    <feComposite in="color" in2="blur" operator="in" result="glow" />
    <feMerge>
      <feMergeNode in="glow" />
      <feMergeNode in="SourceGraphic" />
    </feMerge>
  </filter>
</defs>
```

---

### 11.7 Performance Constraints

Hover detection is one of the most performance-critical systems in the HIC because it fires on every pointer movement. The following constraints ensure smooth 60fps rendering.

**Throttling Strategy:**

- Pointer move events are throttled to one processing cycle per animation frame via `requestAnimationFrame`.
- The rAF callback only fires if there is a pending pointer event; it does not run continuously.
- Hit testing reuses spatial index structures (R-tree for room polygons, simple array scan for floors since floors are axis-aligned and O(1) via Y-coordinate division).

**Performance Budget per Hover Frame:**

```
Total budget per frame at 60fps: 16.67ms

Coordinate transformation:    ~0.1ms
Hit testing (floor, Y-div):   ~0.05ms
Hit testing (room, polygon):  ~0.3ms (worst case, 30 rooms)
CSS class toggle:              ~0.1ms
Tooltip position update:       ~0.2ms
------------------------------------------
Total:                         ~0.75ms  (well within budget)
```

**Spatial Index for Room Hit Testing:**

For floors with many rooms (>15), an R-tree spatial index is built on floor load to accelerate point-in-polygon queries:

```javascript
class SpatialIndex {
  constructor(rooms) {
    this._tree = new RBush();
    const items = rooms.map(room => {
      const bbox = this._computeBBox(room.vertices);
      return { ...bbox, room };
    });
    this._tree.load(items);
  }

  _computeBBox(vertices) {
    let minX = Infinity, minY = Infinity, maxX = -Infinity, maxY = -Infinity;
    for (const v of vertices) {
      if (v.x < minX) minX = v.x;
      if (v.y < minY) minY = v.y;
      if (v.x > maxX) maxX = v.x;
      if (v.y > maxY) maxY = v.y;
    }
    return { minX, minY, maxX, maxY };
  }

  query(point) {
    const candidates = this._tree.search({
      minX: point.x, minY: point.y,
      maxX: point.x, maxY: point.y
    });
    for (const candidate of candidates) {
      if (pointInPolygon(point, candidate.room.vertices)) {
        return candidate.room;
      }
    }
    return null;
  }
}
```

**GPU Acceleration:**

All hover transitions use properties that can be hardware-accelerated by the GPU compositor:

- `opacity` -- composited on GPU, no layout/paint.
- `transform` -- composited on GPU, no layout/paint.
- `filter` (drop-shadow/blur) -- GPU-accelerated in modern browsers.
- `will-change` declarations are applied to elements that frequently change.

Properties to AVOID in hover transitions (trigger layout/paint):

- `width`, `height` -- triggers layout.
- `top`, `left` -- triggers layout.
- `border-width` -- triggers paint. (Use `stroke-width` on SVG which does not trigger HTML layout.)
- `box-shadow` -- triggers paint. (Use `filter: drop-shadow()` instead.)

---

### 11.8 Mobile Long-Press Hover Equivalent

On touch devices, there is no pointer hover. The HIC implements a long-press gesture to serve as a hover equivalent, providing the same informational tooltips and visual highlighting without triggering a click action.

**Long-Press Detection:**

```javascript
class LongPressDetector {
  constructor(element, callback, config = {}) {
    this._threshold = config.threshold || 500; // ms
    this._driftMax = config.driftMax || 10;    // px
    this._timer = null;
    this._startPos = null;
    this._active = false;

    element.addEventListener('touchstart', (e) => this._onStart(e), { passive: false });
    element.addEventListener('touchmove', (e) => this._onMove(e), { passive: false });
    element.addEventListener('touchend', (e) => this._onEnd(e));
    element.addEventListener('touchcancel', () => this._cancel());

    this._callback = callback;
  }

  _onStart(e) {
    if (e.touches.length !== 1) return; // ignore multi-touch
    this._startPos = { x: e.touches[0].clientX, y: e.touches[0].clientY };
    this._timer = setTimeout(() => {
      this._active = true;
      e.preventDefault(); // prevent subsequent click
      // Haptic feedback
      if (navigator.vibrate) navigator.vibrate(10);
      this._callback({
        type: 'long-press',
        clientX: this._startPos.x,
        clientY: this._startPos.y,
        originalEvent: e
      });
    }, this._threshold);
  }

  _onMove(e) {
    if (!this._startPos || !this._timer) return;
    const dx = e.touches[0].clientX - this._startPos.x;
    const dy = e.touches[0].clientY - this._startPos.y;
    if (Math.sqrt(dx * dx + dy * dy) > this._driftMax) {
      this._cancel(); // finger moved too far, this is a pan not a long-press
    }
  }

  _onEnd(e) {
    if (this._active) {
      e.preventDefault(); // prevent click event after long-press
      this._active = false;
    }
    this._cancel();
  }

  _cancel() {
    clearTimeout(this._timer);
    this._timer = null;
    this._startPos = null;
  }
}
```

**Long-Press Behavior Differences from Desktop Hover:**

| Aspect           | Desktop Hover                     | Mobile Long-Press                      |
|------------------|-----------------------------------|----------------------------------------|
| Trigger          | Pointer enters element bounds     | 500ms hold on element                  |
| Duration         | Persists while pointer is over    | Tooltip dismissed on finger lift       |
| Tooltip position | Follows pointer with offset       | Appears above touch point (fixed)      |
| Highlight        | Applied/removed dynamically       | Applied on threshold, removed on lift  |
| Exit             | Pointer leaves element bounds     | Finger lifted or dragged away          |

---

## AGENT 12: ZOOM/PAN AND NAVIGATION STATE

### 12.1 Zoom Level Definitions

The HIC operates at five semantic zoom levels. Each level corresponds to a specific view of the skyscraper and determines which elements are rendered, interactive, and labeled.

| Zoom Level | Multiplier | View Name    | Content Visible                                             | Label Density |
|------------|------------|--------------|-------------------------------------------------------------|---------------|
| Level 0    | 0.1x       | Full Building| Entire skyscraper silhouette, all floors as thin strips     | Floor numbers only |
| Level 1    | 0.5x       | Section      | 5-7 floor section, floor strips with names                  | Floor names, basic stats |
| Level 2    | 1.0x       | Floor        | Single floor cross-section, all rooms visible               | Room names, status |
| Level 3    | 2.0x       | Zone         | Portion of floor, rooms with interior detail beginning      | Room names, door labels |
| Level 4    | 3.0x       | Room         | Single room interior, all objects/documents visible         | All labels, document titles |

**Zoom Level Constants:**

```javascript
const ZOOM_LEVELS = {
  BUILDING: { level: 0, scale: 0.1,  label: 'Full Building' },
  SECTION:  { level: 1, scale: 0.5,  label: 'Section' },
  FLOOR:    { level: 2, scale: 1.0,  label: 'Floor' },
  ZONE:     { level: 3, scale: 2.0,  label: 'Zone' },
  ROOM:     { level: 4, scale: 3.0,  label: 'Room' }
};

// Content visibility thresholds (elements appear/disappear at these scales)
const VISIBILITY_THRESHOLDS = {
  floorNumbers:      0.05,   // always visible
  floorNames:        0.3,    // visible from section view
  floorStats:        0.4,    // visible from section view
  roomOutlines:      0.7,    // visible approaching floor view
  roomNames:         0.9,    // visible at floor view
  roomStatus:        0.9,    // visible at floor view
  roomInterior:      1.5,    // visible approaching zone view
  doorLabels:        1.8,    // visible at zone view
  objectIcons:       2.2,    // visible approaching room view
  documentTitles:    2.5,    // visible at room view
  documentMetadata:  2.8     // visible at deep room view
};
```

**Level-of-Detail Rendering:**

As the zoom level changes continuously (during pinch or scroll), elements fade in and out according to their visibility threshold. The fade range is +/- 0.1 of the threshold scale:

```javascript
function computeElementOpacity(currentScale, threshold) {
  const fadeRange = 0.1;
  if (currentScale >= threshold) return 1.0;
  if (currentScale < threshold - fadeRange) return 0.0;
  return (currentScale - (threshold - fadeRange)) / fadeRange;
}
```

---

### 12.2 Zoom Transition Animation

When the user triggers a discrete zoom action (click a floor, press +/- key, double-click), the view transitions smoothly to the target zoom level and focal point.

**Animation Parameters:**

```javascript
const ZOOM_ANIMATION = {
  DURATION_MS: 300,
  EASING: 'cubic-bezier(0.25, 0.1, 0.25, 1.0)', // ease-in-out equivalent
  DOUBLE_CLICK_DURATION_MS: 200,                   // faster for double-click
  KEYBOARD_DURATION_MS: 150,                        // even faster for key press
  MIN_DURATION_MS: 100,                             // never shorter than this
  MAX_DURATION_MS: 500                              // never longer than this
};
```

**Zoom Transition Implementation:**

```javascript
class ZoomAnimator {
  constructor(viewport) {
    this._viewport = viewport;
    this._animation = null;
  }

  zoomTo(targetScale, targetCenter, durationMs = 300, easing = 'ease-in-out') {
    // Cancel any in-progress animation
    if (this._animation) {
      this._animation.cancel();
    }

    const startScale = this._viewport.currentScale;
    const startCenter = { ...this._viewport.currentCenter };
    const startTime = performance.now();

    return new Promise((resolve) => {
      const animate = (now) => {
        const elapsed = now - startTime;
        const t = Math.min(elapsed / durationMs, 1.0);
        const easedT = this._applyEasing(t, easing);

        // Interpolate scale (logarithmic for perceptual linearity)
        const logStart = Math.log(startScale);
        const logEnd = Math.log(targetScale);
        const currentScale = Math.exp(logStart + (logEnd - logStart) * easedT);

        // Interpolate center position (linear)
        const currentCenter = {
          x: startCenter.x + (targetCenter.x - startCenter.x) * easedT,
          y: startCenter.y + (targetCenter.y - startCenter.y) * easedT
        };

        this._viewport.setView(currentScale, currentCenter);

        if (t < 1.0) {
          this._animation = requestAnimationFrame(animate);
        } else {
          this._animation = null;
          resolve();
        }
      };

      this._animation = requestAnimationFrame(animate);
    });
  }

  _applyEasing(t, easing) {
    // Built-in easing functions
    switch (easing) {
      case 'ease-in-out':
        return t < 0.5
          ? 4 * t * t * t
          : 1 - Math.pow(-2 * t + 2, 3) / 2;
      case 'ease-out':
        return 1 - Math.pow(1 - t, 3);
      case 'ease-in':
        return t * t * t;
      case 'linear':
        return t;
      default:
        return t;
    }
  }
}
```

**Logarithmic Scale Interpolation:** The zoom interpolation uses logarithmic interpolation rather than linear. This ensures that zooming from 0.1x to 1.0x (a 10x change) feels the same speed as zooming from 1.0x to 3.0x (a 3x change in absolute terms but a similar perceptual change). Without logarithmic interpolation, zooming out would feel sluggish and zooming in would feel abrupt.

---

### 12.3 Pan Mechanics

Panning allows the user to move the viewport when zoomed in past the point where the entire view fits on screen.

**Pan Trigger:** Click and drag (mouse button 0 held) on empty space (no interactive element hit).

**Pan Bounds:** The viewport is constrained so the user cannot pan beyond the content area. A soft margin of 50 SVG units is allowed beyond the content bounds to give a sense of space without losing the building entirely.

```javascript
const PAN_CONFIG = {
  SOFT_MARGIN: 50,           // SVG units beyond content bounds
  INERTIA_ENABLED: true,
  INERTIA_FRICTION: 0.92,    // velocity multiplier per frame
  INERTIA_MIN_VELOCITY: 0.5, // px/frame, below this inertia stops
  INERTIA_MAX_VELOCITY: 50   // px/frame, clamp to prevent runaway
};

class PanController {
  constructor(viewport, contentBounds) {
    this._viewport = viewport;
    this._bounds = contentBounds;
    this._isPanning = false;
    this._lastPoint = null;
    this._velocity = { x: 0, y: 0 };
  }

  onPointerDown(svgPoint) {
    this._isPanning = true;
    this._lastPoint = svgPoint;
    this._velocity = { x: 0, y: 0 };
  }

  onPointerMove(svgPoint) {
    if (!this._isPanning) return;

    const dx = svgPoint.x - this._lastPoint.x;
    const dy = svgPoint.y - this._lastPoint.y;

    this._velocity = { x: dx, y: dy };
    this._lastPoint = svgPoint;

    this._viewport.pan(dx, dy);
    this._clampToBounds();
  }

  onPointerUp() {
    this._isPanning = false;
    if (PAN_CONFIG.INERTIA_ENABLED) {
      this._startInertia();
    }
  }

  _clampToBounds() {
    const view = this._viewport.getViewBox();
    const margin = PAN_CONFIG.SOFT_MARGIN;

    const minX = this._bounds.minX - margin;
    const maxX = this._bounds.maxX + margin - view.width;
    const minY = this._bounds.minY - margin;
    const maxY = this._bounds.maxY + margin - view.height;

    this._viewport.clamp(minX, maxX, minY, maxY);
  }

  _startInertia() {
    const step = () => {
      this._velocity.x *= PAN_CONFIG.INERTIA_FRICTION;
      this._velocity.y *= PAN_CONFIG.INERTIA_FRICTION;

      const speed = Math.sqrt(
        this._velocity.x ** 2 + this._velocity.y ** 2
      );

      if (speed < PAN_CONFIG.INERTIA_MIN_VELOCITY) return;

      this._viewport.pan(this._velocity.x, this._velocity.y);
      this._clampToBounds();
      requestAnimationFrame(step);
    };

    requestAnimationFrame(step);
  }
}
```

---

### 12.4 Scroll Wheel Zoom

The scroll wheel (or trackpad scroll gesture) zooms in and out, centered on the cursor position.

**Behavior:**

- Scroll up (positive deltaY): zoom in.
- Scroll down (negative deltaY): zoom out.
- Zoom center: the SVG point under the cursor remains fixed (the view scales around that point).

```javascript
const SCROLL_ZOOM_CONFIG = {
  SCALE_FACTOR: 0.001,       // scale change per deltaY unit
  MIN_SCALE: 0.08,           // cannot zoom out past this
  MAX_SCALE: 4.0,            // cannot zoom in past this
  SMOOTH_SCROLL_LERP: 0.15,  // smoothing factor for trackpad gestures
  DISCRETE_STEP: 0.2         // scale change per mouse wheel "click"
};

function handleScrollZoom(event, viewport, svgElement) {
  event.preventDefault();

  const svgPoint = clientToSVG(event, svgElement);

  // Detect discrete vs smooth scroll
  const isDiscrete = (event.deltaMode === WheelEvent.DOM_DELTA_LINE)
    || Math.abs(event.deltaY) > 50;

  let scaleDelta;
  if (isDiscrete) {
    // Mouse wheel: fixed step
    scaleDelta = event.deltaY > 0
      ? -SCROLL_ZOOM_CONFIG.DISCRETE_STEP
      : SCROLL_ZOOM_CONFIG.DISCRETE_STEP;
  } else {
    // Trackpad: proportional
    scaleDelta = -event.deltaY * SCROLL_ZOOM_CONFIG.SCALE_FACTOR;
  }

  const currentScale = viewport.currentScale;
  const newScale = Math.max(
    SCROLL_ZOOM_CONFIG.MIN_SCALE,
    Math.min(SCROLL_ZOOM_CONFIG.MAX_SCALE, currentScale * (1 + scaleDelta))
  );

  // Zoom centered on cursor position
  // The point under the cursor should remain at the same screen position
  const ratio = newScale / currentScale;
  const newCenterX = svgPoint.x - (svgPoint.x - viewport.currentCenter.x) * ratio;
  const newCenterY = svgPoint.y - (svgPoint.y - viewport.currentCenter.y) * ratio;

  viewport.setView(newScale, { x: newCenterX, y: newCenterY });
}
```

---

### 12.5 Pinch Zoom on Touch Devices

Pinch zoom uses two-finger touch gestures to zoom in and out, centered on the midpoint between the two fingers.

```javascript
class PinchZoomController {
  constructor(viewport, svgElement) {
    this._viewport = viewport;
    this._svgElement = svgElement;
    this._initialDistance = null;
    this._initialScale = null;
    this._initialMidpoint = null;

    svgElement.addEventListener('touchstart', (e) => this._onTouchStart(e), { passive: false });
    svgElement.addEventListener('touchmove', (e) => this._onTouchMove(e), { passive: false });
    svgElement.addEventListener('touchend', (e) => this._onTouchEnd(e));
  }

  _onTouchStart(e) {
    if (e.touches.length === 2) {
      e.preventDefault();
      this._initialDistance = this._getDistance(e.touches[0], e.touches[1]);
      this._initialScale = this._viewport.currentScale;
      this._initialMidpoint = this._getMidpoint(e.touches[0], e.touches[1]);
    }
  }

  _onTouchMove(e) {
    if (e.touches.length !== 2 || !this._initialDistance) return;
    e.preventDefault();

    const currentDistance = this._getDistance(e.touches[0], e.touches[1]);
    const currentMidpoint = this._getMidpoint(e.touches[0], e.touches[1]);

    // Calculate new scale
    const ratio = currentDistance / this._initialDistance;
    const newScale = Math.max(
      SCROLL_ZOOM_CONFIG.MIN_SCALE,
      Math.min(SCROLL_ZOOM_CONFIG.MAX_SCALE, this._initialScale * ratio)
    );

    // Calculate pan offset from midpoint movement
    const svgMidpoint = clientToSVG(
      { clientX: currentMidpoint.x, clientY: currentMidpoint.y },
      this._svgElement
    );

    this._viewport.setView(newScale, svgMidpoint);
  }

  _onTouchEnd(e) {
    if (e.touches.length < 2) {
      this._initialDistance = null;
      this._initialScale = null;
      this._initialMidpoint = null;
    }
  }

  _getDistance(touch1, touch2) {
    const dx = touch1.clientX - touch2.clientX;
    const dy = touch1.clientY - touch2.clientY;
    return Math.sqrt(dx * dx + dy * dy);
  }

  _getMidpoint(touch1, touch2) {
    return {
      x: (touch1.clientX + touch2.clientX) / 2,
      y: (touch1.clientY + touch2.clientY) / 2
    };
  }
}
```

---

### 12.6 Navigation State Object

The navigation state object is the single source of truth for the HIC's current view configuration. All view changes flow through this object. It is serializable to JSON for persistence and URL encoding.

**State Object Definition:**

```json
{
  "version": 1,
  "current_view": "floor",
  "floor_id": 7,
  "room_id": "room-7-03",
  "selected_object_id": null,
  "zoom_level": 1.0,
  "pan_offset": {
    "x": 0,
    "y": -120.5
  },
  "breadcrumb_trail": [
    { "view": "building", "label": "Building", "floor_id": null, "room_id": null },
    { "view": "floor", "label": "Floor 7: Intel Ops", "floor_id": 7, "room_id": null },
    { "view": "room", "label": "Threat Archive", "floor_id": 7, "room_id": "room-7-03" }
  ],
  "timestamp": 1739836200000,
  "session_id": "hic-sess-a4f29c"
}
```

**State Controller:**

```javascript
class NavigationState {
  constructor() {
    this._state = {
      version: 1,
      current_view: 'building',
      floor_id: null,
      room_id: null,
      selected_object_id: null,
      zoom_level: 0.1,
      pan_offset: { x: 0, y: 0 },
      breadcrumb_trail: [
        { view: 'building', label: 'Building', floor_id: null, room_id: null }
      ],
      timestamp: Date.now(),
      session_id: this._generateSessionId()
    };
    this._listeners = new Set();
  }

  get state() {
    return Object.freeze({ ...this._state });
  }

  navigateToFloor(floorId, floorName) {
    this._state.current_view = 'floor';
    this._state.floor_id = floorId;
    this._state.room_id = null;
    this._state.selected_object_id = null;
    this._state.zoom_level = 1.0;
    this._state.pan_offset = { x: 0, y: 0 };
    this._state.breadcrumb_trail = [
      { view: 'building', label: 'Building', floor_id: null, room_id: null },
      { view: 'floor', label: `Floor ${floorId}: ${floorName}`, floor_id: floorId, room_id: null }
    ];
    this._state.timestamp = Date.now();
    this._notify();
  }

  navigateToRoom(roomId, roomName) {
    this._state.current_view = 'room';
    this._state.room_id = roomId;
    this._state.selected_object_id = null;
    this._state.zoom_level = 3.0;
    this._state.pan_offset = { x: 0, y: 0 };
    this._state.breadcrumb_trail.push({
      view: 'room',
      label: roomName,
      floor_id: this._state.floor_id,
      room_id: roomId
    });
    this._state.timestamp = Date.now();
    this._notify();
  }

  navigateUp() {
    if (this._state.breadcrumb_trail.length <= 1) return;
    this._state.breadcrumb_trail.pop();
    const target = this._state.breadcrumb_trail[this._state.breadcrumb_trail.length - 1];
    this._state.current_view = target.view;
    this._state.floor_id = target.floor_id;
    this._state.room_id = target.room_id;
    this._state.selected_object_id = null;
    this._state.zoom_level = ZOOM_LEVELS_BY_VIEW[target.view];
    this._state.pan_offset = { x: 0, y: 0 };
    this._state.timestamp = Date.now();
    this._notify();
  }

  selectObject(objectId) {
    this._state.selected_object_id = objectId;
    this._state.timestamp = Date.now();
    this._notify();
  }

  updatePan(offset) {
    this._state.pan_offset = { ...offset };
    // Do NOT notify listeners for pan (too frequent), use separate pan channel
  }

  updateZoom(level) {
    this._state.zoom_level = level;
    this._state.timestamp = Date.now();
    this._notify();
  }

  subscribe(listener) {
    this._listeners.add(listener);
    return () => this._listeners.delete(listener);
  }

  _notify() {
    const frozenState = this.state;
    for (const listener of this._listeners) {
      listener(frozenState);
    }
  }

  _generateSessionId() {
    return 'hic-sess-' + Math.random().toString(36).substring(2, 8);
  }

  toJSON() {
    return JSON.stringify(this._state);
  }

  static fromJSON(json) {
    const state = new NavigationState();
    Object.assign(state._state, JSON.parse(json));
    return state;
  }
}

const ZOOM_LEVELS_BY_VIEW = {
  building: 0.1,
  section: 0.5,
  floor: 1.0,
  zone: 2.0,
  room: 3.0
};
```

---

### 12.7 URL Hash Encoding

The navigation state is encoded in the URL hash fragment to enable bookmarking and link sharing. The hash is updated on every view change (but not on pan or smooth zoom to avoid excessive history entries).

**Hash Format:**

```
#/building
#/floor/7
#/floor/7/room/room-7-03
#/floor/7/room/room-7-03/obj/doc-threat-intel-v3
#/floor/7?z=1.5&px=-120&py=50
```

**Encoding/Decoding Implementation:**

```javascript
class HashEncoder {
  static encode(state) {
    let hash = '#';

    switch (state.current_view) {
      case 'building':
        hash += '/building';
        break;
      case 'floor':
        hash += `/floor/${state.floor_id}`;
        break;
      case 'room':
        hash += `/floor/${state.floor_id}/room/${state.room_id}`;
        break;
    }

    if (state.selected_object_id) {
      hash += `/obj/${state.selected_object_id}`;
    }

    // Append zoom and pan as query parameters if non-default
    const params = [];
    const defaultZoom = ZOOM_LEVELS_BY_VIEW[state.current_view];
    if (Math.abs(state.zoom_level - defaultZoom) > 0.01) {
      params.push(`z=${state.zoom_level.toFixed(2)}`);
    }
    if (Math.abs(state.pan_offset.x) > 1 || Math.abs(state.pan_offset.y) > 1) {
      params.push(`px=${Math.round(state.pan_offset.x)}`);
      params.push(`py=${Math.round(state.pan_offset.y)}`);
    }

    if (params.length > 0) {
      hash += '?' + params.join('&');
    }

    return hash;
  }

  static decode(hash) {
    const state = {
      current_view: 'building',
      floor_id: null,
      room_id: null,
      selected_object_id: null,
      zoom_level: 0.1,
      pan_offset: { x: 0, y: 0 }
    };

    if (!hash || hash === '#' || hash === '#/') {
      return state;
    }

    // Separate path from query params
    const [path, query] = hash.replace(/^#/, '').split('?');
    const segments = path.split('/').filter(Boolean);

    let i = 0;
    while (i < segments.length) {
      switch (segments[i]) {
        case 'building':
          state.current_view = 'building';
          state.zoom_level = 0.1;
          i++;
          break;
        case 'floor':
          state.current_view = 'floor';
          state.floor_id = parseInt(segments[++i], 10);
          state.zoom_level = 1.0;
          i++;
          break;
        case 'room':
          state.current_view = 'room';
          state.room_id = segments[++i];
          state.zoom_level = 3.0;
          i++;
          break;
        case 'obj':
          state.selected_object_id = segments[++i];
          i++;
          break;
        default:
          i++;
      }
    }

    // Parse query parameters
    if (query) {
      const params = new URLSearchParams(query);
      if (params.has('z')) state.zoom_level = parseFloat(params.get('z'));
      if (params.has('px')) state.pan_offset.x = parseInt(params.get('px'), 10);
      if (params.has('py')) state.pan_offset.y = parseInt(params.get('py'), 10);
    }

    return state;
  }
}
```

---

### 12.8 Browser History Integration

The HIC integrates with the browser's History API so that the Back and Forward buttons navigate the skyscraper's view hierarchy naturally.

**History Push Strategy:**

- View changes (building to floor, floor to room) push a new history entry.
- Object selection pushes a new history entry.
- Pan and continuous zoom do NOT push history entries (would flood the stack).
- Deselection (click outside) does NOT push a history entry; it replaces the current entry.

```javascript
class HistoryIntegration {
  constructor(navigationState) {
    this._navState = navigationState;
    this._ignorePopState = false;

    // Listen for navigation state changes
    navigationState.subscribe((state) => {
      if (this._ignorePopState) {
        this._ignorePopState = false;
        return;
      }
      const hash = HashEncoder.encode(state);
      if (window.location.hash !== hash) {
        window.history.pushState(
          { hicState: state },
          '',
          hash
        );
      }
    });

    // Listen for browser back/forward
    window.addEventListener('popstate', (event) => {
      if (event.state && event.state.hicState) {
        this._ignorePopState = true;
        this._navState.restore(event.state.hicState);
      } else {
        // Decode from URL hash
        this._ignorePopState = true;
        const decoded = HashEncoder.decode(window.location.hash);
        this._navState.restore(decoded);
      }
    });

    // Initialize from current URL hash on page load
    const initialState = HashEncoder.decode(window.location.hash);
    if (initialState.current_view !== 'building') {
      this._ignorePopState = true;
      this._navState.restore(initialState);
    }
  }
}
```

**History State Diagram:**

```
Browser History Stack (example user session):

  [0] #/building                    <-- initial load
  [1] #/floor/3                     <-- clicked Floor 3
  [2] #/floor/3/room/room-3-07     <-- entered Threat Archive room
  [3] #/floor/3/room/room-3-07/obj/doc-42  <-- selected a document
  [4] #/floor/3                     <-- pressed Escape (back to floor)
  [5] #/floor/9                     <-- clicked Floor 9 from building view
       ^
       current position

  User presses Back: --> goes to [4] #/floor/3
  User presses Back: --> goes to [3] #/floor/3/room/room-3-07/obj/doc-42
  User presses Forward: --> goes to [4] #/floor/3
```

---

### 12.9 Keyboard Shortcuts

Full keyboard navigation ensures accessibility and power-user efficiency.

**Keyboard Shortcut Map:**

| Key             | Action                                        | Context          |
|-----------------|-----------------------------------------------|------------------|
| `Arrow Up`      | Pan up by 50 SVG units                        | All views        |
| `Arrow Down`    | Pan down by 50 SVG units                      | All views        |
| `Arrow Left`    | Pan left by 50 SVG units                      | All views        |
| `Arrow Right`   | Pan right by 50 SVG units                     | All views        |
| `+` / `=`       | Zoom in one step (scale * 1.25)               | All views        |
| `-`             | Zoom out one step (scale * 0.8)               | All views        |
| `Escape`        | Zoom out one level / deselect                 | All views        |
| `Enter`         | Zoom into selected element                    | When selected    |
| `Tab`           | Cycle focus to next interactive element        | All views        |
| `Shift+Tab`     | Cycle focus to previous interactive element    | All views        |
| `Home`          | Reset to Building View (zoom level 0)         | All views        |
| `1`-`9`, `0`    | Jump to floor 1-9, 10                         | Building View    |
| `Shift+1`-`0`   | Jump to floor 11-20                           | Building View    |
| `B`             | Toggle bookmark on current view               | All views        |
| `?`             | Show keyboard shortcuts overlay               | All views        |
| `Space`         | Open/activate selected element                | When selected    |
| `Ctrl+C`        | Copy current view URL to clipboard            | All views        |

**Keyboard Handler Implementation:**

```javascript
class KeyboardController {
  constructor(navigationState, viewport) {
    this._navState = navigationState;
    this._viewport = viewport;
    this._shortcutsOverlayVisible = false;

    document.addEventListener('keydown', (e) => this._handleKey(e));
  }

  _handleKey(e) {
    // Ignore if focus is in an input field
    if (['INPUT', 'TEXTAREA', 'SELECT'].includes(document.activeElement?.tagName)) {
      return;
    }

    const PAN_STEP = 50;
    const ZOOM_IN_FACTOR = 1.25;
    const ZOOM_OUT_FACTOR = 0.8;
    const ZOOM_KEY_DURATION = 150;

    switch (e.key) {
      case 'ArrowUp':
        e.preventDefault();
        this._viewport.pan(0, PAN_STEP);
        break;
      case 'ArrowDown':
        e.preventDefault();
        this._viewport.pan(0, -PAN_STEP);
        break;
      case 'ArrowLeft':
        e.preventDefault();
        this._viewport.pan(PAN_STEP, 0);
        break;
      case 'ArrowRight':
        e.preventDefault();
        this._viewport.pan(-PAN_STEP, 0);
        break;
      case '+':
      case '=':
        e.preventDefault();
        this._viewport.zoomAnimated(
          this._viewport.currentScale * ZOOM_IN_FACTOR,
          this._viewport.currentCenter,
          ZOOM_KEY_DURATION
        );
        break;
      case '-':
        e.preventDefault();
        this._viewport.zoomAnimated(
          this._viewport.currentScale * ZOOM_OUT_FACTOR,
          this._viewport.currentCenter,
          ZOOM_KEY_DURATION
        );
        break;
      case 'Escape':
        e.preventDefault();
        this._navState.navigateUp();
        break;
      case 'Enter':
      case ' ':
        e.preventDefault();
        this._activateSelected();
        break;
      case 'Home':
        e.preventDefault();
        this._navState.navigateToBuilding();
        break;
      case 'Tab':
        // Allow default tab behavior but within HIC focus trap
        this._handleTabNavigation(e);
        break;
      case '?':
        e.preventDefault();
        this._toggleShortcutsOverlay();
        break;
      default:
        // Floor number shortcuts
        if (this._navState.state.current_view === 'building') {
          const floorNum = this._keyToFloorNumber(e.key, e.shiftKey);
          if (floorNum !== null) {
            e.preventDefault();
            this._navState.navigateToFloor(floorNum, FLOOR_REGISTRY[floorNum].name);
          }
        }
    }
  }

  _keyToFloorNumber(key, shift) {
    const digit = parseInt(key, 10);
    if (isNaN(digit)) return null;
    const floor = digit === 0 ? 10 : digit;
    return shift ? floor + 10 : floor;
  }

  _activateSelected() {
    const state = this._navState.state;
    if (state.selected_object_id) {
      this._navState.openObject(state.selected_object_id);
    } else if (state.room_id && state.current_view === 'floor') {
      this._navState.navigateToRoom(state.room_id, ROOM_REGISTRY_LOOKUP[state.room_id].name);
    }
  }

  _handleTabNavigation(e) {
    // Custom focus management within the HIC SVG
    // Elements are ordered: floors (building view), rooms (floor view), objects (room view)
    const focusableElements = this._getFocusableElements();
    const currentIndex = focusableElements.indexOf(document.activeElement);

    if (e.shiftKey) {
      const prev = currentIndex > 0 ? currentIndex - 1 : focusableElements.length - 1;
      focusableElements[prev].focus();
    } else {
      const next = currentIndex < focusableElements.length - 1 ? currentIndex + 1 : 0;
      focusableElements[next].focus();
    }
    e.preventDefault();
  }

  _getFocusableElements() {
    return Array.from(
      document.querySelectorAll('[data-hic-focusable="true"]')
    ).sort((a, b) => {
      return (parseInt(a.dataset.hicFocusOrder) || 0)
        - (parseInt(b.dataset.hicFocusOrder) || 0);
    });
  }
}
```

**Focus Indicator Styling:**

```css
[data-hic-focusable="true"]:focus {
  outline: none; /* remove default */
}

[data-hic-focusable="true"]:focus-visible {
  outline: 2px solid #FFFFFF;
  outline-offset: 3px;
  filter: drop-shadow(0 0 6px #FFFFFF) drop-shadow(0 0 12px #00FFCC);
  transition: filter 150ms ease-in-out, outline-offset 150ms ease-in-out;
}
```

---

### 12.10 Breadcrumb Trail

The breadcrumb trail provides persistent orientation within the HIC hierarchy. It is rendered as an overlay bar anchored to the top-left of the viewport, outside the SVG transform matrix (so it does not zoom or pan with the content).

**Visual Design:**

```
+============================================================================+
|  BUILDING  >  FLOOR 7: INTEL OPS  >  THREAT ARCHIVE                       |
+============================================================================+

  - Background: rgba(10, 10, 30, 0.85)
  - Border bottom: 1px solid rgba(0, 255, 204, 0.3)
  - Font: JetBrains Mono, 12px, letter-spacing 1px
  - Inactive crumbs: #808099 (clickable, underline on hover)
  - Active crumb (last): #00FFCC, no underline, font-weight 600
  - Separator: ">" in #404055
  - Padding: 8px 16px
  - Height: 36px
  - Position: fixed top-left, z-index 8000
```

**Breadcrumb Rendering:**

```javascript
class BreadcrumbRenderer {
  constructor(container, navigationState) {
    this._container = container;
    this._navState = navigationState;

    navigationState.subscribe((state) => this._render(state));
    this._render(navigationState.state);
  }

  _render(state) {
    this._container.innerHTML = '';

    state.breadcrumb_trail.forEach((crumb, index) => {
      const isLast = index === state.breadcrumb_trail.length - 1;

      // Separator
      if (index > 0) {
        const sep = document.createElement('span');
        sep.className = 'hic-breadcrumb__separator';
        sep.textContent = '>';
        sep.setAttribute('aria-hidden', 'true');
        this._container.appendChild(sep);
      }

      // Crumb
      const el = document.createElement(isLast ? 'span' : 'button');
      el.className = `hic-breadcrumb__item ${isLast ? 'hic-breadcrumb__item--active' : ''}`;
      el.textContent = crumb.label.toUpperCase();

      if (!isLast) {
        el.setAttribute('aria-label', `Navigate to ${crumb.label}`);
        el.addEventListener('click', () => {
          this._navState.navigateToBreadcrumb(index);
        });
      } else {
        el.setAttribute('aria-current', 'page');
      }

      this._container.appendChild(el);
    });
  }
}
```

**Breadcrumb CSS:**

```css
.hic-breadcrumb {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 36px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  gap: 8px;
  background: rgba(10, 10, 30, 0.85);
  border-bottom: 1px solid rgba(0, 255, 204, 0.3);
  backdrop-filter: blur(8px);
  z-index: 8000;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 12px;
  letter-spacing: 1px;
}

.hic-breadcrumb__item {
  color: #808099;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 3px;
  font-family: inherit;
  font-size: inherit;
  letter-spacing: inherit;
  transition: color 150ms, background 150ms;
}

.hic-breadcrumb__item:hover {
  color: #E0E0FF;
  background: rgba(0, 255, 204, 0.1);
  text-decoration: underline;
  text-underline-offset: 3px;
}

.hic-breadcrumb__item--active {
  color: #00FFCC;
  font-weight: 600;
  cursor: default;
  text-shadow: 0 0 8px rgba(0, 255, 204, 0.5);
}

.hic-breadcrumb__item--active:hover {
  background: none;
  text-decoration: none;
}

.hic-breadcrumb__separator {
  color: #404055;
  font-size: 10px;
  user-select: none;
}
```

**Breadcrumb Navigation State Machine:**

```
BREADCRUMB NAVIGATION STATE FLOW:

  Click "BUILDING"         Click "FLOOR 7"          (current)
  +-----------+            +-----------+            +-----------+
  | BUILDING  | ---------> | FLOOR 7:  | ---------> | THREAT    |
  |           |    >       | INTEL OPS |    >       | ARCHIVE   |
  +-----------+            +-----------+            +-----------+
       ^                        ^                        |
       |                        |                        |
       |  Clicking "BUILDING"   |  Clicking "FLOOR 7"   |
       |  from any depth:       |  from room depth:      |
       |                        |                        |
       |  - Pops all crumbs     |  - Pops room crumb     |
       |    except first        |  - Resets to floor view |
       |  - Resets to 0.1x zoom |  - Resets to 1.0x zoom |
       |  - Clears floor/room   |  - Clears room_id      |
       |  - Pushes history      |  - Pushes history       |
       +------------------------+------------------------+
```

---

## CROSS-AGENT INTEGRATION NOTES

### Interaction Priority Resolution

When multiple interaction systems could respond to a single user action, the following priority order resolves conflicts:

```
PRIORITY (highest first):
  1. Context menu (right-click) -- always takes precedence
  2. Double-click -- takes precedence over single click
  3. Click on interactive element -- takes precedence over pan
  4. Pan (click-drag on empty space)
  5. Hover -- lowest priority, never blocks other interactions
```

### Event Flow Diagram

```
User Input Event
       |
       v
  +---------+    right-click    +----------------+
  | Pointer |  --------------> | Context Menu   |
  | Event   |                  | Controller     |
  | Router  |                  +----------------+
  |         |
  |         |    double-click   +----------------+
  |         |  --------------> | Click          |
  |         |    single-click  | Controller     |
  |         |  --------------> | (Agent 10)     |
  |         |                  +--------+-------+
  |         |                           |
  |         |                  +--------v-------+
  |         |                  | Navigation     |
  |         |                  | State          |
  |         |                  | Controller     |
  |         |                  | (Agent 12)     |
  |         |                  +--------+-------+
  |         |                           |
  |         |    drag            +------v--------+
  |         |  ------------->   | Pan/Zoom      |
  |         |    scroll/pinch   | Controller    |
  |         |  ------------->   | (Agent 12)    |
  |         |                   +---------------+
  |         |
  |         |    move (no btn)  +----------------+
  |         |  --------------> | Hover          |
  |         |                  | Controller     |
  |         |                  | (Agent 11)     |
  +---------+                  +----------------+
```

### Accessibility Summary

All three agents contribute to WCAG 2.1 AA compliance:

- **Agent 10 (Click):** All click targets meet 44x44px minimum, keyboard Enter/Space equivalents for all click actions, focus management on view transitions.
- **Agent 11 (Hover):** Tooltips use `role="tooltip"` and `aria-live="polite"`, hover information is also available via keyboard focus (Tab to element shows tooltip), color is not the only indicator of status.
- **Agent 12 (Navigation):** Full keyboard navigation (Tab, arrows, shortcuts), breadcrumb uses `aria-current="page"`, URL hash enables bookmarking and sharing, browser back/forward work as expected, `prefers-reduced-motion` media query disables all animations when set.

**Reduced Motion Support:**

```css
@media (prefers-reduced-motion: reduce) {
  .hic-floor-strip,
  .hic-room-polygon,
  .hic-room-object,
  .hic-tooltip,
  .hic-breadcrumb__item {
    transition-duration: 0ms !important;
    animation-duration: 0ms !important;
  }
}
```

---

> **End of Stage 4: HIC Interaction Mechanics -- Agents 10-12**
>
> This specification provides implementation-ready details for all user interactions with the Holm Intelligence Complex interface. All timing values, coordinate systems, state machines, and code samples are production-ready and should be used as the authoritative reference for the HIC interaction layer.
