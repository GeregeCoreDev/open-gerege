'use client'
import * as React from 'react'
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuLabel,
} from '@/components/ui/dropdown-menu'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import {
  User as UserIcon,
  Settings,
  Wallet,
  PanelsTopLeft,
  LogOut,
  ChevronDown,
  ShieldHalf,
  ShieldCheck,
} from 'lucide-react'
import { useUserStore } from '@/lib/stores/user'
import { Button } from '../ui/button'
import { useRouter } from 'next/navigation'
import { Link } from '@/i18n/navigation'
import { useLocale } from 'next-intl'
import { useRoleStore } from '@/lib/stores/role'
import { useMenuStore } from '@/lib/stores/menu'
import { logout } from '@/lib/logout'

export default function ProfileDropdown(props: {
  variant?: 'full' | 'avatar' | 'ghost'
  size?: number
}) {
  const variant = props.variant ?? 'full'
  const size = props.size ?? 32

  const user_info = useUserStore((s) => s.user_info)
  const user_name = useUserStore((s) => s.user_name)
  const profile_image = useUserStore((s) => s.profile_image)
  const status = useUserStore((s) => s.status)

  const { roleList, selectRole, selectedRole } = useRoleStore()
  const { menuList, getMenuList, getFirstChildPath, selectRoot, clear } = useMenuStore()

  const router = useRouter()

  const locale = useLocale()

  async function navigateToCP() {
    try {
      // Store-ийн menuList-ээс авах, хоосон бол API дуудах
      let menus = menuList
      if (menus.length === 0) {
        menus = await getMenuList()
      }

      if (menus.length > 0) {
        const firstMenu = menus[0]
        selectRoot(firstMenu.id)
        // Эхний path-тай menu-г олох
        const firstPath = getFirstChildPath(firstMenu.id)
        if (firstPath) {
          const href = toLocaleHref(firstPath)
          router.push(href)
          router.refresh()
        } else {
          router.push(`/${locale}/profile`)
          router.refresh()
        }
      } else {
        router.push(`/${locale}/profile`)
        router.refresh()
      }
    } catch (error) {
      console.error('Failed to navigate to control panel:', error)
      router.push(`/${locale}/profile`)
      router.refresh()
    }
  }

  async function onLogout() {
    await logout()
  }

  const toLocaleHref = (rawPath?: string) => {
    const p = (rawPath || '').startsWith('/') ? rawPath : `/${rawPath || ''}`
    return `/${locale}${p}`.replace(/\/{2,}/g, '/')
  }

  async function onSelectRole(role: App.UserRole) {
    if (role.role_id == selectedRole?.role_id) return

    clear()
    await selectRole(role)

    // Menu жагсаалт авна
    const menus = await getMenuList()
    if (menus.length > 0) {
      const firstMenu = menus[0]
      selectRoot(firstMenu.id)
      // Эхний path-тай menu-г олох
      const firstPath = getFirstChildPath(firstMenu.id)
      if (firstPath) {
        const href = toLocaleHref(firstPath)
        router.push(href)
        router.refresh()
      } else {
        router.push(`/${locale}/profile`)
        router.refresh()
      }
    } else {
      router.push(`/${locale}/profile`)
      router.refresh()
    }
  }

  function ActiveDot({ show }: { show: boolean }) {
    if (!show) return null
    return <span className="ml-auto inline-block h-2 w-2 rounded-full bg-black dark:bg-white" />
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        {variant === 'avatar' ? (
          <Avatar
            className="cursor-pointer rounded-full transition-all"
            style={{ width: size, height: size }}
          >
            <AvatarImage
              src={profile_image || undefined}
              alt={user_name || 'user'}
              className="h-full w-full object-cover object-center"
            />
            <AvatarFallback className="bg-gray-200 text-xs font-semibold text-gray-700 uppercase dark:bg-gray-700 dark:text-gray-200">
              {user_info?.last_name?.[0] || ''}
              {user_info?.first_name?.[0] || ''}
            </AvatarFallback>
          </Avatar>
        ) : (
          <Button className="h-12 gap-2" variant={variant === 'ghost' ? 'ghost' : 'outline'}>
            <Avatar className="rounded-full" style={{ width: size, height: size }}>
              <AvatarImage
                src={profile_image || undefined}
                alt={user_name || 'user'}
                className="h-full w-full object-cover object-center"
              />
              <AvatarFallback className="text-xs uppercase">
                {user_info?.last_name?.[0] || ''}
                {user_info?.first_name?.[0] || ''}
              </AvatarFallback>
            </Avatar>
            {variant === 'full' ||
              (variant === 'ghost' && (
                <>
                  <span className="max-w-56 truncate capitalize">
                    {user_name ?? (status === 'loading' ? 'Loading…' : user_name)}
                  </span>
                  <ChevronDown className="size-6 opacity-70" />
                </>
              ))}
          </Button>
        )}
      </DropdownMenuTrigger>

      <DropdownMenuContent align="end" side="right" className="w-56">
        <DropdownMenuLabel>Хувийн мэдээлэл</DropdownMenuLabel>
        <DropdownMenuItem asChild>
          <Link href="/profile" locale={locale} className="flex items-center gap-2">
            <UserIcon className="size-6" /> <span>Profile</span>
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem
          onSelect={(e) => {
            e.preventDefault()
            navigateToCP()
          }}
          className="flex items-center gap-2"
        >
          <PanelsTopLeft className="size-6" /> <span>Control panel</span>
        </DropdownMenuItem>
        <DropdownMenuItem asChild>
          <Link href="/settings" locale={locale} className="flex items-center gap-2">
            <Settings className="size-6" /> <span>Settings</span>
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem asChild>
          <Link href="/wallet" locale={locale} className="flex items-center gap-2">
            <Wallet className="size-6" /> <span>Wallet</span>
          </Link>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuLabel>Дүр солих</DropdownMenuLabel>
        {roleList.map((child) => {
          const isActive = child.role_id == selectedRole?.role_id
          return (
            <DropdownMenuItem asChild key={child.role_id}>
              <div
                className="flex items-center justify-between gap-2"
                key={child.role_id}
                onClick={() => onSelectRole(child)}
              >
                <div className="flex items-center gap-x-2">
                  {isActive ? (
                    <ShieldCheck className="size-6" />
                  ) : (
                    <ShieldHalf className="size-6" />
                  )}{' '}
                  <p>{child.role.name}</p>
                  {child.role.system && (
                    <p className="font-medium text-gray-700 dark:text-gray-200">
                      ({child.role.system?.code})
                    </p>
                  )}
                </div>
                <ActiveDot show={isActive} />
              </div>
            </DropdownMenuItem>
          )
        })}
        {/* <DropdownMenuSeparator />
        <DropdownMenuItem onClick={navigateToChangeSystem} className="flex items-center gap-2">
          <PanelsTopLeft className="size-6" />
          <span>Change system</span>
        </DropdownMenuItem>
        <DropdownMenuItem
          onClick={navigateToChangeOrganization}
          className="flex items-center gap-2"
        >
          <Building2 className="size-6" />
          <span>Change org</span>
        </DropdownMenuItem> */}
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={onLogout} className="flex items-center gap-2 text-red-600">
          <LogOut className="size-6" />
          <span>Logout</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
