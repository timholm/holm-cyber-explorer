'use client'

import { forwardRef, InputHTMLAttributes } from 'react'
import { clsx } from 'clsx'

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
  icon?: React.ReactNode
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, icon, className, ...props }, ref) => {
    return (
      <div className="w-full">
        {label && (
          <label className="block text-sm text-animus-text-secondary mb-1.5">
            {label}
          </label>
        )}
        <div className="relative">
          {icon && (
            <div className="absolute left-3 top-1/2 -translate-y-1/2 text-animus-text-secondary">
              {icon}
            </div>
          )}
          <input
            ref={ref}
            className={clsx(
              'w-full bg-animus-navy/80 border border-animus-gold/30 rounded-md',
              'px-4 py-2.5 text-white text-sm',
              'placeholder:text-animus-text-secondary',
              'focus:outline-none focus:border-animus-cyan focus:shadow-[0_0_10px_rgba(0,212,255,0.2)]',
              'transition-all duration-300',
              icon && 'pl-10',
              error && 'border-animus-red focus:border-animus-red',
              className
            )}
            {...props}
          />
        </div>
        {error && (
          <p className="mt-1.5 text-xs text-animus-red">{error}</p>
        )}
      </div>
    )
  }
)

Input.displayName = 'Input'
