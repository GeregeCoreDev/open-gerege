// lib/bootstrap/AuthBootstrap.tsx
'use client'

import { useEffect } from 'react'
import { setUnauthorizedHandler } from '@/lib/api'
import { useRouter, usePathname } from '@/i18n/navigation'
import { useUserStore } from '@/lib/stores/user'
import { useOrgStore } from '@/lib/stores/org'

export default function AuthBootstrap() {
  const router = useRouter()
  const pathname = usePathname()

  useEffect(() => {
    setUnauthorizedHandler(async () => {
      useUserStore.getState().clearAll?.()
      useOrgStore.getState().clear?.()

      const parts = (pathname || '').split('/').filter(Boolean)
      const tail = parts[parts.length - 1] || ''
      const isPublic = ['home', 'login', 'signup', 'error', 'msg'].includes(tail)
      if (isPublic) return

      router.push('/home')
    })
    return () => setUnauthorizedHandler(null)
  }, [router, pathname])

  return null
}
