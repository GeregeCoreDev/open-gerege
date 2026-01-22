import { describe, it, expect, beforeEach, vi } from 'vitest'
import { cache, CACHE_TTL, withCache } from '../cache'

describe('MemoryCache', () => {
  beforeEach(() => {
    cache.clear()
  })

  describe('get and set', () => {
    it('stores and retrieves data', () => {
      cache.set('test-key', { value: 42 }, CACHE_TTL.MEDIUM)
      const result = cache.get<{ value: number }>('test-key')
      expect(result).toEqual({ value: 42 })
    })

    it('returns null for non-existent keys', () => {
      const result = cache.get('non-existent')
      expect(result).toBeNull()
    })

    it('returns null for expired entries', async () => {
      cache.set('expired-key', 'data', 10) // 10ms TTL
      await new Promise((resolve) => setTimeout(resolve, 20))
      const result = cache.get('expired-key')
      expect(result).toBeNull()
    })

    it('removes expired entries on get', async () => {
      cache.set('expired-key', 'data', 10)
      await new Promise((resolve) => setTimeout(resolve, 20))
      cache.get('expired-key')
      expect(cache.size).toBe(0)
    })
  })

  describe('invalidate', () => {
    it('removes a specific key', () => {
      cache.set('key1', 'value1', CACHE_TTL.MEDIUM)
      cache.set('key2', 'value2', CACHE_TTL.MEDIUM)

      cache.invalidate('key1')

      expect(cache.get('key1')).toBeNull()
      expect(cache.get('key2')).toBe('value2')
    })
  })

  describe('clear', () => {
    it('removes all entries', () => {
      cache.set('key1', 'value1', CACHE_TTL.MEDIUM)
      cache.set('key2', 'value2', CACHE_TTL.MEDIUM)

      cache.clear()

      expect(cache.size).toBe(0)
    })
  })
})

describe('withCache', () => {
  beforeEach(() => {
    cache.clear()
  })

  it('returns cached data on subsequent calls', async () => {
    const fetcher = vi.fn().mockResolvedValue({ data: 'fresh' })

    // First call - should fetch
    const result1 = await withCache(fetcher, 'test-key', CACHE_TTL.MEDIUM)
    expect(result1).toEqual({ data: 'fresh' })
    expect(fetcher).toHaveBeenCalledTimes(1)

    // Second call - should return cached
    const result2 = await withCache(fetcher, 'test-key', CACHE_TTL.MEDIUM)
    expect(result2).toEqual({ data: 'fresh' })
    expect(fetcher).toHaveBeenCalledTimes(1) // Not called again
  })

  it('fetches fresh data when cache expires', async () => {
    const fetcher = vi.fn().mockResolvedValue({ data: 'fresh' })

    await withCache(fetcher, 'test-key', 10) // 10ms TTL
    await new Promise((resolve) => setTimeout(resolve, 20))

    await withCache(fetcher, 'test-key', 10)
    expect(fetcher).toHaveBeenCalledTimes(2)
  })
})

describe('CACHE_TTL constants', () => {
  it('has correct values', () => {
    expect(CACHE_TTL.SHORT).toBe(60 * 1000)
    expect(CACHE_TTL.MEDIUM).toBe(5 * 60 * 1000)
    expect(CACHE_TTL.LONG).toBe(30 * 60 * 1000)
    expect(CACHE_TTL.HOUR).toBe(60 * 60 * 1000)
  })
})
