<script>
  import { downloadVideo, formatDuration, formatRelativeTime, formatBytes } from '../lib/api.js';
  import DownloadProgressBar from './DownloadProgressBar.svelte';

  let { video, showChannel = false, navigate = null, downloadProgress = null } = $props();

  let downloading = $state(false);
  let error = $state(null);

  // Check if this video has active download progress
  const hasProgress = $derived(downloadProgress && downloadProgress.percentage > 0);

  async function handleDownload(e) {
    e.stopPropagation();
    downloading = true;
    error = null;

    try {
      await downloadVideo(video.id);
    } catch (err) {
      error = err.message;
    } finally {
      downloading = false;
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

  function getStatusText(status) {
    switch (status) {
      case 'completed': return 'Downloaded';
      case 'downloading': return 'Downloading';
      case 'failed': return 'Failed';
      case 'pending': return 'Pending';
      default: return status;
    }
  }

  function handleCardClick() {
    if (navigate) {
      // Navigate to video player
      navigate('video', { id: video.id });
    } else if (video.status === 'completed') {
      // Fallback: open stream in new tab
      window.open(`/api/videos/${video.id}/stream`, '_blank');
    }
  }

  function handleChannelClick(e) {
    e.stopPropagation();
    if (navigate && video.channelId) {
      navigate('channel', { id: video.channelId });
    }
  }
</script>

<div
  class="card group cursor-pointer hover:border-dark-700 transition-all duration-200"
  onclick={handleCardClick}
  onkeydown={(e) => e.key === 'Enter' && handleCardClick()}
  role="button"
  tabindex="0"
>
  <!-- Thumbnail -->
  <div class="relative aspect-video bg-dark-800 overflow-hidden">
    {#if video.thumbnailUrl}
      <img
        src={video.thumbnailUrl}
        alt={video.title}
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
        loading="lazy"
      />
    {:else}
      <div class="w-full h-full flex items-center justify-center">
        <svg class="w-12 h-12 text-dark-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
    {/if}

    <!-- Duration badge -->
    {#if video.duration}
      <div class="absolute bottom-2 right-2 px-1.5 py-0.5 bg-black/80 rounded text-xs text-white font-medium">
        {formatDuration(video.duration)}
      </div>
    {/if}

    <!-- Status badge -->
    <div class="absolute top-2 left-2">
      <span class="badge {getStatusBadgeClass(video.status)}">
        {getStatusText(video.status)}
      </span>
    </div>

    <!-- Play overlay for completed videos -->
    {#if video.status === 'completed'}
      <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
        <div class="w-14 h-14 rounded-full bg-red-600/90 flex items-center justify-center">
          <svg class="w-6 h-6 text-white ml-1" fill="currentColor" viewBox="0 0 24 24">
            <path d="M8 5v14l11-7z"/>
          </svg>
        </div>
      </div>
    {/if}
  </div>

  <!-- Content -->
  <div class="p-4">
    <h3 class="font-medium text-dark-100 line-clamp-2 group-hover:text-red-400 transition-colors" title={video.title}>
      {video.title || video.videoId}
    </h3>

    {#if showChannel && video.channelName}
      <button
        onclick={handleChannelClick}
        class="text-sm text-dark-400 hover:text-dark-200 mt-1 truncate block text-left"
      >
        {video.channelName}
      </button>
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

    <!-- Error message -->
    {#if error}
      <p class="text-red-400 text-xs mt-2">{error}</p>
    {/if}
  </div>

  <!-- Download Progress (if downloading) -->
  {#if video.status === 'downloading' && downloadProgress}
    <div class="px-4 pb-2">
      <DownloadProgressBar
        percentage={downloadProgress.percentage || 0}
        downloadedBytes={downloadProgress.downloaded_bytes || 0}
        totalBytes={downloadProgress.total_bytes || 0}
        speed={downloadProgress.speed || ''}
        eta={downloadProgress.eta || ''}
        showDetails={true}
        size="sm"
      />
    </div>
  {/if}

  <!-- Actions -->
  <div class="px-4 pb-4 flex gap-2">
    {#if video.status === 'completed'}
      <a
        href="/api/videos/{video.id}/download"
        download
        onclick={(e) => e.stopPropagation()}
        class="flex-1 btn btn-secondary text-sm py-1.5 text-center"
      >
        <svg class="w-4 h-4 mr-1.5 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
        </svg>
        Download
      </a>
    {:else if video.status === 'pending' || video.status === 'failed'}
      <button
        onclick={handleDownload}
        class="flex-1 btn btn-primary text-sm py-1.5"
        disabled={downloading}
      >
        {#if downloading}
          <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-1.5 inline-block"></div>
          Starting...
        {:else}
          <svg class="w-4 h-4 mr-1.5 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          {video.status === 'failed' ? 'Retry' : 'Download'}
        {/if}
      </button>
    {:else if video.status === 'downloading'}
      <div class="flex-1 text-center">
        {#if !downloadProgress}
          <div class="btn btn-secondary text-sm py-1.5 cursor-not-allowed opacity-75 w-full">
            <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-dark-400 mr-1.5 inline-block"></div>
            Downloading...
          </div>
        {:else}
          <span class="text-xs text-dark-500">
            {downloadProgress.speed || 'Downloading...'}
          </span>
        {/if}
      </div>
    {/if}
  </div>
</div>

<style>
  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
</style>
