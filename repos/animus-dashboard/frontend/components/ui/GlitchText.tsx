'use client'

import { HTMLAttributes } from 'react'
import { clsx } from 'clsx'

interface GlitchTextProps extends HTMLAttributes<HTMLElement> {
  text: string
  as?: 'span' | 'h1' | 'h2' | 'h3' | 'p'
}

export function GlitchText({ text, as: Component = 'span', className, ...props }: GlitchTextProps) {
  return (
    <Component
      className={clsx('glitch-text', className)}
      data-text={text}
      {...props}
    >
      {text}
    </Component>
  )
}
