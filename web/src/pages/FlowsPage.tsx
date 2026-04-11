import { useEffect, useState } from 'react'

import { apiClient } from '../api/client'
import { FilterBar } from '../components/FilterBar'
import { FlowTable } from '../components/FlowTable'
import { SavedViewsPanel } from '../components/SavedViewsPanel'
import { useGlobalFilters } from '../hooks/useFilters'
import { formatBytes, formatNumber } from '../lib/format'
import { FlowRecord, PageResult } from '../lib/types'

export function FlowsPage() {
  const { filters, setFilters, params } = useGlobalFilters()
  const [tab, setTab] = useState<'active' | 'historical'>('active')
  const [page, setPage] = useState(1)
  const [data, setData] = useState<PageResult<FlowRecord> | null>(null)
  const [selected, setSelected] = useState<FlowRecord | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    setPage(1)
  }, [params, tab])

  useEffect(() => {
    let active = true
    setLoading(true)
    const q = new URLSearchParams(params)
    q.set('page', String(page))
    q.set('page_size', '50')
    q.set('sort_by', 'timestamp_end')
    q.set('sort_dir', 'DESC')
    apiClient
      .get(`/api/flows/${tab}?` + q.toString())
      .then((res) => {
        if (!active) return
        setData(res.data)
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => {
      active = false
    }
  }, [params, page, tab])

  return (
    <div className="space-y-4">
      <FilterBar filters={filters} setFilters={setFilters} />
      <SavedViewsPanel scope="flows" filters={filters} setFilters={setFilters} />
      <div className="panel flex items-center justify-between px-3 py-2">
        <div className="flex gap-2">
          <button className={tab === 'active' ? 'primary' : 'secondary'} onClick={() => setTab('active')}>Active Flows</button>
          <button className={tab === 'historical' ? 'primary' : 'secondary'} onClick={() => setTab('historical')}>Historical Flows</button>
        </div>
        <div className="text-xs text-slate-400">Page {data?.page ?? 1} | Total {data?.total ?? 0}</div>
      </div>

      <div className="grid grid-cols-1 gap-4 xl:grid-cols-[minmax(0,1fr)_360px]">
        <div className="panel p-3">
          <FlowTable rows={data?.data ?? []} onSelect={setSelected} />
          <div className="mt-3 flex items-center justify-between text-sm">
            <button className="secondary" onClick={() => setPage((p) => Math.max(1, p - 1))}>Previous</button>
            <span className="text-slate-400">{loading ? 'Loading...' : `Rows: ${data?.data.length ?? 0}`}</span>
            <button className="secondary" disabled={!data?.has_next} onClick={() => setPage((p) => p + 1)}>Next</button>
          </div>
        </div>

        <aside className="panel space-y-2 p-3">
          <h3 className="text-sm font-semibold">Flow Details</h3>
          {selected ? (
            <div className="space-y-2 text-sm">
              <Row label="Flow ID" value={selected.flow_id} mono />
              <Row label="Time" value={`${new Date(selected.timestamp_start).toLocaleString()} - ${new Date(selected.timestamp_end).toLocaleString()}`} />
              <Row label="Source" value={`${selected.src_ip}:${selected.src_port}`} mono />
              <Row label="Destination" value={`${selected.dst_ip}:${selected.dst_port}`} mono />
              <Row label="Protocol" value={`${selected.l4_protocol_name} (${selected.ip_protocol})`} />
              <Row label="Bytes" value={formatBytes(selected.bytes)} />
              <Row label="Packets" value={formatNumber(selected.packets)} />
              <Row label="Exporter" value={selected.exporter_id} />
              <Row label="Interfaces" value={`${selected.input_interface_alias || selected.input_interface} -> ${selected.output_interface_alias || selected.output_interface}`} />
              <Row label="Direction" value={selected.flow_direction ?? 'unknown'} />
              <Row label="ASN" value={`${selected.src_asn} -> ${selected.dst_asn}`} />
              <Row label="Country" value={`${selected.src_country} -> ${selected.dst_country}`} />
              <Row label="Private/Public" value={`${selected.src_is_private ? 'private' : 'public'} -> ${selected.dst_is_private ? 'private' : 'public'}`} />
              <Row label="Source Type" value={selected.source_type} />
            </div>
          ) : (
            <p className="text-sm text-slate-400">Click a flow row to inspect details.</p>
          )}
        </aside>
      </div>
    </div>
  )
}

function Row({ label, value, mono }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="rounded border border-slate-700/40 bg-slate-900/40 px-2 py-1">
      <div className="label">{label}</div>
      <div className={mono ? 'value break-all' : 'break-all text-slate-200'}>{value}</div>
    </div>
  )
}
