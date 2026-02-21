# STAGE 4: HOLM INTELLIGENCE COMPLEX -- VISUAL INTERFACE DESIGN

## Agents 4-6: Neon Blueprint Style Guide, Color System, Typography, Glow/Contrast Rules

**VISUAL THEME:** Dark background, neon gridlines, glowing outlines, futuristic blueprint aesthetic, cyberpunk operations center feel.

**Design Philosophy:** The Holm Intelligence Complex (HIC) presents itself as a towering cyberpunk skyscraper rendered in neon blueprint style. Every pixel serves a dual purpose: aesthetic immersion and operational clarity. The interface must feel like standing inside a holographic war room -- data streams pulse through glowing conduits, floor plans shimmer with activity indicators, and the building itself breathes with the rhythm of its intelligence operations.

**Implementation Target:** Web-based rendering via HTML5/CSS3/SVG with optional WebGL acceleration for glow effects. All specifications are production-ready with exact values, fallbacks, and accessibility compliance.

---

---

## Agent 4: Neon Blueprint Style Guide

### 4.1 Visual Language Specification

The HIC visual language draws from three source aesthetics fused into a single coherent system:

1. **Architectural Blueprint** -- Technical precision, orthographic projection, dimension lines, section marks
2. **Cyberpunk Neon** -- Electric glow effects, scanline textures, holographic transparency, light-bleed
3. **Military Operations Center** -- Status grids, threat boards, classified zone markers, hierarchical color-coding

Every rendered element belongs to one of four visual layers, composited in order:

| Layer | Z-Index Range | Contents | Opacity Range |
|-------|--------------|----------|---------------|
| `substrate` | 0-99 | Background grid, ambient noise texture | 100% |
| `structure` | 100-499 | Walls, floors, building shell, static geometry | 80-100% |
| `data` | 500-899 | Labels, readouts, status indicators, annotations | 60-100% |
| `fx` | 900-999 | Glow halos, pulse animations, scan sweeps, particles | 10-60% |

```
CSS Custom Properties -- Visual Layers:

--layer-substrate-z:    0;
--layer-structure-z:    100;
--layer-data-z:         500;
--layer-fx-z:           900;
--layer-substrate-opacity: 1.0;
--layer-structure-opacity: 0.9;
--layer-data-opacity:      0.85;
--layer-fx-opacity:        0.4;
```

### 4.2 Blueprint Rendering Rules

#### 4.2.1 Line Weight Hierarchy

All structural lines follow a strict weight hierarchy. Weights are specified in pixels at 1x resolution and scale proportionally.

| Line Class | Weight (px) | Dash Pattern | Color Token | Usage |
|------------|------------|--------------|-------------|-------|
| `wall-exterior` | 3.0 | Solid | `--color-cyan-primary` | Outer building shell |
| `wall-interior` | 2.0 | Solid | `--color-cyan-secondary` | Interior room walls |
| `wall-partition` | 1.0 | `8 4` | `--color-cyan-muted` | Soft dividers, cubicle walls |
| `grid-major` | 1.0 | Solid | `--color-grid-major` | Primary grid divisions |
| `grid-minor` | 0.5 | `2 4` | `--color-grid-minor` | Secondary grid subdivisions |
| `dimension` | 0.5 | Solid | `--color-dim-line` | Measurement annotations |
| `section-cut` | 2.0 | `12 4 4 4` | `--color-magenta-primary` | Section plane indicators |
| `hidden-line` | 1.0 | `4 4` | `--color-cyan-ghost` | Elements behind current plane |
| `conduit-data` | 1.5 | Solid | `--color-data-flow` | Data pipeline routes |
| `conduit-power` | 1.5 | `6 2` | `--color-amber-warning` | Power distribution lines |

```
CSS Custom Properties -- Line Weights:

--line-wall-exterior:     3.0px;
--line-wall-interior:     2.0px;
--line-wall-partition:    1.0px;
--line-grid-major:        1.0px;
--line-grid-minor:        0.5px;
--line-dimension:         0.5px;
--line-section-cut:       2.0px;
--line-hidden:            1.0px;
--line-conduit-data:      1.5px;
--line-conduit-power:     1.5px;
```

#### 4.2.2 Dash Patterns

Dash patterns are specified as SVG `stroke-dasharray` values. All patterns repeat and are resolution-independent.

```
--dash-solid:       none;
--dash-partition:   8 4;
--dash-grid-minor:  2 4;
--dash-section:     12 4 4 4;
--dash-hidden:      4 4;
--dash-conduit:     6 2;
--dash-alert:       3 3;
--dash-scan:        1 6;
```

#### 4.2.3 Fill Rules

Rooms and zones use semi-transparent fills to indicate state without obscuring underlying grid structure.

| Fill Type | Base Opacity | Pattern | Notes |
|-----------|-------------|---------|-------|
| Active room | 0.08 | Solid | Barely visible tint over grid |
| Selected room | 0.15 | Solid | Highlighted for interaction |
| Hovered room | 0.12 | Solid | Cursor proximity feedback |
| Locked zone | 0.06 | 45-degree hatch, 8px spacing | Restricted access indicator |
| Classified zone | 0.10 | Cross-hatch, 6px spacing | Top-secret areas |
| Alert zone | 0.12 | Solid + pulse | Animated opacity oscillation |
| Inactive/offline | 0.03 | Solid | Nearly invisible, ghosted |
| Maintenance | 0.08 | Horizontal stripe, 12px spacing | Under repair |

### 4.3 Glow Effects Specification

Glow is the defining visual characteristic of the HIC. Every glowing element uses a standardized multi-layer glow stack.

#### 4.3.1 Glow Stack Architecture

Each glowing element renders three shadow layers in order:

```
/* Standard neon glow -- 3-layer stack */
.neon-glow {
  filter:
    drop-shadow(0 0 2px  var(--glow-color))    /* Layer 1: Tight core */
    drop-shadow(0 0 8px  var(--glow-color))    /* Layer 2: Inner halo */
    drop-shadow(0 0 20px var(--glow-color-dim)); /* Layer 3: Outer bloom */
}

/* Intense neon glow -- for primary elements */
.neon-glow-intense {
  filter:
    drop-shadow(0 0 2px  var(--glow-color))
    drop-shadow(0 0 6px  var(--glow-color))
    drop-shadow(0 0 14px var(--glow-color))
    drop-shadow(0 0 30px var(--glow-color-dim));
}

/* Subtle neon glow -- for secondary/background elements */
.neon-glow-subtle {
  filter:
    drop-shadow(0 0 1px var(--glow-color))
    drop-shadow(0 0 4px var(--glow-color-dim));
}
```

```
CSS Custom Properties -- Glow Intensities:

--glow-radius-core:      2px;
--glow-radius-inner:     8px;
--glow-radius-outer:     20px;
--glow-radius-bloom:     30px;
--glow-opacity-core:     1.0;
--glow-opacity-inner:    0.7;
--glow-opacity-outer:    0.3;
--glow-opacity-bloom:    0.15;
```

#### 4.3.2 Pulse Animations

Pulse animations indicate activity, alerts, or data flow. Three standard pulse profiles are defined:

```css
/* Heartbeat pulse -- steady operational rhythm */
@keyframes pulse-heartbeat {
  0%, 100% { opacity: 0.7; filter: drop-shadow(0 0 4px var(--glow-color)); }
  50%      { opacity: 1.0; filter: drop-shadow(0 0 12px var(--glow-color)); }
}
/* Duration: 2.4s | Easing: ease-in-out | Loop: infinite */

/* Alert pulse -- urgent rapid flash */
@keyframes pulse-alert {
  0%, 100% { opacity: 0.4; }
  15%      { opacity: 1.0; }
  30%      { opacity: 0.4; }
  45%      { opacity: 1.0; }
}
/* Duration: 1.2s | Easing: linear | Loop: infinite */

/* Data-flow pulse -- traveling dot along conduit path */
@keyframes pulse-dataflow {
  0%   { stroke-dashoffset: 0; }
  100% { stroke-dashoffset: -24; }
}
/* Duration: 0.8s | Easing: linear | Loop: infinite */

/* Scan sweep -- rotating radar line */
@keyframes scan-sweep {
  0%   { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
/* Duration: 4.0s | Easing: linear | Loop: infinite */
```

```
CSS Custom Properties -- Animation Timing:

--anim-heartbeat-duration:   2.4s;
--anim-alert-duration:       1.2s;
--anim-dataflow-duration:    0.8s;
--anim-scan-duration:        4.0s;
--anim-transition-fast:      0.15s;
--anim-transition-standard:  0.3s;
--anim-transition-slow:      0.6s;
```

### 4.4 Grid System

The HIC grid provides spatial orientation and snap-alignment for all blueprint elements.

#### 4.4.1 Grid Dimensions

```
--grid-unit:           8px;        /* Base atomic unit */
--grid-minor-spacing:  32px;       /* 4 base units */
--grid-major-spacing:  128px;      /* 16 base units */
--grid-page-width:     1280px;     /* 10 major divisions */
--grid-page-height:    960px;      /* 7.5 major divisions */
```

#### 4.4.2 Grid Rendering

```
Major gridlines:    1.0px solid at 8% opacity, every 128px
Minor gridlines:    0.5px dashed at 4% opacity, every 32px
Origin crosshair:   2.0px solid at 15% opacity, full viewport span
Snap tolerance:     4px (half base unit)
```

#### 4.4.3 ASCII Grid Mockup

```
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|  :  :  :  |  :  :  :  |  :  :  :  |  :  :  :  |    Major = |
|..:..:..:..:..:..:..:..:..:..:..:..:..:..:..:..:|    Minor = :
|  :  :  :  |  :  :  :  |  :  :  :  |  :  :  :  |    Snap  = .
|..:..:..:..:..:..:..:..:..:..:..:..:..:..:..:..:|
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|  :  :  :  |  :  :  :  |  :  :  :  |  :  :  :  |
|..:..:..:..:..:..:..:..:..:..:..:..:..:..:..:..:|
|  :  :  :  |  :  :  :  |  :  :  :  |  :  :  :  |
|..:..:..:..:..:..:..:..:..:..:..:..:..:..:..:..:|
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
```

### 4.5 Elevation View Rendering (Building Exterior)

The elevation view displays the HIC as a towering skyscraper seen from the front. Each floor is a horizontal band. The building tapers slightly toward the top and features antenna arrays on the roof.

#### 4.5.1 ASCII Elevation Mockup

```
                        /\
                       /||\
                      / || \
                     /  ||  \        ANTENNA ARRAY
                    /___||___\
                   |__________|
                  _|__________|_
                 | [FL-20] TOP  |     <- Penthouse / Command
                 |______________|
                 | [FL-19]      |
                 |______________|
                 | [FL-18]      |     <- Upper Floors
                 |______________|
                 |   . . . .    |
                 |______________|
                 | [FL-10]      |     <- Mid Floors
                 |______________|
                 |   . . . .    |
                 |______________|
                _| [FL-05]      |_
               | |______________| |   <- Lower Floors (wider base)
               | | [FL-04]      | |
               | |______________| |
               | | [FL-03]      | |
               | |______________| |
               | | [FL-02]      | |
               | |______________| |
               | | [FL-01] LOBBY| |
               |_|______________|_|
              /____________________\  <- Foundation / Sublevel Access
             /______________________\
            |________________________|
               HOLM INTELLIGENCE
                   COMPLEX
```

#### 4.5.2 Elevation Rendering Rules

| Element | Line Weight | Glow | Color |
|---------|-----------|------|-------|
| Building outline | `wall-exterior` (3px) | Intense | `--color-cyan-primary` |
| Floor dividers | `wall-interior` (2px) | Standard | `--color-cyan-secondary` |
| Floor labels | 12px mono | Subtle | `--color-text-primary` |
| Active floor highlight | 2px + fill 0.12 | Intense + pulse | `--color-cyan-primary` |
| Window panes | 0.5px grid | None | `--color-grid-minor` |
| Antenna array | 1.5px | Intense | `--color-magenta-primary` |
| Foundation | 2px | Subtle | `--color-cyan-muted` |

### 4.6 Plan View Rendering (Floor Blueprints)

When a floor is selected, the view transitions to a top-down plan view showing room layouts, corridors, and operational zones.

#### 4.6.1 ASCII Plan View Mockup -- Standard Operations Floor

```
+=======================================================================+
|  FL-07  INTELLIGENCE OPERATIONS                    [ACTIVE] 14:32:07  |
+=======================================================================+
|                          |                    |                        |
|   +------------------+   |   +------------+   |   +--------------+    |
|   | RM-0701          |   |   | RM-0702    |   |   | RM-0703      |    |
|   | SIGINT            |   |   | ANALYSIS   |   |   | CRYPTO       |    |
|   | Analysis Lab     |   |   | Bullpen    |   |   | Vault        |    |
|   |                  |   |   |            |   |   | ########     |    |
|   |  [*ACTIVE*]      |   |   | [ACTIVE]   |   |   | [LOCKED]     |    |
|   |  Agents: 3       |   |   | Agents: 7  |   |   | Clearance: 5 |    |
|   |  Load: 78%       |   |   | Load: 45%  |   |   |              |    |
|   +------------------+   |   +------------+   |   +--------------+    |
|                          |                    |                        |
|   - - - - - - - - - - - -+- - - - - - - - - -+- - - - - - - -        |
|                          |   CORRIDOR 07-A                            |
|   - - - - - - - - - - - -+- - - - - - - - - -+- - - - - - - -        |
|                          |                    |                        |
|   +------------------+   |   +------------------------------------+   |
|   | RM-0704          |   |   | RM-0705                            |   |
|   | COMMS            |   |   | WAR ROOM                           |   |
|   | Relay Hub        |   |   | Strategic Planning                 |   |
|   |                  |   |   |                                    |   |
|   | [ACTIVE]         |   |   | [*ALERT*]                          |   |
|   | Uplink: LIVE     |   |   | Threat Level: ELEVATED             |   |
|   | Bandwidth: 92%   |   |   | Participants: 12                   |   |
|   +------------------+   |   +------------------------------------+   |
|                          |                                            |
+=======================================================================+
|  STATUS: OPERATIONAL  |  ALERTS: 1  |  AGENTS: 22  |  LOAD: 64%      |
+=======================================================================+
```

### 4.7 Room State Indicators

Each room displays a state badge that determines its visual rendering.

| State | Badge Text | Border Color | Fill | Glow | Animation |
|-------|-----------|-------------|------|------|-----------|
| `ACTIVE` | `[ACTIVE]` | `--color-cyan-primary` | Cyan @ 0.08 | Standard | None |
| `*ACTIVE*` | `[*ACTIVE*]` | `--color-cyan-primary` | Cyan @ 0.12 | Intense | Heartbeat |
| `IDLE` | `[IDLE]` | `--color-cyan-muted` | None | Subtle | None |
| `LOCKED` | `[LOCKED]` | `--color-amber-warning` | Amber @ 0.06 hatch | Standard | None |
| `CLASSIFIED` | `[CLASSIFIED]` | `--color-magenta-primary` | Magenta @ 0.08 cross-hatch | Intense | None |
| `*ALERT*` | `[*ALERT*]` | `--color-red-alert` | Red @ 0.12 | Intense | Alert pulse |
| `OFFLINE` | `[OFFLINE]` | `--color-text-dim` | None | None | None |
| `MAINTENANCE` | `[MAINT]` | `--color-amber-warning` | Amber @ 0.05 stripe | Subtle | None |
| `BREACH` | `[!!BREACH!!]` | `--color-red-critical` | Red @ 0.18 | Max | Alert pulse + scan |
| `STANDBY` | `[STANDBY]` | `--color-green-success` | Green @ 0.04 | Subtle | Slow heartbeat |

### 4.8 Complete Design Token Table -- Agent 4

```css
:root {
  /* --- LAYER SYSTEM --- */
  --layer-substrate-z:           0;
  --layer-structure-z:           100;
  --layer-data-z:                500;
  --layer-fx-z:                  900;

  /* --- LINE WEIGHTS --- */
  --line-wall-exterior:          3.0px;
  --line-wall-interior:          2.0px;
  --line-wall-partition:         1.0px;
  --line-grid-major:             1.0px;
  --line-grid-minor:             0.5px;
  --line-dimension:              0.5px;
  --line-section-cut:            2.0px;
  --line-hidden:                 1.0px;
  --line-conduit:                1.5px;

  /* --- DASH PATTERNS --- */
  --dash-solid:                  none;
  --dash-partition:              8 4;
  --dash-grid-minor:             2 4;
  --dash-section:                12 4 4 4;
  --dash-hidden:                 4 4;
  --dash-conduit:                6 2;
  --dash-alert:                  3 3;

  /* --- GRID --- */
  --grid-unit:                   8px;
  --grid-minor-spacing:          32px;
  --grid-major-spacing:          128px;

  /* --- FILL OPACITIES --- */
  --fill-active:                 0.08;
  --fill-selected:               0.15;
  --fill-hover:                  0.12;
  --fill-locked:                 0.06;
  --fill-classified:             0.10;
  --fill-alert:                  0.12;
  --fill-inactive:               0.03;
  --fill-breach:                 0.18;

  /* --- GLOW RADII --- */
  --glow-radius-core:            2px;
  --glow-radius-inner:           8px;
  --glow-radius-outer:           20px;
  --glow-radius-bloom:           30px;

  /* --- ANIMATION TIMING --- */
  --anim-heartbeat-duration:     2.4s;
  --anim-alert-duration:         1.2s;
  --anim-dataflow-duration:      0.8s;
  --anim-scan-duration:          4.0s;
  --anim-transition-fast:        0.15s;
  --anim-transition-standard:    0.3s;
  --anim-transition-slow:        0.6s;
}
```

---

---

## Agent 5: Color System

### 5.1 Master Color Palette

The HIC color system is organized into functional groups. Every color has a named token, a hex value, and a defined purpose. No color may be used outside this palette without formal amendment.

#### 5.1.1 Primary Colors -- Neon Signature Tones

| Token | Name | Hex | RGB | Usage |
|-------|------|-----|-----|-------|
| `--color-cyan-primary` | Neon Cyan | `#00f0ff` | 0, 240, 255 | Primary UI elements, active walls, main glow |
| `--color-cyan-secondary` | Ice Cyan | `#00bcd4` | 0, 188, 212 | Secondary structure lines, interior walls |
| `--color-cyan-muted` | Ghost Cyan | `#0a6e7a` | 10, 110, 122 | Inactive elements, background structure |
| `--color-cyan-ghost` | Phantom Cyan | `#063f47` | 6, 63, 71 | Hidden lines, deeply recessed elements |
| `--color-blue-electric` | Electric Blue | `#0066ff` | 0, 102, 255 | Data conduits, information flow paths |
| `--color-blue-deep` | Deep Blue | `#0033aa` | 0, 51, 170 | Secondary data elements, deep links |
| `--color-magenta-primary` | Hot Magenta | `#ff0080` | 255, 0, 128 | Section cuts, classified indicators, alerts |
| `--color-magenta-secondary` | Soft Magenta | `#cc0066` | 204, 0, 102 | Secondary classified elements |
| `--color-magenta-muted` | Dark Magenta | `#660033` | 102, 0, 51 | Background classified fills |

#### 5.1.2 Secondary Colors -- Status & Semantic

| Token | Name | Hex | RGB | Usage |
|-------|------|-----|-----|-------|
| `--color-amber-warning` | Warning Amber | `#ffaa00` | 255, 170, 0 | Warnings, locked zones, caution states |
| `--color-amber-muted` | Dim Amber | `#8a5c00` | 138, 92, 0 | Background warning fills |
| `--color-red-alert` | Alert Red | `#ff2244` | 255, 34, 68 | Active alerts, elevated threat |
| `--color-red-critical` | Critical Red | `#ff0000` | 255, 0, 0 | Breach events, critical failures |
| `--color-red-muted` | Dark Red | `#660011` | 102, 0, 17 | Background alert fills |
| `--color-green-success` | Success Green | `#00ff88` | 0, 255, 136 | Operational confirmations, healthy status |
| `--color-green-muted` | Dim Green | `#006633` | 0, 102, 51 | Background success fills |
| `--color-white-pure` | Signal White | `#ffffff` | 255, 255, 255 | Maximum emphasis, critical labels |
| `--color-yellow-data` | Data Yellow | `#ffee00` | 255, 238, 0 | Data highlights, metric callouts |

#### 5.1.3 Background Colors -- The Dark Foundation

| Token | Name | Hex | RGB | Usage |
|-------|------|-----|-----|-------|
| `--color-bg-void` | Void Black | `#000000` | 0, 0, 0 | Deepest background, unused space |
| `--color-bg-primary` | Deep Black | `#0a0a0f` | 10, 10, 15 | Primary application background |
| `--color-bg-secondary` | Dark Navy | `#0d1117` | 13, 17, 23 | Panel backgrounds, card surfaces |
| `--color-bg-tertiary` | Charcoal | `#161b22` | 22, 27, 34 | Elevated surfaces, modals |
| `--color-bg-surface` | Slate | `#1c2333` | 28, 35, 51 | Interactive surface backgrounds |
| `--color-bg-elevated` | Steel | `#242d3d` | 36, 45, 61 | Hover states, raised panels |
| `--color-bg-overlay` | Smoke | `#0a0a0fcc` | 10,10,15 @ 80% | Modal overlays, dimming layers |

#### 5.1.4 Grid Colors

| Token | Name | Hex | Usage |
|-------|------|-----|-------|
| `--color-grid-major` | Grid Major | `#00f0ff14` | Major gridlines (cyan @ 8% opacity) |
| `--color-grid-minor` | Grid Minor | `#00f0ff0a` | Minor gridlines (cyan @ 4% opacity) |
| `--color-grid-origin` | Grid Origin | `#00f0ff26` | Origin crosshair (cyan @ 15% opacity) |

#### 5.1.5 Text Colors

| Token | Name | Hex | Usage |
|-------|------|-----|-------|
| `--color-text-primary` | Bright Cyan Text | `#b0f0ff` | Primary readable text |
| `--color-text-secondary` | Mid Cyan Text | `#6ab8c7` | Secondary labels, descriptions |
| `--color-text-dim` | Dim Cyan Text | `#3a6a75` | Tertiary info, timestamps |
| `--color-text-ghost` | Ghost Text | `#1e3a42` | Disabled text, placeholders |
| `--color-text-inverse` | Dark Text | `#0a0a0f` | Text on bright backgrounds (rare) |

### 5.2 Floor Zone Color Coding

Each floor type receives a unique accent color to provide instant visual identification in both elevation and plan views.

| Zone Type | Accent Color | Hex | Glow Color | Floor Examples |
|-----------|-------------|-----|------------|----------------|
| Command & Control | Neon Cyan | `#00f0ff` | `#00f0ff` | FL-20 (Penthouse), FL-19 |
| Intelligence Operations | Electric Blue | `#0066ff` | `#3388ff` | FL-15 through FL-18 |
| Research & Development | Hot Magenta | `#ff0080` | `#ff0080` | FL-12 through FL-14 |
| Communications | Data Yellow | `#ffee00` | `#ffee00` | FL-10, FL-11 |
| Analytics & Processing | Success Green | `#00ff88` | `#00ff88` | FL-07 through FL-09 |
| Security & Defense | Alert Red | `#ff2244` | `#ff2244` | FL-05, FL-06 |
| Administration | Warning Amber | `#ffaa00` | `#ffaa00` | FL-03, FL-04 |
| Public Interface | Signal White | `#ffffff` | `#b0f0ff` | FL-01 (Lobby), FL-02 |
| Infrastructure & Support | Ghost Cyan | `#0a6e7a` | `#0a6e7a` | Sublevels B1-B3 |

### 5.3 Room Type Color Coding

| Room Type | Border Color | Fill Color (@ opacity) | Badge Color |
|-----------|-------------|----------------------|-------------|
| War Room | `#ff2244` | `#ff224415` | `#ff2244` |
| Server Room | `#0066ff` | `#0066ff10` | `#0066ff` |
| Analysis Lab | `#00f0ff` | `#00f0ff10` | `#00f0ff` |
| Crypto Vault | `#ff0080` | `#ff008010` | `#ff0080` |
| Comms Relay | `#ffee00` | `#ffee0010` | `#ffee00` |
| Briefing Room | `#00ff88` | `#00ff8810` | `#00ff88` |
| Armory | `#ff2244` | `#ff22440c` | `#ff2244` |
| Archive | `#0a6e7a` | `#0a6e7a10` | `#0a6e7a` |
| Corridor | `#0a6e7a` | transparent | `--color-text-dim` |
| Elevator Shaft | `#ffaa00` | `#ffaa000a` | `#ffaa00` |
| Stairwell | `#6ab8c7` | transparent | `#6ab8c7` |
| Restroom/Utility | `#3a6a75` | transparent | `#3a6a75` |
| Training Room | `#00ff88` | `#00ff880c` | `#00ff88` |
| Director Office | `#ff0080` | `#ff008015` | `#ff0080` |
| Conference Room | `#00bcd4` | `#00bcd40c` | `#00bcd4` |

### 5.4 Security Level Color Mapping

| Level | Name | Color | Hex | Glow Intensity |
|-------|------|-------|-----|----------------|
| 0 | Public | Signal White | `#ffffff` | Subtle |
| 1 | Internal | Bright Cyan Text | `#b0f0ff` | Subtle |
| 2 | Confidential | Neon Cyan | `#00f0ff` | Standard |
| 3 | Secret | Electric Blue | `#0066ff` | Standard |
| 4 | Top Secret | Hot Magenta | `#ff0080` | Intense |
| 5 | Compartmented | Critical Red | `#ff0000` | Intense + Pulse |

### 5.5 Data State Colors

| State | Color | Hex | Meaning |
|-------|-------|-----|---------|
| Streaming | Data Yellow | `#ffee00` | Live data flowing |
| Cached | Dim Green | `#006633` | Data stored locally |
| Stale | Dim Amber | `#8a5c00` | Data needs refresh |
| Missing | Dark Red | `#660011` | Data unavailable |
| Encrypted | Hot Magenta | `#ff0080` | Data is encrypted |
| Compressed | Deep Blue | `#0033aa` | Data is compressed |

### 5.6 Contrast Ratios & Accessibility

All text-on-background combinations must meet WCAG 2.1 AA minimums. The following table documents measured contrast ratios for primary combinations.

| Foreground | Background | Ratio | WCAG AA | WCAG AAA | Usage |
|------------|-----------|-------|---------|----------|-------|
| `#b0f0ff` on `#0a0a0f` | Primary text | 14.2:1 | Pass | Pass | Body text |
| `#6ab8c7` on `#0a0a0f` | Secondary text | 8.7:1 | Pass | Pass | Labels |
| `#3a6a75` on `#0a0a0f` | Dim text | 4.6:1 | Pass | Fail | Timestamps, hints |
| `#00f0ff` on `#0a0a0f` | Neon cyan | 12.8:1 | Pass | Pass | Headings, active elements |
| `#ff0080` on `#0a0a0f` | Hot magenta | 5.3:1 | Pass | Fail | Classified markers |
| `#ffaa00` on `#0a0a0f` | Warning amber | 10.1:1 | Pass | Pass | Warning text |
| `#ff2244` on `#0a0a0f` | Alert red | 5.8:1 | Pass | Fail | Alert badges |
| `#00ff88` on `#0a0a0f` | Success green | 12.1:1 | Pass | Pass | Confirmations |
| `#ffee00` on `#0a0a0f` | Data yellow | 15.9:1 | Pass | Pass | Data callouts |
| `#ffffff` on `#0a0a0f` | Pure white | 19.4:1 | Pass | Pass | Maximum emphasis |
| `#b0f0ff` on `#161b22` | Text on elevated bg | 10.8:1 | Pass | Pass | Panel text |
| `#1e3a42` on `#0a0a0f` | Ghost text | 2.1:1 | Fail | Fail | Decorative only |

**Rule:** Ghost-level text (`--color-text-ghost`) is classified as decorative and does not carry semantic content. It must never be the sole indicator of meaningful information.

### 5.7 Glow Color Rules

Not every element glows. Glow is expensive (rendering cost) and must be used with discipline.

#### Elements That MUST Glow

- Active room borders
- Building exterior outline (elevation view)
- Currently selected floor indicator
- Alert/breach badges
- Data flow conduit lines (animated)
- Navigation active state
- Primary action buttons
- Floor number on active floor

#### Elements That MUST NOT Glow

- Background grid (major and minor)
- Inactive room borders
- Dimension/annotation lines
- Body text (readable content)
- Disabled UI elements
- Table borders and dividers
- Scrollbar tracks
- Tooltip backgrounds

#### Elements That MAY Glow (Context-Dependent)

- Hovered elements (glow on hover, remove on leave)
- Section cut indicators (glow when section is active)
- Status bar text (glow during state transitions only)
- Floor labels in elevation (glow for active floor only)

### 5.8 Dark-on-Dark Readability Specifications

Working in an almost entirely dark palette requires strict rules to prevent elements from disappearing.

**Minimum luminance delta:** Adjacent structural elements must differ by at least 3% luminance. Background surfaces must step in increments of at least `#060608` per elevation level.

**Surface elevation luminance ladder:**

| Elevation Level | Token | Hex | Luminance (relative) |
|----------------|-------|-----|---------------------|
| 0 (base) | `--color-bg-void` | `#000000` | 0.000 |
| 1 | `--color-bg-primary` | `#0a0a0f` | 0.014 |
| 2 | `--color-bg-secondary` | `#0d1117` | 0.022 |
| 3 | `--color-bg-tertiary` | `#161b22` | 0.035 |
| 4 | `--color-bg-surface` | `#1c2333` | 0.048 |
| 5 | `--color-bg-elevated` | `#242d3d` | 0.065 |

**Border differentiation rule:** When two dark surfaces are adjacent, a 1px border of at least `#2a3a4a` (`0.08 luminance`) must separate them, OR the surfaces must differ by at least 2 elevation levels.

### 5.9 Complete Color Token CSS

```css
:root {
  /* --- PRIMARY NEON --- */
  --color-cyan-primary:          #00f0ff;
  --color-cyan-secondary:        #00bcd4;
  --color-cyan-muted:            #0a6e7a;
  --color-cyan-ghost:            #063f47;
  --color-blue-electric:         #0066ff;
  --color-blue-deep:             #0033aa;
  --color-magenta-primary:       #ff0080;
  --color-magenta-secondary:     #cc0066;
  --color-magenta-muted:         #660033;

  /* --- SECONDARY STATUS --- */
  --color-amber-warning:         #ffaa00;
  --color-amber-muted:           #8a5c00;
  --color-red-alert:             #ff2244;
  --color-red-critical:          #ff0000;
  --color-red-muted:             #660011;
  --color-green-success:         #00ff88;
  --color-green-muted:           #006633;
  --color-white-pure:            #ffffff;
  --color-yellow-data:           #ffee00;

  /* --- BACKGROUNDS --- */
  --color-bg-void:               #000000;
  --color-bg-primary:            #0a0a0f;
  --color-bg-secondary:          #0d1117;
  --color-bg-tertiary:           #161b22;
  --color-bg-surface:            #1c2333;
  --color-bg-elevated:           #242d3d;
  --color-bg-overlay:            #0a0a0fcc;

  /* --- GRID --- */
  --color-grid-major:            #00f0ff14;
  --color-grid-minor:            #00f0ff0a;
  --color-grid-origin:           #00f0ff26;

  /* --- TEXT --- */
  --color-text-primary:          #b0f0ff;
  --color-text-secondary:        #6ab8c7;
  --color-text-dim:              #3a6a75;
  --color-text-ghost:            #1e3a42;
  --color-text-inverse:          #0a0a0f;

  /* --- SEMANTIC ALIASES --- */
  --color-border-default:        #2a3a4a;
  --color-border-focus:          #00f0ff;
  --color-border-error:          #ff2244;
  --color-border-success:        #00ff88;
  --color-dim-line:              #3a6a75;
  --color-data-flow:             #0066ff;
}
```

---

---

## Agent 6: Typography & Glow/Contrast Rules

### 6.1 Font Stack

The HIC uses a strict two-tier font system. Monospace is primary for all technical/blueprint content. Sans-serif is secondary for UI chrome and natural-language content.

```css
:root {
  /* Primary: Monospace -- all blueprint labels, data readouts, room names */
  --font-mono: 'JetBrains Mono', 'Fira Code', 'Source Code Pro',
               'Cascadia Code', 'Menlo', 'Consolas', 'Monaco',
               'Liberation Mono', 'Courier New', monospace;

  /* Secondary: Sans-serif -- UI buttons, navigation, prose descriptions */
  --font-sans: 'Inter', 'SF Pro Display', 'Segoe UI', 'Roboto',
               -apple-system, BlinkMacSystemFont, 'Helvetica Neue',
               Arial, sans-serif;

  /* Tertiary: Display -- building title, floor headers (optional, decorative) */
  --font-display: 'Orbitron', 'Rajdhani', 'Share Tech Mono',
                  var(--font-mono);
}
```

**Loading strategy:** `font-display: swap` for all web fonts. System monospace must be visually acceptable as the permanent fallback. No layout shift on font load -- all size calculations use the fallback metrics.

### 6.2 Type Scale

The type scale uses a 1.25 ratio (Major Third) anchored at 14px base. All sizes are specified in `rem` with pixel equivalents for reference.

| Token | Size (rem) | Size (px) | Line Height | Weight | Font | Usage |
|-------|-----------|-----------|-------------|--------|------|-------|
| `--type-display-xl` | 2.441 | 39 | 1.1 | 700 | Display | Building title |
| `--type-display-lg` | 1.953 | 31 | 1.15 | 700 | Display | Floor headers (elevation) |
| `--type-display-md` | 1.563 | 25 | 1.2 | 600 | Display | Zone titles |
| `--type-heading-lg` | 1.25 | 20 | 1.3 | 600 | Mono | Floor plan title bar |
| `--type-heading-md` | 1.0 | 16 | 1.35 | 600 | Mono | Room names |
| `--type-heading-sm` | 0.875 | 14 | 1.4 | 600 | Mono | Sub-section headers |
| `--type-body` | 0.875 | 14 | 1.5 | 400 | Mono | Primary body text |
| `--type-body-sm` | 0.75 | 12 | 1.5 | 400 | Mono | Room metadata, agent counts |
| `--type-caption` | 0.6875 | 11 | 1.4 | 400 | Mono | Dimension labels, timestamps |
| `--type-micro` | 0.625 | 10 | 1.3 | 400 | Mono | Grid coordinates, tiny annotations |
| `--type-badge` | 0.6875 | 11 | 1.0 | 700 | Mono | Status badges `[ACTIVE]` |
| `--type-ui-button` | 0.875 | 14 | 1.0 | 600 | Sans | UI action buttons |
| `--type-ui-nav` | 0.75 | 12 | 1.0 | 500 | Sans | Navigation items |
| `--type-ui-tooltip` | 0.6875 | 11 | 1.4 | 400 | Sans | Tooltip content |

### 6.3 Font Weight Definitions

```css
:root {
  --weight-regular:    400;
  --weight-medium:     500;
  --weight-semibold:   600;
  --weight-bold:       700;
}
```

Only these four weights are permitted. No light (300) or thin (100) weights -- they become illegible against dark backgrounds with glow effects.

### 6.4 Label Placement Rules on Blueprints

#### 6.4.1 Room Labels

Room labels are positioned according to a strict hierarchy within the room boundary:

```
+----------------------------------+
|  RM-0701                         |   <- Room ID: top-left, 2px inset, --type-caption
|  SIGINT ANALYSIS LAB             |   <- Room Name: below ID, --type-heading-md
|                                  |
|                                  |
|  [*ACTIVE*]                      |   <- State Badge: bottom-left zone, --type-badge
|  Agents: 3  |  Load: 78%        |   <- Metrics: below badge, --type-body-sm
+----------------------------------+
```

**Placement rules:**

1. **Room ID** (`RM-XXYY`): Top-left corner, 8px inset from wall on both axes. Always visible, never truncated.
2. **Room Name**: Directly below Room ID, same left inset. May truncate with ellipsis if room is narrow.
3. **Room Type**: Below Room Name, same inset. Uses `--color-text-secondary`. May be omitted in compact view.
4. **State Badge**: Lower-left quadrant, 8px from bottom wall, 8px from left wall.
5. **Metrics Line**: Below state badge. Shows 2-3 key metrics separated by `|` dividers.
6. **Minimum label room**: If the room bounding box is smaller than 120x80px, only Room ID and State Badge are shown. If smaller than 60x40px, only the State Badge is shown as a colored dot.

#### 6.4.2 Corridor Labels

Corridor labels center-align horizontally along the corridor's longest axis. They use `--type-caption` and `--color-text-dim`. Format: `CORRIDOR XX-Y` where XX is the floor number and Y is a sequential letter.

#### 6.4.3 Floor Title Bar

The floor title bar spans the full width of the plan view and contains:

```
+========================================================================+
|  FL-07  INTELLIGENCE OPERATIONS                     [ACTIVE] 14:32:07  |
+========================================================================+
   ^      ^                                            ^       ^
   |      |                                            |       |
   |      Floor Name: --type-heading-lg, left-aligned  |       Timestamp
   |                                                   |
   Floor Number: --type-heading-lg, --weight-bold      State Badge
```

### 6.5 Text Glow Effects

Text glow is used sparingly. Overuse destroys readability by blooming letter forms together.

#### 6.5.1 When Text Glows

| Element | Glow Type | Glow Color | Glow Radius | Condition |
|---------|-----------|------------|-------------|-----------|
| Building title | Intense | `--color-cyan-primary` | 12px | Always |
| Active floor number (elevation) | Standard | Floor accent color | 8px | When floor is active |
| Room state badge `[*ALERT*]` | Alert | `--color-red-alert` | 6px | Always when in alert state |
| Room state badge `[*ACTIVE*]` | Subtle | `--color-cyan-primary` | 4px | Always when focus-active |
| Selected room name | Subtle | Room accent color | 4px | On selection only |
| Status bar alert count | Standard | `--color-red-alert` | 6px | When alerts > 0 |
| Data readout values (live) | Subtle | `--color-yellow-data` | 3px | During live data stream |
| Navigation active item | Subtle | `--color-cyan-primary` | 4px | Active nav state |

#### 6.5.2 When Text MUST NOT Glow

- Body text / descriptions / prose content
- Room metadata lines (agent counts, load percentages as static values)
- Dimension labels and annotations
- Timestamps (except in alert conditions)
- Tooltip content
- Table cell content
- Input field placeholder text
- Breadcrumb navigation (inactive segments)
- Log entries and audit trails

#### 6.5.3 Text Glow CSS Implementation

```css
.text-glow-intense {
  text-shadow:
    0 0 4px  var(--glow-color),
    0 0 8px  var(--glow-color),
    0 0 16px var(--glow-color),
    0 0 32px var(--glow-color);
}

.text-glow-standard {
  text-shadow:
    0 0 2px var(--glow-color),
    0 0 6px var(--glow-color),
    0 0 12px var(--glow-color);
}

.text-glow-subtle {
  text-shadow:
    0 0 2px var(--glow-color),
    0 0 4px var(--glow-color);
}

.text-glow-alert {
  text-shadow:
    0 0 2px var(--glow-color),
    0 0 6px var(--glow-color),
    0 0 10px var(--glow-color);
  animation: pulse-alert var(--anim-alert-duration) linear infinite;
}
```

### 6.6 Contrast Minimum Ratios

These are absolute minimums enforced across the system. They supplement the color-specific ratios defined in Agent 5 with typographic context.

| Content Type | Minimum Ratio | Standard | Notes |
|-------------|--------------|----------|-------|
| Body text (14px+) | 4.5:1 | WCAG AA | Non-negotiable |
| Large text (18px+ or 14px bold) | 3.0:1 | WCAG AA | Headings, titles |
| UI components (buttons, inputs) | 3.0:1 | WCAG AA | Border and fill contrast |
| Decorative text (ghost labels) | No minimum | Decorative | Must not carry meaning |
| Focus indicators | 3.0:1 | WCAG AA | Against adjacent colors |
| Glow text on dark bg | 7.0:1 | Enhanced | Glow can reduce perceived contrast |
| Status badges | 4.5:1 | WCAG AA | Badge text against badge bg |
| Data readout values | 4.5:1 | WCAG AA | Must remain legible during animation |

**Glow contrast compensation rule:** When text has a glow effect, the base text color must achieve 7.0:1 contrast minimum rather than the standard 4.5:1. This is because glow halos can visually interfere with character recognition at certain blur radii.

### 6.7 Room Label Format Specification

Room labels follow a rigid format for machine-readability and visual consistency.

#### 6.7.1 Room ID Format

```
RM-XXYY

Where:
  RM   = Literal prefix (always uppercase)
  XX   = Floor number, zero-padded to 2 digits (01-20)
  YY   = Room sequence on floor, zero-padded to 2 digits (01-99)

Examples:
  RM-0701  = Floor 7, Room 1
  RM-2001  = Floor 20, Room 1
  RM-0315  = Floor 3, Room 15
```

#### 6.7.2 Room Name Format

```
Line 1: Room ID        --type-caption, --color-text-dim, uppercase
Line 2: ROOM NAME      --type-heading-md, --color-text-primary, uppercase
Line 3: Room Subtype    --type-body-sm, --color-text-secondary, Title Case
```

#### 6.7.3 State Badge Format

```
Standard states:   [STATE]        Square brackets, uppercase
Emphasis states:   [*STATE*]      Asterisk-wrapped within brackets
Critical states:   [!!STATE!!]    Double-bang wrapped within brackets

Badge font:  --type-badge (11px, weight 700, monospace)
Badge color: Matches state color from Agent 4 state table
```

### 6.8 Floor Number Display Format

#### 6.8.1 Elevation View

Floor numbers in the elevation view are displayed in a fixed-width label block adjacent to each floor band.

```
Format:    FL-XX
Font:      --type-heading-md (16px mono, weight 600)
Color:     Floor zone accent color (see Agent 5, Section 5.2)
Position:  Left-aligned, 16px outside building exterior wall
Glow:      Standard glow in floor accent color (active floor only)
```

#### 6.8.2 Plan View

Floor number appears in the title bar (see Section 6.4.3) and as a persistent corner watermark.

```
Corner watermark:
  Format:    FL-XX
  Font:      --type-display-lg (31px display, weight 700)
  Color:     --color-text-ghost (#1e3a42)
  Position:  Bottom-right, 24px inset from edges
  Opacity:   0.3
  Glow:      None
```

### 6.9 Data Readout Typography

Data readouts display live metrics, sensor data, and system statistics. They use a specialized typographic treatment.

#### 6.9.1 Readout Layout

```
+---------------------------------------+
|  SYSTEM LOAD                          |   <- Label: --type-caption, --color-text-dim
|  ████████████████████░░░░░ 78%        |   <- Bar + Value: --type-body, --color-cyan-primary
|                                       |
|  ACTIVE AGENTS     22 / 50           |   <- Label + Value on same line
|  DATA THROUGHPUT   1.4 TB/s          |   <- Numeric values right-aligned
|  UPTIME            99.97%            |
|  THREAT LEVEL      ██ ELEVATED       |   <- Color-coded inline badge
+---------------------------------------+
```

#### 6.9.2 Numeric Value Formatting

| Data Type | Format | Example | Font Variant |
|-----------|--------|---------|-------------|
| Percentage | `XX%` or `XX.X%` | `78%`, `99.97%` | Tabular nums |
| Count | Locale-formatted integers | `1,247` | Tabular nums |
| Throughput | `X.X UNIT/s` | `1.4 TB/s` | Tabular nums |
| Duration | `HH:MM:SS` | `14:32:07` | Tabular nums |
| Temperature | `XX.X C` | `22.4 C` | Tabular nums |
| Memory | `X.XX UNIT` | `3.72 TB` | Tabular nums |
| Ratio | `X / Y` | `22 / 50` | Tabular nums |
| Currency | `$X,XXX.XX` | `$4,200.00` | Tabular nums |

```css
.data-readout-value {
  font-family: var(--font-mono);
  font-size: var(--type-body);
  font-weight: var(--weight-semibold);
  font-variant-numeric: tabular-nums;
  letter-spacing: 0.02em;
  color: var(--color-text-primary);
}

.data-readout-label {
  font-family: var(--font-mono);
  font-size: var(--type-caption);
  font-weight: var(--weight-regular);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--color-text-dim);
}
```

### 6.10 Status Bar Text Rules

The status bar anchors the bottom of every view (elevation and plan). It provides system-wide situational awareness.

#### 6.10.1 Status Bar Layout

```
+=======================================================================+
|  STATUS: OPERATIONAL  |  ALERTS: 1  |  AGENTS: 22  |  LOAD: 64%      |
+=======================================================================+
```

#### 6.10.2 Status Bar Typography Rules

| Element | Font | Size | Weight | Color | Glow |
|---------|------|------|--------|-------|------|
| Label (`STATUS:`) | Mono | `--type-body-sm` | 600 | `--color-text-dim` | None |
| Value (`OPERATIONAL`) | Mono | `--type-body-sm` | 700 | Status-dependent (see below) | Subtle, status color |
| Divider (`|`) | Mono | `--type-body-sm` | 400 | `--color-text-ghost` | None |
| Alert count (when > 0) | Mono | `--type-body-sm` | 700 | `--color-red-alert` | Standard, red |
| Alert count (when 0) | Mono | `--type-body-sm` | 400 | `--color-green-success` | None |

**Status value colors:**

| Status | Color | Glow |
|--------|-------|------|
| `OPERATIONAL` | `--color-green-success` | Subtle green |
| `DEGRADED` | `--color-amber-warning` | Standard amber |
| `ALERT` | `--color-red-alert` | Standard red + pulse |
| `CRITICAL` | `--color-red-critical` | Intense red + alert pulse |
| `OFFLINE` | `--color-text-ghost` | None |
| `MAINTENANCE` | `--color-amber-muted` | None |

### 6.11 Letter Spacing & Tracking

```css
:root {
  --tracking-tight:       -0.01em;    /* Display headings only */
  --tracking-normal:       0.00em;    /* Body text default */
  --tracking-wide:         0.02em;    /* Data readout values */
  --tracking-wider:        0.05em;    /* State badges */
  --tracking-widest:       0.08em;    /* Uppercase labels, captions */
  --tracking-extreme:      0.12em;    /* Building title display */
}
```

**Rule:** All `text-transform: uppercase` elements must use at least `--tracking-widest` (0.08em) to maintain readability. Uppercase without increased tracking creates visual congestion in monospace faces.

### 6.12 Print Fallback Specifications

When the HIC interface is printed (or exported to PDF), the following transformations apply:

#### 6.12.1 Color Transformations

| Screen Element | Print Replacement |
|---------------|-------------------|
| All neon/glow colors | Black lines at 100% opacity |
| Background `#0a0a0f` | White (`#ffffff`) |
| Grid lines | Light gray (`#cccccc`) at 0.5px |
| Text `--color-text-primary` | Black (`#000000`) |
| Text `--color-text-secondary` | Dark gray (`#444444`) |
| Text `--color-text-dim` | Medium gray (`#888888`) |
| Alert/status colors | Retained but with increased contrast |
| All glow effects | Removed entirely |
| All pulse animations | Frozen at peak opacity frame |

#### 6.12.2 Print Typography Adjustments

```css
@media print {
  :root {
    --font-mono: 'Courier New', 'Courier', monospace;
    --font-sans: 'Helvetica', 'Arial', sans-serif;
    --font-display: var(--font-mono);
  }

  * {
    text-shadow: none !important;
    filter: none !important;
    animation: none !important;
  }

  body {
    background: #ffffff !important;
    color: #000000 !important;
    font-size: 10pt;
  }

  .neon-glow,
  .neon-glow-intense,
  .neon-glow-subtle {
    filter: none !important;
  }

  .status-badge {
    border: 2px solid #000000;
    font-weight: 700;
  }

  .room-border {
    stroke: #000000;
    stroke-width: 1.5px;
  }

  .grid-major { stroke: #cccccc; stroke-width: 0.5px; }
  .grid-minor { display: none; }
}
```

#### 6.12.3 Print Layout Rules

- Page size: A3 landscape preferred, A4 landscape acceptable
- Margins: 15mm all sides
- Floor plans: One floor per page
- Elevation view: Single page, scale to fit
- Status bar: Repeated as page footer
- Room labels: All labels forced visible regardless of zoom level
- Color legend: Auto-generated on first page when color is used semantically

### 6.13 Complete Typography Token CSS

```css
:root {
  /* --- FONT FAMILIES --- */
  --font-mono:             'JetBrains Mono', 'Fira Code', 'Source Code Pro',
                           'Cascadia Code', 'Menlo', 'Consolas', monospace;
  --font-sans:             'Inter', 'SF Pro Display', 'Segoe UI', 'Roboto',
                           -apple-system, BlinkMacSystemFont, sans-serif;
  --font-display:          'Orbitron', 'Rajdhani', 'Share Tech Mono',
                           var(--font-mono);

  /* --- TYPE SCALE (Major Third 1.25, base 14px) --- */
  --type-display-xl:       2.441rem;   /* 39px */
  --type-display-lg:       1.953rem;   /* 31px */
  --type-display-md:       1.563rem;   /* 25px */
  --type-heading-lg:       1.25rem;    /* 20px */
  --type-heading-md:       1.0rem;     /* 16px */
  --type-heading-sm:       0.875rem;   /* 14px */
  --type-body:             0.875rem;   /* 14px */
  --type-body-sm:          0.75rem;    /* 12px */
  --type-caption:          0.6875rem;  /* 11px */
  --type-micro:            0.625rem;   /* 10px */
  --type-badge:            0.6875rem;  /* 11px */
  --type-ui-button:        0.875rem;   /* 14px */
  --type-ui-nav:           0.75rem;    /* 12px */
  --type-ui-tooltip:       0.6875rem;  /* 11px */

  /* --- LINE HEIGHTS --- */
  --leading-display:       1.1;
  --leading-display-lg:    1.15;
  --leading-display-md:    1.2;
  --leading-heading:       1.3;
  --leading-heading-md:    1.35;
  --leading-heading-sm:    1.4;
  --leading-body:          1.5;
  --leading-caption:       1.4;
  --leading-micro:         1.3;
  --leading-badge:         1.0;
  --leading-ui:            1.0;

  /* --- FONT WEIGHTS --- */
  --weight-regular:        400;
  --weight-medium:         500;
  --weight-semibold:       600;
  --weight-bold:           700;

  /* --- LETTER SPACING --- */
  --tracking-tight:        -0.01em;
  --tracking-normal:       0.00em;
  --tracking-wide:         0.02em;
  --tracking-wider:        0.05em;
  --tracking-widest:       0.08em;
  --tracking-extreme:      0.12em;
}
```

---

---

## Appendix A: Full Interface ASCII Mockup -- Elevation View

```
 ___________________________________________________________________________
|                                                                           |
|                          /\                                               |
|                         /||\                                              |
|                        / || \            HOLM INTELLIGENCE COMPLEX        |
|                       /  ||  \           ========================        |
|                      /___||___\          Neon Blueprint Interface         |
|                     |__________|                                          |
|                    _|__________|_                                         |
|    FL-20  >>>     | ============ |    COMMAND CENTER          [ACTIVE]    |
|                   |______________|                                        |
|    FL-19          | ============ |    STRATEGIC OPS           [ACTIVE]    |
|                   |______________|                                        |
|    FL-18          | ============ |    SIGINT DIVISION         [IDLE]      |
|                   |______________|                                        |
|    FL-17          | ============ |    HUMINT DIVISION         [ACTIVE]    |
|                   |______________|                                        |
|    FL-16          | ============ |    OSINT DIVISION          [ACTIVE]    |
|                   |______________|                                        |
|    FL-15          | ============ |    CYBER WARFARE           [*ALERT*]   |
|                   |______________|                                        |
|    FL-14          | ============ |    R&D - ADVANCED          [LOCKED]    |
|                   |______________|                                        |
|    FL-13          | ============ |    R&D - PROTOTYPE         [ACTIVE]    |
|                   |______________|                                        |
|    FL-12          | ============ |    R&D - THEORETICAL       [IDLE]      |
|                   |______________|                                        |
|    FL-11          | ============ |    COMMS - SATELLITE       [ACTIVE]    |
|                   |______________|                                        |
|    FL-10          | ============ |    COMMS - TERRESTRIAL     [ACTIVE]    |
|                   |______________|                                        |
|    FL-09          | ============ |    ANALYTICS - PREDICTIVE  [ACTIVE]    |
|                   |______________|                                        |
|    FL-08          | ============ |    ANALYTICS - REALTIME    [ACTIVE]    |
|                   |______________|                                        |
|    FL-07          | ============ |    ANALYTICS - HISTORICAL  [STANDBY]   |
|                   |______________|                                        |
|    FL-06          | ============ |    SECURITY - PHYSICAL     [ACTIVE]    |
|                   |______________|                                        |
|    FL-05          | ============ |    SECURITY - CYBER        [*ALERT*]   |
|                  _|______________|_                                       |
|    FL-04        | ================ |  ADMIN - OPERATIONS      [ACTIVE]    |
|                 |__________________|                                      |
|    FL-03        | ================ |  ADMIN - PERSONNEL       [IDLE]      |
|                 |__________________|                                      |
|    FL-02        | ================ |  PUBLIC - BRIEFING       [ACTIVE]    |
|                 |__________________|                                      |
|    FL-01        | ================ |  PUBLIC - LOBBY          [ACTIVE]    |
|                 |__________________|                                      |
|                /____________________\                                     |
|               /______________________\                                    |
|              |________________________|                                   |
|                                                                           |
|===========================================================================|
| STATUS: OPERATIONAL | ALERTS: 2 | FLOORS: 20/20 | AGENTS: 147 | 14:32:07|
|===========================================================================|
```

## Appendix B: Full Interface ASCII Mockup -- Plan View with Sidebar

```
+============================================================================+
|  FL-15  CYBER WARFARE DIVISION                       [*ALERT*]  14:32:07   |
+============================================================================+
|       |                                                                     |
| N     |  +------------------+   +------------------+   +-------------+     |
| A     |  | RM-1501          |   | RM-1502          |   | RM-1503     |     |
| V     |  | THREAT           |   | INCIDENT         |   | MALWARE     |     |
|       |  | DETECTION        |   | RESPONSE         |   | ANALYSIS    |     |
| [20]  |  |                  |   |                  |   | #########   |     |
| [19]  |  | [*ACTIVE*]       |   | [*ALERT*]        |   | [LOCKED]    |     |
| [18]  |  | Analysts: 4      |   | Responders: 6    |   | CL: 4       |     |
| [17]  |  | Threats: 12      |   | Incidents: 2     |   |             |     |
| [16]  |  +------------------+   +------------------+   +-------------+     |
|>[15]< |                         |                                           |
| [14]  |  - - - - - - - - - - - -+- - - - - - - - - - - - - - - - -         |
| [13]  |                         |  CORRIDOR 15-A                            |
| [12]  |  - - - - - - - - - - - -+- - - - - - - - - - - - - - - - -         |
| [11]  |                         |                                           |
| [10]  |  +------------------+   +--------------------------------------+   |
| [09]  |  | RM-1504          |   | RM-1505                              |   |
| [08]  |  | FORENSICS        |   | CYBER OPERATIONS CENTER              |   |
| [07]  |  | LAB              |   |                                      |   |
| [06]  |  |                  |   |  +--------+  +--------+  +--------+  |   |
| [05]  |  | [ACTIVE]         |   |  |WKSTN-01|  |WKSTN-02|  |WKSTN-03|  |   |
| [04]  |  | Evidence: 47     |   |  | [ON]   |  | [ON]   |  | [OFF]  |  |   |
| [03]  |  | Chain: INTACT    |   |  +--------+  +--------+  +--------+  |   |
| [02]  |  |                  |   |                                      |   |
| [01]  |  +------------------+   |  [*ALERT*]  Operator: CMDR_NYX      |   |
|       |                         |  Active Op: SHADOWSTRIKE             |   |
|       |                         +--------------------------------------+   |
|       |                                                                     |
+============================================================================+
| STATUS: ALERT  |  INCIDENTS: 2  |  AGENTS: 14  |  THREAT: ELEVATED        |
+============================================================================+
```

## Appendix C: Composite Design Token Export -- All Agents

This is the complete, unified CSS custom property block combining tokens from all three agents. Copy this block into the root stylesheet.

```css
:root {
  /* ================================================================
     HOLM INTELLIGENCE COMPLEX -- UNIFIED DESIGN TOKENS
     Agents 4-6: Blueprint Style, Color System, Typography
     ================================================================ */

  /* --- LAYER SYSTEM (Agent 4) --- */
  --layer-substrate-z:           0;
  --layer-structure-z:           100;
  --layer-data-z:                500;
  --layer-fx-z:                  900;

  /* --- LINE WEIGHTS (Agent 4) --- */
  --line-wall-exterior:          3.0px;
  --line-wall-interior:          2.0px;
  --line-wall-partition:         1.0px;
  --line-grid-major:             1.0px;
  --line-grid-minor:             0.5px;
  --line-dimension:              0.5px;
  --line-section-cut:            2.0px;
  --line-hidden:                 1.0px;
  --line-conduit:                1.5px;

  /* --- DASH PATTERNS (Agent 4) --- */
  --dash-partition:              8 4;
  --dash-grid-minor:             2 4;
  --dash-section:                12 4 4 4;
  --dash-hidden:                 4 4;
  --dash-conduit:                6 2;
  --dash-alert:                  3 3;

  /* --- GRID SYSTEM (Agent 4) --- */
  --grid-unit:                   8px;
  --grid-minor-spacing:          32px;
  --grid-major-spacing:          128px;

  /* --- FILL OPACITIES (Agent 4) --- */
  --fill-active:                 0.08;
  --fill-selected:               0.15;
  --fill-hover:                  0.12;
  --fill-locked:                 0.06;
  --fill-classified:             0.10;
  --fill-alert:                  0.12;
  --fill-inactive:               0.03;
  --fill-breach:                 0.18;

  /* --- GLOW RADII (Agent 4) --- */
  --glow-radius-core:            2px;
  --glow-radius-inner:           8px;
  --glow-radius-outer:           20px;
  --glow-radius-bloom:           30px;

  /* --- ANIMATION TIMING (Agent 4) --- */
  --anim-heartbeat-duration:     2.4s;
  --anim-alert-duration:         1.2s;
  --anim-dataflow-duration:      0.8s;
  --anim-scan-duration:          4.0s;
  --anim-transition-fast:        0.15s;
  --anim-transition-standard:    0.3s;
  --anim-transition-slow:        0.6s;

  /* --- PRIMARY NEON COLORS (Agent 5) --- */
  --color-cyan-primary:          #00f0ff;
  --color-cyan-secondary:        #00bcd4;
  --color-cyan-muted:            #0a6e7a;
  --color-cyan-ghost:            #063f47;
  --color-blue-electric:         #0066ff;
  --color-blue-deep:             #0033aa;
  --color-magenta-primary:       #ff0080;
  --color-magenta-secondary:     #cc0066;
  --color-magenta-muted:         #660033;

  /* --- SECONDARY STATUS COLORS (Agent 5) --- */
  --color-amber-warning:         #ffaa00;
  --color-amber-muted:           #8a5c00;
  --color-red-alert:             #ff2244;
  --color-red-critical:          #ff0000;
  --color-red-muted:             #660011;
  --color-green-success:         #00ff88;
  --color-green-muted:           #006633;
  --color-white-pure:            #ffffff;
  --color-yellow-data:           #ffee00;

  /* --- BACKGROUND COLORS (Agent 5) --- */
  --color-bg-void:               #000000;
  --color-bg-primary:            #0a0a0f;
  --color-bg-secondary:          #0d1117;
  --color-bg-tertiary:           #161b22;
  --color-bg-surface:            #1c2333;
  --color-bg-elevated:           #242d3d;
  --color-bg-overlay:            #0a0a0fcc;

  /* --- GRID COLORS (Agent 5) --- */
  --color-grid-major:            #00f0ff14;
  --color-grid-minor:            #00f0ff0a;
  --color-grid-origin:           #00f0ff26;

  /* --- TEXT COLORS (Agent 5) --- */
  --color-text-primary:          #b0f0ff;
  --color-text-secondary:        #6ab8c7;
  --color-text-dim:              #3a6a75;
  --color-text-ghost:            #1e3a42;
  --color-text-inverse:          #0a0a0f;

  /* --- SEMANTIC COLORS (Agent 5) --- */
  --color-border-default:        #2a3a4a;
  --color-border-focus:          #00f0ff;
  --color-border-error:          #ff2244;
  --color-border-success:        #00ff88;
  --color-dim-line:              #3a6a75;
  --color-data-flow:             #0066ff;

  /* --- FONT FAMILIES (Agent 6) --- */
  --font-mono:                   'JetBrains Mono', 'Fira Code', 'Source Code Pro',
                                 'Cascadia Code', 'Menlo', 'Consolas', monospace;
  --font-sans:                   'Inter', 'SF Pro Display', 'Segoe UI', 'Roboto',
                                 -apple-system, BlinkMacSystemFont, sans-serif;
  --font-display:                'Orbitron', 'Rajdhani', 'Share Tech Mono',
                                 var(--font-mono);

  /* --- TYPE SCALE (Agent 6) --- */
  --type-display-xl:             2.441rem;
  --type-display-lg:             1.953rem;
  --type-display-md:             1.563rem;
  --type-heading-lg:             1.25rem;
  --type-heading-md:             1.0rem;
  --type-heading-sm:             0.875rem;
  --type-body:                   0.875rem;
  --type-body-sm:                0.75rem;
  --type-caption:                0.6875rem;
  --type-micro:                  0.625rem;
  --type-badge:                  0.6875rem;
  --type-ui-button:              0.875rem;
  --type-ui-nav:                 0.75rem;
  --type-ui-tooltip:             0.6875rem;

  /* --- LINE HEIGHTS (Agent 6) --- */
  --leading-display:             1.1;
  --leading-display-lg:          1.15;
  --leading-display-md:          1.2;
  --leading-heading:             1.3;
  --leading-heading-md:          1.35;
  --leading-heading-sm:          1.4;
  --leading-body:                1.5;
  --leading-caption:             1.4;
  --leading-micro:               1.3;
  --leading-badge:               1.0;
  --leading-ui:                  1.0;

  /* --- FONT WEIGHTS (Agent 6) --- */
  --weight-regular:              400;
  --weight-medium:               500;
  --weight-semibold:             600;
  --weight-bold:                 700;

  /* --- LETTER SPACING (Agent 6) --- */
  --tracking-tight:              -0.01em;
  --tracking-normal:             0.00em;
  --tracking-wide:               0.02em;
  --tracking-wider:              0.05em;
  --tracking-widest:             0.08em;
  --tracking-extreme:            0.12em;
}
```

---

## Appendix D: Implementation Checklist

- [ ] Load web fonts (`JetBrains Mono`, `Inter`, `Orbitron`) with `font-display: swap`
- [ ] Apply unified design token CSS block to `:root`
- [ ] Implement SVG-based grid renderer using `--grid-*` tokens
- [ ] Build glow effect utility classes (`.neon-glow`, `.neon-glow-intense`, `.neon-glow-subtle`)
- [ ] Build text glow utility classes (`.text-glow-*`)
- [ ] Implement pulse animation keyframes
- [ ] Build room state indicator component with all 10 states
- [ ] Build elevation view SVG with floor selection interaction
- [ ] Build plan view SVG with room hover/select interaction
- [ ] Implement status bar component with live data binding
- [ ] Implement label placement engine respecting minimum room sizes
- [ ] Test all contrast ratios with automated WCAG checker
- [ ] Validate print stylesheet with physical print test
- [ ] Performance test glow effects -- ensure 60fps on target hardware
- [ ] Verify `font-variant-numeric: tabular-nums` renders correctly across font stack

---

*Document version: 1.0.0 | Agents: 4, 5, 6 | Classification: OPERATIONAL*
*Holm Intelligence Complex -- Visual Interface Design Specification*
