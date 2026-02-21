import { useState, useEffect } from 'react'
import { formatDistanceToNow } from 'date-fns'
import axios from 'axios'

export default function Queue() {
  const [queue, setQueue] = useState({ pending: [], recent: [] })
  const [loading, setLoading] = useState(true)

  const fetchQueue = async () => {
    try {
      const res = await axios.get('/api/queue')
      setQueue(res.data)
    } catch (err) {
      console.error('Error fetching queue:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchQueue()
    // Poll every 3 seconds
    const interval = setInterval(fetchQueue, 3000)
    return () => clearInterval(interval)
  }, [])

  const getStatusBadge = (status) => {
    switch (status) {
      case 'pending':
        return <span className="px-2 py-0.5 bg-yellow-900/50 text-yellow-400 rounded text-xs">Pending</span>
      case 'downloading':
        return <span className="px-2 py-0.5 bg-blue-900/50 text-blue-400 rounded text-xs animate-pulse">Downloading</span>
      case 'completed':
        return <span className="px-2 py-0.5 bg-green-900/50 text-green-400 rounded text-xs">Completed</span>
      case 'failed':
        return <span className="px-2 py-0.5 bg-red-900/50 text-red-400 rounded text-xs">Failed</span>
      default:
        return <span className="px-2 py-0.5 bg-gray-700 text-gray-400 rounded text-xs">{status}</span>
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center py-12">
        <div className="spinner" style={{ width: 48, height: 48 }} />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Download Queue</h1>

      {/* Active/Pending */}
      <div className="bg-tt-gray rounded-lg overflow-hidden">
        <div className="px-4 py-3 border-b border-gray-700">
          <h2 className="font-semibold">
            Active Downloads
            {queue.pending.length > 0 && (
              <span className="ml-2 text-tt-cyan">{queue.pending.length}</span>
            )}
          </h2>
        </div>

        {queue.pending.length > 0 ? (
          <div className="divide-y divide-gray-700">
            {queue.pending.map(item => (
              <div key={item.id} className="px-4 py-3 flex items-center gap-4">
                {item.status === 'downloading' && (
                  <div className="spinner" />
                )}
                <div className="flex-1 min-w-0">
                  <div className="text-sm truncate text-gray-300">{item.url}</div>
                </div>
                {getStatusBadge(item.status)}
              </div>
            ))}
          </div>
        ) : (
          <div className="px-4 py-8 text-center text-gray-500">
            No active downloads
          </div>
        )}
      </div>

      {/* Recent */}
      <div className="bg-tt-gray rounded-lg overflow-hidden">
        <div className="px-4 py-3 border-b border-gray-700">
          <h2 className="font-semibold">Recent Activity</h2>
        </div>

        {queue.recent.length > 0 ? (
          <div className="divide-y divide-gray-700">
            {queue.recent.map(item => (
              <div key={item.id} className="px-4 py-3">
                <div className="flex items-center gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="text-sm truncate text-gray-300">{item.url}</div>
                    {item.error && (
                      <div className="text-xs text-red-400 mt-1 truncate">
                        {item.error}
                      </div>
                    )}
                  </div>
                  <div className="text-right">
                    {getStatusBadge(item.status)}
                    {item.completed_at && (
                      <div className="text-xs text-gray-500 mt-1">
                        {formatDistanceToNow(new Date(item.completed_at), { addSuffix: true })}
                      </div>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="px-4 py-8 text-center text-gray-500">
            No recent downloads
          </div>
        )}
      </div>

      {/* Tips */}
      <div className="bg-gray-800/50 rounded-lg p-4">
        <h3 className="font-semibold mb-2">💡 Tips</h3>
        <ul className="text-sm text-gray-400 space-y-1">
          <li>• Videos are downloaded in the background</li>
          <li>• Thumbnails and audio are extracted automatically</li>
          <li>• Failed downloads can be retried from the home page</li>
          <li>• The queue updates every 3 seconds</li>
        </ul>
      </div>
    </div>
  )
}
