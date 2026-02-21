import { Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import Home from './pages/Home'
import Browse from './pages/Browse'
import VideoDetail from './pages/VideoDetail'
import Collections from './pages/Collections'
import Queue from './pages/Queue'

function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/browse" element={<Browse />} />
        <Route path="/video/:id" element={<VideoDetail />} />
        <Route path="/collections" element={<Collections />} />
        <Route path="/queue" element={<Queue />} />
      </Routes>
    </Layout>
  )
}

export default App
