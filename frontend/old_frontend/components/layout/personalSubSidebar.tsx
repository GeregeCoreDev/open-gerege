/**
 * 👤 Personal Sub Sidebar Component
 *
 * Хувийн хэсгийн хоёрдахь sidebar
 * - Profile
 * - Settings
 * - Wallet
 *
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import * as React from 'react'
import { cn } from '@/lib/utils'
import { LucideIcon } from '@/lib/utils/icon'
import { Link, usePathname } from '@/i18n/navigation'
import { useTranslations } from 'next-intl'
import { MAIN_SIDEBAR_WIDTH } from './mainSidebar'

export const PERSONAL_SUB_SIDEBAR_WIDTH = 240

// Profile menu items
export const profileMenuList = [
  {
    id: 1,
    path: '/profile',
    icon: 'i-lucide-user',
    nameKey: 'profile',
  },
  {
    id: 2,
    path: '/settings',
    icon: 'i-lucide-cog',
    nameKey: 'settings',
  },
  {
    id: 3,
    path: '/wallet',
    icon: 'i-lucide-wallet',
    nameKey: 'wallet',
  },
]

interface PersonalSubSidebarProps {
  collapsed?: boolean
}

// Normalize path helper
function normalizePath(p: string) {
  let x = p || ''
  x = x.replace(/^\/[a-z]{2}(?:-[A-Z]{2})?(?=\/|$)/, '')
  if (!x.startsWith('/')) x = '/' + x
  return x.replace(/\/+$/, '')
}

export default function PersonalSubSidebar({ collapsed = false }: PersonalSubSidebarProps) {
  const t = useTranslations()
  const pathname = usePathname()
  const currentPath = normalizePath(pathname)

  // Find active item
  const activeItem = React.useMemo(() => {
    return profileMenuList.find((item) => {
      const itemPath = normalizePath(item.path)
      return (
        itemPath !== '/' && (currentPath === itemPath || currentPath.startsWith(itemPath + '/'))
      )
    })
  }, [currentPath])

  return (
    <aside
      className={cn(
        'fixed inset-y-0 z-40 flex flex-col border-r transition-all duration-300',
        collapsed ? 'w-0 overflow-hidden opacity-0' : 'w-[240px] opacity-100',
      )}
      style={{
        left: MAIN_SIDEBAR_WIDTH,
        backgroundColor: 'var(--sidebar)',
        borderColor: 'var(--sidebar-border)',
      }}
    >
      {/* Header */}
      <div className="flex h-14 items-center px-4">
        <p className="truncate text-base font-semibold text-gray-700 dark:text-white">
          {t(activeItem?.nameKey || 'profile')}
        </p>
      </div>

      {/* Menu items */}
      <nav className="flex-1 overflow-y-auto px-3 pb-4">
        {profileMenuList.map((item) => {
          const itemPath = normalizePath(item.path)
          const isActive =
            itemPath !== '/' && (currentPath === itemPath || currentPath.startsWith(itemPath + '/'))

          return (
            <Link
              key={item.id}
              href={item.path}
              className={cn(
                'relative mb-0.5 flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors',
                isActive
                  ? 'bg-primary/10 text-primary dark:bg-primary/20 dark:text-primary font-medium'
                  : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800/50 dark:hover:text-gray-200',
              )}
            >
              {isActive && (
                <span className="bg-primary absolute top-1/2 left-0 h-5 w-1 -translate-y-1/2 rounded-r-full" />
              )}
              <LucideIcon name={item.icon} className="h-[18px] w-[18px] shrink-0" />
              <span className="truncate">{t(item.nameKey)}</span>
            </Link>
          )
        })}
      </nav>
    </aside>
  )
}
