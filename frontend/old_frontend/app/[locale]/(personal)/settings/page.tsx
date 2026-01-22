/**
 * ⚙️ Settings Page - Appearance
 * 
 * Theme тохиргооны хуудас
 * - Theme mode (System/Light/Dark)
 * - Primary color
 * - Light/Dark color schemes
 * - Notification settings
 * - Card skin
 * - Reset theme
 * 
 * @author Gerege Core Team
 */

'use client'

import * as React from 'react'
import { useTranslations } from 'next-intl'
import { useTheme } from 'next-themes'
import { useRouter, usePathname } from '@/i18n/navigation'
import { Card, CardContent } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@/components/ui/select'
import { cn } from '@/lib/utils'

// ============================================================
// Types & Constants
// ============================================================

type ThemeState = {
  mode: 'system' | 'light' | 'dark'
  primaryColor: string
  lightScheme: 'white' | 'gray'
  darkScheme: 'mint' | 'navy' | 'mirage' | 'cinder' | 'black'
  notificationStyle: 'stacked' | 'expanded'
  notificationMaxCount: number
  notificationPosition: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left'
  cardSkin: 'bordered' | 'shadow' | 'flat'
  monochromeMode: boolean
}

const STORAGE_KEY = 'ui:theme'

const DEFAULT_STATE: ThemeState = {
  mode: 'system',
  primaryColor: '#3b82f6',
  lightScheme: 'white',
  darkScheme: 'mint',
  notificationStyle: 'stacked',
  notificationMaxCount: 2,
  notificationPosition: 'bottom-right',
  cardSkin: 'bordered',
  monochromeMode: false,
}

const PRIMARY_COLORS = [
  { id: 'blue', color: '#3b82f6', name: 'Blue' },
  { id: 'purple', color: '#8b5cf6', name: 'Purple' },
  { id: 'orange', color: '#f97316', name: 'Orange' },
  { id: 'red', color: '#ef4444', name: 'Red' },
  { id: 'green', color: '#22c55e', name: 'Green' },
]

const LIGHT_SCHEMES = [
  { id: 'white', name: 'White', bg: '#ffffff', accent: '#6b7280' },
  { id: 'gray', name: 'Gray', bg: '#f9fafb', accent: '#6b7280' },
]

const DARK_SCHEMES = [
  { id: 'mint', name: 'Mint', bg: '#0f172a', accent: '#334155' },
  { id: 'navy', name: 'Navy', bg: '#1e293b', accent: '#475569' },
  { id: 'mirage', name: 'Mirage', bg: '#1f2937', accent: '#4b5563' },
  { id: 'cinder', name: 'Cinder', bg: '#27272a', accent: '#52525b' },
  { id: 'black', name: 'Black', bg: '#171717', accent: '#404040' },
]

// ============================================================
// Storage helpers
// ============================================================

function readThemeState(): ThemeState {
  if (typeof window === 'undefined') return DEFAULT_STATE
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return { ...DEFAULT_STATE, ...JSON.parse(raw) }
  } catch {}
  return DEFAULT_STATE
}

function writeThemeState(s: ThemeState) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(s))
  } catch {}
}

// Helper function to generate color shades
function generateColorShades(hex: string) {
  // Convert hex to RGB
  const r = parseInt(hex.slice(1, 3), 16)
  const g = parseInt(hex.slice(3, 5), 16)
  const b = parseInt(hex.slice(5, 7), 16)

  // Generate shades
  const shades: Record<string, string> = {}
  
  // Lighter shades (50-400)
  const lightShades = [
    { name: '50', mix: 0.95 },
    { name: '100', mix: 0.9 },
    { name: '200', mix: 0.8 },
    { name: '300', mix: 0.7 },
    { name: '400', mix: 0.6 },
  ]
  
  lightShades.forEach(({ name, mix }) => {
    const newR = Math.round(r + (255 - r) * mix)
    const newG = Math.round(g + (255 - g) * mix)
    const newB = Math.round(b + (255 - b) * mix)
    shades[name] = `#${newR.toString(16).padStart(2, '0')}${newG.toString(16).padStart(2, '0')}${newB.toString(16).padStart(2, '0')}`
  })
  
  // Base color (500)
  shades['500'] = hex
  
  // Darker shades (600-900)
  const darkShades = [
    { name: '600', mix: 0.8 },
    { name: '700', mix: 0.6 },
    { name: '800', mix: 0.4 },
    { name: '900', mix: 0.2 },
  ]
  
  darkShades.forEach(({ name, mix }) => {
    const newR = Math.round(r * mix)
    const newG = Math.round(g * mix)
    const newB = Math.round(b * mix)
    shades[name] = `#${newR.toString(16).padStart(2, '0')}${newG.toString(16).padStart(2, '0')}${newB.toString(16).padStart(2, '0')}`
  })
  
  return shades
}

function applyThemeVars(s: ThemeState) {
  const r = document.documentElement
  
  // Apply primary color - set both --primary-hex and --primary for compatibility
  r.style.setProperty('--primary-hex', s.primaryColor)
  r.style.setProperty('--primary', s.primaryColor)
  
  // Generate and apply primary color shades
  const shades = generateColorShades(s.primaryColor)
  Object.entries(shades).forEach(([shade, color]) => {
    r.style.setProperty(`--primary-${shade}`, color)
  })
  
  // Apply light/dark scheme colors
  const lightScheme = LIGHT_SCHEMES.find((x) => x.id === s.lightScheme)
  const darkScheme = DARK_SCHEMES.find((x) => x.id === s.darkScheme)
  
  if (lightScheme) {
    r.style.setProperty('--light-bg', lightScheme.bg)
    r.style.setProperty('--light-accent', lightScheme.accent)
  }
  if (darkScheme) {
    r.style.setProperty('--dark-bg', darkScheme.bg)
    r.style.setProperty('--dark-accent', darkScheme.accent)
  }
  
  // Apply notification settings
  r.style.setProperty('--notification-position', s.notificationPosition)
  r.style.setProperty('--notification-max-count', String(s.notificationMaxCount))
  r.dataset.notificationStyle = s.notificationStyle
  
  // Apply card skin
  r.dataset.cardSkin = s.cardSkin
  
  // Apply monochrome mode
  if (s.monochromeMode) {
    r.classList.add('monochrome-mode')
  } else {
    r.classList.remove('monochrome-mode')
  }
}

// ============================================================
// Component
// ============================================================

export default function SettingsPage() {
  const t = useTranslations()
  const { theme, setTheme } = useTheme()
  const _router = useRouter()
  const _pathname = usePathname()

  const [state, setState] = React.useState<ThemeState>(() => readThemeState())
  const [mounted, setMounted] = React.useState(false)

  React.useEffect(() => {
    setMounted(true)
    applyThemeVars(state)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  // Apply theme vars whenever state changes
  React.useEffect(() => {
    if (mounted) {
      applyThemeVars(state)
    }
  }, [state, mounted])

  const update = (next: Partial<ThemeState>) => {
    const merged = { ...state, ...next }
    setState(merged)
    writeThemeState(merged)

    // Sync mode with next-themes
    if (next.mode) {
      setTheme(next.mode)
    }
    
    // Dispatch custom event for other components to listen
    if (typeof window !== 'undefined') {
      window.dispatchEvent(new CustomEvent('theme-updated'))
    }
  }

  const reset = () => {
    setState(DEFAULT_STATE)
    applyThemeVars(DEFAULT_STATE)
    writeThemeState(DEFAULT_STATE)
    setTheme(DEFAULT_STATE.mode)
    
    // Dispatch custom event for other components to listen
    if (typeof window !== 'undefined') {
      window.dispatchEvent(new CustomEvent('theme-updated'))
    }
  }

  // Sync initial theme mode
  React.useEffect(() => {
    if (mounted && theme) {
      const mode = theme as ThemeState['mode']
      if (mode !== state.mode) {
        setState((prev) => ({ ...prev, mode }))
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [mounted, theme])

  if (!mounted) return null

  return (
    <div className="h-full w-full overflow-auto p-4 sm:p-6">
      <Card className="mx-auto w-full max-w-4xl">
        <CardContent className="p-6">
          {/* Header */}
          <div className="mb-6">
            <h1 className="text-xl font-semibold">{t('appearance')}</h1>
            <p className="text-muted-foreground text-sm">
              {t('appearance_desc')}
            </p>
          </div>

          <Separator className="mb-6" />

          {/* Theme Mode */}
          <section className="mb-8">
            <h2 className="mb-1 text-base font-medium">{t('theme_mode')}</h2>
            <p className="text-muted-foreground mb-4 text-sm">
              {t('theme_mode_desc')}
            </p>
            <div className="flex flex-wrap gap-4">
              {/* System */}
              <ThemeCard
                selected={state.mode === 'system'}
                onClick={() => update({ mode: 'system' })}
                label={t('system')}
              >
                <div className="relative h-full w-full overflow-hidden rounded-md">
                  <div className="absolute inset-0 bg-linear-to-br from-white to-white" style={{ clipPath: 'polygon(0 0, 100% 0, 0 100%)' }}>
                    <div className="flex flex-col gap-1.5 p-3">
                      <div className="h-2 w-16 rounded bg-primary-500" />
                      <div className="h-1.5 w-12 rounded bg-gray-300" />
                      <div className="h-1.5 w-14 rounded bg-gray-300" />
                    </div>
                  </div>
                  <div className="absolute inset-0 bg-gray-900" style={{ clipPath: 'polygon(100% 0, 100% 100%, 0 100%)' }}>
                    <div className="flex flex-col gap-1.5 p-3 pt-12">
                      <div className="h-2 w-16 rounded bg-primary-500" />
                      <div className="h-1.5 w-12 rounded bg-gray-600" />
                    </div>
                  </div>
                </div>
              </ThemeCard>

              {/* Light */}
              <ThemeCard
                selected={state.mode === 'light'}
                onClick={() => update({ mode: 'light' })}
                label={t('light')}
              >
                <div className="flex h-full flex-col gap-1.5 rounded-md bg-white p-3">
                  <div className="h-2 w-16 rounded bg-gray-300" />
                  <div className="h-1.5 w-20 rounded bg-gray-200" />
                  <div className="h-1.5 w-14 rounded bg-gray-200" />
                  <div className="mt-auto flex items-center gap-2">
                    <div className="h-3 w-3 rounded-full bg-primary-500" />
                    <div className="h-1.5 w-12 rounded bg-gray-300" />
                  </div>
                </div>
              </ThemeCard>

              {/* Dark */}
              <ThemeCard
                selected={state.mode === 'dark'}
                onClick={() => update({ mode: 'dark' })}
                label={t('dark')}
              >
                <div className="flex h-full flex-col gap-1.5 rounded-md bg-gray-900 p-3">
                  <div className="h-2 w-16 rounded bg-gray-600" />
                  <div className="h-1.5 w-20 rounded bg-gray-700" />
                  <div className="h-1.5 w-14 rounded bg-gray-700" />
                  <div className="mt-auto flex items-center gap-2">
                    <div className="h-3 w-3 rounded-full bg-primary-500" />
                    <div className="h-1.5 w-12 rounded bg-gray-600" />
                  </div>
                </div>
              </ThemeCard>
            </div>
          </section>

          <Separator className="mb-6" />

          {/* Primary Color */}
          <section className="mb-8">
            <h2 className="mb-1 text-base font-medium">{t('primary_color')}</h2>
            <p className="text-muted-foreground mb-4 text-sm">
              {t('primary_color_desc')}
            </p>
            <div className="flex flex-wrap gap-3">
              {PRIMARY_COLORS.map((c) => (
                <button
                  key={c.id}
                  onClick={() => update({ primaryColor: c.color })}
                  className={cn(
                    'flex h-12 w-12 items-center justify-center rounded-lg border-2 transition-all',
                    state.primaryColor === c.color
                      ? 'border-primary ring-2 ring-primary/20'
                      : 'border-gray-200 hover:border-gray-300 dark:border-gray-700',
                  )}
                >
                  <div
                    className="h-6 w-6 rounded-full"
                    style={{ backgroundColor: c.color }}
                  />
                </button>
              ))}
            </div>
          </section>

          <Separator className="mb-6" />

          {/* Light Color Scheme */}
          <section className="mb-8">
            <h2 className="mb-1 text-base font-medium">{t('light_color_scheme')}</h2>
            <p className="text-muted-foreground mb-4 text-sm">
              {t('light_color_scheme_desc')}
            </p>
            <div className="flex flex-wrap gap-4">
              {LIGHT_SCHEMES.map((scheme) => (
                <SchemeCard
                  key={scheme.id}
                  selected={state.lightScheme === scheme.id}
                  onClick={() => update({ lightScheme: scheme.id as ThemeState['lightScheme'] })}
                  label={t(scheme.id)}
                  bg={scheme.bg}
                  accent={scheme.accent}
                  isDark={false}
                />
              ))}
            </div>
          </section>

          <Separator className="mb-6" />

          {/* Dark Color Scheme */}
          <section className="mb-8">
            <h2 className="mb-1 text-base font-medium">{t('dark_color_scheme')}</h2>
            <p className="text-muted-foreground mb-4 text-sm">
              {t('dark_color_scheme_desc')}
            </p>
            <div className="flex flex-wrap gap-4">
              {DARK_SCHEMES.map((scheme) => (
                <SchemeCard
                  key={scheme.id}
                  selected={state.darkScheme === scheme.id}
                  onClick={() => update({ darkScheme: scheme.id as ThemeState['darkScheme'] })}
                  label={scheme.name}
                  bg={scheme.bg}
                  accent={scheme.accent}
                  isDark={true}
                />
              ))}
            </div>
          </section>

          <Separator className="mb-6" />

          {/* Notification Settings */}
          <section className="mb-8">
            <h2 className="mb-1 text-base font-medium">{t('notification_settings')}</h2>
            <p className="text-muted-foreground mb-4 text-sm">
              {t('notification_settings_desc')}
            </p>

            {/* Notification Group Style */}
            <div className="mb-6">
              <Label className="text-muted-foreground mb-3 block text-sm">{t('notification_group_style')}</Label>
              <div className="flex flex-wrap gap-4">
                <NotificationStyleCard
                  selected={state.notificationStyle === 'stacked'}
                  onClick={() => update({ notificationStyle: 'stacked' })}
                  label={t('stacked')}
                  type="stacked"
                />
                <NotificationStyleCard
                  selected={state.notificationStyle === 'expanded'}
                  onClick={() => update({ notificationStyle: 'expanded' })}
                  label={t('expanded')}
                  type="expanded"
                />
              </div>
            </div>

            {/* Notification Max Count */}
            <div className="mb-6">
              <Label className="text-muted-foreground mb-3 block text-sm">{t('notification_max_count')}</Label>
              <div className="flex gap-2">
                {[1, 2, 3, 4, 5].map((n) => (
                  <button
                    key={n}
                    onClick={() => update({ notificationMaxCount: n })}
                    className={cn(
                      'flex h-10 w-14 items-center justify-center rounded-md border text-sm font-medium transition-all',
                      state.notificationMaxCount === n
                        ? 'border-primary bg-primary/10 text-primary dark:bg-primary/20'
                        : 'border-gray-200 hover:border-gray-300 dark:border-gray-700',
                    )}
                  >
                    {n}
                  </button>
                ))}
              </div>
            </div>

            {/* Notification Position */}
            <div className="mb-6 flex items-center justify-between">
              <Label className="text-sm">{t('notification_position')}:</Label>
              <Select
                value={state.notificationPosition}
                onValueChange={(v) => update({ notificationPosition: v as ThemeState['notificationPosition'] })}
              >
                <SelectTrigger className="w-64">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="top-right">{t('top_right')}</SelectItem>
                  <SelectItem value="top-left">{t('top_left')}</SelectItem>
                  <SelectItem value="bottom-right">{t('bottom_right')}</SelectItem>
                  <SelectItem value="bottom-left">{t('bottom_left')}</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <Separator className="my-6" />

            {/* Card Skin */}
            <div className="mb-6 flex items-center justify-between">
              <Label className="text-sm">{t('card_skin')}:</Label>
              <Select
                value={state.cardSkin}
                onValueChange={(v) => update({ cardSkin: v as ThemeState['cardSkin'] })}
              >
                <SelectTrigger className="w-64">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="bordered">{t('bordered')}</SelectItem>
                  <SelectItem value="shadow">{t('shadow')}</SelectItem>
                  <SelectItem value="flat">{t('flat')}</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Monochrome Mode */}
            <div className="flex items-center justify-between">
              <Label className="text-sm">{t('theme_chrome_mode')}:</Label>
              <div className="flex items-center gap-3">
                <span className="text-muted-foreground text-sm">{t('monochrome_mode')}</span>
                <Switch
                  checked={state.monochromeMode}
                  onCheckedChange={(v) => update({ monochromeMode: v })}
                />
              </div>
            </div>
          </section>

          <Separator className="mb-6" />

          {/* Reset Button */}
          <Button onClick={reset} className="gap-2">
            {t('reset_theme')}
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}

// ============================================================
// Sub Components
// ============================================================

function ThemeCard({
  selected,
  onClick,
  label,
  children,
}: {
  selected: boolean
  onClick: () => void
  label: string
  children: React.ReactNode
}) {
  return (
    <button onClick={onClick} className="text-left">
      <div
        className={cn(
          'h-24 w-32 overflow-hidden rounded-lg border-2 transition-all',
          selected ? 'border-primary ring-2 ring-primary/20' : 'border-gray-200 dark:border-gray-700',
        )}
      >
        {children}
      </div>
      <p className="mt-2 text-center text-sm font-medium">{label}</p>
    </button>
  )
}

function SchemeCard({
  selected,
  onClick,
  label,
  bg,
  accent,
  isDark,
}: {
  selected: boolean
  onClick: () => void
  label: string
  bg: string
  accent: string
  isDark: boolean
}) {
  // Calculate border color for dark schemes - use primary color when selected
  const borderColor = isDark 
    ? (selected ? 'var(--primary)' : '#4b5563') // Use primary for selected, gray-600 for unselected in dark mode
    : (selected ? 'var(--primary)' : '#e5e7eb') // Use primary for selected, gray-200 for unselected in light mode
  
  return (
    <button onClick={onClick} className="text-left">
      <div
        className={cn(
          'h-24 w-28 overflow-hidden rounded-lg border-2 p-3 transition-all',
          selected && 'ring-2 ring-primary/20 dark:ring-primary/30',
        )}
        style={{ 
          backgroundColor: bg,
          borderColor: borderColor,
        }}
      >
        <div className="flex flex-col gap-1.5">
          <div className="h-2 w-16 rounded" style={{ backgroundColor: accent }} />
          <div className="h-1.5 w-12 rounded" style={{ backgroundColor: accent, opacity: 0.6 }} />
          <div className="h-1.5 w-14 rounded" style={{ backgroundColor: accent, opacity: 0.6 }} />
        </div>
        <div className="mt-auto flex items-center gap-2 pt-3">
          <div className="h-3 w-3 rounded-full" style={{ backgroundColor: isDark ? '#64748b' : '#94a3b8' }} />
          <div className="h-1.5 w-10 rounded" style={{ backgroundColor: accent, opacity: 0.5 }} />
        </div>
      </div>
      <p className="mt-2 text-center text-sm font-medium">{label}</p>
    </button>
  )
}

function NotificationStyleCard({
  selected,
  onClick,
  label,
  type,
}: {
  selected: boolean
  onClick: () => void
  label: string
  type: 'stacked' | 'expanded'
}) {
  return (
    <button onClick={onClick} className="text-left">
      <div
        className={cn(
          'h-40 w-52 overflow-hidden rounded-lg border-2 bg-gray-50 p-4 transition-all dark:bg-gray-900',
          selected ? 'border-primary ring-2 ring-primary/20' : 'border-gray-200 dark:border-gray-700',
        )}
      >
        {type === 'stacked' ? (
          <div className="relative">
            <div className="absolute top-0 left-0 right-0 h-12 rounded-md border bg-white p-2 shadow-sm dark:bg-gray-800">
              <div className="h-2 w-20 rounded bg-gray-300 dark:bg-gray-600" />
              <div className="mt-1 h-1.5 w-28 rounded bg-gray-200 dark:bg-gray-700" />
            </div>
            <div className="absolute top-3 left-1 right-1 h-12 rounded-md border bg-white p-2 shadow-sm dark:bg-gray-800">
              <div className="h-2 w-24 rounded bg-gray-300 dark:bg-gray-600" />
              <div className="mt-1 h-1.5 w-32 rounded bg-gray-200 dark:bg-gray-700" />
            </div>
            <div className="absolute top-6 left-2 right-2 h-12 rounded-md border bg-white p-2 shadow-sm dark:bg-gray-800">
              <div className="h-2 w-16 rounded bg-primary-400" />
              <div className="mt-1 h-1.5 w-24 rounded bg-gray-200 dark:bg-gray-700" />
            </div>
          </div>
        ) : (
          <div className="flex flex-col gap-2">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="h-7 rounded-md border bg-white p-1.5 dark:bg-gray-800">
                <div className="h-1.5 w-20 rounded bg-gray-300 dark:bg-gray-600" />
                <div className="mt-0.5 h-1 w-28 rounded bg-gray-200 dark:bg-gray-700" />
              </div>
            ))}
          </div>
        )}
      </div>
      <p className="mt-2 text-center text-sm font-medium">{label}</p>
    </button>
  )
}
