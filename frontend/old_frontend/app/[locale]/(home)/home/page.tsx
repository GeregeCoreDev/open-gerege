/**
 * 🏠 Home Page (/[locale]/(home)/home/page.tsx)
 *
 * Нүүр хуудас - Carousel харуулна
 */

'use client'
import { HomeCarousel } from './components/homeCarousel'

export default function Page() {
  return (
    <main className="space-y-4">
      <HomeCarousel />
    </main>
  )
}
