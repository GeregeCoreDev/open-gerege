import { logout } from './logout'
import { toast } from 'sonner'

/**
 * 🧱 APIError class
 * Represents an error that occurred during an API request.
 * @property status - HTTP status code
 * @property info - Additional error information (optional)
 */
export class APIError extends Error {
  status: number
  info?: unknown
  constructor(message: string, status: number, info?: unknown) {
    super(message)
    this.name = 'APIError'
    this.status = status
    this.info = info
  }
}

/**
 * 🧩 onUnauthorized handler
 * Callback function to execution when a 401 error occurs
 */
let onUnauthorized: null | (() => Promise<void> | void) = null

/**
 * 🔧 Set unauthorized handler
 * @param h - function to run on unauthorized
 */
export const setUnauthorizedHandler = (h: typeof onUnauthorized) => {
  onUnauthorized = h
}

/**
 * 🔍 Query parameter type
 * Used to generate query strings for API requests
 */
export type Query = Record<string, string | number | boolean | null | undefined>

/**
 * 🧮 Query string generator
 * @param query - Query object
 * @returns URLSearchParams formatted query string
 */
function qs(query?: Query) {
  if (!query) return ''
  const p = new URLSearchParams()
  for (const [k, v] of Object.entries(query)) {
    if (v == null) continue
    p.set(k, String(v))
  }
  const s = p.toString()
  return s ? `?${s}` : ''
}

type Method = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
type Body = BodyInit | Record<string, unknown> | undefined

/**
 * ⚙️ API Options
 * Configuration object for API calls
 */
export type ApiOptions = {
  method?: Method
  query?: Query
  headers?: HeadersInit
  body?: Body
  baseURL?: string
  timeoutMs?: number
  cache?: RequestCache
  signal?: AbortSignal
  onRequest?(url: string, init: RequestInit): void | Promise<void>
  onResponse?(res: Response): void | Promise<void>
  onError?(err: unknown): void | Promise<void>
  unwrapData?: boolean
  hasToast?: boolean | 'success' | 'error'
}

/**
 * 📏 Remove trailing slash from base URL
 */
function normalizeBase(base: string) {
  return base.replace(/\/+$/, '')
}

/**
 * 🔗 Helper function to join Base and path
 */
function joinURL(base: string, path: string) {
  const b = normalizeBase(base)
  const p = path.startsWith('/') ? path : `/${path}`
  return `${b}${p}`
}

/**
 * ⏳ Promise wrapper with timeout
 * @param p - async promise
 * @param ms - duration (milliseconds)
 */
async function withTimeout<T>(p: Promise<T>, ms: number) {
  return new Promise<T>((resolve, reject) => {
    const id = setTimeout(() => reject(new APIError('Request timeout', 408)), ms)
    p.then(
      (v) => {
        clearTimeout(id)
        resolve(v)
      },
      (e) => {
        clearTimeout(id)
        reject(e)
      },
    )
  })
}

/**
 * 🧰 Normalize request body to standard format
 * @param body - Request body
 * @param headers - Header configuration
 */
function normalizeBody(body: Body, headers: HeadersInit = {}) {
  if (!body) return { body: undefined as BodyInit | undefined, headers }
  if (
    typeof body === 'string' ||
    (typeof Blob !== 'undefined' && body instanceof Blob) ||
    (typeof FormData !== 'undefined' && body instanceof FormData) ||
    (typeof URLSearchParams !== 'undefined' && body instanceof URLSearchParams)
  ) {
    return { body: body as BodyInit, headers }
  }
  return {
    body: JSON.stringify(body),
    headers: { 'content-type': 'application/json', ...headers },
  }
}

/**
 * 🚀 General API request function
 * @template T - Return data type
 * @param path - API path
 * @param opts - Configuration parameters
 * @returns Promise<T>
 */
export async function api<T>(path: string, opts: ApiOptions = {}): Promise<T> {
  // Always use /api proxy route to hide base URL
  const baseEnv = '/api'
  const base = normalizeBase(opts.baseURL ?? baseEnv)

  const url = path.startsWith('http') ? path : joinURL(base, path) + qs(opts.query)

  const { body, headers } = normalizeBody(opts.body, opts.headers)

  let authHeader: string | undefined
  if (process.env.NODE_ENV !== 'production') {
    const sid = process.env.NEXT_PUBLIC_DEV_SID
    const h = headers as Record<string, string> | undefined
    if (sid && !h?.authorization && !h?.Authorization) {
      authHeader = sid
    }
  }

  const init: RequestInit = {
    method: opts.method ?? 'GET',
    headers: {
      accept: 'application/json',
      ...(authHeader ? { authorization: 'bearer ' + authHeader } : {}),
      ...headers,
    },
    body,
    cache: opts.cache ?? 'no-store',
    signal: opts.signal,
    credentials: 'include',
    next: { revalidate: 0 },
  }

  try {
    // 🔹 Hook before sending request
    await opts.onRequest?.(url, init)

    // 🔹 Fetch with timeout
    const res = await withTimeout(fetch(url, init), opts.timeoutMs ?? 15000)
    await opts.onResponse?.(res)

    const contentType = res.headers.get('content-type') || ''
    const isJSON =
      contentType.includes('application/json') || contentType.includes('application/problem+json')

    // ⚠️ Handle Unauthorized (401) - call logout or custom handler
    if (res.status === 401) {
      try {
        if (onUnauthorized) await onUnauthorized()
        else await logout()
      } finally {
        throw new Error('Unauthorized')
      }
    }

    // ❗ If HTML is returned, treat as middleware or route error
    if (!isJSON) {
      const preview = (await res.text().catch(() => '')).slice(0, 300)
      throw new APIError(
        `Expected JSON but got "${contentType}". Check that requests go to ${base} and that middleware skips /api.`,
        res.status || 500,
        preview,
      )
    }

    // 📦 Parse JSON payload
    const payload = await res.json().catch(() => ({}))

    // ❌ Throw error if Response is not OK
    if (!res.ok) {
      const errorPayload = payload as { message?: string }
      const msg = errorPayload?.message || res.statusText
      throw new APIError(msg, res.status, payload)
    }

    // ✅ Data unwrap configuration
    const unwrap = opts.unwrapData ?? true
    if (unwrap && payload && typeof payload === 'object' && 'data' in payload) {
      return (payload as { data: T }).data
    }
    return payload as T
  } catch (err) {
    // ⚠️ Show toast on error
    const showToast = opts.hasToast ?? true
    if (showToast) {
      let message = 'Something went wrong'
      if (err instanceof APIError) message = err.message
      else if (err instanceof Error) message = err.message

      // Always show error toast on failures (avoid "success" toast in catch)
      toast.error(message)
    }
    await opts.onError?.(err)
    throw err
  }
}

/**
 * 🧩 Helper shortcut functions (for each method)
 */
api.get = <T>(path: string, o?: Omit<ApiOptions, 'method' | 'body'>) =>
  api<T>(path, { ...o, method: 'GET' })

api.post = <T>(path: string, body?: Body, o?: Omit<ApiOptions, 'method'>) =>
  api<T>(path, { ...o, method: 'POST', body })

api.put = <T>(path: string, body?: Body, o?: Omit<ApiOptions, 'method'>) =>
  api<T>(path, { ...o, method: 'PUT', body })

api.patch = <T>(path: string, body?: Body, o?: Omit<ApiOptions, 'method'>) =>
  api<T>(path, { ...o, method: 'PATCH', body })

api.del = <T>(path: string, body?: Body, o?: Omit<ApiOptions, 'method' | 'body'>) =>
  api<T>(path, { ...o, method: 'DELETE', body })

export default api
