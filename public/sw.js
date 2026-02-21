// HOLM Vault — Service Worker
// Cache-first for static assets, network-first for API calls

var CACHE_VERSION = 'holm-chat-v4';

var PRECACHE_URLS = [
  '/',
  '/index.html',
  '/manifest.json',
  '/icons/icon-192.svg',
  '/icons/icon-512.svg',
  '/fonts/JetBrainsMono-Regular.woff2',
  '/fonts/JetBrainsMono-Medium.woff2',
  '/fonts/JetBrainsMono-SemiBold.woff2',
  '/fonts/JetBrainsMono-Bold.woff2',
  '/fonts/JetBrainsMono-Light.woff2'
];

var API_ROUTES = ['/api/docs', '/api/graph', '/api/tags', '/api/comments', '/api/health', '/api/search'];

// ─── Install: precache static assets ────────────────────────────────
self.addEventListener('install', function(event) {
  event.waitUntil(
    caches.open(CACHE_VERSION).then(function(cache) {
      return cache.addAll(PRECACHE_URLS);
    }).then(function() {
      return self.skipWaiting();
    })
  );
});

// ─── Activate: clean up old cache versions ──────────────────────────
self.addEventListener('activate', function(event) {
  event.waitUntil(
    caches.keys().then(function(cacheNames) {
      return Promise.all(
        cacheNames
          .filter(function(name) { return name !== CACHE_VERSION; })
          .map(function(name) { return caches.delete(name); })
      );
    }).then(function() {
      return self.clients.claim();
    })
  );
});

// ─── Fetch handler ──────────────────────────────────────────────────
self.addEventListener('fetch', function(event) {
  var url = new URL(event.request.url);

  // Only handle same-origin requests
  if (url.origin !== location.origin) return;

  // Check if this is an API route
  var isAPI = API_ROUTES.some(function(route) {
    return url.pathname === route || url.pathname.startsWith(route + '/');
  });

  if (isAPI && event.request.method === 'GET') {
    // Network-first for API calls
    event.respondWith(networkFirstStrategy(event.request));
  } else if (event.request.mode === 'navigate') {
    // Navigation requests: try cache first, fallback to network
    event.respondWith(cacheFirstStrategy(event.request));
  } else {
    // Static assets: cache-first
    event.respondWith(cacheFirstStrategy(event.request));
  }
});

// ─── Cache-first strategy ───────────────────────────────────────────
function cacheFirstStrategy(request) {
  return caches.match(request).then(function(cachedResponse) {
    if (cachedResponse) {
      return cachedResponse;
    }
    return fetch(request).then(function(networkResponse) {
      if (networkResponse && networkResponse.status === 200) {
        var responseClone = networkResponse.clone();
        caches.open(CACHE_VERSION).then(function(cache) {
          cache.put(request, responseClone);
        });
      }
      return networkResponse;
    }).catch(function() {
      // Offline fallback for navigation requests
      if (request.mode === 'navigate') {
        return caches.match('/index.html');
      }
      return new Response('Offline', { status: 503, statusText: 'Service Unavailable' });
    });
  });
}

// ─── Network-first strategy (for API routes) ────────────────────────
function networkFirstStrategy(request) {
  return fetch(request).then(function(networkResponse) {
    if (networkResponse && networkResponse.status === 200) {
      var responseClone = networkResponse.clone();
      caches.open(CACHE_VERSION).then(function(cache) {
        cache.put(request, responseClone);
      });
    }
    return networkResponse;
  }).catch(function() {
    return caches.match(request).then(function(cachedResponse) {
      if (cachedResponse) {
        return cachedResponse;
      }
      // No cache available — return a JSON error
      return new Response(
        JSON.stringify({ error: 'offline', message: 'You are offline and no cached data is available.' }),
        {
          status: 503,
          statusText: 'Service Unavailable',
          headers: { 'Content-Type': 'application/json' }
        }
      );
    });
  });
}
