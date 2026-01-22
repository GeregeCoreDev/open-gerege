'use client'

import { memo, useMemo } from 'react'
import { Building2, ChevronDown, Moon, Sun } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Input } from '@/components/ui/input'
import { useOrgStore } from '@/lib/stores/org'
import { useTranslations } from 'next-intl'
import { useTheme } from 'next-themes'
import { useLocale } from 'next-intl'
import { useRouter, usePathname } from '@/i18n/navigation'
import { Search } from 'lucide-react'
import Image from 'next/image'
import NotificationDropdown from './notificationDropdown'

function MainHeaderInner() {
  const t = useTranslations()
  const { theme, setTheme } = useTheme()
  const locale = useLocale()
  const router = useRouter()
  const pathname = usePathname()

  const { organizations, selectedOrganization, selectOrg } = useOrgStore()
  const orgId = selectedOrganization?.id ? String(selectedOrganization.id) : undefined
  const orgMap = useMemo(
    () => new Map(organizations.map((o) => [String(o.id), o])),
    [organizations],
  )

  function toggleTheme() {
    setTheme(theme === 'dark' ? 'light' : 'dark')
  }

  function changeLocale(newLocale: string) {
    router.replace(pathname, { locale: newLocale })
  }

  return (
    <div className="flex h-full w-full items-center justify-end">
      {/* Right side: Search + Icons + Org selector */}
      <div className="flex items-center gap-3">
        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
          <Input placeholder={t('search') + '...'} className="h-9 w-44 pl-9 text-sm" />
        </div>

        {/* Notification */}
        <NotificationDropdown />

        {/* Dark mode toggle */}
        <Button variant="ghost" size="icon" className="h-9 w-9" onClick={toggleTheme}>
          {theme === 'dark' ? (
            <Sun className="h-5 w-5 text-gray-500" />
          ) : (
            <Moon className="h-5 w-5 text-gray-500" />
          )}
        </Button>

        {/* Language selector */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon" className="h-9 w-9">
              <Image
                src={`/flag/${locale}.png`}
                alt="flag"
                width={20}
                height={20}
                className="rounded-sm"
              />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => changeLocale('mn')} className="gap-2">
              <Image src="/flag/mn.png" alt="mn" width={16} height={16} className="rounded-sm" />
              <span>Монгол</span>
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => changeLocale('en')} className="gap-2">
              <Image src="/flag/en.png" alt="en" width={16} height={16} className="rounded-sm" />
              <span>English</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        {/* Organization selector (compact) */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-9 gap-2 px-3 text-sm font-medium">
              <Building2 className="h-4 w-4 text-gray-500" />
              <span className="max-w-32 truncate">{selectedOrganization?.name || t('organization')}</span>
              <ChevronDown className="h-4 w-4 opacity-50" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-64">
            <DropdownMenuLabel className="text-xs uppercase opacity-70">
              {t('organization')}
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuRadioGroup
              value={orgId}
              onValueChange={(v) => {
                const o = orgMap.get(v)
                if (o) selectOrg(o)
              }}
            >
              <div className="max-h-[280px] overflow-auto">
                {organizations.map((o) => (
                  <DropdownMenuRadioItem key={o.id} value={String(o.id)} className="py-2">
                    <span className="truncate font-medium">{o.name}</span>
                  </DropdownMenuRadioItem>
                ))}
              </div>
            </DropdownMenuRadioGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  )
}

const MainHeader = memo(MainHeaderInner)
export default MainHeader
