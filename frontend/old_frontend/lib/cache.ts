/**
 * Simple in-memory cache for API responses
 * Used for caching static or rarely-changing data
 */

type CacheEntry<T> = {
  data: T
  expiresAt: number
}

class MemoryCache {
  private cache = new Map<string, CacheEntry<unknown>>()

  /**
   * Get cached data if it exists and hasn't expired
   */
  get<T>(key: string): T | null {
    const entry = this.cache.get(key)
    if (!entry) return null

    if (Date.now() > entry.expiresAt) {
      this.cache.delete(key)
      return null
    }

    return entry.data as T
  }

  /**
   * Set cache data with TTL (time to live) in milliseconds
   */
  set<T>(key: string, data: T, ttlMs: number): void {
    this.cache.set(key, {
      data,
      expiresAt: Date.now() + ttlMs,
    })
  }

  /**
   * Invalidate a specific cache key
   */
  invalidate(key: string): void {
    this.cache.delete(key)
  }

  /**
   * Clear all cache
   */
  clear(): void {
    this.cache.clear()
  }

  /**
   * Get cache size
   */
  get size(): number {
    return this.cache.size
  }
}

// Singleton instance
export const cache = new MemoryCache()

// Cache TTL constants (in milliseconds)
export const CACHE_TTL = {
  SHORT: 60 * 1000, // 1 minute
  MEDIUM: 5 * 60 * 1000, // 5 minutes
  LONG: 30 * 60 * 1000, // 30 minutes
  HOUR: 60 * 60 * 1000, // 1 hour
} as const

/**
 * Utility to create a cached API call
 * @param fetcher - Function that makes the API call
 * @param cacheKey - Unique key for this data
 * @param ttlMs - Cache duration in milliseconds
 */
export async function withCache<T>(
  fetcher: () => Promise<T>,
  cacheKey: string,
  ttlMs: number = CACHE_TTL.MEDIUM
): Promise<T> {
  const cached = cache.get<T>(cacheKey)
  if (cached !== null) {
    return cached
  }

  const data = await fetcher()
  cache.set(cacheKey, data, ttlMs)
  return data
}
