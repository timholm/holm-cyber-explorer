<script>
  import Dashboard from './routes/Dashboard.svelte';
  import Channels from './routes/Channels.svelte';
  import Channel from './routes/Channel.svelte';
  import Videos from './routes/Videos.svelte';
  import VideoPlayer from './routes/VideoPlayer.svelte';
  import Jobs from './routes/Jobs.svelte';
  import Settings from './routes/Settings.svelte';

  let currentRoute = $state('dashboard');
  let currentParams = $state({});
  let sidebarOpen = $state(false);
  let globalSearchQuery = $state('');

  const routes = [
    { id: 'dashboard', name: 'Dashboard', icon: 'home' },
    { id: 'channels', name: 'Channels', icon: 'users' },
    { id: 'videos', name: 'Videos', icon: 'play' },
    { id: 'jobs', name: 'Jobs', icon: 'activity' },
    { id: 'settings', name: 'Settings', icon: 'settings' }
  ];

  function navigate(route, params = {}) {
    currentRoute = route;
    currentParams = params;
    sidebarOpen = false;
  }

  function handleGlobalSearch(e) {
    if (e.key === 'Enter' && globalSearchQuery.trim()) {
      currentRoute = 'videos';
      currentParams = { search: globalSearchQuery.trim() };
      globalSearchQuery = '';
    }
  }

  function getIcon(icon) {
    const icons = {
      home: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6',
      users: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z',
      play: 'M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
      activity: 'M13 10V3L4 14h7v7l9-11h-7z',
      settings: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z'
    };
    return icons[icon] || icons.home;
  }
</script>

<div class="flex h-screen overflow-hidden">
  <!-- Mobile sidebar overlay -->
  {#if sidebarOpen}
    <div
      class="fixed inset-0 bg-black/50 z-40 lg:hidden"
      onclick={() => sidebarOpen = false}
      onkeydown={(e) => e.key === 'Escape' && (sidebarOpen = false)}
      role="button"
      tabindex="0"
      aria-label="Close sidebar"
    ></div>
  {/if}

  <!-- Sidebar -->
  <aside class="fixed lg:static inset-y-0 left-0 z-50 w-64 bg-dark-900 border-r border-dark-800 transform transition-transform duration-200 ease-in-out {sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}">
    <div class="flex flex-col h-full">
      <!-- Logo -->
      <div class="flex items-center gap-3 px-6 py-5 border-b border-dark-800">
        <div class="w-10 h-10 bg-red-600 rounded-lg flex items-center justify-center">
          <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="currentColor">
            <path d="M19.615 3.184c-3.604-.246-11.631-.245-15.23 0C.488 3.45.029 5.804 0 12c.029 6.185.484 8.549 4.385 8.816 3.6.245 11.626.246 15.23 0C23.512 20.55 23.971 18.196 24 12c-.029-6.185-.484-8.549-4.385-8.816zM9 16V8l8 4-8 4z"/>
          </svg>
        </div>
        <div>
          <h1 class="text-lg font-bold text-dark-100">YT Archive</h1>
          <p class="text-xs text-dark-500">Channel Archiver</p>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 px-4 py-6 space-y-1 overflow-y-auto">
        {#each routes as route}
          <button
            onclick={() => navigate(route.id)}
            class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-left transition-colors duration-200 {currentRoute === route.id ? 'bg-red-600/10 text-red-500' : 'text-dark-400 hover:bg-dark-800 hover:text-dark-100'}"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d={getIcon(route.icon)} />
            </svg>
            <span class="font-medium">{route.name}</span>
          </button>
        {/each}
      </nav>

      <!-- Footer -->
      <div class="px-6 py-4 border-t border-dark-800">
        <p class="text-xs text-dark-500">YouTube Archiver v1.0</p>
      </div>
    </div>
  </aside>

  <!-- Main content -->
  <main class="flex-1 overflow-y-auto">
    <!-- Desktop header with search -->
    <header class="sticky top-0 z-30 hidden lg:flex bg-dark-900/95 backdrop-blur border-b border-dark-800 px-6 py-3 items-center justify-between">
      <div class="relative w-96">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-dark-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          type="text"
          bind:value={globalSearchQuery}
          onkeydown={handleGlobalSearch}
          placeholder="Search videos... (press Enter)"
          class="w-full pl-10 pr-4 py-2 bg-dark-800 border border-dark-700 rounded-lg text-dark-100 placeholder-dark-500 focus:outline-none focus:border-red-500 transition-colors"
        />
      </div>
      <div class="flex items-center gap-4">
        <span class="text-sm text-dark-400">Press ⌘K to search</span>
      </div>
    </header>

    <!-- Mobile header -->
    <header class="sticky top-0 z-30 lg:hidden bg-dark-900/95 backdrop-blur border-b border-dark-800 px-4 py-3">
      <div class="flex items-center gap-4">
        <button
          onclick={() => sidebarOpen = true}
          class="p-2 rounded-lg hover:bg-dark-800 text-dark-400"
          aria-label="Open menu"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
        <h1 class="text-lg font-semibold text-dark-100 flex-1">
          {routes.find(r => r.id === currentRoute)?.name || 'Dashboard'}
        </h1>
        <button
          onclick={() => navigate('videos')}
          class="p-2 rounded-lg hover:bg-dark-800 text-dark-400"
          aria-label="Search"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </button>
      </div>
    </header>

    <!-- Page content -->
    <div class="p-6 lg:p-8">
      {#if currentRoute === 'dashboard'}
        <Dashboard {navigate} />
      {:else if currentRoute === 'channels'}
        <Channels {navigate} />
      {:else if currentRoute === 'channel'}
        <Channel channelId={currentParams.id} {navigate} />
      {:else if currentRoute === 'videos'}
        <Videos {navigate} initialSearch={currentParams.search || ''} />
      {:else if currentRoute === 'video'}
        <VideoPlayer videoId={currentParams.id} {navigate} />
      {:else if currentRoute === 'jobs'}
        <Jobs />
      {:else if currentRoute === 'settings'}
        <Settings />
      {/if}
    </div>
  </main>
</div>
