'use client'

import { useEffect, useRef, useState, useCallback } from 'react'

interface UseWebSocketOptions {
  onOpen?: () => void
  onClose?: () => void
  onError?: (error: Event) => void
  reconnectInterval?: number
  reconnectAttempts?: number
}

interface UseWebSocketReturn<T> {
  data: T | null
  isConnected: boolean
  error: string | null
  send: (message: string | object) => void
  reconnect: () => void
}

export function useWebSocket<T>(
  url: string,
  options: UseWebSocketOptions = {}
): UseWebSocketReturn<T> {
  const {
    onOpen,
    onClose,
    onError,
    reconnectInterval = 3000,
    reconnectAttempts = 5,
  } = options

  const [data, setData] = useState<T | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const wsRef = useRef<WebSocket | null>(null)
  const reconnectCountRef = useRef(0)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)

  const connect = useCallback(() => {
    if (typeof window === 'undefined') return

    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsHost = process.env.NEXT_PUBLIC_WS_URL || `${wsProtocol}//${window.location.host}`
    const fullUrl = `${wsHost}${url}`

    try {
      wsRef.current = new WebSocket(fullUrl)

      wsRef.current.onopen = () => {
        setIsConnected(true)
        setError(null)
        reconnectCountRef.current = 0
        onOpen?.()
      }

      wsRef.current.onclose = () => {
        setIsConnected(false)
        onClose?.()

        // Attempt reconnection
        if (reconnectCountRef.current < reconnectAttempts) {
          reconnectTimeoutRef.current = setTimeout(() => {
            reconnectCountRef.current++
            connect()
          }, reconnectInterval)
        }
      }

      wsRef.current.onerror = (event) => {
        setError('WebSocket connection error')
        onError?.(event)
      }

      wsRef.current.onmessage = (event) => {
        try {
          const parsed = JSON.parse(event.data)
          setData(parsed)
        } catch {
          setData(event.data as T)
        }
      }
    } catch (err) {
      setError(`Failed to connect: ${err}`)
    }
  }, [url, onOpen, onClose, onError, reconnectInterval, reconnectAttempts])

  const send = useCallback((message: string | object) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      const msg = typeof message === 'object' ? JSON.stringify(message) : message
      wsRef.current.send(msg)
    }
  }, [])

  const reconnect = useCallback(() => {
    if (wsRef.current) {
      wsRef.current.close()
    }
    reconnectCountRef.current = 0
    connect()
  }, [connect])

  useEffect(() => {
    connect()

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [connect])

  return { data, isConnected, error, send, reconnect }
}

// Specialized hook for log streaming
export interface LogStreamData {
  entries: Array<{
    timestamp: string
    source: string
    level: 'info' | 'warn' | 'error'
    message: string
    namespace?: string
  }>
}

export function useLogStream(node?: string, namespace?: string) {
  const params = new URLSearchParams()
  if (node) params.set('node', node)
  if (namespace) params.set('namespace', namespace)

  const queryString = params.toString()
  const url = `/ws/logs${queryString ? `?${queryString}` : ''}`

  return useWebSocket<LogStreamData>(url)
}

// Specialized hook for real-time metrics
export interface MetricsData {
  nodes: Array<{
    name: string
    cpu: number
    memory: number
    disk: number
    podCount: number
    status: string
  }>
  cluster: {
    totalCpu: number
    totalMemory: number
    healthyNodes: number
    totalNodes: number
  }
}

export function useMetricsStream() {
  return useWebSocket<MetricsData>('/ws/metrics')
}
