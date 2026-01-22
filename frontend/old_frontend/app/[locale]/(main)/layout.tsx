'use client'

import { useState } from 'react'
import MainHeader from '@/components/layout/mainHeader'
import MainSidebar, { MAIN_SIDEBAR_WIDTH } from '@/components/layout/mainSidebar'
import SubSidebar, { SUB_SIDEBAR_WIDTH } from '@/components/layout/subSidebar'
import { PanelLeftClose, PanelLeftOpen } from 'lucide-react'

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const [subSidebarCollapsed, setSubSidebarCollapsed] = useState(false)

  // Dynamic widths
  const subSidebarWidth = subSidebarCollapsed ? 0 : SUB_SIDEBAR_WIDTH
  const totalSidebarWidth = MAIN_SIDEBAR_WIDTH + subSidebarWidth

  return (
    <div
      className="flex h-screen w-screen overflow-hidden text-gray-900 dark:text-gray-100"
      style={{ backgroundColor: 'var(--background)' }}
    >
      {/* Main Sidebar (Icon Rail - Level 1) */}
      <MainSidebar />

      {/* Sub Sidebar (Level 2, 3) */}
      <SubSidebar collapsed={subSidebarCollapsed} />

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
            className="flex h-8 w-8 items-center justify-center rounded-md bg-white text-primary-500 transition-all hover:bg-gray-50 active:scale-95 dark:bg-gray-800 dark:text-primary-400"
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
