import { Dispatch, SetStateAction } from 'react'

import { GlobalFilters } from '../hooks/useFilters'

type Props = {
  filters: GlobalFilters
  setFilters: Dispatch<SetStateAction<GlobalFilters>>
  compact?: boolean
}

export function FilterBar({ filters, setFilters, compact = false }: Props) {
  const cls = compact ? 'grid-cols-2 md:grid-cols-4' : 'grid-cols-2 md:grid-cols-6'

  return (
    <section className="panel mb-4 px-3 py-3">
      <div className={`grid gap-2 ${cls}`}>
        <Field label="From">
          <input
            type="datetime-local"
            value={filters.from}
            onChange={(e) => setFilters((s) => ({ ...s, from: e.target.value }))}
          />
        </Field>
        <Field label="To">
          <input
            type="datetime-local"
            value={filters.to}
            onChange={(e) => setFilters((s) => ({ ...s, to: e.target.value }))}
          />
        </Field>
        <Field label="Search">
          <input value={filters.search} onChange={(e) => setFilters((s) => ({ ...s, search: e.target.value }))} placeholder="ip/hostname/service" />
        </Field>
        <Field label="Exporter">
          <input value={filters.exporter} onChange={(e) => setFilters((s) => ({ ...s, exporter: e.target.value }))} placeholder="lab-router-a" />
        </Field>
        <Field label="Protocol">
          <select value={filters.protocol} onChange={(e) => setFilters((s) => ({ ...s, protocol: e.target.value }))}>
            <option value="">Any</option>
            <option value="TCP">TCP</option>
            <option value="UDP">UDP</option>
            <option value="ICMP">ICMP</option>
          </select>
        </Field>
        <Field label="Subnet">
          <input value={filters.subnet} onChange={(e) => setFilters((s) => ({ ...s, subnet: e.target.value }))} placeholder="10.10.1.0/24" />
        </Field>
        <Field label="Source IP">
          <input value={filters.src_ip} onChange={(e) => setFilters((s) => ({ ...s, src_ip: e.target.value }))} placeholder="10.10.1.10" />
        </Field>
        <Field label="Destination IP">
          <input value={filters.dst_ip} onChange={(e) => setFilters((s) => ({ ...s, dst_ip: e.target.value }))} placeholder="198.51.100.12" />
        </Field>
        <Field label="Source Port">
          <input value={filters.src_port} onChange={(e) => setFilters((s) => ({ ...s, src_port: e.target.value }))} />
        </Field>
        <Field label="Destination Port">
          <input value={filters.dst_port} onChange={(e) => setFilters((s) => ({ ...s, dst_port: e.target.value }))} />
        </Field>
        <Field label="ASN">
          <input value={filters.asn} onChange={(e) => setFilters((s) => ({ ...s, asn: e.target.value }))} />
        </Field>
        <Field label="Country">
          <input value={filters.country} onChange={(e) => setFilters((s) => ({ ...s, country: e.target.value.toUpperCase() }))} placeholder="US" />
        </Field>
      </div>
    </section>
  )
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <label className="flex flex-col gap-1">
      <span className="label">{label}</span>
      {children}
    </label>
  )
}
