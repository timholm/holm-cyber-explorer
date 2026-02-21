import { Link, useLocation } from 'react-router-dom'

export default function Layout({ children }) {
  const location = useLocation()

  const navItems = [
    { path: '/', label: 'Home', icon: '🏠' },
    { path: '/browse', label: 'Browse', icon: '📼' },
    { path: '/collections', label: 'Collections', icon: '📁' },
    { path: '/queue', label: 'Queue', icon: '⏳' },
  ]

  return (
    <div className="min-h-screen bg-tt-dark">
      {/* Header */}
      <header className="border-b border-gray-800 bg-tt-gray sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-2">
            <span className="text-2xl">📦</span>
            <span className="text-xl font-bold gradient-text">TikTok Archive</span>
          </Link>

          <nav className="flex items-center gap-1">
            {navItems.map(item => (
              <Link
                key={item.path}
                to={item.path}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                  location.pathname === item.path
                    ? 'bg-tt-pink text-white'
                    : 'text-gray-400 hover:text-white hover:bg-gray-800'
                }`}
              >
                <span className="mr-2">{item.icon}</span>
                {item.label}
              </Link>
            ))}
          </nav>
        </div>
      </header>

      {/* Main content */}
      <main className="max-w-7xl mx-auto px-4 py-6">
        {children}
      </main>

      {/* Footer */}
      <footer className="border-t border-gray-800 mt-auto py-4 text-center text-gray-500 text-sm">
        TikTok Archive - Local backup tool
      </footer>
    </div>
  )
}
