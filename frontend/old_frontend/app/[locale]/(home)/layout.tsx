'use client'

import HomeHeader from '@/components/layout/homeHeader'

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="h-screen w-screen overflow-x-hidden overflow-y-auto bg-white dark:bg-gray-800">
      <HomeHeader />
      {children}
    </div>
  )
}
