import cytoscape, { Core } from 'cytoscape'
import { useEffect, useMemo, useRef, useState } from 'react'

import { apiClient } from '../api/client'
import { FilterBar } from '../components/FilterBar'
import { FlowTable } from '../components/FlowTable'
import { useGlobalFilters } from '../hooks/useFilters'
import { formatBytes, formatNumber } from '../lib/format'
import { EdgeDetails, GraphResponse, NodeDetails } from '../lib/types'

const MODES = [
  { value: 'host_to_host', label: 'Host to Host' },
  { value: 'subnet_to_subnet', label: 'Subnet to Subnet' },
  { value: 'service_to_service', label: 'Service to Service' },
  { value: 'internal_to_external', label: 'Internal to External' },
]

const GROUPS = [
  { value: 'none', label: 'No Grouping' },
  { value: 'subnet', label: 'Group by Subnet' },
  { value: 'service', label: 'Group by Service' },
  { value: 'environment', label: 'Group by Environment' },
  { value: 'asn', label: 'Group by ASN' },
]

export function MapPage() {
  const { filters, setFilters, params } = useGlobalFilters()
  const [mode, setMode] = useState('host_to_host')
  const [groupBy, setGroupBy] = useState('none')
  const [minBytes, setMinBytes] = useState('0')
  const [minFlows, setMinFlows] = useState('0')
  const [egoNode, setEgoNode] = useState('')
  const [graph, setGraph] = useState<GraphResponse | null>(null)
  const [nodeDetails, setNodeDetails] = useState<NodeDetails | null>(null)
  const [edgeDetails, setEdgeDetails] = useState<EdgeDetails | null>(null)
  const [selectedNodes, setSelectedNodes] = useState<string[]>([])
  const [selectedEdges, setSelectedEdges] = useState<string[]>([])
  const [selectedNodeMeta, setSelectedNodeMeta] = useState<any | null>(null)

  const containerRef = useRef<HTMLDivElement>(null)
  const cyRef = useRef<Core | null>(null)

  useEffect(() => {
    const q = new URLSearchParams(params)
    q.set('mode', mode)
    q.set('group_by', groupBy)
    if (minBytes !== '0') q.set('min_bytes', minBytes)
    if (minFlows !== '0') q.set('min_flows', minFlows)
    if (egoNode.trim()) q.set('node_id', egoNode.trim())

    let active = true
    apiClient.get('/api/map/graph?' + q.toString()).then((res) => {
      if (!active) return
      setGraph(res.data)
    })
    return () => {
      active = false
    }
  }, [params, mode, groupBy, minBytes, minFlows, egoNode])

  const elements = useMemo(() => {
    if (!graph) return []
    return [
      ...graph.nodes.map((n) => ({
        data: {
          id: n.id,
          label: n.label,
          bytesIn: n.bytes_in,
          bytesOut: n.bytes_out,
          packetsIn: n.packets_in,
          packetsOut: n.packets_out,
          flowCount: n.flow_count,
          private: n.private,
          internal: n.internal,
          collapsed: !!n.collapsed,
          children: n.children ?? [],
          tags: n.tags,
        },
      })),
      ...graph.edges.map((e) => ({
        data: {
          id: e.id,
          source: e.source,
          target: e.destination,
          bytes: e.bytes,
          packets: e.packets,
          flows: e.flows,
          protocols: e.protocols,
          ports: e.top_destination_ports,
        },
      })),
    ]
  }, [graph])

  useEffect(() => {
    if (!containerRef.current) return
    if (cyRef.current) {
      cyRef.current.destroy()
      cyRef.current = null
    }
    const cyStyle: any = [
        {
          selector: 'node',
          style: {
            'background-color': (ele: any) => {
              const internal = Boolean(ele.data('internal'))
              const privateNode = Boolean(ele.data('private'))
              if (internal && privateNode) return '#22c55e'
              if (internal) return '#10b981'
              if (privateNode) return '#f59e0b'
              return '#0ea5e9'
            },
            'border-color': '#f8fafc',
            'border-width': (ele: any) => (ele.data('collapsed') ? 3 : 1),
            width: (ele: any) => Math.max(20, Math.min(80, 16 + Math.log10((ele.data('flowCount') || 1) + 1) * 12)),
            height: (ele: any) => Math.max(20, Math.min(80, 16 + Math.log10((ele.data('flowCount') || 1) + 1) * 12)),
            label: 'data(label)',
            'font-size': 10,
            color: '#e2e8f0',
            'text-wrap': 'wrap',
            'text-max-width': 120,
          },
        },
        {
          selector: 'edge',
          style: {
            width: (ele: any) => Math.max(1, Math.min(12, Math.log10((ele.data('bytes') || 1) + 1))),
            'line-color': '#38bdf8',
            'target-arrow-color': '#38bdf8',
            'target-arrow-shape': 'triangle',
            'curve-style': 'bezier',
            opacity: 0.7,
          },
        },
        {
          selector: '.faded',
          style: {
            opacity: 0.1,
          },
        },
        {
          selector: '.highlight',
          style: {
            opacity: 1,
            'line-color': '#f97316',
            'target-arrow-color': '#f97316',
            'background-color': '#f97316',
          },
        },
      ]
    const cy = cytoscape({
      container: containerRef.current,
      elements,
      style: cyStyle,
      layout: {
        name: 'cose',
        animate: false,
        padding: 30,
      },
    })

    cy.on('tap', 'node', async (evt: any) => {
      const node = evt.target
      const id = String(node.id())
      setSelectedNodeMeta(node.data())
      setEdgeDetails(null)
      setSelectedNodes([id])
      cy.elements().addClass('faded').removeClass('highlight')
      node.removeClass('faded').addClass('highlight')
      const neighbors = node.neighborhood()
      neighbors.removeClass('faded').addClass('highlight')
      const q = new URLSearchParams(params)
      const { data } = await apiClient.get(`/api/map/node/${encodeURIComponent(id)}?` + q.toString())
      setNodeDetails(data)
    })

    cy.on('tap', 'edge', async (evt: any) => {
      const edge = evt.target
      const id = String(edge.id())
      setNodeDetails(null)
      setSelectedEdges([id])
      cy.elements().addClass('faded').removeClass('highlight')
      edge.removeClass('faded').addClass('highlight')
      edge.source().removeClass('faded').addClass('highlight')
      edge.target().removeClass('faded').addClass('highlight')
      const q = new URLSearchParams(params)
      q.set('page', '1')
      q.set('page_size', '30')
      const { data } = await apiClient.get(`/api/map/edge/${encodeURIComponent(id)}?` + q.toString())
      setEdgeDetails(data)
    })

    cy.on('tap', (evt) => {
      if (evt.target === cy) {
        cy.elements().removeClass('faded highlight')
      }
    })

    cyRef.current = cy
    return () => {
      cy.destroy()
      cyRef.current = null
    }
  }, [elements, params])

  const exportSelection = () => {
    if (!graph) return
    const exportData = {
      mode,
      group_by: groupBy,
      selected_nodes: selectedNodes,
      selected_edges: selectedEdges,
      nodes: graph.nodes.filter((n) => selectedNodes.length === 0 || selectedNodes.includes(n.id)),
      edges: graph.edges.filter((e) => selectedEdges.length === 0 || selectedEdges.includes(e.id)),
      exported_at: new Date().toISOString(),
    }
    const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `flowscope-graph-selection-${Date.now()}.json`
    a.click()
    URL.revokeObjectURL(url)
  }

  return (
    <div className="space-y-4">
      <FilterBar filters={filters} setFilters={setFilters} compact />

      <section className="panel space-y-3 px-3 py-3">
        <div className="grid grid-cols-2 gap-2 md:grid-cols-6">
          <label className="flex flex-col gap-1">
            <span className="label">Graph Mode</span>
            <select value={mode} onChange={(e) => setMode(e.target.value)}>
              {MODES.map((m) => (
                <option key={m.value} value={m.value}>{m.label}</option>
              ))}
            </select>
          </label>
          <label className="flex flex-col gap-1">
            <span className="label">Group / Collapse</span>
            <select value={groupBy} onChange={(e) => setGroupBy(e.target.value)}>
              {GROUPS.map((m) => (
                <option key={m.value} value={m.value}>{m.label}</option>
              ))}
            </select>
          </label>
          <label className="flex flex-col gap-1">
            <span className="label">Min Bytes</span>
            <input value={minBytes} onChange={(e) => setMinBytes(e.target.value)} />
          </label>
          <label className="flex flex-col gap-1">
            <span className="label">Min Flows</span>
            <input value={minFlows} onChange={(e) => setMinFlows(e.target.value)} />
          </label>
          <label className="flex flex-col gap-1">
            <span className="label">Isolate Ego Node</span>
            <input value={egoNode} onChange={(e) => setEgoNode(e.target.value)} placeholder="node id" />
          </label>
          <div className="flex items-end gap-2">
            <button className="secondary" onClick={() => setEgoNode('')}>Clear Ego</button>
            <button className="secondary" onClick={exportSelection}>Export JSON</button>
          </div>
        </div>
      </section>

      <div className="grid grid-cols-1 gap-4 xl:grid-cols-[minmax(0,1fr)_390px]">
        <div className="panel p-2">
          <div className="mb-2 px-2 text-xs text-slate-400">
            Click a node to highlight neighbors. Click an edge for drill-down. Drag, zoom, and pan directly in the graph.
          </div>
          <div ref={containerRef} className="h-[680px] w-full rounded-lg border border-slate-700/50 bg-slate-950/40" />
        </div>

        <aside className="space-y-4">
          <div className="panel p-3">
            <h3 className="mb-2 text-sm font-semibold">Suspicious Interactions</h3>
            <ul className="space-y-2 text-sm">
              {(graph?.suspicious_interactions ?? []).map((s, idx) => (
                <li key={`${s.kind}-${idx}`} className="rounded border border-amber-500/40 bg-amber-950/20 p-2">
                  <div className="text-xs uppercase text-amber-300">{s.kind}</div>
                  <div>{s.description}</div>
                </li>
              ))}
              {!(graph?.suspicious_interactions?.length) && <li className="text-slate-400">No suspicious interactions in this window.</li>}
            </ul>
          </div>

          <div className="panel p-3">
            <h3 className="mb-2 text-sm font-semibold">Node Details</h3>
            {selectedNodeMeta?.collapsed && (
              <button
                className="secondary mb-2"
                onClick={() => {
                  setGroupBy('none')
                  const firstChild = selectedNodeMeta.children?.[0]
                  if (firstChild) setEgoNode(firstChild)
                }}
              >
                Expand Selected Group
              </button>
            )}
            {nodeDetails ? (
              <div className="space-y-2 text-sm">
                <button className="secondary" onClick={() => setEgoNode(nodeDetails.node.id)}>Isolate This Node</button>
                <Metric label="Node" value={nodeDetails.node.label} />
                <Metric label="Inbound" value={`${formatBytes(nodeDetails.inbound_totals.bytes)} / ${formatNumber(nodeDetails.inbound_totals.flows)} flows`} />
                <Metric label="Outbound" value={`${formatBytes(nodeDetails.outbound_totals.bytes)} / ${formatNumber(nodeDetails.outbound_totals.flows)} flows`} />
                <div>
                  <div className="label">Top Peers</div>
                  <ul className="mt-1 space-y-1">
                    {nodeDetails.top_peers.slice(0, 6).map((p) => <li key={p.key}>{p.key} ({formatBytes(p.bytes)})</li>)}
                  </ul>
                </div>
              </div>
            ) : (
              <p className="text-sm text-slate-400">Select a node on the graph.</p>
            )}
          </div>

          <div className="panel p-3">
            <h3 className="mb-2 text-sm font-semibold">Edge Details and Matching Flows</h3>
            {edgeDetails ? (
              <div className="space-y-2">
                <Metric label="Edge" value={`${edgeDetails.edge.source} -> ${edgeDetails.edge.destination}`} mono />
                <Metric label="Traffic" value={`${formatBytes(edgeDetails.edge.bytes)} / ${formatNumber(edgeDetails.edge.flows)} flows`} />
                <Metric label="Protocols" value={edgeDetails.edge.protocols.join(', ')} />
                <FlowTable rows={edgeDetails.flows.data} />
              </div>
            ) : (
              <p className="text-sm text-slate-400">Select an edge to drill down to flows.</p>
            )}
          </div>
        </aside>
      </div>
    </div>
  )
}

function Metric({ label, value, mono = false }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="rounded border border-slate-700/40 bg-slate-900/40 px-2 py-1">
      <div className="label">{label}</div>
      <div className={mono ? 'value break-all' : 'break-all text-slate-200'}>{value}</div>
    </div>
  )
}
