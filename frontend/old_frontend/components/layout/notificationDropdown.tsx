'use client'

import { useState } from 'react'
import { Bell, FileText, Mail, ShieldAlert, CheckCircle, FileCheck, Settings } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { cn } from '@/lib/utils'
import { useTranslations } from 'next-intl'

type NotificationType = 'all' | 'message' | 'task' | 'log' | 'security'

interface Notification {
  id: number
  type: 'user' | 'message' | 'security' | 'task' | 'report'
  title: string
  description: string
  time: string
  read: boolean
}

// Demo notifications
const demoNotifications: Notification[] = [
  {
    id: 1,
    type: 'user',
    title: 'User Photo Changed',
    description: 'John Doe changed his avatar ph...',
    time: '2 hours ago',
    read: false,
  },
  {
    id: 2,
    type: 'message',
    title: 'New user registered',
    description: 'Jane Doe has registered',
    time: '2 hours ago',
    read: false,
  },
  {
    id: 3,
    type: 'security',
    title: 'Security alert',
    description: 'New device login detected',
    time: '11 hours ago',
    read: false,
  },
  {
    id: 4,
    type: 'task',
    title: 'Design ERP Completed',
    description: 'Design ERP completed',
    time: 'a day ago',
    read: true,
  },
  {
    id: 5,
    type: 'report',
    title: 'Weekly Report',
    description: 'The weekly report was uploaded',
    time: '2 days ago',
    read: true,
  },
]

function getNotificationIcon(type: Notification['type']) {
  switch (type) {
    case 'user':
      return <FileText className="h-5 w-5 text-gray-500" />
    case 'message':
      return <Mail className="text-primary-500 h-5 w-5" />
    case 'security':
      return <ShieldAlert className="h-5 w-5 text-orange-500" />
    case 'task':
      return <CheckCircle className="h-5 w-5 text-green-500" />
    case 'report':
      return <FileCheck className="text-primary-500 h-5 w-5" />
    default:
      return <Bell className="h-5 w-5 text-gray-500" />
  }
}

export default function NotificationDropdown() {
  const _t = useTranslations()
  const [activeTab, setActiveTab] = useState<NotificationType>('all')
  const [notifications] = useState<Notification[]>(demoNotifications)

  const unreadCount = notifications.filter((n) => !n.read).length

  const tabs: { key: NotificationType; label: string }[] = [
    { key: 'all', label: 'All' },
    { key: 'message', label: 'Message' },
    { key: 'task', label: 'Task' },
    { key: 'log', label: 'Log' },
    { key: 'security', label: 'Security' },
  ]

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="relative h-9 w-9">
          <Bell className="h-5 w-5 text-gray-500" />
          {unreadCount > 0 && (
            <span className="absolute -top-0.5 -right-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-[10px] font-medium text-white">
              {unreadCount}
            </span>
          )}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-80 p-0">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-gray-200 px-4 py-3 dark:border-gray-700">
          <div className="flex items-center gap-2">
            <span className="text-sm font-semibold text-gray-900 dark:text-white">
              Notifications
            </span>
            <span className="bg-primary-100 text-primary-700 dark:bg-primary-900/50 dark:text-primary-400 flex h-5 min-w-5 items-center justify-center rounded-full px-1.5 text-xs font-medium">
              {unreadCount}
            </span>
          </div>
          <Button variant="ghost" size="icon" className="h-7 w-7">
            <Settings className="h-4 w-4 text-gray-500" />
          </Button>
        </div>

        {/* Tabs */}
        <div className="flex gap-1 border-b border-gray-200 px-2 py-2 dark:border-gray-700">
          {tabs.map((tab) => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={cn(
                'rounded-md px-3 py-1.5 text-xs font-medium transition-all',
                activeTab === tab.key
                  ? 'bg-primary-500 text-white dark:bg-primary-500 dark:text-white'
                  : 'text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800',
              )}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Notification List */}
        <div className="max-h-80 overflow-y-auto">
          {notifications.map((notification) => (
            <div
              key={notification.id}
              className={cn(
                'flex cursor-pointer gap-3 border-b border-gray-100 px-4 py-3 transition-colors dark:border-gray-800',
                notification.read
                  ? 'hover:bg-gray-50 dark:hover:bg-gray-800/50'
                  : 'bg-primary/10 hover:bg-primary/20 dark:bg-primary/20 dark:hover:bg-primary/30',
              )}
            >
              <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-gray-100 dark:bg-gray-800">
                {getNotificationIcon(notification.type)}
              </div>
              <div className="min-w-0 flex-1">
                <p className="truncate text-sm font-medium text-gray-900 dark:text-white">
                  {notification.title}
                </p>
                <p className="truncate text-xs text-gray-500 dark:text-gray-400">
                  {notification.description}
                </p>
                <p className="text-primary-600 dark:text-primary-400 mt-1 text-xs">
                  {notification.time}
                </p>
              </div>
            </div>
          ))}
        </div>

        {/* Footer */}
        <div className="border-t border-gray-200 px-4 py-3 dark:border-gray-700">
          <button className="text-primary-600 hover:text-primary-700 dark:text-primary-400 w-full text-center text-sm font-medium">
            Archive all notifications
          </button>
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
