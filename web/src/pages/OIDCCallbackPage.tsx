import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'

import { applyAuthFromQuery } from '../api/client'

export function OIDCCallbackPage({ onSuccess }: { onSuccess: () => void }) {
  const [params] = useSearchParams()
  const navigate = useNavigate()
  const [error, setError] = useState('')

  useEffect(() => {
    const query = '?' + params.toString()
    const err = params.get('error')
    if (err) {
      setError(err)
      return
    }
    const session = applyAuthFromQuery(query)
    if (!session) {
      setError('OIDC callback is missing token or user')
      return
    }
    onSuccess()
    navigate('/overview', { replace: true })
  }, [navigate, onSuccess, params])

  return (
    <div className="mx-auto mt-24 max-w-xl">
      <div className="panel space-y-3 px-6 py-6">
        <h2 className="text-lg font-semibold">OIDC Sign-in</h2>
        {error ? (
          <div className="rounded border border-red-500/50 bg-red-950/20 p-3 text-sm text-red-300">{error}</div>
        ) : (
          <p className="text-sm text-slate-300">Completing sign-in...</p>
        )}
      </div>
    </div>
  )
}
