'use client'

import { useState, useEffect, useCallback, useRef } from 'react'
import api, { Query, ApiOptions } from '@/lib/api'
import { cache, CACHE_TTL } from '@/lib/cache'

type QueryStatus = 'idle' | 'loading' | 'success' | 'error'

// Shared in-flight requests map for deduplication
const globalInFlightRequests = new Map<string, Promise<unknown>>()

type UseApiQueryOptions<T> = {
  /** API endpoint path */
  path: string
  /** Query parameters */
  query?: Query
  /** Enable/disable the query */
  enabled?: boolean
  /** Cache key (defaults to path + query) */
  cacheKey?: string
  /** Cache TTL in milliseconds */
  cacheTtl?: number
  /** Initial data */
  initialData?: T
  /** Callback on success */
  onSuccess?: (data: T) => void
  /** Callback on error */
  onError?: (error: Error) => void
  /** Refetch interval in milliseconds (0 to disable) */
  refetchInterval?: number
  /** Show toast on error */
  hasToast?: boolean
}

type UseApiQueryResult<T> = {
  data: T | undefined
  status: QueryStatus
  isLoading: boolean
  isSuccess: boolean
  isError: boolean
  error: Error | null
  refetch: () => Promise<void>
  invalidate: () => void
}

/**
 * Custom hook for fetching API data with caching support
 *
 * @example
 * const { data, isLoading, error, refetch } = useApiQuery<User[]>({
 *   path: '/users',
 *   query: { page: 1 },
 *   cacheTtl: CACHE_TTL.MEDIUM,
 * })
 */
export function useApiQuery<T>(options: UseApiQueryOptions<T>): UseApiQueryResult<T> {
  const {
    path,
    query,
    enabled = true,
    cacheKey,
    cacheTtl = CACHE_TTL.MEDIUM,
    initialData,
    onSuccess,
    onError,
    refetchInterval = 0,
    hasToast = false,
  } = options

  const [data, setData] = useState<T | undefined>(initialData)
  const [status, setStatus] = useState<QueryStatus>('idle')
  const [error, setError] = useState<Error | null>(null)

  // Generate cache key from path and query
  const getCacheKey = useCallback(() => {
    if (cacheKey) return cacheKey
    const queryString = query ? JSON.stringify(query) : ''
    return `api:${path}:${queryString}`
  }, [path, query, cacheKey])

  const isMountedRef = useRef(true)

  const fetchData = useCallback(async () => {
    if (!enabled) return

    const key = getCacheKey()

    // Check cache first
    const cached = cache.get<T>(key)
    if (cached !== null) {
      setData(cached)
      setStatus('success')
      return
    }

    setStatus('loading')
    setError(null)

    try {
      // Check for in-flight request
      if (globalInFlightRequests.has(key)) {
        const promise = globalInFlightRequests.get(key) as Promise<T>
        const result = await promise
        if (isMountedRef.current) {
          setData(result)
          setStatus('success')
          onSuccess?.(result)
        }
        return
      }

      const promise = api.get<T>(path, { query, hasToast })
      globalInFlightRequests.set(key, promise)

      const result = await promise

      if (isMountedRef.current) {
        setData(result)
        setStatus('success')
        cache.set(key, result, cacheTtl)
        onSuccess?.(result)
      }
    } catch (err) {
      if (isMountedRef.current) {
        const error = err instanceof Error ? err : new Error('Unknown error')
        setError(error)
        setStatus('error')
        onError?.(error)
      }
    } finally {
      globalInFlightRequests.delete(key)
    }
  }, [enabled, getCacheKey, path, query, hasToast, cacheTtl, onSuccess, onError])

  const invalidate = useCallback(() => {
    const key = getCacheKey()
    cache.invalidate(key)
  }, [getCacheKey])

  const refetch = useCallback(async () => {
    invalidate()
    await fetchData()
  }, [invalidate, fetchData])

  // Initial fetch
  useEffect(() => {
    isMountedRef.current = true
    fetchData()

    return () => {
      isMountedRef.current = false
    }
  }, [fetchData])

  // Refetch interval
  useEffect(() => {
    if (refetchInterval <= 0 || !enabled) return

    const intervalId = setInterval(refetch, refetchInterval)
    return () => clearInterval(intervalId)
  }, [refetchInterval, enabled, refetch])

  return {
    data,
    status,
    isLoading: status === 'loading',
    isSuccess: status === 'success',
    isError: status === 'error',
    error,
    refetch,
    invalidate,
  }
}

/**
 * Custom hook for mutating API data (POST, PUT, PATCH, DELETE)
 *
 * @example
 * const { mutate, isLoading } = useApiMutation<User, CreateUserInput>({
 *   path: '/users',
 *   method: 'POST',
 *   onSuccess: (data) => console.log('User created:', data),
 * })
 *
 * // Usage: mutate({ name: 'John' })
 */
type UseMutationOptions<TData, TVariables> = {
  path: string
  method?: 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  onSuccess?: (data: TData, variables: TVariables) => void
  onError?: (error: Error, variables: TVariables) => void
  /** Cache keys to invalidate on success */
  invalidateKeys?: string[]
  hasToast?: boolean
}

type UseMutationResult<TData, TVariables> = {
  mutate: (variables: TVariables) => Promise<TData | undefined>
  data: TData | undefined
  isLoading: boolean
  isSuccess: boolean
  isError: boolean
  error: Error | null
  reset: () => void
}

export function useApiMutation<TData, TVariables = unknown>(
  options: UseMutationOptions<TData, TVariables>
): UseMutationResult<TData, TVariables> {
  const {
    path,
    method = 'POST',
    onSuccess,
    onError,
    invalidateKeys = [],
    hasToast = true,
  } = options

  const [data, setData] = useState<TData | undefined>(undefined)
  const [status, setStatus] = useState<QueryStatus>('idle')
  const [error, setError] = useState<Error | null>(null)

  const mutate = useCallback(
    async (variables: TVariables): Promise<TData | undefined> => {
      setStatus('loading')
      setError(null)

      try {
        const apiOptions: ApiOptions = { method, hasToast }
        let result: TData

        switch (method) {
          case 'POST':
            result = await api.post<TData>(path, variables as Record<string, unknown>, apiOptions)
            break
          case 'PUT':
            result = await api.put<TData>(path, variables as Record<string, unknown>, apiOptions)
            break
          case 'PATCH':
            result = await api.patch<TData>(path, variables as Record<string, unknown>, apiOptions)
            break
          case 'DELETE':
            result = await api.del<TData>(path, variables as Record<string, unknown>, apiOptions)
            break
          default:
            throw new Error(`Unsupported method: ${method}`)
        }

        setData(result)
        setStatus('success')

        // Invalidate cache keys
        invalidateKeys.forEach((key) => cache.invalidate(key))

        onSuccess?.(result, variables)
        return result
      } catch (err) {
        const error = err instanceof Error ? err : new Error('Unknown error')
        setError(error)
        setStatus('error')
        onError?.(error, variables)
        return undefined
      }
    },
    [path, method, hasToast, invalidateKeys, onSuccess, onError]
  )

  const reset = useCallback(() => {
    setData(undefined)
    setStatus('idle')
    setError(null)
  }, [])

  return {
    mutate,
    data,
    isLoading: status === 'loading',
    isSuccess: status === 'success',
    isError: status === 'error',
    error,
    reset,
  }
}
