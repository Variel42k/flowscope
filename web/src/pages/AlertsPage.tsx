import { useEffect, useMemo, useState } from 'react'

import { apiClient, getCurrentRole } from '../api/client'
import { formatBytes, formatNumber } from '../lib/format'
import { AlertEvent, AlertRule, PageResult } from '../lib/types'

const RULE_TYPES: Array<AlertRule['rule_type']> = ['new_edge', 'fanout_external', 'high_byte_edge', 'port_outlier']

export function AlertsPage() {
  const role = getCurrentRole()
  const isAdmin = role.toLowerCase() === 'admin'
  const [rules, setRules] = useState<AlertRule[]>([])
  const [events, setEvents] = useState<PageResult<AlertEvent> | null>(null)
  const [name, setName] = useState('')
  const [ruleType, setRuleType] = useState<AlertRule['rule_type']>('new_edge')
  const [severity, setSeverity] = useState<'low' | 'medium' | 'high'>('medium')
  const [windowMinutes, setWindowMinutes] = useState('15')
  const [threshold, setThreshold] = useState('1')
  const [enabled, setEnabled] = useState(true)
  const [editingID, setEditingID] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [severityFilter, setSeverityFilter] = useState('')

  async function refresh() {
    const q = new URLSearchParams({ page: '1', page_size: '100' })
    if (severityFilter) q.set('severity', severityFilter)
    const [rulesRes, eventsRes] = await Promise.all([
      apiClient.get<{ data: AlertRule[] }>('/api/alerts/rules'),
      apiClient.get<PageResult<AlertEvent>>('/api/alerts/events?' + q.toString()),
    ])
    setRules(rulesRes.data.data ?? [])
    setEvents(eventsRes.data)
  }

  useEffect(() => {
    refresh().catch(() => {})
  }, [severityFilter])

  const selectedRule = useMemo(() => rules.find((r) => r.rule_id === editingID) ?? null, [rules, editingID])

  useEffect(() => {
    if (!selectedRule) return
    setName(selectedRule.name)
    setRuleType(selectedRule.rule_type)
    setSeverity(selectedRule.severity)
    setWindowMinutes(String(selectedRule.window_minutes))
    setThreshold(String(selectedRule.threshold_value))
    setEnabled(selectedRule.enabled)
  }, [selectedRule])

  function clearForm() {
    setEditingID('')
    setName('')
    setRuleType('new_edge')
    setSeverity('medium')
    setWindowMinutes('15')
    setThreshold('1')
    setEnabled(true)
  }

  async function submit() {
    if (!isAdmin) return
    setError('')
    setLoading(true)
    try {
      const payload = {
        name: name.trim(),
        rule_type: ruleType,
        severity,
        window_minutes: Number(windowMinutes),
        threshold_value: Number(threshold),
        enabled,
      }
      if (editingID) {
        await apiClient.put(`/api/alerts/rules/${encodeURIComponent(editingID)}`, payload)
      } else {
        await apiClient.post('/api/alerts/rules', payload)
      }
      clearForm()
      await refresh()
    } catch (err: any) {
      setError(err?.response?.data?.error ?? err.message)
    } finally {
      setLoading(false)
    }
  }

  async function deleteRule(ruleID: string) {
    if (!isAdmin) return
    setError('')
    setLoading(true)
    try {
      await apiClient.delete(`/api/alerts/rules/${encodeURIComponent(ruleID)}`)
      if (editingID === ruleID) clearForm()
      await refresh()
    } catch (err: any) {
      setError(err?.response?.data?.error ?? err.message)
    } finally {
      setLoading(false)
    }
  }

  async function runEvaluate() {
    setError('')
    setLoading(true)
    try {
      await apiClient.post('/api/alerts/evaluate')
      await refresh()
    } catch (err: any) {
      setError(err?.response?.data?.error ?? err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-4">
      <section className="panel space-y-3 px-3 py-3">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div>
            <h2 className="text-sm font-semibold">Alert Rules</h2>
            <p className="text-xs text-slate-400">Heuristic detections: new edge, fanout to externals, high-byte edge, and port outliers.</p>
          </div>
          <button className="secondary" onClick={runEvaluate} disabled={loading}>Evaluate Now</button>
        </div>
        <div className="grid grid-cols-1 gap-2 md:grid-cols-6">
          <label>
            <span className="label">Name</span>
            <input value={name} onChange={(e) => setName(e.target.value)} placeholder="New Edge Observed" />
          </label>
          <label>
            <span className="label">Rule Type</span>
            <select value={ruleType} onChange={(e) => setRuleType(e.target.value as AlertRule['rule_type'])}>
              {RULE_TYPES.map((rt) => (
                <option key={rt} value={rt}>{rt}</option>
              ))}
            </select>
          </label>
          <label>
            <span className="label">Severity</span>
            <select value={severity} onChange={(e) => setSeverity(e.target.value as 'low' | 'medium' | 'high')}>
              <option value="low">low</option>
              <option value="medium">medium</option>
              <option value="high">high</option>
            </select>
          </label>
          <label>
            <span className="label">Window (min)</span>
            <input value={windowMinutes} onChange={(e) => setWindowMinutes(e.target.value)} />
          </label>
          <label>
            <span className="label">Threshold</span>
            <input value={threshold} onChange={(e) => setThreshold(e.target.value)} />
          </label>
          <label className="flex items-end gap-2 text-sm">
            <input type="checkbox" checked={enabled} onChange={(e) => setEnabled(e.target.checked)} />
            <span>Enabled</span>
          </label>
        </div>
        <div className="flex flex-wrap gap-2">
          <button className="primary" disabled={!isAdmin || loading || !name.trim()} onClick={submit}>
            {editingID ? 'Update Rule' : 'Create Rule'}
          </button>
          <button className="secondary" onClick={clearForm}>Clear</button>
          {!isAdmin && <span className="text-xs text-slate-400">Rule CRUD доступен только admin роли.</span>}
        </div>
        {error && <div className="rounded border border-red-500/50 bg-red-950/20 p-2 text-sm text-red-300">{error}</div>}
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="text-xs uppercase text-slate-400">
              <tr>
                <th className="pb-2">Name</th>
                <th className="pb-2">Type</th>
                <th className="pb-2">Severity</th>
                <th className="pb-2">Window</th>
                <th className="pb-2">Threshold</th>
                <th className="pb-2">Enabled</th>
                <th className="pb-2">Actions</th>
              </tr>
            </thead>
            <tbody>
              {rules.map((rule) => (
                <tr key={rule.rule_id} className="border-t border-slate-700/40">
                  <td className="py-2">{rule.name}</td>
                  <td className="py-2 font-mono text-xs">{rule.rule_type}</td>
                  <td className="py-2">{rule.severity}</td>
                  <td className="py-2">{rule.window_minutes}m</td>
                  <td className="py-2">{formatNumber(rule.threshold_value)}</td>
                  <td className="py-2">{rule.enabled ? 'yes' : 'no'}</td>
                  <td className="py-2">
                    <div className="flex gap-2">
                      <button className="secondary" onClick={() => setEditingID(rule.rule_id)}>Edit</button>
                      <button className="secondary" disabled={!isAdmin} onClick={() => deleteRule(rule.rule_id)}>Delete</button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <section className="panel space-y-3 px-3 py-3">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <h2 className="text-sm font-semibold">Alert Events</h2>
          <label className="text-xs text-slate-300">
            Severity:
            <select className="ml-2" value={severityFilter} onChange={(e) => setSeverityFilter(e.target.value)}>
              <option value="">all</option>
              <option value="low">low</option>
              <option value="medium">medium</option>
              <option value="high">high</option>
            </select>
          </label>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="text-xs uppercase text-slate-400">
              <tr>
                <th className="pb-2">Detected</th>
                <th className="pb-2">Rule</th>
                <th className="pb-2">Severity</th>
                <th className="pb-2">Description</th>
                <th className="pb-2">Bytes</th>
                <th className="pb-2">Flows</th>
              </tr>
            </thead>
            <tbody>
              {(events?.data ?? []).map((event) => (
                <tr key={event.event_id} className="border-t border-slate-700/40">
                  <td className="py-2">{new Date(event.detected_at).toLocaleString()}</td>
                  <td className="py-2">{event.rule_name}</td>
                  <td className="py-2">{event.severity}</td>
                  <td className="py-2">{event.description}</td>
                  <td className="py-2">{formatBytes(event.bytes)}</td>
                  <td className="py-2">{formatNumber(event.flows)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  )
}
