'use client'

import { useState } from 'react'
import { FileText, RotateCcw, AlertCircle, CheckCircle, Clock } from 'lucide-react'
import { Card, CardHeader, CardTitle } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Pod } from '@/lib/api'

interface PodListProps {
  pods: Pod[]
  onViewLogs?: (podName: string, namespace: string) => void
}

function getPodStatusIcon(status: string) {
  switch (status.toLowerCase()) {
    case 'running':
      return <CheckCircle className="w-4 h-4 text-animus-green" />
    case 'pending':
      return <Clock className="w-4 h-4 text-animus-gold animate-pulse" />
    case 'failed':
    case 'error':
      return <AlertCircle className="w-4 h-4 text-animus-red" />
    default:
      return <Clock className="w-4 h-4 text-animus-text-secondary" />
  }
}

function getStatusColor(status: string) {
  switch (status.toLowerCase()) {
    case 'running':
      return 'text-animus-green'
    case 'pending':
      return 'text-animus-gold'
    case 'failed':
    case 'error':
      return 'text-animus-red'
    default:
      return 'text-animus-text-secondary'
  }
}

export function PodList({ pods, onViewLogs }: PodListProps) {
  const [filter, setFilter] = useState('')

  const filteredPods = pods.filter(
    (pod) =>
      pod.name.toLowerCase().includes(filter.toLowerCase()) ||
      pod.namespace.toLowerCase().includes(filter.toLowerCase())
  )

  return (
    <Card variant="default" padding="none">
      <CardHeader className="px-4 pt-4">
        <div className="flex items-center justify-between">
          <CardTitle>Running Sequences (Pods)</CardTitle>
          <span className="text-sm text-animus-text-secondary">
            {pods.length} total
          </span>
        </div>
        <div className="mt-3">
          <Input
            placeholder="Filter sequences..."
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
          />
        </div>
      </CardHeader>

      <div className="max-h-96 overflow-y-auto">
        <table className="w-full">
          <thead className="bg-animus-black/50 sticky top-0">
            <tr className="text-left text-xs text-animus-text-secondary uppercase tracking-wider">
              <th className="px-4 py-3">Name</th>
              <th className="px-4 py-3">Namespace</th>
              <th className="px-4 py-3">Status</th>
              <th className="px-4 py-3">Restarts</th>
              <th className="px-4 py-3">Age</th>
              <th className="px-4 py-3">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-animus-gold/10">
            {filteredPods.map((pod) => (
              <tr
                key={`${pod.namespace}/${pod.name}`}
                className="hover:bg-animus-gold/5 transition-colors"
              >
                <td className="px-4 py-3">
                  <div className="flex items-center gap-2">
                    {getPodStatusIcon(pod.status)}
                    <span className="text-sm text-white font-mono truncate max-w-[200px]">
                      {pod.name}
                    </span>
                  </div>
                </td>
                <td className="px-4 py-3">
                  <span className="text-sm text-animus-cyan">{pod.namespace}</span>
                </td>
                <td className="px-4 py-3">
                  <span className={`text-sm font-medium ${getStatusColor(pod.status)}`}>
                    {pod.status}
                  </span>
                </td>
                <td className="px-4 py-3">
                  <span
                    className={`text-sm font-mono ${
                      pod.restarts > 0 ? 'text-animus-gold' : 'text-animus-text-secondary'
                    }`}
                  >
                    {pod.restarts}
                    {pod.restarts > 5 && (
                      <RotateCcw className="inline w-3 h-3 ml-1 text-animus-red" />
                    )}
                  </span>
                </td>
                <td className="px-4 py-3">
                  <span className="text-sm text-animus-text-secondary">{pod.age}</span>
                </td>
                <td className="px-4 py-3">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => onViewLogs?.(pod.name, pod.namespace)}
                  >
                    <FileText className="w-4 h-4" />
                    <span className="ml-1">Logs</span>
                  </Button>
                </td>
              </tr>
            ))}
            {filteredPods.length === 0 && (
              <tr>
                <td colSpan={6} className="px-4 py-8 text-center text-animus-text-secondary">
                  {filter ? 'No sequences match filter' : 'No sequences running'}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </Card>
  )
}
