import { useEffect, useMemo, useState } from 'react'
import ReactECharts from 'echarts-for-react'

import { apiClient } from '../api/client'
import { FilterBar } from '../components/FilterBar'
import { FlowTable } from '../components/FlowTable'
import { SavedViewsPanel } from '../components/SavedViewsPanel'
import { useGlobalFilters } from '../hooks/useFilters'
import { formatBytes, formatNumber } from '../lib/format'
import { FlowRecord, PageResult, SankeyResponse } from '../lib/types'

type SankeyTooltipData = {
  source?: string
  target?: string
  value?: number
  flows?: number
}

type SankeyTooltipParam = {
  dataType?: string
  data?: SankeyTooltipData
  name?: string
}

type SankeyClickEvent = {
  dataType?: string
  data?: SankeyTooltipData
}

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
        formatter: (p: SankeyTooltipParam) => {
          if (p.dataType === 'edge') {
            const source = p.data?.source ?? 'unknown'
            const target = p.data?.target ?? 'unknown'
            const value = typeof p.data?.value === 'number' ? p.data.value : 0
            const flows = typeof p.data?.flows === 'number' ? p.data.flows : 0
            return `${source} -> ${target}<br/>${formatBytes(value)} / ${formatNumber(flows)} flows`
          }
          return p.name ?? ''
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
    click: async (event: SankeyClickEvent) => {
      if (event.dataType !== 'edge') return
      const src = event.data?.source
      const dst = event.data?.target
      if (!src || !dst) return
      setSelection(`${src} -> ${dst}`)
      const q = new URLSearchParams(params)
      q.set('search', src)
      q.set('page_size', '30')
      q.set('sort_by', 'bytes')
      q.set('sort_dir', 'DESC')
      const { data } = await apiClient.get<PageResult<FlowRecord>>('/api/flows/historical?' + q.toString())
      const filtered = data.data.filter((f) => {
        const fields = [f.src_ip, f.dst_ip, f.src_hostname, f.dst_hostname, f.src_service, f.dst_service].filter(Boolean)
        return fields.some((x) => x?.includes(src) || x?.includes(dst))
      })
      setFlows(filtered)
    },
  }

  return (
    <div className="space-y-4">
      <FilterBar filters={filters} setFilters={setFilters} compact />
      <SavedViewsPanel scope="sankey" filters={filters} setFilters={setFilters} />
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
