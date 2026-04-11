import { Dispatch, SetStateAction, useCallback, useEffect, useMemo, useState } from 'react'

import { apiClient, getCurrentRole, getCurrentUser } from '../api/client'
import { GlobalFilters, makeDefaultFilters, normalizeSavedFilters } from '../hooks/useFilters'
import { getErrorMessage } from '../lib/http'
import { SavedView } from '../lib/types'

type Props = {
  scope: 'overview' | 'flows' | 'sankey' | 'map'
  filters: GlobalFilters
  setFilters: Dispatch<SetStateAction<GlobalFilters>>
}

export function SavedViewsPanel({ scope, filters, setFilters }: Props) {
  const [views, setViews] = useState<SavedView[]>([])
  const [selectedID, setSelectedID] = useState('')
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [isShared, setIsShared] = useState(false)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const role = getCurrentRole()
  const user = getCurrentUser()
  const isAdmin = role.toLowerCase() === 'admin'

  const selected = useMemo(() => views.find((v) => v.view_id === selectedID) ?? null, [views, selectedID])

  const loadViews = useCallback(async () => {
    const q = new URLSearchParams({ scope })
    const { data } = await apiClient.get<{ data: SavedView[] }>('/api/views?' + q.toString())
    setViews(data.data ?? [])
  }, [scope])

  useEffect(() => {
    void loadViews()
  }, [loadViews])

  useEffect(() => {
    if (!selected) return
    setName(selected.name)
    setDescription(selected.description ?? '')
    setIsShared(selected.is_shared)
  }, [selected])

  function applySelected() {
    if (!selected) return
    setFilters((prev) => ({ ...prev, ...normalizeSavedFilters(selected.filters) }))
  }

  function resetFilters() {
    setFilters(makeDefaultFilters())
  }

  async function saveNew() {
    setLoading(true)
    setError('')
    try {
      await apiClient.post('/api/views', {
        name: name.trim(),
        description: description.trim(),
        scope,
        is_shared: isShared && isAdmin,
        filters,
      })
      setName('')
      setDescription('')
      setIsShared(false)
      await loadViews()
    } catch (err: unknown) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  async function updateSelected() {
    if (!selected) return
    setLoading(true)
    setError('')
    try {
      await apiClient.put(`/api/views/${encodeURIComponent(selected.view_id)}`, {
        name: name.trim(),
        description: description.trim(),
        scope,
        is_shared: isShared && isAdmin,
        filters,
      })
      await loadViews()
    } catch (err: unknown) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  async function deleteSelected() {
    if (!selected) return
    setLoading(true)
    setError('')
    try {
      await apiClient.delete(`/api/views/${encodeURIComponent(selected.view_id)}`)
      setSelectedID('')
      await loadViews()
    } catch (err: unknown) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="panel mb-3 space-y-3 px-3 py-3">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div>
          <h3 className="text-sm font-semibold">Saved Views</h3>
          <p className="text-xs text-slate-400">Save and reuse filter presets for {scope}. Shared views are visible to all users.</p>
        </div>
        <span className="text-xs text-slate-400">
          User: {user} ({role})
        </span>
      </div>
      <div className="grid grid-cols-1 gap-2 md:grid-cols-6">
        <label className="md:col-span-2">
          <span className="label">Select View</span>
          <select value={selectedID} onChange={(e) => setSelectedID(e.target.value)}>
            <option value="">Choose...</option>
            {views.map((v) => (
              <option key={v.view_id} value={v.view_id}>
                {v.name} {v.is_shared ? '[shared]' : '[private]'}
              </option>
            ))}
          </select>
        </label>
        <label>
          <span className="label">Name</span>
          <input value={name} onChange={(e) => setName(e.target.value)} placeholder="SOC triage window" />
        </label>
        <label className="md:col-span-2">
          <span className="label">Description</span>
          <input value={description} onChange={(e) => setDescription(e.target.value)} placeholder="Optional notes" />
        </label>
        <label className="flex items-end gap-2 text-sm">
          <input
            type="checkbox"
            checked={isShared}
            disabled={!isAdmin}
            onChange={(e) => setIsShared(e.target.checked)}
          />
          <span className={isAdmin ? '' : 'text-slate-500'}>Shared</span>
        </label>
      </div>
      <div className="flex flex-wrap gap-2">
        <button className="secondary" onClick={applySelected} disabled={!selectedID}>Apply</button>
        <button className="secondary" onClick={resetFilters}>Reset Filters</button>
        <button className="primary" onClick={saveNew} disabled={loading || !name.trim()}>Save New</button>
        <button className="secondary" onClick={updateSelected} disabled={loading || !selectedID}>Update Selected</button>
        <button className="secondary" onClick={deleteSelected} disabled={loading || !selectedID}>Delete</button>
      </div>
      {error && <div className="rounded border border-red-500/50 bg-red-950/20 p-2 text-xs text-red-300">{error}</div>}
    </section>
  )
}
