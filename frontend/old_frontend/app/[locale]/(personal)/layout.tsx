'use client'

import { useState } from 'react'
import MainHeader from '@/components/layout/mainHeader'
import MainSidebar, { MAIN_SIDEBAR_WIDTH } from '@/components/layout/mainSidebar'
import PersonalSubSidebar, {
  PERSONAL_SUB_SIDEBAR_WIDTH,
} from '@/components/layout/personalSubSidebar'
import SubSidebar, { SUB_SIDEBAR_WIDTH } from '@/components/layout/subSidebar'
import { PanelLeftClose, PanelLeftOpen } from 'lucide-react'
import { usePathname } from '@/i18n/navigation'
import { useMenuStore } from '@/lib/stores/menu'
import * as React from 'react'

// Normalize path helper
function normalizePath(p: string) {
  let x = p || ''
  x = x.replace(/^\/[a-z]{2}(?:-[A-Z]{2})?(?=\/|$)/, '')
  if (!x.startsWith('/')) x = '/' + x
  return x.replace(/\/+$/, '')
}

export default function PersonalLayout({ children }: { children: React.ReactNode }) {
  const [subSidebarCollapsed, setSubSidebarCollapsed] = useState(false)
  const pathname = usePathname()
  const currentPath = normalizePath(pathname)
  const { selectedRootId } = useMenuStore()

  // Check if we're on a personal page (settings, profile, wallet)
  const isPersonalPage = React.useMemo(() => {
    return (
      currentPath === '/settings' ||
      currentPath === '/profile' ||
      currentPath === '/wallet' ||
      currentPath.startsWith('/settings/') ||
      currentPath.startsWith('/profile/') ||
      currentPath.startsWith('/wallet/')
    )
  }, [currentPath])

  // If a level 1 menu is selected and we're not on a personal page, show SubSidebar instead of PersonalSubSidebar
  // Settings хуудас дээр байхад PersonalSubSidebar харуулах
  const showMainSubSidebar = selectedRootId != null && !isPersonalPage

  // Dynamic widths
  const subSidebarWidth = subSidebarCollapsed
    ? 0
    : showMainSubSidebar
      ? SUB_SIDEBAR_WIDTH
      : PERSONAL_SUB_SIDEBAR_WIDTH
  const totalSidebarWidth = MAIN_SIDEBAR_WIDTH + subSidebarWidth

  return (
    <div
      className="flex h-screen w-screen overflow-hidden text-gray-900 dark:text-gray-100"
      style={{ backgroundColor: 'var(--background)' }}
    >
      {/* Main Sidebar (Icon Rail - Level 1) */}
      <MainSidebar />

      {/* Sub Sidebar - Show main SubSidebar if level 1 menu is selected, otherwise show PersonalSubSidebar */}
      {showMainSubSidebar ? (
        <SubSidebar collapsed={subSidebarCollapsed} />
      ) : (
        <PersonalSubSidebar collapsed={subSidebarCollapsed} />
      )}

      {/* Main Content */}
      <div
        className="flex h-full w-full flex-1 flex-col overflow-hidden transition-all duration-300"
        style={{ paddingLeft: totalSidebarWidth }}
      >
        {/* Header with toggle button */}
        <header
          className="sticky top-0 z-30 flex h-16 shrink-0 items-center gap-4 border-b border-gray-200 px-4 dark:border-gray-800"
          style={{ backgroundColor: 'var(--background)' }}
        >
          {/* Sidebar toggle button */}
          <button
            className="text-primary-500 dark:text-primary-400 flex h-8 w-8 items-center justify-center rounded-md bg-white transition-all hover:bg-gray-50 active:scale-95 dark:bg-gray-800"
            onClick={() => setSubSidebarCollapsed(!subSidebarCollapsed)}
          >
            {subSidebarCollapsed ? (
              <PanelLeftOpen className="h-6 w-6 transition-transform duration-200" />
            ) : (
              <PanelLeftClose className="h-6 w-6 transition-transform duration-200" />
            )}
          </button>
          <div className="flex-1">
            <MainHeader />
          </div>
        </header>

        {/* Content */}
        <main className="flex-1 overflow-y-auto" style={{ backgroundColor: 'var(--background)' }}>
          {children}
        </main>
      </div>
    </div>
  )
}
