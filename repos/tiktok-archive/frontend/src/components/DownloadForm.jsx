import { useState } from 'react'
import axios from 'axios'

export default function DownloadForm({ onSuccess }) {
  const [url, setUrl] = useState('')
  const [urls, setUrls] = useState('')
  const [mode, setMode] = useState('single') // single or bulk
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState(null)

  const handleSingleSubmit = async (e) => {
    e.preventDefault()
    if (!url.trim()) return

    setLoading(true)
    setMessage(null)

    try {
      const res = await axios.post('/api/download', { url: url.trim() })
      if (res.data.status === 'exists') {
        setMessage({ type: 'info', text: 'Video already archived!' })
      } else {
        setMessage({ type: 'success', text: 'Download queued!' })
        setUrl('')
      }
      onSuccess?.()
    } catch (err) {
      setMessage({ type: 'error', text: err.response?.data?.detail || 'Download failed' })
    } finally {
      setLoading(false)
    }
  }

  const handleBulkSubmit = async (e) => {
    e.preventDefault()
    const urlList = urls.split('\n').map(u => u.trim()).filter(u => u)
    if (urlList.length === 0) return

    setLoading(true)
    setMessage(null)

    try {
      const res = await axios.post('/api/download/bulk', { urls: urlList })
      const queued = res.data.queued.filter(q => q.status === 'queued').length
      const exists = res.data.queued.filter(q => q.status === 'exists').length
      setMessage({
        type: 'success',
        text: `${queued} queued, ${exists} already archived`
      })
      setUrls('')
      onSuccess?.()
    } catch (err) {
      setMessage({ type: 'error', text: err.response?.data?.detail || 'Bulk download failed' })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="bg-tt-gray rounded-lg p-4">
      {/* Mode toggle */}
      <div className="flex gap-2 mb-4">
        <button
          onClick={() => setMode('single')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            mode === 'single' ? 'bg-tt-pink text-white' : 'bg-gray-700 text-gray-300'
          }`}
        >
          Single URL
        </button>
        <button
          onClick={() => setMode('bulk')}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            mode === 'bulk' ? 'bg-tt-pink text-white' : 'bg-gray-700 text-gray-300'
          }`}
        >
          Bulk Import
        </button>
      </div>

      {mode === 'single' ? (
        <form onSubmit={handleSingleSubmit}>
          <div className="flex gap-2">
            <input
              type="text"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="Paste TikTok URL here..."
              className="flex-1 bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-white placeholder-gray-500 focus:outline-none focus:border-tt-pink"
            />
            <button
              type="submit"
              disabled={loading || !url.trim()}
              className="px-6 py-2 bg-tt-pink text-white rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed hover:bg-pink-600 transition-colors"
            >
              {loading ? <span className="spinner inline-block" /> : 'Archive'}
            </button>
          </div>
        </form>
      ) : (
        <form onSubmit={handleBulkSubmit}>
          <textarea
            value={urls}
            onChange={(e) => setUrls(e.target.value)}
            placeholder="Paste multiple TikTok URLs (one per line)..."
            rows={5}
            className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-white placeholder-gray-500 focus:outline-none focus:border-tt-pink mb-2"
          />
          <button
            type="submit"
            disabled={loading || !urls.trim()}
            className="px-6 py-2 bg-tt-pink text-white rounded-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed hover:bg-pink-600 transition-colors"
          >
            {loading ? <span className="spinner inline-block" /> : 'Archive All'}
          </button>
        </form>
      )}

      {message && (
        <div className={`mt-3 px-4 py-2 rounded-lg text-sm ${
          message.type === 'success' ? 'bg-green-900/50 text-green-400' :
          message.type === 'error' ? 'bg-red-900/50 text-red-400' :
          'bg-blue-900/50 text-blue-400'
        }`}>
          {message.text}
        </div>
      )}
    </div>
  )
}
