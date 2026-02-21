const API_BASE = process.env.NEXT_PUBLIC_API_URL || ''

export interface Node {
  id: string
  name: string
  status: 'Ready' | 'NotReady' | 'Unknown'
  cpu: number
  memory: number
  disk: number
  podCount: number
  k3sVersion: string
  osVersion: string
  uptime: string
  ip: string
  roles: string[]
  hasUpdate: boolean
}

export interface Pod {
  name: string
  namespace: string
  status: string
  restarts: number
  age: string
  cpu: string
  memory: string
  nodeName: string
}

export interface LogEntry {
  timestamp: string
  source: string
  level: 'info' | 'warn' | 'error'
  message: string
  namespace?: string
}

export interface Script {
  id: string
  name: string
  description: string
  category: string
  targetType: 'single' | 'all'
}

export interface ScriptExecution {
  id: string
  scriptId: string
  status: 'running' | 'completed' | 'failed'
  startedAt: string
  completedAt?: string
  output: string[]
  targetNodes: string[]
}

export interface ClusterMetrics {
  totalNodes: number
  healthyNodes: number
  totalPods: number
  runningPods: number
  cpuUsage: number
  memoryUsage: number
}

async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  })

  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText}`)
  }

  return response.json()
}

export const api = {
  // Nodes
  async getNodes(): Promise<Node[]> {
    return fetchAPI<Node[]>('/api/nodes')
  },

  async getNode(id: string): Promise<Node> {
    return fetchAPI<Node>(`/api/nodes/${id}`)
  },

  async getNodePods(id: string): Promise<Pod[]> {
    return fetchAPI<Pod[]>(`/api/nodes/${id}/pods`)
  },

  // Cluster
  async getClusterMetrics(): Promise<ClusterMetrics> {
    return fetchAPI<ClusterMetrics>('/api/cluster/metrics')
  },

  // Logs
  async getLogs(params: {
    node?: string
    namespace?: string
    query?: string
    limit?: number
    start?: string
    end?: string
  }): Promise<LogEntry[]> {
    const searchParams = new URLSearchParams()
    if (params.node) searchParams.set('node', params.node)
    if (params.namespace) searchParams.set('namespace', params.namespace)
    if (params.query) searchParams.set('query', params.query)
    if (params.limit) searchParams.set('limit', params.limit.toString())
    if (params.start) searchParams.set('start', params.start)
    if (params.end) searchParams.set('end', params.end)

    return fetchAPI<LogEntry[]>(`/api/logs?${searchParams.toString()}`)
  },

  // Scripts
  async getScripts(): Promise<Script[]> {
    return fetchAPI<Script[]>('/api/scripts')
  },

  async runScript(scriptId: string, targetNodes: string[]): Promise<ScriptExecution> {
    return fetchAPI<ScriptExecution>('/api/scripts/run', {
      method: 'POST',
      body: JSON.stringify({ scriptId, targetNodes }),
    })
  },

  async getScriptExecution(executionId: string): Promise<ScriptExecution> {
    return fetchAPI<ScriptExecution>(`/api/scripts/executions/${executionId}`)
  },

  async getScriptExecutions(): Promise<ScriptExecution[]> {
    return fetchAPI<ScriptExecution[]>('/api/scripts/executions')
  },

  // Versions
  async checkUpdates(): Promise<{ node: string; updates: { name: string; current: string; available: string }[] }[]> {
    return fetchAPI('/api/versions/check')
  },
}

// WebSocket connection helper
export function createWebSocket(path: string): WebSocket {
  const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsHost = process.env.NEXT_PUBLIC_WS_URL || `${wsProtocol}//${window.location.host}`
  return new WebSocket(`${wsHost}${path}`)
}
