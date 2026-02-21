import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { format } from 'date-fns'
import axios from 'axios'

export default function VideoDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [video, setVideo] = useState(null)
  const [loading, setLoading] = useState(true)
  const [deleting, setDeleting] = useState(false)

  useEffect(() => {
    const fetchVideo = async () => {
      try {
        const res = await axios.get(`/api/videos/${id}`)
        setVideo(res.data)
      } catch (err) {
        console.error('Error fetching video:', err)
      } finally {
        setLoading(false)
      }
    }
    fetchVideo()
  }, [id])

  const handleDelete = async () => {
    if (!confirm('Are you sure you want to delete this video from your archive?')) return

    setDeleting(true)
    try {
      await axios.delete(`/api/videos/${id}`)
      navigate('/browse')
    } catch (err) {
      alert('Error deleting video')
      setDeleting(false)
    }
  }

  const formatCount = (num) => {
    if (!num) return '0'
    return num.toLocaleString()
  }

  const formatSize = (bytes) => {
    if (!bytes) return 'Unknown'
    if (bytes >= 1024 ** 2) return (bytes / (1024 ** 2)).toFixed(2) + ' MB'
    return (bytes / 1024).toFixed(2) + ' KB'
  }

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <div className="spinner" style={{ width: 48, height: 48 }} />
      </div>
    )
  }

  if (!video) {
    return (
      <div className="text-center py-12">
        <span className="text-4xl block mb-4">❌</span>
        <h2 className="text-xl font-semibold">Video not found</h2>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto">
      <button
        onClick={() => navigate(-1)}
        className="text-gray-400 hover:text-white mb-4 flex items-center gap-2"
      >
        ← Back
      </button>

      <div className="grid md:grid-cols-2 gap-6">
        {/* Video player */}
        <div>
          <div className="bg-black rounded-lg overflow-hidden aspect-[9/16]">
            {video.video_url ? (
              <video
                src={video.video_url}
                controls
                autoPlay
                loop
                className="w-full h-full object-contain"
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-gray-600">
                <span className="text-4xl">🎬</span>
              </div>
            )}
          </div>

          {/* Audio download */}
          {video.audio_url && (
            <a
              href={video.audio_url}
              download
              className="mt-4 block text-center px-4 py-2 bg-gray-700 rounded-lg hover:bg-gray-600 transition-colors"
            >
              🎵 Download Audio (MP3)
            </a>
          )}
        </div>

        {/* Video info */}
        <div className="space-y-4">
          <div>
            <h1 className="text-xl font-bold mb-2">{video.title || 'Untitled'}</h1>
            <a
              href={video.url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-tt-cyan hover:underline text-sm"
            >
              View original on TikTok ↗
            </a>
          </div>

          {/* Creator */}
          <div className="bg-tt-gray rounded-lg p-4">
            <div className="text-lg font-semibold text-tt-pink">
              @{video.uploader || 'unknown'}
            </div>
            {video.uploader_id && (
              <div className="text-sm text-gray-500">ID: {video.uploader_id}</div>
            )}
          </div>

          {/* Stats */}
          <div className="grid grid-cols-3 gap-2">
            <div className="bg-tt-gray rounded-lg p-3 text-center">
              <div className="text-lg font-bold">{formatCount(video.view_count)}</div>
              <div className="text-xs text-gray-500">Views</div>
            </div>
            <div className="bg-tt-gray rounded-lg p-3 text-center">
              <div className="text-lg font-bold">{formatCount(video.like_count)}</div>
              <div className="text-xs text-gray-500">Likes</div>
            </div>
            <div className="bg-tt-gray rounded-lg p-3 text-center">
              <div className="text-lg font-bold">{formatCount(video.comment_count)}</div>
              <div className="text-xs text-gray-500">Comments</div>
            </div>
          </div>

          {/* Description */}
          {video.description && (
            <div className="bg-tt-gray rounded-lg p-4">
              <h3 className="font-semibold mb-2">Description</h3>
              <p className="text-gray-300 text-sm whitespace-pre-wrap">{video.description}</p>
            </div>
          )}

          {/* Music */}
          {(video.music_title || video.music_author) && (
            <div className="bg-tt-gray rounded-lg p-4">
              <h3 className="font-semibold mb-2">🎵 Sound</h3>
              <div className="text-sm">
                <div className="text-white">{video.music_title}</div>
                {video.music_author && (
                  <div className="text-gray-500">{video.music_author}</div>
                )}
              </div>
            </div>
          )}

          {/* Tags */}
          {video.tags?.length > 0 && (
            <div className="bg-tt-gray rounded-lg p-4">
              <h3 className="font-semibold mb-2">Tags</h3>
              <div className="flex flex-wrap gap-2">
                {video.tags.map((tag, i) => (
                  <span
                    key={i}
                    className="px-2 py-1 bg-gray-700 rounded-full text-xs"
                  >
                    #{tag}
                  </span>
                ))}
              </div>
            </div>
          )}

          {/* Metadata */}
          <div className="bg-tt-gray rounded-lg p-4 space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-gray-500">Duration</span>
              <span>{video.duration ? `${Math.floor(video.duration)}s` : 'Unknown'}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">File Size</span>
              <span>{formatSize(video.file_size)}</span>
            </div>
            {video.upload_date && (
              <div className="flex justify-between">
                <span className="text-gray-500">Upload Date</span>
                <span>{format(new Date(video.upload_date), 'MMM d, yyyy')}</span>
              </div>
            )}
            <div className="flex justify-between">
              <span className="text-gray-500">Archived</span>
              <span>{format(new Date(video.archived_at), 'MMM d, yyyy HH:mm')}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">Original Status</span>
              <span className={video.is_available ? 'text-green-400' : 'text-red-400'}>
                {video.is_available ? 'Available' : 'Removed'}
              </span>
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-2">
            {video.video_url && (
              <a
                href={video.video_url}
                download
                className="flex-1 text-center px-4 py-2 bg-tt-pink rounded-lg hover:bg-pink-600 transition-colors font-medium"
              >
                Download Video
              </a>
            )}
            <button
              onClick={handleDelete}
              disabled={deleting}
              className="px-4 py-2 bg-red-900 text-red-400 rounded-lg hover:bg-red-800 transition-colors disabled:opacity-50"
            >
              {deleting ? '...' : 'Delete'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
