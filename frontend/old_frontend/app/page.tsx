/**
 * 🏠 Root Page (/page.tsx)
 * 
 * Энэ нь төслийн үндсэн root хуудас юм.
 * Зорилго: Хэрэглэгчийг хэл тохируулсан нүүр хуудас руу автоматаар чиглүүлэх
 * 
 * Үйл ажиллагаа:
 * 1. Cookie-с хэрэглэгчийн сонгосон хэлийг уншина (NEXT_LOCALE)
 * 2. Хэрэв хүчинтэй хэл биш бол өгөгдмөл хэл (mn) ашиглана
 * 3. Тухайн хэл дээрх нүүр хуудас руу redirect хийнэ: /{locale}/home
 * 
 * Жишээ: / → /mn/home эсвэл /en/home
 * 
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

// 🔧 Server-side динамик рендерлэх тохиргоо
export const prerender = false
export const dynamic = 'force-dynamic'

import { cookies } from 'next/headers'
import { redirect } from 'next/navigation'
import { defaultLocale, isLocale } from '@/i18n/config'

export default async function Page() {
  // 🔹 Cookie-с хэрэглэгчийн хэлний сонголтыг уншина
  const cookieLocale = (await cookies()).get('NEXT_LOCALE')?.value
  
  // 🔹 Хүчинтэй хэл эсэхийг шалгаад өгөгдмөл хэл ашиглана
  const locale = isLocale(cookieLocale ?? '')
    ? (cookieLocale as typeof defaultLocale)
    : defaultLocale

  // 🔹 Хэл тохируулсан нүүр хуудас руу чиглүүлэх
  redirect(`/${locale}/home`)
}
