/**
 * 🛡️ Role Page (/[locale]/(main)/admin/role/page.tsx)
 * 
 * Энэ нь дүр (Role) удирдах админ хуудас юм.
 * Зорилго: Системийн дүрүүдийн CRUD удирдлага, хэрэглэгчийн эрхийн удирдлага
 * 
 * Features:
 * - ✅ Full CRUD operations
 * - ✅ Pagination with server-side meta
 * - ✅ Search by name filter
 * - ✅ Filter by system
 * - ✅ Progress bar loading
 * - ✅ Active/Inactive toggle
 * - ✅ User role assignment dialog
 * - ✅ Form validation (Zod)
 * - ✅ System selection dropdown
 * 
 * Table Columns:
 * - Actions (Users/Edit/Delete)
 * - ID
 * - Code
 * - Name
 * - System (related)
 * - Description
 * - Active status
 * 
 * Form Fields:
 * - Code: Role identifier
 * - Name: Display name
 * - Description: Optional
 * - System: Parent system (dropdown)
 * - is_active: Enable/Disable
 * 
 * Related Entities:
 * - System: Parent entity
 * - User: Many-to-many через UserRole
 * - Module: Roles have module permissions
 * 
 * Related Components:
 * - UserRoleDialog: Assign users to role
 * 
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import React, { useEffect, useMemo, useRef, useState } from 'react'
import { useForm, type SubmitHandler } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useTranslations } from 'next-intl'
import api from '@/lib/api'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Separator } from '@/components/ui/separator'
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
  PaginationLink,
  PaginationEllipsis,
} from '@/components/ui/pagination'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { TooltipArrow } from '@radix-ui/react-tooltip'
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { Progress } from '@/components/ui/progress'
import { Badge } from '@/components/ui/badge'

import { Plus, Pencil, Trash2, Loader2, Search, Users2, Shield } from 'lucide-react'
import RoleUsersDialog from './actions/user'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'
import { appConfig } from '@/config/app.config'
import { Card } from '@/components/ui/card'

/** ---------- Schemas ---------- */
const RoleSchema = z.object({
  system_id: z.number().int().positive({ message: 'Required' }),
  code: z
    .string()
    .min(2, 'Required')
    .max(64, 'Max 64')
    .regex(/^[A-Z0-9._-]+$/, 'Use A-Z, 0-9, dot, underscore, dash'),
  name: z.string().min(2, 'Required').max(255),
  description: z.string().max(255).optional().or(z.literal('')),
  is_active: z.boolean().default(true),
  is_system_role: z.boolean().default(false),
})
type RoleFormIn = z.input<typeof RoleSchema>
type RoleFormOut = z.output<typeof RoleSchema>

type RoleRow = App.Role

/** ---------- Utils ---------- */
const pickItems = <T extends { id: number }>(data: unknown): T[] => {
  if (!data) return []
  if (Array.isArray(data)) return data as T[]
  const d = data as { items?: T[]; data?: T[] }
  if (Array.isArray(d.items)) return d.items
  if (Array.isArray(d.data)) return d.data
  return []
}

export default function RolesPage() {
  const t = useTranslations()

  // ==== List state ====
  const [rows, setRows] = useState<RoleRow[]>([])
  const [loading, setLoading] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)

  // meta/pagination
  const [meta, setMeta] = useState<App.ApiMeta | null>(null)
  const [pageNumber, setPageNumber] = useState(1)
  const [pageSize, setPageSize] = useState<number>(appConfig.pagination.defaultPageSize)
  const [totalPage, setTotalPage] = useState(1)
  const [totalRow, setTotalRow] = useState(0)

  // filters
  const [filterName, setFilterName] = useState('')
  const [systems, setSystems] = useState<App.System[]>([])
  const [systemFilter, setSystemFilter] = useState<'all' | number>('all')

  // progress bar - simplified with hook
  const progress = useLoadingProgress(loading || deleting)

  // systems list
  useEffect(() => {
    ;(async () => {
      try {
        const data = await api.get<App.ListData<App.System> | App.System[]>('/system', {
          query: { page: 1, size: 500, is_active: true },
          cache: 'no-store',
        })
        setSystems(pickItems<App.System>(data))
      } catch {
        setSystems([])
      }
    })()
  }, [])

  // ==== Load roles ====
  async function load(page = pageNumber, size = pageSize, system = systemFilter) {
    setLoading(true)
    setFetchError(null)
    setSystemFilter(system)
    try {
      const data = await api.get<App.ListData<RoleRow>>('/role', {
        query: {
          page,
          size,
          name: filterName || undefined,
          system_id: system === 'all' ? undefined : system,
        },
        cache: 'no-store',
      })
      const m = data.meta
      setRows(data.items ?? [])
      setMeta(m)
      setPageNumber(m?.page ?? page)
      setPageSize(m?.size ?? size)
      setTotalPage(m?.pages ?? 1)
      setTotalRow(m?.total ?? 0)
    } catch {
      setFetchError("Error occurred")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load(1, pageSize)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageSize])

  const [openCreate, setOpenCreate] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [selected, setSelected] = useState<RoleRow | null>(null)

  const createForm = useForm<RoleFormIn>({
    resolver: zodResolver(RoleSchema),
    defaultValues: {
      system_id: 0,
      code: '',
      name: '',
      description: '',
      is_active: true,
      is_system_role: false,
    },
  })
  const editForm = useForm<RoleFormIn>({
    resolver: zodResolver(RoleSchema),
    defaultValues: {
      system_id: 0,
      code: '',
      name: '',
      description: '',
      is_active: true,
      is_system_role: false,
    },
  })

  const onOpenCreate = () => {
    setOpenCreate(true)
    createForm.reset({
      system_id: systemFilter === 'all' ? 0 : Number(systemFilter), // ← системийн шүүлтүүрээс урьдчилан бөглөнө
      code: '',
      name: '',
      description: '',
      is_active: true,
      is_system_role: false,
    })
  }
  const onOpenEdit = (row: RoleRow) => {
    setSelected(row)
    editForm.reset({
      system_id: row.system_id ?? 0,
      code: row.code ?? '',
      name: row.name ?? '',
      description: row.description ?? '',
      is_active: row.is_active,
      is_system_role: Boolean(row.is_system_role),
    })
    setOpenEdit(true)
  }
  const onOpenDelete = (row: RoleRow) => {
    setSelected(row)
    setOpenDelete(true)
  }

  // ==== CRUD (Role) ====
  const onCreate: SubmitHandler<RoleFormIn> = async (valuesIn) => {
    const v: RoleFormOut = RoleSchema.parse(valuesIn)
    try {
      await api.post<RoleRow>('/role', {
        system_id: v.system_id,
        code: v.code,
        name: v.name,
        description: v.description,
        is_active: v.is_active,
        is_system_role: v.is_system_role,
      } as Record<string, unknown>)
      setOpenCreate(false)
      await load(1, pageSize)
    } catch {}
  }

  const onUpdate: SubmitHandler<RoleFormIn> = async (valuesIn) => {
    if (!selected) return
    const v: RoleFormOut = RoleSchema.parse(valuesIn)
    try {
      await api.put<RoleRow>(`/role/${selected.id}`, {
        id: selected.id,
        system_id: v.system_id,
        code: v.code,
        name: v.name,
        description: v.description,
        is_active: v.is_active,
        is_system_role: v.is_system_role,
      } as Record<string, unknown>)
      setOpenEdit(false)
      setSelected(null)
      await load(pageNumber, pageSize)
    } catch {}
  }

  const onDelete = async () => {
    if (!selected) return
    try {
      setDeleting(true)
      await api.del<void>(`/role/${selected.id}`)
      const willBeEmpty = rows.length - 1 <= 0 && pageNumber > 1
      setOpenDelete(false)
      setSelected(null)
      await load(willBeEmpty ? pageNumber - 1 : pageNumber, pageSize)
    } catch {
    } finally {
      setDeleting(false)
    }
  }

  /** ===== Role ⇄ Permissions dialog ===== */
  const [openPermissions, setOpenPermissions] = useState(false)
  const [_permissionsLoading, setPermissionsLoading] = useState(false)
  const [permissionsCatalog, setPermissionsCatalog] = useState<App.Permission[]>([])
  const [permissionsCatalogLoading, setPermissionsCatalogLoading] = useState(false)
  const [permissionsCatalogErr, setPermissionsCatalogErr] = useState<string | null>(null)
  const [currentSelectedPermissions, setCurrentSelectedPermissions] = useState<Set<number>>(
    new Set(),
  )
  const initialAssignedPermissionsRef = useRef<Set<number>>(new Set())
  const [savingPermissions, setSavingPermissions] = useState(false)

  const openPermissionsForRole = async (row: RoleRow) => {
    setCurrentSelectedPermissions(new Set())
    setSelected(row)
    setOpenPermissions(true)
    setPermissionsCatalogErr(null)
    await fetchPermissionsCatalog(row.system_id)
    const assigned = await fetchRolePermissions(row.id)
    const initial = new Set(assigned.map((p) => p.id))
    initialAssignedPermissionsRef.current = initial
    setCurrentSelectedPermissions(new Set(initial))
  }

  const fetchRolePermissions = async (roleId: number) => {
    try {
      setPermissionsLoading(true)
      const data = await api.get<App.ListData<App.Permission> | App.Permission[]>(
        '/role/permissions',
        {
          query: { role_id: roleId, page: 1, size: 2000 },
          cache: 'no-store',
        },
      )
      return pickItems<App.Permission>(data)
    } catch {
      return []
    } finally {
      setPermissionsLoading(false)
    }
  }

  const fetchPermissionsCatalog = async (system_id?: number) => {
    try {
      setPermissionsCatalogLoading(true)
      const systemId = system_id == 7 ? undefined : system_id
      const data = await api.get<App.ListData<App.Permission> | App.Permission[]>('/permission', {
        query: { page: 1, size: 500, system_id: systemId },
        cache: 'no-store',
      })
      const items = pickItems<App.Permission>(data)
      setPermissionsCatalog(items)
    } catch (error) {
      console.error('Failed to load permissions catalog:', error)
      setPermissionsCatalog([])
      setPermissionsCatalogErr("Error occurred")
    } finally {
      setPermissionsCatalogLoading(false)
    }
  }

  const togglePermissionPick = (id: number, on?: boolean) => {
    setCurrentSelectedPermissions((prev) => {
      const next = new Set(prev)
      const willOn = on ?? !next.has(id)
      if (willOn) next.add(id)
      else next.delete(id)
      return next
    })
  }
  const selectAllPermissions = () =>
    setCurrentSelectedPermissions(new Set(permissionsCatalog.map((p) => p.id)))
  const clearAllPermissions = () => setCurrentSelectedPermissions(new Set())

  const savePickedPermissions = async () => {
    if (!selected) return
    try {
      setSavingPermissions(true)
      await api.post('/role/permissions', {
        role_id: selected.id,
        permission_ids: Array.from(currentSelectedPermissions),
      } as Record<string, unknown>)
      initialAssignedPermissionsRef.current = new Set(currentSelectedPermissions)
      setOpenPermissions(false)
    } catch {
      setPermissionsCatalogErr("Error occurred")
    } finally {
      setSavingPermissions(false)
    }
  }

  /** ===== Role ⇄ Users dialog ===== */
  const [openUsers, setOpenUsers] = useState(false)

  // ==== Pagination helpers for main table ====
  const showingFrom = totalRow === 0 ? 0 : (meta?.start_idx ?? 0) + 1
  const showingTo = totalRow === 0 ? 0 : (meta?.end_idx ?? -1) + 1
  const canPrev = meta?.has_prev ?? pageNumber > 1
  const canNext = meta?.has_next ?? pageNumber < totalPage
  const goPage = (n: number) => {
    const target = Math.min(Math.max(1, n), Math.max(1, totalPage))
    if (target !== pageNumber) load(target, pageSize)
  }
  const pageLinks = useMemo(() => {
    const links: (number | 'ellipsis')[] = []
    const tp = totalPage
    if (tp <= 7) {
      for (let i = 1; i <= tp; i++) links.push(i)
      return links
    }
    const w = 2
    links.push(1)
    if (pageNumber > 1 + w + 1) links.push('ellipsis')
    const start = Math.max(2, pageNumber - w)
    const end = Math.min(tp - 1, pageNumber + w)
    for (let i = start; i <= end; i++) links.push(i)
    if (pageNumber < tp - w - 1) links.push('ellipsis')
    links.push(tp)
    return links
  }, [pageNumber, totalPage])

  const isCreating = createForm.formState.isSubmitting
  const isUpdating = editForm.formState.isSubmitting
  const isRowBusy = (rowId: string | number) =>
    (isUpdating && selected?.id === rowId) || (deleting && selected?.id === rowId)

  // ─────────────── Render ───────────────
  const headerCols: Array<number | string> = [260, 80, 160, 160, 'auto', 200, 'auto', 200]
  const bodyCols = headerCols

  return (
    <>
      <div className="h-full w-full">
        <div className="relative flex h-full flex-col overflow-hidden">
          {/* Loading bar */}
          {progress > 0 && (
            <div className="absolute inset-x-0 top-0 z-10">
              <Progress value={progress} className="h-1 rounded-none" aria-label="Уншиж байна" />
            </div>
          )}

          <div className="flex flex-col overflow-hidden px-4 pb-4">
            {/* Header */}
            <div className="flex flex-col gap-3 pt-4 md:flex-row md:items-center md:justify-between">
              <h1 className="text-xl font-semibold text-gray-900 dark:text-white">{t('role')}</h1>
              <Button
                onClick={onOpenCreate}
                className="gap-2"
                disabled={isCreating || isUpdating || deleting}
              >
                {isCreating ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Plus className="h-4 w-4" />
                )}
                <span className="lowercase first-letter:uppercase">
                  {t('create', { name: t('role') })}
                </span>
              </Button>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-3 py-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex flex-1 flex-wrap items-center gap-2">
                <div className="relative">
                  <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
                  <Input
                    value={filterName}
                    onChange={(e) => setFilterName(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && load(1, pageSize)}
                    placeholder={t('search_by_name')}
                    className="h-9 pl-8 sm:w-64"
                  />
                </div>
                {/* System filter */}
                <Select
                  value={systemFilter === 'all' ? 'all' : String(systemFilter)}
                  onValueChange={async (v) => {
                    setSystemFilter(v === 'all' ? 'all' : Number(v))
                    await load(1, pageSize, v === 'all' ? 'all' : Number(v))
                  }}
                >
                  <SelectTrigger className="h-9 w-56">
                    <SelectValue placeholder={t('system')} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">{t('all')}</SelectItem>
                    {systems.map((s) => (
                      <SelectItem key={s.id} value={String(s.id)}>
                        {s.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="flex items-center gap-x-2">
                <p className="text-muted-foreground text-sm">{t('rows')}</p>
              <Select value={String(pageSize)} onValueChange={(v) => setPageSize(Number(v))}>
                <SelectTrigger className="h-9 w-[84px]">
                  <SelectValue placeholder={pageSize} />
                </SelectTrigger>
                <SelectContent>
                  {[5, 10, 20, 50, 100].map((s) => (
                    <SelectItem key={s} value={String(s)}>
                      {s}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              </div>
            </div>

            {/* Table */}
            <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-t-md border border-gray-200 dark:border-gray-800">
          {fetchError ? (
            <div className="p-6 text-sm text-red-600">{fetchError}</div>
          ) : (
            <div className="flex min-h-0 flex-1 flex-col">
              {/* Header table */}
              <div className="min-w-full overflow-x-auto">
                <Table className="w-full table-fixed">
                  <colgroup>
                    {headerCols.map((w, i) => (
                      <col
                        key={i}
                        style={{ width: typeof w === 'number' ? `${w}px` : String(w) }}
                      />
                    ))}
                  </colgroup>
                  <TableHeader>
                    <TableRow className="[&>th]:bg-gray-50 dark:[&>th]:bg-gray-800 [&>th]:z-20">
                      <TableHead className="text-left">{t('actions')}</TableHead>
                      <TableHead>ID</TableHead>
                      <TableHead>{t('system')}</TableHead>
                      <TableHead>{t('code')}</TableHead>
                      <TableHead>{t('name')}</TableHead>
                      <TableHead>{t('description')}</TableHead>
                      <TableHead>{t('is_system_role')}</TableHead>
                      <TableHead>{t('is_active')}</TableHead>
                    </TableRow>
                  </TableHeader>
                </Table>
              </div>

              {loading ? (
                <div className="flex h-32 w-full flex-col items-center justify-center gap-y-6">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  <p className="text-gray-700 dark:text-gray-200">{t('loading')}</p>
                </div>
              ) : rows.length === 0 ? (
                <div className="text-muted-foreground flex h-48 w-full flex-col items-center justify-center gap-y-6">
                  <p className="text-gray-700 dark:text-gray-200">
                    {t('no_information_available')}
                  </p>
                </div>
              ) : (
                <div className="min-h-0 flex-1 overflow-auto">
                  <div className="min-w-full overflow-x-auto">
                    <Table className="w-full table-fixed">
                      <colgroup>
                        {bodyCols.map((w, i) => (
                          <col
                            key={i}
                            style={{ width: typeof w === 'number' ? `${w}px` : String(w) }}
                          />
                        ))}
                      </colgroup>
                      <TableBody>
                        {rows.map((row) => {
                          const busy = isRowBusy(row.id)
                          return (
                            <TableRow key={row.id} className="[&>td]:align-center">
                              <TableCell className="text-left">
                                {/* Permissions */}
                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <Button
                                      variant="default"
                                      size="sm"
                                      className="mr-2 gap-1"
                                      onClick={() => openPermissionsForRole(row)}
                                      aria-label={`Permissions of ${row.name}`}
                                      disabled={deleting || isCreating || isUpdating}
                                    >
                                      <Shield className="h-4 w-4" />
                                    </Button>
                                  </TooltipTrigger>
                                  <TooltipContent>
                                    <p className="lowercase first-letter:uppercase">
                                      {t('permission')}
                                    </p>
                                    <TooltipArrow className="fill-popover" />
                                  </TooltipContent>
                                </Tooltip>

                                {/* Users */}
                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <Button
                                      variant="default"
                                      size="sm"
                                      className="mr-2 gap-1"
                                      onClick={() => {
                                        setOpenUsers(true)
                                        setSelected(row)
                                      }}
                                      aria-label={`Users of ${row.name}`}
                                      disabled={deleting || isCreating || isUpdating}
                                    >
                                      <Users2 className="h-4 w-4" />
                                    </Button>
                                  </TooltipTrigger>
                                  <TooltipContent>
                                    <p className="lowercase first-letter:uppercase">{t('users')}</p>
                                    <TooltipArrow className="fill-popover" />
                                  </TooltipContent>
                                </Tooltip>

                                {/* Edit */}
                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <Button
                                      variant="outline"
                                      size="sm"
                                      className="mr-2 gap-1"
                                      onClick={() => onOpenEdit(row)}
                                      aria-label={`Edit ${row.name}`}
                                      disabled={busy || deleting || isCreating}
                                    >
                                      {isUpdating && selected?.id === row.id ? (
                                        <Loader2 className="h-4 w-4 animate-spin" />
                                      ) : (
                                        <Pencil className="h-4 w-4" />
                                      )}
                                    </Button>
                                  </TooltipTrigger>
                                  <TooltipContent>
                                    <p className="lowercase first-letter:uppercase">
                                      {t('update', { name: t('role') })}
                                    </p>
                                    <TooltipArrow className="fill-popover" />
                                  </TooltipContent>
                                </Tooltip>

                                {/* Delete */}
                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <Button
                                      variant="destructive"
                                      size="sm"
                                      className="gap-1"
                                      onClick={() => onOpenDelete(row)}
                                      aria-label={`Delete ${row.name}`}
                                      disabled={busy || isUpdating || isCreating}
                                    >
                                      {deleting && selected?.id === row.id ? (
                                        <Loader2 className="h-4 w-4 animate-spin" />
                                      ) : (
                                        <Trash2 className="h-4 w-4" />
                                      )}
                                    </Button>
                                  </TooltipTrigger>
                                  <TooltipContent>
                                    <p className="lowercase first-letter:uppercase">
                                      {t('delete', { name: t('role') })}
                                    </p>
                                    <TooltipArrow className="fill-popover" />
                                  </TooltipContent>
                                </Tooltip>
                              </TableCell>

                              <TableCell className="text-muted-foreground">{row.id}</TableCell>
                              <TableCell>
                                {row.system ? (
                                  <span>{row.system.name}</span>
                                ) : (
                                  <span className="opacity-60">—</span>
                                )}
                              </TableCell>
                              <TableCell className="text-xs">
                                <Badge variant="outline" className="rounded-md">
                                  {row.code}
                                </Badge>
                              </TableCell>
                              <TableCell className="font-medium">{row.name}</TableCell>
                              <TableCell className="text-muted-foreground">
                                {row.description || <span className="opacity-60">—</span>}
                              </TableCell>
                              <TableCell>
                                {row.is_system_role ? (
                                  <Badge className="bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300">
                                    {t('yes')}
                                  </Badge>
                                ) : (
                                  <Badge variant="outline" className="text-xs text-muted-foreground">
                                    {t('no')}
                                  </Badge>
                                )}
                              </TableCell>
                              <TableCell>
                                {row.is_active ? (
                                  <Badge className="bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300">
                                    {t('active')}
                                  </Badge>
                                ) : (
                                  <Badge className="bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300">
                                    {t('inactive')}
                                  </Badge>
                                )}
                              </TableCell>
                            </TableRow>
                          )
                        })}
                      </TableBody>
                    </Table>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

            {/* Pagination */}
            <div className="flex items-center justify-between gap-3 rounded-b-md border-r border-b border-l border-gray-200 px-6 py-2 dark:border-gray-800">
              <div className="text-muted-foreground w-auto min-w-72 text-sm">
                {t.rich('showing', {
                  from: () => <span className="font-medium">{showingFrom || 0}</span>,
                  to: () => <span className="font-medium">{showingTo || 0}</span>,
                  total: () => <span className="font-medium">{totalRow || 0}</span>,
                })}
              </div>
              <Pagination className="flex w-full justify-end">
                <PaginationContent>
                  <PaginationItem>
                    <PaginationPrevious
                      aria-disabled={!canPrev}
                      className={!canPrev ? 'pointer-events-none opacity-50' : ''}
                      onClick={() => canPrev && goPage(pageNumber - 1)}
                    />
                  </PaginationItem>
                  {pageLinks.map((p, idx) =>
                    p === 'ellipsis' ? (
                      <PaginationItem key={`e-${idx}`}>
                        <PaginationEllipsis />
                      </PaginationItem>
                    ) : (
                      <PaginationItem key={p}>
                        <PaginationLink isActive={p === pageNumber} onClick={() => goPage(p)}>
                          {p}
                        </PaginationLink>
                      </PaginationItem>
                    ),
                  )}
                  <PaginationItem>
                    <PaginationNext
                      aria-disabled={!canNext}
                      className={!canNext ? 'pointer-events-none opacity-50' : ''}
                      onClick={() => canNext && goPage(pageNumber + 1)}
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          </div>
        </div>

        {/* ===== Role ⇄ Permissions ===== */}
      <Dialog
        open={openPermissions}
        onOpenChange={(v) => {
          setOpenPermissions(v)
          if (!v) setPermissionsCatalogErr(null)
        }}
      >
        <DialogContent className="sm:max-w-4xl">
          <DialogHeader>
            <DialogTitle className="lowercase first-letter:uppercase">
              {selected?.name} - {t('permission')}
            </DialogTitle>
          </DialogHeader>

          {/* Bulk controls */}
          <div className="mb-3 flex items-center justify-between gap-2">
            <div className="text-muted-foreground text-sm">
              {permissionsCatalogLoading
                ? t('loading')
                : t('total') + ': ' + (permissionsCatalog?.length ?? 0)}
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={selectAllPermissions}
                disabled={permissionsCatalogLoading}
              >
                {t('select_all')}
              </Button>
              <Button variant="outline" size="sm" onClick={clearAllPermissions}>
                {t('clear')}
              </Button>
              <span className="rounded bg-emerald-500/10 px-2 py-1 text-sm text-emerald-700">
                {t('selected')}: <b>{currentSelectedPermissions.size}</b>
              </span>
            </div>
          </div>

          {permissionsCatalogErr && (
            <div className="text-sm text-red-600">{permissionsCatalogErr}</div>
          )}
          <Separator />

          {permissionsCatalogLoading ? (
            <div className="py-10 text-center text-sm opacity-70">
              <Loader2 className="mx-auto mb-2 h-4 w-4 animate-spin" />
              {t('loading')}
            </div>
          ) : permissionsCatalog.length === 0 ? (
            <div className="py-10 text-center text-sm opacity-70">
              {t('no_information_available')}
            </div>
          ) : (
            <div className="grid max-h-[60vh] grid-cols-2 gap-2 overflow-y-auto">
              {permissionsCatalog.map((p) => {
                const initially = initialAssignedPermissionsRef.current.has(p.id)
                const isSelected = currentSelectedPermissions.has(p.id)
                return (
                  <Card
                    key={p.id}
                    role="button"
                    onClick={() => togglePermissionPick(p.id)}
                    className="cursor-pointer rounded-md transition hover:shadow-md"
                  >
                    <div className="-my-4 flex items-center justify-between gap-2 px-4">
                      <div className="flex items-center gap-2">
                        <Checkbox
                          checked={isSelected}
                          onCheckedChange={(v) => togglePermissionPick(p.id, Boolean(v))}
                          aria-label={`pick-${p.id}`}
                          className="mt-1"
                        />
                        <div className="min-w-0">
                          <div className="truncate font-medium">{p.name}</div>
                          <div className="text-muted-foreground text-xs">{p.code}</div>
                          {initially && isSelected && (
                            <div className="text-muted-foreground text-xs">( {t('selected')} )</div>
                          )}
                          {initially && !isSelected && (
                            <div className="text-xs text-red-600">
                              ( {t('delete', { name: '' })} )
                            </div>
                          )}
                        </div>
                      </div>
                      <p className="text-sm text-gray-500">{p.module?.name || p.module_id}</p>
                    </div>
                  </Card>
                )
              })}
            </div>
          )}

          <DialogFooter className="gap-2">
            <Button
              variant="outline"
              onClick={() => setOpenPermissions(false)}
              disabled={savingPermissions}
            >
              {t('cancel')}
            </Button>
            <Button onClick={savePickedPermissions} disabled={savingPermissions}>
              {savingPermissions && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {t('save')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* ===== Create / Edit / Delete dialogs ===== */}
      <Dialog open={openCreate} onOpenChange={setOpenCreate}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>{t('create', { name: t('role') })}</DialogTitle>
          </DialogHeader>
          <Form {...createForm}>
            <form
              onSubmit={createForm.handleSubmit(onCreate)}
              className="space-y-4"
              autoComplete="off"
            >
              <FormField
                control={createForm.control}
                name="system_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('system')}</FormLabel>
                    <Select
                      value={field.value ? String(field.value) : ''}
                      onValueChange={(v) => field.onChange(Number(v))}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder={t('system')} />
                      </SelectTrigger>
                      <SelectContent className="max-h-60">
                        {systems.map((s) => (
                          <SelectItem key={s.id} value={String(s.id)}>
                            {s.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={createForm.control}
                name="code"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('code')}</FormLabel>
                    <FormControl>
                      <Input {...field} placeholder="ADMIN, OPS_MANAGER, ..." />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={createForm.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('name')}</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={createForm.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('description')}</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={createForm.control}
                name="is_active"
                render={({ field }) => (
                  <FormItem className="flex items-center space-x-2">
                    <FormControl>
                      <Checkbox
                        checked={field.value}
                        onCheckedChange={(v) => field.onChange(Boolean(v))}
                      />
                    </FormControl>
                    <FormLabel className="text-sm font-normal">{t('is_active')}</FormLabel>
                  </FormItem>
                )}
              />
              <FormField
                control={createForm.control}
                name="is_system_role"
                render={({ field }) => (
                  <FormItem className="flex items-center space-x-2">
                    <FormControl>
                      <Checkbox
                        checked={field.value}
                        onCheckedChange={(v) => field.onChange(Boolean(v))}
                      />
                    </FormControl>
                    <FormLabel className="text-sm font-normal">{t('is_system_role')}</FormLabel>
                  </FormItem>
                )}
              />

              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setOpenCreate(false)}>
                  {t('cancel')}
                </Button>
                <Button type="submit" disabled={createForm.formState.isSubmitting}>
                  {createForm.formState.isSubmitting && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  {t('save')}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      <Dialog open={openEdit} onOpenChange={setOpenEdit}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>{t('update', { name: t('role') })}</DialogTitle>
          </DialogHeader>
          <Form {...editForm}>
            <form
              onSubmit={editForm.handleSubmit(onUpdate)}
              className="space-y-4"
              autoComplete="off"
            >
              <FormField
                control={editForm.control}
                name="system_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('system')}</FormLabel>
                    <Select
                      value={field.value ? String(field.value) : ''}
                      onValueChange={(v) => field.onChange(Number(v))}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder={t('system')} />
                      </SelectTrigger>
                      <SelectContent className="max-h-60">
                        {systems.map((s) => (
                          <SelectItem key={s.id} value={String(s.id)}>
                            {s.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={editForm.control}
                name="code"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('code')}</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={editForm.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('name')}</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={editForm.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('description')}</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={editForm.control}
                name="is_active"
                render={({ field }) => (
                  <FormItem className="flex items-center space-x-2">
                    <FormControl>
                      <Checkbox
                        checked={field.value}
                        onCheckedChange={(v) => field.onChange(Boolean(v))}
                      />
                    </FormControl>
                    <FormLabel className="text-sm font-normal">{t('is_active')}</FormLabel>
                  </FormItem>
                )}
              />
              <FormField
                control={editForm.control}
                name="is_system_role"
                render={({ field }) => (
                  <FormItem className="flex items-center space-x-2">
                    <FormControl>
                      <Checkbox
                        checked={field.value}
                        onCheckedChange={(v) => field.onChange(Boolean(v))}
                      />
                    </FormControl>
                    <FormLabel className="text-sm font-normal">{t('is_system_role')}</FormLabel>
                  </FormItem>
                )}
              />

              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setOpenEdit(false)}>
                  {t('cancel')}
                </Button>
                <Button type="submit" disabled={editForm.formState.isSubmitting}>
                  {editForm.formState.isSubmitting && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  {t('save')}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      <Dialog open={openDelete} onOpenChange={setOpenDelete}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>{t('delete', { name: t('role') })}</DialogTitle>
          </DialogHeader>
          <p className="text-muted-foreground text-sm">
            {t.rich('delete_warning', {
              name: () => <span className="font-medium">{selected?.name}</span>,
            })}
          </p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setOpenDelete(false)} disabled={deleting}>
              {t('cancel')}
            </Button>
            <Button variant="destructive" onClick={onDelete} disabled={deleting}>
              {deleting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              <span className="capitalize">{t('delete', { name: '' })}</span>
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Users dialog */}
      <RoleUsersDialog
        open={openUsers}
        onOpenChange={setOpenUsers}
        role={selected}
        onChanged={() => {}}
      />
      </div>
    </>
  )
}
