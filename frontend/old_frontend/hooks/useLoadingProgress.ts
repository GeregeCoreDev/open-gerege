/**
 * 🔄 useLoadingProgress Hook
 * 
 * Loading progress bar state-ийг удирдах hook.
 * Бүх CRUD page-д давтагдаж байсан 25 мөр кодыг 1 мөрт багасгана.
 * 
 * Features:
 * - ✅ Auto-animating progress (0-90%)
 * - ✅ Completes to 100% when loading ends
 * - ✅ Auto-cleanup after animation
 * - ✅ Optimized with useRef & useEffect
 * 
 * Usage:
 * ```tsx
 * const [loading, setLoading] = useState(false)
 * const progress = useLoadingProgress(loading)
 * 
 * return <Progress value={progress} />
 * ```
 * 
 * @param isLoading - Loading state (true/false)
 * @returns progress value (0-100)
 * 
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import { useState, useEffect, useRef } from 'react'

export function useLoadingProgress(isLoading: boolean): number {
  const [progress, setProgress] = useState(0)
  const progressTimer = useRef<ReturnType<typeof setInterval> | null>(null)

  useEffect(() => {
    // Cleanup previous timer if exists
    if (progressTimer.current) {
      clearInterval(progressTimer.current)
      progressTimer.current = null
    }
    
    let timeoutId: ReturnType<typeof setTimeout> | null = null
    
    if (isLoading) {
      // Start progress animation from 0
      setProgress(0)
      
      // Increment progress randomly (stops at 90%)
      progressTimer.current = setInterval(() => {
        setProgress((p) => Math.min(p + Math.random() * 12 + 8, 90))
      }, 250)
    } else {
      // Complete to 100%
      setProgress(100)
      
      // Reset to 0 after animation
      timeoutId = setTimeout(() => setProgress(0), 300)
    }
    
    // Cleanup on unmount or dependency change
    return () => {
      if (progressTimer.current) {
        clearInterval(progressTimer.current)
        progressTimer.current = null
      }
      if (timeoutId) clearTimeout(timeoutId)
    }
  }, [isLoading])

  return progress
}

