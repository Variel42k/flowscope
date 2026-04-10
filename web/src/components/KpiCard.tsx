export function KpiCard({ title, value, subtitle }: { title: string; value: string; subtitle?: string }) {
  return (
    <div className="panel px-4 py-3">
      <div className="label">{title}</div>
      <div className="mt-2 text-2xl font-semibold">{value}</div>
      {subtitle && <div className="mt-1 text-xs text-slate-400">{subtitle}</div>}
    </div>
  )
}
