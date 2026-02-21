<script>
  import { getVideo, getVideoStreamUrl, getVideoFiles, formatDuration, formatBytes, formatRelativeTime } from '../lib/api.js';

  let { videoId, navigate } = $props();

  let video = $state(null);
  let files = $state([]);
  let loading = $state(true);
  let error = $state(null);
  let showDescription = $state(false);

  async function loadVideo() {
    loading = true;
    error = null;

    try {
      const [videoData, filesData] = await Promise.all([
        getVideo(videoId),
        getVideoFiles(videoId).catch(() => ({ files: [] }))
      ]);
      video = videoData;
      files = filesData.files || [];
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (videoId) {
      loadVideo();
    }
  });

  function handleBack() {
    navigate('videos');
  }

  function handleChannelClick() {
    if (video?.channelId) {
      navigate('channel', { id: video.channelId });
    }
  }
</script>

<div class="max-w-6xl mx-auto space-y-6">
  <!-- Back button -->
  <button
    onclick={handleBack}
    class="flex items-center gap-2 text-dark-400 hover:text-dark-200 transition-colors"
  >
    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
    </svg>
    Back to Videos
  </button>

  {#if loading}
    <div class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-red-500"></div>
    </div>
  {:else if error}
    <div class="card p-6">
      <p class="text-red-400">Error loading video: {error}</p>
      <button onclick={loadVideo} class="btn btn-secondary mt-4">Retry</button>
    </div>
  {:else if video}
    <!-- Video Player -->
    <div class="card overflow-hidden">
      <div class="aspect-video bg-black">
        {#if video.status === 'completed'}
          <video
            controls
            autoplay
            class="w-full h-full"
            poster={video.thumbnailUrl}
          >
            <source src={getVideoStreamUrl(videoId)} type="video/mp4" />
            Your browser does not support the video tag.
          </video>
        {:else}
          <div class="w-full h-full flex items-center justify-center bg-dark-900">
            {#if video.thumbnailUrl}
              <img
                src={video.thumbnailUrl}
                alt={video.title}
                class="max-w-full max-h-full object-contain opacity-50"
              />
            {/if}
            <div class="absolute inset-0 flex items-center justify-center">
              <div class="text-center">
                <div class="w-16 h-16 mx-auto bg-dark-800 rounded-full flex items-center justify-center mb-4">
                  <svg class="w-8 h-8 text-dark-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <p class="text-dark-400">Video not available</p>
                <p class="text-dark-500 text-sm mt-1">Status: {video.status}</p>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>

    <!-- Video Info -->
    <div class="card p-6">
      <h1 class="text-xl font-bold text-dark-100">{video.title}</h1>

      <div class="flex flex-wrap items-center gap-4 mt-4 text-sm text-dark-400">
        {#if video.channelName}
          <button
            onclick={handleChannelClick}
            class="flex items-center gap-2 hover:text-dark-200 transition-colors"
          >
            <div class="w-8 h-8 bg-dark-700 rounded-full flex items-center justify-center">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
            <span class="font-medium">{video.channelName}</span>
          </button>
        {/if}

        {#if video.viewCount}
          <span>{video.viewCount.toLocaleString()} views</span>
        {/if}

        {#if video.publishedAt}
          <span>{formatRelativeTime(video.publishedAt)}</span>
        {/if}

        {#if video.duration}
          <span>{formatDuration(video.duration)}</span>
        {/if}
      </div>

      <!-- Description -->
      {#if video.description}
        <div class="mt-6">
          <button
            onclick={() => showDescription = !showDescription}
            class="flex items-center gap-2 text-dark-400 hover:text-dark-200 transition-colors"
          >
            <span class="font-medium">Description</span>
            <svg
              class="w-4 h-4 transition-transform {showDescription ? 'rotate-180' : ''}"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>

          {#if showDescription}
            <div class="mt-3 p-4 bg-dark-800 rounded-lg">
              <p class="text-dark-300 whitespace-pre-wrap text-sm">{video.description}</p>
            </div>
          {/if}
        </div>
      {/if}
    </div>

    <!-- File Info -->
    {#if video.status === 'completed' && (video.fileSize || files.length > 0)}
      <div class="card p-6">
        <h2 class="text-lg font-semibold text-dark-100 mb-4">File Information</h2>

        <div class="space-y-3">
          {#if video.fileSize}
            <div class="flex justify-between text-sm">
              <span class="text-dark-400">File Size</span>
              <span class="text-dark-200">{formatBytes(video.fileSize)}</span>
            </div>
          {/if}

          {#if video.filePath}
            <div class="flex justify-between text-sm">
              <span class="text-dark-400">File Path</span>
              <span class="text-dark-200 truncate max-w-xs" title={video.filePath}>{video.filePath}</span>
            </div>
          {/if}

          {#if files.length > 0}
            <div class="mt-4">
              <h3 class="text-sm font-medium text-dark-300 mb-2">Available Files</h3>
              <div class="space-y-2">
                {#each files as file}
                  <div class="flex items-center justify-between p-2 bg-dark-800 rounded">
                    <span class="text-sm text-dark-300">{file.name}</span>
                    <span class="text-xs text-dark-500">{formatBytes(file.size)}</span>
                  </div>
                {/each}
              </div>
            </div>
          {/if}
        </div>

        <!-- Download Button -->
        <div class="mt-6">
          <a
            href="/api/videos/{videoId}/stream"
            download
            class="btn btn-primary"
          >
            <svg class="w-4 h-4 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            Download Video
          </a>
        </div>
      </div>
    {/if}
  {/if}
</div>
