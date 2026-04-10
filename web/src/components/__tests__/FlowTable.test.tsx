import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import { FlowTable } from '../FlowTable'

describe('FlowTable', () => {
  it('renders rows and triggers selection callback', () => {
    const onSelect = vi.fn()
    render(
      <FlowTable
        onSelect={onSelect}
        rows={[
          {
            flow_id: '1',
            timestamp_start: '2026-04-10T10:00:00Z',
            timestamp_end: '2026-04-10T10:00:05Z',
            exporter_id: 'exp-a',
            exporter_ip: '172.16.0.1',
            src_ip: '10.10.1.10',
            dst_ip: '10.10.2.20',
            src_port: 51000,
            dst_port: 443,
            ip_protocol: 6,
            l4_protocol_name: 'TCP',
            bytes: 1024,
            packets: 10,
            input_interface: 10,
            output_interface: 11,
            src_asn: 64512,
            dst_asn: 64513,
            src_country: 'US',
            dst_country: 'US',
            src_subnet: '10.10.1.0/24',
            dst_subnet: '10.10.2.0/24',
            src_is_private: true,
            dst_is_private: true,
            source_type: 'seed',
          },
        ]}
      />, 
    )

    expect(screen.getByText('10.10.1.10')).toBeInTheDocument()
    fireEvent.click(screen.getByText('10.10.1.10'))
    expect(onSelect).toHaveBeenCalledTimes(1)
  })
})
