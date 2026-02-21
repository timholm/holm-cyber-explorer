'use client'

import { useEffect, useState } from 'react'
import { Activity, Cpu, HardDrive, MemoryStick, RefreshCw } from 'lucide-react'
import { Navigation } from '@/components/Navigation'
import { NodeCard } from '@/components/NodeCard'
import { ScriptRunner } from '@/components/ScriptRunner'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { GlitchText } from '@/components/ui/GlitchText'
import { ProgressBar } from '@/components/ui/ProgressBar'
import { DNALoader } from '@/components/ui/DNALoader'
import { useNodes, useClusterMetrics } from '@/hooks/useNodes'
import { api, Script } from '@/lib/api'
import { calculateSyncPercentage } from '@/lib/theme'

export default function HomePage() {
  const { nodes, loading, error, refetch, isRealtime } = useNodes()
  const { metrics } = useClusterMetrics()
  const [scripts, setScripts] = useState<Script[]>([])

  useEffect(() => {
    api.getScripts().then(setScripts).catch(console.error)
  }, [])

  const healthyNodes = nodes.filter((n) => n.status === 'Ready').length
  const syncPercentage = calculateSyncPercentage(healthyNodes, nodes.length)

  const clusterScripts = scripts.filter((s) => s.targetType === 'all')

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <DNALoader size="lg" text="Synchronizing with cluster..." />
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Card variant="highlight" className="max-w-md">
          <CardContent className="text-center py-8">
            <p className="text-animus-red text-lg mb-4">Desynchronization Detected</p>
            <p className="text-animus-text-secondary mb-6">{error}</p>
            <Button onClick={refetch}>
              <RefreshCw className="w-4 h-4 mr-2" />
              Retry Connection
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen">
      <Navigation nodes={nodes} />

      <main className="max-w-screen-2xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div>
            <GlitchText
              text="ANIMUS CORE"
              as="h1"
              className="text-3xl font-bold text-animus-gold tracking-widest"
            />
            <p className="text-animus-text-secondary mt-1">
              Cluster Synchronization Overview
            </p>
          </div>
          <div className="flex items-center gap-4">
            {isRealtime && (
              <div className="flex items-center gap-2 text-animus-green text-sm">
                <Activity className="w-4 h-4 animate-pulse" />
                Real-time
              </div>
            )}
            <Button variant="ghost" onClick={refetch}>
              <RefreshCw className="w-4 h-4" />
            </Button>
          </div>
        </div>

        {/* Cluster Overview */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          {/* Sync Status */}
          <Card variant="highlight" className="md:col-span-2">
            <CardContent className="flex items-center gap-6">
              <div className="relative">
                <svg className="w-24 h-24 transform -rotate-90">
                  <circle
                    cx="48"
                    cy="48"
                    r="40"
                    stroke="currentColor"
                    strokeWidth="8"
                    fill="none"
                    className="text-animus-gold/20"
                  />
                  <circle
                    cx="48"
                    cy="48"
                    r="40"
                    stroke="currentColor"
                    strokeWidth="8"
                    fill="none"
                    strokeDasharray={`${(syncPercentage / 100) * 251.2} 251.2`}
                    className="text-animus-gold transition-all duration-1000"
                  />
                </svg>
                <div className="absolute inset-0 flex items-center justify-center">
                  <span className="text-2xl font-bold text-white">{syncPercentage}%</span>
                </div>
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white">
                  Cluster Synchronization
                </h3>
                <p className="text-animus-text-secondary mt-1">
                  {healthyNodes} of {nodes.length} memory cores operational
                </p>
                <div className="mt-3 flex items-center gap-4 text-sm">
                  <span className="text-animus-green">● {healthyNodes} Synced</span>
                  {nodes.length - healthyNodes > 0 && (
                    <span className="text-animus-red">
                      ● {nodes.length - healthyNodes} Desynced
                    </span>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* CPU Usage */}
          <Card variant="glow">
            <CardContent>
              <div className="flex items-center gap-3 mb-3">
                <div className="p-2 rounded-lg bg-animus-cyan/10">
                  <Cpu className="w-5 h-5 text-animus-cyan" />
                </div>
                <div>
                  <p className="text-xs text-animus-text-secondary">Total CPU</p>
                  <p className="text-xl font-bold text-white">
                    {metrics?.cpuUsage ?? '--'}%
                  </p>
                </div>
              </div>
              <ProgressBar value={metrics?.cpuUsage ?? 0} size="md" />
            </CardContent>
          </Card>

          {/* Memory Usage */}
          <Card variant="glow">
            <CardContent>
              <div className="flex items-center gap-3 mb-3">
                <div className="p-2 rounded-lg bg-animus-gold/10">
                  <MemoryStick className="w-5 h-5 text-animus-gold" />
                </div>
                <div>
                  <p className="text-xs text-animus-text-secondary">Total Memory</p>
                  <p className="text-xl font-bold text-white">
                    {metrics?.memoryUsage ?? '--'}%
                  </p>
                </div>
              </div>
              <ProgressBar value={metrics?.memoryUsage ?? 0} size="md" />
            </CardContent>
          </Card>
        </div>

        {/* Node Grid */}
        <div className="mb-8">
          <h2 className="text-xl font-semibold text-white mb-4 flex items-center gap-2">
            <HardDrive className="w-5 h-5 text-animus-gold" />
            Memory Cores
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {nodes.map((node) => (
              <NodeCard key={node.id} node={node} />
            ))}
          </div>
        </div>

        {/* Global Protocols */}
        <div>
          <h2 className="text-xl font-semibold text-white mb-4">
            Global Protocols
          </h2>
          <ScriptRunner
            nodes={nodes}
            scripts={clusterScripts}
          />
        </div>
      </main>
    </div>
  )
}
