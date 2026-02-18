# HIC Interaction Protocols

## Navigation, Input Handling & User Experience

> **Stage 4 Specification** -- Holm Intelligence Complex
> Defines how users move through, interact with, and experience the HIC skyscraper.
> Every UX primitive is mapped to a spatial metaphor within the building.

---

## Table of Contents

1. [Design Philosophy](#design-philosophy)
2. [Navigation Paradigms](#navigation-paradigms)
   - [Street Level View](#street-level-view)
   - [Lobby Entry](#lobby-entry)
   - [Floor Navigation](#floor-navigation)
   - [Room Entry](#room-entry)
   - [Deep Dive](#deep-dive)
   - [Express Elevator](#express-elevator)
   - [Emergency Exit](#emergency-exit)
3. [Click Interactions](#click-interactions)
   - [Single Click: Opening Doors](#single-click-opening-doors)
   - [Floor Marker Clicks](#floor-marker-clicks)
   - [Room Label Clicks](#room-label-clicks)
   - [Neon Sign Activation](#neon-sign-activation)
   - [Window Click: Adjacent Preview](#window-click-adjacent-preview)
4. [Hover Behaviors](#hover-behaviors)
   - [Door Hover: Frosted Glass Preview](#door-hover-frosted-glass-preview)
   - [Neon Sign Hover](#neon-sign-hover)
   - [Floor Plan Hover](#floor-plan-hover)
   - [Cross-Reference Hover](#cross-reference-hover)
   - [Proximity Glow](#proximity-glow)
5. [Zoom and Pan](#zoom-and-pan)
   - [Zoom Out: Building Facade](#zoom-out-building-facade)
   - [Zoom In: Room Detail](#zoom-in-room-detail)
   - [Pan: Scrolling Through Floors](#pan-scrolling-through-floors)
   - [Pinch Zoom: Mobile Touch](#pinch-zoom-mobile-touch)
   - [Scroll Speed Zones](#scroll-speed-zones)
6. [Keyboard Navigation](#keyboard-navigation)
   - [Directional Movement](#directional-movement)
   - [Floor Traversal](#floor-traversal)
   - [Home Key: Lobby Return](#home-key-lobby-return)
   - [Tab Cycling](#tab-cycling)
   - [Slash: Building Intercom](#slash-building-intercom)
   - [Escape: Overlay Dismissal](#escape-overlay-dismissal)
   - [Full Keyboard Map](#full-keyboard-map)
7. [Search as Intercom System](#search-as-intercom-system)
   - [Intercom Architecture](#intercom-architecture)
   - [Voice Query: Search Input](#voice-query-search-input)
   - [Intercom Response: Result Rendering](#intercom-response-result-rendering)
   - [Floor-Specific Search](#floor-specific-search)
   - [Building-Wide Search](#building-wide-search)
   - [Recent Announcements: Search History](#recent-announcements-search-history)
8. [Mobile Experience](#mobile-experience)
   - [Touch Gesture Mapping](#touch-gesture-mapping)
   - [Swipe Navigation](#swipe-navigation)
   - [Long Press Preview](#long-press-preview)
   - [Two-Finger Scroll](#two-finger-scroll)
   - [Responsive Layout: Adaptive Building](#responsive-layout-adaptive-building)
   - [Portrait Mode: Single-Room View](#portrait-mode-single-room-view)
   - [Landscape Mode: Floor Plan Overview](#landscape-mode-floor-plan-overview)
9. [Accessibility](#accessibility)
   - [Screen Reader: Audio Descriptions](#screen-reader-audio-descriptions)
   - [High Contrast: Emergency Lighting](#high-contrast-emergency-lighting)
   - [Keyboard-Only Navigation](#keyboard-only-navigation)
   - [Reduced Motion: Static Signage](#reduced-motion-static-signage)
   - [Focus Indicators: Highlighted Doorways](#focus-indicators-highlighted-doorways)
10. [Loading States](#loading-states)
    - [Building Under Construction](#building-under-construction)
    - [Elevator in Transit](#elevator-in-transit)
    - [Room Furnishing](#room-furnishing)
    - [Power Outage](#power-outage)
11. [Performance](#performance)
    - [Lazy Loading: Render on Approach](#lazy-loading-render-on-approach)
    - [Prefetching: Adjacent Pre-Load](#prefetching-adjacent-pre-load)
    - [Caching: Furnished Rooms](#caching-furnished-rooms)
    - [Offline Mode: Emergency Generator](#offline-mode-emergency-generator)
12. [Animation & Transition Specifications](#animation--transition-specifications)
13. [State Management](#state-management)
14. [Event Architecture](#event-architecture)

---

## Design Philosophy

The HIC skyscraper is not a metaphor layered on top of a website. It **is** the interface. Every interaction a user performs corresponds to a physical action within the building. A click is a hand on a door handle. A hover is a glance through glass. A scroll is footsteps down a corridor. When the metaphor breaks, the interface breaks.

Three principles govern all interaction design:

1. **Spatial Consistency** -- Users should always know where they are in the building. If they are on Floor 7, Room 3, the interface must communicate that unambiguously through visual, textual, and structural cues.

2. **Physical Plausibility** -- Interactions must feel like they could happen in a physical building. You do not teleport through walls. You open doors. You ride elevators. You look through windows. Transitions respect the spatial graph.

3. **Neon Legibility** -- The cyberpunk aesthetic serves function, not decoration. Neon signage exists to guide attention, not to overwhelm it. Glow communicates interactivity. Brightness communicates focus. Darkness communicates the periphery.

```
Metaphor Mapping (Summary)
============================
Website Concept        | HIC Building Equivalent
-----------------------|--------------------------
Homepage               | Street-level exterior view
Domain/section index   | Floor lobby
Article list           | Floor corridor with doors
Article                | Room
Content section        | Area within a room
Link                   | Door / corridor connection
Navigation bar         | Express elevator panel
Breadcrumb trail       | Emergency exit signage
Search                 | Building intercom
Loading spinner        | Construction scaffolding
Error page             | Power outage / blocked room
Hover tooltip          | Frosted glass preview
Modal / overlay        | Security checkpoint window
Sidebar                | Floor directory plaque
Footer                 | Building foundation / basement
```

---

## Navigation Paradigms

Navigation through the HIC follows a strict spatial hierarchy. Users move from the outside in, from the macro to the micro, from the street to the room. Each level of navigation has its own interaction rules, transition animations, and visual identity.

### Street Level View

**Maps to:** Homepage / site index
**Metaphor:** Standing on the sidewalk, looking up at the building

The Street Level View is the first thing a user sees. The full facade of the HIC skyscraper is visible -- every illuminated floor, every neon sign, the full vertical extent of the documentation institution. This is the orientation point. Users have never entered the building. They are deciding where to go.

**Visual composition:**
- The building occupies the center of the viewport, rendered in forced perspective
- Each floor is labeled with its domain name in neon lettering along the facade
- Active floors (recently updated) pulse with brighter illumination
- The entrance at ground level glows with an inviting warmth, distinct from the cooler neon of upper floors
- A subtle rain/particle effect reinforces the cyberpunk atmosphere without obscuring content

**Interactions available at street level:**
- Click any floor label on the facade to ride the elevator directly to that floor (domain)
- Click the entrance to enter the lobby (same as clicking the primary domain)
- Scroll down to see the building foundation (footer content, site metadata)
- Scroll up is blocked -- you cannot go above the rooftop
- Hover over any floor to see it illuminate and display a brief description of that domain

**Technical implementation:**
```css
.hic-street-level {
  perspective: 1200px;
  perspective-origin: 50% 85%;
  overflow-y: auto;
  overflow-x: hidden;
}

.hic-facade-floor {
  transform: rotateX(2deg);
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
}

.hic-facade-floor:hover {
  box-shadow: 0 0 30px var(--neon-domain-color),
              0 0 60px var(--neon-domain-color-dim);
}
```

**State:**
```
location: { level: "street", floor: null, room: null }
history: []
```

---

### Lobby Entry

**Maps to:** Landing on a domain index page
**Metaphor:** Walking through the front doors onto a specific floor

When a user selects a domain, they enter the lobby of that floor. The lobby is the domain index -- it shows the layout of the entire floor, lists all available rooms (articles), and provides wayfinding signage.

**Visual composition:**
- The perspective shifts from exterior to interior
- A floor plan is rendered showing room layout -- each room is a card representing an article
- A directory plaque near the elevator shows all rooms on this floor, sorted by relevance or recency
- The floor's signature neon color dominates the palette
- Corridor lines connect related rooms, visible as glowing pathways on the floor plan

**Interactions available in the lobby:**
- Click any room on the floor plan to enter it (navigate to article)
- Hover over a room to see a frosted-glass preview (article summary, word count, last updated)
- Click the elevator panel to move to a different floor (domain nav)
- Click the exit sign to return to street level (homepage)
- Scroll to pan across the floor plan if the floor has many rooms

**Transition from street to lobby:**
```
Duration: 400ms
Easing: cubic-bezier(0.4, 0.0, 0.2, 1)
Animation sequence:
  1. Facade zooms into the selected floor (scale 1.0 -> 3.0, opacity 1.0 -> 0.0)
  2. Interior fades in from black (opacity 0.0 -> 1.0)
  3. Floor plan elements stagger in from left to right (translateX -20px -> 0, 50ms stagger)
  4. Neon signs flicker on (opacity keyframe: 0, 0.8, 0, 1.0 over 200ms)
```

**State:**
```
location: { level: "lobby", floor: "api-reference", room: null }
history: [{ level: "street" }]
```

---

### Floor Navigation

**Maps to:** Moving between articles within a domain
**Metaphor:** Walking the corridors of a floor, passing doors

Once inside a floor, the user can navigate laterally between rooms. This is corridor navigation -- the user is walking past doors, reading room labels, deciding where to go next. Floor navigation never changes the floor; it moves horizontally through the content on a single level.

**Corridor rendering:**
- Rooms appear as doorways along a corridor
- Each door has a nameplate (article title) and a status indicator (new, updated, deprecated)
- The current room's door is open and illuminated
- Adjacent rooms have subtle light spilling from under their doors
- The corridor has depth -- rooms further from the current position recede visually

**Navigation methods within a floor:**
- Click a door to enter a room (article)
- Use left/right arrow keys to move to the adjacent room
- Click "Next" / "Previous" buttons styled as corridor direction signs
- Use the floor directory (sidebar) to jump to any room on the floor
- Breadcrumb trail shows: Street > Floor Name > Current Room

**Corridor scroll behavior:**
```javascript
// Corridor uses horizontal scroll with snap points
const corridorConfig = {
  scrollDirection: 'horizontal',
  snapType: 'x mandatory',
  snapAlign: 'center',
  scrollPadding: '0 20%',
  overscrollBehavior: 'contain', // Do not bleed into floor changes
};
```

**State:**
```
location: { level: "corridor", floor: "api-reference", room: null, corridor_position: 4 }
history: [{ level: "street" }, { level: "lobby", floor: "api-reference" }]
```

---

### Room Entry

**Maps to:** Opening a specific article
**Metaphor:** Stepping through a doorway into a room

Entering a room is the primary content-viewing state. The room contains the article's full text, code examples, diagrams, and interactive elements. The room's decor reflects the content type -- tutorial rooms have workbenches, reference rooms have filing cabinets, guide rooms have whiteboards.

**Visual composition:**
- The door opens with a brief animation (swing or slide, depending on room type)
- The room interior fills the viewport
- Content is arranged spatially within the room:
  - Heading structure maps to areas within the room (sections are alcoves)
  - Code blocks are terminal screens mounted on walls
  - Images and diagrams are framed displays
  - Tables are workbenches with items laid out
  - Callouts and warnings are neon signs on the walls
- A room map (table of contents) appears as a wall-mounted floor plan of this room

**Interactions available inside a room:**
- Scroll to move through the room's content (the "Deep Dive," detailed below)
- Click internal links to open doors to other rooms
- Click code blocks to copy, expand, or interact with them
- Hover over glossary terms to see definitions (neon tooltip)
- Click the room map to jump to a specific section (scroll to heading)
- Click the door you entered from to return to the corridor

**Transition from corridor to room:**
```
Duration: 300ms
Easing: ease-out
Animation sequence:
  1. Door swings open (rotateY 0deg -> -90deg on the door element)
  2. Interior scales from 0.9 to 1.0 with slight blur clearing (blur 4px -> 0)
  3. Room lights turn on (background lightens, neon elements activate)
  4. Content fades in top-to-bottom (stagger 30ms per section)
```

**State:**
```
location: { level: "room", floor: "api-reference", room: "authentication", scroll_position: 0 }
history: [
  { level: "street" },
  { level: "lobby", floor: "api-reference" },
  { level: "corridor", floor: "api-reference", corridor_position: 4 }
]
```

---

### Deep Dive

**Maps to:** Scrolling through article content
**Metaphor:** Exploring a room's details -- walking around, examining items

The Deep Dive is the scroll experience within a room. As the user scrolls, they move deeper into the room's content. The spatial metaphor shifts from walking through a corridor to examining objects within a space.

**Scroll behavior inside a room:**
- Vertical scroll moves through content sections sequentially
- Scroll position is tracked and persisted -- if you leave and return, you are where you left off
- The room map (table of contents) highlights the current section as you scroll
- Heading elements trigger waypoint events that update the URL hash and room map
- Scroll depth is visualized as a vertical progress bar styled as a room depth meter

**Scroll-triggered events:**
```javascript
const deepDiveConfig = {
  // Sections reveal as you approach them
  revealThreshold: 0.15,       // 15% visible before reveal animation
  revealAnimation: 'fade-up',  // Slide up 20px and fade in
  revealDuration: 300,         // milliseconds
  revealStagger: 50,           // milliseconds between sibling elements

  // Heading waypoints
  waypointOffset: 80,          // pixels from top to trigger waypoint
  waypointCallback: (heading) => {
    updateRoomMap(heading.id);
    updateURLHash(heading.id);
    updateBreadcrumb(heading.textContent);
  },

  // Progress tracking
  progressSelector: '.hic-room-depth-meter',
  progressSmooth: true,
};
```

**Section types and their room-area mappings:**

| Content Element      | Room Area Metaphor        | Scroll Behavior                     |
|----------------------|---------------------------|--------------------------------------|
| H2 heading           | New alcove entrance       | Snaps to top with 80px offset        |
| H3 heading           | Sub-area within alcove    | Smooth scroll, no snap               |
| Code block           | Terminal screen on wall   | Sticky positioning while in view     |
| Blockquote           | Plaque on the wall        | Slight parallax effect               |
| Image / diagram      | Framed display            | Fade-in on scroll reveal             |
| Table                | Workbench layout          | Horizontal scroll if overflowing     |
| Warning callout      | Red neon sign             | Pulse glow on reveal                 |
| Info callout         | Blue neon sign            | Steady glow on reveal                |
| Tip callout          | Green neon sign           | Flicker-on effect on reveal          |

**State:**
```
location: {
  level: "room",
  floor: "api-reference",
  room: "authentication",
  scroll_position: 2340,
  active_section: "oauth-flow",
  depth_percentage: 0.42
}
```

---

### Express Elevator

**Maps to:** Using the main navigation bar to jump between domains
**Metaphor:** Riding the express elevator between floors

The Express Elevator is the primary cross-domain navigation mechanism. It appears as a persistent UI element -- a vertical panel styled as an elevator control panel. Each button on the panel represents a floor (domain). Pressing a button takes you directly to that floor's lobby, regardless of your current location.

**Elevator panel design:**
- Fixed position on the left side of the viewport (desktop) or bottom of the viewport (mobile)
- Each floor button shows the floor number, domain name, and domain icon
- The current floor's button is illuminated in the domain's neon color
- Floors above the current position point upward; floors below point downward
- An indicator light shows the elevator's current position

**Interactions:**
- Click a floor button to go directly to that floor's lobby
- Hover over a floor button to see the floor's description and room count
- The elevator panel can be collapsed/expanded with a toggle
- On mobile, the elevator panel is accessed via a floating action button

**Transition animation (elevator ride):**
```
Duration: 500ms + (50ms * floors_traveled)  // Longer rides for distant floors
Easing: cubic-bezier(0.4, 0.0, 0.6, 1)
Animation sequence:
  1. Current room/lobby fades out (200ms)
  2. Elevator shaft animation plays:
     - Floor indicators count up/down
     - Subtle vertical motion blur
     - Ambient sound cue (optional, off by default)
  3. Destination lobby fades in (200ms)
  4. Door opening animation (100ms)
```

**Elevator panel state management:**
```javascript
const elevatorPanel = {
  floors: [
    { id: 'getting-started', number: 1, label: 'Getting Started', color: '#00ff88' },
    { id: 'tutorials',       number: 2, label: 'Tutorials',       color: '#00ccff' },
    { id: 'api-reference',   number: 3, label: 'API Reference',   color: '#ff6600' },
    { id: 'guides',          number: 4, label: 'Guides',          color: '#ff00ff' },
    { id: 'concepts',        number: 5, label: 'Concepts',        color: '#ffcc00' },
    { id: 'changelog',       number: 6, label: 'Changelog',       color: '#cc00ff' },
  ],
  currentFloor: null,
  isCollapsed: false,
  isAnimating: false,

  goToFloor(floorId) {
    if (this.isAnimating) return;
    const target = this.floors.find(f => f.id === floorId);
    const distance = Math.abs(target.number - this.currentFloor);
    this.isAnimating = true;
    playElevatorTransition(distance).then(() => {
      this.currentFloor = target.number;
      this.isAnimating = false;
      navigateTo(`/${floorId}/`);
    });
  }
};
```

---

### Emergency Exit

**Maps to:** Back button, breadcrumb trail, escape routes
**Metaphor:** Emergency exit signage guiding you out of the building

Users must always be able to leave. The Emergency Exit system ensures that no matter how deep a user is in the HIC, they can retrace their steps or jump to safety. It combines the browser's back button behavior with an explicit breadcrumb trail rendered as emergency exit signage.

**Breadcrumb trail as exit signs:**
- Rendered as a horizontal bar at the top of the viewport
- Styled as illuminated green exit signs (universal emergency signage)
- Each segment shows the path taken: Street > Floor Name > Room Name > Section
- Clicking any segment navigates back to that point
- The trail never exceeds 4 segments (Street / Floor / Room / Section)

**Breadcrumb rendering:**
```html
<nav class="hic-exit-signs" aria-label="Breadcrumb navigation">
  <ol>
    <li>
      <a href="/" class="hic-exit-sign hic-exit-sign--street">
        <span class="hic-exit-icon">EXIT</span>
        <span class="hic-exit-label">Street Level</span>
      </a>
    </li>
    <li>
      <a href="/api-reference/" class="hic-exit-sign hic-exit-sign--floor">
        <span class="hic-exit-icon">&laquo;</span>
        <span class="hic-exit-label">API Reference</span>
      </a>
    </li>
    <li aria-current="page">
      <span class="hic-exit-sign hic-exit-sign--room hic-exit-sign--current">
        <span class="hic-exit-label">Authentication</span>
      </span>
    </li>
  </ol>
</nav>
```

**Back button behavior:**
- The browser back button moves one step back in the navigation history
- If the user entered via a direct link (no history), back goes to the floor lobby
- If the user is at the lobby, back goes to street level
- If the user is at street level, back exits the site (default browser behavior)
- History state is managed with `pushState` to ensure clean back/forward traversal

**Escape hatch hierarchy:**
```
Priority 1: Browser back button -- always works, respects history stack
Priority 2: Breadcrumb click -- explicit navigation to any ancestor level
Priority 3: Escape key -- closes the current overlay or modal, then returns to parent
Priority 4: Logo click -- always returns to street level (homepage)
Priority 5: 404 handler -- if lost, show a "you are here" map and suggest exits
```

---

## Click Interactions

Clicks are the primary interaction verb in the HIC. Every clickable element corresponds to a physical object in the building that you would interact with by touch -- doors, buttons, signs, switches.

### Single Click: Opening Doors

A single click on a navigational element is equivalent to opening a door. The target determines which door you are opening and what lies beyond it.

**Link types and their door equivalents:**

| Link Type             | Door Metaphor              | Click Behavior                              |
|-----------------------|----------------------------|----------------------------------------------|
| Internal article link | Wooden corridor door       | Navigate to article with door-open animation |
| External link         | Fire escape / exterior door| Open in new tab, show "leaving building" cue |
| Anchor link           | Interior room door         | Smooth scroll to section within current room |
| Download link         | Supply closet              | Trigger download with "item retrieved" cue   |
| API endpoint link     | Server room access panel   | Open in interactive API explorer             |
| Email / contact link  | Intercom handset           | Open default mail client                     |

**Click feedback:**
```css
.hic-door {
  cursor: pointer;
  position: relative;
}

.hic-door::after {
  content: '';
  position: absolute;
  inset: -4px;
  border: 2px solid transparent;
  border-radius: 4px;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.hic-door:active::after {
  border-color: var(--neon-accent);
  box-shadow: 0 0 12px var(--neon-accent), inset 0 0 12px var(--neon-accent-dim);
}
```

**Click event handling:**
```javascript
function handleDoorClick(event, doorElement) {
  const target = doorElement.dataset.target;
  const doorType = doorElement.dataset.doorType;

  // Prevent double-clicks from double-navigating
  if (doorElement.classList.contains('hic-door--opening')) return;
  doorElement.classList.add('hic-door--opening');

  // Play door animation
  const animation = getDoorAnimation(doorType);
  animation.play();

  // Navigate after animation completes
  animation.onfinish = () => {
    if (doorType === 'external') {
      window.open(target, '_blank', 'noopener');
      doorElement.classList.remove('hic-door--opening');
    } else {
      navigateToRoom(target);
    }
  };
}
```

---

### Floor Marker Clicks

Floor markers are the primary mechanism for inter-domain navigation from within the building. They appear as illuminated floor indicators in the elevator panel and as floor labels on the building facade.

**Floor marker locations:**
- Elevator panel (persistent navigation)
- Building facade (street level only)
- Floor directory plaques (lobby view)
- Cross-references in content that point to other domains

**Click behavior:**
1. User clicks a floor marker
2. Current room/lobby state is saved to session storage
3. Elevator transition animation begins
4. Destination floor lobby loads
5. Elevator doors open to reveal the new floor

**Floor marker element:**
```javascript
class FloorMarker extends HTMLElement {
  connectedCallback() {
    this.addEventListener('click', () => this.activateElevator());
    this.addEventListener('keydown', (e) => {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        this.activateElevator();
      }
    });
  }

  activateElevator() {
    const floorId = this.getAttribute('floor');
    const currentFloor = document.querySelector('.hic-floor--current');

    if (currentFloor?.id === floorId) return; // Already on this floor

    saveRoomState();
    elevatorPanel.goToFloor(floorId);
  }
}

customElements.define('hic-floor-marker', FloorMarker);
```

---

### Room Label Clicks

Room labels are the article-level navigation elements. They appear as nameplates on doors in the corridor view and as entries in the floor directory.

**Nameplate design:**
- Room title in the domain's neon color
- Subtitle showing article type (tutorial, reference, guide)
- Status badge: a small indicator showing freshness (new, updated, stable, deprecated)
- Word count / reading time estimate displayed as room size

**Click behavior:**
1. Nameplate illuminates with a brief flash
2. Door-open animation plays (300ms)
3. Room content loads
4. Room interior renders with section-by-section reveal

**Nameplate hierarchy in the corridor:**
```
[STATUS] ROOM TITLE                    [SIZE]
         Article subtitle / description
         Tags: tag1, tag2, tag3
```

---

### Neon Sign Activation

Neon signs are interactive elements that are not navigational links. They trigger actions, toggle states, open modals, or activate tools. Examples include code copy buttons, theme toggles, language selectors, and feedback forms.

**Neon sign types:**

| Sign Type         | Action                          | Visual                           |
|-------------------|---------------------------------|-----------------------------------|
| Copy button       | Copy code to clipboard          | Clipboard icon, flash on success  |
| Theme toggle      | Switch light/dark mode          | Sun/moon icon, neon color shift   |
| Language selector | Change code language            | Globe icon, dropdown overlay      |
| Version selector  | Switch documentation version    | Stack icon, version list overlay  |
| Feedback button   | Open feedback form              | Speech bubble icon, modal overlay |
| Expand/collapse   | Toggle content section          | Chevron icon, section animation   |

**Click behavior for neon signs:**
```javascript
function handleNeonSignClick(sign) {
  // Flash effect
  sign.classList.add('hic-neon--flash');
  setTimeout(() => sign.classList.remove('hic-neon--flash'), 200);

  // Execute action
  const action = sign.dataset.action;
  switch (action) {
    case 'copy':
      copyToClipboard(sign.dataset.content);
      showNotification('Copied to clipboard', 'success');
      break;
    case 'toggle-theme':
      toggleBuildingLighting();
      break;
    case 'expand':
      toggleSection(sign.dataset.target);
      break;
    default:
      console.warn(`Unknown neon sign action: ${action}`);
  }
}
```

**Neon flash animation:**
```css
@keyframes neon-flash {
  0%   { filter: brightness(1.0); }
  20%  { filter: brightness(2.5); text-shadow: 0 0 20px currentColor; }
  40%  { filter: brightness(0.8); }
  60%  { filter: brightness(1.8); text-shadow: 0 0 10px currentColor; }
  100% { filter: brightness(1.0); }
}

.hic-neon--flash {
  animation: neon-flash 200ms ease-out;
}
```

---

### Window Click: Adjacent Preview

Windows allow users to preview content in adjacent rooms without fully entering them. Clicking a window opens a preview card -- a constrained view of the target content that can be dismissed or expanded into a full room entry.

**Window locations:**
- Inline links within article content
- "See also" references at the bottom of a room
- Related articles in the floor directory
- Cross-domain references

**Preview card specifications:**
```javascript
const previewCardConfig = {
  width: 400,             // pixels
  maxHeight: 300,         // pixels
  offset: { x: 0, y: 8 },// offset from trigger element
  showDelay: 0,           // show immediately on click
  hideDelay: 200,         // 200ms delay before hiding on mouse leave
  position: 'auto',       // auto-position to stay in viewport
  content: {
    title: true,          // Show article title
    excerpt: true,        // Show first paragraph
    meta: true,           // Show reading time, last updated
    thumbnail: false,     // No thumbnail in preview
    actions: ['open', 'open-new-tab'], // Available actions
  },
};
```

**Preview card rendering:**
```html
<div class="hic-window-preview" role="dialog" aria-label="Room preview">
  <div class="hic-window-preview__glass">
    <h3 class="hic-window-preview__title">Authentication Guide</h3>
    <p class="hic-window-preview__excerpt">
      Learn how to authenticate API requests using OAuth 2.0,
      API keys, or session tokens...
    </p>
    <div class="hic-window-preview__meta">
      <span>5 min read</span>
      <span>Updated 2 days ago</span>
    </div>
    <div class="hic-window-preview__actions">
      <button class="hic-neon-sign" data-action="navigate">Enter Room</button>
      <button class="hic-neon-sign" data-action="new-tab">Open in New Window</button>
    </div>
  </div>
</div>
```

---

## Hover Behaviors

Hover interactions provide ambient information without requiring a click. In the HIC, hovering is equivalent to looking at something -- you get visual feedback and information without committing to an action.

### Door Hover: Frosted Glass Preview

When the cursor hovers over a door (link to another article), the door's frosted glass panel clears slightly to reveal a preview of the room beyond.

**Preview content shown on hover:**
- Article title (large, neon-colored)
- First 2-3 sentences of the article
- Reading time estimate
- Last updated timestamp
- Tags / categories

**Hover timing:**
```javascript
const doorHoverConfig = {
  showDelay: 300,    // Wait 300ms before showing preview (avoid accidental triggers)
  hideDelay: 150,    // Wait 150ms before hiding (allow cursor to move to preview)
  fadeIn: 200,       // Fade-in duration
  fadeOut: 100,       // Fade-out duration
};
```

**Frosted glass effect:**
```css
.hic-door__glass {
  backdrop-filter: blur(12px) saturate(0.5);
  background: rgba(0, 0, 0, 0.6);
  transition: backdrop-filter 0.3s ease, background 0.3s ease;
}

.hic-door:hover .hic-door__glass {
  backdrop-filter: blur(2px) saturate(1.0);
  background: rgba(0, 0, 0, 0.3);
}
```

**Implementation notes:**
- Preview content is fetched via a lightweight API endpoint that returns only title, excerpt, and metadata
- Previews are cached in a `Map` keyed by article URL -- once fetched, they are never re-fetched during the session
- If the preview fails to load within 500ms, the frosted glass effect still plays but no text appears
- On touch devices, door hover is replaced by long-press (see Mobile Experience)

---

### Neon Sign Hover

Hovering over a neon sign brightens it and reveals additional information about what the sign does.

**Hover effects:**
1. **Brightness increase**: The sign's glow intensifies from 100% to 150% brightness
2. **Tooltip reveal**: A small label appears below the sign explaining its function
3. **Cursor change**: Cursor changes to `pointer` to indicate interactivity
4. **Ambient glow**: A soft light pool appears beneath the sign, illuminating nearby elements

```css
.hic-neon-sign {
  filter: brightness(1.0);
  transition: filter 0.2s ease;
}

.hic-neon-sign:hover {
  filter: brightness(1.5);
}

.hic-neon-sign:hover::after {
  content: attr(data-tooltip);
  position: absolute;
  bottom: -28px;
  left: 50%;
  transform: translateX(-50%);
  padding: 4px 8px;
  background: rgba(0, 0, 0, 0.85);
  color: var(--neon-text);
  font-size: 12px;
  border-radius: 4px;
  white-space: nowrap;
  border: 1px solid var(--neon-accent-dim);
  opacity: 0;
  animation: tooltip-reveal 0.2s ease 0.3s forwards;
}

@keyframes tooltip-reveal {
  to { opacity: 1; }
}
```

---

### Floor Plan Hover

Hovering over a room on the floor plan (lobby view) highlights that room and shows its connections to other rooms.

**Hover behavior in the floor plan:**
1. The hovered room's outline brightens to full neon intensity
2. Connection lines from this room to related rooms illuminate
3. Connected rooms receive a dim highlight (secondary glow)
4. The room's metadata appears in a floating tooltip
5. Other, unconnected rooms dim slightly to reduce visual noise

**Connection line rendering:**
```javascript
function highlightConnections(roomElement) {
  const roomId = roomElement.dataset.roomId;
  const connections = roomGraph.getConnections(roomId);

  // Brighten the hovered room
  roomElement.style.setProperty('--room-glow', '1.0');

  // Illuminate connection lines
  connections.forEach(conn => {
    const line = document.querySelector(`[data-connection="${roomId}-${conn.target}"]`);
    if (line) {
      line.style.opacity = '1.0';
      line.style.filter = 'brightness(1.5)';
    }

    // Secondary glow on connected rooms
    const connectedRoom = document.querySelector(`[data-room-id="${conn.target}"]`);
    if (connectedRoom) {
      connectedRoom.style.setProperty('--room-glow', '0.4');
    }
  });

  // Dim unconnected rooms
  document.querySelectorAll('.hic-room-card').forEach(room => {
    if (room !== roomElement && !connections.some(c => c.target === room.dataset.roomId)) {
      room.style.setProperty('--room-glow', '0.1');
    }
  });
}
```

---

### Cross-Reference Hover

Cross-references are inline links within article content that point to other articles (on the same floor or different floors). Hovering over them shows a preview of the target.

**Implementation:**
- Cross-references use the same preview card system as Window Clicks
- The preview card appears on hover after a 400ms delay
- If the cross-reference points to a different floor, the preview card shows a floor badge
- The preview card dismisses on mouse leave with a 200ms grace period

**Cross-reference element:**
```html
<a href="/guides/authentication/"
   class="hic-cross-ref"
   data-floor="guides"
   data-room="authentication"
   aria-describedby="preview-auth">
  authentication guide
</a>
```

---

### Proximity Glow

The Proximity Glow is a subtle ambient effect where interactive elements illuminate as the cursor approaches them, even before a hover state is triggered. This creates a feeling of the building responding to the user's presence.

**Implementation approach:**
```javascript
class ProximityGlow {
  constructor(element, options = {}) {
    this.element = element;
    this.maxDistance = options.maxDistance || 200; // pixels
    this.maxGlow = options.maxGlow || 0.5;       // opacity
    this.glowColor = options.glowColor || 'var(--neon-accent)';

    document.addEventListener('mousemove', (e) => this.onMouseMove(e));
  }

  onMouseMove(event) {
    const rect = this.element.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;

    const distance = Math.hypot(event.clientX - centerX, event.clientY - centerY);

    if (distance < this.maxDistance) {
      const intensity = 1 - (distance / this.maxDistance);
      const glow = intensity * this.maxGlow;
      this.element.style.boxShadow = `0 0 ${20 * intensity}px ${this.glowColor}`;
      this.element.style.setProperty('--proximity-glow', glow);
    } else {
      this.element.style.boxShadow = 'none';
      this.element.style.setProperty('--proximity-glow', '0');
    }
  }
}

// Apply to all interactive elements
document.querySelectorAll('.hic-door, .hic-neon-sign, .hic-floor-marker').forEach(el => {
  new ProximityGlow(el);
});
```

**Performance considerations:**
- Proximity glow uses `requestAnimationFrame` to throttle calculations
- Only elements within the viewport are tracked (IntersectionObserver gates activation)
- On low-performance devices, proximity glow is disabled entirely
- The `mousemove` listener is passive and uses no DOM reads inside the hot path after initial `getBoundingClientRect`

---

## Zoom and Pan

Zoom and pan controls allow users to change their perspective within the HIC -- from the full building overview to an individual paragraph within a room.

### Zoom Out: Building Facade

Zooming out pulls the user's perspective back, revealing more of the building structure. At maximum zoom-out, the user sees the Street Level View -- the full building facade.

**Zoom levels:**

| Level | Name            | Viewport Content                       | Trigger                    |
|-------|-----------------|----------------------------------------|----------------------------|
| 0     | Street Level    | Full building facade                   | Ctrl/Cmd + minus (max out) |
| 1     | Floor Overview  | All rooms on current floor (floor plan)| Ctrl/Cmd + minus           |
| 2     | Corridor View   | 3-5 adjacent rooms visible             | Default for lobby          |
| 3     | Room View       | Single room fills viewport             | Default for articles       |
| 4     | Detail View     | Code blocks and diagrams fill viewport | Ctrl/Cmd + plus            |

**Zoom out behavior:**
```javascript
const zoomLevels = {
  current: 3,
  min: 0,
  max: 4,

  zoomOut() {
    if (this.current <= this.min) return;
    this.current--;
    this.applyZoom();
  },

  zoomIn() {
    if (this.current >= this.max) return;
    this.current++;
    this.applyZoom();
  },

  applyZoom() {
    const viewport = document.querySelector('.hic-viewport');
    viewport.dataset.zoomLevel = this.current;

    // Transition between zoom levels
    viewport.style.transition = 'transform 0.4s cubic-bezier(0.4, 0, 0.2, 1)';

    switch (this.current) {
      case 0:
        navigateTo('/');
        break;
      case 1:
        showFloorOverview();
        break;
      case 2:
        showCorridorView();
        break;
      case 3:
        showRoomView();
        break;
      case 4:
        showDetailView();
        break;
    }
  }
};
```

---

### Zoom In: Room Detail

Zooming in increases the level of detail visible. At maximum zoom, individual code blocks, diagrams, and inline elements fill the viewport. This is useful for examining code examples or reading dense technical content.

**Detail view features:**
- Code blocks expand to fill the viewport width with larger font size
- Diagrams scale up and become pannable
- Line numbers in code blocks become more prominent
- Syntax highlighting intensifies -- neon colors become bolder
- A "zoom out" floating button appears in the corner to return to Room View

**Code block zoom:**
```css
[data-zoom-level="4"] .hic-code-block {
  font-size: 16px;
  line-height: 1.6;
  padding: 24px;
  border: 2px solid var(--neon-accent);
  box-shadow: 0 0 20px var(--neon-accent-dim);
}

[data-zoom-level="4"] .hic-code-block .line-number {
  opacity: 1.0;
  color: var(--neon-accent);
  font-weight: 600;
}
```

---

### Pan: Scrolling Through Floors

Panning moves the viewport horizontally or vertically without changing the zoom level. The behavior depends on the current context.

**Pan behavior by context:**

| Context         | Horizontal Pan                   | Vertical Pan                        |
|-----------------|----------------------------------|--------------------------------------|
| Street Level    | Not applicable (fixed view)      | Scroll up/down the facade            |
| Floor Plan      | Pan across the floor layout      | Pan across the floor layout          |
| Corridor View   | Move between rooms               | Not applicable                       |
| Room View       | Not applicable                   | Scroll through content (Deep Dive)   |
| Detail View     | Pan within enlarged content      | Pan within enlarged content          |

**Scroll direction locking:**
```javascript
const panController = {
  activeAxis: null,
  threshold: 10, // pixels of movement before locking axis

  onPointerMove(event) {
    if (!this.activeAxis) {
      const dx = Math.abs(event.movementX);
      const dy = Math.abs(event.movementY);
      if (dx > this.threshold || dy > this.threshold) {
        this.activeAxis = dx > dy ? 'horizontal' : 'vertical';
      }
    }

    if (this.activeAxis === 'horizontal') {
      this.panHorizontal(event.movementX);
    } else if (this.activeAxis === 'vertical') {
      this.panVertical(event.movementY);
    }
  },

  onPointerUp() {
    this.activeAxis = null;
  }
};
```

---

### Pinch Zoom: Mobile Touch

On touch devices, pinch-to-zoom replaces the keyboard zoom controls. Two-finger pinch maps directly to the zoom level system.

**Pinch zoom implementation:**
```javascript
class PinchZoomHandler {
  constructor(viewport) {
    this.viewport = viewport;
    this.initialDistance = 0;
    this.currentZoom = 3; // Room View default

    viewport.addEventListener('touchstart', (e) => this.onTouchStart(e), { passive: true });
    viewport.addEventListener('touchmove', (e) => this.onTouchMove(e), { passive: false });
    viewport.addEventListener('touchend', () => this.onTouchEnd());
  }

  onTouchStart(event) {
    if (event.touches.length === 2) {
      this.initialDistance = this.getTouchDistance(event.touches);
    }
  }

  onTouchMove(event) {
    if (event.touches.length !== 2) return;
    event.preventDefault(); // Prevent native zoom

    const currentDistance = this.getTouchDistance(event.touches);
    const scale = currentDistance / this.initialDistance;

    if (scale > 1.3 && this.currentZoom < 4) {
      this.currentZoom++;
      this.initialDistance = currentDistance;
      zoomLevels.current = this.currentZoom;
      zoomLevels.applyZoom();
    } else if (scale < 0.7 && this.currentZoom > 0) {
      this.currentZoom--;
      this.initialDistance = currentDistance;
      zoomLevels.current = this.currentZoom;
      zoomLevels.applyZoom();
    }
  }

  getTouchDistance(touches) {
    const dx = touches[0].clientX - touches[1].clientX;
    const dy = touches[0].clientY - touches[1].clientY;
    return Math.hypot(dx, dy);
  }

  onTouchEnd() {
    this.initialDistance = 0;
  }
}
```

---

### Scroll Speed Zones

Different areas of the HIC have different scroll behaviors. The lobby scrolls differently than a room, and the elevator panel has its own scroll physics. This prevents a one-size-fits-all scroll speed from feeling wrong in different contexts.

**Scroll zone configurations:**

| Zone              | Scroll Speed | Scroll Snap | Overscroll       | Momentum          |
|-------------------|-------------|-------------|------------------|--------------------|
| Building facade   | 1.0x        | None        | Bounce (top)     | Standard           |
| Floor plan        | 0.8x        | Proximity   | Contain          | Dampened           |
| Corridor          | 1.0x        | Mandatory   | Contain          | Standard           |
| Room content      | 1.0x        | None        | Bounce (bottom)  | Standard           |
| Code blocks       | 0.6x        | None        | Contain          | Heavy dampening    |
| Elevator panel    | 0.5x        | Mandatory   | None             | No momentum        |

**Implementation:**
```css
.hic-floor-plan {
  scroll-behavior: smooth;
  scroll-snap-type: both proximity;
  overscroll-behavior: contain;
}

.hic-corridor {
  scroll-behavior: smooth;
  scroll-snap-type: x mandatory;
  overscroll-behavior-x: contain;
}

.hic-room-content {
  scroll-behavior: smooth;
  overscroll-behavior-y: auto; /* Allow bounce at bottom */
}

.hic-code-block {
  overflow-x: auto;
  overscroll-behavior: contain;
  /* Scroll speed handled via JS for dampening */
}
```

---

## Keyboard Navigation

Full keyboard navigation ensures the HIC is accessible to all users and efficient for power users. Every keyboard shortcut maps to a physical action within the building.

### Directional Movement

Arrow keys move between rooms on the same floor, like walking left or right down a corridor.

- **Left Arrow**: Move to the previous room (previous article in the domain)
- **Right Arrow**: Move to the next room (next article in the domain)
- When at the first room, Left Arrow does nothing (corridor dead end)
- When at the last room, Right Arrow does nothing (corridor dead end)
- In room view, Left/Right Arrow keys are not captured (allow native text selection)

### Floor Traversal

Page Up and Page Down move between floors, like taking the stairs one floor at a time.

- **Page Up**: Move to the floor above (previous domain in order)
- **Page Down**: Move to the floor below (next domain in order)
- Lands on the floor lobby (not a specific room)
- When on the top floor, Page Up goes to street level
- When on the ground floor, Page Down does nothing

### Home Key: Lobby Return

- **Home**: Return to street level (homepage) from anywhere in the building
- This is the "panic button" -- one keystroke to get completely out
- Confirmation is not required; the action is immediate

### Tab Cycling

Tab moves focus between interactive elements within the current view, following the building's spatial layout.

- **Tab**: Move to the next interactive element (door, sign, button)
- **Shift + Tab**: Move to the previous interactive element
- Tab order follows the visual layout: top-to-bottom, left-to-right
- Skip links are provided at the top of each page to jump past repeated navigation
- Focus wraps around -- tabbing past the last element returns to the first

### Slash: Building Intercom

- **/ (forward slash)**: Open the search overlay (building intercom)
- Only captured when no input element is focused
- The search overlay takes full keyboard focus when open
- See the Search as Intercom System section for full details

### Escape: Overlay Dismissal

Escape closes things in order of specificity, working outward from the most local context.

- **Escape (first press)**: Close the current overlay, modal, or preview card
- **Escape (second press, if no overlay)**: Return to the parent view (room to corridor, corridor to lobby)
- **Escape (in search)**: Close search overlay and return to previous view

### Full Keyboard Map

```
Navigation:
  Arrow Left        Move to previous room on this floor
  Arrow Right       Move to next room on this floor
  Arrow Up          Scroll up within current room / move to section above
  Arrow Down        Scroll down within current room / move to section below
  Page Up           Move to the floor above
  Page Down         Move to the floor below
  Home              Return to street level (homepage)
  End               Jump to bottom of current room

Actions:
  Enter             Open focused door / activate focused sign
  Space             Activate focused button / toggle focused element
  /                 Open search (building intercom)
  Escape            Close overlay / return to parent view
  Tab               Move focus to next interactive element
  Shift + Tab       Move focus to previous interactive element

Shortcuts:
  Ctrl/Cmd + K      Open command palette (building directory)
  Ctrl/Cmd + /      Toggle keyboard shortcut help overlay
  Ctrl/Cmd + +      Zoom in (increase detail level)
  Ctrl/Cmd + -      Zoom out (decrease detail level)
  Ctrl/Cmd + 0      Reset zoom to default (Room View)
  Ctrl/Cmd + F      Browser find (room-specific text search)
  Ctrl/Cmd + C      Copy selected text
  ?                 Show keyboard shortcut overlay (when no input focused)
```

---

## Search as Intercom System

Search is the building's intercom -- a way to communicate with the building itself and ask it to guide you to what you need. The metaphor extends from input (speaking into the intercom) through results (the building announcing where to go) to navigation (following the announcement to a room).

### Intercom Architecture

The search system is structured as a three-layer intercom network:

1. **Room-level intercom**: Searches within the current article (browser find, Ctrl+F)
2. **Floor-level intercom**: Searches within the current domain
3. **Building-level intercom**: Searches the entire documentation site

**Architecture diagram:**
```
                    +----------------------------+
                    |   Building-Wide Intercom    |
                    |   (Global full-text index)  |
                    +----------------------------+
                           |            |
              +------------+            +------------+
              |                                      |
     +--------+--------+               +--------+--------+
     | Floor 1 Intercom |               | Floor N Intercom |
     | (Domain index)   |               | (Domain index)   |
     +---------+--------+               +---------+--------+
               |                                   |
     +---------+---------+               +---------+---------+
     | Room 1  | Room 2  |               | Room 1  | Room 2  |
     | (Local) | (Local) |               | (Local) | (Local) |
     +---------+---------+               +---------+---------+
```

### Voice Query: Search Input

"Speaking into the intercom" is typing a search query. The search input is styled as an intercom panel with a microphone icon and a glowing input field.

**Search overlay design:**
```html
<div class="hic-intercom-overlay" role="dialog" aria-label="Search">
  <div class="hic-intercom-panel">
    <div class="hic-intercom-header">
      <span class="hic-intercom-icon" aria-hidden="true">&#x1f50a;</span>
      <span class="hic-intercom-title">Building Intercom</span>
      <button class="hic-intercom-close" aria-label="Close search">ESC</button>
    </div>
    <div class="hic-intercom-input-container">
      <input
        type="search"
        class="hic-intercom-input"
        placeholder="Where would you like to go?"
        aria-label="Search documentation"
        autocomplete="off"
        autofocus
      />
      <div class="hic-intercom-scope">
        <button class="hic-intercom-scope-btn active" data-scope="building">
          Entire Building
        </button>
        <button class="hic-intercom-scope-btn" data-scope="floor">
          This Floor Only
        </button>
      </div>
    </div>
    <div class="hic-intercom-results" role="listbox" aria-label="Search results">
      <!-- Results render here -->
    </div>
    <div class="hic-intercom-recent" aria-label="Recent searches">
      <!-- Recent searches render here -->
    </div>
  </div>
</div>
```

**Search input behavior:**
- Search triggers on keypress with a 200ms debounce
- Results appear as the user types (live search)
- The first result is auto-highlighted (press Enter to navigate)
- Arrow keys move through results
- Escape closes the overlay and returns to the previous view
- The search input remembers the last query within the session

### Intercom Response: Result Rendering

Search results are rendered as building announcements -- each result is a card showing where in the building the match was found.

**Result card structure:**
```html
<div class="hic-intercom-result" role="option" aria-selected="false">
  <div class="hic-intercom-result__location">
    <span class="hic-intercom-result__floor" style="color: var(--floor-color)">
      Floor 3: API Reference
    </span>
    <span class="hic-intercom-result__room">
      Room: Authentication
    </span>
  </div>
  <div class="hic-intercom-result__content">
    <h4 class="hic-intercom-result__title">OAuth 2.0 Flow</h4>
    <p class="hic-intercom-result__excerpt">
      Configure the <mark>authentication</mark> flow by setting
      the redirect URI and client credentials...
    </p>
  </div>
  <div class="hic-intercom-result__meta">
    <span class="hic-intercom-result__type">Section</span>
    <span class="hic-intercom-result__relevance" aria-label="Relevance: high">
      &#9733;&#9733;&#9733;
    </span>
  </div>
</div>
```

**Result ranking factors:**
1. Title match (highest weight -- the room name matches)
2. Heading match (section name within a room)
3. Content match (body text within a room)
4. Recency (recently updated rooms rank higher)
5. Proximity (rooms on the current floor rank higher)

### Floor-Specific Search

When the user toggles the scope to "This Floor Only," the intercom limits results to the current domain. This is useful when the user knows they are in the right general area but needs to find a specific room.

**Floor search behavior:**
- Results only include articles from the current domain
- The floor name is shown in the search scope indicator
- If no results are found on the current floor, a suggestion appears: "No results on this floor. Try the building-wide intercom?"
- Floor search uses a pre-built index specific to that domain for faster results

### Building-Wide Search

Building-wide search queries the full documentation index. Results are grouped by floor for easy scanning.

**Grouped results display:**
```
Search results for "authentication" (12 results)

  Floor 1: Getting Started (2 results)
    - Quick Start Guide > Setting Up Auth
    - First API Call > Authentication Header

  Floor 3: API Reference (5 results)
    - Authentication > OAuth 2.0 Flow
    - Authentication > API Keys
    - Authentication > Session Tokens
    - Authentication > Token Refresh
    - Errors > 401 Unauthorized

  Floor 4: Guides (3 results)
    - Security Best Practices > Authentication Patterns
    - Multi-Tenant Auth > Tenant Isolation
    - SSO Integration > SAML Configuration

  Floor 5: Concepts (2 results)
    - Identity Management > Auth Architecture
    - Security Model > Token Lifecycle
```

### Recent Announcements: Search History

Recent searches are stored locally and displayed when the search overlay opens with an empty input. This allows quick re-access to previous queries.

**Recent search implementation:**
```javascript
class SearchHistory {
  constructor(maxEntries = 10) {
    this.maxEntries = maxEntries;
    this.storageKey = 'hic-intercom-history';
  }

  getHistory() {
    const raw = localStorage.getItem(this.storageKey);
    return raw ? JSON.parse(raw) : [];
  }

  addEntry(query, resultCount) {
    const history = this.getHistory();
    // Remove duplicate if exists
    const filtered = history.filter(entry => entry.query !== query);
    // Add new entry at the front
    filtered.unshift({
      query,
      resultCount,
      timestamp: Date.now(),
    });
    // Trim to max length
    const trimmed = filtered.slice(0, this.maxEntries);
    localStorage.setItem(this.storageKey, JSON.stringify(trimmed));
  }

  clearHistory() {
    localStorage.removeItem(this.storageKey);
  }

  render(container) {
    const history = this.getHistory();
    if (history.length === 0) {
      container.innerHTML = '<p class="hic-intercom-no-history">No recent announcements.</p>';
      return;
    }

    container.innerHTML = history.map(entry => `
      <button class="hic-intercom-recent-entry" data-query="${entry.query}">
        <span class="hic-intercom-recent-query">${entry.query}</span>
        <span class="hic-intercom-recent-count">${entry.resultCount} results</span>
        <span class="hic-intercom-recent-time">${this.formatTime(entry.timestamp)}</span>
      </button>
    `).join('');
  }

  formatTime(timestamp) {
    const diff = Date.now() - timestamp;
    const minutes = Math.floor(diff / 60000);
    if (minutes < 1) return 'Just now';
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`;
    const days = Math.floor(hours / 24);
    return `${days}d ago`;
  }
}
```

---

## Mobile Experience

The HIC building adapts to mobile devices. Rooms become narrower, corridors become single-file, and the building reorients to fit the screen. Touch gestures replace mouse interactions, and the spatial metaphor is preserved through careful adaptation.

### Touch Gesture Mapping

Every mouse interaction has a touch equivalent. The mapping preserves the building metaphor while feeling natural on a touch device.

| Mouse Action        | Touch Equivalent          | Building Metaphor                  |
|---------------------|---------------------------|------------------------------------|
| Click               | Tap                       | Push open a door                   |
| Hover               | Long press (500ms)        | Peer through a window              |
| Right-click         | Long press + context menu | Examine a sign closely             |
| Scroll              | Single-finger swipe       | Walk down a corridor               |
| Zoom                | Pinch                     | Step closer or further away        |
| Drag                | Touch and drag             | Push furniture / rearrange          |
| Double-click        | Double-tap                | Knock on a door (quick-open)       |

### Swipe Navigation

Swipe gestures provide fluid navigation between rooms and floors. The direction of the swipe maps to the spatial direction of movement.

**Swipe configuration:**
```javascript
const swipeConfig = {
  // Minimum distance to trigger a swipe (pixels)
  threshold: 50,

  // Maximum time for the swipe gesture (milliseconds)
  maxDuration: 300,

  // Swipe directions and their actions
  directions: {
    left: {
      corridor: 'nextRoom',          // Next room in the corridor
      room: 'nextRoom',              // Next article
      lobby: null,                    // No action
    },
    right: {
      corridor: 'previousRoom',      // Previous room in the corridor
      room: 'previousRoom',          // Previous article
      lobby: null,                    // No action
    },
    up: {
      corridor: null,                 // No action
      room: 'scrollContent',          // Continue scrolling
      lobby: 'floorAbove',            // Go to floor above
    },
    down: {
      corridor: null,                 // No action
      room: 'scrollContent',          // Continue scrolling
      lobby: 'floorBelow',            // Go to floor below
    },
  },

  // Edge swipe (from left edge of screen)
  edgeSwipe: {
    threshold: 20,                    // pixels from edge to trigger
    action: 'openElevatorPanel',      // Open navigation panel
  },
};
```

**Swipe feedback:**
- As the user swipes, the current view translates in the swipe direction
- The adjacent content (next room, previous room) is visible behind the current view
- If the swipe distance exceeds the threshold, the navigation completes on release
- If the swipe distance does not exceed the threshold, the view snaps back elastically

### Long Press Preview

Long press replaces hover on touch devices. Pressing and holding on a door or link shows the same preview that a desktop user would see on hover.

**Long press implementation:**
```javascript
class LongPressHandler {
  constructor(element, options = {}) {
    this.element = element;
    this.delay = options.delay || 500;
    this.timer = null;
    this.isActive = false;

    element.addEventListener('touchstart', (e) => this.onTouchStart(e), { passive: true });
    element.addEventListener('touchend', () => this.onTouchEnd());
    element.addEventListener('touchmove', () => this.onTouchEnd()); // Cancel on move
    element.addEventListener('contextmenu', (e) => e.preventDefault()); // Prevent native menu
  }

  onTouchStart(event) {
    this.timer = setTimeout(() => {
      this.isActive = true;
      // Haptic feedback if available
      if (navigator.vibrate) navigator.vibrate(10);
      // Show preview
      showPreviewCard(this.element, event.touches[0]);
    }, this.delay);
  }

  onTouchEnd() {
    clearTimeout(this.timer);
    if (this.isActive) {
      this.isActive = false;
      hidePreviewCard();
    }
  }
}
```

### Two-Finger Scroll

Two-finger scroll is used for panning across floor plans and navigating large diagrams within rooms. It prevents conflict with single-finger scroll (which handles content scrolling).

**Two-finger scroll behavior:**
- Two fingers on the floor plan: pan freely in both axes
- Two fingers on a room diagram: pan within the diagram bounds
- Two fingers on regular content: default system scroll behavior (no override)
- When two-finger scroll is active, momentum is preserved after fingers lift

### Responsive Layout: Adaptive Building

The HIC building physically adapts to the screen size. It does not simply reflow content -- it changes its architecture.

**Breakpoints and building adaptations:**

| Breakpoint      | Screen Width    | Building Adaptation                        |
|-----------------|-----------------|--------------------------------------------|
| Desktop XL      | >= 1440px       | Full building with side panels visible      |
| Desktop         | 1024 - 1439px   | Full building, side panels collapsible      |
| Tablet          | 768 - 1023px    | Building narrows, corridor becomes linear   |
| Mobile Large    | 480 - 767px     | Single-room view, floor plan as overlay     |
| Mobile Small    | < 480px         | Minimal room view, bottom sheet navigation  |

**CSS breakpoint implementation:**
```css
/* Desktop: Full building with all panels */
@media (min-width: 1024px) {
  .hic-building {
    display: grid;
    grid-template-columns: 260px 1fr 300px;
    grid-template-rows: auto 1fr auto;
    grid-template-areas:
      "elevator header header"
      "elevator content sidebar"
      "elevator footer  footer";
  }
}

/* Tablet: Narrower building, collapsible panels */
@media (min-width: 768px) and (max-width: 1023px) {
  .hic-building {
    display: grid;
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr auto;
    grid-template-areas:
      "header"
      "content"
      "footer";
  }

  .hic-elevator-panel {
    position: fixed;
    left: -260px;
    transition: left 0.3s ease;
  }

  .hic-elevator-panel.is-open {
    left: 0;
  }
}

/* Mobile: Single-room view */
@media (max-width: 767px) {
  .hic-building {
    display: flex;
    flex-direction: column;
    min-height: 100dvh;
  }

  .hic-elevator-panel {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: auto;
    max-height: 60vh;
    transform: translateY(calc(100% - 56px));
    transition: transform 0.3s ease;
  }

  .hic-elevator-panel.is-open {
    transform: translateY(0);
  }
}
```

### Portrait Mode: Single-Room View

In portrait mode (typical phone orientation), the building shows a single room at a time. The corridor is not visible. Navigation between rooms uses swipe gestures or the elevator panel.

**Portrait mode characteristics:**
- Only one room (article) is visible at a time
- The room fills the full viewport width
- Room map (table of contents) is accessible via a collapsible top bar
- Floor navigation is accessed via a bottom sheet (swipe up from bottom)
- Breadcrumbs collapse to show only the current floor and room

### Landscape Mode: Floor Plan Overview

In landscape mode (phone rotated or tablet), the building shows the floor plan with room previews visible. This gives users an overview of the floor before entering a specific room.

**Landscape mode characteristics:**
- Floor plan is visible with rooms arranged in a grid
- Each room card shows title and brief excerpt
- Tapping a room card enters the room (switches to portrait-like room view)
- The elevator panel appears as a thin strip on the left edge
- More horizontal space allows side-by-side viewing of room content and room map

---

## Accessibility

The HIC building must be accessible to all visitors. Accessibility is not an afterthought layered on top of the visual design -- it is a structural requirement of the building itself. Every neon sign has a text equivalent. Every door has a label. Every floor has audio guidance.

### Screen Reader: Audio Descriptions

Screen readers provide audio descriptions of each room, floor, and building element. The HIC uses ARIA landmarks, roles, and labels to create a complete audio representation of the building.

**ARIA landmark mapping:**
```html
<body>
  <!-- Street level / page header -->
  <header role="banner" aria-label="HIC Building Entrance">
    <nav role="navigation" aria-label="Express Elevator">
      <!-- Elevator panel -->
    </nav>
  </header>

  <!-- Breadcrumb / exit signs -->
  <nav aria-label="Emergency Exit Signs (Breadcrumb)">
    <ol role="list">
      <!-- Breadcrumb items -->
    </ol>
  </nav>

  <!-- Main content area / current room -->
  <main role="main" aria-label="Current Room: Authentication">
    <!-- Room content -->
  </main>

  <!-- Room map / table of contents -->
  <nav role="navigation" aria-label="Room Map (Table of Contents)">
    <!-- Section links -->
  </nav>

  <!-- Building foundation / footer -->
  <footer role="contentinfo" aria-label="Building Foundation">
    <!-- Footer content -->
  </footer>
</body>
```

**Live region announcements:**
```javascript
// Announce navigation events to screen readers
const announcer = document.createElement('div');
announcer.setAttribute('aria-live', 'polite');
announcer.setAttribute('aria-atomic', 'true');
announcer.classList.add('sr-only');
document.body.appendChild(announcer);

function announce(message) {
  announcer.textContent = '';
  // Force re-announcement by clearing and setting in separate frames
  requestAnimationFrame(() => {
    announcer.textContent = message;
  });
}

// Usage examples:
announce('Entering Floor 3: API Reference');
announce('Opening room: Authentication');
announce('Search results: 12 matches found');
announce('Elevator arriving at Floor 5: Concepts');
```

### High Contrast: Emergency Lighting

High contrast mode is the building's "emergency lighting." When activated, all neon effects are replaced with solid, high-contrast colors. Backgrounds become pure black or pure white. Text contrast ratios meet or exceed WCAG AAA standards (7:1).

**High contrast color scheme:**
```css
@media (prefers-contrast: high) {
  :root {
    --bg-primary: #000000;
    --bg-secondary: #1a1a1a;
    --text-primary: #ffffff;
    --text-secondary: #e0e0e0;
    --border-color: #ffffff;
    --link-color: #66ccff;
    --link-visited: #cc99ff;
    --focus-outline: #ffff00;
    --neon-glow: none;
  }

  /* Remove all glow and blur effects */
  .hic-neon-sign,
  .hic-door,
  .hic-floor-marker {
    text-shadow: none !important;
    box-shadow: none !important;
    filter: none !important;
  }

  /* Solid borders instead of glow */
  .hic-door:hover,
  .hic-neon-sign:hover {
    outline: 3px solid var(--focus-outline);
    outline-offset: 2px;
  }

  /* High contrast code blocks */
  .hic-code-block {
    background: #000000;
    border: 2px solid #ffffff;
    color: #ffffff;
  }
}
```

### Keyboard-Only Navigation

All building features must be accessible without a mouse. This section supplements the Keyboard Navigation section with accessibility-specific requirements.

**Focus management rules:**
1. Focus must be visible at all times -- no element may receive focus without a visible indicator
2. Focus order must follow the visual layout (no tab traps, no illogical jumps)
3. Modal overlays must trap focus within themselves until dismissed
4. After navigation (entering a room), focus moves to the room's heading
5. Skip links must be provided to bypass repeated navigation elements

**Skip links:**
```html
<a href="#main-content" class="hic-skip-link">
  Skip to room content
</a>
<a href="#room-map" class="hic-skip-link">
  Skip to room map
</a>
<a href="#elevator-panel" class="hic-skip-link">
  Skip to elevator panel
</a>
```

**Skip link styling:**
```css
.hic-skip-link {
  position: absolute;
  top: -100%;
  left: 0;
  padding: 12px 24px;
  background: var(--bg-primary);
  color: var(--text-primary);
  border: 2px solid var(--focus-outline);
  z-index: 10000;
  font-size: 16px;
}

.hic-skip-link:focus {
  top: 0;
}
```

**Focus trap for modals:**
```javascript
function trapFocus(modalElement) {
  const focusableSelectors = [
    'a[href]', 'button:not([disabled])', 'input:not([disabled])',
    'select:not([disabled])', 'textarea:not([disabled])',
    '[tabindex]:not([tabindex="-1"])'
  ];

  const focusableElements = modalElement.querySelectorAll(focusableSelectors.join(','));
  const firstFocusable = focusableElements[0];
  const lastFocusable = focusableElements[focusableElements.length - 1];

  modalElement.addEventListener('keydown', (event) => {
    if (event.key !== 'Tab') return;

    if (event.shiftKey) {
      if (document.activeElement === firstFocusable) {
        event.preventDefault();
        lastFocusable.focus();
      }
    } else {
      if (document.activeElement === lastFocusable) {
        event.preventDefault();
        firstFocusable.focus();
      }
    }
  });

  firstFocusable.focus();
}
```

### Reduced Motion: Static Signage

When the user has enabled reduced motion preferences, all neon animations, transition effects, and scroll-triggered reveals are disabled. Signs become static. Doors do not animate. Transitions are instant.

**Reduced motion implementation:**
```css
@media (prefers-reduced-motion: reduce) {
  /* Disable all transitions */
  *, *::before, *::after {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    scroll-behavior: auto !important;
  }

  /* Remove neon flicker */
  .hic-neon-sign {
    animation: none !important;
  }

  /* Disable parallax */
  .hic-parallax {
    transform: none !important;
  }

  /* Static scroll reveals - show everything immediately */
  .hic-scroll-reveal {
    opacity: 1 !important;
    transform: none !important;
  }

  /* Disable proximity glow */
  .hic-door,
  .hic-neon-sign,
  .hic-floor-marker {
    box-shadow: none !important;
  }
}
```

**JavaScript reduced motion check:**
```javascript
const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)');

function getAnimationDuration(defaultMs) {
  return prefersReducedMotion.matches ? 0 : defaultMs;
}

// Listen for changes (user may toggle preference during session)
prefersReducedMotion.addEventListener('change', (event) => {
  if (event.matches) {
    disableAllAnimations();
  } else {
    enableAnimations();
  }
});
```

### Focus Indicators: Highlighted Doorways

Focus indicators are the HIC equivalent of a highlighted doorway -- a bright outline showing exactly where the user is standing. Every focusable element must have a visible focus indicator that meets WCAG 2.2 requirements.

**Focus indicator specifications:**
- Outline: 3px solid, high-contrast color (yellow on dark backgrounds, blue on light)
- Outline offset: 2px (space between the element and the outline)
- Outline must not be obscured by other elements
- Focus indicator must be visible in all color modes (standard, high contrast, dark, light)

**Focus indicator styles:**
```css
/* Base focus indicator */
:focus-visible {
  outline: 3px solid var(--focus-outline, #ffcc00);
  outline-offset: 2px;
  border-radius: 2px;
}

/* Door-specific focus (navigational links) */
.hic-door:focus-visible {
  outline-color: var(--neon-accent);
  box-shadow: 0 0 0 6px rgba(255, 204, 0, 0.2);
}

/* Neon sign focus (interactive elements) */
.hic-neon-sign:focus-visible {
  outline-color: #ffffff;
  outline-width: 3px;
}

/* Remove default outline only when :focus-visible is supported */
:focus:not(:focus-visible) {
  outline: none;
}
```

---

## Loading States

The building is always under some degree of construction. Loading states communicate progress and maintain the spatial metaphor even when content is not yet available.

### Building Under Construction

Skeleton screens are the "construction scaffolding" visible when a floor or room is loading. They show the shape of the content to come without revealing its details.

**Skeleton screen design principles:**
- Skeletons match the layout of the final content
- Animated shimmer effect moves left-to-right (like welding sparks moving along scaffolding)
- Skeletons use the domain's neon color at low opacity for the shimmer
- No text content is shown -- only geometric placeholders

**Skeleton implementation:**
```css
.hic-skeleton {
  background: linear-gradient(
    90deg,
    var(--bg-secondary) 25%,
    var(--bg-tertiary) 50%,
    var(--bg-secondary) 75%
  );
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s ease-in-out infinite;
  border-radius: 4px;
}

@keyframes skeleton-shimmer {
  0%   { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* Respect reduced motion */
@media (prefers-reduced-motion: reduce) {
  .hic-skeleton {
    animation: none;
    background: var(--bg-secondary);
  }
}
```

**Skeleton components:**
```html
<!-- Room skeleton -->
<article class="hic-room hic-room--loading" aria-busy="true" aria-label="Loading room content">
  <div class="hic-skeleton hic-skeleton--title" style="width: 60%; height: 32px;"></div>
  <div class="hic-skeleton hic-skeleton--meta" style="width: 30%; height: 16px; margin-top: 8px;"></div>
  <div class="hic-skeleton hic-skeleton--text" style="width: 100%; height: 16px; margin-top: 24px;"></div>
  <div class="hic-skeleton hic-skeleton--text" style="width: 95%; height: 16px; margin-top: 8px;"></div>
  <div class="hic-skeleton hic-skeleton--text" style="width: 87%; height: 16px; margin-top: 8px;"></div>
  <div class="hic-skeleton hic-skeleton--code" style="width: 100%; height: 120px; margin-top: 24px;"></div>
  <div class="hic-skeleton hic-skeleton--text" style="width: 100%; height: 16px; margin-top: 24px;"></div>
  <div class="hic-skeleton hic-skeleton--text" style="width: 92%; height: 16px; margin-top: 8px;"></div>
</article>
```

### Elevator in Transit

Page transition animations play when the user is moving between floors or rooms. The elevator is "in transit" -- the user is between locations.

**Transition states:**
```
State 1: DEPARTING
  - Current view fades out
  - "Doors closing" effect (viewport edges darken inward)
  - Duration: 150ms

State 2: IN TRANSIT
  - Floor counter animates (shows floors passing)
  - Subtle vertical motion
  - Duration: 200ms + (50ms * floors_traveled)

State 3: ARRIVING
  - "Doors opening" effect (viewport edges lighten outward)
  - New view fades in
  - Duration: 150ms
```

**Transition implementation:**
```javascript
async function elevatorTransition(fromFloor, toFloor) {
  const distance = Math.abs(toFloor - fromFloor);
  const transitDuration = 200 + (50 * distance);

  const viewport = document.querySelector('.hic-viewport');

  // State 1: Departing
  viewport.classList.add('hic-transit--departing');
  await sleep(150);

  // State 2: In Transit
  viewport.classList.remove('hic-transit--departing');
  viewport.classList.add('hic-transit--moving');

  // Animate floor counter
  const counter = document.querySelector('.hic-elevator-counter');
  const direction = toFloor > fromFloor ? 1 : -1;
  let current = fromFloor;
  const counterInterval = setInterval(() => {
    current += direction;
    counter.textContent = current;
    if (current === toFloor) clearInterval(counterInterval);
  }, transitDuration / distance);

  await sleep(transitDuration);

  // State 3: Arriving
  viewport.classList.remove('hic-transit--moving');
  viewport.classList.add('hic-transit--arriving');
  await sleep(150);

  viewport.classList.remove('hic-transit--arriving');
}
```

### Room Furnishing

Content loading progressively within a room -- elements appear in order as they load, like furniture being placed in a room.

**Progressive loading order:**
1. Room title and metadata (instant -- always in initial payload)
2. Text content (fast -- HTML from server)
3. Room map / table of contents (derived from headings, built client-side)
4. Code blocks with syntax highlighting (requires highlight.js processing)
5. Images and diagrams (loaded lazily as they approach the viewport)
6. Interactive elements (hydrated after main content is stable)

**Progressive enhancement stages:**
```javascript
const furnishingStages = [
  {
    name: 'structure',
    description: 'Room walls and layout',
    elements: ['h1', 'h2', 'h3', 'nav'],
    timing: 'immediate',
  },
  {
    name: 'content',
    description: 'Room contents -- text and basic elements',
    elements: ['p', 'ul', 'ol', 'blockquote', 'table'],
    timing: 'immediate',
  },
  {
    name: 'decoration',
    description: 'Code highlighting, syntax colors',
    elements: ['.hic-code-block'],
    timing: 'after DOMContentLoaded',
  },
  {
    name: 'artwork',
    description: 'Images, diagrams, media',
    elements: ['img', 'figure', 'video', '.hic-diagram'],
    timing: 'lazy, on intersection',
  },
  {
    name: 'electronics',
    description: 'Interactive elements, widgets',
    elements: ['.hic-interactive', '.hic-playground', '.hic-feedback'],
    timing: 'after idle callback',
  },
];
```

### Power Outage

Error states represent power outages in the building. When something goes wrong -- a page fails to load, a search returns an error, a network request times out -- the building's power flickers or fails.

**Error state categories:**

| Error               | Building Metaphor        | User-Facing Message                        |
|----------------------|--------------------------|---------------------------------------------|
| 404 Not Found        | Room not found           | "This room doesn't exist. Check the floor directory." |
| 500 Server Error     | Power outage             | "Building power failure. Maintenance has been notified." |
| Network timeout      | Elevator stuck           | "Connection interrupted. Try again shortly." |
| Search failure       | Intercom malfunction     | "Intercom offline. Try a floor-level search." |
| Asset load failure   | Broken furnishing        | "Some room elements failed to load. Refresh to try again." |
| Rate limit           | Building at capacity     | "Building is at capacity. Please wait a moment." |

**Error page design:**
```html
<div class="hic-error hic-error--404" role="alert">
  <div class="hic-error__visual">
    <!-- Flickering neon sign effect -->
    <h1 class="hic-error__code hic-neon--flicker">404</h1>
    <p class="hic-error__title">Room Not Found</p>
  </div>
  <div class="hic-error__content">
    <p>The room you're looking for doesn't exist on this floor.</p>
    <p>It may have been moved, renamed, or decommissioned.</p>
  </div>
  <div class="hic-error__actions">
    <a href="/" class="hic-neon-sign">Return to Street Level</a>
    <button class="hic-neon-sign" onclick="history.back()">Go Back</button>
    <button class="hic-neon-sign" onclick="openSearch()">Use Intercom (Search)</button>
  </div>
  <div class="hic-error__suggestions">
    <h2>Nearby Rooms</h2>
    <!-- Dynamically populated with similar articles -->
  </div>
</div>
```

**Neon flicker effect for error states:**
```css
@keyframes neon-flicker-error {
  0%, 19%, 21%, 23%, 25%, 54%, 56%, 100% {
    text-shadow:
      0 0 4px #ff0040,
      0 0 11px #ff0040,
      0 0 19px #ff0040,
      0 0 40px #ff0040;
    opacity: 1;
  }
  20%, 24%, 55% {
    text-shadow: none;
    opacity: 0.6;
  }
}

.hic-neon--flicker {
  animation: neon-flicker-error 3s infinite alternate;
}
```

---

## Performance

The HIC must load fast and stay fast. A slow building is a broken building. Performance optimization uses the building metaphor to guide implementation: you only furnish rooms people are about to enter, you keep recently visited rooms intact, and you have emergency power for offline access.

### Lazy Loading: Render on Approach

Rooms (articles) and their contents load only when the user approaches them. "Approaching" means the content is near the viewport or the user has signaled intent to navigate.

**Lazy loading strategy:**

| Content Type       | Load Trigger                              | Priority   |
|--------------------|-------------------------------------------|------------|
| Current room text  | Immediate (in initial payload)            | Critical   |
| Current room code  | DOMContentLoaded                          | High       |
| Current room images| IntersectionObserver (200px rootMargin)   | Medium     |
| Adjacent rooms     | User idle + prefetch hint                 | Low        |
| Other floor content| Not loaded until navigated to             | None       |

**IntersectionObserver for room furnishing:**
```javascript
const furnishingObserver = new IntersectionObserver(
  (entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const element = entry.target;
        furnishElement(element);
        furnishingObserver.unobserve(element);
      }
    });
  },
  {
    rootMargin: '200px 0px',  // Start loading 200px before visible
    threshold: 0.01,          // Trigger as soon as 1% is visible
  }
);

// Observe all lazy elements
document.querySelectorAll('[data-furnish]').forEach(el => {
  furnishingObserver.observe(el);
});

function furnishElement(element) {
  const type = element.dataset.furnish;
  switch (type) {
    case 'image':
      element.src = element.dataset.src;
      element.removeAttribute('data-src');
      break;
    case 'code':
      highlightCode(element);
      break;
    case 'interactive':
      hydrateWidget(element);
      break;
  }
}
```

### Prefetching: Adjacent Pre-Load

Adjacent rooms (the next and previous articles in the domain, plus any prominently linked articles) are prefetched after the current room is fully loaded. This ensures that corridor navigation feels instantaneous.

**Prefetch strategy:**
```javascript
function prefetchAdjacentRooms() {
  // Wait until the current room is fully loaded and the browser is idle
  if ('requestIdleCallback' in window) {
    requestIdleCallback(() => {
      performPrefetch();
    }, { timeout: 3000 });
  } else {
    setTimeout(performPrefetch, 2000);
  }
}

function performPrefetch() {
  const adjacentLinks = getAdjacentRoomLinks();

  adjacentLinks.forEach(link => {
    // Use <link rel="prefetch"> for document prefetch
    const prefetchLink = document.createElement('link');
    prefetchLink.rel = 'prefetch';
    prefetchLink.href = link.href;
    prefetchLink.as = 'document';
    document.head.appendChild(prefetchLink);
  });
}

function getAdjacentRoomLinks() {
  const links = [];

  // Previous and next room (corridor neighbors)
  const prevLink = document.querySelector('[rel="prev"]');
  const nextLink = document.querySelector('[rel="next"]');
  if (prevLink) links.push(prevLink);
  if (nextLink) links.push(nextLink);

  // Prominently linked rooms (first 3 internal links in content)
  const contentLinks = document.querySelectorAll('main a[href^="/"]');
  const seen = new Set();
  for (const link of contentLinks) {
    if (seen.has(link.href)) continue;
    seen.add(link.href);
    links.push(link);
    if (seen.size >= 3) break;
  }

  return links;
}
```

### Caching: Furnished Rooms

Once a room has been visited and its content loaded, it stays "furnished" in the cache. Returning to a previously visited room is instant -- no network request, no loading skeleton, no transition delay.

**Caching layers:**

| Layer                | Scope               | Lifetime              | Content Cached                |
|----------------------|---------------------|-----------------------|-------------------------------|
| Memory cache         | Current session     | Until page unload     | Parsed DOM, computed state    |
| Session storage      | Current tab         | Until tab close       | Scroll position, room state   |
| Service worker cache | Persistent          | Until cache refresh   | HTML, CSS, JS, images         |
| HTTP cache           | Browser-managed     | Per cache headers     | All static assets             |

**Memory cache for room state:**
```javascript
class RoomCache {
  constructor(maxSize = 20) {
    this.cache = new Map();
    this.maxSize = maxSize;
  }

  set(roomId, state) {
    // Evict oldest entry if at capacity
    if (this.cache.size >= this.maxSize) {
      const oldestKey = this.cache.keys().next().value;
      this.cache.delete(oldestKey);
    }

    this.cache.set(roomId, {
      html: state.html,
      scrollPosition: state.scrollPosition,
      activeSection: state.activeSection,
      timestamp: Date.now(),
    });
  }

  get(roomId) {
    const entry = this.cache.get(roomId);
    if (!entry) return null;

    // Move to end (LRU)
    this.cache.delete(roomId);
    this.cache.set(roomId, entry);

    return entry;
  }

  has(roomId) {
    return this.cache.has(roomId);
  }

  invalidate(roomId) {
    this.cache.delete(roomId);
  }

  clear() {
    this.cache.clear();
  }
}

const roomCache = new RoomCache();
```

### Offline Mode: Emergency Generator

The service worker acts as the building's emergency generator. When the network is unavailable, previously visited rooms remain accessible. The building does not go dark -- it switches to backup power.

**Service worker strategy:**
```javascript
// sw.js -- HIC Emergency Generator

const CACHE_NAME = 'hic-building-v1';
const CRITICAL_ASSETS = [
  '/',
  '/styles/hic-core.css',
  '/scripts/hic-core.js',
  '/fonts/neon-display.woff2',
  '/images/hic-logo.svg',
  '/offline.html',
];

// Install: Cache critical building infrastructure
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll(CRITICAL_ASSETS);
    })
  );
});

// Fetch: Network-first with cache fallback
self.addEventListener('fetch', (event) => {
  // Only cache GET requests for same-origin documents and assets
  if (event.request.method !== 'GET') return;

  event.respondWith(
    fetch(event.request)
      .then((response) => {
        // Clone and cache the fresh response
        const responseClone = response.clone();
        caches.open(CACHE_NAME).then((cache) => {
          cache.put(event.request, responseClone);
        });
        return response;
      })
      .catch(() => {
        // Network failed: try the cache
        return caches.match(event.request).then((cachedResponse) => {
          if (cachedResponse) {
            return cachedResponse;
          }

          // If the request is for a page, show the offline page
          if (event.request.mode === 'navigate') {
            return caches.match('/offline.html');
          }

          // For other assets, return a generic failure
          return new Response('Offline', { status: 503, statusText: 'Offline' });
        });
      })
  );
});

// Activate: Clean up old caches
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames
          .filter((name) => name !== CACHE_NAME)
          .map((name) => caches.delete(name))
      );
    })
  );
});
```

**Offline page design:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>HIC - Emergency Generator Active</title>
</head>
<body class="hic-offline">
  <div class="hic-offline__panel">
    <h1 class="hic-neon--dim">Emergency Generator Active</h1>
    <p>The building has lost external power (network connection).</p>
    <p>Previously visited rooms are still accessible from backup power.</p>
    <div class="hic-offline__actions">
      <button onclick="location.reload()">Try Reconnecting</button>
      <a href="/">Return to Street Level</a>
    </div>
    <div class="hic-offline__cached-rooms">
      <h2>Available Rooms (Cached)</h2>
      <ul id="cached-room-list">
        <!-- Populated by JS from cache inventory -->
      </ul>
    </div>
  </div>
  <script>
    // List cached pages
    caches.open('hic-building-v1').then(cache => {
      cache.keys().then(keys => {
        const list = document.getElementById('cached-room-list');
        keys
          .filter(req => req.url.endsWith('/') || req.url.endsWith('.html'))
          .forEach(req => {
            const li = document.createElement('li');
            const a = document.createElement('a');
            a.href = new URL(req.url).pathname;
            a.textContent = new URL(req.url).pathname || 'Street Level';
            li.appendChild(a);
            list.appendChild(li);
          });
      });
    });
  </script>
</body>
</html>
```

---

## Animation & Transition Specifications

All animations in the HIC follow a consistent specification to maintain visual coherence and respect user preferences.

### Timing Functions

```css
:root {
  /* Standard easing curves */
  --ease-standard:    cubic-bezier(0.4, 0.0, 0.2, 1.0);   /* Most transitions */
  --ease-decelerate:  cubic-bezier(0.0, 0.0, 0.2, 1.0);   /* Elements entering */
  --ease-accelerate:  cubic-bezier(0.4, 0.0, 1.0, 1.0);   /* Elements exiting */
  --ease-sharp:       cubic-bezier(0.4, 0.0, 0.6, 1.0);   /* Elevator movement */
  --ease-bounce:      cubic-bezier(0.34, 1.56, 0.64, 1.0); /* Playful feedback */

  /* Duration scale */
  --duration-instant:  0ms;
  --duration-fast:     100ms;
  --duration-normal:   200ms;
  --duration-slow:     400ms;
  --duration-elevator: 500ms;   /* Base elevator duration */
}
```

### Transition Catalog

| Transition Name      | Duration      | Easing          | Elements Affected                |
|----------------------|---------------|-----------------|----------------------------------|
| Door open            | 300ms         | decelerate      | Room door element                |
| Room reveal          | 300ms         | standard        | Room content container           |
| Section reveal       | 300ms         | decelerate      | Content sections (staggered)     |
| Neon flash           | 200ms         | linear          | Neon signs on activation         |
| Elevator ride        | 500ms + Nms   | sharp           | Viewport during floor change     |
| Floor plan highlight | 200ms         | standard        | Room cards on hover              |
| Preview card show    | 200ms         | decelerate      | Preview card on hover/click      |
| Preview card hide    | 100ms         | accelerate      | Preview card on dismiss          |
| Search overlay open  | 250ms         | decelerate      | Search overlay container         |
| Search overlay close | 150ms         | accelerate      | Search overlay container         |
| Skeleton shimmer     | 1500ms        | ease-in-out     | Skeleton placeholders (loop)     |
| Error flicker        | 3000ms        | linear          | Error code neon sign (loop)      |

---

## State Management

The HIC maintains a global state object that tracks the user's position, history, preferences, and cached data. This state drives all navigation, rendering, and interaction logic.

### State Schema

```typescript
interface HICState {
  // Current location in the building
  location: {
    level: 'street' | 'lobby' | 'corridor' | 'room';
    floor: string | null;       // Domain ID
    room: string | null;        // Article ID
    section: string | null;     // Heading anchor
    scrollPosition: number;     // Pixels scrolled
    zoomLevel: number;          // 0-4
  };

  // Navigation history stack
  history: Array<{
    level: string;
    floor: string | null;
    room: string | null;
    section: string | null;
    timestamp: number;
  }>;

  // User preferences
  preferences: {
    theme: 'dark' | 'light' | 'system';
    reducedMotion: boolean;
    highContrast: boolean;
    fontSize: 'small' | 'medium' | 'large';
    elevatorPanelCollapsed: boolean;
  };

  // Search state
  search: {
    isOpen: boolean;
    query: string;
    scope: 'building' | 'floor';
    results: SearchResult[];
    selectedIndex: number;
    history: SearchHistoryEntry[];
  };

  // Performance state
  performance: {
    roomCache: Map<string, CachedRoom>;
    prefetchQueue: string[];
    isOffline: boolean;
  };

  // UI state
  ui: {
    isTransitioning: boolean;
    activeOverlay: string | null;
    focusedElement: string | null;
    previewCard: PreviewCardState | null;
  };
}
```

### State Transitions

```
STREET_LEVEL  --[click floor]--> LOBBY
LOBBY         --[click room]---> ROOM
LOBBY         --[click floor]--> LOBBY (different floor)
ROOM          --[click link]---> ROOM (same or different floor)
ROOM          --[click exit]---> LOBBY
ROOM          --[press Home]---> STREET_LEVEL
LOBBY         --[press Home]---> STREET_LEVEL
*             --[press /]------> SEARCH_OVERLAY (overlay, not navigation)
*             --[press Esc]----> previous state or close overlay
```

---

## Event Architecture

The HIC uses a custom event system to decouple interaction handling from navigation logic. All user interactions emit events that are consumed by the navigation controller, animation system, analytics tracker, and accessibility announcer.

### Event Types

```javascript
// Navigation events
'hic:navigate'          // User initiated navigation
'hic:navigate:start'    // Navigation transition beginning
'hic:navigate:complete' // Navigation transition finished
'hic:navigate:error'    // Navigation failed

// Interaction events
'hic:door:click'        // Door (link) clicked
'hic:door:hover'        // Door hover started
'hic:door:leave'        // Door hover ended
'hic:sign:activate'     // Neon sign clicked
'hic:floor:select'      // Floor marker clicked

// Search events
'hic:search:open'       // Search overlay opened
'hic:search:close'      // Search overlay closed
'hic:search:query'      // Search query submitted
'hic:search:select'     // Search result selected

// Lifecycle events
'hic:room:enter'        // Room content loaded and displayed
'hic:room:leave'        // Leaving current room
'hic:room:scroll'       // Scroll position changed within room
'hic:room:section'      // Active section changed (waypoint hit)

// Performance events
'hic:prefetch:start'    // Prefetch initiated
'hic:prefetch:complete' // Prefetch finished
'hic:cache:hit'         // Content served from cache
'hic:cache:miss'        // Content fetched from network
'hic:offline:detected'  // Network connection lost
'hic:online:restored'   // Network connection restored
```

### Event Dispatch Example

```javascript
function emitHICEvent(type, detail = {}) {
  const event = new CustomEvent(type, {
    bubbles: true,
    composed: true,
    detail: {
      timestamp: Date.now(),
      location: getCurrentLocation(),
      ...detail,
    },
  });
  document.dispatchEvent(event);
}

// Usage
emitHICEvent('hic:door:click', {
  targetFloor: 'api-reference',
  targetRoom: 'authentication',
  sourceElement: doorElement,
});

// Listeners
document.addEventListener('hic:door:click', (event) => {
  const { targetFloor, targetRoom } = event.detail;
  navigationController.navigate(targetFloor, targetRoom);
});

document.addEventListener('hic:room:enter', (event) => {
  const { location } = event.detail;
  announce(`Entering room: ${location.room} on floor ${location.floor}`);
  analyticsTracker.trackPageView(location);
});
```

---

## Summary

The HIC interaction model is built on a single principle: **the interface is the building**. Users do not "use a website" -- they walk through a neon-lit skyscraper. Every click is a door. Every hover is a glance. Every scroll is a step forward. Every error is a power outage.

This document defines the protocols for making that experience consistent, performant, and accessible. Implementers should treat this as the ground truth for all interaction decisions. When in doubt, ask: "What would a person do in a physical building?" Then build that.

```
                    _______________
                   |  ___________  |
                   | |  H  I  C  | |
                   | |___________| |
                   |  ___________  |
                   | |  FLOOR 6  | |
                   | |___________| |
                   |  ___________  |
                   | |  FLOOR 5  | |
                   | |___________| |
                   |  ___________  |
                   | |  FLOOR 4  | |
                   | |___________| |
                   |  ___________  |
                   | |  FLOOR 3  | |
                   | |___________| |
                   |  ___________  |
                   | |  FLOOR 2  | |
                   | |___________| |
                   |  ___________  |
                   | |  FLOOR 1  | |
                   | |___________| |
                   |  _____ _____ |
                   | | [ ENTER ] | |
                   |_|___________|_|
                   =================
                    YOU ARE HERE [X]
```
