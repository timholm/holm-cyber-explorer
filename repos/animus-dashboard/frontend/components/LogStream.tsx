'use client'

import { useEffect, useRef, useState } from 'react'
import { Play, Pause, Trash2, Download, Search, Filter } from 'lucide-react'
import { Card, CardHeader, CardTitle } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Select } from '@/components/ui/Select'
import { useLogs, useLogFilters } from '@/hooks/useLogs'
import { LogEntry } from '@/lib/api'
import { clsx } from 'clsx'

interface LogStreamProps {
  node?: string
  defaultNamespace?: string
}

function LogLine({ entry }: { entry: LogEntry }) {
  const levelColors = {
    info: 'text-animus-cyan',
    warn: 'text-animus-gold',
    error: 'text-animus-red',
  }

  return (
    <div className="log-line flex items-start gap-3 py-1.5 px-3 hover:bg-animus-gold/5 font-mono text-xs">
      <span className="log-timestamp text-animus-cyan whitespace-nowrap">
        {new Date(entry.timestamp).toLocaleTimeString()}
      </span>
      <span className="log-source text-animus-gold min-w-[100px] truncate">
        {entry.source}
      </span>
      <span className={clsx('w-12', levelColors[entry.level])}>
        [{entry.level.toUpperCase()}]
      </span>
      <span className="text-white flex-1 break-all">{entry.message}</span>
    </div>
  )
}

export function LogStream({ node, defaultNamespace }: LogStreamProps) {
  const [isStreaming, setIsStreaming] = useState(true)
  const [autoScroll, setAutoScroll] = useState(true)
  const logContainerRef = useRef<HTMLDivElement>(null)

  const {
    namespace,
    setNamespace,
    level,
    setLevel,
    searchQuery,
    setSearchQuery,
    clearFilters,
    hasFilters,
  } = useLogFilters()

  const { logs, loading, isStreaming: connected, clearLogs } = useLogs({
    node,
    namespace: namespace || defaultNamespace,
    query: searchQuery,
    streaming: isStreaming,
  })

  // Filter by level
  const filteredLogs = level
    ? logs.filter((log) => log.level === level)
    : logs

  // Auto-scroll to bottom
  useEffect(() => {
    if (autoScroll && logContainerRef.current) {
      logContainerRef.current.scrollTop = logContainerRef.current.scrollHeight
    }
  }, [filteredLogs, autoScroll])

  const handleScroll = () => {
    if (!logContainerRef.current) return
    const { scrollTop, scrollHeight, clientHeight } = logContainerRef.current
    const isAtBottom = scrollHeight - scrollTop - clientHeight < 50
    setAutoScroll(isAtBottom)
  }

  const handleExport = () => {
    const content = filteredLogs
      .map(
        (log) =>
          `${log.timestamp} [${log.level.toUpperCase()}] ${log.source}: ${log.message}`
      )
      .join('\n')

    const blob = new Blob([content], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `animus-logs-${node || 'all'}-${new Date().toISOString()}.txt`
    a.click()
    URL.revokeObjectURL(url)
  }

  const namespaceOptions = [
    { value: '', label: 'All Namespaces' },
    { value: 'default', label: 'default' },
    { value: 'kube-system', label: 'kube-system' },
    { value: 'monitoring', label: 'monitoring' },
    { value: 'media', label: 'media' },
  ]

  const levelOptions = [
    { value: '', label: 'All Levels' },
    { value: 'info', label: 'Info' },
    { value: 'warn', label: 'Warning' },
    { value: 'error', label: 'Error' },
  ]

  return (
    <Card variant="default" padding="none" className="flex flex-col h-[500px]">
      <CardHeader className="px-4 pt-4 flex-shrink-0">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2">
            Memory Stream (Logs)
            {connected && isStreaming && (
              <span className="w-2 h-2 bg-animus-green rounded-full animate-pulse" />
            )}
          </CardTitle>
          <div className="flex items-center gap-2">
            <Button
              variant={isStreaming ? 'secondary' : 'ghost'}
              size="sm"
              onClick={() => setIsStreaming(!isStreaming)}
            >
              {isStreaming ? (
                <Pause className="w-4 h-4" />
              ) : (
                <Play className="w-4 h-4" />
              )}
            </Button>
            <Button variant="ghost" size="sm" onClick={clearLogs}>
              <Trash2 className="w-4 h-4" />
            </Button>
            <Button variant="ghost" size="sm" onClick={handleExport}>
              <Download className="w-4 h-4" />
            </Button>
          </div>
        </div>

        {/* Filters */}
        <div className="mt-3 flex items-center gap-3">
          <div className="flex-1">
            <Input
              placeholder="Search logs..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              icon={<Search className="w-4 h-4" />}
            />
          </div>
          <div className="w-40">
            <Select
              options={namespaceOptions}
              value={namespace}
              onChange={(e) => setNamespace(e.target.value)}
            />
          </div>
          <div className="w-32">
            <Select
              options={levelOptions}
              value={level}
              onChange={(e) => setLevel(e.target.value)}
            />
          </div>
          {hasFilters && (
            <Button variant="ghost" size="sm" onClick={clearFilters}>
              <Filter className="w-4 h-4" />
              Clear
            </Button>
          )}
        </div>
      </CardHeader>

      <div
        ref={logContainerRef}
        onScroll={handleScroll}
        className="flex-1 overflow-y-auto bg-animus-black/50 border-t border-animus-gold/10"
      >
        {loading && filteredLogs.length === 0 ? (
          <div className="flex items-center justify-center h-full text-animus-text-secondary">
            Loading memory stream...
          </div>
        ) : filteredLogs.length === 0 ? (
          <div className="flex items-center justify-center h-full text-animus-text-secondary">
            No log entries found
          </div>
        ) : (
          <div className="py-2">
            {filteredLogs.map((entry, index) => (
              <LogLine key={`${entry.timestamp}-${index}`} entry={entry} />
            ))}
          </div>
        )}
      </div>

      {!autoScroll && (
        <button
          onClick={() => {
            setAutoScroll(true)
            if (logContainerRef.current) {
              logContainerRef.current.scrollTop = logContainerRef.current.scrollHeight
            }
          }}
          className="absolute bottom-4 right-4 bg-animus-gold text-animus-black px-3 py-1 rounded text-xs font-medium hover:bg-animus-cyan transition-colors"
        >
          ↓ New logs
        </button>
      )}
    </Card>
  )
}
