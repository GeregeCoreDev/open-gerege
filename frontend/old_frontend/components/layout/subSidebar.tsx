/**
 * 📋 Sub Sidebar Component
 *
 * Хоёрдахь sidebar (Level 2, 3 menu)
 * - Title header
 * - Collapsible menu groups
 * - Active state indicators
 *
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import * as React from 'react'
import { cn } from '@/lib/utils'
import { LucideIcon } from '@/lib/utils/icon'
import { Link, usePathname } from '@/i18n/navigation'
import { Collapsible, CollapsibleTrigger, CollapsibleContent } from '@/components/ui/collapsible'
import { ChevronDown } from 'lucide-react'
import { useMenuStore } from '@/lib/stores/menu'
import { MAIN_SIDEBAR_WIDTH } from './mainSidebar'

export const SUB_SIDEBAR_WIDTH = 240

interface SubSidebarProps {
  collapsed?: boolean
}

// Normalize path helper
function normalizePath(p: string) {
  let x = p || ''
  x = x.replace(/^\/[a-z]{2}(?:-[A-Z]{2})?(?=\/|$)/, '')
  if (!x.startsWith('/')) x = '/' + x
  return x.replace(/\/+$/, '')
}

export default function SubSidebar({ collapsed = false }: SubSidebarProps) {
  const pathname = usePathname()
  const currentPath = normalizePath(pathname)

  // Store-оос авах
  const { menuList, selectedRootId, openGroups, toggleGroup } = useMenuStore()

  // Active root
  const activeRoot = React.useMemo(() => {
    if (selectedRootId == null && menuList.length > 0) return menuList[0]
    return menuList.find((r) => r.id === selectedRootId) ?? menuList[0]
  }, [menuList, selectedRootId])

  // Sub sidebar items (level 2, 3)
  const items = activeRoot?.children ?? []

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
          {activeRoot?.name || ''}
        </p>
      </div>

      {/* Menu items */}
      <nav className="flex-1 overflow-y-auto px-3 pb-4">
        {items.map((item) => {
          const hasChildren = (item.children?.length ?? 0) > 0
          const isOpen = openGroups[item.id] ?? true

          // No children - render as link
          if (!hasChildren) {
            const itemPath = normalizePath(item.path || '#')
            const isActive =
              itemPath !== '/' &&
              (currentPath === itemPath || currentPath.startsWith(itemPath + '/'))

            return (
              <Link
                key={item.id}
                href={item.path || '#'}
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
                <LucideIcon name={item.icon || 'File'} className="h-[18px] w-[18px] shrink-0" />
                <span className="truncate">{item.name}</span>
              </Link>
            )
          }

          // Has children - render as collapsible
          return (
            <Collapsible key={item.id} open={isOpen} onOpenChange={() => toggleGroup(item.id)}>
              <CollapsibleTrigger asChild>
                <button className="mb-0.5 flex w-full items-center justify-between rounded-md px-3 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800/50 dark:hover:text-gray-200">
                  <div className="flex items-center gap-3">
                    <LucideIcon
                      name={item.icon || 'Folder'}
                      className="h-[18px] w-[18px] shrink-0"
                    />
                    <span className="truncate">{item.name}</span>
                  </div>
                  <ChevronDown
                    className={cn(
                      'h-4 w-4 shrink-0 text-gray-400 transition-transform duration-200',
                      !isOpen && '-rotate-90',
                    )}
                  />
                </button>
              </CollapsibleTrigger>

              {/* Level 3 children */}
              <CollapsibleContent className="ml-3 pl-4">
                {item.children?.map((child) => {
                  const childPath = normalizePath(child.path || '#')
                  const isChildActive =
                    childPath !== '/' &&
                    (currentPath === childPath || currentPath.startsWith(childPath + '/'))

                  return (
                    <Link
                      key={child.id}
                      href={child.path || '#'}
                      className={cn(
                        'relative mb-0.5 flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors',
                        isChildActive
                          ? 'bg-primary/10 text-primary dark:bg-primary/20 dark:text-primary font-medium'
                          : 'text-gray-500 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-500 dark:hover:bg-gray-800/50 dark:hover:text-gray-200',
                      )}
                    >
                      {isChildActive && (
                        <span className="bg-primary absolute top-1/2 left-0 h-5 w-1 -translate-y-1/2 rounded-r-full" />
                      )}
                      <LucideIcon
                        name={child.icon || 'Circle'}
                        className="h-[16px] w-[16px] shrink-0"
                      />
                      <span className="truncate">{child.name}</span>
                    </Link>
                  )
                })}
              </CollapsibleContent>
            </Collapsible>
          )
        })}
      </nav>
    </aside>
  )
}
