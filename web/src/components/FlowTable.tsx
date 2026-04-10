import { useMemo } from 'react'
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table'

import { FlowRecord } from '../lib/types'

const col = createColumnHelper<FlowRecord>()

export function FlowTable({
  rows,
  onSelect,
}: {
  rows: FlowRecord[]
  onSelect?: (row: FlowRecord) => void
}) {
  const columns = useMemo(
    () => [
      col.accessor('timestamp_end', { header: 'End', cell: (v) => new Date(v.getValue()).toLocaleString() }),
      col.accessor('src_ip', { header: 'Source' }),
      col.accessor('dst_ip', { header: 'Destination' }),
      col.accessor('l4_protocol_name', { header: 'Protocol' }),
      col.accessor('dst_port', { header: 'Dst Port' }),
      col.accessor('bytes', { header: 'Bytes' }),
      col.accessor('packets', { header: 'Packets' }),
      col.accessor('exporter_id', { header: 'Exporter' }),
      col.accessor('flow_direction', { header: 'Direction' }),
    ],
    [],
  )

  const table = useReactTable({ data: rows, columns, getCoreRowModel: getCoreRowModel() })

  return (
    <div className="overflow-x-auto rounded-lg border border-slate-700/60">
      <table className="w-full text-left text-sm">
        <thead className="bg-slate-800/80 text-slate-300">
          {table.getHeaderGroups().map((hg) => (
            <tr key={hg.id}>
              {hg.headers.map((h) => (
                <th key={h.id} className="px-2 py-2 font-medium">
                  {h.isPlaceholder ? null : flexRender(h.column.columnDef.header, h.getContext())}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody>
          {table.getRowModel().rows.map((row) => (
            <tr
              key={row.id}
              className="cursor-pointer border-t border-slate-700/40 bg-slate-950/40 hover:bg-slate-800/60"
              onClick={() => onSelect?.(row.original)}
            >
              {row.getVisibleCells().map((cell) => (
                <td key={cell.id} className="px-2 py-2">
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
