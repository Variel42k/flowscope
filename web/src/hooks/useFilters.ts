import { useMemo, useState } from 'react'

export type GlobalFilters = {
  from: string
  to: string
  search: string
  exporter: string
  protocol: string
  src_ip: string
  dst_ip: string
  subnet: string
  src_port: string
  dst_port: string
  asn: string
  country: string
  min_bytes: string
  min_flows: string
}

function defaultFromISO() {
  const date = new Date(Date.now() - 60 * 60 * 1000)
  return date.toISOString().slice(0, 19)
}

function defaultToISO() {
  return new Date().toISOString().slice(0, 19)
}

export const FILTER_KEYS: Array<keyof GlobalFilters> = [
  'from',
  'to',
  'search',
  'exporter',
  'protocol',
  'src_ip',
  'dst_ip',
  'subnet',
  'src_port',
  'dst_port',
  'asn',
  'country',
  'min_bytes',
  'min_flows',
]

export function makeDefaultFilters(): GlobalFilters {
  return {
    from: defaultFromISO(),
    to: defaultToISO(),
    search: '',
    exporter: '',
    protocol: '',
    src_ip: '',
    dst_ip: '',
    subnet: '',
    src_port: '',
    dst_port: '',
    asn: '',
    country: '',
    min_bytes: '0',
    min_flows: '0',
  }
}

export function normalizeSavedFilters(input: Record<string, string>): Partial<GlobalFilters> {
  const out: Partial<GlobalFilters> = {}
  for (const key of FILTER_KEYS) {
    const value = input[key]
    if (typeof value === 'string') {
      out[key] = value
    }
  }
  return out
}

export function useGlobalFilters() {
  const [filters, setFilters] = useState<GlobalFilters>(makeDefaultFilters())

  const params = useMemo(() => {
    const p = new URLSearchParams()
    Object.entries(filters).forEach(([k, v]) => {
      const trimmed = v.trim()
      if (!trimmed || trimmed === '0') return
      if (k === 'from' || k === 'to') {
        p.set(k, `${trimmed}Z`)
      } else {
        p.set(k, trimmed)
      }
    })
    return p
  }, [filters])

  return { filters, setFilters, params }
}
