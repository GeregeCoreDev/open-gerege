/**
 * 🔄 Callback Page (/[locale]/callback/page.tsx)
 *
 * Энэ нь нэвтрэлтийн дараах буцах хуудас (OAuth callback, SSO callback гэх мэт)
 * Зорилго: Хэрэглэгч нэвтэрсний дараа зөв хуудас руу чиглүүлэх
 *
 * Үйл явц:
 * 1. Хэрэглэгчийн профайл ачаална (loadProfile)
 * 2. Дүрийн жагсаалт татна (getRoleList)
 * 4. Эхний системийн эхний модуль руу чиглүүлнэ
 * 5. Хэрэв систем, модуль байхгүй бол /profile руу чиглүүлнэ
 *
 * Жишээ navigation flow:
 * - Систем олдсон: /mn/admin/dashboard
 * - Систем олдоогүй: /mn/profile
 *
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

// 🔧 Dynamic рендерлэх тохиргоо
export const prerender = false
export const dynamic = 'force-dynamic'

import { useRouter } from 'next/navigation'
import { useLocale } from 'next-intl'
import { useEffect } from 'react'
import { useUserStore } from '@/lib/stores/user'
import { useRoleStore } from '@/lib/stores/role'
import { useMenuStore } from '@/lib/stores/menu'

export default function CallbackPage() {
  const router = useRouter()
  const locale = useLocale()
  const loadProfile = useUserStore((s) => s.loadProfile)
  const { getRoleList, roleList } = useRoleStore()
  const { getMenuList, getFirstChildPath, selectRoot } = useMenuStore()

  /**
   * 🔗 Path-г хэл тохируулсан URL болгох туслах функц
   * @param rawPath - Анхны path (жишээ: "/admin/dashboard")
   * @returns Хэл тохируулсан URL (жишээ: "/mn/admin/dashboard")
   */
  const toLocaleHref = (rawPath?: string) => {
    const p = (rawPath || '').startsWith('/') ? rawPath : `/${rawPath || ''}`
    return `/${locale}${p}`.replace(/\/{2,}/g, '/')
  }

  /**
   * ⚙️ Анхны ачаалал хийх функц
   * Хэрэглэгчийн мэдээлэл, дүр, систем ачаалаад зөв хуудас руу чиглүүлнэ
   */
  async function init() {
    // 🔹 Хэрэглэгчийн профайл ачаална
    await loadProfile()
      .catch(() => {})
      .finally(async () => {
        // 🔹 Хэрэв дүр аль хэдийн байвал дахин татахгүй
        if (roleList.length > 0) return

        // 🔹 Дүрийн жагсаалт татна
        await getRoleList()

        // 🔹 Menu жагсаалт авна
        const menus = await getMenuList()

        if (menus.length > 0) {
          // ✅ Эхний root menu-г сонгоно
          const firstRoot = menus[0]
          selectRoot(firstRoot.id)

          // 🔍 Эхний path-тай menu-г олох
          const firstPath = getFirstChildPath(firstRoot.id)
          if (firstPath) {
            // 🎯 Эхний menu руу чиглүүлнэ
            const href = toLocaleHref(firstPath)
            router.push(href)
            router.refresh()
          } else {
            // ⚠️ Path байхгүй бол профайл руу
            router.push(`/${locale}/profile`)
            router.refresh()
          }
        } else {
          // ⚠️ Menu байхгүй бол профайл руу чиглүүлнэ
          router.push(`/${locale}/profile`)
          router.refresh()
        }
      })
  }

  // 🔄 Component mount болоход init функц ажиллана
  useEffect(() => {
    init()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  // 🔹 UI харуулахгүй (зөвхөн redirect логик)
  return null
}
