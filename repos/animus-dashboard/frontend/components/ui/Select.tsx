'use client'

import { forwardRef, SelectHTMLAttributes } from 'react'
import { clsx } from 'clsx'
import { ChevronDown } from 'lucide-react'

interface SelectOption {
  value: string
  label: string
}

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label?: string
  options: SelectOption[]
  placeholder?: string
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(
  ({ label, options, placeholder, className, ...props }, ref) => {
    return (
      <div className="w-full">
        {label && (
          <label className="block text-sm text-animus-text-secondary mb-1.5">
            {label}
          </label>
        )}
        <div className="relative">
          <select
            ref={ref}
            className={clsx(
              'w-full appearance-none bg-animus-navy/80 border border-animus-gold/30 rounded-md',
              'px-4 py-2.5 pr-10 text-white text-sm',
              'focus:outline-none focus:border-animus-cyan focus:shadow-[0_0_10px_rgba(0,212,255,0.2)]',
              'transition-all duration-300',
              'cursor-pointer',
              className
            )}
            {...props}
          >
            {placeholder && (
              <option value="" className="bg-animus-navy">
                {placeholder}
              </option>
            )}
            {options.map((option) => (
              <option
                key={option.value}
                value={option.value}
                className="bg-animus-navy"
              >
                {option.label}
              </option>
            ))}
          </select>
          <ChevronDown
            className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-animus-text-secondary pointer-events-none"
          />
        </div>
      </div>
    )
  }
)

Select.displayName = 'Select'
