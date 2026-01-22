'use client'

import React, { useEffect, useMemo, useRef, useState } from 'react'
import { useForm, type SubmitHandler } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import { useTranslations } from 'next-intl'
import api from '@/lib/api'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
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
import { Checkbox } from '@/components/ui/checkbox'
import { Progress } from '@/components/ui/progress'
import { Badge } from '@/components/ui/badge'
import {
  Plus,
  Pencil,
  Trash2,
  Loader2,
  List,
  Search,
  Layers,
} from 'lucide-react'

/** ---------- Types ---------- */
type RoleRow = App.Role & { code: string }

/** ---------- Schemas (system_id байхгүй — store-оос авч илгээнэ) ---------- */
const RoleSchema = z.object({
  code: z
    .string()
    .min(2, 'Required')
    .max(64, 'Max 64')
    .regex(/^[A-Z0-9._-]+$/, 'Use A-Z, 0-9, dot, underscore, dash'),
  name: z.string().min(2, 'Required').max(255),
  description: z.string().max(255).optional().or(z.literal('')),
})
type RoleFormIn = z.input<typeof RoleSchema>
type RoleFormOut = z.output<typeof RoleSchema>

/** ---------- Small util ---------- */
const pickItems = <T extends { id: number }>(data: unknown): T[] => {
  if (!data) return []
  if (Array.isArray(data)) return data as T[]
  const d = data as { items?: T[]; data?: T[] }
  if (Array.isArray(d.items)) return d.items
  if (Array.isArray(d.data)) return d.data
  return []
}

export default function RolesBySystem() {
  const t = useTranslations()
  const selectedSystem = null as App.System | null

  // ==== List state ====
  const [rows, setRows] = useState<RoleRow[]>([])
  const [loading, setLoading] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)

  // meta/pagination
  const [meta, setMeta] = useState<App.ApiMeta | null>(null)
  const [pageNumber, setPageNumber] = useState(1)
  const [pageSize, setPageSize] = useState(50)
  const [totalPage, setTotalPage] = useState(1)
  const [totalRow, setTotalRow] = useState(0)

  // filters
  const [filterName, setFilterName] = useState('')

  // progress bar
  const [progress, setProgress] = useState(0)
  const progressTimer = useRef<ReturnType<typeof setInterval> | null>(null)
  useEffect(() => {
    if (progressTimer.current) clearInterval(progressTimer.current)
    let timeoutId: ReturnType<typeof setTimeout> | null = null
    if (loading || deleting) {
      setProgress(0)
      progressTimer.current = setInterval(
        () => setProgress((p) => Math.min(p + Math.random() * 12 + 8, 90)),
        250,
      )
    } else {
      setProgress(100)
      timeoutId = setTimeout(() => setProgress(0), 300)
    }
    return () => {
      if (progressTimer.current) clearInterval(progressTimer.current)
      if (timeoutId) clearTimeout(timeoutId)
    }
  }, [loading, deleting])

  // ==== Load roles (selectedSystem-с хамаарна) ====
  async function load(page = pageNumber, size = pageSize) {
    if (!selectedSystem?.id) {
      setRows([])
      setMeta(null)
      setTotalPage(1)
      setTotalRow(0)
      return
    }
    setLoading(true)
    setFetchError(null)
    try {
      const data = await api.get<App.ListData<RoleRow>>('/role', {
        query: {
          page,
          size,
          name: filterName || undefined, // backend name/code-оор хайж чаддаг бол OK
          system_id: selectedSystem.id,
        },
      })
      const m = data.meta
      setRows(data.items ?? [])
      setMeta(m)
      setPageNumber(m?.page ?? page)
      setPageSize(m?.size ?? size)
      setTotalPage(m?.pages ?? 1)
      setTotalRow(m?.total ?? 0)
    } catch {
      setFetchError('Failed to load')
    } finally {
      setLoading(false)
    }
  }

  // selectedSystem, filters өөрчлөгдөхөд
  useEffect(() => {
    const id = setTimeout(() => {
      setPageNumber(1)
      load(1, pageSize)
    }, 250)
    return () => clearTimeout(id)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filterName, pageSize, selectedSystem?.id])

  // ==== Dialogs & forms ====
  const [openCreate, setOpenCreate] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [selected, setSelected] = useState<RoleRow | null>(null)

  const createForm = useForm<RoleFormIn>({
    resolver: zodResolver(RoleSchema),
    defaultValues: { code: '', name: '', description: '' },
  })
  const editForm = useForm<RoleFormIn>({
    resolver: zodResolver(RoleSchema),
    defaultValues: { code: '', name: '', description: '' },
  })

  const onOpenCreate = () => {
    setOpenCreate(true)
    createForm.reset({ code: '', name: '', description: '' })
  }
  const onOpenEdit = (row: RoleRow) => {
    setSelected(row)
    editForm.reset({
      code: row.code ?? '',
      name: row.name ?? '',
      description: row.description ?? '',
    })
    setOpenEdit(true)
  }
  const onOpenDelete = (row: RoleRow) => {
    setSelected(row)
    setOpenDelete(true)
  }

  // ==== CRUD (Role) — system_id-г store-оос ====
  const onCreate: SubmitHandler<RoleFormIn> = async (valuesIn) => {
    if (!selectedSystem?.id) return
    const v: RoleFormOut = RoleSchema.parse(valuesIn)
    try {
      await api.post<RoleRow>('/role', {
        system_id: selectedSystem.id,
        code: v.code,
        name: v.name,
        description: v.description || undefined,
      } as Record<string, unknown>)
      setOpenCreate(false)
      await load(1, pageSize)
    } catch {}
  }

  const onUpdate: SubmitHandler<RoleFormIn> = async (valuesIn) => {
    if (!selected || !selectedSystem?.id) return
    const v: RoleFormOut = RoleSchema.parse(valuesIn)
    try {
      await api.put<RoleRow>(`/role/${selected.id}`, {
        id: selected.id,
        system_id: selectedSystem.id,
        code: v.code,
        name: v.name,
        description: v.description || undefined,
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

  const [openModules, setOpenModules] = useState(false)
  const [_roleModules, setRoleModules] = useState<App.Module[]>([])
  const [_modulesLoading, setModulesLoading] = useState(false)

  const [catalog, setCatalog] = useState<App.Module[]>([])
  const [catalogLoading, setCatalogLoading] = useState(false)
  const [catalogErr, setCatalogErr] = useState<string | null>(null)

  const [currentSelected, setCurrentSelected] = useState<Set<number>>(new Set())
  const initialAssignedRef = useRef<Set<number>>(new Set())
  const [saving, setSaving] = useState(false)

  const openModulesForRole = async (row: RoleRow) => {
    setSelected(row)
    setOpenModules(true)
    setCatalogErr(null)

    // Эхлээд role-д оноогдсон модуль авч, анхны set-ийг тогтооно
    const assignedList = await fetchRoleModules(row.id)
    const initial = new Set(assignedList.map((m) => m.id))
    initialAssignedRef.current = initial
    setCurrentSelected(new Set(initial)) // UI-д харагдах анхны төлөв

    // Дараа нь тухайн системийн бүх модуль каталог
    await fetchCatalog(row.system_id)
  }

  const fetchRoleModules = async (roleId: number) => {
    try {
      setModulesLoading(true)
      const data = await api.get<App.ListData<App.Module> | App.Module[]>('/rolemodule', {
        query: { id: roleId, page: 1, size: 2000 },
      })
      const list = pickItems<App.Module>(data)
      setRoleModules(list)
      return list // ✅ буцааж байна
    } finally {
      setModulesLoading(false)
    }
  }

  const fetchCatalog = async (system_id?: number | null) => {
    try {
      setCatalogLoading(true)
      const data = await api.get<App.ListData<App.Module> | App.Module[]>('/module', {
        query: { page: 1, size: 100, system_id: system_id ?? undefined },
      })
      setCatalog(pickItems<App.Module>(data))
    } catch {
      setCatalog([])
      setCatalogErr('Failed to load modules')
    } finally {
      setCatalogLoading(false)
    }
  }

  const togglePick = (id: number, on?: boolean) => {
    setCurrentSelected((prev) => {
      const next = new Set(prev)
      const willOn = on ?? !next.has(id)
      if (willOn) next.add(id)
      else next.delete(id)
      return next
    })
  }

  const selectAllVisible = () => {
    const allVisibleIds = catalog.map((m) => m.id)
    setCurrentSelected(new Set(allVisibleIds))
  }

  const clearAll = () => setCurrentSelected(new Set())

  const savePicked = async () => {
    if (!selected) return
    const initial = initialAssignedRef.current
    const desired = currentSelected

    const added = Array.from(desired).filter((id) => !initial.has(id))
    const removed = Array.from(initial).filter((id) => !desired.has(id))

    // Хэрвээ өөрчлөлтгүй бол шууд хаана
    if (added.length === 0 && removed.length === 0) {
      setOpenModules(false)
      return
    }

    try {
      setSaving(true)
      if (added.length > 0) {
        await api.post('/rolemodule', { id: selected.id, module_ids: added })
      }
      if (removed.length > 0) {
        // Таны API batch DELETE дэмждэггүй тул тус бүрийг устгана
        await Promise.all(
          removed.map((mid) => api.del('/rolemodule', { id: selected.id, module_id: mid })),
        )
      }
      // Шинэ оноолтыг уншиж UI-г шинэчилнэ
      const newAssigned = await fetchRoleModules(selected.id)
      const newSet = new Set(newAssigned.map((m) => m.id))
      initialAssignedRef.current = newSet
      setCurrentSelected(new Set(newSet))
      setOpenModules(false)
    } catch {
      setCatalogErr(t('failed_to_update'))
    } finally {
      setSaving(false)
    }
  }

  // ==== Pagination helpers ====
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

  // ============================ Render ============================
  if (!selectedSystem?.id) {
    return (
      <div className="flex h-[60vh] items-center justify-center">
        <div className="text-center">
          <Layers className="mx-auto mb-3 h-8 w-8 opacity-50" />
          <p className="text-muted-foreground">no system</p>
        </div>
      </div>
    )
  }

  return (
    <div className="h-full w-full overflow-hidden p-4 sm:p-6">
      <Card className="relative flex h-full flex-col gap-0 overflow-hidden border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-900">
        {/* Loading bar */}
        {progress > 0 && (
          <div className="absolute inset-x-0 top-0">
            <Progress value={progress} className="h-1 rounded-none" aria-label="Loading" />
          </div>
        )}

        <CardHeader className="flex flex-col gap-3 pb-4 md:flex-row md:items-center md:justify-between">
          <CardTitle className="text-lg font-medium">
            {t('role')}{' '}
            <span className="text-muted-foreground ml-2 text-sm">
              {t('system')}: <b>{selectedSystem?.name ?? selectedSystem.id}</b>
            </span>
          </CardTitle>
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
        </CardHeader>

        <Separator />

        <CardContent>
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div className="flex gap-2 py-2">
              <div className="relative">
                <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
                <Input
                  value={filterName}
                  onChange={(e) => setFilterName(e.target.value)}
                  placeholder={t('search_by_name')}
                  className="h-9 pl-8 sm:w-64"
                />
              </div>
            </div>

            {/* Rows/page — shadcn Select + i18n label (таны preference) */}
            <div className="flex items-center gap-x-2">
              <p className="text-muted-foreground text-sm">{t('rows')}</p>
              <div className="border-input bg-background h-9 w-[84px] rounded-md border px-2">
                <select
                  className="h-full w-full bg-transparent text-sm outline-none"
                  value={String(pageSize)}
                  onChange={(e) => setPageSize(Number(e.target.value))}
                >
                  {[5, 10, 20, 50, 100].map((s) => (
                    <option key={s} value={s}>
                      {s}
                    </option>
                  ))}
                </select>
              </div>
            </div>
          </div>
        </CardContent>

        <Separator />

        {/* Table */}
        <CardContent className="flex min-h-0 flex-1 flex-col overflow-hidden p-0 py-0">
          {fetchError ? (
            <div className="p-6 text-sm text-red-600">{fetchError}</div>
          ) : (
            <div className="flex min-h-0 flex-1 flex-col">
              {/* Header */}
              <div className="min-w-full overflow-x-auto">
                <Table className="w-full table-fixed">
                  <colgroup>
                    <col style={{ width: '200px' }} />
                    <col style={{ width: '120px' }} />
                    <col style={{ width: '160px' }} />
                    <col />
                    <col style={{ width: '200px' }} />
                    <col />
                  </colgroup>
                  <TableHeader>
                    <TableRow className="[&>th]:bg-background [&>th]:z-20">
                      <TableHead className="text-left">{t('actions')}</TableHead>
                      <TableHead>ID</TableHead>
                      <TableHead>{t('code')}</TableHead>
                      <TableHead>{t('name')}</TableHead>
                      <TableHead>{t('system')}</TableHead>
                      <TableHead>{t('description')}</TableHead>
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
                        <col style={{ width: '200px' }} />
                        <col style={{ width: '120px' }} />
                        <col style={{ width: '160px' }} />
                        <col />
                        <col style={{ width: '200px' }} />
                        <col />
                      </colgroup>
                      <TableBody>
                        {rows.map((row) => {
                          const busy = isRowBusy(row.id)
                          return (
                            <TableRow key={row.id} className="[&>td]:align-center">
                              <TableCell className="text-left">
                                {/* Modules */}
                                <Tooltip>
                                  <TooltipTrigger asChild>
                                    <Button
                                      variant="default"
                                      size="sm"
                                      className="mr-2 gap-1"
                                      onClick={() => {
                                        openModulesForRole(row)
                                        setSelected(row)
                                      }}
                                      aria-label={`Modules of ${row.name}`}
                                      disabled={deleting || isCreating || isUpdating}
                                    >
                                      <List className="h-4 w-4" />
                                    </Button>
                                  </TooltipTrigger>
                                  <TooltipContent>
                                    <p className="lowercase first-letter:uppercase">
                                      {t('module')}
                                    </p>
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
                              <TableCell className="text-xs">
                                <Badge variant="outline" className="rounded-md">
                                  {row.code}
                                </Badge>
                              </TableCell>
                              <TableCell className="font-medium">{row.name}</TableCell>
                              <TableCell>
                                {row.system ? (
                                  <span>{row.system.name}</span>
                                ) : (
                                  <span className="opacity-60">—</span>
                                )}
                              </TableCell>
                              <TableCell className="text-muted-foreground">
                                {row.description || <span className="opacity-60">—</span>}
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
        </CardContent>

        <Separator />

        {/* Pagination footer */}
        <CardFooter className="flex items-center justify-between gap-3 pt-4">
          <div className="text-muted-foreground w-auto min-w-72 text-sm">
            {t.rich('showing', {
              from: () => <span className="font-medium">{(meta?.start_idx ?? -1) + 1 || 0}</span>,
              to: () => <span className="font-medium">{(meta?.end_idx ?? -1) + 1 || 0}</span>,
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
        </CardFooter>
      </Card>

      <Dialog
        open={openModules}
        onOpenChange={(v) => {
          setOpenModules(v)
          if (!v) {
            setCatalogErr(null)
          }
        }}
      >
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle className="lowercase first-letter:uppercase">
              {t.rich('role_module', {
                name: () => <span className="font-semibold uppercase">{selected?.name}</span>,
              })}
            </DialogTitle>
          </DialogHeader>

          {/* Bulk controls */}
          <div className="mb-3 flex items-center justify-between gap-2">
            <div className="text-muted-foreground text-sm">
              {catalogLoading ? t('loading') : t('total') + ': ' + (catalog?.length ?? 0)}
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={selectAllVisible}
                disabled={catalogLoading}
              >
                {t('select_all')}
              </Button>
              <Button variant="outline" size="sm" onClick={clearAll}>
                {t('clear')}
              </Button>
              <span className="rounded bg-emerald-500/10 px-2 py-1 text-sm text-emerald-700">
                {t('selected')}: <b>{currentSelected.size}</b>
              </span>
            </div>
          </div>

          {catalogErr && <div className="text-sm text-red-600">{catalogErr}</div>}

          <Separator />
          {catalogLoading ? (
            <div className="py-10 text-center text-sm opacity-70">
              <Loader2 className="mx-auto mb-2 h-4 w-4 animate-spin" />
              {t('loading')}
            </div>
          ) : catalog.length === 0 ? (
            <div className="py-10 text-center text-sm opacity-70">
              {t('no_information_available')}
            </div>
          ) : (
            <div className="grid max-h-[60vh] grid-cols-2 gap-2 overflow-y-auto">
              {catalog.map((m) => {
                const isInitiallyAssigned = initialAssignedRef.current.has(m.id)
                const isSelected = currentSelected.has(m.id)

                return (
                  <Card
                    key={m.id}
                    role="button"
                    onClick={() => togglePick(m.id)}
                    className={`cursor-pointer rounded-md transition hover:shadow-md`}
                  >
                    <CardContent className="-my-4 flex items-center justify-between gap-2 px-4">
                      <div className="flex items-center gap-2">
                        <Checkbox
                          checked={isSelected}
                          onCheckedChange={(v) => togglePick(m.id, Boolean(v))}
                          aria-label={`pick-${m.id}`}
                          className="mt-1"
                        />
                        <div className="min-w-0">
                          <div className="truncate font-medium">{m.name}</div>
                          <div className="text-muted-foreground text-xs">{m.code}</div>

                          {isInitiallyAssigned && isSelected && (
                            <div className="text-muted-foreground text-xs">( {t('selected')} )</div>
                          )}
                          {isInitiallyAssigned && !isSelected && (
                            <div className="text-xs text-red-600">( {t('delete')} )</div>
                          )}
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                )
              })}
            </div>
          )}

          <DialogFooter className="gap-2">
            <Button variant="outline" onClick={() => setOpenModules(false)} disabled={saving}>
              {t('cancel')}
            </Button>
            <Button onClick={savePicked} disabled={saving}>
              {saving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
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
            <form onSubmit={createForm.handleSubmit(onCreate)} className="space-y-4">
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

              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setOpenCreate(false)}>
                  {t('cancel')}
                </Button>
                <Button type="submit" disabled={createForm.formState.isSubmitting}>
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
            <form onSubmit={editForm.handleSubmit(onUpdate)} className="space-y-4">
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

              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setOpenEdit(false)}>
                  {t('cancel')}
                </Button>
                <Button type="submit" disabled={editForm.formState.isSubmitting}>
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
            <Button variant="outline" onClick={() => setOpenDelete(false)}>
              {t('cancel')}
            </Button>
            <Button variant="destructive" onClick={onDelete} disabled={deleting}>
              {deleting ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
              <p className="capitalize">{t('delete', { name: '' })}</p>
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
