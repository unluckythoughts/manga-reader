const BASE_URL = '/api/v1'

// Shape of the go-microservice v2 response envelope
interface Envelope<T> {
  ok: boolean
  id: string
  error?: string
  data?: T
}

function getAuthHeaders(): Record<string, string> {
  const token = localStorage.getItem('auth_token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders()
    },
    body: body !== undefined ? JSON.stringify(body) : undefined
  })

  const envelope = await res.json() as Envelope<T>

  if (!res.ok || !envelope.ok) {
    throw new Error(envelope.error || `HTTP ${res.status}`)
  }

  return envelope.data as T
}

export const apiClient = {
  get: <T>(path: string) => request<T>('GET', path),
  post: <T>(path: string, body: unknown) => request<T>('POST', path, body),
  put: <T>(path: string, body: unknown) => request<T>('PUT', path, body),
  patch: <T>(path: string, body?: unknown) => request<T>('PATCH', path, body),
  delete: <T>(path: string) => request<T>('DELETE', path)
}
