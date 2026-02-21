import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import axios from 'axios'
import DownloadForm from '../components/DownloadForm'
import VideoCard from '../components/VideoCard'

export default function Home() {
  const [stats, setStats] = useState(null)
  const [recentVideos, setRecentVideos] = useState([])
  const [loading, setLoading] = useState(true)

  const fetchData = async () => {
    try {
      const [statsRes, videosRes] = await Promise.all([
        axios.get('/api/stats'),
        axios.get('/api/videos?limit=8'),
      ])
      setStats(statsRes.data)
      setRecentVideos(videosRes.data.videos)
    } catch (err) {
      console.error('Error fetching data:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  const formatSize = (bytes) => {
    if (bytes >= 1024 ** 3) return (bytes / (1024 ** 3)).toFixed(2) + ' GB'
    if (bytes >= 1024 ** 2) return (bytes / (1024 ** 2)).toFixed(2) + ' MB'
    return (bytes / 1024).toFixed(2) + ' KB'
  }

  return (
    <div className="space-y-8">
      {/* Hero section */}
      <div className="text-center py-8">
        <h1 className="text-4xl font-bold mb-4">
          <span className="gradient-text">Archive Your TikToks</span>
        </h1>
        <p className="text-gray-400 max-w-lg mx-auto">
          Save TikTok videos locally before they disappear. Download, organize,
          and browse your personal archive.
        </p>
      </div>

      {/* Download form */}
      <DownloadForm onSuccess={fetchData} />

      {/* Stats cards */}
      {stats && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="bg-tt-gray rounded-lg p-4">
            <div className="text-3xl font-bold text-tt-pink">{stats.total_videos}</div>
            <div className="text-sm text-gray-400">Videos Archived</div>
          </div>
          <div className="bg-tt-gray rounded-lg p-4">
            <div className="text-3xl font-bold text-tt-cyan">{stats.total_collections}</div>
            <div className="text-sm text-gray-400">Collections</div>
          </div>
          <div className="bg-tt-gray rounded-lg p-4">
            <div className="text-3xl font-bold text-white">{formatSize(stats.total_size_bytes)}</div>
            <div className="text-sm text-gray-400">Total Storage</div>
          </div>
          <div className="bg-tt-gray rounded-lg p-4">
            <div className="text-3xl font-bold text-yellow-400">{stats.pending_downloads}</div>
            <div className="text-sm text-gray-400">In Queue</div>
          </div>
        </div>
      )}

      {/* Top uploaders */}
      {stats?.top_uploaders?.length > 0 && (
        <div className="bg-tt-gray rounded-lg p-4">
          <h2 className="text-lg font-semibold mb-3">Most Archived Creators</h2>
          <div className="flex flex-wrap gap-2">
            {stats.top_uploaders.slice(0, 8).map((u, i) => (
              <Link
                key={i}
                to={`/browse?uploader=${encodeURIComponent(u.name)}`}
                className="px-3 py-1 bg-gray-700 rounded-full text-sm hover:bg-gray-600 transition-colors"
              >
                @{u.name} <span className="text-gray-400">({u.count})</span>
              </Link>
            ))}
          </div>
        </div>
      )}

      {/* Recent videos */}
      {recentVideos.length > 0 && (
        <div>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold">Recently Archived</h2>
            <Link to="/browse" className="text-tt-pink text-sm hover:underline">
              View all →
            </Link>
          </div>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
            {recentVideos.map(video => (
              <VideoCard key={video.id} video={video} />
            ))}
          </div>
        </div>
      )}

      {/* Empty state */}
      {!loading && recentVideos.length === 0 && (
        <div className="text-center py-12">
          <span className="text-6xl block mb-4">📦</span>
          <h2 className="text-xl font-semibold mb-2">Your archive is empty</h2>
          <p className="text-gray-400">
            Paste a TikTok URL above to start archiving
          </p>
        </div>
      )}
    </div>
  )
}
