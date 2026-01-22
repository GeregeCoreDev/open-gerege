/**
 * 📦 Module Page (/[locale]/(main)/admin/module/page.tsx)
 *
 * Энэ нь модуль удирдах админ хуудас юм.
 * Зорилго: Sidebar navigation модулиудын CRUD удирдлага
 *
 * Features:
 * - ✅ Full CRUD operations
 * - ✅ Pagination with server-side meta
 * - ✅ Triple filter: Name + System + Module Group
 * - ✅ Progress bar loading
 * - ✅ Active/Inactive toggle
 * - ✅ Translation key support
 * - ✅ Form validation (Zod)
 * - ✅ Hierarchical selection (System → Group)
 *
 * Table Columns:
 * - Actions (Edit/Delete)
 * - Name
 * - Translation Key
 * - System → Group
 * - Active status
 *
 * Form Fields:
 * - code: Module code (required)
 * - key: Translation key (required)
 * - name: Display name (required)
 * - description: Optional
 * - system_id: Parent system (dropdown)
 * - is_active: Enable/Disable
 *
 * Related Entities:
 * - System: Grandparent entity
 * - Permission: Child entity
 *
 * Navigation Structure:
 * System → Module
 * (sidebar hierarchy)
 *
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

'use client'

import React, { useEffect, useMemo, useState } from 'react'
import { useForm, type SubmitHandler } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
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
import { Textarea } from '@/components/ui/textarea'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationPrevious,
  PaginationNext,
  PaginationLink,
  PaginationEllipsis,
} from '@/components/ui/pagination'
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from '@/components/ui/select'
import { Plus, Pencil, Trash2, Loader2, Shield } from 'lucide-react'
import { Progress } from '@/components/ui/progress'
import api from '@/lib/api'
import { useTranslations } from 'next-intl'
import { Badge } from '@/components/ui/badge'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'
import { appConfig } from '@/config/app.config'
import { PermissionsManager } from './permissions-manager'
import { LucideIcon } from '@/lib/utils/icon'
import { VisuallyHidden } from '@radix-ui/react-visually-hidden'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { TooltipArrow } from '@radix-ui/react-tooltip'

// --------- Zod schema (Create/Update payload) ----------
const ModuleSchema = z.object({
  code: z.string().min(1, 'Required').max(255),
  name: z.string().min(1, 'Required').max(255),
  description: z.string().max(255).optional().or(z.literal('')),
  is_active: z.boolean().nullable().optional(),
  system_id: z.coerce.number().int().min(1, 'Required'),
})
type ModuleForm = z.input<typeof ModuleSchema>

export default function ModulesPage() {
  const t = useTranslations()

  // table rows
  const [rows, setRows] = useState<App.Module[]>([])
  const [systems, setSystems] = useState<App.System[]>([])
  const [loading, setLoading] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)

  // meta for pagination
  const [meta, setMeta] = useState<App.ApiMeta | null>(null)

  // pagination ui
  const [pageNumber, setPageNumber] = useState(1)
  const [pageSize, setPageSize] = useState<number>(appConfig.pagination.defaultPageSize)
  const [totalPage, setTotalPage] = useState(1)
  const [totalRow, setTotalRow] = useState(0)

  // filters
  const [filterName, setFilterName] = useState('')
  const [filterSystem, setFilterSystem] = useState<number | 'all'>('all')

  // progress bar - simplified with hook
  const progress = useLoadingProgress(loading)

  // dialogs
  const [openCreate, setOpenCreate] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [openPermissions, setOpenPermissions] = useState(false)
  const [selected, setSelected] = useState<App.Module | null>(null)
  const [_selectedModuleForActions, _setSelectedModuleForActions] = useState<App.Module | null>(null)
  const [selectedModuleForPermissions, setSelectedModuleForPermissions] =
    useState<App.Module | null>(null)

  // ---------- forms ----------
  const createForm = useForm<ModuleForm>({
    resolver: zodResolver(ModuleSchema),
    defaultValues: {
      code: '',
      name: '',
      description: '',
      is_active: true,
      system_id: 0,
    },
  })

  const editForm = useForm<ModuleForm>({
    resolver: zodResolver(ModuleSchema),
    defaultValues: {
      code: '',
      name: '',
      description: '',
      is_active: true,
      system_id: 0,
    },
  })

  // ---------- load systems + module groups ----------
  async function loadSystems() {
    try {
      const data = await api.get<App.ListData<App.System>>('/system', {
        query: { page: 1, size: 500, is_active: true },
        cache: 'no-store',
      })
      setSystems(data.items ?? [])
    } catch {
      setSystems([])
    }
  }

  // ---------- main table load ----------
  async function load(page = pageNumber, size = pageSize, system = filterSystem) {
    setLoading(true)
    setFetchError(null)
    try {
      const data = await api.get<App.ListData<App.Module>>('/module', {
        query: {
          page,
          size,
          name: filterName || undefined,
          system_id: system === 'all' ? undefined : system,
        },
      })
      const m = data.meta
      setRows(data.items ?? [])
      setMeta(m)
      setPageNumber(m.page)
      setPageSize(m.size)
      setTotalPage(m.pages)
      setTotalRow(m.total)
    } catch {
      setFetchError("Error occurred")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadSystems()
  }, [])

  useEffect(() => {
    load(1, pageSize)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageSize])

  // ---------- CRUD ----------
  const onCreate: SubmitHandler<ModuleForm> = async (valuesIn) => {
    try {
      const v = ModuleSchema.parse(valuesIn)
      await api.post<App.Module>('/module', v)
      setOpenCreate(false)
      await load(1, pageSize)

      createForm.reset({
        code: '',
        name: '',
        description: '',
        is_active: true,
        system_id: 0,
      })
    } catch {}
  }

  const onUpdate: SubmitHandler<ModuleForm> = async (valuesIn) => {
    if (!selected) return
    try {
      const v = ModuleSchema.parse(valuesIn)
      await api.put<App.Module>(`/module/${selected.id}`, v)
      setOpenEdit(false)
      setSelected(null)
      await load(pageNumber, pageSize)
    } catch {}
  }

  const onDelete = async () => {
    if (!selected) return
    try {
      setDeleting(true)
      await api.del<void>(`/module/${selected.id}`)
      const newCount = rows.length - 1
      const willBeEmpty = newCount <= 0 && pageNumber > 1
      setOpenDelete(false)
      setSelected(null)
      await load(willBeEmpty ? pageNumber - 1 : pageNumber, pageSize)
    } catch {
    } finally {
      setDeleting(false)
    }
  }

  const onOpenPermissions = (module: App.Module) => {
    setSelectedModuleForPermissions(module)
    setOpenPermissions(true)
  }

  // ---------- UI helpers ----------
  const canPrev = meta?.has_prev ?? pageNumber > 1
  const canNext = meta?.has_next ?? pageNumber < totalPage
  const goPage = (n: number) => {
    const target = Math.min(Math.max(1, n), Math.max(1, totalPage))
    if (target !== pageNumber) load(target, pageSize, filterSystem)
  }

  const pageLinks = useMemo(() => {
    const links: (number | 'ellipsis')[] = []
    const tp = totalPage
    if (tp <= 7) {
      for (let i = 1; i <= tp; i++) links.push(i)
      return links
    }
    const windowSize = 2
    links.push(1)
    if (pageNumber > 1 + windowSize + 1) links.push('ellipsis')
    const start = Math.max(2, pageNumber - windowSize)
    const end = Math.min(tp - 1, pageNumber + windowSize)
    for (let i = start; i <= end; i++) links.push(i)
    if (pageNumber < tp - windowSize - 1) links.push('ellipsis')
    links.push(tp)
    return links
  }, [pageNumber, totalPage])

  const showingFrom = totalRow === 0 ? 0 : (meta?.start_idx ?? 0) + 1
  const showingTo = totalRow === 0 ? 0 : (meta?.end_idx ?? -1) + 1

  const ActiveBadge = ({ value }: { value: boolean | null | undefined }) => {
    const on = value === true
    const off = value === false
    return (
      <span
        className={[
          'inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs',
          on
            ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300'
            : off
              ? 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300'
              : 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-300',
        ].join(' ')}
      >
        <span
          className={[
            'h-1.5 w-1.5 rounded-full',
            on ? 'bg-emerald-500' : off ? 'bg-rose-500' : 'bg-gray-400',
          ].join(' ')}
        />
        {on ? t('active') : off ? t('inactive') : ''}
      </span>
    )
  }

  // col widths
  const headerCols = [120, 40, 120, 240, 200, 200, 180, 0]
  const bodyCols = headerCols

  return (
    <>
      <div className="h-full w-full">
        <div className="relative flex h-full flex-col overflow-hidden">
          {progress > 0 && (
            <div className="absolute inset-x-0 top-0 z-10">
              <Progress value={progress} className="h-1 rounded-none" aria-label="Уншиж байна" />
            </div>
          )}

          <div className="flex flex-col overflow-hidden px-4 pb-4">
            {/* Header */}
            <div className="flex flex-col gap-3 pt-4 md:flex-row md:items-center md:justify-between">
              <h1 className="text-xl font-semibold text-gray-900 dark:text-white">{t('module')}</h1>
              <Button onClick={() => setOpenCreate(true)} className="gap-2" disabled={loading}>
                <Plus className="h-4 w-4" />
                <span>{t('create', { name: t('module') })}</span>
              </Button>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-3 py-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex flex-1 flex-wrap items-center gap-2">
                <Input
                  value={filterName}
                  onChange={(e) => setFilterName(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && load(1, pageSize)}
                  placeholder={t('search_by_name')}
                  className="h-9 sm:w-56"
                />
                <Select
                  value={filterSystem === 'all' ? 'all' : String(filterSystem)}
                  onValueChange={async (v) => {
                    const sys = v === 'all' ? 'all' : Number(v)
                    setFilterSystem(sys)

                    await load(1, pageSize, sys)
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
            </div>

            {/* Table Content */}
            <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-t-md border border-gray-200 dark:border-gray-800">
              {fetchError ? (
                <div className="p-6 text-sm text-red-600">{fetchError}</div>
              ) : (
                <div className="flex min-h-0 flex-1 flex-col">
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
                        <TableRow className="[&>th]:z-20 [&>th]:bg-gray-50 dark:[&>th]:bg-gray-800">
                          <TableHead></TableHead>
                          <TableHead>ID</TableHead>
                          <TableHead>{t('is_active')}</TableHead>
                          <TableHead>{t('code')}</TableHead>
                          <TableHead>{t('name')}</TableHead>
                          <TableHead>{t('description')}</TableHead>
                          <TableHead>{t('system')}</TableHead>
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
                      <LucideIcon name="i-lucide-archive-x" className="h-6 w-6" />
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
                            {rows.map((row) => (
                              <TableRow key={row.id}>
                                <TableCell>
                                  <Tooltip>
                                    <TooltipTrigger asChild>
                                      <Button
                                        variant="outline"
                                        size="sm"
                                        className="mr-2"
                                        onClick={() => onOpenPermissions(row)}
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
                                  <Button
                                    variant="outline"
                                    size="sm"
                                    className="mr-2"
                                    onClick={() => {
                                      setSelected(row)
                                      editForm.reset(row)
                                      setOpenEdit(true)
                                    }}
                                  >
                                    <Pencil className="h-4 w-4" />
                                  </Button>
                                  <Button
                                    variant="destructive"
                                    size="sm"
                                    onClick={() => {
                                      setSelected(row)
                                      setOpenDelete(true)
                                    }}
                                  >
                                    <Trash2 className="h-4 w-4" />
                                  </Button>
                                </TableCell>
                                <TableCell>{row.id}</TableCell>
                                <TableCell>
                                  <ActiveBadge value={row.is_active} />
                                </TableCell>
                                <TableCell>
                                  <Badge variant="outline">{row.code}</Badge>
                                </TableCell>
                                <TableCell>{row.name}</TableCell>
                                <TableCell>{row.description || '—'}</TableCell>
                                <TableCell>{row.system?.name ?? '-'}</TableCell>
                              </TableRow>
                            ))}
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
              <div className="text-muted-foreground w-72 text-sm">
                {t.rich('showing', {
                  from: () => <span className="font-medium">{showingFrom}</span>,
                  to: () => <span className="font-medium">{showingTo}</span>,
                  total: () => <span className="font-medium">{totalRow}</span>,
                })}
              </div>
              <Pagination className="flex justify-end">
                <PaginationContent>
                  <PaginationItem>
                    <PaginationPrevious
                      onClick={() => canPrev && goPage(pageNumber - 1)}
                      className={!canPrev ? 'pointer-events-none opacity-50' : ''}
                    />
                  </PaginationItem>
                  {pageLinks.map((p, i) =>
                    p === 'ellipsis' ? (
                      <PaginationItem key={`e-${i}`}>
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
                      onClick={() => canNext && goPage(pageNumber + 1)}
                      className={!canNext ? 'pointer-events-none opacity-50' : ''}
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          </div>
        </div>

        {/* ---------- Create Dialog ---------- */}
        <Dialog open={openCreate} onOpenChange={setOpenCreate}>
          <DialogContent className="sm:max-w-lg">
            <DialogHeader>
              <DialogTitle>{t('create', { name: t('module') })}</DialogTitle>
            </DialogHeader>
            <Form {...createForm}>
              <form onSubmit={createForm.handleSubmit(onCreate)} className="space-y-4">
                <div className="grid gap-4 sm:grid-cols-2">
                  <FormField
                    control={createForm.control}
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
                  {/* system select */}
                  <FormField
                    control={createForm.control}
                    name="system_id"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('system')}</FormLabel>
                        <FormControl>
                          <Select
                            value={field.value ? String(field.value) : ''}
                            onValueChange={(v) => field.onChange(Number(v))}
                          >
                            <SelectTrigger className="w-full">
                              <SelectValue placeholder={t('select_system')} />
                            </SelectTrigger>
                            <SelectContent>
                              {systems.map((s) => (
                                <SelectItem key={s.id} value={String(s.id)}>
                                  {s.name}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* status */}
                  <FormField
                    control={createForm.control}
                    name="is_active"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('status')}</FormLabel>
                        <FormControl>
                          <Select
                            value={field.value ? 'true' : 'false'}
                            onValueChange={(v) => field.onChange(v === 'true')}
                          >
                            <SelectTrigger className="w-full">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="true">{t('active')}</SelectItem>
                              <SelectItem value="false">{t('inactive')}</SelectItem>
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <FormField
                  control={createForm.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('description')}</FormLabel>
                      <FormControl>
                        <Textarea rows={3} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <DialogFooter>
                  <Button variant="outline" onClick={() => setOpenCreate(false)}>
                    {t('cancel')}
                  </Button>
                  <Button type="submit">{t('save')}</Button>
                </DialogFooter>
              </form>
            </Form>
          </DialogContent>
        </Dialog>

        {/* ---------- Edit Dialog ---------- */}
        <Dialog open={openEdit} onOpenChange={setOpenEdit}>
          <DialogContent className="sm:max-w-lg">
            <DialogHeader>
              <DialogTitle>{t('update', { name: t('module') })}</DialogTitle>
            </DialogHeader>
            <Form {...editForm}>
              <form onSubmit={editForm.handleSubmit(onUpdate)} className="space-y-4">
                <div className="grid gap-4 sm:grid-cols-2">
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

                  {/* system select */}
                  <FormField
                    control={editForm.control}
                    name="system_id"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('system')}</FormLabel>
                        <FormControl>
                          <Select
                            value={field.value ? String(field.value) : ''}
                            onValueChange={(v) => field.onChange(Number(v))}
                          >
                            <SelectTrigger className="w-full">
                              <SelectValue placeholder={t('select_system')} />
                            </SelectTrigger>
                            <SelectContent>
                              {systems.map((s) => (
                                <SelectItem key={s.id} value={String(s.id)}>
                                  {s.name}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={editForm.control}
                    name="is_active"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('status')}</FormLabel>
                        <FormControl>
                          <Select
                            value={field.value ? 'true' : 'false'}
                            onValueChange={(v) => field.onChange(v === 'true')}
                          >
                            <SelectTrigger className="w-full">
                              <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="true">{t('active')}</SelectItem>
                              <SelectItem value="false">{t('inactive')}</SelectItem>
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
                <FormField
                  control={editForm.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('description')}</FormLabel>
                      <FormControl>
                        <Textarea rows={3} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <DialogFooter>
                  <Button variant="outline" onClick={() => setOpenEdit(false)}>
                    {t('cancel')}
                  </Button>
                  <Button type="submit">{t('save')}</Button>
                </DialogFooter>
              </form>
            </Form>
          </DialogContent>
        </Dialog>

        {/* ---------- Delete Dialog ---------- */}
        <Dialog open={openDelete} onOpenChange={setOpenDelete}>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>{t('delete', { name: t('module') })}</DialogTitle>
              <DialogDescription>
                {t.rich('delete_warning', {
                  name: () => <span className="font-medium">{selected?.name}</span>,
                })}
              </DialogDescription>
            </DialogHeader>
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

        {/* ---------- Permissions Dialog ---------- */}
        <Dialog open={openPermissions} onOpenChange={setOpenPermissions}>
          <VisuallyHidden>
            <DialogTitle>{t('permission')}</DialogTitle>
          </VisuallyHidden>
          <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-5xl">
            <PermissionsManager module={selectedModuleForPermissions} />
          </DialogContent>
        </Dialog>
      </div>
    </>
  )
}
