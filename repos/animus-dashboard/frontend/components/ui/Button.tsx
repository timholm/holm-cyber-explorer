'use client'

import { forwardRef, ButtonHTMLAttributes } from 'react'
import { clsx } from 'clsx'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  isLoading?: boolean
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      children,
      variant = 'primary',
      size = 'md',
      isLoading = false,
      className,
      disabled,
      ...props
    },
    ref
  ) => {
    const baseStyles =
      'relative inline-flex items-center justify-center font-medium rounded-md transition-all duration-300 overflow-hidden focus:outline-none focus:ring-2 focus:ring-animus-gold focus:ring-offset-2 focus:ring-offset-animus-black'

    const variants = {
      primary:
        'bg-gradient-to-r from-animus-gold/20 to-animus-cyan/10 border border-animus-gold text-animus-gold hover:from-animus-gold/40 hover:to-animus-cyan/20 hover:shadow-[0_0_20px_rgba(201,162,39,0.3)]',
      secondary:
        'bg-animus-navy border border-animus-cyan/30 text-animus-cyan hover:border-animus-cyan hover:shadow-[0_0_15px_rgba(0,212,255,0.2)]',
      danger:
        'bg-animus-red/20 border border-animus-red text-animus-red hover:bg-animus-red/40 hover:shadow-[0_0_15px_rgba(255,51,102,0.3)]',
      ghost:
        'bg-transparent border border-transparent text-animus-text-secondary hover:text-animus-gold hover:border-animus-gold/30',
    }

    const sizes = {
      sm: 'text-xs px-3 py-1.5',
      md: 'text-sm px-4 py-2',
      lg: 'text-base px-6 py-3',
    }

    return (
      <button
        ref={ref}
        disabled={disabled || isLoading}
        className={clsx(
          baseStyles,
          variants[variant],
          sizes[size],
          (disabled || isLoading) && 'opacity-50 cursor-not-allowed',
          className
        )}
        {...props}
      >
        {/* Shimmer effect */}
        <span className="absolute inset-0 overflow-hidden">
          <span className="absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-white/10 to-transparent group-hover:translate-x-full transition-transform duration-500" />
        </span>

        {/* Content */}
        <span className="relative flex items-center gap-2">
          {isLoading && (
            <svg
              className="animate-spin h-4 w-4"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              />
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              />
            </svg>
          )}
          {children}
        </span>
      </button>
    )
  }
)

Button.displayName = 'Button'
