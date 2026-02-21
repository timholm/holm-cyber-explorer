'use client'

import { HTMLAttributes, forwardRef } from 'react'
import { clsx } from 'clsx'

interface ProgressBarProps extends HTMLAttributes<HTMLDivElement> {
  value: number
  max?: number
  variant?: 'default' | 'warning' | 'danger'
  showLabel?: boolean
  size?: 'sm' | 'md' | 'lg'
}

export const ProgressBar = forwardRef<HTMLDivElement, ProgressBarProps>(
  (
    {
      value,
      max = 100,
      variant = 'default',
      showLabel = false,
      size = 'md',
      className,
      ...props
    },
    ref
  ) => {
    const percentage = Math.min(Math.max((value / max) * 100, 0), 100)

    // Auto-determine variant based on percentage
    const autoVariant =
      variant === 'default'
        ? percentage >= 90
          ? 'danger'
          : percentage >= 70
          ? 'warning'
          : 'default'
        : variant

    const heights = {
      sm: 'h-1',
      md: 'h-2',
      lg: 'h-3',
    }

    const gradients = {
      default: 'bg-gradient-to-r from-animus-gold to-animus-cyan',
      warning: 'bg-gradient-to-r from-animus-gold to-yellow-500',
      danger: 'bg-gradient-to-r from-animus-red to-orange-500',
    }

    return (
      <div ref={ref} className={clsx('w-full', className)} {...props}>
        {showLabel && (
          <div className="flex justify-between text-xs text-animus-text-secondary mb-1">
            <span>{value}%</span>
            <span>{max}%</span>
          </div>
        )}
        <div
          className={clsx(
            'w-full bg-white/10 rounded-full overflow-hidden',
            heights[size]
          )}
        >
          <div
            className={clsx(
              'h-full rounded-full transition-all duration-500',
              gradients[autoVariant]
            )}
            style={{ width: `${percentage}%` }}
          />
        </div>
      </div>
    )
  }
)

ProgressBar.displayName = 'ProgressBar'
