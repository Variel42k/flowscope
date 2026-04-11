import axios from 'axios'

type ErrorPayload = {
  error?: string
}

export function getErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const payload = error.response?.data as ErrorPayload | undefined
    if (payload && typeof payload.error === 'string' && payload.error.trim() !== '') {
      return payload.error
    }
    if (typeof error.message === 'string' && error.message.trim() !== '') {
      return error.message
    }
    return 'Request failed'
  }
  if (error instanceof Error && error.message.trim() !== '') {
    return error.message
  }
  return 'Unexpected error'
}
