/**
 * 👥 User Page (/[locale]/(main)/admin/user/page.tsx)
 * 
 * Энэ нь хэрэглэгч удирдах админ хуудас юм.
 * Зорилго: Системийн бүх хэрэглэгчдийн CRUD удирдлага
 * 
 * Features:
 * - ✅ Full CRUD operations
 * - ✅ Pagination with server-side meta
 * - ✅ Search by name filter
 * - ✅ Progress bar with smooth animation
 * - ✅ Core integration: Find user from central system
 * - ✅ Two-step creation: Find → Confirm → Save
 * - ✅ User roles dialog integration
 * - ✅ Form validation (Zod)
 * - ✅ Responsive table with fixed columns
 * - ✅ Duplicate request prevention
 * 
 * Table Columns:
 * - Actions (Roles/Edit/Delete buttons)
 * - ID
 * - Registration Number
 * - Name (last_name + first_name)
 * - Phone Number
 * - Email
 * - Created Date
 * 
 * Create Workflow:
 * 1. Enter registration number
 * 2. Search from Core system (/user/find-from-core)
 * 3. Display found user info
 * 4. Add phone/email (optional)
 * 5. Save to local system
 * 
 * Related Components:
 * - UserRolesDialog: Manage user's system roles
 * 
 * API Endpoints:
 * - GET /user - List users with pagination
 * - POST /user/find-from-core - Search user from Core
 * - POST /user - Create user in local system
 * - PUT /user - Update phone/email
 * - DELETE /user/:id - Delete user
 * 
 * @author Sengum Soronzonbold
 * @company Gerege Core Team
 */

// app/[locale]/user/page.tsx
'use client'

import React, { useEffect, useMemo, useRef, useState } from 'react'
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
import { Plus, Pencil, Trash2, Loader2, Search, ShieldCheck } from 'lucide-react'
import { Progress } from '@/components/ui/progress'
import api from '@/lib/api'
import { useTranslations } from 'next-intl'
import UserRolesDialog from './actions/role'
import { useLoadingProgress } from '@/hooks/useLoadingProgress'
import { appConfig } from '@/config/app.config'

type UserRow = App.User & { created_date?: string }

/** ===== Schemas ===== */
const CreateUserSchema = z.object({
  reg_no: z.string().min(1, 'Required'),
})
type CreateUserIn = z.input<typeof CreateUserSchema>

const PhoneEmailSchema = z.object({
  phone_no: z.string().optional().default(''),
  email: z.string().email().or(z.literal('')).optional().default(''),
})
type EditUserIn = z.input<typeof PhoneEmailSchema>
type EditUserOut = z.output<typeof PhoneEmailSchema>

export default function UsersPage() {
  const t = useTranslations()

  /** ========= state ========= */
  const [rows, setRows] = useState<UserRow[]>([])
  const [loading, setLoading] = useState(false)
  const [fetchError, setFetchError] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)

  const [filterName, setFilterName] = useState('')
  const [pageNumber, setPageNumber] = useState(1)
  const [pageSize, setPageSize] = useState<number>(appConfig.pagination.defaultPageSize)
  const [totalPage, setTotalPage] = useState(1)
  const [totalRow, setTotalRow] = useState(0)

  const [openRoles, setOpenRoles] = useState(false)

  // progress bar - simplified with hook
  const progress = useLoadingProgress(loading)

  const lastReqId = useRef(0)

  const load = async (page = pageNumber, size = pageSize, name = filterName) => {
    const reqId = ++lastReqId.current
    setLoading(true)
    setFetchError(null)

    try {
      const data = await api.get<App.ListData<UserRow>>('/user', {
        query: { page, size, name: name || undefined },
      })
      if (reqId !== lastReqId.current) return

      const m = data.meta
      setRows(data.items ?? [])
      setPageNumber(m?.page ?? page)
      setPageSize(m?.size ?? size)
      setTotalPage(m?.pages ?? 1)
      setTotalRow(m?.total ?? 0)
    } catch {
      if (reqId === lastReqId.current) {
        setFetchError("Error occurred")
      }
    } finally {
      if (reqId === lastReqId.current) setLoading(false)
    }
  }

  useEffect(() => {
    load(pageNumber, pageSize, filterName)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageNumber, pageSize, filterName])

  const [openCreate, setOpenCreate] = useState(false)
  const [openEdit, setOpenEdit] = useState(false)
  const [openDelete, setOpenDelete] = useState(false)
  const [selected, setSelected] = useState<UserRow | null>(null)

  const createForm = useForm<CreateUserIn>({
    resolver: zodResolver(CreateUserSchema),
    defaultValues: { reg_no: '' },
  })
  const createExtraForm = useForm<EditUserIn>({
    resolver: zodResolver(PhoneEmailSchema),
    defaultValues: { phone_no: '', email: '' },
  })

  const editForm = useForm<EditUserIn>({
    resolver: zodResolver(PhoneEmailSchema),
    defaultValues: { phone_no: '', email: '' },
  })

  /** Core-оос хайлт */
  const [finding, setFinding] = useState(false)
  const [foundUser, setFoundUser] = useState<Partial<UserRow> | null>(null)
  const [findError, setFindError] = useState<string | null>(null)

  const doFindFromCore = async () => {
    setFindError(null)
    setFoundUser(null)
    const { reg_no } = createForm.getValues()
    const parsed = CreateUserSchema.safeParse({ reg_no })
    if (!parsed.success) {
      createForm.trigger('reg_no')
      return
    }
    setFinding(true)
    try {
      const data = await api.post<UserRow>('/user/find-from-core', { search_text: reg_no })
      setFoundUser(data)
      createExtraForm.reset({
        phone_no: data.phone_no ?? '',
        email: data.email ?? '',
      })
    } catch {
      setFindError("Error occurred")
    } finally {
      setFinding(false)
    }
  }

  const canCreate = !!foundUser && !!foundUser?.reg_no

  /** open handlers */
  const onOpenEdit = (r: UserRow) => {
    setSelected(r)
    editForm.reset({
      phone_no: r.phone_no ?? '',
      email: r.email ?? '',
    })
    setOpenEdit(true)
  }
  const onOpenDelete = (r: UserRow) => {
    setSelected(r)
    setOpenDelete(true)
  }

  /** CRUD */
  const onCreate: SubmitHandler<CreateUserIn> = async () => {
    if (!canCreate) {
      setFindError('Core-с хэрэглэгчээ эхлээд олж баталгаажуулна уу.')
      return
    }
    const extra = PhoneEmailSchema.parse(createExtraForm.getValues())
    try {
      await api.post<UserRow>('/user', { ...foundUser, ...extra })
      setOpenCreate(false)
      createForm.reset({ reg_no: '' })
      createExtraForm.reset({ phone_no: '', email: '' })
      setFoundUser(null)
      // эхний хуудсанд буцаад ганц load
      setPageNumber(1)
      load(1, pageSize, filterName)
    } catch {
      setFindError("Error occurred")
    }
  }

  const onUpdate: SubmitHandler<EditUserIn> = async (valuesIn) => {
    if (!selected) return
    const values: EditUserOut = PhoneEmailSchema.parse(valuesIn)
    const payload = { id: selected.id, ...values }
    try {
      await api.put<UserRow>('/user', payload)
      setOpenEdit(false)
      setSelected(null)
      load(pageNumber, pageSize, filterName)
    } catch {}
  }

  const onDelete = async () => {
    if (!selected) return
    try {
      setDeleting(true)
      await api.del<void>(`/user/${selected.id}`)
      const newCount = rows.length - 1
      const willBeEmpty = newCount <= 0 && pageNumber > 1
      const nextPage = willBeEmpty ? pageNumber - 1 : pageNumber
      setOpenDelete(false)
      setSelected(null)
      setPageNumber(nextPage) // effect нь ачаална
    } catch {
      setFetchError("Error occurred")
    } finally {
      setDeleting(false)
    }
  }

  /** ========= pagination / filters ========= */
  const goPage = (n: number) => {
    const target = Math.min(Math.max(1, n), Math.max(1, totalPage))
    if (target !== pageNumber) setPageNumber(target)
  }

  const changePageSize = (val: number) => {
    const s = Math.max(5, Math.min(200, val || 50))
    setPageSize(s)
    setPageNumber(1)
  }

  /** ========= computed ========= */
  const pageLinks = useMemo(() => {
    const links: (number | 'ellipsis')[] = []
    if (totalPage <= 7) {
      for (let i = 1; i <= totalPage; i++) links.push(i)
      return links
    }
    const windowSize = 2
    links.push(1)
    if (pageNumber > 1 + windowSize + 1) links.push('ellipsis')
    const start = Math.max(2, pageNumber - windowSize)
    const end = Math.min(totalPage - 1, pageNumber + windowSize)
    for (let i = start; i <= end; i++) links.push(i)
    if (pageNumber < totalPage - windowSize - 1) links.push('ellipsis')
    links.push(totalPage)
    return links
  }, [pageNumber, totalPage])

  const showingFrom = (pageNumber - 1) * pageSize + (rows.length ? 1 : 0)
  const showingTo = (pageNumber - 1) * pageSize + rows.length
  const canPrev = pageNumber > 1
  const canNext = pageNumber < totalPage

  const isUpdating = editForm.formState.isSubmitting
  const isRowBusy = (id: number) =>
    (isUpdating && selected?.id === id) || (deleting && selected?.id === id)

  /** ========= render ========= */
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
              <h1 className="text-xl font-semibold text-gray-900 dark:text-white">{t('all_user')}</h1>
              <Button
                onClick={() => {
                  setOpenCreate(true)
                  setFoundUser(null)
                  setFindError(null)
                  createForm.reset({ reg_no: '' })
                  createExtraForm.reset({ phone_no: '', email: '' })
                }}
                className="gap-2"
                disabled={isUpdating || deleting}
              >
                <Plus className="h-4 w-4" />
                <span className="lowercase first-letter:uppercase">
                  {t('create', { name: t('user') })}
                </span>
              </Button>
            </div>

            {/* Filters */}
            <div className="flex flex-col gap-3 py-3 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex flex-1 flex-wrap items-center gap-2">
                <Input
                  value={filterName}
                  onChange={(e) => {
                    setFilterName(e.target.value)
                    setPageNumber(1)
                  }}
                  placeholder={t('search_by_name')}
                  className="h-9 sm:w-64"
                />
              </div>

              <div className="flex items-center gap-x-2">
                <p className="text-muted-foreground text-sm">{t('rows')}</p>
                <Select value={String(pageSize)} onValueChange={(v) => changePageSize(Number(v))}>
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
              <div className="min-w-full overflow-x-auto">
                <Table className="w-full table-fixed">
                  <colgroup>
                    <col style={{ width: '120px' }} /> {/* actions */}
                    <col style={{ width: '80px' }} /> {/* id */}
                    <col style={{ width: '160px' }} /> {/* reg_no */}
                    <col style={{ width: '220px' }} /> {/* name */}
                    <col style={{ width: '220px' }} /> {/* phone */}
                    <col style={{ width: '360px' }} /> {/* email */}
                    <col style={{ width: '160px' }} /> {/* created_date */}
                  </colgroup>
                  <TableHeader>
                    <TableRow className="[&>th]:bg-gray-50 dark:[&>th]:bg-gray-800 [&>th]:sticky [&>th]:top-0 [&>th]:z-20">
                      <TableHead className="text-right"></TableHead>
                      <TableHead>ID</TableHead>
                      <TableHead>{t('reg_no')}</TableHead>
                      <TableHead>{t('name')}</TableHead>
                      <TableHead>{t('phone_no')}</TableHead>
                      <TableHead>{t('email')}</TableHead>
                      <TableHead>{t('created_date')}</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {loading ? (
                      <TableRow>
                        <TableCell colSpan={7} className="py-6 text-center">
                          <Loader2 className="mx-auto h-4 w-4 animate-spin" />
                        </TableCell>
                      </TableRow>
                    ) : rows.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={7} className="text-muted-foreground py-6 text-center">
                          {t('no_information_available')}
                        </TableCell>
                      </TableRow>
                    ) : (
                      rows.map((r) => {
                        const busy = isRowBusy(r.id)
                        return (
                          <TableRow key={r.id} className="[&>td]:align-middle">
                            <TableCell className="flex justify-start gap-x-2">
                              <Button
                                variant="default"
                                size="sm"
                                onClick={() => {
                                  setOpenRoles(true)
                                  setSelected(r)
                                }}
                                disabled={busy}
                                aria-label={`Roles ${r.id}`}
                              >
                                <ShieldCheck className="h-4 w-4" />
                              </Button>
                              <Button
                                variant="outline"
                                size="sm"
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
                              {r.reg_no || <span className="opacity-60">—</span>}
                            </TableCell>
                            <TableCell className="font-medium capitalize">
                              {[r.last_name, r.first_name].filter(Boolean).join(' ')}
                            </TableCell>
                            <TableCell className="text-muted-foreground">
                              {r.phone_no || <span className="opacity-60">—</span>}
                            </TableCell>
                            <TableCell className="text-muted-foreground">
                              {r.email || <span className="opacity-60">—</span>}
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

        {/* ---------- Create Dialog ---------- */}
      <Dialog open={openCreate} onOpenChange={setOpenCreate}>
        <DialogContent className="sm:max-w-xl">
          <DialogHeader>
            <DialogTitle className="lowercase first-letter:uppercase">
              {t('create', { name: t('user') })}
            </DialogTitle>
            <DialogDescription>{t('fill_the_fields_and_save_to_create')}</DialogDescription>
          </DialogHeader>

          <Form {...createForm}>
            <form
              onSubmit={createForm.handleSubmit(onCreate)}
              className="space-y-4 pt-2"
              autoComplete="off"
            >
              <div className="grid grid-cols-1 gap-3 sm:grid-cols-[1fr_auto]">
                <FormField
                  control={createForm.control}
                  name="reg_no"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('reg_no')}</FormLabel>
                      <FormControl>
                        <Input placeholder={t('reg_no')} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <div className="flex items-end">
                  <Button
                    type="button"
                    variant="secondary"
                    onClick={doFindFromCore}
                    disabled={finding}
                    className="w-full sm:w-auto"
                  >
                    {finding ? (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    ) : (
                      <Search className="mr-2 h-4 w-4" />
                    )}
                    {t('search')}
                  </Button>
                </div>
              </div>

              {findError && (
                <div className="rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
                  {findError}
                </div>
              )}
              {foundUser && (
                <div className="rounded-lg border p-3">
                  <div className="grid gap-2 sm:grid-cols-2">
                    <div>
                      <p className="text-muted-foreground text-xs">{t('reg_no')}</p>
                      <p className="font-medium">{foundUser?.reg_no}</p>
                    </div>
                    <div>
                      <p className="text-muted-foreground text-xs">{t('name')}</p>
                      <p className="font-medium">
                        {[
                          foundUser?.family_name,
                          foundUser?.last_name,
                          foundUser?.first_name,
                        ]
                          .filter(Boolean)
                          .join(' ')}
                      </p>
                    </div>
                  </div>
                </div>
              )}

              {/* Phone/Email to save */}
              <Form {...createExtraForm}>
                <div className="grid grid-cols-1 gap-3">
                  <FormField
                    control={createExtraForm.control}
                    name="phone_no"
                    render={({ field }) => (
                      <FormItem className="w-full">
                        <FormLabel>{t('phone_no')}</FormLabel>
                        <FormControl>
                          <Input placeholder={t('phone_no')} disabled={!foundUser} {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={createExtraForm.control}
                    name="email"
                    render={({ field }) => (
                      <FormItem className="w-full">
                        <FormLabel>{t('email')}</FormLabel>
                        <FormControl>
                          <Input
                            placeholder={t('email')}
                            type="email"
                            disabled={!foundUser}
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              </Form>

              <DialogFooter>
                <Button type="button" variant="outline" onClick={() => setOpenCreate(false)}>
                  {t('cancel')}
                </Button>
                <Button type="submit" disabled={!canCreate}>
                  {t('save')}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* ---------- Edit Dialog ---------- */}
      <Dialog open={openEdit} onOpenChange={setOpenEdit}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="lowercase first-letter:uppercase">
              {t('update', { name: t('user') })}
            </DialogTitle>
            <DialogDescription>{t('update_fields_and_save_your_changes')}</DialogDescription>
          </DialogHeader>

          <Form {...editForm}>
            <form
              onSubmit={editForm.handleSubmit(onUpdate)}
              className="space-y-4 pt-2"
              autoComplete="off"
            >
              <div className="grid grid-cols-1 gap-2 rounded-md border p-3 sm:grid-cols-2">
                <div>
                  <p className="text-muted-foreground text-xs">{t('reg_no')}</p>
                  <p className="font-medium">{selected?.reg_no || '—'}</p>
                </div>
                <div>
                  <p className="text-muted-foreground text-xs">{t('name')}</p>
                  <p className="font-medium">
                    {[selected?.family_name, selected?.last_name, selected?.first_name]
                      .filter(Boolean)
                      .join(' ') || '—'}
                  </p>
                </div>
              </div>

              <FormField
                control={editForm.control}
                name="phone_no"
                render={({ field }) => (
                  <FormItem className="w-full">
                    <FormLabel>{t('phone_no')}</FormLabel>
                    <FormControl>
                      <Input placeholder={t('phone_no')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={editForm.control}
                name="email"
                render={({ field }) => (
                  <FormItem className="w-full">
                    <FormLabel>{t('email')}</FormLabel>
                    <FormControl>
                      <Input placeholder={t('email')} type="email" {...field} />
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
              {t('delete', { name: t('user') })}
            </DialogTitle>
            <DialogDescription className="pt-2 text-base">
              {t.rich('delete_warning', {
                name: () => (
                  <span className="font-medium">
                    {[selected?.family_name, selected?.last_name, selected?.first_name]
                      .filter(Boolean)
                      .join(' ') || selected?.reg_no}
                  </span>
                ),
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

      <UserRolesDialog
        open={openRoles}
        onOpenChange={setOpenRoles}
        user={selected}
        onChanged={() => {
          /* need -> refresh parent if required */
        }}
      />
      </div>
    </>
  )
}
