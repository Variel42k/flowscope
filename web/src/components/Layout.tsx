import { NavLink } from 'react-router-dom'

import { logout } from '../api/client'

type Props = {
  children: React.ReactNode
  onLoggedOut: () => void
}

export function Layout({ children, onLoggedOut }: Props) {
  const user = localStorage.getItem('flowscope_user') ?? 'admin'
  const role = localStorage.getItem('flowscope_role') ?? 'viewer'

  return (
    <div className="mx-auto flex min-h-screen w-full max-w-[1700px] flex-col px-3 pb-6 pt-4 md:px-6">
      <header className="panel mb-4 flex flex-col items-start justify-between gap-3 px-4 py-3 md:flex-row md:items-center">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">FlowScope</h1>
          <p className="text-sm text-slate-400">Network Flow Observability Platform</p>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-xs text-slate-300">Signed in as {user} ({role})</span>
          <button
            className="secondary"
            onClick={() => {
              logout()
              onLoggedOut()
            }}
          >
            Logout
          </button>
        </div>
      </header>
      <nav className="panel mb-4 flex flex-wrap gap-2 px-2 py-2">
        {[
          ['/overview', 'Overview'],
          ['/flows', 'Flows'],
          ['/sankey', 'Sankey'],
          ['/map', 'Interaction Map'],
          ['/alerts', 'Alerts'],
          ['/views', 'Saved Views'],
        ].map(([to, label]) => (
          <NavLink
            key={to}
            to={to}
            className={({ isActive }) =>
              `rounded-lg px-3 py-2 text-sm ${isActive ? 'bg-accent text-slate-950' : 'bg-slate-800 text-slate-200 hover:bg-slate-700'}`
            }
          >
            {label}
          </NavLink>
        ))}
      </nav>
      <main className="flex-1">{children}</main>
    </div>
  )
}
