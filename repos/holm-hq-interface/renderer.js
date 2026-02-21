/* ============================================================
   HOLM HQ — Renderer Engine
   Handles all SVG drawing: building, floors, rooms
   ============================================================ */

const Renderer = (() => {
  const NS = 'http://www.w3.org/2000/svg';

  function el(tag, attrs = {}) {
    const e = document.createElementNS(NS, tag);
    for (const [k, v] of Object.entries(attrs)) {
      e.setAttribute(k, v);
    }
    return e;
  }

  /* ---- Draw the full skyscraper overview ---- */
  function drawBuilding(svg, mapData) {
    svg.innerHTML = '';

    const bw = mapData.building.width;
    const floorH = 60;
    const gap = 6;
    const totalFloors = mapData.floors.length;
    const totalH = totalFloors * (floorH + gap) + 140;

    svg.setAttribute('viewBox', `0 0 ${bw + 80} ${totalH}`);
    svg.setAttribute('width', bw + 80);
    svg.setAttribute('height', totalH);

    const defs = el('defs');

    // Glow filter
    const filter = el('filter', { id: 'neonGlow', x: '-20%', y: '-20%', width: '140%', height: '140%' });
    const blur = el('feGaussianBlur', { stdDeviation: '3', result: 'glow' });
    const merge = el('feMerge');
    const mn1 = el('feMergeNode', { in: 'glow' });
    const mn2 = el('feMergeNode', { in: 'SourceGraphic' });
    merge.appendChild(mn1);
    merge.appendChild(mn2);
    filter.appendChild(blur);
    filter.appendChild(merge);
    defs.appendChild(filter);
    svg.appendChild(defs);

    // Spire
    const spireGroup = el('g', { class: 'spire-group' });
    const cx = (bw + 80) / 2;
    const spireBase = 80;
    const spireTop = 20;

    spireGroup.appendChild(el('line', {
      class: 'spire-line',
      x1: cx, y1: spireBase, x2: cx, y2: spireTop,
      filter: 'url(#neonGlow)'
    }));
    spireGroup.appendChild(el('circle', {
      class: 'spire-dot',
      cx: cx, cy: spireTop, r: 3,
      filter: 'url(#neonGlow)'
    }));

    // Small antenna branches
    spireGroup.appendChild(el('line', {
      class: 'spire-line', x1: cx, y1: 50, x2: cx - 15, y2: 40
    }));
    spireGroup.appendChild(el('line', {
      class: 'spire-line', x1: cx, y1: 50, x2: cx + 15, y2: 40
    }));

    svg.appendChild(spireGroup);

    // Building outline
    const buildingY = 80;
    const buildingH = totalFloors * (floorH + gap) + 20;
    const outline = el('rect', {
      x: 38, y: buildingY, width: bw + 4, height: buildingH,
      fill: 'none',
      stroke: 'rgba(0, 255, 204, 0.15)',
      'stroke-width': 1
    });
    svg.appendChild(outline);

    // Floors (bottom-up visually, so reverse draw order)
    const floorGroup = el('g', { class: 'building-group' });

    for (let i = 0; i < totalFloors; i++) {
      const floor = mapData.floors[i];
      const fy = buildingY + 10 + (totalFloors - 1 - i) * (floorH + gap);

      const g = el('g', { 'data-floor-id': floor.id, 'data-floor-index': String(i) });

      const rect = el('rect', {
        class: 'floor-rect',
        x: 40, y: fy, width: bw, height: floorH,
        rx: 2, ry: 2,
        stroke: floor.color,
        'data-floor-id': floor.id
      });

      const label = el('text', {
        class: 'floor-label',
        x: 50, y: fy + floorH / 2 + 4,
        fill: '#ffffff',
        'font-weight': 'bold'
      });
      label.textContent = floor.label;

      // Floor number on the right
      const num = el('text', {
        class: 'floor-label',
        x: bw + 30, y: fy + floorH / 2 + 4,
        'text-anchor': 'end',
        fill: floor.color,
        'font-weight': 'bold'
      });
      num.textContent = `F${String(floor.level).padStart(2, '0')}`;

      // Decorative lines inside floor
      for (let li = 0; li < 3; li++) {
        const dx = 50 + li * (bw / 3);
        g.appendChild(el('line', {
          x1: dx, y1: fy + 1, x2: dx, y2: fy + floorH - 1,
          stroke: floor.color, 'stroke-width': 0.3, opacity: '0.3'
        }));
      }

      g.appendChild(rect);
      g.appendChild(label);
      g.appendChild(num);
      floorGroup.appendChild(g);
    }

    svg.appendChild(floorGroup);

    // Foundation
    const foundY = buildingY + buildingH;
    svg.appendChild(el('line', {
      x1: 20, y1: foundY, x2: bw + 60, y2: foundY,
      stroke: 'rgba(0, 255, 204, 0.3)', 'stroke-width': 2
    }));
    // Ground hash
    for (let hx = 20; hx < bw + 60; hx += 12) {
      svg.appendChild(el('line', {
        x1: hx, y1: foundY, x2: hx - 6, y2: foundY + 8,
        stroke: 'rgba(0, 255, 204, 0.15)', 'stroke-width': 1
      }));
    }

    return { floorH, gap, buildingY, totalH, bw };
  }

  /* ---- Draw a single floor blueprint ---- */
  function drawFloor(svg, floor, mapData) {
    svg.innerHTML = '';

    const padding = 40;
    const bpW = 400;
    const bpH = 260;
    const totalW = bpW + padding * 2;
    const totalH = bpH + padding * 2 + 60;

    svg.setAttribute('viewBox', `0 0 ${totalW} ${totalH}`);
    svg.setAttribute('width', totalW);
    svg.setAttribute('height', totalH);

    const defs = el('defs');
    const filter = el('filter', { id: 'roomGlow', x: '-10%', y: '-10%', width: '120%', height: '120%' });
    const blur = el('feGaussianBlur', { stdDeviation: '2', result: 'glow' });
    const merge = el('feMerge');
    merge.appendChild(el('feMergeNode', { in: 'glow' }));
    merge.appendChild(el('feMergeNode', { in: 'SourceGraphic' }));
    filter.appendChild(blur);
    filter.appendChild(merge);
    defs.appendChild(filter);
    svg.appendChild(defs);

    // Floor title
    const title = el('text', {
      x: String(totalW / 2),
      y: '30',
      fill: floor.color,
      'font-family': "'Courier New', monospace",
      'font-size': '18',
      'font-weight': 'bold',
      'text-anchor': 'middle',
      'letter-spacing': '3',
      filter: 'url(#roomGlow)'
    });
    title.textContent = floor.label.toUpperCase();
    svg.appendChild(title);

    // Blueprint border
    svg.appendChild(el('rect', {
      x: padding, y: padding + 20, width: bpW, height: bpH,
      fill: 'none',
      stroke: floor.color,
      'stroke-width': 1,
      'stroke-dasharray': '4,4',
      opacity: '0.4'
    }));

    // Corner markers
    const corners = [
      [padding, padding + 20],
      [padding + bpW, padding + 20],
      [padding, padding + 20 + bpH],
      [padding + bpW, padding + 20 + bpH]
    ];
    corners.forEach(([cx, cy]) => {
      svg.appendChild(el('line', {
        x1: cx - 8, y1: cy, x2: cx + 8, y2: cy,
        stroke: floor.color, 'stroke-width': 1, opacity: '0.6'
      }));
      svg.appendChild(el('line', {
        x1: cx, y1: cy - 8, x2: cx, y2: cy + 8,
        stroke: floor.color, 'stroke-width': 1, opacity: '0.6'
      }));
    });

    // Rooms
    const roomGroup = el('g', { class: 'room-group' });
    const iconMap = {
      terminal: '\u2588', shield: '\u2666', map: '\u2302', signal: '\u2637',
      network: '\u2641', archive: '\u2610', database: '\u2261', tools: '\u2692',
      blueprint: '\u2316', grid: '\u2630', book: '\u2261', search: '\u2315',
      flask: '\u2697', chart: '\u2584', vr: '\u2609', screen: '\u25a3',
      cross: '\u271a', pill: '\u2716', bed: '\u2302', users: '\u2603',
      food: '\u2616', radar: '\u25ce', lock: '\u2302', brain: '\u2609',
      eye: '\u25c9'
    };

    floor.rooms.forEach((room, idx) => {
      const rx = padding + room.x;
      const ry = padding + 20 + room.y;

      const g = el('g', { 'data-room-id': room.id, 'data-room-index': String(idx) });

      const rect = el('rect', {
        class: 'room-rect',
        x: rx, y: ry, width: room.w, height: room.h,
        rx: 1, ry: 1,
        stroke: floor.color,
        'data-room-id': room.id
      });

      const icon = el('text', {
        class: 'room-icon',
        x: rx + room.w / 2,
        y: ry + room.h / 2 - 8,
        fill: floor.color,
        opacity: '0.8'
      });
      icon.textContent = iconMap[room.icon] || '\u25a0';

      const label = el('text', {
        class: 'room-label',
        x: rx + room.w / 2,
        y: ry + room.h / 2 + 12,
        fill: '#ffffff',
        'font-weight': 'bold'
      });
      label.textContent = room.label;

      g.appendChild(rect);
      g.appendChild(icon);
      g.appendChild(label);
      roomGroup.appendChild(g);
    });

    svg.appendChild(roomGroup);

    // Dimension annotations
    const dimY = padding + 20 + bpH + 20;
    const dimText = el('text', {
      x: String(totalW / 2), y: String(dimY + 15),
      fill: floor.color, opacity: '1',
      'font-family': "'Courier New', monospace",
      'font-size': '12',
      'font-weight': 'bold',
      'text-anchor': 'middle',
      'letter-spacing': '2'
    });
    dimText.textContent = `${floor.rooms.length} ROOMS \u2014 LEVEL ${floor.level}`;
    svg.appendChild(dimText);

    return { bpW, bpH, totalW, totalH };
  }

  /* ---- Draw minimap ---- */
  function drawMinimap(svg, mapData, viewState) {
    svg.innerHTML = '';
    const bw = mapData.building.width;
    const floorH = 60;
    const gap = 6;
    const totalFloors = mapData.floors.length;
    const totalH = totalFloors * (floorH + gap) + 140;

    svg.setAttribute('viewBox', `0 0 ${bw + 80} ${totalH}`);

    // Simplified floors
    const buildingY = 80;
    for (let i = 0; i < totalFloors; i++) {
      const floor = mapData.floors[i];
      const fy = buildingY + 10 + (totalFloors - 1 - i) * (floorH + gap);
      svg.appendChild(el('rect', {
        x: 40, y: fy, width: bw, height: floorH,
        fill: 'none', stroke: floor.color, 'stroke-width': 0.8, opacity: '0.5'
      }));
    }

    // Viewport rectangle
    if (viewState) {
      const vpRect = el('rect', {
        id: 'minimap-viewport-rect',
        x: viewState.x, y: viewState.y,
        width: viewState.w, height: viewState.h
      });
      svg.appendChild(vpRect);
    }
  }

  return { drawBuilding, drawFloor, drawMinimap, el };
})();
