/**
 * 🎯 Action Page (/[locale]/(main)/admin/action/page.tsx)
 *
 * Энэ нь үйлдэл удирдах админ хуудас юм.
 * Зорилго: Бүх үйлдлүүд CRUD удирдлага
 *
 * Features:
 * - ✅ Full CRUD operations
 * - ✅ Pagination with server-side meta
 * - ✅ Search by name/code filter
 * - ✅ Progress bar loading
 * - ✅ Active/Inactive toggle
 * - ✅ Form validation (Zod)
 *
 * Table Columns:
 * - Actions (Edit/Delete)
 * - ID
 * - Code
 * - Name
 * - Description
 * - Active status
 *
 * Form Fields:
 * - Code: Action identifier
 * - Name: Display name
 * - Description: Optional
 * - is_active: Enable/Disable
 *
 * @author Gerege Core Team
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
import { Plus, Pencil, Trash2, Loader2 } from 'lucide-react'
import { Progress } from '@/components/ui/progress'
import api from '@/lib/api'
import { useTranslations } from 'next-intl'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'
import { appConfig } from '@/config/app.config'
import { LucideIcon } from '@/lib/utils/icon'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { TooltipArrow } from '@radix-ui/react-tooltip'
import { Badge } from '@/components/ui/badge'

const ActionSchema = z.object({
  code: z.string().min(1, 'Required').max(255),
  name: z.string().min(1, 'Required').max(255),
  description: z.string().max(255).optional().or(z.literal('')),
  is_active: z.boolean().nullable().optional(),
})

type ActionForm = z.input<typeof ActionSchema>

export default function ActionsPage() {
  const t = useTranslations()

  const [rows, setRows] = useState<App.Action[]>([])
  const [loading, setLoading] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)

  // keep last meta for pagination UI
  const [meta, setMeta] = useState<App.ApiMeta | null>(null)

  // progress bar - simplified with hook
  const progress = useLoadingProgress(loading)

  // pagination
  const [pageNumber, setPageNumber] = useState(1)
  const [pageSize, setPageSize] = useState<number>(appConfig.pagination.defaultPageSize)
  const [totalPage, setTotalPage] = useState(1)
  const [totalRow, setTotalRow] = useState(0)

  // filters
  const [filterName, setFilterName] = useState('')

  async function load(page = pageNumber, size = pageSize) {
    setLoading(true)
    setFetchError(null)
    try {
      const data = await api.get<App.ListData<App.Action>>('/action', {
        query: {
          page,
          size,
          name: filterName || undefined,
        },
      })
      const m = data.meta
      setMeta(m)
      setRows(data.items ?? [])
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

  // dialogs
  const [openCreate, setOpenCreate] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [selected, setSelected] = useState<App.Action | null>(null)

  // forms
  const createForm = useForm<ActionForm>({
    resolver: zodResolver(ActionSchema),
    defaultValues: {
      code: '',
      name: '',
      description: '',
      is_active: true,
    },
  })
  const editForm = useForm<ActionForm>({
    resolver: zodResolver(ActionSchema),
    defaultValues: {
      code: '',
      name: '',
      description: '',
      is_active: true,
    },
  })

  const onOpenEdit = (row: App.Action) => {
    setSelected(row)
    editForm.reset({
      code: row.code ?? '',
      name: row.name ?? '',
      description: row.description ?? '',
      is_active: row.is_active ?? null,
    })
    setOpenEdit(true)
  }
  const onOpenDelete = (row: App.Action) => {
    setSelected(row)
    setOpenDelete(true)
  }

  // create / update / delete
  const onCreate: SubmitHandler<ActionForm> = async (valuesIn) => {
    const values = ActionSchema.parse(valuesIn)
    try {
      const payload = {
        code: values.code,
        name: values.name,
        description: values.description || undefined,
        is_active: values.is_active ?? undefined,
      }
      await api.post<App.Action>('/action', payload as Record<string, unknown>)
      setOpenCreate(false)
      createForm.reset({
        code: '',
        name: '',
        description: '',
        is_active: true,
      })
      await load(1, pageSize)
    } catch {}
  }

  const onUpdate: SubmitHandler<ActionForm> = async (valuesIn) => {
    if (!selected) return
    const values = ActionSchema.parse(valuesIn)
    try {
      const payload = {
        id: selected.id,
        code: values.code,
        name: values.name,
        description: values.description || undefined,
        is_active: values.is_active ?? undefined,
      }
      await api.put<App.Action>(`/action/${selected.id}`, payload as Record<string, unknown>)
      setOpenEdit(false)
      setSelected(null)
      await load(pageNumber, pageSize)
    } catch {}
  }

  const onDelete = async () => {
    if (!selected) return
    try {
      setDeleting(true)
      await api.del<void>(`/action/${selected.id}`)
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

  // pagination helpers
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

  // showing ranges — серверийн meta-г ашиглая
  const showingFrom = totalRow === 0 ? 0 : (meta?.start_idx ?? 0) + 1
  const showingTo = totalRow === 0 ? 0 : (meta?.end_idx ?? -1) + 1

  // helpers
  const isCreating = createForm.formState.isSubmitting
  const isUpdating = editForm.formState.isSubmitting
  const isRowBusy = (rowId: string | number) =>
    (isUpdating && selected?.id === rowId) || (deleting && selected?.id === rowId)

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
              <h1 className="text-xl font-semibold text-gray-900 dark:text-white">{t('action')}</h1>
              <Button
                onClick={() => setOpenCreate(true)}
                className="gap-2"
                disabled={isCreating || isUpdating || deleting}
              >
                {isCreating ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Plus className="h-4 w-4" />
                )}
                <span className="lowercase first-letter:uppercase">
                  {t('create', { name: t('action') })}
                </span>
              </Button>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-3 py-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex flex-1 flex-wrap items-center gap-2">
                <Input
                  value={filterName}
                  onChange={(e) => setFilterName(e.target.value)}
                  placeholder={t('search_by_name')}
                  className="h-9 sm:w-56"
                />
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

            {/* Table Content */}
            <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-t-md border border-gray-200 dark:border-gray-800">
              {fetchError ? (
                <div className="p-6 text-sm text-red-600">{fetchError}</div>
              ) : (
                <div className="flex min-h-0 flex-1 flex-col">
                  <div className="min-w-full overflow-x-auto">
                    <Table className="w-full table-fixed">
                      <colgroup>
                        <col style={{ width: '160px' }} />
                        <col style={{ width: '80px' }} />
                        <col style={{ width: '160px' }} />
                        <col />
                        <col style={{ width: '240px' }} />
                        <col style={{ width: '180px' }} />
                        <col style={{ width: '120px' }} />
                      </colgroup>
                      <TableHeader>
                        <TableRow className="[&>th]:z-20 [&>th]:bg-gray-50 dark:[&>th]:bg-gray-800">
                          <TableHead></TableHead>
                          <TableHead>ID</TableHead>
                          <TableHead>{t('code')}</TableHead>
                          <TableHead>{t('name')}</TableHead>
                          <TableHead>{t('description')}</TableHead>
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
                            <col style={{ width: '160px' }} />
                            <col style={{ width: '80px' }} />
                            <col style={{ width: '160px' }} />
                            <col />
                            <col style={{ width: '240px' }} />
                            <col style={{ width: '180px' }} />
                            <col style={{ width: '120px' }} />
                          </colgroup>
                          <TableBody>
                            {rows.map((row) => {
                              const busy = isRowBusy(row.id)
                              return (
                                <TableRow key={row.id} className="[&>td]:align-center">
                                  <TableCell className="text-left">
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
                                          {t('update', { name: t('action') })}
                                        </p>
                                        <TooltipArrow className="fill-popover" />
                                      </TooltipContent>
                                    </Tooltip>

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
                                          {t('delete', { name: t('action') })}
                                        </p>
                                        <TooltipArrow className="fill-popover" />
                                      </TooltipContent>
                                    </Tooltip>
                                  </TableCell>

                                  <TableCell className="text-muted-foreground">{row.id}</TableCell>
                                  <TableCell>
                                    <Badge variant="outline">{row.code}</Badge>
                                  </TableCell>
                                  <TableCell className="font-medium">{row.name}</TableCell>
                                  <TableCell className="text-muted-foreground">
                                    {row.description || <span className="opacity-60">—</span>}
                                  </TableCell>
                                  <TableCell>
                                    <ActiveBadge value={row.is_active} />
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

            {/* Footer / Pagination */}
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
                        <PaginationLink
                          className="cursor-pointer"
                          isActive={p === pageNumber}
                          onClick={() => goPage(p)}
                        >
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

        {/* ---------- Create Dialog ---------- */}
        <Dialog open={openCreate} onOpenChange={setOpenCreate}>
          <DialogContent className="sm:max-w-lg">
            <DialogHeader>
              <DialogTitle className="lowercase first-letter:uppercase">
                {t('create', { name: t('action') })}
              </DialogTitle>
              <DialogDescription>{t('fill_the_fields_and_save_to_create')}</DialogDescription>
            </DialogHeader>

            <Form {...createForm}>
              <form
                onSubmit={createForm.handleSubmit(onCreate)}
                className="space-y-4 pt-2"
                autoComplete="off"
              >
                <div className="grid gap-4 sm:grid-cols-2">
                  <FormField
                    control={createForm.control}
                    name="code"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('code')}</FormLabel>
                        <FormControl>
                          <Input placeholder={t('code')} {...field} />
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
                          <Input placeholder={t('name')} {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={createForm.control}
                    name="is_active"
                    render={({ field }) => (
                      <FormItem className="w-full">
                        <FormLabel>{t('status')}</FormLabel>
                        <FormControl className="w-full">
                          <Select
                            value={
                              field.value === null || field.value === undefined
                                ? 'true'
                                : field.value
                                  ? 'true'
                                  : 'false'
                            }
                            onValueChange={(v) => field.onChange(v === 'true')}
                          >
                            <SelectTrigger className="h-9 w-full">
                              <SelectValue placeholder={t('status')} />
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
                        <Textarea
                          placeholder={`${t('description')}, ${t('optional')}`}
                          rows={4}
                          {...field}
                        />
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

        {/* ---------- Edit Dialog ---------- */}
        <Dialog open={openEdit} onOpenChange={setOpenEdit}>
          <DialogContent className="sm:max-w-lg">
            <DialogHeader>
              <DialogTitle className="lowercase first-letter:uppercase">
                {t('update', { name: t('action') })}
              </DialogTitle>
              <DialogDescription>{t('update_fields_and_save_your_changes')}</DialogDescription>
            </DialogHeader>

            <Form {...editForm}>
              <form
                onSubmit={editForm.handleSubmit(onUpdate)}
                className="space-y-4 pt-2"
                autoComplete="off"
              >
                <div className="grid w-full gap-4 sm:grid-cols-2">
                  <FormField
                    control={editForm.control}
                    name="code"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t('code')}</FormLabel>
                        <FormControl>
                          <Input placeholder={t('code')} {...field} />
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
                          <Input placeholder={t('name')} {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={editForm.control}
                    name="is_active"
                    render={({ field }) => (
                      <FormItem className="w-full">
                        <FormLabel>{t('status')}</FormLabel>
                        <FormControl className="w-full">
                          <Select
                            value={
                              field.value === null || field.value === undefined
                                ? 'null'
                                : field.value
                                  ? 'true'
                                  : 'false'
                            }
                            onValueChange={(v) =>
                              field.onChange(v === 'null' ? null : v === 'true' ? true : false)
                            }
                          >
                            <SelectTrigger className="h-9 w-full">
                              <SelectValue placeholder={t('status')} />
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
                        <Textarea
                          placeholder={`${t('description')}, ${t('optional')}`}
                          rows={4}
                          {...field}
                        />
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

        {/* ---------- Delete Dialog ---------- */}
        <Dialog open={openDelete} onOpenChange={setOpenDelete}>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle className="lowercase first-letter:uppercase">
                {t('delete', { name: t('action') })}
              </DialogTitle>
              <DialogDescription className="pt-2 text-base">
                {t.rich('delete_warning', {
                  name: () => <span className="font-medium">{selected?.name}</span>,
                })}
              </DialogDescription>
            </DialogHeader>

            <DialogFooter className="pt-2">
              <Button variant="outline" onClick={() => setOpenDelete(false)} disabled={deleting}>
                {t('cancel')}
              </Button>
              <Button variant="destructive" onClick={onDelete} disabled={deleting}>
                {deleting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                <p className="capitalize">{t('delete', { name: '' })}</p>
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </>
  )
}
