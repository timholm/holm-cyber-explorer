'use client'

import { useState, useEffect, useCallback } from 'react'
import { api, Node, Pod, ClusterMetrics } from '@/lib/api'
import { useMetricsStream } from './useWebSocket'

export function useNodes() {
  const [nodes, setNodes] = useState<Node[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const { data: metricsData, isConnected } = useMetricsStream()

  const fetchNodes = useCallback(async () => {
    try {
      setLoading(true)
      const data = await api.getNodes()
      setNodes(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch nodes')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchNodes()
  }, [fetchNodes])

  // Update nodes with real-time metrics when available
  useEffect(() => {
    if (metricsData?.nodes) {
      setNodes((prevNodes) =>
        prevNodes.map((node) => {
          const metrics = metricsData.nodes.find((m) => m.name === node.name)
          if (metrics) {
            return {
              ...node,
              cpu: metrics.cpu,
              memory: metrics.memory,
              disk: metrics.disk,
              podCount: metrics.podCount,
              status: metrics.status as Node['status'],
            }
          }
          return node
        })
      )
    }
  }, [metricsData])

  return {
    nodes,
    loading,
    error,
    refetch: fetchNodes,
    isRealtime: isConnected,
  }
}

export function useNode(nodeId: string) {
  const [node, setNode] = useState<Node | null>(null)
  const [pods, setPods] = useState<Pod[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const { data: metricsData } = useMetricsStream()

  const fetchNodeData = useCallback(async () => {
    try {
      setLoading(true)
      const [nodeData, podsData] = await Promise.all([
        api.getNode(nodeId),
        api.getNodePods(nodeId),
      ])
      setNode(nodeData)
      setPods(podsData)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch node data')
    } finally {
      setLoading(false)
    }
  }, [nodeId])

  useEffect(() => {
    fetchNodeData()
  }, [fetchNodeData])

  // Update node with real-time metrics
  useEffect(() => {
    if (metricsData?.nodes && node) {
      const metrics = metricsData.nodes.find((m) => m.name === node.name)
      if (metrics) {
        setNode((prev) =>
          prev
            ? {
                ...prev,
                cpu: metrics.cpu,
                memory: metrics.memory,
                disk: metrics.disk,
                podCount: metrics.podCount,
              }
            : prev
        )
      }
    }
  }, [metricsData, node])

  return {
    node,
    pods,
    loading,
    error,
    refetch: fetchNodeData,
  }
}

export function useClusterMetrics() {
  const [metrics, setMetrics] = useState<ClusterMetrics | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const { data: metricsData } = useMetricsStream()

  const fetchMetrics = useCallback(async () => {
    try {
      setLoading(true)
      const data = await api.getClusterMetrics()
      setMetrics(data)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch metrics')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchMetrics()
  }, [fetchMetrics])

  // Update with real-time data
  useEffect(() => {
    if (metricsData?.cluster) {
      setMetrics((prev) =>
        prev
          ? {
              ...prev,
              cpuUsage: metricsData.cluster.totalCpu,
              memoryUsage: metricsData.cluster.totalMemory,
              healthyNodes: metricsData.cluster.healthyNodes,
              totalNodes: metricsData.cluster.totalNodes,
            }
          : prev
      )
    }
  }, [metricsData])

  return {
    metrics,
    loading,
    error,
    refetch: fetchMetrics,
  }
}
