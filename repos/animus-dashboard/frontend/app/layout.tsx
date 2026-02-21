import type { Metadata } from 'next'
import '@/styles/globals.css'

export const metadata: Metadata = {
  title: 'ANIMUS v2.0 - Cluster Management System',
  description: 'Assassin\'s Creed themed Kubernetes dashboard for Raspberry Pi cluster',
  icons: {
    icon: '/favicon.ico',
  },
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-animus-black hex-bg antialiased">
        {/* Scan lines overlay */}
        <div className="scan-lines" aria-hidden="true" />

        {/* Main content */}
        <div className="relative z-10">
          {children}
        </div>
      </body>
    </html>
  )
}
