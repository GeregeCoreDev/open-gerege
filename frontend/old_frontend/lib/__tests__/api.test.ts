import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import api, { APIError, setUnauthorizedHandler } from '../api'

// Mock logout
vi.mock('../logout', () => ({
  logout: vi.fn(),
}))

// Mock toast
vi.mock('sonner', () => ({
  toast: {
    error: vi.fn(),
    success: vi.fn(),
  },
}))

const server = setupServer()

beforeEach(() => {
  server.listen({ onUnhandledRequest: 'error' })
})

afterEach(() => {
  server.resetHandlers()
  vi.clearAllMocks()
})

describe('APIError', () => {
  it('creates an error with status and info', () => {
    const error = new APIError('Test error', 404, { detail: 'Not found' })
    expect(error.message).toBe('Test error')
    expect(error.status).toBe(404)
    expect(error.info).toEqual({ detail: 'Not found' })
    expect(error.name).toBe('APIError')
  })
})

describe('api', () => {
  describe('GET requests', () => {
    it('makes a successful GET request', async () => {
      server.use(
        http.get('/api/users', () => {
          return HttpResponse.json({ data: [{ id: 1, name: 'Test User' }] })
        })
      )

      const result = await api.get<{ id: number; name: string }[]>('/users')
      expect(result).toEqual([{ id: 1, name: 'Test User' }])
    })

    it('handles query parameters', async () => {
      let capturedUrl = ''
      server.use(
        http.get('/api/users', ({ request }) => {
          capturedUrl = request.url
          return HttpResponse.json({ data: [] })
        })
      )

      await api.get('/users', { query: { page: 1, limit: 10 } })
      expect(capturedUrl).toContain('page=1')
      expect(capturedUrl).toContain('limit=10')
    })

    it('skips null/undefined query values', async () => {
      let capturedUrl = ''
      server.use(
        http.get('/api/users', ({ request }) => {
          capturedUrl = request.url
          return HttpResponse.json({ data: [] })
        })
      )

      await api.get('/users', { query: { page: 1, filter: null, sort: undefined } })
      expect(capturedUrl).toContain('page=1')
      expect(capturedUrl).not.toContain('filter')
      expect(capturedUrl).not.toContain('sort')
    })
  })

  describe('POST requests', () => {
    it('makes a successful POST request with JSON body', async () => {
      let capturedBody: unknown
      server.use(
        http.post('/api/users', async ({ request }) => {
          capturedBody = await request.json()
          return HttpResponse.json({ data: { id: 1, ...capturedBody as object } })
        })
      )

      const result = await api.post<{ id: number; name: string }>('/users', { name: 'New User' })
      expect(capturedBody).toEqual({ name: 'New User' })
      expect(result).toEqual({ id: 1, name: 'New User' })
    })

    it('sends correct content-type header for JSON', async () => {
      let capturedContentType: string | null = null
      server.use(
        http.post('/api/test', ({ request }) => {
          capturedContentType = request.headers.get('content-type')
          return HttpResponse.json({ data: {} })
        })
      )

      await api.post('/test', { key: 'value' })
      expect(capturedContentType).toBe('application/json')
    })
  })

  describe('PUT requests', () => {
    it('makes a successful PUT request', async () => {
      server.use(
        http.put('/api/users/1', () => {
          return HttpResponse.json({ data: { id: 1, name: 'Updated User' } })
        })
      )

      const result = await api.put<{ id: number; name: string }>('/users/1', { name: 'Updated User' })
      expect(result).toEqual({ id: 1, name: 'Updated User' })
    })
  })

  describe('PATCH requests', () => {
    it('makes a successful PATCH request', async () => {
      server.use(
        http.patch('/api/users/1', () => {
          return HttpResponse.json({ data: { id: 1, status: 'active' } })
        })
      )

      const result = await api.patch<{ id: number; status: string }>('/users/1', { status: 'active' })
      expect(result).toEqual({ id: 1, status: 'active' })
    })
  })

  describe('DELETE requests', () => {
    it('makes a successful DELETE request', async () => {
      server.use(
        http.delete('/api/users/1', () => {
          return HttpResponse.json({ data: { success: true } })
        })
      )

      const result = await api.del<{ success: boolean }>('/users/1')
      expect(result).toEqual({ success: true })
    })
  })

  describe('Error handling', () => {
    it('throws APIError on non-ok response', async () => {
      server.use(
        http.get('/api/error', () => {
          return HttpResponse.json(
            { message: 'Resource not found' },
            { status: 404 }
          )
        })
      )

      await expect(api.get('/error', { hasToast: false })).rejects.toThrow(APIError)
    })

    it('throws error on non-JSON response', async () => {
      server.use(
        http.get('/api/html', () => {
          return new HttpResponse('<html></html>', {
            headers: { 'content-type': 'text/html' },
          })
        })
      )

      await expect(api.get('/html', { hasToast: false })).rejects.toThrow('Expected JSON')
    })

    it('handles 401 unauthorized', async () => {
      const mockHandler = vi.fn()
      setUnauthorizedHandler(mockHandler)

      server.use(
        http.get('/api/protected', () => {
          return HttpResponse.json({}, { status: 401 })
        })
      )

      await expect(api.get('/protected', { hasToast: false })).rejects.toThrow('Unauthorized')
      expect(mockHandler).toHaveBeenCalled()

      // Reset handler
      setUnauthorizedHandler(null)
    })
  })

  describe('Timeout handling', () => {
    it('throws timeout error when request takes too long', async () => {
      server.use(
        http.get('/api/slow', async () => {
          await new Promise((resolve) => setTimeout(resolve, 200))
          return HttpResponse.json({ data: {} })
        })
      )

      await expect(
        api.get('/slow', { timeoutMs: 50, hasToast: false })
      ).rejects.toThrow('Request timeout')
    })
  })

  describe('Data unwrapping', () => {
    it('unwraps data by default', async () => {
      server.use(
        http.get('/api/wrapped', () => {
          return HttpResponse.json({ data: { value: 'test' }, meta: {} })
        })
      )

      const result = await api.get<{ value: string }>('/wrapped')
      expect(result).toEqual({ value: 'test' })
    })

    it('returns full payload when unwrapData is false', async () => {
      server.use(
        http.get('/api/wrapped', () => {
          return HttpResponse.json({ data: { value: 'test' }, meta: { total: 100 } })
        })
      )

      const result = await api.get<{ data: { value: string }; meta: { total: number } }>(
        '/wrapped',
        { unwrapData: false }
      )
      expect(result).toEqual({ data: { value: 'test' }, meta: { total: 100 } })
    })
  })

  describe('Callbacks', () => {
    it('calls onRequest before making request', async () => {
      const onRequest = vi.fn()
      server.use(
        http.get('/api/test', () => {
          return HttpResponse.json({ data: {} })
        })
      )

      await api.get('/test', { onRequest })
      expect(onRequest).toHaveBeenCalledWith(
        expect.stringContaining('/api/test'),
        expect.objectContaining({ method: 'GET' })
      )
    })

    it('calls onResponse after receiving response', async () => {
      const onResponse = vi.fn()
      server.use(
        http.get('/api/test', () => {
          return HttpResponse.json({ data: {} })
        })
      )

      await api.get('/test', { onResponse })
      expect(onResponse).toHaveBeenCalledWith(expect.any(Response))
    })

    it('calls onError when error occurs', async () => {
      const onError = vi.fn()
      server.use(
        http.get('/api/error', () => {
          return HttpResponse.json({ message: 'Error' }, { status: 500 })
        })
      )

      await expect(api.get('/error', { onError, hasToast: false })).rejects.toThrow()
      expect(onError).toHaveBeenCalled()
    })
  })

  describe('FormData handling', () => {
    it('handles FormData body without setting content-type', async () => {
      let capturedContentType: string | null = null
      server.use(
        http.post('/api/upload', ({ request }) => {
          capturedContentType = request.headers.get('content-type')
          return HttpResponse.json({ data: { success: true } })
        })
      )

      const formData = new FormData()
      formData.append('file', 'test')

      await api.post('/upload', formData)
      // FormData should not have explicit content-type (browser sets it with boundary)
      expect(capturedContentType).not.toBe('application/json')
    })
  })
})

describe('External URL handling', () => {
  it('uses full URL when path starts with http', async () => {
    server.use(
      http.get('https://external.api.com/data', () => {
        return HttpResponse.json({ data: { external: true } })
      })
    )

    const result = await api.get<{ external: boolean }>('https://external.api.com/data')
    expect(result).toEqual({ external: true })
  })
})
