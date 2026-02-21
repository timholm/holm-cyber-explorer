'use client'

import { useState, useEffect, useRef } from 'react'
import { Play, Square, Terminal, CheckCircle, XCircle, Clock, ChevronDown, ChevronRight } from 'lucide-react'
import { Card, CardHeader, CardTitle } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { Select } from '@/components/ui/Select'
import { api, Script, ScriptExecution, Node } from '@/lib/api'
import { useWebSocket } from '@/hooks/useWebSocket'
import { clsx } from 'clsx'

interface ScriptRunnerProps {
  nodes: Node[]
  scripts: Script[]
  selectedNode?: string
}

interface ExecutionOutput {
  type: 'stdout' | 'stderr' | 'status'
  line: string
  timestamp: string
}

export function ScriptRunner({ nodes, scripts, selectedNode }: ScriptRunnerProps) {
  const [selectedScript, setSelectedScript] = useState('')
  const [targetNodes, setTargetNodes] = useState<string[]>(selectedNode ? [selectedNode] : [])
  const [isRunning, setIsRunning] = useState(false)
  const [currentExecution, setCurrentExecution] = useState<ScriptExecution | null>(null)
  const [output, setOutput] = useState<ExecutionOutput[]>([])
  const [expandedSections, setExpandedSections] = useState<{ [key: string]: boolean }>({
    system: true,
    kubernetes: true,
    monitoring: true,
    custom: true,
  })

  const outputRef = useRef<HTMLDivElement>(null)

  // WebSocket for real-time script output
  const { data: wsData, isConnected } = useWebSocket<ExecutionOutput>(
    currentExecution ? `/ws/scripts/${currentExecution.id}` : ''
  )

  // Append WebSocket data to output
  useEffect(() => {
    if (wsData) {
      setOutput((prev) => [...prev, wsData])
    }
  }, [wsData])

  // Auto-scroll output
  useEffect(() => {
    if (outputRef.current) {
      outputRef.current.scrollTop = outputRef.current.scrollHeight
    }
  }, [output])

  // Group scripts by category
  const scriptsByCategory = scripts.reduce((acc, script) => {
    const category = script.category || 'custom'
    if (!acc[category]) {
      acc[category] = []
    }
    acc[category].push(script)
    return acc
  }, {} as { [key: string]: Script[] })

  const handleRunScript = async () => {
    if (!selectedScript || targetNodes.length === 0) return

    try {
      setIsRunning(true)
      setOutput([])

      const execution = await api.runScript(selectedScript, targetNodes)
      setCurrentExecution(execution)

      // Poll for completion if WebSocket fails
      const pollInterval = setInterval(async () => {
        const updated = await api.getScriptExecution(execution.id)
        if (updated.status !== 'running') {
          setIsRunning(false)
          setCurrentExecution(updated)
          clearInterval(pollInterval)
        }
      }, 2000)

    } catch (error) {
      setIsRunning(false)
      setOutput([{
        type: 'stderr',
        line: `Error: ${error instanceof Error ? error.message : 'Failed to run script'}`,
        timestamp: new Date().toISOString(),
      }])
    }
  }

  const handleStopScript = async () => {
    if (!currentExecution) return
    // API call to stop script would go here
    setIsRunning(false)
  }

  const toggleSection = (category: string) => {
    setExpandedSections((prev) => ({
      ...prev,
      [category]: !prev[category],
    }))
  }

  const categoryLabels: { [key: string]: string } = {
    system: 'System Maintenance',
    kubernetes: 'Kubernetes Operations',
    monitoring: 'Monitoring',
    custom: 'Custom Protocols',
  }

  const nodeOptions = [
    { value: 'all', label: 'All Nodes' },
    ...nodes.map((node) => ({ value: node.id, label: node.name })),
  ]

  return (
    <Card variant="default" padding="none">
      <CardHeader className="px-4 pt-4">
        <CardTitle className="flex items-center gap-2">
          <Terminal className="w-5 h-5 text-animus-gold" />
          Protocols (Scripts)
        </CardTitle>
      </CardHeader>

      <div className="p-4 border-b border-animus-gold/10">
        <div className="flex items-center gap-4">
          <div className="w-48">
            <Select
              options={nodeOptions}
              value={targetNodes[0] || ''}
              onChange={(e) => {
                if (e.target.value === 'all') {
                  setTargetNodes(nodes.map((n) => n.id))
                } else {
                  setTargetNodes([e.target.value])
                }
              }}
              placeholder="Select target..."
            />
          </div>
          <Button
            variant="primary"
            onClick={handleRunScript}
            disabled={!selectedScript || targetNodes.length === 0 || isRunning}
            isLoading={isRunning}
          >
            <Play className="w-4 h-4" />
            Execute
          </Button>
          {isRunning && (
            <Button variant="danger" onClick={handleStopScript}>
              <Square className="w-4 h-4" />
              Stop
            </Button>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-0 border-b border-animus-gold/10">
        {/* Script Selection */}
        <div className="border-r border-animus-gold/10 max-h-80 overflow-y-auto">
          {Object.entries(scriptsByCategory).map(([category, categoryScripts]) => (
            <div key={category}>
              <button
                onClick={() => toggleSection(category)}
                className="w-full flex items-center justify-between px-4 py-2 bg-animus-black/30 hover:bg-animus-gold/5 transition-colors"
              >
                <span className="text-sm font-medium text-animus-text-secondary uppercase tracking-wider">
                  {categoryLabels[category] || category}
                </span>
                {expandedSections[category] ? (
                  <ChevronDown className="w-4 h-4 text-animus-text-secondary" />
                ) : (
                  <ChevronRight className="w-4 h-4 text-animus-text-secondary" />
                )}
              </button>
              {expandedSections[category] && (
                <div className="py-1">
                  {categoryScripts.map((script) => (
                    <button
                      key={script.id}
                      onClick={() => setSelectedScript(script.id)}
                      className={clsx(
                        'w-full text-left px-4 py-2 hover:bg-animus-gold/5 transition-colors',
                        selectedScript === script.id && 'bg-animus-gold/10 border-l-2 border-animus-gold'
                      )}
                    >
                      <div className="text-sm text-white">{script.name}</div>
                      <div className="text-xs text-animus-text-secondary mt-0.5">
                        {script.description}
                      </div>
                    </button>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>

        {/* Output Terminal */}
        <div className="flex flex-col">
          <div className="px-4 py-2 bg-animus-black/30 flex items-center justify-between">
            <span className="text-xs text-animus-text-secondary font-mono">OUTPUT</span>
            {currentExecution && (
              <div className="flex items-center gap-2">
                {currentExecution.status === 'running' && (
                  <Clock className="w-4 h-4 text-animus-gold animate-pulse" />
                )}
                {currentExecution.status === 'completed' && (
                  <CheckCircle className="w-4 h-4 text-animus-green" />
                )}
                {currentExecution.status === 'failed' && (
                  <XCircle className="w-4 h-4 text-animus-red" />
                )}
                <span className="text-xs text-animus-text-secondary">
                  {currentExecution.status}
                </span>
              </div>
            )}
          </div>
          <div
            ref={outputRef}
            className="flex-1 min-h-[200px] max-h-80 overflow-y-auto bg-animus-black/50 p-3 font-mono text-xs"
          >
            {output.length === 0 ? (
              <div className="text-animus-text-secondary">
                Select a protocol and target to execute...
              </div>
            ) : (
              output.map((line, index) => (
                <div
                  key={index}
                  className={clsx(
                    'py-0.5',
                    line.type === 'stderr' && 'text-animus-red',
                    line.type === 'stdout' && 'text-white',
                    line.type === 'status' && 'text-animus-cyan'
                  )}
                >
                  {line.line}
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </Card>
  )
}
