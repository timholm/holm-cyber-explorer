export const theme = {
  colors: {
    primary: {
      black: '#0a0a0f',
      navy: '#12121a',
    },
    accent: {
      gold: '#c9a227',
      cyan: '#00d4ff',
    },
    status: {
      synchronized: '#00ff88',
      warning: '#c9a227',
      desynchronized: '#ff3366',
    },
    text: {
      primary: '#ffffff',
      secondary: '#8888aa',
    },
  },
  terminology: {
    nodes: 'Memory Cores',
    pods: 'Sequences',
    logs: 'Memory Stream',
    scripts: 'Protocols',
    health: 'Synchronization',
    errors: 'Desynchronization Events',
    dashboard: 'ANIMUS v2.0',
  },
} as const

export type NodeStatus = 'synchronized' | 'warning' | 'desynchronized'

export function getStatusColor(status: NodeStatus): string {
  switch (status) {
    case 'synchronized':
      return theme.colors.status.synchronized
    case 'warning':
      return theme.colors.status.warning
    case 'desynchronized':
      return theme.colors.status.desynchronized
  }
}

export function getStatusLabel(status: NodeStatus): string {
  switch (status) {
    case 'synchronized':
      return 'SYNC'
    case 'warning':
      return 'WARN'
    case 'desynchronized':
      return 'DESYNC'
  }
}

export function calculateSyncPercentage(healthy: number, total: number): number {
  if (total === 0) return 0
  return Math.round((healthy / total) * 100)
}
