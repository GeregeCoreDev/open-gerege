'use client'

import * as Icons from 'lucide-react'
import { cn } from '@/lib/utils'

/**
 * 🧩 toLucideKey функц
 * Nuxt UI-ийн "i-lucide-*" форматтай icon нэрийг `lucide-react` компонентын нэршилд хөрвүүлнэ.
 *
 * Жишээ:
 * - "i-lucide-settings-2" → "Settings2"
 *
 * @param icon - "i-lucide-" урд залгасан icon-ийн нэр
 * @returns PascalCase хэлбэрийн Lucide icon нэр (ж: Settings2)
 */
function toLucideKey(icon?: string) {
  if (!icon) return undefined
  const raw = icon.replace(/^i-lucide-/, '') // "settings-2"
  const parts = raw.split(/[^a-zA-Z0-9]+/g).filter(Boolean) // ['settings','2']
  const pascal = parts.map((p) => p.charAt(0).toUpperCase() + p.slice(1)).join('') // "Settings2"
  return pascal
}

/**
 * 🧱 LucideIcon компонент
 * Icon-ийн нэрийг динамикаар авч `lucide-react` сангийн ижил нэртэй icon-г харуулна.
 * Хэрэв тухайн icon байхгүй бол `Circle` icon-ийг fallback байдлаар харуулна.
 *
 * @param name - i-lucide-* форматтай icon-ийн нэр (ж: i-lucide-settings-2)
 * @param className - Tailwind class нэр (icon-ийн хэмжээ, өнгө, margin гэх мэт)
 * @returns React элемент (icon)
 *
 * 💡 Ашиглах жишээ:
 * ```tsx
 * <LucideIcon name="i-lucide-settings-2" className="text-gray-500" />
 * ```
 */
type IconComponent = React.ComponentType<{ className?: string }>
type IconsMap = Record<string, IconComponent>

export function LucideIcon({ name, className }: { name?: string; className?: string }) {
  const key = toLucideKey(name)
  const iconsMap = Icons as unknown as IconsMap
  const Fallback = iconsMap['Circle']
  const Icon = key ? iconsMap[key] : undefined
  const Comp = Icon ?? Fallback

  return <Comp className={cn('h-4 w-4', className)} aria-hidden="true" />
}
