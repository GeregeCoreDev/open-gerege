/**
 * 🔐 Scope Page (/[locale]/(main)/admin/scope/page.tsx)
 * 
 * Энэ нь OAuth/API scope удирдах хуудас юм.
 * Зорилго: Client applications-ийн эрхийн scope бүртгэл удирдлага
 * 
 * Features:
 * - ✅ Full CRUD operations
 * - ✅ Pagination with server-side meta
 * - ✅ Search by name with debounce
 * - ✅ URL state management (filter, page, page_size)
 * - ✅ Browser back/forward support
 * - ✅ Progress bar indicator
 * - ✅ Form validation (Zod)
 * - ✅ Responsive design
 * 
 * Table Columns:
 * - Actions (Edit/Delete)
 * - ID
 * - Owner Client ID
 * - Key (scope identifier)
 * - Description
 * - Created Date
 * 
 * URL Parameters:
 * - ?name=... - Search filter
 * - ?page=... - Current page
 * - ?page_size=... - Items per page
 * 
 * API Endpoints:
 * - GET /client/scope - List scopes
 * - POST /client/scope - Create scope
 * - PUT /client/scope - Update scope
 * - DELETE /client/scope - Delete scope
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
import { Plus, Pencil, Trash2, Loader2 } from 'lucide-react'
import { Progress } from '@/components/ui/progress'
import api from '@/lib/api'
import { useTranslations } from 'next-intl'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'

// i18n-aware navigation
import { usePathname, useRouter } from '@/i18n/navigation'
import { useSearchParams } from 'next/navigation'

type ScopeRow = App.Scope

/** ===== Zod Schemas ===== */
const ScopeSchema = z.object({
  owner_client_id: z.string().min(1, 'Required'),
  key: z.string().min(1, 'Required'),
  description: z.string().optional().default(''),
})
type ScopeIn = z.input<typeof ScopeSchema>
type ScopeOut = z.output<typeof ScopeSchema>

export default function ScopePage() {
  const t = useTranslations()
  const pathname = usePathname()
  const router = useRouter()
  const sp = useSearchParams()

  /** ===== URL -> initial state ===== */
  const initialName = sp.get('name') ?? ''
  const initialPage = Math.max(1, Number(sp.get('page') ?? '1') || 1)
  const initialPageSize = Math.max(1, Number(sp.get('page_size') ?? '50') || 50)

  /** ===== list state ===== */
  const [rows, setRows] = useState<ScopeRow[]>([])
  const [meta, setMeta] = useState<App.ApiMeta | null>(null)
  const [loading, setLoading] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)

  const [filterName, setFilterName] = useState(initialName)
  const [pageNumber, setPageNumber] = useState(initialPage)
  const [pageSize, setPageSize] = useState(initialPageSize)
  const [totalPage, setTotalPage] = useState(1)
  const [totalRow, setTotalRow] = useState(0)

  /** ===== progress bar ===== */
  const progress = useLoadingProgress(loading || deleting)

  /** ===== URL helpers ===== */
  const readUrl = React.useCallback(() => {
    const s = new URLSearchParams(Array.from(sp.entries()))
    const name = s.get('name') ?? ''
    const page = Math.max(1, Number(s.get('page') ?? '1') || 1)
    const page_size = Math.max(1, Number(s.get('page_size') ?? '50') || 50)
    return { name, page, page_size }
  }, [sp])

  function makeHref(basePath: string, q: { name?: string; page?: number; page_size?: number }) {
    const qs = new URLSearchParams()
    if (q.name) qs.set('name', q.name)
    if (q.page && q.page > 1) qs.set('page', String(q.page))
    if (q.page_size && q.page_size !== 50) qs.set('page_size', String(q.page_size))
    const s = qs.toString()
    return s ? `${basePath}?${s}` : basePath
  }

  const writeUrl = React.useCallback(
    (next: { name?: string; page?: number; page_size?: number }) => {
      const cur = readUrl()
      const href = makeHref(pathname, {
        name: next.name ?? cur.name,
        page: next.page ?? cur.page,
        page_size: next.page_size ?? cur.page_size,
      })
      router.replace(href, { scroll: false })
    },
    [pathname, readUrl, router],
  )

  /** ===== API load (ListData) ===== */
  async function load(page = pageNumber, size = pageSize, name = filterName) {
    setLoading(true)
    setFetchError(null)
    try {
      // ШИНЭ: ListData<T> буцаана, query нь page/size
      const data = await api.get<App.ListData<ScopeRow>>('/client/scope', {
        query: { page, size, name: name || undefined },
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

  // initial
  useEffect(() => {
    load(initialPage, initialPageSize, initialName)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  // react to URL changes (back/forward or manual typing)
  useEffect(() => {
    const { name, page, page_size } = readUrl()
    setFilterName(name)
    setPageNumber(page)
    setPageSize(page_size)
    load(page, page_size, name)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sp])

  /** ===== dialogs & forms ===== */
  const [openCreate, setOpenCreate] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [selected, setSelected] = useState<ScopeRow | null>(null)

  const createForm = useForm<ScopeIn>({
    resolver: zodResolver(ScopeSchema),
    defaultValues: { owner_client_id: '', key: '', description: '' },
  })
  const editForm = useForm<ScopeIn>({
    resolver: zodResolver(ScopeSchema),
    defaultValues: { owner_client_id: '', key: '', description: '' },
  })

  const onOpenCreate = () => {
    setOpenCreate(true)
    createForm.reset({ owner_client_id: '', key: '', description: '' })
  }
  const onOpenEdit = (r: ScopeRow) => {
    setSelected(r)
    editForm.reset({
      owner_client_id: r.owner_client_id ?? '',
      key: r.key ?? '',
      description: r.description ?? '',
    })
    setOpenEdit(true)
  }
  const onOpenDelete = (r: ScopeRow) => {
    setSelected(r)
    setOpenDelete(true)
  }

  /** ===== CRUD ===== */
  const onCreate: SubmitHandler<ScopeIn> = async (valuesIn) => {
    const values: ScopeOut = ScopeSchema.parse(valuesIn)
    try {
      await api.post<ScopeRow>('/client/scope', values)
      setOpenCreate(false)
      writeUrl({ page: 1 })
      load(1, pageSize, filterName)
    } catch {}
  }

  const onUpdate: SubmitHandler<ScopeIn> = async (valuesIn) => {
    if (!selected) return
    const payload = { id: selected.id, ...ScopeSchema.parse(valuesIn) }
    try {
      await api.put<ScopeRow>('/client/scope', payload)
      setOpenEdit(false)
      setSelected(null)
      load(pageNumber, pageSize, filterName)
    } catch {}
  }

  const onDelete = async () => {
    if (!selected) return
    try {
      setDeleting(true)
      await api.del<void>('/client/scope', { ...selected })
      const newCount = rows.length - 1
      const willBeEmpty = newCount <= 0 && pageNumber > 1
      const nextPage = willBeEmpty ? pageNumber - 1 : pageNumber
      setOpenDelete(false)
      setSelected(null)
      writeUrl({ page: nextPage })
      load(nextPage, pageSize, filterName)
    } catch {
    } finally {
      setDeleting(false)
    }
  }

  /** ===== filter & pagination ===== */
  useEffect(() => {
    const id = setTimeout(() => {
      writeUrl({ name: filterName || undefined, page: 1 })
      load(1, pageSize, filterName)
    }, 300)
    return () => clearTimeout(id)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filterName])

  const goPage = (n: number) => {
    const target = Math.min(Math.max(1, n), Math.max(1, totalPage))
    if (target !== pageNumber) {
      writeUrl({ page: target })
      load(target, pageSize, filterName)
    }
  }

  const changePageSize = (val: number) => {
    const s = Math.max(1, val || 50)
    writeUrl({ page_size: s, page: 1 })
    load(1, s, filterName)
  }

  /** ===== computed ===== */
  const pageLinks = useMemo(() => {
    const links: (number | 'ellipsis')[] = []
    if (totalPage <= 7) {
      for (let i = 1; i <= totalPage; i++) links.push(i)
      return links
    }
    const w = 2
    links.push(1)
    if (pageNumber > 1 + w + 1) links.push('ellipsis')
    const start = Math.max(2, pageNumber - w)
    const end = Math.min(totalPage - 1, pageNumber + w)
    for (let i = start; i <= end; i++) links.push(i)
    if (pageNumber < totalPage - w - 1) links.push('ellipsis')
    links.push(totalPage)
    return links
  }, [pageNumber, totalPage])

  const canPrev = meta?.has_prev ?? pageNumber > 1
  const canNext = meta?.has_next ?? pageNumber < totalPage
  const showingFrom = totalRow === 0 ? 0 : (meta?.start_idx ?? 0) + 1
  const showingTo = totalRow === 0 ? 0 : (meta?.end_idx ?? -1) + 1

  const isUpdating = editForm.formState.isSubmitting
  const isRowBusy = (id: number) =>
    (isUpdating && selected?.id === id) || (deleting && selected?.id === id)

  /** ===== render ===== */
  return (
    <>
      <div className="h-full w-full">
        <div className="relative flex h-full flex-col overflow-hidden">
          {progress > 0 && (
            <div className="absolute inset-x-0 top-0 z-10">
              <Progress value={progress} className="h-1 rounded-none" aria-label="Loading" />
            </div>
          )}

          <div className="flex flex-col overflow-hidden px-4 pb-4">
            {/* Header */}
            <div className="flex flex-col gap-3 pt-4 md:flex-row md:items-center md:justify-between">
              <h1 className="text-xl font-semibold text-gray-900 dark:text-white">{t('scope')}</h1>
              <Button onClick={onOpenCreate} className="gap-2" disabled={isUpdating || deleting}>
                <Plus className="h-4 w-4" />
                <span className="lowercase first-letter:uppercase">
                  {t('create', { name: t('scope') })}
                </span>
              </Button>
            </div>

        <Separator />

            {/* Filters */}
            <div className="flex flex-col gap-3 py-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex flex-1 flex-wrap items-center gap-2">
                <Input
                  value={filterName}
                  onChange={(e) => setFilterName(e.target.value)}
                  placeholder={t('search_by_name')}
                  className="h-9 sm:w-64"
                />
              </div>

              <div className="flex items-center gap-x-2">
                <p className="text-muted-foreground text-sm">{t('rows')}</p>
                <Input
                  type="number"
                  min={5}
                  max={200}
                  step={5}
                  value={pageSize}
                  onChange={(e) => changePageSize(Number(e.target.value))}
                  className="h-9 w-[90px]"
                />
              </div>
            </div>

            {/* Table */}
            <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-t-md border border-gray-200 dark:border-gray-800">
          {fetchError ? (
            <div className="p-6 text-sm text-red-600">{fetchError}</div>
          ) : (
            <div className="flex min-h-0 flex-1 flex-col">
              <div className="min-w-full overflow-x-auto">
                <Table className="w-full table-fixed">
                  <colgroup>
                    <col style={{ width: 120 }} />
                    <col style={{ width: 80 }} />
                    <col style={{ width: 240 }} />
                    <col style={{ width: 400 }} />
                    <col style={{ width: 160 }} />
                  </colgroup>
                  <TableHeader>
                    <TableRow className="[&>th]:bg-gray-50 dark:[&>th]:bg-gray-800 [&>th]:sticky [&>th]:top-0 [&>th]:z-20">
                      <TableHead className="text-right"></TableHead>
                      <TableHead>ID</TableHead>
                      <TableHead>{t('owner_client_id')}</TableHead>
                      <TableHead>{t('translation_key')}</TableHead>
                      <TableHead>{t('description')}</TableHead>
                      <TableHead>{t('created_date')}</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {loading ? (
                      <TableRow>
                        <TableCell colSpan={6} className="py-6 text-center">
                          <Loader2 className="mx-auto h-4 w-4 animate-spin" />
                        </TableCell>
                      </TableRow>
                    ) : rows.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={6} className="text-muted-foreground py-6 text-center">
                          {t('no_information_available')}
                        </TableCell>
                      </TableRow>
                    ) : (
                      rows.map((r) => {
                        const busy = isRowBusy(r.id)
                        return (
                          <TableRow key={r.id} className="[&>td]:align-top">
                            <TableCell className="text-right">
                              <Button
                                variant="outline"
                                size="sm"
                                className="mr-2"
                                onClick={() => onOpenEdit(r)}
                                disabled={busy}
                                aria-label={`Edit ${r.id}`}
                              >
                                <Pencil className="h-4 w-4" />
                              </Button>
                              <Button
                                variant="destructive"
                                size="sm"
                                onClick={() => onOpenDelete(r)}
                                disabled={busy}
                                aria-label={`Delete ${r.id}`}
                              >
                                <Trash2 className="h-4 w-4" />
                              </Button>
                            </TableCell>
                            <TableCell className="text-muted-foreground">{r.id}</TableCell>
                            <TableCell className="text-muted-foreground">
                              {r.owner_client_id}
                            </TableCell>
                            <TableCell className="font-medium wrap-break-word">{r.key}</TableCell>
                            <TableCell className="text-muted-foreground wrap-break-word">
                              {r.description || <span className="opacity-60">—</span>}
                            </TableCell>
                            <TableCell className="text-muted-foreground">
                              {r.created_date?.slice(0, 10) || '-'}
                            </TableCell>
                          </TableRow>
                        )
                      })
                    )}
                  </TableBody>
                </Table>
              </div>
            </div>
          )}
        </div>

            {/* Pagination */}
            <div className="flex items-center justify-between gap-3 rounded-b-md border-r border-b border-l border-gray-200 px-6 py-2 dark:border-gray-800">
              <div className="text-muted-foreground w-auto min-w-72 text-sm">
                {t.rich('showing', {
                  from: () => <span className="font-medium">{showingFrom}</span>,
                  to: () => <span className="font-medium">{showingTo}</span>,
                  total: () => <span className="font-medium">{totalRow}</span>,
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

        {/* ===== Create ===== */}
      <Dialog open={openCreate} onOpenChange={setOpenCreate}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="lowercase first-letter:uppercase">
              {t('create', { name: t('scope') })}
            </DialogTitle>
            <DialogDescription>{t('fill_the_fields_and_save_to_create')}</DialogDescription>
          </DialogHeader>

          <Form {...createForm}>
            <form
              onSubmit={createForm.handleSubmit(onCreate)}
              className="space-y-4"
              autoComplete="off"
            >
              <FormField
                control={createForm.control}
                name="owner_client_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('owner_client_id')}</FormLabel>
                    <FormControl>
                      <Input placeholder={t('owner_client_id')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={createForm.control}
                name="key"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('translation_key')}</FormLabel>
                    <FormControl>
                      <Input placeholder={t('translation_key')} {...field} />
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
                      <Input placeholder={`${t('description')} (${t('optional')})`} {...field} />
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

      {/* ===== Edit ===== */}
      <Dialog open={openEdit} onOpenChange={setOpenEdit}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="lowercase first-letter:uppercase">
              {t('update', { name: t('scope') })}
            </DialogTitle>
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
                name="owner_client_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Owner client ID</FormLabel>
                    <FormControl>
                      <Input placeholder="Owner client ID" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={editForm.control}
                name="key"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('translation_key')}</FormLabel>
                    <FormControl>
                      <Input placeholder={t('translation_key')} {...field} />
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
                      <Input placeholder={`${t('description')} (${t('optional')})`} {...field} />
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

      {/* ===== Delete ===== */}
      <Dialog open={openDelete} onOpenChange={setOpenDelete}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="lowercase first-letter:uppercase">
              {t('delete', { name: t('scope') })}
            </DialogTitle>
            <DialogDescription className="pt-2 text-base">
              {t.rich('delete_warning', {
                name: () => <span className="font-medium">{selected?.key}</span>,
              })}
            </DialogDescription>
          </DialogHeader>

          <DialogFooter className="pt-2">
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
      </div>
    </>
  )
}
