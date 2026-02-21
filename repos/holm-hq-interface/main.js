/* ============================================================
   HOLM HQ — Main Application Controller
   Pan/zoom, navigation, keyboard, touch, info panel, state
   ============================================================ */

(function () {
  'use strict';

  /* ---- State ---- */
  const state = {
    view: 'building',        // 'building' | 'floor' | 'room'
    currentFloor: null,
    currentRoom: null,
    zoom: 1,
    minZoom: 0.3,
    maxZoom: 3,
    panX: 0,
    panY: 0,
    isPanning: false,
    panStartX: 0,
    panStartY: 0,
    panStartPanX: 0,
    panStartPanY: 0,
    mapData: null,
    buildingMeta: null,
    selectedFloorIndex: -1,
    touchDist: 0
  };

  /* ---- DOM refs ---- */
  const $ = (s) => document.querySelector(s);
  const viewport = $('#viewport');
  const world = $('#canvas-world');
  const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
  svg.style.display = 'block';
  world.appendChild(svg);

  const minimapSvg = $('#minimap svg') || (() => {
    const s = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    $('#minimap').appendChild(s);
    return s;
  })();

  const infoPanel = $('#info-panel');
  const infoPanelTitle = $('#info-panel-title');
  const infoPanelBody = $('#info-panel-body');
  const infoPanelClose = $('#info-panel-close');
  const backBtn = $('#back-button');
  const helpOverlay = $('#help-overlay');
  const zoomDisplay = $('#zoom-display');
  const breadcrumb = $('#hud-breadcrumb');
  const floorIndicator = $('#floor-indicator');
  const tooltip = $('#tooltip');
  const loadingScreen = $('#loading-screen');
  const loadingBar = $('#loading-bar');
  const loadingStatus = $('#loading-status');

  /* ---- Transform helper ---- */
  function applyTransform() {
    world.style.transform = `translate(${state.panX}px, ${state.panY}px) scale(${state.zoom})`;
    if (zoomDisplay) zoomDisplay.textContent = `${Math.round(state.zoom * 100)}%`;
    updateMinimap();
  }

  function centerView(contentW, contentH) {
    const vw = window.innerWidth;
    const vh = window.innerHeight;
    const scaleX = (vw * 0.7) / contentW;
    const scaleY = (vh * 0.7) / contentH;
    state.zoom = Math.min(scaleX, scaleY, state.maxZoom);
    state.zoom = Math.max(state.zoom, state.minZoom);
    state.panX = (vw - contentW * state.zoom) / 2;
    state.panY = (vh - contentH * state.zoom) / 2;
    applyTransform();
  }

  /* ---- Zoom ---- */
  function zoomAt(delta, cx, cy) {
    const oldZoom = state.zoom;
    state.zoom *= delta > 0 ? 1.08 : 0.92;
    state.zoom = Math.max(state.minZoom, Math.min(state.maxZoom, state.zoom));

    const scale = state.zoom / oldZoom;
    state.panX = cx - (cx - state.panX) * scale;
    state.panY = cy - (cy - state.panY) * scale;
    applyTransform();
  }

  /* ---- Minimap ---- */
  function updateMinimap() {
    if (state.view !== 'building' || !state.buildingMeta) return;
    const bm = state.buildingMeta;
    const vw = window.innerWidth;
    const vh = window.innerHeight;

    const viewX = (-state.panX / state.zoom);
    const viewY = (-state.panY / state.zoom);
    const viewW = vw / state.zoom;
    const viewH = vh / state.zoom;

    Renderer.drawMinimap(minimapSvg, state.mapData, {
      x: viewX, y: viewY, w: viewW, h: viewH
    });
  }

  /* ---- Breadcrumb ---- */
  function updateBreadcrumb() {
    let html = '<span data-nav="building">HQ</span>';
    if (state.currentFloor) {
      html += '<span class="sep">/</span>';
      html += `<span data-nav="floor">${state.currentFloor.label}</span>`;
    }
    if (state.currentRoom) {
      html += '<span class="sep">/</span>';
      html += `<span data-nav="room">${state.currentRoom.label}</span>`;
    }
    breadcrumb.innerHTML = html;

    breadcrumb.querySelectorAll('span[data-nav]').forEach(s => {
      s.addEventListener('click', () => {
        const nav = s.getAttribute('data-nav');
        if (nav === 'building') navigateToBuilding();
        else if (nav === 'floor' && state.currentFloor) navigateToFloor(state.currentFloor);
      });
    });
  }

  /* ---- Floor indicator pips ---- */
  function buildFloorIndicator() {
    if (!state.mapData) return;
    floorIndicator.innerHTML = '';
    state.mapData.floors.slice().reverse().forEach((f, i) => {
      const pip = document.createElement('div');
      pip.className = 'floor-indicator-pip';
      pip.title = f.label;
      pip.setAttribute('data-floor-id', f.id);
      pip.addEventListener('click', () => navigateToFloor(f));
      floorIndicator.appendChild(pip);
    });
  }

  function updateFloorPips(activeId) {
    floorIndicator.querySelectorAll('.floor-indicator-pip').forEach(p => {
      p.classList.toggle('active', p.getAttribute('data-floor-id') === activeId);
    });
  }

  /* ---- Info Panel ---- */
  function showInfoPanel(title, fields) {
    infoPanelTitle.textContent = title;
    infoPanelBody.innerHTML = '';
    fields.forEach(f => {
      const div = document.createElement('div');
      div.className = 'field';
      div.innerHTML = `<div class="field-label">${f.label}</div><div class="field-value">${f.value}</div>`;
      infoPanelBody.appendChild(div);
    });

    // Random status bar for visual flair
    const statusDiv = document.createElement('div');
    statusDiv.className = 'field';
    statusDiv.innerHTML = `
      <div class="field-label">System Status</div>
      <div class="status-bar"><div class="status-bar-fill" style="width: ${60 + Math.random() * 35}%"></div></div>
    `;
    infoPanelBody.appendChild(statusDiv);

    infoPanel.classList.add('visible');
  }

  function hideInfoPanel() {
    infoPanel.classList.remove('visible');
  }

  /* ---- Tooltip ---- */
  function showTooltip(text, x, y) {
    tooltip.textContent = text;
    tooltip.style.left = (x + 14) + 'px';
    tooltip.style.top = (y - 8) + 'px';
    tooltip.classList.add('visible');
  }

  function hideTooltip() {
    tooltip.classList.remove('visible');
  }

  /* ---- Navigation ---- */
  function navigateToBuilding() {
    state.view = 'building';
    state.currentFloor = null;
    state.currentRoom = null;
    state.selectedFloorIndex = -1;
    hideInfoPanel();
    backBtn.classList.remove('visible');
    floorIndicator.classList.remove('visible');

    svg.classList.add('view-transition');
    setTimeout(() => svg.classList.remove('view-transition'), 400);

    state.buildingMeta = Renderer.drawBuilding(svg, state.mapData);
    centerView(state.buildingMeta.bw + 80, state.buildingMeta.totalH);
    updateBreadcrumb();
    bindBuildingEvents();
    $('#minimap').style.display = 'block';
  }

  function navigateToFloor(floor) {
    state.view = 'floor';
    state.currentFloor = floor;
    state.currentRoom = null;
    hideInfoPanel();
    backBtn.classList.add('visible');
    floorIndicator.classList.add('visible');
    updateFloorPips(floor.id);

    svg.classList.add('view-transition');
    setTimeout(() => svg.classList.remove('view-transition'), 400);

    const meta = Renderer.drawFloor(svg, floor, state.mapData);
    centerView(meta.totalW, meta.totalH);
    updateBreadcrumb();
    bindFloorEvents(floor);
    $('#minimap').style.display = 'none';

    // Update selected floor index
    state.selectedFloorIndex = state.mapData.floors.indexOf(floor);
  }

  function selectRoom(floor, room) {
    state.view = 'room';
    state.currentRoom = room;
    updateBreadcrumb();

    // Highlight room
    svg.querySelectorAll('.room-rect').forEach(r => r.classList.remove('selected'));
    const roomEl = svg.querySelector(`[data-room-id="${room.id}"]`);
    if (roomEl) roomEl.classList.add('selected');

    showInfoPanel(room.label, [
      { label: 'Room ID', value: room.id },
      { label: 'Floor', value: floor.label },
      { label: 'Position', value: `X:${room.x} Y:${room.y}` },
      { label: 'Dimensions', value: `${room.w} x ${room.h}` },
      { label: 'Type', value: room.icon.toUpperCase() },
      { label: 'Access Level', value: `LEVEL-${floor.level}` },
      { label: 'Contents', value: '<em style="color:#667788;">[ PLACEHOLDER ]</em>' }
    ]);
  }

  /* ---- Event bindings ---- */
  function bindBuildingEvents() {
    svg.querySelectorAll('.floor-rect').forEach(rect => {
      rect.addEventListener('click', (e) => {
        e.stopPropagation();
        const fid = rect.getAttribute('data-floor-id');
        const floor = state.mapData.floors.find(f => f.id === fid);
        if (floor) navigateToFloor(floor);
      });
      rect.addEventListener('mouseenter', (e) => {
        const fid = rect.getAttribute('data-floor-id');
        const floor = state.mapData.floors.find(f => f.id === fid);
        if (floor) showTooltip(floor.label, e.clientX, e.clientY);
      });
      rect.addEventListener('mousemove', (e) => {
        tooltip.style.left = (e.clientX + 14) + 'px';
        tooltip.style.top = (e.clientY - 8) + 'px';
      });
      rect.addEventListener('mouseleave', hideTooltip);
    });
  }

  function bindFloorEvents(floor) {
    svg.querySelectorAll('.room-rect').forEach(rect => {
      rect.addEventListener('click', (e) => {
        e.stopPropagation();
        const rid = rect.getAttribute('data-room-id');
        const room = floor.rooms.find(r => r.id === rid);
        if (room) selectRoom(floor, room);
      });
      rect.addEventListener('mouseenter', (e) => {
        const rid = rect.getAttribute('data-room-id');
        const room = floor.rooms.find(r => r.id === rid);
        if (room) showTooltip(room.label, e.clientX, e.clientY);
      });
      rect.addEventListener('mousemove', (e) => {
        tooltip.style.left = (e.clientX + 14) + 'px';
        tooltip.style.top = (e.clientY - 8) + 'px';
      });
      rect.addEventListener('mouseleave', hideTooltip);
    });
  }

  /* ---- Pan/Zoom input ---- */

  // Mouse
  viewport.addEventListener('mousedown', (e) => {
    if (e.button !== 0) return;
    if (e.target.closest('.floor-rect, .room-rect')) return;
    state.isPanning = true;
    state.panStartX = e.clientX;
    state.panStartY = e.clientY;
    state.panStartPanX = state.panX;
    state.panStartPanY = state.panY;
    viewport.classList.add('grabbing');
  });

  window.addEventListener('mousemove', (e) => {
    if (!state.isPanning) return;
    state.panX = state.panStartPanX + (e.clientX - state.panStartX);
    state.panY = state.panStartPanY + (e.clientY - state.panStartY);
    applyTransform();
  });

  window.addEventListener('mouseup', () => {
    state.isPanning = false;
    viewport.classList.remove('grabbing');
  });

  // Wheel zoom
  viewport.addEventListener('wheel', (e) => {
    e.preventDefault();
    zoomAt(-e.deltaY, e.clientX, e.clientY);
  }, { passive: false });

  // Touch
  viewport.addEventListener('touchstart', (e) => {
    if (e.touches.length === 1) {
      const t = e.touches[0];
      if (e.target.closest('.floor-rect, .room-rect')) return;
      state.isPanning = true;
      state.panStartX = t.clientX;
      state.panStartY = t.clientY;
      state.panStartPanX = state.panX;
      state.panStartPanY = state.panY;
    } else if (e.touches.length === 2) {
      state.isPanning = false;
      const dx = e.touches[0].clientX - e.touches[1].clientX;
      const dy = e.touches[0].clientY - e.touches[1].clientY;
      state.touchDist = Math.hypot(dx, dy);
    }
  }, { passive: true });

  viewport.addEventListener('touchmove', (e) => {
    if (e.touches.length === 1 && state.isPanning) {
      const t = e.touches[0];
      state.panX = state.panStartPanX + (t.clientX - state.panStartX);
      state.panY = state.panStartPanY + (t.clientY - state.panStartY);
      applyTransform();
    } else if (e.touches.length === 2) {
      e.preventDefault();
      const dx = e.touches[0].clientX - e.touches[1].clientX;
      const dy = e.touches[0].clientY - e.touches[1].clientY;
      const newDist = Math.hypot(dx, dy);
      const cx = (e.touches[0].clientX + e.touches[1].clientX) / 2;
      const cy = (e.touches[0].clientY + e.touches[1].clientY) / 2;
      const delta = newDist - state.touchDist;
      zoomAt(delta, cx, cy);
      state.touchDist = newDist;
    }
  }, { passive: false });

  viewport.addEventListener('touchend', () => {
    state.isPanning = false;
  });

  /* ---- Keyboard ---- */
  document.addEventListener('keydown', (e) => {
    const key = e.key.toLowerCase();

    // Help
    if (key === '?' || (key === 'h' && !e.ctrlKey && !e.metaKey)) {
      helpOverlay.classList.toggle('visible');
      return;
    }

    // Close help/panel on Escape
    if (key === 'escape') {
      if (helpOverlay.classList.contains('visible')) {
        helpOverlay.classList.remove('visible');
      } else if (infoPanel.classList.contains('visible')) {
        hideInfoPanel();
      } else if (state.view === 'floor' || state.view === 'room') {
        navigateToBuilding();
      }
      return;
    }

    // Back
    if (key === 'backspace' || key === 'b') {
      if (state.view === 'room') {
        navigateToFloor(state.currentFloor);
      } else if (state.view === 'floor') {
        navigateToBuilding();
      }
      return;
    }

    // Zoom
    if (key === '=' || key === '+') {
      zoomAt(1, window.innerWidth / 2, window.innerHeight / 2);
      return;
    }
    if (key === '-') {
      zoomAt(-1, window.innerWidth / 2, window.innerHeight / 2);
      return;
    }
    if (key === '0') {
      if (state.view === 'building' && state.buildingMeta) {
        centerView(state.buildingMeta.bw + 80, state.buildingMeta.totalH);
      }
      return;
    }

    // Arrow pan
    const panStep = 40;
    if (key === 'arrowleft')  { state.panX += panStep; applyTransform(); return; }
    if (key === 'arrowright') { state.panX -= panStep; applyTransform(); return; }
    if (key === 'arrowup')    { state.panY += panStep; applyTransform(); return; }
    if (key === 'arrowdown')  { state.panY -= panStep; applyTransform(); return; }

    // Number keys to jump to floor
    if (state.view === 'building' || state.view === 'floor') {
      const num = parseInt(key);
      if (num >= 1 && num <= 9 && state.mapData.floors[num - 1]) {
        navigateToFloor(state.mapData.floors[num - 1]);
        return;
      }
    }

    // Floor navigation with J/K in floor view
    if (state.view === 'floor' || state.view === 'room') {
      if (key === 'j' && state.selectedFloorIndex < state.mapData.floors.length - 1) {
        navigateToFloor(state.mapData.floors[state.selectedFloorIndex + 1]);
        return;
      }
      if (key === 'k' && state.selectedFloorIndex > 0) {
        navigateToFloor(state.mapData.floors[state.selectedFloorIndex - 1]);
        return;
      }
    }
  });

  /* ---- Panel close ---- */
  infoPanelClose.addEventListener('click', hideInfoPanel);

  /* ---- Back button ---- */
  backBtn.addEventListener('click', () => {
    if (state.view === 'room') navigateToFloor(state.currentFloor);
    else if (state.view === 'floor') navigateToBuilding();
  });

  /* ---- Help overlay click-outside ---- */
  helpOverlay.addEventListener('click', (e) => {
    if (e.target === helpOverlay) helpOverlay.classList.remove('visible');
  });

  /* ---- Loading sequence ---- */
  async function boot() {
    const steps = [
      { pct: 15, msg: 'LOADING MAP DATA...' },
      { pct: 40, msg: 'INITIALIZING RENDERER...' },
      { pct: 65, msg: 'BUILDING SCHEMATIC...' },
      { pct: 85, msg: 'BINDING SYSTEMS...' },
      { pct: 100, msg: 'ONLINE' }
    ];

    for (const step of steps) {
      loadingBar.style.width = step.pct + '%';
      loadingStatus.textContent = step.msg;
      await new Promise(r => setTimeout(r, 300 + Math.random() * 200));
    }

    // Load map data
    try {
      const resp = await fetch('map-data.json');
      state.mapData = await resp.json();
    } catch (err) {
      // Fallback: try to use inline data
      console.error('Failed to fetch map-data.json, attempting inline fallback');
      if (window.__HOLM_MAP_DATA) {
        state.mapData = window.__HOLM_MAP_DATA;
      } else {
        loadingStatus.textContent = 'ERROR: MAP DATA NOT FOUND';
        return;
      }
    }

    // Build
    buildFloorIndicator();
    navigateToBuilding();

    // Dismiss loading
    await new Promise(r => setTimeout(r, 400));
    loadingScreen.classList.add('fade-out');
    setTimeout(() => loadingScreen.remove(), 700);
  }

  /* ---- Window resize handler ---- */
  window.addEventListener('resize', () => {
    if (state.view === 'building' && state.buildingMeta) {
      // Soft re-center
      applyTransform();
    }
  });

  /* ---- Init ---- */
  boot();
})();
