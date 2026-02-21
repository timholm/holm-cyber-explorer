'use client'

import { HTMLAttributes, forwardRef } from 'react'
import { clsx } from 'clsx'

interface DNALoaderProps extends HTMLAttributes<HTMLDivElement> {
  size?: 'sm' | 'md' | 'lg'
  text?: string
}

export const DNALoader = forwardRef<HTMLDivElement, DNALoaderProps>(
  ({ size = 'md', text, className, ...props }, ref) => {
    const sizes = {
      sm: { bar: 'w-1 h-4', gap: 'gap-1' },
      md: { bar: 'w-1 h-6', gap: 'gap-1' },
      lg: { bar: 'w-1.5 h-8', gap: 'gap-1.5' },
    }

    const s = sizes[size]

    return (
      <div
        ref={ref}
        className={clsx('flex flex-col items-center gap-3', className)}
        {...props}
      >
        <div className={clsx('flex items-center', s.gap)}>
          {[...Array(7)].map((_, i) => (
            <span
              key={i}
              className={clsx(
                s.bar,
                'bg-animus-gold rounded-full',
                'animate-[dna-wave_1s_ease-in-out_infinite]'
              )}
              style={{ animationDelay: `${i * 0.1}s` }}
            />
          ))}
        </div>
        {text && (
          <span className="text-sm text-animus-text-secondary font-mono">
            {text}
          </span>
        )}
      </div>
    )
  }
)

DNALoader.displayName = 'DNALoader'
