<script>
  import { getJobs, cancelJob, getDownloadsProgress, formatBytes, formatRelativeTime } from '../lib/api.js';
  import DownloadProgressBar from '../components/DownloadProgressBar.svelte';

  let jobs = $state([]);
  let activeDownloads = $state([]);
  let loading = $state(true);
  let error = $state(null);
  let cancellingId = $state(null);

  async function loadJobs() {
    loading = jobs.length === 0;
    error = null;

    try {
      const response = await getJobs();
      jobs = response.jobs || response || [];
    } catch (err) {
      error = err.message;
      jobs = [];
    } finally {
      loading = false;
    }
  }

  async function loadActiveDownloads() {
    try {
      const response = await getDownloadsProgress();
      activeDownloads = response.downloads || [];
    } catch (err) {
      // Silently fail for progress updates - don't overwrite error state
      console.warn('Failed to fetch download progress:', err);
    }
  }

  async function loadAll() {
    await Promise.all([loadJobs(), loadActiveDownloads()]);
  }

  async function handleCancel(job) {
    cancellingId = job.id;

    try {
      await cancelJob(job.id);
      await loadJobs();
    } catch (err) {
      error = err.message;
    } finally {
      cancellingId = null;
    }
  }

  $effect(() => {
    loadAll();
    // Poll for updates every 2 seconds
    const interval = setInterval(loadAll, 2000);
    return () => clearInterval(interval);
  });

  function getStatusBadgeClass(status) {
    switch (status) {
      case 'completed': return 'badge-success';
      case 'running': return 'badge-warning';
      case 'downloading': return 'badge-warning';
      case 'failed':
      case 'cancelled': return 'badge-error';
      case 'queued': return 'badge-info';
      default: return 'badge-neutral';
    }
  }

  function getStatusIcon(status) {
    switch (status) {
      case 'completed':
        return 'M5 13l4 4L19 7';
      case 'running':
      case 'downloading':
        return 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15';
      case 'failed':
      case 'cancelled':
        return 'M6 18L18 6M6 6l12 12';
      case 'queued':
        return 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z';
      default:
        return 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z';
    }
  }

  const runningJobs = $derived(jobs.filter(j => j.status === 'running'));
  const queuedJobs = $derived(jobs.filter(j => j.status === 'queued'));
  const completedJobs = $derived(jobs.filter(j => j.status === 'completed'));
  const failedJobs = $derived(jobs.filter(j => j.status === 'failed' || j.status === 'cancelled'));

  // Combine active downloads count
  const activeDownloadCount = $derived(activeDownloads.length + runningJobs.length);
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
    <div class="hidden lg:block">
      <h1 class="text-2xl font-bold text-dark-100">Jobs</h1>
      <p class="text-dark-400 mt-1">Monitor download workers and progress</p>
    </div>
    <button onclick={loadAll} class="btn btn-secondary">
      <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
      </svg>
      Refresh
    </button>
  </div>

  {#if loading && jobs.length === 0}
    <div class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-red-500"></div>
    </div>
  {:else if error && jobs.length === 0}
    <div class="card p-6">
      <p class="text-red-400">Error loading jobs: {error}</p>
      <button onclick={loadAll} class="btn btn-secondary mt-4">Retry</button>
    </div>
  {:else}
    <!-- Summary Cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="card p-4">
        <div class="flex items-center justify-between">
          <span class="text-dark-400 text-sm">Active Downloads</span>
          <span class="text-yellow-400 font-bold text-xl">{activeDownloads.length}</span>
        </div>
      </div>
      <div class="card p-4">
        <div class="flex items-center justify-between">
          <span class="text-dark-400 text-sm">Queued</span>
          <span class="text-blue-400 font-bold text-xl">{queuedJobs.length}</span>
        </div>
      </div>
      <div class="card p-4">
        <div class="flex items-center justify-between">
          <span class="text-dark-400 text-sm">Completed</span>
          <span class="text-green-400 font-bold text-xl">{completedJobs.length}</span>
        </div>
      </div>
      <div class="card p-4">
        <div class="flex items-center justify-between">
          <span class="text-dark-400 text-sm">Failed</span>
          <span class="text-red-400 font-bold text-xl">{failedJobs.length}</span>
        </div>
      </div>
    </div>

    {#if jobs.length === 0 && activeDownloads.length === 0}
      <div class="card p-12 text-center">
        <div class="w-16 h-16 mx-auto bg-dark-800 rounded-full flex items-center justify-center mb-4">
          <svg class="w-8 h-8 text-dark-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-dark-200 mb-2">No jobs</h3>
        <p class="text-dark-400">Download jobs will appear here when you sync channels or download videos</p>
      </div>
    {:else}
      <!-- Active Downloads (Real-time progress) -->
      {#if activeDownloads.length > 0}
        <div class="card">
          <div class="px-6 py-4 border-b border-dark-800 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-dark-100">Active Downloads</h2>
            <div class="flex items-center gap-2">
              <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
              <span class="text-sm text-dark-400">Live</span>
            </div>
          </div>
          <div class="divide-y divide-dark-800">
            {#each activeDownloads as download}
              <div class="p-6">
                <div class="flex items-start justify-between gap-4 mb-4">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2">
                      <h3 class="font-medium text-dark-100 truncate">{download.video_id}</h3>
                      <span class="badge badge-warning text-xs">{download.status}</span>
                    </div>
                    <p class="text-sm text-dark-400 mt-1">Worker: {download.worker_id || 'Unknown'}</p>
                  </div>
                </div>

                <DownloadProgressBar
                  percentage={download.percentage || 0}
                  downloadedBytes={download.downloaded_bytes || 0}
                  totalBytes={download.total_bytes || 0}
                  speed={download.speed || ''}
                  eta={download.eta || ''}
                  fragment={download.fragment || ''}
                />
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Running Jobs (from job queue) -->
      {#if runningJobs.length > 0}
        <div class="card">
          <div class="px-6 py-4 border-b border-dark-800">
            <h2 class="text-lg font-semibold text-dark-100">Running Jobs</h2>
          </div>
          <div class="divide-y divide-dark-800">
            {#each runningJobs as job}
              <div class="p-6">
                <div class="flex items-start justify-between gap-4 mb-4">
                  <div class="flex-1 min-w-0">
                    <h3 class="font-medium text-dark-100 truncate">{job.title || job.videoId}</h3>
                    <p class="text-sm text-dark-400 mt-1">{job.channelName || 'Unknown channel'}</p>
                  </div>
                  <button
                    onclick={() => handleCancel(job)}
                    class="btn btn-ghost text-red-400 hover:text-red-300 hover:bg-red-900/20 shrink-0"
                    disabled={cancellingId === job.id}
                  >
                    {#if cancellingId === job.id}
                      <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-red-400"></div>
                    {:else}
                      Cancel
                    {/if}
                  </button>
                </div>

                <DownloadProgressBar
                  percentage={job.progress || 0}
                  downloadedBytes={job.downloadedBytes || 0}
                  totalBytes={job.totalBytes || 0}
                  speed={job.speed || ''}
                  eta={job.eta || ''}
                />
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Queued Jobs -->
      {#if queuedJobs.length > 0}
        <div class="card">
          <div class="px-6 py-4 border-b border-dark-800">
            <h2 class="text-lg font-semibold text-dark-100">Queued</h2>
          </div>
          <div class="divide-y divide-dark-800">
            {#each queuedJobs as job}
              <div class="px-6 py-4 flex items-center justify-between gap-4">
                <div class="flex items-center gap-3 min-w-0">
                  <div class="w-8 h-8 rounded-full bg-blue-900/30 flex items-center justify-center shrink-0">
                    <svg class="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getStatusIcon('queued')} />
                    </svg>
                  </div>
                  <div class="min-w-0">
                    <p class="text-dark-200 truncate">{job.title || job.videoId}</p>
                    <p class="text-xs text-dark-500">{job.channelName || 'Unknown channel'}</p>
                  </div>
                </div>
                <button
                  onclick={() => handleCancel(job)}
                  class="btn btn-ghost text-sm text-dark-400 hover:text-red-400"
                  disabled={cancellingId === job.id}
                >
                  Cancel
                </button>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Recent Jobs -->
      {#if completedJobs.length > 0 || failedJobs.length > 0}
        <div class="card">
          <div class="px-6 py-4 border-b border-dark-800">
            <h2 class="text-lg font-semibold text-dark-100">Recent Jobs</h2>
          </div>
          <div class="divide-y divide-dark-800">
            {#each [...completedJobs, ...failedJobs].slice(0, 20) as job}
              <div class="px-6 py-4 flex items-center justify-between gap-4">
                <div class="flex items-center gap-3 min-w-0">
                  <div class="w-8 h-8 rounded-full {job.status === 'completed' ? 'bg-green-900/30' : 'bg-red-900/30'} flex items-center justify-center shrink-0">
                    <svg class="w-4 h-4 {job.status === 'completed' ? 'text-green-400' : 'text-red-400'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d={getStatusIcon(job.status)} />
                    </svg>
                  </div>
                  <div class="min-w-0">
                    <p class="text-dark-200 truncate">{job.title || job.videoId}</p>
                    <p class="text-xs text-dark-500">
                      {job.channelName || 'Unknown channel'}
                      {#if job.completedAt}
                        &middot; {formatRelativeTime(job.completedAt)}
                      {/if}
                    </p>
                  </div>
                </div>
                <span class="badge {getStatusBadgeClass(job.status)} shrink-0">
                  {job.status}
                </span>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    {/if}
  {/if}
</div>
