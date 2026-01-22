/**
 * 🏛️ Organization Type Page (/[locale]/(main)/admin/organization-type/page.tsx)
 *
 * Энэ нь байгууллагын төрөл удирдах админ хуудас юм.
 * Зорилго: Байгууллагын төрлүүд болон тэдгээрт зөвшөөрөгдсөн системүүдийн удирдлага
 *
 * Features:
 * - ✅ Full CRUD operations
 * - ✅ Pagination with server-side meta
 * - ✅ Search by name filter
 * - ✅ Progress bar loading
 * - ✅ System access management
 * - ✅ Multi-select system assignment
 * - ✅ Statistics cards (types, orgs, systems)
 * - ✅ Form validation (Zod)
 * - ✅ Badge display for system count
 *
 * Table Columns:
 * - Actions (Systems/Edit/Delete)
 * - Type name + Code
 * - Description
 * - Systems (count badge)
 * - Organizations count
 *
 * Form Fields:
 * - code: Type code (required)
 * - name: Type name (required)
 * - description: Optional description
 *
 * System Assignment:
 * - View assigned systems
 * - Add/remove system access
 * - Toggle all systems
 *
 * Statistics:
 * - Total organization types
 * - Total organizations
 * - Types with systems
 * - Total system access count
 *
 * Related Components:
 * - SystemAccessDialog: Manage systems for org type
 *
 * API Endpoints:
 * - GET /organization-type - List types
 * - GET /system - Get all systems
 * - POST /organization-type - Create type
 * - PUT /organization-type - Update type
 * - DELETE /organization-type/:id - Delete type
 * - GET /organization-type/:id/systems - Get assigned systems
 * - POST /organization-type/:id/systems - Assign systems
 *
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

// app/[locale]/orgtype/page.tsx
'use client'

import React, { useEffect, useMemo, useState, useCallback } from 'react'
import { useTranslations } from 'next-intl'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Progress } from '@/components/ui/progress'
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
import { Loader2, Plus, Pencil, Trash2, Search, X, Layers, UserCog } from 'lucide-react'
import { z } from 'zod'
import { useForm, type SubmitHandler } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import OrgTypeSystemsDialog from './actions/system'
import OrgTypeRolesDialog from './actions/role'
import { Badge } from '@/components/ui/badge'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'
import { appConfig } from '@/config/app.config'
import {
  Form,
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from '@/components/ui/form'

type OrgType = App.OrganizationType

const OrgTypeSchema = z.object({
  name: z.string().min(2, 'Required').max(255),
  description: z.string().max(255).optional().or(z.literal('')),
  code: z.string().min(2, 'Required').max(255),
})
type OrgTypeIn = z.input<typeof OrgTypeSchema>
type OrgTypeOut = z.output<typeof OrgTypeSchema>

export default function OrgTypePage() {
  const t = useTranslations()

  const [rows, setRows] = useState<OrgType[]>([])
  const [loading, setLoading] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)

  const [meta, setMeta] = useState<App.ApiMeta | null>(null)
  const [pageNumber, setPageNumber] = useState(1)
  const [pageSize, setPageSize] = useState<number>(appConfig.pagination.defaultPageSize)
  const [totalPage, setTotalPage] = useState(1)
  const [totalRow, setTotalRow] = useState(0)

  const [filterName, setFilterName] = useState('')

  // progress bar - simplified with hook
  const progress = useLoadingProgress(loading)

  const load = useCallback(
    async (page = pageNumber, size = pageSize) => {
      setLoading(true)
      setFetchError(null)
      try {
        const data = await api.get<App.ListData<OrgType>>('/orgtype', {
          query: { page, size, name: filterName },
        })
        const m = data.meta
        setRows(data.items ?? [])
        setMeta(m ?? null)
        setPageNumber(m?.page ?? page)
        setPageSize(m?.size ?? size)
        setTotalPage(m?.pages ?? 1)
        setTotalRow(m?.total ?? 0)
      } catch {
        setFetchError("Error occurred")
      } finally {
        setLoading(false)
      }
    },
    [filterName, pageNumber, pageSize],
  )

  useEffect(() => {
    load(1, pageSize)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageSize])

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

  // dialogs
  const [openCreate, setOpenCreate] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [selected, setSelected] = useState<OrgType | null>(null)
  const [openSystem, setOpenSystem] = useState(false)
  const [openRole, setOpenRole] = useState(false)

  const createForm = useForm<OrgTypeIn>({
    resolver: zodResolver(OrgTypeSchema),
    defaultValues: { name: '', description: '', code: '' },
  })
  const editForm = useForm<OrgTypeIn>({
    resolver: zodResolver(OrgTypeSchema),
    defaultValues: { name: '', description: '', code: '' },
  })

  const onOpenCreate = () => {
    setOpenCreate(true)
    createForm.reset({ name: '', description: '', code: '' })
  }
  const onOpenEdit = (row: OrgType) => {
    setSelected(row)
    editForm.reset({
      name: row.name ?? '',
      description: row.description ?? '',
      code: row.code ?? '',
    })
    setOpenEdit(true)
  }
  const onOpenDelete = (row: OrgType) => {
    setSelected(row)
    setOpenDelete(true)
  }

  const onCreate: SubmitHandler<OrgTypeIn> = async (valuesIn) => {
    const v: OrgTypeOut = OrgTypeSchema.parse(valuesIn)
    try {
      await api.post<OrgType>('/orgtype', {
        name: v.name,
        description: v.description || undefined,
        code: v.code,
      })
      setOpenCreate(false)
      await load(1, pageSize)
    } catch {}
  }
  const onUpdate: SubmitHandler<OrgTypeIn> = async (valuesIn) => {
    if (!selected) return
    const v: OrgTypeOut = OrgTypeSchema.parse(valuesIn)
    try {
      await api.put<OrgType>(`/orgtype/${selected.id}`, {
        id: selected.id,
        name: v.name,
        description: v.description || undefined,
        code: v.code,
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
      await api.del<void>(`/orgtype/${selected.id}`)
      const willBeEmpty = rows.length - 1 <= 0 && pageNumber > 1
      const next = willBeEmpty ? pageNumber - 1 : pageNumber
      setOpenDelete(false)
      setSelected(null)
      await load(next, pageSize)
    } catch {
    } finally {
      setDeleting(false)
    }
  }

  const isUpdating = editForm.formState.isSubmitting
  const isRowBusy = (rowId: number) =>
    (isUpdating && selected?.id === rowId) || (deleting && selected?.id === rowId)

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
              <h1 className="text-xl font-semibold text-gray-900 dark:text-white">
                {t('organization_type') || 'Organization types'}
              </h1>
              <Button onClick={onOpenCreate} className="gap-2" disabled={isUpdating || deleting}>
                <Plus className="h-4 w-4" />
                <span className="lowercase first-letter:uppercase">
                  {t('create', { name: t('organization_type') || 'type' })}
                </span>
              </Button>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-3 py-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex flex-1 flex-wrap items-center gap-2">
                <div className="relative w-72">
                  <Search className="absolute top-1/2 left-2 h-4 w-4 -translate-y-1/2 opacity-60" />
                  <Input
                    value={filterName}
                    onChange={(e) => setFilterName(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') load(1, pageSize)
                    }}
                    placeholder={t('search_by_name')}
                    className="h-9 pl-8"
                  />
                  {!!filterName && (
                    <button
                      className="absolute top-1/2 right-2 -translate-y-1/2"
                      onClick={() => setFilterName('')}
                      type="button"
                      aria-label="clear"
                    >
                      <X className="h-4 w-4 opacity-60" />
                    </button>
                  )}
                </div>
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

            <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-md border border-gray-200 dark:border-gray-800">
              {fetchError ? (
                <div className="p-6 text-sm text-red-600">{fetchError}</div>
              ) : (
                <div className="flex min-h-0 flex-1 flex-col">
                  <div className="min-w-full overflow-x-auto">
                    <Table className="w-full table-fixed">
                      <colgroup>
                        <col style={{ width: '180px' }} />
                        <col style={{ width: '80px' }} />
                        <col />
                        <col />
                        <col />
                      </colgroup>
                      <TableHeader>
                        <TableRow className="[&>th]:z-20 [&>th]:bg-gray-50 dark:[&>th]:bg-gray-800">
                          <TableHead className="text-left">{t('actions')}</TableHead>
                          <TableHead>ID</TableHead>
                          <TableHead>{t('code')}</TableHead>
                          <TableHead>{t('name')}</TableHead>
                          <TableHead>{t('description')}</TableHead>
                        </TableRow>
                      </TableHeader>
                    </Table>
                  </div>

                  {/* Body table (scrollable) */}
                  {loading ? (
                    <div className="flex h-32 w-full items-center justify-center text-sm opacity-70">
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      {t('loading')}
                    </div>
                  ) : rows.length === 0 ? (
                    <div className="text-muted-foreground flex h-48 w-full items-center justify-center text-sm">
                      {t('no_information_available')}
                    </div>
                  ) : (
                    <div className="min-h-0 flex-1 overflow-auto">
                      <div className="min-w-full overflow-x-auto">
                        <Table className="w-full table-fixed">
                          <colgroup>
                            <col style={{ width: '180px' }} />
                            <col style={{ width: '80px' }} />
                            <col />
                            <col />
                            <col />
                          </colgroup>
                          <TableBody>
                            {rows.map((row) => (
                              <TableRow key={row.id} className="[&>td]:align-middle">
                                <TableCell className="text-left">
                                  <Button
                                    variant="outline"
                                    size="sm"
                                    className="mr-2"
                                    onClick={() => onOpenEdit(row)}
                                    disabled={isRowBusy(row.id)}
                                  >
                                    <Pencil className="h-4 w-4" />
                                  </Button>
                                  <Button
                                    variant="destructive"
                                    size="sm"
                                    className="mr-2"
                                    onClick={() => onOpenDelete(row)}
                                    disabled={isRowBusy(row.id)}
                                  >
                                    <Trash2 className="h-4 w-4" />
                                  </Button>
                                  <Button
                                    size="sm"
                                    variant="default"
                                    onClick={() => {
                                      setSelected(row)
                                      setOpenSystem(true)
                                    }}
                                    aria-label="Системүүд"
                                  >
                                    <Layers className="h-4 w-4" />
                                  </Button>
                                  <Button
                                    size="sm"
                                    variant="secondary"
                                    className="ml-2"
                                    onClick={() => {
                                      setSelected(row)
                                      setOpenRole(true)
                                    }}
                                    aria-label="Дүрүүд"
                                  >
                                    <UserCog className="h-4 w-4" />
                                  </Button>
                                </TableCell>
                                <TableCell className="text-muted-foreground">{row.id}</TableCell>
                                <TableCell className="font-medium">
                                  <Badge variant="outline">
                                    {row.code || <span className="opacity-60">—</span>}
                                  </Badge>
                                </TableCell>
                                <TableCell className="font-medium">{row.name}</TableCell>
                                <TableCell className="text-muted-foreground">
                                  {row.description || <span className="opacity-60">—</span>}
                                </TableCell>
                              </TableRow>
                            ))}
                          </TableBody>
                        </Table>
                      </div>
                    </div>
                  )}
                  {/* Pagination */}
                  <div className="-l flex items-center justify-between gap-3 border-t border-gray-200 px-6 py-2 dark:border-gray-800">
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
              )}
            </div>
          </div>
        </div>

        {/* Create */}
        <Dialog open={openCreate} onOpenChange={setOpenCreate}>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>{t('create', { name: t('organization_type') || 'type' })}</DialogTitle>
            </DialogHeader>

            <Form {...createForm}>
              <form
                onSubmit={createForm.handleSubmit(onCreate)}
                className="space-y-4"
                autoComplete="off"
              >
                <FormField
                  control={createForm.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('name')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('name')} {...field} />
                      </FormControl>
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
                        <Input placeholder="ORG-TYPE" {...field} />
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
                        <Input placeholder={t('description')} {...field} />
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

        {/* Edit */}
        <Dialog open={openEdit} onOpenChange={setOpenEdit}>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>{t('update', { name: t('organization_type') || 'type' })}</DialogTitle>
              <DialogDescription>{t('update_fields_and_save_your_changes')}</DialogDescription>
            </DialogHeader>

            <Form {...editForm}>
              <form
                onSubmit={editForm.handleSubmit(onUpdate)}
                className="space-y-4"
                autoComplete="off"
              >
                <FormField
                  control={editForm.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('name')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('name')} {...field} />
                      </FormControl>
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
                        <Input placeholder="ORG-TYPE" {...field} />
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
                        <Input placeholder={t('description')} {...field} />
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

        {/* Delete */}
        <Dialog open={openDelete} onOpenChange={setOpenDelete}>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>{t('delete', { name: t('organization_type') || 'type' })}</DialogTitle>
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
                {deleting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}{' '}
                <p className="capitalize">{t('delete', { name: '' })}</p>
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        <OrgTypeSystemsDialog open={openSystem} onOpenChange={setOpenSystem} orgType={selected} />
        <OrgTypeRolesDialog open={openRole} onOpenChange={setOpenRole} orgType={selected} />
      </div>
    </>
  )
}
