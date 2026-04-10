import { useEffect, useState } from 'react'
import { Navigate, Route, Routes } from 'react-router-dom'

import { isLoggedIn } from './api/client'
import { Layout } from './components/Layout'
import { LoginForm } from './components/LoginForm'
import { FlowsPage } from './pages/FlowsPage'
import { MapPage } from './pages/MapPage'
import { OverviewPage } from './pages/OverviewPage'
import { SankeyPage } from './pages/SankeyPage'

export function App() {
  const [authed, setAuthed] = useState(isLoggedIn())

  useEffect(() => {
    const timer = window.setInterval(() => {
      if (!isLoggedIn()) {
        setAuthed(false)
      }
    }, 1500)
    return () => window.clearInterval(timer)
  }, [])

  if (!authed) {
    return <LoginForm onSuccess={() => setAuthed(true)} />
  }

  return (
    <Layout onLoggedOut={() => setAuthed(false)}>
      <Routes>
        <Route path="/" element={<Navigate to="/overview" replace />} />
        <Route path="/overview" element={<OverviewPage />} />
        <Route path="/flows" element={<FlowsPage />} />
        <Route path="/sankey" element={<SankeyPage />} />
        <Route path="/map" element={<MapPage />} />
      </Routes>
    </Layout>
  )
}
