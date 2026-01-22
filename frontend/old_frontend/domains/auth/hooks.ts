/**
 * Auth Domain Hooks
 */

'use client'

import { useCallback, useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { logout as logoutApi } from './api'
import type { AuthStatus } from './types'
import {
  SSO_CONFIG,
  getSSOEmbedUrl,
  getSSORedirectUrl,
  parseSSOMessage,
} from './api'

/**
 * Auth status hook
 */
export function useAuthStatus(): {
  status: AuthStatus
  isAuthenticated: boolean
  isLoading: boolean
} {
  const [status, setStatus] = useState<AuthStatus>('loading')

  useEffect(() => {
    const hasSession = document.cookie.includes('session')
    setStatus(hasSession ? 'authenticated' : 'unauthenticated')
  }, [])

  return {
    status,
    isAuthenticated: status === 'authenticated',
    isLoading: status === 'loading',
  }
}

/**
 * Logout hook
 */
export function useLogout() {
  const router = useRouter()
  const [isLoading, setIsLoading] = useState(false)

  const handleLogout = useCallback(
    async (everywhere?: boolean) => {
      setIsLoading(true)
      try {
        await logoutApi({ everywhere })
        if (typeof window !== 'undefined') {
          localStorage.removeItem('user-store')
          localStorage.removeItem('orgs-store')
          localStorage.removeItem('systems-store')
          localStorage.removeItem('role-store')
          localStorage.removeItem('menu-store')
        }
        router.push('/')
      } catch (error) {
        console.error('Logout failed:', error)
      } finally {
        setIsLoading(false)
      }
    },
    [router]
  )

  return {
    logout: handleLogout,
    isLoading,
  }
}

/**
 * Protected route hook
 */
export function useRequireAuth(redirectTo = '/') {
  const router = useRouter()
  const { status, isAuthenticated } = useAuthStatus()

  useEffect(() => {
    if (status !== 'loading' && !isAuthenticated) {
      router.push(redirectTo)
    }
  }, [status, isAuthenticated, router, redirectTo])

  return { isAuthenticated, isLoading: status === 'loading' }
}

// ==========================================
// SSO Hooks
// ==========================================

/**
 * SSO Message listener hook
 */
export function useSSOMessageListener(
  onSuccess: (session: string) => void,
  onCancel?: () => void,
  onError?: (error: string) => void
) {
  useEffect(() => {
    const handleMessage = (event: MessageEvent) => {
      const message = parseSSOMessage(event)
      if (!message) return

      switch (message.type) {
        case 'SSO_AUTH_SUCCESS':
          const session = message.token || message.sid
          if (session) {
            onSuccess(session)
          }
          break
        case 'SSO_AUTH_CANCEL':
          onCancel?.()
          break
        case 'SSO_AUTH_ERROR':
          onError?.(message.error || 'Unknown error')
          break
      }
    }

    window.addEventListener('message', handleMessage)
    return () => window.removeEventListener('message', handleMessage)
  }, [onSuccess, onCancel, onError])
}

/**
 * SSO Login hook
 */
export function useSSOLogin() {
  const [isOpen, setIsOpen] = useState(false)
  const [isLoading, setIsLoading] = useState(false)

  const openPopup = useCallback(() => {
    setIsOpen(true)
  }, [])

  const closePopup = useCallback(() => {
    setIsOpen(false)
  }, [])

  const redirectToSSO = useCallback(() => {
    setIsLoading(true)
    window.location.href = getSSORedirectUrl()
  }, [])

  const handleSuccess = useCallback((session: string) => {
    setIsLoading(true)
    setIsOpen(false)
    // Redirect to frontend proxy /api/auth/verify (which forwards to backend)
    // This ensures cookies are set on the frontend domain
    window.location.href = `/api/auth/verify?sid=${session}`
  }, [])

  const handleCancel = useCallback(() => {
    setIsOpen(false)
    console.warn('SSO login canceled')
  }, [])

  const handleError = useCallback((error: string) => {
    setIsOpen(false)
    console.error('SSO login error:', error)
  }, [])

  // SSO message listener
  useSSOMessageListener(handleSuccess, handleCancel, handleError)

  return {
    isOpen,
    isLoading,
    openPopup,
    closePopup,
    redirectToSSO,
    handleSuccess,
    handleCancel,
    handleError,
    ssoEmbedUrl: getSSOEmbedUrl(),
    ssoConfig: SSO_CONFIG,
  }
}
