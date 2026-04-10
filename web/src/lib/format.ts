export function formatBytes(n: number) {
  if (!Number.isFinite(n)) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let value = n
  let i = 0
  while (value >= 1024 && i < units.length - 1) {
    value /= 1024
    i++
  }
  return `${value.toFixed(i < 2 ? 0 : 2)} ${units[i]}`
}

export function formatNumber(n: number) {
  return new Intl.NumberFormat().format(n)
}
