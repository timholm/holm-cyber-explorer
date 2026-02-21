'use client'

import { HTMLAttributes, forwardRef } from 'react'
import { clsx } from 'clsx'

type Status = 'synchronized' | 'warning' | 'desynchronized'

interface StatusIndicatorProps extends HTMLAttributes<HTMLDivElement> {
  status: Status
  label?: string
  pulse?: boolean
  size?: 'sm' | 'md' | 'lg'
}

const statusConfig = {
  synchronized: {
    color: 'bg-animus-green',
    glow: 'shadow-[0_0_10px_rgba(0,255,136,0.5)]',
    text: 'text-animus-green',
    label: 'SYNC',
  },
  warning: {
    color: 'bg-animus-gold',
    glow: 'shadow-[0_0_10px_rgba(201,162,39,0.5)]',
    text: 'text-animus-gold',
    label: 'WARN',
  },
  desynchronized: {
    color: 'bg-animus-red',
    glow: 'shadow-[0_0_10px_rgba(255,51,102,0.5)]',
    text: 'text-animus-red',
    label: 'DESYNC',
  },
}

const sizes = {
  sm: 'w-2 h-2',
  md: 'w-3 h-3',
  lg: 'w-4 h-4',
}

export const StatusIndicator = forwardRef<HTMLDivElement, StatusIndicatorProps>(
  ({ status, label, pulse = true, size = 'md', className, ...props }, ref) => {
    const config = statusConfig[status]

    return (
      <div
        ref={ref}
        className={clsx('inline-flex items-center gap-2', className)}
        {...props}
      >
        <span
          className={clsx(
            'rounded-full',
            sizes[size],
            config.color,
            config.glow,
            pulse && 'animate-pulse'
          )}
        />
        {label !== undefined ? (
          <span className={clsx('text-xs font-mono uppercase', config.text)}>
            {label}
          </span>
        ) : (
          <span className={clsx('text-xs font-mono uppercase', config.text)}>
            {config.label}
          </span>
        )}
      </div>
    )
  }
)

StatusIndicator.displayName = 'StatusIndicator'

export function getStatusFromHealth(
  ready: boolean,
  cpuPercent?: number,
  memPercent?: number
): Status {
  if (!ready) return 'desynchronized'
  if (cpuPercent && cpuPercent > 90) return 'warning'
  if (memPercent && memPercent > 90) return 'warning'
  return 'synchronized'
}
