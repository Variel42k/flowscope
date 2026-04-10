import { useEffect, useMemo, useState } from 'react'
import ReactECharts from 'echarts-for-react'

import { apiClient } from '../api/client'
import { FilterBar } from '../components/FilterBar'
import { KpiCard } from '../components/KpiCard'
import { useGlobalFilters } from '../hooks/useFilters'
import { formatBytes, formatNumber } from '../lib/format'
import { GraphResponse, TopItem } from '../lib/types'

export function OverviewPage() {
  const { filters, setFilters, params } = useGlobalFilters()
  const [talkers, setTalkers] = useState<TopItem[]>([])
  const [protocols, setProtocols] = useState<TopItem[]>([])
  const [interfaces, setInterfaces] = useState<TopItem[]>([])
  const [exporters, setExporters] = useState<any[]>([])
  const [graph, setGraph] = useState<GraphResponse | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    let active = true
    setLoading(true)
    Promise.all([
      apiClient.get('/api/talkers/top?' + params.toString()),
      apiClient.get('/api/protocols/top?' + params.toString()),
      apiClient.get('/api/interfaces/top?' + params.toString()),
      apiClient.get('/api/exporters?' + params.toString()),
      apiClient.get('/api/map/graph?' + new URLSearchParams({ ...Object.fromEntries(params), mode: 'host_to_host' }).toString()),
    ])
      .then(([t, p, i, e, g]) => {
        if (!active) return
        setTalkers(t.data.data ?? [])
        setProtocols(p.data.data ?? [])
        setInterfaces(i.data.data ?? [])
        setExporters(e.data.data ?? [])
        setGraph(g.data)
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => {
      active = false
    }
  }, [params])

  const totals = useMemo(() => {
    const bytes = talkers.reduce((acc, t) => acc + (t.bytes ?? 0), 0)
    const flows = talkers.reduce((acc, t) => acc + (t.flows ?? 0), 0)
    return {
      bytes,
      flows,
      edges: graph?.edges.length ?? 0,
      exporters: exporters.length,
    }
  }, [talkers, exporters, graph])

  const topProtocolChart = {
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: protocols.map((p) => p.key) },
    yAxis: { type: 'value' },
    series: [{ type: 'bar', data: protocols.map((p) => p.bytes), itemStyle: { color: '#22c55e' } }],
    grid: { top: 30, left: 50, right: 20, bottom: 40 },
  }

  const topTalkersChart = {
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: talkers.map((t) => t.key).slice(0, 10) },
    yAxis: { type: 'value' },
    series: [{ type: 'bar', data: talkers.map((t) => t.bytes).slice(0, 10), itemStyle: { color: '#0ea5e9' } }],
    grid: { top: 30, left: 50, right: 20, bottom: 80 },
  }

  return (
    <div className="space-y-4">
      <FilterBar filters={filters} setFilters={setFilters} />
      <div className="grid grid-cols-2 gap-3 md:grid-cols-4">
        <KpiCard title="Total Throughput" value={formatBytes(totals.bytes)} subtitle="filtered window" />
        <KpiCard title="Total Flows" value={formatNumber(totals.flows)} subtitle="top-talker sum" />
        <KpiCard title="Active Exporters" value={formatNumber(totals.exporters)} />
        <KpiCard title="Top Edges" value={formatNumber(totals.edges)} subtitle="graph edge count" />
      </div>

      <div className="grid grid-cols-1 gap-3 lg:grid-cols-2">
        <div className="panel px-3 py-3">
          <h3 className="mb-2 text-sm font-semibold">Top Protocols by Bytes</h3>
          <ReactECharts option={topProtocolChart} style={{ height: 280 }} />
        </div>
        <div className="panel px-3 py-3">
          <h3 className="mb-2 text-sm font-semibold">Top Talkers by Bytes</h3>
          <ReactECharts option={topTalkersChart} style={{ height: 280 }} />
        </div>
      </div>

      <div className="grid grid-cols-1 gap-3 lg:grid-cols-3">
        <div className="panel px-3 py-3">
          <h3 className="mb-2 text-sm font-semibold">Top Interfaces</h3>
          <ul className="space-y-1 text-sm">
            {interfaces.slice(0, 10).map((item) => (
              <li key={item.key} className="flex justify-between border-b border-slate-700/30 pb-1">
                <span>{item.key}</span>
                <span className="font-mono">{formatBytes(item.bytes)}</span>
              </li>
            ))}
          </ul>
        </div>
        <div className="panel px-3 py-3">
          <h3 className="mb-2 text-sm font-semibold">Recent Suspicious Interactions</h3>
          <ul className="space-y-2 text-sm">
            {(graph?.suspicious_interactions ?? []).slice(0, 10).map((s, idx) => (
              <li key={`${s.kind}-${idx}`} className="rounded border border-amber-500/40 bg-amber-950/20 p-2">
                <div className="text-xs uppercase text-amber-300">{s.kind}</div>
                <div>{s.description}</div>
              </li>
            ))}
            {!graph?.suspicious_interactions?.length && <li className="text-slate-400">No heuristic anomalies in current view.</li>}
          </ul>
        </div>
        <div className="panel px-3 py-3">
          <h3 className="mb-2 text-sm font-semibold">Top Edges</h3>
          <ul className="space-y-1 text-sm">
            {graph?.edges.slice(0, 10).map((edge) => (
              <li key={edge.id} className="rounded border border-slate-700/40 px-2 py-1">
                <div className="truncate text-xs text-slate-400">{edge.source} -&gt; {edge.destination}</div>
                <div className="font-mono">{formatBytes(edge.bytes)} / {formatNumber(edge.flows)} flows</div>
              </li>
            ))}
          </ul>
        </div>
      </div>

      {loading && <div className="text-sm text-slate-400">Loading overview...</div>}
    </div>
  )
}
