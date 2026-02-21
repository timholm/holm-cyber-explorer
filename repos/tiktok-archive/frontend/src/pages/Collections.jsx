import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import axios from 'axios'

export default function Collections() {
  const [collections, setCollections] = useState([])
  const [loading, setLoading] = useState(true)
  const [showCreate, setShowCreate] = useState(false)
  const [newName, setNewName] = useState('')
  const [newDesc, setNewDesc] = useState('')
  const [creating, setCreating] = useState(false)

  const fetchCollections = async () => {
    try {
      const res = await axios.get('/api/collections')
      setCollections(res.data.collections)
    } catch (err) {
      console.error('Error fetching collections:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchCollections()
  }, [])

  const handleCreate = async (e) => {
    e.preventDefault()
    if (!newName.trim()) return

    setCreating(true)
    try {
      await axios.post('/api/collections', {
        name: newName.trim(),
        description: newDesc.trim() || null,
      })
      setNewName('')
      setNewDesc('')
      setShowCreate(false)
      fetchCollections()
    } catch (err) {
      alert('Error creating collection')
    } finally {
      setCreating(false)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Collections</h1>
        <button
          onClick={() => setShowCreate(!showCreate)}
          className="px-4 py-2 bg-tt-pink rounded-lg hover:bg-pink-600 transition-colors font-medium"
        >
          {showCreate ? 'Cancel' : '+ New Collection'}
        </button>
      </div>

      {/* Create form */}
      {showCreate && (
        <form onSubmit={handleCreate} className="bg-tt-gray rounded-lg p-4 space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Name</label>
            <input
              type="text"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              placeholder="My Collection"
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-white placeholder-gray-500 focus:outline-none focus:border-tt-pink"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Description (optional)</label>
            <textarea
              value={newDesc}
              onChange={(e) => setNewDesc(e.target.value)}
              placeholder="What's this collection for?"
              rows={2}
              className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-white placeholder-gray-500 focus:outline-none focus:border-tt-pink"
            />
          </div>
          <button
            type="submit"
            disabled={creating || !newName.trim()}
            className="px-6 py-2 bg-tt-pink text-white rounded-lg font-medium disabled:opacity-50 hover:bg-pink-600 transition-colors"
          >
            {creating ? 'Creating...' : 'Create Collection'}
          </button>
        </form>
      )}

      {/* Collections list */}
      {loading ? (
        <div className="flex justify-center py-12">
          <div className="spinner" style={{ width: 48, height: 48 }} />
        </div>
      ) : collections.length > 0 ? (
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {collections.map(collection => (
            <Link
              key={collection.id}
              to={`/browse?collection_id=${collection.id}`}
              className="bg-tt-gray rounded-lg p-4 hover:bg-gray-800 transition-colors block"
            >
              <div className="flex items-start justify-between mb-2">
                <h3 className="font-semibold text-lg">{collection.name}</h3>
                <span className="text-tt-cyan font-mono text-sm">
                  {collection.video_count}
                </span>
              </div>
              {collection.description && (
                <p className="text-gray-400 text-sm mb-2 line-clamp-2">
                  {collection.description}
                </p>
              )}
              <div className="text-xs text-gray-500">
                Created {new Date(collection.created_at).toLocaleDateString()}
              </div>
            </Link>
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <span className="text-4xl block mb-4">📁</span>
          <h2 className="text-xl font-semibold mb-2">No collections yet</h2>
          <p className="text-gray-400">
            Create a collection to organize your archived videos
          </p>
        </div>
      )}
    </div>
  )
}
