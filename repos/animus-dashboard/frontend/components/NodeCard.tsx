'use client'

import Link from 'next/link'
import { Cpu, HardDrive, MemoryStick, Box } from 'lucide-react'
import { Card } from '@/components/ui/Card'
import { StatusIndicator, getStatusFromHealth } from '@/components/ui/StatusIndicator'
import { ProgressBar } from '@/components/ui/ProgressBar'
import { Node } from '@/lib/api'

interface NodeCardProps {
  node: Node
}

export function NodeCard({ node }: NodeCardProps) {
  const status = getStatusFromHealth(
    node.status === 'Ready',
    node.cpu,
    node.memory
  )

  return (
    <Link href={`/node/${node.id}`}>
      <Card
        variant="glow"
        className="cursor-pointer hover:translate-y-[-2px] transition-all duration-300"
      >
        <div className="flex items-start justify-between mb-4">
          <div>
            <h3 className="text-lg font-semibold text-white font-mono">
              {node.name}
            </h3>
            <p className="text-xs text-animus-text-secondary mt-0.5">
              {node.ip}
            </p>
          </div>
          <StatusIndicator status={status} size="md" />
        </div>

        <div className="space-y-3">
          {/* CPU */}
          <div className="flex items-center gap-3">
            <Cpu className="w-4 h-4 text-animus-cyan" />
            <div className="flex-1">
              <div className="flex justify-between text-xs mb-1">
                <span className="text-animus-text-secondary">CPU</span>
                <span className="text-white font-mono">{node.cpu}%</span>
              </div>
              <ProgressBar value={node.cpu} size="sm" />
            </div>
          </div>

          {/* Memory */}
          <div className="flex items-center gap-3">
            <MemoryStick className="w-4 h-4 text-animus-gold" />
            <div className="flex-1">
              <div className="flex justify-between text-xs mb-1">
                <span className="text-animus-text-secondary">Memory</span>
                <span className="text-white font-mono">{node.memory}%</span>
              </div>
              <ProgressBar value={node.memory} size="sm" />
            </div>
          </div>

          {/* Disk */}
          <div className="flex items-center gap-3">
            <HardDrive className="w-4 h-4 text-animus-green" />
            <div className="flex-1">
              <div className="flex justify-between text-xs mb-1">
                <span className="text-animus-text-secondary">Disk</span>
                <span className="text-white font-mono">{node.disk}%</span>
              </div>
              <ProgressBar value={node.disk} size="sm" />
            </div>
          </div>
        </div>

        {/* Pod count */}
        <div className="mt-4 pt-3 border-t border-animus-gold/10 flex items-center justify-between">
          <div className="flex items-center gap-2 text-animus-text-secondary">
            <Box className="w-4 h-4" />
            <span className="text-sm">{node.podCount} sequences</span>
          </div>
          {node.hasUpdate && (
            <span className="text-xs text-animus-cyan animate-pulse">
              Update available
            </span>
          )}
        </div>
      </Card>
    </Link>
  )
}
