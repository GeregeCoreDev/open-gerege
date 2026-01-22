'use client'

import Image from 'next/image'
import logo from '@/public/logo/logo.png'
import { cn } from '@/lib/utils'
import { useTranslations, useLocale } from 'next-intl'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Button } from '@/components/ui/button'
import { Moon, Sun, Check } from 'lucide-react'
import { useRouter, usePathname } from '@/i18n/navigation'
import { defaultLocale, isLocale, type Locale } from '@/i18n/config'
import FlagIcon from '../flag-icon/flagIcon'
import { useMemo, useCallback, useState, useEffect } from 'react'
import ProfileDropdown from './profileDropDown'
import { useUserStore } from '@/lib/stores/user'
import { SignInPopUp } from './SignInPopup' 

export default function HomeHeader() {
  const t = useTranslations()
  const router = useRouter()
  const pathname = usePathname()
  const [darkMode, setDarkMode] = useState(false)

  const user = useUserStore((s) => s.user_info)

  useEffect(() => {
    if (darkMode) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [darkMode])

  const l = useLocale()
  const currentLocale: Locale = isLocale(l) ? l : defaultLocale

  const languages = useMemo(
    () =>
      [
        { code: 'mn', name: 'Монгол' },
        { code: 'en', name: 'English' },
      ] as const,
    [],
  )

  const currentLanguage = useMemo(
    () => languages.find((lang) => lang.code === currentLocale) ?? languages[0],
    [languages, currentLocale],
  )

  const changeLanguage = useCallback(
    (locale: Locale) => {
      router.replace(pathname, { locale })
    },
    [router, pathname],
  )

  return (
    <div className="w-full border-b border-slate-200 dark:border-gray-700" style={{ backgroundColor: 'var(--background)' }}>
      <div
        className={cn('mx-auto hidden h-20 w-full max-w-7xl items-center justify-between lg:flex')}
      >
        <div className="flex items-center gap-x-2">
          <Image src={logo} width={32} height={32} alt="logo" />
          <p className="text-lg font-medium capitalize">{t('system_name')}</p>
        </div>

        <div className="flex items-center gap-x-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" className="h-12 gap-2">
                <FlagIcon currentLanguage={currentLanguage} />
                <span>{currentLanguage.name}</span>
              </Button>
            </DropdownMenuTrigger>

            <DropdownMenuContent align="end" className="w-36">
              {languages.map((lang) => (
                <DropdownMenuItem
                  key={lang.code}
                  onClick={() => changeLanguage(lang.code)}
                  className="flex items-center justify-between gap-x-2"
                >
                  <FlagIcon currentLanguage={lang} />
                  <div className="flex items-center gap-2">
                    <p>{lang.name}</p>
                    {currentLocale === lang.code ? (
                      <Check className="size-6" />
                    ) : (
                      <div className="size-4" />
                    )}
                  </div>
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>

          <Button variant="outline" className="h-12 w-12" onClick={() => setDarkMode(!darkMode)}>
            {darkMode ? <Sun className="size-6" /> : <Moon className="size-6" />}
            <span className="sr-only">Toggle theme</span>
          </Button>

          {Boolean(user) ? <ProfileDropdown variant="full" /> : <SignInPopUp />}
        </div>
      </div>
    </div>
  )
}
