import { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import axios from 'axios'
import VideoCard from '../components/VideoCard'

export default function Browse() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [videos, setVideos] = useState([])
  const [loading, setLoading] = useState(true)
  const [pagination, setPagination] = useState({ total: 0, page: 1, pages: 1 })

  const search = searchParams.get('search') || ''
  const uploader = searchParams.get('uploader') || ''
  const sort = searchParams.get('sort') || 'archived_at'
  const page = parseInt(searchParams.get('page') || '1', 10)

  const fetchVideos = async () => {
    setLoading(true)
    try {
      const params = new URLSearchParams()
      if (search) params.set('search', search)
      if (uploader) params.set('uploader', uploader)
      params.set('sort', sort)
      params.set('page', page.toString())
      params.set('limit', '20')

      const res = await axios.get(`/api/videos?${params}`)
      setVideos(res.data.videos)
      setPagination({
        total: res.data.total,
        page: res.data.page,
        pages: res.data.pages,
      })
    } catch (err) {
      console.error('Error fetching videos:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchVideos()
  }, [search, uploader, sort, page])

  const updateParams = (updates) => {
    const newParams = new URLSearchParams(searchParams)
    Object.entries(updates).forEach(([k, v]) => {
      if (v) newParams.set(k, v)
      else newParams.delete(k)
    })
    setSearchParams(newParams)
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Browse Archive</h1>

      {/* Filters */}
      <div className="flex flex-wrap gap-4 bg-tt-gray rounded-lg p-4">
        <div className="flex-1 min-w-[200px]">
          <input
            type="text"
            placeholder="Search videos..."
            value={search}
            onChange={(e) => updateParams({ search: e.target.value, page: '1' })}
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-white placeholder-gray-500 focus:outline-none focus:border-tt-pink"
          />
        </div>

        <div className="flex-1 min-w-[200px]">
          <input
            type="text"
            placeholder="Filter by creator..."
            value={uploader}
            onChange={(e) => updateParams({ uploader: e.target.value, page: '1' })}
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-white placeholder-gray-500 focus:outline-none focus:border-tt-pink"
          />
        </div>

        <select
          value={sort}
          onChange={(e) => updateParams({ sort: e.target.value, page: '1' })}
          className="bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-white focus:outline-none focus:border-tt-pink"
        >
          <option value="archived_at">Recently Archived</option>
          <option value="upload_date">Upload Date</option>
          <option value="view_count">Most Views</option>
          <option value="like_count">Most Likes</option>
          <option value="duration">Duration</option>
        </select>
      </div>

      {/* Results count */}
      <div className="text-sm text-gray-400">
        {pagination.total} videos found
        {(search || uploader) && (
          <button
            onClick={() => updateParams({ search: '', uploader: '' })}
            className="ml-2 text-tt-pink hover:underline"
          >
            Clear filters
          </button>
        )}
      </div>

      {/* Video grid */}
      {loading ? (
        <div className="flex justify-center py-12">
          <div className="spinner" style={{ width: 48, height: 48 }} />
        </div>
      ) : videos.length > 0 ? (
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
          {videos.map(video => (
            <VideoCard key={video.id} video={video} />
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <span className="text-4xl block mb-4">🔍</span>
          <p className="text-gray-400">No videos found</p>
        </div>
      )}

      {/* Pagination */}
      {pagination.pages > 1 && (
        <div className="flex justify-center gap-2">
          <button
            onClick={() => updateParams({ page: (page - 1).toString() })}
            disabled={page <= 1}
            className="px-4 py-2 bg-gray-700 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition-colors"
          >
            Previous
          </button>
          <span className="px-4 py-2 text-gray-400">
            Page {page} of {pagination.pages}
          </span>
          <button
            onClick={() => updateParams({ page: (page + 1).toString() })}
            disabled={page >= pagination.pages}
            className="px-4 py-2 bg-gray-700 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-600 transition-colors"
          >
            Next
          </button>
        </div>
      )}
    </div>
  )
}
