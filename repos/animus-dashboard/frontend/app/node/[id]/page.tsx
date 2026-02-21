'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import {
  ArrowLeft,
  Cpu,
  HardDrive,
  MemoryStick,
  Clock,
  Server,
  Package,
  RefreshCw,
  ArrowUpCircle,
} from 'lucide-react'
import { Navigation } from '@/components/Navigation'
import { PodList } from '@/components/PodList'
import { LogStream } from '@/components/LogStream'
import { ScriptRunner } from '@/components/ScriptRunner'
import { MetricsChart, useMetricsHistory } from '@/components/MetricsChart'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { GlitchText } from '@/components/ui/GlitchText'
import { ProgressBar } from '@/components/ui/ProgressBar'
import { StatusIndicator, getStatusFromHealth } from '@/components/ui/StatusIndicator'
import { DNALoader } from '@/components/ui/DNALoader'
import { useNode, useNodes } from '@/hooks/useNodes'
import { api, Script } from '@/lib/api'

export default function NodeDetailPage() {
  const params = useParams()
  const router = useRouter()
  const nodeId = params.id as string

  const { nodes } = useNodes()
  const { node, pods, loading, error, refetch } = useNode(nodeId)
  const [scripts, setScripts] = useState<Script[]>([])

  const cpuHistory = useMetricsHistory(node?.cpu ?? 0)
  const memHistory = useMetricsHistory(node?.memory ?? 0)

  useEffect(() => {
    api.getScripts().then(setScripts).catch(console.error)
  }, [])

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <DNALoader size="lg" text={`Accessing memory sequence: ${nodeId}...`} />
      </div>
    )
  }

  if (error || !node) {
    return (
      <div className="min-h-screen">
        <Navigation nodes={nodes} />
        <main className="max-w-screen-2xl mx-auto px-4 py-8">
          <Card variant="highlight" className="max-w-md mx-auto">
            <CardContent className="text-center py-8">
              <p className="text-animus-red text-lg mb-4">
                Memory Sequence Not Found
              </p>
              <p className="text-animus-text-secondary mb-6">
                {error || `Node "${nodeId}" does not exist in the cluster.`}
              </p>
              <Button onClick={() => router.push('/')}>
                <ArrowLeft className="w-4 h-4 mr-2" />
                Return to Animus Core
              </Button>
            </CardContent>
          </Card>
        </main>
      </div>
    )
  }

  const status = getStatusFromHealth(
    node.status === 'Ready',
    node.cpu,
    node.memory
  )

  const nodeScripts = scripts.filter((s) => s.targetType === 'single')

  return (
    <div className="min-h-screen">
      <Navigation nodes={nodes} />

      <main className="max-w-screen-2xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-4">
            <Button variant="ghost" onClick={() => router.push('/')}>
              <ArrowLeft className="w-4 h-4" />
            </Button>
            <div>
              <div className="flex items-center gap-3">
                <GlitchText
                  text={`MEMORY SEQUENCE: ${node.name.toUpperCase()}`}
                  as="h1"
                  className="text-2xl font-bold text-animus-gold tracking-wider"
                />
                <StatusIndicator status={status} size="lg" />
              </div>
              <div className="flex items-center gap-4 mt-1 text-sm text-animus-text-secondary">
                <span className="flex items-center gap-1">
                  <Server className="w-3 h-3" />
                  {node.ip}
                </span>
                <span className="flex items-center gap-1">
                  <Clock className="w-3 h-3" />
                  Uptime: {node.uptime}
                </span>
                <span className="flex items-center gap-1">
                  <Package className="w-3 h-3" />
                  K3s {node.k3sVersion}
                </span>
              </div>
            </div>
          </div>
          <Button variant="ghost" onClick={refetch}>
            <RefreshCw className="w-4 h-4" />
          </Button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Column - Metrics and Version */}
          <div className="space-y-6">
            {/* System Metrics */}
            <Card variant="glow">
              <CardHeader>
                <CardTitle>System Metrics</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* CPU */}
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      <Cpu className="w-4 h-4 text-animus-cyan" />
                      <span className="text-sm text-animus-text-secondary">CPU</span>
                    </div>
                    <span className="text-lg font-mono text-white">{node.cpu}%</span>
                  </div>
                  <ProgressBar value={node.cpu} size="md" />
                </div>

                {/* Memory */}
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      <MemoryStick className="w-4 h-4 text-animus-gold" />
                      <span className="text-sm text-animus-text-secondary">Memory</span>
                    </div>
                    <span className="text-lg font-mono text-white">{node.memory}%</span>
                  </div>
                  <ProgressBar value={node.memory} size="md" />
                </div>

                {/* Disk */}
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      <HardDrive className="w-4 h-4 text-animus-green" />
                      <span className="text-sm text-animus-text-secondary">Disk</span>
                    </div>
                    <span className="text-lg font-mono text-white">{node.disk}%</span>
                  </div>
                  <ProgressBar value={node.disk} size="md" />
                </div>
              </CardContent>
            </Card>

            {/* Version Status */}
            <Card variant="default">
              <CardHeader>
                <CardTitle>Version Status</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-center justify-between py-2 border-b border-animus-gold/10">
                  <div>
                    <p className="text-sm text-white">K3s</p>
                    <p className="text-xs text-animus-text-secondary">{node.k3sVersion}</p>
                  </div>
                  {node.hasUpdate ? (
                    <span className="flex items-center gap-1 text-animus-cyan text-xs">
                      <ArrowUpCircle className="w-4 h-4" />
                      Update Available
                    </span>
                  ) : (
                    <span className="text-animus-green text-xs">✓ Latest</span>
                  )}
                </div>
                <div className="flex items-center justify-between py-2">
                  <div>
                    <p className="text-sm text-white">OS</p>
                    <p className="text-xs text-animus-text-secondary">{node.osVersion}</p>
                  </div>
                  <span className="text-animus-green text-xs">✓ Up to date</span>
                </div>
              </CardContent>
            </Card>

            {/* Mini Charts */}
            <div className="space-y-4">
              <MetricsChart
                title="CPU History"
                data={cpuHistory}
                color="cyan"
                type="area"
              />
              <MetricsChart
                title="Memory History"
                data={memHistory}
                color="gold"
                type="area"
              />
            </div>
          </div>

          {/* Right Column - Pods, Logs, Scripts */}
          <div className="lg:col-span-2 space-y-6">
            {/* Pod List */}
            <PodList pods={pods} />

            {/* Log Stream */}
            <LogStream node={node.name} />

            {/* Node Protocols */}
            <ScriptRunner
              nodes={[node]}
              scripts={nodeScripts}
              selectedNode={node.id}
            />
          </div>
        </div>
      </main>
    </div>
  )
}
