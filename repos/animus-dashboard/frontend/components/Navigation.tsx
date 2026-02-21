'use client'

import { useState } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { Home, Server, ScrollText, Settings, Menu, X, LogOut, User } from 'lucide-react'
import { GlitchText } from '@/components/ui/GlitchText'
import { Node } from '@/lib/api'
import { clsx } from 'clsx'

interface NavigationProps {
  nodes: Node[]
  user?: { name: string; email: string }
}

export function Navigation({ nodes, user }: NavigationProps) {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const pathname = usePathname()

  const mainLinks = [
    { href: '/', label: 'Animus Core', icon: Home },
    { href: '/scripts', label: 'Protocols', icon: ScrollText },
    { href: '/settings', label: 'Settings', icon: Settings },
  ]

  const isActive = (href: string) => {
    if (href === '/') return pathname === '/'
    return pathname.startsWith(href)
  }

  const isNodeActive = (nodeId: string) => pathname === `/node/${nodeId}`

  return (
    <nav className="bg-animus-navy border-b border-animus-gold/20">
      <div className="max-w-screen-2xl mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-animus-gold to-animus-cyan flex items-center justify-center">
              <span className="text-animus-black font-bold text-lg">A</span>
            </div>
            <div>
              <GlitchText
                text="ANIMUS v2.0"
                className="text-lg font-bold text-animus-gold tracking-wider"
              />
              <p className="text-xs text-animus-text-secondary -mt-0.5">
                Cluster Management System
              </p>
            </div>
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center gap-1">
            {/* Main links */}
            {mainLinks.map((link) => {
              const Icon = link.icon
              return (
                <Link
                  key={link.href}
                  href={link.href}
                  className={clsx(
                    'nav-tab flex items-center gap-2',
                    isActive(link.href) && 'active'
                  )}
                >
                  <Icon className="w-4 h-4" />
                  {link.label}
                </Link>
              )
            })}

            {/* Divider */}
            <div className="h-8 w-px bg-animus-gold/20 mx-2" />

            {/* Node tabs - scrollable */}
            <div className="flex items-center gap-1 overflow-x-auto max-w-lg scrollbar-hide">
              {nodes.map((node) => (
                <Link
                  key={node.id}
                  href={`/node/${node.id}`}
                  className={clsx(
                    'nav-tab whitespace-nowrap flex items-center gap-2',
                    isNodeActive(node.id) && 'active'
                  )}
                >
                  <Server className="w-3 h-3" />
                  {node.name}
                </Link>
              ))}
            </div>
          </div>

          {/* User menu */}
          <div className="hidden md:flex items-center gap-4">
            {user && (
              <div className="flex items-center gap-3">
                <div className="text-right">
                  <p className="text-sm text-white">{user.name}</p>
                  <p className="text-xs text-animus-text-secondary">{user.email}</p>
                </div>
                <div className="w-8 h-8 rounded-full bg-animus-gold/20 flex items-center justify-center">
                  <User className="w-4 h-4 text-animus-gold" />
                </div>
                <button className="p-2 text-animus-text-secondary hover:text-animus-red transition-colors">
                  <LogOut className="w-4 h-4" />
                </button>
              </div>
            )}
          </div>

          {/* Mobile menu button */}
          <button
            className="md:hidden p-2 text-animus-gold"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          >
            {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>
      </div>

      {/* Mobile menu */}
      {isMobileMenuOpen && (
        <div className="md:hidden border-t border-animus-gold/20 bg-animus-navy">
          <div className="px-4 py-4 space-y-2">
            {mainLinks.map((link) => {
              const Icon = link.icon
              return (
                <Link
                  key={link.href}
                  href={link.href}
                  onClick={() => setIsMobileMenuOpen(false)}
                  className={clsx(
                    'flex items-center gap-3 px-4 py-3 rounded-lg transition-colors',
                    isActive(link.href)
                      ? 'bg-animus-gold/10 text-animus-cyan'
                      : 'text-animus-text-secondary hover:bg-animus-gold/5'
                  )}
                >
                  <Icon className="w-5 h-5" />
                  {link.label}
                </Link>
              )
            })}

            <div className="pt-2 border-t border-animus-gold/10">
              <p className="px-4 py-2 text-xs text-animus-text-secondary uppercase tracking-wider">
                Memory Cores
              </p>
              {nodes.map((node) => (
                <Link
                  key={node.id}
                  href={`/node/${node.id}`}
                  onClick={() => setIsMobileMenuOpen(false)}
                  className={clsx(
                    'flex items-center gap-3 px-4 py-3 rounded-lg transition-colors',
                    isNodeActive(node.id)
                      ? 'bg-animus-gold/10 text-animus-cyan'
                      : 'text-animus-text-secondary hover:bg-animus-gold/5'
                  )}
                >
                  <Server className="w-4 h-4" />
                  {node.name}
                </Link>
              ))}
            </div>
          </div>
        </div>
      )}
    </nav>
  )
}
