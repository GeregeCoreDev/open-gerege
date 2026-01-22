import { defaultLocale } from '@/i18n/config'
import api from './api'
import { useUserStore } from './stores/user'

function readLocale() {
  // Prefer cookie set by next-intl
  if (typeof document !== 'undefined') {
    const match = document.cookie.split('; ').find((row) => row.startsWith('NEXT_LOCALE='))
    const fromCookie = match?.split('=')[1]
    if (fromCookie) return decodeURIComponent(fromCookie)

    const fromStorage = localStorage.getItem('locale')
    if (fromStorage) return fromStorage
  }

  return defaultLocale
}

/**
 * Logout helper for client side.
 * Clears session, resets stores, and navigates back to locale home.
 */
export async function logout() {
  // Guard server-side usage
  if (typeof window === 'undefined') return

  const { clearAll } = useUserStore.getState()

  // Read locale before wiping storage
  const locale = readLocale()

  try {
    await api.post('/auth/logout')
  } catch (error) {
    // Ignore error, proceed to clear client session
    console.error('Logout API failed:', error)
  }

  localStorage.clear()
  sessionStorage.clear()
  clearAll()

  window.location.replace(`/${locale}/home`)
}
