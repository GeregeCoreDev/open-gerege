// components/ui/sidebar.tsx
'use client'

import * as React from 'react'
import { Slot } from '@radix-ui/react-slot'
import { cn } from '@/lib/utils'
import { PanelLeft } from 'lucide-react'

/* ──────────────────────────────────────────────────────────
 * Context: open/close state (collapsed sidebar)
 * ────────────────────────────────────────────────────────── */
type SidebarCtx = {
  open: boolean
  setOpen(v: boolean): void
  toggle(): void
}
const Ctx = React.createContext<SidebarCtx | null>(null)

export function SidebarProvider({
  defaultOpen = true,
  children,
}: {
  defaultOpen?: boolean
  children: React.ReactNode
}) {
  const [open, setOpen] = React.useState(defaultOpen)

  // persist optional
  React.useEffect(() => {
    if (typeof localStorage !== 'undefined') {
      const raw = localStorage.getItem('sidebar-open')
      if (raw != null) setOpen(raw !== '0')
    }
  }, [])
  React.useEffect(() => {
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('sidebar-open', open ? '1' : '0')
    }
  }, [open])

  const value = React.useMemo<SidebarCtx>(
    () => ({ open, setOpen, toggle: () => setOpen((v) => !v) }),
    [open],
  )

  return <Ctx.Provider value={value}>{children}</Ctx.Provider>
}

export function useSidebar() {
  const v = React.useContext(Ctx)
  if (!v) throw new Error('useSidebar must be used inside <SidebarProvider>')
  return v
}

/* ──────────────────────────────────────────────────────────
 * Root + regions
 * ────────────────────────────────────────────────────────── */
export function Sidebar({
  className,
  children,
}: {
  className?: string
  children: React.ReactNode
}) {
  const { open } = useSidebar()
  return (
    <aside
      className={cn(
        // container
        'relative hidden h-full shrink-0 overflow-hidden lg:flex',
        // glass + border
        // width transition
        open ? 'w-64 xl:w-72' : 'w-10',
        'transition-[width] duration-300 ease-[cubic-bezier(0.2,0.8,0.2,1)]',
        className,
      )}
    >
      <div className="flex h-full w-full flex-col">{children}</div>
    </aside>
  )
}

export function SidebarHeader({
  className,
  children,
}: {
  className?: string
  children: React.ReactNode
}) {
  return <div className={cn('flex items-center gap-2', className)}>{children}</div>
}

export function SidebarContent({
  className,
  children,
}: {
  className?: string
  children: React.ReactNode
}) {
  return <div className={cn('flex-1 overflow-y-auto', className)}>{children}</div>
}

export function SidebarFooter({
  className,
  children,
}: {
  className?: string
  children: React.ReactNode
}) {
  return (
    <div
      className={cn(
        'border-t px-2 py-2',
        'border-slate-200/70 dark:border-slate-700/70',
        className,
      )}
    >
      {children}
    </div>
  )
}

/* ──────────────────────────────────────────────────────────
 * Groups (section blocks)
 * ────────────────────────────────────────────────────────── */
export function SidebarGroup({
  className,
  children,
}: {
  className?: string
  children: React.ReactNode
}) {
  return <div className={cn('mb-1', className)}>{children}</div>
}

export function SidebarGroupLabel({
  className,
  children,
}: {
  className?: string
  children: React.ReactNode
}) {
  const { open } = useSidebar()
  return (
    <div
      className={cn(
        'flex items-center gap-2 px-2 py-2 text-xs font-semibold tracking-wide uppercase',
        'dark:text-slate-400',
        open ? 'justify-start' : 'justify-center',
        className,
      )}
    >
      {children}
    </div>
  )
}

/* ──────────────────────────────────────────────────────────
 * Menu + items
 * ────────────────────────────────────────────────────────── */
export function SidebarMenu({
  className,
  children,
}: {
  className?: string
  children: React.ReactNode
}) {
  return <ul className={cn('flex flex-col gap-0.5', className)}>{children}</ul>
}

export function SidebarMenuItem({
  className,
  children,
}: {
  className?: string
  children: React.ReactNode
}) {
  return <li className={cn('relative', className)}>{children}</li>
}

export function SidebarMenuButton({
  asChild,
  className,
  active,
  ...props
}: React.ButtonHTMLAttributes<HTMLButtonElement> & {
  asChild?: boolean
  active?: boolean
}) {
  const Comp = asChild ? Slot : 'button'
  return (
    <Comp
      className={cn(
        'group relative flex h-9 w-full items-center gap-2 rounded-lg px-3 text-sm transition-colors',
        // hover
        'hover:bg-primary-500/5 hover:text-primary-700 dark:hover:text-primary-300',
        // active
        active
          ? 'bg-primary-500/10 text-primary-700 dark:text-primary-300'
          : 'text-slate-800 dark:text-slate-200',
        className,
      )}
      {...props}
    />
  )
}

/* ──────────────────────────────────────────────────────────
 * Trigger (collapse toggle)
 * ────────────────────────────────────────────────────────── */
export function SidebarTrigger({ className }: { className?: string }) {
  const { toggle } = useSidebar()
  return (
    <button
      type="button"
      onClick={toggle}
      className={cn(
        'hidden h-8 w-8 items-center justify-center rounded-md border lg:inline-flex',
        'border-slate-200/70 bg-white/90 backdrop-blur dark:border-slate-700/70 dark:bg-slate-900/85',
        'text-slate-700 hover:bg-slate-50 dark:text-slate-200 dark:hover:bg-slate-800',
        className,
      )}
      aria-label="Toggle sidebar"
    >
      {/* Tailwind Iconify (optional). Хэрэв байхгүй бол дүрсээ солиорой */}
      <PanelLeft className="size-4" />
    </button>
  )
}

/* ──────────────────────────────────────────────────────────
 * (Optional) Rail component — одоохондоо хэрэггүй тул null
 * ────────────────────────────────────────────────────────── */
export function SidebarRail() {
  return null
}
