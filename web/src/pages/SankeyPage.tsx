import { useEffect, useMemo, useState } from 'react'
import ReactECharts from 'echarts-for-react'

import { apiClient } from '../api/client'
import { FilterBar } from '../components/FilterBar'
import { FlowTable } from '../components/FlowTable'
import { useGlobalFilters } from '../hooks/useFilters'
import { formatBytes, formatNumber } from '../lib/format'
import { FlowRecord, SankeyResponse } from '../lib/types'

export function SankeyPage() {
  const { filters, setFilters, params } = useGlobalFilters()
  const [sankey, setSankey] = useState<SankeyResponse>({ nodes: [], links: [] })
  const [flows, setFlows] = useState<FlowRecord[]>([])
  const [selection, setSelection] = useState('')

  useEffect(() => {
    let active = true
    apiClient.get('/api/sankey?' + params.toString()).then((res) => {
      if (active) setSankey(res.data)
    })
    return () => {
      active = false
    }
  }, [params])

  const option = useMemo(
    () => ({
      tooltip: {
        trigger: 'item',
        formatter: (p: any) => {
          if (p.dataType === 'edge') {
            return `${p.data.source} -> ${p.data.target}<br/>${formatBytes(p.data.value)} / ${formatNumber(p.data.flows)} flows`
          }
          return p.name
        },
      },
      series: [
        {
          type: 'sankey',
          data: sankey.nodes,
          links: sankey.links,
          emphasis: { focus: 'adjacency' },
          lineStyle: { color: 'gradient', curveness: 0.5 },
          nodeAlign: 'left',
        },
      ],
    }),
    [sankey],
  )

  const onEvents = {
    click: async (event: any) => {
      if (event.dataType !== 'edge') return
      const src = event.data.source as string
      const dst = event.data.target as string
      setSelection(`${src} -> ${dst}`)
      const q = new URLSearchParams(params)
      q.set('search', src)
      q.set('page_size', '30')
      q.set('sort_by', 'bytes')
      q.set('sort_dir', 'DESC')
      const { data } = await apiClient.get('/api/flows/historical?' + q.toString())
      const filtered = (data.data as FlowRecord[]).filter((f) => {
        const fields = [f.src_ip, f.dst_ip, f.src_hostname, f.dst_hostname, f.src_service, f.dst_service].filter(Boolean)
        return fields.some((x) => x?.includes(src) || x?.includes(dst))
      })
      setFlows(filtered)
    },
  }

  return (
    <div className="space-y-4">
      <FilterBar filters={filters} setFilters={setFilters} compact />
      <div className="panel p-3">
        <h2 className="mb-2 text-sm font-semibold">Traffic Sankey</h2>
        <p className="mb-2 text-xs text-slate-400">Source - protocol/service - destination. Click a band to drill into matching flows.</p>
        <ReactECharts option={option} style={{ height: 500 }} onEvents={onEvents} />
      </div>
      <div className="panel p-3">
        <h3 className="mb-2 text-sm font-semibold">Drill-down Flows {selection && <span className="text-slate-400">for {selection}</span>}</h3>
        <FlowTable rows={flows} />
      </div>
    </div>
  )
}
