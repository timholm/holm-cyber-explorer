import { Link } from 'react-router-dom'
import { formatDistanceToNow } from 'date-fns'

export default function VideoCard({ video }) {
  const formatDuration = (seconds) => {
    if (!seconds) return '0:00'
    const mins = Math.floor(seconds / 60)
    const secs = Math.floor(seconds % 60)
    return `${mins}:${secs.toString().padStart(2, '0')}`
  }

  const formatCount = (num) => {
    if (!num) return '0'
    if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
    if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
    return num.toString()
  }

  return (
    <Link to={`/video/${video.id}`} className="video-card block bg-tt-gray rounded-lg overflow-hidden">
      {/* Thumbnail */}
      <div className="relative aspect-[9/16] bg-gray-900">
        {video.thumbnail_url ? (
          <img
            src={video.thumbnail_url}
            alt={video.title || 'Video thumbnail'}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-gray-600">
            <span className="text-4xl">🎬</span>
          </div>
        )}

        {/* Duration badge */}
        <div className="absolute bottom-2 right-2 bg-black/80 px-2 py-0.5 rounded text-xs">
          {formatDuration(video.duration)}
        </div>

        {/* Play overlay on hover */}
        <div className="absolute inset-0 flex items-center justify-center opacity-0 hover:opacity-100 transition-opacity bg-black/30">
          <div className="w-12 h-12 rounded-full bg-white/90 flex items-center justify-center">
            <span className="text-tt-pink text-xl ml-1">▶</span>
          </div>
        </div>
      </div>

      {/* Info */}
      <div className="p-3">
        <h3 className="font-medium text-sm line-clamp-2 mb-2">
          {video.title || 'Untitled'}
        </h3>

        <div className="flex items-center gap-2 text-xs text-gray-400 mb-2">
          <span className="font-medium text-tt-cyan">@{video.uploader || 'unknown'}</span>
        </div>

        <div className="flex items-center gap-3 text-xs text-gray-500">
          {video.view_count && (
            <span>👁 {formatCount(video.view_count)}</span>
          )}
          {video.like_count && (
            <span>❤️ {formatCount(video.like_count)}</span>
          )}
        </div>

        {video.archived_at && (
          <div className="mt-2 text-xs text-gray-600">
            Archived {formatDistanceToNow(new Date(video.archived_at), { addSuffix: true })}
          </div>
        )}
      </div>
    </Link>
  )
}
