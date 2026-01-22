/**
 * 🌍 Locale Page (/[locale]/page.tsx)
 * 
 * Энэ нь хэл тохируулсан root хуудас юм.
 * Зорилго: Хэрэглэгчийг /mn эсвэл /en гэх мэт хэлний root-оос нүүр хуудас руу чиглүүлэх
 * 
 * Үйл ажиллагаа:
 * 1. Cookie-с хэрэглэгчийн сонгосон хэлийг уншина (NEXT_LOCALE)
 * 2. Хүчинтэй хэл эсэхийг шалгана, өгөгдмөл хэл (mn) ашиглана
 * 3. Нүүр хуудас руу redirect хийнэ: /{locale}/home
 * 
 * Жишээ: /mn → /mn/home, /en → /en/home
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

  // 🔹 Нүүр хуудас руу чиглүүлэх
  redirect(`/${locale}/home`)
}
