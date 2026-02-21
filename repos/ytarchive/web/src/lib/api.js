const API_BASE = '/api';

class ApiError extends Error {
  constructor(message, status, data = null) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.data = data;
  }
}

async function request(endpoint, options = {}) {
  const url = `${API_BASE}${endpoint}`;
  const config = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    },
    ...options
  };

  try {
    const response = await fetch(url, config);

    if (!response.ok) {
      let errorData = null;
      try {
        errorData = await response.json();
      } catch {
        // Response may not be JSON
      }
      throw new ApiError(
        errorData?.error || errorData?.message || `HTTP ${response.status}`,
        response.status,
        errorData
      );
    }

    // Handle empty responses
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      return await response.json();
    }
    return null;
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    throw new ApiError(error.message || 'Network error', 0);
  }
}

// Channels API
export async function getChannels() {
  return request('/channels');
}

export async function getChannel(id) {
  return request(`/channels/${id}`);
}

export async function addChannel(url) {
  return request('/channels', {
    method: 'POST',
    body: JSON.stringify({ youtube_url: url })
  });
}

export async function deleteChannel(id) {
  return request(`/channels/${id}`, {
    method: 'DELETE'
  });
}

export async function triggerSync(channelId) {
  return request(`/channels/${channelId}/sync`, {
    method: 'POST'
  });
}

// Search API (FTS5 full-text search)
export async function searchVideos(query, params = {}) {
  const searchParams = new URLSearchParams();
  searchParams.set('q', query);
  if (params.channel) searchParams.set('channel', params.channel);
  if (params.status) searchParams.set('status', params.status);
  if (params.limit) searchParams.set('limit', params.limit);

  return request(`/search?${searchParams.toString()}`);
}

// Videos API
export async function getChannelVideos(channelId, params = {}) {
  const searchParams = new URLSearchParams();
  if (params.status) searchParams.set('status', params.status);
  if (params.limit) searchParams.set('limit', params.limit);
  if (params.offset) searchParams.set('offset', params.offset);

  const query = searchParams.toString();
  return request(`/channels/${channelId}/videos${query ? `?${query}` : ''}`);
}

export async function getAllVideos(params = {}) {
  const searchParams = new URLSearchParams();
  if (params.search) searchParams.set('search', params.search);
  if (params.status) searchParams.set('status', params.status);
  if (params.limit) searchParams.set('limit', params.limit);
  if (params.offset) searchParams.set('offset', params.offset);

  const query = searchParams.toString();
  return request(`/videos${query ? `?${query}` : ''}`);
}

export async function getVideo(id) {
  return request(`/videos/${id}`);
}

export async function downloadVideo(id) {
  return request(`/videos/${id}/download`, {
    method: 'POST'
  });
}

// Jobs API
export async function getJobs() {
  return request('/jobs');
}

export async function getJob(id) {
  return request(`/jobs/${id}`);
}

export async function cancelJob(id) {
  return request(`/jobs/${id}/cancel`, {
    method: 'POST'
  });
}

export async function getProgress() {
  return request('/jobs/progress');
}

// Get real-time download progress for all active downloads
export async function getDownloadsProgress() {
  return request('/downloads/progress');
}

// Stats API
export async function getStats() {
  return request('/stats');
}

export async function getRecentActivity() {
  return request('/activity');
}

// Utility to format bytes
export function formatBytes(bytes) {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Utility to format duration
export function formatDuration(seconds) {
  if (!seconds) return '0:00';
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = Math.floor(seconds % 60);

  if (h > 0) {
    return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
  }
  return `${m}:${s.toString().padStart(2, '0')}`;
}

// Utility to format relative time
export function formatRelativeTime(dateString) {
  const date = new Date(dateString);
  const now = new Date();
  const diffInSeconds = Math.floor((now - date) / 1000);

  if (diffInSeconds < 60) return 'just now';
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`;
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h ago`;
  if (diffInSeconds < 604800) return `${Math.floor(diffInSeconds / 86400)}d ago`;

  return date.toLocaleDateString();
}

// Cookies API
export async function getCookies() {
  return request('/cookies');
}

export async function saveCookies(cookies) {
  return request('/cookies', {
    method: 'POST',
    body: JSON.stringify({ cookies })
  });
}

export async function deleteCookies() {
  return request('/cookies', {
    method: 'DELETE'
  });
}

// Video metadata and streaming
export function getVideoStreamUrl(videoId) {
  return `${API_BASE}/videos/${videoId}/stream`;
}

export function getVideoThumbnailUrl(videoId) {
  return `${API_BASE}/videos/${videoId}/thumbnail`;
}

export async function getVideoMetadata(videoId) {
  return request(`/videos/${videoId}/metadata`);
}

export async function getVideoFiles(videoId) {
  return request(`/videos/${videoId}/files`);
}

export { ApiError };
