'use client'

import { useState } from 'react'
import { Settings, Key, Bell, Palette, Save, RefreshCw } from 'lucide-react'
import { Navigation } from '@/components/Navigation'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Select } from '@/components/ui/Select'
import { GlitchText } from '@/components/ui/GlitchText'
import { useNodes } from '@/hooks/useNodes'

export default function SettingsPage() {
  const { nodes } = useNodes()

  const [lokiUrl, setLokiUrl] = useState('http://loki.monitoring:3100')
  const [refreshInterval, setRefreshInterval] = useState('5')
  const [theme, setTheme] = useState('animus')
  const [notifications, setNotifications] = useState(true)

  const handleSave = () => {
    // Save settings to localStorage or API
    localStorage.setItem('animus-settings', JSON.stringify({
      lokiUrl,
      refreshInterval,
      theme,
      notifications,
    }))
  }

  return (
    <div className="min-h-screen">
      <Navigation nodes={nodes} />

      <main className="max-w-screen-xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-center gap-4 mb-8">
          <div className="p-3 rounded-lg bg-animus-gold/10">
            <Settings className="w-8 h-8 text-animus-gold" />
          </div>
          <div>
            <GlitchText
              text="SETTINGS"
              as="h1"
              className="text-3xl font-bold text-animus-gold tracking-widest"
            />
            <p className="text-animus-text-secondary mt-1">
              Configure your Animus dashboard
            </p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Connection Settings */}
          <Card variant="default">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Key className="w-5 h-5" />
                Connection Settings
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <Input
                label="Loki URL"
                value={lokiUrl}
                onChange={(e) => setLokiUrl(e.target.value)}
                placeholder="http://loki.monitoring:3100"
              />
              <Select
                label="Metrics Refresh Interval"
                value={refreshInterval}
                onChange={(e) => setRefreshInterval(e.target.value)}
                options={[
                  { value: '1', label: '1 second' },
                  { value: '5', label: '5 seconds' },
                  { value: '10', label: '10 seconds' },
                  { value: '30', label: '30 seconds' },
                  { value: '60', label: '1 minute' },
                ]}
              />
              <div className="pt-2">
                <p className="text-xs text-animus-text-secondary">
                  API Endpoint: {process.env.NEXT_PUBLIC_API_URL || 'Not configured'}
                </p>
              </div>
            </CardContent>
          </Card>

          {/* Appearance */}
          <Card variant="default">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Palette className="w-5 h-5" />
                Appearance
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <Select
                label="Theme"
                value={theme}
                onChange={(e) => setTheme(e.target.value)}
                options={[
                  { value: 'animus', label: 'Animus (Gold/Cyan)' },
                  { value: 'assassin', label: 'Assassin (Red/White)' },
                  { value: 'templar', label: 'Templar (Red/Gold)' },
                  { value: 'modern', label: 'Modern Day' },
                ]}
              />
              <div className="flex items-center justify-between py-2">
                <div>
                  <p className="text-sm text-white">Scan Lines Effect</p>
                  <p className="text-xs text-animus-text-secondary">
                    CRT-style overlay effect
                  </p>
                </div>
                <button
                  className="relative w-12 h-6 rounded-full bg-animus-gold/20 border border-animus-gold/30 transition-colors"
                  onClick={() => {}}
                >
                  <span className="absolute left-1 top-1 w-4 h-4 rounded-full bg-animus-gold transition-transform translate-x-6" />
                </button>
              </div>
              <div className="flex items-center justify-between py-2">
                <div>
                  <p className="text-sm text-white">Glitch Effects</p>
                  <p className="text-xs text-animus-text-secondary">
                    Text hover animations
                  </p>
                </div>
                <button
                  className="relative w-12 h-6 rounded-full bg-animus-gold/20 border border-animus-gold/30 transition-colors"
                  onClick={() => {}}
                >
                  <span className="absolute left-1 top-1 w-4 h-4 rounded-full bg-animus-gold transition-transform translate-x-6" />
                </button>
              </div>
            </CardContent>
          </Card>

          {/* Notifications */}
          <Card variant="default">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Bell className="w-5 h-5" />
                Notifications
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between py-2">
                <div>
                  <p className="text-sm text-white">Desktop Notifications</p>
                  <p className="text-xs text-animus-text-secondary">
                    Get notified of critical events
                  </p>
                </div>
                <button
                  className={`relative w-12 h-6 rounded-full border transition-colors ${
                    notifications
                      ? 'bg-animus-gold/20 border-animus-gold'
                      : 'bg-white/10 border-white/20'
                  }`}
                  onClick={() => setNotifications(!notifications)}
                >
                  <span
                    className={`absolute top-1 w-4 h-4 rounded-full transition-all ${
                      notifications
                        ? 'left-7 bg-animus-gold'
                        : 'left-1 bg-white/50'
                    }`}
                  />
                </button>
              </div>
              <div className="flex items-center justify-between py-2">
                <div>
                  <p className="text-sm text-white">Node Alerts</p>
                  <p className="text-xs text-animus-text-secondary">
                    Alert when nodes go offline
                  </p>
                </div>
                <button
                  className="relative w-12 h-6 rounded-full bg-animus-gold/20 border border-animus-gold transition-colors"
                  onClick={() => {}}
                >
                  <span className="absolute left-7 top-1 w-4 h-4 rounded-full bg-animus-gold transition-transform" />
                </button>
              </div>
              <div className="flex items-center justify-between py-2">
                <div>
                  <p className="text-sm text-white">Script Completion</p>
                  <p className="text-xs text-animus-text-secondary">
                    Notify when scripts finish
                  </p>
                </div>
                <button
                  className="relative w-12 h-6 rounded-full bg-animus-gold/20 border border-animus-gold transition-colors"
                  onClick={() => {}}
                >
                  <span className="absolute left-7 top-1 w-4 h-4 rounded-full bg-animus-gold transition-transform" />
                </button>
              </div>
            </CardContent>
          </Card>

          {/* Cluster Info */}
          <Card variant="glow">
            <CardHeader>
              <CardTitle>Cluster Information</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex justify-between py-2 border-b border-animus-gold/10">
                <span className="text-animus-text-secondary">Total Nodes</span>
                <span className="text-white font-mono">{nodes.length}</span>
              </div>
              <div className="flex justify-between py-2 border-b border-animus-gold/10">
                <span className="text-animus-text-secondary">Dashboard Version</span>
                <span className="text-white font-mono">v2.0.0</span>
              </div>
              <div className="flex justify-between py-2 border-b border-animus-gold/10">
                <span className="text-animus-text-secondary">API Status</span>
                <span className="text-animus-green font-mono">Connected</span>
              </div>
              <div className="flex justify-between py-2">
                <span className="text-animus-text-secondary">WebSocket</span>
                <span className="text-animus-green font-mono">Active</span>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Save Button */}
        <div className="mt-8 flex justify-end gap-4">
          <Button variant="ghost">
            <RefreshCw className="w-4 h-4 mr-2" />
            Reset to Defaults
          </Button>
          <Button variant="primary" onClick={handleSave}>
            <Save className="w-4 h-4 mr-2" />
            Save Changes
          </Button>
        </div>
      </main>
    </div>
  )
}
