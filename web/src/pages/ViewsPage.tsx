import { useEffect, useState } from 'react'

import { apiClient, getCurrentRole, getCurrentUser } from '../api/client'
import { SavedView } from '../lib/types'

export function ViewsPage() {
  const [views, setViews] = useState<SavedView[]>([])
  const [scope, setScope] = useState('')
  const [error, setError] = useState('')
  const user = getCurrentUser()
  const role = getCurrentRole()
  const isAdmin = role.toLowerCase() === 'admin'

  async function load() {
    try {
      setError('')
      const q = new URLSearchParams()
      if (scope) q.set('scope', scope)
      const { data } = await apiClient.get<{ data: SavedView[] }>('/api/views?' + q.toString())
      setViews(data.data ?? [])
    } catch (err: any) {
      setError(err?.response?.data?.error ?? err.message)
    }
  }

  useEffect(() => {
    load().catch(() => {})
  }, [scope])

  async function deleteView(viewID: string) {
    try {
      setError('')
      await apiClient.delete(`/api/views/${encodeURIComponent(viewID)}`)
      await load()
    } catch (err: any) {
      setError(err?.response?.data?.error ?? err.message)
    }
  }

  return (
    <div className="space-y-4">
      <section className="panel space-y-3 px-3 py-3">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div>
            <h2 className="text-sm font-semibold">Saved Views Registry</h2>
            <p className="text-xs text-slate-400">Inventory of private/shared filter presets. Use page-level "Saved Views" blocks to apply quickly.</p>
          </div>
          <div className="text-xs text-slate-400">Signed in as {user} ({role})</div>
        </div>
        <label className="inline-flex items-center gap-2 text-sm">
          <span className="label">Scope</span>
          <select value={scope} onChange={(e) => setScope(e.target.value)}>
            <option value="">all</option>
            <option value="overview">overview</option>
            <option value="flows">flows</option>
            <option value="sankey">sankey</option>
            <option value="map">map</option>
            <option value="global">global</option>
          </select>
        </label>
        {error && <div className="rounded border border-red-500/50 bg-red-950/20 p-2 text-sm text-red-300">{error}</div>}
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="text-xs uppercase text-slate-400">
              <tr>
                <th className="pb-2">Name</th>
                <th className="pb-2">Scope</th>
                <th className="pb-2">Owner</th>
                <th className="pb-2">Shared</th>
                <th className="pb-2">Filters</th>
                <th className="pb-2">Updated</th>
                <th className="pb-2">Actions</th>
              </tr>
            </thead>
            <tbody>
              {views.map((view) => (
                <tr key={view.view_id} className="border-t border-slate-700/40">
                  <td className="py-2">{view.name}</td>
                  <td className="py-2">{view.scope}</td>
                  <td className="py-2">{view.owner_user}</td>
                  <td className="py-2">{view.is_shared ? 'yes' : 'no'}</td>
                  <td className="py-2">
                    <code className="text-xs">{Object.keys(view.filters ?? {}).length} keys</code>
                  </td>
                  <td className="py-2">{new Date(view.updated_at).toLocaleString()}</td>
                  <td className="py-2">
                    <button
                      className="secondary"
                      disabled={!isAdmin && view.owner_user !== user}
                      onClick={() => deleteView(view.view_id)}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  )
}
