import { FormEvent, useState } from 'react'

import { login, startOIDCLogin } from '../api/client'
import { getErrorMessage } from '../lib/http'

export function LoginForm({ onSuccess }: { onSuccess: () => void }) {
  const [username, setUsername] = useState('admin')
  const [password, setPassword] = useState('admin123')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const oidcEnabled = String(import.meta.env.VITE_OIDC_ENABLED ?? 'false').toLowerCase() === 'true'

  async function submit(e: FormEvent) {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      await login(username, password)
      onSuccess()
    } catch (err: unknown) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="mx-auto mt-24 max-w-md">
      <form className="panel space-y-4 px-6 py-6" onSubmit={submit}>
        <h2 className="text-xl font-semibold">Login to FlowScope</h2>
        <p className="text-sm text-slate-400">Use your local admin credentials from environment variables.</p>
        <label className="block">
          <span className="label">Username</span>
          <input className="mt-1 w-full" value={username} onChange={(e) => setUsername(e.target.value)} />
        </label>
        <label className="block">
          <span className="label">Password</span>
          <input className="mt-1 w-full" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        </label>
        {error && <div className="rounded border border-red-500/60 bg-red-950/40 p-2 text-sm text-red-300">{error}</div>}
        <button className="primary w-full" disabled={loading} type="submit">
          {loading ? 'Signing in...' : 'Sign in'}
        </button>
        {oidcEnabled && (
          <button
            className="secondary w-full"
            disabled={loading}
            type="button"
            onClick={async () => {
              setError('')
              setLoading(true)
              try {
                await startOIDCLogin()
              } catch (err: unknown) {
                setError(getErrorMessage(err))
                setLoading(false)
              }
            }}
          >
            Sign in with OIDC
          </button>
        )}
      </form>
    </div>
  )
}
