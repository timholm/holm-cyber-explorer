'use client'

import { useState, useEffect, useCallback, useRef } from 'react'
import { api, LogEntry } from '@/lib/api'
import { useLogStream } from './useWebSocket'

interface UseLogsOptions {
  node?: string
  namespace?: string
  query?: string
  limit?: number
  streaming?: boolean
}

export function useLogs(options: UseLogsOptions = {}) {
  const { node, namespace, query, limit = 100, streaming = true } = options

  const [logs, setLogs] = useState<LogEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const { data: streamData, isConnected } = useLogStream(
    streaming ? node : undefined,
    streaming ? namespace : undefined
  )

  const logsRef = useRef<LogEntry[]>([])

  // Fetch initial logs
  const fetchLogs = useCallback(async () => {
    try {
      setLoading(true)
      const data = await api.getLogs({ node, namespace, query, limit })
      setLogs(data)
      logsRef.current = data
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch logs')
    } finally {
      setLoading(false)
    }
  }, [node, namespace, query, limit])

  useEffect(() => {
    fetchLogs()
  }, [fetchLogs])

  // Append streaming logs
  useEffect(() => {
    if (streamData?.entries && streamData.entries.length > 0) {
      setLogs((prevLogs) => {
        const newLogs = [...streamData.entries, ...prevLogs]
        // Keep only the most recent logs
        const trimmedLogs = newLogs.slice(0, limit * 2)
        logsRef.current = trimmedLogs
        return trimmedLogs
      })
    }
  }, [streamData, limit])

  // Filter logs by query
  const filteredLogs = query
    ? logs.filter(
        (log) =>
          log.message.toLowerCase().includes(query.toLowerCase()) ||
          log.source.toLowerCase().includes(query.toLowerCase())
      )
    : logs

  const clearLogs = useCallback(() => {
    setLogs([])
    logsRef.current = []
  }, [])

  return {
    logs: filteredLogs,
    allLogs: logs,
    loading,
    error,
    isStreaming: isConnected && streaming,
    refetch: fetchLogs,
    clearLogs,
  }
}

export function useLogFilters() {
  const [namespace, setNamespace] = useState<string>('')
  const [level, setLevel] = useState<string>('')
  const [searchQuery, setSearchQuery] = useState('')

  const clearFilters = useCallback(() => {
    setNamespace('')
    setLevel('')
    setSearchQuery('')
  }, [])

  return {
    namespace,
    setNamespace,
    level,
    setLevel,
    searchQuery,
    setSearchQuery,
    clearFilters,
    hasFilters: Boolean(namespace || level || searchQuery),
  }
}
