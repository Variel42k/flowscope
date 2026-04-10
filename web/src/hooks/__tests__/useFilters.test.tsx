import { renderHook, act } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { useGlobalFilters } from '../useFilters'

describe('useGlobalFilters', () => {
  it('serializes active filter params', () => {
    const { result } = renderHook(() => useGlobalFilters())
    act(() => {
      result.current.setFilters((prev) => ({
        ...prev,
        search: '10.10.1.10',
        exporter: 'lab-router-a',
        protocol: 'TCP',
        min_bytes: '1000',
      }))
    })
    const params = result.current.params
    expect(params.get('search')).toBe('10.10.1.10')
    expect(params.get('exporter')).toBe('lab-router-a')
    expect(params.get('protocol')).toBe('TCP')
    expect(params.get('min_bytes')).toBe('1000')
    expect(params.get('from')).toBeTruthy()
    expect(params.get('to')).toBeTruthy()
  })
})
