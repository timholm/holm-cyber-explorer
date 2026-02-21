<script>
  import { getChannel, getChannelVideos, triggerSync } from '../lib/api.js';
  import VideoCard from '../components/VideoCard.svelte';

  let { channelId, navigate } = $props();

  let channel = $state(null);
  let videos = $state([]);
  let loading = $state(true);
  let error = $state(null);
  let syncing = $state(false);
  let filter = $state('all');

  const filters = [
    { id: 'all', label: 'All' },
    { id: 'completed', label: 'Completed' },
    { id: 'pending', label: 'Pending' },
    { id: 'downloading', label: 'Downloading' },
    { id: 'failed', label: 'Failed' }
  ];

  async function loadData() {
    loading = true;
    error = null;

    try {
      const [channelData, videosData] = await Promise.all([
        getChannel(channelId),
        getChannelVideos(channelId)
      ]);

      channel = channelData.channel || channelData;
      videos = videosData.videos || videosData || [];
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  async function handleSync() {
    syncing = true;
    try {
      await triggerSync(channelId);
      // Reload data after sync
      await loadData();
    } catch (err) {
      error = err.message;
    } finally {
      syncing = false;
    }
  }

  $effect(() => {
    if (channelId) {
      loadData();
    }
  });

  const filteredVideos = $derived(
    filter === 'all'
      ? videos
      : videos.filter(v => v.status === filter)
  );

  const statusCounts = $derived({
    all: videos.length,
    completed: videos.filter(v => v.status === 'completed').length,
    pending: videos.filter(v => v.status === 'pending').length,
    downloading: videos.filter(v => v.status === 'downloading').length,
    failed: videos.filter(v => v.status === 'failed').length
  });
</script>

<div class="space-y-6">
  <!-- Back button -->
  <button onclick={() => navigate('channels')} class="btn btn-ghost -ml-4">
    <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
    </svg>
    Back to Channels
  </button>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-red-500"></div>
    </div>
  {:else if error}
    <div class="card p-6">
      <p class="text-red-400">Error loading channel: {error}</p>
      <button onclick={loadData} class="btn btn-secondary mt-4">Retry</button>
    </div>
  {:else if channel}
    <!-- Channel Header -->
    <div class="card overflow-hidden">
      <!-- Banner -->
      <div class="h-32 sm:h-48 bg-gradient-to-r from-dark-800 to-dark-700 relative">
        {#if channel.bannerUrl}
          <img
            src={channel.bannerUrl}
            alt=""
            class="w-full h-full object-cover"
          />
        {/if}
      </div>

      <!-- Channel Info -->
      <div class="relative px-6 pb-6">
        <!-- Avatar -->
        <div class="absolute -top-12 left-6">
          <div class="w-24 h-24 rounded-full border-4 border-dark-900 overflow-hidden bg-dark-700">
            {#if channel.avatarUrl}
              <img
                src={channel.avatarUrl}
                alt={channel.name}
                class="w-full h-full object-cover"
              />
            {:else}
              <div class="w-full h-full flex items-center justify-center text-dark-400 text-3xl font-bold">
                {channel.name?.[0]?.toUpperCase() || '?'}
              </div>
            {/if}
          </div>
        </div>

        <!-- Info and Actions -->
        <div class="pt-16 flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
          <div>
            <h1 class="text-2xl font-bold text-dark-100">{channel.name}</h1>
            {#if channel.description}
              <p class="text-dark-400 mt-2 max-w-2xl">{channel.description}</p>
            {/if}
            <div class="flex items-center gap-4 mt-3 text-sm text-dark-500">
              <span>{channel.subscriberCount?.toLocaleString() || 0} subscribers</span>
              <span>{videos.length} videos</span>
            </div>
          </div>

          <button
            onclick={handleSync}
            class="btn btn-primary shrink-0"
            disabled={syncing}
          >
            {#if syncing}
              <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2 inline-block"></div>
              Syncing...
            {:else}
              <svg class="w-4 h-4 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Sync Channel
            {/if}
          </button>
        </div>
      </div>
    </div>

    <!-- Filter Tabs -->
    <div class="flex gap-2 overflow-x-auto pb-2">
      {#each filters as f}
        <button
          onclick={() => filter = f.id}
          class="px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors {filter === f.id ? 'bg-red-600 text-white' : 'bg-dark-800 text-dark-400 hover:bg-dark-700 hover:text-dark-200'}"
        >
          {f.label}
          <span class="ml-1.5 px-1.5 py-0.5 rounded-full text-xs {filter === f.id ? 'bg-red-700' : 'bg-dark-700'}">
            {statusCounts[f.id]}
          </span>
        </button>
      {/each}
    </div>

    <!-- Videos Grid -->
    {#if filteredVideos.length === 0}
      <div class="card p-12 text-center">
        <p class="text-dark-400">
          {#if filter === 'all'}
            No videos found. Try syncing the channel.
          {:else}
            No {filter} videos found.
          {/if}
        </p>
      </div>
    {:else}
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {#each filteredVideos as video}
          <VideoCard {video} />
        {/each}
      </div>
    {/if}
  {/if}
</div>
