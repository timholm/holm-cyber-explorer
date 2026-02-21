'use client'

import { useEffect, useState } from 'react'
import { ScrollText, History, CheckCircle, XCircle, Clock, Play } from 'lucide-react'
import { Navigation } from '@/components/Navigation'
import { ScriptRunner } from '@/components/ScriptRunner'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { GlitchText } from '@/components/ui/GlitchText'
import { DNALoader } from '@/components/ui/DNALoader'
import { useNodes } from '@/hooks/useNodes'
import { api, Script, ScriptExecution } from '@/lib/api'
import { clsx } from 'clsx'

export default function ScriptsPage() {
  const { nodes, loading: nodesLoading } = useNodes()
  const [scripts, setScripts] = useState<Script[]>([])
  const [executions, setExecutions] = useState<ScriptExecution[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      api.getScripts(),
      api.getScriptExecutions(),
    ])
      .then(([scriptsData, executionsData]) => {
        setScripts(scriptsData)
        setExecutions(executionsData)
      })
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  const refreshExecutions = async () => {
    const data = await api.getScriptExecutions()
    setExecutions(data)
  }

  if (loading || nodesLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <DNALoader size="lg" text="Loading protocols..." />
      </div>
    )
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <CheckCircle className="w-4 h-4 text-animus-green" />
      case 'failed':
        return <XCircle className="w-4 h-4 text-animus-red" />
      case 'running':
        return <Clock className="w-4 h-4 text-animus-gold animate-pulse" />
      default:
        return null
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'text-animus-green'
      case 'failed':
        return 'text-animus-red'
      case 'running':
        return 'text-animus-gold'
      default:
        return 'text-animus-text-secondary'
    }
  }

  return (
    <div className="min-h-screen">
      <Navigation nodes={nodes} />

      <main className="max-w-screen-2xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-center gap-4 mb-8">
          <div className="p-3 rounded-lg bg-animus-gold/10">
            <ScrollText className="w-8 h-8 text-animus-gold" />
          </div>
          <div>
            <GlitchText
              text="PROTOCOLS"
              as="h1"
              className="text-3xl font-bold text-animus-gold tracking-widest"
            />
            <p className="text-animus-text-secondary mt-1">
              Execute and manage Ansible playbooks across the cluster
            </p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Script Runner */}
          <div className="lg:col-span-2">
            <ScriptRunner nodes={nodes} scripts={scripts} />
          </div>

          {/* Execution History */}
          <div>
            <Card variant="default">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <CardTitle className="flex items-center gap-2">
                    <History className="w-5 h-5" />
                    Execution History
                  </CardTitle>
                  <Button variant="ghost" size="sm" onClick={refreshExecutions}>
                    Refresh
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-3 max-h-96 overflow-y-auto">
                  {executions.length === 0 ? (
                    <p className="text-animus-text-secondary text-sm text-center py-8">
                      No executions yet
                    </p>
                  ) : (
                    executions.map((execution) => {
                      const script = scripts.find((s) => s.id === execution.scriptId)
                      return (
                        <div
                          key={execution.id}
                          className="p-3 rounded-lg bg-animus-black/50 border border-animus-gold/10 hover:border-animus-gold/30 transition-colors"
                        >
                          <div className="flex items-start justify-between">
                            <div>
                              <p className="text-sm text-white font-medium">
                                {script?.name || execution.scriptId}
                              </p>
                              <p className="text-xs text-animus-text-secondary mt-0.5">
                                {execution.targetNodes.join(', ')}
                              </p>
                            </div>
                            <div className="flex items-center gap-2">
                              {getStatusIcon(execution.status)}
                              <span
                                className={clsx(
                                  'text-xs font-medium',
                                  getStatusColor(execution.status)
                                )}
                              >
                                {execution.status}
                              </span>
                            </div>
                          </div>
                          <div className="mt-2 flex items-center justify-between text-xs text-animus-text-secondary">
                            <span>
                              {new Date(execution.startedAt).toLocaleString()}
                            </span>
                            {execution.completedAt && (
                              <span>
                                Duration:{' '}
                                {Math.round(
                                  (new Date(execution.completedAt).getTime() -
                                    new Date(execution.startedAt).getTime()) /
                                    1000
                                )}
                                s
                              </span>
                            )}
                          </div>
                        </div>
                      )
                    })
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Quick Stats */}
            <Card variant="glow" className="mt-6">
              <CardContent>
                <h3 className="text-sm font-medium text-animus-text-secondary mb-4">
                  Quick Stats
                </h3>
                <div className="grid grid-cols-3 gap-4 text-center">
                  <div>
                    <p className="text-2xl font-bold text-animus-green">
                      {executions.filter((e) => e.status === 'completed').length}
                    </p>
                    <p className="text-xs text-animus-text-secondary">Completed</p>
                  </div>
                  <div>
                    <p className="text-2xl font-bold text-animus-red">
                      {executions.filter((e) => e.status === 'failed').length}
                    </p>
                    <p className="text-xs text-animus-text-secondary">Failed</p>
                  </div>
                  <div>
                    <p className="text-2xl font-bold text-animus-gold">
                      {executions.filter((e) => e.status === 'running').length}
                    </p>
                    <p className="text-xs text-animus-text-secondary">Running</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </main>
    </div>
  )
}
