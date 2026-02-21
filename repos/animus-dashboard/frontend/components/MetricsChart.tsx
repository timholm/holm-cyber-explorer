'use client'

import { useState, useEffect } from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Area,
  AreaChart,
} from 'recharts'
import { Card, CardHeader, CardTitle } from '@/components/ui/Card'

interface MetricsChartProps {
  title: string
  data: Array<{ time: string; value: number }>
  color?: 'gold' | 'cyan' | 'green'
  unit?: string
  type?: 'line' | 'area'
}

const colors = {
  gold: '#c9a227',
  cyan: '#00d4ff',
  green: '#00ff88',
}

export function MetricsChart({
  title,
  data,
  color = 'cyan',
  unit = '%',
  type = 'area',
}: MetricsChartProps) {
  const strokeColor = colors[color]

  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="bg-animus-navy border border-animus-gold/30 rounded px-3 py-2 shadow-lg">
          <p className="text-xs text-animus-text-secondary">{label}</p>
          <p className="text-sm font-mono" style={{ color: strokeColor }}>
            {payload[0].value}
            {unit}
          </p>
        </div>
      )
    }
    return null
  }

  const ChartComponent = type === 'area' ? AreaChart : LineChart

  return (
    <Card variant="default" padding="sm">
      <CardHeader className="pb-2">
        <CardTitle className="text-sm">{title}</CardTitle>
      </CardHeader>
      <div className="h-40">
        <ResponsiveContainer width="100%" height="100%">
          <ChartComponent data={data}>
            <defs>
              <linearGradient id={`gradient-${color}`} x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor={strokeColor} stopOpacity={0.3} />
                <stop offset="95%" stopColor={strokeColor} stopOpacity={0} />
              </linearGradient>
            </defs>
            <CartesianGrid strokeDasharray="3 3" stroke="#1a1a2a" />
            <XAxis
              dataKey="time"
              stroke="#8888aa"
              tick={{ fontSize: 10 }}
              tickLine={false}
            />
            <YAxis
              stroke="#8888aa"
              tick={{ fontSize: 10 }}
              tickLine={false}
              domain={[0, 100]}
            />
            <Tooltip content={<CustomTooltip />} />
            {type === 'area' ? (
              <Area
                type="monotone"
                dataKey="value"
                stroke={strokeColor}
                strokeWidth={2}
                fill={`url(#gradient-${color})`}
              />
            ) : (
              <Line
                type="monotone"
                dataKey="value"
                stroke={strokeColor}
                strokeWidth={2}
                dot={false}
              />
            )}
          </ChartComponent>
        </ResponsiveContainer>
      </div>
    </Card>
  )
}

// Hook to generate chart data from real-time metrics
export function useMetricsHistory(currentValue: number, maxPoints: number = 20) {
  const [history, setHistory] = useState<Array<{ time: string; value: number }>>([])

  useEffect(() => {
    const now = new Date().toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })

    setHistory((prev) => {
      const newHistory = [...prev, { time: now, value: currentValue }]
      if (newHistory.length > maxPoints) {
        return newHistory.slice(-maxPoints)
      }
      return newHistory
    })
  }, [currentValue, maxPoints])

  return history
}
