<script>
  import { getStats, getJobs, getRecentActivity, getDownloadsProgress, formatBytes, formatRelativeTime } from '../lib/api.js';
  import DownloadProgressBar from '../components/DownloadProgressBar.svelte';

  let { navigate } = $props();

  let stats = $state(null);
  let jobs = $state([]);
  let activity = $state([]);
  let activeDownloads = $state([]);
  let loading = $state(true);
  let error = $state(null);

  async function loadData() {
    loading = true;
    error = null;

    try {
      const [statsData, jobsData, activityData, progressData] = await Promise.all([
        getStats().catch(() => ({
          totalChannels: 0,
          totalVideos: 0,
          completedVideos: 0,
          pendingVideos: 0,
          failedVideos: 0,
          storageUsed: 0
        })),
        getJobs().catch(() => ({ jobs: [] })),
        getRecentActivity().catch(() => ({ activity: [] })),
        getDownloadsProgress().catch(() => ({ downloads: [] }))
      ]);

      stats = statsData;
      jobs = jobsData.jobs || jobsData || [];
      activity = activityData.activity || activityData || [];
      activeDownloads = progressData.downloads || [];
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  async function loadProgress() {
    try {
      const progressData = await getDownloadsProgress();
      activeDownloads = progressData.downloads || [];
    } catch (err) {
      // Silently fail
    }
  }

  $effect(() => {
    loadData();
    // Poll for updates every 10 seconds for stats, but faster for downloads
    const statsInterval = setInterval(loadData, 10000);
    const progressInterval = setInterval(loadProgress, 2000);
    return () => {
      clearInterval(statsInterval);
      clearInterval(progressInterval);
    };
  });

  const statCards = $derived([
    { label: 'Total Channels', value: stats?.totalChannels || stats?.total_channels || 0, icon: 'users', color: 'text-blue-400' },
    { label: 'Total Videos', value: stats?.totalVideos || stats?.total_videos || 0, icon: 'play', color: 'text-green-400' },
    { label: 'Active Downloads', value: activeDownloads.length || 0, icon: 'download', color: 'text-yellow-400' },
    { label: 'Storage Used', value: formatBytes(stats?.storageUsed || stats?.storage_used || 0), icon: 'database', color: 'text-purple-400' }
  ]);

  function getIconPath(icon) {
    const icons = {
      users: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z',
      play: 'M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
      download: 'M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4',
      database: 'M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4'
    };
    return icons[icon] || icons.play;
  }

  function getStatusBadgeClass(status) {
    switch (status) {
      case 'completed': return 'badge-success';
      case 'downloading':
      case 'running': return 'badge-warning';
      case 'failed': return 'badge-error';
      default: return 'badge-neutral';
    }
  }
</script>

<div class="space-y-8">
  <!-- Header -->
  <div class="hidden lg:block">
    <h1 class="text-2xl font-bold text-dark-100">Dashboard</h1>
    <p class="text-dark-400 mt-1">Overview of your YouTube archive</p>
  </div>

  {#if loading && !stats}
    <div class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-red-500"></div>
    </div>
  {:else if error}
    <div class="card p-6">
      <p class="text-red-400">Error loading dashboard: {error}</p>
      <button onclick={loadData} class="btn btn-secondary mt-4">Retry</button>
    </div>
  {:else}
    <!-- Stats Grid -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      {#each statCards as card}
        <div class="card p-5">
          <div class="flex items-start justify-between">
            <div>
              <p class="text-dark-400 text-sm">{card.label}</p>
              <p class="text-2xl font-bold text-dark-100 mt-1">{card.value}</p>
            </div>
            <div class="{card.color} p-2 bg-dark-800 rounded-lg">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d={getIconPath(card.icon)} />
              </svg>
            </div>
          </div>
        </div>
      {/each}
    </div>

    <!-- Active Downloads -->
    {#if activeDownloads.length > 0}
      <div class="card">
        <div class="px-6 py-4 border-b border-dark-800 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <h2 class="text-lg font-semibold text-dark-100">Active Downloads</h2>
            <div class="flex items-center gap-1.5">
              <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
              <span class="text-xs text-dark-400">Live</span>
            </div>
          </div>
          <button onclick={() => navigate('jobs')} class="text-sm text-red-500 hover:text-red-400">
            View all
          </button>
        </div>
        <div class="divide-y divide-dark-800">
          {#each activeDownloads.slice(0, 5) as download}
            <div class="px-6 py-4">
              <div class="flex items-center justify-between mb-3">
                <p class="text-dark-100 font-medium truncate flex-1 mr-4">{download.video_id}</p>
                <span class="badge badge-warning">{download.status || 'downloading'}</span>
              </div>
              <DownloadProgressBar
                percentage={download.percentage || 0}
                downloadedBytes={download.downloaded_bytes || 0}
                totalBytes={download.total_bytes || 0}
                speed={download.speed || ''}
                eta={download.eta || ''}
                size="sm"
              />
            </div>
          {/each}
        </div>
      </div>
    {:else if jobs?.length > 0}
      <div class="card">
        <div class="px-6 py-4 border-b border-dark-800 flex items-center justify-between">
          <h2 class="text-lg font-semibold text-dark-100">Recent Jobs</h2>
          <button onclick={() => navigate('jobs')} class="text-sm text-red-500 hover:text-red-400">
            View all
          </button>
        </div>
        <div class="divide-y divide-dark-800">
          {#each jobs.slice(0, 5) as job}
            <div class="px-6 py-4">
              <div class="flex items-center justify-between mb-2">
                <p class="text-dark-100 font-medium truncate flex-1 mr-4">{job.title || job.videoId}</p>
                <span class="badge {getStatusBadgeClass(job.status)}">{job.status}</span>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Quick Stats -->
    <div class="grid lg:grid-cols-2 gap-6">
      <!-- Video Status Breakdown -->
      <div class="card p-6">
        <h2 class="text-lg font-semibold text-dark-100 mb-4">Video Status</h2>
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <span class="text-dark-400">Completed</span>
            <span class="text-green-400 font-medium">{stats?.completedVideos || stats?.downloaded_videos || 0}</span>
          </div>
          <div class="w-full bg-dark-800 rounded-full h-2">
            <div
              class="bg-green-500 h-2 rounded-full transition-all duration-500"
              style="width: {(stats?.totalVideos || stats?.total_videos) ? ((stats?.completedVideos || stats?.downloaded_videos || 0) / (stats?.totalVideos || stats?.total_videos) * 100) : 0}%"
            ></div>
          </div>

          <div class="flex items-center justify-between">
            <span class="text-dark-400">Pending</span>
            <span class="text-yellow-400 font-medium">{stats?.pendingVideos || stats?.pending_videos || 0}</span>
          </div>
          <div class="w-full bg-dark-800 rounded-full h-2">
            <div
              class="bg-yellow-500 h-2 rounded-full transition-all duration-500"
              style="width: {(stats?.totalVideos || stats?.total_videos) ? ((stats?.pendingVideos || stats?.pending_videos || 0) / (stats?.totalVideos || stats?.total_videos) * 100) : 0}%"
            ></div>
          </div>

          <div class="flex items-center justify-between">
            <span class="text-dark-400">Failed</span>
            <span class="text-red-400 font-medium">{stats?.failedVideos || stats?.failed_videos || 0}</span>
          </div>
          <div class="w-full bg-dark-800 rounded-full h-2">
            <div
              class="bg-red-500 h-2 rounded-full transition-all duration-500"
              style="width: {(stats?.totalVideos || stats?.total_videos) ? ((stats?.failedVideos || stats?.failed_videos || 0) / (stats?.totalVideos || stats?.total_videos) * 100) : 0}%"
            ></div>
          </div>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="card p-6">
        <h2 class="text-lg font-semibold text-dark-100 mb-4">Recent Activity</h2>
        {#if activity?.length > 0}
          <div class="space-y-3">
            {#each activity.slice(0, 5) as item}
              <div class="flex items-start gap-3">
                <div class="w-2 h-2 rounded-full mt-2 {
                  item.type === 'download_complete' ? 'bg-green-500' :
                  item.type === 'download_failed' ? 'bg-red-500' :
                  item.type === 'sync_complete' ? 'bg-blue-500' : 'bg-dark-500'
                }"></div>
                <div class="flex-1 min-w-0">
                  <p class="text-dark-200 text-sm truncate">{item.message}</p>
                  <p class="text-dark-500 text-xs">{formatRelativeTime(item.timestamp)}</p>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <p class="text-dark-500 text-sm">No recent activity</p>
        {/if}
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="card p-6">
      <h2 class="text-lg font-semibold text-dark-100 mb-4">Quick Actions</h2>
      <div class="flex flex-wrap gap-3">
        <button onclick={() => navigate('channels')} class="btn btn-primary">
          <svg class="w-4 h-4 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Add Channel
        </button>
        <button onclick={() => navigate('videos')} class="btn btn-secondary">
          Browse Videos
        </button>
        <button onclick={() => navigate('jobs')} class="btn btn-secondary">
          View Jobs
        </button>
      </div>
    </div>
  {/if}
</div>
