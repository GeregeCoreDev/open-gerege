/**
 * 📊 Admin Dashboard Page (/[locale]/(main)/admin/dashboard/page.tsx)
 *
 * Энэ нь админ системийн нүүр хуудас юм.
 * Зорилго: Системийн ерөнхий статистик, хяналтын самбар
 *
 * Features:
 * - 📈 Statistics cards:
 *   - Total users (200)
 *   - Total roles (5)
 *   - Organizations (3)
 *   - Systems (4)
 * - 🎨 Icon-based cards with colors
 * - 📱 Responsive grid layout (4 columns)
 *
 * Planned Features:
 * - Real-time statistics from API
 * - Charts and graphs
 * - Recent activities
 * - Quick actions
 * - System health monitoring
 *
 * Current: Static mock data
 * TODO: Connect to real API endpoints
 *
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

import { Separator } from '@/components/ui/separator'
import { Building, Database, Shield, Users } from 'lucide-react'

export default async function AdminDashboardPage() {
  return (
    <>
      <div className="h-full w-full p-6">
        <div className="flex flex-col gap-3 border-b border-gray-200 px-6 py-4 md:flex-row md:items-center md:justify-between dark:border-gray-800">
          <h1 className="text-xl font-semibold text-gray-900 dark:text-white">Хянах самбар</h1>
        </div>

        <Separator />

        <div>
          <div className="grid grid-cols-4 gap-4 pt-4">
            <div className="flex h-40 w-full flex-col justify-between gap-4 rounded-xl border p-4">
              <div className="flex items-center justify-between">
                <p>Нийт хэрэглэгчид</p>
                <Users className="h-5 w-5 text-primary-500" />
              </div>
              <p className="text-xl">200</p>
            </div>
            <div className="flex h-40 w-full flex-col justify-between gap-4 rounded-xl border p-4">
              <div className="flex items-center justify-between">
                <p>Нийт дүрүүд</p>
                <Shield className="h-5 w-5 text-green-500" />
              </div>
              <p className="text-xl">5</p>
            </div>
            <div className="flex h-40 w-full flex-col justify-between gap-4 rounded-xl border p-4">
              <div className="flex items-center justify-between">
                <p>Байгууллагууд</p>
                <Building className="h-5 w-5 text-purple-500" />
              </div>
              <p className="text-xl">3</p>
            </div>
            <div className="flex h-40 w-full flex-col justify-between gap-4 rounded-xl border p-4">
              <div className="flex items-center justify-between">
                <p>Системүүд</p>
                <Database className="h-5 w-5 text-orange-500" />
              </div>
              <p className="text-xl">4</p>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
