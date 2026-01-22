/**
 * 🎯 Main Sidebar Component (Icon Rail)
 *
 * Зүүн талын icon sidebar (Level 1 menu)
 * - Logo
 * - Root menu icons
 * - Settings link
 * - Profile dropdown
 *
 * @author Gerege Core Team
 */

'use client'

import * as React from 'react'
import { cn } from '@/lib/utils'
import { LucideIcon } from '@/lib/utils/icon'
import { Link, usePathname, useRouter } from '@/i18n/navigation'
import { Settings } from 'lucide-react'
import ProfileDropdown from '@/components/layout/profileDropDown'
import { Tooltip, TooltipContent, TooltipTrigger, TooltipProvider } from '@/components/ui/tooltip'
import Image from 'next/image'
import logo from '@/public/logo/logo.png'
import { useTranslations } from 'next-intl'
import { useMenuStore } from '@/lib/stores/menu'

export const MAIN_SIDEBAR_WIDTH = 72

// Normalize path helper
function normalizePath(p: string) {
  let x = p || ''
  x = x.replace(/^\/[a-z]{2}(?:-[A-Z]{2})?(?=\/|$)/, '')
  if (!x.startsWith('/')) x = '/' + x
  return x.replace(/\/+$/, '')
}

const MainSidebar = React.memo(function MainSidebar() {
  const t = useTranslations()
  const pathname = usePathname()
  const router = useRouter()
  const currentPath = normalizePath(pathname)

  // Store-оос авах
  const { menuList, selectedRootId, selectRoot, getFirstChildPath, clearSelectedRoot } =
    useMenuStore()

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

  const activeRoot = React.useMemo(() => {
    if (selectedRootId == null && menuList.length > 0) return menuList[0]
    return menuList.find((r) => r.id === selectedRootId) ?? menuList[0]
  }, [menuList, selectedRootId])

  return (
    <TooltipProvider>
      <aside
        className="fixed inset-y-0 left-0 z-50 flex w-[72px] flex-col items-center border-r py-4"
        style={{
          backgroundColor: 'var(--sidebar)',
          borderColor: 'var(--sidebar-border)',
        }}
      >
        {/* Logo */}
        <div className="mb-6">
          <Link href="/home">
            <Image src={logo} width={36} height={36} alt="logo" className="rounded-md" />
          </Link>
        </div>

        {/* Menu Icons */}
        <div className="flex flex-1 flex-col items-center gap-3">
          {menuList.map((item) => {
            // Settings button идэвхжсэн үед дээд талын menu идэвхжих ёсгүй
            const isActive = !isPersonalPage && activeRoot?.id === item.id

            return (
              <Tooltip key={item.id}>
                <TooltipTrigger asChild>
                  <button
                    onClick={() => {
                      selectRoot(item.id)
                      // Эхний child-д navigate хийх (personal page дээр байсан ч)
                      const firstChildPath = getFirstChildPath(item.id)
                      if (firstChildPath) {
                        router.push(firstChildPath)
                      }
                    }}
                    className={cn(
                      'flex h-11 w-11 cursor-pointer items-center justify-center rounded-xl transition-all',
                      isActive
                        ? 'bg-primary/20 text-primary dark:bg-primary/30 dark:text-primary'
                        : 'text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200',
                    )}
                  >
                    <LucideIcon name={item.icon || 'Folder'} className="h-6 w-6" />
                  </button>
                </TooltipTrigger>
                <TooltipContent side="right">
                  <p>{item.name}</p>
                </TooltipContent>
              </Tooltip>
            )
          })}
        </div>

        {/* Bottom section: Settings + Profile */}
        <div className="mt-auto flex flex-col items-center gap-4 pb-4">
          <Tooltip>
            <TooltipTrigger asChild>
              <Link
                href="/settings"
                onClick={() => {
                  // Settings хуудас руу шилжих үед selectedRootId-г null болгох
                  clearSelectedRoot()
                }}
                className={cn(
                  'flex h-11 w-11 items-center justify-center rounded-xl transition-all',
                  isPersonalPage
                    ? 'bg-primary/20 text-primary dark:bg-primary/30 dark:text-primary'
                    : 'text-gray-400 hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-gray-800',
                )}
              >
                <Settings className="h-6 w-6" />
              </Link>
            </TooltipTrigger>
            <TooltipContent side="right">
              <p>{t('settings')}</p>
            </TooltipContent>
          </Tooltip>

          <ProfileDropdown variant="avatar" size={42} />
        </div>
      </aside>
    </TooltipProvider>
  )
})

export default MainSidebar
