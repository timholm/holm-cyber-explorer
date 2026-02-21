<script>
  import { getAllVideos, searchVideos, getDownloadsProgress, formatDuration, formatRelativeTime, formatBytes, downloadVideo } from '../lib/api.js';
  import VideoCard from '../components/VideoCard.svelte';

  let { navigate, initialSearch = '' } = $props();

  let videos = $state([]);
  let downloadProgressMap = $state({});
  let loading = $state(true);
  let error = $state(null);
  let searchQuery = $state('');
  let searchTimeout = $state(null);
  let statusFilter = $state('');
  let useFullTextSearch = $state(true);
  let viewMode = $state('grid'); // 'grid' | 'list'
  let sortBy = $state('newest'); // 'newest' | 'oldest' | 'title' | 'duration'

  async function loadVideos(search = '') {
    loading = videos.length === 0;
    error = null;

    try {
      let response;
      if (search && useFullTextSearch) {
        // Use FTS5 full-text search for better results
        response = await searchVideos(search, { status: statusFilter, limit: 100 });
        videos = response.results || [];
      } else {
        // Fall back to basic search
        response = await getAllVideos({ search, status: statusFilter, limit: 100 });
        videos = response.videos || response || [];
      }
    } catch (err) {
      error = err.message;
      videos = [];
    } finally {
      loading = false;
    }
  }

  async function loadDownloadProgress() {
    try {
      const response = await getDownloadsProgress();
      const downloads = response.downloads || [];
      // Create a map of video_id -> progress for quick lookup
      const newMap = {};
      for (const download of downloads) {
        if (download.video_id) {
          newMap[download.video_id] = download;
        }
      }
      downloadProgressMap = newMap;
    } catch (err) {
      // Silently fail for progress updates
      console.warn('Failed to fetch download progress:', err);
    }
  }

  async function loadAll() {
    await Promise.all([loadVideos(searchQuery), loadDownloadProgress()]);
  }

  function handleSearchInput(e) {
    const value = e.target.value;
    searchQuery = value;

    // Debounce search
    if (searchTimeout) {
      clearTimeout(searchTimeout);
    }

    searchTimeout = setTimeout(() => {
      loadVideos(value);
    }, 300);
  }

  // Check if any videos are downloading
  const hasDownloadingVideos = $derived(
    videos.some(v => v.status === 'downloading') || Object.keys(downloadProgressMap).length > 0
  );

  $effect(() => {
    // Set initial search query if provided
    if (initialSearch) {
      searchQuery = initialSearch;
    }
    loadAll();
    // Poll for progress updates every 2 seconds if there are downloading videos
    const interval = setInterval(() => {
      if (hasDownloadingVideos) {
        loadDownloadProgress();
      }
    }, 2000);
    return () => clearInterval(interval);
  });

  $effect(() => {
    return () => {
      if (searchTimeout) {
        clearTimeout(searchTimeout);
      }
    };
  });

  // Get download progress for a specific video
  function getProgressForVideo(videoId) {
    return downloadProgressMap[videoId] || null;
  }

  // Sort videos based on selected sort option
  const sortedVideos = $derived(() => {
    const sorted = [...videos];
    switch (sortBy) {
      case 'oldest':
        return sorted.sort((a, b) => new Date(a.publishedAt || 0) - new Date(b.publishedAt || 0));
      case 'title':
        return sorted.sort((a, b) => (a.title || '').localeCompare(b.title || ''));
      case 'duration':
        return sorted.sort((a, b) => (b.duration || 0) - (a.duration || 0));
      case 'newest':
      default:
        return sorted.sort((a, b) => new Date(b.publishedAt || 0) - new Date(a.publishedAt || 0));
    }
  });

  function handleVideoClick(video) {
    if (navigate) {
      navigate('video', { id: video.id });
    }
  }

  function getStatusBadgeClass(status) {
    switch (status) {
      case 'completed': return 'badge-success';
      case 'downloading': return 'badge-warning';
      case 'failed': return 'badge-error';
      case 'pending': return 'badge-info';
      default: return 'badge-neutral';
    }
  }
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
    <div class="hidden lg:block">
      <h1 class="text-2xl font-bold text-dark-100">Videos</h1>
      <p class="text-dark-400 mt-1">Browse all archived videos</p>
    </div>
    {#if Object.keys(downloadProgressMap).length > 0}
      <div class="flex items-center gap-2 text-sm text-dark-400">
        <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
        <span>{Object.keys(downloadProgressMap).length} active download{Object.keys(downloadProgressMap).length === 1 ? '' : 's'}</span>
      </div>
    {/if}
  </div>

  <!-- Search and Filters -->
  <div class="flex flex-col gap-4">
    <div class="flex flex-col sm:flex-row gap-4">
      <div class="relative flex-1">
        <svg class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-dark-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          type="text"
          value={searchQuery}
          oninput={handleSearchInput}
          placeholder="Search videos (full-text search)..."
          class="input pl-12"
        />
      </div>

      <select
        bind:value={statusFilter}
        onchange={() => loadVideos(searchQuery)}
        class="input w-full sm:w-40"
      >
        <option value="">All Status</option>
        <option value="completed">Downloaded</option>
        <option value="pending">Pending</option>
        <option value="downloading">Downloading</option>
        <option value="failed">Failed</option>
      </select>

      <select
        bind:value={sortBy}
        class="input w-full sm:w-40"
      >
        <option value="newest">Newest First</option>
        <option value="oldest">Oldest First</option>
        <option value="title">Title A-Z</option>
        <option value="duration">Longest</option>
      </select>
    </div>

    <!-- View Mode Toggle -->
    <div class="flex items-center justify-between">
      <p class="text-dark-400 text-sm">
        {videos.length} video{videos.length === 1 ? '' : 's'} found
        {#if searchQuery}
          for "{searchQuery}"
        {/if}
      </p>
      <div class="flex items-center gap-2 bg-dark-800 rounded-lg p-1">
        <button
          onclick={() => viewMode = 'grid'}
          class="p-2 rounded-lg transition-colors {viewMode === 'grid' ? 'bg-dark-700 text-dark-100' : 'text-dark-500 hover:text-dark-300'}"
          title="Grid view"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
        </button>
        <button
          onclick={() => viewMode = 'list'}
          class="p-2 rounded-lg transition-colors {viewMode === 'list' ? 'bg-dark-700 text-dark-100' : 'text-dark-500 hover:text-dark-300'}"
          title="List view"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
      </div>
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-red-500"></div>
    </div>
  {:else if error}
    <div class="card p-6">
      <p class="text-red-400">Error loading videos: {error}</p>
      <button onclick={() => loadVideos(searchQuery)} class="btn btn-secondary mt-4">Retry</button>
    </div>
  {:else if videos.length === 0}
    <div class="card p-12 text-center">
      <div class="w-16 h-16 mx-auto bg-dark-800 rounded-full flex items-center justify-center mb-4">
        <svg class="w-8 h-8 text-dark-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <h3 class="text-lg font-medium text-dark-200 mb-2">
        {#if searchQuery}
          No videos found
        {:else}
          No videos yet
        {/if}
      </h3>
      <p class="text-dark-400 mb-6">
        {#if searchQuery}
          Try a different search term
        {:else}
          Add a channel to start archiving videos
        {/if}
      </p>
      {#if !searchQuery}
        <button onclick={() => navigate('channels')} class="btn btn-primary">
          Add Channel
        </button>
      {/if}
    </div>
  {:else}
    <!-- Videos Display -->
    {#if viewMode === 'grid'}
      <!-- Grid View -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {#each sortedVideos() as video (video.id)}
          <VideoCard
            {video}
            showChannel={true}
            {navigate}
            downloadProgress={getProgressForVideo(video.id)}
          />
        {/each}
      </div>
    {:else}
      <!-- List View -->
      <div class="space-y-2">
        {#each sortedVideos() as video (video.id)}
          <div
            class="card flex gap-4 p-3 hover:border-dark-700 cursor-pointer transition-all"
            onclick={() => handleVideoClick(video)}
            onkeydown={(e) => e.key === 'Enter' && handleVideoClick(video)}
            role="button"
            tabindex="0"
          >
            <!-- Thumbnail -->
            <div class="relative w-40 sm:w-48 flex-shrink-0 aspect-video bg-dark-800 rounded-lg overflow-hidden">
              {#if video.thumbnailUrl}
                <img
                  src={video.thumbnailUrl}
                  alt={video.title}
                  class="w-full h-full object-cover"
                  loading="lazy"
                />
              {:else}
                <div class="w-full h-full flex items-center justify-center">
                  <svg class="w-8 h-8 text-dark-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
              {/if}
              {#if video.duration}
                <div class="absolute bottom-1 right-1 px-1.5 py-0.5 bg-black/80 rounded text-xs text-white font-medium">
                  {formatDuration(video.duration)}
                </div>
              {/if}
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0 flex flex-col justify-between py-1">
              <div>
                <h3 class="font-medium text-dark-100 line-clamp-2 hover:text-red-400 transition-colors">
                  {video.title || video.videoId}
                </h3>
                {#if video.channelName}
                  <p class="text-sm text-dark-400 mt-1">{video.channelName}</p>
                {/if}
                <div class="flex items-center gap-2 mt-2 text-xs text-dark-500">
                  {#if video.viewCount}
                    <span>{video.viewCount.toLocaleString()} views</span>
                    <span>&middot;</span>
                  {/if}
                  {#if video.publishedAt}
                    <span>{formatRelativeTime(video.publishedAt)}</span>
                  {/if}
                </div>
              </div>
              <div class="flex items-center gap-3 mt-2">
                <span class="badge {getStatusBadgeClass(video.status)}">{video.status}</span>
                {#if video.fileSize}
                  <span class="text-xs text-dark-500">{formatBytes(video.fileSize)}</span>
                {/if}
              </div>
            </div>

            <!-- Actions -->
            <div class="flex items-center gap-2 flex-shrink-0">
              {#if video.status === 'completed'}
                <a
                  href="/api/videos/{video.id}/download"
                  download
                  onclick={(e) => e.stopPropagation()}
                  class="btn btn-secondary text-sm py-1.5 px-3"
                  title="Download"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                  </svg>
                </a>
              {/if}
            </div>
          </div>
        {/each}
      </div>
    {/if}
  {/if}
</div>

<style>
  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
</style>
